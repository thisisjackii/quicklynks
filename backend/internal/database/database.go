// File: /quicklynks/backend/internal/database/database.go
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thisisjackii/quicklynks/backend/config"
	"github.com/thisisjackii/quicklynks/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection pool.
var DB *gorm.DB

// ConnectDB initializes the database connection and runs migrations.
func ConnectDB(cfg config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslMode)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Don't log ErrRecordNotFound
			Colorful:                  true,        // Enable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Link{}, &models.Click{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
