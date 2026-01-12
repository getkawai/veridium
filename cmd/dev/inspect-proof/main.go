package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
)

func main() {
	fmt.Println("🔍 Inspecting Merkle proof...")
	fmt.Println("")

	accountID := constant.GetCfAccountId()
	apiToken := constant.GetCfApiToken()

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatalf("Failed to create Cloudflare client: %v", err)
	}

	ctx := context.Background()
	proofsNS := constant.GetCfKvProofsNamespaceId()

	key := "0xab48220e6721754b906c30463142dc0a8f5ebba2:1"

	value, err := api.GetWorkersKV(ctx, cloudflare.AccountIdentifier(accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: proofsNS,
		Key:         key,
	})
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}

	// Pretty print JSON
	var data map[string]interface{}
	if err := json.Unmarshal(value, &data); err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}
