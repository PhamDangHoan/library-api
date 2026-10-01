package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library-api/internal/services"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		status := http.StatusInternalServerError
		message := "Internal server error"

		switch {
		case errors.Is(err, services.ErrInvalidBookData):
			status = http.StatusBadRequest
			message = "Invalid book data"

		case errors.Is(err, services.ErrDuplicateISBN):
			status = http.StatusConflict
			message = "Book with this ISBN already exists"

		case errors.Is(err, gorm.ErrRecordNotFound):
			status = http.StatusNotFound
			message = "Book not found"

		case errors.Is(err, services.ErrEmailAlreadyExists):
			status = http.StatusConflict
			message = "Email already exists"

		case errors.Is(err, services.ErrInvalidUserData):
			status = http.StatusBadRequest
			message = "Invalid user data"

		case errors.Is(err, services.ErrInvalidCredentials):
			status = http.StatusUnauthorized
			message = "Invalid email or password"

		case errors.Is(err, services.ErrInvalidBorrowData):
			status = http.StatusBadRequest
			message = "Invalid borrow data"

		case errors.Is(err, services.ErrAlreadyBorrowed):
			status = http.StatusConflict
			message = "You have already borrowed this book"

		case errors.Is(err, services.ErrBookUnavailable):
			status = http.StatusConflict
			message = "Book is not available"

		case errors.Is(err, services.ErrBorrowNotFound):
			status = http.StatusNotFound
			message = "Borrow record not found"

		case errors.Is(err, services.ErrAlreadyReturned):
			status = http.StatusConflict
			message = "This book has already been returned"

		case errors.Is(err, services.ErrBorrowAccessDenied):
			status = http.StatusForbidden
			message = "You do not have permission to return this borrow"

		case errors.Is(err, services.ErrAuthorNotFound):
			status = http.StatusNotFound
			message = "Author not found"

		default:
			message = err.Error()
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": message,
		})
	}
}
