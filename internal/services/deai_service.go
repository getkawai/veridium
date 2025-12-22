package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common" // Added this import
	"github.com/kawai-network/veridium/pkg/jarvis/contracts"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
)

// DeAIService handles interactions with the Veridium smart contracts
type DeAIService struct {
	reader *reader.EthReader
	wallet *WalletService
}

// NewDeAIService creates a new instance of DeAIService
func NewDeAIService(wallet *WalletService) *DeAIService {
	// Initialize EthReader with default public nodes for BSC Testnet
	nodes := map[string]string{
		"bsctestnet":  "https://bsc-testnet-rpc.publicnode.com",
		"bsctestnet2": "https://data-seed-prebsc-1-s1.binance.org:8545",
	}
	ethReader := reader.NewEthReaderGeneric(nodes, nil)

	return &DeAIService{
		reader: ethReader,
		wallet: wallet,
	}
}

// GetVaultBalance returns the balance of the payment vault contract in USDT
// This is just a placeholder example, real logic depends on what we want to read
func (s *DeAIService) GetVaultBalance() (string, error) {
	vault, err := contracts.Vault("PaymentVault", s.reader)
	if err != nil {
		return "", fmt.Errorf("failed to load vault: %w", err)
	}

	// For now, let's just return the owner address as a test
	owner, err := vault.Owner(nil)
	if err != nil {
		return "", fmt.Errorf("failed to get owner: %w", err)
	}

	return owner.Hex(), nil
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
		chainId := big.NewInt(97)
		opts, err := s.wallet.GetTransactOpts(chainId)
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
	chainId := big.NewInt(97)
	opts, err := s.wallet.GetTransactOpts(chainId)
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
	chainId := big.NewInt(97)
	opts, err := s.wallet.GetTransactOpts(chainId)
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
	chainId := big.NewInt(97)
	opts, err := s.wallet.GetTransactOpts(chainId)
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

	chainId := big.NewInt(97)
	opts, err := s.wallet.GetTransactOpts(chainId)
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

	chainId := big.NewInt(97)
	opts, err := s.wallet.GetTransactOpts(chainId)
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
	chainId := big.NewInt(97)
	opts, err := s.wallet.GetTransactOpts(chainId)
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
