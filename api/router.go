package api

import (
	"personal-finance-tracker-api/api/handlers"
	"personal-finance-tracker-api/api/middleware"
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
	userHandler *handlers.UserHandler,
) *gin.Engine {
	r := gin.Default()

	// Custom Logrus Middleware
	r.Use(func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

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

	// Base path for the API
	api := r.Group("/api/v1")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.RegisterUser)
			users.POST("/login", userHandler.LoginUser)
		}

		// Protected routes group: Apply AuthMiddleware to these routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())

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
