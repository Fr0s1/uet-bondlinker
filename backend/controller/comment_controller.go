package controller

import (
	"net/http"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"
	"socialnet/util"

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
	postID, err := middleware.ParseUUIDParam(c, "id")
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	// Check if post exists
	_, err = cc.repo.Post.FindByID(postID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	var filter model.Pagination
	if !middleware.BindQuery(c, &filter) {
		return
	}

	// Get comments for post
	comments, err := cc.repo.Comment.FindByPostID(postID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Comments retrieved successfully", comments)
}

// CreateComment adds a comment to a post
func (cc *CommentController) CreateComment(c *gin.Context) {
	postID, err := middleware.ParseUUIDParam(c, "id")
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	// Check if post exists
	_, err = cc.repo.Post.FindByID(postID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	var input model.CommentCreate
	if !middleware.BindJSON(c, &input) {
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
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// Get the created comment with author details
	createdComment, err := cc.repo.Comment.FindByID(comment.ID)
	if err != nil {
		util.RespondWithSuccess(c, http.StatusCreated, "Comment created successfully", gin.H{
			"id": comment.ID.String(),
		})
		return
	}

	util.RespondWithSuccess(c, http.StatusCreated, "Comment created successfully", createdComment)
}

// UpdateComment updates an existing comment
func (cc *CommentController) UpdateComment(c *gin.Context) {
	commentID, err := middleware.ParseUUIDParam(c, "commentId")
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid comment ID format")
		return
	}

	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	// Check if comment exists
	comment, err := cc.repo.Comment.FindByID(commentID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Comment not found")
		return
	}

	// Check if user owns the comment
	if !middleware.CheckResourceOwnership(c, comment.UserID, userID) {
		return
	}

	var input model.CommentUpdate
	if !middleware.BindJSON(c, &input) {
		return
	}

	// Update comment fields
	comment.Content = input.Content

	// Save updated comment to database
	err = cc.repo.Comment.Update(comment)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to update comment")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Comment updated successfully", comment)
}

// DeleteComment deletes a comment
func (cc *CommentController) DeleteComment(c *gin.Context) {
	commentID, err := middleware.ParseUUIDParam(c, "commentId")
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid comment ID format")
		return
	}

	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	// Check if comment exists
	comment, err := cc.repo.Comment.FindByID(commentID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Comment not found")
		return
	}

	// Check if user owns the comment
	if !middleware.CheckResourceOwnership(c, comment.UserID, userID) {
		return
	}

	// Delete comment from database
	err = cc.repo.Comment.Delete(commentID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Comment deleted successfully", nil)
}
