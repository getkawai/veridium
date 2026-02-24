package services

import (
	"context"

	"github.com/kawai-network/x/store"
)

// ReferralService provides referral system functionality for the Wails desktop app
type ReferralService struct {
	kvStore *store.KVStore
}

// NewReferralService creates a new referral service
func NewReferralService(kvStore *store.KVStore) *ReferralService {
	return &ReferralService{
		kvStore: kvStore,
	}
}

// ReferralStats represents referral statistics
type ReferralStats struct {
	Code                 string  `json:"code"`
	Total_Referrals      int     `json:"total_referrals"`
	Total_Earnings_USDT  float64 `json:"total_earnings_usdt"`
	Total_Earnings_Kawai string  `json:"total_earnings_kawai"`
}

// TrialClaim represents the result of claiming trial
type TrialClaim struct {
	Balance_Added_USDT   float64 `json:"balance_added_usdt"`
	Balance_Added_Kawai  string  `json:"balance_added_kawai"`
	Has_Referral         bool    `json:"has_referral"`
	Referral_Bonus_USDT  float64 `json:"referral_bonus_usdt,omitempty"`
	Referral_Bonus_Kawai string  `json:"referral_bonus_kawai,omitempty"`
}

// CreateReferralCode creates a new referral code for the user
// This is a Wails-exposed method for the frontend
func (s *ReferralService) CreateReferralCode(userAddress string) (string, error) {
	ctx := context.Background()

	referralData, err := s.kvStore.CreateReferralCode(ctx, userAddress)
	if err != nil {
		return "", err
	}

	return referralData.Code, nil
}

// GetReferralStats retrieves referral statistics for the user
// This is a Wails-exposed method for the frontend
func (s *ReferralService) GetReferralStats(userAddress string) (*ReferralStats, error) {
	ctx := context.Background()

	referralData, err := s.kvStore.GetReferralCodeByAddress(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// Convert micro USDT to USDT
	totalUSDT := float64(referralData.TotalEarningsUSDT) / 1_000_000

	return &ReferralStats{
		Code:                 referralData.Code,
		Total_Referrals:      referralData.TotalReferrals,
		Total_Earnings_USDT:  totalUSDT,
		Total_Earnings_Kawai: referralData.TotalEarningsKawai,
	}, nil
}

// ClaimFreeTrialWithReferral claims free trial with optional referral code
// This is a Wails-exposed method for the frontend
func (s *ReferralService) ClaimFreeTrialWithReferral(
	address string,
	machineID string,
	referralCode string,
) (*TrialClaim, error) {
	ctx := context.Background()

	usdtAmount, kawaiAmount, err := s.kvStore.ClaimFreeTrialWithReferral(ctx, address, machineID, referralCode)
	if err != nil {
		return nil, err
	}

	// Convert micro USDT to USDT
	usdtFloat := float64(usdtAmount) / 1_000_000

	hasReferral := referralCode != ""
	response := &TrialClaim{
		Balance_Added_USDT:  usdtFloat,
		Balance_Added_Kawai: kawaiAmount,
		Has_Referral:        hasReferral,
	}

	// Add referral bonus info if applicable
	if hasReferral {
		response.Referral_Bonus_USDT = 10.0                     // 10 USDT for new user with referral
		response.Referral_Bonus_Kawai = "200000000000000000000" // 200 KAWAI
	}

	return response, nil
}

// GetReferralBonusAmounts returns the bonus amounts for referral system
// This is a helper method for the frontend to display bonus information
func (s *ReferralService) GetReferralBonusAmounts() map[string]interface{} {
	return map[string]interface{}{
		"base_trial_usdt":       5.0,                     // 5 USDT without referral
		"base_trial_kawai":      "100000000000000000000", // 100 KAWAI
		"referral_trial_usdt":   10.0,                    // 10 USDT with referral
		"referral_trial_kawai":  "200000000000000000000", // 200 KAWAI
		"referrer_reward_usdt":  5.0,                     // 5 USDT for referrer
		"referrer_reward_kawai": "100000000000000000000", // 100 KAWAI
	}
}
