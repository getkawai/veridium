-- Generation Topics

-- name: GetGenerationTopic :one
SELECT * FROM generation_topics WHERE id = ?;

-- name: ListGenerationTopics :many
SELECT * FROM generation_topics
ORDER BY created_at DESC;

-- name: CreateGenerationTopic :one
INSERT INTO generation_topics (
    id, title, cover_url, created_at, updated_at
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateGenerationTopic :one
UPDATE generation_topics
SET title = ?,
    cover_url = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteGenerationTopic :exec
DELETE FROM generation_topics WHERE id = ?;

-- Generation Batches

-- name: GetGenerationBatch :one
SELECT * FROM generation_batches WHERE id = ?;

-- name: ListGenerationBatches :many
SELECT * FROM generation_batches
WHERE generation_topic_id = ?
ORDER BY created_at DESC;

-- name: CreateGenerationBatch :one
INSERT INTO generation_batches (
    id, generation_topic_id, provider, model, prompt,
    width, height, ratio, config, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteGenerationBatch :exec
DELETE FROM generation_batches WHERE id = ?;

-- Generations

-- name: GetGeneration :one
SELECT * FROM generations WHERE id = ?;

-- name: ListGenerations :many
SELECT * FROM generations
WHERE generation_batch_id = ?
ORDER BY created_at ASC;

-- name: CreateGeneration :one
INSERT INTO generations (
    id, generation_batch_id, async_task_id, file_id,
    seed, asset, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateGeneration :one
UPDATE generations
SET async_task_id = ?,
    file_id = ?,
    asset = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteGeneration :exec
DELETE FROM generations WHERE id = ?;

-- Complex queries with JOINs for optimization

-- name: GetGenerationBatchWithGenerations :one
SELECT 
    gb.*,
    GROUP_CONCAT(g.id) as generation_ids
FROM generation_batches gb
LEFT JOIN generations g ON gb.id = g.generation_batch_id
WHERE gb.id = ?
GROUP BY gb.id;

-- name: ListGenerationBatchesWithGenerations :many
SELECT 
    gb.id as batch_id,
    gb.generation_topic_id,
    gb.provider,
    gb.model,
    gb.prompt,
    gb.width,
    gb.height,
    gb.ratio,
    gb.config,
    gb.created_at as batch_created_at,
    gb.updated_at as batch_updated_at,
    g.id as gen_id,
    g.async_task_id,
    g.file_id,
    g.seed,
    g.asset,
    g.created_at as gen_created_at,
    g.updated_at as gen_updated_at,
    at.id as task_id,
    at.status as task_state,
    at.error as task_error
FROM generation_batches gb
LEFT JOIN generations g ON gb.id = g.generation_batch_id
LEFT JOIN async_tasks at ON g.async_task_id = at.id
WHERE gb.generation_topic_id = ?
ORDER BY gb.created_at ASC, g.created_at ASC, g.id ASC;

-- name: GetGenerationTopicWithBatches :one
SELECT 
    gt.*,
    COUNT(DISTINCT gb.id) as batch_count
FROM generation_topics gt
LEFT JOIN generation_batches gb ON gt.id = gb.generation_topic_id
WHERE gt.id = ?
GROUP BY gt.id;

-- name: ListGenerationTopicsWithCounts :many
SELECT 
    gt.*,
    COUNT(DISTINCT gb.id) as batch_count,
    COUNT(DISTINCT g.id) as generation_count
FROM generation_topics gt
LEFT JOIN generation_batches gb ON gt.id = gb.generation_topic_id
LEFT JOIN generations g ON gb.id = g.generation_batch_id
GROUP BY gt.id
ORDER BY gt.updated_at DESC;

-- Delete queries with cascade information

-- name: GetGenerationBatchAssets :many
SELECT g.asset
FROM generations g
WHERE g.generation_batch_id = ?;

-- name: GetGenerationTopicAssets :many
SELECT g.asset, gt.cover_url
FROM generation_topics gt
LEFT JOIN generation_batches gb ON gt.id = gb.generation_topic_id
LEFT JOIN generations g ON gb.id = g.generation_batch_id
WHERE gt.id = ?;

-- name: GetGenerationWithAsyncTask :one
SELECT 
    g.*,
    at.status as async_task_status,
    at.error as async_task_error
FROM generations g
LEFT JOIN async_tasks at ON g.async_task_id = at.id
WHERE g.id = ?;
