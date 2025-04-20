package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message represents a chat message between users
type Message struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Content        string    `json:"content" gorm:"type:text;not null"`
	ConversationID uuid.UUID `json:"conversationId" gorm:"type:uuid;not null"`
	SenderID       uuid.UUID `json:"senderId" gorm:"type:uuid;not null"`
	RecipientID    uuid.UUID `json:"recipientId" gorm:"type:uuid;not null"`
	IsRead         bool      `json:"isRead" gorm:"default:false"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// Conversation represents a chat conversation between two users
type Conversation struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID1   uuid.UUID `json:"userId1" gorm:"type:uuid;not null"`
	UserID2   uuid.UUID `json:"userId2" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// ConversationResponse is the API response for a conversation
type ConversationResponse struct {
	ID          uuid.UUID     `json:"id"`
	Recipient   User          `json:"recipient"`
	LastMessage *MessageBrief `json:"lastMessage"`
}

// MessageBrief is a simplified message struct for conversation listings
type MessageBrief struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	IsRead    bool      `json:"isRead"`
}

// MessageFilter defines filters for querying messages
type MessageFilter struct {
	Pagination
	ConversationID uuid.UUID `form:"conversationId"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type WsMessageType string

const (
	WsMessageTypeMessage WsMessageType = "message"
	WsMessageTypeTyping  WsMessageType = "typing"
)

type WsMessage[T any] struct {
	ToUserId uuid.UUID     `json:"toUserId"`
	Type     WsMessageType `json:"type"`
	Payload  T             `json:"payload"`
}

func NewWsMessage[T any](toUserId uuid.UUID, messageType WsMessageType, payload T) WsMessage[T] {
	return WsMessage[T]{
		ToUserId: toUserId,
		Type:     messageType,
		Payload:  payload,
	}
}
