
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
	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	var filter model.NotificationFilter
	if !middleware.BindQuery(c, &filter) {
		return
	}

	// Query notifications from database
	notifications, err := nc.repo.Notification.FindByUserID(userID, filter)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Notifications retrieved successfully", notifications)
}

// GetUnreadCount returns the count of unread notifications
func (nc *NotificationController) GetUnreadCount(c *gin.Context) {
	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	// Get unread count
	count, err := nc.repo.Notification.CountUnread(userID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to count notifications")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Unread notification count retrieved", gin.H{"count": count})
}

// MarkAsRead marks a notification as read
func (nc *NotificationController) MarkAsRead(c *gin.Context) {
	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	notificationID, err := middleware.ParseUUIDParam(c, "id")
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	// Get notification to verify ownership
	notification, err := nc.repo.Notification.FindByID(notificationID)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Notification not found")
		return
	}

	// Verify notification belongs to user
	if !middleware.CheckResourceOwnership(c, notification.UserID, userID) {
		return
	}

	// Mark as read
	err = nc.repo.Notification.MarkAsRead(notificationID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Notification marked as read", nil)
}

// MarkAllAsRead marks all notifications for the current user as read
func (nc *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID, ok := middleware.RequireAuthentication(c)
	if !ok {
		return
	}

	// Mark all as read
	err := nc.repo.Notification.MarkAllAsRead(userID)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to mark notifications as read")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "All notifications marked as read", nil)
}
