package config

import (
	"os"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Email    EmailConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host        string
	Port        string
	Env         string
	CorsOrigins []string
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	Secret     string
	ExpiryTime time.Duration
}

// AWSConfig holds AWS-specific configuration
type AWSConfig struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	CdnURL          string
}

// EmailConfig holds email-specific configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	FromEmail    string
	Password     string
	FrontendURL  string
	VerifyExpiry time.Duration
	ResetExpiry  time.Duration
}

// New creates a new configuration from environment variables
func New() *Config {
	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	verifyExpiry, err := time.ParseDuration(getEnv("EMAIL_VERIFY_EXPIRY", "48h"))
	if err != nil {
		verifyExpiry = 48 * time.Hour
	}

	resetExpiry, err := time.ParseDuration(getEnv("PASSWORD_RESET_EXPIRY", "15m"))
	if err != nil {
		resetExpiry = 15 * time.Minute
	}

	return &Config{
		Server: ServerConfig{
			Host:        getEnv("HOST", "0.0.0.0"),
			Port:        getEnv("PORT", "8080"),
			Env:         getEnv("ENV", "development"),
			CorsOrigins: strings.Split(getEnv("ALLOWED_CORS_ORIGINS", "http://localhost:5173"), ","),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "socialnet"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "default_jwt_secret"),
			ExpiryTime: jwtExpiry,
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			Bucket:          getEnv("AWS_BUCKET", "socialnet-uploads"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Endpoint:        getEnv("AWS_ENDPOINT", ""),
			CdnURL:          getEnv("AWS_CDN_URL", ""),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
			SMTPPort:     getEnv("EMAIL_SMTP_PORT", "587"),
			FromEmail:    getEnv("EMAIL_FROM", "noreply@socialnet.com"),
			Password:     getEnv("EMAIL_PASSWORD", ""),
			FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
			VerifyExpiry: verifyExpiry,
			ResetExpiry:  resetExpiry,
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
