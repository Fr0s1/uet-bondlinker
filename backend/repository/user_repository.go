
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
	if err != nil {
		return nil, err
	}
	
	// Count followers and following
	followers, _ := r.CountFollowers(id)
	following, _ := r.CountFollowing(id)
	user.Followers = followers
	user.Following = following
	
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	
	// Count followers and following
	followers, _ := r.CountFollowers(user.ID)
	following, _ := r.CountFollowing(user.ID)
	user.Followers = followers
	user.Following = following
	
	return &user, nil
}

// Update updates a user in the database
func (r *UserRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// FindAll finds all users matching the query with pagination
func (r *UserRepo) FindAll(query string, limit, offset int) ([]model.User, error) {
	var users []model.User
	db := r.db
	
	if query != "" {
		db = db.Where("name ILIKE ? OR username ILIKE ?", "%"+query+"%", "%"+query+"%")
	}
	
	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	
	// Count followers and following for each user
	for i := range users {
		followers, _ := r.CountFollowers(users[i].ID)
		following, _ := r.CountFollowing(users[i].ID)
		users[i].Followers = followers
		users[i].Following = following
	}
	
	return users, nil
}

// Follow creates a follow relationship between users
func (r *UserRepo) Follow(followerID, followingID uuid.UUID) error {
	follow := model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return r.db.Create(&follow).Error
}

// Unfollow removes a follow relationship between users
func (r *UserRepo) Unfollow(followerID, followingID uuid.UUID) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&model.Follow{}).Error
}

// IsFollowing checks if a user is following another user
func (r *UserRepo) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}

// GetFollowers gets all users who follow the specified user
func (r *UserRepo) GetFollowers(userID uuid.UUID, limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Limit(limit).Offset(offset).Order("follows.created_at DESC").
		Find(&users).Error
	
	if err != nil {
		return nil, err
	}
	
	// Count followers and following for each user
	for i := range users {
		followers, _ := r.CountFollowers(users[i].ID)
		following, _ := r.CountFollowing(users[i].ID)
		users[i].Followers = followers
		users[i].Following = following
		
		// Mark as followed
		isFollowed, _ := r.IsFollowing(userID, users[i].ID)
		users[i].IsFollowed = &isFollowed
	}
	
	return users, nil
}

// GetFollowing gets all users the specified user follows
func (r *UserRepo) GetFollowing(userID uuid.UUID, limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.Joins("JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", userID).
		Limit(limit).Offset(offset).Order("follows.created_at DESC").
		Find(&users).Error
	
	if err != nil {
		return nil, err
	}
	
	// Count followers and following for each user
	for i := range users {
		followers, _ := r.CountFollowers(users[i].ID)
		following, _ := r.CountFollowing(users[i].ID)
		users[i].Followers = followers
		users[i].Following = following
		
		// Mark as followed (all are followed)
		isFollowed := true
		users[i].IsFollowed = &isFollowed
	}
	
	return users, nil
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
