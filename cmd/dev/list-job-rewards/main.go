package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	fmt.Println("📋 Listing job_rewards keys...")
	fmt.Println("")

	accountID := constant.GetCfAccountId()
	apiToken := constant.GetCfApiToken()

	client, err := store.NewKVClient(apiToken, accountID)
	if err != nil {
		log.Fatalf("Failed to create KV client: %v", err)
	}

	ctx := context.Background()
	contributorsNS := constant.GetCfKvContributorsNamespaceId()

	keys, err := client.ListAllKeys(ctx, contributorsNS, "job_rewards:")
	if err != nil {
		log.Fatalf("Failed to list keys: %v", err)
	}

	jobRewardsCount := 0
	contributorAddresses := make(map[string]bool)

	for _, key := range keys {
		jobRewardsCount++
		parts := strings.Split(key, ":")
		if len(parts) >= 2 {
			contributorAddresses[parts[1]] = true
		}
		fmt.Printf("  %s\n", key)
	}

	fmt.Println("")
	fmt.Printf("Total job_rewards: %d\n", jobRewardsCount)
	fmt.Printf("Unique contributors: %d\n", len(contributorAddresses))
	fmt.Println("")
	fmt.Println("Contributors:")
	for addr := range contributorAddresses {
		fmt.Printf("  %s\n", addr)
	}
}
