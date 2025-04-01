
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

// UserController handles user-related requests
type UserController struct {
	db  *sql.DB
	cfg *config.Config
}

// NewUserController creates a new UserController
func NewUserController(db *sql.DB, cfg *config.Config) *UserController {
	return &UserController{
		db:  db,
		cfg: cfg,
	}
}

// GetUsers returns a list of users
func (uc *UserController) GetUsers(c *gin.Context) {
	var currentUserID string
	loggedIn := false

	// Check if user is authenticated
	if id, err := middleware.GetUserID(c); err == nil {
		currentUserID = id
		loggedIn = true
	}

	// Parse query parameters
	query := c.DefaultQuery("q", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Query users from database
	rows, err := uc.db.Query(
		`SELECT u.id, u.name, u.username, u.email, u.bio, u.avatar, u.location, u.website, 
		u.created_at, u.updated_at, 
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following,
		CASE WHEN $4 = true THEN
			(SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $5 AND following_id = u.id))
		ELSE
			NULL
		END as is_followed
		FROM users u
		WHERE ($1 = '' OR u.name ILIKE '%' || $1 || '%' OR u.username ILIKE '%' || $1 || '%')
		ORDER BY u.created_at DESC
		LIMIT $2 OFFSET $3`,
		query, limit, offset, loggedIn, currentUserID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	// Parse results
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Username, &user.Email,
			&user.Bio, &user.Avatar, &user.Location, &user.Website,
			&user.CreatedAt, &user.UpdatedAt,
			&user.Followers, &user.Following, &user.IsFollowed,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
			return
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// GetUser returns a specific user by ID
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	var currentUserID string
	loggedIn := false

	// Check if user is authenticated
	if id, err := middleware.GetUserID(c); err == nil {
		currentUserID = id
		loggedIn = true
	}

	// Query user from database
	var user model.User
	err := uc.db.QueryRow(
		`SELECT u.id, u.name, u.username, u.email, u.bio, u.avatar, u.location, u.website, 
		u.created_at, u.updated_at, 
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following,
		CASE WHEN $2 = true THEN
			(SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $3 AND following_id = u.id))
		ELSE
			NULL
		END as is_followed
		FROM users u
		WHERE u.id = $1`,
		id, loggedIn, currentUserID,
	).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email,
		&user.Bio, &user.Avatar, &user.Location, &user.Website,
		&user.CreatedAt, &user.UpdatedAt,
		&user.Followers, &user.Following, &user.IsFollowed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByUsername returns a specific user by username
func (uc *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	var currentUserID string
	loggedIn := false

	// Check if user is authenticated
	if id, err := middleware.GetUserID(c); err == nil {
		currentUserID = id
		loggedIn = true
	}

	// Query user from database
	var user model.User
	err := uc.db.QueryRow(
		`SELECT u.id, u.name, u.username, u.email, u.bio, u.avatar, u.location, u.website, 
		u.created_at, u.updated_at, 
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following,
		CASE WHEN $2 = true THEN
			(SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $3 AND following_id = u.id))
		ELSE
			NULL
		END as is_followed
		FROM users u
		WHERE u.username = $1`,
		username, loggedIn, currentUserID,
	).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email,
		&user.Bio, &user.Avatar, &user.Location, &user.Website,
		&user.CreatedAt, &user.UpdatedAt,
		&user.Followers, &user.Following, &user.IsFollowed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetCurrentUser returns the authenticated user
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Query user from database
	var user model.User
	err = uc.db.QueryRow(
		`SELECT u.id, u.name, u.username, u.email, u.bio, u.avatar, u.location, u.website, 
		u.created_at, u.updated_at, 
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following
		FROM users u
		WHERE u.id = $1`,
		userID,
	).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email,
		&user.Bio, &user.Avatar, &user.Location, &user.Website,
		&user.CreatedAt, &user.UpdatedAt,
		&user.Followers, &user.Following,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user's profile
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := middleware.GetUserID(c)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	// Check if user is updating their own profile
	if id != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another user's profile"})
		return
	}
	
	var input model.UserUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Build the SQL query dynamically
	query := "UPDATE users SET updated_at = $1"
	args := []interface{}{time.Now()}
	argCount := 2
	
	if input.Name != nil {
		query += ", name = $" + strconv.Itoa(argCount)
		args = append(args, *input.Name)
		argCount++
	}
	
	if input.Bio != nil {
		query += ", bio = $" + strconv.Itoa(argCount)
		args = append(args, *input.Bio)
		argCount++
	}
	
	if input.Avatar != nil {
		query += ", avatar = $" + strconv.Itoa(argCount)
		args = append(args, *input.Avatar)
		argCount++
	}
	
	if input.Location != nil {
		query += ", location = $" + strconv.Itoa(argCount)
		args = append(args, *input.Location)
		argCount++
	}
	
	if input.Website != nil {
		query += ", website = $" + strconv.Itoa(argCount)
		args = append(args, *input.Website)
		argCount++
	}
	
	query += " WHERE id = $" + strconv.Itoa(argCount) + " RETURNING id"
	args = append(args, id)
	
	// Execute the update
	var updatedID string
	err = uc.db.QueryRow(query, args...).Scan(&updatedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	
	// Get the updated user
	var user model.User
	err = uc.db.QueryRow(
		`SELECT u.id, u.name, u.username, u.email, u.bio, u.avatar, u.location, u.website, 
		u.created_at, u.updated_at, 
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following
		FROM users u
		WHERE u.id = $1`,
		id,
	).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email,
		&user.Bio, &user.Avatar, &user.Location, &user.Website,
		&user.CreatedAt, &user.UpdatedAt,
		&user.Followers, &user.Following,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

// FollowUser creates a follow relationship between users
func (uc *UserController) FollowUser(c *gin.Context) {
	followingID := c.Param("id")
	followerID, err := middleware.GetUserID(c)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	// Check if trying to follow oneself
	if followerID == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}
	
	// Check if user to follow exists
	var exists bool
	err = uc.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", followingID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to follow not found"})
		return
	}
	
	// Check if already following
	err = uc.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)",
		followerID, followingID,
	).Scan(&exists)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Already following this user"})
		return
	}
	
	// Create follow relationship
	_, err = uc.db.Exec(
		"INSERT INTO follows (follower_id, following_id, created_at) VALUES ($1, $2, $3)",
		followerID, followingID, time.Now(),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully followed user"})
}

// UnfollowUser removes a follow relationship between users
func (uc *UserController) UnfollowUser(c *gin.Context) {
	followingID := c.Param("id")
	followerID, err := middleware.GetUserID(c)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	// Check if follow relationship exists
	var exists bool
	err = uc.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)",
		followerID, followingID,
	).Scan(&exists)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not following this user"})
		return
	}
	
	// Remove follow relationship
	_, err = uc.db.Exec(
		"DELETE FROM follows WHERE follower_id = $1 AND following_id = $2",
		followerID, followingID,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}

// GetFollowers returns users who follow the authenticated user
func (uc *UserController) GetFollowers(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	// Query followers from database
	rows, err := uc.db.Query(
		`SELECT u.id, u.name, u.username, u.email, u.bio, u.avatar, u.location, u.website, 
		u.created_at, u.updated_at, 
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following,
		(SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = u.id)) as is_followed
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}
	defer rows.Close()
	
	// Parse results
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Username, &user.Email,
			&user.Bio, &user.Avatar, &user.Location, &user.Website,
			&user.CreatedAt, &user.UpdatedAt,
			&user.Followers, &user.Following, &user.IsFollowed,
		)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
			return
		}
		
		users = append(users, user)
	}
	
	c.JSON(http.StatusOK, users)
}

// GetFollowing returns users that the authenticated user follows
func (uc *UserController) GetFollowing(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	// Query following from database
	rows, err := uc.db.Query(
		`SELECT u.id, u.name, u.username, u.email, u.bio, u.avatar, u.location, u.website, 
		u.created_at, u.updated_at, 
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following,
		true as is_followed
		FROM users u
		JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch following"})
		return
	}
	defer rows.Close()
	
	// Parse results
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Username, &user.Email,
			&user.Bio, &user.Avatar, &user.Location, &user.Website,
			&user.CreatedAt, &user.UpdatedAt,
			&user.Followers, &user.Following, &user.IsFollowed,
		)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
			return
		}
		
		users = append(users, user)
	}
	
	c.JSON(http.StatusOK, users)
}
