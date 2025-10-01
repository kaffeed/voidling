package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DiscordToken       string
	DatabasePath       string
	LogLevel           string
	GuildID            string // Optional: specific guild for slash command registration
	CoordinatorRoleID  string // Optional: specific role ID for Coordinator permissions
}

// Load loads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable is required")
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		// Default to local data directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath = filepath.Join(homeDir, ".voidbound", "voidbound.db")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	guildID := os.Getenv("DISCORD_GUILD_ID")
	coordinatorRoleID := os.Getenv("COORDINATOR_ROLE_ID")

	return &Config{
		DiscordToken:      token,
		DatabasePath:      dbPath,
		LogLevel:          logLevel,
		GuildID:           guildID,
		CoordinatorRoleID: coordinatorRoleID,
	}, nil
}
