package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("🧹 COMPLETE KV CLEANUP - DELETE ALL DATA")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("⚠️  DANGER: This will DELETE ALL DATA from ALL KV namespaces!")
	fmt.Println("")
	fmt.Println("Namespaces to be cleaned:")
	fmt.Println("  • Contributors (job rewards, balances)")
	fmt.Println("  • Proofs (Merkle proofs)")
	fmt.Println("  • Settlements (settlement periods)")
	fmt.Println("  • Cashback")
	fmt.Println("  • Holders")
	fmt.Println("  • Users")
	fmt.Println("  • P2P Marketplace")
	fmt.Println("")

	// Check for confirmation flag
	if len(os.Args) < 2 || os.Args[1] != "--confirm-delete-all" {
		fmt.Println("To proceed, run:")
		fmt.Println("  go run cmd/dev/cleanup-kv-all/main.go --confirm-delete-all")
		fmt.Println("")
		os.Exit(0)
	}

	fmt.Println("🚀 Starting COMPLETE cleanup...")
	fmt.Println("")

	// Get Cloudflare credentials
	accountID := constant.GetCfAccountId()
	apiToken := constant.GetCfApiToken()

	// Initialize Cloudflare client
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatalf("Failed to create Cloudflare client: %v", err)
	}

	ctx := context.Background()

	// Get namespace IDs
	contributorsNS := constant.GetCfKvContributorsNamespaceId()
	proofsNS := constant.GetCfKvProofsNamespaceId()
	settlementsNS := constant.GetCfKvSettlementsNamespaceId()
	cashbackNS := constant.GetCfKvCashbackNamespaceId()
	holdersNS := constant.GetCfKvHolderNamespaceId()
	usersNS := constant.GetCfKvUsersNamespaceId()
	marketplaceNS := constant.GetCfKvP2pMarketplaceNamespaceId()

	// Cleanup each namespace
	cleanupNamespace(ctx, api, accountID, contributorsNS, "Contributors")
	cleanupNamespace(ctx, api, accountID, proofsNS, "Proofs")
	cleanupNamespace(ctx, api, accountID, settlementsNS, "Settlements")
	cleanupNamespace(ctx, api, accountID, cashbackNS, "Cashback")
	cleanupNamespace(ctx, api, accountID, holdersNS, "Holders")
	cleanupNamespace(ctx, api, accountID, usersNS, "Users")
	cleanupNamespace(ctx, api, accountID, marketplaceNS, "P2P Marketplace")

	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("✅ COMPLETE CLEANUP FINISHED!")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("All KV data has been deleted. You can now start fresh.")
}

func cleanupNamespace(ctx context.Context, api *cloudflare.API, accountID, namespaceID, name string) {
	fmt.Printf("🗑️  Cleaning %s namespace...\n", name)

	deletedCount := 0
	cursor := ""

	for {
		// List keys in this namespace
		params := cloudflare.ListWorkersKVsParams{
			NamespaceID: namespaceID,
			Limit:       1000, // Max per request
		}
		if cursor != "" {
			params.Cursor = cursor
		}

		resp, err := api.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(accountID), params)
		if err != nil {
			log.Printf("   ⚠️  Failed to list keys: %v", err)
			break
		}

		if len(resp.Result) == 0 {
			break
		}

		// Delete each key
		for _, key := range resp.Result {
			_, err := api.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(accountID), cloudflare.DeleteWorkersKVEntryParams{
				NamespaceID: namespaceID,
				Key:         key.Name,
			})
			if err != nil {
				log.Printf("   ⚠️  Failed to delete key %s: %v", key.Name, err)
				continue
			}
			deletedCount++
		}

		// Check if there are more keys
		if resp.ResultInfo.Cursor == "" {
			break
		}
		cursor = resp.ResultInfo.Cursor
	}

	fmt.Printf("   ✅ Deleted %d keys from %s\n", deletedCount, name)
}
