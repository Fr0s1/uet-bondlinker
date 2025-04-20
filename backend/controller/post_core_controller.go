package controller

import (
	"net/http"
	"socialnet/util"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostController handles post-related requests
type PostController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewPostController creates a new PostController
func NewPostController(repo *repository.Repository, cfg *config.Config) *PostController {
	return &PostController{
		repo: repo,
		cfg:  cfg,
	}
}

// GetPosts returns a list of posts
func (pc *PostController) GetPosts(c *gin.Context) {
	var filter model.PostFilter
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

	// Query posts from database
	posts, err := pc.repo.Post.FindAll(filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	if posts, err = pc.repo.Post.FillLikeInfo(currentUserID, posts); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch like info")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", posts)
}

// GetPost returns a specific post by ID
func (pc *PostController) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	var currentUserID *uuid.UUID

	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}

	// Query post from database
	post, err := pc.repo.Post.FindByID(id)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	// If user is authenticated, check if post is liked
	if currentUserID != nil {
		isLiked, _ := pc.repo.Post.IsLiked(*currentUserID, post.ID)
		post.IsLiked = &isLiked
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", post)
}

// CreatePost creates a new post
func (pc *PostController) CreatePost(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	var input model.PostCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create new post
	post := model.Post{
		ID:      uuid.New(),
		UserID:  userID,
		Content: input.Content,
		Image:   input.Image,
	}

	// Save post to database
	err = pc.repo.Post.Create(&post)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Get the created post with author details
	createdPost, err := pc.repo.Post.FindByID(post.ID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "failed to fetch post")
		return
	}

	// Set is_liked to true since user just created it
	isLiked := false
	createdPost.IsLiked = &isLiked

	c.JSON(http.StatusCreated, createdPost)
}

// UpdatePost updates an existing post
func (pc *PostController) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	// Check if post exists
	post, err := pc.repo.Post.FindByID(id)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	// Check if user owns the post
	if post.UserID != userID {
		util.RespondWithError(c, http.StatusForbidden, "Cannot update another user's post")
		return
	}

	var input model.PostUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update post fields
	post.Content = input.Content
	post.Image = input.Image

	// Save updated post to database
	err = pc.repo.Post.Update(post)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to update post")
		return
	}

	// Check if post is liked
	isLiked, _ := pc.repo.Post.IsLiked(userID, post.ID)
	post.IsLiked = &isLiked

	util.RespondWithSuccess(c, http.StatusOK, "success", post)
}

// DeletePost deletes a post
func (pc *PostController) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	// Check if post exists
	post, err := pc.repo.Post.FindByID(id)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	// Check if user owns the post
	if post.UserID != userID {
		util.RespondWithError(c, http.StatusForbidden, "Cannot delete another user's post")
		return
	}

	// Delete post from database
	err = pc.repo.Post.Delete(id)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Post deleted successfully", nil)
}

// GetFeed returns posts from users that the authenticated user follows
func (pc *PostController) GetFeed(c *gin.Context) {
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

	// Get feed posts (posts from followed users and own posts)
	posts, err := pc.repo.Post.FindFeed(userID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch feed")
		return
	}

	if posts, err = pc.repo.Post.FillLikeInfo(&userID, posts); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch like info")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", posts)
}

// GetTrending returns trending posts based on engagement
func (pc *PostController) GetTrending(c *gin.Context) {
	var filter model.Pagination
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

	// Get trending posts
	posts, err := pc.repo.Post.FindTrending(filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch trending posts")
		return
	}

	if posts, err = pc.repo.Post.FillLikeInfo(currentUserID, posts); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch like info")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", posts)
}

// GetSuggestedPosts returns posts that might interest the user
func (pc *PostController) GetSuggestedPosts(c *gin.Context) {
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

	// Get suggested posts
	posts, err := pc.repo.Post.GetSuggestedPosts(userID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch suggested posts")
		return
	}

	if posts, err = pc.repo.Post.FillLikeInfo(&userID, posts); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch like info")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", posts)
}
