package store

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// =============================================================================
// REFERRAL SYSTEM
// =============================================================================

const (
	// Bonus amounts (in micro USDT)
	BaseTrialAmount     = 5_000_000  // 5 USDT without referral
	ReferralTrialAmount = 10_000_000 // 10 USDT with referral
	ReferrerReward      = 5_000_000  // 5 USDT for referrer

	// KAWAI Token Bonus (in wei - 18 decimals)
	BaseTrialKawai      = 100_000_000_000_000_000_000 // 100 KAWAI without referral
	ReferralTrialKawai  = 200_000_000_000_000_000_000 // 200 KAWAI with referral
	ReferrerKawaiReward = 100_000_000_000_000_000_000 // 100 KAWAI for referrer
)

// ReferralData stores referral information
type ReferralData struct {
	Code               string    `json:"code"`                 // Unique 6-char code
	OwnerAddress       string    `json:"owner_address"`        // Who owns this code
	TotalReferrals     int       `json:"total_referrals"`      // Count of successful referrals
	TotalEarningsUSDT  int64     `json:"total_earnings_usdt"`  // Total USDT earned (micro)
	TotalEarningsKawai string    `json:"total_earnings_kawai"` // Total KAWAI earned (wei as string)
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ReferralClaim tracks individual referral claims
type ReferralClaim struct {
	ReferralCode        string     `json:"referral_code"`
	ReferredUser        string     `json:"referred_user"`         // Address of new user
	ReferrerRewardUSDT  int64      `json:"referrer_reward_usdt"`  // USDT earned by referrer (micro)
	ReferrerRewardKawai string     `json:"referrer_reward_kawai"` // KAWAI earned by referrer (wei)
	NewUserBonusUSDT    int64      `json:"new_user_bonus_usdt"`   // USDT received by new user
	NewUserBonusKawai   string     `json:"new_user_bonus_kawai"`  // KAWAI received by new user
	Status              string     `json:"status"`                // "pending", "completed"
	CreatedAt           time.Time  `json:"created_at"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
}

// GenerateReferralCode creates a unique 6-character alphanumeric code
func GenerateReferralCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Exclude confusing chars (0,O,1,I)
	const length = 6

	code := make([]byte, length)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}

	return string(code), nil
}

// CreateReferralCode creates a new referral code for a user
func (s *KVStore) CreateReferralCode(ctx context.Context, address string) (*ReferralData, error) {
	// Check if user already has a code
	existingCode, err := s.GetReferralCodeByAddress(ctx, address)
	if err == nil && existingCode != nil {
		return existingCode, nil // Already has a code
	}

	// Generate unique code
	var code string
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		code, err = GenerateReferralCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}

		// Check if code already exists
		existing, _ := s.GetReferralData(ctx, code)
		if existing == nil {
			break // Code is unique
		}
	}

	// Create referral data
	referralData := &ReferralData{
		Code:               code,
		OwnerAddress:       address,
		TotalReferrals:     0,
		TotalEarningsUSDT:  0,
		TotalEarningsKawai: "0",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Save to KV
	data, err := json.Marshal(referralData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal referral data: %w", err)
	}

	// Store by code
	keyByCode := fmt.Sprintf("referral:code:%s", code)
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.usersNamespaceID,
		Key:         keyByCode,
		Value:       data,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save referral data: %w", err)
	}

	// Store by address (for lookup)
	keyByAddress := fmt.Sprintf("referral:address:%s", strings.ToLower(address))
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.usersNamespaceID,
		Key:         keyByAddress,
		Value:       []byte(code),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save address mapping: %w", err)
	}

	return referralData, nil
}

// GetReferralData retrieves referral data by code
func (s *KVStore) GetReferralData(ctx context.Context, code string) (*ReferralData, error) {
	key := fmt.Sprintf("referral:code:%s", strings.ToUpper(code))

	data, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.usersNamespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("referral code not found: %w", err)
	}

	var referralData ReferralData
	if err := json.Unmarshal(data, &referralData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal referral data: %w", err)
	}

	return &referralData, nil
}

// GetReferralCodeByAddress retrieves user's referral code by their address
func (s *KVStore) GetReferralCodeByAddress(ctx context.Context, address string) (*ReferralData, error) {
	// Get code from address mapping
	keyByAddress := fmt.Sprintf("referral:address:%s", strings.ToLower(address))

	codeBytes, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.usersNamespaceID,
		Key:         keyByAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("no referral code for this address: %w", err)
	}

	code := string(codeBytes)
	return s.GetReferralData(ctx, code)
}

// ClaimFreeTrialWithReferral claims free trial with optional referral code
// Returns: (usdtAmount, kawaiAmount, error)
func (s *KVStore) ClaimFreeTrialWithReferral(ctx context.Context, address string, machineID string, referralCode string) (int64, string, error) {
	// Determine bonus amounts
	usdtBonus := int64(BaseTrialAmount)
	kawaiBonus := "100000000000000000000" // 100 KAWAI (as string to avoid overflow)
	hasReferral := referralCode != ""

	if hasReferral {
		// Validate referral code
		referralData, err := s.GetReferralData(ctx, referralCode)
		if err != nil {
			return 0, "0", fmt.Errorf("invalid referral code: %w", err)
		}

		// Prevent self-referral
		if strings.EqualFold(referralData.OwnerAddress, address) {
			return 0, "0", fmt.Errorf("cannot use your own referral code")
		}

		usdtBonus = ReferralTrialAmount
		kawaiBonus = "200000000000000000000" // 200 KAWAI
	}

	// Claim trial (USDT + KAWAI)
	err := s.claimTrialWithDualReward(ctx, address, machineID, usdtBonus, kawaiBonus)
	if err != nil {
		return 0, "0", err
	}

	// If referral was used, reward the referrer
	if hasReferral {
		referrerKawai := "100000000000000000000" // 100 KAWAI for referrer
		if err := s.rewardReferrer(ctx, referralCode, address, usdtBonus, referrerKawai); err != nil {
			// Log error but don't fail the claim
			fmt.Printf("Warning: Failed to reward referrer: %v\n", err)
		}
	}

	return usdtBonus, kawaiBonus, nil
}

// claimTrialWithDualReward is internal method to claim trial with USDT + KAWAI
func (s *KVStore) claimTrialWithDualReward(ctx context.Context, address string, machineID string, usdtAmount int64, kawaiAmount string) error {
	// Pre-check Machine ID
	if machineID != "" {
		keyMachine := fmt.Sprintf("trial_machine:%s", machineID)
		valMachine, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.usersNamespaceID,
			Key:         keyMachine,
		})
		if err == nil && string(valMachine) == "true" {
			return fmt.Errorf("free trial already claimed by this device")
		}
	}

	// Atomic Read-Modify-Write Loop
	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		currentData, err := s.GetUserBalance(ctx, address)
		if err != nil {
			return err
		}

		if currentData.TrialClaimed {
			return fmt.Errorf("free trial already claimed by this address")
		}

		// Update USDT balance
		currentUSDTBalance := new(big.Int)
		currentUSDTBalance.SetString(currentData.USDTBalance, 10)
		newUSDTBalance := new(big.Int).Add(currentUSDTBalance, big.NewInt(usdtAmount))
		currentData.USDTBalance = newUSDTBalance.String()

		// Update KAWAI balance
		currentKawaiBalance := new(big.Int)
		if currentData.KawaiBalance != "" {
			currentKawaiBalance.SetString(currentData.KawaiBalance, 10)
		}
		kawaiToAdd := new(big.Int)
		kawaiToAdd.SetString(kawaiAmount, 10)
		newKawaiBalance := new(big.Int).Add(currentKawaiBalance, kawaiToAdd)
		currentData.KawaiBalance = newKawaiBalance.String()

		currentData.TrialClaimed = true

		data, err := json.Marshal(currentData)
		if err != nil {
			return fmt.Errorf("failed to marshal balance data: %w", err)
		}

		key := fmt.Sprintf("balance:%s", address)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         key,
			Value:       data,
		})

		if err == nil {
			// Mark machine ID
			if machineID != "" {
				keyMachine := fmt.Sprintf("trial_machine:%s", machineID)
				_, _ = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
					NamespaceID: s.usersNamespaceID,
					Key:         keyMachine,
					Value:       []byte("true"),
				})
			}
			return nil
		}

		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	return fmt.Errorf("failed to claim trial after %d retries", maxRetries)
}

// rewardReferrer adds USDT + KAWAI reward to referrer's balance
func (s *KVStore) rewardReferrer(ctx context.Context, referralCode string, referredUser string, newUserUSDT int64, newUserKawai string) error {
	// Get referral data
	referralData, err := s.GetReferralData(ctx, referralCode)
	if err != nil {
		return err
	}

	// Add reward to referrer's balance
	referrerAddress := referralData.OwnerAddress

	// Atomic balance update
	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		currentData, err := s.GetUserBalance(ctx, referrerAddress)
		if err != nil {
			return err
		}

		// Update USDT balance
		currentUSDTBalance := new(big.Int)
		currentUSDTBalance.SetString(currentData.USDTBalance, 10)
		newUSDTBalance := new(big.Int).Add(currentUSDTBalance, big.NewInt(ReferrerReward))
		currentData.USDTBalance = newUSDTBalance.String()

		// Update KAWAI balance
		currentKawaiBalance := new(big.Int)
		if currentData.KawaiBalance != "" {
			currentKawaiBalance.SetString(currentData.KawaiBalance, 10)
		}
		referrerKawai := new(big.Int)
		referrerKawai.SetString("100000000000000000000", 10) // 100 KAWAI
		newKawaiBalance := new(big.Int).Add(currentKawaiBalance, referrerKawai)
		currentData.KawaiBalance = newKawaiBalance.String()

		data, err := json.Marshal(currentData)
		if err != nil {
			return fmt.Errorf("failed to marshal balance data: %w", err)
		}

		key := fmt.Sprintf("balance:%s", referrerAddress)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         key,
			Value:       data,
		})

		if err != nil {
			if attempt < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return fmt.Errorf("failed to update referrer balance after %d retries: %w", maxRetries, err)
		}

		// Update referral stats
		referralData.TotalReferrals++
		referralData.TotalEarningsUSDT += ReferrerReward

		// Update KAWAI earnings
		currentKawaiEarnings := new(big.Int)
		if referralData.TotalEarningsKawai != "" {
			currentKawaiEarnings.SetString(referralData.TotalEarningsKawai, 10)
		}
		referrerKawaiEarnings := new(big.Int)
		referrerKawaiEarnings.SetString("100000000000000000000", 10) // 100 KAWAI
		newKawaiEarnings := new(big.Int).Add(currentKawaiEarnings, referrerKawaiEarnings)
		referralData.TotalEarningsKawai = newKawaiEarnings.String()
		referralData.UpdatedAt = time.Now()

		statsData, err2 := json.Marshal(referralData)
		if err2 != nil {
			return fmt.Errorf("failed to marshal referral stats: %w", err2)
		}

		keyByCode := fmt.Sprintf("referral:code:%s", referralCode)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         keyByCode,
			Value:       statsData,
		})
		if err != nil {
			return fmt.Errorf("failed to update referral stats: %w", err)
		}

		// Record claim
		now := time.Now()
		claim := ReferralClaim{
			ReferralCode:        referralCode,
			ReferredUser:        referredUser,
			ReferrerRewardUSDT:  ReferrerReward,
			ReferrerRewardKawai: "100000000000000000000", // 100 KAWAI
			NewUserBonusUSDT:    newUserUSDT,
			NewUserBonusKawai:   newUserKawai,
			Status:              "completed",
			CreatedAt:           time.Now(),
			CompletedAt:         &now,
		}

		claimData, err := json.Marshal(claim)
		if err != nil {
			return fmt.Errorf("failed to marshal claim data: %w", err)
		}

		claimKey := fmt.Sprintf("referral:claim:%s:%s", referralCode, referredUser)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         claimKey,
			Value:       claimData,
		})
		if err != nil {
			return fmt.Errorf("failed to record claim: %w", err)
		}

		// Success - exit loop
		return nil
	}

	return fmt.Errorf("failed to reward referrer after %d retries", maxRetries)
}
