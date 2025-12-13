-- name: GetAgent :one
SELECT * FROM agents WHERE id = ?;



-- name: ListAgents :many
SELECT * FROM agents
ORDER BY updated_at DESC
LIMIT ? OFFSET ?;

-- name: SearchAgents :many
SELECT * FROM agents
WHERE (title LIKE ? OR description LIKE ?)
ORDER BY updated_at DESC
LIMIT ?;

-- name: CreateAgent :one
INSERT INTO agents (
    id, title, description, tags, avatar, background_color,
    plugins, chat_config, few_shots, model,
    params, provider, system_role, tts, virtual, opening_message,
    opening_questions, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
WHERE id = ?
RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agents WHERE id = ?;

-- Agent to Session relationships

-- name: LinkAgentToSession :exec
INSERT INTO agents_to_sessions (agent_id, session_id)
VALUES (?, ?);

-- name: UnlinkAgentFromSession :exec
DELETE FROM agents_to_sessions
WHERE agent_id = ? AND session_id = ?;

-- name: GetSessionAgents :many
SELECT a.* FROM agents a
INNER JOIN agents_to_sessions ats ON a.id = ats.agent_id
WHERE ats.session_id = ?;

-- name: GetAgentSessions :many
SELECT s.* FROM sessions s
INNER JOIN agents_to_sessions ats ON s.id = ats.session_id
WHERE ats.agent_id = ?;

-- name: GetOrphanedAgents :many
SELECT a.* FROM agents a
LEFT JOIN agents_to_sessions ats ON a.id = ats.agent_id
WHERE ats.agent_id IS NULL;

-- Agent Files

-- name: LinkAgentToFile :exec
INSERT INTO agents_files (file_id, agent_id, enabled, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: UnlinkAgentFromFile :exec
DELETE FROM agents_files
WHERE file_id = ? AND agent_id = ?;

-- name: GetAgentFiles :many
SELECT f.* FROM files f
INNER JOIN agents_files af ON f.id = af.file_id
WHERE af.agent_id = ?;

-- Agent Knowledge Bases

-- name: LinkAgentToKnowledgeBase :exec
INSERT INTO agents_knowledge_bases (agent_id, knowledge_base_id, enabled, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: UnlinkAgentFromKnowledgeBase :exec
DELETE FROM agents_knowledge_bases
WHERE agent_id = ? AND knowledge_base_id = ?;

-- name: GetAgentKnowledgeBases :many
SELECT kb.*, akb.enabled
FROM knowledge_bases kb
INNER JOIN agents_knowledge_bases akb ON kb.id = akb.knowledge_base_id
WHERE akb.agent_id = ?
ORDER BY akb.created_at DESC;

-- name: GetAgentFilesWithEnabled :many
SELECT f.*, af.enabled
FROM files f
INNER JOIN agents_files af ON f.id = af.file_id
WHERE af.agent_id = ?
ORDER BY af.created_at DESC;

-- name: ToggleAgentKnowledgeBase :exec
UPDATE agents_knowledge_bases
SET enabled = ?
WHERE agent_id = ? AND knowledge_base_id = ?;

-- name: ToggleAgentFile :exec
UPDATE agents_files
SET enabled = ?
WHERE agent_id = ? AND file_id = ?;

-- name: GetAgentFileIds :many
SELECT file_id FROM agents_files
WHERE agent_id = ? AND file_id IN (sqlc.slice('fileIds'));

-- name: BatchLinkAgentToFiles :exec
INSERT INTO agents_files (agent_id, file_id, enabled, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetAgentBySessionId :one
SELECT a.* FROM agents a
INNER JOIN agents_to_sessions ats ON a.id = ats.agent_id
WHERE ats.session_id = ?
LIMIT 1;

-- name: DuplicateAgentForSession :one
-- Duplicate an agent for a new session
-- Parameters: new_agent_id, new_session_id, source_session_id, user_id, created_at, updated_at
INSERT INTO agents (
    id, title, description, tags, avatar, background_color,
    plugins, chat_config, few_shots, model,
    params, provider, system_role, tts, virtual, opening_message,
    opening_questions, created_at, updated_at
)
SELECT 
    ? as id,           -- new_agent_id
    title,
    description,
    tags,
    avatar,
    background_color,
    plugins,
    chat_config,
    few_shots,
    model,
    params,
    provider,
    system_role,
    tts,
    virtual,
    opening_message,
    opening_questions,
    ? as created_at,   -- new created_at
    ? as updated_at    -- new updated_at
FROM agents a
INNER JOIN agents_to_sessions ats ON a.id = ats.agent_id
WHERE ats.session_id = ?
LIMIT 1
RETURNING *;

-- name: LinkDuplicatedAgentToSession :exec
-- Link a duplicated agent to a new session
INSERT INTO agents_to_sessions (agent_id, session_id)
VALUES (?, ?);
