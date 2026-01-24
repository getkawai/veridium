package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/types"
)

func main() {
	ctx := context.Background()

	fmt.Println("🚀 Uploading All Mining Merkle Roots")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatalf("Failed to initialize KV store: %v", err)
	}

	// Get all mining settlement periods
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		log.Fatalf("Failed to list settlement periods: %v", err)
	}

	// Filter for mining periods (kawai type)
	var miningPeriods []*store.SettlementPeriod
	for _, period := range periods {
		if period.RewardType == types.RewardTypeMining {
			miningPeriods = append(miningPeriods, period)
		}
	}

	if len(miningPeriods) == 0 {
		log.Fatalf("No mining settlement periods found")
	}

	fmt.Printf("Found %d mining settlement periods:\n", len(miningPeriods))
	for i, period := range miningPeriods {
		fmt.Printf("  %d. Period %d: %s KAWAI (Root: %s)\n",
			i+1, period.PeriodID, period.TotalAmount, period.MerkleRoot)
	}
	fmt.Println()

	// Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Monad: %v", err)
	}
	defer client.Close()

	// Load MiningRewardDistributor contract
	distributorAddr := common.HexToAddress(constant.MiningRewardDistributorAddress)
	distributor, err := miningdistributor.NewMiningRewardDistributor(distributorAddr, client)
	if err != nil {
		log.Fatalf("Failed to load MiningRewardDistributor: %v", err)
	}

	// Get private key
	privateKeyHex := constant.GetAdminPrivateKey()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Context = ctx

	fmt.Printf("⚠️  About to upload %d Merkle roots to MiningRewardDistributor\n", len(miningPeriods))
	fmt.Printf("Contract: %s\n", distributorAddr.Hex())
	fmt.Println()

	// Upload each period
	successCount := 0
	for i, period := range miningPeriods {
		fmt.Printf("📤 Uploading %d/%d: Period %d\n", i+1, len(miningPeriods), period.PeriodID)

		// Parse Merkle root
		merkleRootHex := period.MerkleRoot
		if strings.HasPrefix(merkleRootHex, "0x") {
			merkleRootHex = merkleRootHex[2:]
		}
		merkleRootBytes := common.Hex2Bytes(merkleRootHex)
		if len(merkleRootBytes) != 32 {
			fmt.Printf("   ❌ Invalid Merkle root length: expected 32 bytes, got %d\n", len(merkleRootBytes))
			continue
		}
		var merkleRoot [32]byte
		copy(merkleRoot[:], merkleRootBytes)

		// Upload Merkle root
		fmt.Printf("   🌳 Uploading root: %s\n", period.MerkleRoot)
		tx, err := distributor.SetMerkleRoot(auth, merkleRoot)
		if err != nil {
			fmt.Printf("   ❌ Failed to upload: %v\n", err)
			continue
		}

		fmt.Printf("   ✅ Transaction sent: %s\n", tx.Hash().Hex())
		fmt.Printf("   ⏳ Waiting for confirmation...\n")

		// Wait for confirmation
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			fmt.Printf("   ❌ Failed to confirm: %v\n", err)
			continue
		}

		if receipt.Status != 1 {
			fmt.Printf("   ❌ Transaction failed with status: %d\n", receipt.Status)
			continue
		}

		fmt.Printf("   ✅ Confirmed in block %d\n", receipt.BlockNumber.Uint64())
		fmt.Println()
		successCount++
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("✅ Upload completed: %d/%d successful\n", successCount, len(miningPeriods))

	if successCount == len(miningPeriods) {
		fmt.Println("🎉 All mining Merkle roots uploaded successfully!")
		fmt.Println("📝 Users can now claim mining rewards via UI")
	} else {
		fmt.Printf("⚠️  %d uploads failed - check logs above\n", len(miningPeriods)-successCount)
	}
}
