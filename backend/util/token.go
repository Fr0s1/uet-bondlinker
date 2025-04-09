
package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID string, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// EmailVerificationClaims represents claims for email verification tokens
type EmailVerificationClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateEmailVerificationToken generates a JWT token for email verification
func GenerateEmailVerificationToken(userID string, secret string, expiry time.Duration) (string, error) {
	claims := &EmailVerificationClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "email_verification",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseEmailVerificationToken parses and validates an email verification token
func ParseEmailVerificationToken(tokenString string, secret string) (*EmailVerificationClaims, error) {
	claims := &EmailVerificationClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Subject != "email_verification" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// PasswordResetClaims represents claims for password reset tokens
type PasswordResetClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GeneratePasswordResetToken generates a JWT token for password reset
func GeneratePasswordResetToken(userID string, secret string, expiry time.Duration) (string, error) {
	claims := &PasswordResetClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "password_reset",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParsePasswordResetToken parses and validates a password reset token
func ParsePasswordResetToken(tokenString string, secret string) (*PasswordResetClaims, error) {
	claims := &PasswordResetClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Subject != "password_reset" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
