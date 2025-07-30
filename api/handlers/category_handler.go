package handlers

import (
	"fmt"
	"net/http"
	"personal-finance-tracker-api/api/middleware"
	"personal-finance-tracker-api/api/responses"
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var categoryValidate *validator.Validate

func init() {
	categoryValidate = validator.New()
}

// CategoryHandler holds the service for business logic access
type CategoryHandler struct {
	Service services.CategoryService
}

// NewCategoryHandler creates a new handler for categories
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service: service}
}

// CreateCategory handles the creation of a new category
// @Summary Create a new category
// @Description Add a new category to the system, associated with the authenticated user
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object"
// @Success 201 {object} models.Category
// @Failure 400 {object} responses.ValidationErrorResponse "Invalid input or validation error"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized (missing or invalid token)"
// @Failure 409 {object} responses.ErrorResponse "Conflict error (e.g., category name already exists for this user)"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		logrus.Error("CreateCategory: UserID not found in context, authentication middleware error.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Authenticated user ID not found.",
		})
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"userID": userID,
		}).Warn("CreateCategory: Invalid JSON format or data type mismatch.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid JSON format or data type mismatch.",
		})
		return
	}

	// Set the UserID from the authenticated context
	category.UserID = userID

	// Perform validation using the 'categoryValidate' instance
	if err := categoryValidate.Struct(category); err != nil {
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
				"category":         category,
				"userID":           userID,
			}).Warn("CreateCategory: Input validation error.")
			c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse{
				Error:  "Validation Error",
				Fields: fields,
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"category": category,
			"userID":   userID,
		}).Warn("CreateCategory: Unknown input validation error.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Validation failed: " + err.Error(),
		})
		return
	}

	// Capture both returned values
	createdCategory, err := h.Service.CreateCategory(c.Request.Context(), &category)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"category":  category,
			"errorType": appErrors.GetType(err),
			"userID":    userID,
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
		"userID":       userID,
	}).Info("CreateCategory: Category created successfully.")
	c.JSON(http.StatusCreated, createdCategory)
}

// GetCategories handles listing all categories with pagination
// @Summary Get all categories
// @Description Retrieve a list of all transaction categories with optional pagination, filtered by authenticated user
// @Tags categories
// @Produce json
// @Param limit query int false "Maximum number of categories to retrieve" default(100)
// @Param offset query int false "Number of categories to skip" default(0)
// @Param name query string false "Search categories by name (case-insensitive)"
// @Success 200 {array} models.Category
// @Failure 400 {object} responses.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized (missing or invalid token)"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		logrus.Error("GetCategories: UserID not found in context, authentication middleware error.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Authenticated user ID not found.",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		logrus.WithFields(logrus.Fields{
			"limitStr": limitStr,
			"error":    err,
			"userID":   userID,
		}).Warn("GetCategories: Invalid limit parameter, defaulting to 100.")
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		logrus.WithFields(logrus.Fields{
			"offsetStr": offsetStr,
			"error":     err,
			"userID":    userID,
		}).Warn("GetCategories: Invalid offset parameter, defaulting to 0.")
		offset = 0
	}

	// Filtering parameters
	var categoryName *string
	if nameStr := c.Query("name"); nameStr != "" {
		categoryName = &nameStr
	}

	categories, err := h.Service.GetCategories(c.Request.Context(), userID, limit, offset, categoryName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"errorType": appErrors.GetType(err),
			"userID":    userID,
		}).Error("GetCategories: Failed to retrieve categories via service.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve categories.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"count":  len(categories),
		"limit":  limit,
		"offset": offset,
		"userID": userID,
	}).Info("GetCategories: Categories retrieved successfully with pagination and user filter.")
	c.JSON(http.StatusOK, categories)
}
