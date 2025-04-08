
package repository

import (
	"socialnet/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepo implements UserRepository
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

// Create adds a new user to the database
func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepo) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

// FindByEmail finds a user by email
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindByUsername finds a user by username
func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// Update updates a user in the database
func (r *UserRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// FindAll finds all users matching the query with pagination
func (r *UserRepo) FindAll(filter model.UserFilter) ([]model.User, error) {
	var users []model.User
	db := r.db
	
	if filter.Query != "" {
		db = db.Where("name ILIKE ? OR username ILIKE ?", "%"+filter.Query+"%", "%"+filter.Query+"%")
	}
	
	err := db.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&users).Error
	return users, err
}

// Follow creates a follow relationship between users and increments counters
func (r *UserRepo) Follow(followerID, followingID uuid.UUID) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create follow relationship
	follow := model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	if err := tx.Create(&follow).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Increment follower's following_count
	if err := tx.Model(&model.User{}).Where("id = ?", followerID).Update("following_count", gorm.Expr("following_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Increment followed user's followers_count
	if err := tx.Model(&model.User{}).Where("id = ?", followingID).Update("followers_count", gorm.Expr("followers_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Unfollow removes a follow relationship between users and decrements counters
func (r *UserRepo) Unfollow(followerID, followingID uuid.UUID) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Remove follow relationship
	if err := tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&model.Follow{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Decrement follower's following_count
	if err := tx.Model(&model.User{}).Where("id = ?", followerID).Update("following_count", gorm.Expr("following_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Decrement followed user's followers_count
	if err := tx.Model(&model.User{}).Where("id = ?", followingID).Update("followers_count", gorm.Expr("followers_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// IsFollowing checks if a user is following another user
func (r *UserRepo) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}

// GetFollowers gets all users who follow the specified user
func (r *UserRepo) GetFollowers(userID uuid.UUID, filter model.FollowFilter) ([]model.User, error) {
	var users []model.User
	// Use join to avoid N+1 query
	err := r.db.Select("users.*").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Limit(filter.Limit).Offset(filter.Offset).
		Order("follows.created_at DESC").
		Find(&users).Error
	
	return users, err
}

// GetFollowing gets all users the specified user follows
func (r *UserRepo) GetFollowing(userID uuid.UUID, filter model.FollowFilter) ([]model.User, error) {
	var users []model.User
	// Use join to avoid N+1 query
	err := r.db.Select("users.*").
		Joins("JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", userID).
		Limit(filter.Limit).Offset(filter.Offset).
		Order("follows.created_at DESC").
		Find(&users).Error
	
	return users, err
}

// CountFollowers counts the number of followers for a user
func (r *UserRepo) CountFollowers(userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return int(count), err
}

// CountFollowing counts the number of users the specified user follows
func (r *UserRepo) CountFollowing(userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return int(count), err
}
