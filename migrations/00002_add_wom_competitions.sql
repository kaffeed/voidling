-- +goose Up
-- +goose StatementBegin

-- WOM competitions table (replaces trackable_events)
CREATE TABLE wom_competitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wom_competition_id INTEGER NOT NULL UNIQUE,
    verification_code TEXT NOT NULL,
    discord_thread_id TEXT NOT NULL,
    metric TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('BOSS_OF_THE_WEEK', 'SKILL_OF_THE_WEEK')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wom_competitions_wom_id ON wom_competitions(wom_competition_id);
CREATE INDEX idx_wom_competitions_thread_id ON wom_competitions(discord_thread_id);
CREATE INDEX idx_wom_competitions_type ON wom_competitions(type);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wom_competitions;

-- +goose StatementEnd
