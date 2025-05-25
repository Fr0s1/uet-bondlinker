package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents a post in the system
type Post struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID      `json:"userId" gorm:"type:uuid;not null"`
	Content       string         `json:"content" gorm:"type:text;not null"`
	Image         *string        `json:"image,omitempty"`
	LikesCount    int            `json:"likes" gorm:"default:0"`
	CommentsCount int            `json:"comments" gorm:"default:0"`
	SharesCount   int            `json:"shares" gorm:"default:0"`
	SharedPostID  *uuid.UUID     `json:"sharedPostId,omitempty" gorm:"type:uuid"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	IsLiked       *bool          `json:"isLiked,omitempty" gorm:"-"`

	// Relations
	Author       *User     `json:"author,omitempty" gorm:"foreignKey:UserID"`
	LikesList    []Like    `json:"-" gorm:"foreignKey:PostID"`
	CommentsList []Comment `json:"-" gorm:"foreignKey:PostID"`
	SharedPost   *Post     `json:"sharedPost,omitempty" gorm:"foreignKey:SharedPostID"`
}

// TableName specifies the table name for Post model
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate will set a UUID rather than numeric ID.
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Increment user's post count on creation
	err := tx.Model(&User{}).Where("id = ?", p.UserID).Update("posts_count", gorm.Expr("posts_count + 1")).Error
	if err != nil {
		return err
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

// PostShare represents data needed to share a post
type PostShare struct {
	Content string `json:"content"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"userId" gorm:"type:uuid;not null"`
	PostID    uuid.UUID      `json:"postId" gorm:"type:uuid;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Author *User `json:"author,omitempty" gorm:"foreignKey:UserID"`
	Post   *Post `json:"-" gorm:"foreignKey:PostID"`
}

// TableName specifies the table name for Comment model
func (Comment) TableName() string {
	return "comments"
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
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `json:"postId" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`

	// Relations
	User User `json:"-" gorm:"foreignKey:UserID"`
	Post Post `json:"-" gorm:"foreignKey:PostID"`
}

// TableName specifies the table name for Like model
func (Like) TableName() string {
	return "likes"
}
