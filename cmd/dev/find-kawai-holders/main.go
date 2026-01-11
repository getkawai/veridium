package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/blockchain"
)

func main() {
	ctx := context.Background()

	// Initialize holder scanner
	scanner, err := blockchain.NewHolderScanner()
	if err != nil {
		log.Fatalf("Failed to create scanner: %v", err)
	}

	fmt.Println("🔍 Scanning for KAWAI holders...")
	fmt.Printf("Token Address: %s\n", constant.KawaiTokenAddress)
	fmt.Println()

	// Get total supply
	totalSupply, err := scanner.GetTotalSupply(ctx)
	if err != nil {
		log.Fatalf("Failed to get total supply: %v", err)
	}
	fmt.Printf("📊 Total KAWAI Supply: %s\n\n", totalSupply.String())

	// Get current block
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RPC: %v", err)
	}

	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("Failed to get current block: %v", err)
	}

	startBlock := currentBlock
	if currentBlock > 90 {
		startBlock = currentBlock - 90
	}

	fmt.Printf("🔍 Scanning blocks %d to %d (last 90 blocks)\n\n", startBlock, currentBlock)

	holders, err := scanner.ScanHoldersFromBlock(ctx, startBlock)
	if err != nil {
		log.Fatalf("Failed to scan holders: %v", err)
	}

	fmt.Printf("📋 Found %d unique addresses in recent transfers\n\n", len(holders))

	// Check balances
	fmt.Println("💰 Checking balances...")
	holdersWithBalance := 0

	for i, addr := range holders {
		if i >= 10 {
			fmt.Printf("... and %d more addresses\n", len(holders)-10)
			break
		}

		balance, err := scanner.GetBalance(ctx, addr)
		if err != nil {
			fmt.Printf("  ❌ %s - Error: %v\n", addr.Hex(), err)
			continue
		}

		if balance.Cmp(big.NewInt(0)) > 0 {
			holdersWithBalance++
			fmt.Printf("  ✅ %s - Balance: %s KAWAI\n", addr.Hex(), balance.String())
		} else {
			fmt.Printf("  ⚪ %s - Balance: 0\n", addr.Hex())
		}
	}

	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("  Total addresses scanned: %d\n", len(holders))
	fmt.Printf("  Holders with balance: %d\n", holdersWithBalance)
}
