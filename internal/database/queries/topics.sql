-- name: GetTopic :one
SELECT * FROM topics WHERE id = ?;

-- name: ListTopics :many
SELECT * FROM topics
WHERE (COALESCE(session_id, '') = COALESCE(?, ''))
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: CreateTopic :one
INSERT INTO topics (
    id, title, favorite, session_id, group_id,
    history_summary, metadata, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTopic :one
UPDATE topics
SET title = ?,
    history_summary = ?,
    metadata = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTopic :exec
DELETE FROM topics WHERE id = ?;

-- name: CountTopicsBySession :one
SELECT COUNT(*) FROM topics
WHERE session_id = ?;

-- name: CountTopics :one
SELECT COUNT(*) FROM topics;

-- name: CountTopicsByDateRange :one
SELECT COUNT(*) FROM topics
WHERE created_at >= ?
  AND created_at <= ?;

-- name: ListAllTopics :many
SELECT * FROM topics
ORDER BY updated_at DESC;

-- name: SearchTopicsByTitle :many
SELECT * FROM topics
WHERE title LIKE ?
  AND (? = '' OR session_id = ? OR group_id = ?)
ORDER BY updated_at DESC;

-- name: SearchTopicsByMessageContent :many
SELECT DISTINCT t.*
FROM topics t
INNER JOIN messages m ON t.id = m.topic_id
WHERE m.content LIKE ?
  AND (? = '' OR t.session_id = ? OR t.group_id = ?)
ORDER BY t.updated_at DESC;

-- name: RankTopics :many
SELECT
    t.id,
    t.title,
    t.session_id,
    COUNT(m.id) as count
FROM topics t
LEFT JOIN messages m ON t.id = m.topic_id
GROUP BY t.id, t.title, t.session_id
HAVING COUNT(m.id) > 0
ORDER BY count DESC, t.updated_at DESC
LIMIT ?;

-- name: BatchDeleteTopics :exec
DELETE FROM topics
WHERE id IN (sqlc.slice('ids'));

-- name: DeleteTopicsBySession :exec
DELETE FROM topics
WHERE session_id = ?;

-- name: DeleteTopicsByGroup :exec
DELETE FROM topics
WHERE group_id = ?;

-- name: DeleteAllTopics :exec
DELETE FROM topics;

-- name: UpdateMessagesTopicId :exec
UPDATE messages
SET topic_id = ?
WHERE id IN (sqlc.slice('ids'));

-- name: GetMessagesByTopicId :many
SELECT * FROM messages
WHERE topic_id = ?
ORDER BY created_at ASC;

-- name: ToggleTopicFavorite :exec
UPDATE topics SET favorite = ?, updated_at = ?
WHERE id = ?;

-- name: DuplicateTopic :one
-- Duplicate a topic with a new ID and title
INSERT INTO topics (
    id, title, favorite, session_id, group_id,
    history_summary, metadata, created_at, updated_at
)
SELECT 
    ? as id,                -- new_topic_id
    ? as title,             -- new_title
    t.favorite,
    t.session_id,
    t.group_id,
    t.history_summary,
    t.metadata,
    ? as created_at,        -- new created_at
    ? as updated_at         -- new updated_at
FROM topics t
WHERE t.id = ?
RETURNING *;
