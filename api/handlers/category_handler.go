package handlers

import (
	"net/http"
	"personal-finance-tracker-api/api/responses"
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/models"
	"personal-finance-tracker-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CategoryHandler holds the repository for database access
type CategoryHandler struct {
	Repo repository.Repository
}

// NewCategoryHandler creates a new handler for categories
func NewCategoryHandler(repo repository.Repository) *CategoryHandler {
	return &CategoryHandler{Repo: repo}
}

// CreateCategory handles the creation of a new category
// @Summary Create a new category
// @Description Add a new category to the system
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("CreateCategory: Invalid JSON format or data type mismatch")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid JSON format or data type mismatch.",
		})
		return
	}

	err := h.Repo.CreateCategory(c.Request.Context(), &category)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"category":  category,
			"errorType": appErrors.GetType(err),
		}).Error("CreateCategory: Failed to create category in repository.")

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
		"categoryID":   category.ID,
		"categoryName": category.Name,
	}).Info("CreateCategory: Category successfully created")

	c.JSON(http.StatusCreated, category)
}

// GetCategories handles listing all categories
// @Summary Get all categories
// @Description Retrieve a list of all transaction categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.Repo.GetCategories(c.Request.Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"errorType": appErrors.GetType(err),
		}).Error("GetCategories: Failed to retrieve categories from repository")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to retrieve categories.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"count": len(categories),
	}).Info("GetCategories: Categories retrieved successfully")

	c.JSON(http.StatusOK, categories)
}
