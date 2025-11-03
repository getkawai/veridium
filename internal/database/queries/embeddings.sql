-- name: GetEmbeddingsItem :one
SELECT * FROM embeddings WHERE id = ? AND user_id = ?;

-- name: ListEmbeddingsItems :many
SELECT * FROM embeddings
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CreateEmbeddingsItem :one
INSERT INTO embeddings (
    id, chunk_id, embeddings, model, client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: BulkCreateEmbeddingsItems :exec
INSERT INTO embeddings (
    id, chunk_id, embeddings, model, client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(chunk_id) DO NOTHING;

-- name: DeleteEmbeddingsItem :exec
DELETE FROM embeddings WHERE id = ? AND user_id = ?;

-- name: CountEmbeddingsItems :one
SELECT COUNT(*) FROM embeddings WHERE user_id = ?;
