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
	fmt.Println("🧹 Cleanup Old Merkle Proofs")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("⚠️  WARNING: This will delete OLD Merkle proofs from Cloudflare KV!")
	fmt.Println("")
	fmt.Println("This will keep ONLY the 2 most recent periods:")
	fmt.Println("  ✅ Period 1768232683 (fresh - just uploaded)")
	fmt.Println("  ✅ Period 1768169460 (from earlier today)")
	fmt.Println("")
	fmt.Println("This will DELETE all older periods:")
	fmt.Println("  ❌ Period 1768139780 and older (8 periods)")
	fmt.Println("")

	// Check for confirmation flag
	if len(os.Args) < 2 || os.Args[1] != "--confirm" {
		fmt.Println("To proceed, run:")
		fmt.Println("  go run cmd/dev/cleanup-old-proofs/main.go --confirm")
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

	// Test address
	testAddress := "0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"

	// Old periods to delete (keep only 1768232683 and 1768169460)
	oldPeriods := []int64{
		1768139780,
		1768137242,
		1768137123,
		1768136095,
		1768135359,
		1768130418,
		1767650263,
		1767557168,
	}

	deleted := 0
	failed := 0

	for _, periodID := range oldPeriods {
		fmt.Printf("   Deleting proof for period %d... ", periodID)
		err := kv.DeleteMerkleProof(ctx, testAddress, periodID)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			failed++
		} else {
			fmt.Printf("✅ Deleted\n")
			deleted++
		}
	}

	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("✅ Cleanup completed!\n")
	fmt.Printf("   Deleted: %d proofs\n", deleted)
	fmt.Printf("   Failed:  %d proofs\n", failed)
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("Remaining proofs:")
	fmt.Println("  ✅ Period 1768232683: 126 KAWAI (fresh)")
	fmt.Println("  ✅ Period 1768169460: 126 KAWAI (earlier)")
	fmt.Println("")
	fmt.Println("Total claimable: 252 KAWAI")
	fmt.Println("")
	fmt.Println("Now refresh the UI to see clean data!")
}
