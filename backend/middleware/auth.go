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
		accessToken, err := resolveAuthToken(c)
		if err != nil {
			util.RespondWithError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Validate the token
		userID, err := util.ValidateToken(accessToken, cfg.JWT.Secret)
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

func resolveAuthToken(c *gin.Context) (string, error) {
	tokenFromQuery := c.Query("token")
	if tokenFromQuery != "" {
		return tokenFromQuery, nil
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if the Authorization header has the correct format
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return headerParts[1], nil
}

// GetUserID retrieves the user ID from the Gin context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", errors.New("user ID not found in context")
	}
	return userID.(string), nil
}
