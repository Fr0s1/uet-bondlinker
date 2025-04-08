
package repository

import (
	"socialnet/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CommentRepo implements CommentRepository
type CommentRepo struct {
	db *gorm.DB
}

// NewCommentRepo creates a new CommentRepo
func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db}
}

// Create adds a new comment to the database
func (r *CommentRepo) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// FindByID finds a comment by ID
func (r *CommentRepo) FindByID(id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("Author").First(&comment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// Update updates a comment in the database
func (r *CommentRepo) Update(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

// Delete deletes a comment from the database
func (r *CommentRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Comment{}, "id = ?", id).Error
}

// FindByPostID finds all comments for a post with pagination
func (r *CommentRepo) FindByPostID(postID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("Author").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error
	
	if err != nil {
		return nil, err
	}
	
	return comments, nil
}
