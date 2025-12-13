-- Plugins

-- name: GetPlugin :one
SELECT * FROM user_installed_plugins
WHERE identifier = ?;

-- name: ListPlugins :many
SELECT * FROM user_installed_plugins
ORDER BY created_at DESC;

-- name: CreatePlugin :one
INSERT INTO user_installed_plugins (
    identifier, type, manifest, custom_params, settings, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpsertPlugin :one
INSERT INTO user_installed_plugins (
    identifier, type, manifest, custom_params, settings, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(identifier) DO UPDATE SET
    type = excluded.type,
    manifest = excluded.manifest,
    custom_params = excluded.custom_params,
    settings = excluded.settings,
    updated_at = excluded.updated_at
RETURNING *;

-- name: UpdatePlugin :exec
UPDATE user_installed_plugins
SET type = ?,
    manifest = ?,
    custom_params = ?,
    settings = ?,
    updated_at = ?
WHERE identifier = ?;

-- name: DeletePlugin :exec
DELETE FROM user_installed_plugins
WHERE identifier = ?;

-- name: DeleteAllPlugins :exec
DELETE FROM user_installed_plugins;
