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
	// Use transaction to handle comment creation and counter update
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create comment
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Increment post's comments_count
	if err := tx.Model(&model.Post{}).Where("id = ?", comment.PostID).Update("comments_count", gorm.Expr("comments_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// FindByID finds a comment by ID
func (r *CommentRepo) FindByID(id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment
	// Preload author to avoid N+1 query
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
	// Use transaction to handle comment deletion and counter update
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get comment to get post ID for counter update
	var comment model.Comment
	if err := tx.First(&comment, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete comment
	if err := tx.Delete(&model.Comment{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Decrement post's comments_count
	if err := tx.Model(&model.Post{}).Where("id = ?", comment.PostID).Update("comments_count", gorm.Expr("comments_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// FindByPostID finds all comments for a post with pagination
func (r *CommentRepo) FindByPostID(postID uuid.UUID, filter model.Pagination) ([]model.Comment, error) {
	var comments []model.Comment
	// Preload author to avoid N+1 query
	err := r.db.Preload("Author").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Limit(filter.Limit).Offset(filter.Offset).
		Find(&comments).Error

	return comments, err
}
