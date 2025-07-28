package repository

import (
	"personal-finance-tracker-api/internal/models"

	"gorm.io/gorm"
)

// Repository defines the interface for database operations
type Repository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactions() ([]models.Transaction, error)
	CreateCategory(category *models.Category) error
	GetCategories() ([]models.Category, error)
}

// GormRepository is an implementation of Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// CreateTransaction adds a new transaction to the database
func (r *GormRepository) CreateTransaction(t *models.Transaction) error {
	return r.db.Create(t).Error
}

// GetTransactions retrieves all transactions from the database
func (r *GormRepository) GetTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Category").Order("date desc").Find(&transactions).Error
	return transactions, err
}

// CreateCategory adds a new category to the database
func (r *GormRepository) CreateCategory(c *models.Category) error {
	return r.db.Create(c).Error
}

// Getcategories retrieves all categories
func (r *GormRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}
