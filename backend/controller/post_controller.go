
package controller

import (
	"net/http"
	"strconv"

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
	var currentUserID *uuid.UUID
	
	// Check if user is authenticated
	if userIDStr, err := middleware.GetUserID(c); err == nil {
		userID, _ := uuid.Parse(userIDStr)
		currentUserID = &userID
	}
	
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	var filterUserID *uuid.UUID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			filterUserID = &userID
		}
	}
	
	// Query posts from database
	posts, err := pc.repo.Post.FindAll(filterUserID, limit, offset)
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

// LikePost adds a like to a post
func (pc *PostController) LikePost(c *gin.Context) {
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
	_, err = pc.repo.Post.FindByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Check if already liked
	isLiked, err := pc.repo.Post.IsLiked(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if isLiked {
		c.JSON(http.StatusConflict, gin.H{"error": "Post already liked"})
		return
	}
	
	// Add like to database
	err = pc.repo.Post.Like(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}
	
	// Get updated like count
	likeCount, err := pc.repo.Post.CountLikes(postID)
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
func (pc *PostController) UnlikePost(c *gin.Context) {
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
	isLiked, err := pc.repo.Post.IsLiked(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if !isLiked {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not liked"})
		return
	}
	
	// Remove like from database
	err = pc.repo.Post.Unlike(userID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		return
	}
	
	// Get updated like count
	likeCount, err := pc.repo.Post.CountLikes(postID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
		"likes":   likeCount,
	})
}

// GetFeed returns posts from users that the authenticated user follows
func (pc *PostController) GetFeed(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	userID, _ := uuid.Parse(userIDStr)
	
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	// Get feed posts (posts from followed users and own posts)
	posts, err := pc.repo.Post.FindFeed(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feed"})
		return
	}
	
	c.JSON(http.StatusOK, posts)
}

// GetComments returns comments for a specific post
func (pc *PostController) GetComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}
	
	// Check if post exists
	_, err = pc.repo.Post.FindByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	// Get comments for post
	comments, err := pc.repo.Comment.FindByPostID(postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	
	c.JSON(http.StatusOK, comments)
}

// CreateComment adds a comment to a post
func (pc *PostController) CreateComment(c *gin.Context) {
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
	_, err = pc.repo.Post.FindByID(postID)
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
	err = pc.repo.Comment.Create(&comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	
	// Get the created comment with author details
	createdComment, err := pc.repo.Comment.FindByID(comment.ID)
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
func (pc *PostController) UpdateComment(c *gin.Context) {
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
	comment, err := pc.repo.Comment.FindByID(commentID)
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
	err = pc.repo.Comment.Update(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}
	
	c.JSON(http.StatusOK, comment)
}

// DeleteComment deletes a comment
func (pc *PostController) DeleteComment(c *gin.Context) {
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
	comment, err := pc.repo.Comment.FindByID(commentID)
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
	err = pc.repo.Comment.Delete(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
