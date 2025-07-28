package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	APIPort     string
	DatabaseURL string
}

// New loads configuration from environment variables
func New() *Config {
	// godotenv.Load() will ignore the error if the .env file doesn't exist
	// This is useful for production environments where env vars are set directly
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
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

	return &Config{
		APIPort:     getEnv("API_PORT", "8080"),
		DatabaseURL: databaseUrl,
	}
}

// getEnv retrieves and environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Defaulting to %s for %s", fallback, key)
	return fallback
}
