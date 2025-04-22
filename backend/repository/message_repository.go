package repository

import (
	"cmp"
	"errors"
	"slices"
	"socialnet/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrConversationNotFound = gorm.ErrRecordNotFound

// MessageRepository handles database operations for messages
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db}
}

// UpsertConversation creates a new conversation between two users or returns existing one
func (r *MessageRepository) UpsertConversation(senderId, recipientId uuid.UUID) (*model.Conversation, error) {
	// Check if users exist
	var sender, recipient model.User
	if err := r.db.First(&sender, "id = ?", senderId).Error; err != nil {
		return nil, errors.New("sender not found")
	}
	if err := r.db.First(&recipient, "id = ?", recipientId).Error; err != nil {
		return nil, errors.New("recipient not found")
	}

	ids := []uuid.UUID{senderId, recipientId}
	slices.SortStableFunc(ids, func(a, b uuid.UUID) int {
		return cmp.Compare(a.String(), b.String())
	})

	// Check if conversation already exists
	var conversation model.Conversation
	err := r.db.Where(
		"user_id1 = ? AND user_id2 = ?",
		ids[0], ids[1],
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
		UserID1: ids[0],
		UserID2: ids[1],
	}

	if err := r.db.Create(&conversation).Error; err != nil {
		return nil, err
	}

	conversation.Recipient = recipient

	return &conversation, nil
}

// FindConversations gets all conversations for a user
func (r *MessageRepository) FindConversations(userID uuid.UUID) ([]model.ConversationResponse, error) {
	// Find all conversations where the user is either user1 or user2
	var conversations []model.Conversation
	if err := r.db.Preload("LastMessage").Where(
		"(user_id1 = ? OR user_id2 = ?) and last_message_id is not null", userID, userID,
	).Order("updated_at DESC").Find(&conversations).Error; err != nil {
		return nil, err
	}

	if len(conversations) == 0 {
		return nil, nil
	}

	userIds := make([]uuid.UUID, 0, len(conversations))
	lastMessageIds := make([]uuid.UUID, 0, len(conversations))
	for _, conv := range conversations {
		if conv.UserID1 == userID {
			userIds = append(userIds, conv.UserID2)
		} else {
			userIds = append(userIds, conv.UserID1)
		}
	}

	var userInfos []model.User
	if err := r.db.Model(&model.User{}).Where("id in ?", userIds).Find(&userInfos).Error; err != nil {
		return nil, err
	}

	var userMap = make(map[uuid.UUID]model.User)
	for _, user := range userInfos {
		userMap[user.ID] = user
	}

	var lastMessageMap = make(map[uuid.UUID]model.Message)
	if len(lastMessageIds) > 0 {
		var lastMessages []model.Message
		if err := r.db.Model(&model.Message{}).Where("id in ?", lastMessageIds).Find(&lastMessages).Error; err != nil {
			return nil, err
		}

		for _, message := range lastMessageMap {
			lastMessageMap[message.ID] = message
		}
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

		// Get last message in this conversation
		var lastMessageBrief *model.MessageBrief
		if conv.LastMessage != nil {
			lastMessageBrief = &model.MessageBrief{
				Content:   conv.LastMessage.Content,
				CreatedAt: conv.LastMessage.CreatedAt,
				IsRead:    conv.LastMessage.IsRead || conv.LastMessage.SenderID == userID, // Message is read if user is sender
			}
		}

		results = append(results, model.ConversationResponse{
			ID:          conv.ID,
			Recipient:   userMap[otherUserID],
			LastMessage: lastMessageBrief,
		})
	}

	return results, nil
}

// FindConversation gets a specific conversation
func (r *MessageRepository) FindConversation(id uuid.UUID, senderId uuid.UUID) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := r.db.Preload("LastMessage").First(&conversation, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if conversation.UserID1 != senderId && conversation.UserID2 != senderId {
		return nil, ErrConversationNotFound
	}

	var recipientId uuid.UUID
	if conversation.UserID1 == senderId {
		recipientId = conversation.UserID2
	} else {
		recipientId = conversation.UserID1
	}

	var recipient model.User
	if err := r.db.First(&recipient, "id = ?", recipientId).Error; err != nil {
		return nil, err
	}

	conversation.Recipient = recipient

	return &conversation, nil
}

// CreateMessage creates a new message
func (r *MessageRepository) CreateMessage(message *model.Message) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(message).Error
		if err != nil {
			return err
		}

		return tx.Model(&model.Conversation{}).Where("id = ?", message.ConversationID).Updates(map[string]any{
			"last_message_id": message.ID,
			"updated_at":      gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
	})
}

// FindMessages gets all messages in a conversation
func (r *MessageRepository) FindMessages(conversationID uuid.UUID, filter model.MessageFilter) ([]model.Message, error) {
	var messages []model.Message
	// Get messages between these users
	query := r.db.Where(
		"conversation_id = ?", conversationID,
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
	conv, err := r.FindConversation(conversationID, userID)
	if err != nil {
		return err
	}

	return r.db.Model(&model.Message{}).
		Where("conversation_id = ? AND sender_id = ? AND recipient_id = ? AND is_read = ?", conversationID, conv.Recipient.ID, userID, false).
		Update("is_read", true).Error
}
