
package repository

import (
	"errors"
	"socialnet/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MessageRepository handles database operations for messages
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db}
}

// CreateConversation creates a new conversation between two users or returns existing one
func (r *MessageRepository) CreateConversation(userID1, userID2 uuid.UUID) (*model.Conversation, error) {
	// Check if users exist
	var user1, user2 model.User
	if err := r.db.First(&user1, "id = ?", userID1).Error; err != nil {
		return nil, errors.New("sender not found")
	}
	if err := r.db.First(&user2, "id = ?", userID2).Error; err != nil {
		return nil, errors.New("recipient not found")
	}

	// Check if conversation already exists
	var conversation model.Conversation
	err := r.db.Where(
		"(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)",
		userID1, userID2, userID2, userID1,
	).First(&conversation).Error

	if err == nil {
		// Conversation exists, return it
		return &conversation, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new conversation
	conversation = model.Conversation{
		UserID1: userID1,
		UserID2: userID2,
	}

	if err := r.db.Create(&conversation).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

// FindConversations gets all conversations for a user
func (r *MessageRepository) FindConversations(userID uuid.UUID) ([]model.ConversationResponse, error) {
	// Find all conversations where the user is either user1 or user2
	var conversations []model.Conversation
	if err := r.db.Where(
		"user_id1 = ? OR user_id2 = ?", userID, userID,
	).Order("updated_at DESC").Find(&conversations).Error; err != nil {
		return nil, err
	}

	var results []model.ConversationResponse
	for _, conv := range conversations {
		// Determine which user is the other person in the conversation
		var otherUserID uuid.UUID
		if conv.UserID1 == userID {
			otherUserID = conv.UserID2
		} else {
			otherUserID = conv.UserID1
		}

		// Get other user's info
		var otherUser model.User
		if err := r.db.Model(&model.User{}).Where("id = ?", otherUserID).First(&otherUser).Error; err != nil {
			continue // Skip this conversation if user not found
		}

		// Get last message in this conversation
		var lastMessage model.Message
		lastMessageErr := r.db.Where(
			"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
			conv.UserID1, conv.UserID2, conv.UserID2, conv.UserID1,
		).Order("created_at DESC").First(&lastMessage).Error

		var lastMessageBrief *model.MessageBrief
		if lastMessageErr == nil {
			lastMessageBrief = &model.MessageBrief{
				Content:   lastMessage.Content,
				CreatedAt: lastMessage.CreatedAt,
				IsRead:    lastMessage.IsRead || lastMessage.SenderID == userID, // Message is read if user is sender
			}
		}

		results = append(results, model.ConversationResponse{
			ID: conv.ID,
			Recipient: otherUser,
			LastMessage: lastMessageBrief,
		})
	}

	return results, nil
}

// FindConversation gets a specific conversation
func (r *MessageRepository) FindConversation(id uuid.UUID) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := r.db.First(&conversation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

// FindConversationWithUser finds or creates a conversation with another user
func (r *MessageRepository) FindConversationWithUser(userID, otherUserID uuid.UUID) (*model.Conversation, error) {
	return r.CreateConversation(userID, otherUserID)
}

// GetConversationUsers gets the users in a conversation
func (r *MessageRepository) GetConversationUsers(conversationID uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	var conversation model.Conversation
	if err := r.db.First(&conversation, "id = ?", conversationID).Error; err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	return conversation.UserID1, conversation.UserID2, nil
}

// CreateMessage creates a new message
func (r *MessageRepository) CreateMessage(message *model.Message) error {
	return r.db.Create(message).Error
}

// FindMessages gets all messages in a conversation
func (r *MessageRepository) FindMessages(conversationID uuid.UUID, filter model.MessageFilter) ([]model.Message, error) {
	var messages []model.Message
	
	// Get the users in the conversation
	user1, user2, err := r.GetConversationUsers(conversationID)
	if err != nil {
		return nil, err
	}

	// Get messages between these users
	query := r.db.Where(
		"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		user1, user2, user2, user1,
	).Order("created_at DESC")

	// Apply pagination if specified
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(50) // Default limit
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	// Reverse the messages order to be chronological
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// MarkMessagesAsRead marks all messages in a conversation as read for a user
func (r *MessageRepository) MarkMessagesAsRead(conversationID, userID uuid.UUID) error {
	// Get the users in the conversation
	user1, user2, err := r.GetConversationUsers(conversationID)
	if err != nil {
		return err
	}

	// Verify user is part of conversation
	if userID != user1 && userID != user2 {
		return errors.New("user not part of this conversation")
	}

	// Mark all messages from the other user as read
	var otherUserID uuid.UUID
	if userID == user1 {
		otherUserID = user2
	} else {
		otherUserID = user1
	}

	return r.db.Model(&model.Message{}).
		Where("sender_id = ? AND recipient_id = ? AND is_read = ?", otherUserID, userID, false).
		Update("is_read", true).Error
}
