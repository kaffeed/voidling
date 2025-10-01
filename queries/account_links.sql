-- name: GetAccountLinkByDiscordID :one
SELECT * FROM account_links
WHERE discord_member_id = ? AND is_active = 1
LIMIT 1;

-- name: GetAccountLinkByID :one
SELECT * FROM account_links
WHERE id = ?
LIMIT 1;

-- name: CreateAccountLink :one
INSERT INTO account_links (discord_member_id, runescape_name, is_active)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeactivateAccountLink :exec
UPDATE account_links
SET is_active = 0, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeactivateAllAccountLinksForUser :exec
UPDATE account_links
SET is_active = 0, updated_at = CURRENT_TIMESTAMP
WHERE discord_member_id = ?;

-- name: GetExistingAccountLink :one
SELECT * FROM account_links
WHERE discord_member_id = ? AND LOWER(runescape_name) = LOWER(?)
LIMIT 1;

-- name: ActivateAccountLink :exec
UPDATE account_links
SET is_active = 1, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: GetAllAccountLinksForUser :many
SELECT * FROM account_links
WHERE discord_member_id = ?
ORDER BY created_at DESC;

-- name: GetAccountLinkByUsername :one
SELECT * FROM account_links
WHERE runescape_name = ? AND is_active = 1
LIMIT 1;
