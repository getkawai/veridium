-- name: GetThread :one
SELECT * FROM threads WHERE id = ?;

-- name: ListThreadsByTopic :many
SELECT * FROM threads
WHERE topic_id = ?
ORDER BY last_active_at DESC;

-- name: CreateThread :one
INSERT INTO threads (
    id, title, type, status, topic_id, source_message_id,
    parent_thread_id, last_active_at,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING
RETURNING *;

-- name: UpdateThread :one
UPDATE threads
SET title = ?,
    status = ?,
    last_active_at = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteThread :exec
DELETE FROM threads WHERE id = ?;

-- name: ListAllThreads :many
SELECT * FROM threads
ORDER BY updated_at DESC;

-- name: DeleteAllThreads :exec
DELETE FROM threads;
