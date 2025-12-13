-- Chunks

-- name: GetChunk :one
SELECT * FROM chunks WHERE id = ?;

-- name: ListChunks :many
SELECT * FROM chunks
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateChunk :one
INSERT INTO chunks (
    id, document_id, text, abstract, metadata, chunk_index, type
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetChunksByDocumentID :many
SELECT * FROM chunks
WHERE document_id = ?
ORDER BY chunk_index ASC;

-- name: UpdateChunk :one
UPDATE chunks
SET text = ?,
    abstract = ?,
    metadata = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteChunk :exec
DELETE FROM chunks WHERE id = ?;

-- Unstructured Chunks

-- name: GetUnstructuredChunk :one
SELECT * FROM unstructured_chunks WHERE id = ?;

-- name: ListUnstructuredChunksByFile :many
SELECT * FROM unstructured_chunks
WHERE file_id = ?
ORDER BY chunk_index ASC;

-- name: CreateUnstructuredChunk :one
INSERT INTO unstructured_chunks (
    id, text, metadata, chunk_index, type, parent_id, composite_id,
    file_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteUnstructuredChunk :exec
DELETE FROM unstructured_chunks WHERE id = ?;

-- File Chunks

-- name: LinkFileToChunk :exec
INSERT INTO file_chunks (file_id, chunk_id)
VALUES (?, ?);

-- name: UnlinkFileFromChunk :exec
DELETE FROM file_chunks
WHERE file_id = ? AND chunk_id = ?;

-- name: GetFileChunks :many
SELECT c.* FROM chunks c
INNER JOIN file_chunks fc ON c.id = fc.chunk_id
WHERE fc.file_id = ?
ORDER BY c.chunk_index ASC
LIMIT ? OFFSET ?;

-- name: GetFileChunksWithMetadata :many
SELECT 
    c.id,
    c.text,
    c.abstract,
    c.metadata,
    c.chunk_index,
    c.type,
    c.created_at,
    c.updated_at
FROM chunks c
INNER JOIN file_chunks fc ON c.id = fc.chunk_id
WHERE fc.file_id = ?
ORDER BY c.chunk_index ASC
LIMIT ? OFFSET ?;

-- name: GetChunksTextByFileId :many
SELECT c.id, c.text, c.metadata, c.type
FROM chunks c
INNER JOIN file_chunks fc ON c.id = fc.chunk_id
WHERE fc.file_id = ?;

-- name: CountChunksByFileId :one
SELECT COUNT(*) as count
FROM file_chunks
WHERE file_id = ?;

-- name: CountChunksByFileIds :many
SELECT file_id, COUNT(*) as count
FROM file_chunks
GROUP BY file_id;

-- name: GetOrphanedChunks :many
SELECT c.id as chunk_id
FROM chunks c
LEFT JOIN file_chunks fc ON c.id = fc.chunk_id
WHERE fc.file_id IS NULL;

-- name: BatchDeleteChunks :exec
DELETE FROM chunks
WHERE id IN (sqlc.slice('ids'));

-- name: GetChunksByIDs :many
SELECT c.*, d.file_id
FROM chunks c
LEFT JOIN documents d ON c.document_id = d.id
WHERE c.id IN (sqlc.slice('ids'));

-- RAG Evaluation

-- name: GetRagEvalDataset :one
SELECT * FROM rag_eval_datasets WHERE id = ?;

-- name: ListRagEvalDatasets :many
SELECT * FROM rag_eval_datasets
ORDER BY created_at DESC;

-- name: CreateRagEvalDataset :one
INSERT INTO rag_eval_datasets (
    id, name, description
) VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateRagEvalDataset :exec
UPDATE rag_eval_datasets
SET name = COALESCE(?, name),
    description = COALESCE(?, description),
    updated_at = ?
WHERE id = ?;

-- name: DeleteRagEvalDataset :exec
DELETE FROM rag_eval_datasets WHERE id = ?;

-- name: GetRagEvalDatasetRecord :one
SELECT * FROM rag_eval_dataset_records WHERE id = ?;

-- name: ListRagEvalDatasetRecords :many
SELECT * FROM rag_eval_dataset_records
WHERE dataset_id = ?
ORDER BY created_at ASC;

-- name: CreateRagEvalDatasetRecord :one
INSERT INTO rag_eval_dataset_records (
    id, dataset_id, query, reference_answer, reference_contexts,
    metadata
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRagEvalDatasetRecord :exec
UPDATE rag_eval_dataset_records
SET "query" = COALESCE(?, "query"),
    reference_answer = COALESCE(?, reference_answer),
    reference_contexts = COALESCE(?, reference_contexts),
    metadata = COALESCE(?, metadata),
    updated_at = ?
WHERE id = ?;

-- name: DeleteRagEvalDatasetRecord :exec
DELETE FROM rag_eval_dataset_records WHERE id = ?;

-- RAG Eval Evaluations

-- name: CreateRagEvalEvaluation :one
INSERT INTO rag_eval_evaluations (
    id, name, dataset_id, config, status
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRagEvalEvaluation :one
SELECT * FROM rag_eval_evaluations WHERE id = ? LIMIT 1;

-- name: ListRagEvalEvaluationsByDataset :many
SELECT * FROM rag_eval_evaluations
WHERE dataset_id = ?
ORDER BY created_at DESC;

-- name: UpdateRagEvalEvaluation :exec
UPDATE rag_eval_evaluations
SET name = COALESCE(?, name),
    config = COALESCE(?, config),
    status = COALESCE(?, status),
    updated_at = ?
WHERE id = ?;

-- name: DeleteRagEvalEvaluation :exec
DELETE FROM rag_eval_evaluations WHERE id = ?;

-- RAG Eval Evaluation Records

-- name: CreateRagEvalEvaluationRecord :one
INSERT INTO rag_eval_evaluation_records (
    id, evaluation_id, dataset_record_id,
    retrieved_contexts, generated_answer, metrics
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRagEvalEvaluationRecord :one
SELECT * FROM rag_eval_evaluation_records WHERE id = ? LIMIT 1;

-- name: ListRagEvalEvaluationRecordsByEvaluation :many
SELECT * FROM rag_eval_evaluation_records
WHERE evaluation_id = ?
ORDER BY created_at DESC;

-- name: UpdateRagEvalEvaluationRecord :exec
UPDATE rag_eval_evaluation_records
SET retrieved_contexts = COALESCE(?, retrieved_contexts),
    generated_answer = COALESCE(?, generated_answer),
    metrics = COALESCE(?, metrics),
    updated_at = ?
WHERE id = ?;

-- name: DeleteRagEvalEvaluationRecord :exec
DELETE FROM rag_eval_evaluation_records WHERE id = ?;

-- name: BatchInsertRagEvalEvaluationRecords :exec
INSERT INTO rag_eval_evaluation_records (
    id, evaluation_id, dataset_record_id,
    retrieved_contexts, generated_answer, metrics,
    created_at, updated_at
)
SELECT 
    json_extract(value, '$.id'),
    json_extract(value, '$.evaluation_id'),
    json_extract(value, '$.dataset_record_id'),
    json_extract(value, '$.retrieved_contexts'),
    json_extract(value, '$.generated_answer'),
    json_extract(value, '$.metrics'),
    CAST(json_extract(value, '$.created_at') AS INTEGER),
    CAST(json_extract(value, '$.updated_at') AS INTEGER)
FROM json_each(?);
