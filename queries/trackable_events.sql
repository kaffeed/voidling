-- name: CreateTrackableEvent :one
INSERT INTO trackable_events (type, activity, is_active)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetTrackableEventByID :one
SELECT * FROM trackable_events
WHERE id = ?
LIMIT 1;

-- name: GetActiveTrackableEvents :many
SELECT * FROM trackable_events
WHERE is_active = 1
ORDER BY created_at DESC;

-- name: GetActiveTrackableEventsByType :many
SELECT * FROM trackable_events
WHERE type = ? AND is_active = 1
ORDER BY created_at DESC;

-- name: GetLastActiveEventByType :one
SELECT * FROM trackable_events
WHERE type = ? AND is_active = 1
ORDER BY id DESC
LIMIT 1;

-- name: DeactivateTrackableEvent :exec
UPDATE trackable_events
SET is_active = 0
WHERE id = ?;

-- name: CreateTrackableParticipation :one
INSERT INTO trackable_event_participations (event_id, account_link_id, starting_point)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetTrackableParticipation :one
SELECT * FROM trackable_event_participations
WHERE event_id = ? AND account_link_id = ?
LIMIT 1;

-- name: GetTrackableParticipationsByEvent :many
SELECT tep.*, al.discord_member_id, al.runescape_name
FROM trackable_event_participations tep
JOIN account_links al ON tep.account_link_id = al.id
WHERE tep.event_id = ?
ORDER BY tep.created_at;

-- name: UpdateTrackableParticipationEndPoint :exec
UPDATE trackable_event_participations
SET end_point = ?
WHERE id = ?;

-- name: GetEventWinners :many
SELECT
    al.runescape_name,
    al.discord_member_id,
    tep.starting_point,
    tep.end_point,
    (tep.end_point - tep.starting_point) as progress
FROM trackable_event_participations tep
JOIN account_links al ON tep.account_link_id = al.id
WHERE tep.event_id = ? AND tep.end_point IS NOT NULL
ORDER BY progress DESC
LIMIT 3;

-- name: GetAllEventWinnersByType :many
SELECT
    te.id as event_id,
    al.runescape_name,
    (tep.end_point - tep.starting_point) as progress
FROM trackable_events te
JOIN trackable_event_participations tep ON te.id = tep.event_id
JOIN account_links al ON tep.account_link_id = al.id
WHERE te.type = ? AND tep.end_point IS NOT NULL AND te.is_active = 0
ORDER BY te.created_at DESC;

-- name: CreateTrackableProgress :exec
INSERT INTO trackable_event_progress (participation_id, progress, fetched_at)
VALUES (?, ?, ?);

-- name: GetProgressForParticipation :many
SELECT * FROM trackable_event_progress
WHERE participation_id = ?
ORDER BY fetched_at DESC;
