package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"library-api/internal/models"
	"library-api/internal/repositories"
	"library-api/pkg/cache"

	"context"
	"fmt"
)

type BookService interface {
	GetAllBooks(params repositories.BookListParams) (*repositories.BookListResult, error)
	GetBookByID(id uint) (*models.Book, error)

	CreateBook(book *models.Book, authorIDs []uint) error
	UpdateBook(book *models.Book, authorIDs []uint) error

	DeleteBook(id uint) error
}

type bookService struct {
	bookRepository   repositories.BookRepository
	authorRepository repositories.AuthorRepository
	cache            cache.Cache
}

func NewBookService(
	bookRepository repositories.BookRepository,
	authorRepository repositories.AuthorRepository,
	redisCache cache.Cache,
) BookService {
	return &bookService{
		bookRepository:   bookRepository,
		authorRepository: authorRepository,
		cache:            redisCache,
	}
}

func (s *bookService) GetAllBooks(
	params repositories.BookListParams,
) (*repositories.BookListResult, error) {

	ctx := context.Background()

	cacheKey := buildBookCacheKey(params)

	var cachedResult repositories.BookListResult

	err := s.cache.Get(
		ctx,
		cacheKey,
		&cachedResult,
	)

	if err == nil {
		return &cachedResult, nil
	}

	result, err := s.bookRepository.GetAll(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *bookService) GetBookByID(id uint) (*models.Book, error) {
	if id == 0 {
		return nil, errors.New("invalid book ID")
	}

	return s.bookRepository.GetByID(id)
}

func (s *bookService) getAuthors(authorIDs []uint) ([]models.Author, error) {
	if len(authorIDs) == 0 {
		return []models.Author{}, nil
	}

	authors := make([]models.Author, 0, len(authorIDs))

	for _, authorID := range authorIDs {
		author, err := s.authorRepository.GetByID(authorID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAuthorNotFound
			}

			return nil, err
		}

		authors = append(authors, *author)
	}

	return authors, nil
}

func (s *bookService) CreateBook(
	book *models.Book,
	authorIDs []uint,
) error {

	book.Title = strings.TrimSpace(book.Title)
	book.ISBN = strings.TrimSpace(book.ISBN)

	if book.Title == "" || book.ISBN == "" {
		return ErrInvalidBookData
	}

	if book.Quantity < 0 ||
		book.Available < 0 ||
		book.Available > book.Quantity {
		return ErrInvalidBookData
	}

	existingBook, err := s.bookRepository.FindByISBN(book.ISBN)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingBook != nil {
		return ErrDuplicateISBN
	}

	authors, err := s.getAuthors(authorIDs)
	if err != nil {
		return err
	}

	err = s.bookRepository.Create(book, authors)
	if err != nil {
		return err
	}

	ctx := context.Background()

	_ = s.cache.DeleteByPattern(
		ctx,
		"library:books:*",
	)

	return nil
}

func (s *bookService) UpdateBook(
	book *models.Book,
	authorIDs []uint,
) error {

	existingBook, err := s.bookRepository.GetByID(book.ID)

	if err != nil {
		return err
	}

	book.Title = strings.TrimSpace(book.Title)
	book.ISBN = strings.TrimSpace(book.ISBN)

	if book.Title == "" || book.ISBN == "" {
		return ErrInvalidBookData
	}

	if book.Quantity < 0 ||
		book.Available < 0 ||
		book.Available > book.Quantity {
		return ErrInvalidBookData
	}

	bookWithSameISBN, err :=
		s.bookRepository.FindByISBN(book.ISBN)

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if bookWithSameISBN != nil &&
		bookWithSameISBN.ID != book.ID {
		return ErrDuplicateISBN
	}

	existingBook.Title = book.Title
	existingBook.ISBN = book.ISBN
	existingBook.Description = book.Description
	existingBook.Quantity = book.Quantity
	existingBook.Available = book.Available

	authors, err := s.getAuthors(authorIDs)
	if err != nil {
		return err
	}

	err = s.bookRepository.Update(existingBook, authors)
	if err != nil {
		return err
	}

	ctx := context.Background()

	_ = s.cache.DeleteByPattern(
		ctx,
		"library:books:*",
	)

	return nil
}

func (s *bookService) DeleteBook(id uint) error {
	if id == 0 {
		return errors.New("invalid book ID")
	}

	err := s.bookRepository.Delete(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}

	if err != nil {
		return err
	}

	ctx := context.Background()

	_ = s.cache.DeleteByPattern(
		ctx,
		"library:books:*",
	)

	return nil
}

func buildBookCacheKey(params repositories.BookListParams) string {
	return fmt.Sprintf(
		"library:books:page=%d&limit=%d&search=%s&isbn=%s&sort=%s&order=%s",
		params.Page,
		params.Limit,
		params.Search,
		params.ISBN,
		params.Sort,
		params.Order,
	)
}
