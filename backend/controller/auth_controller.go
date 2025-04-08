
package controller

import (
	"net/http"
	"socialnet/config"
	"socialnet/model"
	"socialnet/repository"
	"socialnet/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthController handles authentication-related requests
type AuthController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewAuthController creates a new AuthController
func NewAuthController(repo *repository.Repository, cfg *config.Config) *AuthController {
	return &AuthController{
		repo: repo,
		cfg:  cfg,
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
	existingUser, _ := ac.repo.User.FindByEmail(input.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	// Check if username already exists
	existingUser, _ = ac.repo.User.FindByUsername(input.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	user := model.User{
		ID:       uuid.New(),
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Save user to database
	err = ac.repo.User.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID.String(), ac.cfg.JWT.Secret, ac.cfg.JWT.ExpiryTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
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
	user, err := ac.repo.User.FindByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if !util.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID.String(), ac.cfg.JWT.Secret, ac.cfg.JWT.ExpiryTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.AuthResponse{
		Token: token,
		User:  *user,
	})
}
