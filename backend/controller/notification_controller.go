
package controller

import (
	"net/http"
	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"
	"socialnet/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationController handles notification-related requests
type NotificationController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewNotificationController creates a new NotificationController
func NewNotificationController(repo *repository.Repository, cfg *config.Config) *NotificationController {
	return &NotificationController{
		repo: repo,
		cfg:  cfg,
	}
}

// GetNotifications returns a list of notifications for the current user
func (nc *NotificationController) GetNotifications(c *gin.Context) {
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

	var filter model.NotificationFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query notifications from database
	notifications, err := nc.repo.Notification.FindByUserID(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUnreadCount returns the count of unread notifications
func (nc *NotificationController) GetUnreadCount(c *gin.Context) {
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

	// Get unread count
	count, err := nc.repo.Notification.CountUnread(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// MarkAsRead marks a notification as read
func (nc *NotificationController) MarkAsRead(c *gin.Context) {
	userIDStr, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Get notification to verify ownership
	notification, err := nc.repo.Notification.FindByID(notificationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	// Verify notification belongs to user
	userID, _ := uuid.Parse(userIDStr)
	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot access this notification"})
		return
	}

	// Mark as read
	err = nc.repo.Notification.MarkAsRead(notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead marks all notifications for the current user as read
func (nc *NotificationController) MarkAllAsRead(c *gin.Context) {
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

	// Mark all as read
	err = nc.repo.Notification.MarkAllAsRead(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}
