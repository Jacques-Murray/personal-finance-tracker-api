package api

import (
	"personal-finance-tracker-api/api/handlers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "personal-finance-tracker-api/docs"
)

// SetupRouter configures the API routes and returns a Gin engine
// It now accepts handler instances directly
func SetupRouter(
	transactionHandler *handlers.TransactionHandler,
	categoryHandler *handlers.CategoryHandler,
) *gin.Engine {
	r := gin.Default()

	// Custom Logrus Middleware
	r.Use(func(c *gin.Context) {
		startTime := time.Now()

		c.Next() // Process the request

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		logrus.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency":    latency,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
		}).Info("Request completed")
	})

	// Handlers are now passed in, no longer created here
	// transactionHandler := handlers.NewTransactionHandler(repo)
	// categoryHandler := handlers.NewCategoryHandler(repo)

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
