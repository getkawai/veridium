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

