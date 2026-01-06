package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("🔑 Generating Test Wallets for Mining Rewards Testing")
	fmt.Println("======================================================")
	fmt.Println("")

	// Generate 3 test wallets
	for i := 1; i <= 3; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA)

		fmt.Printf("Test Wallet #%d:\n", i)
		fmt.Printf("  Address:     %s\n", address.Hex())
		fmt.Printf("  Private Key: %s\n", privateKeyHex)
		fmt.Println("")
	}

	fmt.Println("======================================================")
	fmt.Println("📝 Instructions:")
	fmt.Println("")
	fmt.Println("1. Copy the addresses above")
	fmt.Println("2. Update cmd/test-inject-mining-data/main.go with these addresses")
	fmt.Println("3. Re-run: make test-inject-mining-data")
	fmt.Println("4. Re-run: go run cmd/mining-settlement/main.go generate --type kawai")
	fmt.Println("5. Re-run: make upload-mining-root PERIOD=<new_period_id>")
	fmt.Println("6. Import private key to MetaMask for testing")
	fmt.Println("")
	fmt.Println("⚠️  IMPORTANT: These are TEST wallets only!")
	fmt.Println("   Do NOT use for production or real funds!")
}

