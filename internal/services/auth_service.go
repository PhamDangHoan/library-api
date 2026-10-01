package services

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"library-api/internal/models"
	"library-api/internal/repositories"
	"library-api/pkg/jwt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidUserData    = errors.New("invalid user data")
)

type AuthService interface {
	Register(name string, email string, password string) (*models.User, error)
	Login(email string, password string) (*models.User, string, error)
}

type authService struct {
	userRepository repositories.UserRepository
	jwtSecret      string
	jwtExpireHours int
}

func NewAuthService(
	userRepository repositories.UserRepository,
	jwtSecret string,
	jwtExpireHours int,
) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

func (s *authService) Register(
	name string,
	email string,
	password string,
) (*models.User, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" || email == "" || password == "" {
		return nil, ErrInvalidUserData
	}

	if len(password) < 6 {
		return nil, ErrInvalidUserData
	}

	existingUser, err := s.userRepository.FindByEmail(email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(
	email string,
	password string,
) (*models.User, string, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return nil, "", ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}

		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
		s.jwtSecret,
		s.jwtExpireHours,
	)

	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
