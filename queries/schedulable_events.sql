-- name: CreateSchedulableEvent :one
INSERT INTO schedulable_events (type, activity, location, scheduled_at, discord_event_id)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetSchedulableEventByID :one
SELECT * FROM schedulable_events
WHERE id = ?
LIMIT 1;

-- name: GetSchedulableEvents :many
SELECT * FROM schedulable_events
ORDER BY scheduled_at DESC;

-- name: GetUpcomingSchedulableEvents :many
SELECT * FROM schedulable_events
WHERE scheduled_at > ?
ORDER BY scheduled_at ASC;

-- name: GetSchedulableEventsInTimeRange :many
SELECT * FROM schedulable_events
WHERE scheduled_at >= ? AND scheduled_at < ?
ORDER BY scheduled_at ASC;

-- name: CreateSchedulableParticipation :one
INSERT INTO schedulable_event_participations (event_id, account_link_id, notified)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSchedulableParticipation :one
SELECT * FROM schedulable_event_participations
WHERE event_id = ? AND account_link_id = ?
LIMIT 1;

-- name: GetSchedulableParticipationsByEvent :many
SELECT sep.*, al.discord_member_id, al.runescape_name
FROM schedulable_event_participations sep
JOIN account_links al ON sep.account_link_id = al.id
WHERE sep.event_id = ?
ORDER BY sep.created_at;

-- name: GetUnnotifiedParticipations :many
SELECT sep.*, al.discord_member_id, al.runescape_name, se.activity, se.location, se.scheduled_at, se.type
FROM schedulable_event_participations sep
JOIN account_links al ON sep.account_link_id = al.id
JOIN schedulable_events se ON sep.event_id = se.id
WHERE sep.notified = 0 AND se.scheduled_at >= ? AND se.scheduled_at < ?
ORDER BY se.scheduled_at ASC;

-- name: MarkParticipationAsNotified :exec
UPDATE schedulable_event_participations
SET notified = 1
WHERE id = ?;

-- name: DeleteSchedulableEvent :exec
DELETE FROM schedulable_events
WHERE id = ?;
