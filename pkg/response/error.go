package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Message: message,
	})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Success: false,
		Message: message,
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: message,
	})
}
