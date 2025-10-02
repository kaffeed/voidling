-- name: GetUserTimezone :one
SELECT * FROM user_timezone_preferences
WHERE discord_user_id = ?
LIMIT 1;

-- name: CreateUserTimezone :one
INSERT INTO user_timezone_preferences (discord_user_id, timezone)
VALUES (?, ?)
RETURNING *;

-- name: UpdateUserTimezone :exec
UPDATE user_timezone_preferences
SET timezone = ?, updated_at = CURRENT_TIMESTAMP
WHERE discord_user_id = ?;

-- name: UpsertUserTimezone :exec
INSERT INTO user_timezone_preferences (discord_user_id, timezone)
VALUES (?, ?)
ON CONFLICT(discord_user_id) DO UPDATE SET
    timezone = excluded.timezone,
    updated_at = CURRENT_TIMESTAMP;

-- name: DeleteUserTimezone :exec
DELETE FROM user_timezone_preferences
WHERE discord_user_id = ?;
