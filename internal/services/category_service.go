package services

import (
	"context"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/repository"
)

// CategoryService defines the interface for category-related business logic
type CategoryService interface {
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
	GetCategories(ctx context.Context, userID uint, limit, offset int) ([]models.Category, error)
}

// categoryService implements the CategoryService interface
type categoryService struct {
	repo repository.Repository
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(repo repository.Repository) CategoryService {
	return &categoryService{repo: repo}
}

// CreateCategory handles the creation of a new category, applying business rules if any
func (s *categoryService) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	// Execute the creation within a database transaction
	err := s.repo.Transaction(func(txRepo repository.Repository) error {
		// Use txRepo for operations within this transaction
		if err := txRepo.CreateCategory(ctx, category); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return category, nil
}

// GetCategories retrieves a list of categories, applying business rules if any
func (s *categoryService) GetCategories(ctx context.Context, userID uint, limit, offset int) ([]models.Category, error) {
	categories, err := s.repo.GetCategories(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	// Example: Further processing or filtering of categories based on business rules
	return categories, nil
}
