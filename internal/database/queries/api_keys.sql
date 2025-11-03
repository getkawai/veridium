-- name: GetAPIKey :one
SELECT * FROM api_keys WHERE id = ? AND user_id = ?;

-- name: GetAPIKeyByKey :one
SELECT * FROM api_keys WHERE key = ?;

-- name: ListAPIKeys :many
SELECT * FROM api_keys
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CreateAPIKey :one
INSERT INTO api_keys (
    name, key, enabled, expires_at, last_used_at, user_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateAPIKey :one
UPDATE api_keys
SET name = ?,
    enabled = ?,
    expires_at = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys
SET last_used_at = ?
WHERE id = ?;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE id = ? AND user_id = ?;

-- name: DeleteAllAPIKeys :exec
DELETE FROM api_keys WHERE user_id = ?;

