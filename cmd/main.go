package main

import (
	"fmt"
	"os"

	"personal-finance-tracker-api/api"
	"personal-finance-tracker-api/api/handlers" // Import handlers package
	"personal-finance-tracker-api/config"
	"personal-finance-tracker-api/internal/repository"
	"personal-finance-tracker-api/internal/services" // Import services package

	"github.com/sirupsen/logrus"
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
	// Initialize Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Load application configuration
	cfg := config.New()

	// Initialize database connection
	db := repository.InitDB(cfg.DatabaseURL)

	// Create repository instance
	repo := repository.NewGormRepository(db)

	// Create service instances, injecting the repository
	transactionService := services.NewTransactionService(repo)
	categoryService := services.NewCategoryService(repo)

	// Create handler instances, injecting the services
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Set up the router, passing the initialized handlers
	router := api.SetupRouter(transactionHandler, categoryHandler) // Changed signature

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.APIPort)
	logrus.WithFields(logrus.Fields{
		"address": serverAddr,
		"port":    cfg.APIPort,
	}).Info("Server starting")

	if err := router.Run(serverAddr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to start server")
	}
}
