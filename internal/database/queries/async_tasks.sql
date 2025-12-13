-- name: GetAsyncTask :one
SELECT * FROM async_tasks WHERE id = ?;

-- name: ListAsyncTasks :many
SELECT * FROM async_tasks
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListAsyncTasksByStatus :many
SELECT * FROM async_tasks
WHERE status = ?
ORDER BY created_at DESC;

-- name: CreateAsyncTask :one
INSERT INTO async_tasks (
    id, type, status, error, duration, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateAsyncTask :one
UPDATE async_tasks
SET status = ?,
    error = ?,
    duration = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteAsyncTask :exec
DELETE FROM async_tasks WHERE id = ?;

-- name: GetAsyncTasksByIds :many
SELECT * FROM async_tasks
WHERE id IN (sqlc.slice('ids')) AND type = ?;

-- name: GetTimeoutTasks :many
SELECT id FROM async_tasks
WHERE id IN (sqlc.slice('ids'))
  AND status = ?
  AND created_at < ?;

-- name: UpdateTimeoutTasks :exec
UPDATE async_tasks
SET status = ?,
    error = ?,
    updated_at = ?
WHERE id IN (sqlc.slice('ids'));
