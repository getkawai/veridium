-- name: CreateRagEvalDataset :one
INSERT INTO rag_eval_datasets (
  id,
  name,
  description,
  user_id,
  created_at,
  updated_at
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRagEvalDataset :one
SELECT * FROM rag_eval_datasets
WHERE id = ? AND user_id = ?
LIMIT 1;

-- name: ListRagEvalDatasets :many
SELECT * FROM rag_eval_datasets
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: UpdateRagEvalDataset :exec
UPDATE rag_eval_datasets
SET
  name = COALESCE(?, name),
  description = COALESCE(?, description),
  updated_at = ?
WHERE id = ? AND user_id = ?;

-- name: DeleteRagEvalDataset :exec
DELETE FROM rag_eval_datasets
WHERE id = ? AND user_id = ?;

-- name: DeleteAllRagEvalDatasets :exec
DELETE FROM rag_eval_datasets
WHERE user_id = ?;

