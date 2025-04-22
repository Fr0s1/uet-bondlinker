package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"
	"socialnet/util"
)

// AuthController handles authentication-related requests
type AuthController struct {
	repo         *repository.Repository
	cfg          *config.Config
	emailService *util.EmailService
}

// NewAuthController creates a new AuthController
func NewAuthController(repo *repository.Repository, cfg *config.Config) *AuthController {
	return &AuthController{
		repo:         repo,
		cfg:          cfg,
		emailService: util.NewEmailService(cfg),
	}
}

// Register creates a new user account
func (ac *AuthController) Register(c *gin.Context) {
	var input model.UserRegistration

	if !middleware.BindJSON(c, &input) {
		return
	}

	// Check if email already exists
	existingUser, _ := ac.repo.User.FindByEmail(input.Email)
	if existingUser != nil {
		util.RespondWithError(c, http.StatusConflict, "Email already in use")
		return
	}

	// Check if username already exists
	existingUser, _ = ac.repo.User.FindByUsername(input.Username)
	if existingUser != nil {
		util.RespondWithError(c, http.StatusConflict, "Username already in use")
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
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
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate verification token
	verificationToken, err := util.GenerateEmailVerificationToken(user.ID.String(), ac.cfg.JWT.Secret, ac.cfg.Email.VerifyExpiry)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to generate verification token")
		return
	}

	// Send welcome email with verification link
	go ac.emailService.SendWelcomeEmail(user.Name, user.Email, verificationToken)

	// Generate JWT token for automatic login
	token, err := util.GenerateToken(user.ID.String(), ac.cfg.JWT.Secret, ac.cfg.JWT.ExpiryTime)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	util.RespondWithSuccess(c, http.StatusCreated, "User registered successfully", model.AuthResponse{
		Token: token,
		User:  user,
	})
}

// LoginInput represents login request data
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FCMToken string `json:"fcmToken"`
	Device   string `json:"device"`
}

// Login authenticates a user and returns a JWT token
func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Verify credentials and get user
	user, err := ac.repo.User.FindByEmail(input.Email)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !util.CheckPasswordHash(input.Password, user.Password) {
		util.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID.String(), ac.cfg.JWT.Secret)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Save FCM token if provided
	if input.FCMToken != "" && input.Device != "" {
		err = ac.repo.User.SaveFCMToken(user.ID, input.FCMToken, input.Device)
		if err != nil {
			util.RespondWithError(c, http.StatusInternalServerError, "Error saving FCM token")
			return
		}
	}

	util.RespondWithSuccess(c, http.StatusOK, "Login successful", model.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// LogoutInput represents logout request data
type LogoutInput struct {
	FCMToken string `json:"fcmToken"`
}

// Logout removes the FCM token
func (ac *AuthController) Logout(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var input LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}

	if input.FCMToken != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}

		err = ac.repo.User.RemoveFCMToken(uuid, input.FCMToken)
		if err != nil {
			util.RespondWithError(c, http.StatusInternalServerError, "Error removing FCM token")
			return
		}
	}

	util.RespondWithSuccess(c, http.StatusOK, "Logout successful", nil)
}

// ChangePassword changes a user's password
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if !middleware.BindJSON(c, &input) {
		return
	}

	// Get user ID from token
	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	// Get user from database
	user, err := ac.repo.User.FindByID(userID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify current password
	if !util.CheckPasswordHash(input.CurrentPassword, user.Password) {
		util.RespondWithError(c, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash the new password
	hashedPassword, err := util.HashPassword(input.NewPassword)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update the user's password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	err = ac.repo.User.Update(user)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Password updated successfully", nil)
}

// VerifyEmail verifies a user's email using the verification token
func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		util.RespondWithError(c, http.StatusBadRequest, "Verification token is required")
		return
	}

	// Validate the token
	claims, err := util.ParseEmailVerificationToken(token, ac.cfg.JWT.Secret)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired verification token")
		return
	}

	// Get the user ID from the token
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID in token")
		return
	}

	// Mark the user's email as verified
	user, err := ac.repo.User.FindByID(userID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Update the user's verified status
	user.EmailVerified = true
	err = ac.repo.User.Update(user)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Email verified successfully", nil)
}

// ForgotPassword initiates the password reset process
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if !middleware.BindJSON(c, &input) {
		return
	}

	// Check if user exists
	user, err := ac.repo.User.FindByEmail(input.Email)
	if err != nil {
		// Don't reveal that the email doesn't exist for security reasons
		util.RespondWithSuccess(c, http.StatusOK, "If your email is registered, you will receive a password reset link", nil)
		return
	}

	// Generate password reset token
	resetToken, err := util.GeneratePasswordResetToken(user.ID.String(), ac.cfg.JWT.Secret, ac.cfg.Email.ResetExpiry)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to generate reset token")
		return
	}

	// Send password reset email
	go ac.emailService.SendPasswordResetEmail(user.Name, user.Email, resetToken)

	util.RespondWithSuccess(c, http.StatusOK, "If your email is registered, you will receive a password reset link", nil)
}

// ResetPassword resets a user's password using the reset token
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if !middleware.BindJSON(c, &input) {
		return
	}

	// Validate the token
	claims, err := util.ParsePasswordResetToken(input.Token, ac.cfg.JWT.Secret)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired reset token")
		return
	}

	// Get the user ID from the token
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID in token")
		return
	}

	// Get the user
	user, err := ac.repo.User.FindByID(userID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Hash the new password
	hashedPassword, err := util.HashPassword(input.NewPassword)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update the user's password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	err = ac.repo.User.Update(user)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Password has been reset successfully", nil)
}
