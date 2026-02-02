package services_test

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/otcmarket"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// Shared Mock Implementations for Testing
// ============================================================================

// MockWalletProvider is a mock implementation of WalletProvider interface
type MockWalletProvider struct {
	mock.Mock
}

func (m *MockWalletProvider) HasWallet() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockWalletProvider) GetWallets() []services.WalletInfo {
	args := m.Called()
	return args.Get(0).([]services.WalletInfo)
}

func (m *MockWalletProvider) GenerateMnemonic() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) CreateWallet(password string, mnemonic string, description string) (string, error) {
	args := m.Called(password, mnemonic, description)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) SetupWallet(password string, mnemonic string, name string) (string, error) {
	args := m.Called(password, mnemonic, name)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) SwitchWallet(address string, password string) (string, error) {
	args := m.Called(address, password)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) UnlockWallet(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) LockWallet() {
	m.Called()
}

func (m *MockWalletProvider) DeleteWallet(address string) error {
	args := m.Called(address)
	return args.Error(0)
}

func (m *MockWalletProvider) ExportKeystore(address string) (string, error) {
	args := m.Called(address)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) ImportKeystore(keystoreJSON string, password string, description string) (string, error) {
	args := m.Called(keystoreJSON, password, description)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) ImportPrivateKey(privateKeyHex string, password string, description string) (string, error) {
	args := m.Called(privateKeyHex, password, description)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) UpdateWalletDescription(address string, description string) error {
	args := m.Called(address, description)
	return args.Error(0)
}

func (m *MockWalletProvider) GetStatus() services.WalletStatus {
	args := m.Called()
	return args.Get(0).(services.WalletStatus)
}

func (m *MockWalletProvider) GetCurrentAddress() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockWalletProvider) SignMessage(message string) (string, error) {
	args := m.Called(message)
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) GetTransactOpts(chainId *big.Int) (*bind.TransactOpts, error) {
	args := m.Called(chainId)
	return args.Get(0).(*bind.TransactOpts), args.Error(1)
}

func (m *MockWalletProvider) GetAPIKey() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockWalletProvider) AutoClaimTrialIfNeeded(referralCode string) (bool, float64, string, error) {
	args := m.Called(referralCode)
	return args.Bool(0), args.Get(1).(float64), args.String(2), args.Error(3)
}

func (m *MockWalletProvider) IsUnlocked() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockWalletProvider) GetCurrentAccountAddress() string {
	args := m.Called()
	return args.String(0)
}

// MockBlockchainProvider is a mock implementation of BlockchainProvider interface
type MockBlockchainProvider struct {
	mock.Mock
}

func (m *MockBlockchainProvider) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockBlockchainProvider) GetChainID() *big.Int {
	args := m.Called()
	return args.Get(0).(*big.Int)
}

func (m *MockBlockchainProvider) GetKawaiTokenBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockBlockchainProvider) GetUSDTBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockBlockchainProvider) ValidateTradeBalance(ctx context.Context, buyer common.Address, stablecoinAmount *big.Int) error {
	args := m.Called(ctx, buyer, stablecoinAmount)
	return args.Error(0)
}

func (m *MockBlockchainProvider) ValidateOrderCreationBalance(ctx context.Context, seller common.Address, tokenAmount *big.Int) error {
	args := m.Called(ctx, seller, tokenAmount)
	return args.Error(0)
}

func (m *MockBlockchainProvider) MarketplaceCreateOrder(ctx context.Context, transactOpts *bind.TransactOpts, tokenAmount, stablecoinPrice *big.Int) (*types.Transaction, error) {
	args := m.Called(ctx, transactOpts, tokenAmount, stablecoinPrice)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceBuyOrder(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int) (*types.Transaction, error) {
	args := m.Called(ctx, transactOpts, orderID)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceBuyOrderPartial(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int, amount *big.Int) (*types.Transaction, error) {
	args := m.Called(ctx, transactOpts, orderID, amount)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceCancelOrder(ctx context.Context, transactOpts *bind.TransactOpts, orderID *big.Int) (*types.Transaction, error) {
	args := m.Called(ctx, transactOpts, orderID)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceGetOrder(ctx context.Context, orderID *big.Int) (otcmarket.OTCMarketOrder, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(otcmarket.OTCMarketOrder), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceGetOrdersCount(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceGetOrdersBySeller(ctx context.Context, seller common.Address, offset, limit *big.Int) ([]otcmarket.OTCMarketOrder, error) {
	args := m.Called(ctx, seller, offset, limit)
	return args.Get(0).([]otcmarket.OTCMarketOrder), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceGetActiveOrders(ctx context.Context, offset, limit *big.Int) ([]otcmarket.OTCMarketOrder, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]otcmarket.OTCMarketOrder), args.Error(1)
}

func (m *MockBlockchainProvider) MarketplaceGetOrders(ctx context.Context, orderIDs []*big.Int) ([]otcmarket.OTCMarketOrder, error) {
	args := m.Called(ctx, orderIDs)
	return args.Get(0).([]otcmarket.OTCMarketOrder), args.Error(1)
}

func (m *MockBlockchainProvider) GetEthClient() *ethclient.Client {
	args := m.Called()
	return args.Get(0).(*ethclient.Client)
}
