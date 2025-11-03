-- name: CreateRagEvalDatasetRecord :one
INSERT INTO rag_eval_dataset_records (
  id,
  dataset_id,
  "query",
  reference_answer,
  reference_contexts,
  metadata,
  user_id,
  created_at,
  updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRagEvalDatasetRecord :one
SELECT * FROM rag_eval_dataset_records
WHERE id = ? AND user_id = ?
LIMIT 1;

-- name: ListRagEvalDatasetRecordsByDataset :many
SELECT * FROM rag_eval_dataset_records
WHERE dataset_id = ? AND user_id = ?
ORDER BY created_at DESC;

-- name: UpdateRagEvalDatasetRecord :exec
UPDATE rag_eval_dataset_records
SET
  "query" = COALESCE(?, "query"),
  reference_answer = COALESCE(?, reference_answer),
  reference_contexts = COALESCE(?, reference_contexts),
  metadata = COALESCE(?, metadata),
  updated_at = ?
WHERE id = ? AND user_id = ?;

-- name: DeleteRagEvalDatasetRecord :exec
DELETE FROM rag_eval_dataset_records
WHERE id = ? AND user_id = ?;

-- name: DeleteRagEvalDatasetRecordsByDataset :exec
DELETE FROM rag_eval_dataset_records
WHERE dataset_id = ? AND user_id = ?;

-- name: BatchInsertRagEvalDatasetRecords :exec
INSERT INTO rag_eval_dataset_records (
  id,
  dataset_id,
  "query",
  reference_answer,
  reference_contexts,
  metadata,
  user_id,
  created_at,
  updated_at
)
SELECT 
  json_extract(value, '$.id') AS id,
  json_extract(value, '$.dataset_id') AS dataset_id,
  json_extract(value, '$.query') AS "query",
  json_extract(value, '$.reference_answer') AS reference_answer,
  json_extract(value, '$.reference_contexts') AS reference_contexts,
  json_extract(value, '$.metadata') AS metadata,
  json_extract(value, '$.user_id') AS user_id,
  CAST(json_extract(value, '$.created_at') AS INTEGER) AS created_at,
  CAST(json_extract(value, '$.updated_at') AS INTEGER) AS updated_at
FROM json_each(?);

