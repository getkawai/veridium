package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kawai-network/x/store"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <user_address>")
	}

	userAddress := os.Args[1]

	fmt.Println("🐛 Debug Cashback Query")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("User Address: %s\n", userAddress)
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	ctx := context.Background()

	// Get settled periods
	settledPeriods, err := kv.GetSettledCashbackPeriods(ctx)
	if err != nil {
		log.Fatal("Failed to get settled periods:", err)
	}

	fmt.Printf("✅ Settled periods: %v\n", settledPeriods)
	fmt.Println()

	// Manually check each period
	for _, period := range settledPeriods {
		fmt.Printf("🔍 Checking period %d...\n", period)

		// Check merkle root
		rootKey := fmt.Sprintf("cashback_period:%d:merkle_root", period)
		fmt.Printf("   Root key: %s\n", rootKey)
		rootData, err := kv.GetCashbackData(ctx, rootKey)
		if err != nil {
			fmt.Printf("   ❌ Merkle root not found: %v\n", err)
			continue
		}
		var merkleRoot string
		if err := json.Unmarshal(rootData, &merkleRoot); err != nil {
			fmt.Printf("   ❌ Failed to unmarshal root: %v\n", err)
			continue
		}
		fmt.Printf("   ✅ Merkle root: %s\n", merkleRoot)

		// Check proof
		proofKey := fmt.Sprintf("cashback_proof:%d:%s", period, userAddress)
		fmt.Printf("   Proof key: %s\n", proofKey)
		proofData, err := kv.GetCashbackData(ctx, proofKey)
		if err != nil {
			fmt.Printf("   ❌ Proof not found: %v\n", err)
			continue
		}

		var proofRecord struct {
			Proof   []string `json:"proof"`
			Amount  string   `json:"amount"`
			Claimed bool     `json:"claimed"`
		}
		if err := json.Unmarshal(proofData, &proofRecord); err != nil {
			fmt.Printf("   ❌ Failed to unmarshal proof: %v\n", err)
			continue
		}

		fmt.Printf("   ✅ Proof found!\n")
		fmt.Printf("      Amount: %s\n", proofRecord.Amount)
		fmt.Printf("      Claimed: %v\n", proofRecord.Claimed)
		fmt.Printf("      Proof length: %d\n", len(proofRecord.Proof))
		fmt.Println()
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
}
