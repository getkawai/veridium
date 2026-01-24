package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/cashbackdistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/referraldistributor"
	"github.com/kawai-network/veridium/pkg/alert"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/types"
)

var autoConfirm bool // Global flag for auto-confirmation

func main() {
	// Subcommands
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	allCmd := flag.NewFlagSet("all", flag.ExitOnError)

	// Flags for generate command
	var rewardType string
	generateCmd.StringVar(&rewardType, "type", "mining", "Reward type: mining, cashback, or referral")
	generateCmd.BoolVar(&autoConfirm, "auto-confirm", false, "Auto-confirm all prompts (for testing)")

	// Flags for upload command
	var uploadType string
	uploadCmd.StringVar(&uploadType, "type", "mining", "Reward type: mining, cashback, or referral")

	// Flags for all command
	allCmd.BoolVar(&autoConfirm, "auto-confirm", false, "Auto-confirm all prompts (for testing)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Initialize KV Store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatalf("Failed to initialize KV store: %v", err)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "generate":
		generateCmd.Parse(os.Args[2:])
		if err := generateSettlement(ctx, kv, rewardType); err != nil {
			log.Fatalf("Generate failed: %v", err)
		}

	case "upload":
		uploadCmd.Parse(os.Args[2:])
		if err := uploadMerkleRoot(ctx, kv, uploadType); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}

	case "status":
		statusCmd.Parse(os.Args[2:])
		if err := showStatus(ctx, kv); err != nil {
			log.Fatalf("Status failed: %v", err)
		}

	case "all":
		allCmd.Parse(os.Args[2:])
		if err := settleAll(ctx, kv); err != nil {
			log.Fatalf("Settle all failed: %v", err)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Unified Reward Settlement Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  reward-settlement generate --type <type>  Generate Merkle tree for specific reward type")
	fmt.Println("  reward-settlement upload --type <type>    Upload Merkle root to contract")
	fmt.Println("  reward-settlement status                  Show settlement status for all types")
	fmt.Println("  reward-settlement all                     Settle all reward types at once")
	fmt.Println("")
	fmt.Println("Reward Types:")
	fmt.Println("  mining     - Mining rewards (9-field Merkle tree)")
	fmt.Println("  cashback   - Deposit cashback rewards (3-field Merkle tree)")
	fmt.Println("  referral   - Referral commission rewards (3-field Merkle tree)")
	fmt.Println("  revenue    - Revenue sharing (USDT dividends, 3-field Merkle tree)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  reward-settlement generate --type mining")
	fmt.Println("  reward-settlement generate --type revenue")
	fmt.Println("  reward-settlement all")
	fmt.Println("")
	fmt.Println("Workflow:")
	fmt.Println("  1. Run 'generate' weekly to create Merkle trees")
	fmt.Println("  2. Run 'upload' to upload Merkle roots to contracts")
	fmt.Println("  3. Users can claim rewards via frontend")
	fmt.Println("")
	fmt.Println("  Or use 'all' to settle all reward types at once (recommended)")
}

func generateSettlement(ctx context.Context, kv *store.KVStore, rewardType string) error {
	log.Printf("🌳 Generating %s settlement...", rewardType)
	log.Println("")

	switch types.RewardType(rewardType) {
	case types.RewardTypeMining:
		return generateMiningSettlement(ctx, kv)
	case types.RewardTypeCashback:
		return generateCashbackSettlement(ctx, kv)
	case types.RewardTypeReferral:
		return generateReferralSettlement(ctx, kv)
	case types.RewardTypeRevenue:
		return generateRevenueSettlement(ctx, kv)
	default:
		return fmt.Errorf("unknown reward type: %s (must be mining, cashback, referral, or revenue)", rewardType)
	}
}

func generateMiningSettlement(ctx context.Context, kv store.Store) error {
	alerter := alert.NewTelegramAlert()

	log.Println("📊 Mining Rewards Settlement")
	log.Println("─────────────────────────────")

	// Send start notification
	alerter.SendAlert("INFO", "Settlement", "🔄 Starting mining settlement...")

	// Generate mining settlement with 9-field Merkle leaves
	period, err := kv.GenerateMiningSettlement(ctx, types.RewardTypeMining)
	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Mining settlement failed!\nError: %v", err))
		return fmt.Errorf("failed to generate mining settlement: %w", err)
	}

	log.Printf("✅ Mining settlement generated!")
	log.Println("")
	log.Printf("Period ID:     %d", period.PeriodID)
	log.Printf("Merkle Root:   %s", period.MerkleRoot)
	log.Printf("Contributors:  %d", period.ContributorCount)
	log.Printf("Total Amount:  %s KAWAI", period.TotalAmount)
	log.Printf("Proofs Saved:  %d", period.ProofsSaved)
	log.Printf("Status:        %s", period.Status)
	log.Println("")
	log.Printf("📝 Next: reward-settlement upload --type mining")

	// Send success notification
	alerter.SendAlert("SUCCESS", "Settlement",
		fmt.Sprintf("✅ Mining settlement complete!\nPeriod: %d\nContributors: %d\nTotal: %s KAWAI\nMerkle Root: %s",
			period.PeriodID, period.ContributorCount, period.TotalAmount, period.MerkleRoot[:16]+"..."))

	return nil
}

func generateCashbackSettlement(ctx context.Context, kv *store.KVStore) error {
	alerter := alert.NewTelegramAlert()

	log.Println("📊 Cashback Rewards Settlement")
	log.Println("─────────────────────────────")

	alerter.SendAlert("INFO", "Settlement", "🔄 Starting cashback settlement...")

	// Initialize cashback settlement
	settlement, err := blockchain.NewCashbackSettlement(kv, constant.GetAdminPrivateKey())
	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Cashback settlement init failed!\nError: %v", err))
		return fmt.Errorf("failed to initialize cashback settlement: %w", err)
	}

	// Get current period from contract (not from timestamp calculation)
	currentPeriod, err := settlement.GetCurrentPeriod(ctx)
	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Failed to get current period!\nError: %v", err))
		return fmt.Errorf("failed to get current period from contract: %w", err)
	}

	// Settle CURRENT period (not previous), because setMerkleRoot sets root for currentPeriod
	settlementPeriod := currentPeriod

	log.Printf("Current Period:    %d", currentPeriod)
	log.Printf("Settling Period:   %d", settlementPeriod)
	log.Println("")

	// Run settlement
	if err := settlement.SettleCashback(ctx, settlementPeriod); err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Cashback settlement failed!\nPeriod: %d\nError: %v", settlementPeriod, err))
		return fmt.Errorf("failed to settle cashback: %w", err)
	}

	log.Println("")
	log.Printf("✅ Cashback settlement completed!")
	log.Println("")
	log.Printf("📝 Next: reward-settlement upload --type cashback")

	alerter.SendAlert("SUCCESS", "Settlement",
		fmt.Sprintf("✅ Cashback settlement complete!\nPeriod: %d", settlementPeriod))

	return nil
}

func generateReferralSettlement(ctx context.Context, kv *store.KVStore) error {
	log.Println("📊 Referral Rewards Settlement")
	log.Println("─────────────────────────────")
	log.Println("")

	// Get current period (same as mining, weekly)
	// For referral, we settle based on mining periods
	// Get latest mining settlement period
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list settlement periods: %w", err)
	}

	if len(periods) == 0 {
		return fmt.Errorf("no mining settlement found - run mining settlement first")
	}

	// Get latest KAWAI period
	var latest *store.SettlementPeriod
	for i := len(periods) - 1; i >= 0; i-- {
		if periods[i].RewardType == types.RewardTypeMining {
			latest = periods[i]
			break
		}
	}

	if latest == nil {
		return fmt.Errorf("no KAWAI settlement found - run mining settlement first")
	}

	settlementPeriod := uint64(latest.PeriodID)

	log.Printf("Mining Period:     %d", settlementPeriod)
	log.Printf("Settling Referral: Period %d", settlementPeriod)
	log.Println("")

	// Initialize referral settlement
	settlement := blockchain.NewReferralSettlement(kv, constant.GetAdminPrivateKey())

	// Run settlement
	if err := settlement.SettleReferral(ctx, settlementPeriod); err != nil {
		return fmt.Errorf("failed to settle referral: %w", err)
	}

	log.Println("")
	log.Printf("✅ Referral settlement completed!")
	log.Println("")
	log.Printf("📝 Next: reward-settlement upload --type referral")

	return nil
}

func uploadMerkleRoot(ctx context.Context, kv *store.KVStore, rewardType string) error {
	log.Printf("🚀 Uploading %s Merkle root to contract...", rewardType)
	log.Println("")

	switch types.RewardType(rewardType) {
	case types.RewardTypeMining:
		return uploadMiningRoot(ctx, kv)
	case types.RewardTypeCashback:
		return uploadCashbackRoot(ctx, kv)
	case types.RewardTypeReferral:
		return uploadReferralRoot(ctx, kv)
	case types.RewardTypeRevenue:
		return uploadRevenueRoot(ctx, kv)
	default:
		return fmt.Errorf("unknown reward type: %s", rewardType)
	}
}

func uploadMiningRoot(ctx context.Context, kv store.Store) error {
	alerter := alert.NewTelegramAlert()

	// Get latest mining period
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list periods: %w", err)
	}

	var latest *store.SettlementPeriod
	for i := len(periods) - 1; i >= 0; i-- {
		if periods[i].RewardType == types.RewardTypeMining {
			latest = periods[i]
			break
		}
	}

	if latest == nil {
		return fmt.Errorf("no mining settlement found - run 'generate --type mining' first")
	}

	log.Printf("Period ID:     %d", latest.PeriodID)
	log.Printf("Merkle Root:   %s", latest.MerkleRoot)
	log.Printf("Total Amount:  %s KAWAI", latest.TotalAmount)
	log.Println("")

	alerter.SendAlert("INFO", "Settlement",
		fmt.Sprintf("📤 Uploading mining merkle root...\nPeriod: %d", latest.PeriodID))

	// Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Failed to connect to RPC!\nError: %v", err))
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// Load MiningRewardDistributor contract
	distributorAddr := common.HexToAddress(constant.MiningRewardDistributorAddress)
	distributor, err := miningdistributor.NewMiningRewardDistributor(distributorAddr, client)
	if err != nil {
		return fmt.Errorf("failed to load MiningRewardDistributor: %w", err)
	}

	// Get private key
	privateKeyHex := constant.GetAdminPrivateKey()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Context = ctx

	// Parse Merkle root
	merkleRootHex := latest.MerkleRoot
	if strings.HasPrefix(merkleRootHex, "0x") {
		merkleRootHex = merkleRootHex[2:]
	}
	merkleRootBytes := common.Hex2Bytes(merkleRootHex)
	if len(merkleRootBytes) != 32 {
		return fmt.Errorf("invalid Merkle root length: expected 32 bytes, got %d", len(merkleRootBytes))
	}
	var merkleRoot [32]byte
	copy(merkleRoot[:], merkleRootBytes)

	log.Printf("⚠️  About to upload Merkle root to MiningRewardDistributor")
	log.Printf("    Period ID: %d", latest.PeriodID)
	log.Printf("    Merkle Root: %s", latest.MerkleRoot)
	if !confirm("Continue with upload?") {
		return fmt.Errorf("upload cancelled by user")
	}
	log.Println("")

	// Get current on-chain period to determine upload strategy
	currentPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period: %w", err)
	}

	log.Printf("📊 Contract currentPeriod: %d, Settlement period: %d", currentPeriod.Int64(), latest.PeriodID)

	// Production-grade upload strategy based on period relationship
	if latest.PeriodID == currentPeriod.Int64() {
		// Update current period's root
		log.Printf("🌳 [MINING] Updating Merkle root for current period %d", latest.PeriodID)
		tx, err := distributor.SetMerkleRoot(auth, merkleRoot)
		if err != nil {
			alerter.SendAlert("ERROR", "Settlement",
				fmt.Sprintf("❌ Mining merkle root upload failed!\nPeriod: %d\nError: %v", latest.PeriodID, err))
			return fmt.Errorf("failed to upload Merkle root: %w", err)
		}
		log.Printf("✅ [MINING] SetMerkleRoot transaction sent: %s", tx.Hash().Hex())
		log.Println("⏳ [MINING] Waiting for confirmation...")

		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			alerter.SendAlert("ERROR", "Settlement",
				fmt.Sprintf("❌ Mining tx confirmation failed!\nTx: %s\nError: %v", tx.Hash().Hex(), err))
			return fmt.Errorf("failed to wait for confirmation: %w", err)
		}
		if receipt.Status != 1 {
			alerter.SendAlert("ERROR", "Settlement",
				fmt.Sprintf("❌ Mining tx reverted!\nTx: %s\nStatus: %d", tx.Hash().Hex(), receipt.Status))
			return fmt.Errorf("transaction failed with status: %d", receipt.Status)
		}
		log.Printf("✅ [MINING] Merkle root uploaded successfully in block %d", receipt.BlockNumber.Uint64())

		alerter.SendAlert("SUCCESS", "Settlement",
			fmt.Sprintf("✅ Mining merkle root uploaded!\nPeriod: %d\nTx: %s\nBlock: %d\nGas: %d",
				latest.PeriodID, tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed))

	} else if latest.PeriodID == currentPeriod.Int64()+1 {
		// Advance to next period
		log.Printf("🌳 [MINING] Advancing to period %d with new Merkle root", latest.PeriodID)
		tx, err := distributor.AdvancePeriod(auth, merkleRoot)
		if err != nil {
			alerter.SendAlert("ERROR", "Settlement",
				fmt.Sprintf("❌ Mining period advance failed!\nPeriod: %d\nError: %v", latest.PeriodID, err))
			return fmt.Errorf("failed to advance period: %w", err)
		}
		log.Printf("✅ [MINING] AdvancePeriod transaction sent: %s", tx.Hash().Hex())
		log.Println("⏳ [MINING] Waiting for confirmation...")

		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			alerter.SendAlert("ERROR", "Settlement",
				fmt.Sprintf("❌ Mining tx confirmation failed!\nTx: %s\nError: %v", tx.Hash().Hex(), err))
			return fmt.Errorf("failed to wait for confirmation: %w", err)
		}
		if receipt.Status != 1 {
			alerter.SendAlert("ERROR", "Settlement",
				fmt.Sprintf("❌ Mining tx reverted!\nTx: %s\nStatus: %d", tx.Hash().Hex(), receipt.Status))
			return fmt.Errorf("transaction failed with status: %d", receipt.Status)
		}
		log.Printf("✅ [MINING] Period advanced successfully in block %d", receipt.BlockNumber.Uint64())

		alerter.SendAlert("SUCCESS", "Settlement",
			fmt.Sprintf("✅ Mining period advanced!\nPeriod: %d\nTx: %s\nBlock: %d\nGas: %d",
				latest.PeriodID, tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed))

	} else {
		return fmt.Errorf("period mismatch: settlement period %d, contract period %d (expected %d or %d)",
			latest.PeriodID, currentPeriod.Int64(), currentPeriod.Int64(), currentPeriod.Int64()+1)
	}
	log.Println("")
	log.Printf("✅ Mining root upload completed!")
	log.Println("")
	log.Printf("📝 Next: Users can now claim mining rewards via UI")

	return nil
}

func uploadCashbackRoot(ctx context.Context, kv *store.KVStore) error {
	alerter := alert.NewTelegramAlert()

	// 1. Get latest cashback period
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list periods: %w", err)
	}

	var latest *store.SettlementPeriod
	for i := len(periods) - 1; i >= 0; i-- {
		// Note: Cashback generates periods with RewardType="cashback"
		if periods[i].RewardType == types.RewardTypeCashback {
			latest = periods[i]
			break
		}
	}

	if latest == nil {
		return fmt.Errorf("no cashback settlement found - run 'generate --type cashback' first")
	}

	log.Printf("Period ID:     %d", latest.PeriodID)
	log.Printf("Merkle Root:   %s", latest.MerkleRoot)
	log.Printf("Total Records: %d", latest.ProofsSaved)
	log.Println("")

	alerter.SendAlert("INFO", "Settlement",
		fmt.Sprintf("📤 Uploading cashback merkle root...\nPeriod: %d", latest.PeriodID))

	// 2. Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Failed to connect to RPC!\nError: %v", err))
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// 3. Load DepositCashbackDistributor contract
	// We use the generated binding
	distributorAddr := common.HexToAddress(constant.CashbackDistributorAddress)
	distributor, err := cashbackdistributor.NewDepositCashbackDistributor(distributorAddr, client)
	if err != nil {
		return fmt.Errorf("failed to load DepositCashbackDistributor: %w", err)
	}

	// 4. Setup Transactor
	privateKeyHex := constant.GetAdminPrivateKey()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Context = ctx

	// 5. Prepare Merkle Root
	merkleRootHex := latest.MerkleRoot
	if strings.HasPrefix(merkleRootHex, "0x") {
		merkleRootHex = merkleRootHex[2:]
	}
	merkleRootBytes := common.Hex2Bytes(merkleRootHex)
	if len(merkleRootBytes) != 32 {
		return fmt.Errorf("invalid Merkle root length: expected 32 bytes, got %d", len(merkleRootBytes))
	}
	var merkleRoot [32]byte
	copy(merkleRoot[:], merkleRootBytes)

	// User Confirmation
	log.Printf("⚠️  About to upload Merkle root to DepositCashbackDistributor")
	log.Printf("    Period ID: %d", latest.PeriodID)
	log.Printf("    Merkle Root: %s", latest.MerkleRoot)
	if !confirm("Continue with upload?") {
		return fmt.Errorf("upload cancelled by user")
	}
	log.Println("")

	// 6. Check Current Contract Period to decide logic
	// The contract has currentPeriod() method
	contractPeriodBig, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period from contract: %w", err)
	}
	contractPeriod := contractPeriodBig.Uint64()

	log.Printf("📊 Contract currentPeriod: %d, Settlement period: %d", contractPeriod, latest.PeriodID)

	var tx *ethtypes.Transaction

	if uint64(latest.PeriodID) == contractPeriod {
		// Update current period (SetMerkleRoot)
		log.Printf("🌳 [CASHBACK] Updating Merkle root for current period %d", latest.PeriodID)
		tx, err = distributor.SetMerkleRoot(auth, merkleRoot)
	} else if uint64(latest.PeriodID) == contractPeriod+1 {
		// Advance to next period (AdvancePeriod)
		log.Printf("🌳 [CASHBACK] Advancing to period %d with new Merkle root", latest.PeriodID)
		tx, err = distributor.AdvancePeriod(auth, merkleRoot)
	} else {
		return fmt.Errorf("period mismatch: settlement period %d, contract period %d", latest.PeriodID, contractPeriod)
	}

	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Cashback upload failed!\nPeriod: %d\nError: %v", latest.PeriodID, err))
		return fmt.Errorf("transaction failed: %w", err)
	}

	log.Printf("✅ [CASHBACK] Transaction sent: %s", tx.Hash().Hex())
	log.Println("⏳ [CASHBACK] Waiting for confirmation...")

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for confirmation: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("transaction reverted with status: %d", receipt.Status)
	}

	log.Printf("✅ [CASHBACK] Merkle root uploaded successfully in block %d", receipt.BlockNumber.Uint64())

	alerter.SendAlert("SUCCESS", "Settlement",
		fmt.Sprintf("✅ Cashback merkle root uploaded!\nPeriod: %d\nTx: %s",
			latest.PeriodID, tx.Hash().Hex()))

	log.Println("")
	log.Printf("📝 Next: Users can now claim cashback rewards via UI")

	return nil
}

func uploadReferralRoot(ctx context.Context, kv *store.KVStore) error {
	alerter := alert.NewTelegramAlert()

	// 1. Get latest mining period (Referral settlement piggybacks on mining periods)
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to list periods: %w", err)
	}

	var latestMining *store.SettlementPeriod
	for i := len(periods) - 1; i >= 0; i-- {
		// Referral tracks mining periods (usually "kawai" type for mining rewards)
		if periods[i].RewardType == types.RewardTypeMining {
			latestMining = periods[i]
			break
		}
	}

	if latestMining == nil {
		return fmt.Errorf("no mining settlement found - cannot determine referral period")
	}

	periodID := uint64(latestMining.PeriodID)

	// 2. Read Referral Merkle Root from KV
	// Key format from pkg/blockchain/referral_settlement.go: referral_period:%d:merkle_root
	rootKey := fmt.Sprintf("referral_period:%d:merkle_root", periodID)
	rootDataJSON, err := kv.GetMarketplaceData(ctx, rootKey)
	if err != nil {
		return fmt.Errorf("failed to get referral merkle root for period %d: %w", periodID, err)
	}

	var merkleRootHex string
	if err := json.Unmarshal(rootDataJSON, &merkleRootHex); err != nil {
		return fmt.Errorf("failed to unmarshal referral merkle root: %w", err)
	}

	log.Printf("Period ID:     %d", periodID)
	log.Printf("Merkle Root:   %s", merkleRootHex)
	log.Println("")

	alerter.SendAlert("INFO", "Settlement",
		fmt.Sprintf("📤 Uploading referral merkle root...\nPeriod: %d", periodID))

	// 3. Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Failed to connect to RPC!\nError: %v", err))
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// 4. Load ReferralRewardDistributor contract
	distributorAddr := common.HexToAddress(constant.ReferralDistributorAddress)
	distributor, err := referraldistributor.NewReferralRewardDistributor(distributorAddr, client)
	if err != nil {
		return fmt.Errorf("failed to load ReferralRewardDistributor: %w", err)
	}

	// 5. Setup Transactor
	privateKeyHex := constant.GetAdminPrivateKey()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Context = ctx

	// 6. Prepare Merkle Root
	if strings.HasPrefix(merkleRootHex, "0x") {
		merkleRootHex = merkleRootHex[2:]
	}
	merkleRootBytes := common.Hex2Bytes(merkleRootHex)
	if len(merkleRootBytes) != 32 {
		return fmt.Errorf("invalid Merkle root length: expected 32 bytes, got %d", len(merkleRootBytes))
	}
	var merkleRoot [32]byte
	copy(merkleRoot[:], merkleRootBytes)

	// User Confirmation
	log.Printf("⚠️  About to upload Merkle root to ReferralRewardDistributor")
	log.Printf("    Period ID: %d", periodID)
	log.Printf("    Merkle Root: %s", merkleRootHex)
	if !confirm("Continue with upload?") {
		return fmt.Errorf("upload cancelled by user")
	}
	log.Println("")

	// 7. Check Current Contract Period to decide logic
	contractPeriodBig, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period from contract: %w", err)
	}
	contractPeriod := contractPeriodBig.Uint64()

	log.Printf("📊 Contract currentPeriod: %d, Settlement period: %d", contractPeriod, periodID)

	var tx *ethtypes.Transaction

	if periodID == contractPeriod {
		// Update current period (SetMerkleRoot)
		log.Printf("🌳 [REFERRAL] Updating Merkle root for current period %d", periodID)
		tx, err = distributor.SetMerkleRoot(auth, merkleRoot)
	} else if periodID == contractPeriod+1 {
		// Advance to next period (AdvancePeriod)
		log.Printf("🌳 [REFERRAL] Advancing to period %d with new Merkle root", periodID)
		tx, err = distributor.AdvancePeriod(auth, merkleRoot)
	} else {
		return fmt.Errorf("period mismatch: settlement period %d, contract period %d", periodID, contractPeriod)
	}

	if err != nil {
		alerter.SendAlert("ERROR", "Settlement",
			fmt.Sprintf("❌ Referral upload failed!\nPeriod: %d\nError: %v", periodID, err))
		return fmt.Errorf("transaction failed: %w", err)
	}

	log.Printf("✅ [REFERRAL] Transaction sent: %s", tx.Hash().Hex())
	log.Println("⏳ [REFERRAL] Waiting for confirmation...")

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for confirmation: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("transaction reverted with status: %d", receipt.Status)
	}

	log.Printf("✅ [REFERRAL] Merkle root uploaded successfully in block %d", receipt.BlockNumber.Uint64())

	alerter.SendAlert("SUCCESS", "Settlement",
		fmt.Sprintf("✅ Referral merkle root uploaded!\nPeriod: %d\nTx: %s",
			periodID, tx.Hash().Hex()))

	log.Println("")
	log.Printf("📝 Next: Users can now claim referral rewards via UI")

	return nil
}

func generateRevenueSettlement(ctx context.Context, kv *store.KVStore) error {
	log.Println("📊 Revenue Sharing Settlement (USDT Dividends)")
	log.Println("─────────────────────────────")
	log.Println("")

	// Initialize revenue settlement
	settlement, err := blockchain.NewRevenueSettlement(kv)
	if err != nil {
		return fmt.Errorf("failed to initialize revenue settlement: %w", err)
	}

	// Step 1: Generate settlement
	log.Println("Step 1: Generating revenue settlement...")
	log.Println("")

	// Get current period (shared across all reward types: mining, cashback, referral, revenue)
	// Period system: Weekly periods starting Jan 1, 2025, incrementing every Monday 00:00 UTC
	// Settlement always processes previous period (currentPeriod - 1)
	currentPeriod := settlement.GetCurrentPeriod()
	settlementPeriod := currentPeriod - 1

	if settlementPeriod < 1 {
		return fmt.Errorf("no period to settle yet (current period: %d)", currentPeriod)
	}

	log.Printf("Current Period:    %d", currentPeriod)
	log.Printf("Settling Period:   %d", settlementPeriod)
	log.Println("")

	merkleRoot, err := settlement.SettleRevenue(ctx, settlementPeriod)
	if err != nil {
		return fmt.Errorf("settlement generation failed: %w", err)
	}

	log.Printf("✅ Settlement generated successfully")
	log.Printf("Merkle Root: 0x%x", merkleRoot)
	log.Println("")

	// Step 2: Get amount
	log.Println("Step 2: Getting vault balance...")
	log.Println("")

	amount, err := settlement.GetPaymentVaultBalance(ctx)
	if err != nil {
		return fmt.Errorf("failed to get vault balance: %w", err)
	}

	log.Printf("Total Revenue: %s USDT", amount.String())
	log.Println("")

	// Step 3: Withdraw USDT
	log.Println("Step 3: Withdrawing USDT to distributor...")
	log.Println("")

	log.Printf("⚠️  About to withdraw %s USDT to USDT_Distributor", amount.String())
	if !confirm("Continue with withdrawal?") {
		return fmt.Errorf("withdrawal cancelled by user")
	}
	log.Println("")

	if err := settlement.WithdrawToDistributor(ctx, amount); err != nil {
		return fmt.Errorf("withdrawal failed: %w", err)
	}

	log.Printf("✅ USDT withdrawn successfully")
	log.Println("")

	// Step 4: Upload merkle root
	log.Println("Step 4: Uploading merkle root...")
	log.Println("")

	log.Printf("⚠️  About to upload merkle root: 0x%x", merkleRoot)
	if !confirm("Continue with upload?") {
		return fmt.Errorf("upload cancelled by user")
	}
	log.Println("")

	if err := settlement.UploadMerkleRoot(ctx, merkleRoot); err != nil {
		return fmt.Errorf("merkle root upload failed: %w", err)
	}

	log.Printf("✅ Merkle root uploaded successfully")
	log.Println("")
	log.Printf("✅ Revenue settlement completed!")
	log.Println("")
	log.Printf("📝 Next: reward-settlement upload --type revenue")

	return nil
}

func uploadRevenueRoot(ctx context.Context, kv *store.KVStore) error {
	log.Println("⚠️  Revenue root upload already done during generate")
	log.Println("")
	log.Println("Revenue settlement uploads the Merkle root automatically")
	log.Println("during the generate step (after withdrawal confirmation).")
	log.Println("")
	log.Println("If you need to re-upload, use the manual command:")
	log.Println("")
	log.Println("📋 Manual upload command:")
	log.Printf("   cast send 0xE964B52D496F37749bd0caF287A356afdC10836C \\")
	log.Printf("     'setMerkleRoot(bytes32)' <MERKLE_ROOT> \\")
	log.Printf("     --rpc-url $RPC_URL --private-key <ADMIN_PRIVATE_KEY>")

	return nil
}

// Helper function for user confirmation
func confirm(prompt string) bool {
	if autoConfirm {
		log.Printf("✓ Auto-confirmed: %s", prompt)
		return true
	}

	fmt.Printf("%s (y/n): ", prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func showStatus(ctx context.Context, kv *store.KVStore) error {
	log.Println("📊 Reward Settlement Status")
	log.Println("═══════════════════════════════════════════════════════════════")
	log.Println("")

	// Mining status
	log.Println("⛏️  MINING REWARDS")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := showMiningStatus(ctx, kv); err != nil {
		log.Printf("Error: %v", err)
	}
	log.Println("")

	// Cashback status
	log.Println("💰 CASHBACK REWARDS")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := showCashbackStatus(ctx, kv); err != nil {
		log.Printf("Error: %v", err)
	}
	log.Println("")

	// Referral status
	log.Println("🤝 REFERRAL REWARDS")
	log.Println("───────────────────────────────────────────────────────────────")
	log.Println("Status: ✅ fully operational")
	log.Println("Use 'generate --type referral' and 'upload --type referral'")
	log.Println("")

	// Revenue status
	log.Println("💸 REVENUE SHARING")
	log.Println("───────────────────────────────────────────────────────────────")
	log.Println("Status: ✅ fully operational (upload integrated in generate)")
	log.Println("")

	log.Println("═══════════════════════════════════════════════════════════════")
	log.Println("💡 Tip: Run 'reward-settlement all' to settle all types at once")

	return nil
}

func showMiningStatus(ctx context.Context, kv store.Store) error {
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		return err
	}

	var kawaiPeriods []*store.SettlementPeriod
	for _, p := range periods {
		if p.RewardType == types.RewardTypeMining {
			kawaiPeriods = append(kawaiPeriods, p)
		}
	}

	if len(kawaiPeriods) == 0 {
		log.Println("No settlements found")
		return nil
	}

	// Show last 3
	start := 0
	if len(kawaiPeriods) > 3 {
		start = len(kawaiPeriods) - 3
	}

	fmt.Printf("%-12s | %-20s | %-15s | %s\n", "Period ID", "Status", "Total Amount", "Proofs")
	log.Println("───────────────────────────────────────────────────────────────")
	for _, p := range kawaiPeriods[start:] {
		fmt.Printf("%-12d | %-20s | %-15s | %d\n",
			p.PeriodID, p.Status, p.TotalAmount+" KAWAI", p.ProofsSaved)
	}

	log.Printf("\nTotal Settlements: %d", len(kawaiPeriods))
	return nil
}

func showCashbackStatus(ctx context.Context, kv *store.KVStore) error {
	currentPeriod := kv.GetCurrentPeriod()
	log.Printf("Current Period: %d", currentPeriod)
	log.Printf("Next Settlement: Period %d (in %s)", currentPeriod-1, "next Monday")
	log.Println("")
	log.Println("Recent Settlements: (TODO: implement)")
	return nil
}

func settleAll(ctx context.Context, kv *store.KVStore) error {
	log.Println("🚀 Settling All Reward Types")
	log.Println("═══════════════════════════════════════════════════════════════")
	log.Println("")

	success := 0
	failed := 0

	// 1. Mining
	log.Println("1️⃣  Mining Rewards")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := generateMiningSettlement(ctx, kv); err != nil {
		log.Printf("❌ Mining settlement failed: %v", err)
		failed++
	} else {
		success++
	}
	log.Println("")

	// 2. Cashback
	log.Println("2️⃣  Cashback Rewards")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := generateCashbackSettlement(ctx, kv); err != nil {
		log.Printf("❌ Cashback settlement failed: %v", err)
		failed++
	} else {
		success++
	}
	log.Println("")

	// 3. Referral
	log.Println("3️⃣  Referral Rewards")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := generateReferralSettlement(ctx, kv); err != nil {
		log.Printf("⚠️  Referral settlement skipped: %v", err)
		// Don't count as failed since it's not implemented yet
	}
	log.Println("")

	// 4. Revenue Sharing
	log.Println("4️⃣  Revenue Sharing (USDT Dividends)")
	log.Println("───────────────────────────────────────────────────────────────")
	if err := generateRevenueSettlement(ctx, kv); err != nil {
		log.Printf("❌ Revenue settlement failed: %v", err)
		failed++
	} else {
		success++
	}
	log.Println("")

	// Summary
	log.Println("═══════════════════════════════════════════════════════════════")
	log.Printf("✅ Successful: %d", success)
	log.Printf("❌ Failed: %d", failed)
	log.Println("")

	if failed > 0 {
		return fmt.Errorf("some settlements failed")
	}

	log.Println("🎉 All settlements completed successfully!")
	log.Println("")
	log.Println("📝 Next steps:")
	log.Println("  1. Upload Merkle roots to contracts")
	log.Println("  2. Users can claim rewards via frontend")

	return nil
}
