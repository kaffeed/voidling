-- +goose Up
ALTER TABLE guild_config ADD COLUMN event_notification_channel_id INTEGER;

-- +goose Down
ALTER TABLE guild_config DROP COLUMN event_notification_channel_id;
