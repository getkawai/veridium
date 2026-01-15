package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/cashbackdistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/referraldistributor"
)

var (
	action      = flag.String("action", "", "Action: pause, unpause, status")
	contractStr = flag.String("contract", "all", "Contract: mining, cashback, referral, all")
	dryRun      = flag.Bool("dry-run", false, "Dry run (don't send transactions)")
)

func main() {
	flag.Parse()

	if *action == "" {
		fmt.Println("Usage: go run main.go -action <pause|unpause|status> [-contract <mining|cashback|referral|all>] [-dry-run]")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run main.go -action status")
		fmt.Println("  go run main.go -action pause -contract all")
		fmt.Println("  go run main.go -action unpause -contract mining")
		fmt.Println("  go run main.go -action pause -contract all -dry-run")
		os.Exit(1)
	}

	// Connect to RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RPC: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Get private key
	privateKeyHex := constant.GetObfuscatedTemp()
	if privateKeyHex == "" {
		log.Fatal("Private key not found in constants")
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Get public address first (before creating transactor)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 100000

	fmt.Printf("🔐 Admin Address: %s\n", fromAddress.Hex())
	fmt.Printf("🌐 RPC: %s\n", constant.MonadRpcUrl)
	fmt.Printf("⛓️  Chain ID: %s\n\n", chainID.String())

	// Execute action
	switch *action {
	case "status":
		checkStatus(ctx, client)
	case "pause":
		if *dryRun {
			fmt.Println("🔍 DRY RUN MODE - No transactions will be sent\n")
		}
		pauseContracts(ctx, client, auth, *contractStr, *dryRun)
	case "unpause":
		if *dryRun {
			fmt.Println("🔍 DRY RUN MODE - No transactions will be sent\n")
		}
		unpauseContracts(ctx, client, auth, *contractStr, *dryRun)
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func checkStatus(ctx context.Context, client *ethclient.Client) {
	fmt.Println("📊 Checking pause status...\n")

	// Mining Distributor
	miningAddr := common.HexToAddress(constant.MiningRewardDistributorAddr)
	mining, err := miningdistributor.NewMiningRewardDistributor(miningAddr, client)
	if err != nil {
		log.Printf("❌ Failed to connect to Mining Distributor: %v", err)
	} else {
		paused, err := mining.Paused(&bind.CallOpts{Context: ctx})
		if err != nil {
			log.Printf("❌ Failed to check Mining Distributor status: %v", err)
		} else {
			status := "✅ Active"
			if paused {
				status = "🚨 PAUSED"
			}
			fmt.Printf("Mining Distributor (%s): %s\n", constant.MiningRewardDistributorAddr, status)
		}
	}

	// Cashback Distributor
	cashbackAddr := common.HexToAddress(constant.CashbackDistributorAddress)
	cashback, err := cashbackdistributor.NewDepositCashbackDistributor(cashbackAddr, client)
	if err != nil {
		log.Printf("❌ Failed to connect to Cashback Distributor: %v", err)
	} else {
		paused, err := cashback.Paused(&bind.CallOpts{Context: ctx})
		if err != nil {
			log.Printf("❌ Failed to check Cashback Distributor status: %v", err)
		} else {
			status := "✅ Active"
			if paused {
				status = "🚨 PAUSED"
			}
			fmt.Printf("Cashback Distributor (%s): %s\n", constant.CashbackDistributorAddress, status)
		}
	}

	// Referral Distributor (if deployed)
	if constant.ReferralDistributorAddress != "" {
		referralAddr := common.HexToAddress(constant.ReferralDistributorAddress)
		referral, err := referraldistributor.NewReferralRewardDistributor(referralAddr, client)
		if err != nil {
			log.Printf("❌ Failed to connect to Referral Distributor: %v", err)
		} else {
			paused, err := referral.Paused(&bind.CallOpts{Context: ctx})
			if err != nil {
				log.Printf("❌ Failed to check Referral Distributor status: %v", err)
			} else {
				status := "✅ Active"
				if paused {
					status = "🚨 PAUSED"
				}
				fmt.Printf("Referral Distributor (%s): %s\n", constant.ReferralDistributorAddress, status)
			}
		}
	}
}

func pauseContracts(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, contractStr string, dryRun bool) {
	fmt.Println("🚨 PAUSING DISTRIBUTORS...\n")

	contracts := getContracts(contractStr)

	for _, contract := range contracts {
		switch contract {
		case "mining":
			pauseMining(ctx, client, auth, dryRun)
		case "cashback":
			pauseCashback(ctx, client, auth, dryRun)
		case "referral":
			pauseReferral(ctx, client, auth, dryRun)
		}
	}

	if !dryRun {
		fmt.Println("\n✅ All specified distributors paused!")
		fmt.Println("💡 Run with -action status to verify")
	}
}

func unpauseContracts(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, contractStr string, dryRun bool) {
	fmt.Println("🔓 UNPAUSING DISTRIBUTORS...\n")

	contracts := getContracts(contractStr)

	for _, contract := range contracts {
		switch contract {
		case "mining":
			unpauseMining(ctx, client, auth, dryRun)
		case "cashback":
			unpauseCashback(ctx, client, auth, dryRun)
		case "referral":
			unpauseReferral(ctx, client, auth, dryRun)
		}
	}

	if !dryRun {
		fmt.Println("\n✅ All specified distributors unpaused!")
		fmt.Println("💡 Run with -action status to verify")
	}
}

func pauseMining(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, dryRun bool) {
	addr := common.HexToAddress(constant.MiningRewardDistributorAddr)
	mining, err := miningdistributor.NewMiningRewardDistributor(addr, client)
	if err != nil {
		log.Printf("❌ Mining Distributor: Failed to connect: %v", err)
		return
	}

	// Check if already paused
	paused, err := mining.Paused(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Printf("❌ Mining Distributor: Failed to check status: %v", err)
		return
	}

	if paused {
		fmt.Printf("⚠️  Mining Distributor: Already paused\n")
		return
	}

	if dryRun {
		fmt.Printf("🔍 Mining Distributor: Would pause (dry-run)\n")
		return
	}

	tx, err := mining.Pause(auth)
	if err != nil {
		log.Printf("❌ Mining Distributor: Failed to pause: %v", err)
		return
	}

	fmt.Printf("✅ Mining Distributor: Paused (tx: %s)\n", tx.Hash().Hex())
	auth.Nonce = new(big.Int).Add(auth.Nonce, big.NewInt(1))
}

func unpauseMining(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, dryRun bool) {
	addr := common.HexToAddress(constant.MiningRewardDistributorAddr)
	mining, err := miningdistributor.NewMiningRewardDistributor(addr, client)
	if err != nil {
		log.Printf("❌ Mining Distributor: Failed to connect: %v", err)
		return
	}

	// Check if already unpaused
	paused, err := mining.Paused(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Printf("❌ Mining Distributor: Failed to check status: %v", err)
		return
	}

	if !paused {
		fmt.Printf("⚠️  Mining Distributor: Already active\n")
		return
	}

	if dryRun {
		fmt.Printf("🔍 Mining Distributor: Would unpause (dry-run)\n")
		return
	}

	tx, err := mining.Unpause(auth)
	if err != nil {
		log.Printf("❌ Mining Distributor: Failed to unpause: %v", err)
		return
	}

	fmt.Printf("✅ Mining Distributor: Unpaused (tx: %s)\n", tx.Hash().Hex())
	auth.Nonce = new(big.Int).Add(auth.Nonce, big.NewInt(1))
}

func pauseCashback(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, dryRun bool) {
	addr := common.HexToAddress(constant.CashbackDistributorAddress)
	cashback, err := cashbackdistributor.NewDepositCashbackDistributor(addr, client)
	if err != nil {
		log.Printf("❌ Cashback Distributor: Failed to connect: %v", err)
		return
	}

	paused, err := cashback.Paused(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Printf("❌ Cashback Distributor: Failed to check status: %v", err)
		return
	}

	if paused {
		fmt.Printf("⚠️  Cashback Distributor: Already paused\n")
		return
	}

	if dryRun {
		fmt.Printf("🔍 Cashback Distributor: Would pause (dry-run)\n")
		return
	}

	tx, err := cashback.Pause(auth)
	if err != nil {
		log.Printf("❌ Cashback Distributor: Failed to pause: %v", err)
		return
	}

	fmt.Printf("✅ Cashback Distributor: Paused (tx: %s)\n", tx.Hash().Hex())
	auth.Nonce = new(big.Int).Add(auth.Nonce, big.NewInt(1))
}

func unpauseCashback(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, dryRun bool) {
	addr := common.HexToAddress(constant.CashbackDistributorAddress)
	cashback, err := cashbackdistributor.NewDepositCashbackDistributor(addr, client)
	if err != nil {
		log.Printf("❌ Cashback Distributor: Failed to connect: %v", err)
		return
	}

	paused, err := cashback.Paused(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Printf("❌ Cashback Distributor: Failed to check status: %v", err)
		return
	}

	if !paused {
		fmt.Printf("⚠️  Cashback Distributor: Already active\n")
		return
	}

	if dryRun {
		fmt.Printf("🔍 Cashback Distributor: Would unpause (dry-run)\n")
		return
	}

	tx, err := cashback.Unpause(auth)
	if err != nil {
		log.Printf("❌ Cashback Distributor: Failed to unpause: %v", err)
		return
	}

	fmt.Printf("✅ Cashback Distributor: Unpaused (tx: %s)\n", tx.Hash().Hex())
	auth.Nonce = new(big.Int).Add(auth.Nonce, big.NewInt(1))
}

func pauseReferral(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, dryRun bool) {
	if constant.ReferralDistributorAddress == "" {
		fmt.Printf("⚠️  Referral Distributor: Not deployed\n")
		return
	}

	addr := common.HexToAddress(constant.ReferralDistributorAddress)
	referral, err := referraldistributor.NewReferralRewardDistributor(addr, client)
	if err != nil {
		log.Printf("❌ Referral Distributor: Failed to connect: %v", err)
		return
	}

	paused, err := referral.Paused(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Printf("❌ Referral Distributor: Failed to check status: %v", err)
		return
	}

	if paused {
		fmt.Printf("⚠️  Referral Distributor: Already paused\n")
		return
	}

	if dryRun {
		fmt.Printf("🔍 Referral Distributor: Would pause (dry-run)\n")
		return
	}

	tx, err := referral.Pause(auth)
	if err != nil {
		log.Printf("❌ Referral Distributor: Failed to pause: %v", err)
		return
	}

	fmt.Printf("✅ Referral Distributor: Paused (tx: %s)\n", tx.Hash().Hex())
	auth.Nonce = new(big.Int).Add(auth.Nonce, big.NewInt(1))
}

func unpauseReferral(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, dryRun bool) {
	if constant.ReferralDistributorAddress == "" {
		fmt.Printf("⚠️  Referral Distributor: Not deployed\n")
		return
	}

	addr := common.HexToAddress(constant.ReferralDistributorAddress)
	referral, err := referraldistributor.NewReferralRewardDistributor(addr, client)
	if err != nil {
		log.Printf("❌ Referral Distributor: Failed to connect: %v", err)
		return
	}

	paused, err := referral.Paused(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Printf("❌ Referral Distributor: Failed to check status: %v", err)
		return
	}

	if !paused {
		fmt.Printf("⚠️  Referral Distributor: Already active\n")
		return
	}

	if dryRun {
		fmt.Printf("🔍 Referral Distributor: Would unpause (dry-run)\n")
		return
	}

	tx, err := referral.Unpause(auth)
	if err != nil {
		log.Printf("❌ Referral Distributor: Failed to unpause: %v", err)
		return
	}

	fmt.Printf("✅ Referral Distributor: Unpaused (tx: %s)\n", tx.Hash().Hex())
	auth.Nonce = new(big.Int).Add(auth.Nonce, big.NewInt(1))
}

func getContracts(contractStr string) []string {
	if contractStr == "all" {
		return []string{"mining", "cashback", "referral"}
	}
	return []string{contractStr}
}
