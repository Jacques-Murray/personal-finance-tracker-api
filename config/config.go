package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	APIPort     string
	DatabaseURL string
	JWTSecret   string
}

// Global variable to hold the loaded configuration
var appConfig *Config

// New loads configuration from environment variables
func New() *Config {
	if appConfig != nil {
		return appConfig
	}

	if err := godotenv.Load(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Error loading .env file. Environment variables will be used directly.")
	}

	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "finance_tracker")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Create the database connection string
	databaseUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	appConfig = &Config{
		APIPort:     getEnv("API_PORT", "8080"),
		DatabaseURL: databaseUrl,
		JWTSecret:   getEnv("JWT_SECRET", "supersecretjwtkey"),
	}

	// Warn if using default JWT secret in production
	if appConfig.JWTSecret == "supersecretjwtkey" {
		logrus.Warn("Using default JWT_SECRET. Please set a strong, unique JWT_SECRET environment variable in production.")
	}

	return appConfig
}

// getEnv retrieves and environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	logrus.WithFields(logrus.Fields{
		"key": key,
	}).Info("Defaulting to fallback value for environment variable")
	return fallback
}

// GetJWTSecret provides access to the loaded JWT secret
func GetJWTSecret() string {
	if appConfig == nil {
		logrus.Fatal("Configuration not loaded. Call config.New() first.")
	}
	return appConfig.JWTSecret
}
