-- name: CreateWarning :one
INSERT INTO warnings (guild_id, user_id, moderator_id, reason)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetWarningsByUser :many
SELECT * FROM warnings
WHERE guild_id = ? AND user_id = ?
ORDER BY created_at DESC;

-- name: GetWarningsByGuild :many
SELECT * FROM warnings
WHERE guild_id = ?
ORDER BY created_at DESC;

-- name: GetWarningByID :one
SELECT * FROM warnings
WHERE id = ?
LIMIT 1;

-- name: SetGuildWarningChannel :one
INSERT INTO guild_warning_channels (guild_id, channel_id)
VALUES (?, ?)
ON CONFLICT(guild_id) DO UPDATE SET channel_id = excluded.channel_id
RETURNING *;

-- name: GetGuildWarningChannel :one
SELECT * FROM guild_warning_channels
WHERE guild_id = ?
LIMIT 1;

-- name: DeleteGuildWarningChannel :exec
DELETE FROM guild_warning_channels
WHERE guild_id = ?;
