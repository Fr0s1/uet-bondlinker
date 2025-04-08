package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Username       string         `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Email          string         `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Password       string         `json:"-" gorm:"column:password_hash;size:255;not null"`
	Bio            *string        `json:"bio,omitempty" gorm:"type:text"`
	Avatar         *string        `json:"avatar,omitempty" gorm:"size:255"`
	Location       *string        `json:"location,omitempty" gorm:"size:100"`
	Website        *string        `json:"website,omitempty" gorm:"size:255"`
	FollowersCount int            `json:"followers" gorm:"default:0"`
	FollowingCount int            `json:"following" gorm:"default:0"`
	PostsCount     int            `json:"posts_count" gorm:"default:0"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	IsFollowed     *bool          `json:"is_followed,omitempty" gorm:"-"`

	// Relations
	Posts    []Post    `json:"-" gorm:"foreignKey:UserID"`
	Likes    []Like    `json:"-" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"-" gorm:"foreignKey:UserID"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// UserRegistration represents data needed to register a new user
type UserRegistration struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLogin represents data needed to log in
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserUpdate represents data that can be updated for a user
type UserUpdate struct {
	Name     *string `json:"name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Location *string `json:"location,omitempty"`
	Website  *string `json:"website,omitempty"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Follow represents a follow relationship between users
type Follow struct {
	FollowerID  uuid.UUID `json:"follower_id" gorm:"type:uuid;primaryKey"`
	FollowingID uuid.UUID `json:"following_id" gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relations
	Follower  User `json:"-" gorm:"foreignKey:FollowerID"`
	Following User `json:"-" gorm:"foreignKey:FollowingID"`
}
