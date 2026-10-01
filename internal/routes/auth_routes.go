package routes

import (
	"library-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.RouterGroup,
	authHandler *handlers.AuthHandler,
) {
	auth := router.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}
