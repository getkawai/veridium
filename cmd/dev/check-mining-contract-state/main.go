package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/types"
)

func main() {
	if err := checkMiningContractState(); err != nil {
		log.Fatalf("Failed to check contract state: %v", err)
	}
}

func checkMiningContractState() error {
	ctx := context.Background()

	fmt.Println("🔍 Checking MiningRewardDistributor Contract State")
	fmt.Println("═══════════════════════════════════════════════════")
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

	fmt.Printf("Contract Address: %s\n", constant.MiningRewardDistributorAddress)
	fmt.Println()

	// 1. Check current period
	currentPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period: %w", err)
	}

	fmt.Printf("📊 Contract State:\n")
	fmt.Printf("   Current Period: %s\n", currentPeriod.String())
	fmt.Println()

	// 2. Check period Merkle roots for periods 1-10
	fmt.Printf("🌳 Period Merkle Roots:\n")
	for i := 1; i <= 10; i++ {
		periodBig := big.NewInt(int64(i))
		root, err := distributor.PeriodMerkleRoots(nil, periodBig)
		if err != nil {
			fmt.Printf("   Period %d: ERROR - %v\n", i, err)
			continue
		}

		rootHex := common.Bytes2Hex(root[:])
		if rootHex == "0000000000000000000000000000000000000000000000000000000000000000" {
			fmt.Printf("   Period %d: (empty)\n", i)
		} else {
			fmt.Printf("   Period %d: 0x%s\n", i, rootHex)
		}
	}
	fmt.Println()

	// 3. Check current merkle root (for current period)
	currentRoot, err := distributor.MerkleRoot(nil)
	if err != nil {
		return fmt.Errorf("failed to get current merkle root: %w", err)
	}

	currentRootHex := common.Bytes2Hex(currentRoot[:])
	fmt.Printf("🎯 Current Merkle Root: 0x%s\n", currentRootHex)
	fmt.Println()

	// 4. Get stats
	stats, err := distributor.GetStats(nil)
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	fmt.Printf("📈 Contract Stats:\n")
	fmt.Printf("   Period:             %s\n", stats.Period.String())
	fmt.Printf("   Contributor Rewards: %s\n", stats.ContributorRewards.String())
	fmt.Printf("   Developer Rewards:   %s\n", stats.DeveloperRewards.String())
	fmt.Printf("   User Rewards:        %s\n", stats.UserRewards.String())
	fmt.Printf("   Affiliator Rewards:  %s\n", stats.AffiliatorRewards.String())
	fmt.Println()

	// 5. Compare with our settlement data
	fmt.Printf("🔄 Comparing with Settlement Data:\n")

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

	fmt.Printf("   Settlement Periods Found: %d\n", len(kawaiPeriods))
	for i, period := range kawaiPeriods {
		fmt.Printf("   Settlement %d: Period ID %d, Root: %s\n",
			i+1, period.PeriodID, period.MerkleRoot)
	}
	fmt.Println()

	// 6. Recommendations
	fmt.Printf("💡 Analysis:\n")
	if currentPeriod.Uint64() == 0 {
		fmt.Printf("   ❌ Contract is at period 0 - no periods have been advanced yet\n")
		fmt.Printf("   📝 Need to use advancePeriod() instead of setMerkleRoot()\n")
	} else {
		fmt.Printf("   ✅ Contract is at period %d\n", currentPeriod.Uint64())
		if len(kawaiPeriods) > int(currentPeriod.Uint64()) {
			fmt.Printf("   ⚠️  We have %d settlement periods but contract only has %d periods\n",
				len(kawaiPeriods), currentPeriod.Uint64())
			fmt.Printf("   📝 Need to advance contract to match settlement data\n")
		}
	}

	fmt.Println()
	fmt.Printf("🔧 Next Steps:\n")
	fmt.Printf("   1. Use advancePeriod() to advance contract periods 1, 2, 3, 4\n")
	fmt.Printf("   2. Map settlement period IDs to sequential contract periods\n")
	fmt.Printf("   3. Update claiming logic to use sequential period numbers\n")

	return nil
}
