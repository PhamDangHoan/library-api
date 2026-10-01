package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"library-api/internal/models"
	"library-api/internal/repositories"
)

var (
	ErrAuthorNotFound    = errors.New("author not found")
	ErrInvalidAuthorData = errors.New("invalid author data")
)

type AuthorService interface {
	GetAllAuthors() ([]models.Author, error)
	GetAuthorByID(id uint) (*models.Author, error)
	CreateAuthor(author *models.Author) error
	UpdateAuthor(author *models.Author) error
	DeleteAuthor(id uint) error
}

type authorService struct {
	authorRepository repositories.AuthorRepository
}

func NewAuthorService(
	authorRepository repositories.AuthorRepository,
) AuthorService {
	return &authorService{
		authorRepository: authorRepository,
	}
}

func (s *authorService) GetAllAuthors() ([]models.Author, error) {
	return s.authorRepository.GetAll()
}

func (s *authorService) GetAuthorByID(id uint) (*models.Author, error) {
	if id == 0 {
		return nil, ErrInvalidAuthorData
	}

	author, err := s.authorRepository.GetByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAuthorNotFound
	}

	if err != nil {
		return nil, err
	}

	return author, nil
}

func (s *authorService) CreateAuthor(
	author *models.Author,
) error {

	author.Name = strings.TrimSpace(author.Name)
	author.Bio = strings.TrimSpace(author.Bio)

	if author.Name == "" {
		return ErrInvalidAuthorData
	}

	return s.authorRepository.Create(author)
}

func (s *authorService) UpdateAuthor(
	author *models.Author,
) error {

	author.Name = strings.TrimSpace(author.Name)
	author.Bio = strings.TrimSpace(author.Bio)

	if author.ID == 0 || author.Name == "" {
		return ErrInvalidAuthorData
	}

	existingAuthor, err := s.authorRepository.GetByID(author.ID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAuthorNotFound
	}

	if err != nil {
		return err
	}

	existingAuthor.Name = author.Name
	existingAuthor.Bio = author.Bio

	return s.authorRepository.Update(existingAuthor)
}

func (s *authorService) DeleteAuthor(id uint) error {

	if id == 0 {
		return ErrInvalidAuthorData
	}

	err := s.authorRepository.Delete(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAuthorNotFound
	}

	return err
}
