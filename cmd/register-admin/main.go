package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/x/constant"
)

func main() {
	// CLI flags
	dryRun := flag.Bool("dry-run", false, "Show what would be registered without actually doing it")
	specificAddress := flag.String("address", "", "Register a specific address instead of all treasury addresses")
	flag.Parse()

	ctx := context.Background()

	// Initialize KV Store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		slog.Error("Failed to connect to KV", "error", err)
		os.Exit(1)
	}
	slog.Info("✓ Connected to Cloudflare KV")

	var addresses []string
	if *specificAddress != "" {
		// Validate and normalize the address
		if !common.IsHexAddress(*specificAddress) {
			slog.Error("Invalid Ethereum address", "address", *specificAddress)
			os.Exit(1)
		}
		addresses = []string{common.HexToAddress(*specificAddress).Hex()}
	} else {
		// Get all treasury addresses
		addresses = constant.GetTreasuryAddresses()
	}

	if len(addresses) == 0 {
		slog.Error("No addresses to register")
		os.Exit(1)
	}

	slog.Info("Admin addresses to register", "count", len(addresses))

	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, addr := range addresses {
		// Normalize address
		normalizedAddr := common.HexToAddress(addr).Hex()

		fmt.Printf("\n[%d/%d] Processing: %s\n", i+1, len(addresses), normalizedAddr)

		// Check if already exists
		existing, err := kv.GetContributor(ctx, normalizedAddr)
		if err == nil && existing != nil {
			if existing.IsAdmin {
				fmt.Printf("  ⏭️  Already registered as admin\n")
				skipCount++
				continue
			} else {
				fmt.Printf("  ⚠️  Exists as regular contributor, will update to admin\n")
			}
		}

		if *dryRun {
			fmt.Printf("  [DRY-RUN] Would register/update as admin\n")
			successCount++
			continue
		}

		// Register or update as admin
		contributorData := &store.ContributorData{
			WalletAddress:      normalizedAddr,
			EndpointURL:        "", // Admin doesn't need endpoint
			HardwareSpecs:      "Admin Account",
			RegisteredAt:       time.Now(),
			LastSeen:           time.Now(),
			Status:             store.ContributorStatusAdmin,
			AccumulatedRewards: "0",
			AccumulatedUSDT:    "0",
			IsActive:           true,
			IsAdmin:            true,
		}

		// If exists, preserve existing balances
		if existing != nil {
			contributorData.AccumulatedRewards = existing.AccumulatedRewards
			contributorData.AccumulatedUSDT = existing.AccumulatedUSDT
			contributorData.RegisteredAt = existing.RegisteredAt
		}

		err = kv.SaveContributor(ctx, contributorData)
		if err != nil {
			fmt.Printf("  ❌ Failed: %v\n", err)
			errorCount++
			continue
		}

		fmt.Printf("  ✅ Successfully registered/updated as admin\n")
		successCount++
	}

	// Summary
	separator := "=================================================="
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("Registration Summary:\n")
	fmt.Printf("  Total addresses: %d\n", len(addresses))
	fmt.Printf("  ✅ Success: %d\n", successCount)
	fmt.Printf("  ⏭️  Skipped (already admin): %d\n", skipCount)
	fmt.Printf("  ❌ Failed: %d\n", errorCount)
	fmt.Printf("%s\n", separator)

	if *dryRun {
		fmt.Printf("\n⚠️  This was a DRY RUN. No changes were made.\n")
		fmt.Printf("Run without --dry-run to actually register the addresses.\n")
	}

	if errorCount > 0 {
		os.Exit(1)
	}
}
