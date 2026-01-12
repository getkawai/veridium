package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/dev/inject-cashback-data/main.go <user_address> [deposit_amount_usdt]")
		fmt.Println("Example: go run cmd/dev/inject-cashback-data/main.go 0x123... 1000")
		fmt.Println("")
		fmt.Println("Default: 3 deposits (100, 500, 1000 USDT)")
		os.Exit(1)
	}

	userAddress := os.Args[1]

	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	ctx := context.Background()

	fmt.Println("💰 Injecting cashback data for:", userAddress)
	fmt.Println()

	// Get current period and inject into previous period (for settlement)
	currentPeriod := kv.GetCurrentPeriod()
	period := currentPeriod - 1 // Inject into previous period so settlement can process it
	fmt.Printf("📅 Current period: %d\n", currentPeriod)
	fmt.Printf("📅 Injecting into period: %d (for settlement)\n", period)
	fmt.Println()

	// If custom amount provided, use it
	if len(os.Args) >= 3 {
		var depositUSDT int64
		fmt.Sscanf(os.Args[2], "%d", &depositUSDT)

		if err := injectDeposit(ctx, kv, userAddress, depositUSDT, period, 1); err != nil {
			log.Fatal("Failed to inject deposit:", err)
		}

		fmt.Println()
		fmt.Println("✅ Cashback data injected!")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Run settlement: make settle-cashback")
		fmt.Println("  2. Upload Merkle root: go run cmd/reward-settlement/main.go upload --type cashback")
		fmt.Println("  3. Test claiming via UI")
		return
	}

	// Default: 3 deposits with different tiers
	deposits := []struct {
		amount int64
		desc   string
	}{
		{100, "Tier 2 (100 USDT)"},
		{500, "Tier 3 (500 USDT)"},
		{1000, "Tier 4 (1000 USDT)"},
	}

	for i, deposit := range deposits {
		if err := injectDeposit(ctx, kv, userAddress, deposit.amount, period, i+1); err != nil {
			log.Fatal("Failed to inject deposit:", err)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println()
	fmt.Println("✅ Cashback data injected!")
	fmt.Println()
	fmt.Println("📊 Summary:")
	stats, err := kv.GetCashbackStats(ctx, userAddress)
	if err != nil {
		log.Printf("Warning: Failed to get stats: %v", err)
	} else {
		totalKAWAI := new(big.Int)
		totalKAWAI.SetString(stats.TotalCashback, 10)

		// Convert to human-readable (18 decimals)
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		wholePart := new(big.Int).Div(totalKAWAI, divisor)

		fmt.Printf("   Total Deposits: %d\n", stats.TotalDeposits)
		fmt.Printf("   Total Cashback: ~%s KAWAI\n", wholePart.String())
		fmt.Printf("   Pending: %s wei\n", stats.PendingCashback)
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run settlement: make settle-cashback")
	fmt.Println("  2. Upload Merkle root: go run cmd/reward-settlement/main.go upload --type cashback")
	fmt.Println("  3. Test claiming via UI")
}

func injectDeposit(ctx context.Context, kv *store.KVStore, userAddress string, depositUSDT int64, period uint64, depositNum int) error {
	// Convert USDT to wei (6 decimals)
	depositAmount := big.NewInt(depositUSDT)
	depositAmount.Mul(depositAmount, big.NewInt(1_000_000)) // 1e6

	// Generate fake tx hash
	txHash := fmt.Sprintf("0x%064d", time.Now().Unix()+int64(depositNum))

	// Track cashback
	if err := kv.TrackCashback(ctx, userAddress, txHash, depositAmount, period); err != nil {
		return fmt.Errorf("failed to track cashback: %w", err)
	}

	fmt.Printf("✅ Deposit %d: %d USDT\n", depositNum, depositUSDT)

	return nil
}
