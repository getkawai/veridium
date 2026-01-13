package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/kawai-network/veridium/pkg/jarvis/contracts"
	"github.com/kawai-network/veridium/pkg/jarvis/networks"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
	"github.com/kawai-network/veridium/pkg/store"
)

// monadChainID is the chain ID for Monad Testnet
var monadChainID = big.NewInt(int64(networks.MonadTestnet.GetChainID()))

// ClaimableReward represents a single claimable reward proof
type ClaimableReward struct {
	Index       uint64   `json:"index"`
	Amount      string   `json:"amount"`        // BigInt as string (raw, no decimals)
	Proof       []string `json:"proof"`         // Merkle proof (hex strings)
	MerkleRoot  string   `json:"merkle_root"`   // Root hash for verification
	PeriodID    int64    `json:"period_id"`     // Settlement period identifier
	RewardType  string   `json:"reward_type"`   // "kawai" or "usdt"
	ClaimStatus string   `json:"claim_status"`  // "unclaimed", "pending", "confirmed", "failed"
	ClaimTxHash string   `json:"claim_tx_hash"` // Transaction hash if claimed
	CreatedAt   string   `json:"created_at"`    // When proof was generated
	ClaimedAt   string   `json:"claimed_at"`    // When claimed (if confirmed)
	Formatted   string   `json:"formatted"`     // Human-readable amount
	Decimals    int      `json:"decimals"`      // Token decimals (18 for KAWAI, 6 for USDT)

	// Mining-specific fields (for 9-field ClaimMiningReward)
	ContributorAmount string `json:"contributor_amount,omitempty"` // Contributor's share
	DeveloperAmount   string `json:"developer_amount,omitempty"`   // Developer's share
	UserAmount        string `json:"user_amount,omitempty"`        // User's cashback
	AffiliatorAmount  string `json:"affiliator_amount,omitempty"`  // Affiliator's commission
	DeveloperAddress  string `json:"developer_address,omitempty"`  // Developer address
	UserAddress       string `json:"user_address,omitempty"`       // User address
	AffiliatorAddress string `json:"affiliator_address,omitempty"` // Affiliator address
}

// ClaimableRewardsResponse represents the response from GetClaimableRewards
type ClaimableRewardsResponse struct {
	Address                      string             `json:"address"`
	UnclaimedProofs              []*ClaimableReward `json:"unclaimed_proofs"`
	PendingProofs                []*ClaimableReward `json:"pending_proofs"`
	ConfirmedProofs              []*ClaimableReward `json:"confirmed_proofs"` // NEW: For Recent Activity
	TotalKawaiClaimable          string             `json:"total_kawai_claimable"`
	TotalKawaiClaimableFormatted string             `json:"total_kawai_claimable_formatted"`
	TotalUSDTClaimable           string             `json:"total_usdt_claimable"`
	TotalUSDTClaimableFormatted  string             `json:"total_usdt_claimable_formatted"`
	CurrentKawaiAccumulating     string             `json:"current_kawai_accumulating"`
	CurrentUSDTAccumulating      string             `json:"current_usdt_accumulating"`
}

// ClaimResult represents the result of a claim transaction
type ClaimResult struct {
	TxHash     string `json:"tx_hash"`
	PeriodID   int64  `json:"period_id"`
	RewardType string `json:"reward_type"`
	Amount     string `json:"amount"`
	Status     string `json:"status"` // "submitted", "pending", "confirmed", "failed"
}

// DeAIService handles interactions with the Veridium smart contracts
type DeAIService struct {
	reader *reader.EthReader
	wallet *WalletService
	kv     store.Store // Cloudflare KV store for off-chain data
}

// NewDeAIService creates a new instance of DeAIService
func NewDeAIService(wallet *WalletService, kv store.Store) *DeAIService {
	// Initialize EthReader with Monad Testnet nodes from jarvis network config
	ethReader := reader.NewEthReaderGeneric(networks.MonadTestnet.GetDefaultNodes(), nil)

	return &DeAIService{
		reader: ethReader,
		wallet: wallet,
		kv:     kv,
	}
}

// GetVaultBalance returns the USDT balance of the current wallet
func (s *DeAIService) GetVaultBalance() (string, error) {
	// Check if wallet is unlocked
	if s.wallet.currentAccount == nil {
		return "", fmt.Errorf("wallet is locked")
	}

	// 1. Get User Address
	userAddr := s.wallet.currentAccount.Address()

	// 2. Load USDT
	usdtAddr, err := contracts.ResolveAddress("MockUSDT")
	if err != nil {
		return "", fmt.Errorf("USDT address not found: %w", err)
	}
	usdt, err := contracts.KawaiToken(usdtAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load USDT: %w", err)
	}

	// 3. Get Balance
	bal, err := usdt.BalanceOf(nil, userAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	// 4. Format (assuming 6 decimals)
	fBalance := new(big.Float).SetInt(bal)
	fBalance.Quo(fBalance, big.NewFloat(1000000))

	return fBalance.Text('f', 2), nil
}

// GetKawaiBalance returns the KAWAI token balance of the current wallet
func (s *DeAIService) GetKawaiBalance() (string, error) {
	// Check if wallet is unlocked
	if s.wallet.currentAccount == nil {
		return "", fmt.Errorf("wallet is locked")
	}

	// 1. Get User Address
	userAddr := s.wallet.currentAccount.Address()

	// 2. Load KAWAI Token
	kawaiAddr, err := contracts.ResolveAddress("KawaiToken")
	if err != nil {
		return "", fmt.Errorf("KAWAI address not found: %w", err)
	}
	kawai, err := contracts.KawaiToken(kawaiAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load KAWAI: %w", err)
	}

	// 3. Get Balance
	bal, err := kawai.BalanceOf(nil, userAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	// 4. Return raw balance (18 decimals)
	return bal.String(), nil
}

// GetKawaiBalanceFormatted returns the KAWAI token balance formatted (with decimals)
func (s *DeAIService) GetKawaiBalanceFormatted() (string, error) {
	balStr, err := s.GetKawaiBalance()
	if err != nil {
		return "", err
	}

	bal := new(big.Int)
	bal.SetString(balStr, 10)

	// Format (18 decimals)
	fBalance := new(big.Float).SetInt(bal)
	fBalance.Quo(fBalance, big.NewFloat(1e18))

	return fBalance.Text('f', 4), nil
}

// GetKawaiTotalSupply returns the total supply of KAWAI tokens
func (s *DeAIService) GetKawaiTotalSupply() (string, error) {
	// Load KAWAI Token
	kawaiAddr, err := contracts.ResolveAddress("KawaiToken")
	if err != nil {
		return "", fmt.Errorf("KAWAI address not found: %w", err)
	}
	kawai, err := contracts.KawaiToken(kawaiAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load KAWAI: %w", err)
	}

	// Get Total Supply
	supply, err := kawai.TotalSupply(nil)
	if err != nil {
		return "", fmt.Errorf("failed to get total supply: %w", err)
	}

	return supply.String(), nil
}

// GetRevenueShareStats returns revenue sharing statistics for the current wallet
func (s *DeAIService) GetRevenueShareStats() (map[string]interface{}, error) {
	// Check if wallet is unlocked
	if s.wallet.currentAccount == nil {
		return nil, fmt.Errorf("wallet is locked")
	}

	// Get KAWAI balance
	balanceStr, err := s.GetKawaiBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Get total supply
	supplyStr, err := s.GetKawaiTotalSupply()
	if err != nil {
		return nil, fmt.Errorf("failed to get total supply: %w", err)
	}

	// Calculate share percentage
	balance := new(big.Float)
	if _, ok := balance.SetString(balanceStr); !ok {
		return nil, fmt.Errorf("invalid balance format: %s", balanceStr)
	}

	supply := new(big.Float)
	if _, ok := supply.SetString(supplyStr); !ok {
		return nil, fmt.Errorf("invalid supply format: %s", supplyStr)
	}

	sharePercentage := new(big.Float)
	if supply.Cmp(big.NewFloat(0)) > 0 {
		sharePercentage.Quo(balance, supply)
		sharePercentage.Mul(sharePercentage, big.NewFloat(100))
	}

	// Format balance for display
	balanceFormatted := new(big.Float).Quo(balance, big.NewFloat(1e18))

	return map[string]interface{}{
		"kawai_balance":           balanceStr,
		"kawai_balance_formatted": balanceFormatted.Text('f', 4),
		"total_supply":            supplyStr,
		"share_percentage":        sharePercentage.Text('f', 6),
	}, nil
}

// DepositToVault deposits USDT into the vault for service credits
func (s *DeAIService) DepositToVault(amountStr string) (string, error) {
	// Check if wallet is unlocked
	if s.wallet.currentAccount == nil {
		return "", fmt.Errorf("wallet is locked")
	}

	// 1. Convert amount to big.Int
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// 2. Resolve Addresses
	usdtAddr, err := contracts.ResolveAddress("MockUSDT")
	if err != nil {
		return "", fmt.Errorf("MockUSDT address not found: %w", err)
	}
	vaultAddr, err := contracts.ResolveAddress("PaymentVault")
	if err != nil {
		return "", fmt.Errorf("PaymentVault address not found: %w", err)
	}

	// 3. Check Allowance
	ctx := context.Background()
	allowance, err := s.GetUSDTAllowance(s.wallet.currentAccount.AddressHex(), "PaymentVault")
	if err != nil {
		return "", fmt.Errorf("failed to check allowance: %w", err)
	}

	allowanceBig := new(big.Int)
	allowanceBig.SetString(allowance, 10)

	// 4. Approve if allowance < amount
	if allowanceBig.Cmp(amount) < 0 {
		fmt.Println("Allowance insufficient, approving...")
		chainId := monadChainID
		opts, err := s.wallet.getTransactOpts(chainId)
		if err != nil {
			return "", fmt.Errorf("failed to get opts: %w", err)
		}

		usdt, err := contracts.KawaiToken(usdtAddr.Hex(), s.reader)
		if err != nil {
			return "", fmt.Errorf("failed to load USDT: %w", err)
		}

		tx, err := usdt.Approve(opts, vaultAddr, amount)
		if err != nil {
			return "", fmt.Errorf("approve failed: %w", err)
		}

		fmt.Printf("Approval tx sent: %s. Waiting for mining...\n", tx.Hash().Hex())

		// Wait for mining
		receipt, err := bind.WaitMined(ctx, s.reader.Client(), tx)
		if err != nil {
			return "", fmt.Errorf("failed to wait for approval mining: %w", err)
		}
		if receipt.Status == 0 {
			return "", fmt.Errorf("approval transaction failed")
		}
		fmt.Println("Approval confirmed!")
	}

	// 5. Deposit
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction opts: %w", err)
	}

	vault, err := contracts.Vault("PaymentVault", s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load vault: %w", err)
	}

	tx, err := vault.Deposit(opts, amount)
	if err != nil {
		return "", fmt.Errorf("deposit failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// GetUSDTAllowance returns the current allowance of owner to spender
func (s *DeAIService) GetUSDTAllowance(ownerStr string, spenderStr string) (string, error) {
	owner := common.HexToAddress(ownerStr)
	spender, err := contracts.ResolveAddress(spenderStr)
	if err != nil {
		return "0", fmt.Errorf("invalid spender: %w", err)
	}

	usdtAddr, err := contracts.ResolveAddress("MockUSDT")
	if err != nil {
		return "0", fmt.Errorf("USDT not found")
	}

	usdt, err := contracts.KawaiToken(usdtAddr.Hex(), s.reader)
	if err != nil {
		return "0", err
	}

	allowance, err := usdt.Allowance(nil, owner, spender)
	if err != nil {
		return "0", err
	}

	return allowance.String(), nil
}

// ApproveUSDT approves a spender to spend MockUSDT
func (s *DeAIService) ApproveUSDT(spenderStr string, amountStr string) (string, error) {
	// 1. Parse inputs
	spender, err := contracts.ResolveAddress(spenderStr)
	if err != nil {
		return "", fmt.Errorf("invalid spender address: %w", err)
	}
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// 2. Get Opts
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 3. Load MockUSDT
	usdtAddr, err := contracts.ResolveAddress("MockUSDT")
	if err != nil {
		return "", fmt.Errorf("MockUSDT address not found: %w", err)
	}

	usdt, err := contracts.KawaiToken(usdtAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load USDT contract: %w", err)
	}

	// 4. Approve
	tx, err := usdt.Approve(opts, spender, amount)
	if err != nil {
		return "", fmt.Errorf("approval failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// ApproveToken approves a spender to spend a specific token
func (s *DeAIService) ApproveToken(tokenName string, spenderStr string, amountStr string) (string, error) {
	// 1. Resolve Addresses
	tokenAddr, err := contracts.ResolveAddress(tokenName)
	if err != nil {
		return "", fmt.Errorf("token address not found: %w", err)
	}
	spender, err := contracts.ResolveAddress(spenderStr)
	if err != nil {
		return "", fmt.Errorf("invalid spender address: %w", err)
	}

	// 2. Parse Amount
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// 3. Get Opts
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 4. Load Token
	token, err := contracts.KawaiToken(tokenAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load token: %w", err)
	}

	// 5. Approve
	tx, err := token.Approve(opts, spender, amount)
	if err != nil {
		return "", fmt.Errorf("approval failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// CreateSellOrder creates a sell order in the OTC Market
func (s *DeAIService) CreateSellOrder(tokenAmountStr string, priceStr string) (string, error) {
	tokenAmount := new(big.Int)
	tokenAmount, ok := tokenAmount.SetString(tokenAmountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid token amount")
	}
	price := new(big.Int)
	price, ok = price.SetString(priceStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid price format")
	}

	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	escrowAddr, err := contracts.ResolveAddress("Escrow")
	if err != nil {
		return "", fmt.Errorf("Escrow address not found: %w", err)
	}

	escrow, err := contracts.Escrow(escrowAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load Escrow: %w", err)
	}

	tx, err := escrow.CreateOrder(opts, tokenAmount, price)
	if err != nil {
		return "", fmt.Errorf("create order failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// BuyOrder buys an order from the OTC Market
func (s *DeAIService) BuyOrder(orderIdStr string) (string, error) {
	orderId := new(big.Int)
	orderId, ok := orderId.SetString(orderIdStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid order id")
	}

	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	escrowAddr, err := contracts.ResolveAddress("Escrow")
	if err != nil {
		return "", fmt.Errorf("Escrow address not found: %w", err)
	}

	escrow, err := contracts.Escrow(escrowAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load Escrow: %w", err)
	}

	tx, err := escrow.BuyOrder(opts, orderId)
	if err != nil {
		return "", fmt.Errorf("buy order failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// MintTestTokens mints MockUSDT and KawaiTokens to the caller (for testing only)
func (s *DeAIService) MintTestTokens() (string, error) {
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 1. Mint USDT
	usdtAddr, err := contracts.ResolveAddress("MockUSDT")
	if err != nil {
		return "", fmt.Errorf("token address not found")
	}
	usdt, _ := contracts.KawaiToken(usdtAddr.Hex(), s.reader) // Using KawaiToken wrapper for mint

	// Mint 1000 USDT
	amount := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1000000)) // 1000 * 10^6
	tx1, err := usdt.Mint(opts, opts.From, amount)
	if err != nil {
		return "", fmt.Errorf("mint usdt failed: %w", err)
	}

	return tx1.Hash().Hex(), nil
}

// TransferUSDT sends USDT from the current wallet to a recipient
func (s *DeAIService) TransferUSDT(to string, amountStr string) (string, error) {
	// 1. Resolve Addresses
	usdtAddr, err := contracts.ResolveAddress("MockUSDT")
	if err != nil {
		return "", fmt.Errorf("USDT address not found: %w", err)
	}
	recipient := common.HexToAddress(to)

	// 2. Parse Amount
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// 3. Get Opts
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 4. Load Contract
	usdt, err := contracts.KawaiToken(usdtAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load USDT: %w", err)
	}

	// 5. Transfer
	tx, err := usdt.Transfer(opts, recipient, amount)
	if err != nil {
		return "", fmt.Errorf("transfer failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// TransferNative sends native coin (MON, ETH) from the current wallet to a recipient
func (s *DeAIService) TransferNative(to string, amountStr string) (string, error) {
	// 1. Parse address
	recipient := common.HexToAddress(to)

	// 2. Parse Amount (input is in ETH string, e.g., "0.1")
	// Convert to Wei (10^18)
	val, ok := new(big.Float).SetString(amountStr)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}
	wei := new(big.Float).Mul(val, big.NewFloat(1e18))
	amount := new(big.Int)
	wei.Int(amount) // Convert float to int

	// 3. Get Opts (Wait for signing)
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 4. Create Transaction
	// Native transfer is just a transaction with value
	nonce, err := s.reader.Client().PendingNonceAt(context.Background(), opts.From)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := s.reader.Client().SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	gasLimit := uint64(21000) // Standard transfer gas limit

	// Create transaction
	tx := ethtypes.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, nil)

	// 5. Sign Transaction
	signedTx, err := opts.Signer(opts.From, tx)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	// 6. Send Transaction
	err = s.reader.Client().SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send tx: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

// TransferToken sends any ERC20 token from the current wallet to a recipient
func (s *DeAIService) TransferToken(tokenAddress string, to string, amountStr string) (string, error) {
	// 1. Validate inputs
	if !common.IsHexAddress(tokenAddress) {
		return "", fmt.Errorf("invalid token address")
	}
	if !common.IsHexAddress(to) {
		return "", fmt.Errorf("invalid recipient address")
	}

	// 2. Parse Amount (Raw integer string, handled by caller or assumed raw)
	// For this generic function, let's assume raw amount for now to overlap with other implementations,
	// OR better, let the frontend pass the raw string.
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	recipient := common.HexToAddress(to)

	// 3. Get Opts
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 4. Load Contract Generic
	// We use KawaiToken wrapper because it satisfies standard ERC20 interface
	token, err := contracts.KawaiToken(tokenAddress, s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load token contract: %w", err)
	}

	// 5. Transfer
	tx, err := token.Transfer(opts, recipient, amount)
	if err != nil {
		return "", fmt.Errorf("transfer failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// =============================================================================
// REWARDS CLAIM METHODS
// =============================================================================

// getContractPeriodForTimestamp is no longer needed - contract now supports timestamp-based periods
// Contract's setMerkleRootForPeriod() accepts any period ID (timestamp)
// This function is kept for backward compatibility but just returns the timestamp as-is
func (s *DeAIService) getContractPeriodForTimestamp(timestamp int64) int64 {
	// Contract now uses timestamp-based periods directly via setMerkleRootForPeriod()
	return timestamp
}

// GetClaimableRewards fetches all claimable rewards for the current wallet
// Uses Cloudflare KV store directly for off-chain Merkle proof data
func (s *DeAIService) GetClaimableRewards() (*ClaimableRewardsResponse, error) {
	if s.wallet.currentAccount == nil {
		return nil, fmt.Errorf("no wallet connected")
	}

	if s.kv == nil {
		return nil, fmt.Errorf("KV store not initialized")
	}

	userAddr := s.wallet.currentAccount.AddressHex()
	ctx := context.Background()

	// Get claimable rewards from KV store
	claimableData, err := s.kv.GetClaimableRewards(ctx, userAddr)
	if err != nil {
		// Return empty response if no data found
		return &ClaimableRewardsResponse{
			Address:                      userAddr,
			UnclaimedProofs:              []*ClaimableReward{},
			PendingProofs:                []*ClaimableReward{},
			ConfirmedProofs:              []*ClaimableReward{}, // NEW: Empty confirmed proofs
			TotalKawaiClaimable:          "0",
			TotalKawaiClaimableFormatted: "0.0000",
			TotalUSDTClaimable:           "0",
			TotalUSDTClaimableFormatted:  "0.00",
			CurrentKawaiAccumulating:     "0",
			CurrentUSDTAccumulating:      "0",
		}, nil
	}

	// Convert from map[string]interface{} to ClaimableRewardsResponse
	result := &ClaimableRewardsResponse{
		Address:                  userAddr,
		UnclaimedProofs:          []*ClaimableReward{},
		PendingProofs:            []*ClaimableReward{},
		ConfirmedProofs:          []*ClaimableReward{}, // NEW: Initialize confirmed proofs
		TotalKawaiClaimable:      getStringFromMap(claimableData, "total_kawai_claimable", "0"),
		TotalUSDTClaimable:       getStringFromMap(claimableData, "total_usdt_claimable", "0"),
		CurrentKawaiAccumulating: getStringFromMap(claimableData, "current_kawai_accumulating", "0"),
		CurrentUSDTAccumulating:  getStringFromMap(claimableData, "current_usdt_accumulating", "0"),
	}

	// Convert unclaimed proofs
	if unclaimedRaw, ok := claimableData["unclaimed_proofs"]; ok {
		if unclaimedList, ok := unclaimedRaw.([]*store.MerkleProofData); ok {
			for _, proof := range unclaimedList {
				claimable := s.convertMerkleProofToClaimable(proof)
				result.UnclaimedProofs = append(result.UnclaimedProofs, claimable)
			}
		}
	}

	// Convert pending proofs
	if pendingRaw, ok := claimableData["pending_proofs"]; ok {
		if pendingList, ok := pendingRaw.([]*store.MerkleProofData); ok {
			for _, proof := range pendingList {
				claimable := s.convertMerkleProofToClaimable(proof)
				result.PendingProofs = append(result.PendingProofs, claimable)
			}
		}
	}

	// NEW: Convert confirmed proofs for Recent Activity
	if confirmedRaw, ok := claimableData["confirmed_proofs"]; ok {
		if confirmedList, ok := confirmedRaw.([]*store.MerkleProofData); ok {
			for _, proof := range confirmedList {
				claimable := s.convertMerkleProofToClaimable(proof)
				result.ConfirmedProofs = append(result.ConfirmedProofs, claimable)
			}
		}
	}

	// Format total amounts
	result.TotalKawaiClaimableFormatted = s.formatRewardAmount(result.TotalKawaiClaimable, "kawai")
	result.TotalUSDTClaimableFormatted = s.formatRewardAmount(result.TotalUSDTClaimable, "usdt")

	return result, nil
}

// convertMerkleProofToClaimable converts store.MerkleProofData to ClaimableReward
func (s *DeAIService) convertMerkleProofToClaimable(proof *store.MerkleProofData) *ClaimableReward {
	decimals := 18
	if proof.RewardType == "usdt" {
		decimals = 6
	}

	claimable := &ClaimableReward{
		Index:       proof.Index,
		Amount:      proof.Amount,
		Proof:       proof.Proof,
		MerkleRoot:  proof.MerkleRoot,
		PeriodID:    proof.PeriodID,
		RewardType:  proof.RewardType,
		ClaimStatus: string(proof.ClaimStatus),
		ClaimTxHash: proof.ClaimTxHash,
		CreatedAt:   proof.CreatedAt.Format("2006-01-02T15:04:05Z"),
		ClaimedAt:   proof.ClaimedAt.Format("2006-01-02T15:04:05Z"),
		Formatted:   s.formatRewardAmount(proof.Amount, proof.RewardType),
		Decimals:    decimals,
	}

	// Add mining-specific fields for kawai rewards (9-field format)
	if proof.RewardType == "kawai" {
		claimable.ContributorAmount = proof.ContributorAmount
		claimable.DeveloperAmount = proof.DeveloperAmount
		claimable.UserAmount = proof.UserAmount
		claimable.AffiliatorAmount = proof.AffiliatorAmount
		claimable.DeveloperAddress = proof.DeveloperAddress
		claimable.UserAddress = proof.UserAddress
		claimable.AffiliatorAddress = proof.AffiliatorAddress
	}

	return claimable
}

// getStringFromMap safely extracts a string from map[string]interface{}
func getStringFromMap(m map[string]interface{}, key string, defaultVal string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultVal
}

// ClaimKawaiReward claims KAWAI rewards using a Merkle proof
func (s *DeAIService) ClaimKawaiReward(periodID int64, index uint64, amount string, proof []string) (*ClaimResult, error) {
	return s.claimReward("kawai", periodID, index, amount, proof)
}

// ClaimUSDTReward claims USDT rewards using a Merkle proof
func (s *DeAIService) ClaimUSDTReward(periodID int64, index uint64, amount string, proof []string) (*ClaimResult, error) {
	return s.claimReward("usdt", periodID, index, amount, proof)
}

// ClaimCashbackReward claims deposit cashback rewards using a Merkle proof
// Uses the DepositCashbackDistributor contract with period-based claims
func (s *DeAIService) ClaimCashbackReward(period uint64, kawaiAmount string, proof []string) (*ClaimResult, error) {
	if s.wallet.currentAccount == nil {
		return nil, fmt.Errorf("no wallet connected")
	}

	// 1. Input validation
	// Note: Empty proof is valid for single-leaf Merkle trees
	if kawaiAmount == "" || kawaiAmount == "0" {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// 2. Load DepositCashbackDistributor contract
	distributor, err := contracts.CashbackDistributor("CashbackDistributor", s.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to load cashback distributor: %w", err)
	}

	// 3. Parse amount
	amount := new(big.Int)
	amount, ok := amount.SetString(kawaiAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount format")
	}
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	// 3. Convert proof strings to [32]byte array
	merkleProof := make([][32]byte, len(proof))
	for i, p := range proof {
		// Remove 0x prefix if present
		if len(p) >= 2 && p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		if len(proofBytes) != 32 {
			return nil, fmt.Errorf("invalid proof element at index %d: expected 32 bytes, got %d", i, len(proofBytes))
		}
		copy(merkleProof[i][:], proofBytes)
	}

	// 4. Get transaction options
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}

	// 5. Mark claim as pending BEFORE submitting transaction (prevents double-claim UX issue)
	ctx := context.Background()
	if s.kv != nil {
		kvStore, ok := s.kv.(*store.KVStore)
		if ok {
			if err := kvStore.MarkCashbackPending(ctx, s.wallet.currentAccount.AddressHex(), period, ""); err != nil {
				return nil, fmt.Errorf("failed to mark cashback claim as pending: %w", err)
			}
		}
	}

	// 6. Submit claim transaction
	tx, err := distributor.ClaimCashback(opts, big.NewInt(int64(period)), amount, merkleProof)
	if err != nil {
		// Rollback pending status on transaction failure
		if s.kv != nil {
			kvStore, ok := s.kv.(*store.KVStore)
			if ok {
				kvStore.MarkCashbackFailed(ctx, s.wallet.currentAccount.AddressHex(), period, err.Error())
			}
		}
		return nil, fmt.Errorf("claim transaction failed: %w", err)
	}

	txHash := tx.Hash().Hex()

	// 7. Wait for transaction confirmation
	receipt, err := bind.WaitMined(ctx, s.reader.Client(), tx)
	if err != nil {
		// Keep pending status - will be checked later
		return nil, fmt.Errorf("failed to wait for transaction confirmation: %w", err)
	}

	// 8. Check transaction status
	if receipt.Status != 1 {
		// Mark as failed if transaction reverted
		if s.kv != nil {
			kvStore, ok := s.kv.(*store.KVStore)
			if ok {
				kvStore.MarkCashbackFailed(ctx, s.wallet.currentAccount.AddressHex(), period, "transaction reverted")
			}
		}
		return nil, fmt.Errorf("transaction reverted (status: %d)", receipt.Status)
	}

	// 9. Mark claim as completed in KV store
	if s.kv != nil {
		kvStore, ok := s.kv.(*store.KVStore)
		if ok {
			if err := kvStore.MarkCashbackClaimed(ctx, s.wallet.currentAccount.AddressHex(), period); err != nil {
				// Log warning - tx was successful, but KV update failed
				fmt.Printf("Warning: failed to mark cashback claim in KV: %v\n", err)
			}
		}
	}

	return &ClaimResult{
		TxHash:     txHash,
		PeriodID:   int64(period),
		RewardType: "cashback",
		Amount:     kawaiAmount,
		Status:     "confirmed",
	}, nil
}

// ClaimMiningReward claims mining rewards with referral-based splits
// Uses the new MiningRewardDistributor contract with 9-field Merkle leaves
// ClaimMiningReward claims mining rewards with referral splits
// Maps timestamp-based settlement periods to sequential contract periods
func (s *DeAIService) ClaimMiningReward(
	period int64,
	contributorAmount string,
	developerAmount string,
	userAmount string,
	affiliatorAmount string,
	developerAddress string,
	userAddress string,
	affiliatorAddress string,
	proof []string,
) (*ClaimResult, error) {
	if s.wallet.currentAccount == nil {
		return nil, fmt.Errorf("no wallet connected")
	}

	// Contract now supports timestamp-based periods directly via setMerkleRootForPeriod()
	// No mapping needed - use timestamp period ID as-is
	contractPeriod := period

	fmt.Printf("🔄 Claiming mining reward: period %d (timestamp-based)\n", period)

	// Load MiningRewardDistributor contract
	distributor, err := contracts.MiningRewardDistributor("MiningRewardDistributor", s.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to load mining distributor: %w", err)
	}

	// Parse amounts
	contribAmt := new(big.Int)
	contribAmt, ok := contribAmt.SetString(contributorAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid contributor amount format")
	}

	devAmt := new(big.Int)
	devAmt, ok = devAmt.SetString(developerAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid developer amount format")
	}

	userAmt := new(big.Int)
	userAmt, ok = userAmt.SetString(userAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid user amount format")
	}

	affAmt := new(big.Int)
	affAmt, ok = affAmt.SetString(affiliatorAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid affiliator amount format")
	}

	// 4. Convert proof strings to [32]byte array
	merkleProof := make([][32]byte, len(proof))
	for i, p := range proof {
		// Remove 0x prefix if present
		if len(p) >= 2 && p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		if len(proofBytes) != 32 {
			return nil, fmt.Errorf("invalid proof element at index %d: expected 32 bytes, got %d", i, len(proofBytes))
		}
		copy(merkleProof[i][:], proofBytes)
	}

	// 5. Parse addresses
	devAddr := common.HexToAddress(developerAddress)
	usrAddr := common.HexToAddress(userAddress)
	affAddr := common.HexToAddress(affiliatorAddress)

	// 6. Get transaction options
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}

	// 7. Mark claim as pending BEFORE submitting transaction (prevents double-claim UX issue)
	ctx := context.Background()
	if s.kv != nil {
		if err := s.kv.MarkClaimPending(ctx, s.wallet.currentAccount.AddressHex(), period, ""); err != nil {
			return nil, fmt.Errorf("failed to mark claim as pending: %w", err)
		}
	}

	// 8. Submit claim transaction using the mapped contract period
	// claimReward(uint256 period, uint256 contributorAmount, uint256 developerAmount,
	//             uint256 userAmount, uint256 affiliatorAmount, address developer,
	//             address user, address affiliator, bytes32[] calldata merkleProof)
	tx, err := distributor.ClaimReward(
		opts,
		big.NewInt(contractPeriod), // Use sequential contract period
		contribAmt,
		devAmt,
		userAmt,
		affAmt,
		devAddr,
		usrAddr,
		affAddr,
		merkleProof,
	)
	if err != nil {
		// Rollback pending status on transaction failure
		if s.kv != nil {
			s.kv.MarkClaimFailed(ctx, s.wallet.currentAccount.AddressHex(), period, err.Error())
		}
		return nil, fmt.Errorf("mining claim transaction failed: %w", err)
	}

	txHash := tx.Hash().Hex()

	// 9. Wait for transaction confirmation
	receipt, err := bind.WaitMined(ctx, s.reader.Client(), tx)
	if err != nil {
		// Keep pending status - auto-confirm will check later
		return nil, fmt.Errorf("failed to wait for transaction confirmation: %w", err)
	}

	// 10. Check transaction status
	if receipt.Status != 1 {
		// Mark as failed if transaction reverted
		if s.kv != nil {
			s.kv.MarkClaimFailed(ctx, s.wallet.currentAccount.AddressHex(), period, "transaction reverted")
		}
		return nil, fmt.Errorf("transaction reverted (status: %d)", receipt.Status)
	}

	// 11. Update with actual tx hash (status remains pending, auto-confirm will update to confirmed)
	if s.kv != nil {
		if err := s.kv.MarkClaimPending(ctx, s.wallet.currentAccount.AddressHex(), period, txHash); err != nil {
			// Log warning - tx was successful, auto-confirm will fix status later
			fmt.Printf("Warning: failed to update mining claim tx hash in KV: %v\n", err)
		}
	}

	return &ClaimResult{
		TxHash:     txHash,
		PeriodID:   period,
		RewardType: "mining",
		Amount:     contributorAmount,
		Status:     "confirmed",
	}, nil
}

// claimReward is the internal implementation for claiming rewards
func (s *DeAIService) claimReward(rewardType string, periodID int64, index uint64, amountStr string, proofStrings []string) (*ClaimResult, error) {
	if s.wallet.currentAccount == nil {
		return nil, fmt.Errorf("no wallet connected")
	}

	// 1. Resolve distributor address
	var distributorName string
	if rewardType == "kawai" {
		distributorName = "KAWAI_Distributor"
	} else {
		distributorName = "USDT_Distributor"
	}

	// 2. Load MerkleDistributor contract
	distributor, err := contracts.MerkleDistributor(distributorName, s.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to load distributor: %w", err)
	}

	// 3. Parse amount
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount format")
	}

	// 4. Convert proof strings to [32]byte array
	merkleProof := make([][32]byte, len(proofStrings))
	for i, p := range proofStrings {
		// Remove 0x prefix if present
		if len(p) >= 2 && p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		if len(proofBytes) != 32 {
			return nil, fmt.Errorf("invalid proof element at index %d: expected 32 bytes, got %d", i, len(proofBytes))
		}
		copy(merkleProof[i][:], proofBytes)
	}

	// 5. Get transaction options
	chainId := monadChainID
	opts, err := s.wallet.getTransactOpts(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}

	// 6. Submit claim transaction
	tx, err := distributor.Claim(opts, big.NewInt(int64(index)), opts.From, amount, merkleProof)
	if err != nil {
		return nil, fmt.Errorf("claim transaction failed: %w", err)
	}

	txHash := tx.Hash().Hex()

	// 7. Wait for transaction confirmation
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, s.reader.Client(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction confirmation: %w", err)
	}

	// 8. Check transaction status
	if receipt.Status != 1 {
		return nil, fmt.Errorf("transaction reverted (status: %d)", receipt.Status)
	}

	// 9. Mark claim as completed in KV store (for tracking)
	if s.kv != nil {
		if err := s.kv.MarkClaimPending(ctx, s.wallet.currentAccount.AddressHex(), periodID, txHash); err != nil {
			// Log warning but don't fail - the TX was successful
			fmt.Printf("Warning: failed to mark claim in KV: %v\n", err)
		}
	}

	return &ClaimResult{
		TxHash:     txHash,
		PeriodID:   periodID,
		RewardType: rewardType,
		Amount:     amountStr,
		Status:     "confirmed",
	}, nil
}

// IsRewardClaimed checks if a specific reward has already been claimed on-chain
func (s *DeAIService) IsRewardClaimed(rewardType string, index uint64) (bool, error) {
	// 1. Resolve distributor address
	var distributorName string
	if rewardType == "kawai" {
		distributorName = "KAWAI_Distributor"
	} else {
		distributorName = "USDT_Distributor"
	}

	// 2. Load MerkleDistributor contract
	distributor, err := contracts.MerkleDistributor(distributorName, s.reader)
	if err != nil {
		return false, fmt.Errorf("failed to load distributor: %w", err)
	}

	// 3. Check if claimed
	claimed, err := distributor.IsClaimed(nil, big.NewInt(int64(index)))
	if err != nil {
		return false, fmt.Errorf("failed to check claim status: %w", err)
	}

	return claimed, nil
}

// GetDistributorMerkleRoot returns the current Merkle root from a distributor contract
func (s *DeAIService) GetDistributorMerkleRoot(rewardType string) (string, error) {
	var distributorName string
	if rewardType == "kawai" {
		distributorName = "KAWAI_Distributor"
	} else {
		distributorName = "USDT_Distributor"
	}

	distributor, err := contracts.MerkleDistributor(distributorName, s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load distributor: %w", err)
	}

	root, err := distributor.MerkleRoot(nil)
	if err != nil {
		return "", fmt.Errorf("failed to get merkle root: %w", err)
	}

	return fmt.Sprintf("0x%x", root), nil
}

// WaitForClaimConfirmation waits for a claim transaction to be mined
func (s *DeAIService) WaitForClaimConfirmation(txHash string) (bool, error) {
	ctx := context.Background()

	// Parse transaction hash
	hash := common.HexToHash(txHash)

	// Wait for receipt (with timeout handled by context)
	receipt, err := bind.WaitMined(ctx, s.reader.Client(), ethtypes.NewTx(&ethtypes.LegacyTx{}))
	if err != nil {
		// Try alternative: query receipt directly
		receipt, err = s.reader.Client().TransactionReceipt(ctx, hash)
		if err != nil {
			return false, fmt.Errorf("failed to get transaction receipt: %w", err)
		}
	}

	// Check status (1 = success, 0 = failed)
	return receipt.Status == 1, nil
}

// ConfirmRewardClaim confirms a reward claim after the transaction is confirmed on-chain
func (s *DeAIService) ConfirmRewardClaim(periodID int64) error {
	if s.wallet.currentAccount == nil {
		return fmt.Errorf("no wallet connected")
	}

	if s.kv == nil {
		return fmt.Errorf("KV store not initialized")
	}

	ctx := context.Background()
	userAddr := s.wallet.currentAccount.AddressHex()

	return s.kv.ConfirmClaim(ctx, userAddr, periodID)
}

// MarkClaimFailed marks a claim as failed after the transaction reverts
func (s *DeAIService) MarkClaimFailed(periodID int64, reason string) error {
	if s.wallet.currentAccount == nil {
		return fmt.Errorf("no wallet connected")
	}

	if s.kv == nil {
		return fmt.Errorf("KV store not initialized")
	}

	ctx := context.Background()
	userAddr := s.wallet.currentAccount.AddressHex()

	return s.kv.MarkClaimFailed(ctx, userAddr, periodID, reason)
}

// formatRewardAmount formats raw token amount to human-readable format
func (s *DeAIService) formatRewardAmount(rawAmount string, rewardType string) string {
	amount := new(big.Int)
	amount, ok := amount.SetString(rawAmount, 10)
	if !ok {
		return "0.00"
	}

	var decimals int64
	if rewardType == "usdt" {
		decimals = 6
	} else {
		decimals = 18
	}

	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil))
	formatted := new(big.Float).Quo(new(big.Float).SetInt(amount), divisor)

	// Format with appropriate precision
	if rewardType == "usdt" {
		return formatted.Text('f', 2)
	}
	return formatted.Text('f', 4)
}
