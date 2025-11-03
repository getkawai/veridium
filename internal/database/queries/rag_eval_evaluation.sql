-- name: CreateRagEvalEvaluation :one
INSERT INTO rag_eval_evaluations (
  id,
  name,
  dataset_id,
  config,
  status,
  user_id,
  created_at,
  updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRagEvalEvaluation :one
SELECT * FROM rag_eval_evaluations
WHERE id = ? AND user_id = ?
LIMIT 1;

-- name: ListRagEvalEvaluationsByDataset :many
SELECT * FROM rag_eval_evaluations
WHERE dataset_id = ? AND user_id = ?
ORDER BY created_at DESC;

-- name: UpdateRagEvalEvaluation :exec
UPDATE rag_eval_evaluations
SET
  name = COALESCE(?, name),
  config = COALESCE(?, config),
  status = COALESCE(?, status),
  updated_at = ?
WHERE id = ? AND user_id = ?;

-- name: UpdateRagEvalEvaluationStatus :exec
UPDATE rag_eval_evaluations
SET
  status = ?,
  updated_at = ?
WHERE id = ? AND user_id = ?;

-- name: DeleteRagEvalEvaluation :exec
DELETE FROM rag_eval_evaluations
WHERE id = ? AND user_id = ?;

-- name: DeleteRagEvalEvaluationsByDataset :exec
DELETE FROM rag_eval_evaluations
WHERE dataset_id = ? AND user_id = ?;

