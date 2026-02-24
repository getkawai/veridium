package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/x/store"
	"github.com/kawai-network/x/constant"
)

func main() {
	fmt.Println("📋 Listing all KV keys...")
	fmt.Println("")

	accountID := constant.GetCfAccountId()
	apiToken := constant.GetCfApiToken()

	client, err := store.NewKVClient(apiToken, accountID)
	if err != nil {
		log.Fatalf("Failed to create KV client: %v", err)
	}

	ctx := context.Background()
	proofsNS := constant.GetCfKvProofsNamespaceId()

	fmt.Println("🔍 Proofs namespace keys:")
	fmt.Println("")

	keys, err := client.ListAllKeys(ctx, proofsNS, "")
	if err != nil {
		log.Fatalf("Failed to list keys: %v", err)
	}

	if len(keys) == 0 {
		fmt.Println("❌ No keys found!")
		return
	}

	for i, key := range keys {
		fmt.Printf("%d. %s\n", i+1, key)
	}

	fmt.Println("")
	fmt.Printf("Total: %d keys\n", len(keys))
}
