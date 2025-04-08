package main

import (
	"fmt"
	"log"
	"os"

	"socialnet/config"
	"socialnet/database"
	"socialnet/router"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Connect to the database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup router
	r := router.SetupRouter(db, cfg)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server running on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
