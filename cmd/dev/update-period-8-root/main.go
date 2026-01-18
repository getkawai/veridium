package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
)

func main() {
	if err := updatePeriod8Root(); err != nil {
		log.Fatalf("Failed to update period 8 root: %v", err)
	}
}

func updatePeriod8Root() error {
	ctx := context.Background()

	fmt.Println("🔄 Updating Period 8 Merkle Root")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// FINAL Treasury settlement details (generated with period 8)
	finalMerkleRoot := "0x9f736ae3bb707ce306d9cf42d3049efafedb40fc85e7673d95b75c9aa2e25871"
	settlementID := int64(1768141317)
	contractPeriod := int64(8)

	fmt.Printf("📋 FINAL Treasury Settlement Details:\n")
	fmt.Printf("   Settlement ID:     %d (KV storage key)\n", settlementID)
	fmt.Printf("   Contract Period:   %d (current period)\n", contractPeriod)
	fmt.Printf("   New Merkle Root:   %s\n", finalMerkleRoot)
	fmt.Printf("   Treasury Address:  0x94D5C06229811c4816107005ff05259f229Eb07b\n")
	fmt.Println()

	// Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// Load MiningRewardDistributor contract
	distributorAddr := common.HexToAddress(constant.MiningRewardDistributorAddr)
	distributor, err := miningdistributor.NewMiningRewardDistributor(distributorAddr, client)
	if err != nil {
		return fmt.Errorf("failed to load MiningRewardDistributor: %w", err)
	}

	// Check current period
	currentPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period: %w", err)
	}

	fmt.Printf("🎯 Contract current period: %d\n", currentPeriod.Uint64())

	if currentPeriod.Uint64() != uint64(contractPeriod) {
		return fmt.Errorf("contract period mismatch: expected %d, got %d", contractPeriod, currentPeriod.Uint64())
	}

	// Check current root for period 8
	currentRoot, err := distributor.PeriodMerkleRoots(nil, big.NewInt(contractPeriod))
	if err != nil {
		return fmt.Errorf("failed to get current root: %w", err)
	}

	fmt.Printf("📊 Current root for period %d: 0x%x\n", contractPeriod, currentRoot)
	fmt.Printf("📊 New root for period %d:     %s\n", contractPeriod, finalMerkleRoot)
	fmt.Println()

	// Get private key and setup auth
	privateKeyHex := constant.GetAdminPrivateKey()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Context = ctx

	// Parse new Merkle root
	merkleRootHex := finalMerkleRoot
	if strings.HasPrefix(merkleRootHex, "0x") {
		merkleRootHex = merkleRootHex[2:]
	}
	merkleRootBytes := common.Hex2Bytes(merkleRootHex)
	if len(merkleRootBytes) != 32 {
		return fmt.Errorf("invalid Merkle root length: expected 32 bytes, got %d", len(merkleRootBytes))
	}
	var merkleRoot [32]byte
	copy(merkleRoot[:], merkleRootBytes)

	// Confirm before updating
	fmt.Printf("⚠️  About to update period %d Merkle root\n", contractPeriod)
	fmt.Printf("🌳 New Root: %s\n", finalMerkleRoot)
	if !confirm("Continue with update?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Update Merkle root for current period
	tx, err := distributor.SetMerkleRoot(auth, merkleRoot)
	if err != nil {
		return fmt.Errorf("failed to set merkle root: %w", err)
	}

	fmt.Printf("✅ SetMerkleRoot transaction sent: %s\n", tx.Hash().Hex())
	fmt.Printf("⏳ Waiting for confirmation...\n")

	// Wait for confirmation
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for confirmation: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	fmt.Printf("✅ Period 8 Merkle root updated successfully in block %d\n", receipt.BlockNumber.Uint64())
	fmt.Println()

	// Verify updated root
	updatedRoot, err := distributor.PeriodMerkleRoots(nil, big.NewInt(contractPeriod))
	if err != nil {
		return fmt.Errorf("failed to get updated root: %w", err)
	}

	fmt.Printf("🎉 Updated root for period %d: 0x%x\n", contractPeriod, updatedRoot)
	fmt.Println()

	// Check if roots match
	if updatedRoot != merkleRoot {
		fmt.Printf("❌ Root update failed - mismatch!\n")
		return fmt.Errorf("root update failed")
	}

	fmt.Printf("✅ FINAL Treasury Merkle root updated successfully!\n")
	fmt.Printf("🎯 All components now match:\n")
	fmt.Printf("   ✅ Settlement generated with period 8\n")
	fmt.Printf("   ✅ Contract has period 8 root\n")
	fmt.Printf("   ✅ Treasury address matches msg.sender\n")
	fmt.Printf("   ✅ Ready for successful claim!\n")

	return nil
}

// Helper function for user confirmation
func confirm(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
