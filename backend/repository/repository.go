package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"socialnet/model"
)

// UserRepository handles database operations related to users
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	FindAll(filter model.UserFilter) ([]model.User, error)
	Follow(followerID, followingID uuid.UUID) error
	Unfollow(followerID, followingID uuid.UUID) error
	IsFollowing(followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(userID uuid.UUID, filter model.FollowFilter) ([]model.User, error)
	GetFollowing(userID uuid.UUID, filter model.FollowFilter) ([]model.User, error)
	CountFollowers(userID uuid.UUID) (int, error)
	CountFollowing(userID uuid.UUID) (int, error)
	SearchUsers(query string, filter model.Pagination) ([]model.User, error)
	GetSuggestedUsers(userID uuid.UUID, filter model.Pagination) ([]model.User, error)
}

// PostRepository handles database operations related to posts
type PostRepository interface {
	Create(post *model.Post) error
	FindByID(id uuid.UUID) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id uuid.UUID) error
	FindAll(filter model.PostFilter) ([]model.Post, error)
	FindFeed(userID uuid.UUID, filter model.Pagination) ([]model.Post, error)
	FindTrending(filter model.Pagination) ([]model.Post, error)
	SearchPosts(query string, filter model.Pagination) ([]model.Post, error)
	Like(userID, postID uuid.UUID) (int, error)
	Unlike(userID, postID uuid.UUID) (int, error)
	IsLiked(userID, postID uuid.UUID) (bool, error)
	CountLikes(postID uuid.UUID) (int, error)
	Share(userID, postID uuid.UUID, content string) (*model.Post, error)
	GetSuggestedPosts(userID uuid.UUID, filter model.Pagination) ([]model.Post, error)
}

// CommentRepository handles database operations related to comments
type CommentRepository interface {
	Create(comment *model.Comment) error
	FindByID(id uuid.UUID) (*model.Comment, error)
	Update(comment *model.Comment) error
	Delete(id uuid.UUID) error
	FindByPostID(postID uuid.UUID, filter model.Pagination) ([]model.Comment, error)
}

// Repository holds all repositories
type Repository struct {
	User         UserRepository
	Post         PostRepository
	Comment      CommentRepository
	Message      *MessageRepository
	Notification *NotificationRepository
}

// NewRepository creates a new Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Post:         NewPostRepository(db),
		Comment:      NewCommentRepository(db),
		Message:      NewMessageRepository(db),
		Notification: NewNotificationRepository(db),
	}
}
