-- name: GetWalletTransaction :one
SELECT * FROM wallet_transactions
WHERE id = ?;

-- name: GetWalletTransactionByHash :one
SELECT * FROM wallet_transactions
WHERE tx_hash = ?;

-- name: ListWalletTransactions :many
SELECT * FROM wallet_transactions
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByAddress :many
SELECT * FROM wallet_transactions
WHERE from_address = ? OR to_address = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByFromAddress :many
SELECT * FROM wallet_transactions
WHERE from_address = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByToAddress :many
SELECT * FROM wallet_transactions
WHERE to_address = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByType :many
SELECT * FROM wallet_transactions
WHERE tx_type = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByStatus :many
SELECT * FROM wallet_transactions
WHERE status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByNetwork :many
SELECT * FROM wallet_transactions
WHERE network = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByAddressAndType :many
SELECT * FROM wallet_transactions
WHERE (from_address = ? OR to_address = ?)
  AND tx_type = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWalletTransactionsByDateRange :many
SELECT * FROM wallet_transactions
WHERE created_at >= ?
  AND created_at <= ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountWalletTransactions :one
SELECT COUNT(*) FROM wallet_transactions;

-- name: CountWalletTransactionsByAddress :one
SELECT COUNT(*) FROM wallet_transactions
WHERE from_address = ? OR to_address = ?;

-- name: CountWalletTransactionsByType :one
SELECT COUNT(*) FROM wallet_transactions
WHERE tx_type = ?;

-- name: CountWalletTransactionsByStatus :one
SELECT COUNT(*) FROM wallet_transactions
WHERE status = ?;

-- name: CreateWalletTransaction :one
INSERT INTO wallet_transactions (
  tx_hash, tx_type, from_address, to_address, amount,
  token_address, token_symbol, token_decimals,
  status, description, block_number, network, chain_id,
  gas_used, gas_price, nonce, metadata,
  created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateWalletTransaction :one
UPDATE wallet_transactions
SET 
  tx_type = ?,
  from_address = ?,
  to_address = ?,
  amount = ?,
  token_address = ?,
  token_symbol = ?,
  token_decimals = ?,
  status = ?,
  description = ?,
  block_number = ?,
  network = ?,
  chain_id = ?,
  gas_used = ?,
  gas_price = ?,
  nonce = ?,
  metadata = ?,
  updated_at = ?
WHERE id = ?
RETURNING *;

-- name: UpdateWalletTransactionStatus :one
UPDATE wallet_transactions
SET 
  status = ?,
  block_number = ?,
  updated_at = ?
WHERE tx_hash = ?
RETURNING *;

-- name: DeleteWalletTransaction :exec
DELETE FROM wallet_transactions
WHERE id = ?;

-- name: DeleteWalletTransactionByHash :exec
DELETE FROM wallet_transactions
WHERE tx_hash = ?;

-- name: GetLatestWalletTransactionsByAddress :many
SELECT * FROM wallet_transactions
WHERE from_address = ? OR to_address = ?
ORDER BY created_at DESC
LIMIT ?;

-- name: GetWalletTransactionsSummaryByAddress :one
SELECT 
  COUNT(*) as total_transactions,
  SUM(CASE WHEN status = 'Success' THEN 1 ELSE 0 END) as successful_transactions,
  SUM(CASE WHEN status = 'Failed' THEN 1 ELSE 0 END) as failed_transactions,
  SUM(CASE WHEN status = 'Pending' THEN 1 ELSE 0 END) as pending_transactions
FROM wallet_transactions
WHERE from_address = ? OR to_address = ?;

-- name: GetWalletTransactionsByTokenSymbol :many
SELECT * FROM wallet_transactions
WHERE token_symbol = ?
  AND (from_address = ? OR to_address = ?)
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPendingWalletTransactions :many
SELECT * FROM wallet_transactions
WHERE status = 'Pending'
ORDER BY created_at DESC;

-- name: GetFailedWalletTransactions :many
SELECT * FROM wallet_transactions
WHERE status = 'Failed'
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

