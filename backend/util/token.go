
package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT token
func GenerateToken(userID, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
