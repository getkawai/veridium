package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/pkg/admin"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	snapshotCmd := flag.NewFlagSet("snapshot", flag.ExitOnError)
	auditCmd := flag.NewFlagSet("audit", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'snapshot' or 'audit' subcommands")
		os.Exit(1)
	}

	// Shared Config (In prod, use env vars)
	accountID := "ceab218751d33cd804878196ad7bef74"
	apiToken := "OP8BZQhyeJxrovCPKt15eUOSC6i5LXTVECGRSMc1"
	namespaceID := "55bd9f26233940dabb65f2a1992bfae9"

	// Initialize Store
	kv, err := store.NewKVStore(apiToken, accountID, namespaceID)
	if err != nil {
		log.Fatalf("Failed to initialize KV store: %v", err)
	}

	// Initialize Blockchain Client (Optional for some commands)
	rpcURL := os.Getenv("BSC_RPC_URL")
	var chainClient *blockchain.Client
	if rpcURL != "" {
		chainClient, _ = blockchain.NewClient(blockchain.Config{
			RPCUrl: rpcURL,
		})
	}

	mgr := admin.NewAdminManager(chainClient, kv)
	ctx := context.Background()

	switch os.Args[1] {
	case "snapshot":
		snapshotCmd.Parse(os.Args[2:])
		if err := mgr.CalculateDividends(ctx); err != nil {
			log.Fatalf("Snapshot failed: %v", err)
		}
	case "audit":
		auditCmd.Parse(os.Args[2:])
		if err := mgr.AuditWorkers(ctx); err != nil {
			log.Fatalf("Audit failed: %v", err)
		}
	default:
		fmt.Println("expected 'snapshot' or 'audit' subcommands")
		os.Exit(1)
	}
}
