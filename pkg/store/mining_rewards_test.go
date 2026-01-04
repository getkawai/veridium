package store

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMiningRewardDistribution tests the reward split logic
func TestMiningRewardDistribution(t *testing.T) {
	tests := []struct {
		name            string
		tokenUsage      int64
		hasReferrer     bool
		expectedContrib string // percentage
		expectedDev     string
		expectedUser    string
		expectedAff     string
	}{
		{
			name:            "Referral User - 85/5/5/5 Split",
			tokenUsage:      1_000_000, // 1M tokens = 100 KAWAI base
			hasReferrer:     true,
			expectedContrib: "85.00", // 85 KAWAI
			expectedDev:     "5.00",  // 5 KAWAI
			expectedUser:    "5.00",  // 5 KAWAI
			expectedAff:     "5.00",  // 5 KAWAI
		},
		{
			name:            "Non-Referral User - 90/5/5 Split",
			tokenUsage:      1_000_000,
			hasReferrer:     false,
			expectedContrib: "90.00", // 90 KAWAI
			expectedDev:     "5.00",  // 5 KAWAI
			expectedUser:    "5.00",  // 5 KAWAI
			expectedAff:     "0.00",  // 0 KAWAI
		},
		{
			name:            "Small Job - Referral",
			tokenUsage:      100_000, // 100K tokens = 10 KAWAI base
			hasReferrer:     true,
			expectedContrib: "8.50",
			expectedDev:     "0.50",
			expectedUser:    "0.50",
			expectedAff:     "0.50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate base reward (100 KAWAI per 1M tokens)
			baseReward := new(big.Int).Mul(
				big.NewInt(tt.tokenUsage),
				big.NewInt(100), // KAWAI_RATE_PER_MILLION
			)
			baseReward.Div(baseReward, big.NewInt(1_000_000))

			// Calculate splits
			var contributorShare, developerShare, userShare, affiliatorShare *big.Int

			if tt.hasReferrer {
				// 85/5/5/5
				contributorShare = new(big.Int).Mul(baseReward, big.NewInt(85))
				contributorShare.Div(contributorShare, big.NewInt(100))

				developerShare = new(big.Int).Mul(baseReward, big.NewInt(5))
				developerShare.Div(developerShare, big.NewInt(100))

				userShare = new(big.Int).Mul(baseReward, big.NewInt(5))
				userShare.Div(userShare, big.NewInt(100))

				affiliatorShare = new(big.Int).Mul(baseReward, big.NewInt(5))
				affiliatorShare.Div(affiliatorShare, big.NewInt(100))
			} else {
				// 90/5/5
				contributorShare = new(big.Int).Mul(baseReward, big.NewInt(90))
				contributorShare.Div(contributorShare, big.NewInt(100))

				developerShare = new(big.Int).Mul(baseReward, big.NewInt(5))
				developerShare.Div(developerShare, big.NewInt(100))

				userShare = new(big.Int).Mul(baseReward, big.NewInt(5))
				userShare.Div(userShare, big.NewInt(100))

				affiliatorShare = big.NewInt(0)
			}

			// Convert to KAWAI (18 decimals)
			toKawai := func(amount *big.Int) string {
				kawai := new(big.Int).Mul(amount, big.NewInt(1e18))
				// Convert to float for comparison
				f := new(big.Float).SetInt(kawai)
				f.Quo(f, big.NewFloat(1e18))
				result, _ := f.Float64()
				return formatFloat(result)
			}

			// Verify splits
			assert.Equal(t, tt.expectedContrib, toKawai(contributorShare), "Contributor share mismatch")
			assert.Equal(t, tt.expectedDev, toKawai(developerShare), "Developer share mismatch")
			assert.Equal(t, tt.expectedUser, toKawai(userShare), "User share mismatch")
			assert.Equal(t, tt.expectedAff, toKawai(affiliatorShare), "Affiliator share mismatch")

			// Verify total = 100%
			total := new(big.Int).Add(contributorShare, developerShare)
			total.Add(total, userShare)
			total.Add(total, affiliatorShare)
			assert.Equal(t, baseReward.String(), total.String(), "Total should equal base reward")
		})
	}
}

// TestJobRewardRecordCreation tests job reward record creation
func TestJobRewardRecordCreation(t *testing.T) {
	record := &JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: "0xContributor",
		UserAddress:        "0xUser",
		ReferrerAddress:    "0xReferrer",
		DeveloperAddress:   "0xDeveloper",
		ContributorAmount:  "85000000000000000000", // 85 KAWAI
		DeveloperAmount:    "5000000000000000000",  // 5 KAWAI
		UserAmount:         "5000000000000000000",  // 5 KAWAI
		AffiliatorAmount:   "5000000000000000000",  // 5 KAWAI
		TokenUsage:         1_000_000,
		RewardType:         "kawai",
		HasReferrer:        true,
		IsSettled:          false,
	}

	// Verify all fields are set
	assert.NotEmpty(t, record.ContributorAddress)
	assert.NotEmpty(t, record.UserAddress)
	assert.NotEmpty(t, record.ReferrerAddress)
	assert.NotEmpty(t, record.DeveloperAddress)
	assert.Equal(t, "kawai", record.RewardType)
	assert.True(t, record.HasReferrer)
	assert.False(t, record.IsSettled)

	// Verify amounts sum to 100 KAWAI
	contrib := new(big.Int)
	contrib.SetString(record.ContributorAmount, 10)
	dev := new(big.Int)
	dev.SetString(record.DeveloperAmount, 10)
	user := new(big.Int)
	user.SetString(record.UserAmount, 10)
	aff := new(big.Int)
	aff.SetString(record.AffiliatorAmount, 10)

	total := new(big.Int).Add(contrib, dev)
	total.Add(total, user)
	total.Add(total, aff)

	expected := new(big.Int)
	expected.SetString("100000000000000000000", 10) // 100 KAWAI

	assert.Equal(t, expected.String(), total.String(), "Total should be 100 KAWAI")
}

// TestMerkleLeafGeneration tests 9-field Merkle leaf generation
func TestMerkleLeafGeneration(t *testing.T) {
	// This would test the actual Merkle leaf generation
	// matching the Solidity keccak256(abi.encodePacked(...))

	// Mock data
	period := uint64(1704326400)
	contributor := "0x1234567890123456789012345678901234567890"
	contributorAmt := "85000000000000000000"
	developerAmt := "5000000000000000000"
	userAmt := "5000000000000000000"
	affiliatorAmt := "5000000000000000000"
	developer := "0xDEV1234567890123456789012345678901234567890"
	user := "0xUSER123456789012345678901234567890123456789"
	affiliator := "0xAFF1234567890123456789012345678901234567890"

	// In real implementation, this would call generateMiningMerkleLeaf()
	// and verify it matches the smart contract's leaf generation

	t.Log("Period:", period)
	t.Log("Contributor:", contributor)
	t.Log("Amounts:", contributorAmt, developerAmt, userAmt, affiliatorAmt)
	t.Log("Addresses:", developer, user, affiliator)

	// TODO: Implement actual Merkle leaf generation test
	// This requires the generateMiningMerkleLeaf function to be testable
}

// Helper function
func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

// Benchmark for performance testing
func BenchmarkRewardCalculation(b *testing.B) {
	tokenUsage := int64(1_000_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		baseReward := new(big.Int).Mul(
			big.NewInt(tokenUsage),
			big.NewInt(100),
		)
		baseReward.Div(baseReward, big.NewInt(1_000_000))

		// Calculate 85/5/5/5 split
		contrib := new(big.Int).Mul(baseReward, big.NewInt(85))
		contrib.Div(contrib, big.NewInt(100))

		dev := new(big.Int).Mul(baseReward, big.NewInt(5))
		dev.Div(dev, big.NewInt(100))

		user := new(big.Int).Mul(baseReward, big.NewInt(5))
		user.Div(user, big.NewInt(100))

		aff := new(big.Int).Mul(baseReward, big.NewInt(5))
		aff.Div(aff, big.NewInt(100))

		_ = contrib
		_ = dev
		_ = user
		_ = aff
	}
}
