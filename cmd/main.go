package main

import (
	"fmt"
	"log"

	"personal-finance-tracker-api/api"
	"personal-finance-tracker-api/config"
	"personal-finance-tracker-api/internal/repository"
)

// @title Personal Finance Tracker API
// @version 1.0
// @description This is a RESTful API for a personal finance tracking application.
// @termsOfService https://jacquesmurray.site/terms/

// @contact.name Jacques Murray
// @contact.url https://jacquesmurray.site/support
// @contact.email support@jacquesmurray.site

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load application configuration
	cfg := config.New()

	// Initialize database connection
	db := repository.InitDB(cfg.DatabaseURL)

	// Create a new repository instance
	repo := repository.NewGormRepository(db)

	// Set up the router
	router := api.SetupRouter(repo)

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.APIPort)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
