package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/contracts/otcmarket"
)

// BlockchainProvider defines the interface for blockchain operations.
// This abstraction allows for testing without real blockchain connections
// and breaks import cycles between internal/app and internal/services.
type BlockchainProvider interface {
	// Connection State
	IsConnected() bool
	GetChainID() *big.Int

	// Token Operations
	GetKawaiTokenBalance(ctx context.Context, address common.Address) (*big.Int, error)

	// Stablecoin Operations
	GetUSDTBalance(ctx context.Context, address common.Address) (*big.Int, error)

	// Balance Validation
	ValidateTradeBalance(ctx context.Context, buyer common.Address, stablecoinAmount *big.Int) error
	ValidateOrderCreationBalance(ctx context.Context, seller common.Address, tokenAmount *big.Int) error

	// Marketplace Operations
	MarketplaceCreateOrder(ctx context.Context, transactOpts *bind.TransactOpts, tokenAmount, stablecoinPrice *big.Int) (*types.Transaction, error)
	MarketplaceBuyOrder(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int) (*types.Transaction, error)
	MarketplaceBuyOrderPartial(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int, amount *big.Int) (*types.Transaction, error)
	MarketplaceCancelOrder(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int) (*types.Transaction, error)

	// Marketplace Query Operations
	MarketplaceGetOrder(ctx context.Context, orderID *big.Int) (otcmarket.OTCMarketOrder, error)
	MarketplaceGetOrdersCount(ctx context.Context) (*big.Int, error)
	MarketplaceGetOrdersBySeller(ctx context.Context, seller common.Address, offset, limit *big.Int) ([]otcmarket.OTCMarketOrder, error)
	MarketplaceGetActiveOrders(ctx context.Context, offset, limit *big.Int) ([]otcmarket.OTCMarketOrder, error)
	MarketplaceGetOrders(ctx context.Context, orderIDs []*big.Int) ([]otcmarket.OTCMarketOrder, error)

	// Low-level access (for advanced use cases)
	GetEthClient() *ethclient.Client
}

// WalletProvider defines the interface for wallet management operations.
// This abstraction allows for testing without real wallet infrastructure
// and breaks import cycles between internal/app and internal/services.
type WalletProvider interface {
	// HasWallet checks if a wallet already exists
	HasWallet() bool

	// GetWallets returns a list of all stored wallets
	GetWallets() []WalletInfo

	// GenerateMnemonic creates a new 12-word bip39 mnemonic
	GenerateMnemonic() (string, error)

	// CreateWallet creates a new wallet (supports multiple wallets)
	CreateWallet(password string, mnemonic string, description string) (string, error)

	// SetupWallet creates a new keystore from a password and mnemonic (first wallet only)
	SetupWallet(password string, mnemonic string, name string) (string, error)

	// SwitchWallet switches to a different wallet by address
	SwitchWallet(address string, password string) (string, error)

	// UnlockWallet decrypts the keystore and loads it into memory
	UnlockWallet(password string) (string, error)

	// LockWallet clears the private key from memory
	LockWallet()

	// DeleteWallet removes a wallet from storage
	DeleteWallet(address string) error

	// ExportKeystore returns the keystore JSON content for a wallet
	ExportKeystore(address string) (string, error)

	// ImportKeystore imports a keystore from JSON content
	ImportKeystore(keystoreJSON string, password string, description string) (string, error)

	// ImportPrivateKey imports a wallet from a private key hex string
	ImportPrivateKey(privateKeyHex string, password string, description string) (string, error)

	// UpdateWalletDescription updates the description for a wallet
	UpdateWalletDescription(address string, description string) error

	// GetStatus returns the current state of the wallet
	GetStatus() WalletStatus

	// GetCurrentAddress returns the currently unlocked wallet address
	GetCurrentAddress() string

	// SignMessage signs a message with the current wallet
	SignMessage(message string) (string, error)

	// GetTransactOpts returns transaction options for the current wallet
	// This is used for blockchain transactions
	GetTransactOpts(chainId *big.Int) (*bind.TransactOpts, error)

	// GetAPIKey returns the API key for the current wallet (generating one if needed)
	GetAPIKey() (string, error)

	// AutoClaimTrialIfNeeded attempts to auto-claim trial tokens for new wallets
	AutoClaimTrialIfNeeded(referralCode string) (bool, float64, string, error)

	// IsUnlocked returns true if the wallet is currently unlocked
	IsUnlocked() bool

	// GetCurrentAccountAddress returns the current account address if unlocked
	GetCurrentAccountAddress() string
}
