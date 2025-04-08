
package database

import (
	"fmt"
	"log"
	"socialnet/config"
	"socialnet/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB establishes a connection to the database using GORM
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigrations runs all database migrations using GORM
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	// AutoMigrate will create tables, foreign keys, constraints, etc.
	err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
		&model.Follow{},
		&model.Like{},
	)
	
	if err != nil {
		return err
	}
	
	// Check if counters need to be initialized
	if err := updateCountersIfNeeded(db); err != nil {
		log.Printf("Warning: Failed to update counters: %v", err)
	}
	
	log.Println("Migrations completed successfully.")
	return nil
}

// updateCountersIfNeeded initializes counter fields if they're new
func updateCountersIfNeeded(db *gorm.DB) error {
	// Update followers/following counts for each user
	if err := db.Exec(`
		UPDATE users u
		SET followers_count = (
			SELECT COUNT(*) FROM follows f WHERE f.following_id = u.id
		)
		WHERE EXISTS (
			SELECT 1 FROM follows f WHERE f.following_id = u.id
		) AND (followers_count IS NULL OR followers_count = 0)
	`).Error; err != nil {
		return err
	}
	
	// Update following count
	if err := db.Exec(`
		UPDATE users u
		SET following_count = (
			SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id
		)
		WHERE EXISTS (
			SELECT 1 FROM follows f WHERE f.follower_id = u.id
		) AND (following_count IS NULL OR following_count = 0)
	`).Error; err != nil {
		return err
	}
	
	// Update post counts
	if err := db.Exec(`
		UPDATE users u
		SET posts_count = (
			SELECT COUNT(*) FROM posts p WHERE p.user_id = u.id
		)
		WHERE EXISTS (
			SELECT 1 FROM posts p WHERE p.user_id = u.id
		) AND (posts_count IS NULL OR posts_count = 0)
	`).Error; err != nil {
		return err
	}
	
	// Update likes count for posts
	if err := db.Exec(`
		UPDATE posts p
		SET likes_count = (
			SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id
		)
		WHERE EXISTS (
			SELECT 1 FROM likes l WHERE l.post_id = p.id
		) AND (likes_count IS NULL OR likes_count = 0)
	`).Error; err != nil {
		return err
	}
	
	// Update comments count for posts
	if err := db.Exec(`
		UPDATE posts p
		SET comments_count = (
			SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id
		)
		WHERE EXISTS (
			SELECT 1 FROM comments c WHERE c.post_id = p.id
		) AND (comments_count IS NULL OR comments_count = 0)
	`).Error; err != nil {
		return err
	}
	
	return nil
}
