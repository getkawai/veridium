-- OIDC Clients

-- name: GetOIDCClient :one
SELECT * FROM oidc_clients WHERE id = ?;

-- name: ListOIDCClients :many
SELECT * FROM oidc_clients
ORDER BY created_at DESC;

-- name: CreateOIDCClient :one
INSERT INTO oidc_clients (
    id, name, description, client_secret, redirect_uris, grants,
    response_types, scopes, token_endpoint_auth_method, application_type,
    client_uri, logo_uri, policy_uri, tos_uri, is_first_party,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateOIDCClient :one
UPDATE oidc_clients
SET name = ?,
    description = ?,
    redirect_uris = ?,
    grants = ?,
    response_types = ?,
    scopes = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteOIDCClient :exec
DELETE FROM oidc_clients WHERE id = ?;

-- OIDC Sessions

-- name: GetOIDCSession :one
SELECT * FROM oidc_sessions WHERE id = ?;

-- name: CreateOIDCSession :one
INSERT INTO oidc_sessions (
    id, data, expires_at, user_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteOIDCSession :exec
DELETE FROM oidc_sessions WHERE id = ?;

-- name: DeleteExpiredOIDCSessions :exec
DELETE FROM oidc_sessions WHERE expires_at < ?;

-- OIDC Authorization Codes

-- name: GetOIDCAuthorizationCode :one
SELECT * FROM oidc_authorization_codes WHERE id = ?;

-- name: CreateOIDCAuthorizationCode :one
INSERT INTO oidc_authorization_codes (
    id, data, expires_at, consumed_at, user_id, client_id, grant_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ConsumeOIDCAuthorizationCode :exec
UPDATE oidc_authorization_codes
SET consumed_at = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteOIDCAuthorizationCode :exec
DELETE FROM oidc_authorization_codes WHERE id = ?;

-- OIDC Access Tokens

-- name: GetOIDCAccessToken :one
SELECT * FROM oidc_access_tokens WHERE id = ?;

-- name: CreateOIDCAccessToken :one
INSERT INTO oidc_access_tokens (
    id, data, expires_at, consumed_at, user_id, client_id, grant_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteOIDCAccessToken :exec
DELETE FROM oidc_access_tokens WHERE id = ?;

-- OIDC Refresh Tokens

-- name: GetOIDCRefreshToken :one
SELECT * FROM oidc_refresh_tokens WHERE id = ?;

-- name: CreateOIDCRefreshToken :one
INSERT INTO oidc_refresh_tokens (
    id, data, expires_at, consumed_at, user_id, client_id, grant_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteOIDCRefreshToken :exec
DELETE FROM oidc_refresh_tokens WHERE id = ?;

-- OIDC Grants

-- name: GetOIDCGrant :one
SELECT * FROM oidc_grants WHERE id = ?;

-- name: CreateOIDCGrant :one
INSERT INTO oidc_grants (
    id, data, expires_at, consumed_at, user_id, client_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteOIDCGrant :exec
DELETE FROM oidc_grants WHERE id = ?;

-- OIDC Consents

-- name: GetOIDCConsent :one
SELECT * FROM oidc_consents WHERE user_id = ? AND client_id = ?;

-- name: CreateOIDCConsent :one
INSERT INTO oidc_consents (
    user_id, client_id, scopes, expires_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateOIDCConsent :one
UPDATE oidc_consents
SET scopes = ?,
    expires_at = ?,
    updated_at = ?
WHERE user_id = ? AND client_id = ?
RETURNING *;

-- name: DeleteOIDCConsent :exec
DELETE FROM oidc_consents WHERE user_id = ? AND client_id = ?;

-- OIDC Interactions

-- name: GetOIDCInteraction :one
SELECT * FROM oidc_interactions WHERE id = ?;

-- name: CreateOIDCInteraction :one
INSERT INTO oidc_interactions (
    id, data, expires_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteOIDCInteraction :exec
DELETE FROM oidc_interactions WHERE id = ?;

-- OIDC Device Codes

-- name: GetOIDCDeviceCode :one
SELECT * FROM oidc_device_codes WHERE id = ?;

-- name: GetOIDCDeviceCodeByUserCode :one
SELECT * FROM oidc_device_codes WHERE user_code = ?;

-- name: CreateOIDCDeviceCode :one
INSERT INTO oidc_device_codes (
    id, data, expires_at, consumed_at, user_id, client_id, grant_id,
    user_code, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteOIDCDeviceCode :exec
DELETE FROM oidc_device_codes WHERE id = ?;

-- OAuth Handoffs

-- name: GetOAuthHandoff :one
SELECT * FROM oauth_handoffs WHERE id = ?;

-- name: CreateOAuthHandoff :one
INSERT INTO oauth_handoffs (
    id, client, payload, created_at, updated_at
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteOAuthHandoff :exec
DELETE FROM oauth_handoffs WHERE id = ?;

-- name: GetOAuthHandoffByClient :one
SELECT * FROM oauth_handoffs
WHERE id = ? AND client = ? AND created_at > ?;

-- name: CleanupExpiredOAuthHandoffs :exec
DELETE FROM oauth_handoffs WHERE created_at < ?;

