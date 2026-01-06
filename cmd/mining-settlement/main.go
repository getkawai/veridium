package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	// Subcommands
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	// Flags for generate command
	var rewardType string
	generateCmd.StringVar(&rewardType, "type", "kawai", "Reward type: kawai or usdt")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Initialize KV Store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatalf("Failed to initialize KV store: %v", err)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "generate":
		generateCmd.Parse(os.Args[2:])
		if err := generateMerkleTree(ctx, kv, rewardType); err != nil {
			log.Fatalf("Generate failed: %v", err)
		}

	case "upload":
		uploadCmd.Parse(os.Args[2:])
		if err := uploadMerkleRoot(ctx, kv); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}

	case "status":
		statusCmd.Parse(os.Args[2:])
		if err := showStatus(ctx, kv); err != nil {
			log.Fatalf("Status failed: %v", err)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Mining Reward Settlement Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  mining-settlement generate [--type kawai|usdt]  Generate Merkle tree from balances")
	fmt.Println("  mining-settlement upload                        Upload Merkle root to contract")
	fmt.Println("  mining-settlement status                        Show current settlement status")
	fmt.Println("")
	fmt.Println("Workflow:")
	fmt.Println("  1. Run 'generate' weekly to create Merkle tree from accumulated rewards")
	fmt.Println("  2. Run 'upload' to upload the Merkle root to MiningRewardDistributor contract")
	fmt.Println("  3. Contributors can then claim their rewards via the frontend")
}

func generateMerkleTree(ctx context.Context, kv store.Store, rewardType string) error {
	log.Printf("🌳 Generating Merkle tree for %s mining rewards...", rewardType)
	log.Println("")

	// Generate mining settlement with 9-field Merkle leaves
	period, err := kv.GenerateMiningSettlement(ctx, rewardType)
	if err != nil {
		return fmt.Errorf("failed to generate settlement: %w", err)
	}

	log.Printf("✅ Merkle tree generated successfully!")
	log.Println("")
	log.Printf("📊 Settlement Details:")
	log.Printf("   Period ID: %d", period.PeriodID)
	log.Printf("   Merkle Root: %s", period.MerkleRoot)
	log.Printf("   Contributors: %d", period.ContributorCount)
	log.Printf("   Total Amount: %s %s", period.TotalAmount, rewardType)
	log.Printf("   Proofs Saved: %d", period.ProofsSaved)
	log.Printf("   Status: %s", period.Status)
	log.Println("")
	log.Printf("📝 Next step: Run 'mining-settlement upload' to upload root to contract")
	log.Println("")
	log.Printf("💡 Contributors can claim rewards via frontend after root is uploaded")

	return nil
}

func uploadMerkleRoot(ctx context.Context, kv store.Store) error {
	log.Println("🚀 Uploading Merkle root to MiningRewardDistributor contract...")

	// Get all settlement periods
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list periods: %w", err)
	}

	if len(periods) == 0 {
		return fmt.Errorf("no settlement periods found - run 'generate' first")
	}

	// Get latest KAWAI period
	var latest *store.SettlementPeriod
	for i := len(periods) - 1; i >= 0; i-- {
		if periods[i].RewardType == "kawai" {
			latest = periods[i]
			break
		}
	}

	if latest == nil {
		return fmt.Errorf("no KAWAI settlement periods found - run 'generate --type kawai' first")
	}
	if latest.Status != store.SettlementStatusProofsSaved && latest.Status != store.SettlementStatusCompleted {
		return fmt.Errorf("latest period not ready for upload (status: %s)", latest.Status)
	}

	log.Printf("📊 Latest Period:")
	log.Printf("   Period ID: %d", latest.PeriodID)
	log.Printf("   Merkle Root: %s", latest.MerkleRoot)
	log.Printf("   Total Amount: %s", latest.TotalAmount)
	log.Printf("")

	// TODO: Implement contract interaction
	log.Println("⚠️  Contract upload not yet implemented")
	log.Println("📋 Manual upload command:")
	log.Printf("   cast send 0xa0dDC59DAcBA9201CC9Ef613707d287b77b2723F \\")
	log.Printf("     'setMerkleRoot(bytes32)' %s \\", latest.MerkleRoot)
	log.Printf("     --rpc-url $RPC_URL --private-key <ADMIN_PRIVATE_KEY>")
	log.Println("")
	log.Println("   Or use 'advancePeriod(bytes32)' for new period")

	return nil
}

func showStatus(ctx context.Context, kv store.Store) error {
	log.Println("📊 Mining Reward Settlement Status")
	log.Println("")

	// Get all settlement periods
	allPeriods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list settlement periods: %w", err)
	}

	if len(allPeriods) == 0 {
		log.Println("No settlement periods found")
		log.Println("")
		log.Println("💡 Tip: Run 'generate --type kawai' to create first settlement period")
		return nil
	}

	// Filter by reward type and show last 5 of each
	var kawaiPeriods, usdtPeriods []*store.SettlementPeriod
	for _, p := range allPeriods {
		if p.RewardType == "kawai" {
			kawaiPeriods = append(kawaiPeriods, p)
		} else if p.RewardType == "usdt" {
			usdtPeriods = append(usdtPeriods, p)
		}
	}

	// Show KAWAI settlements (last 5)
	if len(kawaiPeriods) > 0 {
		log.Println("KAWAI Settlements:")
		log.Println("─────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-12s | %-20s | %-15s | %-10s | %s\n", "Period ID", "Status", "Total Amount", "Proofs", "Created")
		log.Println("─────────────────────────────────────────────────────────────────────")

		// Show last 5
		start := 0
		if len(kawaiPeriods) > 5 {
			start = len(kawaiPeriods) - 5
		}
		for _, p := range kawaiPeriods[start:] {
			created := time.Unix(p.PeriodID, 0).Format("2006-01-02 15:04")
			fmt.Printf("%-12d | %-20s | %-15s | %-10d | %s\n",
				p.PeriodID, p.Status, p.TotalAmount, p.ProofsSaved, created)
		}
		log.Println("")
	} else {
		log.Println("No KAWAI settlements found")
		log.Println("")
	}

	// Show USDT settlements (last 5)
	if len(usdtPeriods) > 0 {
		log.Println("USDT Settlements:")
		log.Println("─────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-12s | %-20s | %-15s | %-10s | %s\n", "Period ID", "Status", "Total Amount", "Proofs", "Created")
		log.Println("─────────────────────────────────────────────────────────────────────")

		// Show last 5
		start := 0
		if len(usdtPeriods) > 5 {
			start = len(usdtPeriods) - 5
		}
		for _, p := range usdtPeriods[start:] {
			created := time.Unix(p.PeriodID, 0).Format("2006-01-02 15:04")
			fmt.Printf("%-12d | %-20s | %-15s | %-10d | %s\n",
				p.PeriodID, p.Status, p.TotalAmount, p.ProofsSaved, created)
		}
		log.Println("")
	} else {
		log.Println("No USDT settlements found")
		log.Println("")
	}

	log.Println("💡 Tip: Run 'generate --type kawai' weekly to create new settlement period")

	return nil
}
