
package database

import (
	"database/sql"
	"fmt"

	"socialnet/config"

	_ "github.com/lib/pq"
)

// ConnectDB establishes a connection to the database
func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Create users table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			bio TEXT,
			avatar VARCHAR(255),
			location VARCHAR(100),
			website VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	// Create posts table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			image VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	// Create follows table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS follows (
			follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (follower_id, following_id)
		)
	`); err != nil {
		return err
	}

	// Create likes table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, post_id)
		)
	`); err != nil {
		return err
	}

	// Create comments table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	return nil
}
