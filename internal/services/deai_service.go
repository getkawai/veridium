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
)

// monadChainID is the chain ID for Monad Testnet
var monadChainID = big.NewInt(int64(networks.MonadTestnet.GetChainID()))

// DeAIService handles interactions with the Veridium smart contracts
type DeAIService struct {
	reader *reader.EthReader
	wallet *WalletService
}

// NewDeAIService creates a new instance of DeAIService
func NewDeAIService(wallet *WalletService) *DeAIService {
	// Initialize EthReader with Monad Testnet nodes from jarvis network config
	ethReader := reader.NewEthReaderGeneric(networks.MonadTestnet.GetDefaultNodes(), nil)

	return &DeAIService{
		reader: ethReader,
		wallet: wallet,
	}
}

// GetVaultBalance returns the USDT balance of the current wallet
func (s *DeAIService) GetVaultBalance() (string, error) {
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

// DepositToVault deposits USDT into the vault for service credits
func (s *DeAIService) DepositToVault(amountStr string) (string, error) {
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
