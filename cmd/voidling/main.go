package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	_ "time/tzdata" // Embed timezone database for Windows

	"github.com/kaffeed/voidling/config"
	"github.com/kaffeed/voidling/internal/bot"
	"github.com/kaffeed/voidling/internal/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

func run() error {
	// Configure log to show file and line number

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.LogLevel == "debug" {
		log.Printf("Set debug log level..")
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	log.Printf("Starting voidling...")
	log.Printf("Database: %s", cfg.DatabasePath)
	log.Printf("Log level: %s", cfg.LogLevel)

	// Ensure database directory exists
	dbDir := filepath.Dir(cfg.DatabasePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Run migrations
	log.Println("Running database migrations...")
	if err := runMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create database queries
	queries := database.New(db)

	// Create bot instance
	log.Println("Creating bot instance...")
	b, err := bot.New(cfg, queries, db)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	// Start bot
	log.Println("Starting bot...")
	if err := b.Start(); err != nil {
		return fmt.Errorf("failed to start bot: %w", err)
	}

	log.Println("Bot started successfully!")
	log.Println("Using Wise Old Man API for competition tracking")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")

	// Stop bot
	if err := b.Stop(); err != nil {
		log.Printf("Error stopping bot: %v", err)
	}

	log.Println("Shutdown complete")
	return nil
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	// Get migrations directory
	// In production, migrations should be embedded or in a known location
	migrationsDir := "migrations"

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try relative to executable
		ex, err := os.Executable()
		if err == nil {
			migrationsDir = filepath.Join(filepath.Dir(ex), "migrations")
		}
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}

	log.Printf("Database migrated to version: %d", version)
	return nil
}
