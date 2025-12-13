-- AI Providers

-- name: GetAIProvider :one
SELECT * FROM ai_providers WHERE id = ?;

-- name: ListAIProviders :many
SELECT * FROM ai_providers
ORDER BY sort ASC;

-- name: ListEnabledAIProviders :many
SELECT * FROM ai_providers
WHERE enabled = 1
ORDER BY sort ASC;

-- name: CreateAIProvider :one
INSERT INTO ai_providers (
    id, name, sort, enabled, fetch_on_client, check_model,
    logo, description, key_vaults, source, settings, config,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
WHERE id = ?
RETURNING *;

-- name: DeleteAIProvider :exec
DELETE FROM ai_providers WHERE id = ?;

-- name: DeleteAllAIProviders :exec
DELETE FROM ai_providers;

-- name: UpsertAIProvider :one
INSERT INTO ai_providers (
    id, name, sort, enabled, fetch_on_client, check_model,
    logo, description, key_vaults, source, settings, config,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    name = excluded.name,
    sort = excluded.sort,
    enabled = excluded.enabled,
    fetch_on_client = excluded.fetch_on_client,
    check_model = excluded.check_model,
    logo = excluded.logo,
    description = excluded.description,
    key_vaults = excluded.key_vaults,
    settings = excluded.settings,
    config = excluded.config,
    updated_at = excluded.updated_at
RETURNING *;

-- name: UpsertAIProviderConfig :one
INSERT INTO ai_providers (
    id, key_vaults, config, fetch_on_client, check_model,
    source, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    key_vaults = excluded.key_vaults,
    config = excluded.config,
    fetch_on_client = excluded.fetch_on_client,
    check_model = excluded.check_model,
    updated_at = excluded.updated_at
RETURNING *;

-- name: ToggleAIProviderEnabled :one
INSERT INTO ai_providers (
    id, enabled, source, created_at, updated_at
) VALUES (?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    enabled = excluded.enabled,
    updated_at = excluded.updated_at
RETURNING *;

-- name: GetAIProviderListSimple :many
SELECT 
    id,
    name,
    logo,
    description,
    enabled,
    sort,
    source
FROM ai_providers
ORDER BY sort ASC, updated_at DESC;

-- name: GetAIProviderDetail :one
SELECT 
    id,
    name,
    logo,
    description,
    enabled,
    source,
    key_vaults,
    settings,
    config,
    fetch_on_client,
    check_model
FROM ai_providers
WHERE id = ?;

-- name: GetAIProviderRuntimeConfigs :many
SELECT 
    id,
    key_vaults,
    settings,
    config,
    fetch_on_client
FROM ai_providers;

-- name: DeleteModelsByProvider :exec
DELETE FROM ai_models
WHERE provider_id = ?;

-- AI Models

-- name: GetAIModel :one
SELECT * FROM ai_models
WHERE id = ? AND provider_id = ?;

-- name: ListAIModels :many
SELECT * FROM ai_models
ORDER BY sort ASC;

-- name: ListAIModelsByProvider :many
SELECT * FROM ai_models
WHERE provider_id = ?
ORDER BY sort ASC;

-- name: ListEnabledAIModels :many
SELECT * FROM ai_models
WHERE enabled = 1
ORDER BY sort ASC;

-- name: CreateAIModel :one
INSERT INTO ai_models (
    id, display_name, description, organization, enabled, provider_id,
    type, sort, pricing, parameters, config, abilities,
    context_window_tokens, source, released_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
WHERE id = ? AND provider_id = ?
RETURNING *;

-- name: DeleteAIModel :exec
DELETE FROM ai_models
WHERE id = ? AND provider_id = ?;

-- name: DeleteAllAIModels :exec
DELETE FROM ai_models;

-- name: UpsertAIModel :one
INSERT INTO ai_models (
    id, display_name, description, organization, enabled, provider_id,
    type, sort, pricing, parameters, config, abilities,
    context_window_tokens, source, released_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id, provider_id) DO UPDATE SET
    display_name = excluded.display_name,
    description = excluded.description,
    enabled = excluded.enabled,
    sort = excluded.sort,
    pricing = excluded.pricing,
    parameters = excluded.parameters,
    config = excluded.config,
    abilities = excluded.abilities,
    updated_at = excluded.updated_at
RETURNING *;

-- name: ToggleAIModelEnabled :one
INSERT INTO ai_models (
    id, provider_id, enabled, type, source, updated_at, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id, provider_id) DO UPDATE SET
    enabled = excluded.enabled,
    type = COALESCE(excluded.type, ai_models.type),
    updated_at = excluded.updated_at
RETURNING *;

-- name: UpdateAIModelSort :one
INSERT INTO ai_models (
    id, provider_id, sort, type, enabled, source, updated_at, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id, provider_id) DO UPDATE SET
    sort = excluded.sort,
    type = COALESCE(excluded.type, ai_models.type),
    updated_at = excluded.updated_at
RETURNING *;

-- name: BatchUpdateAIModelEnabled :exec
UPDATE ai_models
SET enabled = ?
WHERE provider_id = ? AND id IN (sqlc.slice('ids'));

-- name: DeleteAIModelsByProviderAndSource :exec
DELETE FROM ai_models
WHERE provider_id = ? AND source = ?;

-- name: DeleteAIModelsByProvider :exec
DELETE FROM ai_models
WHERE provider_id = ?;
