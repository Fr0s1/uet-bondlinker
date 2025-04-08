
package controller

import (
	"net/http"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostInteractionController handles post interaction-related requests (likes)
type PostInteractionController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewPostInteractionController creates a new PostInteractionController
func NewPostInteractionController(repo *repository.Repository, cfg *config.Config) *PostInteractionController {
	return &PostInteractionController{
		repo: repo,
		cfg:  cfg,
	}
}

// LikePost adds a like to a post
func (pic *PostInteractionController) LikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}
	
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	userID, _ := uuid.Parse(userIDStr)
	
	// Check if post exists
	_, err = pic.repo.Post.FindByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Check if already liked
	isLiked, err := pic.repo.Post.IsLiked(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if isLiked {
		c.JSON(http.StatusConflict, gin.H{"error": "Post already liked"})
		return
	}
	
	// Add like to database
	err = pic.repo.Post.Like(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}
	
	// Get updated like count
	likeCount, err := pic.repo.Post.CountLikes(postID)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"message": "Post liked successfully"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post liked successfully",
		"likes":   likeCount,
	})
}

// UnlikePost removes a like from a post
func (pic *PostInteractionController) UnlikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}
	
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	userID, _ := uuid.Parse(userIDStr)
	
	// Check if like exists
	isLiked, err := pic.repo.Post.IsLiked(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if !isLiked {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not liked"})
		return
	}
	
	// Remove like from database
	err = pic.repo.Post.Unlike(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		return
	}
	
	// Get updated like count
	likeCount, err := pic.repo.Post.CountLikes(postID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
		"likes":   likeCount,
	})
}
