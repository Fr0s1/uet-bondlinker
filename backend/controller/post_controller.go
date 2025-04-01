
package controller

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostController handles post-related requests
type PostController struct {
	db  *sql.DB
	cfg *config.Config
}

// NewPostController creates a new PostController
func NewPostController(db *sql.DB, cfg *config.Config) *PostController {
	return &PostController{
		db:  db,
		cfg: cfg,
	}
}

// GetPosts returns a list of posts
func (pc *PostController) GetPosts(c *gin.Context) {
	var currentUserID string
	loggedIn := false

	// Check if user is authenticated
	if id, err := middleware.GetUserID(c); err == nil {
		currentUserID = id
		loggedIn = true
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	userID := c.Query("user_id")

	// Build query
	baseQuery := `
		SELECT p.id, p.user_id, p.content, p.image, p.created_at, p.updated_at,
		(SELECT COUNT(*) FROM likes WHERE post_id = p.id) as likes,
		(SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments,
		CASE WHEN $1 = true THEN
			(SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $2 AND post_id = p.id))
		ELSE
			NULL
		END as is_liked,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
	`

	var args []interface{}
	args = append(args, loggedIn, currentUserID)

	if userID != "" {
		baseQuery += " WHERE p.user_id = $3"
		args = append(args, userID)
	}

	baseQuery += " ORDER BY p.created_at DESC LIMIT $"
	baseQuery += strconv.Itoa(len(args)+1) + " OFFSET $"
	baseQuery += strconv.Itoa(len(args)+2)
	
	args = append(args, limit, offset)

	// Execute query
	rows, err := pc.db.Query(baseQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	// Parse results
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var authorID, authorName, authorUsername string
		var authorAvatar sql.NullString
		var isLiked sql.NullBool

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Content, &post.Image, 
			&post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Comments,
			&isLiked, &authorID, &authorName, &authorUsername, &authorAvatar,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse post data"})
			return
		}

		// Set author information
		post.Author = &model.User{
			ID:       authorID,
			Name:     authorName,
			Username: authorUsername,
		}

		if authorAvatar.Valid {
			avatar := authorAvatar.String
			post.Author.Avatar = &avatar
		}

		// Set isLiked if available
		if isLiked.Valid {
			liked := isLiked.Bool
			post.IsLiked = &liked
		}

		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost returns a specific post by ID
func (pc *PostController) GetPost(c *gin.Context) {
	id := c.Param("id")
	var currentUserID string
	loggedIn := false

	// Check if user is authenticated
	if id, err := middleware.GetUserID(c); err == nil {
		currentUserID = id
		loggedIn = true
	}

	// Query post from database
	var post model.Post
	var authorID, authorName, authorUsername string
	var authorAvatar sql.NullString
	var isLiked sql.NullBool

	err := pc.db.QueryRow(
		`SELECT p.id, p.user_id, p.content, p.image, p.created_at, p.updated_at,
		(SELECT COUNT(*) FROM likes WHERE post_id = p.id) as likes,
		(SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments,
		CASE WHEN $1 = true THEN
			(SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $2 AND post_id = p.id))
		ELSE
			NULL
		END as is_liked,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $3`,
		loggedIn, currentUserID, id,
	).Scan(
		&post.ID, &post.UserID, &post.Content, &post.Image, 
		&post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Comments,
		&isLiked, &authorID, &authorName, &authorUsername, &authorAvatar,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	// Set author information
	post.Author = &model.User{
		ID:       authorID,
		Name:     authorName,
		Username: authorUsername,
	}

	if authorAvatar.Valid {
		avatar := authorAvatar.String
		post.Author.Avatar = &avatar
	}

	// Set isLiked if available
	if isLiked.Valid {
		liked := isLiked.Bool
		post.IsLiked = &liked
	}

	c.JSON(http.StatusOK, post)
}

// CreatePost creates a new post
func (pc *PostController) CreatePost(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var input model.PostCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate UUID for the post
	postID := uuid.New().String()
	now := time.Now()

	// Save post to database
	_, err = pc.db.Exec(
		"INSERT INTO posts (id, user_id, content, image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		postID, userID, input.Content, input.Image, now, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Fetch the created post with author details
	var post model.Post
	var authorID, authorName, authorUsername string
	var authorAvatar sql.NullString

	err = pc.db.QueryRow(
		`SELECT p.id, p.user_id, p.content, p.image, p.created_at, p.updated_at,
		0 as likes, 0 as comments, true as is_liked,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1`,
		postID,
	).Scan(
		&post.ID, &post.UserID, &post.Content, &post.Image, 
		&post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Comments,
		&post.IsLiked, &authorID, &authorName, &authorUsername, &authorAvatar,
	)

	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"id": postID,
			"message": "Post created successfully",
		})
		return
	}

	// Set author information
	post.Author = &model.User{
		ID:       authorID,
		Name:     authorName,
		Username: authorUsername,
	}

	if authorAvatar.Valid {
		avatar := authorAvatar.String
		post.Author.Avatar = &avatar
	}

	c.JSON(http.StatusCreated, post)
}

// UpdatePost updates an existing post
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if post exists and belongs to user
	var postOwnerID string
	err = pc.db.QueryRow("SELECT user_id FROM posts WHERE id = $1", id).Scan(&postOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if postOwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another user's post"})
		return
	}

	var input model.PostUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update post in database
	now := time.Now()
	_, err = pc.db.Exec(
		"UPDATE posts SET content = $1, image = $2, updated_at = $3 WHERE id = $4",
		input.Content, input.Image, now, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Fetch the updated post with author details
	var post model.Post
	var authorID, authorName, authorUsername string
	var authorAvatar sql.NullString
	var isLiked bool

	err = pc.db.QueryRow(
		`SELECT p.id, p.user_id, p.content, p.image, p.created_at, p.updated_at,
		(SELECT COUNT(*) FROM likes WHERE post_id = p.id) as likes,
		(SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments,
		(SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = p.id)) as is_liked,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $2`,
		userID, id,
	).Scan(
		&post.ID, &post.UserID, &post.Content, &post.Image, 
		&post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Comments,
		&isLiked, &authorID, &authorName, &authorUsername, &authorAvatar,
	)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id": id,
			"message": "Post updated successfully",
		})
		return
	}

	// Set author information
	post.Author = &model.User{
		ID:       authorID,
		Name:     authorName,
		Username: authorUsername,
	}

	if authorAvatar.Valid {
		avatar := authorAvatar.String
		post.Author.Avatar = &avatar
	}

	post.IsLiked = &isLiked

	c.JSON(http.StatusOK, post)
}

// DeletePost deletes a post
func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if post exists and belongs to user
	var postOwnerID string
	err = pc.db.QueryRow("SELECT user_id FROM posts WHERE id = $1", id).Scan(&postOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if postOwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete another user's post"})
		return
	}

	// Delete post from database
	_, err = pc.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// LikePost adds a like to a post
func (pc *PostController) LikePost(c *gin.Context) {
	postID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if post exists
	var exists bool
	err = pc.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", postID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if already liked
	err = pc.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)",
		userID, postID,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Post already liked"})
		return
	}

	// Add like to database
	_, err = pc.db.Exec(
		"INSERT INTO likes (user_id, post_id, created_at) VALUES ($1, $2, $3)",
		userID, postID, time.Now(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	// Get updated like count
	var likeCount int
	err = pc.db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = $1", postID).Scan(&likeCount)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"message": "Post liked successfully"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post liked successfully",
		"likes": likeCount,
	})
}

// UnlikePost removes a like from a post
func (pc *PostController) UnlikePost(c *gin.Context) {
	postID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if like exists
	var exists bool
	err = pc.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)",
		userID, postID,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not liked"})
		return
	}

	// Remove like from database
	_, err = pc.db.Exec(
		"DELETE FROM likes WHERE user_id = $1 AND post_id = $2",
		userID, postID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		return
	}

	// Get updated like count
	var likeCount int
	err = pc.db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = $1", postID).Scan(&likeCount)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
		"likes": likeCount,
	})
}

// GetFeed returns posts from users that the authenticated user follows
func (pc *PostController) GetFeed(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get feed posts (posts from followed users and own posts)
	rows, err := pc.db.Query(
		`SELECT p.id, p.user_id, p.content, p.image, p.created_at, p.updated_at,
		(SELECT COUNT(*) FROM likes WHERE post_id = p.id) as likes,
		(SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments,
		(SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = p.id)) as is_liked,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id IN (
			SELECT following_id FROM follows WHERE follower_id = $1
		) OR p.user_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feed"})
		return
	}
	defer rows.Close()

	// Parse results
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var authorID, authorName, authorUsername string
		var authorAvatar sql.NullString
		var isLiked bool

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Content, &post.Image, 
			&post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Comments,
			&isLiked, &authorID, &authorName, &authorUsername, &authorAvatar,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse post data"})
			return
		}

		// Set author information
		post.Author = &model.User{
			ID:       authorID,
			Name:     authorName,
			Username: authorUsername,
		}

		if authorAvatar.Valid {
			avatar := authorAvatar.String
			post.Author.Avatar = &avatar
		}

		post.IsLiked = &isLiked

		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

// GetComments returns comments for a specific post
func (pc *PostController) GetComments(c *gin.Context) {
	postID := c.Param("id")
	var currentUserID string
	loggedIn := false

	// Check if user is authenticated
	if id, err := middleware.GetUserID(c); err == nil {
		currentUserID = id
		loggedIn = true
	}

	// Check if post exists
	var exists bool
	err := pc.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", postID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get comments for post
	rows, err := pc.db.Query(
		`SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, c.updated_at,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3`,
		postID, limit, offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	// Parse results
	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		var authorID, authorName, authorUsername string
		var authorAvatar sql.NullString

		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.PostID, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt,
			&authorID, &authorName, &authorUsername, &authorAvatar,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse comment data"})
			return
		}

		// Set author information
		comment.Author = &model.User{
			ID:       authorID,
			Name:     authorName,
			Username: authorUsername,
		}

		if authorAvatar.Valid {
			avatar := authorAvatar.String
			comment.Author.Avatar = &avatar
		}

		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}

// CreateComment adds a comment to a post
func (pc *PostController) CreateComment(c *gin.Context) {
	postID := c.Param("id")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if post exists
	var exists bool
	err = pc.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", postID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var input model.CommentCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate UUID for the comment
	commentID := uuid.New().String()
	now := time.Now()

	// Save comment to database
	_, err = pc.db.Exec(
		"INSERT INTO comments (id, user_id, post_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		commentID, userID, postID, input.Content, now, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Fetch the created comment with author details
	var comment model.Comment
	var authorID, authorName, authorUsername string
	var authorAvatar sql.NullString

	err = pc.db.QueryRow(
		`SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, c.updated_at,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1`,
		commentID,
	).Scan(
		&comment.ID, &comment.UserID, &comment.PostID, &comment.Content,
		&comment.CreatedAt, &comment.UpdatedAt,
		&authorID, &authorName, &authorUsername, &authorAvatar,
	)

	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"id": commentID,
			"message": "Comment created successfully",
		})
		return
	}

	// Set author information
	comment.Author = &model.User{
		ID:       authorID,
		Name:     authorName,
		Username: authorUsername,
	}

	if authorAvatar.Valid {
		avatar := authorAvatar.String
		comment.Author.Avatar = &avatar
	}

	c.JSON(http.StatusCreated, comment)
}

// UpdateComment updates an existing comment
func (pc *PostController) UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if comment exists and belongs to user
	var commentOwnerID string
	err = pc.db.QueryRow("SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&commentOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if commentOwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another user's comment"})
		return
	}

	var input model.CommentUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update comment in database
	now := time.Now()
	_, err = pc.db.Exec(
		"UPDATE comments SET content = $1, updated_at = $2 WHERE id = $3",
		input.Content, now, commentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// Fetch the updated comment with author details
	var comment model.Comment
	var authorID, authorName, authorUsername string
	var authorAvatar sql.NullString

	err = pc.db.QueryRow(
		`SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, c.updated_at,
		u.id as author_id, u.name as author_name, u.username as author_username, 
		u.avatar as author_avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1`,
		commentID,
	).Scan(
		&comment.ID, &comment.UserID, &comment.PostID, &comment.Content,
		&comment.CreatedAt, &comment.UpdatedAt,
		&authorID, &authorName, &authorUsername, &authorAvatar,
	)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id": commentID,
			"message": "Comment updated successfully",
		})
		return
	}

	// Set author information
	comment.Author = &model.User{
		ID:       authorID,
		Name:     authorName,
		Username: authorUsername,
	}

	if authorAvatar.Valid {
		avatar := authorAvatar.String
		comment.Author.Avatar = &avatar
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment deletes a comment
func (pc *PostController) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Check if comment exists and belongs to user
	var commentOwnerID string
	err = pc.db.QueryRow("SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&commentOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if commentOwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete another user's comment"})
		return
	}

	// Delete comment from database
	_, err = pc.db.Exec("DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
