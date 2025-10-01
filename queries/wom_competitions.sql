-- name: CreateWOMCompetition :one
INSERT INTO wom_competitions (wom_competition_id, verification_code, discord_thread_id, metric, type)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetWOMCompetitionByID :one
SELECT * FROM wom_competitions
WHERE id = ?
LIMIT 1;

-- name: GetWOMCompetitionByWOMID :one
SELECT * FROM wom_competitions
WHERE wom_competition_id = ?
LIMIT 1;

-- name: GetWOMCompetitionByThreadID :one
SELECT * FROM wom_competitions
WHERE discord_thread_id = ?
LIMIT 1;

-- name: GetWOMCompetitionsByType :many
SELECT * FROM wom_competitions
WHERE type = ?
ORDER BY created_at DESC;

-- name: GetLatestWOMCompetitionByType :one
SELECT * FROM wom_competitions
WHERE type = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteWOMCompetition :exec
DELETE FROM wom_competitions
WHERE id = ?;
