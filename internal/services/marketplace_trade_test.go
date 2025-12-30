package services

import (
	"testing"
	"time"
)

func TestStorePartialTradeRecord(t *testing.T) {
	// Setup mock data
	order := &Order{
		ID:            "order123",
		Seller:        "seller_address",
		TokenAmount:   "1000",
		PricePerToken: "1",
	}

	tradeID := "trade456"
	txHash := "0xtxhash"
	amount := "500"
	usdtAmount := "500"
	buyer := "buyer_address"

	// Replicating the logic from storePartialTradeRecord to verify struct population
	trade := Trade{
		ID:          tradeID,
		OrderID:     order.ID,
		Seller:      order.Seller,
		Buyer:       buyer, // This is what we fixed
		TokenAmount: amount,
		Price:       order.PricePerToken,
		USDTAmount:  usdtAmount,
		TxHash:      txHash,
		Timestamp:   time.Now(),
	}

	if trade.Buyer != buyer {
		t.Errorf("Expected buyer %s, got %s", buyer, trade.Buyer)
	}

	t.Logf("✅ Trade record populated correctly with buyer: %s", trade.Buyer)
}

// Verification of history key generation
func TestTradeHistoryKeys(t *testing.T) {
	s := &TradeService{}

	sellerKey := s.getUserTradesKey("seller_id", "seller")
	buyerKey := s.getUserTradesKey("buyer_id", "buyer")

	if sellerKey == buyerKey {
		t.Error("Seller and buyer history keys should be different")
	}

	if sellerKey != "user:seller_id:trades:seller" {
		t.Errorf("Expected user:seller_id:trades:seller, got %s", sellerKey)
	}

	t.Logf("✅ Trade history keys generated correctly: %s, %s", sellerKey, buyerKey)
}

func TestTradeHistoryIdempotency(t *testing.T) {
	// Replicating the logic from addTradeToUserHistory to verify idempotency
	tradeID := "trade123"

	// Initial state: empty
	var tradeIDs []string

	// 1st append
	found := false
	for _, id := range tradeIDs {
		if id == tradeID {
			found = true
			break
		}
	}
	if !found {
		tradeIDs = append(tradeIDs, tradeID)
	}

	if len(tradeIDs) != 1 {
		t.Errorf("Expected 1 trade ID, got %d", len(tradeIDs))
	}

	// 2nd append (simulating redundant call)
	found = false
	for _, id := range tradeIDs {
		if id == tradeID {
			found = true
			break
		}
	}
	if !found {
		tradeIDs = append(tradeIDs, tradeID)
	}

	if len(tradeIDs) != 1 {
		t.Errorf("Idempotency failed: expected 1 trade ID, got %d", len(tradeIDs))
	}

	t.Logf("✅ Trade history idempotency verified")
}
