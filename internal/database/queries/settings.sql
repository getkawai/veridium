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

-- name: UpdateUserSettings :one
UPDATE user_settings
SET tts = COALESCE(?, tts),
    hotkey = COALESCE(?, hotkey),
    key_vaults = COALESCE(?, key_vaults),
    general = COALESCE(?, general),
    language_model = COALESCE(?, language_model),
    system_agent = COALESCE(?, system_agent),
    default_agent = COALESCE(?, default_agent),
    tool = COALESCE(?, tool),
    image = COALESCE(?, image)
WHERE id = ?
RETURNING *;
