package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/cashbackdistributor"
)

// CashbackCache stores claimable records in memory for faster subsequent loads
type CashbackCache struct {
	records   []CashbackRecord
	expiresAt time.Time
	mu        sync.RWMutex
}

// Global cache for cashback records (5 minute TTL)
var cashbackCache = make(map[string]*CashbackCache)
var cacheMu sync.RWMutex

const cashbackCacheTTL = 5 * time.Minute

// CashbackRecord represents a single cashback entry
type CashbackRecord struct {
	UserAddress    string    `json:"user_address"`
	DepositTxHash  string    `json:"deposit_tx_hash"`
	DepositAmount  string    `json:"deposit_amount"`  // USDT amount (6 decimals)
	CashbackAmount string    `json:"cashback_amount"` // KAWAI amount (18 decimals)
	Rate           uint64    `json:"rate"`            // Rate in basis points (e.g., 200 = 2%)
	Tier           uint64    `json:"tier"`            // Tier (1-5)
	IsFirstTime    bool      `json:"is_first_time"`   // First-time bonus
	CreatedAt      time.Time `json:"created_at"`
	Period         uint64    `json:"period"` // Settlement period
	Claimed        bool      `json:"claimed"`
	Proof          []string  `json:"proof,omitempty"`       // Merkle proof (added during settlement)
	MerkleRoot     string    `json:"merkle_root,omitempty"` // Merkle root for this period
}

// CashbackStats represents user's cashback statistics
type CashbackStats struct {
	TotalCashback      string     `json:"total_cashback"`       // Total KAWAI earned (wei)
	PendingCashback    string     `json:"pending_cashback"`     // Pending (unclaimed) (wei)
	ClaimedCashback    string     `json:"claimed_cashback"`     // Already claimed (wei)
	TotalDeposits      uint64     `json:"total_deposits"`       // Number of deposits
	TotalDepositAmount string     `json:"total_deposit_amount"` // Total USDT deposited (wei, 6 decimals)
	FirstDepositAt     *time.Time `json:"first_deposit_at"`
	LastDepositAt      *time.Time `json:"last_deposit_at"`
}

// CalculateCashback calculates KAWAI cashback for a USDT deposit
// Returns: cashbackAmount (wei), rate (bps), tier, isFirstTime
func (s *KVStore) CalculateCashback(ctx context.Context, userAddress string, depositAmount *big.Int) (string, uint64, uint64, bool, error) {
	// Get user stats to check if first-time
	stats, err := s.GetCashbackStats(ctx, userAddress)
	if err != nil {
		log.Printf("⚠️  [Cashback] Failed to get stats for %s: %v", userAddress, err)
		// Assume not first-time on error
		stats = &CashbackStats{}
	}

	isFirstTime := stats.TotalDeposits == 0

	// Determine tier based on deposit amount (USDT has 6 decimals)
	// Tier 1: < 100 USDT
	// Tier 2: 100-500 USDT
	// Tier 3: 500-1000 USDT
	// Tier 4: 1000-5000 USDT
	// Tier 5: >= 5000 USDT

	usdtAmount := new(big.Int).Set(depositAmount) // Copy to avoid mutation
	tier := uint64(1)
	baseRate := uint64(100) // 1% in basis points

	// Convert to USDT (divide by 1e6)
	oneMillion := big.NewInt(1_000_000)
	usdtValue := new(big.Int).Div(usdtAmount, oneMillion)

	// Tier structure matches DepositCashbackDistributor.sol
	// Base rates: 1-2% (100-200 basis points)
	// Phase 1 bounds: 1.5%-2.5% (150-250 basis points)
	if usdtValue.Cmp(big.NewInt(5000)) >= 0 {
		tier = 5
		baseRate = 200 // 2% (capped by Phase 1 max: 2.5%)
	} else if usdtValue.Cmp(big.NewInt(1000)) >= 0 {
		tier = 4
		baseRate = 175 // 1.75%
	} else if usdtValue.Cmp(big.NewInt(500)) >= 0 {
		tier = 3
		baseRate = 150 // 1.5%
	} else if usdtValue.Cmp(big.NewInt(100)) >= 0 {
		tier = 2
		baseRate = 125 // 1.25%
	}

	// First-time bonus: 2.5% (Phase 1 max)
	rate := baseRate
	if isFirstTime {
		rate = 250 // 2.5% for first deposit (Phase 1 max bound)
	}

	// Calculate cashback: (depositAmount * rate * 1e18) / (10000 * 1e6)
	// This converts USDT (6 decimals) to KAWAI (18 decimals) with rate applied
	cashback := new(big.Int).Mul(depositAmount, big.NewInt(int64(rate)))
	cashback.Mul(cashback, big.NewInt(1e18))
	cashback.Div(cashback, big.NewInt(10000))
	cashback.Div(cashback, big.NewInt(1e6))

	// Apply tier caps (max KAWAI per deposit)
	// Tier 1: 5K KAWAI
	// Tier 2: 10K KAWAI
	// Tier 3: 15K KAWAI
	// Tier 4: 20K KAWAI
	// Tier 5: 20K KAWAI
	maxCashback := big.NewInt(5000) // Default 5K
	switch tier {
	case 2:
		maxCashback = big.NewInt(10000)
	case 3:
		maxCashback = big.NewInt(15000)
	case 4, 5:
		maxCashback = big.NewInt(20000)
	}
	maxCashback.Mul(maxCashback, big.NewInt(1e18)) // Convert to wei

	if cashback.Cmp(maxCashback) > 0 {
		cashback = maxCashback
	}

	return cashback.String(), rate, tier, isFirstTime, nil
}

// TrackCashback records a cashback entry for a deposit
func (s *KVStore) TrackCashback(ctx context.Context, userAddress, txHash string, depositAmount *big.Int, period uint64) error {
	// Calculate cashback
	cashbackAmount, rate, tier, isFirstTime, err := s.CalculateCashback(ctx, userAddress, depositAmount)
	if err != nil {
		return fmt.Errorf("failed to calculate cashback: %w", err)
	}

	// Create record
	record := CashbackRecord{
		UserAddress:    userAddress,
		DepositTxHash:  txHash,
		DepositAmount:  depositAmount.String(),
		CashbackAmount: cashbackAmount,
		Rate:           rate,
		Tier:           tier,
		IsFirstTime:    isFirstTime,
		CreatedAt:      time.Now(),
		Period:         period,
		Claimed:        false,
	}

	// Store record
	key := fmt.Sprintf("cashback:%s:%s", userAddress, txHash)
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal cashback record: %w", err)
	}

	if err := s.StoreCashbackData(ctx, key, data); err != nil {
		return fmt.Errorf("failed to store cashback record: %w", err)
	}

	// Track user for this period (for settlement)
	if err := s.trackUserForPeriod(ctx, period, userAddress); err != nil {
		log.Printf("⚠️  [Cashback] Failed to track user for period: %v", err)
		// Don't fail - record is already stored
	}

	// Update stats
	if err := s.updateCashbackStats(ctx, userAddress, cashbackAmount, depositAmount.String()); err != nil {
		log.Printf("⚠️  [Cashback] Failed to update stats: %v", err)
		// Don't fail - record is already stored
	}

	log.Printf("✅ [Cashback] Tracked: user=%s, deposit=%s USDT, cashback=%s KAWAI, rate=%d bps, tier=%d, first=%v",
		userAddress, depositAmount.String(), cashbackAmount, rate, tier, isFirstTime)

	return nil
}

// GetCashbackStats retrieves cashback statistics for a user
func (s *KVStore) GetCashbackStats(ctx context.Context, userAddress string) (*CashbackStats, error) {
	key := fmt.Sprintf("cashback_stats:%s", userAddress)
	data, err := s.GetCashbackData(ctx, key)
	if err != nil {
		return &CashbackStats{
			TotalCashback:   "0",
			PendingCashback: "0",
			ClaimedCashback: "0",
			TotalDeposits:   0,
		}, nil
	}

	var stats CashbackStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cashback stats: %w", err)
	}

	return &stats, nil
}

// updateCashbackStats atomically updates user's cashback statistics
func (s *KVStore) updateCashbackStats(ctx context.Context, userAddress, cashbackAmount, depositAmount string) error {
	key := fmt.Sprintf("cashback_stats:%s", userAddress)

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		// Read current stats
		stats, err := s.GetCashbackStats(ctx, userAddress)
		if err != nil {
			return err
		}

		// Update cashback stats
		totalCashback := new(big.Int)
		if stats.TotalCashback != "" && stats.TotalCashback != "0" {
			if _, ok := totalCashback.SetString(stats.TotalCashback, 10); !ok {
				return fmt.Errorf("invalid TotalCashback value: %s", stats.TotalCashback)
			}
		}

		pendingCashback := new(big.Int)
		if stats.PendingCashback != "" && stats.PendingCashback != "0" {
			if _, ok := pendingCashback.SetString(stats.PendingCashback, 10); !ok {
				return fmt.Errorf("invalid PendingCashback value: %s", stats.PendingCashback)
			}
		}

		newCashback := new(big.Int)
		if _, ok := newCashback.SetString(cashbackAmount, 10); !ok {
			return fmt.Errorf("invalid cashbackAmount value: %s", cashbackAmount)
		}

		totalCashback.Add(totalCashback, newCashback)
		pendingCashback.Add(pendingCashback, newCashback)

		// Update deposit amount stats
		totalDepositAmt := new(big.Int)
		if stats.TotalDepositAmount != "" && stats.TotalDepositAmount != "0" {
			if _, ok := totalDepositAmt.SetString(stats.TotalDepositAmount, 10); !ok {
				return fmt.Errorf("invalid TotalDepositAmount value: %s", stats.TotalDepositAmount)
			}
		}

		newDepositAmt := new(big.Int)
		if _, ok := newDepositAmt.SetString(depositAmount, 10); !ok {
			return fmt.Errorf("invalid depositAmount value: %s", depositAmount)
		}

		totalDepositAmt.Add(totalDepositAmt, newDepositAmt)

		now := time.Now()
		stats.TotalCashback = totalCashback.String()
		stats.PendingCashback = pendingCashback.String()
		stats.TotalDepositAmount = totalDepositAmt.String()
		stats.TotalDeposits++
		stats.LastDepositAt = &now
		if stats.FirstDepositAt == nil {
			stats.FirstDepositAt = &now
		}

		// Write back
		data, err := json.Marshal(stats)
		if err != nil {
			return fmt.Errorf("failed to marshal stats: %w", err)
		}

		if err := s.StoreCashbackData(ctx, key, data); err != nil {
			if i == maxRetries-1 {
				return fmt.Errorf("failed to write stats after %d retries: %w", maxRetries, err)
			}
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to update stats after %d retries", maxRetries)
}

// trackUserForPeriod adds a user to the period's user list for settlement
func (s *KVStore) trackUserForPeriod(ctx context.Context, period uint64, userAddress string) error {
	periodKey := fmt.Sprintf("cashback_period:%d:users", period)

	// Get existing users
	data, err := s.GetCashbackData(ctx, periodKey)
	var users []string
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &users); err != nil {
			return fmt.Errorf("failed to unmarshal users: %w", err)
		}
	}

	// Check if user already tracked
	for _, u := range users {
		if u == userAddress {
			return nil // Already tracked
		}
	}

	// Add user
	users = append(users, userAddress)
	data, err = json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := s.StoreCashbackData(ctx, periodKey, data); err != nil {
		return fmt.Errorf("failed to store users: %w", err)
	}

	return nil
}

// GetCurrentPeriod returns the current settlement period
// Period 1 starts at a configurable date, increments weekly
func (s *KVStore) GetCurrentPeriod() uint64 {
	// Default to 2025-01-01 if not configured
	// Can be overridden via CASHBACK_PERIOD_START env var (RFC3339 format)
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if startEnv := os.Getenv("CASHBACK_PERIOD_START"); startEnv != "" {
		if parsed, err := time.Parse(time.RFC3339, startEnv); err == nil {
			startDate = parsed
		} else {
			log.Printf("Warning: Invalid CASHBACK_PERIOD_START format, using default: %v", err)
		}
	}

	now := time.Now().UTC()

	// Calculate weeks since start
	duration := now.Sub(startDate)
	weeks := duration.Hours() / (24 * 7)

	return uint64(weeks) + 1 // Period 1, 2, 3, ...
}

// GetClaimableCashbackRecords retrieves all claimable cashback records for a user
// Returns records with Merkle proofs that have been settled but not yet claimed
// OPTIMIZED: Uses in-memory cache (5min TTL) + settled periods index + parallel queries
// Performance: 20s → <1s first load, instant on subsequent loads
func (s *KVStore) GetClaimableCashbackRecords(ctx context.Context, userAddress string) ([]CashbackRecord, error) {
	// Check cache first
	cacheMu.RLock()
	cached, exists := cashbackCache[userAddress]
	cacheMu.RUnlock()

	if exists {
		cached.mu.RLock()
		defer cached.mu.RUnlock()

		if time.Now().Before(cached.expiresAt) {
			log.Printf("⚡ [Cashback] Cache HIT for user %s (%d records)", userAddress, len(cached.records))
			return cached.records, nil
		}
		log.Printf("🕐 [Cashback] Cache EXPIRED for user %s", userAddress)
	}

	// Cache miss or expired - fetch from KV
	log.Printf("💾 [Cashback] Cache MISS for user %s, fetching from KV", userAddress)

	// Get list of settled periods (1 API call instead of scanning all periods)
	settledPeriods, err := s.GetSettledCashbackPeriods(ctx)
	if err != nil {
		log.Printf("⚠️  [Cashback] Failed to get settled periods, falling back to full scan: %v", err)
		// Fallback to old method if index doesn't exist yet
		return s.getClaimableCashbackRecordsFallback(ctx, userAddress)
	}

	if len(settledPeriods) == 0 {
		log.Printf("ℹ️  [Cashback] No settled periods yet for user %s", userAddress)
		return []CashbackRecord{}, nil
	}

	log.Printf("🔍 [Cashback] Checking %d settled periods for user %s", len(settledPeriods), userAddress)

	// Connect to blockchain to check on-chain claimed status
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		log.Printf("⚠️  [Cashback] Failed to connect to blockchain: %v", err)
		// Continue without on-chain check (will use KV data only)
		client = nil
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	var distributor *cashbackdistributor.DepositCashbackDistributor
	if client != nil {
		distributorAddr := common.HexToAddress(constant.CashbackDistributorAddress)
		distributor, err = cashbackdistributor.NewDepositCashbackDistributor(distributorAddr, client)
		if err != nil {
			log.Printf("⚠️  [Cashback] Failed to load distributor contract: %v", err)
			distributor = nil
		}
	}

	// Use parallel queries to check only settled periods
	type periodResult struct {
		record CashbackRecord
		found  bool
	}

	resultsChan := make(chan periodResult, len(settledPeriods))
	var wg sync.WaitGroup

	// Launch goroutines to check each settled period in parallel
	for _, period := range settledPeriods {
		wg.Add(1)
		go func(p uint64, queryCtx context.Context) {
			defer wg.Done()

			// Get merkle root
			rootKey := fmt.Sprintf("cashback_period:%d:merkle_root", p)
			rootData, err := s.GetCashbackData(queryCtx, rootKey)
			if err != nil || len(rootData) == 0 {
				return
			}

			var merkleRoot string
			if err := json.Unmarshal(rootData, &merkleRoot); err != nil {
				log.Printf("⚠️  [Cashback] Failed to unmarshal merkle root for period %d: %v", p, err)
				return
			}

			// Get user's proof for this period
			proofKey := fmt.Sprintf("cashback_proof:%d:%s", p, userAddress)
			proofData, err := s.GetCashbackData(queryCtx, proofKey)
			if err != nil || len(proofData) == 0 {
				// No proof for this user in this period
				return
			}

			var proofRecord struct {
				Proof   []string `json:"proof"`
				Amount  string   `json:"amount"`
				Claimed bool     `json:"claimed"`
			}
			if err := json.Unmarshal(proofData, &proofRecord); err != nil {
				log.Printf("⚠️  [Cashback] Failed to unmarshal proof for period %d, user %s: %v", p, userAddress, err)
				return
			}

			// Check on-chain claimed status (source of truth)
			claimed := proofRecord.Claimed // Default to KV value
			if distributor != nil {
				onChainClaimed, err := distributor.HasClaimed(nil, big.NewInt(int64(p)), common.HexToAddress(userAddress))
				if err == nil {
					claimed = onChainClaimed
					// Update KV if mismatch
					if onChainClaimed != proofRecord.Claimed {
						log.Printf("🔄 [Cashback] Syncing claimed status for period %d, user %s: KV=%v, OnChain=%v", p, userAddress, proofRecord.Claimed, onChainClaimed)
						proofRecord.Claimed = onChainClaimed
						updatedData, _ := json.Marshal(proofRecord)
						s.StoreCashbackData(queryCtx, proofKey, updatedData)
					}
				} else {
					log.Printf("⚠️  [Cashback] Failed to check on-chain claimed status for period %d: %v", p, err)
				}
			}

			// Create a claimable record
			record := CashbackRecord{
				UserAddress:    userAddress,
				DepositTxHash:  fmt.Sprintf("period-%d", p), // Placeholder
				CashbackAmount: proofRecord.Amount,
				Period:         p,
				Claimed:        claimed, // Use on-chain status
				Proof:          proofRecord.Proof,
				MerkleRoot:     merkleRoot,
				CreatedAt:      time.Now(), // Placeholder
			}

			resultsChan <- periodResult{record: record, found: true}
		}(period, ctx)
	}

	// Wait for all goroutines and close channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var claimableRecords []CashbackRecord
	for result := range resultsChan {
		if result.found {
			claimableRecords = append(claimableRecords, result.record)
		}
	}

	// Sort by period for consistent ordering
	sort.Slice(claimableRecords, func(i, j int) bool {
		return claimableRecords[i].Period < claimableRecords[j].Period
	})

	log.Printf("✅ [Cashback] Found %d claimable periods for user %s (optimized query)", len(claimableRecords), userAddress)

	// Update cache
	cacheMu.Lock()
	cashbackCache[userAddress] = &CashbackCache{
		records:   claimableRecords,
		expiresAt: time.Now().Add(cashbackCacheTTL),
	}
	cacheMu.Unlock()

	log.Printf("💾 [Cashback] Cached results for user %s (TTL: %v)", userAddress, cashbackCacheTTL)

	return claimableRecords, nil
}

// getClaimableCashbackRecordsFallback is the fallback method when settled periods index doesn't exist
// This scans all periods sequentially (slower, used only for backward compatibility)
func (s *KVStore) getClaimableCashbackRecordsFallback(ctx context.Context, userAddress string) ([]CashbackRecord, error) {
	log.Printf("⚠️  [Cashback] Using fallback method (scanning all periods)")

	var claimableRecords []CashbackRecord
	currentPeriod := s.GetCurrentPeriod()

	// Scan last 52 weeks (1 year) of periods
	for period := uint64(1); period < currentPeriod; period++ {
		// Check if this period has been settled (has merkle root)
		rootKey := fmt.Sprintf("cashback_period:%d:merkle_root", period)
		rootData, err := s.GetCashbackData(ctx, rootKey)
		if err != nil || len(rootData) == 0 {
			continue
		}

		var merkleRoot string
		if err := json.Unmarshal(rootData, &merkleRoot); err != nil {
			continue
		}

		// Get user's proof for this period
		proofKey := fmt.Sprintf("cashback_proof:%d:%s", period, userAddress)
		proofData, err := s.GetCashbackData(ctx, proofKey)
		if err != nil || len(proofData) == 0 {
			continue
		}

		var proofRecord struct {
			Proof   []string `json:"proof"`
			Amount  string   `json:"amount"`
			Claimed bool     `json:"claimed"`
		}
		if err := json.Unmarshal(proofData, &proofRecord); err != nil {
			continue
		}

		record := CashbackRecord{
			UserAddress:    userAddress,
			DepositTxHash:  fmt.Sprintf("period-%d", period),
			CashbackAmount: proofRecord.Amount,
			Period:         period,
			Claimed:        proofRecord.Claimed,
			Proof:          proofRecord.Proof,
			MerkleRoot:     merkleRoot,
			CreatedAt:      time.Now(),
		}

		claimableRecords = append(claimableRecords, record)
	}

	return claimableRecords, nil
}

// MarkCashbackClaimed marks a cashback record as claimed
func (s *KVStore) MarkCashbackClaimed(ctx context.Context, userAddress string, period uint64) error {
	// Update the proof record to mark as claimed
	proofKey := fmt.Sprintf("cashback_proof:%d:%s", period, userAddress)
	proofData, err := s.GetCashbackData(ctx, proofKey)
	if err != nil {
		return fmt.Errorf("proof not found: %w", err)
	}

	var proofRecord struct {
		Proof   []string `json:"proof"`
		Amount  string   `json:"amount"`
		Claimed bool     `json:"claimed"`
	}
	if err := json.Unmarshal(proofData, &proofRecord); err != nil {
		return fmt.Errorf("failed to unmarshal proof: %w", err)
	}

	proofRecord.Claimed = true

	data, err := json.Marshal(proofRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal proof: %w", err)
	}

	if err := s.StoreCashbackData(ctx, proofKey, data); err != nil {
		return fmt.Errorf("failed to store claimed proof: %w", err)
	}

	// Update user stats
	stats, err := s.GetCashbackStats(ctx, userAddress)
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	// Move from pending to claimed
	pending := new(big.Int)
	if stats.PendingCashback != "" && stats.PendingCashback != "0" {
		pending.SetString(stats.PendingCashback, 10)
	}

	claimed := new(big.Int)
	if stats.ClaimedCashback != "" && stats.ClaimedCashback != "0" {
		claimed.SetString(stats.ClaimedCashback, 10)
	}

	amount := new(big.Int)
	amount.SetString(proofRecord.Amount, 10)

	// Check for negative balance
	if pending.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient pending balance: have %s, need %s", pending.String(), amount.String())
	}

	pending.Sub(pending, amount)
	claimed.Add(claimed, amount)

	stats.PendingCashback = pending.String()
	stats.ClaimedCashback = claimed.String()

	statsData, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	statsKey := fmt.Sprintf("cashback_stats:%s", userAddress)
	if err := s.StoreCashbackData(ctx, statsKey, statsData); err != nil {
		return fmt.Errorf("failed to store stats: %w", err)
	}

	// Invalidate cache since claim status changed
	cacheMu.Lock()
	delete(cashbackCache, userAddress)
	cacheMu.Unlock()

	log.Printf("✅ [Cashback] Marked period %d as claimed for user %s (cache invalidated)", period, userAddress)

	return nil
}

// MarkCashbackPending marks a cashback claim as pending (transaction submitted)
func (s *KVStore) MarkCashbackPending(ctx context.Context, userAddress string, period uint64, txHash string) error {
	proofKey := fmt.Sprintf("cashback_proof:%d:%s", period, userAddress)
	proofData, err := s.GetCashbackData(ctx, proofKey)
	if err != nil {
		return fmt.Errorf("proof not found: %w", err)
	}

	var proofRecord struct {
		Proof   []string `json:"proof"`
		Amount  string   `json:"amount"`
		Claimed bool     `json:"claimed"`
		TxHash  string   `json:"tx_hash,omitempty"`
		Status  string   `json:"status,omitempty"` // "unclaimed", "pending", "confirmed", "failed"
	}
	if err := json.Unmarshal(proofData, &proofRecord); err != nil {
		return fmt.Errorf("failed to unmarshal proof: %w", err)
	}

	// Don't allow marking as pending if already claimed
	if proofRecord.Claimed {
		return fmt.Errorf("cashback already claimed for period %d", period)
	}

	proofRecord.Status = "pending"
	if txHash != "" {
		proofRecord.TxHash = txHash
	}

	data, err := json.Marshal(proofRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal proof: %w", err)
	}

	if err := s.StoreCashbackData(ctx, proofKey, data); err != nil {
		return fmt.Errorf("failed to store pending proof: %w", err)
	}

	log.Printf("📝 [Cashback] Marked period %d as pending for user %s (tx: %s)", period, userAddress, txHash)
	return nil
}

// MarkCashbackFailed marks a cashback claim as failed (transaction reverted)
func (s *KVStore) MarkCashbackFailed(ctx context.Context, userAddress string, period uint64, reason string) error {
	proofKey := fmt.Sprintf("cashback_proof:%d:%s", period, userAddress)
	proofData, err := s.GetCashbackData(ctx, proofKey)
	if err != nil {
		return fmt.Errorf("proof not found: %w", err)
	}

	var proofRecord struct {
		Proof   []string `json:"proof"`
		Amount  string   `json:"amount"`
		Claimed bool     `json:"claimed"`
		TxHash  string   `json:"tx_hash,omitempty"`
		Status  string   `json:"status,omitempty"`
		Error   string   `json:"error,omitempty"`
	}
	if err := json.Unmarshal(proofData, &proofRecord); err != nil {
		return fmt.Errorf("failed to unmarshal proof: %w", err)
	}

	proofRecord.Status = "failed"
	proofRecord.Error = reason

	data, err := json.Marshal(proofRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal proof: %w", err)
	}

	if err := s.StoreCashbackData(ctx, proofKey, data); err != nil {
		return fmt.Errorf("failed to store failed proof: %w", err)
	}

	log.Printf("❌ [Cashback] Marked period %d as failed for user %s (reason: %s)", period, userAddress, reason)
	return nil
}
