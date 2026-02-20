// Package app provides the core application context and service initialization.
// This file defines the service interfaces used throughout the application
// to enable dependency injection and easier testing.
package app

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/contracts/otcmarket"
	"github.com/kawai-network/veridium/internal/services"
)

// ============================================================================
// DATABASE INTERFACE
// ============================================================================

// Querier defines the minimal interface for database query operations.
// This is what most services actually need from the database.
type Querier interface {
	// Queries returns the generated queries interface for raw SQL operations
	Queries() *db.Queries
}

// Database defines the full interface for database operations.
// This abstraction allows for easy mocking in tests and future database swaps.
type Database interface {
	Querier

	// DB returns the underlying database connection for custom queries
	DB() *sql.DB

	// Close gracefully closes the database connection
	Close() error

	// WithTx executes a function within a database transaction
	WithTx(ctx context.Context, fn func(*db.Queries) error) error
}

// databaseWrapper is a minimal struct that implements the Database interface.
// It demonstrates the interface contract and can be used for testing or as a
// reference implementation. Each field is a function that implements one
// interface method, allowing flexible composition.
type databaseWrapper struct {
	queriesFn func() *db.Queries
	dbFn      func() *sql.DB
	closeFn   func() error
	withTxFn  func(context.Context, func(*db.Queries) error) error
}

func (d *databaseWrapper) Queries() *db.Queries { return d.queriesFn() }
func (d *databaseWrapper) DB() *sql.DB          { return d.dbFn() }
func (d *databaseWrapper) Close() error         { return d.closeFn() }
func (d *databaseWrapper) WithTx(ctx context.Context, fn func(*db.Queries) error) error {
	return d.withTxFn(ctx, fn)
}

// Compile-time check: ensure databaseWrapper implements Database
var _ Database = (*databaseWrapper)(nil)

// ============================================================================
// WALLET PROVIDER INTERFACE
// ============================================================================

// WalletProvider defines the interface for wallet management operations.
// This abstraction allows for testing without real wallet infrastructure.
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
	GetStatus() services.WalletStatus

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

// WalletInfo represents a wallet in the list
// Alias to services.WalletInfo for compatibility
type WalletInfo = services.WalletInfo

// WalletStatus represents the current state of the wallet
// Alias to services.WalletStatus for compatibility
type WalletStatus = services.WalletStatus

// ============================================================================
// BLOCKCHAIN PROVIDER INTERFACE
// ============================================================================

// BlockchainProvider defines the interface for blockchain operations.
// This abstraction allows for testing without real blockchain connections
type BlockchainProvider interface {
	// Connection State
	IsConnected() bool
	GetChainID() *big.Int

	// Token Operations
	GetTotalSupply(ctx context.Context) (*big.Int, error)
	GetMaxSupply(ctx context.Context) (*big.Int, error)
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
