package controller

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"socialnet/util"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles user-related requests
type UserController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewUserController creates a new UserController
func NewUserController(repo *repository.Repository, cfg *config.Config) *UserController {
	return &UserController{
		repo: repo,
		cfg:  cfg,
	}
}

// GetUsers returns a list of users
func (uc *UserController) GetUsers(c *gin.Context) {
	var filter model.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var currentUserID *uuid.UUID

	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}

	// Query users from database
	users, err := uc.repo.User.FindAll(filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// If user is authenticated, check follow status for each user
	if currentUserID != nil {
		for i := range users {
			isFollowed, _ := uc.repo.User.IsFollowing(*currentUserID, users[i].ID)
			users[i].IsFollowed = &isFollowed
		}
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", users)
}

// GetUser returns a specific user by ID
func (uc *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var currentUserID *uuid.UUID

	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}

	// Query user from database
	user, err := uc.repo.User.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.RespondWithError(c, http.StatusNotFound, "User not found")
			return
		}
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// If user is authenticated, check follow status
	if currentUserID != nil {
		isFollowed, _ := uc.repo.User.IsFollowing(*currentUserID, user.ID)
		user.IsFollowed = &isFollowed
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", user)
}

// GetUserByUsername returns a specific user by username
func (uc *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	var currentUserID *uuid.UUID

	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}

	// Query user from database
	user, err := uc.repo.User.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.RespondWithError(c, http.StatusNotFound, "User not found")
			return
		}
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// If user is authenticated, check follow status
	if currentUserID != nil {
		isFollowed, _ := uc.repo.User.IsFollowing(*currentUserID, user.ID)
		user.IsFollowed = &isFollowed
	}

	util.RespondWithSuccess(c, http.StatusOK, "User found", user)
}

// GetCurrentUser returns the authenticated user
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Query user from database
	user, err := uc.repo.User.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.RespondWithError(c, http.StatusNotFound, "User not found")
			return
		}
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "User found", user)
}

// UpdateUser updates a user's profile
func (uc *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	// Check if user is updating their own profile
	if id != userID {
		util.RespondWithError(c, http.StatusForbidden, "Cannot update another user's profile")
		return
	}

	var input model.UserUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get existing user
	user, err := uc.repo.User.FindByID(id)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Update fields if provided
	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Bio != nil {
		user.Bio = input.Bio
	}

	if input.Avatar != nil {
		user.Avatar = input.Avatar
	}

	if input.Cover != nil {
		user.Cover = input.Cover
	}

	if input.Location != nil {
		user.Location = input.Location
	}

	if input.Website != nil {
		user.Website = input.Website
	}

	// Update user in database
	err = uc.repo.User.Update(user)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", user)
}

// FollowUser creates a follow relationship between users
func (uc *UserController) FollowUser(c *gin.Context) {
	followingIDStr := c.Param("id")
	followingID, err := uuid.Parse(followingIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	followerIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	followerID, _ := uuid.Parse(followerIDStr)

	// Check if trying to follow oneself
	if followerID == followingID {
		util.RespondWithError(c, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	// Check if user to follow exists
	_, err = uc.repo.User.FindByID(followingID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "User to follow not found")
		return
	}

	// Check if already following
	isFollowing, err := uc.repo.User.IsFollowing(followerID, followingID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	if isFollowing {
		util.RespondWithError(c, http.StatusConflict, "Already following this user")
		return
	}

	// Create follow relationship
	err = uc.repo.User.Follow(followerID, followingID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to follow user")
		return
	}

	util.RespondWithSuccess(c, http.StatusCreated, "Successfully followed user", nil)
}

// UnfollowUser removes a follow relationship between users
func (uc *UserController) UnfollowUser(c *gin.Context) {
	followingIDStr := c.Param("id")
	followingID, err := uuid.Parse(followingIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	followerIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	followerID, _ := uuid.Parse(followerIDStr)

	// Check if follow relationship exists
	isFollowing, err := uc.repo.User.IsFollowing(followerID, followingID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	if !isFollowing {
		util.RespondWithError(c, http.StatusNotFound, "Not following this user")
		return
	}

	// Remove follow relationship
	err = uc.repo.User.Unfollow(followerID, followingID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to unfollow user")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", gin.H{"message": "Successfully unfollowed user"})
}

// GetFollowers returns users who follow the specified user
func (uc *UserController) GetFollowers(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var filter model.FollowFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Query followers from database
	followers, err := uc.repo.User.GetFollowers(userID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch followers")
		return
	}

	// Check if the requesting user is authenticated
	if currentUserIDStr, err := middleware.GetUserID(c); err == nil {
		currentUserID, _ := uuid.Parse(currentUserIDStr)

		// Mark which users the current user is following
		for i := range followers {
			isFollowed, _ := uc.repo.User.IsFollowing(currentUserID, followers[i].ID)
			followers[i].IsFollowed = &isFollowed
		}
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", followers)
}

// GetFollowing returns users that the specified user follows
func (uc *UserController) GetFollowing(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var filter model.FollowFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Query following from database
	following, err := uc.repo.User.GetFollowing(userID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch following")
		return
	}

	// Check if the requesting user is authenticated
	if currentUserIDStr, err := middleware.GetUserID(c); err == nil {
		currentUserID, _ := uuid.Parse(currentUserIDStr)

		// Mark which users the current user is following
		for i := range following {
			isFollowed, _ := uc.repo.User.IsFollowing(currentUserID, following[i].ID)
			following[i].IsFollowed = &isFollowed
		}
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", following)
}

// GetSuggestedUsers returns users that might interest the user
func (uc *UserController) GetSuggestedUsers(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	var filter model.Pagination
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get suggested users
	suggestedUsers, err := uc.repo.User.GetSuggestedUsers(userID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch suggested users")
		return
	}

	// Mark which users the current user is following
	for i := range suggestedUsers {
		isFollowed, _ := uc.repo.User.IsFollowing(userID, suggestedUsers[i].ID)
		suggestedUsers[i].IsFollowed = &isFollowed
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", suggestedUsers)
}
