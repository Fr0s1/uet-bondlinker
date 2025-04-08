
package controller

import (
	"net/http"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchController handles search-related requests
type SearchController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewSearchController creates a new SearchController
func NewSearchController(repo *repository.Repository, cfg *config.Config) *SearchController {
	return &SearchController{
		repo: repo,
		cfg:  cfg,
	}
}

// SearchUsers searches for users by name or username
func (sc *SearchController) SearchUsers(c *gin.Context) {
	var filter model.SearchFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentUserID *uuid.UUID

	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}

	// Search users in database
	users, err := sc.repo.User.SearchUsers(filter.Query, filter.Pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// If user is authenticated, check follow status for each user
	if currentUserID != nil {
		for i := range users {
			isFollowed, _ := sc.repo.User.IsFollowing(*currentUserID, users[i].ID)
			users[i].IsFollowed = &isFollowed
		}
	}

	c.JSON(http.StatusOK, users)
}

// SearchPosts searches for posts by content
func (sc *SearchController) SearchPosts(c *gin.Context) {
	var filter model.SearchFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentUserID *uuid.UUID

	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}

	// Search posts in database
	posts, err := sc.repo.Post.SearchPosts(filter.Query, filter.Pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}

	// If user is authenticated, check if posts are liked
	if currentUserID != nil {
		for i := range posts {
			isLiked, _ := sc.repo.Post.IsLiked(*currentUserID, posts[i].ID)
			posts[i].IsLiked = &isLiked
		}
	}

	c.JSON(http.StatusOK, posts)
}

// Search performs a combined search across users and posts
func (sc *SearchController) Search(c *gin.Context) {
	var filter model.SearchFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentUserID *uuid.UUID

	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}

	// Search users
	users, err := sc.repo.User.SearchUsers(filter.Query, filter.Pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// If user is authenticated, check follow status for each user
	if currentUserID != nil {
		for i := range users {
			isFollowed, _ := sc.repo.User.IsFollowing(*currentUserID, users[i].ID)
			users[i].IsFollowed = &isFollowed
		}
	}

	// Search posts
	posts, err := sc.repo.Post.SearchPosts(filter.Query, filter.Pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}

	// If user is authenticated, check if posts are liked
	if currentUserID != nil {
		for i := range posts {
			isLiked, _ := sc.repo.Post.IsLiked(*currentUserID, posts[i].ID)
			posts[i].IsLiked = &isLiked
		}
	}

	// Return combined results
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"posts": posts,
	})
}
