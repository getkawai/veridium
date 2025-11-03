-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserWithSettings :one
SELECT 
    u.id,
    u.username,
    u.email,
    u.avatar,
    u.first_name,
    u.last_name,
    u.is_onboarded,
    u.preference,
    u.created_at,
    u.updated_at,
    us.tts as settings_tts,
    us.hotkey as settings_hotkey,
    us.key_vaults as settings_key_vaults,
    us.general as settings_general,
    us.language_model as settings_language_model,
    us.system_agent as settings_system_agent,
    us.default_agent as settings_default_agent,
    us.tool as settings_tool,
    us.image as settings_image
FROM users u
LEFT JOIN user_settings us ON u.id = us.id
WHERE u.id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: CreateUser :one
INSERT INTO users (
    id, username, email, avatar, phone, first_name, last_name,
    is_onboarded, clerk_created_at, email_verified_at, preference,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: EnsureUserExists :exec
INSERT INTO users (id, created_at, updated_at)
VALUES (?, ?, ?)
ON CONFLICT(id) DO NOTHING;

-- name: UpdateUser :one
UPDATE users
SET username = ?,
    email = ?,
    avatar = ?,
    phone = ?,
    first_name = ?,
    last_name = ?,
    preference = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: UpdateUserOnboarding :exec
UPDATE users
SET is_onboarded = ?, updated_at = ?
WHERE id = ?;

-- name: UpdateUserPreference :one
UPDATE users
SET preference = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: DeleteUserSettings :exec
DELETE FROM user_settings WHERE id = ?;

-- User Settings

-- name: GetUserSettings :one
SELECT * FROM user_settings WHERE id = ?;

-- name: UpsertUserSettings :one
INSERT INTO user_settings (
    id, tts, hotkey, key_vaults, general, language_model,
    system_agent, default_agent, tool, image
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    tts = excluded.tts,
    hotkey = excluded.hotkey,
    key_vaults = excluded.key_vaults,
    general = excluded.general,
    language_model = excluded.language_model,
    system_agent = excluded.system_agent,
    default_agent = excluded.default_agent,
    tool = excluded.tool,
    image = excluded.image
RETURNING *;

-- name: UpdateUserSettingsTTS :exec
UPDATE user_settings SET tts = ? WHERE id = ?;

-- name: UpdateUserSettingsHotkey :exec
UPDATE user_settings SET hotkey = ? WHERE id = ?;

-- name: UpdateUserSettingsGeneral :exec
UPDATE user_settings SET general = ? WHERE id = ?;

-- User Installed Plugins

-- name: ListUserPlugins :many
SELECT * FROM user_installed_plugins
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: GetUserPlugin :one
SELECT * FROM user_installed_plugins
WHERE user_id = ? AND identifier = ?;

-- name: InstallUserPlugin :one
INSERT INTO user_installed_plugins (
    user_id, identifier, type, manifest, settings, custom_params,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUserPlugin :one
UPDATE user_installed_plugins
SET settings = ?, custom_params = ?, updated_at = ?
WHERE user_id = ? AND identifier = ?
RETURNING *;

-- name: UninstallUserPlugin :exec
DELETE FROM user_installed_plugins
WHERE user_id = ? AND identifier = ?;

