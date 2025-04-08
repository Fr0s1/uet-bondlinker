
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

// FindByID finds a post by ID
func (r *PostRepo) FindByID(id uuid.UUID) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("Author").First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	
	// Count likes and comments
	likes, _ := r.CountLikes(id)
	comments, _ := r.CountComments(id)
	post.Likes = likes
	post.Comments = comments
	
	return &post, nil
}

// Update updates a post in the database
func (r *PostRepo) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

// Delete deletes a post from the database
func (r *PostRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Post{}, "id = ?", id).Error
}

// FindAll finds all posts with pagination
func (r *PostRepo) FindAll(userID *uuid.UUID, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	query := r.db.Preload("Author").Order("created_at DESC").Limit(limit).Offset(offset)
	
	// Filter by user if userID is provided
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	
	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}
	
	// Populate likes, comments, and isLiked for each post
	for i := range posts {
		likes, _ := r.CountLikes(posts[i].ID)
		comments, _ := r.CountComments(posts[i].ID)
		posts[i].Likes = likes
		posts[i].Comments = comments
		
		// Check if the post is liked by the current user (if authenticated)
		if userID != nil {
			isLiked, _ := r.IsLiked(*userID, posts[i].ID)
			posts[i].IsLiked = &isLiked
		}
	}
	
	return posts, nil
}

// FindFeed finds posts for a user's feed (posts from followed users and own posts)
func (r *PostRepo) FindFeed(userID uuid.UUID, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	
	// Get posts from followed users and own posts
	err := r.db.Preload("Author").
		Joins("LEFT JOIN follows ON posts.user_id = follows.following_id").
		Where("follows.follower_id = ? OR posts.user_id = ?", userID, userID).
		Order("posts.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error
	
	if err != nil {
		return nil, err
	}
	
	// Populate likes, comments, and isLiked for each post
	for i := range posts {
		likes, _ := r.CountLikes(posts[i].ID)
		comments, _ := r.CountComments(posts[i].ID)
		posts[i].Likes = likes
		posts[i].Comments = comments
		
		isLiked, _ := r.IsLiked(userID, posts[i].ID)
		posts[i].IsLiked = &isLiked
	}
	
	return posts, nil
}

// Like adds a like to a post
func (r *PostRepo) Like(userID, postID uuid.UUID) error {
	like := model.Like{
		UserID: userID,
		PostID: postID,
	}
	return r.db.Create(&like).Error
}

// Unlike removes a like from a post
func (r *PostRepo) Unlike(userID, postID uuid.UUID) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Like{}).Error
}

// IsLiked checks if a post is liked by a user
func (r *PostRepo) IsLiked(userID, postID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

// CountLikes counts the number of likes for a post
func (r *PostRepo) CountLikes(postID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&model.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return int(count), err
}

// CountComments counts the number of comments for a post
func (r *PostRepo) CountComments(postID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return int(count), err
}
