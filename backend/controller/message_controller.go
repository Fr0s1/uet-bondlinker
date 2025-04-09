
package controller

import (
	"net/http"
	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MessageController handles message-related requests
type MessageController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewMessageController creates a new MessageController
func NewMessageController(repo *repository.Repository, cfg *config.Config) *MessageController {
	return &MessageController{
		repo: repo,
		cfg:  cfg,
	}
}

// GetConversations returns all conversations for the current user
func (mc *MessageController) GetConversations(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	conversations, err := mc.repo.Message.FindConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

// GetConversation returns a specific conversation
func (mc *MessageController) GetConversation(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Get conversation to check if user is part of it
	conv, err := mc.repo.Message.FindConversation(convID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Verify user is part of conversation
	if conv.UserID1 != userID && conv.UserID2 != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this conversation"})
		return
	}

	// Determine which user is the other person in the conversation
	var otherUserID uuid.UUID
	if conv.UserID1 == userID {
		otherUserID = conv.UserID2
	} else {
		otherUserID = conv.UserID1
	}

	// Get other user's info
	otherUser, err := mc.repo.User.FindByID(otherUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": conv.ID,
		"recipient": otherUser,
	})
}

// CreateConversation creates a new conversation with another user
func (mc *MessageController) CreateConversation(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get recipient ID from request body
	var input struct {
		RecipientID string `json:"recipientId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipientID, err := uuid.Parse(input.RecipientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient ID"})
		return
	}

	// Can't create conversation with yourself
	if userID == recipientID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create conversation with yourself"})
		return
	}

	// Create or get existing conversation
	conversation, err := mc.repo.Message.CreateConversation(userID, recipientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	// Get recipient details
	recipient, err := mc.repo.User.FindByID(recipientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipient details"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": conversation.ID,
		"recipient": recipient,
		"lastMessage": nil,
	})
}

// GetMessages returns all messages in a conversation
func (mc *MessageController) GetMessages(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Get conversation to check if user is part of it
	conv, err := mc.repo.Message.FindConversation(convID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Verify user is part of conversation
	if conv.UserID1 != userID && conv.UserID2 != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view these messages"})
		return
	}

	// Parse filter parameters
	var filter model.MessageFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filter.ConversationID = convID

	// Get messages
	messages, err := mc.repo.Message.FindMessages(convID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// CreateMessage sends a new message in a conversation
func (mc *MessageController) CreateMessage(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Get conversation to check if user is part of it
	conv, err := mc.repo.Message.FindConversation(convID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Verify user is part of conversation
	if conv.UserID1 != userID && conv.UserID2 != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to send messages in this conversation"})
		return
	}

	// Determine recipient ID
	var recipientID uuid.UUID
	if conv.UserID1 == userID {
		recipientID = conv.UserID2
	} else {
		recipientID = conv.UserID1
	}

	// Parse message content
	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create message
	message := model.Message{
		SenderID:    userID,
		RecipientID: recipientID,
		Content:     input.Content,
		IsRead:      false,
	}

	if err := mc.repo.Message.CreateMessage(&message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// MarkConversationAsRead marks all messages in a conversation as read
func (mc *MessageController) MarkConversationAsRead(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Mark messages as read
	if err := mc.repo.Message.MarkMessagesAsRead(convID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark messages as read"})
		return
	}

	c.Status(http.StatusNoContent)
}
