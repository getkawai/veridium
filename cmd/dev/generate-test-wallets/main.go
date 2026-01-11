package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("🔑 Generating Test Wallets for Claiming")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Generate 3 test wallets
	for i := 1; i <= 3; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Failed to generate private key: %v", err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hexutil.Encode(privateKeyBytes)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatalf("Failed to cast public key to ECDSA")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA)

		fmt.Printf("🔐 Test Wallet %d:\n", i)
		fmt.Printf("   Address:     %s\n", address.Hex())
		fmt.Printf("   Private Key: %s\n", privateKeyHex)
		fmt.Println()
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("📝 Next Steps:")
	fmt.Println("1. Copy these addresses")
	fmt.Println("2. Send MON tokens: make send-test-mon ADDR=<address>")
	fmt.Println("3. Mint KAWAI tokens to these addresses")
	fmt.Println("4. Inject mining data for these addresses")
	fmt.Println("5. Run settlement to generate claimable rewards")
	fmt.Println("6. Import private keys to UI for claiming")
}
