
package model

import "time"

// Post represents a post in the system
type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Author    *User     `json:"author,omitempty"`
	Content   string    `json:"content"`
	Image     *string   `json:"image,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
	IsLiked   *bool     `json:"is_liked,omitempty"`
}

// PostCreate represents data needed to create a new post
type PostCreate struct {
	Content string  `json:"content" binding:"required"`
	Image   *string `json:"image,omitempty"`
}

// PostUpdate represents data that can be updated for a post
type PostUpdate struct {
	Content string  `json:"content" binding:"required"`
	Image   *string `json:"image,omitempty"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Author    *User     `json:"author,omitempty"`
	PostID    string    `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentCreate represents data needed to create a new comment
type CommentCreate struct {
	Content string `json:"content" binding:"required"`
}

// CommentUpdate represents data that can be updated for a comment
type CommentUpdate struct {
	Content string `json:"content" binding:"required"`
}
