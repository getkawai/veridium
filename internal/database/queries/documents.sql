-- name: GetDocument :one
SELECT * FROM documents WHERE id = ?;

-- name: GetDocumentByFileID :one
SELECT * FROM documents WHERE file_id = ? LIMIT 1;

-- name: ListDocuments :many
SELECT * FROM documents
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateDocument :one
INSERT INTO documents (
    title, content, file_type, filename, total_char_count,
    total_line_count, metadata, pages, source_type, source,
    file_id, editor_data
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateDocument :one
UPDATE documents
SET title = ?,
    content = ?,
    metadata = ?,
    editor_data = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteDocument :exec
DELETE FROM documents WHERE id = ?;

-- name: DeleteAllDocuments :exec
DELETE FROM documents;

-- Document Chunks

-- name: LinkDocumentToChunk :exec
INSERT INTO document_chunks (document_id, chunk_id, page_index)
VALUES (?, ?, ?);

-- name: UnlinkDocumentFromChunk :exec
DELETE FROM document_chunks
WHERE document_id = ? AND chunk_id = ?;

-- name: GetDocumentChunks :many
SELECT c.* FROM chunks c
INNER JOIN document_chunks dc ON c.id = dc.chunk_id
WHERE dc.document_id = ?
ORDER BY dc.page_index ASC, c.chunk_index ASC;

-- Topic Documents

-- name: LinkTopicToDocument :exec
INSERT INTO topic_documents (document_id, topic_id, created_at)
VALUES (?, ?, ?);

-- name: UnlinkTopicFromDocument :exec
DELETE FROM topic_documents
WHERE document_id = ? AND topic_id = ?;

-- name: GetTopicDocuments :many
SELECT d.* FROM documents d
INNER JOIN topic_documents td ON d.id = td.document_id
WHERE td.topic_id = ?;
