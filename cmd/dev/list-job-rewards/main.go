package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
)

func main() {
	fmt.Println("📋 Listing job_rewards keys...")
	fmt.Println("")

	accountID := constant.GetCfAccountId()
	apiToken := constant.GetCfApiToken()

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatalf("Failed to create Cloudflare client: %v", err)
	}

	ctx := context.Background()
	contributorsNS := constant.GetCfKvContributorsNamespaceId()

	params := cloudflare.ListWorkersKVsParams{
		NamespaceID: contributorsNS,
		Limit:       100,
	}

	resp, err := api.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(accountID), params)
	if err != nil {
		log.Fatalf("Failed to list keys: %v", err)
	}

	jobRewardsCount := 0
	contributorAddresses := make(map[string]bool)

	for _, key := range resp.Result {
		if strings.HasPrefix(key.Name, "job_rewards:") {
			jobRewardsCount++
			parts := strings.Split(key.Name, ":")
			if len(parts) >= 2 {
				contributorAddresses[parts[1]] = true
			}
			fmt.Printf("  %s\n", key.Name)
		}
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
