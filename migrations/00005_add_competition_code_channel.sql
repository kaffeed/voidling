-- +goose Up
ALTER TABLE guild_config ADD COLUMN competition_code_channel_id INTEGER;

-- +goose Down
ALTER TABLE guild_config DROP COLUMN competition_code_channel_id;
