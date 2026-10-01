package routes

import (
	"github.com/gin-gonic/gin"

	"library-api/internal/handlers"
	"library-api/internal/middlewares"
)

func RegisterBorrowRoutes(
	router *gin.RouterGroup,
	borrowHandler *handlers.BorrowHandler,
	jwtSecret string,
) {
	borrows := router.Group("/borrows")

	borrows.Use(
		middlewares.AuthMiddleware(jwtSecret),
	)

	borrows.POST("", borrowHandler.BorrowBook)

	borrows.GET(
		"/my",
		borrowHandler.GetMyBorrows,
	)

	borrows.POST(
		"/:id/return",
		borrowHandler.ReturnBook,
	)
}
