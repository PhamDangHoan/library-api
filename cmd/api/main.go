// @title Library API
// @version 1.0
// @description RESTful Library Management API.
// @description Built with Go, Gin, GORM and MySQL.
// @host localhost:8080
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"

	"github.com/gin-gonic/gin"

	_ "library-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"library-api/configs"
	"library-api/internal/handlers"
	"library-api/internal/middlewares"
	"library-api/internal/models"
	"library-api/internal/repositories"
	"library-api/internal/routes"
	"library-api/internal/services"

	"library-api/pkg/cache"
)

func main() {
	// ========================================
	// 1. Load configuration
	// ========================================

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ========================================
	// 2. Connect database
	// ========================================

	db, err := configs.ConnectDatabase(config)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// ========================================
	// 2. Connect Redis
	// ========================================
	redisClient, err := configs.ConnectRedis(config)
	if err != nil {
		log.Fatalf("failed to connect Redis: %v", err)
	}

	defer redisClient.Close()

	log.Println("Redis connected successfully")

	// ========================================
	// 3. Auto migrate database
	// ========================================

	err = db.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Author{},
		&models.Borrow{},
		&models.Return{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// ========================================
	// 4. Create Gin router
	// ========================================

	router := gin.Default()

	router.Use(
		middlewares.RateLimit(redisClient),
	)

	// Don't trust any proxy by default.
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to configure trusted proxies: %v", err)
	}

	// Centralized error handling
	router.Use(middlewares.ErrorHandler())

	// ========================================
	// 5. Swagger
	// ========================================

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// ========================================
	// 6. Health check
	// ========================================

	// HealthCheck godoc
	// @Summary Health check
	// @Description Check whether the Library API is running
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Library API is running",
		})
	})

	// ========================================
	// 7. API v1
	// ========================================

	api := router.Group("/api/v1")

	// ========================================
	// 8. Authentication
	// ========================================

	userRepository := repositories.NewUserRepository(db)

	authService := services.NewAuthService(
		userRepository,
		config.JWTSecret,
		config.JWTExpireHours,
	)

	authHandler := handlers.NewAuthHandler(authService)

	routes.RegisterAuthRoutes(
		api,
		authHandler,
	)

	// ========================================
	// 9. Book
	// ========================================

	bookRepository := repositories.NewBookRepository(db)

	authorRepository := repositories.NewAuthorRepository(db)

	redisCache := cache.NewRedisCache(redisClient)

	bookService := services.NewBookService(
		bookRepository,
		authorRepository,
		redisCache,
	)

	bookHandler := handlers.NewBookHandler(bookService)

	routes.RegisterBookRoutes(
		api,
		bookHandler,
		config.JWTSecret,
	)

	// ========================================
	// 10. Author
	// ========================================

	authorService := services.NewAuthorService(
		authorRepository,
	)

	authorHandler := handlers.NewAuthorHandler(
		authorService,
	)

	routes.RegisterAuthorRoutes(
		api,
		authorHandler,
		config.JWTSecret,
	)

	// ========================================
	// 11. Borrow / Return
	// ========================================

	borrowRepository := repositories.NewBorrowRepository(db)

	borrowService := services.NewBorrowService(
		borrowRepository,
	)

	borrowHandler := handlers.NewBorrowHandler(
		borrowService,
	)

	routes.RegisterBorrowRoutes(
		api,
		borrowHandler,
		config.JWTSecret,
	)

	// ========================================
	// 12. Start server
	// ========================================

	log.Printf("Server running on port %s", config.AppPort)

	if err := router.Run(":" + config.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
