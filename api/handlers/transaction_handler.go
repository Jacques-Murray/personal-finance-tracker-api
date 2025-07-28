package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/repository"

	"github.com/gin-gonic/gin"
)

// TransactionHandler holds the repository for database access
type TransactionHandler struct {
	Repo repository.Repository
}

// NewTransactionHandler creates a new handler for transactions
func NewTransactionHandler(repo repository.Repository) *TransactionHandler {
	return &TransactionHandler{Repo: repo}
}

// CreateTransaction handles the creation of a new transaction
// @Summary Create a new transaction
// @Description Add a new income or expense transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction object"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if transaction.Type != models.Income && transaction.Type != models.Expense {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	if err := h.Repo.CreateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactions handles listing all transactions
// @Summary Get all transactions
// @Description Retrieve a list of all transactions, ordered by date
// @Tags transactions
// @Produce json
// @Success 200 {array} models.Transaction
// @Failure 500 {object} map[string]string
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	transactions, err := h.Repo.GetTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// ExportTransactionsCSV handles exporting transactions to a CSV file
// @Summary Export transactions to CSV
// @Description Download a CSV file containing all transaction data
// @Tags transactions
// @Produce text/csv
// @Success 200 {file} file
// @Failure 500 {object} map[string]string
// @Router /transactions/export/csv [get]
func (h *TransactionHandler) ExportTransactionsCSV(c *gin.Context) {
	transactions, err := h.Repo.GetTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	// Set headers for CSV download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=transactions.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ID", "Description", "Amount", "Type", "Date", "Category"}
	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
		return
	}

	// Write transaction data
	for _, t := range transactions {
		record := []string{
			fmt.Sprintf("%d", t.ID),
			t.Description,
			fmt.Sprintf("%.2f", t.Amount),
			string(t.Type),
			t.Date.Format("2006-01-02"),
			t.Category.Name,
		}
		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV record"})
			return
		}
	}
}
