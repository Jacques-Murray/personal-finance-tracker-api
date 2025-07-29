package repository

import (
	"context"
	"fmt"
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/models"

	"github.com/lib/pq"
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
	result := r.db.WithContext(ctx).Create(t)
	if result.Error != nil {
		if pqErr, ok := result.Error.(*pq.Error); ok {
			// Check for unique constraint violation (e.g., if description + date were unique)
			// Adjust this logic if you have specific unique constraints for transactions
			if pqErr.Code.Name() == "unique_violation" {
				return appErrors.NewConflictError("Transaction already exists with given details", result.Error)
			}
			// Check for foreign key violation (e.g., category_id does not exist)
			if pqErr.Code.Name() == "foreign_key_violation" {
				return appErrors.NewNotFoundError("Invalid category ID for transaction", result.Error)
			}
		}
		// Generic internal error for other DB issues
		return appErrors.NewInternalError("Failed to create transaction due to database error", result.Error)
	}
	return nil
}

// GetTransactions retrieves all transactions from the database
func (r *GormRepository) GetTransactions(ctx context.Context) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).Preload("Category").Order("date desc").Find(&transactions).Error
	if err != nil {
		// gorm.ErrRecordNotFound is typically for single record queries, Find returns empty slice
		return nil, appErrors.NewInternalError("Failed to retrieve transactions from database", err)
	}
	return transactions, nil
}

// CreateCategory adds a new category to the database
func (r *GormRepository) CreateCategory(ctx context.Context, c *models.Category) error {
	result := r.db.WithContext(ctx).Create(c)
	if result.Error != nil {
		if pqErr, ok := result.Error.(*pq.Error); ok {
			// Check for unique constraint violation for category name
			if pqErr.Code.Name() == "unique_violation" {
				return appErrors.NewAlreadyExistsError(fmt.Sprintf("Category with name '%s' already exists", c.Name), result.Error)
			}
		}
		// Generic internal error for other DB issues
		return appErrors.NewInternalError("Failed to create category due to database error", result.Error)
	}
	return nil
}

// GetCategories retrieves all categories, preloading their parent category
func (r *GormRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Preload("Parent").Find(&categories).Error
	if err != nil {
		// gorm.ErrRecordNotFound is typically for single record queries, Find returns empty slice
		return nil, appErrors.NewInternalError("Failed to retrieve categories from database", err)
	}
	return categories, nil
}
