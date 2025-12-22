package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"

	"github.com/kawai-network/veridium/pkg/jarvis/accounts"
	"github.com/kawai-network/veridium/pkg/jarvis/accounts/types"
	"github.com/kawai-network/veridium/pkg/jarvis/util/account"
)

// WalletInfo represents a wallet in the list
type WalletInfo struct {
	Address     string `json:"address"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

// WalletStatus represents the current state of the wallet
type WalletStatus struct {
	HasWallet bool         `json:"hasWallet"`
	IsLocked  bool         `json:"isLocked"`
	Address   string       `json:"address"`
	Wallets   []WalletInfo `json:"wallets"`
}

// WalletService handles local wallet management using jarvis/accounts
type WalletService struct {
	currentAccount *account.Account
	address        string
}

// NewWalletService creates a new wallet service
func NewWalletService(dataDir string) *WalletService {
	return &WalletService{}
}

// HasWallet checks if a wallet already exists
func (s *WalletService) HasWallet() bool {
	accs := accounts.GetAccounts()
	return len(accs) > 0
}

// GetWallets returns a list of all stored wallets
func (s *WalletService) GetWallets() []WalletInfo {
	accs := accounts.GetAccounts()
	result := make([]WalletInfo, 0, len(accs))
	for addr, acc := range accs {
		result = append(result, WalletInfo{
			Address:     addr,
			Description: acc.Desc,
			IsActive:    s.address == addr,
		})
	}
	return result
}

// GenerateMnemonic creates a new 12-word bip39 mnemonic
func (s *WalletService) GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// CreateWallet creates a new wallet (supports multiple wallets)
func (s *WalletService) CreateWallet(password string, mnemonic string, description string) (string, error) {
	// Validate mnemonic
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", errors.New("invalid mnemonic")
	}

	// Derive private key from mnemonic
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %v", err)
	}

	privBytes := crypto.FromECDSA(masterKey)
	privHex := fmt.Sprintf("%x", privBytes)
	address := crypto.PubkeyToAddress(masterKey.PublicKey).Hex()

	// Check if wallet already exists
	accs := accounts.GetAccounts()
	if _, exists := accs[address]; exists {
		return "", errors.New("wallet with this address already exists")
	}

	// Store keystore
	keystorePath, err := accounts.StorePrivateKeyWithKeystore(privHex, password)
	if err != nil {
		return "", fmt.Errorf("failed to store keystore: %v", err)
	}

	// Verify keystore
	verifiedAddr, err := accounts.VerifyKeystore(keystorePath)
	if err != nil {
		return "", fmt.Errorf("failed to verify keystore: %v", err)
	}

	// Use default description if empty
	if description == "" {
		description = fmt.Sprintf("Wallet %d", len(accs)+1)
	}

	// Store account metadata
	accDesc := types.AccDesc{
		Address: verifiedAddr,
		Kind:    "keystore",
		Keypath: keystorePath,
		Desc:    description,
	}
	if err := accounts.StoreAccountRecord(accDesc); err != nil {
		return "", fmt.Errorf("failed to store account record: %v", err)
	}

	// Auto-unlock the new wallet
	acc, err := accounts.UnlockKeystoreAccountWithPassword(accDesc, password)
	if err != nil {
		return "", fmt.Errorf("failed to unlock account: %v", err)
	}

	s.currentAccount = acc
	s.address = address

	return address, nil
}

// SetupWallet creates a new keystore from a password and mnemonic (legacy, wraps CreateWallet)
func (s *WalletService) SetupWallet(password string, mnemonic string) (string, error) {
	if s.HasWallet() {
		return "", errors.New("wallet already exists")
	}
	return s.CreateWallet(password, mnemonic, "Veridium Wallet")
}

// SwitchWallet switches to a different wallet by address
func (s *WalletService) SwitchWallet(address string, password string) (string, error) {
	accs := accounts.GetAccounts()
	accDesc, exists := accs[address]
	if !exists {
		return "", errors.New("wallet not found")
	}

	acc, err := accounts.UnlockKeystoreAccountWithPassword(accDesc, password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	s.currentAccount = acc
	s.address = acc.AddressHex()

	return s.address, nil
}

// UnlockWallet decrypts the keystore and loads it into memory
func (s *WalletService) UnlockWallet(password string) (string, error) {
	accs := accounts.GetAccounts()
	if len(accs) == 0 {
		return "", errors.New("no wallet found")
	}

	// Get the first account if none is active
	var accDesc types.AccDesc
	for _, acc := range accs {
		accDesc = acc
		break
	}

	acc, err := accounts.UnlockKeystoreAccountWithPassword(accDesc, password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	s.currentAccount = acc
	s.address = acc.AddressHex()

	return s.address, nil
}

// LockWallet clears the private key from memory
func (s *WalletService) LockWallet() {
	s.currentAccount = nil
	s.address = ""
}

// DeleteWallet removes a wallet from storage
func (s *WalletService) DeleteWallet(address string) error {
	accs := accounts.GetAccounts()
	accDesc, exists := accs[address]
	if !exists {
		return errors.New("wallet not found")
	}

	// Cannot delete active wallet
	if s.address == address {
		return errors.New("cannot delete active wallet, switch to another wallet first")
	}

	// Delete keystore file
	if err := os.Remove(accDesc.Keypath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete keystore: %v", err)
	}

	// Delete account record (metadata file)
	homeDir, _ := os.UserHomeDir()
	metadataPath := fmt.Sprintf("%s/.jarvis/%s.json", homeDir, address)
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata: %v", err)
	}

	return nil
}

// ExportKeystore returns the keystore JSON content for a wallet
func (s *WalletService) ExportKeystore(address string) (string, error) {
	accs := accounts.GetAccounts()
	accDesc, exists := accs[address]
	if !exists {
		return "", errors.New("wallet not found")
	}

	content, err := os.ReadFile(accDesc.Keypath)
	if err != nil {
		return "", fmt.Errorf("failed to read keystore: %v", err)
	}

	return string(content), nil
}

// ImportKeystore imports a keystore from JSON content
func (s *WalletService) ImportKeystore(keystoreJSON string, password string, description string) (string, error) {
	// Validate keystore by trying to decrypt it
	var keystoreData map[string]interface{}
	if err := json.Unmarshal([]byte(keystoreJSON), &keystoreData); err != nil {
		return "", errors.New("invalid keystore JSON format")
	}

	// Get address from keystore
	addressRaw, ok := keystoreData["address"].(string)
	if !ok {
		return "", errors.New("keystore missing address field")
	}
	address := "0x" + addressRaw

	// Check if wallet already exists
	accs := accounts.GetAccounts()
	if _, exists := accs[address]; exists {
		return "", errors.New("wallet with this address already exists")
	}

	// Save keystore to file
	homeDir, _ := os.UserHomeDir()
	keystoreDir := fmt.Sprintf("%s/.jarvis/keystores", homeDir)
	os.MkdirAll(keystoreDir, 0755)
	keystorePath := fmt.Sprintf("%s/%s.json", keystoreDir, address)

	if err := os.WriteFile(keystorePath, []byte(keystoreJSON), 0600); err != nil {
		return "", fmt.Errorf("failed to save keystore: %v", err)
	}

	// Use default description if empty
	if description == "" {
		description = fmt.Sprintf("Imported Wallet %d", len(accs)+1)
	}

	// Store account metadata
	accDesc := types.AccDesc{
		Address: address,
		Kind:    "keystore",
		Keypath: keystorePath,
		Desc:    description,
	}
	if err := accounts.StoreAccountRecord(accDesc); err != nil {
		return "", fmt.Errorf("failed to store account record: %v", err)
	}

	// Verify by unlocking
	acc, err := accounts.UnlockKeystoreAccountWithPassword(accDesc, password)
	if err != nil {
		// Rollback: delete the keystore file
		os.Remove(keystorePath)
		homeDir, _ := os.UserHomeDir()
		os.Remove(fmt.Sprintf("%s/.jarvis/%s.json", homeDir, address))
		return "", errors.New("invalid password for keystore")
	}

	s.currentAccount = acc
	s.address = address

	return address, nil
}

// UpdateWalletDescription updates the description of a wallet
func (s *WalletService) UpdateWalletDescription(address string, description string) error {
	accs := accounts.GetAccounts()
	accDesc, exists := accs[address]
	if !exists {
		return errors.New("wallet not found")
	}

	accDesc.Desc = description
	return accounts.StoreAccountRecord(accDesc)
}

// GetStatus returns the current wallet status
func (s *WalletService) GetStatus() WalletStatus {
	return WalletStatus{
		HasWallet: s.HasWallet(),
		IsLocked:  s.currentAccount == nil,
		Address:   s.address,
		Wallets:   s.GetWallets(),
	}
}

// SignMessage signs a message with the private key
func (s *WalletService) SignMessage(message string) (string, error) {
	if s.currentAccount == nil {
		return "", errors.New("wallet is locked")
	}

	return s.currentAccount.SignMessage(message)
}

// GetTransactOpts creates a bind.TransactOpts for the current account
func (s *WalletService) GetTransactOpts(chainId *big.Int) (*bind.TransactOpts, error) {
	if s.currentAccount == nil {
		return nil, errors.New("wallet is locked")
	}

	return &bind.TransactOpts{
		From: s.currentAccount.Address(),
		Signer: func(addr common.Address, tx *ethtypes.Transaction) (*ethtypes.Transaction, error) {
			_, signedTx, err := s.currentAccount.SignTx(tx, chainId)
			return signedTx, err
		},
	}, nil
}
