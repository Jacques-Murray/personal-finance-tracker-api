package handlers

import (
	"fmt"
	"net/http"
	"personal-finance-tracker-api/api/responses"
	"personal-finance-tracker-api/config"
	appErrors "personal-finance-tracker-api/internal/errors"
	"personal-finance-tracker-api/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var userValidate *validator.Validate

func init() {
	userValidate = validator.New()
}

// UserHandler holds the user service for business logic access
type UserHandler struct {
	UserService services.UserService // User service dependency
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// RegisterUserRequest represents the request body for user registration
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"` // Basic password validation
}

// LoginUserRequest represents the request body for user login
type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response body for a successful login
type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterUser handles new user registration
// @Summary Register a new user
// @Description Register a new user with a username and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterUserRequest true "User registration details"
// @Success 201 {object} models.User "User registered successfully"
// @Failure 400 {object} responses.ValidationErrorResponse "Invalid input or validation error"
// @Failure 409 {object} responses.ErrorResponse "Conflict (username already exists)"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("RegisterUser: Invalid JSON format or data type mismatch for registration.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid JSON format or data type mismatch.",
		})
		return
	}

	// Input validation for registration request
	if err := userValidate.Struct(req); err != nil {
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
				"username":         req.Username,
			}).Warn("RegisterUser: Input validation error for registration.")
			c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse{
				Error:  "Validation Error",
				Fields: fields,
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"username": req.Username,
		}).Warn("RegisterUser: Unknown input validation error for registration.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Validation failed: " + err.Error(),
		})
		return
	}

	// Call the user service to register the user
	user, err := h.UserService.RegisterUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"username":  req.Username,
			"errorType": appErrors.GetType(err),
		}).Error("RegisterUser: Failed to register user via service.")

		if appErrors.IsType(err, appErrors.TypeAlreadyExists) {
			c.JSON(http.StatusConflict, responses.ErrorResponse{
				Error:   "Conflict",
				Details: err.Error(),
			})
			return
		}
		// Catch other potential business logic validation errors from service layer
		if appErrors.IsType(err, appErrors.TypeValidation) {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{
				Error:   "Bad Request",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to register user.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID":   user.ID,
		"username": user.Username,
	}).Info("RegisterUser: User registered successfully.")
	c.JSON(http.StatusCreated, user)
}

// LoginUser handles user login and issues a JWT token
// @Summary Log in a user
// @Description Authenticate a user and return an authentication token
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginUserRequest true "User login details"
// @Success 200 {object} LoginResponse "Authentication successful with JWT token"
// @Failure 400 {object} responses.ValidationErrorResponse "Invalid input"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("LoginUser: Invalid JSON format or data type mismatch for login.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Invalid JSON format or data type mismatch.",
		})
		return
	}

	if err := userValidate.Struct(req); err != nil {
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
				"username":         req.Username,
			}).Warn("LoginUser: Input validation error for login.")
			c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse{
				Error:  "Validation Error",
				Fields: fields,
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"username": req.Username,
		}).Warn("LoginUser: Unknown input validation error for login.")
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Bad Request",
			Details: "Validation failed: " + err.Error(),
		})
		return
	}

	user, err := h.UserService.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"username":  req.Username,
			"errorType": appErrors.GetType(err),
		}).Error("LoginUser: Failed to authenticate user")

		if appErrors.IsType(err, appErrors.TypeUnauthorized) {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Error:   "Unauthorized",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to log in.",
		})
		return
	}

	// Generate JWT token
	// Define custom claims
	claims := jwt.MapClaims{
		"authorized": true,
		"userID":     user.ID,
		"username":   user.Username,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"userID":   user.ID,
			"username": user.Username,
		}).Error("LoginUser: Failed to generate JWT token.")
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Failed to generate authentication token.",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID":   user.ID,
		"username": user.Username,
	}).Info("LoginUser: User logged in successfully and JWT generated.")
	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}
