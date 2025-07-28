package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"personal-finance-tracker-api/api/responses"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// Declare a global validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

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
		logrus.WithFields(logrus.Fields{
			"error":   err.Error(),
			"payload": c.Request.Body,
		}).Warn("CreateTransaction: Invalid JSON format or data type mismatch")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid JSON format or data type mismatch.",
		})
		return
	}

	// Perform validation using the 'validate' instance
	if err := validate.Struct(transaction); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var fields []responses.ValidationFieldError
			for _, fieldErr := range validationErrors {
				fields = append(fields, responses.ValidationFieldError{
					Field:   fieldErr.Field(),
					Tag:     fieldErr.Tag(),
					Message: fmt.Sprintf("Validation failed on '%s' for tag '%s'", fieldErr.Field(), fieldErr.Tag()),
				})
			}
			logrus.WithFields(logrus.Fields{
				"validationErrors": fields,
				"transaction":      transaction,
			}).Warn("CreateTransaction: Validation error")
			c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse{
				Error:  "Validation Error",
				Fields: fields,
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"error":       err.Error(),
			"transaction": transaction,
		}).Warn("CreateTransaction: Unknown validation error")

		// Fallback for other types of validation errors
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Validation failed: " + err.Error(),
		})
		return
	}

	if err := h.Repo.CreateTransaction(c.Request.Context(), &transaction); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err.Error(),
			"transaction": transaction,
		}).Error("CreateTransaction: Failed to create transaction in repository")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to create transaction.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"transactionID": transaction.ID,
		"amount":        transaction.Amount,
		"type":          transaction.Type,
	}).Info("CreateTransaction: Transaction created successfully")

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
	transactions, err := h.Repo.GetTransactions(c.Request.Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("GetTransactions: Failed to retrieve transactions from repository")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve transactions.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"transactionCount": len(transactions),
	}).Info("GetTransactions: Retrieved transactions successfully")

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
	transactions, err := h.Repo.GetTransactions(c.Request.Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("ExportTransactionsCSV: Failed to retrieve transactions for CSV export")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve transactions.",
		})
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
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("ExportTransactionsCSV: Failed to write CSV header")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to write CSV header.",
		})
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
			logrus.WithFields(logrus.Fields{
				"error":         err.Error(),
				"transactionID": t.ID,
			}).Error("ExportTransactionsCSV: Failed to write CSV record")
		}
	}
	logrus.Info("ExportTransactionsCSV: Transactions exported successfully")
}
