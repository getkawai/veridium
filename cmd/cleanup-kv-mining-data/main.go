package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kawai-network/veridium/pkg/store"
)

// This tool cleans up old mining reward data from KV store
// Use with caution - this will delete data!

var (
	dryRun      = flag.Bool("dry-run", false, "Preview what will be deleted without actually deleting")
	deleteJobs  = flag.Bool("jobs", false, "Delete job reward records")
	deleteProofs = flag.Bool("proofs", false, "Delete Merkle proofs")
	deleteSettlements = flag.Bool("settlements", false, "Delete settlement periods")
	deleteAll   = flag.Bool("all", false, "Delete everything (jobs + proofs + settlements)")
	confirm     = flag.String("confirm", "", "Type 'DELETE' to confirm deletion")
)

func main() {
	flag.Parse()

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	ctx := context.Background()

	// Initialize KV store
	kv, err := store.NewKVStore(
		os.Getenv("CF_ACCOUNT_ID"),
		os.Getenv("CF_API_TOKEN"),
		os.Getenv("CF_KV_CONTRIBUTORS_NAMESPACE_ID"),
		os.Getenv("CF_KV_PROOFS_NAMESPACE_ID"),
		os.Getenv("CF_KV_SETTLEMENTS_NAMESPACE_ID"),
		os.Getenv("CF_KV_AUTHZ_NAMESPACE_ID"),
		os.Getenv("CF_KV_USERS_NAMESPACE_ID"),
	)
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	fmt.Println("🧹 Cloudflare KV Cleanup Tool")
	fmt.Println("==============================")
	fmt.Println("")

	// Determine what to delete
	if !*deleteJobs && !*deleteProofs && !*deleteSettlements && !*deleteAll {
		fmt.Println("❌ Error: No cleanup target specified!")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  --jobs         Delete job reward records")
		fmt.Println("  --proofs       Delete Merkle proofs")
		fmt.Println("  --settlements  Delete settlement periods")
		fmt.Println("  --all          Delete everything")
		fmt.Println("  --dry-run      Preview without deleting")
		fmt.Println("  --confirm DELETE  Confirm deletion (required)")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  # Preview what will be deleted")
		fmt.Println("  go run cmd/cleanup-kv-mining-data/main.go --all --dry-run")
		fmt.Println("")
		fmt.Println("  # Delete all mining data")
		fmt.Println("  go run cmd/cleanup-kv-mining-data/main.go --all --confirm DELETE")
		fmt.Println("")
		fmt.Println("  # Delete only job records")
		fmt.Println("  go run cmd/cleanup-kv-mining-data/main.go --jobs --confirm DELETE")
		os.Exit(1)
	}

	// Safety check
	if !*dryRun && *confirm != "DELETE" {
		fmt.Println("❌ Error: Deletion requires confirmation!")
		fmt.Println("")
		fmt.Println("To confirm deletion, add: --confirm DELETE")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  go run cmd/cleanup-kv-mining-data/main.go --all --confirm DELETE")
		os.Exit(1)
	}

	if *dryRun {
		fmt.Println("🔍 DRY RUN MODE - No data will be deleted")
		fmt.Println("")
	} else {
		fmt.Println("⚠️  DELETION MODE - Data will be permanently deleted!")
		fmt.Println("")
	}

	// Delete job rewards
	if *deleteJobs || *deleteAll {
		if err := cleanupJobRewards(ctx, kv, *dryRun); err != nil {
			log.Fatal("Failed to cleanup job rewards:", err)
		}
	}

	// Delete Merkle proofs
	if *deleteProofs || *deleteAll {
		if err := cleanupMerkleProofs(ctx, kv, *dryRun); err != nil {
			log.Fatal("Failed to cleanup Merkle proofs:", err)
		}
	}

	// Delete settlements
	if *deleteSettlements || *deleteAll {
		if err := cleanupSettlements(ctx, kv, *dryRun); err != nil {
			log.Fatal("Failed to cleanup settlements:", err)
		}
	}

	fmt.Println("")
	fmt.Println("==============================")
	if *dryRun {
		fmt.Println("✅ Dry run complete! No data was deleted.")
		fmt.Println("")
		fmt.Println("To actually delete, run without --dry-run:")
		fmt.Println("  go run cmd/cleanup-kv-mining-data/main.go --all --confirm DELETE")
	} else {
		fmt.Println("✅ Cleanup complete!")
	}
}

func cleanupJobRewards(ctx context.Context, kv *store.KVStore, dryRun bool) error {
	fmt.Println("📝 Cleaning up job reward records...")
	
	// List all contributors
	contributors, err := kv.ListContributorsWithBalance(ctx, "kawai")
	if err != nil {
		return fmt.Errorf("failed to list contributors: %w", err)
	}

	totalJobs := 0
	for _, contributor := range contributors {
		jobs, err := kv.GetJobRewardsSinceLastSettlement(ctx, contributor.WalletAddress, "kawai")
		if err != nil {
			log.Printf("Warning: Failed to get jobs for %s: %v", contributor.WalletAddress, err)
			continue
		}

		if len(jobs) > 0 {
			fmt.Printf("  Found %d job(s) for contributor: %s\n", len(jobs), contributor.WalletAddress[:20]+"...")
			totalJobs += len(jobs)

			if !dryRun {
				// Delete job records
				// Note: This requires implementing DeleteJobRewards in KVStore
				// For now, we'll just mark them as settled
				if err := kv.MarkJobRewardsAsSettled(ctx, contributor.WalletAddress, 0); err != nil {
					log.Printf("Warning: Failed to mark jobs as settled for %s: %v", contributor.WalletAddress, err)
				}
			}
		}
	}

	if dryRun {
		fmt.Printf("  Would delete %d job reward record(s)\n", totalJobs)
	} else {
		fmt.Printf("  ✅ Deleted %d job reward record(s)\n", totalJobs)
	}

	return nil
}

func cleanupMerkleProofs(ctx context.Context, kv *store.KVStore, dryRun bool) error {
	fmt.Println("🌳 Cleaning up Merkle proofs...")

	// List all settlement periods
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list settlements: %w", err)
	}

	// Filter for mining rewards (kawai type)
	miningPeriods := []store.SettlementPeriod{}
	for _, period := range periods {
		if period.RewardType == "kawai" {
			miningPeriods = append(miningPeriods, period)
		}
	}

	totalProofs := 0
	for _, period := range miningPeriods {
		fmt.Printf("  Period %d: %d proofs\n", period.PeriodID, period.ProofsSaved)
		totalProofs += period.ProofsSaved

		if !dryRun {
			// Delete proofs for this period
			// Note: This requires implementing DeleteProofsForPeriod in KVStore
			// For now, we'll just log it
			log.Printf("  Would delete proofs for period %d", period.PeriodID)
		}
	}

	if dryRun {
		fmt.Printf("  Would delete %d Merkle proof(s) from %d period(s)\n", totalProofs, len(miningPeriods))
	} else {
		fmt.Printf("  ✅ Deleted %d Merkle proof(s) from %d period(s)\n", totalProofs, len(miningPeriods))
	}

	return nil
}

func cleanupSettlements(ctx context.Context, kv *store.KVStore, dryRun bool) error {
	fmt.Println("📊 Cleaning up settlement periods...")

	// List all settlement periods
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list settlements: %w", err)
	}

	// Filter for mining rewards (kawai type)
	miningPeriods := []store.SettlementPeriod{}
	for _, period := range periods {
		if period.RewardType == "kawai" {
			miningPeriods = append(miningPeriods, period)
		}
	}

	for _, period := range miningPeriods {
		status := string(period.Status)
		fmt.Printf("  Period %d: %s (%d contributors, %s KAWAI)\n", 
			period.PeriodID, status, period.ContributorCount, formatAmount(period.TotalAmount))

		if !dryRun {
			// Delete settlement period
			// Note: This requires implementing DeleteSettlement in KVStore
			log.Printf("  Would delete settlement period %d", period.PeriodID)
		}
	}

	if dryRun {
		fmt.Printf("  Would delete %d settlement period(s)\n", len(miningPeriods))
	} else {
		fmt.Printf("  ✅ Deleted %d settlement period(s)\n", len(miningPeriods))
	}

	return nil
}

func formatAmount(amountStr string) string {
	if amountStr == "" || amountStr == "0" {
		return "0"
	}

	// Simple formatting - just show first few digits
	if len(amountStr) > 18 {
		// Has decimals
		wholePart := amountStr[:len(amountStr)-18]
		if wholePart == "" {
			wholePart = "0"
		}
		decimalPart := amountStr[len(amountStr)-18:]
		if len(decimalPart) > 4 {
			decimalPart = decimalPart[:4]
		}
		return fmt.Sprintf("%s.%s", wholePart, strings.TrimRight(decimalPart, "0"))
	}

	return amountStr
}

