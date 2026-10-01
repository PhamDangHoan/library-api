package routes

import (
	"github.com/gin-gonic/gin"

	"library-api/internal/handlers"
	"library-api/internal/middlewares"
)

func RegisterBookRoutes(
	router *gin.RouterGroup,
	bookHandler *handlers.BookHandler,
	jwtSecret string,
) {
	books := router.Group("/books")

	// Public routes
	books.GET("", bookHandler.GetBooks)
	books.GET("/:id", bookHandler.GetBook)

	// Authentication required
	protected := books.Group("")
	protected.Use(middlewares.AuthMiddleware(jwtSecret))

	// user + admin
	protected.POST("", bookHandler.CreateBook)
	protected.PUT("/:id", bookHandler.UpdateBook)

	// admin only
	admin := protected.Group("")
	admin.Use(middlewares.RequireRole("admin"))

	admin.DELETE("/:id", bookHandler.DeleteBook)
}
