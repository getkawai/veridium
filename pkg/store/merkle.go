package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// ClaimStatus represents the status of a claim
type ClaimStatus string

const (
	ClaimStatusUnclaimed ClaimStatus = "unclaimed"
	ClaimStatusPending   ClaimStatus = "pending"
	ClaimStatusConfirmed ClaimStatus = "confirmed"
	ClaimStatusFailed    ClaimStatus = "failed"
)

// MerkleProofData represents a Merkle proof for a specific settlement period
type MerkleProofData struct {
	Index         uint64      `json:"index"`
	Amount        string      `json:"amount"`                   // BigInt as string
	Proof         []string    `json:"proof"`                    // Hex strings
	MerkleRoot    string      `json:"merkle_root"`              // Root hash for this settlement period
	PeriodID      int64       `json:"period_id"`                // Settlement period timestamp
	CreatedAt     time.Time   `json:"created_at"`               // When proof was generated
	RewardType    string      `json:"reward_type,omitempty"`    // "kawai" or "usdt"
	ClaimStatus   ClaimStatus `json:"claim_status,omitempty"`   // Claim status tracking
	ClaimTxHash   string      `json:"claim_tx_hash,omitempty"`  // Transaction hash when claiming
	ClaimAttempts int         `json:"claim_attempts,omitempty"` // Number of claim attempts
	ClaimedAt     time.Time   `json:"claimed_at,omitempty"`     // When claimed successfully
	Address       string      `json:"address,omitempty"`        // Contributor address (for listing all proofs)
}

// SettlementStatus represents the status of a settlement
type SettlementStatus string

const (
	SettlementStatusPending       SettlementStatus = "pending"
	SettlementStatusProofsSaved   SettlementStatus = "proofs_saved"
	SettlementStatusBalancesReset SettlementStatus = "balances_reset"
	SettlementStatusCompleted     SettlementStatus = "completed"
	SettlementStatusFailed        SettlementStatus = "failed"
)

// SettlementPeriod represents a weekly settlement cycle
type SettlementPeriod struct {
	PeriodID         int64            `json:"period_id"`   // Unix timestamp of settlement
	MerkleRoot       string           `json:"merkle_root"` // Root hash
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	TotalAmount      string           `json:"total_amount"`                // Total rewards distributed
	RewardType       string           `json:"reward_type,omitempty"`       // "kawai" or "usdt"
	Status           SettlementStatus `json:"status,omitempty"`            // Settlement status
	ContributorCount int              `json:"contributor_count,omitempty"` // Number of contributors
	ProofsSaved      int              `json:"proofs_saved,omitempty"`      // Number of proofs saved
	BalancesReset    int              `json:"balances_reset,omitempty"`    // Number of balances reset
	StartedAt        time.Time        `json:"started_at,omitempty"`        // When settlement started
	CompletedAt      time.Time        `json:"completed_at,omitempty"`      // When settlement completed
	Error            string           `json:"error,omitempty"`             // Error message if failed
}

// =============================================================================
// MERKLE PROOF OPERATIONS (Proofs Namespace)
// =============================================================================
// Key format: {address}:{periodID}
// Example: 0x742d35cc6634c0532925a3b844bc454e4438f44e:1704067200000000000

// SaveMerkleProof stores proof for a specific settlement period (DEPRECATED)
func (s *KVStore) SaveMerkleProof(ctx context.Context, address string, data *MerkleProofData) error {
	slog.Warn("SaveMerkleProof is deprecated, use SaveMerkleProofForPeriod instead")
	if data.PeriodID == 0 {
		data.PeriodID = time.Now().Unix()
	}
	return s.SaveMerkleProofForPeriod(ctx, address, data.PeriodID, data)
}

// GetMerkleProof retrieves the latest proof for an address (DEPRECATED)
func (s *KVStore) GetMerkleProof(ctx context.Context, address string) (*MerkleProofData, error) {
	slog.Warn("GetMerkleProof is deprecated, use GetMerkleProofForPeriod or ListMerkleProofs instead")

	proofs, err := s.ListMerkleProofs(ctx, address)
	if err != nil {
		return nil, err
	}

	if len(proofs) == 0 {
		return nil, fmt.Errorf("no proofs found for address %s", address)
	}

	return proofs[0], nil
}

// SaveMerkleProofForPeriod stores proof for a specific settlement period
// Key format: {address}:{periodID}
func (s *KVStore) SaveMerkleProofForPeriod(ctx context.Context, address string, periodID int64, data *MerkleProofData) error {
	key := ProofKey(address, periodID)

	// Ensure fields are set
	data.PeriodID = periodID
	data.Address = address
	if data.CreatedAt.IsZero() {
		data.CreatedAt = time.Now()
	}

	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal proof data: %w", err)
	}

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.proofsNamespaceID,
		Key:         key,
		Value:       value,
	})
	if err != nil {
		return fmt.Errorf("failed to write proof to KV: %w", err)
	}

	slog.Info("Saved Merkle proof", "address", address, "period", periodID, "amount", data.Amount)
	return nil
}

// GetMerkleProofForPeriod retrieves proof for a specific settlement period
// Key format: {address}:{periodID}
func (s *KVStore) GetMerkleProofForPeriod(ctx context.Context, address string, periodID int64) (*MerkleProofData, error) {
	key := ProofKey(address, periodID)

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.proofsNamespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get proof from KV: %w", err)
	}

	var data MerkleProofData
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proof data: %w", err)
	}

	return &data, nil
}

// ListMerkleProofs retrieves all proofs for a contributor (sorted by period, newest first)
// Key prefix: {address}:
func (s *KVStore) ListMerkleProofs(ctx context.Context, address string) ([]*MerkleProofData, error) {
	prefix := ProofPrefixForAddress(address)

	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.proofsNamespaceID,
		Prefix:      prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list proof keys: %w", err)
	}

	var proofs []*MerkleProofData
	for _, key := range resp.Result {
		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.proofsNamespaceID,
			Key:         key.Name,
		})
		if err != nil {
			slog.Warn("Failed to get proof for key", "key", key.Name, "error", err)
			continue
		}

		var data MerkleProofData
		if err := json.Unmarshal(value, &data); err != nil {
			slog.Warn("Failed to unmarshal proof for key", "key", key.Name, "error", err)
			continue
		}

		proofs = append(proofs, &data)
	}

	// Sort by period ID (newest first)
	for i := 0; i < len(proofs)-1; i++ {
		for j := i + 1; j < len(proofs); j++ {
			if proofs[i].PeriodID < proofs[j].PeriodID {
				proofs[i], proofs[j] = proofs[j], proofs[i]
			}
		}
	}

	return proofs, nil
}

// DeleteMerkleProof deletes a proof for a specific period
// Key format: {address}:{periodID}
func (s *KVStore) DeleteMerkleProof(ctx context.Context, address string, periodID int64) error {
	key := ProofKey(address, periodID)

	_, err := s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: s.proofsNamespaceID,
		Key:         key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete proof from KV: %w", err)
	}

	slog.Info("Deleted Merkle proof", "address", address, "period", periodID)
	return nil
}

// CleanupOldProofs deletes proofs older than specified duration
func (s *KVStore) CleanupOldProofs(ctx context.Context, olderThan time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-olderThan)
	cutoffPeriodID := cutoffTime.Unix()

	// List all proof keys (no prefix in proofs namespace)
	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.proofsNamespaceID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list proof keys: %w", err)
	}

	deletedCount := 0
	for _, key := range resp.Result {
		// Extract period ID from key format "address:periodID"
		_, periodID, err := ParseProofKey(key.Name)
		if err != nil {
			continue
		}

		if periodID < cutoffPeriodID {
			_, err := s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
				NamespaceID: s.proofsNamespaceID,
				Key:         key.Name,
			})
			if err != nil {
				slog.Warn("Failed to delete old proof", "key", key.Name, "error", err)
				continue
			}
			deletedCount++
		}
	}

	slog.Info("Cleaned up old proofs", "count", deletedCount, "older_than", olderThan)
	return deletedCount, nil
}

// =============================================================================
// SETTLEMENT PERIOD OPERATIONS (Settlements Namespace)
// =============================================================================
// Key format: {periodID}
// Example: 1704067200000000000

// SaveSettlementPeriod stores metadata about a settlement period
// Key format: {periodID}
func (s *KVStore) SaveSettlementPeriod(ctx context.Context, period *SettlementPeriod) error {
	key := SettlementKey(period.PeriodID)

	value, err := json.Marshal(period)
	if err != nil {
		return fmt.Errorf("failed to marshal settlement period: %w", err)
	}

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.settlementsNamespaceID,
		Key:         key,
		Value:       value,
	})
	if err != nil {
		return fmt.Errorf("failed to write settlement period to KV: %w", err)
	}

	slog.Info("Saved settlement period", "period", period.PeriodID, "root", period.MerkleRoot, "total", period.TotalAmount)
	return nil
}

// GetSettlementPeriod retrieves metadata about a settlement period
// Key format: {periodID}
func (s *KVStore) GetSettlementPeriod(ctx context.Context, periodID int64) (*SettlementPeriod, error) {
	key := SettlementKey(periodID)

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.settlementsNamespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get settlement period from KV: %w", err)
	}

	var period SettlementPeriod
	if err := json.Unmarshal(value, &period); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settlement period: %w", err)
	}

	return &period, nil
}

// ListSettlementPeriods retrieves all settlement periods (sorted by period ID, newest first)
func (s *KVStore) ListSettlementPeriods(ctx context.Context) ([]*SettlementPeriod, error) {
	// List all keys (no prefix needed - entire namespace is settlements)
	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.settlementsNamespaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list settlement keys: %w", err)
	}

	var periods []*SettlementPeriod
	for _, key := range resp.Result {
		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.settlementsNamespaceID,
			Key:         key.Name,
		})
		if err != nil {
			slog.Warn("Failed to get settlement for key", "key", key.Name, "error", err)
			continue
		}

		var period SettlementPeriod
		if err := json.Unmarshal(value, &period); err != nil {
			slog.Warn("Failed to unmarshal settlement for key", "key", key.Name, "error", err)
			continue
		}

		periods = append(periods, &period)
	}

	// Sort by period ID (newest first)
	for i := 0; i < len(periods)-1; i++ {
		for j := i + 1; j < len(periods); j++ {
			if periods[i].PeriodID < periods[j].PeriodID {
				periods[i], periods[j] = periods[j], periods[i]
			}
		}
	}

	return periods, nil
}
