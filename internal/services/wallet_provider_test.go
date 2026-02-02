package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: MockWalletProvider is now in test_mocks.go

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
