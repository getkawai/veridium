package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/kawai-network/x/store"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <address> <period_id>", os.Args[0])
	}

	address := os.Args[1]
	periodIDStr := os.Args[2]

	periodID, err := strconv.ParseInt(periodIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid period ID: %v", err)
	}

	if err := confirmClaim(address, periodID); err != nil {
		log.Fatalf("Failed to confirm claim: %v", err)
	}
}

func confirmClaim(address string, periodID int64) error {
	ctx := context.Background()

	fmt.Printf("🔄 Confirming Claim in KV Store\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Address:   %s\n", address)
	fmt.Printf("Period ID: %d\n", periodID)
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Confirm the claim
	err = kv.ConfirmClaim(ctx, address, periodID)
	if err != nil {
		return fmt.Errorf("failed to confirm claim: %w", err)
	}

	fmt.Printf("✅ Claim confirmed successfully!\n")
	fmt.Printf("   The reward should now show as 'confirmed' in the UI\n")
	fmt.Printf("   and be moved from 'pending' to transaction history\n")

	return nil
}
