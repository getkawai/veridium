package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/joho/godotenv"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/store"
)

// This tool injects test mining reward data into KV store
// for testing without manual UI interaction

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	ctx := context.Background()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	fmt.Println("🧪 Injecting Test Mining Reward Data")
	fmt.Println("=====================================")
	fmt.Println("")

	// Test scenario 1: Referral user (85/5/5/5)
	fmt.Println("📝 Scenario 1: Referral User")
	if err := injectReferralUserReward(ctx, kv); err != nil {
		log.Fatal("Failed to inject referral user reward:", err)
	}
	fmt.Println("✅ Injected referral user reward")
	fmt.Println("")

	// Test scenario 2: Non-referral user (90/5/5)
	fmt.Println("📝 Scenario 2: Non-Referral User")
	if err := injectNonReferralUserReward(ctx, kv); err != nil {
		log.Fatal("Failed to inject non-referral user reward:", err)
	}
	fmt.Println("✅ Injected non-referral user reward")
	fmt.Println("")

	// Test scenario 3: Multiple jobs from same user
	fmt.Println("📝 Scenario 3: Multiple Jobs (Same User)")
	if err := injectMultipleJobs(ctx, kv); err != nil {
		log.Fatal("Failed to inject multiple jobs:", err)
	}
	fmt.Println("✅ Injected multiple jobs")
	fmt.Println("")

	fmt.Println("=====================================")
	fmt.Println("🎉 Test Data Injection Complete!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("  1. Run: mining-settlement generate --reward-type kawai")
	fmt.Println("  2. Check generated Merkle proofs")
	fmt.Println("  3. Upload Merkle root to contract")
	fmt.Println("  4. Test claim in UI")
}

func injectReferralUserReward(ctx context.Context, kv *store.KVStore) error {
	// Calculate 85/5/5/5 split for 1M tokens (100 KAWAI base)
	baseReward := big.NewInt(100) // 100 KAWAI

	contributorAmt := new(big.Int).Mul(baseReward, big.NewInt(85))
	contributorAmt.Div(contributorAmt, big.NewInt(100))
	contributorAmt.Mul(contributorAmt, big.NewInt(1e18)) // 85 KAWAI

	developerAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
	developerAmt.Div(developerAmt, big.NewInt(100))
	developerAmt.Mul(developerAmt, big.NewInt(1e18)) // 5 KAWAI

	userAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
	userAmt.Div(userAmt, big.NewInt(100))
	userAmt.Mul(userAmt, big.NewInt(1e18)) // 5 KAWAI

	affiliatorAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
	affiliatorAmt.Div(affiliatorAmt, big.NewInt(100))
	affiliatorAmt.Mul(affiliatorAmt, big.NewInt(1e18)) // 5 KAWAI

	record := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: "0xTestContributor1111111111111111111111111111",
		UserAddress:        "0xTestUser11111111111111111111111111111111",
		ReferrerAddress:    "0xTestReferrer111111111111111111111111111",
		DeveloperAddress:   constant.GetRandomTreasuryAddress(),
		ContributorAmount:  contributorAmt.String(),
		DeveloperAmount:    developerAmt.String(),
		UserAmount:         userAmt.String(),
		AffiliatorAmount:   affiliatorAmt.String(),
		TokenUsage:         1_000_000,
		RewardType:         "kawai",
		HasReferrer:        true,
		IsSettled:          false,
	}

	fmt.Printf("  Contributor: %s (85 KAWAI)\n", record.ContributorAddress[:20]+"...")
	fmt.Printf("  User: %s (5 KAWAI cashback)\n", record.UserAddress[:20]+"...")
	fmt.Printf("  Referrer: %s (5 KAWAI commission)\n", record.ReferrerAddress[:20]+"...")
	fmt.Printf("  Developer: %s (5 KAWAI)\n", record.DeveloperAddress[:20]+"...")

	return kv.SaveJobReward(ctx, record)
}

func injectNonReferralUserReward(ctx context.Context, kv *store.KVStore) error {
	// Calculate 90/5/5 split for 1M tokens (100 KAWAI base)
	baseReward := big.NewInt(100)

	contributorAmt := new(big.Int).Mul(baseReward, big.NewInt(90))
	contributorAmt.Div(contributorAmt, big.NewInt(100))
	contributorAmt.Mul(contributorAmt, big.NewInt(1e18)) // 90 KAWAI

	developerAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
	developerAmt.Div(developerAmt, big.NewInt(100))
	developerAmt.Mul(developerAmt, big.NewInt(1e18)) // 5 KAWAI

	userAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
	userAmt.Div(userAmt, big.NewInt(100))
	userAmt.Mul(userAmt, big.NewInt(1e18)) // 5 KAWAI

	record := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: "0xTestContributor2222222222222222222222222222",
		UserAddress:        "0xTestUser22222222222222222222222222222222",
		ReferrerAddress:    "", // No referrer
		DeveloperAddress:   constant.GetRandomTreasuryAddress(),
		ContributorAmount:  contributorAmt.String(),
		DeveloperAmount:    developerAmt.String(),
		UserAmount:         userAmt.String(),
		AffiliatorAmount:   "0", // No affiliator
		TokenUsage:         1_000_000,
		RewardType:         "kawai",
		HasReferrer:        false,
		IsSettled:          false,
	}

	fmt.Printf("  Contributor: %s (90 KAWAI)\n", record.ContributorAddress[:20]+"...")
	fmt.Printf("  User: %s (5 KAWAI cashback)\n", record.UserAddress[:20]+"...")
	fmt.Printf("  Developer: %s (5 KAWAI)\n", record.DeveloperAddress[:20]+"...")

	return kv.SaveJobReward(ctx, record)
}

func injectMultipleJobs(ctx context.Context, kv *store.KVStore) error {
	// Inject 3 jobs from same contributor to test aggregation
	contributorAddr := "0xTestContributor3333333333333333333333333333"
	userAddr := "0xTestUser33333333333333333333333333333333"
	referrerAddr := "0xTestReferrer333333333333333333333333333"

	for i := 0; i < 3; i++ {
		baseReward := big.NewInt(50) // 50 KAWAI per job

		contributorAmt := new(big.Int).Mul(baseReward, big.NewInt(85))
		contributorAmt.Div(contributorAmt, big.NewInt(100))
		contributorAmt.Mul(contributorAmt, big.NewInt(1e18))

		developerAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
		developerAmt.Div(developerAmt, big.NewInt(100))
		developerAmt.Mul(developerAmt, big.NewInt(1e18))

		userAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
		userAmt.Div(userAmt, big.NewInt(100))
		userAmt.Mul(userAmt, big.NewInt(1e18))

		affiliatorAmt := new(big.Int).Mul(baseReward, big.NewInt(5))
		affiliatorAmt.Div(affiliatorAmt, big.NewInt(100))
		affiliatorAmt.Mul(affiliatorAmt, big.NewInt(1e18))

		record := &store.JobRewardRecord{
			Timestamp:          time.Now().Add(time.Duration(i) * time.Minute),
			ContributorAddress: contributorAddr,
			UserAddress:        userAddr,
			ReferrerAddress:    referrerAddr,
			DeveloperAddress:   constant.GetRandomTreasuryAddress(),
			ContributorAmount:  contributorAmt.String(),
			DeveloperAmount:    developerAmt.String(),
			UserAmount:         userAmt.String(),
			AffiliatorAmount:   affiliatorAmt.String(),
			TokenUsage:         500_000, // 500K tokens per job
			RewardType:         "kawai",
			HasReferrer:        true,
			IsSettled:          false,
		}

		if err := kv.SaveJobReward(ctx, record); err != nil {
			return err
		}

		fmt.Printf("  Job %d: 42.5 KAWAI (85/5/5/5 split)\n", i+1)
	}

	fmt.Printf("  Total for contributor: 127.5 KAWAI (3 jobs aggregated)\n")

	return nil
}

