package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"personal-finance-tracker-api/api/responses"
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/services" // Import services package

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// TransactionHandler holds the service for business logic access
type TransactionHandler struct {
	Service services.TransactionService // Changed from Repo to Service
}

// NewTransactionHandler creates a new handler for transactions
func NewTransactionHandler(service services.TransactionService) *TransactionHandler { // Changed parameter
	return &TransactionHandler{Service: service} // Changed Repo to Service
}

// CreateTransaction handles the creation of a new transaction
// @Summary Create a new transaction
// @Description Add a new income or expense transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction object"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} responses.ValidationErrorResponse "Invalid input or validation error"
// @Failure 409 {object} responses.ErrorResponse "Conflict error (e.g., transaction already exists)"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("CreateTransaction: Invalid JSON format or data type mismatch.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid JSON format or data type mismatch.",
		})
		return
	}

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
			}).Warn("CreateTransaction: Input validation error.")
			c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse{
				Error:  "Validation Error",
				Fields: fields,
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"error":       err.Error(),
			"transaction": transaction,
		}).Warn("CreateTransaction: Unknown input validation error.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Validation failed: " + err.Error(),
		})
		return
	}

	// Call service layer instead of repository
	createdTransaction, err := h.Service.CreateTransaction(c.Request.Context(), &transaction)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err.Error(),
			"transaction": transaction,
			"errorType":   appErrors.GetType(err),
		}).Error("CreateTransaction: Failed to create transaction via service.")

		if appErrors.IsType(err, appErrors.TypeConflict) {
			c.JSON(http.StatusConflict, responses.ErrorResponse{
				Error:   "Conflict",
				Details: err.Error(),
			})
			return
		}
		if appErrors.IsType(err, appErrors.TypeValidation) {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Error:   "Bad Request",
				Details: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to create transaction.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"transactionID": createdTransaction.ID,
		"amount":        createdTransaction.Amount,
		"type":          createdTransaction.Type,
	}).Info("CreateTransaction: Transaction created successfully.")
	c.JSON(http.StatusCreated, createdTransaction)
}

// GetTransactions handles listing all transactions
// @Summary Get all transactions
// @Description Retrieve a list of all transactions, ordered by date
// @Tags transactions
// @Produce json
// @Success 200 {array} models.Transaction
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	// Call service layer instead of repository
	transactions, err := h.Service.GetTransactions(c.Request.Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"errorType": appErrors.GetType(err),
		}).Error("GetTransactions: Failed to retrieve transactions via service.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve transactions.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"count": len(transactions),
	}).Info("GetTransactions: Transactions retrieved successfully.")
	c.JSON(http.StatusOK, transactions)
}

// ExportTransactionsCSV handles exporting transactions to a CSV file
// @Summary Export transactions to CSV
// @Description Download a CSV file containing all transaction data
// @Tags transactions
// @Produce text/csv
// @Success 200 {file} file
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /transactions/export/csv [get]
func (h *TransactionHandler) ExportTransactionsCSV(c *gin.Context) {
	// Call service layer instead of repository
	transactions, err := h.Service.ExportTransactionsCSV(c.Request.Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"errorType": appErrors.GetType(err),
		}).Error("ExportTransactionsCSV: Failed to retrieve transactions for CSV export via service.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve transactions.",
		})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=transactions.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	header := []string{"ID", "Description", "Amount", "Type", "Date", "Category"}
	if err := writer.Write(header); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("ExportTransactionsCSV: Failed to write CSV header.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to write CSV header.",
		})
		return
	}

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
			}).Error("ExportTransactionsCSV: Failed to write CSV record. Stopping export.")
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Error:   "Internal Server Error",
				Details: "Failed to write CSV record during export.",
			})
			return
		}
	}
	logrus.Info("ExportTransactionsCSV: Transactions exported successfully.")
}
