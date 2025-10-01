-- +goose Up
-- +goose StatementBegin

-- Account links table
CREATE TABLE account_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    discord_member_id INTEGER NOT NULL,
    runescape_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_account_links_discord_member_id ON account_links(discord_member_id);
CREATE INDEX idx_account_links_is_active ON account_links(is_active);
CREATE UNIQUE INDEX idx_account_links_discord_member_active ON account_links(discord_member_id, is_active) WHERE is_active = 1;

-- Trackable events table (Boss of the Week, Skill of the Week)
CREATE TABLE trackable_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK(type IN ('BossOfTheWeek', 'SkillOfTheWeek')),
    activity TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_trackable_events_type ON trackable_events(type);
CREATE INDEX idx_trackable_events_is_active ON trackable_events(is_active);

-- Trackable event participations (join table with extra data)
CREATE TABLE trackable_event_participations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    account_link_id INTEGER NOT NULL,
    starting_point INTEGER NOT NULL,
    end_point INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES trackable_events(id) ON DELETE CASCADE,
    FOREIGN KEY (account_link_id) REFERENCES account_links(id) ON DELETE CASCADE,
    UNIQUE(event_id, account_link_id)
);

CREATE INDEX idx_trackable_participations_event_id ON trackable_event_participations(event_id);
CREATE INDEX idx_trackable_participations_account_id ON trackable_event_participations(account_link_id);

-- Schedulable events table (Mass events, Wildy Wednesday)
CREATE TABLE schedulable_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK(type IN ('Mass', 'WildyWednesday')),
    activity TEXT NOT NULL,
    location TEXT NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_schedulable_events_type ON schedulable_events(type);
CREATE INDEX idx_schedulable_events_scheduled_at ON schedulable_events(scheduled_at);

-- Schedulable event participations (join table)
CREATE TABLE schedulable_event_participations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    account_link_id INTEGER NOT NULL,
    notified BOOLEAN NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES schedulable_events(id) ON DELETE CASCADE,
    FOREIGN KEY (account_link_id) REFERENCES account_links(id) ON DELETE CASCADE,
    UNIQUE(event_id, account_link_id)
);

CREATE INDEX idx_schedulable_participations_event_id ON schedulable_event_participations(event_id);
CREATE INDEX idx_schedulable_participations_account_id ON schedulable_event_participations(account_link_id);
CREATE INDEX idx_schedulable_participations_notified ON schedulable_event_participations(notified);

-- Guild warning channels table
CREATE TABLE guild_warning_channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id INTEGER NOT NULL UNIQUE,
    channel_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Warnings table
CREATE TABLE warnings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    moderator_id INTEGER NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_warnings_guild_user ON warnings(guild_id, user_id);
CREATE INDEX idx_warnings_created_at ON warnings(created_at);

-- Trackable event progress snapshots (for progress tracking over time)
CREATE TABLE trackable_event_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    participation_id INTEGER NOT NULL,
    progress INTEGER NOT NULL,
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (participation_id) REFERENCES trackable_event_participations(id) ON DELETE CASCADE
);

CREATE INDEX idx_trackable_progress_participation_id ON trackable_event_progress(participation_id);
CREATE INDEX idx_trackable_progress_fetched_at ON trackable_event_progress(fetched_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS trackable_event_progress;
DROP TABLE IF EXISTS warnings;
DROP TABLE IF EXISTS guild_warning_channels;
DROP TABLE IF EXISTS schedulable_event_participations;
DROP TABLE IF EXISTS schedulable_events;
DROP TABLE IF EXISTS trackable_event_participations;
DROP TABLE IF EXISTS trackable_events;
DROP TABLE IF EXISTS account_links;

-- +goose StatementEnd
