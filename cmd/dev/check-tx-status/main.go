package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/contracts"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <tx_hash>", os.Args[0])
	}

	txHash := os.Args[1]
	if err := checkTransactionStatus(txHash); err != nil {
		log.Fatalf("Failed to check transaction status: %v", err)
	}
}

func checkTransactionStatus(txHashStr string) error {
	ctx := context.Background()

	fmt.Printf("🔍 Checking Transaction Status\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("TX Hash: %s\n", txHashStr)
	fmt.Println()

	// Connect to Monad RPC
	client, err := ethclient.Dial(contracts.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// Parse transaction hash
	txHash := common.HexToHash(txHashStr)

	// Get transaction receipt
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		fmt.Printf("❌ Transaction not found or not mined yet: %v\n", err)
		return nil
	}

	fmt.Printf("✅ Transaction Found!\n")
	fmt.Printf("   Block Number: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("   Gas Used:     %d\n", receipt.GasUsed)
	fmt.Printf("   Status:       %d ", receipt.Status)

	if receipt.Status == 1 {
		fmt.Printf("(✅ SUCCESS)\n")
	} else {
		fmt.Printf("(❌ FAILED)\n")
	}

	fmt.Printf("   Explorer:     https://testnet.monadexplorer.com/tx/%s\n", txHashStr)
	fmt.Println()

	// If successful, update KV store to mark as confirmed
	if receipt.Status == 1 {
		fmt.Printf("🔄 Updating KV store to mark claim as confirmed...\n")

		// Initialize KV store
		kv, err := store.NewMultiNamespaceKVStore()
		if err != nil {
			return fmt.Errorf("failed to initialize KV store: %w", err)
		}

		// For now, we'll just print that we would update it
		// In a real implementation, we'd need the address and period ID
		fmt.Printf("✅ Transaction confirmed! KV store should be updated.\n")
		fmt.Printf("   (Manual update needed - requires address and period ID)\n")

		_ = kv // Suppress unused variable warning
	}

	return nil
}
