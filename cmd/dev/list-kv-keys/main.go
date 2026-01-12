package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
)

func main() {
	fmt.Println("📋 Listing all KV keys...")
	fmt.Println("")

	accountID := constant.GetCfAccountId()
	apiToken := constant.GetCfApiToken()

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatalf("Failed to create Cloudflare client: %v", err)
	}

	ctx := context.Background()
	proofsNS := constant.GetCfKvProofsNamespaceId()

	fmt.Println("🔍 Proofs namespace keys:")
	fmt.Println("")

	params := cloudflare.ListWorkersKVsParams{
		NamespaceID: proofsNS,
		Limit:       100,
	}

	resp, err := api.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(accountID), params)
	if err != nil {
		log.Fatalf("Failed to list keys: %v", err)
	}

	if len(resp.Result) == 0 {
		fmt.Println("❌ No keys found!")
		return
	}

	for i, key := range resp.Result {
		fmt.Printf("%d. %s\n", i+1, key.Name)
	}

	fmt.Println("")
	fmt.Printf("Total: %d keys\n", len(resp.Result))
}
