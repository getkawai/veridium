package services_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWalletProvider is a mock implementation of WalletProvider interface
// This demonstrates how the new interface enables easy testing
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

// TestWalletProviderInterface demonstrates how to mock the wallet provider
func TestWalletProviderInterface(t *testing.T) {
	// Create mock wallet
	mockWallet := new(MockWalletProvider)

	// Setup expectations
	mockWallet.On("HasWallet").Return(true)
	mockWallet.On("GetCurrentAddress").Return("0x1234567890abcdef")
	mockWallet.On("IsUnlocked").Return(true)
	mockWallet.On("GetCurrentAccountAddress").Return("0x1234567890abcdef")

	// Test the mock
	assert.True(t, mockWallet.HasWallet())
	assert.Equal(t, "0x1234567890abcdef", mockWallet.GetCurrentAddress())
	assert.True(t, mockWallet.IsUnlocked())
	assert.Equal(t, "0x1234567890abcdef", mockWallet.GetCurrentAccountAddress())

	// Verify all expectations were met
	mockWallet.AssertExpectations(t)
}

// Example: Testing DeAIService with mocked wallet
func TestDeAIServiceWithMockWallet(t *testing.T) {
	// This test demonstrates how to test DeAIService with a mocked wallet
	// without needing actual wallet infrastructure

	mockWallet := new(MockWalletProvider)
	mockWallet.On("IsUnlocked").Return(true)
	mockWallet.On("GetCurrentAccountAddress").Return("0x1234567890abcdef")

	// Now you can create DeAIService with the mock wallet
	// service := services.NewDeAIService(mockWallet, nil)

	// And test service methods without real blockchain
	assert.True(t, mockWallet.IsUnlocked())
	assert.Equal(t, "0x1234567890abcdef", mockWallet.GetCurrentAccountAddress())
}
