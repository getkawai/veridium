-- name: GetMessage :one
SELECT * FROM messages WHERE id = ? AND user_id = ?;

-- name: ListMessages :many
SELECT * FROM messages
WHERE user_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListMessagesBySession :many
SELECT * FROM messages
WHERE user_id = ? AND session_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListMessagesByTopic :many
SELECT * FROM messages
WHERE user_id = ? AND topic_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListMessagesByGroup :many
SELECT * FROM messages
WHERE user_id = ? AND group_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: CountMessages :one
SELECT COUNT(*) FROM messages WHERE user_id = ?;

-- name: CountMessagesByDateRange :one
SELECT COUNT(*) FROM messages
WHERE user_id = ?
  AND created_at >= ?
  AND created_at <= ?;

-- name: CountMessageWords :one
SELECT SUM(LENGTH(content)) as total_length
FROM messages
WHERE user_id = ?;

-- name: CountMessageWordsByDateRange :one
SELECT SUM(LENGTH(content)) as total_length
FROM messages
WHERE user_id = ?
  AND created_at >= ?
  AND created_at <= ?;

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

-- name: BatchDeleteMessages :exec
DELETE FROM messages
WHERE user_id = ? AND id IN (sqlc.slice('ids'));

-- name: DeleteAllMessages :exec
DELETE FROM messages WHERE user_id = ?;

-- name: DeleteMessagesByGroup :exec
DELETE FROM messages WHERE group_id = ? AND user_id = ?;

-- Batch queries - Note: These will be wrapped with JSON parsing in Go
-- For now, we'll create a simple version and handle batching in Go layer

-- name: GetMessageByToolCallId :one
SELECT mp.id 
FROM message_plugins mp
WHERE mp.tool_call_id = ? AND mp.user_id = ?;

-- name: GetDocumentByFileId :one
SELECT d.file_id, d.content
FROM documents d
WHERE d.file_id = ? AND d.user_id = ?;

-- name: GetMessagesWithRelations :many
SELECT 
    m.id,
    m.role,
    m.content,
    m.reasoning,
    m.search,
    m.metadata,
    m.error,
    m.model,
    m.provider,
    m.created_at,
    m.updated_at,
    m.topic_id,
    m.parent_id,
    m.thread_id,
    m.group_id,
    m.agent_id,
    m.target_id,
    m.tools,
    m.favorite,
    mp.tool_call_id,
    mp.api_name as plugin_api_name,
    mp.arguments as plugin_arguments,
    mp.identifier as plugin_identifier,
    mp.type as plugin_type,
    mp.state as plugin_state,
    mp.error as plugin_error,
    mt.content as translate_content,
    mt.from as translate_from,
    mt.to as translate_to,
    mts.id as tts_id,
    mts.content_md5 as tts_content_md5,
    mts.file_id as tts_file_id,
    mts.voice as tts_voice
FROM messages m
LEFT JOIN message_plugins mp ON m.id = mp.id
LEFT JOIN message_translates mt ON m.id = mt.id
LEFT JOIN message_tts mts ON m.id = mts.id
WHERE m.user_id = ?
ORDER BY m.created_at ASC
LIMIT ? OFFSET ?;

-- name: GetMessagesWithRelationsBySession :many
SELECT 
    m.id,
    m.role,
    m.content,
    m.reasoning,
    m.search,
    m.metadata,
    m.error,
    m.model,
    m.provider,
    m.created_at,
    m.updated_at,
    m.topic_id,
    m.parent_id,
    m.thread_id,
    m.group_id,
    m.agent_id,
    m.target_id,
    m.tools,
    m.favorite,
    mp.tool_call_id,
    mp.api_name as plugin_api_name,
    mp.arguments as plugin_arguments,
    mp.identifier as plugin_identifier,
    mp.type as plugin_type,
    mp.state as plugin_state,
    mp.error as plugin_error,
    mt.content as translate_content,
    mt.from as translate_from,
    mt.to as translate_to,
    mts.id as tts_id,
    mts.content_md5 as tts_content_md5,
    mts.file_id as tts_file_id,
    mts.voice as tts_voice
FROM messages m
LEFT JOIN message_plugins mp ON m.id = mp.id
LEFT JOIN message_translates mt ON m.id = mt.id
LEFT JOIN message_tts mts ON m.id = mts.id
WHERE m.user_id = ? AND m.session_id = ?
ORDER BY m.created_at ASC
LIMIT ? OFFSET ?;

-- name: GetMessageHeatmaps :many
SELECT 
    DATE(created_at / 1000, 'unixepoch') as date,
    COUNT(*) as count
FROM messages
WHERE user_id = ? 
    AND created_at >= ? 
    AND created_at <= ?
GROUP BY date
ORDER BY date DESC;

-- name: ToggleMessageFavorite :exec
UPDATE messages SET favorite = ?, updated_at = ?
WHERE id = ? AND user_id = ?;

-- name: RankModels :many
SELECT model as id, COUNT(*) as count
FROM messages
WHERE user_id = ? AND model IS NOT NULL AND model != ''
GROUP BY model
HAVING COUNT(*) > 0
ORDER BY count DESC, model ASC
LIMIT ?;

-- name: SearchMessagesByKeyword :many
SELECT * FROM messages
WHERE user_id = ? AND content LIKE ?
ORDER BY created_at DESC
LIMIT ?;

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

-- name: UpsertMessageTTS :one
INSERT INTO message_tts (
    id, content_md5, file_id, voice, client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    content_md5 = excluded.content_md5,
    file_id = excluded.file_id,
    voice = excluded.voice
RETURNING *;

-- name: DeleteMessageTTS :exec
DELETE FROM message_tts WHERE id = ? AND user_id = ?;

-- Message Translates

-- name: GetMessageTranslate :one
SELECT * FROM message_translates WHERE id = ? AND user_id = ?;

-- name: CreateMessageTranslate :one
INSERT INTO message_translates (
    id, content, "from", "to", client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpsertMessageTranslate :one
INSERT INTO message_translates (
    id, content, "from", "to", client_id, user_id
) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    content = excluded.content,
    "from" = excluded."from",
    "to" = excluded."to"
RETURNING *;

-- name: DeleteMessageTranslate :exec
DELETE FROM message_translates WHERE id = ? AND user_id = ?;

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

-- name: DeleteMessageQuery :exec
DELETE FROM message_queries WHERE id = ? AND user_id = ?;

-- Message Query Chunks

-- name: LinkMessageQueryToChunk :exec
INSERT INTO message_query_chunks (message_id, query_id, chunk_id, similarity, user_id)
VALUES (?, ?, ?, ?, ?);

-- name: GetMessageQueryChunks :many
SELECT 
    mqc.message_id,
    mqc.similarity,
    c.id,
    c.text,
    f.id as file_id,
    f.name as filename,
    f.file_type,
    f.url as file_url
FROM message_query_chunks mqc
LEFT JOIN chunks c ON mqc.chunk_id = c.id
LEFT JOIN file_chunks fc ON c.id = fc.chunk_id
LEFT JOIN files f ON fc.file_id = f.id
WHERE mqc.message_id IN (sqlc.slice('messageIds')) AND mqc.user_id = ?;

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

