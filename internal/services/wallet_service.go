package services

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"

	"github.com/kawai-network/veridium/pkg/jarvis/accounts"
	"github.com/kawai-network/veridium/pkg/jarvis/accounts/types"
	"github.com/kawai-network/veridium/pkg/jarvis/util/account"
)

// WalletStatus represents the current state of the wallet
type WalletStatus struct {
	HasWallet bool   `json:"hasWallet"`
	IsLocked  bool   `json:"isLocked"`
	Address   string `json:"address"`
}

// WalletService handles local wallet management using jarvis/accounts
type WalletService struct {
	currentAccount *account.Account
	address        string
}

// NewWalletService creates a new wallet service
func NewWalletService(dataDir string) *WalletService {
	// dataDir is not used anymore, jarvis/accounts manages its own directory
	return &WalletService{}
}

// HasWallet checks if a wallet already exists
func (s *WalletService) HasWallet() bool {
	accs := accounts.GetAccounts()
	return len(accs) > 0
}

// GenerateMnemonic creates a new 12-word bip39 mnemonic
func (s *WalletService) GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// SetupWallet creates a new keystore from a password and mnemonic
func (s *WalletService) SetupWallet(password string, mnemonic string) (string, error) {
	if s.HasWallet() {
		return "", errors.New("wallet already exists")
	}

	// 1. Derive private key from mnemonic
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := crypto.ToECDSA(seed[:32]) // Simple derivation for demo, usually use BIP44 path
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %v", err)
	}

	privBytes := crypto.FromECDSA(masterKey)
	privHex := fmt.Sprintf("%x", privBytes)
	address := crypto.PubkeyToAddress(masterKey.PublicKey).Hex()

	// 2. Store keystore using jarvis/accounts
	keystorePath, err := accounts.StorePrivateKeyWithKeystore(privHex, password)
	if err != nil {
		return "", fmt.Errorf("failed to store keystore: %v", err)
	}

	// 3. Verify keystore
	verifiedAddr, err := accounts.VerifyKeystore(keystorePath)
	if err != nil {
		return "", fmt.Errorf("failed to verify keystore: %v", err)
	}

	// 4. Store account metadata
	accDesc := types.AccDesc{
		Address: verifiedAddr,
		Kind:    "keystore",
		Keypath: keystorePath,
		Desc:    "Veridium Wallet",
	}
	if err := accounts.StoreAccountRecord(accDesc); err != nil {
		return "", fmt.Errorf("failed to store account record: %v", err)
	}

	// 5. Unlock the account
	acc, err := accounts.UnlockKeystoreAccountWithPassword(accDesc, password)
	if err != nil {
		return "", fmt.Errorf("failed to unlock account: %v", err)
	}

	// 6. Set state
	s.currentAccount = acc
	s.address = address

	return address, nil
}

// UnlockWallet decrypts the keystore and loads it into memory
func (s *WalletService) UnlockWallet(password string) (string, error) {
	// Get all accounts
	accs := accounts.GetAccounts()
	if len(accs) == 0 {
		return "", errors.New("no wallet found")
	}

	// Get the first account (assuming single wallet for now)
	var accDesc types.AccDesc
	for _, acc := range accs {
		accDesc = acc
		break
	}

	// Unlock the account
	acc, err := accounts.UnlockKeystoreAccountWithPassword(accDesc, password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Set state
	s.currentAccount = acc
	s.address = acc.AddressHex()

	return s.address, nil
}

// LockWallet clears the private key from memory
func (s *WalletService) LockWallet() {
	s.currentAccount = nil
	s.address = ""
}

// GetStatus returns the current wallet status
func (s *WalletService) GetStatus() WalletStatus {
	return WalletStatus{
		HasWallet: s.HasWallet(),
		IsLocked:  s.currentAccount == nil,
		Address:   s.address,
	}
}

// SignMessage signs a message with the private key
func (s *WalletService) SignMessage(message string) (string, error) {
	if s.currentAccount == nil {
		return "", errors.New("wallet is locked")
	}

	return s.currentAccount.SignMessage(message)
}
