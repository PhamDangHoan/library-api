package repositories

import (
	"errors"

	"gorm.io/gorm"

	"library-api/internal/models"
)

type AuthorRepository interface {
	GetAll() ([]models.Author, error)
	GetByID(id uint) (*models.Author, error)
	Create(author *models.Author) error
	Update(author *models.Author) error
	Delete(id uint) error
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{
		db: db,
	}
}

func (r *authorRepository) GetAll() ([]models.Author, error) {
	var authors []models.Author

	result := r.db.
		Preload("Books").
		Find(&authors)

	if result.Error != nil {
		return nil, result.Error
	}

	return authors, nil
}

func (r *authorRepository) GetByID(id uint) (*models.Author, error) {
	var author models.Author

	result := r.db.
		Preload("Books").
		First(&author, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &author, nil
}

func (r *authorRepository) Create(author *models.Author) error {
	return r.db.Create(author).Error
}

func (r *authorRepository) Update(author *models.Author) error {
	return r.db.Save(author).Error
}

func (r *authorRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Author{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
