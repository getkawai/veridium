package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/pkg/jarvis/binding"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
	"github.com/kawai-network/contracts"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <address>")
		fmt.Println("Example: go run main.go 0x1234...")
		os.Exit(1)
	}

	address := os.Args[1]

	// Connect to Monad RPC
	nodes := map[string]string{"monad": contracts.MonadRpcUrl}
	ethReader := reader.NewEthReaderGeneric(nodes, nil)

	// Load KAWAI token contract
	kawaiToken, err := binding.KawaiToken("KawaiToken", ethReader)
	if err != nil {
		log.Fatalf("Failed to load KAWAI token contract: %v", err)
	}

	// Get balance
	balance, err := kawaiToken.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	// Convert to KAWAI (divide by 1e18)
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	kawaiBalance := new(big.Int).Div(balance, divisor)

	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Address: %s\n", address)
	fmt.Printf("KAWAI Balance: %s KAWAI\n", kawaiBalance.String())
	fmt.Printf("Wei Balance: %s wei\n", balance.String())
	fmt.Println("═══════════════════════════════════════")
}
