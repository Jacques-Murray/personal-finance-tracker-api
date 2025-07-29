package handlers

import (
	"net/http"
	"personal-finance-tracker-api/api/responses"
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/services" // Import services package

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CategoryHandler holds the service for business logic access
type CategoryHandler struct {
	Service services.CategoryService // Changed from Repo to Service
}

// NewCategoryHandler creates a new handler for categories
func NewCategoryHandler(service services.CategoryService) *CategoryHandler { // Changed parameter
	return &CategoryHandler{Service: service} // Changed Repo to Service
}

// CreateCategory handles the creation of a new category
// @Summary Create a new category
// @Description Add a new category to the system
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object"
// @Success 201 {object} models.Category
// @Failure 400 {object} responses.ErrorResponse "Invalid input"
// @Failure 409 {object} responses.ErrorResponse "Conflict error (e.g., category name already exists)"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("CreateCategory: Invalid JSON format or data type mismatch.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid JSON format or data type mismatch.",
		})
		return
	}

	// You could add validation here for Category as well, similar to Transaction
	// For example, if category name is required and has a min/max length:
	// if err := validate.Struct(category); err != nil { ... }

	// Call service layer instead of repository
	createdCategory, err := h.Service.CreateCategory(c.Request.Context(), &category)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"category":  category,
			"errorType": appErrors.GetType(err),
		}).Error("CreateCategory: Failed to create category via service.")

		if appErrors.IsType(err, appErrors.TypeAlreadyExists) {
			c.JSON(http.StatusConflict, responses.ErrorResponse{
				Error:   "Conflict",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to create category.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"categoryID":   createdCategory.ID,
		"categoryName": createdCategory.Name,
	}).Info("CreateCategory: Category created successfully.")
	c.JSON(http.StatusCreated, createdCategory)
}

// GetCategories handles listing all categories
// @Summary Get all categories
// @Description Retrieve a list of all transaction categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	// Call service layer instead of repository
	categories, err := h.Service.GetCategories(c.Request.Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"errorType": appErrors.GetType(err),
		}).Error("GetCategories: Failed to retrieve categories via service.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve categories.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"count": len(categories),
	}).Info("GetCategories: Categories retrieved successfully.")
	c.JSON(http.StatusOK, categories)
}
