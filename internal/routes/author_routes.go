package routes

import (
	"github.com/gin-gonic/gin"

	"library-api/internal/handlers"
	"library-api/internal/middlewares"
)

func RegisterAuthorRoutes(
	router *gin.RouterGroup,
	authorHandler *handlers.AuthorHandler,
	jwtSecret string,
) {
	authors := router.Group("/authors")

	// Public
	authors.GET("", authorHandler.GetAuthors)
	authors.GET("/:id", authorHandler.GetAuthor)

	// Authenticated
	protected := authors.Group("")
	protected.Use(middlewares.AuthMiddleware(jwtSecret))

	protected.POST("", authorHandler.CreateAuthor)
	protected.PUT("/:id", authorHandler.UpdateAuthor)

	// Admin only
	admin := protected.Group("")
	admin.Use(middlewares.RequireRole("admin"))

	admin.DELETE("/:id", authorHandler.DeleteAuthor)
}
