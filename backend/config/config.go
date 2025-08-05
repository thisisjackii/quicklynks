# File: backend/config/config.go
package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	Port              string `mapstructure:"PORT"`
	DatabaseURL       string `mapstructure:"DATABASE_URL"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	CORSAllowedOrigin string `mapstructure:"CORS_ALLOWED_ORIGIN"`
}

// Load reads configuration from file or environment variables.
func Load() (config Config, err error) {
	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "quicklynks.db")
	viper.SetDefault("JWT_SECRET", "skibididopdopdop")
	viper.SetDefault("CORS_ALLOWED_ORIGIN", "http://localhost:5173")

	// Look for a .env file in the current directory
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Automatically read environment variables
	viper.AutomaticEnv()

	// Read the config file if it exists
	// Ignore error if file not found, as env vars might be used instead
	_ = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	return
}