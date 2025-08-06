// File: /quicklynks/backend/internal/routes/routes.go
package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thisisjackii/quicklynks/backend/config"
	"github.com/thisisjackii/quicklynks/backend/internal/controllers"
	"github.com/thisisjackii/quicklynks/backend/internal/middleware"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/thisisjackii/quicklynks/backend/internal/docs" // Import generated docs
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg config.Config) {
	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow SvelteKit dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Instantiate controllers
	authController := controllers.AuthController{DB: db, Cfg: cfg}
	linkController := controllers.LinkController{DB: db}

	// --- Public Routes ---

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Swagger Docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redirect route
	router.GET("/:shortCode", linkController.RedirectLink)

	// --- API v1 Routes ---
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Link routes (protected)
		links := api.Group("/links")
		links.Use(middleware.AuthMiddleware(cfg)) // Apply auth middleware
		{
			links.POST("", linkController.CreateLink)
			links.GET("", linkController.GetUserLinks)
		}
	}
}
