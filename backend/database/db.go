
package database

import (
	"fmt"
	"socialnet/config"
	"socialnet/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connects to the database
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name, cfg.Database.SSLMode)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	// Enable the uuid-ossp extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Auto-migrate models
	return db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Post{},
		&model.Like{},
		&model.Comment{},
		&model.Share{},
		&model.Message{},
		&model.Conversation{},
	)
}
