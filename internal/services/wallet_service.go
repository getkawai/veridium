package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

// WalletStatus represents the current state of the wallet
type WalletStatus struct {
	HasWallet bool   `json:"hasWallet"`
	IsLocked  bool   `json:"isLocked"`
	Address   string `json:"address"`
}

// WalletService handles local wallet management
type WalletService struct {
	keystoreDir  string
	keystorePath string
	privateKey   []byte
	address      string
}

// NewWalletService creates a new wallet service
func NewWalletService(dataDir string) *WalletService {
	keystoreDir := filepath.Join(dataDir, "keystore")
	return &WalletService{
		keystoreDir:  keystoreDir,
		keystorePath: filepath.Join(keystoreDir, "wallet.json"),
	}
}

// HasWallet checks if a wallet already exists
func (s *WalletService) HasWallet() bool {
	_, err := os.Stat(s.keystorePath)
	return err == nil
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
	address := crypto.PubkeyToAddress(masterKey.PublicKey).Hex()

	// 2. Encrypt and save
	if err := s.saveKeystore(password, privBytes, address); err != nil {
		return "", err
	}

	// 3. Set state
	s.privateKey = privBytes
	s.address = address

	return address, nil
}

// UnlockWallet decrypts the keystore and loads it into memory
func (s *WalletService) UnlockWallet(password string) (string, error) {
	data, err := os.ReadFile(s.keystorePath)
	if err != nil {
		return "", err
	}

	var keystore struct {
		EncryptedKey string `json:"encryptedKey"`
		Address      string `json:"address"`
		Salt         string `json:"salt"`
		Nonce        string `json:"nonce"`
	}

	if err := json.Unmarshal(data, &keystore); err != nil {
		return "", err
	}

	// Derive key from password
	salt, _ := hex.DecodeString(keystore.Salt)
	key := deriveKey(password, salt)

	// Decrypt
	encrypted, _ := hex.DecodeString(keystore.EncryptedKey)
	nonce, _ := hex.DecodeString(keystore.Nonce)

	privBytes, err := decrypt(encrypted, key, nonce)
	if err != nil {
		return "", errors.New("invalid password")
	}

	s.privateKey = privBytes
	s.address = keystore.Address

	return s.address, nil
}

// GetStatus returns the current wallet status
func (s *WalletService) GetStatus() WalletStatus {
	return WalletStatus{
		HasWallet: s.HasWallet(),
		IsLocked:  s.privateKey == nil,
		Address:   s.address,
	}
}

// SignMessage signs a message with the private key
func (s *WalletService) SignMessage(message string) (string, error) {
	if s.privateKey == nil {
		return "", errors.New("wallet is locked")
	}

	privKey, err := crypto.ToECDSA(s.privateKey)
	if err != nil {
		return "", err
	}

	// Hash the message
	hash := crypto.Keccak256Hash([]byte(message))
	sig, err := crypto.Sign(hash.Bytes(), privKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(sig), nil
}

// Helper: Save keystore to file
func (s *WalletService) saveKeystore(password string, privKey []byte, address string) error {
	salt := make([]byte, 16)
	rand.Read(salt)

	key := deriveKey(password, salt)
	encrypted, nonce, err := encrypt(privKey, key)
	if err != nil {
		return err
	}

	keystore := struct {
		EncryptedKey string `json:"encryptedKey"`
		Address      string `json:"address"`
		Salt         string `json:"salt"`
		Nonce        string `json:"nonce"`
	}{
		EncryptedKey: hex.EncodeToString(encrypted),
		Address:      address,
		Salt:         hex.EncodeToString(salt),
		Nonce:        hex.EncodeToString(nonce),
	}

	os.MkdirAll(s.keystoreDir, 0700)
	data, _ := json.MarshalIndent(keystore, "", "  ")
	return os.WriteFile(s.keystorePath, data, 0600)
}

// Simple key derivation (for production, use scrypt or PBKDF2)
func deriveKey(password string, salt []byte) []byte {
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	return h.Sum(nil)
}

func encrypt(data, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return ciphertext, nonce, nil
}

func decrypt(data, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}
