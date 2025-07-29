package middleware

import (
	"fmt"
	"net/http"
	"personal-finance-tracker-api/api/responses"
	"personal-finance-tracker-api/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware is a Gin middleware to authenticate requests using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logrus.Warn("AuthMiddleware: Missing Authorization header")
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Error:   "Unauthorized",
				Details: "Missing authentication token.",
			})
			c.Abort()
			return
		}

		// Check for "Bearer " prefix and remove it
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			logrus.Warn("AuthMiddleware: Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Error:   "Unauthorized",
				Details: "Invalid token format. Expected 'Bearer [token]'.",
			})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logrus.Warn("AuthMiddleware: Unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.GetJWTSecret()), nil
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("AuthMiddleware: Invalid or expired token")
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Error:   "Unauthorized",
				Details: "Invalid or expired authentication token.",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Token is valid, store claims in context for subsequent handlers
			if userID, ok := claims["userID"].(float64); ok {
				c.Set("userID", uint(userID))
			}
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}
			// You can store the entire claims map if needed
			c.Set("authClaims", claims)

			logrus.WithFields(logrus.Fields{
				"userID":   claims["userID"],
				"username": claims["username"],
				"path":     c.Request.URL.Path,
			}).Info("AuthMiddleware: Token validated successfully")

			c.Next()
		} else {
			logrus.Warn("AuthMiddleware: Invalid token claims or token not valid")
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Error:   "Unauthorized",
				Details: "Invalid authentication token claims.",
			})
			c.Abort()
		}
	}
}

// GetUserIDFromContext is a helper to retrieve userID from Gin context
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uint); ok {
			return id, true
		}
	}
	return 0, false
}

// GetUsernameFromContext is a helper to retrieve username from Gin context
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name, true
		}
	}
	return "", false
}
