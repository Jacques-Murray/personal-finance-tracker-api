package repository

import (
	"context"
	"fmt" // Import fmt for error messages
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/models"
	"time"

	"github.com/lib/pq" // Import for PostgreSQL specific error handling
	"gorm.io/gorm"
)

// Repository defines the interface for database operations
type Repository interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactions(ctx context.Context, userID uint, limit, offset int, startDate, endDate *time.Time, transactionType *models.TransactionType, description *string) ([]models.Transaction, error)
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategories(ctx context.Context, userID uint, limit, offset int) ([]models.Category, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	Transaction(txFunc func(txRepo Repository) error) error
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
			if pqErr.Code.Name() == "unique_violation" {
				return appErrors.NewConflictError("Transaction already exists with given details", result.Error)
			}
			if pqErr.Code.Name() == "foreign_key_violation" {
				return appErrors.NewValidationError("Invalid category ID or User ID for transaction", result.Error)
			}
		}
		return appErrors.NewInternalError("Failed to create transaction due to database error", result.Error)
	}
	return nil
}

// GetTransactions retrieves all transactions from the database with pagination
func (r *GormRepository) GetTransactions(ctx context.Context, userID uint, limit, offset int, startDate, endDate *time.Time, transactionType *models.TransactionType, description *string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Category").Order("date desc")

	// Apply date range filters
	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	// Apply transaction type filter
	if transactionType != nil && *transactionType != "" {
		query = query.Where("type = ?", *transactionType)
	}

	// Apply description filter
	if description != nil && *description != "" {
		query = query.Where("description ILIKE ?", "%"+*description+"%")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&transactions).Error
	if err != nil {
		return nil, appErrors.NewInternalError("Failed to retrieve transactions from database", err)
	}
	return transactions, nil
}

// CreateCategory adds a new category to the database
func (r *GormRepository) CreateCategory(ctx context.Context, c *models.Category) error {
	result := r.db.WithContext(ctx).Create(c)
	if result.Error != nil {
		if pqErr, ok := result.Error.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return appErrors.NewAlreadyExistsError(fmt.Sprintf("Category with name '%s' already exists", c.Name), result.Error)
			}
			if pqErr.Code.Name() == "foreign_key_violation" {
				return appErrors.NewValidationError("Invalid User ID for category", result.Error)
			}
		}
		return appErrors.NewInternalError("Failed to create category due to database error", result.Error)
	}
	return nil
}

// GetCategories retrieves all categories, preloading their parent category
func (r *GormRepository) GetCategories(ctx context.Context, userID uint, limit, offset int) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Parent")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, appErrors.NewInternalError("Failed to retrieve categories from database", err)
	}
	return categories, nil
}

// CreateUser adds a new user to the database
func (r *GormRepository) CreateUser(ctx context.Context, u *models.User) error {
	result := r.db.WithContext(ctx).Create(u)
	if result.Error != nil {
		if pqErr, ok := result.Error.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return appErrors.NewAlreadyExistsError(fmt.Sprintf("User with username '%s' already exists", u.Username), result.Error)
			}
		}
		return appErrors.NewInternalError("Failed to create user due to database error", result.Error)
	}
	return nil
}

// GetUserByUsername retrieves a user by their username
func (r *GormRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.NewNotFoundError(fmt.Sprintf("User '%s' not found", username), err)
		}
		return nil, appErrors.NewInternalError(fmt.Sprintf("Failed to retrieve user '%s' due to database error", username), err)
	}
	return &user, nil
}

// Transaction executes a function within a database transaction.
func (r *GormRepository) Transaction(txFunc func(txRepo Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txGormRepo := NewGormRepository(tx)
		return txFunc(txGormRepo)
	})
}
