package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// MarketplaceError represents different types of marketplace errors
// Requirements: 7.5 - Comprehensive error types and messages
type MarketplaceError struct {
	Type    MarketplaceErrorType   `json:"type"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *MarketplaceError) Error() string {
	return fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)
}

// MarketplaceErrorType represents the category of marketplace error
type MarketplaceErrorType string

const (
	ErrorTypeValidation    MarketplaceErrorType = "VALIDATION"
	ErrorTypeBalance       MarketplaceErrorType = "BALANCE"
	ErrorTypeOrder         MarketplaceErrorType = "ORDER"
	ErrorTypeAuthorization MarketplaceErrorType = "AUTHORIZATION"
	ErrorTypeBlockchain    MarketplaceErrorType = "BLOCKCHAIN"
	ErrorTypeStorage       MarketplaceErrorType = "STORAGE"
	ErrorTypeNetwork       MarketplaceErrorType = "NETWORK"
	ErrorTypeInternal      MarketplaceErrorType = "INTERNAL"
)

// Error constructors for consistent error handling
// Requirements: 7.5 - Proper error propagation to frontend

func NewValidationError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeValidation,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewBalanceError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeBalance,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewOrderError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeOrder,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewAuthorizationError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeAuthorization,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewBlockchainError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeBlockchain,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewStorageError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeStorage,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewNetworkError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeNetwork,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewInternalError(code, message string, details map[string]interface{}) *MarketplaceError {
	return &MarketplaceError{
		Type:    ErrorTypeInternal,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Helper function to wrap standard errors as MarketplaceError
func WrapError(err error, errorType MarketplaceErrorType, code, message string) *MarketplaceError {
	details := map[string]interface{}{
		"original_error": err.Error(),
	}
	return &MarketplaceError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Logging helper for marketplace operations
// Requirements: 7.5 - Add logging for debugging and audit purposes
func logMarketplaceOperation(operation, userAddress, details string, err error) {
	if err != nil {
		log.Printf("❌ [Marketplace:%s] User: %s, Details: %s, Error: %v", operation, userAddress, details, err)
	} else {
		log.Printf("✅ [Marketplace:%s] User: %s, Details: %s", operation, userAddress, details)
	}
}

func logMarketplaceInfo(operation, message string) {
	log.Printf("ℹ️  [Marketplace:%s] %s", operation, message)
}

func logMarketplaceWarning(operation, message string) {
	log.Printf("⚠️  [Marketplace:%s] %s", operation, message)
}

// MarketplaceService provides P2P marketplace operations for KAWAI token trading.
// This service enables Contributors to sell their earned KAWAI tokens to Investors
// using USDT through secure smart contract escrow and atomic swaps.
//
// The service integrates with:
// - Cloudflare KV for off-chain data storage and caching (dedicated P2P marketplace namespace)
// - OTCMarket smart contract for atomic swaps and escrow (task 3.1)
// - WalletService for transaction signing and authorization
//
// Key features:
// - Order creation and management (tasks 2.1-2.4)
// - Secure trade execution with atomic swaps (tasks 5.1-5.6)
// - Real-time market data and analytics (tasks 6.1-6.4)
// - Order history and tracking (tasks 9.1-9.4)
// - Price discovery through supply and demand (task 5.1)
type MarketplaceService struct {
	kvStore           *store.KVStore
	blockchainClient  *blockchain.Client
	walletService     *WalletService
	orderService      *OrderService
	tradeService      *TradeService
	marketDataService *MarketDataService
	eventListener     *MarketplaceEventListener
	app               *application.App // Wails v3 application for event emission
}

// Order represents a sell order in the marketplace
type Order struct {
	ID              string    `json:"id"`
	Seller          string    `json:"seller"`
	TokenAmount     string    `json:"tokenAmount"`   // KAWAI token amount (in wei)
	USDTPrice       string    `json:"usdtPrice"`     // Total USDT price (in wei)
	PricePerToken   string    `json:"pricePerToken"` // USDT price per KAWAI token
	Status          string    `json:"status"`        // active, filled, cancelled
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	TxHash          string    `json:"txHash"`          // Transaction hash for order creation
	RemainingAmount string    `json:"remainingAmount"` // Remaining tokens for partial fills
}

// TradeResult represents the result of a trade execution
type TradeResult struct {
	Success     bool   `json:"success"`
	TxHash      string `json:"txHash"`
	OrderID     string `json:"orderID"`
	TokenAmount string `json:"tokenAmount"` // Amount of KAWAI tokens traded
	USDTAmount  string `json:"usdtAmount"`  // Amount of USDT paid
	Buyer       string `json:"buyer"`       // Buyer's wallet address
	Seller      string `json:"seller"`      // Seller's wallet address
	Error       string `json:"error,omitempty"`
}

// MarketStats represents current marketplace statistics
type MarketStats struct {
	LowestAskPrice    string  `json:"lowestAskPrice"`    // Lowest ask price in USDT per token
	HighestBidPrice   string  `json:"highestBidPrice"`   // Highest recent trade price
	Volume24h         string  `json:"volume24h"`         // 24-hour trading volume in USDT
	PriceChange24h    string  `json:"priceChange24h"`    // 24-hour price change percentage
	ActiveOrdersCount int     `json:"activeOrdersCount"` // Number of active orders
	RecentTrades      []Trade `json:"recentTrades"`      // Recent completed trades
}

// Trade represents a completed trade transaction
type Trade struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"orderID"`
	Buyer       string    `json:"buyer"`
	Seller      string    `json:"seller"`
	TokenAmount string    `json:"tokenAmount"`
	USDTAmount  string    `json:"usdtAmount"`
	Price       string    `json:"price"` // Price per token
	Timestamp   time.Time `json:"timestamp"`
	TxHash      string    `json:"txHash"`
}

// OrderHistory represents a user's order history with execution details
// Requirements: 4.1, 4.2 - Order history with status and timestamps
type OrderHistory struct {
	Orders []OrderHistoryEntry `json:"orders"`
	Trades []TradeHistoryEntry `json:"trades"`
	Total  int                 `json:"total"`
}

// OrderHistoryEntry represents an order in the user's history
// Requirements: 4.1, 4.4, 4.5 - Order details with status and timestamps
type OrderHistoryEntry struct {
	Order
	StatusHistory []OrderStatusChange `json:"statusHistory"`
	TradeCount    int                 `json:"tradeCount"`   // Number of trades for this order
	FilledAmount  string              `json:"filledAmount"` // Total amount filled across all trades
}

// TradeHistoryEntry represents a trade in the user's history
// Requirements: 4.2 - Trade history with execution details
type TradeHistoryEntry struct {
	Trade
	OrderDetails OrderSummary `json:"orderDetails"` // Summary of the original order
}

// OrderStatusChange represents a status change event for an order
// Requirements: 4.5 - Order status tracking with timestamps
type OrderStatusChange struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	TxHash    string    `json:"txHash,omitempty"` // Transaction hash if status change was on-chain
}

// OrderSummary represents a summary of an order for trade history
type OrderSummary struct {
	ID            string `json:"id"`
	TokenAmount   string `json:"tokenAmount"`
	USDTPrice     string `json:"usdtPrice"`
	PricePerToken string `json:"pricePerToken"`
}

// Real-time Event Types
// Requirements: 4.3, 4.5 - Real-time status updates and order detail metadata

// MarketplaceRealtimeEvent represents a real-time marketplace event
type MarketplaceRealtimeEvent struct {
	Type      string      `json:"type"`      // Event type: "order_status_changed", "order_created", "trade_completed"
	Timestamp time.Time   `json:"timestamp"` // Event timestamp
	Data      interface{} `json:"data"`      // Event-specific data
	UserID    string      `json:"userId"`    // User ID for targeted events
}

// OrderStatusUpdateEvent represents an order status change event
type OrderStatusUpdateEvent struct {
	OrderID     string    `json:"orderID"`
	OldStatus   string    `json:"oldStatus"`
	NewStatus   string    `json:"newStatus"`
	Timestamp   time.Time `json:"timestamp"`
	TxHash      string    `json:"txHash,omitempty"`
	OrderDetail Order     `json:"orderDetail"` // Full order details for metadata display
}

// OrderCreatedEvent represents a new order creation event
type OrderCreatedEvent struct {
	Order     Order     `json:"order"`
	Timestamp time.Time `json:"timestamp"`
}

// TradeCompletedEvent represents a completed trade event
type TradeCompletedEvent struct {
	Trade       Trade     `json:"trade"`
	OrderDetail Order     `json:"orderDetail"`
	Timestamp   time.Time `json:"timestamp"`
}

// OrderResult represents the result of order creation
type OrderResult struct {
	Success bool   `json:"success"`
	OrderID string `json:"orderID,omitempty"`
	TxHash  string `json:"txHash,omitempty"`
	Error   string `json:"error,omitempty"`
}

// MarketDepth represents order book depth information
type MarketDepth struct {
	Asks []PriceLevel `json:"asks"` // Sell orders by price level
	Bids []PriceLevel `json:"bids"` // Buy orders by price level (future feature)
}

// PriceLevel represents aggregated orders at a specific price
type PriceLevel struct {
	Price  string `json:"price"`  // Price per token
	Amount string `json:"amount"` // Total amount at this price
	Count  int    `json:"count"`  // Number of orders at this price
}

// PricePoint represents a price data point for trend analysis
type PricePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Price     string    `json:"price"`
	Volume    string    `json:"volume"`
}

// MarketDataService provides market analytics and pricing information
// Requirements: 5.1, 5.2, 5.3, 5.5
type MarketDataService struct {
	kvStore      *store.KVStore
	orderService *OrderService
}

// NewMarketDataService creates a new MarketDataService instance
func NewMarketDataService(kvStore *store.KVStore, orderService *OrderService) *MarketDataService {
	return &MarketDataService{
		kvStore:      kvStore,
		orderService: orderService,
	}
}

// CalculateMarketStats calculates current market statistics
// Requirements: 5.1, 5.2
func (s *MarketDataService) CalculateMarketStats() (*MarketStats, error) {
	// Get all active orders for analysis
	activeOrders, err := s.orderService.GetActiveOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to get active orders: %w", err)
	}

	// Calculate lowest ask price (Requirement 5.1)
	lowestAskPrice := s.calculateLowestAskPrice(activeOrders)

	// Get recent trades for highest bid price and volume calculation
	recentTrades, err := s.getRecentTrades(24 * time.Hour)
	if err != nil {
		log.Printf("⚠️  Failed to get recent trades: %v", err)
		recentTrades = []Trade{} // Continue with empty trades
	}

	// Calculate highest recent trade price (Requirement 5.1)
	highestBidPrice := s.calculateHighestRecentPrice(recentTrades)

	// Calculate 24-hour volume and price change (Requirement 5.2)
	volume24h, priceChange24h := s.calculate24HourMetrics(recentTrades)

	stats := &MarketStats{
		LowestAskPrice:    lowestAskPrice,
		HighestBidPrice:   highestBidPrice,
		Volume24h:         volume24h,
		PriceChange24h:    priceChange24h,
		ActiveOrdersCount: len(activeOrders),
		RecentTrades:      recentTrades,
	}

	log.Printf("✅ Calculated market stats: %d active orders, lowest ask: %s, highest bid: %s",
		len(activeOrders), lowestAskPrice, highestBidPrice)
	return stats, nil
}

// GetPriceTrends returns price trend data for the specified timeframe
// Requirements: 5.3
func (s *MarketDataService) GetPriceTrends(timeframe string) ([]PricePoint, error) {
	// Parse timeframe duration
	duration, err := s.parseTimeframe(timeframe)
	if err != nil {
		return nil, fmt.Errorf("invalid timeframe: %w", err)
	}

	// Get trades within the timeframe
	trades, err := s.getRecentTrades(duration)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades for trends: %w", err)
	}

	// Convert trades to price points
	pricePoints := s.convertTradesToPricePoints(trades)

	log.Printf("✅ Generated %d price points for timeframe: %s", len(pricePoints), timeframe)
	return pricePoints, nil
}

// GetMarketDepth calculates market depth for order book visualization
// Requirements: 5.5
func (s *MarketDataService) GetMarketDepth() (*MarketDepth, error) {
	// Get all active orders
	activeOrders, err := s.orderService.GetActiveOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to get active orders for market depth: %w", err)
	}

	// Group orders by price level
	priceLevels := s.groupOrdersByPriceLevel(activeOrders)

	depth := &MarketDepth{
		Asks: priceLevels,
		Bids: []PriceLevel{}, // Future feature - buy orders not implemented yet
	}

	log.Printf("✅ Calculated market depth with %d price levels", len(priceLevels))
	return depth, nil
}

// Helper methods for market statistics calculation

// calculateLowestAskPrice finds the lowest ask price among active orders
func (s *MarketDataService) calculateLowestAskPrice(orders []Order) string {
	if len(orders) == 0 {
		return "0"
	}

	var lowestPrice *big.Int
	for _, order := range orders {
		if order.Status != "active" {
			continue
		}

		pricePerToken, ok := new(big.Int).SetString(order.PricePerToken, 10)
		if !ok {
			continue
		}

		if lowestPrice == nil || pricePerToken.Cmp(lowestPrice) < 0 {
			lowestPrice = pricePerToken
		}
	}

	if lowestPrice == nil {
		return "0"
	}
	return lowestPrice.String()
}

// calculateHighestRecentPrice finds the highest price from recent trades
func (s *MarketDataService) calculateHighestRecentPrice(trades []Trade) string {
	if len(trades) == 0 {
		return "0"
	}

	var highestPrice *big.Int
	for _, trade := range trades {
		price, ok := new(big.Int).SetString(trade.Price, 10)
		if !ok {
			continue
		}

		if highestPrice == nil || price.Cmp(highestPrice) > 0 {
			highestPrice = price
		}
	}

	if highestPrice == nil {
		return "0"
	}
	return highestPrice.String()
}

// calculate24HourMetrics calculates 24-hour volume and price change
func (s *MarketDataService) calculate24HourMetrics(trades []Trade) (volume24h, priceChange24h string) {
	if len(trades) == 0 {
		return "0", "0"
	}

	// Calculate total volume
	totalVolume := big.NewInt(0)
	var firstPrice, lastPrice *big.Int

	// Sort trades by timestamp (oldest first)
	sortedTrades := make([]Trade, len(trades))
	copy(sortedTrades, trades)

	// Simple bubble sort by timestamp
	for i := 0; i < len(sortedTrades)-1; i++ {
		for j := 0; j < len(sortedTrades)-i-1; j++ {
			if sortedTrades[j].Timestamp.After(sortedTrades[j+1].Timestamp) {
				sortedTrades[j], sortedTrades[j+1] = sortedTrades[j+1], sortedTrades[j]
			}
		}
	}

	for i, trade := range sortedTrades {
		// Add to volume
		usdtAmount, ok := new(big.Int).SetString(trade.USDTAmount, 10)
		if ok {
			totalVolume.Add(totalVolume, usdtAmount)
		}

		// Track first and last prices for change calculation
		price, ok := new(big.Int).SetString(trade.Price, 10)
		if ok {
			if i == 0 {
				firstPrice = price
			}
			if i == len(sortedTrades)-1 {
				lastPrice = price
			}
		}
	}

	volume24h = totalVolume.String()

	// Calculate price change percentage
	if firstPrice != nil && lastPrice != nil && firstPrice.Cmp(big.NewInt(0)) > 0 {
		// Price change = ((lastPrice - firstPrice) / firstPrice) * 100
		change := new(big.Int).Sub(lastPrice, firstPrice)
		change.Mul(change, big.NewInt(100))
		change.Div(change, firstPrice)
		priceChange24h = change.String()
	} else {
		priceChange24h = "0"
	}

	return volume24h, priceChange24h
}

// groupOrdersByPriceLevel groups orders by price level for market depth
func (s *MarketDataService) groupOrdersByPriceLevel(orders []Order) []PriceLevel {
	priceMap := make(map[string]*PriceLevel)

	for _, order := range orders {
		if order.Status != "active" {
			continue
		}

		price := order.PricePerToken
		remainingAmount, ok := new(big.Int).SetString(order.RemainingAmount, 10)
		if !ok {
			continue
		}

		if level, exists := priceMap[price]; exists {
			// Add to existing price level
			currentAmount, ok := new(big.Int).SetString(level.Amount, 10)
			if ok {
				currentAmount.Add(currentAmount, remainingAmount)
				level.Amount = currentAmount.String()
				level.Count++
			}
		} else {
			// Create new price level
			priceMap[price] = &PriceLevel{
				Price:  price,
				Amount: remainingAmount.String(),
				Count:  1,
			}
		}
	}

	// Convert map to sorted slice
	priceLevels := make([]PriceLevel, 0, len(priceMap))
	for _, level := range priceMap {
		priceLevels = append(priceLevels, *level)
	}

	// Sort by price (ascending)
	for i := 0; i < len(priceLevels)-1; i++ {
		for j := 0; j < len(priceLevels)-i-1; j++ {
			price1, ok1 := new(big.Int).SetString(priceLevels[j].Price, 10)
			price2, ok2 := new(big.Int).SetString(priceLevels[j+1].Price, 10)
			if ok1 && ok2 && price1.Cmp(price2) > 0 {
				priceLevels[j], priceLevels[j+1] = priceLevels[j+1], priceLevels[j]
			}
		}
	}

	return priceLevels
}

// convertTradesToPricePoints converts trade records to price points for trend analysis
func (s *MarketDataService) convertTradesToPricePoints(trades []Trade) []PricePoint {
	pricePoints := make([]PricePoint, len(trades))

	for i, trade := range trades {
		pricePoints[i] = PricePoint{
			Timestamp: trade.Timestamp,
			Price:     trade.Price,
			Volume:    trade.USDTAmount,
		}
	}

	return pricePoints
}

// parseTimeframe parses timeframe string to duration
func (s *MarketDataService) parseTimeframe(timeframe string) (time.Duration, error) {
	switch timeframe {
	case "1h":
		return time.Hour, nil
	case "24h":
		return 24 * time.Hour, nil
	case "7d":
		return 7 * 24 * time.Hour, nil
	case "30d":
		return 30 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported timeframe: %s", timeframe)
	}
}

// getRecentTrades retrieves recent trades within the specified duration
func (s *MarketDataService) getRecentTrades(duration time.Duration) ([]Trade, error) {
	// This is a simplified implementation for MVP
	// For now, return empty trades - this will be populated by the event listener
	// when trades are actually executed and stored

	trades := []Trade{}
	log.Printf("⚠️  getRecentTrades: MVP implementation - returning empty trades for duration %v", duration)

	return trades, nil
}

// TradeService manages trade execution and atomic swaps
type TradeService struct {
	blockchainClient   *blockchain.Client
	orderService       *OrderService
	walletService      *WalletService
	kvStore            *store.KVStore
	marketplaceService *MarketplaceService // For real-time event emission
}

// NewTradeService creates a new TradeService instance
func NewTradeService(blockchainClient *blockchain.Client, orderService *OrderService, walletService *WalletService, kvStore *store.KVStore) *TradeService {
	return &TradeService{
		blockchainClient:   blockchainClient,
		orderService:       orderService,
		walletService:      walletService,
		kvStore:            kvStore,
		marketplaceService: nil, // Will be set later to avoid circular dependency
	}
}

// SetMarketplaceService sets the marketplace service for real-time event emission
// Requirements: 4.3 - Real-time status updates
func (s *TradeService) SetMarketplaceService(marketplaceService *MarketplaceService) {
	s.marketplaceService = marketplaceService
}

// ValidateTradeConditions validates all conditions required for trade execution
// Requirements: 3.1 - Buyer balance validation
func (s *TradeService) ValidateTradeConditions(orderID string, buyer string) error {
	// Get the order to validate
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Check if order is still active
	if order.Status != "active" {
		return fmt.Errorf("order is not active: status is %s", order.Status)
	}

	// Check if buyer is not the seller (can't buy your own order)
	if order.Seller == buyer {
		return fmt.Errorf("cannot buy your own order")
	}

	// Parse USDT price for balance validation
	usdtPrice, ok := new(big.Int).SetString(order.USDTPrice, 10)
	if !ok {
		return fmt.Errorf("invalid USDT price format in order")
	}

	// Requirement 3.1: Validate buyer has sufficient USDT balance
	if s.blockchainClient != nil {
		ctx := context.Background()
		buyerAddr := common.HexToAddress(buyer)

		if err := s.blockchainClient.ValidateTradeBalance(ctx, buyerAddr, usdtPrice); err != nil {
			return fmt.Errorf("buyer balance validation failed: %w", err)
		}

		log.Printf("✅ Validated buyer %s has sufficient USDT balance of %s", buyer, usdtPrice.String())
	} else {
		log.Printf("⚠️  Blockchain client not available - skipping balance validation for buyer %s", buyer)
	}

	return nil
}

// ExecuteTrade executes an atomic swap trade for the specified order
// Requirements: 3.1, 3.2, 3.3, 3.4
// Supports both full and partial order execution
func (s *TradeService) ExecuteTrade(orderID string, buyer string) (*TradeResult, error) {
	return s.ExecutePartialTrade(orderID, buyer, "")
}

// ExecutePartialTrade executes a trade for a specific amount (empty amount means full order)
// Requirements: 2.4, 3.5 - Partial order handling
func (s *TradeService) ExecutePartialTrade(orderID string, buyer string, requestedAmount string) (*TradeResult, error) {
	// Validate trade conditions first
	if err := s.ValidateTradeConditions(orderID, buyer); err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("trade validation failed: %v", err),
		}, nil
	}

	// Get the order
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("order not found: %v", err),
		}, nil
	}

	// Parse remaining amount in order
	remainingTokenAmount, ok := new(big.Int).SetString(order.RemainingAmount, 10)
	if !ok {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   "invalid remaining token amount format in order",
		}, nil
	}

	// Determine trade amount (full or partial)
	var tradeTokenAmount *big.Int
	if requestedAmount == "" {
		// Full order execution
		tradeTokenAmount = new(big.Int).Set(remainingTokenAmount)
	} else {
		// Partial order execution
		tradeTokenAmount, ok = new(big.Int).SetString(requestedAmount, 10)
		if !ok {
			return &TradeResult{
				Success: false,
				OrderID: orderID,
				Error:   "invalid requested amount format",
			}, nil
		}

		// Validate requested amount doesn't exceed remaining amount
		if tradeTokenAmount.Cmp(remainingTokenAmount) > 0 {
			return &TradeResult{
				Success: false,
				OrderID: orderID,
				Error:   fmt.Sprintf("requested amount %s exceeds remaining amount %s", tradeTokenAmount.String(), remainingTokenAmount.String()),
			}, nil
		}

		// Validate minimum trade amount (must be > 0)
		if tradeTokenAmount.Cmp(big.NewInt(0)) <= 0 {
			return &TradeResult{
				Success: false,
				OrderID: orderID,
				Error:   "trade amount must be greater than zero",
			}, nil
		}
	}

	// Calculate proportional USDT amount for the trade
	originalTokenAmount, ok := new(big.Int).SetString(order.TokenAmount, 10)
	if !ok {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   "invalid original token amount format in order",
		}, nil
	}

	originalUSDTPrice, ok := new(big.Int).SetString(order.USDTPrice, 10)
	if !ok {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   "invalid USDT price format in order",
		}, nil
	}

	// Calculate USDT amount: (tradeTokenAmount * originalUSDTPrice) / originalTokenAmount
	tradeUSDTAmount := new(big.Int).Mul(tradeTokenAmount, originalUSDTPrice)
	tradeUSDTAmount = tradeUSDTAmount.Div(tradeUSDTAmount, originalTokenAmount)

	// Execute atomic swap on blockchain if client is available
	var txHash string
	if s.blockchainClient != nil && s.walletService != nil {
		ctx := context.Background()

		// Get transaction options from wallet service
		transactOpts, err := s.walletService.getTransactOpts(s.blockchainClient.ChainID)
		if err != nil {
			return &TradeResult{
				Success: false,
				OrderID: orderID,
				Error:   fmt.Sprintf("failed to get transaction options: %v", err),
			}, nil
		}

		// Parse order ID to big.Int for contract call
		contractOrderID, ok := new(big.Int).SetString(orderID, 10)
		if !ok {
			return &TradeResult{
				Success: false,
				OrderID: orderID,
				Error:   "invalid order ID format for contract call",
			}, nil
		}

		// Requirements 3.2, 3.3: Execute atomic swap via smart contract
		// The smart contract ensures atomicity - either both transfers succeed or both fail
		tx, err := s.blockchainClient.MarketplaceBuyOrder(ctx, transactOpts, contractOrderID)
		if err != nil {
			return &TradeResult{
				Success: false,
				OrderID: orderID,
				Error:   fmt.Sprintf("atomic swap execution failed: %v", err),
			}, nil
		}

		txHash = tx.Hash().Hex()
		log.Printf("✅ Executed atomic swap for order %s with tx %s", orderID, txHash)
	} else {
		log.Printf("⚠️  Blockchain client or wallet service not available - simulating trade execution for order %s", orderID)
		txHash = "0x" + orderID // Simulate transaction hash for testing
	}

	// Process trade completion with partial amount handling
	if err := s.ProcessPartialTradeCompletion(orderID, txHash, tradeTokenAmount, tradeUSDTAmount, buyer); err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("trade completion processing failed: %v", err),
		}, nil
	}

	// Requirement 3.4: Return successful trade result (trade completion event is handled by event listener)
	return &TradeResult{
		Success:     true,
		TxHash:      txHash,
		OrderID:     orderID,
		TokenAmount: tradeTokenAmount.String(),
		USDTAmount:  tradeUSDTAmount.String(),
		Buyer:       buyer,
		Seller:      order.Seller,
	}, nil
}

// ProcessTradeCompletion processes the completion of a trade and updates order status
// Requirements: 3.4, 3.5
func (s *TradeService) ProcessTradeCompletion(orderID, txHash string, buyerAddress string) error {
	// Get the order
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// For backward compatibility, assume full order execution
	remainingAmount, ok := new(big.Int).SetString(order.RemainingAmount, 10)
	if !ok {
		return fmt.Errorf("invalid remaining amount format in order")
	}

	return s.ProcessPartialTradeCompletion(orderID, txHash, remainingAmount, nil, buyerAddress)
}

// ProcessPartialTradeCompletion processes the completion of a partial trade
// Requirements: 2.4, 3.5 - Partial order handling and status management
func (s *TradeService) ProcessPartialTradeCompletion(orderID, txHash string, tradeTokenAmount, tradeUSDTAmount *big.Int, buyerAddress string) error {
	// Get the order
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Parse current remaining amount
	currentRemainingAmount, ok := new(big.Int).SetString(order.RemainingAmount, 10)
	if !ok {
		return fmt.Errorf("invalid remaining amount format in order")
	}

	// Calculate new remaining amount after this trade
	newRemainingAmount := new(big.Int).Sub(currentRemainingAmount, tradeTokenAmount)

	// Validate the trade amount doesn't exceed remaining amount
	if newRemainingAmount.Cmp(big.NewInt(0)) < 0 {
		return fmt.Errorf("trade amount %s exceeds remaining amount %s", tradeTokenAmount.String(), currentRemainingAmount.String())
	}

	// Update order based on remaining amount
	if newRemainingAmount.Cmp(big.NewInt(0)) == 0 {
		// Requirement 2.4: Order is fully filled, mark as completed
		if err := s.orderService.UpdateOrderStatus(orderID, "filled"); err != nil {
			return fmt.Errorf("failed to update order status to filled: %w", err)
		}
		log.Printf("✅ Order %s fully filled and marked as completed", orderID)
	} else {
		// Requirement 3.5: Partial fill - update remaining amount and keep active
		if err := s.updateOrderRemainingAmount(orderID, newRemainingAmount); err != nil {
			return fmt.Errorf("failed to update remaining amount: %w", err)
		}
		log.Printf("✅ Order %s partially filled: %s tokens remaining", orderID, newRemainingAmount.String())
	}

	// Store trade record
	if err := s.storePartialTradeRecord(order, txHash, tradeTokenAmount, tradeUSDTAmount, buyerAddress); err != nil {
		log.Printf("⚠️  Failed to store trade record: %v", err)
		// Don't fail the trade completion for this
	} else {
		// Emit real-time trade completion event if trade was stored successfully
		// Requirements: 4.3 - Event-driven status updates for active orders
		if s.marketplaceService != nil {
			// Get the stored trade record to emit with complete information
			// For now, we'll create a temporary trade object for the event
			trade := &Trade{
				ID:          fmt.Sprintf("temp_%s_%d", orderID, time.Now().Unix()),
				OrderID:     orderID,
				Buyer:       buyerAddress, // Now filled synchronously
				Seller:      order.Seller,
				TokenAmount: tradeTokenAmount.String(),
				USDTAmount:  tradeUSDTAmount.String(),
				Price:       order.PricePerToken,
				Timestamp:   time.Now(),
				TxHash:      txHash,
			}
			s.marketplaceService.emitTradeCompleted(trade, order)
		}
	}

	log.Printf("✅ Processed partial trade completion for order %s", orderID)
	return nil
}

// updateOrderRemainingAmount updates the remaining amount for a partially filled order
func (s *TradeService) updateOrderRemainingAmount(orderID string, newRemainingAmount *big.Int) error {
	// Get existing order
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order for remaining amount update: %w", err)
	}

	// Update remaining amount and timestamp
	order.RemainingAmount = newRemainingAmount.String()
	order.UpdatedAt = time.Now()

	// Store updated order
	if err := s.orderService.StoreOrder(order); err != nil {
		return fmt.Errorf("failed to store updated order: %w", err)
	}

	// Ensure real-time updates for partial fill amount changes
	// Requirements: 2.4 - Real-time remaining amount updates for partial fills
	if s.marketplaceService != nil {
		if err := s.marketplaceService.ensureRealTimeOrderUpdates(orderID); err != nil {
			log.Printf("⚠️  Failed to ensure real-time updates for order %s: %v", orderID, err)
		}
	}

	return nil
}

// storeTradeRecord stores a completed trade record in KV store (backward compatibility)
func (s *TradeService) storeTradeRecord(order *Order, txHash string, buyerAddress string) error {
	// For backward compatibility, assume full order execution
	tokenAmount, ok := new(big.Int).SetString(order.RemainingAmount, 10)
	if !ok {
		return fmt.Errorf("invalid remaining amount format")
	}

	usdtAmount, ok := new(big.Int).SetString(order.USDTPrice, 10)
	if !ok {
		return fmt.Errorf("invalid USDT price format")
	}

	return s.storePartialTradeRecord(order, txHash, tokenAmount, usdtAmount, buyerAddress)
}

// storePartialTradeRecord stores a completed trade record with specific amounts
func (s *TradeService) storePartialTradeRecord(order *Order, txHash string, tradeTokenAmount, tradeUSDTAmount *big.Int, buyerAddress string) error {
	ctx := context.Background()

	// Generate trade ID
	tradeID, err := s.generateTradeID()
	if err != nil {
		return fmt.Errorf("failed to generate trade ID: %w", err)
	}

	// If USDT amount is nil, calculate it proportionally
	if tradeUSDTAmount == nil {
		originalTokenAmount, ok := new(big.Int).SetString(order.TokenAmount, 10)
		if !ok {
			return fmt.Errorf("invalid original token amount format")
		}

		originalUSDTPrice, ok := new(big.Int).SetString(order.USDTPrice, 10)
		if !ok {
			return fmt.Errorf("invalid USDT price format")
		}

		// Calculate proportional USDT amount
		tradeUSDTAmount = new(big.Int).Mul(tradeTokenAmount, originalUSDTPrice)
		tradeUSDTAmount = tradeUSDTAmount.Div(tradeUSDTAmount, originalTokenAmount)
	}

	// Create trade record with buyer information
	trade := Trade{
		ID:          tradeID,
		OrderID:     order.ID,
		Buyer:       buyerAddress, // Now filled synchronously at trade execution
		Seller:      order.Seller,
		TokenAmount: tradeTokenAmount.String(),
		USDTAmount:  tradeUSDTAmount.String(),
		Price:       order.PricePerToken,
		Timestamp:   time.Now(),
		TxHash:      txHash,
	}

	// Serialize trade to JSON
	tradeJSON, err := json.Marshal(trade)
	if err != nil {
		return fmt.Errorf("failed to marshal trade: %w", err)
	}

	// Store trade record
	key := s.getTradeKey(tradeID)
	if err := s.kvStore.StoreMarketplaceData(ctx, key, tradeJSON); err != nil {
		return fmt.Errorf("failed to store trade record: %w", err)
	}

	// Add trade to order's trade history
	if err := s.addTradeToOrderHistory(order.ID, tradeID); err != nil {
		log.Printf("⚠️  Failed to add trade to order history: %v", err)
	}

	// Add trade to seller's trade history
	if err := s.addTradeToUserHistory(order.Seller, "seller", tradeID); err != nil {
		log.Printf("⚠️  Failed to add trade to seller history: %v", err)
	}

	// Add trade to buyer's trade history (synchronously)
	if err := s.addTradeToUserHistory(buyerAddress, "buyer", tradeID); err != nil {
		log.Printf("⚠️  Failed to add trade to buyer history: %v", err)
	}

	log.Printf("✅ Stored trade record %s for order %s (Buyer: %s, Seller: %s)", tradeID, order.ID, buyerAddress, order.Seller)
	return nil
}

// addTradeToOrderHistory adds a trade to an order's trade history
func (s *TradeService) addTradeToOrderHistory(orderID, tradeID string) error {
	ctx := context.Background()
	key := s.getOrderTradesKey(orderID)

	// Get existing trade IDs
	var tradeIDs []string
	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err == nil {
		if err := json.Unmarshal(value, &tradeIDs); err != nil {
			return fmt.Errorf("failed to unmarshal existing trade IDs: %w", err)
		}
	}

	// Add new trade ID
	tradeIDs = append(tradeIDs, tradeID)

	// Store updated trade IDs
	tradeIDsJSON, err := json.Marshal(tradeIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal trade IDs: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, key, tradeIDsJSON); err != nil {
		return fmt.Errorf("failed to store trade IDs: %w", err)
	}

	return nil
}

// addTradeToUserHistory adds a trade to a user's trade history
func (s *TradeService) addTradeToUserHistory(walletAddress, role, tradeID string) error {
	ctx := context.Background()
	key := s.getUserTradesKey(walletAddress, role)

	// Get existing trade IDs
	var tradeIDs []string
	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err == nil {
		if err := json.Unmarshal(value, &tradeIDs); err != nil {
			return fmt.Errorf("failed to unmarshal existing trade IDs: %w", err)
		}
	}

	// Check if trade ID already exists (Idempotency)
	for _, id := range tradeIDs {
		if id == tradeID {
			return nil // Already recorded
		}
	}

	// Add new trade ID
	tradeIDs = append(tradeIDs, tradeID)

	// Store updated trade IDs
	tradeIDsJSON, err := json.Marshal(tradeIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal trade IDs: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, key, tradeIDsJSON); err != nil {
		return fmt.Errorf("failed to store trade IDs: %w", err)
	}

	return nil
}

// Key generation methods for trade history
func (s *TradeService) getOrderTradesKey(orderID string) string {
	return fmt.Sprintf("order:%s:trades", orderID)
}

func (s *TradeService) getUserTradesKey(walletAddress, role string) string {
	return fmt.Sprintf("user:%s:trades:%s", walletAddress, role)
}

// generateTradeID generates a unique trade ID
func (s *TradeService) generateTradeID() (string, error) {
	// Generate 16 random bytes and encode as hex
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return "trade_" + hex.EncodeToString(bytes), nil
}

// getTradeKey generates a KV store key for a trade record
func (s *TradeService) getTradeKey(tradeID string) string {
	return fmt.Sprintf("trade:%s", tradeID)
}

// OrderService handles order lifecycle management and validation
type OrderService struct {
	kvStore            *store.KVStore
	blockchainClient   *blockchain.Client
	walletService      *WalletService
	marketplaceService *MarketplaceService // For real-time event emission
}

// NewOrderService creates a new OrderService instance
func NewOrderService(kvStore *store.KVStore, blockchainClient *blockchain.Client, walletService *WalletService) *OrderService {
	return &OrderService{
		kvStore:          kvStore,
		blockchainClient: blockchainClient,
		walletService:    walletService,
		// marketDataService will be set later to avoid circular dependency
	}
}

// SetMarketplaceService sets the marketplace service for real-time event emission
// Requirements: 4.3 - Real-time status updates
func (s *OrderService) SetMarketplaceService(marketplaceService *MarketplaceService) {
	s.marketplaceService = marketplaceService
}

// ValidateOrderCreation validates order creation parameters
// ValidateOrderCreation validates order creation parameters
// Requirements: 1.1, 1.2
func (s *OrderService) ValidateOrderCreation(seller string, tokenAmount, usdtPrice *big.Int) error {
	// Requirement 1.2: Validate required fields
	if seller == "" {
		return NewValidationError(
			"EMPTY_SELLER_ADDRESS",
			"Seller wallet address is required",
			map[string]interface{}{
				"provided_seller": seller,
			},
		)
	}
	if tokenAmount == nil || tokenAmount.Cmp(big.NewInt(0)) <= 0 {
		return NewValidationError(
			"INVALID_TOKEN_AMOUNT",
			"Token amount must be greater than zero",
			map[string]interface{}{
				"provided_amount": tokenAmount,
			},
		)
	}
	if usdtPrice == nil || usdtPrice.Cmp(big.NewInt(0)) <= 0 {
		return NewValidationError(
			"INVALID_USDT_PRICE",
			"USDT price must be greater than zero",
			map[string]interface{}{
				"provided_price": usdtPrice,
			},
		)
	}

	// Requirement 1.1: Validate seller has sufficient KAWAI token balance
	if s.blockchainClient != nil {
		ctx := context.Background()
		sellerAddr := common.HexToAddress(seller)

		if err := s.blockchainClient.ValidateOrderCreationBalance(ctx, sellerAddr, tokenAmount); err != nil {
			return NewBalanceError(
				"INSUFFICIENT_KAWAI_BALANCE",
				"Seller has insufficient KAWAI token balance",
				map[string]interface{}{
					"seller":           seller,
					"required_amount":  tokenAmount.String(),
					"blockchain_error": err.Error(),
				},
			)
		}

		logMarketplaceInfo("ValidateOrderCreation", fmt.Sprintf("Validated seller %s has sufficient KAWAI balance of %s", seller, tokenAmount.String()))
	} else {
		logMarketplaceWarning("ValidateOrderCreation", fmt.Sprintf("Blockchain client not available - skipping balance validation for seller %s", seller))
	}

	return nil
}

// CreateOrder creates a new sell order and stores it in KV store
// Requirements: 1.1, 1.2, 1.4
func (s *OrderService) CreateOrder(seller string, tokenAmount, usdtPrice *big.Int) (*Order, error) {
	// Validate order creation parameters
	if err := s.ValidateOrderCreation(seller, tokenAmount, usdtPrice); err != nil {
		return nil, err // Already wrapped as MarketplaceError
	}

	// Generate unique order ID
	orderID, err := s.generateOrderID()
	if err != nil {
		return nil, NewInternalError(
			"ORDER_ID_GENERATION_FAILED",
			"Failed to generate unique order ID",
			map[string]interface{}{
				"seller":         seller,
				"token_amount":   tokenAmount.String(),
				"usdt_price":     usdtPrice.String(),
				"original_error": err.Error(),
			},
		)
	}

	// Calculate price per token
	pricePerToken := new(big.Int).Div(usdtPrice, tokenAmount)

	// Create order object
	order := &Order{
		ID:              orderID,
		Seller:          seller,
		TokenAmount:     tokenAmount.String(),
		USDTPrice:       usdtPrice.String(),
		PricePerToken:   pricePerToken.String(),
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		RemainingAmount: tokenAmount.String(), // Initially, all tokens are remaining
	}

	// Store order in KV store first
	if err := s.StoreOrder(order); err != nil {
		return nil, WrapError(err, ErrorTypeStorage, "ORDER_STORAGE_FAILED", "Failed to store order in KV store")
	}

	// Add to active orders index
	if err := s.addToActiveOrdersIndex(orderID); err != nil {
		// If adding to index fails, try to clean up the stored order
		s.DeleteOrder(orderID)
		return nil, WrapError(err, ErrorTypeStorage, "ACTIVE_INDEX_UPDATE_FAILED", "Failed to add order to active orders index")
	}

	// Add to user orders index
	if err := s.addToUserOrdersIndex(seller, orderID); err != nil {
		// If adding to user index fails, clean up
		s.removeFromActiveOrdersIndex(orderID)
		s.DeleteOrder(orderID)
		return nil, WrapError(err, ErrorTypeStorage, "USER_INDEX_UPDATE_FAILED", "Failed to add order to user orders index")
	}

	// Create order on blockchain if client is available
	if s.blockchainClient != nil && s.walletService != nil {
		ctx := context.Background()

		// Get transaction options from wallet service
		transactOpts, err := s.walletService.getTransactOpts(s.blockchainClient.ChainID)
		if err != nil {
			// Clean up local storage if blockchain transaction fails
			s.removeFromUserOrdersIndex(seller, orderID)
			s.removeFromActiveOrdersIndex(orderID)
			s.DeleteOrder(orderID)
			return nil, NewBlockchainError(
				"TRANSACTION_OPTIONS_FAILED",
				"Failed to get transaction options from wallet service",
				map[string]interface{}{
					"seller":         seller,
					"order_id":       orderID,
					"chain_id":       s.blockchainClient.ChainID.String(),
					"original_error": err.Error(),
				},
			)
		}

		// Create order on smart contract
		tx, err := s.blockchainClient.MarketplaceCreateOrder(ctx, transactOpts, tokenAmount, usdtPrice)
		if err != nil {
			// Clean up local storage if blockchain transaction fails
			s.removeFromUserOrdersIndex(seller, orderID)
			s.removeFromActiveOrdersIndex(orderID)
			s.DeleteOrder(orderID)
			return nil, NewBlockchainError(
				"CONTRACT_ORDER_CREATION_FAILED",
				"Failed to create order on smart contract",
				map[string]interface{}{
					"seller":         seller,
					"order_id":       orderID,
					"token_amount":   tokenAmount.String(),
					"usdt_price":     usdtPrice.String(),
					"original_error": err.Error(),
				},
			)
		}

		// Update order with transaction hash
		order.TxHash = tx.Hash().Hex()
		order.UpdatedAt = time.Now()

		if err := s.StoreOrder(order); err != nil {
			logMarketplaceWarning("CreateOrder", fmt.Sprintf("Failed to update order %s with transaction hash: %v", orderID, err))
		}

		logMarketplaceInfo("CreateOrder", fmt.Sprintf("Created order %s on blockchain with tx %s", orderID, tx.Hash().Hex()))
	} else {
		logMarketplaceWarning("CreateOrder", fmt.Sprintf("Blockchain client or wallet service not available - order %s created locally only", orderID))
	}

	// Emit real-time order creation event
	// Requirements: 4.3 - Event-driven status updates for active orders
	if s.marketplaceService != nil {
		s.marketplaceService.emitOrderCreated(order)
	}

	logMarketplaceOperation("CreateOrder", seller, fmt.Sprintf("orderID=%s, tokenAmount=%s, usdtPrice=%s, txHash=%s", orderID, tokenAmount.String(), usdtPrice.String(), order.TxHash), nil)
	return order, nil
}

// StoreOrder stores an order in the KV store
func (s *OrderService) StoreOrder(order *Order) error {
	ctx := context.Background()

	// Serialize order to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	// Store in KV using marketplace methods
	key := s.getOrderKey(order.ID)
	if err := s.kvStore.StoreMarketplaceData(ctx, key, orderJSON); err != nil {
		log.Printf("❌ [Marketplace:KV] Failed to store order %s: %v", order.ID, err)
		return fmt.Errorf("failed to store order: %w", err)
	}

	// Update user's order index
	if err := s.addOrderToUserIndex(order.ID, order.Seller); err != nil {
		log.Printf("⚠️  Failed to update user index for order %s: %v", order.ID, err)
		// Don't fail the whole operation if index update fails
	}

	return nil
}

// addOrderToUserIndex adds an order ID to the user's order index
func (s *OrderService) addOrderToUserIndex(orderID, userAddr string) error {
	ctx := context.Background()
	userOrdersKey := fmt.Sprintf("user:%s:orders", strings.ToLower(userAddr))

	// Get existing order IDs
	var orderIDs []string
	value, err := s.kvStore.GetMarketplaceData(ctx, userOrdersKey)
	if err == nil {
		if err := json.Unmarshal(value, &orderIDs); err != nil {
			return fmt.Errorf("failed to unmarshal user orders index: %w", err)
		}
	}

	// Check if order already exists (idempotency)
	for _, id := range orderIDs {
		if id == orderID {
			return nil // Already indexed
		}
	}

	// Add new order ID
	orderIDs = append(orderIDs, orderID)

	// Save updated index
	data, err := json.Marshal(orderIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal user orders index: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, userOrdersKey, data); err != nil {
		return fmt.Errorf("failed to store user orders index: %w", err)
	}

	return nil
}

// GetOrder retrieves an order by ID from KV store
func (s *OrderService) GetOrder(orderID string) (*Order, error) {
	ctx := context.Background()

	key := s.getOrderKey(orderID)
	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from KV: %w", err)
	}

	var order Order
	if err := json.Unmarshal(value, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return &order, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(orderID, status string) error {
	// Get existing order
	order, err := s.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order for status update: %w", err)
	}

	// Update status and timestamp
	oldStatus := order.Status
	order.Status = status
	order.UpdatedAt = time.Now()

	// Store updated order
	if err := s.StoreOrder(order); err != nil {
		return fmt.Errorf("failed to store updated order: %w", err)
	}

	// Add status change to history if status actually changed
	if oldStatus != status {
		if err := s.addOrderStatusChange(orderID, status, ""); err != nil {
			log.Printf("⚠️  Failed to add status change to history for order %s: %v", orderID, err)
		}

		// Emit real-time status update event
		// Requirements: 4.3 - Event-driven status updates for active orders
		if s.marketplaceService != nil {
			s.marketplaceService.emitOrderStatusUpdate(orderID, oldStatus, status, "")

			// Handle order status change with proper cleanup and real-time updates
			// Requirements: 2.5 - Remove completed/cancelled orders from active listings
			if err := s.marketplaceService.handleOrderStatusChange(orderID, oldStatus, status); err != nil {
				log.Printf("⚠️  Failed to handle order status change for order %s: %v", orderID, err)
			}
		}
	}

	// If order is no longer active, remove from active orders index
	if status != "active" {
		if err := s.removeFromActiveOrdersIndex(orderID); err != nil {
			log.Printf("⚠️  Failed to remove order %s from active index: %v", orderID, err)
		}
	}

	log.Printf("✅ Updated order %s status from %s to %s", orderID, oldStatus, status)
	return nil
}

// addOrderStatusChange adds a status change to an order's history
// This method is used by OrderService to track status changes
func (s *OrderService) addOrderStatusChange(orderID, status, txHash string) error {
	ctx := context.Background()
	key := s.getOrderStatusHistoryKey(orderID)

	// Get existing status history
	var statusHistory []OrderStatusChange
	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err == nil {
		if err := json.Unmarshal(value, &statusHistory); err != nil {
			log.Printf("⚠️  Failed to unmarshal existing status history: %v", err)
			statusHistory = []OrderStatusChange{} // Start fresh if unmarshal fails
		}
	}

	// Add new status change
	statusChange := OrderStatusChange{
		Status:    status,
		Timestamp: time.Now(),
		TxHash:    txHash,
	}

	statusHistory = append(statusHistory, statusChange)

	// Store updated status history
	statusHistoryJSON, err := json.Marshal(statusHistory)
	if err != nil {
		return fmt.Errorf("failed to marshal status history: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, key, statusHistoryJSON); err != nil {
		return fmt.Errorf("failed to store status history: %w", err)
	}

	return nil
}

// getOrderStatusHistoryKey generates a KV store key for order status history
func (s *OrderService) getOrderStatusHistoryKey(orderID string) string {
	return fmt.Sprintf("order:%s:status_history", orderID)
}

// GetOrdersByStatus returns all orders with a specific status
func (s *OrderService) GetOrdersByStatus(status string) ([]Order, error) {
	if status == "active" {
		return s.GetActiveOrders()
	}

	// For non-active orders, we'd need to scan all orders
	// This is a simplified implementation - in production, we'd want better indexing
	return nil, fmt.Errorf("getting orders by status '%s' not implemented yet", status)
}

// GetActiveOrders returns all active orders from the active orders index
func (s *OrderService) GetActiveOrders() ([]Order, error) {
	ctx := context.Background()

	// Get active orders index
	activeOrdersKey := s.getActiveOrdersKey()
	value, err := s.kvStore.GetMarketplaceData(ctx, activeOrdersKey)
	if err != nil {
		// If index doesn't exist, return empty list
		return []Order{}, nil
	}

	// Parse order IDs from index
	var orderIDs []string
	if err := json.Unmarshal(value, &orderIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal active orders index: %w", err)
	}

	// Fetch each order
	orders := make([]Order, 0, len(orderIDs))
	for _, orderID := range orderIDs {
		order, err := s.GetOrder(orderID)
		if err != nil {
			log.Printf("⚠️  Failed to get order %s from active index: %v", orderID, err)
			continue
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

// GetUserOrders returns all orders (active and inactive) for a specific user
func (s *OrderService) GetUserOrders(userAddr string) ([]Order, error) {
	ctx := context.Background()

	// Get user's order IDs from user index
	userOrdersKey := fmt.Sprintf("user:%s:orders", strings.ToLower(userAddr))
	value, err := s.kvStore.GetMarketplaceData(ctx, userOrdersKey)
	if err != nil {
		// If index doesn't exist, return empty list
		return []Order{}, nil
	}

	// Parse order IDs from index
	var orderIDs []string
	if err := json.Unmarshal(value, &orderIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user orders index: %w", err)
	}

	// Fetch each order
	orders := make([]Order, 0, len(orderIDs))
	for _, orderID := range orderIDs {
		order, err := s.GetOrder(orderID)
		if err != nil {
			log.Printf("⚠️  Failed to get order %s from user index: %v", orderID, err)
			continue
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

// DeleteOrder removes an order from KV store
func (s *OrderService) DeleteOrder(orderID string) error {
	ctx := context.Background()

	key := s.getOrderKey(orderID)
	if err := s.kvStore.DeleteMarketplaceData(ctx, key); err != nil {
		return fmt.Errorf("failed to delete order from KV: %w", err)
	}

	return nil
}

// CancelOrder cancels an existing order
// Requirements: 1.5, 8.3
func (s *OrderService) CancelOrder(orderID string, requester string) error {
	// Get existing order
	order, err := s.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Requirement 8.3: Verify the requester is the original order creator
	if order.Seller != requester {
		return fmt.Errorf("unauthorized: only the order creator can cancel the order")
	}

	// Check if order is still active
	if order.Status != "active" {
		return fmt.Errorf("cannot cancel order: order status is %s", order.Status)
	}

	// Cancel order on blockchain if client is available
	if s.blockchainClient != nil && s.walletService != nil {
		ctx := context.Background()

		// Get transaction options from wallet service
		transactOpts, err := s.walletService.getTransactOpts(s.blockchainClient.ChainID)
		if err != nil {
			return fmt.Errorf("failed to get transaction options: %w", err)
		}

		// Parse order ID to big.Int for contract call
		contractOrderID, ok := new(big.Int).SetString(orderID, 10)
		if !ok {
			return fmt.Errorf("invalid order ID format for contract call")
		}

		// Cancel order on smart contract
		tx, err := s.blockchainClient.MarketplaceCancelOrder(ctx, transactOpts, contractOrderID)
		if err != nil {
			return fmt.Errorf("failed to cancel order on blockchain: %w", err)
		}

		log.Printf("✅ Cancelled order %s on blockchain with tx %s", orderID, tx.Hash().Hex())
	} else {
		log.Printf("⚠️  Blockchain client or wallet service not available - cancelling order %s locally only", orderID)
	}

	// Update order status to cancelled
	if err := s.UpdateOrderStatus(orderID, "cancelled"); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	log.Printf("✅ Cancelled order %s", orderID)
	return nil
}

// Helper methods

// generateOrderID generates a unique order ID
func (s *OrderService) generateOrderID() (string, error) {
	// Generate 16 random bytes and encode as hex
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// addToActiveOrdersIndex adds an order ID to the active orders index
func (s *OrderService) addToActiveOrdersIndex(orderID string) error {
	ctx := context.Background()

	activeOrdersKey := s.getActiveOrdersKey()

	// Get existing index
	var orderIDs []string
	value, err := s.kvStore.GetMarketplaceData(ctx, activeOrdersKey)
	if err == nil {
		// Index exists, parse it
		if err := json.Unmarshal(value, &orderIDs); err != nil {
			return fmt.Errorf("failed to unmarshal active orders index: %w", err)
		}
	}
	// If error getting index, start with empty list

	// Add new order ID
	orderIDs = append(orderIDs, orderID)

	// Store updated index
	indexJSON, err := json.Marshal(orderIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal active orders index: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, activeOrdersKey, indexJSON); err != nil {
		log.Printf("❌ [Marketplace:KV] Failed to write active orders index for order %s: %v", orderID, err)
		return fmt.Errorf("failed to write active orders index: %w", err)
	}

	return nil
}

// removeFromActiveOrdersIndex removes an order ID from the active orders index
func (s *OrderService) removeFromActiveOrdersIndex(orderID string) error {
	ctx := context.Background()

	activeOrdersKey := s.getActiveOrdersKey()

	// Get existing index
	value, err := s.kvStore.GetMarketplaceData(ctx, activeOrdersKey)
	if err != nil {
		// Index doesn't exist, nothing to remove
		return nil
	}

	var orderIDs []string
	if err := json.Unmarshal(value, &orderIDs); err != nil {
		return fmt.Errorf("failed to unmarshal active orders index: %w", err)
	}

	// Remove order ID from list
	filteredIDs := make([]string, 0, len(orderIDs))
	for _, id := range orderIDs {
		if id != orderID {
			filteredIDs = append(filteredIDs, id)
		}
	}

	// Store updated index
	indexJSON, err := json.Marshal(filteredIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal active orders index: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, activeOrdersKey, indexJSON); err != nil {
		log.Printf("❌ [Marketplace:KV] Failed to remove order %s from active orders index: %v", orderID, err)
		return fmt.Errorf("failed to write active orders index: %w", err)
	}

	return nil
}

// addToUserOrdersIndex adds an order ID to a user's orders index
func (s *OrderService) addToUserOrdersIndex(userAddress, orderID string) error {
	ctx := context.Background()

	userOrdersKey := s.getUserOrdersKey(userAddress)

	// Get existing index
	var orderIDs []string
	value, err := s.kvStore.GetMarketplaceData(ctx, userOrdersKey)
	if err == nil {
		// Index exists, parse it
		if err := json.Unmarshal(value, &orderIDs); err != nil {
			return fmt.Errorf("failed to unmarshal user orders index: %w", err)
		}
	}
	// If error getting index, start with empty list

	// Add new order ID
	orderIDs = append(orderIDs, orderID)

	// Store updated index
	indexJSON, err := json.Marshal(orderIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal user orders index: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, userOrdersKey, indexJSON); err != nil {
		log.Printf("❌ [Marketplace:KV] Failed to write user orders index for user %s, order %s: %v", userAddress, orderID, err)
		return fmt.Errorf("failed to write user orders index: %w", err)
	}

	return nil
}

// removeFromUserOrdersIndex removes an order ID from a user's orders index
func (s *OrderService) removeFromUserOrdersIndex(userAddress, orderID string) error {
	ctx := context.Background()

	userOrdersKey := s.getUserOrdersKey(userAddress)

	// Get existing index
	value, err := s.kvStore.GetMarketplaceData(ctx, userOrdersKey)
	if err != nil {
		// Index doesn't exist, nothing to remove
		return nil
	}

	var orderIDs []string
	if err := json.Unmarshal(value, &orderIDs); err != nil {
		return fmt.Errorf("failed to unmarshal user orders index: %w", err)
	}

	// Remove order ID from list
	filteredIDs := make([]string, 0, len(orderIDs))
	for _, id := range orderIDs {
		if id != orderID {
			filteredIDs = append(filteredIDs, id)
		}
	}

	// Store updated index
	indexJSON, err := json.Marshal(filteredIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal user orders index: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, userOrdersKey, indexJSON); err != nil {
		log.Printf("❌ [Marketplace:KV] Failed to remove order %s from user %s index: %v", orderID, userAddress, err)
		return fmt.Errorf("failed to write user orders index: %w", err)
	}

	return nil
}

// Key generation methods
func (s *OrderService) getOrderKey(orderID string) string {
	return fmt.Sprintf("order:%s", orderID)
}

func (s *OrderService) getActiveOrdersKey() string {
	return "orders:active"
}

func (s *OrderService) getUserOrdersKey(walletAddress string) string {
	return fmt.Sprintf("orders:user:%s", walletAddress)
}

// NewMarketplaceService creates a new marketplace service instance
func NewMarketplaceService(kvStore *store.KVStore, blockchainClient *blockchain.Client, walletService *WalletService) *MarketplaceService {
	// Create OrderService
	orderService := NewOrderService(kvStore, blockchainClient, walletService)

	// Create TradeService
	tradeService := NewTradeService(blockchainClient, orderService, walletService, kvStore)

	// Create MarketDataService
	marketDataService := NewMarketDataService(kvStore, orderService)

	// Create event listener
	var eventListener *MarketplaceEventListener
	if blockchainClient != nil {
		eventListener = NewMarketplaceEventListener(blockchainClient, orderService, walletService, kvStore)
	}

	service := &MarketplaceService{
		kvStore:           kvStore,
		blockchainClient:  blockchainClient, // Can be nil initially, will be set up in task 3.1
		walletService:     walletService,
		orderService:      orderService,
		tradeService:      tradeService,
		marketDataService: marketDataService,
		eventListener:     eventListener,
		app:               nil, // Will be set by SetApp method
	}

	// Set marketplace service reference in sub-services for event emission
	orderService.SetMarketplaceService(service)
	tradeService.SetMarketplaceService(service)

	if blockchainClient == nil {
		log.Printf("⚠️  MarketplaceService: Initialized without blockchain client (will be added in task 3.1)")
	} else {
		log.Printf("✅ MarketplaceService: Initialized with blockchain client")

		// Start event listener if blockchain client is available
		if eventListener != nil {
			ctx := context.Background()
			if err := eventListener.Start(ctx); err != nil {
				log.Printf("⚠️  Failed to start event listener: %v", err)
			}
		}
	}

	log.Printf("✅ MarketplaceService: Initialized with dedicated P2P marketplace namespace")
	return service
}

// SetApp sets the Wails application for event emission
// Requirements: 4.3 - Real-time updates through Wails events
func (s *MarketplaceService) SetApp(app *application.App) {
	s.app = app
	log.Printf("✅ MarketplaceService: Wails application set for real-time events")
}

// Real-time Event Emission Methods
// Requirements: 4.3, 4.5 - Event-driven status updates and WebSocket-like updates

// emitOrderStatusUpdate emits a real-time order status update event
func (s *MarketplaceService) emitOrderStatusUpdate(orderID, oldStatus, newStatus, txHash string) {
	if s.app == nil {
		log.Printf("⚠️  Cannot emit order status update: Wails application not set")
		return
	}

	// Get full order details for metadata display
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		log.Printf("⚠️  Failed to get order details for status update event: %v", err)
		return
	}

	// Create event data
	event := MarketplaceRealtimeEvent{
		Type:      "order_status_changed",
		Timestamp: time.Now(),
		UserID:    order.Seller, // Target the order owner
		Data: OrderStatusUpdateEvent{
			OrderID:     orderID,
			OldStatus:   oldStatus,
			NewStatus:   newStatus,
			Timestamp:   time.Now(),
			TxHash:      txHash,
			OrderDetail: *order,
		},
	}

	// Emit to all connected clients for order book updates
	s.app.Event.Emit("marketplace:order_status_changed", event)

	// Emit targeted event to order owner
	s.app.Event.Emit(fmt.Sprintf("marketplace:user:%s:order_status_changed", order.Seller), event)

	log.Printf("📡 Emitted order status update: %s %s -> %s", orderID, oldStatus, newStatus)
}

// emitOrderCreated emits a real-time order creation event
func (s *MarketplaceService) emitOrderCreated(order *Order) {
	if s.app == nil {
		log.Printf("⚠️  Cannot emit order created: Wails application not set")
		return
	}

	// Create event data
	event := MarketplaceRealtimeEvent{
		Type:      "order_created",
		Timestamp: time.Now(),
		UserID:    order.Seller,
		Data: OrderCreatedEvent{
			Order:     *order,
			Timestamp: time.Now(),
		},
	}

	// Emit to all connected clients (for order book updates)
	s.app.Event.Emit("marketplace:order_created", event)

	// Emit targeted event to order creator
	s.app.Event.Emit(fmt.Sprintf("marketplace:user:%s:order_created", order.Seller), event)

	log.Printf("📡 Emitted order created: %s", order.ID)
}

// emitTradeCompleted emits a real-time trade completion event
func (s *MarketplaceService) emitTradeCompleted(trade *Trade, order *Order) {
	if s.app == nil {
		log.Printf("⚠️  Cannot emit trade completed: Wails application not set")
		return
	}

	// Create event data
	event := MarketplaceRealtimeEvent{
		Type:      "trade_completed",
		Timestamp: time.Now(),
		Data: TradeCompletedEvent{
			Trade:       *trade,
			OrderDetail: *order,
			Timestamp:   time.Now(),
		},
	}

	// Emit to all connected clients (for market data updates)
	s.app.Event.Emit("marketplace:trade_completed", event)

	// Emit targeted events to both buyer and seller
	if trade.Buyer != "" {
		s.app.Event.Emit(fmt.Sprintf("marketplace:user:%s:trade_completed", trade.Buyer), event)
	}
	s.app.Event.Emit(fmt.Sprintf("marketplace:user:%s:trade_completed", trade.Seller), event)

	log.Printf("📡 Emitted trade completed: %s", trade.ID)
}

// emitOrderPartiallyFilled emits an order partially filled event
func (s *MarketplaceService) emitOrderPartiallyFilled(order *Order, amountFilled, buyer string) {
	if s.app == nil {
		log.Printf("⚠️  Cannot emit order partially filled: Wails application not set")
		return
	}

	// Create event data
	event := map[string]interface{}{
		"orderID":         order.ID,
		"amountFilled":    amountFilled,
		"remainingAmount": order.RemainingAmount,
		"buyer":           buyer,
		"seller":          order.Seller,
		"timestamp":       time.Now(),
	}

	// Emit to all connected clients (for order book updates)
	s.app.Event.Emit("marketplace:order_partially_filled", event)

	// Emit targeted events to both buyer and seller
	if buyer != "" {
		s.app.Event.Emit(fmt.Sprintf("marketplace:user:%s:order_partially_filled", buyer), event)
	}
	s.app.Event.Emit(fmt.Sprintf("marketplace:user:%s:order_partially_filled", order.Seller), event)

	log.Printf("📡 Emitted order partially filled: %s (%s filled, %s remaining)", 
		order.ID, amountFilled, order.RemainingAmount)
}

// emitMarketDataUpdate emits a market data update event
func (s *MarketplaceService) emitMarketDataUpdate() {
	if s.app == nil {
		log.Printf("⚠️  Cannot emit market data update: Wails application not set")
		return
	}

	// Get current market stats
	stats, err := s.marketDataService.CalculateMarketStats()
	if err != nil {
		log.Printf("⚠️  Failed to get market stats for update event: %v", err)
		return
	}

	// Create event data
	event := MarketplaceRealtimeEvent{
		Type:      "market_data_updated",
		Timestamp: time.Now(),
		Data:      stats,
	}

	// Emit to all connected clients
	s.app.Event.Emit("marketplace:market_data_updated", event)

	log.Printf("📡 Emitted market data update")
}

// GetActiveOrders returns all active sell orders with optional sorting and filtering
// This is a Wails-exposed method for the frontend
// Requirements: 7.1, 7.2 - Validate input data and return structured data with filtering and sorting
func (s *MarketplaceService) GetActiveOrders(sortBy string, filterBy map[string]interface{}) ([]Order, error) {
	// Requirement 7.1: Validate input data
	if err := s.validateGetActiveOrdersInput(sortBy, filterBy); err != nil {
		logMarketplaceOperation("GetActiveOrders", "anonymous", fmt.Sprintf("sortBy=%s, filters=%d", sortBy, len(filterBy)), err)
		return nil, err // Already wrapped as MarketplaceError
	}

	// Get all active orders from OrderService
	orders, err := s.orderService.GetActiveOrders()
	if err != nil {
		wrappedErr := WrapError(err, ErrorTypeStorage, "ORDER_RETRIEVAL_FAILED", "Failed to retrieve active orders from storage")
		logMarketplaceOperation("GetActiveOrders", "anonymous", "retrieving from storage", wrappedErr)
		return nil, wrappedErr
	}

	// Validate order information completeness for each order (Requirements: 2.2)
	validOrders := make([]Order, 0, len(orders))
	for _, order := range orders {
		if err := s.validateOrderInformationCompleteness(order); err != nil {
			logMarketplaceWarning("GetActiveOrders", fmt.Sprintf("Order %s failed completeness validation: %v", order.ID, err))
			continue // Skip orders with incomplete information
		}
		validOrders = append(validOrders, order)
	}
	orders = validOrders

	// Apply filtering if specified
	if len(filterBy) > 0 {
		orders = s.applyOrderFilters(orders, filterBy)
		logMarketplaceInfo("GetActiveOrders", fmt.Sprintf("Applied %d filters, %d orders remaining", len(filterBy), len(orders)))
	}

	// Apply sorting if specified
	if sortBy != "" {
		orders = s.sortOrders(orders, sortBy)
		logMarketplaceInfo("GetActiveOrders", fmt.Sprintf("Applied sorting: %s", sortBy))
	}

	logMarketplaceOperation("GetActiveOrders", "anonymous", fmt.Sprintf("Retrieved %d orders with sorting '%s' and %d filters", len(orders), sortBy, len(filterBy)), nil)
	return orders, nil
}

// CreateSellOrder creates a new sell order for KAWAI tokens
// This is a Wails-exposed method for the frontend
// Requirements: 8.1, 8.2 - Use connected wallet address and ensure only wallet owner can create orders
func (s *MarketplaceService) CreateSellOrder(tokenAmount, usdtPrice string) (*OrderResult, error) {
	// Requirement 8.1: Use connected wallet address for marketplace operations
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		err := NewAuthorizationError(
			"NO_WALLET_CONNECTED",
			"No wallet connected",
			map[string]interface{}{
				"operation": "CreateSellOrder",
			},
		)
		logMarketplaceOperation("CreateSellOrder", "none", fmt.Sprintf("tokenAmount=%s, usdtPrice=%s", tokenAmount, usdtPrice), err)
		return &OrderResult{
			Success: false,
			Error:   err.Message,
		}, nil
	}

	// Requirement 8.1: Validate wallet address format
	if err := s.validateWalletAddress(currentAddress); err != nil {
		logMarketplaceOperation("CreateSellOrder", currentAddress, fmt.Sprintf("tokenAmount=%s, usdtPrice=%s", tokenAmount, usdtPrice), err)
		return &OrderResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Requirement 7.1: Validate input data
	if err := s.validateCreateSellOrderInput(tokenAmount, usdtPrice); err != nil {
		logMarketplaceOperation("CreateSellOrder", currentAddress, fmt.Sprintf("tokenAmount=%s, usdtPrice=%s", tokenAmount, usdtPrice), err)
		return &OrderResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Parse token amount and USDT price
	tokenAmountBig, ok := new(big.Int).SetString(tokenAmount, 10)
	if !ok {
		err := NewValidationError(
			"INVALID_TOKEN_AMOUNT_PARSE",
			"Failed to parse token amount",
			map[string]interface{}{
				"provided_amount": tokenAmount,
			},
		)
		logMarketplaceOperation("CreateSellOrder", currentAddress, fmt.Sprintf("tokenAmount=%s, usdtPrice=%s", tokenAmount, usdtPrice), err)
		return &OrderResult{
			Success: false,
			Error:   err.Message,
		}, nil
	}

	usdtPriceBig, ok := new(big.Int).SetString(usdtPrice, 10)
	if !ok {
		err := NewValidationError(
			"INVALID_USDT_PRICE_PARSE",
			"Failed to parse USDT price",
			map[string]interface{}{
				"provided_price": usdtPrice,
			},
		)
		logMarketplaceOperation("CreateSellOrder", currentAddress, fmt.Sprintf("tokenAmount=%s, usdtPrice=%s", tokenAmount, usdtPrice), err)
		return &OrderResult{
			Success: false,
			Error:   err.Message,
		}, nil
	}

	// Requirement 8.2: Ensure only connected wallet owner can create orders for their tokens
	// This is enforced by using currentAddress from the wallet service
	order, err := s.orderService.CreateOrder(currentAddress, tokenAmountBig, usdtPriceBig)
	if err != nil {
		// Wrap the error appropriately based on its type
		var wrappedErr *MarketplaceError
		if marketplaceErr, ok := err.(*MarketplaceError); ok {
			wrappedErr = marketplaceErr
		} else {
			wrappedErr = WrapError(err, ErrorTypeInternal, "ORDER_CREATION_FAILED", "Failed to create sell order")
		}

		logMarketplaceOperation("CreateSellOrder", currentAddress, fmt.Sprintf("tokenAmount=%s, usdtPrice=%s", tokenAmount, usdtPrice), wrappedErr)
		return &OrderResult{
			Success: false,
			Error:   wrappedErr.Message,
		}, nil
	}

	logMarketplaceOperation("CreateSellOrder", currentAddress, fmt.Sprintf("orderID=%s, tokenAmount=%s, usdtPrice=%s, txHash=%s", order.ID, tokenAmount, usdtPrice, order.TxHash), nil)
	return &OrderResult{
		Success: true,
		OrderID: order.ID,
		TxHash:  order.TxHash,
	}, nil
}

// BuyOrder executes a buy order for the specified order ID
// This is a Wails-exposed method for the frontend
// Requirements: 8.1, 8.5 - Use connected wallet address and handle wallet operations
func (s *MarketplaceService) BuyOrder(orderID string) (*TradeResult, error) {
	// Requirement 8.1: Use connected wallet address for marketplace operations
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   "no wallet connected",
		}, nil
	}

	// Requirement 8.1: Validate wallet address format
	if err := s.validateWalletAddress(currentAddress); err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("invalid wallet address: %v", err),
		}, nil
	}

	// Validate order ID format
	if err := s.validateOrderID(orderID); err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("invalid order ID: %v", err),
		}, nil
	}

	// Requirement 8.5: Execute trade using TradeService (handles wallet operations)
	result, err := s.tradeService.ExecuteTrade(orderID, currentAddress)
	if err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("failed to execute trade: %v", err),
		}, nil
	}

	log.Printf("✅ Executed buy order %s for wallet %s", orderID, currentAddress)
	return result, nil
}

// BuyPartialOrder executes a partial buy order for the specified order ID and amount
// This is a Wails-exposed method for the frontend
// Requirements: 8.1, 8.5 - Use connected wallet address and handle wallet operations
func (s *MarketplaceService) BuyPartialOrder(orderID string, tokenAmount string) (*TradeResult, error) {
	// Requirement 8.1: Use connected wallet address for marketplace operations
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   "no wallet connected",
		}, nil
	}

	// Requirement 8.1: Validate wallet address format
	if err := s.validateWalletAddress(currentAddress); err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("invalid wallet address: %v", err),
		}, nil
	}

	// Validate order ID and token amount
	if err := s.validateOrderID(orderID); err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("invalid order ID: %v", err),
		}, nil
	}

	if err := s.validateTokenAmount(tokenAmount); err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("invalid token amount: %v", err),
		}, nil
	}

	// Requirement 8.5: Execute partial trade using TradeService (handles wallet operations)
	result, err := s.tradeService.ExecutePartialTrade(orderID, currentAddress, tokenAmount)
	if err != nil {
		return &TradeResult{
			Success: false,
			OrderID: orderID,
			Error:   fmt.Sprintf("failed to execute partial trade: %v", err),
		}, nil
	}

	log.Printf("✅ Executed partial buy order %s (%s tokens) for wallet %s", orderID, tokenAmount, currentAddress)
	return result, nil
}

// CancelOrder cancels an existing sell order
// This is a Wails-exposed method for the frontend
// Requirements: 8.1, 8.3 - Use connected wallet address and verify requester is order creator
func (s *MarketplaceService) CancelOrder(orderID string) error {
	// Requirement 8.1: Use connected wallet address for marketplace operations
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return fmt.Errorf("no wallet connected")
	}

	// Requirement 8.1: Validate wallet address format
	if err := s.validateWalletAddress(currentAddress); err != nil {
		return fmt.Errorf("invalid wallet address: %w", err)
	}

	// Validate order ID format
	if err := s.validateOrderID(orderID); err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	// Requirement 8.3: Cancel order using OrderService (includes authorization checks)
	// The OrderService.CancelOrder method verifies the requester is the original order creator
	if err := s.orderService.CancelOrder(orderID, currentAddress); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	log.Printf("✅ Cancelled order %s for wallet %s", orderID, currentAddress)
	return nil
}

// GetUserOrders returns all orders for a specific wallet address
// This is a Wails-exposed method for the frontend
// Requirements: 7.2 - Return structured data with filtering capabilities
func (s *MarketplaceService) GetUserOrders(walletAddress string) ([]Order, error) {
	// Requirement 8.4: Validate wallet address and apply access control
	if err := s.validateWalletAddress(walletAddress); err != nil {
		return nil, fmt.Errorf("wallet address validation failed: %w", err)
	}

	// Requirement 8.4: Filter results based on connected wallet address for authorization
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return nil, fmt.Errorf("no wallet connected")
	}

	// Only allow users to access their own orders unless they're requesting their own
	if walletAddress != currentAddress {
		return nil, fmt.Errorf("unauthorized: can only access your own orders")
	}

	// Get user orders from KV store index
	ctx := context.Background()
	userOrdersKey := s.orderService.getUserOrdersKey(walletAddress)

	value, err := s.kvStore.GetMarketplaceData(ctx, userOrdersKey)
	if err != nil {
		// If index doesn't exist, return empty list
		return []Order{}, nil
	}

	// Parse order IDs from index
	var orderIDs []string
	if err := json.Unmarshal(value, &orderIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user orders index: %w", err)
	}

	// Fetch each order
	orders := make([]Order, 0, len(orderIDs))
	for _, orderID := range orderIDs {
		order, err := s.orderService.GetOrder(orderID)
		if err != nil {
			log.Printf("⚠️  Failed to get order %s from user index: %v", orderID, err)
			continue
		}
		orders = append(orders, *order)
	}

	// Sort orders by creation date (newest first)
	orders = s.sortOrders(orders, "date_desc")

	log.Printf("✅ Retrieved %d orders for user %s", len(orders), walletAddress)
	return orders, nil
}

// GetMarketStats returns current marketplace statistics and analytics
// This is a Wails-exposed method for the frontend
func (s *MarketplaceService) GetMarketStats() (*MarketStats, error) {
	return s.marketDataService.CalculateMarketStats()
}

// GetOrderHistory returns complete order history for a user
// This is a Wails-exposed method for the frontend
// Requirements: 4.1, 4.2 - Order history with status and trade details
func (s *MarketplaceService) GetOrderHistory(walletAddress string) (*OrderHistory, error) {
	// Requirement 8.4: Validate wallet address and apply access control
	if err := s.validateWalletAddress(walletAddress); err != nil {
		return nil, fmt.Errorf("wallet address validation failed: %w", err)
	}

	// Requirement 8.4: Filter results based on connected wallet address for authorization
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return nil, fmt.Errorf("no wallet connected")
	}

	// Only allow users to access their own order history
	if walletAddress != currentAddress {
		return nil, fmt.Errorf("unauthorized: can only access your own order history")
	}

	// Get user orders
	orders, err := s.GetUserOrders(walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	// Build order history entries with status history and trade details
	orderHistoryEntries := make([]OrderHistoryEntry, 0, len(orders))
	allTrades := []TradeHistoryEntry{}

	for _, order := range orders {
		// Get status history for this order
		statusHistory, err := s.getOrderStatusHistory(order.ID)
		if err != nil {
			log.Printf("⚠️  Failed to get status history for order %s: %v", order.ID, err)
			statusHistory = []OrderStatusChange{} // Continue with empty status history
		}

		// Get trades for this order
		orderTrades, err := s.getOrderTrades(order.ID)
		if err != nil {
			log.Printf("⚠️  Failed to get trades for order %s: %v", order.ID, err)
			orderTrades = []TradeHistoryEntry{} // Continue with empty trades
		}

		// Calculate filled amount across all trades
		filledAmount := s.calculateFilledAmount(orderTrades)

		// Create order history entry
		orderHistoryEntry := OrderHistoryEntry{
			Order:         order,
			StatusHistory: statusHistory,
			TradeCount:    len(orderTrades),
			FilledAmount:  filledAmount,
		}

		orderHistoryEntries = append(orderHistoryEntries, orderHistoryEntry)
		allTrades = append(allTrades, orderTrades...)
	}

	// Sort trades by timestamp (newest first)
	s.sortTradesByTimestamp(allTrades, false)

	history := &OrderHistory{
		Orders: orderHistoryEntries,
		Trades: allTrades,
		Total:  len(orderHistoryEntries),
	}

	log.Printf("✅ Retrieved order history for user %s: %d orders, %d trades", walletAddress, len(orderHistoryEntries), len(allTrades))
	return history, nil
}

// GetOrderStatusHistory returns the status change history for a specific order
// This is a Wails-exposed method for the frontend
// Requirements: 4.5 - Order status tracking with timestamps
func (s *MarketplaceService) GetOrderStatusHistory(orderID string) ([]OrderStatusChange, error) {
	// Validate order ID
	if err := s.validateOrderID(orderID); err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	// Get the order to verify ownership
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Verify user can access this order
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return nil, fmt.Errorf("no wallet connected")
	}

	if order.Seller != currentAddress {
		return nil, fmt.Errorf("unauthorized: can only access your own orders")
	}

	return s.getOrderStatusHistory(orderID)
}

// GetTradeHistory returns trade history for a user
// This is a Wails-exposed method for the frontend
// Requirements: 4.2 - Trade history with execution details
func (s *MarketplaceService) GetTradeHistory(walletAddress string) ([]TradeHistoryEntry, error) {
	// Requirement 8.4: Validate wallet address and apply access control
	if err := s.validateWalletAddress(walletAddress); err != nil {
		return nil, fmt.Errorf("wallet address validation failed: %w", err)
	}

	// Requirement 8.4: Filter results based on connected wallet address for authorization
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return nil, fmt.Errorf("no wallet connected")
	}

	// Only allow users to access their own trade history
	if walletAddress != currentAddress {
		return nil, fmt.Errorf("unauthorized: can only access your own trade history")
	}

	return s.getUserTradeHistory(walletAddress)
}

// GetOrderDetails returns detailed information about a specific order
// This is a Wails-exposed method for the frontend
// Requirements: 4.5 - Order detail metadata display
func (s *MarketplaceService) GetOrderDetails(orderID string) (*OrderHistoryEntry, error) {
	// Validate order ID
	if err := s.validateOrderID(orderID); err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	// Get the order
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Verify user can access this order
	currentAddress := s.walletService.GetCurrentAddress()
	if currentAddress == "" {
		return nil, fmt.Errorf("no wallet connected")
	}

	if order.Seller != currentAddress {
		return nil, fmt.Errorf("unauthorized: can only access your own orders")
	}

	// Get status history for this order
	statusHistory, err := s.getOrderStatusHistory(order.ID)
	if err != nil {
		log.Printf("⚠️  Failed to get status history for order %s: %v", order.ID, err)
		statusHistory = []OrderStatusChange{} // Continue with empty status history
	}

	// Get trades for this order
	orderTrades, err := s.getOrderTrades(order.ID)
	if err != nil {
		log.Printf("⚠️  Failed to get trades for order %s: %v", order.ID, err)
		orderTrades = []TradeHistoryEntry{} // Continue with empty trades
	}

	// Calculate filled amount across all trades
	filledAmount := s.calculateFilledAmount(orderTrades)

	// Create detailed order entry
	orderDetails := &OrderHistoryEntry{
		Order:         *order,
		StatusHistory: statusHistory,
		TradeCount:    len(orderTrades),
		FilledAmount:  filledAmount,
	}

	log.Printf("✅ Retrieved order details for %s", orderID)
	return orderDetails, nil
}

// validateOrderCreation validates order creation parameters
func (s *MarketplaceService) validateOrderCreation(seller string, tokenAmount, usdtPrice *big.Int) error {
	// TODO: Implement validation logic in task 2.1
	return fmt.Errorf("not implemented yet")
}

// validateGetActiveOrdersInput validates input parameters for GetActiveOrders
// validateGetActiveOrdersInput validates input parameters for GetActiveOrders
// Requirements: 7.1 - Input data validation
func (s *MarketplaceService) validateGetActiveOrdersInput(sortBy string, filterBy map[string]interface{}) error {
	// Validate sortBy parameter
	validSortOptions := []string{"", "price_asc", "price_desc", "amount_asc", "amount_desc", "date_asc", "date_desc"}
	validSort := false
	for _, option := range validSortOptions {
		if sortBy == option {
			validSort = true
			break
		}
	}
	if !validSort {
		return NewValidationError(
			"INVALID_SORT_OPTION",
			fmt.Sprintf("Invalid sort option: %s", sortBy),
			map[string]interface{}{
				"provided_option": sortBy,
				"valid_options":   validSortOptions,
			},
		)
	}

	// Validate filterBy parameters
	if filterBy != nil {
		validFilterKeys := []string{"minPrice", "maxPrice", "minAmount", "maxAmount", "seller"}
		for key := range filterBy {
			isValidKey := false
			for _, validFilterKey := range validFilterKeys {
				if key == validFilterKey {
					isValidKey = true
					break
				}
			}
			if !isValidKey {
				return NewValidationError(
					"INVALID_FILTER_KEY",
					fmt.Sprintf("Invalid filter key: %s", key),
					map[string]interface{}{
						"provided_key": key,
						"valid_keys":   validFilterKeys,
					},
				)
			}
		}
	}

	return nil
}

// validateWalletAddress validates a wallet address format
// Requirements: 8.1, 8.4 - Wallet address validation for authorization
func (s *MarketplaceService) validateWalletAddress(address string) error {
	if address == "" {
		return NewValidationError(
			"EMPTY_WALLET_ADDRESS",
			"Wallet address cannot be empty",
			map[string]interface{}{
				"provided_address": address,
			},
		)
	}

	// Basic Ethereum address validation (42 characters, starts with 0x)
	if len(address) != 42 || !strings.HasPrefix(address, "0x") {
		return NewValidationError(
			"INVALID_ADDRESS_FORMAT",
			"Invalid wallet address format",
			map[string]interface{}{
				"provided_address":   address,
				"expected_length":    42,
				"actual_length":      len(address),
				"expected_prefix":    "0x",
				"has_correct_prefix": strings.HasPrefix(address, "0x"),
			},
		)
	}

	// Validate hex characters
	if _, err := hex.DecodeString(address[2:]); err != nil {
		return NewValidationError(
			"INVALID_ADDRESS_HEX",
			"Invalid wallet address hex format",
			map[string]interface{}{
				"provided_address": address,
				"hex_error":        err.Error(),
			},
		)
	}

	return nil
}

// validateCreateSellOrderInput validates input parameters for CreateSellOrder
// Requirements: 7.1 - Input data validation
func (s *MarketplaceService) validateCreateSellOrderInput(tokenAmount, usdtPrice string) error {
	if tokenAmount == "" {
		return NewValidationError(
			"EMPTY_TOKEN_AMOUNT",
			"Token amount cannot be empty",
			map[string]interface{}{
				"provided_amount": tokenAmount,
			},
		)
	}
	if usdtPrice == "" {
		return NewValidationError(
			"EMPTY_USDT_PRICE",
			"USDT price cannot be empty",
			map[string]interface{}{
				"provided_price": usdtPrice,
			},
		)
	}

	// Validate token amount is a valid positive number
	if err := s.validateTokenAmount(tokenAmount); err != nil {
		return err // Already wrapped as MarketplaceError
	}

	// Validate USDT price is a valid positive number
	if err := s.validateUSDTPrice(usdtPrice); err != nil {
		return err // Already wrapped as MarketplaceError
	}

	return nil
}

// validateOrderID validates an order ID format
// Requirements: 7.1 - Input data validation
func (s *MarketplaceService) validateOrderID(orderID string) error {
	if orderID == "" {
		return NewValidationError(
			"EMPTY_ORDER_ID",
			"Order ID cannot be empty",
			map[string]interface{}{
				"provided_id": orderID,
			},
		)
	}

	// Order IDs should be hex strings (generated by generateOrderID)
	if len(orderID) != 32 { // 16 bytes = 32 hex characters
		return NewValidationError(
			"INVALID_ORDER_ID_LENGTH",
			"Invalid order ID length",
			map[string]interface{}{
				"provided_id":     orderID,
				"expected_length": 32,
				"actual_length":   len(orderID),
			},
		)
	}

	// Validate hex format
	if _, err := hex.DecodeString(orderID); err != nil {
		return NewValidationError(
			"INVALID_ORDER_ID_HEX",
			"Invalid order ID hex format",
			map[string]interface{}{
				"provided_id": orderID,
				"hex_error":   err.Error(),
			},
		)
	}

	return nil
}

// validateTokenAmount validates a token amount string
// Requirements: 7.1 - Input data validation
func (s *MarketplaceService) validateTokenAmount(tokenAmount string) error {
	if tokenAmount == "" {
		return NewValidationError(
			"EMPTY_TOKEN_AMOUNT",
			"Token amount cannot be empty",
			map[string]interface{}{
				"provided_amount": tokenAmount,
			},
		)
	}

	amount, ok := new(big.Int).SetString(tokenAmount, 10)
	if !ok {
		return NewValidationError(
			"INVALID_TOKEN_AMOUNT_FORMAT",
			"Invalid token amount format",
			map[string]interface{}{
				"provided_amount": tokenAmount,
			},
		)
	}

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return NewValidationError(
			"INVALID_TOKEN_AMOUNT_VALUE",
			"Token amount must be greater than zero",
			map[string]interface{}{
				"provided_amount": tokenAmount,
				"parsed_value":    amount.String(),
			},
		)
	}

	// Check for reasonable upper limit (prevent overflow attacks)
	maxAmount := new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil) // 10^30
	if amount.Cmp(maxAmount) > 0 {
		return NewValidationError(
			"TOKEN_AMOUNT_TOO_LARGE",
			"Token amount exceeds maximum allowed value",
			map[string]interface{}{
				"provided_amount": tokenAmount,
				"parsed_value":    amount.String(),
				"max_allowed":     maxAmount.String(),
			},
		)
	}

	return nil
}

// ensureRealTimeOrderUpdates ensures that order display updates happen in real-time
// Requirements: 2.4 - Real-time remaining amount updates for partial fills
func (s *MarketplaceService) ensureRealTimeOrderUpdates(orderID string) error {
	// Get the current order to check its state
	order, err := s.orderService.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order for real-time update: %w", err)
	}

	// Validate order information completeness before emitting updates
	if err := s.validateOrderInformationCompleteness(*order); err != nil {
		logMarketplaceWarning("ensureRealTimeOrderUpdates", fmt.Sprintf("Order %s failed completeness validation: %v", orderID, err))
		return fmt.Errorf("order information incomplete for real-time update: %w", err)
	}

	// Emit real-time update event for order status/amount changes
	// Requirements: 4.3 - Event-driven status updates for active orders
	if s.app != nil {
		// Create a real-time update event for the order
		event := MarketplaceRealtimeEvent{
			Type:      "order_updated",
			Timestamp: time.Now(),
			UserID:    order.Seller,
			Data: struct {
				OrderID         string `json:"orderID"`
				RemainingAmount string `json:"remainingAmount"`
				Status          string `json:"status"`
				UpdatedAt       string `json:"updatedAt"`
			}{
				OrderID:         order.ID,
				RemainingAmount: order.RemainingAmount,
				Status:          order.Status,
				UpdatedAt:       order.UpdatedAt.Format(time.RFC3339),
			},
		}

		// Emit to all connected clients for order book updates
		s.app.Event.Emit("marketplace:order_updated", event)

		// Emit targeted event to order owner
		s.app.Event.Emit(fmt.Sprintf("marketplace:user:%s:order_updated", order.Seller), event)

		logMarketplaceInfo("ensureRealTimeOrderUpdates", fmt.Sprintf("Emitted real-time update for order %s", orderID))
	}

	return nil
}

// validateActiveOrderManagement validates that active order management is working correctly
// Requirements: 2.4, 2.5 - Active order management with proper cleanup
func (s *MarketplaceService) validateActiveOrderManagement() error {
	// Get all active orders
	activeOrders, err := s.orderService.GetActiveOrders()
	if err != nil {
		return fmt.Errorf("failed to get active orders for validation: %w", err)
	}

	// Validate each active order
	for _, order := range activeOrders {
		// Requirement 2.5: Ensure only active orders are in active listings
		if order.Status != "active" {
			logMarketplaceWarning("validateActiveOrderManagement", fmt.Sprintf("Non-active order %s found in active listings (status: %s)", order.ID, order.Status))

			// Clean up by removing from active index
			if err := s.orderService.removeFromActiveOrdersIndex(order.ID); err != nil {
				logMarketplaceWarning("validateActiveOrderManagement", fmt.Sprintf("Failed to remove non-active order %s from active index: %v", order.ID, err))
			}
			continue
		}

		// Validate order information completeness
		if err := s.validateOrderInformationCompleteness(order); err != nil {
			logMarketplaceWarning("validateActiveOrderManagement", fmt.Sprintf("Active order %s has incomplete information: %v", order.ID, err))
		}

		// Requirement 2.4: Validate remaining amount is properly maintained for partial fills
		remainingAmount, ok := new(big.Int).SetString(order.RemainingAmount, 10)
		if !ok {
			logMarketplaceWarning("validateActiveOrderManagement", fmt.Sprintf("Active order %s has invalid remaining amount format: %s", order.ID, order.RemainingAmount))
			continue
		}

		// Remaining amount should be > 0 for active orders
		if remainingAmount.Cmp(big.NewInt(0)) <= 0 {
			logMarketplaceWarning("validateActiveOrderManagement", fmt.Sprintf("Active order %s has zero or negative remaining amount: %s", order.ID, order.RemainingAmount))

			// This order should be marked as filled
			if err := s.orderService.UpdateOrderStatus(order.ID, "filled"); err != nil {
				logMarketplaceWarning("validateActiveOrderManagement", fmt.Sprintf("Failed to update order %s status to filled: %v", order.ID, err))
			}
		}

		// Validate remaining amount doesn't exceed original amount
		originalAmount, ok := new(big.Int).SetString(order.TokenAmount, 10)
		if ok && remainingAmount.Cmp(originalAmount) > 0 {
			logMarketplaceWarning("validateActiveOrderManagement", fmt.Sprintf("Active order %s has remaining amount (%s) greater than original amount (%s)", order.ID, order.RemainingAmount, order.TokenAmount))
		}
	}

	logMarketplaceInfo("validateActiveOrderManagement", fmt.Sprintf("Validated %d active orders", len(activeOrders)))
	return nil
}

// cleanupInactiveOrders removes orders that are no longer active from active listings
// Requirements: 2.5 - Remove completed/cancelled orders from active listings
func (s *MarketplaceService) cleanupInactiveOrders() error {
	// Skip cleanup if KV store is not available (e.g., in tests)
	if s.kvStore == nil {
		logMarketplaceWarning("cleanupInactiveOrders", "KV store not available - skipping cleanup")
		return nil
	}

	// Get all active orders from index
	activeOrders, err := s.orderService.GetActiveOrders()
	if err != nil {
		return fmt.Errorf("failed to get active orders for cleanup: %w", err)
	}

	cleanedCount := 0
	for _, order := range activeOrders {
		// If order is not actually active, remove it from active index
		if order.Status != "active" {
			if err := s.orderService.removeFromActiveOrdersIndex(order.ID); err != nil {
				logMarketplaceWarning("cleanupInactiveOrders", fmt.Sprintf("Failed to remove inactive order %s from active index: %v", order.ID, err))
				continue
			}

			logMarketplaceInfo("cleanupInactiveOrders", fmt.Sprintf("Removed inactive order %s (status: %s) from active listings", order.ID, order.Status))
			cleanedCount++
		}
	}

	if cleanedCount > 0 {
		logMarketplaceInfo("cleanupInactiveOrders", fmt.Sprintf("Cleaned up %d inactive orders from active listings", cleanedCount))
	}

	return nil
}

// refreshActiveOrderAmounts ensures all active orders have correct remaining amounts
// Requirements: 2.4 - Partial fill amount updates in real-time
func (s *MarketplaceService) refreshActiveOrderAmounts() error {
	// Get all active orders
	activeOrders, err := s.orderService.GetActiveOrders()
	if err != nil {
		return fmt.Errorf("failed to get active orders for amount refresh: %w", err)
	}

	refreshedCount := 0
	for _, order := range activeOrders {
		// Ensure real-time updates are properly handled
		if err := s.ensureRealTimeOrderUpdates(order.ID); err != nil {
			logMarketplaceWarning("refreshActiveOrderAmounts", fmt.Sprintf("Failed to ensure real-time updates for order %s: %v", order.ID, err))
			continue
		}
		refreshedCount++
	}

	logMarketplaceInfo("refreshActiveOrderAmounts", fmt.Sprintf("Refreshed real-time updates for %d active orders", refreshedCount))
	return nil
}

// MaintainActiveOrderManagement performs comprehensive active order management maintenance
// This method can be called periodically to ensure proper order management
// Requirements: 2.4, 2.5 - Complete active order management
func (s *MarketplaceService) MaintainActiveOrderManagement() error {
	logMarketplaceInfo("MaintainActiveOrderManagement", "Starting active order management maintenance")

	// Step 1: Clean up inactive orders from active listings
	if err := s.cleanupInactiveOrders(); err != nil {
		logMarketplaceWarning("MaintainActiveOrderManagement", fmt.Sprintf("Failed to cleanup inactive orders: %v", err))
	}

	// Step 2: Validate active order management
	if err := s.validateActiveOrderManagement(); err != nil {
		logMarketplaceWarning("MaintainActiveOrderManagement", fmt.Sprintf("Active order validation failed: %v", err))
	}

	// Step 3: Refresh active order amounts and real-time updates
	if err := s.refreshActiveOrderAmounts(); err != nil {
		logMarketplaceWarning("MaintainActiveOrderManagement", fmt.Sprintf("Failed to refresh active order amounts: %v", err))
	}

	logMarketplaceInfo("MaintainActiveOrderManagement", "Completed active order management maintenance")
	return nil
}

// handleOrderStatusChange handles order status changes and ensures proper cleanup
// Requirements: 2.5 - Remove completed/cancelled orders from active listings
func (s *MarketplaceService) handleOrderStatusChange(orderID, oldStatus, newStatus string) error {
	// If order is no longer active, ensure it's removed from active listings
	if newStatus != "active" && oldStatus == "active" {
		// The order should already be removed from active index by UpdateOrderStatus,
		// but let's ensure it's properly handled
		if err := s.orderService.removeFromActiveOrdersIndex(orderID); err != nil {
			logMarketplaceWarning("handleOrderStatusChange", fmt.Sprintf("Failed to remove order %s from active index: %v", orderID, err))
		}

		logMarketplaceInfo("handleOrderStatusChange", fmt.Sprintf("Order %s removed from active listings (status: %s -> %s)", orderID, oldStatus, newStatus))
	}

	// Ensure real-time updates are sent
	if err := s.ensureRealTimeOrderUpdates(orderID); err != nil {
		logMarketplaceWarning("handleOrderStatusChange", fmt.Sprintf("Failed to ensure real-time updates for order %s: %v", orderID, err))
	}

	return nil
}

// validateOrderInformationCompleteness validates that an order contains all required information
// Requirements: 2.2 - Order display information completeness
func (s *MarketplaceService) validateOrderInformationCompleteness(order Order) error {
	// Requirement 2.2: Order display must show token amount, USDT price, price per token, and seller address

	// Validate token amount is present and valid
	if order.TokenAmount == "" {
		return NewValidationError(
			"MISSING_TOKEN_AMOUNT",
			"Order is missing token amount",
			map[string]interface{}{
				"order_id": order.ID,
			},
		)
	}

	if _, ok := new(big.Int).SetString(order.TokenAmount, 10); !ok {
		return NewValidationError(
			"INVALID_TOKEN_AMOUNT_FORMAT",
			"Order has invalid token amount format",
			map[string]interface{}{
				"order_id":     order.ID,
				"token_amount": order.TokenAmount,
			},
		)
	}

	// Validate USDT price is present and valid
	if order.USDTPrice == "" {
		return NewValidationError(
			"MISSING_USDT_PRICE",
			"Order is missing USDT price",
			map[string]interface{}{
				"order_id": order.ID,
			},
		)
	}

	if _, ok := new(big.Int).SetString(order.USDTPrice, 10); !ok {
		return NewValidationError(
			"INVALID_USDT_PRICE_FORMAT",
			"Order has invalid USDT price format",
			map[string]interface{}{
				"order_id":   order.ID,
				"usdt_price": order.USDTPrice,
			},
		)
	}

	// Validate price per token is present and valid
	if order.PricePerToken == "" {
		return NewValidationError(
			"MISSING_PRICE_PER_TOKEN",
			"Order is missing price per token",
			map[string]interface{}{
				"order_id": order.ID,
			},
		)
	}

	if _, ok := new(big.Int).SetString(order.PricePerToken, 10); !ok {
		return NewValidationError(
			"INVALID_PRICE_PER_TOKEN_FORMAT",
			"Order has invalid price per token format",
			map[string]interface{}{
				"order_id":        order.ID,
				"price_per_token": order.PricePerToken,
			},
		)
	}

	// Validate seller address is present and valid
	if order.Seller == "" {
		return NewValidationError(
			"MISSING_SELLER_ADDRESS",
			"Order is missing seller address",
			map[string]interface{}{
				"order_id": order.ID,
			},
		)
	}

	// Validate seller address format (basic Ethereum address validation)
	if err := s.validateWalletAddress(order.Seller); err != nil {
		return NewValidationError(
			"INVALID_SELLER_ADDRESS",
			"Order has invalid seller address format",
			map[string]interface{}{
				"order_id":         order.ID,
				"seller_address":   order.Seller,
				"validation_error": err.Error(),
			},
		)
	}

	// Validate remaining amount is present and valid (for partial fill tracking)
	if order.RemainingAmount == "" {
		return NewValidationError(
			"MISSING_REMAINING_AMOUNT",
			"Order is missing remaining amount",
			map[string]interface{}{
				"order_id": order.ID,
			},
		)
	}

	if _, ok := new(big.Int).SetString(order.RemainingAmount, 10); !ok {
		return NewValidationError(
			"INVALID_REMAINING_AMOUNT_FORMAT",
			"Order has invalid remaining amount format",
			map[string]interface{}{
				"order_id":         order.ID,
				"remaining_amount": order.RemainingAmount,
			},
		)
	}

	// Validate order ID is present
	if order.ID == "" {
		return NewValidationError(
			"MISSING_ORDER_ID",
			"Order is missing order ID",
			map[string]interface{}{
				"seller": order.Seller,
			},
		)
	}

	// Validate status is present and valid
	validStatuses := []string{"active", "filled", "cancelled"}
	statusValid := false
	for _, validStatus := range validStatuses {
		if order.Status == validStatus {
			statusValid = true
			break
		}
	}

	if !statusValid {
		return NewValidationError(
			"INVALID_ORDER_STATUS",
			"Order has invalid status",
			map[string]interface{}{
				"order_id":       order.ID,
				"current_status": order.Status,
				"valid_statuses": validStatuses,
			},
		)
	}

	return nil
}

// validateUSDTPrice validates a USDT price string
// Requirements: 7.1 - Input data validation
func (s *MarketplaceService) validateUSDTPrice(usdtPrice string) error {
	if usdtPrice == "" {
		return NewValidationError(
			"EMPTY_USDT_PRICE",
			"USDT price cannot be empty",
			map[string]interface{}{
				"provided_price": usdtPrice,
			},
		)
	}

	price, ok := new(big.Int).SetString(usdtPrice, 10)
	if !ok {
		return NewValidationError(
			"INVALID_USDT_PRICE_FORMAT",
			"Invalid USDT price format",
			map[string]interface{}{
				"provided_price": usdtPrice,
			},
		)
	}

	if price.Cmp(big.NewInt(0)) <= 0 {
		return NewValidationError(
			"INVALID_USDT_PRICE_VALUE",
			"USDT price must be greater than zero",
			map[string]interface{}{
				"provided_price": usdtPrice,
				"parsed_value":   price.String(),
			},
		)
	}

	// Check for reasonable upper limit (prevent overflow attacks)
	maxPrice := new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil) // 10^30
	if price.Cmp(maxPrice) > 0 {
		return NewValidationError(
			"USDT_PRICE_TOO_LARGE",
			"USDT price exceeds maximum allowed value",
			map[string]interface{}{
				"provided_price": usdtPrice,
				"parsed_value":   price.String(),
				"max_allowed":    maxPrice.String(),
			},
		)
	}

	return nil
}

// applyOrderFilters applies filtering criteria to orders
// Requirements: 7.2 - Filtering capabilities
func (s *MarketplaceService) applyOrderFilters(orders []Order, filterBy map[string]interface{}) []Order {
	filtered := make([]Order, 0, len(orders))

	for _, order := range orders {
		if s.orderMatchesFilters(order, filterBy) {
			filtered = append(filtered, order)
		}
	}

	return filtered
}

// orderMatchesFilters checks if an order matches the specified filters
func (s *MarketplaceService) orderMatchesFilters(order Order, filterBy map[string]interface{}) bool {
	// Parse order values for comparison
	pricePerToken, ok := new(big.Int).SetString(order.PricePerToken, 10)
	if !ok {
		return false
	}

	remainingAmount, ok := new(big.Int).SetString(order.RemainingAmount, 10)
	if !ok {
		return false
	}

	// Check each filter
	for key, value := range filterBy {
		switch key {
		case "minPrice":
			if minPrice, ok := s.parseFilterValue(value); ok {
				if pricePerToken.Cmp(minPrice) < 0 {
					return false
				}
			}
		case "maxPrice":
			if maxPrice, ok := s.parseFilterValue(value); ok {
				if pricePerToken.Cmp(maxPrice) > 0 {
					return false
				}
			}
		case "minAmount":
			if minAmount, ok := s.parseFilterValue(value); ok {
				if remainingAmount.Cmp(minAmount) < 0 {
					return false
				}
			}
		case "maxAmount":
			if maxAmount, ok := s.parseFilterValue(value); ok {
				if remainingAmount.Cmp(maxAmount) > 0 {
					return false
				}
			}
		case "seller":
			if sellerFilter, ok := value.(string); ok {
				if order.Seller != sellerFilter {
					return false
				}
			}
		}
	}

	return true
}

// parseFilterValue parses a filter value to big.Int
func (s *MarketplaceService) parseFilterValue(value interface{}) (*big.Int, bool) {
	switch v := value.(type) {
	case string:
		if result, ok := new(big.Int).SetString(v, 10); ok {
			return result, true
		}
	case float64:
		return big.NewInt(int64(v)), true
	case int:
		return big.NewInt(int64(v)), true
	case int64:
		return big.NewInt(v), true
	}
	return nil, false
}

// sortOrders sorts orders based on the specified criteria
// Requirements: 7.2 - Sorting capabilities
func (s *MarketplaceService) sortOrders(orders []Order, sortBy string) []Order {
	if sortBy == "" {
		return orders
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]Order, len(orders))
	copy(sorted, orders)

	switch sortBy {
	case "price_asc":
		s.sortOrdersByPrice(sorted, true)
	case "price_desc":
		s.sortOrdersByPrice(sorted, false)
	case "amount_asc":
		s.sortOrdersByAmount(sorted, true)
	case "amount_desc":
		s.sortOrdersByAmount(sorted, false)
	case "date_asc":
		s.sortOrdersByDate(sorted, true)
	case "date_desc":
		s.sortOrdersByDate(sorted, false)
	}

	return sorted
}

// sortOrdersByPrice sorts orders by price per token
func (s *MarketplaceService) sortOrdersByPrice(orders []Order, ascending bool) {
	for i := 0; i < len(orders)-1; i++ {
		for j := 0; j < len(orders)-i-1; j++ {
			price1, ok1 := new(big.Int).SetString(orders[j].PricePerToken, 10)
			price2, ok2 := new(big.Int).SetString(orders[j+1].PricePerToken, 10)

			if !ok1 || !ok2 {
				continue
			}

			shouldSwap := false
			if ascending {
				shouldSwap = price1.Cmp(price2) > 0
			} else {
				shouldSwap = price1.Cmp(price2) < 0
			}

			if shouldSwap {
				orders[j], orders[j+1] = orders[j+1], orders[j]
			}
		}
	}
}

// sortOrdersByAmount sorts orders by remaining amount
func (s *MarketplaceService) sortOrdersByAmount(orders []Order, ascending bool) {
	for i := 0; i < len(orders)-1; i++ {
		for j := 0; j < len(orders)-i-1; j++ {
			amount1, ok1 := new(big.Int).SetString(orders[j].RemainingAmount, 10)
			amount2, ok2 := new(big.Int).SetString(orders[j+1].RemainingAmount, 10)

			if !ok1 || !ok2 {
				continue
			}

			shouldSwap := false
			if ascending {
				shouldSwap = amount1.Cmp(amount2) > 0
			} else {
				shouldSwap = amount1.Cmp(amount2) < 0
			}

			if shouldSwap {
				orders[j], orders[j+1] = orders[j+1], orders[j]
			}
		}
	}
}

// sortOrdersByDate sorts orders by creation date
func (s *MarketplaceService) sortOrdersByDate(orders []Order, ascending bool) {
	for i := 0; i < len(orders)-1; i++ {
		for j := 0; j < len(orders)-i-1; j++ {
			shouldSwap := false
			if ascending {
				shouldSwap = orders[j].CreatedAt.After(orders[j+1].CreatedAt)
			} else {
				shouldSwap = orders[j].CreatedAt.Before(orders[j+1].CreatedAt)
			}

			if shouldSwap {
				orders[j], orders[j+1] = orders[j+1], orders[j]
			}
		}
	}
}

// Note: Key generation methods are handled by individual services (OrderService, MarketDataService)
// All services use the dedicated P2P marketplace namespace via KVStore.StoreMarketplaceData() methods

// Order History Helper Methods
// Requirements: 4.1, 4.2, 4.4, 4.5 - Order and trade history management

// getOrderStatusHistory retrieves the status change history for an order
func (s *MarketplaceService) getOrderStatusHistory(orderID string) ([]OrderStatusChange, error) {
	ctx := context.Background()
	key := s.getOrderStatusHistoryKey(orderID)

	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err != nil {
		// If no status history exists, create default from current order
		order, orderErr := s.orderService.GetOrder(orderID)
		if orderErr != nil {
			return []OrderStatusChange{}, nil // Return empty if order doesn't exist
		}

		// Create default status history based on current order state
		defaultHistory := []OrderStatusChange{
			{
				Status:    "active",
				Timestamp: order.CreatedAt,
				TxHash:    order.TxHash,
			},
		}

		// Add current status if different from active
		if order.Status != "active" {
			defaultHistory = append(defaultHistory, OrderStatusChange{
				Status:    order.Status,
				Timestamp: order.UpdatedAt,
			})
		}

		return defaultHistory, nil
	}

	var statusHistory []OrderStatusChange
	if err := json.Unmarshal(value, &statusHistory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status history: %w", err)
	}

	return statusHistory, nil
}

// addOrderStatusChange adds a status change to an order's history
func (s *MarketplaceService) addOrderStatusChange(orderID, status, txHash string) error {
	// Get existing status history
	statusHistory, err := s.getOrderStatusHistory(orderID)
	if err != nil {
		return fmt.Errorf("failed to get existing status history: %w", err)
	}

	// Add new status change
	statusChange := OrderStatusChange{
		Status:    status,
		Timestamp: time.Now(),
		TxHash:    txHash,
	}

	statusHistory = append(statusHistory, statusChange)

	// Store updated status history
	ctx := context.Background()
	key := s.getOrderStatusHistoryKey(orderID)

	statusHistoryJSON, err := json.Marshal(statusHistory)
	if err != nil {
		return fmt.Errorf("failed to marshal status history: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, key, statusHistoryJSON); err != nil {
		return fmt.Errorf("failed to store status history: %w", err)
	}

	return nil
}

// getOrderTrades retrieves all trades for a specific order
func (s *MarketplaceService) getOrderTrades(orderID string) ([]TradeHistoryEntry, error) {
	ctx := context.Background()

	// Get trade IDs for this order
	key := s.getOrderTradesKey(orderID)
	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err != nil {
		// No trades for this order
		return []TradeHistoryEntry{}, nil
	}

	var tradeIDs []string
	if err := json.Unmarshal(value, &tradeIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trade IDs: %w", err)
	}

	// Fetch each trade and build history entries
	trades := make([]TradeHistoryEntry, 0, len(tradeIDs))
	for _, tradeID := range tradeIDs {
		trade, err := s.getTrade(tradeID)
		if err != nil {
			log.Printf("⚠️  Failed to get trade %s: %v", tradeID, err)
			continue
		}

		// Get order details for the trade
		order, err := s.orderService.GetOrder(trade.OrderID)
		if err != nil {
			log.Printf("⚠️  Failed to get order details for trade %s: %v", tradeID, err)
			continue
		}

		orderSummary := OrderSummary{
			ID:            order.ID,
			TokenAmount:   order.TokenAmount,
			USDTPrice:     order.USDTPrice,
			PricePerToken: order.PricePerToken,
		}

		tradeHistoryEntry := TradeHistoryEntry{
			Trade:        *trade,
			OrderDetails: orderSummary,
		}

		trades = append(trades, tradeHistoryEntry)
	}

	return trades, nil
}

// getUserTradeHistory retrieves all trades for a user (as buyer or seller)
func (s *MarketplaceService) getUserTradeHistory(walletAddress string) ([]TradeHistoryEntry, error) {
	// Get user trade IDs (both as buyer and seller)
	buyerTradeIDs, err := s.getUserTradeIDs(walletAddress, "buyer")
	if err != nil {
		log.Printf("⚠️  Failed to get buyer trades for %s: %v", walletAddress, err)
		buyerTradeIDs = []string{}
	}

	sellerTradeIDs, err := s.getUserTradeIDs(walletAddress, "seller")
	if err != nil {
		log.Printf("⚠️  Failed to get seller trades for %s: %v", walletAddress, err)
		sellerTradeIDs = []string{}
	}

	// Combine and deduplicate trade IDs
	allTradeIDs := append(buyerTradeIDs, sellerTradeIDs...)
	uniqueTradeIDs := s.deduplicateStrings(allTradeIDs)

	// Fetch each trade and build history entries
	trades := make([]TradeHistoryEntry, 0, len(uniqueTradeIDs))
	for _, tradeID := range uniqueTradeIDs {
		trade, err := s.getTrade(tradeID)
		if err != nil {
			log.Printf("⚠️  Failed to get trade %s: %v", tradeID, err)
			continue
		}

		// Get order details for the trade
		order, err := s.orderService.GetOrder(trade.OrderID)
		if err != nil {
			log.Printf("⚠️  Failed to get order details for trade %s: %v", tradeID, err)
			continue
		}

		orderSummary := OrderSummary{
			ID:            order.ID,
			TokenAmount:   order.TokenAmount,
			USDTPrice:     order.USDTPrice,
			PricePerToken: order.PricePerToken,
		}

		tradeHistoryEntry := TradeHistoryEntry{
			Trade:        *trade,
			OrderDetails: orderSummary,
		}

		trades = append(trades, tradeHistoryEntry)
	}

	// Sort by timestamp (newest first)
	s.sortTradesByTimestamp(trades, false)

	return trades, nil
}

// getTrade retrieves a trade by ID
func (s *MarketplaceService) getTrade(tradeID string) (*Trade, error) {
	ctx := context.Background()
	key := s.getTradeKey(tradeID)

	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade from KV: %w", err)
	}

	var trade Trade
	if err := json.Unmarshal(value, &trade); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trade: %w", err)
	}

	return &trade, nil
}

// getUserTradeIDs retrieves trade IDs for a user in a specific role (buyer or seller)
func (s *MarketplaceService) getUserTradeIDs(walletAddress, role string) ([]string, error) {
	ctx := context.Background()
	key := s.getUserTradesKey(walletAddress, role)

	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err != nil {
		// No trades for this user in this role
		return []string{}, nil
	}

	var tradeIDs []string
	if err := json.Unmarshal(value, &tradeIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trade IDs: %w", err)
	}

	return tradeIDs, nil
}

// addTradeToOrderHistory adds a trade to an order's trade history
func (s *MarketplaceService) addTradeToOrderHistory(orderID, tradeID string) error {
	ctx := context.Background()
	key := s.getOrderTradesKey(orderID)

	// Get existing trade IDs
	var tradeIDs []string
	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err == nil {
		if err := json.Unmarshal(value, &tradeIDs); err != nil {
			return fmt.Errorf("failed to unmarshal existing trade IDs: %w", err)
		}
	}

	// Add new trade ID
	tradeIDs = append(tradeIDs, tradeID)

	// Store updated trade IDs
	tradeIDsJSON, err := json.Marshal(tradeIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal trade IDs: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, key, tradeIDsJSON); err != nil {
		return fmt.Errorf("failed to store trade IDs: %w", err)
	}

	return nil
}

// addTradeToUserHistory adds a trade to a user's trade history
func (s *MarketplaceService) addTradeToUserHistory(walletAddress, role, tradeID string) error {
	ctx := context.Background()
	key := s.getUserTradesKey(walletAddress, role)

	// Get existing trade IDs
	var tradeIDs []string
	value, err := s.kvStore.GetMarketplaceData(ctx, key)
	if err == nil {
		if err := json.Unmarshal(value, &tradeIDs); err != nil {
			return fmt.Errorf("failed to unmarshal existing trade IDs: %w", err)
		}
	}

	// Add new trade ID
	tradeIDs = append(tradeIDs, tradeID)

	// Store updated trade IDs
	tradeIDsJSON, err := json.Marshal(tradeIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal trade IDs: %w", err)
	}

	if err := s.kvStore.StoreMarketplaceData(ctx, key, tradeIDsJSON); err != nil {
		return fmt.Errorf("failed to store trade IDs: %w", err)
	}

	return nil
}

// calculateFilledAmount calculates the total filled amount from trades
func (s *MarketplaceService) calculateFilledAmount(trades []TradeHistoryEntry) string {
	totalFilled := big.NewInt(0)

	for _, trade := range trades {
		tokenAmount, ok := new(big.Int).SetString(trade.TokenAmount, 10)
		if ok {
			totalFilled.Add(totalFilled, tokenAmount)
		}
	}

	return totalFilled.String()
}

// sortTradesByTimestamp sorts trades by timestamp
func (s *MarketplaceService) sortTradesByTimestamp(trades []TradeHistoryEntry, ascending bool) {
	for i := 0; i < len(trades)-1; i++ {
		for j := 0; j < len(trades)-i-1; j++ {
			shouldSwap := false
			if ascending {
				shouldSwap = trades[j].Timestamp.After(trades[j+1].Timestamp)
			} else {
				shouldSwap = trades[j].Timestamp.Before(trades[j+1].Timestamp)
			}

			if shouldSwap {
				trades[j], trades[j+1] = trades[j+1], trades[j]
			}
		}
	}
}

// deduplicateStrings removes duplicate strings from a slice
func (s *MarketplaceService) deduplicateStrings(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, str := range input {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

// Key generation methods for order history
func (s *MarketplaceService) getOrderStatusHistoryKey(orderID string) string {
	return fmt.Sprintf("order:%s:status_history", orderID)
}

func (s *MarketplaceService) getOrderTradesKey(orderID string) string {
	return fmt.Sprintf("order:%s:trades", orderID)
}

func (s *MarketplaceService) getUserTradesKey(walletAddress, role string) string {
	return fmt.Sprintf("user:%s:trades:%s", walletAddress, role)
}

func (s *MarketplaceService) getTradeKey(tradeID string) string {
	return fmt.Sprintf("trade:%s", tradeID)
}
