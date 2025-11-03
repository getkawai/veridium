-- AI Providers

-- name: GetAIProvider :one
SELECT * FROM ai_providers WHERE id = ? AND user_id = ?;

-- name: ListAIProviders :many
SELECT * FROM ai_providers
WHERE user_id = ?
ORDER BY sort ASC;

-- name: ListEnabledAIProviders :many
SELECT * FROM ai_providers
WHERE user_id = ? AND enabled = 1
ORDER BY sort ASC;

-- name: CreateAIProvider :one
INSERT INTO ai_providers (
    id, name, user_id, sort, enabled, fetch_on_client, check_model,
    logo, description, key_vaults, source, settings, config,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateAIProvider :one
UPDATE ai_providers
SET name = ?,
    sort = ?,
    enabled = ?,
    fetch_on_client = ?,
    check_model = ?,
    logo = ?,
    description = ?,
    key_vaults = ?,
    settings = ?,
    config = ?,
    updated_at = ?
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteAIProvider :exec
DELETE FROM ai_providers WHERE id = ? AND user_id = ?;

-- AI Models

-- name: GetAIModel :one
SELECT * FROM ai_models
WHERE id = ? AND provider_id = ? AND user_id = ?;

-- name: ListAIModels :many
SELECT * FROM ai_models
WHERE user_id = ?
ORDER BY sort ASC;

-- name: ListAIModelsByProvider :many
SELECT * FROM ai_models
WHERE provider_id = ? AND user_id = ?
ORDER BY sort ASC;

-- name: ListEnabledAIModels :many
SELECT * FROM ai_models
WHERE user_id = ? AND enabled = 1
ORDER BY sort ASC;

-- name: CreateAIModel :one
INSERT INTO ai_models (
    id, display_name, description, organization, enabled, provider_id,
    type, sort, user_id, pricing, parameters, config, abilities,
    context_window_tokens, source, released_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateAIModel :one
UPDATE ai_models
SET display_name = ?,
    description = ?,
    enabled = ?,
    sort = ?,
    pricing = ?,
    parameters = ?,
    config = ?,
    abilities = ?,
    updated_at = ?
WHERE id = ? AND provider_id = ? AND user_id = ?
RETURNING *;

-- name: DeleteAIModel :exec
DELETE FROM ai_models
WHERE id = ? AND provider_id = ? AND user_id = ?;

