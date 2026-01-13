package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <period> <user_address>")
	}

	period := os.Args[1]
	userAddress := os.Args[2]

	fmt.Println("🔍 Inspecting Cashback Proof")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("Period: %s\n", period)
	fmt.Printf("User Address: %s\n", userAddress)
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	ctx := context.Background()

	// Try to get proof
	proofKey := fmt.Sprintf("cashback_proof:%s:%s", period, userAddress)
	fmt.Printf("🔑 Proof Key: %s\n", proofKey)
	fmt.Println()

	data, err := kv.GetCashbackData(ctx, proofKey)
	if err != nil {
		log.Printf("❌ Failed to get proof: %v", err)
		return
	}

	fmt.Println("✅ Proof found!")
	fmt.Println()
	fmt.Println("📄 Raw Data:")
	fmt.Println(string(data))
	fmt.Println()

	// Try to unmarshal
	var proofRecord struct {
		Proof   []string `json:"proof"`
		Amount  string   `json:"amount"`
		Claimed bool     `json:"claimed"`
	}
	if err := json.Unmarshal(data, &proofRecord); err != nil {
		log.Printf("❌ Failed to unmarshal: %v", err)
		return
	}

	fmt.Println("✅ Parsed Proof:")
	fmt.Printf("   Amount: %s\n", proofRecord.Amount)
	fmt.Printf("   Claimed: %v\n", proofRecord.Claimed)
	fmt.Printf("   Proof length: %d\n", len(proofRecord.Proof))
	for i, p := range proofRecord.Proof {
		fmt.Printf("   Proof[%d]: %s\n", i, p)
	}
}
