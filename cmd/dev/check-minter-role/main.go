package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/jarvis/contracts"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
)

// MINTER_ROLE = keccak256("MINTER_ROLE")
var MINTER_ROLE = common.HexToHash("0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6")

func main() {
	// Connect to Monad RPC
	nodes := map[string]string{"monad": constant.MonadRpcUrl}
	ethReader := reader.NewEthReaderGeneric(nodes, nil)

	// Load KAWAI token contract
	kawaiToken, err := contracts.KawaiToken("KawaiToken", ethReader)
	if err != nil {
		log.Fatalf("Failed to load KAWAI token contract: %v", err)
	}

	// Define distributors to check (updated 2026-01-12)
	distributors := map[string]string{
		"MiningRewardDistributor":    constant.MiningRewardDistributorAddr,
		"DepositCashbackDistributor": constant.CashbackDistributorAddress,
		"ReferralRewardDistributor":  constant.KawaiDistributorAddr,
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("🔐 MINTER_ROLE Status Check")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("")

	allGranted := true

	for name, address := range distributors {
		addr := common.HexToAddress(address)
		hasRole, err := kawaiToken.HasRole(nil, MINTER_ROLE, addr)
		if err != nil {
			log.Printf("❌ Failed to check %s: %v", name, err)
			allGranted = false
			continue
		}

		status := "❌ NOT GRANTED"
		if hasRole {
			status = "✅ GRANTED"
		} else {
			allGranted = false
		}

		fmt.Printf("%-30s %s\n", name+":", status)
		fmt.Printf("   Address: %s\n", address)
		fmt.Println("")
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	if allGranted {
		fmt.Println("✅ All distributors have MINTER_ROLE!")
		fmt.Println("   Ready for reward claims.")
	} else {
		fmt.Println("⚠️  Some distributors are missing MINTER_ROLE!")
		fmt.Println("   Run: ./GRANT_ALL_MINTER_ROLES.sh")
	}
	fmt.Println("═══════════════════════════════════════════════════════════")

	if !allGranted {
		os.Exit(1)
	}
}
