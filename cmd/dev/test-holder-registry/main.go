package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	ctx := context.Background()

	// Initialize KV Store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatalf("Failed to connect to KV: %v", err)
	}
	fmt.Println("✓ Connected to Cloudflare KV")

	// Initialize Holder Registry
	holderRegistry := blockchain.NewHolderRegistry(kv)

	// Test 1: Register a test holder (admin address)
	testAddress := common.HexToAddress(constant.GetAdminAddress())
	fmt.Printf("\n📝 Registering test holder: %s\n", testAddress.Hex())

	err = holderRegistry.RegisterHolder(ctx, testAddress, "test")
	if err != nil {
		log.Fatalf("Failed to register holder: %v", err)
	}
	fmt.Println("✅ Holder registered successfully")

	// Test 2: Get holder count
	count, err := holderRegistry.GetHolderCount(ctx)
	if err != nil {
		log.Fatalf("Failed to get holder count: %v", err)
	}
	fmt.Printf("\n📊 Total registered holders: %d\n", count)

	// Test 3: List all holders
	holders, err := holderRegistry.GetAllHolders(ctx)
	if err != nil {
		log.Fatalf("Failed to list holders: %v", err)
	}

	fmt.Printf("\n📋 Registered Holders:\n")
	for i, addr := range holders {
		info, err := holderRegistry.GetHolderInfo(ctx, addr)
		if err != nil {
			fmt.Printf("  %d. %s (error getting info: %v)\n", i+1, addr.Hex(), err)
			continue
		}
		fmt.Printf("  %d. %s (source: %s, registered: %d)\n", i+1, addr.Hex(), info.Source, info.Registered)
	}

	// Test 4: Export holders as JSON
	json, err := holderRegistry.ExportHolders(ctx)
	if err != nil {
		log.Fatalf("Failed to export holders: %v", err)
	}
	fmt.Printf("\n📄 Holder Registry JSON:\n%s\n", json)

	// Test 5: Check if holder has KAWAI balance
	scanner, err := blockchain.NewHolderScanner()
	if err != nil {
		log.Fatalf("Failed to create scanner: %v", err)
	}

	balance, err := scanner.GetBalance(ctx, testAddress)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("\n💰 KAWAI Balance: %s\n", balance.String())

	if balance.Cmp(common.Big0) == 0 {
		fmt.Println("\n⚠️  WARNING: Test address has 0 KAWAI balance")
		fmt.Println("   Revenue settlement will skip this holder")
		fmt.Println("   To test properly, use an address with KAWAI tokens")
	}

	fmt.Println("\n✅ Holder registry test complete!")
	os.Exit(0)
}
