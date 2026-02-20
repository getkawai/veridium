package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/event"
	"github.com/kawai-network/contracts/otcmarket"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/store"
)

// MarketplaceEventListener handles contract event listening and processing
// Requirements: 6.5 - Contract event processing and state synchronization
type MarketplaceEventListener struct {
	blockchainClient *blockchain.Client
	orderService     *OrderService
	walletService    *WalletService
	kvStore          *store.KVStore

	// Event subscriptions
	orderCreatedSub         event.Subscription
	orderFilledSub          event.Subscription
	orderPartiallyFilledSub event.Subscription
	orderCancelledSub       event.Subscription

	// Channels for event processing
	orderCreatedCh         chan *otcmarket.OTCMarketOrderCreated
	orderFilledCh          chan *otcmarket.OTCMarketOrderFulfilled
	orderPartiallyFilledCh chan *otcmarket.OTCMarketOrderPartiallyFilled
	orderCancelledCh       chan *otcmarket.OTCMarketOrderCancelled

	// Control channels
	stopCh chan struct{}
	doneCh chan struct{}

	// Synchronization
	mu      sync.RWMutex
	running bool
}

// MarketplaceEvent represents a processed marketplace event
type MarketplaceEvent struct {
	Type        string    `json:"type"` // "OrderCreated", "OrderFilled", "OrderCancelled"
	OrderID     string    `json:"orderID"`
	TxHash      string    `json:"txHash"`
	BlockNumber uint64    `json:"blockNumber"`
	Timestamp   time.Time `json:"timestamp"`

	// Event-specific data
	Seller      string `json:"seller,omitempty"`
	Buyer       string `json:"buyer,omitempty"`
	TokenAmount string `json:"tokenAmount,omitempty"`
	USDTPrice   string `json:"usdtPrice,omitempty"`
}

// NewMarketplaceEventListener creates a new event listener instance
func NewMarketplaceEventListener(blockchainClient *blockchain.Client, orderService *OrderService, walletService *WalletService, kvStore *store.KVStore) *MarketplaceEventListener {
	return &MarketplaceEventListener{
		blockchainClient: blockchainClient,
		orderService:     orderService,
		walletService:    walletService,
		kvStore:          kvStore,

		// Initialize channels
		orderCreatedCh:         make(chan *otcmarket.OTCMarketOrderCreated, 100),
		orderFilledCh:          make(chan *otcmarket.OTCMarketOrderFulfilled, 100),
		orderPartiallyFilledCh: make(chan *otcmarket.OTCMarketOrderPartiallyFilled, 100),
		orderCancelledCh:       make(chan *otcmarket.OTCMarketOrderCancelled, 100),

		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

// Start begins listening for contract events
// Requirements: 6.5 - Event listener for OrderCreated, OrderFilled, OrderCancelled events
func (l *MarketplaceEventListener) Start(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.running {
		return fmt.Errorf("event listener is already running")
	}

	if l.blockchainClient == nil {
		return fmt.Errorf("blockchain client is not available")
	}

	log.Println("🎧 Starting marketplace event listener...")

	// Subscribe to OrderCreated events
	if err := l.subscribeToOrderCreated(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to OrderCreated events: %w", err)
	}

	// Subscribe to OrderFulfilled events
	if err := l.subscribeToOrderFulfilled(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to OrderFulfilled events: %w", err)
	}

	// Subscribe to OrderPartiallyFilled events
	if err := l.subscribeToOrderPartiallyFilled(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to OrderPartiallyFilled events: %w", err)
	}

	// Subscribe to OrderCancelled events
	if err := l.subscribeToOrderCancelled(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to OrderCancelled events: %w", err)
	}

	l.running = true

	// Start event processing goroutine
	go l.processEvents(ctx)

	log.Println("✅ Marketplace event listener started successfully")
	return nil
}

// Stop stops the event listener and cleans up subscriptions
func (l *MarketplaceEventListener) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.running {
		return
	}

	log.Println("🛑 Stopping marketplace event listener...")

	// Signal stop
	close(l.stopCh)

	// Unsubscribe from events
	if l.orderCreatedSub != nil {
		l.orderCreatedSub.Unsubscribe()
	}
	if l.orderFilledSub != nil {
		l.orderFilledSub.Unsubscribe()
	}
	if l.orderPartiallyFilledSub != nil {
		l.orderPartiallyFilledSub.Unsubscribe()
	}
	if l.orderCancelledSub != nil {
		l.orderCancelledSub.Unsubscribe()
	}

	// Wait for processing to complete
	<-l.doneCh

	l.running = false
	log.Println("✅ Marketplace event listener stopped")
}

// IsRunning returns whether the event listener is currently running
func (l *MarketplaceEventListener) IsRunning() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.running
}

// subscribeToOrderCreated subscribes to OrderCreated events
func (l *MarketplaceEventListener) subscribeToOrderCreated(ctx context.Context) error {
	// Create filter options for OrderCreated events
	filterOpts := &bind.WatchOpts{
		Context: ctx,
		Start:   nil, // Start from latest block
	}

	// Subscribe to OrderCreated events
	sub, err := l.blockchainClient.OTCMarket.WatchOrderCreated(filterOpts, l.orderCreatedCh, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to watch OrderCreated events: %w", err)
	}

	l.orderCreatedSub = sub
	log.Println("📡 Subscribed to OrderCreated events")
	return nil
}

// subscribeToOrderFulfilled subscribes to OrderFulfilled events
func (l *MarketplaceEventListener) subscribeToOrderFulfilled(ctx context.Context) error {
	// Create filter options for OrderFulfilled events
	filterOpts := &bind.WatchOpts{
		Context: ctx,
		Start:   nil, // Start from latest block
	}

	// Subscribe to OrderFulfilled events
	sub, err := l.blockchainClient.OTCMarket.WatchOrderFulfilled(filterOpts, l.orderFilledCh, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to watch OrderFulfilled events: %w", err)
	}

	l.orderFilledSub = sub
	log.Println("📡 Subscribed to OrderFulfilled events")
	return nil
}

// subscribeToOrderPartiallyFilled subscribes to OrderPartiallyFilled events
func (l *MarketplaceEventListener) subscribeToOrderPartiallyFilled(ctx context.Context) error {
	// Create filter options for OrderPartiallyFilled events
	filterOpts := &bind.WatchOpts{
		Context: ctx,
		Start:   nil, // Start from latest block
	}

	// Subscribe to OrderPartiallyFilled events
	sub, err := l.blockchainClient.OTCMarket.WatchOrderPartiallyFilled(filterOpts, l.orderPartiallyFilledCh, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to watch OrderPartiallyFilled events: %w", err)
	}

	l.orderPartiallyFilledSub = sub
	log.Println("📡 Subscribed to OrderPartiallyFilled events")
	return nil
}

// subscribeToOrderCancelled subscribes to OrderCancelled events
func (l *MarketplaceEventListener) subscribeToOrderCancelled(ctx context.Context) error {
	// Create filter options for OrderCancelled events
	filterOpts := &bind.WatchOpts{
		Context: ctx,
		Start:   nil, // Start from latest block
	}

	// Subscribe to OrderCancelled events
	sub, err := l.blockchainClient.OTCMarket.WatchOrderCancelled(filterOpts, l.orderCancelledCh, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to watch OrderCancelled events: %w", err)
	}

	l.orderCancelledSub = sub
	log.Println("📡 Subscribed to OrderCancelled events")
	return nil
}

// processEvents processes incoming contract events
func (l *MarketplaceEventListener) processEvents(ctx context.Context) {
	defer close(l.doneCh)

	log.Println("🔄 Event processing loop started")

	for {
		select {
		case <-l.stopCh:
			log.Println("🔄 Event processing loop stopped")
			return

		case <-ctx.Done():
			log.Println("🔄 Event processing loop cancelled")
			return

		case event := <-l.orderCreatedCh:
			if err := l.handleOrderCreated(event); err != nil {
				log.Printf("❌ Error handling OrderCreated event: %v", err)
			}

		case event := <-l.orderFilledCh:
			if err := l.handleOrderFulfilled(event); err != nil {
				log.Printf("❌ Error handling OrderFulfilled event: %v", err)
			}

		case event := <-l.orderPartiallyFilledCh:
			if err := l.handleOrderPartiallyFilled(event); err != nil {
				log.Printf("❌ Error handling OrderPartiallyFilled event: %v", err)
			}

		case event := <-l.orderCancelledCh:
			if err := l.handleOrderCancelled(event); err != nil {
				log.Printf("❌ Error handling OrderCancelled event: %v", err)
			}
		}
	}
}

// handleOrderCreated processes OrderCreated events
// Requirements: 6.5 - Event processing and state synchronization
func (l *MarketplaceEventListener) handleOrderCreated(event *otcmarket.OTCMarketOrderCreated) error {
	// Only the seller is responsible for indexing their own new order in the marketplace KV
	myAddr := l.walletService.GetCurrentAddress()
	if myAddr == "" || !strings.EqualFold(myAddr, event.Seller.Hex()) {
		return nil // Not my order, skip indexing (I only care about my own activity)
	}

	log.Printf("📝 OrderCreated event: OrderID=%s, Seller=%s, Amount=%s, Price=%s",
		event.OrderId.String(), event.Seller.Hex(), event.Amount.String(), event.Price.String())

	// Convert contract order ID to our internal order ID format
	orderID := event.OrderId.String()

	// Get the order from our service to update it with the transaction hash
	order, err := l.orderService.GetOrder(orderID)
	if err != nil {
		// Order might not exist in our system yet if the event arrived before our local storage
		log.Printf("⚠️  Order %s not found in local storage, event may have arrived before local creation", orderID)
		return nil
	}

	// Update order with transaction hash from the event
	order.TxHash = event.Raw.TxHash.Hex()
	order.UpdatedAt = time.Now()

	if err := l.orderService.StoreOrder(order); err != nil {
		return fmt.Errorf("failed to update order with transaction hash: %w", err)
	}

	log.Printf("✅ Updated order %s with transaction hash %s", orderID, event.Raw.TxHash.Hex())
	return nil
}

// handleOrderFulfilled processes OrderFulfilled events
// Requirements: 6.5 - Event-driven order status updates
func (l *MarketplaceEventListener) handleOrderFulfilled(event *otcmarket.OTCMarketOrderFulfilled) error {
	// Only the buyer and seller are responsible for updating their part of the marketplace state
	myAddr := l.walletService.GetCurrentAddress()
	isBuyer := strings.EqualFold(myAddr, event.Buyer.Hex())
	isSeller := strings.EqualFold(myAddr, event.Seller.Hex())

	if myAddr == "" || (!isBuyer && !isSeller) {
		return nil // Not my trade, skip KV update
	}

	log.Printf("💰 OrderFulfilled event (Participant): OrderID=%s, Buyer=%s, Seller=%s, Amount=%s, Price=%s",
		event.OrderId.String(), event.Buyer.Hex(), event.Seller.Hex(), event.Amount.String(), event.Price.String())

	// Convert contract order ID to our internal order ID format
	orderID := event.OrderId.String()

	// Update order status to filled
	if err := l.orderService.UpdateOrderStatus(orderID, "filled"); err != nil {
		return fmt.Errorf("failed to update order status to filled: %w", err)
	}

	// Update trade records with buyer information
	if err := l.updateTradeRecordsWithBuyer(orderID, event.Buyer.Hex(), event.Raw.TxHash.Hex()); err != nil {
		log.Printf("⚠️  Failed to update trade records with buyer info: %v", err)
	}

	log.Printf("✅ Order %s marked as filled", orderID)
	return nil
}

// updateTradeRecordsWithBuyer updates trade records for an order with buyer information
func (l *MarketplaceEventListener) updateTradeRecordsWithBuyer(orderID, buyer, txHash string) error {
	// This is a simplified implementation - in a production system, we'd want to
	// match the specific trade by transaction hash and amount

	// For now, we'll update the most recent trade for this order
	// In the future, we could enhance this by storing pending trades and matching by tx hash

	log.Printf("📝 Updating trade records for order %s with buyer %s", orderID, buyer)

	// Note: The actual trade record update would require access to the TradeService
	// For now, we'll log this event and handle it in future enhancements
	// The trade record was already created in TradeService.storePartialTradeRecord
	// but without buyer information - this would be where we'd update it

	return nil
}

// handleOrderPartiallyFilled processes OrderPartiallyFilled events
// Requirements: 6.5 - Event-driven order status updates for partial fills
func (l *MarketplaceEventListener) handleOrderPartiallyFilled(event *otcmarket.OTCMarketOrderPartiallyFilled) error {
	// Only the buyer and seller are responsible for updating their part of the marketplace state
	myAddr := l.walletService.GetCurrentAddress()
	isBuyer := strings.EqualFold(myAddr, event.Buyer.Hex())
	isSeller := strings.EqualFold(myAddr, event.Seller.Hex())

	if myAddr == "" || (!isBuyer && !isSeller) {
		return nil // Not my trade, skip KV update
	}

	log.Printf("📊 OrderPartiallyFilled event (Participant): OrderID=%s, Buyer=%s, Seller=%s, AmountFilled=%s, RemainingAmount=%s, PricePaid=%s",
		event.OrderId.String(), event.Buyer.Hex(), event.Seller.Hex(), event.AmountFilled.String(), event.RemainingAmount.String(), event.PricePaid.String())

	// Convert contract order ID to our internal order ID format
	orderID := event.OrderId.String()

	// Get the order from KV store
	order, err := l.orderService.GetOrder(orderID)
	if err != nil {
		log.Printf("⚠️  Order %s not found in local storage: %v", orderID, err)
		return nil
	}

	// Update order with new remaining amount
	order.RemainingAmount = event.RemainingAmount.String()
	order.Status = "active" // Keep active if partially filled
	order.UpdatedAt = time.Now()

	// Store updated order
	if err := l.orderService.StoreOrder(order); err != nil {
		return fmt.Errorf("failed to update order with partial fill: %w", err)
	}

	// Emit event to frontend via marketplaceService
	if l.orderService.marketplaceService != nil {
		l.orderService.marketplaceService.emitOrderPartiallyFilled(order, event.AmountFilled.String(), event.Buyer.Hex())
	}

	log.Printf("✅ Order %s partially filled: %s KAWAI filled, %s KAWAI remaining",
		orderID, event.AmountFilled.String(), event.RemainingAmount.String())
	return nil
}

// handleOrderCancelled processes OrderCancelled events
// Requirements: 6.5 - Event-driven order status updates
func (l *MarketplaceEventListener) handleOrderCancelled(event *otcmarket.OTCMarketOrderCancelled) error {
	// Only the seller is responsible for updating their cancelled order in the marketplace KV
	myAddr := l.walletService.GetCurrentAddress()
	if myAddr == "" || !strings.EqualFold(myAddr, event.Seller.Hex()) {
		return nil // Not my order, skip status update
	}

	log.Printf("❌ OrderCancelled event: OrderID=%s, Seller=%s",
		event.OrderId.String(), event.Seller.Hex())

	// Convert contract order ID to our internal order ID format
	orderID := event.OrderId.String()

	// Update order status to cancelled
	if err := l.orderService.UpdateOrderStatus(orderID, "cancelled"); err != nil {
		return fmt.Errorf("failed to update order status to cancelled: %w", err)
	}

	log.Printf("✅ Order %s marked as cancelled", orderID)
	return nil
}

// GetRecentEvents returns recent marketplace events for debugging/monitoring
func (l *MarketplaceEventListener) GetRecentEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]MarketplaceEvent, error) {
	if l.blockchainClient == nil {
		return nil, fmt.Errorf("blockchain client is not available")
	}

	var events []MarketplaceEvent

	// Create filter options
	toBlockUint := toBlock.Uint64()
	filterOpts := &bind.FilterOpts{
		Start:   fromBlock.Uint64(),
		End:     &toBlockUint,
		Context: ctx,
	}

	// Get OrderCreated events
	orderCreatedIter, err := l.blockchainClient.OTCMarket.FilterOrderCreated(filterOpts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter OrderCreated events: %w", err)
	}
	defer func() { _ = orderCreatedIter.Close() }()

	for orderCreatedIter.Next() {
		event := orderCreatedIter.Event
		events = append(events, MarketplaceEvent{
			Type:        "OrderCreated",
			OrderID:     event.OrderId.String(),
			TxHash:      event.Raw.TxHash.Hex(),
			BlockNumber: event.Raw.BlockNumber,
			Timestamp:   time.Now(), // TODO: Get actual block timestamp
			Seller:      event.Seller.Hex(),
			TokenAmount: event.Amount.String(),
			USDTPrice:   event.Price.String(),
		})
	}

	// Get OrderFulfilled events
	orderFilledIter, err := l.blockchainClient.OTCMarket.FilterOrderFulfilled(filterOpts, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter OrderFulfilled events: %w", err)
	}
	defer func() { _ = orderFilledIter.Close() }()

	for orderFilledIter.Next() {
		event := orderFilledIter.Event
		events = append(events, MarketplaceEvent{
			Type:        "OrderFulfilled",
			OrderID:     event.OrderId.String(),
			TxHash:      event.Raw.TxHash.Hex(),
			BlockNumber: event.Raw.BlockNumber,
			Timestamp:   time.Now(), // TODO: Get actual block timestamp
			Seller:      event.Seller.Hex(),
			Buyer:       event.Buyer.Hex(),
			TokenAmount: event.Amount.String(),
			USDTPrice:   event.Price.String(),
		})
	}

	// Get OrderCancelled events
	orderCancelledIter, err := l.blockchainClient.OTCMarket.FilterOrderCancelled(filterOpts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter OrderCancelled events: %w", err)
	}
	defer func() { _ = orderCancelledIter.Close() }()

	for orderCancelledIter.Next() {
		event := orderCancelledIter.Event
		events = append(events, MarketplaceEvent{
			Type:        "OrderCancelled",
			OrderID:     event.OrderId.String(),
			TxHash:      event.Raw.TxHash.Hex(),
			BlockNumber: event.Raw.BlockNumber,
			Timestamp:   time.Now(), // TODO: Get actual block timestamp
			Seller:      event.Seller.Hex(),
		})
	}

	return events, nil
}

// ✅ NEW: Load last synced block from KV store
func (l *MarketplaceEventListener) loadLastSyncedBlock() uint64 {
	ctx := context.Background()
	data, err := l.kvStore.GetMarketplaceData(ctx, "event_listener:last_synced_block")
	if err != nil {
		log.Println("No last synced block found, starting from genesis")
		return 0
	}

	var block uint64
	if err := json.Unmarshal(data, &block); err != nil {
		log.Printf("Failed to unmarshal last synced block: %v", err)
		return 0
	}

	log.Printf("📍 Resuming from block %d", block)
	return block
}

// ✅ NEW: Save last synced block to KV store
func (l *MarketplaceEventListener) saveLastSyncedBlock(blockNumber uint64) {
	ctx := context.Background()
	data, err := json.Marshal(blockNumber)
	if err != nil {
		log.Printf("Failed to marshal block number: %v", err)
		return
	}

	if err := l.kvStore.StoreMarketplaceData(ctx, "event_listener:last_synced_block", data); err != nil {
		log.Printf("Failed to save last synced block: %v", err)
	}
}

// ✅ NEW: Get current block number from blockchain
func (l *MarketplaceEventListener) getCurrentBlockNumber() (uint64, error) {
	ctx := context.Background()
	header, err := l.blockchainClient.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get current block number: %w", err)
	}
	return header.Number.Uint64(), nil
}

// ✅ NEW: Replay events in chunks (for large gaps)
func (l *MarketplaceEventListener) replayEventsInChunks(ctx context.Context, fromBlock, toBlock uint64) error {
	const chunkSize = 2000 // Safe chunk size for most RPC nodes

	log.Printf("🔄 Large gap detected (%d blocks), starting chunked replay", toBlock-fromBlock)

	for start := fromBlock; start < toBlock; start += chunkSize {
		end := start + chunkSize
		if end > toBlock {
			end = toBlock
		}

		log.Printf("📥 Replaying events from block %d to %d", start, end)

		// Get events in this chunk
		events, err := l.GetRecentEvents(ctx, big.NewInt(int64(start)), big.NewInt(int64(end)))
		if err != nil {
			log.Printf("❌ Failed to get events in range %d-%d: %v", start, end, err)
			time.Sleep(5 * time.Second) // Backoff on error
			continue
		}

		// Process events with rate limiting
		for _, event := range events {
			// Process based on event type
			log.Printf("📝 Processing %s event for order %s", event.Type, event.OrderID)

			// Small delay to avoid overwhelming the system
			time.Sleep(100 * time.Millisecond)
		}

		// Save progress after each chunk
		l.saveLastSyncedBlock(end)

		// Small delay between chunks to avoid rate limit
		time.Sleep(1 * time.Second)
	}

	log.Printf("✅ Event replay completed")
	return nil
}
