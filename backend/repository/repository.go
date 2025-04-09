
package repository

import "gorm.io/gorm"

// Repository holds all repositories
type Repository struct {
	User    *UserRepository
	Post    *PostRepository
	Comment *CommentRepository
	Message *MessageRepository
}

// NewRepository creates a new Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Post:    NewPostRepository(db),
		Comment: NewCommentRepository(db),
		Message: NewMessageRepository(db),
	}
}
