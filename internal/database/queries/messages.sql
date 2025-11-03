-- name: GetMessage :one
SELECT * FROM messages WHERE id = ? AND user_id = ?;

-- name: ListMessages :many
SELECT * FROM messages
WHERE user_id = ? AND session_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListMessagesByTopic :many
SELECT * FROM messages
WHERE user_id = ? AND topic_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListMessagesByThread :many
SELECT * FROM messages
WHERE user_id = ? AND thread_id = ?
ORDER BY created_at ASC;

-- name: CreateMessage :one
INSERT INTO messages (
    id, role, content, reasoning, search, metadata, model, provider,
    favorite, error, tools, trace_id, observation_id, client_id,
    user_id, session_id, topic_id, thread_id, parent_id, quota_id,
    agent_id, group_id, target_id, message_group_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateMessage :one
UPDATE messages
SET content = ?,
    reasoning = ?,
    metadata = ?,
    favorite = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = ? AND user_id = ?;

-- name: DeleteMessagesBySession :exec
DELETE FROM messages WHERE session_id = ? AND user_id = ?;

-- name: DeleteMessagesByTopic :exec
DELETE FROM messages WHERE topic_id = ? AND user_id = ?;

-- name: ToggleMessageFavorite :exec
UPDATE messages SET favorite = ?, updated_at = ?
WHERE id = ? AND user_id = ?;

-- Message Plugins

-- name: GetMessagePlugin :one
SELECT * FROM message_plugins WHERE id = ? AND user_id = ?;

-- name: CreateMessagePlugin :one
INSERT INTO message_plugins (
    id, tool_call_id, type, api_name, arguments, identifier,
    state, error, client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateMessagePlugin :one
UPDATE message_plugins
SET state = ?, error = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- Message TTS

-- name: GetMessageTTS :one
SELECT * FROM message_tts WHERE id = ? AND user_id = ?;

-- name: CreateMessageTTS :one
INSERT INTO message_tts (
    id, content_md5, file_id, voice, client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- Message Translates

-- name: GetMessageTranslate :one
SELECT * FROM message_translates WHERE id = ? AND user_id = ?;

-- name: CreateMessageTranslate :one
INSERT INTO message_translates (
    id, content, "from", "to", client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- Message Queries (RAG)

-- name: GetMessageQuery :one
SELECT * FROM message_queries WHERE id = ? AND user_id = ?;

-- name: CreateMessageQuery :one
INSERT INTO message_queries (
    id, message_id, rewrite_query, user_query, client_id, user_id, embeddings_id
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListMessageQueriesByMessage :many
SELECT * FROM message_queries
WHERE message_id = ? AND user_id = ?;

-- Message Files

-- name: LinkMessageToFile :exec
INSERT INTO messages_files (file_id, message_id, user_id)
VALUES (?, ?, ?);

-- name: UnlinkMessageFromFile :exec
DELETE FROM messages_files
WHERE file_id = ? AND message_id = ? AND user_id = ?;

-- name: GetMessageFiles :many
SELECT f.* FROM files f
INNER JOIN messages_files mf ON f.id = mf.file_id
WHERE mf.message_id = ? AND mf.user_id = ?;

-- Message Chunks

-- name: LinkMessageToChunk :exec
INSERT INTO message_chunks (message_id, chunk_id, user_id)
VALUES (?, ?, ?);

-- name: GetMessageChunks :many
SELECT c.* FROM chunks c
INNER JOIN message_chunks mc ON c.id = mc.chunk_id
WHERE mc.message_id = ? AND mc.user_id = ?;

