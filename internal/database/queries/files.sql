-- name: GetFile :one
SELECT * FROM files WHERE id = ?;

-- name: ListFiles :many
SELECT * FROM files
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateFile :one
INSERT INTO files (
    file_type, file_hash, name, size, url, source, metadata
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateFile :one
UPDATE files
SET name = ?,
    metadata = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: UpdateFileChunkStats :exec
UPDATE files
SET chunk_count = ?,
    chunking_status = ?,
    embedding_status = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteFile :exec
DELETE FROM files WHERE id = ?;

-- Global Files

-- name: GetGlobalFile :one
SELECT * FROM global_files WHERE hash_id = ?;

-- name: GetGlobalFileByHash :one
SELECT * FROM global_files WHERE hash_id = ?;

-- name: CreateGlobalFile :one
INSERT INTO global_files (
    hash_id, file_type, size, url, metadata, creator
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;


-- Knowledge Bases

-- name: GetKnowledgeBase :one
SELECT * FROM knowledge_bases WHERE id = ?;

-- name: ListKnowledgeBases :many
SELECT * FROM knowledge_bases
ORDER BY created_at DESC;

-- name: CreateKnowledgeBase :one
INSERT INTO knowledge_bases (
    id, name, description, avatar, type,
    is_public, settings
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateKnowledgeBase :one
UPDATE knowledge_bases
SET name = ?,
    description = ?,
    avatar = ?,
    settings = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteKnowledgeBase :exec
DELETE FROM knowledge_bases WHERE id = ?;

-- name: DeleteAllKnowledgeBases :exec
DELETE FROM knowledge_bases;

-- name: BatchLinkKnowledgeBaseToFiles :exec
INSERT INTO knowledge_base_files (knowledge_base_id, file_id)
VALUES (?, ?);

-- Knowledge Base Files

-- name: LinkKnowledgeBaseToFile :exec
INSERT INTO knowledge_base_files (knowledge_base_id, file_id)
VALUES (?, ?);

-- name: ListKnowledgeBaseFiles :many
SELECT * FROM knowledge_base_files
WHERE knowledge_base_id = ?;

-- name: UnlinkKnowledgeBaseFromFile :exec
DELETE FROM knowledge_base_files
WHERE knowledge_base_id = ? AND file_id = ?;

-- name: BatchUnlinkKnowledgeBaseFromFiles :exec
DELETE FROM knowledge_base_files
WHERE knowledge_base_id = ? AND file_id = ?;

-- name: GetKnowledgeBaseFiles :many
SELECT f.* FROM files f
INNER JOIN knowledge_base_files kbf ON f.id = kbf.file_id
WHERE kbf.knowledge_base_id = ?;

-- Files to Sessions

-- name: LinkFileToSession :exec
INSERT INTO files_to_sessions (file_id, session_id)
VALUES (?, ?);

-- name: UnlinkFileFromSession :exec
DELETE FROM files_to_sessions
WHERE file_id = ? AND session_id = ?;

-- name: GetSessionFiles :many
SELECT f.* FROM files f
INNER JOIN files_to_sessions fts ON f.id = fts.file_id
WHERE fts.session_id = ?;

-- Complex file queries

-- name: CountFilesByHash :one
SELECT COUNT(*) as count
FROM files
WHERE file_hash = ?;

-- name: GetFilesByHash :many
SELECT * FROM files
WHERE file_hash = ?;

-- name: GetFilesByIds :many
SELECT * FROM files;

-- name: GetFilesByNames :many
SELECT * FROM files
ORDER BY created_at DESC;

-- name: CountFilesUsage :one
SELECT COALESCE(SUM(size), 0) as total_size
FROM files;

-- name: DeleteAllFiles :exec
DELETE FROM files;

-- name: DeleteGlobalFile :exec
DELETE FROM global_files WHERE hash_id = ?;

-- name: GetFileChunkIds :many
SELECT chunk_id FROM file_chunks
WHERE file_id = ?;

-- File query with filters
-- name: QueryFiles :many
SELECT 
    f.id,
    f.name,
    f.file_type,
    f.size,
    f.url,
    f.chunk_count,
    f.chunking_status,
    f.embedding_status,
    f.created_at,
    f.updated_at
FROM files f
ORDER BY f.created_at DESC;

-- name: QueryFilesByKnowledgeBase :many
SELECT 
    f.id,
    f.name,
    f.file_type,
    f.size,
    f.url,
    f.chunk_count,
    f.chunking_status,
    f.embedding_status,
    f.created_at,
    f.updated_at
FROM files f
INNER JOIN knowledge_base_files kbf ON f.id = kbf.file_id
WHERE kbf.knowledge_base_id = ?
ORDER BY f.created_at DESC;
