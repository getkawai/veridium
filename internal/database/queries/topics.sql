-- name: GetTopic :one
SELECT * FROM topics WHERE id = ? AND user_id = ?;

-- name: ListTopics :many
SELECT * FROM topics
WHERE user_id = ? AND session_id = ?
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: CreateTopic :one
INSERT INTO topics (
    id, title, favorite, session_id, group_id, user_id, client_id,
    history_summary, metadata, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTopic :one
UPDATE topics
SET title = ?,
    history_summary = ?,
    metadata = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteTopic :exec
DELETE FROM topics WHERE id = ? AND user_id = ?;

-- name: ToggleTopicFavorite :exec
UPDATE topics SET favorite = ?, updated_at = ?
WHERE id = ? AND user_id = ?;

-- Threads

-- name: GetThread :one
SELECT * FROM threads WHERE id = ? AND user_id = ?;

-- name: ListThreadsByTopic :many
SELECT * FROM threads
WHERE topic_id = ? AND user_id = ?
ORDER BY last_active_at DESC;

-- name: CreateThread :one
INSERT INTO threads (
    id, title, type, status, topic_id, source_message_id,
    parent_thread_id, client_id, user_id, last_active_at,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateThread :one
UPDATE threads
SET title = ?,
    status = ?,
    last_active_at = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteThread :exec
DELETE FROM threads WHERE id = ? AND user_id = ?;

