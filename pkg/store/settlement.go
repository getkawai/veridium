package store

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// SettlementSnapshot represents a snapshot of contributor balances for settlement
type SettlementSnapshot struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

// SettlementConfig configures the settlement process
type SettlementConfig struct {
	BatchSize      int           // Number of contributors to process per batch (default: 50)
	BatchDelay     time.Duration // Delay between batches (default: 100ms)
	MaxRetries     int           // Max retries for failed operations (default: 3)
	EnableRollback bool          // Enable rollback on failure (default: true)
}

// DefaultSettlementConfig returns default configuration
func DefaultSettlementConfig() *SettlementConfig {
	return &SettlementConfig{
		BatchSize:      50,
		BatchDelay:     100 * time.Millisecond,
		MaxRetries:     3,
		EnableRollback: true,
	}
}

// GenerateUniquePeriodID generates a unique period ID using nanoseconds + random suffix
func GenerateUniquePeriodID() int64 {
	baseID := time.Now().UnixNano()

	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomSuffix := int64(randomBytes[0]) + int64(randomBytes[1])<<8

	return baseID + randomSuffix
}

// GetSettlementSnapshots returns current balances for all contributors (before settlement)
// Results are sorted by address (lowercase) for consistent Merkle tree generation
func (s *KVStore) GetSettlementSnapshots(ctx context.Context, rewardType string) ([]*SettlementSnapshot, error) {
	contributors, err := s.ListContributors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list contributors: %w", err)
	}

	snapshots := make([]*SettlementSnapshot, 0)

	for _, c := range contributors {
		var balance string
		if rewardType == "kawai" {
			balance = c.AccumulatedRewards
		} else {
			balance = c.AccumulatedUSDT
		}

		// Skip if balance is zero or empty
		if balance == "" || balance == "0" {
			continue
		}

		snapshots = append(snapshots, &SettlementSnapshot{
			Address: c.WalletAddress,
			Amount:  balance,
		})
	}

	// CRITICAL: Sort by address (lowercase) for consistent Merkle tree ordering
	sort.Slice(snapshots, func(i, j int) bool {
		return strings.ToLower(snapshots[i].Address) < strings.ToLower(snapshots[j].Address)
	})

	slog.Info("Generated snapshots for settlement", "count", len(snapshots), "type", rewardType, "sorted", true)
	return snapshots, nil
}

// PerformSettlement executes a complete settlement cycle with rollback support
func (s *KVStore) PerformSettlement(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData) (*SettlementPeriod, error) {
	return s.PerformSettlementWithConfig(ctx, periodID, merkleRoot, rewardType, proofs, DefaultSettlementConfig())
}

// PerformSettlementWithConfig executes settlement with custom configuration
func (s *KVStore) PerformSettlementWithConfig(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData, config *SettlementConfig) (*SettlementPeriod, error) {
	slog.Info("Starting settlement", "period", periodID, "type", rewardType, "contributors", len(proofs))

	// Check if period already exists (prevent collision)
	existing, err := s.GetSettlementPeriod(ctx, periodID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("period %d already exists with status %s", periodID, existing.Status)
	}

	// Create initial settlement record
	period := &SettlementPeriod{
		PeriodID:         periodID,
		MerkleRoot:       merkleRoot,
		RewardType:       rewardType,
		Status:           SettlementStatusPending,
		ContributorCount: len(proofs),
		StartedAt:        time.Now(),
		StartDate:        time.Unix(periodID/1e9-7*24*3600, 0),
		EndDate:          time.Unix(periodID/1e9, 0),
	}

	// Save initial status
	if err := s.SaveSettlementPeriod(ctx, period); err != nil {
		return nil, fmt.Errorf("failed to save initial settlement: %w", err)
	}

	// Get snapshots for total calculation
	snapshots, err := s.GetSettlementSnapshots(ctx, rewardType)
	if err != nil {
		period.Status = SettlementStatusFailed
		period.Error = fmt.Sprintf("failed to get snapshots: %v", err)
		s.SaveSettlementPeriod(ctx, period)
		return nil, err
	}

	// Calculate total
	totalAmount := new(big.Int)
	for _, snap := range snapshots {
		amount := new(big.Int)
		amount.SetString(snap.Amount, 10)
		totalAmount.Add(totalAmount, amount)
	}
	period.TotalAmount = totalAmount.String()

	// Track saved proofs for rollback
	savedProofs := make([]string, 0)

	// Step 1: Save all proofs first (before resetting any balances)
	slog.Info("Saving proofs in batches", "total_proofs", len(proofs), "batch_size", config.BatchSize)

	addresses := make([]string, 0, len(proofs))
	for addr := range proofs {
		addresses = append(addresses, addr)
	}

	// Sort addresses for consistent ordering
	sort.Strings(addresses)

	for i := 0; i < len(addresses); i += config.BatchSize {
		end := i + config.BatchSize
		if end > len(addresses) {
			end = len(addresses)
		}

		batch := addresses[i:end]

		for _, addr := range batch {
			proof := proofs[addr]
			proof.PeriodID = periodID
			proof.MerkleRoot = merkleRoot
			proof.RewardType = rewardType
			proof.CreatedAt = time.Now()
			proof.ClaimStatus = ClaimStatusUnclaimed
			proof.Address = addr

			var lastErr error
			for retry := 0; retry < config.MaxRetries; retry++ {
				if err := s.SaveMerkleProofForPeriod(ctx, addr, periodID, proof); err != nil {
					lastErr = err
					time.Sleep(50 * time.Millisecond)
					continue
				}
				savedProofs = append(savedProofs, addr)
				lastErr = nil
				break
			}

			if lastErr != nil {
				slog.Warn("Failed to save proof", "address", addr, "retries", config.MaxRetries, "error", lastErr)

				// Rollback if enabled
				if config.EnableRollback {
					slog.Warn("Rolling back saved proofs", "count", len(savedProofs))
					for _, savedAddr := range savedProofs {
						s.DeleteMerkleProof(ctx, savedAddr, periodID)
					}

					period.Status = SettlementStatusFailed
					period.Error = fmt.Sprintf("failed to save proof for %s: %v", addr, lastErr)
					s.SaveSettlementPeriod(ctx, period)
					return nil, fmt.Errorf("settlement failed, rolled back: %w", lastErr)
				}
			}
		}

		// Delay between batches to avoid rate limiting
		if i+config.BatchSize < len(addresses) {
			time.Sleep(config.BatchDelay)
		}
	}

	period.ProofsSaved = len(savedProofs)
	period.Status = SettlementStatusProofsSaved
	s.SaveSettlementPeriod(ctx, period)

	slog.Info("Saved proofs successfully", "count", len(savedProofs))

	// Step 2: Reset all balances (only after ALL proofs are saved)
	slog.Info("Resetting balances", "count", len(snapshots))

	resetCount := 0
	for i := 0; i < len(snapshots); i += config.BatchSize {
		end := i + config.BatchSize
		if end > len(snapshots) {
			end = len(snapshots)
		}

		batch := snapshots[i:end]

		for _, snapshot := range batch {
			var lastErr error
			for retry := 0; retry < config.MaxRetries; retry++ {
				if err := s.ResetAccumulatedRewards(ctx, snapshot.Address, rewardType); err != nil {
					lastErr = err
					time.Sleep(50 * time.Millisecond)
					continue
				}
				resetCount++
				lastErr = nil
				break
			}

			if lastErr != nil {
				slog.Warn("Failed to reset balance (manual reset needed)", "address", snapshot.Address, "error", lastErr)
			}
		}

		if i+config.BatchSize < len(snapshots) {
			time.Sleep(config.BatchDelay)
		}
	}

	period.BalancesReset = resetCount
	period.Status = SettlementStatusBalancesReset
	s.SaveSettlementPeriod(ctx, period)

	slog.Info("Reset balances", "count", resetCount)

	// Step 3: Mark as completed
	period.Status = SettlementStatusCompleted
	period.CompletedAt = time.Now()

	if err := s.SaveSettlementPeriod(ctx, period); err != nil {
		return nil, fmt.Errorf("failed to save final settlement status: %w", err)
	}

	slog.Info("Completed settlement", "period", periodID, "proofs", period.ProofsSaved, "resets", period.BalancesReset, "total", period.TotalAmount)

	return period, nil
}

// ResumeSettlement attempts to resume a failed or interrupted settlement
func (s *KVStore) ResumeSettlement(ctx context.Context, periodID int64, proofs map[string]*MerkleProofData, config *SettlementConfig) (*SettlementPeriod, error) {
	period, err := s.GetSettlementPeriod(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settlement period: %w", err)
	}

	slog.Info("Resuming settlement", "period", periodID, "status", period.Status)

	switch period.Status {
	case SettlementStatusCompleted:
		return period, nil // Already done

	case SettlementStatusPending:
		// Restart from beginning
		return s.PerformSettlementWithConfig(ctx, periodID, period.MerkleRoot, period.RewardType, proofs, config)

	case SettlementStatusProofsSaved:
		// Resume from balance reset
		snapshots, err := s.GetSettlementSnapshots(ctx, period.RewardType)
		if err != nil {
			return nil, err
		}

		resetCount := 0
		for _, snapshot := range snapshots {
			if err := s.ResetAccumulatedRewards(ctx, snapshot.Address, period.RewardType); err != nil {
				slog.Warn("Failed to reset balance", "address", snapshot.Address, "error", err)
				continue
			}
			resetCount++
		}

		period.BalancesReset = resetCount
		period.Status = SettlementStatusCompleted
		period.CompletedAt = time.Now()
		s.SaveSettlementPeriod(ctx, period)

		return period, nil

	case SettlementStatusFailed:
		// Clear error and retry
		period.Error = ""
		return s.PerformSettlementWithConfig(ctx, periodID, period.MerkleRoot, period.RewardType, proofs, config)

	default:
		return nil, fmt.Errorf("unknown settlement status: %s", period.Status)
	}
}

// PerformSettlementParallel executes settlement with parallel processing
func (s *KVStore) PerformSettlementParallel(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData, workers int) (*SettlementPeriod, error) {
	slog.Info("Starting parallel settlement", "workers", workers)

	if workers <= 0 {
		workers = 5
	}

	// Check if period already exists
	existing, err := s.GetSettlementPeriod(ctx, periodID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("period %d already exists", periodID)
	}

	// Create initial settlement record
	period := &SettlementPeriod{
		PeriodID:         periodID,
		MerkleRoot:       merkleRoot,
		RewardType:       rewardType,
		Status:           SettlementStatusPending,
		ContributorCount: len(proofs),
		StartedAt:        time.Now(),
	}
	s.SaveSettlementPeriod(ctx, period)

	// Get snapshots
	snapshots, err := s.GetSettlementSnapshots(ctx, rewardType)
	if err != nil {
		period.Status = SettlementStatusFailed
		period.Error = err.Error()
		s.SaveSettlementPeriod(ctx, period)
		return nil, err
	}

	// Calculate total
	totalAmount := new(big.Int)
	for _, snap := range snapshots {
		amount := new(big.Int)
		amount.SetString(snap.Amount, 10)
		totalAmount.Add(totalAmount, amount)
	}
	period.TotalAmount = totalAmount.String()

	// Create address list
	addresses := make([]string, 0, len(proofs))
	for addr := range proofs {
		addresses = append(addresses, addr)
	}

	// Parallel proof saving
	var wg sync.WaitGroup
	addrChan := make(chan string, len(addresses))
	errChan := make(chan error, len(addresses))
	savedChan := make(chan string, len(addresses))

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr := range addrChan {
				proof := proofs[addr]
				proof.PeriodID = periodID
				proof.MerkleRoot = merkleRoot
				proof.RewardType = rewardType
				proof.CreatedAt = time.Now()
				proof.ClaimStatus = ClaimStatusUnclaimed
				proof.Address = addr

				if err := s.SaveMerkleProofForPeriod(ctx, addr, periodID, proof); err != nil {
					errChan <- fmt.Errorf("failed for %s: %w", addr, err)
					continue
				}
				savedChan <- addr
			}
		}()
	}

	// Send work
	for _, addr := range addresses {
		addrChan <- addr
	}
	close(addrChan)

	// Wait for completion
	wg.Wait()
	close(errChan)
	close(savedChan)

	// Collect results
	savedProofs := make([]string, 0)
	for addr := range savedChan {
		savedProofs = append(savedProofs, addr)
	}

	errors := make([]error, 0)
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		slog.Error("Errors during parallel proof saving", "count", len(errors))
		// Rollback
		for _, addr := range savedProofs {
			s.DeleteMerkleProof(ctx, addr, periodID)
		}
		period.Status = SettlementStatusFailed
		period.Error = fmt.Sprintf("%d proofs failed", len(errors))
		s.SaveSettlementPeriod(ctx, period)
		return nil, errors[0]
	}

	period.ProofsSaved = len(savedProofs)
	period.Status = SettlementStatusProofsSaved
	s.SaveSettlementPeriod(ctx, period)

	// Reset balances (sequential to avoid race conditions)
	resetCount := 0
	for _, snapshot := range snapshots {
		if err := s.ResetAccumulatedRewards(ctx, snapshot.Address, rewardType); err != nil {
			slog.Warn("Failed to reset balance", "address", snapshot.Address, "error", err)
			continue
		}
		resetCount++
	}

	period.BalancesReset = resetCount
	period.Status = SettlementStatusCompleted
	period.CompletedAt = time.Now()
	s.SaveSettlementPeriod(ctx, period)

	return period, nil
}

// GetClaimableRewards returns all unclaimed proofs for a contributor with total amount
func (s *KVStore) GetClaimableRewards(ctx context.Context, address string) (map[string]interface{}, error) {
	proofs, err := s.ListMerkleProofs(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to list proofs: %w", err)
	}

	totalKawaiClaimable := new(big.Int)
	totalUSDTClaimable := new(big.Int)
	claimableProofs := make([]*MerkleProofData, 0)
	pendingProofs := make([]*MerkleProofData, 0)

	for _, proof := range proofs {
		// Skip confirmed claims
		if proof.ClaimStatus == ClaimStatusConfirmed {
			continue
		}

		amount := new(big.Int)
		amount.SetString(proof.Amount, 10)

		if proof.ClaimStatus == ClaimStatusPending {
			pendingProofs = append(pendingProofs, proof)
		} else {
			claimableProofs = append(claimableProofs, proof)

			if proof.RewardType == "usdt" {
				totalUSDTClaimable.Add(totalUSDTClaimable, amount)
			} else {
				totalKawaiClaimable.Add(totalKawaiClaimable, amount)
			}
		}
	}

	// Get current accumulating balance
	contributor, err := s.GetContributor(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributor: %w", err)
	}

	result := map[string]interface{}{
		"address":                    address,
		"unclaimed_proofs":           claimableProofs,
		"pending_proofs":             pendingProofs,
		"total_kawai_claimable":      totalKawaiClaimable.String(),
		"total_usdt_claimable":       totalUSDTClaimable.String(),
		"current_kawai_accumulating": contributor.AccumulatedRewards,
		"current_usdt_accumulating":  contributor.AccumulatedUSDT,
	}

	return result, nil
}

// MarkClaimPending marks a proof as pending claim (transaction submitted)
func (s *KVStore) MarkClaimPending(ctx context.Context, address string, periodID int64, txHash string) error {
	proof, err := s.GetMerkleProofForPeriod(ctx, address, periodID)
	if err != nil {
		return fmt.Errorf("failed to get proof: %w", err)
	}

	proof.ClaimStatus = ClaimStatusPending
	proof.ClaimTxHash = txHash
	proof.ClaimAttempts++

	if err := s.SaveMerkleProofForPeriod(ctx, address, periodID, proof); err != nil {
		return fmt.Errorf("failed to update proof: %w", err)
	}

	slog.Info("Marked proof as pending", "address", address, "period", periodID, "tx", txHash)
	return nil
}

// ConfirmClaim marks a proof as successfully claimed
func (s *KVStore) ConfirmClaim(ctx context.Context, address string, periodID int64) error {
	proof, err := s.GetMerkleProofForPeriod(ctx, address, periodID)
	if err != nil {
		return fmt.Errorf("failed to get proof: %w", err)
	}

	proof.ClaimStatus = ClaimStatusConfirmed
	proof.ClaimedAt = time.Now()

	if err := s.SaveMerkleProofForPeriod(ctx, address, periodID, proof); err != nil {
		return fmt.Errorf("failed to update proof: %w", err)
	}

	slog.Info("Confirmed claim", "address", address, "period", periodID)
	return nil
}

// MarkClaimFailed marks a claim as failed (transaction reverted)
func (s *KVStore) MarkClaimFailed(ctx context.Context, address string, periodID int64, reason string) error {
	proof, err := s.GetMerkleProofForPeriod(ctx, address, periodID)
	if err != nil {
		return fmt.Errorf("failed to get proof: %w", err)
	}

	proof.ClaimStatus = ClaimStatusFailed
	proof.ClaimTxHash = "" // Clear for retry

	if err := s.SaveMerkleProofForPeriod(ctx, address, periodID, proof); err != nil {
		return fmt.Errorf("failed to update proof: %w", err)
	}

	slog.Warn("Marked claim as failed", "address", address, "period", periodID, "reason", reason)
	return nil
}

// RetryFailedClaim resets a failed claim to unclaimed status for retry
func (s *KVStore) RetryFailedClaim(ctx context.Context, address string, periodID int64) error {
	proof, err := s.GetMerkleProofForPeriod(ctx, address, periodID)
	if err != nil {
		return fmt.Errorf("failed to get proof: %w", err)
	}

	if proof.ClaimStatus != ClaimStatusFailed {
		return fmt.Errorf("proof is not in failed status: %s", proof.ClaimStatus)
	}

	proof.ClaimStatus = ClaimStatusUnclaimed
	proof.ClaimTxHash = ""

	if err := s.SaveMerkleProofForPeriod(ctx, address, periodID, proof); err != nil {
		return fmt.Errorf("failed to update proof: %w", err)
	}

	slog.Info("Reset failed claim for retry", "address", address, "period", periodID)
	return nil
}

// MarkProofAsClaimed is a convenience method that confirms the claim
func (s *KVStore) MarkProofAsClaimed(ctx context.Context, address string, periodID int64) error {
	return s.ConfirmClaim(ctx, address, periodID)
}

// GetPendingClaims returns all proofs with pending claim status
func (s *KVStore) GetPendingClaims(ctx context.Context) ([]*MerkleProofData, error) {
	// List all proof keys in proofs namespace
	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.proofsNamespaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list proof keys: %w", err)
	}

	pendingProofs := make([]*MerkleProofData, 0)

	for _, key := range resp.Result {
		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.proofsNamespaceID,
			Key:         key.Name,
		})
		if err != nil {
			continue
		}

		var proof MerkleProofData
		if err := json.Unmarshal(value, &proof); err != nil {
			continue
		}

		if proof.ClaimStatus == ClaimStatusPending {
			pendingProofs = append(pendingProofs, &proof)
		}
	}

	return pendingProofs, nil
}

// EnsureAdminExists ensures admin account exists (auto-register if needed)
func (s *KVStore) EnsureAdminExists(ctx context.Context, adminAddress string) error {
	_, err := s.GetContributor(ctx, adminAddress)
	if err == nil {
		return nil // Already exists
	}

	// Auto-register admin
	admin := &ContributorData{
		WalletAddress: adminAddress,
		RegisteredAt:  time.Now(),
		LastSeen:      time.Now(),
		Status:        ContributorStatusAdmin,
		IsActive:      true,
		IsAdmin:       true,
	}

	if err := s.SaveContributor(ctx, admin); err != nil {
		return fmt.Errorf("failed to register admin: %w", err)
	}

	slog.Info("Auto-registered admin account", "address", adminAddress)
	return nil
}
