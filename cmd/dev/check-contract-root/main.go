package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
)

func main() {
	if err := checkContractRoot(); err != nil {
		log.Fatalf("Failed to check contract root: %v", err)
	}
}

func checkContractRoot() error {
	fmt.Println("🔍 Contract Root Check")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// Load MiningRewardDistributor contract
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

	fmt.Printf("📊 Contract current period: %d\n", currentPeriod.Uint64())
	fmt.Println()

	// Check roots for all periods
	for period := int64(1); period <= int64(currentPeriod.Uint64()); period++ {
		root, err := distributor.PeriodMerkleRoots(nil, big.NewInt(period))
		if err != nil {
			fmt.Printf("❌ Failed to get root for period %d: %v\n", period, err)
			continue
		}

		fmt.Printf("📋 Period %d:\n", period)
		fmt.Printf("   Root: 0x%x\n", root)

		// Check if root is zero (not set)
		var zeroRoot [32]byte
		if root == zeroRoot {
			fmt.Printf("   Status: NOT SET\n")
		} else {
			fmt.Printf("   Status: SET\n")
		}
		fmt.Println()
	}

	// Expected roots from our settlements
	expectedRoots := map[int64]string{
		7: "0xff77aeb0d8b803ac73709f80e3aab1ec566dbab8a0a3ea8182242062cf3ee19e", // Our FIXED settlement
	}

	fmt.Printf("🎯 Expected vs Actual:\n")
	for period, expectedRoot := range expectedRoots {
		actualRoot, err := distributor.PeriodMerkleRoots(nil, big.NewInt(period))
		if err != nil {
			fmt.Printf("❌ Period %d: Failed to get actual root\n", period)
			continue
		}

		expectedBytes := common.Hex2Bytes(expectedRoot[2:])
		var expectedRoot32 [32]byte
		copy(expectedRoot32[:], expectedBytes)

		fmt.Printf("   Period %d:\n", period)
		fmt.Printf("     Expected: %s\n", expectedRoot)
		fmt.Printf("     Actual:   0x%x\n", actualRoot)

		if actualRoot == expectedRoot32 {
			fmt.Printf("     Match:    ✅ YES\n")
		} else {
			fmt.Printf("     Match:    ❌ NO\n")
		}
		fmt.Println()
	}

	return nil
}
