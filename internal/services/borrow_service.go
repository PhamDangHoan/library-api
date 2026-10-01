package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"library-api/internal/models"
	"library-api/internal/repositories"
)

var (
	ErrBorrowNotFound     = errors.New("borrow not found")
	ErrAlreadyBorrowed    = errors.New("book is already borrowed by this user")
	ErrBookUnavailable    = errors.New("book is not available")
	ErrAlreadyReturned    = errors.New("borrow has already been returned")
	ErrInvalidBorrowData  = errors.New("invalid borrow data")
	ErrBorrowAccessDenied = errors.New("you do not have permission to return this borrow")
)

type BorrowService interface {
	BorrowBook(
		userID uint,
		bookID uint,
		dueDate *time.Time,
	) (*models.Borrow, error)

	ReturnBook(
		userID uint,
		borrowID uint,
	) (*models.Borrow, error)

	GetMyBorrows(userID uint) ([]models.Borrow, error)
}

type borrowService struct {
	borrowRepository repositories.BorrowRepository
}

func NewBorrowService(
	borrowRepository repositories.BorrowRepository,
) BorrowService {
	return &borrowService{
		borrowRepository: borrowRepository,
	}
}

func (s *borrowService) BorrowBook(
	userID uint,
	bookID uint,
	dueDate *time.Time,
) (*models.Borrow, error) {

	if userID == 0 || bookID == 0 {
		return nil, ErrInvalidBorrowData
	}

	existingBorrow, err := s.borrowRepository.GetActiveBorrow(
		userID,
		bookID,
	)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingBorrow != nil {
		return nil, ErrAlreadyBorrowed
	}

	borrow, err := s.borrowRepository.CreateBorrow(
		userID,
		bookID,
		dueDate,
	)

	if err != nil {
		if err.Error() == "book is not available" {
			return nil, ErrBookUnavailable
		}

		return nil, err
	}

	return borrow, nil
}

func (s *borrowService) ReturnBook(
	userID uint,
	borrowID uint,
) (*models.Borrow, error) {

	if userID == 0 || borrowID == 0 {
		return nil, ErrInvalidBorrowData
	}

	borrow, err := s.borrowRepository.GetBorrowByID(borrowID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBorrowNotFound
		}

		return nil, err
	}

	if borrow.UserID != userID {
		return nil, ErrBorrowAccessDenied
	}

	if borrow.Status != "borrowed" {
		return nil, ErrAlreadyReturned
	}

	return s.borrowRepository.ReturnBorrow(borrowID)
}

func (s *borrowService) GetMyBorrows(
	userID uint,
) ([]models.Borrow, error) {

	if userID == 0 {
		return nil, ErrInvalidBorrowData
	}

	return s.borrowRepository.GetUserBorrows(userID)
}
