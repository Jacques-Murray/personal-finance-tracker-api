package api

import (
	"personal-finance-tracker-api/api/handlers"
	"personal-finance-tracker-api/internal/repository"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "personal-finance-tracker-api/docs"
)

// SetupRouter configures the API routes and returns a Gin engine
func SetupRouter(repo repository.Repository) *gin.Engine {
	r := gin.Default()

	// Create handlers
	transactionHandler := handlers.NewTransactionHandler(repo)
	categoryHandler := handlers.NewCategoryHandler(repo)

	// Base path for the API
	api := r.Group("/api/v1")
	{
		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.GET("/export/csv", transactionHandler.ExportTransactionsCSV)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.GetCategories)
		}
	}

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
