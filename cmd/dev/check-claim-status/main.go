package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/y/jarvis/binding"
	"github.com/kawai-network/y/jarvis/util/reader"
	"github.com/kawai-network/contracts"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <contract_type> <period> <address>")
		fmt.Println("")
		fmt.Println("Contract Types:")
		fmt.Println("  mining   - Check MiningRewardDistributor")
		fmt.Println("  cashback - Check DepositCashbackDistributor")
		fmt.Println("  referral - Check ReferralRewardDistributor")
		fmt.Println("")
		fmt.Println("Example: go run main.go mining 1767549424 0x1234...")
		os.Exit(1)
	}

	contractType := os.Args[1]
	period := os.Args[2]
	address := os.Args[3]

	// Parse period
	periodBig := new(big.Int)
	periodBig, ok := periodBig.SetString(period, 10)
	if !ok {
		log.Fatalf("Invalid period: %s", period)
	}

	// Connect to Monad RPC
	nodes := map[string]string{"monad": contracts.MonadRpcUrl}
	ethReader := reader.NewEthReaderGeneric(nodes, nil)

	addr := common.HexToAddress(address)

	var hasClaimed bool

	switch contractType {
	case "mining":
		distributor, err := binding.MiningRewardDistributor("MiningRewardDistributor", ethReader)
		if err != nil {
			log.Fatalf("Failed to load mining distributor: %v", err)
		}
		hasClaimed, err = distributor.HasClaimed(nil, periodBig, addr)
		if err != nil {
			log.Fatalf("Failed to check claim status: %v", err)
		}

	case "cashback":
		distributor, err := binding.CashbackDistributor("CashbackDistributor", ethReader)
		if err != nil {
			log.Fatalf("Failed to load cashback distributor: %v", err)
		}
		hasClaimed, err = distributor.HasClaimed(nil, periodBig, addr)
		if err != nil {
			log.Fatalf("Failed to check claim status: %v", err)
		}

	case "referral":
		// Note: ReferralRewardDistributor uses MerkleDistributor base
		// which has hasClaimed(uint256 index, address account)
		// For period-based, we need to check the contract implementation
		fmt.Println("⚠️  Referral claim status check not yet implemented")
		fmt.Println("   (Requires contract method verification)")
		return

	default:
		log.Fatalf("Unknown contract type: %s (use: mining, cashback, or referral)", contractType)
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Contract: %s\n", contractType)
	fmt.Printf("Period: %s\n", period)
	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Has Claimed: %v\n", hasClaimed)
	if hasClaimed {
		fmt.Println("Status: ✅ Already Claimed")
	} else {
		fmt.Println("Status: ⏳ Not Claimed Yet")
	}
	fmt.Println("═══════════════════════════════════════")
}
