
package repository

import (
	"socialnet/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db}
}

// Create adds a new notification to the database
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

// FindByID finds a notification by ID
func (r *NotificationRepository) FindByID(id uuid.UUID) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// FindByUserID finds notifications for a user with filters
func (r *NotificationRepository) FindByUserID(userID uuid.UUID, filter model.NotificationFilter) ([]model.Notification, error) {
	var notifications []model.Notification
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")

	// Apply read filter if provided
	if filter.IsRead != nil {
		query = query.Where("is_read = ?", *filter.IsRead)
	}

	// Apply pagination
	query = query.Limit(filter.Limit).Offset(filter.Offset)

	// Include sender information
	query = query.Preload("Sender")

	err := query.Find(&notifications).Error
	return notifications, err
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllAsRead marks all notifications for a user as read
func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

// Delete removes a notification from the database
func (r *NotificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Notification{}, "id = ?", id).Error
}

// CountUnread counts unread notifications for a user
func (r *NotificationRepository) CountUnread(userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return int(count), err
}

// CreateFollowNotification creates a follow notification
func (r *NotificationRepository) CreateFollowNotification(followerID, followingID uuid.UUID) error {
	// Get follower details
	var follower model.User
	if err := r.db.First(&follower, "id = ?", followerID).Error; err != nil {
		return err
	}

	// Create notification
	notification := model.Notification{
		UserID:   followingID,
		SenderID: &followerID,
		Type:     model.NotificationTypeFollow,
		Message:  follower.Name + " started following you",
	}

	return r.Create(&notification)
}

// CreateLikeNotification creates a like notification
func (r *NotificationRepository) CreateLikeNotification(userID, postOwnerID uuid.UUID, postID uuid.UUID) error {
	// Don't notify yourself
	if userID == postOwnerID {
		return nil
	}

	// Get user details
	var user model.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	// Create notification
	entityType := "post"
	notification := model.Notification{
		UserID:          postOwnerID,
		SenderID:        &userID,
		Type:            model.NotificationTypeLike,
		Message:         user.Name + " liked your post",
		RelatedEntityID: &postID,
		EntityType:      &entityType,
	}

	return r.Create(&notification)
}

// CreateCommentNotification creates a comment notification
func (r *NotificationRepository) CreateCommentNotification(commenterID, postOwnerID uuid.UUID, postID, commentID uuid.UUID) error {
	// Don't notify yourself
	if commenterID == postOwnerID {
		return nil
	}

	// Get commenter details
	var commenter model.User
	if err := r.db.First(&commenter, "id = ?", commenterID).Error; err != nil {
		return err
	}

	// Create notification
	entityType := "comment"
	notification := model.Notification{
		UserID:          postOwnerID,
		SenderID:        &commenterID,
		Type:            model.NotificationTypeComment,
		Message:         commenter.Name + " commented on your post",
		RelatedEntityID: &commentID,
		EntityType:      &entityType,
	}

	return r.Create(&notification)
}

// CreateMessageNotification creates a message notification
func (r *NotificationRepository) CreateMessageNotification(senderID, recipientID uuid.UUID, conversationID uuid.UUID) error {
	// Get sender details
	var sender model.User
	if err := r.db.First(&sender, "id = ?", senderID).Error; err != nil {
		return err
	}

	// Create notification
	entityType := "conversation"
	notification := model.Notification{
		UserID:          recipientID,
		SenderID:        &senderID,
		Type:            model.NotificationTypeMessage,
		Message:         "New message from " + sender.Name,
		RelatedEntityID: &conversationID,
		EntityType:      &entityType,
	}

	return r.Create(&notification)
}
