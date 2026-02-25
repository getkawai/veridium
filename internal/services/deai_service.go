package services

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/kawai-network/x/alert"
	"github.com/kawai-network/y/config"
	"github.com/kawai-network/y/jarvis/binding"
	"github.com/kawai-network/y/jarvis/networks"
	"github.com/kawai-network/y/jarvis/util/reader"
	"github.com/kawai-network/x/store"
	"github.com/kawai-network/y/types"
	"github.com/kawai-network/x/constant"
	"github.com/kawai-network/contracts"
)

// isUserError checks if an error is a user-caused error (not system error)
// User errors should not trigger alerts as they are expected
func isUserError(err error) bool {
	if err == nil {
		return false
	}

	userErrorPatterns := []string{
		"Already claimed",
		"already claimed",
		"Invalid proof",
		"invalid proof",
		"Insufficient balance",
		"insufficient balance",
		"Not eligible",
		"not eligible",
		"no wallet connected",
		"amount must be",
		"invalid amount",
	}

	errStr := err.Error()
	for _, pattern := range userErrorPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// ClaimableReward represents a single claimable reward proof
type ClaimableReward struct {
	Index       uint64           `json:"index"`
	Amount      string           `json:"amount"`        // BigInt as string (raw, no decimals)
	Proof       []string         `json:"proof"`         // Merkle proof (hex strings)
	MerkleRoot  string           `json:"merkle_root"`   // Root hash for verification
	PeriodID    int64            `json:"period_id"`     // Settlement period identifier
	RewardType  types.RewardType `json:"reward_type"`   // "mining", "cashback", "referral", "stablecoin"
	ClaimStatus string           `json:"claim_status"`  // "unclaimed", "pending", "confirmed", "failed"
	ClaimTxHash string           `json:"claim_tx_hash"` // Transaction hash if claimed
	CreatedAt   string           `json:"created_at"`    // When proof was generated
	ClaimedAt   string           `json:"claimed_at"`    // When claimed (if confirmed)
	Formatted   string           `json:"formatted"`     // Human-readable amount
	Decimals    int              `json:"decimals"`      // Token decimals (18 for KAWAI, 6 for USDT)

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
	TxHash     string           `json:"tx_hash"`
	PeriodID   int64            `json:"period_id"`
	RewardType types.RewardType `json:"reward_type"`
	Amount     string           `json:"amount"`
	Status     string           `json:"status"` // "submitted", "pending", "confirmed", "failed"
}

// RevenueShareStatsResponse represents revenue sharing statistics
type RevenueShareStatsResponse struct {
	KawaiBalance          string `json:"kawai_balance"`           // Raw balance (18 decimals)
	KawaiBalanceFormatted string `json:"kawai_balance_formatted"` // Human-readable balance
	TotalSupply           string `json:"total_supply"`            // Total KAWAI supply
	SharePercentage       string `json:"share_percentage"`        // User's share percentage
}

// WalletOperations defines the minimal wallet interface needed by DeAIService
type WalletOperations interface {
	IsUnlocked() bool
	GetCurrentAccountAddress() string
	GetCurrentAddress() string
	GetTransactOpts(chainId *big.Int) (*bind.TransactOpts, error)
}

// DeAIService handles interactions with the Veridium smart contracts
type DeAIService struct {
	reader              *reader.EthReader
	wallet              WalletOperations
	kv                  store.Store // Cloudflare KV store for off-chain data
	alert               *alert.DiscordAlert
	chainID             *big.Int             // Chain ID for the current Monad environment (testnet or mainnet)
	contributorSelector *ContributorSelector // Contributor discovery and selection
}

// NewDeAIService creates a new instance of DeAIService
func NewDeAIService(wallet WalletOperations, kv store.Store) *DeAIService {
	// Initialize EthReader with Monad nodes based on current environment
	var ethReader *reader.EthReader
	if config.IsTestnet() {
		ethReader = reader.NewEthReaderGeneric(networks.MonadTestnet.GetDefaultNodes(), nil)
	} else {
		ethReader = reader.NewEthReaderGeneric(networks.MonadMainnet.GetDefaultNodes(), nil)
	}
	discordAlerter := &alert.DiscordAlert{
		WebhookURL: constant.GetDiscordClaimFailure(),
		Client:     &http.Client{Timeout: 10 * time.Second},
	}

	return &DeAIService{
		reader:              ethReader,
		wallet:              wallet,
		kv:                  kv,
		alert:               discordAlerter,
		chainID:             big.NewInt(int64(config.GetChainID())),
		contributorSelector: NewContributorSelector(kv),
	}
}

// GetVaultBalance returns the stablecoin balance of the current wallet
// Note: Uses MockStablecoin on testnet, USDC on mainnet
func (s *DeAIService) GetVaultBalance() (string, error) {
	// Check if wallet is unlocked
	if !s.wallet.IsUnlocked() {
		return "", fmt.Errorf("wallet is locked")
	}

	// 1. Get User Address
	userAddr := s.wallet.GetCurrentAccountAddress()

	// 2. Load stablecoin contract (MockStablecoin on testnet, USDC on mainnet)
	// Use Jarvis wrapper for cleaner code
	stablecoin, err := binding.Stablecoin(contracts.StablecoinAddress, s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load stablecoin contract: %w", err)
	}

	// 3. Get Balance
	bal, err := stablecoin.BalanceOf(nil, common.HexToAddress(userAddr))
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	// 4. Format (assuming 6 decimals for both USDT and USDC)
	fBalance := new(big.Float).SetInt(bal)
	fBalance.Quo(fBalance, big.NewFloat(1000000))

	return fBalance.Text('f', 2), nil
}

// GetKawaiBalance returns the KAWAI token balance of the current wallet
func (s *DeAIService) GetKawaiBalance() (string, error) {
	// Check if wallet is unlocked
	if !s.wallet.IsUnlocked() {
		return "", fmt.Errorf("wallet is locked")
	}

	// 1. Get User Address
	userAddr := s.wallet.GetCurrentAccountAddress()

	// 2. Load KAWAI Token
	kawaiAddr, err := binding.ResolveAddress("KawaiToken")
	if err != nil {
		return "", fmt.Errorf("KAWAI address not found: %w", err)
	}
	kawai, err := binding.KawaiToken(kawaiAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load KAWAI: %w", err)
	}

	// 3. Get Balance
	bal, err := kawai.BalanceOf(nil, common.HexToAddress(userAddr))
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
	kawaiAddr, err := binding.ResolveAddress("KawaiToken")
	if err != nil {
		return "", fmt.Errorf("KAWAI address not found: %w", err)
	}
	kawai, err := binding.KawaiToken(kawaiAddr.Hex(), s.reader)
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
func (s *DeAIService) GetRevenueShareStats() (*RevenueShareStatsResponse, error) {
	// Check if wallet is unlocked
	if !s.wallet.IsUnlocked() {
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

	return &RevenueShareStatsResponse{
		KawaiBalance:          balanceStr,
		KawaiBalanceFormatted: balanceFormatted.Text('f', 4),
		TotalSupply:           supplyStr,
		SharePercentage:       sharePercentage.Text('f', 6),
	}, nil
}

// DepositToVault deposits stablecoin into the vault for service credits
// Note: Uses MockStablecoin on testnet, USDC on mainnet
func (s *DeAIService) DepositToVault(amountStr string) (string, error) {
	// Check if wallet is unlocked
	if !s.wallet.IsUnlocked() {
		return "", fmt.Errorf("wallet is locked")
	}

	// 1. Convert amount to big.Int
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// 2. Resolve Addresses
	// Use constant which automatically switches based on environment
	stablecoinAddr := common.HexToAddress(contracts.StablecoinAddress)
	vaultAddr, err := binding.ResolveAddress("PaymentVault")
	if err != nil {
		return "", fmt.Errorf("PaymentVault address not found: %w", err)
	}

	// 3. Check Allowance
	ctx := context.Background()
	allowance, err := s.GetUSDTAllowance(s.wallet.GetCurrentAccountAddress(), "PaymentVault")
	if err != nil {
		return "", fmt.Errorf("failed to check allowance: %w", err)
	}

	allowanceBig := new(big.Int)
	allowanceBig.SetString(allowance, 10)

	// 4. Approve if allowance < amount
	if allowanceBig.Cmp(amount) < 0 {
		fmt.Println("Allowance insufficient, approving...")
		chainId := s.chainID
		opts, err := s.wallet.GetTransactOpts(chainId)
		if err != nil {
			return "", fmt.Errorf("failed to get opts: %w", err)
		}

		stablecoin, err := binding.KawaiToken(stablecoinAddr.Hex(), s.reader)
		if err != nil {
			return "", fmt.Errorf("failed to load stablecoin contract: %w", err)
		}

		tx, err := stablecoin.Approve(opts, vaultAddr, amount)
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
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction opts: %w", err)
	}

	vault, err := binding.Vault("PaymentVault", s.reader)
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
// Note: Function name kept for backward compatibility, but works with any stablecoin
func (s *DeAIService) GetUSDTAllowance(ownerStr string, spenderStr string) (string, error) {
	owner := common.HexToAddress(ownerStr)
	spender, err := binding.ResolveAddress(spenderStr)
	if err != nil {
		return "0", fmt.Errorf("invalid spender: %w", err)
	}

	// Use constant which automatically switches based on environment
	stablecoinAddr := common.HexToAddress(contracts.StablecoinAddress)

	stablecoin, err := binding.KawaiToken(stablecoinAddr.Hex(), s.reader)
	if err != nil {
		return "0", err
	}

	allowance, err := stablecoin.Allowance(nil, owner, spender)
	if err != nil {
		return "0", err
	}

	return allowance.String(), nil
}

// ApproveUSDT approves a spender to spend stablecoin (MockStablecoin on testnet, USDC on mainnet)
// Note: Function name kept for backward compatibility
func (s *DeAIService) ApproveUSDT(spenderStr string, amountStr string) (string, error) {
	// 1. Parse inputs
	spender, err := binding.ResolveAddress(spenderStr)
	if err != nil {
		return "", fmt.Errorf("invalid spender address: %w", err)
	}
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// 2. Get Opts
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 3. Load stablecoin contract
	// Use constant which automatically switches based on environment
	stablecoinAddr := common.HexToAddress(contracts.StablecoinAddress)

	stablecoin, err := binding.KawaiToken(stablecoinAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load stablecoin contract: %w", err)
	}

	// 4. Approve
	tx, err := stablecoin.Approve(opts, spender, amount)
	if err != nil {
		return "", fmt.Errorf("approval failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// ApproveToken approves a spender to spend a specific token
func (s *DeAIService) ApproveToken(tokenName string, spenderStr string, amountStr string) (string, error) {
	// 1. Resolve Addresses
	tokenAddr, err := binding.ResolveAddress(tokenName)
	if err != nil {
		return "", fmt.Errorf("token address not found: %w", err)
	}
	spender, err := binding.ResolveAddress(spenderStr)
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
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 4. Load Token
	token, err := binding.KawaiToken(tokenAddr.Hex(), s.reader)
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

	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	escrowAddr, err := binding.ResolveAddress("OTCMarket")
	if err != nil {
		return "", fmt.Errorf("OTCMarket address not found: %w", err)
	}

	escrow, err := binding.OTCMarket(escrowAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load OTCMarket: %w", err)
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

	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	escrowAddr, err := binding.ResolveAddress("OTCMarket")
	if err != nil {
		return "", fmt.Errorf("OTCMarket address not found: %w", err)
	}

	escrow, err := binding.OTCMarket(escrowAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load OTCMarket: %w", err)
	}

	tx, err := escrow.BuyOrder(opts, orderId)
	if err != nil {
		return "", fmt.Errorf("buy order failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// MintTestTokens mints test stablecoin (MockStablecoin) to the caller (for testing only)
// WARNING: This function only works on testnet with MockStablecoin. It will FAIL on mainnet with USDC.
// USDC on mainnet does not have a public mint() function.
func (s *DeAIService) MintTestTokens() (string, error) {
	// Safety check: Only allow on testnet
	// Testnet uses MockStablecoin which has mint(), mainnet uses USDC which doesn't
	if !config.IsTestnet() {
		return "", fmt.Errorf("MintTestTokens is only available on testnet. On mainnet, you must acquire USDC through exchanges or bridges")
	}

	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 1. Mint stablecoin (MockStablecoin on testnet only)
	// Use constant which automatically switches based on environment
	stablecoinAddr := common.HexToAddress(contracts.StablecoinAddress)
	stablecoin, _ := binding.KawaiToken(stablecoinAddr.Hex(), s.reader) // Using KawaiToken wrapper for mint

	// Mint 1000 stablecoin (6 decimals)
	amount := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1000000)) // 1000 * 10^6
	tx1, err := stablecoin.Mint(opts, opts.From, amount)
	if err != nil {
		return "", fmt.Errorf("mint stablecoin failed: %w", err)
	}

	return tx1.Hash().Hex(), nil
}

// TransferUSDT sends stablecoin from the current wallet to a recipient
// Note: Function name kept for backward compatibility, works with MockStablecoin (testnet) or USDC (mainnet)
func (s *DeAIService) TransferUSDT(to string, amountStr string) (string, error) {
	// 1. Resolve Addresses
	// Use constant which automatically switches based on environment
	stablecoinAddr := common.HexToAddress(contracts.StablecoinAddress)
	recipient := common.HexToAddress(to)

	// 2. Parse Amount
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount format")
	}

	// 3. Get Opts
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 4. Load Contract
	stablecoin, err := binding.KawaiToken(stablecoinAddr.Hex(), s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load stablecoin contract: %w", err)
	}

	// 5. Transfer
	tx, err := stablecoin.Transfer(opts, recipient, amount)
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
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
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
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get opts: %w", err)
	}

	// 4. Load Contract Generic
	// We use KawaiToken wrapper because it satisfies standard ERC20 interface
	token, err := binding.KawaiToken(tokenAddress, s.reader)
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
	if !s.wallet.IsUnlocked() {
		return nil, fmt.Errorf("no wallet connected")
	}

	if s.kv == nil {
		return nil, fmt.Errorf("KV store not initialized")
	}

	userAddr := s.wallet.GetCurrentAccountAddress()
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
	result.TotalKawaiClaimableFormatted = s.formatRewardAmount(result.TotalKawaiClaimable, types.RewardTypeMining)
	result.TotalUSDTClaimableFormatted = s.formatRewardAmount(result.TotalUSDTClaimable, types.RewardTypeRevenue)

	return result, nil
}

// convertMerkleProofToClaimable converts store.MerkleProofData to ClaimableReward
func (s *DeAIService) convertMerkleProofToClaimable(proof *store.MerkleProofData) *ClaimableReward {
	rewardType := types.RewardType(proof.RewardType)
	decimals := rewardType.Decimals()

	claimable := &ClaimableReward{
		Index:       proof.Index,
		Amount:      proof.Amount,
		Proof:       proof.Proof,
		MerkleRoot:  proof.MerkleRoot,
		PeriodID:    proof.PeriodID,
		RewardType:  rewardType,
		ClaimStatus: string(proof.ClaimStatus),
		ClaimTxHash: proof.ClaimTxHash,
		CreatedAt:   proof.CreatedAt.Format("2006-01-02T15:04:05Z"),
		ClaimedAt:   proof.ClaimedAt.Format("2006-01-02T15:04:05Z"),
		Formatted:   s.formatRewardAmount(proof.Amount, rewardType),
		Decimals:    decimals,
	}

	// Add mining-specific fields for kawai rewards (9-field format)
	if rewardType.IsKawaiReward() {
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

// ClaimUSDTReward claims USDT rewards using a Merkle proof (revenue sharing)
func (s *DeAIService) ClaimUSDTReward(periodID int64, index uint64, amountStr string, proof []string) (*ClaimResult, error) {
	if !s.wallet.IsUnlocked() {
		return nil, fmt.Errorf("no wallet connected")
	}

	claimerAddr := s.wallet.GetCurrentAccountAddress()
	ctx := context.Background()

	// Mark claim as pending BEFORE submitting transaction
	if s.kv != nil {
		if err := s.kv.MarkClaimPending(ctx, claimerAddr, periodID, ""); err != nil {
			return nil, fmt.Errorf("failed to mark claim as pending: %w", err)
		}
	}

	// Load RevenueDistributor contract for revenue sharing
	distributor, err := binding.RevenueDistributor("RevenueDistributor", s.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to load revenue distributor: %w", err)
	}

	// Parse amount
	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount format")
	}

	// Convert proof strings to [32]byte array
	merkleProof := make([][32]byte, len(proof))
	for i, p := range proof {
		if len(p) >= 2 && p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		if len(proofBytes) != 32 {
			return nil, fmt.Errorf("invalid proof element at index %d: expected 32 bytes, got %d", i, len(proofBytes))
		}
		copy(merkleProof[i][:], proofBytes)
	}

	// Get transaction options
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}

	// Submit claim transaction (fix integer truncation risk)
	idx := new(big.Int).SetUint64(index)
	tx, err := distributor.Claim(opts, idx, opts.From, amount, merkleProof)
	if err != nil {
		if s.kv != nil {
			_ = s.kv.MarkClaimFailed(ctx, claimerAddr, periodID, err.Error())
		}
		if !isUserError(err) {
			s.alert.SendAlert("WARNING", "Claim",
				fmt.Sprintf("⚠️ USDT claim failed\n\nClaimer: %s\nPeriod: %d\nIndex: %d\nAmount: %s\nError: %v",
					claimerAddr, periodID, index, amountStr, err))
		}
		return nil, fmt.Errorf("claim transaction failed: %w", err)
	}

	txHash := tx.Hash().Hex()

	// Wait for confirmation
	receipt, err := bind.WaitMined(ctx, s.reader.Client(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction confirmation: %w", err)
	}

	if receipt.Status != 1 {
		if s.kv != nil {
			_ = s.kv.MarkClaimFailed(ctx, claimerAddr, periodID, "transaction reverted")
		}
		s.alert.SendAlert("ERROR", "Claim",
			fmt.Sprintf("❌ USDT claim reverted!\n\nClaimer: %s\nPeriod: %d\nTx: %s",
				claimerAddr, periodID, txHash))
		return nil, fmt.Errorf("transaction reverted (status: %d)", receipt.Status)
	}

	return &ClaimResult{
		TxHash:     txHash,
		PeriodID:   periodID,
		RewardType: types.RewardTypeRevenue,
		Amount:     amountStr,
		Status:     "confirmed",
	}, nil
}

// ClaimCashbackReward claims deposit cashback rewards using a Merkle proof
// Uses the DepositCashbackDistributor contract with period-based claims
func (s *DeAIService) ClaimCashbackReward(period uint64, kawaiAmount string, proof []string) (*ClaimResult, error) {

	if !s.wallet.IsUnlocked() {
		return nil, fmt.Errorf("no wallet connected")
	}

	userAddr := s.wallet.GetCurrentAccountAddress()

	// 1. Input validation
	// Note: Empty proof is valid for single-leaf Merkle trees
	if kawaiAmount == "" || kawaiAmount == "0" {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// 2. Load DepositCashbackDistributor contract
	distributor, err := binding.CashbackDistributor("CashbackDistributor", s.reader)
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
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}

	// 5. Mark claim as pending BEFORE submitting transaction (prevents double-claim UX issue)
	ctx := context.Background()
	if s.kv != nil {
		kvStore, ok := s.kv.(*store.KVStore)
		if ok {
			if err := kvStore.MarkCashbackPending(ctx, userAddr, period, ""); err != nil {
				return nil, fmt.Errorf("failed to mark cashback claim as pending: %w", err)
			}
		}
	}

	// 6. Submit claim transaction
	tx, err := distributor.ClaimCashback(opts, new(big.Int).SetUint64(period), amount, merkleProof)
	if err != nil {
		// Rollback pending status on transaction failure
		if s.kv != nil {
			kvStore, ok := s.kv.(*store.KVStore)
			if ok {
				_ = kvStore.MarkCashbackFailed(ctx, userAddr, period, err.Error())
			}
		}

		// Alert on unexpected errors only
		if !isUserError(err) {
			s.alert.SendAlert("WARNING", "Claim",
				fmt.Sprintf("⚠️ Cashback claim failed\n\nUser: %s\nPeriod: %d\nAmount: %s KAWAI\nError: %v",
					userAddr, period, kawaiAmount, err))
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
				_ = kvStore.MarkCashbackFailed(ctx, userAddr, period, "transaction reverted")
			}
		}

		// Alert on transaction revert
		s.alert.SendAlert("ERROR", "Claim",
			fmt.Sprintf("❌ Cashback claim reverted!\n\nUser: %s\nPeriod: %d\nAmount: %s KAWAI\nTx: %s\nGas Used: %d",
				userAddr, period, kawaiAmount, txHash, receipt.GasUsed))

		return nil, fmt.Errorf("transaction reverted (status: %d)", receipt.Status)
	}

	// 9. Mark claim as completed in KV store
	if s.kv != nil {
		kvStore, ok := s.kv.(*store.KVStore)
		if ok {
			if err := kvStore.MarkCashbackClaimed(ctx, userAddr, period); err != nil {
				// Log warning - tx was successful, but KV update failed
				fmt.Printf("Warning: failed to mark cashback claim in KV: %v\n", err)
			}
		}
	}

	return &ClaimResult{
		TxHash:     txHash,
		PeriodID:   int64(period),
		RewardType: types.RewardTypeCashback,
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

	if !s.wallet.IsUnlocked() {
		return nil, fmt.Errorf("no wallet connected")
	}

	claimerAddr := s.wallet.GetCurrentAccountAddress()

	// Contract now supports timestamp-based periods directly via setMerkleRootForPeriod()
	// No mapping needed - use timestamp period ID as-is
	contractPeriod := period

	fmt.Printf("🔄 Claiming mining reward: period %d (timestamp-based)\n", period)

	// Load MiningRewardDistributor contract
	distributor, err := binding.MiningRewardDistributor("MiningRewardDistributor", s.reader)
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
	chainId := s.chainID
	opts, err := s.wallet.GetTransactOpts(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}

	// 7. Mark claim as pending BEFORE submitting transaction (prevents double-claim UX issue)
	ctx := context.Background()
	if s.kv != nil {
		if err := s.kv.MarkClaimPending(ctx, claimerAddr, period, ""); err != nil {
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
			_ = s.kv.MarkClaimFailed(ctx, claimerAddr, period, err.Error())
		}

		// Alert on unexpected errors only
		if !isUserError(err) {
			s.alert.SendAlert("WARNING", "Claim",
				fmt.Sprintf("⚠️ Mining claim failed\n\nClaimer: %s\nPeriod: %d\nContributor: %s KAWAI\nError: %v",
					claimerAddr, period, contributorAmount, err))
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
			s.kv.MarkClaimFailed(ctx, claimerAddr, period, "transaction reverted")
		}

		// Alert on transaction revert
		s.alert.SendAlert("ERROR", "Claim",
			fmt.Sprintf("❌ Mining claim reverted!\n\nClaimer: %s\nPeriod: %d\nContributor: %s KAWAI\nTx: %s\nGas Used: %d",
				claimerAddr, period, contributorAmount, txHash, receipt.GasUsed))

		return nil, fmt.Errorf("transaction reverted (status: %d)", receipt.Status)
	}

	// 11. Update with actual tx hash (status remains pending, auto-confirm will update to confirmed)
	if s.kv != nil {
		if err := s.kv.MarkClaimPending(ctx, claimerAddr, period, txHash); err != nil {
			// Log warning - tx was successful, auto-confirm will fix status later
			fmt.Printf("Warning: failed to update mining claim tx hash in KV: %v\n", err)
		}
	}

	return &ClaimResult{
		TxHash:     txHash,
		PeriodID:   period,
		RewardType: types.RewardTypeMining,
		Amount:     contributorAmount,
		Status:     "confirmed",
	}, nil
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
	if !s.wallet.IsUnlocked() {
		return fmt.Errorf("no wallet connected")
	}

	if s.kv == nil {
		return fmt.Errorf("KV store not initialized")
	}

	ctx := context.Background()
	userAddr := s.wallet.GetCurrentAccountAddress()

	return s.kv.ConfirmClaim(ctx, userAddr, periodID)
}

// MarkClaimFailed marks a claim as failed after the transaction reverts
func (s *DeAIService) MarkClaimFailed(periodID int64, reason string) error {
	if !s.wallet.IsUnlocked() {
		return fmt.Errorf("no wallet connected")
	}

	if s.kv == nil {
		return fmt.Errorf("KV store not initialized")
	}

	ctx := context.Background()
	userAddr := s.wallet.GetCurrentAccountAddress()

	return s.kv.MarkClaimFailed(ctx, userAddr, periodID, reason)
}

// formatRewardAmount formats raw token amount to human-readable format
func (s *DeAIService) formatRewardAmount(rawAmount string, rewardType types.RewardType) string {
	amount := new(big.Int)
	amount, ok := amount.SetString(rawAmount, 10)
	if !ok {
		return "0.00"
	}

	decimals := int64(rewardType.Decimals())

	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil))
	formatted := new(big.Float).Quo(new(big.Float).SetInt(amount), divisor)

	// Format with appropriate precision
	if rewardType.IsStablecoinReward() {
		return formatted.Text('f', 2)
	}
	return formatted.Text('f', 4)
}

// =============================================================================
// CONTRIBUTOR DISCOVERY METHODS
// =============================================================================

// ContributorInfo represents contributor information for frontend display
type ContributorInfo struct {
	WalletAddress   string   `json:"wallet_address"`
	EndpointURL     string   `json:"endpoint_url"`
	Region          string   `json:"region"`
	Status          string   `json:"status"`
	LastSeen        string   `json:"last_seen"`
	AvailableModels []string `json:"available_models"`
	ActiveRequests  int64    `json:"active_requests"`
	TotalRequests   int64    `json:"total_requests"`
	AvgResponseTime float64  `json:"avg_response_time"`
	SuccessRate     float64  `json:"success_rate"`
	CPUCores        int      `json:"cpu_cores"`
	TotalRAM        int64    `json:"total_ram"`
	AvailableRAM    int64    `json:"available_ram"`
	GPUModel        string   `json:"gpu_model"`
	GPUMemory       int64    `json:"gpu_memory"`
	Score           float64  `json:"score,omitempty"`
}

// GetAvailableContributors returns all online contributors with their scores
func (s *DeAIService) GetAvailableContributors() ([]*ContributorInfo, error) {
	ctx := context.Background()

	scores, err := s.contributorSelector.GetAvailableContributors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributors: %w", err)
	}

	contributors := make([]*ContributorInfo, 0, len(scores))
	for _, score := range scores {
		c := score.Contributor
		contributors = append(contributors, &ContributorInfo{
			WalletAddress:   c.WalletAddress,
			EndpointURL:     c.EndpointURL,
			Region:          c.Region,
			Status:          string(c.Status),
			LastSeen:        c.LastSeen.Format(time.RFC3339),
			AvailableModels: c.AvailableModels,
			ActiveRequests:  c.ActiveRequests,
			TotalRequests:   c.TotalRequests,
			AvgResponseTime: c.AvgResponseTime,
			SuccessRate:     c.SuccessRate,
			CPUCores:        c.CPUCores,
			TotalRAM:        c.TotalRAM,
			AvailableRAM:    c.AvailableRAM,
			GPUModel:        c.GPUModel,
			GPUMemory:       c.GPUMemory,
			Score:           score.Score,
		})
	}

	return contributors, nil
}

// SelectBestContributor selects the best contributor based on criteria
func (s *DeAIService) SelectBestContributor(preferredRegion string, requiredModel string, minRAM int64, minGPUMemory int64) (*ContributorInfo, error) {
	ctx := context.Background()

	criteria := &SelectionCriteria{
		PreferredRegion: preferredRegion,
		RequiredModel:   requiredModel,
		MinRAM:          minRAM,
		MinGPUMemory:    minGPUMemory,
		MaxLoad:         10, // Max 10 concurrent requests per contributor
	}

	contributor, err := s.contributorSelector.SelectBestContributor(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to select contributor: %w", err)
	}

	return &ContributorInfo{
		WalletAddress:   contributor.WalletAddress,
		EndpointURL:     contributor.EndpointURL,
		Region:          contributor.Region,
		Status:          string(contributor.Status),
		LastSeen:        contributor.LastSeen.Format(time.RFC3339),
		AvailableModels: contributor.AvailableModels,
		ActiveRequests:  contributor.ActiveRequests,
		TotalRequests:   contributor.TotalRequests,
		AvgResponseTime: contributor.AvgResponseTime,
		SuccessRate:     contributor.SuccessRate,
		CPUCores:        contributor.CPUCores,
		TotalRAM:        contributor.TotalRAM,
		AvailableRAM:    contributor.AvailableRAM,
		GPUModel:        contributor.GPUModel,
		GPUMemory:       contributor.GPUMemory,
	}, nil
}

// GetContributorStats returns statistics about the contributor network
func (s *DeAIService) GetContributorStats() (map[string]interface{}, error) {
	ctx := context.Background()

	stats, err := s.contributorSelector.GetContributorStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return stats, nil
}
