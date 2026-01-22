package main

import (
	"fmt"

	"github.com/kawai-network/veridium/pkg/jarvis/contracts"
)

func main() {
	fmt.Println("🔍 Testing Contract Address Resolution")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Test contract names that should resolve
	testNames := []string{
		"USDT_Distributor",
		"KawaiToken",
		"MockStablecoin",
		"PaymentVault",
		"OTCMarket",
		"MiningRewardDistributor",
		"CashbackDistributor",
		"ReferralDistributor",
	}

	for _, name := range testNames {
		fmt.Printf("Testing: %s\n", name)
		addr, err := contracts.ResolveAddress(name)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
		} else {
			fmt.Printf("  ✅ Resolved: %s\n", addr.Hex())
		}
		fmt.Println()
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("✅ Contract resolution test complete!")
}
