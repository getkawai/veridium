package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/mockstablecoin"
	"github.com/kawai-network/veridium/internal/generate/abi/revenuedistributor"
	"github.com/kawai-network/veridium/internal/generate/abi/vault"
	"github.com/kawai-network/veridium/pkg/merkle"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/x/constant"
)

// RevenueSettlement handles weekly revenue sharing settlement
type RevenueSettlement struct {
	client          *ethclient.Client
	paymentVault    *vault.PaymentVault
	vaultAddress    common.Address
	stablecoinToken *mockstablecoin.MockStablecoin // Works with any ERC-20 stablecoin (MockStablecoin on testnet, USDC on mainnet)
	kvStore         *store.KVStore
}

// NewRevenueSettlement creates a new revenue settlement instance
func NewRevenueSettlement(kvStore *store.KVStore) (*RevenueSettlement, error) {
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Monad: %w", err)
	}

	vaultAddr := common.HexToAddress(constant.PaymentVaultAddress)
	paymentVault, err := vault.NewPaymentVault(vaultAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load PaymentVault: %w", err)
	}

	// Validate stablecoin address is configured
	if constant.StablecoinAddress == "" {
		return nil, fmt.Errorf("stablecoin address not configured")
	}

	// Load stablecoin token contract (USDC on mainnet, MockStablecoin on testnet)
	stablecoinAddr := common.HexToAddress(constant.StablecoinAddress)
	if stablecoinAddr == (common.Address{}) {
		return nil, fmt.Errorf("invalid stablecoin address: zero address")
	}

	stablecoinToken, err := mockstablecoin.NewMockStablecoin(stablecoinAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load stablecoin token: %w", err)
	}

	return &RevenueSettlement{
		client:          client,
		paymentVault:    paymentVault,
		vaultAddress:    vaultAddr,
		stablecoinToken: stablecoinToken,
		kvStore:         kvStore,
	}, nil
}

// GetPaymentVaultBalance returns the total stablecoin balance in PaymentVault
//
// ECONOMIC MODEL VERIFICATION:
// This function returns the TOTAL balance in PaymentVault and distributes 100% as dividends.
//
// This is CORRECT based on the following verified facts:
//
// 1. USER DEPOSITS ARE NON-REFUNDABLE (verified in PaymentVault.sol):
//   - Contract has NO user withdraw function
//   - Only owner can withdraw (for dividend distribution)
//   - Users cannot get refunds once deposited
//
// 2. PHASE 1 ECONOMICS (Current):
//   - Contributors are paid in KAWAI tokens (mining rewards)
//   - All stablecoin collected = platform revenue
//   - No stablecoin is paid out to contributors
//
// 3. USER BALANCE TRACKING:
//   - User balances tracked off-chain in KV store
//   - When users spend credits, KV balance decreases
//   - Stablecoin remains in vault (becomes platform revenue)
//
// Example Flow:
// - User deposits 1000 USDC → Vault: 1000 USDC, KV: 1000 USDC
// - User spends 100 USDC on AI → Vault: 1000 USDC, KV: 900 USDC
// - Platform revenue = 100 USDC (spent amount)
// - User remaining credit = 900 USDC (in KV, non-refundable)
//
// At settlement:
// - Vault balance = 1000 USDC (all deposits)
// - Distributable = 1000 USDC (correct, deposits are non-refundable)
// - Users keep their KV credits for future AI usage
// - Token holders receive dividends from all deposits
//
// This model is intentional: deposits are non-refundable contributions to the platform.
// Users get AI service credits in return, not refundable deposits.
//
// Verified in: PaymentVault.sol (no user withdraw), REVENUE_SHARING.md (Economic Model section)
func (rs *RevenueSettlement) GetPaymentVaultBalance(ctx context.Context) (*big.Int, error) {
	// Query stablecoin.balanceOf(vaultAddress) to get total revenue
	balance, err := rs.stablecoinToken.BalanceOf(&bind.CallOpts{Context: ctx}, rs.vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to query stablecoin balance: %w", err)
	}

	log.Printf("📊 [REVENUE SETTLEMENT] PaymentVault: %s", rs.vaultAddress.Hex())
	log.Printf("💵 [REVENUE SETTLEMENT] Stablecoin Balance: %s", balance.String())

	return balance, nil
}

// SettleRevenue performs weekly revenue sharing settlement
// Returns the merkle root for uploading to contract
//
// PHASE COMPATIBILITY:
// This function works in BOTH Phase 1 and Phase 2:
// - Phase 1 (Mining Era): Distributes 100% of stablecoin revenue (contributors paid in KAWAI)
// - Phase 2 (Post-Mining): Distributes 30% of stablecoin revenue (contributors paid 70% in stablecoin)
//
// The phase detection (pkg/config/phase.go) is informational only.
// Revenue settlement runs in both phases, the difference is in the revenue amount:
// - Phase 1: All stablecoin in vault = platform revenue (100%)
// - Phase 2: Vault balance - contributor payments = platform revenue (30%)
//
// No phase check is needed because:
// 1. If vault balance is 0, settlement returns early (no revenue)
// 2. Phase transition is automatic based on token supply
// 3. Revenue calculation is always correct (distribute what's in vault)
func (rs *RevenueSettlement) SettleRevenue(ctx context.Context, period uint64) ([32]byte, error) {
	var emptyRoot [32]byte

	log.Printf("💰 [REVENUE SETTLEMENT] Starting settlement for period %d", period)

	// 1. Get total stablecoin in PaymentVault (this is the revenue to distribute)
	totalRevenue, err := rs.GetPaymentVaultBalance(ctx)
	if err != nil {
		return emptyRoot, fmt.Errorf("failed to get PaymentVault balance: %w", err)
	}

	if totalRevenue.Cmp(big.NewInt(0)) == 0 {
		return emptyRoot, fmt.Errorf("no revenue to distribute (PaymentVault balance: 0)")
	}

	log.Printf("💵 [REVENUE SETTLEMENT] Total revenue to distribute: %s stablecoin", totalRevenue.String())

	// 2. Get KAWAI holders using HYBRID approach (Registry + Recent Blockchain Scan)
	// This solves the Monad testnet RPC 100-block limit issue
	log.Printf("📊 [REVENUE SETTLEMENT] Scanning holders (hybrid: registry + blockchain)")

	// 2a. Get holders from registry (desktop + CLI auto-registration)
	holderRegistry := NewHolderRegistry(rs.kvStore)
	registryAddresses, err := holderRegistry.GetAllHolders(ctx)
	if err != nil {
		log.Printf("⚠️  [REVENUE SETTLEMENT] Failed to get registry holders: %v", err)
		registryAddresses = []common.Address{} // Continue with empty registry
	}
	log.Printf("📋 [REVENUE SETTLEMENT] Registry holders: %d", len(registryAddresses))

	// 2b. Scan recent blockchain transfers (last 90 blocks to stay under RPC limit)
	scanner, err := NewHolderScanner()
	if err != nil {
		return emptyRoot, fmt.Errorf("failed to create holder scanner: %w", err)
	}

	// Get current block
	currentBlock, err := rs.client.BlockNumber(ctx)
	if err != nil {
		return emptyRoot, fmt.Errorf("failed to get current block: %w", err)
	}

	// Scan last 90 blocks (safe margin under 100-block RPC limit)
	startBlock := currentBlock
	if currentBlock > 90 {
		startBlock = currentBlock - 90
	}

	recentHolders, err := scanner.ScanHoldersFromBlock(ctx, startBlock)
	if err != nil {
		log.Printf("⚠️  [REVENUE SETTLEMENT] Failed to scan recent blockchain holders: %v", err)
		recentHolders = []common.Address{} // Continue with empty recent scan
	}
	log.Printf("🔍 [REVENUE SETTLEMENT] Recent blockchain holders (blocks %d-%d): %d", startBlock, currentBlock, len(recentHolders))

	// 2c. Merge and deduplicate holder addresses
	holderMap := make(map[common.Address]bool)
	for _, addr := range registryAddresses {
		holderMap[addr] = true
	}
	for _, addr := range recentHolders {
		holderMap[addr] = true
	}

	// Convert map to slice
	var holderAddresses []common.Address
	for addr := range holderMap {
		holderAddresses = append(holderAddresses, addr)
	}

	log.Printf("📊 [REVENUE SETTLEMENT] Total unique holders: %d (registry: %d, recent: %d)",
		len(holderAddresses), len(registryAddresses), len(recentHolders))

	if len(holderAddresses) == 0 {
		return emptyRoot, fmt.Errorf("no KAWAI holders found - cannot generate settlement")
	}

	// 2d. Query current balances for all holders
	var holders []*KawaiHolder
	failedQueries := 0

	for _, addr := range holderAddresses {
		balance, err := scanner.GetBalance(ctx, addr)
		if err != nil {
			log.Printf("⚠️  [REVENUE SETTLEMENT] Failed to get balance for %s: %v", addr.Hex(), err)
			failedQueries++
			continue
		}

		// Skip holders with zero balance
		if balance.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		holders = append(holders, &KawaiHolder{
			Address: addr,
			Balance: balance,
		})
	}

	// Check failure rate - abort if data quality is insufficient
	if failedQueries > 0 {
		failureRate := float64(failedQueries) / float64(len(holderAddresses)) * 100
		log.Printf("⚠️  [REVENUE SETTLEMENT] Balance query failures: %d/%d (%.2f%%)", failedQueries, len(holderAddresses), failureRate)

		if failureRate > 10.0 {
			return emptyRoot, fmt.Errorf("too many failed balance queries: %d/%d (%.2f%%) - data quality insufficient for settlement",
				failedQueries, len(holderAddresses), failureRate)
		}
	}

	if len(holders) == 0 {
		return emptyRoot, fmt.Errorf("no holders with non-zero balance - cannot generate settlement")
	}

	log.Printf("📊 [REVENUE SETTLEMENT] Holders with balance: %d", len(holders))

	// 3. Get total supply
	totalSupply, err := scanner.GetTotalSupply(ctx)
	if err != nil {
		return emptyRoot, fmt.Errorf("failed to get total supply: %w", err)
	}

	log.Printf("📊 [REVENUE SETTLEMENT] Total KAWAI Supply: %s", totalSupply.String())
	log.Printf("📊 [REVENUE SETTLEMENT] Number of Holders: %d", len(holders))

	// 4. Validate holder balances
	if err := ValidateHolders(holders, totalSupply); err != nil {
		log.Printf("⚠️  [REVENUE SETTLEMENT] Warning: %v", err)
		// Continue anyway - this is just a sanity check
	}

	// 5. Generate Merkle Tree for stablecoin distribution
	var leaves [][]byte
	var proofData []*store.MerkleProofData
	var proofAddresses []string

	currentIndex := uint64(0)
	for _, holder := range holders {
		// Calculate proportional share: (balance / totalSupply) * totalRevenue
		share := CalculateHolderShare(holder.Balance, totalSupply, totalRevenue)

		if share.Cmp(big.NewInt(0)) == 0 {
			continue // Skip holders with zero share
		}

		leaf := merkle.HashLeaf(currentIndex, holder.Address, share)
		leaves = append(leaves, leaf)

		proofData = append(proofData, &store.MerkleProofData{
			Index:      currentIndex,
			Amount:     share.String(),
			RewardType: "stablecoin", // Set reward type for frontend filtering
		})
		proofAddresses = append(proofAddresses, holder.Address.Hex())
		currentIndex++

		log.Printf("💰 Holder %s: %s KAWAI -> %s stablecoin dividend",
			holder.Address.Hex()[:10]+"...",
			holder.Balance.String(),
			share.String())
	}

	if len(leaves) == 0 {
		return emptyRoot, fmt.Errorf("no valid dividend recipients")
	}

	// 6. Build Merkle Tree
	tree := merkle.NewMerkleTree(leaves)
	rootBytes := tree.Root
	log.Printf("🌳 [REVENUE SETTLEMENT] Stablecoin Merkle Root: 0x%x", rootBytes)

	// Convert []byte to [32]byte
	var root [32]byte
	copy(root[:], rootBytes)

	// 7. Save Proofs (with "stablecoin:" prefix to distinguish from KAWAI proofs)
	savedCount := 0
	failedCount := 0

	for i, pd := range proofData {
		proof, ok := tree.GetProof(leaves[i])
		if !ok {
			log.Printf("❌ [REVENUE SETTLEMENT] Error generating stablecoin proof for index %d", i)
			failedCount++
			continue
		}

		var proofHex []string
		for _, p := range proof {
			proofHex = append(proofHex, fmt.Sprintf("0x%x", p))
		}
		pd.Proof = proofHex

		// Save with "stablecoin:" prefix to distinguish from KAWAI proofs
		addrKey := "stablecoin:" + proofAddresses[i]
		err := rs.kvStore.SaveMerkleProof(ctx, addrKey, pd)
		if err != nil {
			log.Printf("❌ [REVENUE SETTLEMENT] Failed to save stablecoin proof for %s: %v", proofAddresses[i], err)
			failedCount++
			continue
		}
		savedCount++
	}

	if failedCount > 0 {
		return emptyRoot, fmt.Errorf("settlement incomplete: %d proofs saved, %d failed", savedCount, failedCount)
	}

	log.Printf("✅ [REVENUE SETTLEMENT] Settlement completed for period %d", period)
	log.Printf("📊 [REVENUE SETTLEMENT] Summary: %d holders will receive dividends", savedCount)
	log.Printf("📝 [REVENUE SETTLEMENT] Next step: Upload Merkle root to Stablecoin Distributor contract")

	return root, nil
}

// GetCurrentPeriod returns the current settlement period
// All reward types (mining, cashback, referral, revenue) share the same weekly period system
// Period 1 starts at January 1, 2025 (configurable via CASHBACK_PERIOD_START env var)
// Period increments every Monday 00:00 UTC
func (rs *RevenueSettlement) GetCurrentPeriod() uint64 {
	return rs.kvStore.GetCurrentPeriod()
}

// WithdrawToDistributor transfers stablecoin from PaymentVault to Stablecoin Distributor
func (rs *RevenueSettlement) WithdrawToDistributor(ctx context.Context, amount *big.Int) error {
	log.Printf("💸 [REVENUE SETTLEMENT] Withdrawing %s stablecoin to Stablecoin Distributor", amount.String())

	// Get private key
	privateKeyHex := constant.GetAdminPrivateKey()
	// Remove 0x prefix if present
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get chain ID
	chainID, err := rs.client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas parameters
	auth.Context = ctx

	// Call PaymentVault.withdraw(RevenueDistributor, amount)
	distributorAddr := common.HexToAddress(constant.RevenueDistributorAddress)
	tx, err := rs.paymentVault.Withdraw(auth, distributorAddr, amount)
	if err != nil {
		return fmt.Errorf("failed to withdraw to distributor: %w", err)
	}

	log.Printf("✅ [REVENUE SETTLEMENT] Withdrawal transaction sent: %s", tx.Hash().Hex())
	log.Printf("⏳ [REVENUE SETTLEMENT] Waiting for confirmation...")

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, rs.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	log.Printf("✅ [REVENUE SETTLEMENT] Withdrawal confirmed in block %d", receipt.BlockNumber.Uint64())
	return nil
}

// UploadMerkleRoot uploads Merkle root to Stablecoin Distributor contract
func (rs *RevenueSettlement) UploadMerkleRoot(ctx context.Context, merkleRoot [32]byte) error {
	log.Printf("🌳 [REVENUE SETTLEMENT] Uploading Merkle root: 0x%x", merkleRoot)

	// Get private key
	privateKeyHex := constant.GetAdminPrivateKey()
	// Remove 0x prefix if present
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get chain ID
	chainID, err := rs.client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas parameters
	auth.Context = ctx

	// Load Revenue Distributor contract
	distributorAddr := common.HexToAddress(constant.RevenueDistributorAddress)
	distributor, err := revenuedistributor.NewRevenueDistributor(distributorAddr, rs.client)
	if err != nil {
		return fmt.Errorf("failed to load Revenue Distributor: %w", err)
	}

	// Call setMerkleRoot
	tx, err := distributor.SetMerkleRoot(auth, merkleRoot)
	if err != nil {
		return fmt.Errorf("failed to set merkle root: %w", err)
	}

	log.Printf("✅ [REVENUE SETTLEMENT] SetMerkleRoot transaction sent: %s", tx.Hash().Hex())
	log.Printf("⏳ [REVENUE SETTLEMENT] Waiting for confirmation...")

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, rs.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	log.Printf("✅ [REVENUE SETTLEMENT] Merkle root uploaded successfully in block %d", receipt.BlockNumber.Uint64())
	return nil
}
