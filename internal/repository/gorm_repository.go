package repository

import (
	"context"
	"personal-finance-tracker-api/internal/models"

	"gorm.io/gorm"
)

// Repository defines the interface for database operations
type Repository interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactions(ctx context.Context) ([]models.Transaction, error)
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategories(ctx context.Context) ([]models.Category, error)
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
func (r *GormRepository) CreateTransaction(ctx context.Context, t *models.Transaction) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// GetTransactions retrieves all transactions from the database
func (r *GormRepository) GetTransactions(ctx context.Context) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Order("date desc").Find(&transactions).Error
	return transactions, err
}

// CreateCategory adds a new category to the database
func (r *GormRepository) CreateCategory(ctx context.Context, c *models.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

// GetCategories retrieves all categories, preloading their parent category
func (r *GormRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Preload("Parent").Find(&categories).Error
	return categories, err
}
