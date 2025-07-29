package services

import (
	"context"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/repository"
)

// TransactionService defines the interface for transaction-related business logic
type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	GetTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error)
	ExportTransactionsCSV(ctx context.Context) ([]models.Transaction, error)
}

// transactionService implements the TransactionService interface
type transactionService struct {
	repo repository.Repository
}

// NewTransactionService creates a new instance of TransactionService
func NewTransactionService(repo repository.Repository) TransactionService {
	return &transactionService{repo: repo}
}

// CreateTransaction handles the creation of a new transaction, applying business rules if any
func (s *transactionService) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	// Example: Here you could add more complex business logic before saving,
	// such as checking user balance, applying limits, etc.
	// For now, it directly calls the repository.
	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

// GetTransactions retrieves a list of transactions, applying business rules if any
func (s *transactionService) GetTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	transactions, err := s.repo.GetTransactions(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	// Example: Further processing or filtering of transactions based on business rules
	return transactions, nil
}

// ExportTransactionsCSV retrieves transactions for CSV export
func (s *transactionService) ExportTransactionsCSV(ctx context.Context) ([]models.Transaction, error) {
	transactions, err := s.repo.GetTransactions(ctx, 0, 0)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
