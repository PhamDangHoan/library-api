package tests

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"library-api/internal/models"
	"library-api/internal/repositories"
	"library-api/internal/services"
)

// ============================================================
// MOCK BOOK REPOSITORY
// ============================================================

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) GetAll(
	params repositories.BookListParams,
) (*repositories.BookListResult, error) {
	args := m.Called(params)

	var result *repositories.BookListResult

	if args.Get(0) != nil {
		result = args.Get(0).(*repositories.BookListResult)
	}

	return result, args.Error(1)
}

func (m *MockBookRepository) GetByID(
	id uint,
) (*models.Book, error) {
	args := m.Called(id)

	var book *models.Book

	if args.Get(0) != nil {
		book = args.Get(0).(*models.Book)
	}

	return book, args.Error(1)
}

func (m *MockBookRepository) FindByISBN(
	isbn string,
) (*models.Book, error) {
	args := m.Called(isbn)

	var book *models.Book

	if args.Get(0) != nil {
		book = args.Get(0).(*models.Book)
	}

	return book, args.Error(1)
}

func (m *MockBookRepository) Create(
	book *models.Book,
	authors []models.Author,
) error {
	args := m.Called(book, authors)

	return args.Error(0)
}

func (m *MockBookRepository) Update(
	book *models.Book,
	authors []models.Author,
) error {
	args := m.Called(book, authors)

	return args.Error(0)
}

func (m *MockBookRepository) Delete(
	id uint,
) error {
	args := m.Called(id)

	return args.Error(0)
}

// ============================================================
// MOCK AUTHOR REPOSITORY
// ============================================================

type MockAuthorRepository struct {
	mock.Mock
}

func (m *MockAuthorRepository) GetAll() ([]models.Author, error) {
	args := m.Called()

	var authors []models.Author

	if args.Get(0) != nil {
		authors = args.Get(0).([]models.Author)
	}

	return authors, args.Error(1)
}

func (m *MockAuthorRepository) GetByID(
	id uint,
) (*models.Author, error) {
	args := m.Called(id)

	var author *models.Author

	if args.Get(0) != nil {
		author = args.Get(0).(*models.Author)
	}

	return author, args.Error(1)
}

func (m *MockAuthorRepository) Create(
	author *models.Author,
) error {
	args := m.Called(author)

	return args.Error(0)
}

func (m *MockAuthorRepository) Update(
	author *models.Author,
) error {
	args := m.Called(author)

	return args.Error(0)
}

func (m *MockAuthorRepository) Delete(
	id uint,
) error {
	args := m.Called(id)

	return args.Error(0)
}

// ============================================================
// MOCK CACHE
// ============================================================

type MockCache struct {
	data map[string]interface{}
}

func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]interface{}),
	}
}

func (m *MockCache) Get(
	ctx context.Context,
	key string,
	dest interface{},
) error {
	// Simulate Redis cache MISS.
	return redis.Nil
}

func (m *MockCache) Set(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
) error {
	m.data[key] = value

	return nil
}

func (m *MockCache) Delete(
	ctx context.Context,
	key string,
) error {
	delete(m.data, key)

	return nil
}

func (m *MockCache) DeleteByPattern(
	ctx context.Context,
	pattern string,
) error {
	// For unit tests, clearing the cache map is enough
	// to simulate cache invalidation.
	m.data = make(map[string]interface{})

	return nil
}

// ============================================================
// TEST: CREATE BOOK - SUCCESS
// ============================================================

func TestCreateBook_Success(t *testing.T) {
	bookRepo := new(MockBookRepository)
	authorRepo := new(MockAuthorRepository)
	mockCache := NewMockCache()

	service := services.NewBookService(
		bookRepo,
		authorRepo,
		mockCache,
	)

	book := &models.Book{
		Title:       "Clean Code",
		ISBN:        "9780132350884",
		Description: "Software craftsmanship",
		Quantity:    10,
		Available:   10,
	}

	author := models.Author{
		ID:   1,
		Name: "Robert C. Martin",
	}

	bookRepo.
		On("FindByISBN", "9780132350884").
		Return(nil, gorm.ErrRecordNotFound).
		Once()

	authorRepo.
		On("GetByID", uint(1)).
		Return(&author, nil).
		Once()

	bookRepo.
		On("Create", book, []models.Author{author}).
		Return(nil).
		Once()

	err := service.CreateBook(
		book,
		[]uint{1},
	)

	require.NoError(t, err)

	assert.Equal(t, "Clean Code", book.Title)
	assert.Equal(t, "9780132350884", book.ISBN)
	assert.Equal(t, 10, book.Quantity)
	assert.Equal(t, 10, book.Available)

	bookRepo.AssertExpectations(t)
	authorRepo.AssertExpectations(t)
}

// ============================================================
// TEST: CREATE BOOK - DUPLICATE ISBN
// ============================================================

func TestCreateBook_DuplicateISBN(t *testing.T) {
	bookRepo := new(MockBookRepository)
	authorRepo := new(MockAuthorRepository)
	mockCache := NewMockCache()

	service := services.NewBookService(
		bookRepo,
		authorRepo,
		mockCache,
	)

	book := &models.Book{
		Title:     "Clean Architecture",
		ISBN:      "9780134494166",
		Quantity:  10,
		Available: 10,
	}

	existingBook := &models.Book{
		ID:    1,
		Title: "Existing Book",
		ISBN:  "9780134494166",
	}

	bookRepo.
		On("FindByISBN", "9780134494166").
		Return(existingBook, nil).
		Once()

	err := service.CreateBook(
		book,
		[]uint{},
	)

	require.Error(t, err)

	assert.ErrorIs(
		t,
		err,
		services.ErrDuplicateISBN,
	)

	bookRepo.AssertExpectations(t)
	authorRepo.AssertExpectations(t)

	bookRepo.AssertNotCalled(
		t,
		"Create",
		mock.Anything,
		mock.Anything,
	)
}

// ============================================================
// TEST: CREATE BOOK - INVALID DATA
// ============================================================

func TestCreateBook_InvalidData(t *testing.T) {
	bookRepo := new(MockBookRepository)
	authorRepo := new(MockAuthorRepository)
	mockCache := NewMockCache()

	service := services.NewBookService(
		bookRepo,
		authorRepo,
		mockCache,
	)

	book := &models.Book{
		Title:     "Invalid Book",
		ISBN:      "123456789",
		Quantity:  5,
		Available: 10,
	}

	err := service.CreateBook(
		book,
		[]uint{},
	)

	require.Error(t, err)

	assert.ErrorIs(
		t,
		err,
		services.ErrInvalidBookData,
	)

	bookRepo.AssertNotCalled(
		t,
		"FindByISBN",
		mock.Anything,
	)

	bookRepo.AssertNotCalled(
		t,
		"Create",
		mock.Anything,
		mock.Anything,
	)
}