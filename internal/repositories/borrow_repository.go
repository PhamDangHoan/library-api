package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"library-api/internal/models"
)

type BorrowRepository interface {
	CreateBorrow(
		userID uint,
		bookID uint,
		dueDate *time.Time,
	) (*models.Borrow, error)

	GetActiveBorrow(
		userID uint,
		bookID uint,
	) (*models.Borrow, error)

	GetBorrowByID(id uint) (*models.Borrow, error)

	ReturnBorrow(
		borrowID uint,
	) (*models.Borrow, error)

	GetUserBorrows(userID uint) ([]models.Borrow, error)
}

type borrowRepository struct {
	db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) BorrowRepository {
	return &borrowRepository{
		db: db,
	}
}

func (r *borrowRepository) CreateBorrow(
	userID uint,
	bookID uint,
	dueDate *time.Time,
) (*models.Borrow, error) {

	var borrow models.Borrow

	err := r.db.Transaction(func(tx *gorm.DB) error {

		var book models.Book

		result := tx.First(&book, bookID)

		if result.Error != nil {
			return result.Error
		}

		if book.Available <= 0 {
			return errors.New("book is not available")
		}

		book.Available--

		if err := tx.Save(&book).Error; err != nil {
			return err
		}

		borrow = models.Borrow{
			UserID:     userID,
			BookID:     bookID,
			BorrowedAt: time.Now(),
			DueDate:    dueDate,
			Status:     "borrowed",
		}

		if err := tx.Create(&borrow).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &borrow, nil
}

func (r *borrowRepository) GetActiveBorrow(
	userID uint,
	bookID uint,
) (*models.Borrow, error) {

	var borrow models.Borrow

	result := r.db.
		Where(
			"user_id = ? AND book_id = ? AND status = ?",
			userID,
			bookID,
			"borrowed",
		).
		First(&borrow)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &borrow, nil
}

func (r *borrowRepository) GetBorrowByID(
	id uint,
) (*models.Borrow, error) {

	var borrow models.Borrow

	result := r.db.
		Preload("Book").
		Preload("User").
		Preload("Return").
		First(&borrow, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &borrow, nil
}

func (r *borrowRepository) ReturnBorrow(
	borrowID uint,
) (*models.Borrow, error) {

	var borrow models.Borrow

	err := r.db.Transaction(func(tx *gorm.DB) error {

		result := tx.First(&borrow, borrowID)

		if result.Error != nil {
			return result.Error
		}

		if borrow.Status != "borrowed" {
			return errors.New("borrow has already been returned")
		}

		var book models.Book

		if err := tx.First(&book, borrow.BookID).Error; err != nil {
			return err
		}

		book.Available++

		if book.Available > book.Quantity {
			book.Available = book.Quantity
		}

		if err := tx.Save(&book).Error; err != nil {
			return err
		}

		borrow.Status = "returned"

		if err := tx.Save(&borrow).Error; err != nil {
			return err
		}

		bookReturn := models.Return{
			BorrowID:   borrow.ID,
			ReturnedAt: time.Now(),
		}

		if err := tx.Create(&bookReturn).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &borrow, nil
}

func (r *borrowRepository) GetUserBorrows(
	userID uint,
) ([]models.Borrow, error) {

	var borrows []models.Borrow

	result := r.db.
		Preload("Book").
		Preload("Return").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&borrows)

	if result.Error != nil {
		return nil, result.Error
	}

	return borrows, nil
}
