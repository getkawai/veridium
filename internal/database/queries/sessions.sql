-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ?;

-- name: GetSessionBySlug :one
SELECT * FROM sessions
WHERE slug = ?;

-- name: GetSessionByIdOrSlug :one
SELECT * FROM sessions
WHERE (id = ? OR slug = ?);

-- name: ListSessions :many
SELECT * FROM sessions
WHERE slug != 'inbox'
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: ListAllSessions :many
SELECT * FROM sessions
ORDER BY updated_at DESC;

-- name: CountSessions :one
SELECT COUNT(*) FROM sessions
WHERE slug != 'inbox';

-- name: CountSessionsByDateRange :one
SELECT COUNT(*) FROM sessions
WHERE slug != 'inbox'
  AND created_at >= ?
  AND created_at <= ?;

-- name: CreateSession :one
INSERT INTO sessions (
    id, slug, title, description, avatar, background_color,
    type, group_id, pinned,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSession :one
UPDATE sessions
SET title = ?,
    description = ?,
    avatar = ?,
    background_color = ?,
    group_id = ?,
    pinned = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;

-- name: BatchDeleteSessions :exec
DELETE FROM sessions
WHERE id IN (sqlc.slice('ids'));

-- name: SearchSessions :many
SELECT * FROM sessions
WHERE (title LIKE ? OR description LIKE ?)
ORDER BY updated_at DESC
LIMIT ?;

-- name: SearchSessionsByKeyword :many
SELECT * FROM sessions
WHERE (title LIKE '%' || ? || '%' OR description LIKE '%' || ? || '%')
ORDER BY updated_at DESC
LIMIT 100;

-- name: GetSessionWithGroup :one
SELECT 
    s.*,
    sg.id as group_id,
    sg.name as group_name,
    sg.sort as group_sort
FROM sessions s
LEFT JOIN session_groups sg ON s.group_id = sg.id
WHERE s.id = ?;

-- name: ListSessionsByGroup :many
SELECT * FROM sessions
WHERE group_id = ?
ORDER BY updated_at DESC;

-- name: MoveSessionToGroup :exec
UPDATE sessions
SET group_id = ?, updated_at = ?
WHERE id = ?;

-- name: PinSession :exec
UPDATE sessions
SET pinned = ?, updated_at = ?
WHERE id = ?;

-- name: GetSessionRank :many
SELECT 
    s.id,
    a.title,
    a.avatar,
    a.background_color,
    COUNT(t.id) as topic_count
FROM sessions s
LEFT JOIN agents_to_sessions ats ON s.id = ats.session_id
LEFT JOIN agents a ON ats.agent_id = a.id
LEFT JOIN topics t ON s.id = t.session_id
WHERE s.slug != 'inbox'
GROUP BY s.id, a.title, a.avatar, a.background_color
ORDER BY topic_count DESC, s.updated_at DESC
LIMIT ?;

-- name: DuplicateSession :one
-- Duplicate a session by creating a new session with the same data but new IDs
-- Parameters: new_session_id, new_title, created_at, updated_at, source_session_id
INSERT INTO sessions (
    id, slug, title, description, avatar, background_color,
    type, group_id, pinned,
    created_at, updated_at
)
SELECT 
    ? as id,                -- new_session_id
    NULL as slug,           -- no slug for duplicated sessions
    ? as title,             -- new_title
    s.description,
    s.avatar,
    s.background_color,
    s.type,
    s.group_id,
    s.pinned,
    ? as created_at,        -- new created_at
    ? as updated_at         -- new updated_at
FROM sessions s
WHERE s.id = ?
RETURNING *;
