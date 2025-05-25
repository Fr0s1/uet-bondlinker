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
	Avatar         *string        `json:"avatar,omitempty" gorm:"size:1000"`
	Cover          *string        `json:"cover,omitempty" gorm:"size:1000"`
	Location       *string        `json:"location,omitempty" gorm:"size:100"`
	Website        *string        `json:"website,omitempty" gorm:"size:255"`
	EmailVerified  bool           `json:"emailVerified" gorm:"default:false"`
	FollowersCount int            `json:"followers" gorm:"default:0"`
	FollowingCount int            `json:"following" gorm:"default:0"`
	PostsCount     int            `json:"postsCount" gorm:"default:0"`
	CreatedAt      time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	IsFollowed     *bool          `json:"isFollowed,omitempty" gorm:"-"`

	// Relations
	Posts    []Post    `json:"-" gorm:"foreignKey:UserID"`
	Likes    []Like    `json:"-" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"-" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
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

type SaveFCMTokenReq struct {
	Token  string `json:"token" binding:"required"`
	Device string `json:"device" binding:"required"`
}

// UserUpdate represents data that can be updated for a user
type UserUpdate struct {
	Name     *string `json:"name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Cover    *string `json:"cover,omitempty"`
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
	FollowerID  uuid.UUID `json:"followerId" gorm:"type:uuid;primaryKey"`
	FollowingID uuid.UUID `json:"followingId" gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`

	// Relations
	Follower  User `json:"-" gorm:"foreignKey:FollowerID"`
	Following User `json:"-" gorm:"foreignKey:FollowingID"`
}

// TableName specifies the table name for Follow model
func (Follow) TableName() string {
	return "follows"
}

// FCMToken represents a Firebase Cloud Messaging token for a user
type FCMToken struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"userId" gorm:"type:uuid;not null"`
	Token     string         `json:"token" gorm:"not null"`
	Device    string         `json:"device"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for FCMToken model
func (FCMToken) TableName() string {
	return "fcm_tokens"
}

// BeforeCreate will set a UUID rather than numeric ID.
func (t *FCMToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
