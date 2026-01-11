package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	if err := checkJobRewards(); err != nil {
		log.Fatalf("Failed to check job rewards: %v", err)
	}
}

func checkJobRewards() error {
	ctx := context.Background()

	fmt.Println("🔍 Checking Job Rewards Data")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Get all unsettled job rewards
	jobRewardsByContributor, err := kv.GetAllUnsettledJobRewards(ctx, "kawai")
	if err != nil {
		return fmt.Errorf("failed to get unsettled job rewards: %w", err)
	}

	fmt.Printf("📊 Found job rewards for %d contributors:\n", len(jobRewardsByContributor))
	fmt.Println()

	for contributorAddr, jobs := range jobRewardsByContributor {
		fmt.Printf("👤 Contributor: %s\n", contributorAddr)
		fmt.Printf("   Jobs: %d\n", len(jobs))

		if len(jobs) > 0 {
			// Show first job as example
			job := jobs[0]
			fmt.Printf("   Example Job:\n")
			fmt.Printf("     User Address:      %s\n", job.UserAddress)
			fmt.Printf("     Referrer Address:  %s\n", job.ReferrerAddress)
			fmt.Printf("     Developer Address: %s\n", job.DeveloperAddress)
			fmt.Printf("     Contributor Amount: %s KAWAI\n", job.ContributorAmount)
			fmt.Printf("     Has Referrer:      %t\n", job.HasReferrer)
			fmt.Printf("     Timestamp:         %s\n", job.Timestamp.Format("2006-01-02 15:04:05"))
		}
		fmt.Println()
	}

	// Check if we have any real user addresses
	fmt.Printf("🔍 Address Analysis:\n")
	realAddresses := 0
	testAddresses := 0

	for _, jobs := range jobRewardsByContributor {
		for _, job := range jobs {
			if job.UserAddress == "0xTestUser3333333333333333333333333333333" ||
				job.ReferrerAddress == "0xTestReferrer333333333333333333333333333" {
				testAddresses++
			} else {
				realAddresses++
			}
		}
	}

	fmt.Printf("   Real addresses: %d\n", realAddresses)
	fmt.Printf("   Test addresses: %d\n", testAddresses)
	fmt.Println()

	if testAddresses > 0 {
		fmt.Printf("⚠️  Found test addresses in job rewards!\n")
		fmt.Printf("📝 This explains the 'Invalid user address' error\n")
		fmt.Printf("💡 Solutions:\n")
		fmt.Printf("   1. Generate new job rewards with real user addresses\n")
		fmt.Printf("   2. Or create test settlement with the actual claiming address\n")
		fmt.Println()
	}

	// Check what address should be used for testing
	testAddress := "0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"
	fmt.Printf("🎯 Test Address: %s\n", testAddress)
	fmt.Printf("📝 For mining claims to work, the Merkle proof must contain this address as UserAddress\n")

	return nil
}
