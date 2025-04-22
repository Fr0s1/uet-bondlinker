package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"
	"socialnet/util"
	"socialnet/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MessageController handles message-related requests
type MessageController struct {
	repo  *repository.Repository
	wsHub *websocket.Hub
	cfg   *config.Config
}

// NewMessageController creates a new MessageController
func NewMessageController(repo *repository.Repository, wsHub *websocket.Hub, cfg *config.Config) *MessageController {
	return &MessageController{
		repo:  repo,
		wsHub: wsHub,
		cfg:   cfg,
	}
}

// GetConversations returns all conversations for the current user
func (mc *MessageController) GetConversations(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	conversations, err := mc.repo.Message.FindConversations(userID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch conversations")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", conversations)
}

// GetConversation returns a specific conversation
func (mc *MessageController) GetConversation(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Get conversation to check if user is part of it
	conversation, err := mc.repo.Message.FindConversation(convID, userID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Conversation not found")
		return
	}

	var lastMessageBrief *model.MessageBrief
	if conversation.LastMessage != nil {
		lastMessageBrief = &model.MessageBrief{
			Content:   conversation.LastMessage.Content,
			CreatedAt: conversation.LastMessage.CreatedAt,
			IsRead:    conversation.LastMessage.IsRead || conversation.LastMessage.SenderID == userID, // Message is read if user is sender
		}
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", model.ConversationResponse{
		ID:          conversation.ID,
		Recipient:   conversation.Recipient,
		LastMessage: lastMessageBrief,
	})
}

// CreateConversation creates a new conversation with another user
func (mc *MessageController) CreateConversation(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get recipient ID from request body
	var input struct {
		RecipientID string `json:"recipientId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	recipientID, err := uuid.Parse(input.RecipientID)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid recipient ID")
		return
	}

	// Can't create conversation with yourself
	if userID == recipientID {
		util.RespondWithError(c, http.StatusBadRequest, "Cannot create conversation with yourself")
		return
	}

	// Create or get existing conversation
	conversation, err := mc.repo.Message.UpsertConversation(userID, recipientID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to create conversation")
		return
	}

	util.RespondWithSuccess(c, http.StatusCreated, "success", model.ConversationResponse{
		ID:        conversation.ID,
		Recipient: conversation.Recipient,
	})
}

// GetMessages returns all messages in a conversation
func (mc *MessageController) GetMessages(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Get conversation to check if user is part of it
	_, err = mc.repo.Message.FindConversation(convID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			util.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch conversation", err)
		return
	}

	// Parse filter parameters
	var filter model.MessageFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	filter.ConversationID = convID

	// Get messages
	messages, err := mc.repo.Message.FindMessages(convID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch messages", err)
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "success", messages)
}

// CreateMessage sends a new message in a conversation
func (mc *MessageController) CreateMessage(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Get conversation to check if user is part of it
	conv, err := mc.repo.Message.FindConversation(convID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			util.RespondWithError(c, http.StatusNotFound, "Conversation not found")
			return
		}
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch conversation", err)
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
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create message
	message := model.Message{
		ConversationID: convID,
		SenderID:       userID,
		RecipientID:    recipientID,
		Content:        input.Content,
		IsRead:         false,
	}

	if err := mc.repo.Message.CreateMessage(&message); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to send message")
		return
	}

	// After successfully creating the message, send FCM notification
	tokens, err := mc.repo.User.GetUserFCMTokens(recipientID)
	if err != nil {
		log.Printf("Error getting FCM tokens: %v", err)
	} else if len(tokens) > 0 {
		// Send notification to all user's devices
		notification := model.WsMessage{
			Type: model.WsMessageTypeMessage,
			Data: message,
		}
		notificationBytes, _ := json.Marshal(notification)
		mc.wsHub.SendToUser(recipientID, notificationBytes)
	}

	util.RespondWithSuccess(c, http.StatusCreated, "success", message)
}

// MarkConversationAsRead marks all messages in a conversation as read
func (mc *MessageController) MarkConversationAsRead(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	convIDStr := c.Param("id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Mark messages as read
	if err := mc.repo.Message.MarkMessagesAsRead(convID, userID); err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			util.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to mark messages as read", err)
		return
	}

	c.Status(http.StatusNoContent)
}
