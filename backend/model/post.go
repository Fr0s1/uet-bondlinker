
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents a post in the system
type Post struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Image     *string        `json:"image,omitempty" gorm:"size:255"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Calculated fields
	Likes    int   `json:"likes" gorm:"-"`
	Comments int   `json:"comments" gorm:"-"`
	IsLiked  *bool `json:"is_liked,omitempty" gorm:"-"`
	
	// Relations
	Author    *User     `json:"author,omitempty" gorm:"foreignKey:UserID"`
	LikesList []Like    `json:"-" gorm:"foreignKey:PostID"`
	CommentsList []Comment `json:"-" gorm:"foreignKey:PostID"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
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
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	PostID    uuid.UUID      `json:"post_id" gorm:"type:uuid;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relations
	Author *User `json:"author,omitempty" gorm:"foreignKey:UserID"`
	Post   *Post `json:"-" gorm:"foreignKey:PostID"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// CommentCreate represents data needed to create a new comment
type CommentCreate struct {
	Content string `json:"content" binding:"required"`
}

// CommentUpdate represents data that can be updated for a comment
type CommentUpdate struct {
	Content string `json:"content" binding:"required"`
}

// Like represents a like on a post
type Like struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `json:"post_id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// Relations
	User User `json:"-" gorm:"foreignKey:UserID"`
	Post Post `json:"-" gorm:"foreignKey:PostID"`
}
