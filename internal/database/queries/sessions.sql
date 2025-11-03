-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ? AND user_id = ?;

-- name: GetSessionBySlug :one
SELECT * FROM sessions
WHERE slug = ? AND user_id = ?;

-- name: ListSessions :many
SELECT * FROM sessions
WHERE user_id = ? AND slug != 'inbox'
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: CountSessions :one
SELECT COUNT(*) FROM sessions
WHERE user_id = ? AND slug != 'inbox';

-- name: CreateSession :one
INSERT INTO sessions (
    id, slug, title, description, avatar, background_color,
    type, user_id, group_id, client_id, pinned,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ? AND user_id = ?;

-- name: SearchSessions :many
SELECT * FROM sessions
WHERE user_id = ? 
  AND (title LIKE ? OR description LIKE ?)
ORDER BY updated_at DESC
LIMIT ?;

-- name: GetSessionWithGroup :one
SELECT 
    s.*,
    sg.id as group_id,
    sg.name as group_name,
    sg.sort as group_sort
FROM sessions s
LEFT JOIN session_groups sg ON s.group_id = sg.id
WHERE s.id = ? AND s.user_id = ?;

-- name: ListSessionsByGroup :many
SELECT * FROM sessions
WHERE user_id = ? AND group_id = ?
ORDER BY updated_at DESC;

-- name: MoveSessionToGroup :exec
UPDATE sessions
SET group_id = ?, updated_at = ?
WHERE id = ? AND user_id = ?;

-- name: PinSession :exec
UPDATE sessions
SET pinned = ?, updated_at = ?
WHERE id = ? AND user_id = ?;

