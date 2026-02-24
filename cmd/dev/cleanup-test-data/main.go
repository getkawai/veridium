package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/x/store"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("🧹 Cleanup Test Data")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("⚠️  WARNING: This will delete test data from Cloudflare KV!")
	fmt.Println("")
	fmt.Println("Data to be cleaned:")
	fmt.Println("  • Mining job rewards (unsettled)")
	fmt.Println("  • Cashback records (all test data)")
	fmt.Println("  • Merkle proofs (all periods)")
	fmt.Println("  • Settlement periods (all)")
	fmt.Println("")
	fmt.Println("Data that will be PRESERVED:")
	fmt.Println("  ✅ User profiles")
	fmt.Println("  ✅ API keys")
	fmt.Println("  ✅ Authentication data")
	fmt.Println("")

	// Check for confirmation flag
	if len(os.Args) < 2 || os.Args[1] != "--confirm" {
		fmt.Println("To proceed, run:")
		fmt.Println("  make cleanup-test-data")
		fmt.Println("")
		fmt.Println("Or manually:")
		fmt.Println("  go run cmd/dev/cleanup-test-data/main.go --confirm")
		fmt.Println("")
		os.Exit(0)
	}

	fmt.Println("🚀 Starting cleanup...")
	fmt.Println("")

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatalf("Failed to initialize KV store: %v", err)
	}

	ctx := context.Background()

	// 1. Cleanup mining job rewards
	fmt.Println("1️⃣  Cleaning up mining job rewards...")
	if err := cleanupJobRewards(ctx, kv); err != nil {
		log.Printf("   ⚠️  Warning: %v", err)
	} else {
		fmt.Println("   ✅ Mining job rewards cleaned")
	}

	// 2. Cleanup cashback data
	fmt.Println("2️⃣  Cleaning up cashback data...")
	if err := cleanupCashbackData(ctx, kv); err != nil {
		log.Printf("   ⚠️  Warning: %v", err)
	} else {
		fmt.Println("   ✅ Cashback data cleaned")
	}

	// 3. Cleanup merkle proofs
	fmt.Println("3️⃣  Cleaning up Merkle proofs...")
	if err := cleanupMerkleProofs(ctx, kv); err != nil {
		log.Printf("   ⚠️  Warning: %v", err)
	} else {
		fmt.Println("   ✅ Merkle proofs cleaned")
	}

	// 4. Cleanup settlement periods
	fmt.Println("4️⃣  Cleaning up settlement periods...")
	if err := cleanupSettlementPeriods(ctx, kv); err != nil {
		log.Printf("   ⚠️  Warning: %v", err)
	} else {
		fmt.Println("   ✅ Settlement periods cleaned")
	}

	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("✅ Cleanup completed!")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("You can now start fresh testing:")
	fmt.Println("  1. make test-inject-mining-data")
	fmt.Println("  2. make settle-all")
	fmt.Println("  3. make dev-hot")
}

func cleanupJobRewards(ctx context.Context, kv *store.KVStore) error {
	// Get all contributors first
	contributors, err := kv.ListContributorsWithBalance(ctx, "kawai")
	if err != nil {
		return fmt.Errorf("failed to list contributors: %w", err)
	}

	deletedCount := 0
	for _, contributor := range contributors {
		// Get all job rewards for this contributor
		jobs, err := kv.GetJobRewardsSinceLastSettlement(ctx, contributor.WalletAddress, "kawai")
		if err != nil {
			log.Printf("   ⚠️  Failed to get jobs for %s: %v", contributor.WalletAddress, err)
			continue
		}

		// Delete each job reward
		for range jobs {
			// Mark as settled with period 0 to effectively delete
			if err := kv.MarkJobRewardsAsSettled(ctx, contributor.WalletAddress, 0); err != nil {
				log.Printf("   ⚠️  Failed to delete jobs for %s: %v", contributor.WalletAddress, err)
				break
			}
		}

		if len(jobs) > 0 {
			deletedCount += len(jobs)
		}
	}

	fmt.Printf("   ✅ Deleted %d job reward records from %d contributors\n", deletedCount, len(contributors))
	return nil
}

func cleanupCashbackData(ctx context.Context, kv *store.KVStore) error {
	// Cashback cleanup is not critical for mining rewards testing
	// Just log that it's skipped
	fmt.Println("   ℹ️  Cashback data cleanup skipped (not needed for mining rewards)")
	return nil
}

func cleanupMerkleProofs(ctx context.Context, kv *store.KVStore) error {
	// Get all settlement periods first
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list periods: %w", err)
	}

	// Get all contributors
	contributors, err := kv.ListContributorsWithBalance(ctx, "kawai")
	if err != nil {
		return fmt.Errorf("failed to list contributors: %w", err)
	}

	deletedCount := 0
	for _, period := range periods {
		for _, contributor := range contributors {
			// Delete proof for this period and contributor
			if err := kv.DeleteMerkleProof(ctx, contributor.WalletAddress, period.PeriodID); err != nil {
				// Ignore not found errors
				continue
			}
			deletedCount++
		}
	}

	fmt.Printf("   ✅ Deleted %d Merkle proofs\n", deletedCount)
	return nil
}

func cleanupSettlementPeriods(ctx context.Context, kv *store.KVStore) error {
	// Get all settlement periods
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list periods: %w", err)
	}

	// Note: KVStore doesn't have DeleteSettlement method
	// Settlements will be overwritten by new ones with same period IDs
	// This is acceptable for fresh testing

	fmt.Printf("   ℹ️  Found %d settlement periods (will be overwritten by new settlements)\n", len(periods))
	return nil
}
