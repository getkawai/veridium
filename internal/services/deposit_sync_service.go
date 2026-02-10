package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/vault"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/x/constant"
)

// DepositSyncService handles manual deposit synchronization from user client
type DepositSyncService struct {
	client       *ethclient.Client
	vault        *vault.PaymentVault
	vaultAddress common.Address
	kvStore      *store.KVStore
}

// NewDepositSyncService creates a new deposit sync service
func NewDepositSyncService(kvStore *store.KVStore) (*DepositSyncService, error) {
	// Connect to blockchain
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	// Load PaymentVault contract
	vaultAddr := common.HexToAddress(constant.PaymentVaultAddress)
	vaultContract, err := vault.NewPaymentVault(vaultAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load PaymentVault contract: %w", err)
	}

	return &DepositSyncService{
		client:       client,
		vault:        vaultContract,
		vaultAddress: vaultAddr,
		kvStore:      kvStore,
	}, nil
}

// SyncDepositRequest represents a sync request from user client
type SyncDepositRequest struct {
	TxHash      string `json:"txHash"`
	UserAddress string `json:"userAddress"`
}

// SyncDepositResponse represents the sync response
type SyncDepositResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Amount      string `json:"amount,omitempty"`
	NewBalance  string `json:"newBalance,omitempty"`
	BlockNumber uint64 `json:"blockNumber,omitempty"`
	AlreadySync bool   `json:"alreadySync,omitempty"`
}

// SyncDeposit verifies and syncs a deposit transaction
// Uses idempotency pattern: safe to call multiple times with same txHash
func (s *DepositSyncService) SyncDeposit(ctx context.Context, req SyncDepositRequest) (*SyncDepositResponse, error) {
	log.Printf("💰 [DepositSync] Sync request: txHash=%s, user=%s", req.TxHash, req.UserAddress)

	// 1. Validate inputs
	if req.TxHash == "" {
		return &SyncDepositResponse{
			Success: false,
			Message: "Transaction hash is required",
		}, nil
	}

	txHash := common.HexToHash(req.TxHash)
	userAddr := common.HexToAddress(req.UserAddress)

	// 2. Get transaction receipt from blockchain (source of truth) with retry and timeout
	receipt, err := s.getTransactionReceiptWithRetry(ctx, txHash)
	if err != nil {
		log.Printf("❌ [DepositSync] Failed to get transaction receipt: %v", err)
		return &SyncDepositResponse{
			Success: false,
			Message: fmt.Sprintf("Transaction not found or not confirmed: %v", err),
		}, nil
	}

	// 3. Verify transaction was successful
	if receipt.Status != 1 {
		log.Printf("❌ [DepositSync] Transaction failed on blockchain: %s", req.TxHash)
		return &SyncDepositResponse{
			Success: false,
			Message: "Transaction failed on blockchain",
		}, nil
	}

	// 4. Parse Deposited event from logs
	var depositAmount *big.Int
	var depositUser common.Address
	found := false

	for _, vLog := range receipt.Logs {
		// Check if log is from PaymentVault
		if vLog.Address != s.vaultAddress {
			continue
		}

		// Try to parse as Deposited event
		event, err := s.vault.ParseDeposited(*vLog)
		if err != nil {
			continue
		}

		// Found Deposited event
		depositUser = event.User
		depositAmount = event.Amount
		found = true
		break
	}

	if !found {
		log.Printf("❌ [DepositSync] No Deposited event found in transaction: %s", req.TxHash)
		return &SyncDepositResponse{
			Success: false,
			Message: "No deposit event found in transaction",
		}, nil
	}

	// 5. Verify user address matches
	if depositUser.Hex() != userAddr.Hex() {
		log.Printf("❌ [DepositSync] User address mismatch: event=%s, request=%s",
			depositUser.Hex(), userAddr.Hex())
		return &SyncDepositResponse{
			Success: false,
			Message: "User address does not match deposit event",
		}, nil
	}

	// 6. IDEMPOTENCY CHECK: Check if already processed
	// This prevents double-spending even if called multiple times
	processedKey := fmt.Sprintf("processed_tx:%s", req.TxHash)
	existing, err := s.kvStore.GetMarketplaceData(ctx, processedKey)
	if err == nil && len(existing) > 0 {
		log.Printf("⚠️  [DepositSync] Transaction already processed: %s", req.TxHash)

		// Get current balance
		balance, _ := s.kvStore.GetUserBalance(ctx, req.UserAddress)

		return &SyncDepositResponse{
			Success:     true,
			Message:     "Deposit already synced",
			AlreadySync: true,
			Amount:      depositAmount.String(),
			NewBalance:  balance.USDTBalance,
			BlockNumber: receipt.BlockNumber.Uint64(),
		}, nil
	}

	// 6.b LOCK CHECK: Check if currently processing (Pending Lock)
	pendingKey := fmt.Sprintf("pending_sync:%s", req.TxHash)
	pending, err := s.kvStore.GetMarketplaceData(ctx, pendingKey)
	if err == nil && len(pending) > 0 {
		log.Printf("⚠️  [DepositSync] Transaction is currently being processed: %s", req.TxHash)
		return &SyncDepositResponse{
			Success: false,
			Message: "Transaction is being processed, please wait...",
		}, nil
	}

	// 6.c ACQUIRE LOCK: Set pending flag with TTL (300s = 5 minutes)
	// TTL is set to 5 minutes to handle:
	// - Network delays (blockchain RPC calls)
	// - KV store latency (multiple read/write operations)
	// - Cashback calculation and tracking
	// Lock is released via defer when function completes successfully
	// TTL ensures lock doesn't persist indefinitely if process crashes
	if err := s.kvStore.StoreMarketplaceDataWithTTL(ctx, pendingKey, []byte("processing"), 300); err != nil {
		log.Printf("❌ [DepositSync] Failed to acquire lock: %v", err)
		return &SyncDepositResponse{
			Success: false,
			Message: "System busy, please try again",
		}, nil
	}
	// Always release lock when function exits
	defer func() { _ = s.kvStore.DeleteMarketplaceData(ctx, pendingKey) }()

	// 7. Update KV Store balance (atomic operation)
	// Note: AddBalanceAtomic has retry logic to handle concurrent updates
	if err := s.kvStore.AddBalanceAtomic(ctx, req.UserAddress, depositAmount); err != nil {
		log.Printf("❌ [DepositSync] Failed to update balance: %v", err)
		return &SyncDepositResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update balance: %v", err),
		}, nil
	}

	// 8. Track cashback for this deposit
	period := s.kvStore.GetCurrentPeriod()
	if err := s.kvStore.TrackCashback(ctx, req.UserAddress, req.TxHash, depositAmount, period); err != nil {
		log.Printf("⚠️  [DepositSync] Failed to track cashback: %v", err)
		// Don't fail - balance already updated, cashback can be manually added if needed
	}

	// 9. Mark transaction as processed (after successful balance update)
	// This makes the operation idempotent
	if err := s.kvStore.StoreMarketplaceData(ctx, processedKey, []byte("completed")); err != nil {
		log.Printf("⚠️  [DepositSync] Failed to mark transaction as processed: %v", err)
		// Don't fail - balance already updated, user can retry if needed
	}

	// 10. Get new balance
	balance, _ := s.kvStore.GetUserBalance(ctx, req.UserAddress)

	log.Printf("✅ [DepositSync] Deposit synced: user=%s, amount=%s USDT, block=%d",
		req.UserAddress, depositAmount.String(), receipt.BlockNumber.Uint64())

	return &SyncDepositResponse{
		Success:     true,
		Message:     "Deposit synced successfully",
		Amount:      depositAmount.String(),
		NewBalance:  balance.USDTBalance,
		BlockNumber: receipt.BlockNumber.Uint64(),
	}, nil
}

// Close closes the service and releases resources
func (s *DepositSyncService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

// getTransactionReceiptWithRetry retrieves transaction receipt with retry logic and timeout
// Retries up to 3 times with exponential backoff for transient network errors
func (s *DepositSyncService) getTransactionReceiptWithRetry(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, constant.BlockchainReceiptTimeout)
	defer cancel()

	var lastErr error
	backoff := constant.BlockchainInitialBackoff

	for attempt := 0; attempt < constant.BlockchainMaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			select {
			case <-time.After(backoff):
				backoff *= 2 // Exponential backoff
				if backoff > constant.BlockchainMaxBackoff {
					backoff = constant.BlockchainMaxBackoff
				}
			case <-timeoutCtx.Done():
				return nil, fmt.Errorf("timeout while retrying: %w", timeoutCtx.Err())
			}
		}

		receipt, err := s.client.TransactionReceipt(timeoutCtx, txHash)
		if err == nil {
			return receipt, nil
		}

		lastErr = err
		log.Printf("⚠️  [DepositSync] Attempt %d/%d failed: %v", attempt+1, constant.BlockchainMaxRetries, err)
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", constant.BlockchainMaxRetries, lastErr)
}
