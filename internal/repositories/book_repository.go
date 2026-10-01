package repositories

import (
	"errors"
	"library-api/internal/models"

	"gorm.io/gorm"
)

type BookListParams struct {
	Page   int
	Limit  int
	Search string
	ISBN   string
	Sort   string
	Order  string
}

type BookListResult struct {
	Books      []models.Book
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type BookRepository interface {
	GetAll(params BookListParams) (*BookListResult, error)
	GetByID(id uint) (*models.Book, error)
	FindByISBN(isbn string) (*models.Book, error)

	Create(book *models.Book, authors []models.Author) error
	Update(book *models.Book, authors []models.Author) error

	Delete(id uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) GetAll(params BookListParams) (*BookListResult, error) {
	var books []models.Book
	var total int64

	query := r.db.Model(&models.Book{})

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where(
			"title LIKE ? OR description LIKE ?",
			search,
			search,
		)
	}

	if params.ISBN != "" {
		query = query.Where("isbn = ?", params.ISBN)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	sortFields := map[string]string{
		"id":         "id",
		"title":      "title",
		"isbn":       "isbn",
		"quantity":   "quantity",
		"available":  "available",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	sortColumn, ok := sortFields[params.Sort]
	if !ok {
		sortColumn = "created_at"
	}

	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.
		Preload("Authors").
		Order(sortColumn + " " + order).
		Offset(offset).
		Limit(params.Limit).
		Find(&books).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &BookListResult{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *bookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book

	result := r.db.
		Preload("Authors").
		First(&book, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &book, nil
}

func (r *bookRepository) Create(
	book *models.Book,
	authors []models.Author,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		// 1. Tạo Book
		if err := tx.Create(book).Error; err != nil {
			return err
		}

		// 2. Gắn Authors vào Book
		if len(authors) > 0 {
			if err := tx.
				Model(book).
				Association("Authors").
				Replace(authors); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *bookRepository) Update(
	book *models.Book,
	authors []models.Author,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		// 1. Cập nhật Book
		if err := tx.Save(book).Error; err != nil {
			return err
		}

		// 2. Thay thế toàn bộ Authors
		if err := tx.
			Model(book).
			Association("Authors").
			Replace(authors); err != nil {
			return err
		}

		return nil
	})
}

func (r *bookRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Book{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *bookRepository) FindByISBN(isbn string) (*models.Book, error) {
	var book models.Book

	result := r.db.
		Where("isbn = ?", isbn).
		First(&book)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &book, nil
}
