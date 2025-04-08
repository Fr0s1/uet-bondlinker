package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken creates a new JWT token
func GenerateToken(userID, secret string, expiry time.Duration) (string, error) {
	tokenID := uuid.New().String()
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,                 // Subject (the user ID)
		"jti": tokenID,                // JWT ID (unique identifier for this token)
		"iat": now.Unix(),             // Issued At
		"nbf": now.Unix(),             // Not Before
		"exp": now.Add(expiry).Unix(), // Expiration Time
		"iss": "socialnet",            // Issuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates a JWT token and returns the user ID
func ValidateToken(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return "", jwt.ErrTokenExpired
			}
		}

		// Return the user ID (subject)
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", jwt.ErrTokenSignatureInvalid
}
