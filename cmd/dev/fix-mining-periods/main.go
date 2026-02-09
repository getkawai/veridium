package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/types"
	"github.com/kawai-network/x/constant"
)

func main() {
	if err := fixMiningPeriods(); err != nil {
		log.Fatalf("Failed to fix mining periods: %v", err)
	}
}

func fixMiningPeriods() error {
	ctx := context.Background()

	fmt.Println("🔧 Fixing Mining Contract Periods")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Get settlement periods
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list settlement periods: %w", err)
	}

	var kawaiPeriods []*store.SettlementPeriod
	for _, p := range periods {
		if p.RewardType == types.RewardTypeMining {
			kawaiPeriods = append(kawaiPeriods, p)
		}
	}

	if len(kawaiPeriods) == 0 {
		return fmt.Errorf("no KAWAI settlement periods found")
	}

	// Sort by period ID (oldest first)
	sort.Slice(kawaiPeriods, func(i, j int) bool {
		return kawaiPeriods[i].PeriodID < kawaiPeriods[j].PeriodID
	})

	fmt.Printf("📊 Found %d KAWAI settlement periods:\n", len(kawaiPeriods))
	for i, period := range kawaiPeriods {
		fmt.Printf("   %d. Period ID %d -> Contract Period %d\n",
			i+1, period.PeriodID, i+1)
		fmt.Printf("      Root: %s\n", period.MerkleRoot)
		fmt.Printf("      Amount: %s KAWAI\n", period.TotalAmount)
	}
	fmt.Println()

	// Connect to contract
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	distributorAddr := common.HexToAddress(constant.MiningRewardDistributorAddress)
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
	fmt.Printf("📝 Need to advance to period: %d\n", len(kawaiPeriods))
	fmt.Println()

	if currentPeriod.Uint64() >= uint64(len(kawaiPeriods)) {
		fmt.Printf("✅ Contract is already at or beyond target period\n")
		return nil
	}

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

	// Advance periods
	for contractPeriod := currentPeriod.Uint64() + 1; contractPeriod <= uint64(len(kawaiPeriods)); contractPeriod++ {
		settlementIndex := contractPeriod - 1 // 0-based index
		settlement := kawaiPeriods[settlementIndex]

		fmt.Printf("🚀 Advancing to contract period %d...\n", contractPeriod)
		fmt.Printf("   Settlement Period ID: %d\n", settlement.PeriodID)
		fmt.Printf("   Merkle Root: %s\n", settlement.MerkleRoot)

		// Parse Merkle root
		merkleRootHex := settlement.MerkleRoot
		if strings.HasPrefix(merkleRootHex, "0x") {
			merkleRootHex = merkleRootHex[2:]
		}
		merkleRootBytes := common.Hex2Bytes(merkleRootHex)
		if len(merkleRootBytes) != 32 {
			return fmt.Errorf("invalid Merkle root length: expected 32 bytes, got %d", len(merkleRootBytes))
		}
		var merkleRoot [32]byte
		copy(merkleRoot[:], merkleRootBytes)

		// Confirm before advancing
		fmt.Printf("⚠️  About to advance contract to period %d\n", contractPeriod)
		if !confirm(fmt.Sprintf("Continue with period %d?", contractPeriod)) {
			return fmt.Errorf("operation cancelled by user")
		}

		// Advance period
		tx, err := distributor.AdvancePeriod(auth, merkleRoot)
		if err != nil {
			return fmt.Errorf("failed to advance period %d: %w", contractPeriod, err)
		}

		fmt.Printf("✅ AdvancePeriod transaction sent: %s\n", tx.Hash().Hex())
		fmt.Printf("⏳ Waiting for confirmation...\n")

		// Wait for confirmation
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			return fmt.Errorf("failed to wait for confirmation: %w", err)
		}

		if receipt.Status != 1 {
			return fmt.Errorf("transaction failed with status: %d", receipt.Status)
		}

		fmt.Printf("✅ Period %d advanced successfully in block %d\n",
			contractPeriod, receipt.BlockNumber.Uint64())
		fmt.Println()
	}

	// Verify final state
	finalPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get final period: %w", err)
	}

	fmt.Printf("🎉 Contract successfully advanced to period %d\n", finalPeriod.Uint64())
	fmt.Println()

	// Create period mapping for reference
	fmt.Printf("📋 Period Mapping (for claiming logic):\n")
	for i, settlement := range kawaiPeriods {
		contractPeriod := i + 1
		fmt.Printf("   Settlement Period %d -> Contract Period %d\n",
			settlement.PeriodID, contractPeriod)
	}
	fmt.Println()

	fmt.Printf("✅ Mining periods fixed successfully!\n")
	fmt.Printf("📝 Next: Update claiming logic to use sequential period numbers\n")

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
