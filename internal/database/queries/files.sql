-- name: GetFile :one
SELECT * FROM files WHERE id = ? AND user_id = ?;

-- name: ListFiles :many
SELECT * FROM files
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateFile :one
INSERT INTO files (
    id, user_id, file_type, file_hash, name, size, url, source,
    client_id, metadata, chunk_task_id, embedding_task_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateFile :one
UPDATE files
SET name = ?,
    metadata = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files WHERE id = ? AND user_id = ?;

-- Global Files

-- name: GetGlobalFile :one
SELECT * FROM global_files WHERE hash_id = ?;

-- name: GetGlobalFileByHash :one
SELECT * FROM global_files WHERE hash_id = ?;

-- name: CreateGlobalFile :one
INSERT INTO global_files (
    hash_id, file_type, size, url, metadata, creator, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;


-- Knowledge Bases

-- name: GetKnowledgeBase :one
SELECT * FROM knowledge_bases WHERE id = ? AND user_id = ?;

-- name: ListKnowledgeBases :many
SELECT * FROM knowledge_bases
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CreateKnowledgeBase :one
INSERT INTO knowledge_bases (
    id, name, description, avatar, type, user_id, client_id,
    is_public, settings, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateKnowledgeBase :one
UPDATE knowledge_bases
SET name = ?,
    description = ?,
    avatar = ?,
    settings = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteKnowledgeBase :exec
DELETE FROM knowledge_bases WHERE id = ? AND user_id = ?;

-- name: DeleteAllKnowledgeBases :exec
DELETE FROM knowledge_bases WHERE user_id = ?;

-- name: BatchLinkKnowledgeBaseToFiles :exec
INSERT INTO knowledge_base_files (knowledge_base_id, file_id, user_id, created_at)
VALUES (?, ?, ?, ?);

-- Knowledge Base Files

-- name: LinkKnowledgeBaseToFile :exec
INSERT INTO knowledge_base_files (knowledge_base_id, file_id, user_id, created_at)
VALUES (?, ?, ?, ?);

-- name: ListKnowledgeBaseFiles :many
SELECT * FROM knowledge_base_files
WHERE knowledge_base_id = ? AND user_id = ?;

-- name: UnlinkKnowledgeBaseFromFile :exec
DELETE FROM knowledge_base_files
WHERE knowledge_base_id = ? AND file_id = ? AND user_id = ?;

-- name: BatchUnlinkKnowledgeBaseFromFiles :exec
DELETE FROM knowledge_base_files
WHERE knowledge_base_id = ? AND file_id = ?;

-- name: GetKnowledgeBaseFiles :many
SELECT f.* FROM files f
INNER JOIN knowledge_base_files kbf ON f.id = kbf.file_id
WHERE kbf.knowledge_base_id = ? AND kbf.user_id = ?;

-- Files to Sessions

-- name: LinkFileToSession :exec
INSERT INTO files_to_sessions (file_id, session_id, user_id)
VALUES (?, ?, ?);

-- name: UnlinkFileFromSession :exec
DELETE FROM files_to_sessions
WHERE file_id = ? AND session_id = ? AND user_id = ?;

-- name: GetSessionFiles :many
SELECT f.* FROM files f
INNER JOIN files_to_sessions fts ON f.id = fts.file_id
WHERE fts.session_id = ? AND fts.user_id = ?;

-- Complex file queries

-- name: CountFilesByHash :one
SELECT COUNT(*) as count
FROM files
WHERE file_hash = ?;

-- name: GetFilesByHash :many
SELECT * FROM files
WHERE file_hash = ? AND user_id = ?;

-- name: GetFilesByIds :many
SELECT * FROM files
WHERE user_id = ?;

-- name: GetFilesByNames :many
SELECT * FROM files
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CountFilesUsage :one
SELECT COALESCE(SUM(size), 0) as total_size
FROM files
WHERE user_id = ?;

-- name: DeleteAllFiles :exec
DELETE FROM files WHERE user_id = ?;

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
    f.created_at,
    f.updated_at,
    f.chunk_task_id,
    f.embedding_task_id
FROM files f
WHERE f.user_id = ?
ORDER BY f.created_at DESC;

-- name: QueryFilesByKnowledgeBase :many
SELECT 
    f.id,
    f.name,
    f.file_type,
    f.size,
    f.url,
    f.created_at,
    f.updated_at,
    f.chunk_task_id,
    f.embedding_task_id
FROM files f
INNER JOIN knowledge_base_files kbf ON f.id = kbf.file_id
WHERE kbf.knowledge_base_id = ? AND f.user_id = ?
ORDER BY f.created_at DESC;

