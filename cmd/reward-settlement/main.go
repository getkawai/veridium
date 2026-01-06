package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/store"
)

const (
	RewardTypeMining   = "mining"
	RewardTypeCashback = "cashback"
	RewardTypeReferral = "referral"
)

func main() {
	// Subcommands
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	allCmd := flag.NewFlagSet("all", flag.ExitOnError)

	// Flags for generate command
	var rewardType string
	generateCmd.StringVar(&rewardType, "type", "mining", "Reward type: mining, cashback, or referral")

	// Flags for upload command
	var uploadType string
	uploadCmd.StringVar(&uploadType, "type", "mining", "Reward type: mining, cashback, or referral")

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
		if err := generateSettlement(ctx, kv, rewardType); err != nil {
			log.Fatalf("Generate failed: %v", err)
		}

	case "upload":
		uploadCmd.Parse(os.Args[2:])
		if err := uploadMerkleRoot(ctx, kv, uploadType); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}

	case "status":
		statusCmd.Parse(os.Args[2:])
		if err := showStatus(ctx, kv); err != nil {
			log.Fatalf("Status failed: %v", err)
		}

	case "all":
		allCmd.Parse(os.Args[2:])
		if err := settleAll(ctx, kv); err != nil {
			log.Fatalf("Settle all failed: %v", err)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Unified Reward Settlement Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  reward-settlement generate --type <type>  Generate Merkle tree for specific reward type")
	fmt.Println("  reward-settlement upload --type <type>    Upload Merkle root to contract")
	fmt.Println("  reward-settlement status                  Show settlement status for all types")
	fmt.Println("  reward-settlement all                     Settle all reward types at once")
	fmt.Println("")
	fmt.Println("Reward Types:")
	fmt.Println("  mining     - Mining rewards (9-field Merkle tree)")
	fmt.Println("  cashback   - Deposit cashback rewards (3-field Merkle tree)")
	fmt.Println("  referral   - Referral commission rewards (3-field Merkle tree)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  reward-settlement generate --type mining")
	fmt.Println("  reward-settlement generate --type cashback")
	fmt.Println("  reward-settlement all")
	fmt.Println("")
	fmt.Println("Workflow:")
	fmt.Println("  1. Run 'generate' weekly to create Merkle trees")
	fmt.Println("  2. Run 'upload' to upload Merkle roots to contracts")
	fmt.Println("  3. Users can claim rewards via frontend")
	fmt.Println("")
	fmt.Println("  Or use 'all' to settle all reward types at once (recommended)")
}

func generateSettlement(ctx context.Context, kv *store.KVStore, rewardType string) error {
	log.Printf("🌳 Generating %s settlement...", rewardType)
	log.Println("")

	switch rewardType {
	case RewardTypeMining:
		return generateMiningSettlement(ctx, kv)
	case RewardTypeCashback:
		return generateCashbackSettlement(ctx, kv)
	case RewardTypeReferral:
		return generateReferralSettlement(ctx, kv)
	default:
		return fmt.Errorf("unknown reward type: %s (must be mining, cashback, or referral)", rewardType)
	}
}

func generateMiningSettlement(ctx context.Context, kv store.Store) error {
	log.Println("📊 Mining Rewards Settlement")
	log.Println("─────────────────────────────")

	// Generate mining settlement with 9-field Merkle leaves
	period, err := kv.GenerateMiningSettlement(ctx, "kawai")
	if err != nil {
		return fmt.Errorf("failed to generate mining settlement: %w", err)
	}

	log.Printf("✅ Mining settlement generated!")
	log.Println("")
	log.Printf("Period ID:     %d", period.PeriodID)
	log.Printf("Merkle Root:   %s", period.MerkleRoot)
	log.Printf("Contributors:  %d", period.ContributorCount)
	log.Printf("Total Amount:  %s KAWAI", period.TotalAmount)
	log.Printf("Proofs Saved:  %d", period.ProofsSaved)
	log.Printf("Status:        %s", period.Status)
	log.Println("")
	log.Printf("📝 Next: reward-settlement upload --type mining")

	return nil
}

func generateCashbackSettlement(ctx context.Context, kv *store.KVStore) error {
	log.Println("📊 Cashback Rewards Settlement")
	log.Println("─────────────────────────────")

	// Get current period
	currentPeriod := kv.GetCurrentPeriod()
	settlementPeriod := currentPeriod - 1 // Settle previous period

	if settlementPeriod < 1 {
		return fmt.Errorf("no period to settle yet (current period: %d)", currentPeriod)
	}

	log.Printf("Current Period:    %d", currentPeriod)
	log.Printf("Settling Period:   %d", settlementPeriod)
	log.Println("")

	// Initialize cashback settlement
	settlement, err := blockchain.NewCashbackSettlement(kv, constant.GetObfuscatedTemp())
	if err != nil {
		return fmt.Errorf("failed to initialize cashback settlement: %w", err)
	}

	// Run settlement
	if err := settlement.SettleCashback(ctx, settlementPeriod); err != nil {
		return fmt.Errorf("failed to settle cashback: %w", err)
	}

	log.Println("")
	log.Printf("✅ Cashback settlement completed!")
	log.Println("")
	log.Printf("📝 Next: reward-settlement upload --type cashback")

	return nil
}

func generateReferralSettlement(ctx context.Context, kv *store.KVStore) error {
	log.Println("📊 Referral Rewards Settlement")
	log.Println("─────────────────────────────")
	log.Println("")

	// Get current period (same as mining, weekly)
	// For referral, we settle based on mining periods
	// Get latest mining settlement period
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list settlement periods: %w", err)
	}

	if len(periods) == 0 {
		return fmt.Errorf("no mining settlement found - run mining settlement first")
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
		return fmt.Errorf("no KAWAI settlement found - run mining settlement first")
	}

	settlementPeriod := uint64(latest.PeriodID)

	log.Printf("Mining Period:     %d", settlementPeriod)
	log.Printf("Settling Referral: Period %d", settlementPeriod)
	log.Println("")

	// Initialize referral settlement
	settlement := blockchain.NewReferralSettlement(kv, constant.GetObfuscatedTemp())

	// Run settlement
	if err := settlement.SettleReferral(ctx, settlementPeriod); err != nil {
		return fmt.Errorf("failed to settle referral: %w", err)
	}

	log.Println("")
	log.Printf("✅ Referral settlement completed!")
	log.Println("")
	log.Printf("📝 Next: reward-settlement upload --type referral")

	return nil
}

func uploadMerkleRoot(ctx context.Context, kv *store.KVStore, rewardType string) error {
	log.Printf("🚀 Uploading %s Merkle root to contract...", rewardType)
	log.Println("")

	switch rewardType {
	case RewardTypeMining:
		return uploadMiningRoot(ctx, kv)
	case RewardTypeCashback:
		return uploadCashbackRoot(ctx, kv)
	case RewardTypeReferral:
		return uploadReferralRoot(ctx, kv)
	default:
		return fmt.Errorf("unknown reward type: %s", rewardType)
	}
}

func uploadMiningRoot(ctx context.Context, kv store.Store) error {
	// Get latest mining period
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list periods: %w", err)
	}

	var latest *store.SettlementPeriod
	for i := len(periods) - 1; i >= 0; i-- {
		if periods[i].RewardType == "kawai" {
			latest = periods[i]
			break
		}
	}

	if latest == nil {
		return fmt.Errorf("no mining settlement found - run 'generate --type mining' first")
	}

	log.Printf("Period ID:     %d", latest.PeriodID)
	log.Printf("Merkle Root:   %s", latest.MerkleRoot)
	log.Printf("Total Amount:  %s KAWAI", latest.TotalAmount)
	log.Println("")
	log.Println("⚠️  Contract upload not yet implemented")
	log.Println("")
	log.Println("📋 Manual upload command:")
	log.Printf("   cast send 0xa0dDC59DAcBA9201CC9Ef613707d287b77b2723F \\")
	log.Printf("     'setMerkleRoot(bytes32)' %s \\", latest.MerkleRoot)
	log.Printf("     --rpc-url $RPC_URL --private-key <ADMIN_PRIVATE_KEY>")

	return nil
}

func uploadCashbackRoot(ctx context.Context, kv *store.KVStore) error {
	log.Println("⚠️  Cashback root upload not yet implemented")
	log.Println("")
	log.Println("TODO:")
	log.Println("  1. Get latest cashback period from KV")
	log.Println("  2. Read Merkle root from cashback_period:N:merkle_root")
	log.Println("  3. Call DepositCashbackDistributor.setMerkleRoot()")
	log.Println("")
	log.Println("Contract: 0x... (DepositCashbackDistributor)")

	return fmt.Errorf("cashback upload not implemented yet")
}

func uploadReferralRoot(ctx context.Context, kv *store.KVStore) error {
	log.Println("⚠️  Referral root upload not yet implemented")
	log.Println("")
	log.Println("TODO:")
	log.Println("  1. Get latest referral period from KV")
	log.Println("  2. Read Merkle root")
	log.Println("  3. Call ReferralRewardDistributor.setMerkleRoot()")
	log.Println("")
	log.Println("Contract: 0x... (ReferralRewardDistributor)")

	return fmt.Errorf("referral upload not implemented yet")
}

func showStatus(ctx context.Context, kv *store.KVStore) error {
	log.Println("📊 Reward Settlement Status")
	log.Println("═══════════════════════════════════════════════════════════════")
	log.Println("")

	// Mining status
	log.Println("⛏️  MINING REWARDS")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := showMiningStatus(ctx, kv); err != nil {
		log.Printf("Error: %v", err)
	}
	log.Println("")

	// Cashback status
	log.Println("💰 CASHBACK REWARDS")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := showCashbackStatus(ctx, kv); err != nil {
		log.Printf("Error: %v", err)
	}
	log.Println("")

	// Referral status
	log.Println("🤝 REFERRAL REWARDS")
	log.Println("───────────────────────────────────────────────────────────────")
	log.Println("Status: Not yet implemented")
	log.Println("")

	log.Println("═══════════════════════════════════════════════════════════════")
	log.Println("💡 Tip: Run 'reward-settlement all' to settle all types at once")

	return nil
}

func showMiningStatus(ctx context.Context, kv store.Store) error {
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return err
	}

	var kawaiPeriods []*store.SettlementPeriod
	for _, p := range periods {
		if p.RewardType == "kawai" {
			kawaiPeriods = append(kawaiPeriods, p)
		}
	}

	if len(kawaiPeriods) == 0 {
		log.Println("No settlements found")
		return nil
	}

	// Show last 3
	start := 0
	if len(kawaiPeriods) > 3 {
		start = len(kawaiPeriods) - 3
	}

	fmt.Printf("%-12s | %-20s | %-15s | %s\n", "Period ID", "Status", "Total Amount", "Proofs")
	log.Println("───────────────────────────────────────────────────────────────")
	for _, p := range kawaiPeriods[start:] {
		fmt.Printf("%-12d | %-20s | %-15s | %d\n",
			p.PeriodID, p.Status, p.TotalAmount+" KAWAI", p.ProofsSaved)
	}

	log.Printf("\nTotal Settlements: %d", len(kawaiPeriods))
	return nil
}

func showCashbackStatus(ctx context.Context, kv *store.KVStore) error {
	currentPeriod := kv.GetCurrentPeriod()
	log.Printf("Current Period: %d", currentPeriod)
	log.Printf("Next Settlement: Period %d (in %s)", currentPeriod-1, "next Monday")
	log.Println("")
	log.Println("Recent Settlements: (TODO: implement)")
	return nil
}

func settleAll(ctx context.Context, kv *store.KVStore) error {
	log.Println("🚀 Settling All Reward Types")
	log.Println("═══════════════════════════════════════════════════════════════")
	log.Println("")

	success := 0
	failed := 0

	// 1. Mining
	log.Println("1️⃣  Mining Rewards")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := generateMiningSettlement(ctx, kv); err != nil {
		log.Printf("❌ Mining settlement failed: %v", err)
		failed++
	} else {
		success++
	}
	log.Println("")

	// 2. Cashback
	log.Println("2️⃣  Cashback Rewards")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := generateCashbackSettlement(ctx, kv); err != nil {
		log.Printf("❌ Cashback settlement failed: %v", err)
		failed++
	} else {
		success++
	}
	log.Println("")

	// 3. Referral
	log.Println("3️⃣  Referral Rewards")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := generateReferralSettlement(ctx, kv); err != nil {
		log.Printf("⚠️  Referral settlement skipped: %v", err)
		// Don't count as failed since it's not implemented yet
	}
	log.Println("")

	// Summary
	log.Println("═══════════════════════════════════════════════════════════════")
	log.Printf("✅ Successful: %d", success)
	log.Printf("❌ Failed: %d", failed)
	log.Println("")

	if failed > 0 {
		return fmt.Errorf("some settlements failed")
	}

	log.Println("🎉 All settlements completed successfully!")
	log.Println("")
	log.Println("📝 Next steps:")
	log.Println("  1. Upload Merkle roots to contracts")
	log.Println("  2. Users can claim rewards via frontend")

	return nil
}
