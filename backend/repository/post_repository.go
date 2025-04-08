
package repository

import (
	"socialnet/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostRepo implements PostRepository
type PostRepo struct {
	db *gorm.DB
}

// NewPostRepo creates a new PostRepo
func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{db}
}

// Create adds a new post to the database
func (r *PostRepo) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// FindByID finds a post by ID with author preloaded
func (r *PostRepo) FindByID(id uuid.UUID) (*model.Post, error) {
	var post model.Post
	// Preload author to avoid N+1 query
	err := r.db.Preload("Author").First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	
	return &post, nil
}

// Update updates a post in the database
func (r *PostRepo) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

// Delete deletes a post from the database
func (r *PostRepo) Delete(id uuid.UUID) error {
	// Use transaction to handle deletion and counter updates
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get post to get user ID for counter update
	var post model.Post
	if err := tx.First(&post, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete post
	if err := tx.Delete(&model.Post{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Decrement user's post count
	if err := tx.Model(&model.User{}).Where("id = ?", post.UserID).Update("posts_count", gorm.Expr("posts_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// FindAll finds all posts with pagination and author preloaded
func (r *PostRepo) FindAll(filter model.PostFilter) ([]model.Post, error) {
	var posts []model.Post
	query := r.db.Preload("Author").Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset)
	
	// Filter by user if userID is provided
	if filter.UserID != "" {
		userID, err := uuid.Parse(filter.UserID)
		if err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}
	
	err := query.Find(&posts).Error
	return posts, err
}

// FindFeed finds posts for a user's feed (posts from followed users and own posts)
func (r *PostRepo) FindFeed(userID uuid.UUID, filter model.Pagination) ([]model.Post, error) {
	var posts []model.Post
	
	// Get posts from followed users and own posts using a single join
	err := r.db.Preload("Author").
		Distinct("posts.*").
		Select("posts.*").
		Table("posts").
		Joins("LEFT JOIN follows ON posts.user_id = follows.following_id AND follows.follower_id = ?", userID).
		Where("follows.follower_id = ? OR posts.user_id = ?", userID, userID).
		Order("posts.created_at DESC").
		Limit(filter.Limit).Offset(filter.Offset).
		Find(&posts).Error
	
	return posts, err
}

// Like adds a like to a post
func (r *PostRepo) Like(userID, postID uuid.UUID) error {
	// Use transaction to handle like creation and counter update
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	like := model.Like{
		UserID: userID,
		PostID: postID,
	}
	
	// Create like
	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Increment post's likes_count
	if err := tx.Model(&model.Post{}).Where("id = ?", postID).Update("likes_count", gorm.Expr("likes_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

// Unlike removes a like from a post
func (r *PostRepo) Unlike(userID, postID uuid.UUID) error {
	// Use transaction to handle like removal and counter update
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	
	// Delete like
	result := tx.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Like{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	
	// If like was found and deleted, decrement counter
	if result.RowsAffected > 0 {
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).Update("likes_count", gorm.Expr("likes_count - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit().Error
}

// IsLiked checks if a post is liked by a user
func (r *PostRepo) IsLiked(userID, postID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

// CountLikes returns the number of likes for a post
func (r *PostRepo) CountLikes(postID uuid.UUID) (int, error) {
	var post model.Post
	if err := r.db.Select("likes_count").First(&post, "id = ?", postID).Error; err != nil {
		return 0, err
	}
	return post.LikesCount, nil
}
