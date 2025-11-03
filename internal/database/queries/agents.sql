-- name: GetAgent :one
SELECT * FROM agents WHERE id = ? AND user_id = ?;

-- name: GetAgentBySlug :one
SELECT * FROM agents WHERE slug = ? AND user_id = ?;

-- name: ListAgents :many
SELECT * FROM agents
WHERE user_id = ?
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: SearchAgents :many
SELECT * FROM agents
WHERE user_id = ? AND (title LIKE ? OR description LIKE ?)
ORDER BY updated_at DESC
LIMIT ?;

-- name: CreateAgent :one
INSERT INTO agents (
    id, slug, title, description, tags, avatar, background_color,
    plugins, client_id, user_id, chat_config, few_shots, model,
    params, provider, system_role, tts, virtual, opening_message,
    opening_questions, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateAgent :one
UPDATE agents
SET title = ?,
    description = ?,
    tags = ?,
    avatar = ?,
    background_color = ?,
    plugins = ?,
    chat_config = ?,
    few_shots = ?,
    model = ?,
    params = ?,
    provider = ?,
    system_role = ?,
    tts = ?,
    opening_message = ?,
    opening_questions = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agents WHERE id = ? AND user_id = ?;

-- Agent to Session relationships

-- name: LinkAgentToSession :exec
INSERT INTO agents_to_sessions (agent_id, session_id, user_id)
VALUES (?, ?, ?);

-- name: UnlinkAgentFromSession :exec
DELETE FROM agents_to_sessions
WHERE agent_id = ? AND session_id = ? AND user_id = ?;

-- name: GetSessionAgents :many
SELECT a.* FROM agents a
INNER JOIN agents_to_sessions ats ON a.id = ats.agent_id
WHERE ats.session_id = ? AND ats.user_id = ?;

-- name: GetAgentSessions :many
SELECT s.* FROM sessions s
INNER JOIN agents_to_sessions ats ON s.id = ats.session_id
WHERE ats.agent_id = ? AND ats.user_id = ?;

-- name: GetOrphanedAgents :many
SELECT a.* FROM agents a
LEFT JOIN agents_to_sessions ats ON a.id = ats.agent_id
WHERE a.user_id = ? AND ats.agent_id IS NULL;

-- Agent Files

-- name: LinkAgentToFile :exec
INSERT INTO agents_files (file_id, agent_id, enabled, user_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UnlinkAgentFromFile :exec
DELETE FROM agents_files
WHERE file_id = ? AND agent_id = ? AND user_id = ?;

-- name: GetAgentFiles :many
SELECT f.* FROM files f
INNER JOIN agents_files af ON f.id = af.file_id
WHERE af.agent_id = ? AND af.user_id = ?;

-- Agent Knowledge Bases

-- name: LinkAgentToKnowledgeBase :exec
INSERT INTO agents_knowledge_bases (agent_id, knowledge_base_id, user_id, enabled, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UnlinkAgentFromKnowledgeBase :exec
DELETE FROM agents_knowledge_bases
WHERE agent_id = ? AND knowledge_base_id = ? AND user_id = ?;

-- name: GetAgentKnowledgeBases :many
SELECT kb.* FROM knowledge_bases kb
INNER JOIN agents_knowledge_bases akb ON kb.id = akb.knowledge_base_id
WHERE akb.agent_id = ? AND akb.user_id = ?;

