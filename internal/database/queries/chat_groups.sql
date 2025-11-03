-- name: GetChatGroup :one
SELECT * FROM chat_groups WHERE id = ? AND user_id = ?;

-- name: ListChatGroups :many
SELECT * FROM chat_groups
WHERE user_id = ?
ORDER BY updated_at DESC;

-- name: CreateChatGroup :one
INSERT INTO chat_groups (
    id, title, description, config, client_id, user_id, pinned,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateChatGroup :one
UPDATE chat_groups
SET title = ?,
    description = ?,
    config = ?,
    pinned = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteChatGroup :exec
DELETE FROM chat_groups WHERE id = ? AND user_id = ?;

-- Chat Group Agents

-- name: LinkChatGroupToAgent :exec
INSERT INTO chat_groups_agents (
    chat_group_id, agent_id, user_id, enabled, sort_order, role,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UnlinkChatGroupFromAgent :exec
DELETE FROM chat_groups_agents
WHERE chat_group_id = ? AND agent_id = ? AND user_id = ?;

-- name: GetChatGroupAgents :many
SELECT a.* FROM agents a
INNER JOIN chat_groups_agents cga ON a.id = cga.agent_id
WHERE cga.chat_group_id = ? AND cga.user_id = ?
ORDER BY cga.sort_order ASC;

-- name: UpdateChatGroupAgentOrder :exec
UPDATE chat_groups_agents
SET sort_order = ?, updated_at = ?
WHERE chat_group_id = ? AND agent_id = ? AND user_id = ?;

-- Message Groups

-- name: GetMessageGroup :one
SELECT * FROM message_groups WHERE id = ? AND user_id = ?;

-- name: ListMessageGroupsByTopic :many
SELECT * FROM message_groups
WHERE topic_id = ? AND user_id = ?
ORDER BY created_at ASC;

-- name: CreateMessageGroup :one
INSERT INTO message_groups (
    id, title, description, topic_id, user_id, parent_group_id,
    client_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteMessageGroup :exec
DELETE FROM message_groups WHERE id = ? AND user_id = ?;

-- name: DeleteAllChatGroups :exec
DELETE FROM chat_groups WHERE user_id = ?;

-- Complex queries with JOINs

-- name: ListChatGroupsWithAgents :many
SELECT 
    cg.id as group_id,
    cg.title as group_title,
    cg.description as group_description,
    cg.config as group_config,
    cg.pinned as group_pinned,
    cg.created_at as group_created_at,
    cg.updated_at as group_updated_at,
    a.id as agent_id,
    a.title as agent_title,
    a.description as agent_description,
    a.avatar as agent_avatar,
    a.background_color as agent_bg_color,
    a.chat_config as agent_chat_config,
    a.params as agent_params,
    a.system_role as agent_system_role,
    a.tts as agent_tts,
    a.model as agent_model,
    a.provider as agent_provider,
    a.created_at as agent_created_at,
    a.updated_at as agent_updated_at,
    cga.sort_order as agent_sort_order,
    cga.enabled as agent_enabled,
    cga.role as agent_role
FROM chat_groups cg
LEFT JOIN chat_groups_agents cga ON cg.id = cga.chat_group_id
LEFT JOIN agents a ON cga.agent_id = a.id
WHERE cg.user_id = ?
ORDER BY cg.updated_at DESC, cga.sort_order ASC;

-- name: GetChatGroupWithAgents :many
SELECT 
    cg.id as group_id,
    cg.title as group_title,
    cg.description as group_description,
    cg.config as group_config,
    cg.pinned as group_pinned,
    cg.created_at as group_created_at,
    cg.updated_at as group_updated_at,
    a.id as agent_id,
    a.title as agent_title,
    a.description as agent_description,
    a.avatar as agent_avatar,
    a.background_color as agent_bg_color,
    a.chat_config as agent_chat_config,
    a.params as agent_params,
    a.system_role as agent_system_role,
    a.tts as agent_tts,
    a.model as agent_model,
    a.provider as agent_provider,
    a.created_at as agent_created_at,
    a.updated_at as agent_updated_at,
    cga.sort_order as agent_sort_order,
    cga.enabled as agent_enabled,
    cga.role as agent_role
FROM chat_groups cg
LEFT JOIN chat_groups_agents cga ON cg.id = cga.chat_group_id
LEFT JOIN agents a ON cga.agent_id = a.id
WHERE cg.id = ? AND cg.user_id = ?
ORDER BY cga.sort_order ASC;

-- name: GetChatGroupAgentLinks :many
SELECT * FROM chat_groups_agents
WHERE chat_group_id = ? AND user_id = ?
ORDER BY sort_order ASC;

-- name: GetEnabledChatGroupAgentLinks :many
SELECT * FROM chat_groups_agents
WHERE chat_group_id = ? AND user_id = ? AND enabled = 1
ORDER BY sort_order ASC;

-- name: UpdateChatGroupAgentLink :one
UPDATE chat_groups_agents
SET sort_order = ?,
    role = ?,
    enabled = ?,
    updated_at = ?
WHERE chat_group_id = ? AND agent_id = ? AND user_id = ?
RETURNING *;

-- name: BatchLinkChatGroupToAgents :exec
INSERT INTO chat_groups_agents (
    chat_group_id, agent_id, user_id, enabled, sort_order, role,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

