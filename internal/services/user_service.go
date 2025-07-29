package services

import (
	"context"
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/repository"

	"golang.org/x/crypto/bcrypt" // Import bcrypt for password hashing
)

// UserService defines the interface for user-related business logic
type UserService interface {
	RegisterUser(ctx context.Context, username, password string) (*models.User, error)
	// AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) // Will be implemented later
}

// userService implements the UserService interface
type userService struct {
	repo repository.Repository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.Repository) UserService {
	return &userService{repo: repo}
}

// RegisterUser handles new user registration, including password hashing
func (s *userService) RegisterUser(ctx context.Context, username, password string) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.NewInternalError("Failed to hash password", err)
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err // Error will be wrapped by repository already
	}

	return user, nil
}

// You can add more user-related methods here, such as AuthenticateUser, GetUserByID, etc.
// AuthenticateUser will compare a provided password with the stored hash.
/*
func (s *userService) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if appErrors.IsType(err, appErrors.TypeNotFound) {
			return nil, appErrors.NewUnauthorizedError("Invalid credentials", nil)
		}
		return nil, err // Propagate other errors
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, appErrors.NewUnauthorizedError("Invalid credentials", nil)
	}

	return user, nil
}
*/
