
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
	
	log.Println("Migrations completed successfully.")
	return nil
}
