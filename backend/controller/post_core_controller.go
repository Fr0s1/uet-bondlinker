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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// If user is authenticated, check if posts are liked
	if currentUserID != nil {
		for i := range posts {
			isLiked, _ := pc.repo.Post.IsLiked(*currentUserID, posts[i].ID)
			posts[i].IsLiked = &isLiked
		}
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost returns a specific post by ID
func (pc *PostController) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// If user is authenticated, check if post is liked
	if currentUserID != nil {
		isLiked, _ := pc.repo.Post.IsLiked(*currentUserID, post.ID)
		post.IsLiked = &isLiked
	}

	c.JSON(http.StatusOK, post)
}

// CreatePost creates a new post
func (pc *PostController) CreatePost(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	var input model.PostCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Get the created post with author details
	createdPost, err := pc.repo.Post.FindByID(post.ID)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"id":      post.ID.String(),
			"message": "Post created successfully",
		})
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
	post, err := pc.repo.Post.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user owns the post
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another user's post"})
		return
	}

	var input model.PostUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update post fields
	post.Content = input.Content
	post.Image = input.Image

	// Save updated post to database
	err = pc.repo.Post.Update(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Check if post is liked
	isLiked, _ := pc.repo.Post.IsLiked(userID, post.ID)
	post.IsLiked = &isLiked

	c.JSON(http.StatusOK, post)
}

// DeletePost deletes a post
func (pc *PostController) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
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
	post, err := pc.repo.Post.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user owns the post
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete another user's post"})
		return
	}

	// Delete post from database
	err = pc.repo.Post.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetFeed returns posts from users that the authenticated user follows
func (pc *PostController) GetFeed(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	var filter model.Pagination
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get feed posts (posts from followed users and own posts)
	posts, err := pc.repo.Post.FindFeed(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feed"})
		return
	}

	// Check if posts are liked by the current user
	for i := range posts {
		isLiked, _ := pc.repo.Post.IsLiked(userID, posts[i].ID)
		posts[i].IsLiked = &isLiked
	}

	c.JSON(http.StatusOK, posts)
}
