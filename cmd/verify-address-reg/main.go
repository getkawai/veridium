package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/kawai-network/veridium/pkg/jarvis/db"
)

func main() {
	names := []string{"MockUSDT", "KawaiToken", "Escrow", "PaymentVault"}
	expected := map[string]string{
		"MockUSDT":     "0x312C4fC3598AC9B54375eD12BbF55af83f86f862",
		"KawaiToken":   "0xD85758a0BC00a22a95E9201551ADC1b1E59A7A83",
		"Escrow":       "0x0F0E32877f8eC14d12E500D7642b2109A02Dd466",
		"PaymentVault": "0xa6Fc4FaF4CD7a4E3f300D164a37CB45d35bf28eD",
	}

	fmt.Println("=== Verifying Jarvis Address Resolution ===")
	for _, name := range names {
		desc, err := db.GetAddress(name)
		if err != nil {
			log.Fatalf("FAILED to resolve %s: %v", name, err)
		}

		fmt.Printf("✓ Resolved %s -> %s\n", name, desc.Address)
		if strings.ToLower(desc.Address) != strings.ToLower(expected[name]) {
			log.Fatalf("  ERROR: Expected %s but got %s", expected[name], desc.Address)
		}
	}
	fmt.Println("\n✓ All addresses resolved correctly!")
}
