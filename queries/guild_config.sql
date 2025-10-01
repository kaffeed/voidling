-- name: GetGuildConfig :one
SELECT * FROM guild_config
WHERE guild_id = ?
LIMIT 1;

-- name: CreateGuildConfig :one
INSERT INTO guild_config (guild_id, coordinator_role_id)
VALUES (?, ?)
RETURNING *;

-- name: UpdateCoordinatorRole :exec
UPDATE guild_config
SET coordinator_role_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE guild_id = ?;

-- name: UpsertGuildConfig :exec
INSERT INTO guild_config (guild_id, coordinator_role_id)
VALUES (?, ?)
ON CONFLICT(guild_id) DO UPDATE SET
    coordinator_role_id = excluded.coordinator_role_id,
    updated_at = CURRENT_TIMESTAMP;
