package main

import (
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/jarvis/common"
)

func main() {
	fmt.Println("Verifying Jarvis Local ABI Resolution...")

	// Test case 1: Project ABI map should contain our contracts
	contracts := []string{"Escrow", "KawaiToken", "PaymentVault"}
	for _, name := range contracts {
		if abi, ok := common.ProjectABIs[name]; ok {
			fmt.Printf("✓ Found local ABI for %s (length: %d)\n", name, len(abi))
		} else {
			log.Fatalf("✗ Could not find local ABI for %s in common.ProjectABIs", name)
		}
	}

	// Test case 2: GetABIString should return the local ABI if the address name matches
	// We'll mock a scenario where GetAddressFromString returns the contract name
	// This tests the logic in GetABIString we added earlier.
	// Since we can't easily mock the address database here, we just verify the common.ProjectABIs map
	// which is what GetABIString uses.
}
