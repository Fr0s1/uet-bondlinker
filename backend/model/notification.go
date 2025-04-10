package model

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeFollow      NotificationType = "follow"
	NotificationTypeLike        NotificationType = "like"
	NotificationTypeComment     NotificationType = "comment"
	NotificationTypeShare       NotificationType = "share"
	NotificationTypeMessage     NotificationType = "message"
	NotificationTypeSystemAlert NotificationType = "system_alert"
)

// Notification represents a notification in the system
type Notification struct {
	ID              uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"`
	SenderID        *uuid.UUID       `json:"sender_id,omitempty" gorm:"type:uuid"`
	Type            NotificationType `json:"type" gorm:"not null"`
	Message         string           `json:"message" gorm:"not null"`
	RelatedEntityID *uuid.UUID       `json:"related_entity_id,omitempty" gorm:"type:uuid"`
	EntityType      *string          `json:"entity_type,omitempty"`
	IsRead          bool             `json:"is_read" gorm:"default:false"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	User   User `json:"-" gorm:"foreignKey:UserID"`
	Sender User `json:"-" gorm:"foreignKey:SenderID"`
}

// TableName specifies the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}

// NotificationFilter represents filter options for notifications
type NotificationFilter struct {
	Pagination
	IsRead *bool `form:"is_read"`
}
