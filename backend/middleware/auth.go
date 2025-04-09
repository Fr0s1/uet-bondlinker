
package middleware

import (
	"errors"
	"net/http"
	"strings"

	"socialnet/config"
	"socialnet/util"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies JWT tokens in request headers
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.RespondWithError(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			util.RespondWithError(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		// Validate the token
		userID, err := util.ValidateToken(headerParts[1], cfg.JWT.Secret)
		if err != nil {
			util.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}

// GetUserID retrieves the user ID from the Gin context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", errors.New("user ID not found in context")
	}
	return userID.(string), nil
}
