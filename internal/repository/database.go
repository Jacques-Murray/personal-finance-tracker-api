package repository

import (
	"personal-finance-tracker-api/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and performs auto-migration
func InitDB(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to connect to database")
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Category{}, &models.Transaction{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to migrate database schema")
	}

	logrus.Info("Database connection successful and schema migrated")
	return db
}
