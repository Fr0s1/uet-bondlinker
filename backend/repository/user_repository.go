
package repository

import (
	"socialnet/model"
	"socialnet/util"

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
	err := r.db.Where("id = ?", id).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (r *UserRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// FindAll finds all users with pagination and search
func (r *UserRepo) FindAll(filter model.UserFilter) ([]model.User, error) {
	var users []model.User
	query := r.db.Limit(filter.Limit).Offset(filter.Offset)

	// Add search filter if query is provided
	if filter.Query != "" {
		query = query.Where("name ILIKE ? OR username ILIKE ?", "%"+filter.Query+"%", "%"+filter.Query+"%")
	}

	err := query.Scan(&users).Error
	return users, err
}

// Follow adds a follow relationship between users
func (r *UserRepo) Follow(followerID, followingID uuid.UUID) error {
	// Use transaction to handle follow creation and counter updates
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	follow := model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	// Create follow relationship
	if err := tx.Create(&follow).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Increment follower's following count
	if err := tx.Model(&model.User{}).Where("id = ?", followerID).Update("following_count", gorm.Expr("following_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Increment followed user's followers count
	if err := tx.Model(&model.User{}).Where("id = ?", followingID).Update("followers_count", gorm.Expr("followers_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Unfollow removes a follow relationship between users
func (r *UserRepo) Unfollow(followerID, followingID uuid.UUID) error {
	// Use transaction to handle follow removal and counter updates
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete follow relationship
	result := tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&model.Follow{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// If relationship was found and deleted, update counters
	if result.RowsAffected > 0 {
		// Decrement follower's following count
		if err := tx.Model(&model.User{}).Where("id = ?", followerID).Update("following_count", gorm.Expr("following_count - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Decrement followed user's followers count
		if err := tx.Model(&model.User{}).Where("id = ?", followingID).Update("followers_count", gorm.Expr("followers_count - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// IsFollowing checks if a user is following another user
func (r *UserRepo) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}

// GetFollowers returns users who follow the specified user
func (r *UserRepo) GetFollowers(userID uuid.UUID, filter model.FollowFilter) ([]model.User, error) {
	var users []model.User
	err := r.db.Table("users").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Limit(filter.Limit).Offset(filter.Offset).
		Scan(&users).Error
	return users, err
}

// GetFollowing returns users that the specified user follows
func (r *UserRepo) GetFollowing(userID uuid.UUID, filter model.FollowFilter) ([]model.User, error) {
	var users []model.User
	err := r.db.Table("users").
		Joins("JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", userID).
		Limit(filter.Limit).Offset(filter.Offset).
		Scan(&users).Error
	return users, err
}

// CountFollowers returns the number of followers for a user
func (r *UserRepo) CountFollowers(userID uuid.UUID) (int, error) {
	var user model.User
	if err := r.db.Select("followers_count").Where("id = ?", userID).Scan(&user).Error; err != nil {
		return 0, err
	}
	return user.FollowersCount, nil
}

// CountFollowing returns the number of users that a user follows
func (r *UserRepo) CountFollowing(userID uuid.UUID) (int, error) {
	var user model.User
	if err := r.db.Select("following_count").Where("id = ?", userID).Scan(&user).Error; err != nil {
		return 0, err
	}
	return user.FollowingCount, nil
}

// SearchUsers searches for users by name or username
func (r *UserRepo) SearchUsers(query string, filter model.Pagination) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("name ILIKE ? OR username ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(filter.Limit).Offset(filter.Offset).
		Scan(&users).Error
	return users, err
}

// GetSuggestedUsers returns users that might interest the given user
func (r *UserRepo) GetSuggestedUsers(userID uuid.UUID, filter model.Pagination) ([]model.User, error) {
	var users []model.User

	// Find users followed by users that the current user follows (friends of friends)
	err := r.db.Distinct("users.*").
		Table("users").
		Joins("JOIN follows f1 ON users.id = f1.following_id").
		Joins("JOIN follows f2 ON f2.follower_id = f1.following_id AND f2.following_id != ?", userID).
		Where("f1.follower_id = ?", userID).
		Where("users.id != ?", userID).
		Where("NOT EXISTS (SELECT 1 FROM follows WHERE follower_id = ? AND following_id = users.id)", userID).
		Limit(filter.Limit).Offset(filter.Offset).
		Scan(&users).Error

	return users, err
}
