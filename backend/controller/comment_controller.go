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

// CommentController handles comment-related requests
type CommentController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewCommentController creates a new CommentController
func NewCommentController(repo *repository.Repository, cfg *config.Config) *CommentController {
	return &CommentController{
		repo: repo,
		cfg:  cfg,
	}
}

// GetComments returns comments for a specific post
func (cc *CommentController) GetComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	// Check if post exists
	_, err = cc.repo.Post.FindByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var filter model.Pagination
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get comments for post
	comments, err := cc.repo.Comment.FindByPostID(postID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// CreateComment adds a comment to a post
func (cc *CommentController) CreateComment(c *gin.Context) {
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
	_, err = cc.repo.Post.FindByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var input model.CommentCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new comment
	comment := model.Comment{
		ID:      uuid.New(),
		UserID:  userID,
		PostID:  postID,
		Content: input.Content,
	}

	// Save comment to database
	err = cc.repo.Comment.Create(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Get the created comment with author details
	createdComment, err := cc.repo.Comment.FindByID(comment.ID)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"id":      comment.ID.String(),
			"message": "Comment created successfully",
		})
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

// UpdateComment updates an existing comment
func (cc *CommentController) UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	// Check if comment exists
	comment, err := cc.repo.Comment.FindByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if user owns the comment
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another user's comment"})
		return
	}

	var input model.CommentUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update comment fields
	comment.Content = input.Content

	// Save updated comment to database
	err = cc.repo.Comment.Update(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment deletes a comment
func (cc *CommentController) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	// Check if comment exists
	comment, err := cc.repo.Comment.FindByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if user owns the comment
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete another user's comment"})
		return
	}

	// Delete comment from database
	err = cc.repo.Comment.Delete(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
