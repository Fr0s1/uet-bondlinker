
package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Password   string     `json:"-"` // Password is never sent to client
	Bio        *string    `json:"bio,omitempty"`
	Avatar     *string    `json:"avatar,omitempty"`
	Location   *string    `json:"location,omitempty"`
	Website    *string    `json:"website,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Followers  int        `json:"followers,omitempty"`
	Following  int        `json:"following,omitempty"`
	IsFollowed *bool      `json:"is_followed,omitempty"`
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
