package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
)

// TransactionType definitions
type TransactionType string

const (
	TxTypeDeposit TransactionType = "DEPOSIT"
	TxTypeSend    TransactionType = "SEND"
	TxTypeReceive TransactionType = "RECEIVE"
	TxTypeApprove TransactionType = "APPROVE"
	TxTypeBuy     TransactionType = "BUY"
	TxTypeSell    TransactionType = "SELL"
	TxTypeSwap    TransactionType = "SWAP"
)

// TransactionRecord represents a single transaction history entry (for API compatibility)
type TransactionRecord struct {
	Hash        string          `json:"hash"`
	Type        TransactionType `json:"type"`
	Amount      string          `json:"amount"` // Display string e.g. "100 USDT"
	To          string          `json:"to"`
	From        string          `json:"from"`
	Timestamp   int64           `json:"timestamp"`
	Status      string          `json:"status"` // "Pending", "Success", "Failed"
	Description string          `json:"description"`
	// Extended fields from database
	TokenSymbol  string `json:"tokenSymbol,omitempty"`
	TokenAddress string `json:"tokenAddress,omitempty"`
	BlockNumber  int64  `json:"blockNumber,omitempty"`
	Network      string `json:"network,omitempty"`
	GasUsed      string `json:"gasUsed,omitempty"`
	GasPrice     string `json:"gasPrice,omitempty"`
}

// HistoryService handles transaction history storage using SQLite database
type HistoryService struct {
	db *database.Service
}

// NewHistoryService creates a new history manager with database backend
func NewHistoryService(dbService *database.Service) *HistoryService {
	return &HistoryService{
		db: dbService,
	}
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
	ctx := context.Background()
	now := time.Now().UnixMilli()

	params := db.CreateWalletTransactionParams{
		TxHash:      txHash,
		TxType:      string(txType),
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
		Status:      "Success", // Assuming added after successful broadcast
		Description: sql.NullString{String: desc, Valid: desc != ""},
		Network:     "BSC", // Default to BSC, can be parameterized later
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := s.db.Queries().CreateWalletTransaction(ctx, params)
	return err
}

// AddTransactionWithDetails records a new transaction with full details
func (s *HistoryService) AddTransactionWithDetails(
	txHash string,
	txType TransactionType,
	amount string,
	from string,
	to string,
	desc string,
	tokenAddress string,
	tokenSymbol string,
	tokenDecimals int64,
	blockNumber int64,
	network string,
	chainID int64,
	gasUsed string,
	gasPrice string,
	nonce int64,
	status string,
) error {
	ctx := context.Background()
	now := time.Now().UnixMilli()

	params := db.CreateWalletTransactionParams{
		TxHash:        txHash,
		TxType:        string(txType),
		FromAddress:   from,
		ToAddress:     to,
		Amount:        amount,
		TokenAddress:  sql.NullString{String: tokenAddress, Valid: tokenAddress != ""},
		TokenSymbol:   sql.NullString{String: tokenSymbol, Valid: tokenSymbol != ""},
		TokenDecimals: sql.NullInt64{Int64: tokenDecimals, Valid: tokenDecimals > 0},
		Status:        status,
		Description:   sql.NullString{String: desc, Valid: desc != ""},
		BlockNumber:   sql.NullInt64{Int64: blockNumber, Valid: blockNumber > 0},
		Network:       network,
		ChainID:       sql.NullInt64{Int64: chainID, Valid: chainID > 0},
		GasUsed:       sql.NullString{String: gasUsed, Valid: gasUsed != ""},
		GasPrice:      sql.NullString{String: gasPrice, Valid: gasPrice != ""},
		Nonce:         sql.NullInt64{Int64: nonce, Valid: nonce >= 0},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	_, err := s.db.Queries().CreateWalletTransaction(ctx, params)
	return err
}

// GetTransactions returns the list of transactions (latest 100 by default)
func (s *HistoryService) GetTransactions() []TransactionRecord {
	ctx := context.Background()

	txs, err := s.db.Queries().ListWalletTransactions(ctx, db.ListWalletTransactionsParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		return []TransactionRecord{}
	}

	return convertDBTransactionsToRecords(txs)
}

// GetTransactionsByAddress returns transactions for a specific address
func (s *HistoryService) GetTransactionsByAddress(address string, limit int, offset int) []TransactionRecord {
	ctx := context.Background()

	txs, err := s.db.Queries().ListWalletTransactionsByAddress(ctx, db.ListWalletTransactionsByAddressParams{
		FromAddress: address,
		ToAddress:   address,
		Limit:       int64(limit),
		Offset:      int64(offset),
	})
	if err != nil {
		return []TransactionRecord{}
	}

	return convertDBTransactionsToRecords(txs)
}

// GetTransactionByHash returns a specific transaction by hash
func (s *HistoryService) GetTransactionByHash(txHash string) (*TransactionRecord, error) {
	ctx := context.Background()

	tx, err := s.db.Queries().GetWalletTransactionByHash(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	records := convertDBTransactionsToRecords([]db.WalletTransaction{tx})
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to convert transaction")
	}

	return &records[0], nil
}

// UpdateTransactionStatus updates the status of a transaction
func (s *HistoryService) UpdateTransactionStatus(txHash string, status string, blockNumber int64) error {
	ctx := context.Background()
	now := time.Now().UnixMilli()

	_, err := s.db.Queries().UpdateWalletTransactionStatus(ctx, db.UpdateWalletTransactionStatusParams{
		Status:      status,
		BlockNumber: sql.NullInt64{Int64: blockNumber, Valid: blockNumber > 0},
		UpdatedAt:   now,
		TxHash:      txHash,
	})

	return err
}

// GetPendingTransactions returns all pending transactions
func (s *HistoryService) GetPendingTransactions() []TransactionRecord {
	ctx := context.Background()

	txs, err := s.db.Queries().GetPendingWalletTransactions(ctx)
	if err != nil {
		return []TransactionRecord{}
	}

	return convertDBTransactionsToRecords(txs)
}

// ClearHistory deletes all records (use with caution)
func (s *HistoryService) ClearHistory() error {
	// Note: This would require a new query in wallet_transactions.sql
	// For now, we'll skip implementation as it's destructive
	return fmt.Errorf("ClearHistory not implemented for safety - use SQL directly if needed")
}

// Helper function to convert database records to API records
func convertDBTransactionsToRecords(dbTxs []db.WalletTransaction) []TransactionRecord {
	records := make([]TransactionRecord, 0, len(dbTxs))

	for _, tx := range dbTxs {
		record := TransactionRecord{
			Hash:      tx.TxHash,
			Type:      TransactionType(tx.TxType),
			Amount:    tx.Amount,
			To:        tx.ToAddress,
			From:      tx.FromAddress,
			Timestamp: tx.CreatedAt / 1000, // Convert milliseconds to seconds for compatibility
			Status:    tx.Status,
			Network:   tx.Network,
		}

		if tx.Description.Valid {
			record.Description = tx.Description.String
		}
		if tx.TokenSymbol.Valid {
			record.TokenSymbol = tx.TokenSymbol.String
		}
		if tx.TokenAddress.Valid {
			record.TokenAddress = tx.TokenAddress.String
		}
		if tx.BlockNumber.Valid {
			record.BlockNumber = tx.BlockNumber.Int64
		}
		if tx.GasUsed.Valid {
			record.GasUsed = tx.GasUsed.String
		}
		if tx.GasPrice.Valid {
			record.GasPrice = tx.GasPrice.String
		}

		records = append(records, record)
	}

	return records
}
