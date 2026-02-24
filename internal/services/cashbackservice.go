package services

import (
	"context"
	"time"

	"github.com/kawai-network/x/store"
)

// CashbackService exposes cashback functionality to the Wails frontend
type CashbackService struct {
	kvStore *store.KVStore
}

// NewCashbackService creates a new cashback service
func NewCashbackService(kvStore *store.KVStore) *CashbackService {
	return &CashbackService{
		kvStore: kvStore,
	}
}

// CashbackStatsResponse represents cashback statistics for the frontend
type CashbackStatsResponse struct {
	Total_Cashback       string `json:"total_cashback"`       // Total KAWAI earned (wei)
	Pending_Cashback     string `json:"pending_cashback"`     // Pending (unclaimed) (wei)
	Claimed_Cashback     string `json:"claimed_cashback"`     // Already claimed (wei)
	Total_Deposits       uint64 `json:"total_deposits"`       // Number of deposits
	Total_Deposit_Amount string `json:"total_deposit_amount"` // Total USDT deposited (wei, 6 decimals)
	First_Deposit_At     string `json:"first_deposit_at"`     // ISO timestamp
	Last_Deposit_At      string `json:"last_deposit_at"`      // ISO timestamp
}

// formatTimestamp converts a time pointer to ISO 8601 string
func formatTimestamp(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z07:00")
}

// GetCashbackStats retrieves cashback statistics for the current user
func (s *CashbackService) GetCashbackStats(userAddress string) (*CashbackStatsResponse, error) {
	ctx := context.Background()

	stats, err := s.kvStore.GetCashbackStats(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	response := &CashbackStatsResponse{
		Total_Cashback:       stats.TotalCashback,
		Pending_Cashback:     stats.PendingCashback,
		Claimed_Cashback:     stats.ClaimedCashback,
		Total_Deposits:       stats.TotalDeposits,
		Total_Deposit_Amount: stats.TotalDepositAmount,
		First_Deposit_At:     formatTimestamp(stats.FirstDepositAt),
		Last_Deposit_At:      formatTimestamp(stats.LastDepositAt),
	}

	return response, nil
}

// GetCurrentPeriod returns the current settlement period
func (s *CashbackService) GetCurrentPeriod() uint64 {
	return s.kvStore.GetCurrentPeriod()
}

// ClaimableCashbackRecord represents a claimable cashback record for the frontend
type ClaimableCashbackRecord struct {
	Period        uint64   `json:"period"`
	Amount        string   `json:"amount"`          // KAWAI amount (wei)
	Proof         []string `json:"proof"`           // Merkle proof
	MerkleRoot    string   `json:"merkle_root"`     // Merkle root for verification
	Claimed       bool     `json:"claimed"`         // Whether already claimed
	DepositTxHash string   `json:"deposit_tx_hash"` // Original deposit tx (if available)
}

// GetClaimableCashback retrieves all claimable cashback records for a user
func (s *CashbackService) GetClaimableCashback(userAddress string) ([]ClaimableCashbackRecord, error) {
	ctx := context.Background()

	records, err := s.kvStore.GetClaimableCashbackRecords(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// Convert to frontend format
	var result []ClaimableCashbackRecord
	for _, record := range records {
		result = append(result, ClaimableCashbackRecord{
			Period:        record.Period,
			Amount:        record.CashbackAmount,
			Proof:         record.Proof,
			MerkleRoot:    record.MerkleRoot,
			Claimed:       record.Claimed,
			DepositTxHash: record.DepositTxHash,
		})
	}

	return result, nil
}
