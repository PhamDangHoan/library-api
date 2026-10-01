package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"library-api/internal/services"
	"library-api/internal/validators"
)

type BorrowHandler struct {
	service services.BorrowService
}

func NewBorrowHandler(
	service services.BorrowService,
) *BorrowHandler {
	return &BorrowHandler{
		service: service,
	}
}

// POST /api/v1/borrows
// BorrowBook godoc
// @Summary Borrow a book
// @Description Borrow a book for the authenticated user
// @Tags Borrow / Return
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param borrow body validators.BorrowBookRequest true "Borrow data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/borrows [post]
func (h *BorrowHandler) BorrowBook(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User is not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user identity",
		})
		return
	}

	var request validators.BorrowBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	borrow, err := h.service.BorrowBook(
		userID,
		request.BookID,
		request.DueDate,
	)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Book borrowed successfully",
		"data":    borrow,
	})
}

// ReturnBook godoc
// @Summary Return a borrowed book
// @Description Return a book belonging to the authenticated user
// @Tags Borrow / Return
// @Produce json
// @Security BearerAuth
// @Param id path int true "Borrow ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/borrows/{id}/return [post]
func (h *BorrowHandler) ReturnBook(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User is not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user identity",
		})
		return
	}

	borrowID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || borrowID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid borrow ID",
		})
		return
	}

	borrow, err := h.service.ReturnBook(
		userID,
		uint(borrowID),
	)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Book returned successfully",
		"data":    borrow,
	})
}

// GetMyBorrows godoc
// @Summary Get my borrowed books
// @Description Get borrowing records of the authenticated user
// @Tags Borrow / Return
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/borrows/my [get]
func (h *BorrowHandler) GetMyBorrows(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User is not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user identity",
		})
		return
	}

	borrows, err := h.service.GetMyBorrows(userID)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    borrows,
	})
}
