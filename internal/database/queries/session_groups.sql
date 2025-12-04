-- name: GetSessionGroup :one
SELECT * FROM session_groups
WHERE id = ? AND user_id = ?;

-- name: ListSessionGroups :many
SELECT * FROM session_groups
WHERE user_id = ?
ORDER BY sort ASC, created_at DESC;

-- name: CreateSessionGroup :one
INSERT INTO session_groups (
    id, name, sort, user_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSessionGroup :one
UPDATE session_groups
SET name = ?, sort = ?, updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteSessionGroup :exec
DELETE FROM session_groups
WHERE id = ? AND user_id = ?;

-- name: CountSessionsInGroup :one
SELECT COUNT(*) FROM sessions
WHERE group_id = ? AND user_id = ?;

-- name: GetSessionGroupWithSessions :one
SELECT 
    sg.*,
    COUNT(s.id) as session_count
FROM session_groups sg
LEFT JOIN sessions s ON sg.id = s.group_id
WHERE sg.id = ? AND sg.user_id = ?
GROUP BY sg.id;

-- name: DeleteAllSessionGroups :exec
DELETE FROM session_groups WHERE user_id = ?;

-- name: UpdateSessionGroupOrder :exec
UPDATE session_groups
SET sort = ?, updated_at = ?
WHERE id = ? AND user_id = ?;

