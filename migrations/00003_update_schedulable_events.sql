-- +goose Up
-- Update schedulable_events to store Discord event ID instead of managing our own events
ALTER TABLE schedulable_events ADD COLUMN discord_event_id TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE schedulable_events DROP COLUMN discord_event_id;
