
package controller

import (
	"database/sql"
	"net/http"
	"time"

	"socialnet/config"
	"socialnet/model"
	"socialnet/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthController handles authentication-related requests
type AuthController struct {
	db  *sql.DB
	cfg *config.Config
}

// NewAuthController creates a new AuthController
func NewAuthController(db *sql.DB, cfg *config.Config) *AuthController {
	return &AuthController{
		db:  db,
		cfg: cfg,
	}
}

// Register creates a new user account
func (ac *AuthController) Register(c *gin.Context) {
	var input model.UserRegistration

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var count int
	err := ac.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", input.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	// Check if username already exists
	err = ac.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", input.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate UUID for the user
	userID := uuid.New().String()

	// Save user to database
	now := time.Now()
	_, err = ac.db.Exec(
		`INSERT INTO users (id, name, username, email, password_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, input.Name, input.Username, input.Email, hashedPassword, now, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(userID, ac.cfg.JWT.Secret, ac.cfg.JWT.ExpiryTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return response
	user := model.User{
		ID:        userID,
		Name:      input.Name,
		Username:  input.Username,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	c.JSON(http.StatusCreated, model.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login authenticates a user and returns a JWT token
func (ac *AuthController) Login(c *gin.Context) {
	var input model.UserLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	var user model.User
	var passwordHash string

	err := ac.db.QueryRow(
		`SELECT id, name, username, email, password_hash, created_at, updated_at 
		FROM users WHERE email = $1`,
		input.Email,
	).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email,
		&passwordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify password
	if !util.CheckPasswordHash(input.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID, ac.cfg.JWT.Secret, ac.cfg.JWT.ExpiryTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.AuthResponse{
		Token: token,
		User:  user,
	})
}
