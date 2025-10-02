-- +goose Up
-- +goose StatementBegin

-- Add event notification role to guild_config
ALTER TABLE guild_config ADD COLUMN event_notification_role_id INTEGER;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE guild_config DROP COLUMN event_notification_role_id;

-- +goose StatementEnd
