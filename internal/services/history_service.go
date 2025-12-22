package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// TransactionType definitions
type TransactionType string

const (
	TxTypeDeposit TransactionType = "DEPOSIT"
	TxTypeSend    TransactionType = "SEND"
	TxTypeReceive TransactionType = "RECEIVE" // For future use if we have an indexer
	TxTypeApprove TransactionType = "APPROVE"
	TxTypeBuy     TransactionType = "BUY"
	TxTypeSell    TransactionType = "SELL"
)

// TransactionRecord represents a single transaction history entry
type TransactionRecord struct {
	Hash        string          `json:"hash"`
	Type        TransactionType `json:"type"`
	Amount      string          `json:"amount"` // Display string e.g. "100 USDT"
	To          string          `json:"to"`
	From        string          `json:"from"`
	Timestamp   int64           `json:"timestamp"`
	Status      string          `json:"status"` // "Pending", "Success", "Failed"
	Description string          `json:"description"`
}

// HistoryService handles local transaction history storage
type HistoryService struct {
	mu           sync.RWMutex
	transactions []TransactionRecord
	filePath     string
}

// NewHistoryService creates a new history manager
func NewHistoryService() *HistoryService {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".jarvis", "transactions.json")

	s := &HistoryService{
		filePath: path,
	}
	s.load()
	return s
}

// Ensure loaded from disk
func (s *HistoryService) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Dir(s.filePath)
	os.MkdirAll(dir, 0755)

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		s.transactions = []TransactionRecord{}
		return
	}

	json.Unmarshal(data, &s.transactions)
}

// save to disk
func (s *HistoryService) save() error {
	data, err := json.MarshalIndent(s.transactions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

// AddTransaction records a new transaction
func (s *HistoryService) AddTransaction(
	txHash string,
	txType TransactionType,
	amount string,
	from string,
	to string,
	desc string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec := TransactionRecord{
		Hash:        txHash,
		Type:        txType,
		Amount:      amount,
		From:        from,
		To:          to,
		Timestamp:   time.Now().Unix(),
		Status:      "Success", // Assuming added after successful broadcast for now
		Description: desc,
	}

	// Prepend
	s.transactions = append([]TransactionRecord{rec}, s.transactions...)

	// Keep only last 100
	if len(s.transactions) > 100 {
		s.transactions = s.transactions[:100]
	}

	return s.save()
}

// GetTransactions returns the list of transactions
func (s *HistoryService) GetTransactions() []TransactionRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return copy to be safe
	result := make([]TransactionRecord, len(s.transactions))
	copy(result, s.transactions)

	// Sort by timestamp desc (just in case)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp > result[j].Timestamp
	})

	return result
}

// ClearHistory deletes all records
func (s *HistoryService) ClearHistory() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions = []TransactionRecord{}
	return s.save()
}
