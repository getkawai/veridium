package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/x/store"
	"github.com/kawai-network/x/constant"
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

	// Initialize KV client
	client, err := store.NewKVClient(apiToken, accountID)
	if err != nil {
		log.Fatalf("Failed to create KV client: %v", err)
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
	cleanupNamespace(ctx, client, contributorsNS, "Contributors")
	cleanupNamespace(ctx, client, proofsNS, "Proofs")
	cleanupNamespace(ctx, client, settlementsNS, "Settlements")
	cleanupNamespace(ctx, client, cashbackNS, "Cashback")
	cleanupNamespace(ctx, client, holdersNS, "Holders")
	cleanupNamespace(ctx, client, usersNS, "Users")
	cleanupNamespace(ctx, client, marketplaceNS, "P2P Marketplace")

	fmt.Println("")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("✅ COMPLETE CLEANUP FINISHED!")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("All KV data has been deleted. You can now start fresh.")
}

func cleanupNamespace(ctx context.Context, client *store.KVClient, namespaceID, name string) {
	fmt.Printf("🗑️  Cleaning %s namespace...\n", name)

	deletedCount := 0
	cursor := ""

	for {
		// List keys in this namespace
		result, err := client.ListKeys(ctx, namespaceID, "", cursor)
		if err != nil {
			log.Printf("   ⚠️  Failed to list keys: %v", err)
			break
		}

		if len(result.Result) == 0 {
			break
		}

		// Delete each key
		for _, key := range result.Result {
			err := client.DeleteValue(ctx, namespaceID, key.Name)
			if err != nil {
				log.Printf("   ⚠️  Failed to delete key %s: %v", key.Name, err)
				continue
			}
			deletedCount++
		}

		// Check if there are more keys
		if result.ResultInfo.Cursor == "" {
			break
		}
		cursor = result.ResultInfo.Cursor
	}

	fmt.Printf("   ✅ Deleted %d keys from %s\n", deletedCount, name)
}
