package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"library-api/internal/models"
	"library-api/internal/services"
	"library-api/internal/validators"
)

type AuthorHandler struct {
	service services.AuthorService
}

func NewAuthorHandler(
	service services.AuthorService,
) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

// GetAuthors godoc
// @Summary Get all authors
// @Description Get all authors
// @Tags Authors
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/authors [get]
func (h *AuthorHandler) GetAuthors(c *gin.Context) {
	authors, err := h.service.GetAllAuthors()

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    authors,
	})
}

// GetAuthor godoc
// @Summary Get author by ID
// @Description Get an author by ID
// @Tags Authors
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/authors/{id} [get]
func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid author ID",
		})
		return
	}

	author, err := h.service.GetAuthorByID(uint(id))

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    author,
	})
}

// CreateAuthor godoc
// @Summary Create an author
// @Description Create a new author
// @Tags Authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param author body validators.CreateAuthorRequest true "Author data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/authors [post]
func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var request validators.CreateAuthorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	author := &models.Author{
		Name: strings.TrimSpace(request.Name),
		Bio:  strings.TrimSpace(request.Bio),
	}

	if err := h.service.CreateAuthor(author); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Author created successfully",
		"data":    author,
	})
}

// UpdateAuthor godoc
// @Summary Update an author
// @Description Update author information
// @Tags Authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Author ID"
// @Param author body validators.UpdateAuthorRequest true "Author data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/authors/{id} [put]
func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid author ID",
		})
		return
	}

	var request validators.UpdateAuthorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	author := &models.Author{
		ID:   uint(id),
		Name: strings.TrimSpace(request.Name),
		Bio:  strings.TrimSpace(request.Bio),
	}

	if err := h.service.UpdateAuthor(author); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Author updated successfully",
		"data":    author,
	})
}

// DeleteAuthor godoc
// @Summary Delete an author
// @Description Soft delete an author
// @Tags Authors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Author ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/authors/{id} [delete]
func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid author ID",
		})
		return
	}

	if err := h.service.DeleteAuthor(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Author deleted successfully",
	})
}
