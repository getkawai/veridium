-- name: GetDocument :one
SELECT * FROM documents WHERE id = ? AND user_id = ?;

-- name: ListDocuments :many
SELECT * FROM documents
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateDocument :one
INSERT INTO documents (
    id, title, content, file_type, filename, total_char_count,
    total_line_count, metadata, pages, source_type, source,
    file_id, user_id, client_id, editor_data, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateDocument :one
UPDATE documents
SET title = ?,
    content = ?,
    metadata = ?,
    editor_data = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteDocument :exec
DELETE FROM documents WHERE id = ? AND user_id = ?;

-- name: DeleteAllDocuments :exec
DELETE FROM documents WHERE user_id = ?;

-- Document Chunks

-- name: LinkDocumentToChunk :exec
INSERT INTO document_chunks (document_id, chunk_id, page_index, user_id)
VALUES (?, ?, ?, ?);

-- name: UnlinkDocumentFromChunk :exec
DELETE FROM document_chunks
WHERE document_id = ? AND chunk_id = ? AND user_id = ?;

-- name: GetDocumentChunks :many
SELECT c.* FROM chunks c
INNER JOIN document_chunks dc ON c.id = dc.chunk_id
WHERE dc.document_id = ? AND dc.user_id = ?
ORDER BY dc.page_index ASC, c.chunk_index ASC;

-- Topic Documents

-- name: LinkTopicToDocument :exec
INSERT INTO topic_documents (document_id, topic_id, user_id, created_at)
VALUES (?, ?, ?, ?);

-- name: UnlinkTopicFromDocument :exec
DELETE FROM topic_documents
WHERE document_id = ? AND topic_id = ? AND user_id = ?;

-- name: GetTopicDocuments :many
SELECT d.* FROM documents d
INNER JOIN topic_documents td ON d.id = td.document_id
WHERE td.topic_id = ? AND td.user_id = ?;

