package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <recipient_address> [amount_in_mon]")
		fmt.Println("Example: go run main.go 0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E 0.1")
		os.Exit(1)
	}

	recipient := os.Args[1]
	amount := "0.1" // Default 0.1 MON
	if len(os.Args) > 2 {
		amount = os.Args[2]
	}

	if err := sendMON(recipient, amount); err != nil {
		log.Fatalf("Failed to send MON: %v", err)
	}
}

func sendMON(recipientAddr, amountStr string) error {
	ctx := context.Background()

	log.Println("💰 Sending MON Testnet Tokens")
	log.Println("═══════════════════════════════════════════════════════════")

	// Connect to Monad RPC
	rpcURL := constant.MonadRpcUrl
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	// Get private key
	privateKeyHex := constant.GetObfuscatedTemp()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get sender address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("failed to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Parse recipient address
	if !common.IsHexAddress(recipientAddr) {
		return fmt.Errorf("invalid recipient address: %s", recipientAddr)
	}
	toAddress := common.HexToAddress(recipientAddr)

	// Parse amount (convert MON to wei)
	amountFloat := new(big.Float)
	amountFloat.SetString(amountStr)

	// 1 MON = 10^18 wei
	weiPerMON := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	amountWei := new(big.Float).Mul(amountFloat, weiPerMON)

	amountBigInt, _ := amountWei.Int(nil)

	// Get nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	// Get chain ID
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Check sender balance
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	log.Printf("From:      %s", fromAddress.Hex())
	log.Printf("To:        %s", toAddress.Hex())
	log.Printf("Amount:    %s MON", amountStr)
	log.Printf("Wei:       %s", amountBigInt.String())
	log.Printf("Balance:   %s wei", balance.String())
	log.Printf("Nonce:     %d", nonce)
	log.Printf("Gas Price: %s", gasPrice.String())
	log.Printf("Chain ID:  %s", chainID.String())
	log.Println("")

	// Check if sender has enough balance
	gasLimit := uint64(21000) // Standard transfer gas limit
	totalCost := new(big.Int).Add(amountBigInt, new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))))

	if balance.Cmp(totalCost) < 0 {
		return fmt.Errorf("insufficient balance: need %s wei, have %s wei", totalCost.String(), balance.String())
	}

	// Create transaction
	tx := types.NewTransaction(nonce, toAddress, amountBigInt, gasLimit, gasPrice, nil)

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	log.Println("⏳ Sending transaction...")
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("✅ Transaction sent: %s", signedTx.Hash().Hex())
	log.Println("")

	// Wait for confirmation
	log.Println("⏳ Waiting for confirmation...")
	receipt, err := waitForConfirmation(ctx, client, signedTx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	log.Printf("✅ MON sent successfully!")
	log.Println("═══════════════════════════════════════════════════════════")
	log.Printf("Transaction Hash: %s", receipt.TxHash.Hex())
	log.Printf("Block Number:     %d", receipt.BlockNumber.Uint64())
	log.Printf("Gas Used:         %d", receipt.GasUsed)
	log.Printf("Explorer:         https://explorer.monad.xyz/tx/%s", receipt.TxHash.Hex())
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("")
	log.Printf("✅ %s now has %s MON for gas fees", toAddress.Hex(), amountStr)
	log.Printf("✅ Ready to test claiming flows!")

	return nil
}

func waitForConfirmation(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		// If error is "not found", continue waiting
		if strings.Contains(err.Error(), "not found") {
			continue
		}

		return nil, err
	}
}
