package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"library-api/internal/models"
	"library-api/internal/repositories"
	"library-api/internal/services"
	"library-api/internal/validators"
)

type BookHandler struct {
	service services.BookService
}

func NewBookHandler(service services.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

// GET /api/v1/books
// GetBooks godoc
// @Summary Get books
// @Description Get books with pagination, filtering and sorting
// @Tags Books
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of books per page" default(10)
// @Param search query string false "Search by title or description"
// @Param isbn query string false "Filter by ISBN"
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order: asc or desc"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/books [get]

func (h *BookHandler) GetBooks(c *gin.Context) {

	page := 1
	limit := 10

	if value := c.Query("page"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid page",
			})
			return
		}

		page = parsed
	}

	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 1 || parsed > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Limit must be between 1 and 100",
			})
			return
		}

		limit = parsed
	}

	params := repositories.BookListParams{
		Page:   page,
		Limit:  limit,
		Search: strings.TrimSpace(c.Query("search")),
		ISBN:   strings.TrimSpace(c.Query("isbn")),
		Sort:   strings.TrimSpace(c.Query("sort")),
		Order:  strings.ToLower(strings.TrimSpace(c.Query("order"))),
	}

	result, err := h.service.GetAllBooks(params)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result.Books,
		"pagination": gin.H{
			"page":        result.Page,
			"limit":       result.Limit,
			"total":       result.Total,
			"total_pages": result.TotalPages,
		},
	})
}

// GET /api/v1/books/:id
// GetBook godoc
// @Summary Get book by ID
// @Description Get a single book by its ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid book ID",
		})
		return
	}

	book, err := h.service.GetBookByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    book,
	})
}

// POST /api/v1/books
// CreateBook godoc
// @Summary Create a book
// @Description Create a new book and optionally associate authors
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param book body validators.CreateBookRequest true "Book data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var request validators.CreateBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	request.Title = strings.TrimSpace(request.Title)
	request.ISBN = strings.TrimSpace(request.ISBN)
	request.Description = strings.TrimSpace(request.Description)

	book := &models.Book{
		Title:       request.Title,
		ISBN:        request.ISBN,
		Description: request.Description,
		Quantity:    request.Quantity,
		Available:   request.Available,
	}

	if err := h.service.CreateBook(
		book,
		request.AuthorIDs,
	); err != nil {
		c.Error(err)
		return
	}

	// Lấy lại Book để response có Authors
	createdBook, err := h.service.GetBookByID(book.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Book created successfully",
		"data":    createdBook,
	})
}

// PUT /api/v1/books/:id
// UpdateBook godoc
// @Summary Update a book
// @Description Update book information and replace its authors
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Param book body validators.UpdateBookRequest true "Book data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/v1/books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid book ID",
		})
		return
	}

	var request validators.UpdateBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	request.Title = strings.TrimSpace(request.Title)
	request.ISBN = strings.TrimSpace(request.ISBN)
	request.Description = strings.TrimSpace(request.Description)

	book := &models.Book{
		ID:          uint(id),
		Title:       request.Title,
		ISBN:        request.ISBN,
		Description: request.Description,
		Quantity:    request.Quantity,
		Available:   request.Available,
	}

	if err := h.service.UpdateBook(
		book,
		request.AuthorIDs,
	); err != nil {
		c.Error(err)
		return
	}

	updatedBook, err := h.service.GetBookByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Book updated successfully",
		"data":    updatedBook,
	})
}

// DELETE /api/v1/books/:id
// DeleteBook godoc
// @Summary Delete a book
// @Description Soft delete a book. Admin role required.
// @Tags Books
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid book ID",
		})
		return
	}

	if err := h.service.DeleteBook(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Book deleted successfully",
	})
}
