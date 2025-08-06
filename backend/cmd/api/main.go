// File: /quicklynks/backend/cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/thisisjackii/quicklynks/backend/config"
	"github.com/thisisjackii/quicklynks/backend/internal/database"
	"github.com/thisisjackii/quicklynks/backend/internal/routes"
)

// @title quicklynks API
// @version 1.0
// @description This is the API for the quicklynks URL shortener application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if cfg.GinMode == "release" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Connect to database
	database.ConnectDB(cfg)
	logger.Info().Msg("Database connection successful")

	// Set Gin mode
	gin.SetMode(cfg.GinMode)
	router := gin.New()

	// Use structured logger middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log := logger.Info()
		if param.ErrorMessage != "" {
			log = logger.Error().Err(fmt.Errorf(param.ErrorMessage))
		}
		log.Str("client_ip", param.ClientIP).
			Str("method", param.Method).
			Int("status_code", param.StatusCode).
			Str("path", param.Path).
			Str("latency", param.Latency.String()).
			Msg("request processed")
		return ""
	}))
	router.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(router, database.DB, cfg)

	// Start server
	serverAddr := fmt.Sprintf("0.0.0.0:%s", cfg.ServerPort)
	logger.Info().Msgf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
