-- Generation Topics

-- name: GetGenerationTopic :one
SELECT * FROM generation_topics WHERE id = ? AND user_id = ?;

-- name: ListGenerationTopics :many
SELECT * FROM generation_topics
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CreateGenerationTopic :one
INSERT INTO generation_topics (
    id, user_id, title, cover_url, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateGenerationTopic :one
UPDATE generation_topics
SET title = ?,
    cover_url = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteGenerationTopic :exec
DELETE FROM generation_topics WHERE id = ? AND user_id = ?;

-- Generation Batches

-- name: GetGenerationBatch :one
SELECT * FROM generation_batches WHERE id = ? AND user_id = ?;

-- name: ListGenerationBatches :many
SELECT * FROM generation_batches
WHERE generation_topic_id = ? AND user_id = ?
ORDER BY created_at DESC;

-- name: CreateGenerationBatch :one
INSERT INTO generation_batches (
    id, user_id, generation_topic_id, provider, model, prompt,
    width, height, ratio, config, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteGenerationBatch :exec
DELETE FROM generation_batches WHERE id = ? AND user_id = ?;

-- Generations

-- name: GetGeneration :one
SELECT * FROM generations WHERE id = ? AND user_id = ?;

-- name: ListGenerations :many
SELECT * FROM generations
WHERE generation_batch_id = ? AND user_id = ?
ORDER BY created_at ASC;

-- name: CreateGeneration :one
INSERT INTO generations (
    id, user_id, generation_batch_id, async_task_id, file_id,
    seed, asset, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateGeneration :one
UPDATE generations
SET async_task_id = ?,
    file_id = ?,
    asset = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteGeneration :exec
DELETE FROM generations WHERE id = ? AND user_id = ?;

-- name: GetGenerationWithAsyncTask :one
SELECT 
    g.*,
    at.status as async_task_status,
    at.error as async_task_error
FROM generations g
LEFT JOIN async_tasks at ON g.async_task_id = at.id
WHERE g.id = ? AND g.user_id = ?;

