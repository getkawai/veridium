package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/x/store"
)

// ReconciliationService handles periodic reconciliation between blockchain and KV store
// This ensures data consistency even if events are missed or KV store becomes out of sync
type ReconciliationService struct {
	blockchainClient *blockchain.Client
	orderService     *OrderService
	walletService    *WalletService
	kvStore          *store.KVStore

	// Configuration
	interval time.Duration
}

// NewReconciliationService creates a new reconciliation service
func NewReconciliationService(
	blockchainClient *blockchain.Client,
	orderService *OrderService,
	walletService *WalletService,
	kvStore *store.KVStore,
) *ReconciliationService {
	return &ReconciliationService{
		blockchainClient: blockchainClient,
		orderService:     orderService,
		walletService:    walletService,
		kvStore:          kvStore,
		interval:         5 * time.Minute, // Reconcile every 5 minutes
	}
}

// Start begins the periodic reconciliation process
func (s *ReconciliationService) Start(ctx context.Context) {
	log.Println("🔄 Starting reconciliation service...")
	go s.run(ctx)
	log.Println("✅ Reconciliation service started")
}

// run is the main reconciliation loop
func (s *ReconciliationService) run(ctx context.Context) {
	// Run immediately on start
	if err := s.Reconcile(ctx); err != nil {
		log.Printf("❌ Initial reconciliation failed: %v", err)
	}

	// Then run periodically
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Reconciliation service stopping (context cancelled)...")
			return
		case <-ticker.C:
			if err := s.Reconcile(ctx); err != nil {
				log.Printf("❌ Reconciliation failed: %v", err)
			}
		}
	}
}

// Reconcile performs a full reconciliation between blockchain and KV store
func (s *ReconciliationService) Reconcile(ctx context.Context) error {
	myAddr := s.walletService.GetCurrentAddress()

	if myAddr == "" {
		// No wallet connected, skip reconciliation
		return nil
	}

	log.Println("🔄 Starting reconciliation...")
	startTime := time.Now()

	// 1. Get my orders from blockchain
	blockchainOrders, err := s.getOrdersFromBlockchain(ctx, myAddr)
	if err != nil {
		return fmt.Errorf("failed to get orders from blockchain: %w", err)
	}

	// 2. Get my orders from KV
	kvOrders, err := s.getOrdersFromKV(myAddr)
	if err != nil {
		return fmt.Errorf("failed to get orders from KV: %w", err)
	}

	// 3. Create maps for easy lookup
	blockchainOrderMap := make(map[string]*Order)
	for i := range blockchainOrders {
		blockchainOrderMap[blockchainOrders[i].ID] = &blockchainOrders[i]
	}

	kvOrderMap := make(map[string]*Order)
	for i := range kvOrders {
		kvOrderMap[kvOrders[i].ID] = &kvOrders[i]
	}

	// 4. Reconcile differences
	stats := s.reconcileDifferences(blockchainOrderMap, kvOrderMap)

	duration := time.Since(startTime)
	log.Printf("✅ Reconciliation completed in %v: %d orders checked, %d added, %d updated, %d removed",
		duration, stats.Checked, stats.Added, stats.Updated, stats.Removed)

	return nil
}

// ReconciliationStats tracks reconciliation statistics
type ReconciliationStats struct {
	Checked int
	Added   int
	Updated int
	Removed int
}

// getOrdersFromBlockchain fetches all orders for a user from blockchain
func (s *ReconciliationService) getOrdersFromBlockchain(ctx context.Context, userAddr string) ([]Order, error) {
	addr := common.HexToAddress(userAddr)

	// Get orders from blockchain (paginated)
	var allOrders []Order
	offset := big.NewInt(0)
	limit := big.NewInt(100)

	for {
		blockchainOrders, err := s.blockchainClient.MarketplaceGetOrdersBySeller(ctx, addr, offset, limit)
		if err != nil {
			return nil, err
		}

		if len(blockchainOrders) == 0 {
			break
		}

		// Convert to our Order struct
		for _, bcOrder := range blockchainOrders {
			order := Order{
				ID:              bcOrder.Id.String(),
				Seller:          bcOrder.Seller.Hex(),
				TokenAmount:     bcOrder.TokenAmount.String(),
				USDTPrice:       bcOrder.PriceInUSDT.String(),
				RemainingAmount: bcOrder.RemainingAmount.String(),
				Status:          s.convertBlockchainStatus(bcOrder.IsActive, bcOrder.RemainingAmount),
				UpdatedAt:       time.Now(),
			}

			// Calculate price per token
			if tokenAmount, ok := new(big.Int).SetString(order.TokenAmount, 10); ok {
				if usdtPrice, ok := new(big.Int).SetString(order.USDTPrice, 10); ok {
					pricePerToken := new(big.Int).Div(usdtPrice, tokenAmount)
					order.PricePerToken = pricePerToken.String()
				}
			}

			allOrders = append(allOrders, order)
		}

		// Move to next page
		offset.Add(offset, limit)

		// Safety check: don't fetch more than 1000 orders
		if offset.Cmp(big.NewInt(1000)) > 0 {
			break
		}
	}

	return allOrders, nil
}

// getOrdersFromKV fetches all orders for a user from KV store
func (s *ReconciliationService) getOrdersFromKV(userAddr string) ([]Order, error) {
	orders, err := s.orderService.GetUserOrders(userAddr)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// convertBlockchainStatus converts blockchain order state to our status string
func (s *ReconciliationService) convertBlockchainStatus(isActive bool, remainingAmount *big.Int) string {
	if !isActive {
		// Check if it was filled or cancelled
		if remainingAmount.Cmp(big.NewInt(0)) == 0 {
			return "filled"
		}
		return "cancelled"
	}
	return "active"
}

// reconcileDifferences compares blockchain and KV orders and fixes mismatches
func (s *ReconciliationService) reconcileDifferences(
	blockchainOrders map[string]*Order,
	kvOrders map[string]*Order,
) ReconciliationStats {
	stats := ReconciliationStats{}

	// Check all blockchain orders
	for orderID, bcOrder := range blockchainOrders {
		stats.Checked++

		kvOrder, existsInKV := kvOrders[orderID]

		if !existsInKV {
			// Order exists on blockchain but not in KV - add it
			log.Printf("🔧 Adding missing order %s to KV", orderID)
			if err := s.orderService.StoreOrder(bcOrder); err != nil {
				log.Printf("❌ Failed to add order %s: %v", orderID, err)
			} else {
				stats.Added++
			}
			continue
		}

		// Order exists in both - check for mismatches
		if s.hasOrderMismatch(bcOrder, kvOrder) {
			log.Printf("🔧 Fixing order %s: KV(status=%s, remaining=%s) → Blockchain(status=%s, remaining=%s)",
				orderID,
				kvOrder.Status, kvOrder.RemainingAmount,
				bcOrder.Status, bcOrder.RemainingAmount)

			// Blockchain is truth - update KV
			kvOrder.Status = bcOrder.Status
			kvOrder.RemainingAmount = bcOrder.RemainingAmount
			kvOrder.UpdatedAt = time.Now()

			if err := s.orderService.StoreOrder(kvOrder); err != nil {
				log.Printf("❌ Failed to update order %s: %v", orderID, err)
			} else {
				stats.Updated++
			}
		}
	}

	// Check for orders in KV that don't exist on blockchain
	// (This shouldn't happen normally, but could if blockchain was reset)
	for orderID := range kvOrders {
		if _, existsOnBlockchain := blockchainOrders[orderID]; !existsOnBlockchain {
			log.Printf("⚠️  Order %s exists in KV but not on blockchain (keeping in KV)", orderID)
			// We keep it in KV for now, as it might be a local-only order
			// In production, you might want to mark it as "orphaned" or remove it
		}
	}

	return stats
}

// hasOrderMismatch checks if there's a mismatch between blockchain and KV order
func (s *ReconciliationService) hasOrderMismatch(bcOrder, kvOrder *Order) bool {
	// Check status mismatch
	if bcOrder.Status != kvOrder.Status {
		return true
	}

	// Check remaining amount mismatch
	if bcOrder.RemainingAmount != kvOrder.RemainingAmount {
		return true
	}

	return false
}

// ReconcileNow triggers an immediate reconciliation (for manual trigger)
func (s *ReconciliationService) ReconcileNow(ctx context.Context) error {
	log.Println("🔄 Manual reconciliation triggered...")
	return s.Reconcile(ctx)
}
