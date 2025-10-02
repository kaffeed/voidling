-- +goose Up
-- +goose StatementBegin

-- Add default timezone to guild_config
ALTER TABLE guild_config ADD COLUMN default_timezone TEXT DEFAULT 'UTC';

-- Create user timezone preferences table
CREATE TABLE user_timezone_preferences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    discord_user_id INTEGER NOT NULL UNIQUE,
    timezone TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_timezone_preferences_user_id ON user_timezone_preferences(discord_user_id);

-- Add timezone to schedulable_events for reference
ALTER TABLE schedulable_events ADD COLUMN timezone TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_timezone_preferences;
ALTER TABLE guild_config DROP COLUMN default_timezone;
ALTER TABLE schedulable_events DROP COLUMN timezone;

-- +goose StatementEnd
