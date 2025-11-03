-- NextAuth Accounts

-- name: GetNextAuthAccount :one
SELECT * FROM nextauth_accounts
WHERE provider = ? AND provider_account_id = ?;

-- name: ListNextAuthAccountsByUser :many
SELECT * FROM nextauth_accounts
WHERE user_id = ?;

-- name: CreateNextAuthAccount :one
INSERT INTO nextauth_accounts (
    access_token, expires_at, id_token, provider, provider_account_id,
    refresh_token, scope, session_state, token_type, type, user_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateNextAuthAccount :one
UPDATE nextauth_accounts
SET access_token = ?,
    expires_at = ?,
    id_token = ?,
    refresh_token = ?,
    scope = ?,
    session_state = ?
WHERE provider = ? AND provider_account_id = ?
RETURNING *;

-- name: DeleteNextAuthAccount :exec
DELETE FROM nextauth_accounts
WHERE provider = ? AND provider_account_id = ?;

-- NextAuth Sessions

-- name: GetNextAuthSession :one
SELECT * FROM nextauth_sessions WHERE session_token = ?;

-- name: CreateNextAuthSession :one
INSERT INTO nextauth_sessions (
    expires, session_token, user_id
) VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateNextAuthSession :one
UPDATE nextauth_sessions
SET expires = ?
WHERE session_token = ?
RETURNING *;

-- name: DeleteNextAuthSession :exec
DELETE FROM nextauth_sessions WHERE session_token = ?;

-- name: DeleteExpiredNextAuthSessions :exec
DELETE FROM nextauth_sessions WHERE expires < ?;

-- NextAuth Verification Tokens

-- name: GetNextAuthVerificationToken :one
SELECT * FROM nextauth_verificationtokens
WHERE identifier = ? AND token = ?;

-- name: CreateNextAuthVerificationToken :one
INSERT INTO nextauth_verificationtokens (
    expires, identifier, token
) VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteNextAuthVerificationToken :exec
DELETE FROM nextauth_verificationtokens
WHERE identifier = ? AND token = ?;

-- name: DeleteExpiredNextAuthVerificationTokens :exec
DELETE FROM nextauth_verificationtokens WHERE expires < ?;

-- NextAuth Authenticators

-- name: GetNextAuthAuthenticator :one
SELECT * FROM nextauth_authenticators
WHERE user_id = ? AND credential_id = ?;

-- name: ListNextAuthAuthenticatorsByUser :many
SELECT * FROM nextauth_authenticators
WHERE user_id = ?;

-- name: CreateNextAuthAuthenticator :one
INSERT INTO nextauth_authenticators (
    counter, credential_backed_up, credential_device_type, credential_id,
    credential_public_key, provider_account_id, transports, user_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateNextAuthAuthenticator :one
UPDATE nextauth_authenticators
SET counter = ?,
    credential_backed_up = ?
WHERE user_id = ? AND credential_id = ?
RETURNING *;

-- name: DeleteNextAuthAuthenticator :exec
DELETE FROM nextauth_authenticators
WHERE user_id = ? AND credential_id = ?;

