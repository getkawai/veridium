package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/contracts/kawaitoken"
	"github.com/kawai-network/contracts"
)

// KawaiHolder represents a KAWAI token holder
type KawaiHolder struct {
	Address common.Address
	Balance *big.Int
}

// HolderScanner scans KAWAI token holders from blockchain
type HolderScanner struct {
	client       *ethclient.Client
	tokenAddress common.Address
	kawaiToken   *kawaitoken.KawaiToken
}

// NewHolderScanner creates a new holder scanner
func NewHolderScanner() (*HolderScanner, error) {
	client, err := ethclient.Dial(contracts.MonadRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Monad: %w", err)
	}

	tokenAddr := common.HexToAddress(contracts.KawaiTokenAddress)
	kawaiToken, err := kawaitoken.NewKawaiToken(tokenAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load KawaiToken: %w", err)
	}

	return &HolderScanner{
		client:       client,
		tokenAddress: tokenAddr,
		kawaiToken:   kawaiToken,
	}, nil
}

// ScanHolders scans all KAWAI token holders from Transfer events
// This is more efficient than iterating all addresses
func (hs *HolderScanner) ScanHolders(ctx context.Context, fromBlock, toBlock *big.Int) ([]*KawaiHolder, error) {
	log.Printf("📊 [HOLDER SCANNER] Scanning KAWAI holders from block %s to %s", fromBlock.String(), toBlock.String())

	// Get all Transfer events with block range filter
	filterOpts := &bind.FilterOpts{
		Start:   fromBlock.Uint64(),
		End:     &[]uint64{toBlock.Uint64()}[0],
		Context: ctx,
	}

	transferIterator, err := hs.kawaiToken.FilterTransfer(filterOpts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter Transfer events: %w", err)
	}
	defer transferIterator.Close()

	// Collect unique addresses
	addressSet := make(map[common.Address]bool)
	for transferIterator.Next() {
		event := transferIterator.Event

		// Add both sender and receiver
		if event.From != (common.Address{}) { // Not mint
			addressSet[event.From] = true
		}
		if event.To != (common.Address{}) { // Not burn
			addressSet[event.To] = true
		}
	}

	if err := transferIterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating Transfer events: %w", err)
	}

	log.Printf("📊 [HOLDER SCANNER] Found %d unique addresses", len(addressSet))

	// Get current balance for each address
	holders := make([]*KawaiHolder, 0, len(addressSet))
	failedQueries := 0

	for addr := range addressSet {
		balance, err := hs.kawaiToken.BalanceOf(nil, addr)
		if err != nil {
			log.Printf("⚠️  [HOLDER SCANNER] Failed to get balance for %s: %v", addr.Hex(), err)
			failedQueries++
			continue
		}

		// Only include addresses with non-zero balance
		if balance.Cmp(big.NewInt(0)) > 0 {
			holders = append(holders, &KawaiHolder{
				Address: addr,
				Balance: balance,
			})
		}
	}

	// Check if too many queries failed (>10% failure rate is unacceptable)
	if failedQueries > 0 {
		failureRate := float64(failedQueries) / float64(len(addressSet)) * 100
		log.Printf("⚠️  [HOLDER SCANNER] Warning: %d balance queries failed (%.2f%%)", failedQueries, failureRate)

		if failedQueries > len(addressSet)/10 {
			return nil, fmt.Errorf("too many failed balance queries: %d/%d (%.2f%%) - data quality insufficient for settlement",
				failedQueries, len(addressSet), failureRate)
		}
	}

	log.Printf("✅ [HOLDER SCANNER] Found %d holders with non-zero balance", len(holders))
	return holders, nil
}

// ScanHoldersLatest scans holders using latest block
// NOTE: This is an ADMIN-ONLY operation used during weekly settlement via CLI tool.
// It is NOT used in the desktop app.
//
// Currently scans from configured start block to latest block.
// Start block is configured in internal/constant/blockchain.go:
// - Testnet: block 0 (scans entire history - acceptable for small blockchain)
// - Mainnet: set to token deployment block to optimize performance
//
// Alternative optimization approaches for future:
// - Option 1: Cache holder addresses and only scan new blocks since last settlement
// - Option 2: Use indexed subgraph for holder queries
// - Option 3: Maintain incremental holder list in KV store
// - Option 4: Implement chunked scanning with configurable block ranges
func (hs *HolderScanner) ScanHoldersLatest(ctx context.Context) ([]*KawaiHolder, error) {
	// Get latest block
	latestBlock, err := hs.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	// Use configured start block from constants
	startBlock := big.NewInt(contracts.HolderScanStartBlock)

	if contracts.HolderScanStartBlock > 0 {
		log.Printf("📊 [HOLDER SCANNER] Using configured start block: %d", contracts.HolderScanStartBlock)
	}

	// Scan from start block to latest
	return hs.ScanHolders(ctx, startBlock, big.NewInt(int64(latestBlock)))
}

// ScanHoldersFromBlock scans holders from a specific start block to latest
// This is used for hybrid holder scanning (registry + recent blockchain scan)
// to work around Monad testnet's 100-block RPC limit
func (hs *HolderScanner) ScanHoldersFromBlock(ctx context.Context, startBlock uint64) ([]common.Address, error) {
	// Get latest block
	latestBlock, err := hs.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	log.Printf("📊 [HOLDER SCANNER] Scanning recent transfers from block %d to %d", startBlock, latestBlock)

	// Get all Transfer events with block range filter
	filterOpts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &latestBlock,
		Context: ctx,
	}

	transferIterator, err := hs.kawaiToken.FilterTransfer(filterOpts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter Transfer events: %w", err)
	}
	defer transferIterator.Close()

	// Collect unique addresses
	addressSet := make(map[common.Address]bool)
	for transferIterator.Next() {
		event := transferIterator.Event

		// Add both sender and receiver
		if event.From != (common.Address{}) { // Not mint
			addressSet[event.From] = true
		}
		if event.To != (common.Address{}) { // Not burn
			addressSet[event.To] = true
		}
	}

	if err := transferIterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating Transfer events: %w", err)
	}

	// Convert map to slice
	addresses := make([]common.Address, 0, len(addressSet))
	for addr := range addressSet {
		addresses = append(addresses, addr)
	}

	log.Printf("📊 [HOLDER SCANNER] Found %d unique addresses in recent transfers", len(addresses))
	return addresses, nil
}

// GetBalance returns the current KAWAI balance for an address
func (hs *HolderScanner) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	return hs.kawaiToken.BalanceOf(&bind.CallOpts{Context: ctx}, address)
}

// GetTotalSupply returns the total KAWAI supply
func (hs *HolderScanner) GetTotalSupply(ctx context.Context) (*big.Int, error) {
	return hs.kawaiToken.TotalSupply(nil)
}

// GetMaxSupply returns the max KAWAI supply
func (hs *HolderScanner) GetMaxSupply(ctx context.Context) (*big.Int, error) {
	return hs.kawaiToken.Cap(nil)
}

// CalculateHolderShare calculates a holder's proportional share of total profit
func CalculateHolderShare(holderBalance, totalSupply, totalProfit *big.Int) *big.Int {
	// share = (holderBalance / totalSupply) * totalProfit
	// To avoid precision loss, we do: (holderBalance * totalProfit) / totalSupply

	// Safety check: prevent division by zero
	if totalSupply.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}

	share := new(big.Int).Mul(holderBalance, totalProfit)
	share.Div(share, totalSupply)

	return share
}

// ValidateHolders validates that holder balances sum to total supply
func ValidateHolders(holders []*KawaiHolder, totalSupply *big.Int) error {
	sum := big.NewInt(0)
	for _, holder := range holders {
		sum.Add(sum, holder.Balance)
	}

	if sum.Cmp(totalSupply) != 0 {
		return fmt.Errorf("holder balances (%s) do not match total supply (%s)", sum.String(), totalSupply.String())
	}

	return nil
}
