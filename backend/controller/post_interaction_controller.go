package controller

import (
	"net/http"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"
	"socialnet/util"

	"github.com/gin-gonic/gin"
)

// PostInteractionController handles post interaction-related requests (likes, shares)
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
	_, err = pic.repo.Post.FindByID(postID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	// Check if already liked
	isLiked, err := pic.repo.Post.IsLiked(userID, postID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, util.ErrorMessages.DatabaseError)
		return
	}

	if isLiked {
		util.RespondWithError(c, http.StatusConflict, "Post already liked")
		return
	}

	// Add like to database
	likeCount, err := pic.repo.Post.Like(userID, postID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to like post")
		return
	}

	util.RespondWithSuccess(c, http.StatusCreated, "Post liked successfully", gin.H{
		"likes": likeCount,
	})
}

// UnlikePost removes a like from a post
func (pic *PostInteractionController) UnlikePost(c *gin.Context) {
	postID, err := middleware.ParseUUIDParam(c, "id")
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	// Check if like exists
	isLiked, err := pic.repo.Post.IsLiked(userID, postID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, util.ErrorMessages.DatabaseError)
		return
	}

	if !isLiked {
		util.RespondWithError(c, http.StatusNotFound, "Post not liked")
		return
	}

	// Remove like from database
	likeCount, err := pic.repo.Post.Unlike(userID, postID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to unlike post")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Post unliked successfully", gin.H{
		"likes": likeCount,
	})
}

// SharePost shares an existing post
func (pic *PostInteractionController) SharePost(c *gin.Context) {
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
	_, err = pic.repo.Post.FindByID(postID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	var input model.PostShare
	if !middleware.BindJSON(c, &input) {
		return
	}

	// Share the post
	sharedPost, err := pic.repo.Post.Share(userID, postID, input.Content)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to share post")
		return
	}

	util.RespondWithSuccess(c, http.StatusCreated, "Post shared successfully", sharedPost)
}
