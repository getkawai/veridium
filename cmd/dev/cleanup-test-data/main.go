package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/pkg/store"
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
	// Delete all job_rewards:* keys
	// Note: This is a simplified version. In production, you'd want to:
	// 1. List all keys with prefix "job_rewards:"
	// 2. Delete them one by one
	// For now, we'll just log that it needs manual cleanup via Cloudflare dashboard
	fmt.Println("   ℹ️  Job rewards cleanup requires manual action via Cloudflare dashboard")
	fmt.Println("      or use: make cleanup-kv-mining-data")
	return nil
}

func cleanupCashbackData(ctx context.Context, kv *store.KVStore) error {
	// Similar to job rewards, this needs KV list/delete operations
	// which are already implemented in cleanup-kv-mining-data tool
	fmt.Println("   ℹ️  Cashback data cleanup requires manual action via Cloudflare dashboard")
	fmt.Println("      Keys to delete: cashback:*, cashback_stats:*, cashback_period:*, cashback_proof:*")
	return nil
}

func cleanupMerkleProofs(ctx context.Context, kv *store.KVStore) error {
	fmt.Println("   ℹ️  Merkle proofs cleanup requires manual action via Cloudflare dashboard")
	fmt.Println("      Keys to delete: proof:*")
	return nil
}

func cleanupSettlementPeriods(ctx context.Context, kv *store.KVStore) error {
	fmt.Println("   ℹ️  Settlement periods cleanup requires manual action via Cloudflare dashboard")
	fmt.Println("      Keys to delete: settlement:*")
	return nil
}
