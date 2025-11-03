-- Chunks

-- name: GetChunk :one
SELECT * FROM chunks WHERE id = ? AND user_id = ?;

-- name: ListChunks :many
SELECT * FROM chunks
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateChunk :one
INSERT INTO chunks (
    id, text, abstract, metadata, chunk_index, type, client_id,
    user_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateChunk :one
UPDATE chunks
SET text = ?,
    abstract = ?,
    metadata = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteChunk :exec
DELETE FROM chunks WHERE id = ? AND user_id = ?;

-- Embeddings

-- name: GetEmbedding :one
SELECT * FROM embeddings WHERE id = ? AND user_id = ?;

-- name: GetEmbeddingByChunk :one
SELECT * FROM embeddings WHERE chunk_id = ? AND user_id = ?;

-- name: CreateEmbedding :one
INSERT INTO embeddings (
    id, chunk_id, embeddings, model, client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteEmbedding :exec
DELETE FROM embeddings WHERE id = ? AND user_id = ?;

-- Unstructured Chunks

-- name: GetUnstructuredChunk :one
SELECT * FROM unstructured_chunks WHERE id = ? AND user_id = ?;

-- name: ListUnstructuredChunksByFile :many
SELECT * FROM unstructured_chunks
WHERE file_id = ? AND user_id = ?
ORDER BY chunk_index ASC;

-- name: CreateUnstructuredChunk :one
INSERT INTO unstructured_chunks (
    id, text, metadata, chunk_index, type, parent_id, composite_id,
    client_id, user_id, file_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteUnstructuredChunk :exec
DELETE FROM unstructured_chunks WHERE id = ? AND user_id = ?;

-- File Chunks

-- name: LinkFileToChunk :exec
INSERT INTO file_chunks (file_id, chunk_id, created_at, user_id)
VALUES (?, ?, ?, ?);

-- name: UnlinkFileFromChunk :exec
DELETE FROM file_chunks
WHERE file_id = ? AND chunk_id = ? AND user_id = ?;

-- name: GetFileChunks :many
SELECT c.* FROM chunks c
INNER JOIN file_chunks fc ON c.id = fc.chunk_id
WHERE fc.file_id = ? AND fc.user_id = ?
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
WHERE fc.file_id = ? AND c.user_id = ?
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
WHERE file_id = ? AND user_id = ?;

-- name: CountChunksByFileIds :many
SELECT file_id, COUNT(*) as count
FROM file_chunks
WHERE user_id = ?
GROUP BY file_id;

-- name: GetOrphanedChunks :many
SELECT c.id as chunk_id
FROM chunks c
LEFT JOIN file_chunks fc ON c.id = fc.chunk_id
WHERE fc.file_id IS NULL;

-- name: BatchDeleteChunks :exec
DELETE FROM chunks
WHERE id IN (sqlc.slice('ids')) AND user_id = ?;

-- Semantic search - fetch chunks with embeddings for JS similarity calculation
-- name: GetChunksWithEmbeddings :many
SELECT 
    c.id,
    c.text,
    c.metadata,
    c.chunk_index,
    c.type,
    e.embeddings as chunk_embedding,
    fc.file_id,
    f.name as file_name
FROM chunks c
LEFT JOIN embeddings e ON c.id = e.chunk_id
LEFT JOIN file_chunks fc ON c.id = fc.chunk_id
LEFT JOIN files f ON fc.file_id = f.id
WHERE c.user_id = ?;

-- name: GetChunksWithEmbeddingsByFileIds :many
SELECT 
    c.id,
    c.text,
    c.metadata,
    c.chunk_index,
    c.type,
    e.embeddings as chunk_embedding,
    fc.file_id,
    f.name as file_name
FROM chunks c
LEFT JOIN embeddings e ON c.id = e.chunk_id
LEFT JOIN file_chunks fc ON c.id = fc.chunk_id
LEFT JOIN files f ON fc.file_id = f.id
WHERE fc.file_id = ? AND fc.user_id = ?
ORDER BY c.chunk_index ASC;

-- RAG Evaluation

-- name: GetRagEvalDataset :one
SELECT * FROM rag_eval_datasets WHERE id = ? AND user_id = ?;

-- name: ListRagEvalDatasets :many
SELECT * FROM rag_eval_datasets
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CreateRagEvalDataset :one
INSERT INTO rag_eval_datasets (
    id, name, description, user_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteRagEvalDataset :exec
DELETE FROM rag_eval_datasets WHERE id = ? AND user_id = ?;

-- name: GetRagEvalDatasetRecord :one
SELECT * FROM rag_eval_dataset_records WHERE id = ? AND user_id = ?;

-- name: ListRagEvalDatasetRecords :many
SELECT * FROM rag_eval_dataset_records
WHERE dataset_id = ? AND user_id = ?
ORDER BY created_at ASC;

-- name: CreateRagEvalDatasetRecord :one
INSERT INTO rag_eval_dataset_records (
    id, dataset_id, query, reference_answer, reference_contexts,
    metadata, user_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteRagEvalDatasetRecord :exec
DELETE FROM rag_eval_dataset_records WHERE id = ? AND user_id = ?;

