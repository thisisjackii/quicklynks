// File: /quicklynks/backend/config/config.go
package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	GinMode    string `mapstructure:"GIN_MODE"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	SecretKey  string `mapstructure:"SECRET_KEY"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBSslMode  string `mapstructure:"DB_SSL_MODE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// Attempt to load .env file. This is for local development.
	// In production, environment variables are set directly.
	godotenv.Load()

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// If the config file is not found, we can continue as env vars might be set.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
