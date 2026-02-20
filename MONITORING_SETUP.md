# Monitoring & Alerts - Production Ready ✅

**Status:** Ready to integrate  
**Location:** `x/alert/telegram.go`  
**Last Updated:** 2026-01-13

---

## 🎯 Overview

Kamu sudah punya Telegram alert system yang siap pakai! Tinggal test dan integrate ke critical points di aplikasi.

**Quick Start:**
```bash
# Test alert system (2 minutes)
make test-telegram-alert

# Check your Telegram for 4 test messages
```

---

## 📱 Telegram Alert System

### Current Implementation
```go
// x/alert/telegram.go
type TelegramAlert struct {
    BotToken string
    ChatID   string
    Client   *http.Client
}

// Usage:
alert := alert.NewTelegramAlert()
alert.SendAlert("ERROR", "Settlement", "Failed to upload merkle root")
```

### Alert Levels
- `ERROR` 🚨 - Critical issues (immediate action)
- `WARNING` ⚠️ - Important issues (review soon)
- `SUCCESS` ✅ - Successful operations
- `INFO` ℹ️ - General information

---

## 🔧 Integration Points

### 1. Settlement Monitoring

**File:** `cmd/reward-settlement/main.go`

Add alerts for settlement operations:

```go
import "github.com/kawai-network/x/alert"

func generateMiningSettlement(ctx context.Context, kv store.Store) error {
    alerter := alert.NewTelegramAlert()
    
    // Start notification
    alerter.SendAlert("INFO", "Settlement", "🔄 Starting mining settlement...")
    
    period, err := kv.GenerateMiningSettlement(ctx, "kawai")
    if err != nil {
        // Error notification
        alerter.SendAlert("ERROR", "Settlement", 
            fmt.Sprintf("❌ Mining settlement failed: %v", err))
        return err
    }
    
    // Success notification
    alerter.SendAlert("SUCCESS", "Settlement", 
        fmt.Sprintf("✅ Mining settlement complete!\nPeriod: %d\nContributors: %d\nTotal: %s KAWAI", 
            period.PeriodID, period.ContributorCount, period.TotalAmount))
    
    return nil
}

func uploadMiningRoot(ctx context.Context, kv store.Store) error {
    alerter := alert.NewTelegramAlert()
    
    // Upload notification
    alerter.SendAlert("INFO", "Settlement", "📤 Uploading mining merkle root...")
    
    // ... upload logic ...
    
    if receipt.Status != 1 {
        alerter.SendAlert("ERROR", "Settlement", 
            fmt.Sprintf("❌ Merkle root upload failed!\nTx: %s\nStatus: %d", 
                tx.Hash().Hex(), receipt.Status))
        return fmt.Errorf("transaction failed")
    }
    
    alerter.SendAlert("SUCCESS", "Settlement", 
        fmt.Sprintf("✅ Merkle root uploaded!\nTx: %s\nBlock: %d\nGas: %d", 
            tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed))
    
    return nil
}
```

---

### 2. Claim Monitoring

**File:** `internal/services/deai_service.go`

Add alerts for failed claims:

```go
import "github.com/kawai-network/x/alert"

func (s *DeAIService) ClaimCashbackReward(period uint64, kawaiAmount string, proof []string) (*ClaimResult, error) {
    alerter := alert.NewTelegramAlert()
    
    // ... existing claim logic ...
    
    // Alert on transaction failure
    if err != nil {
        // Only alert for non-user errors
        if !strings.Contains(err.Error(), "Already claimed") {
            alerter.SendAlert("WARNING", "Claim", 
                fmt.Sprintf("⚠️ Cashback claim failed\nUser: %s\nPeriod: %d\nError: %v", 
                    s.wallet.currentAccount.AddressHex(), period, err))
        }
        return nil, err
    }
    
    // Alert on transaction revert
    if receipt.Status != 1 {
        alerter.SendAlert("ERROR", "Claim", 
            fmt.Sprintf("❌ Cashback claim reverted!\nUser: %s\nPeriod: %d\nTx: %s", 
                s.wallet.currentAccount.AddressHex(), period, tx.Hash().Hex()))
        return nil, fmt.Errorf("transaction reverted")
    }
    
    return &ClaimResult{...}, nil
}
```

---

### 3. Backend Health Monitoring

**File:** `main.go` (or create `pkg/monitoring/health.go`)

Add periodic health checks:

```go
package main

import (
    "context"
    "time"
    "github.com/kawai-network/x/alert"
)

func startHealthMonitor() {
    alerter := alert.NewTelegramAlert()
    ticker := time.NewTicker(1 * time.Hour)
    
    go func() {
        for range ticker.C {
            if err := checkHealth(); err != nil {
                alerter.SendAlert("ERROR", "Health", 
                    fmt.Sprintf("❌ Health check failed: %v", err))
            }
        }
    }()
}

func checkHealth() error {
    // Check RPC connection
    client, err := ethclient.Dial(constant.MonadRpcUrl)
    if err != nil {
        return fmt.Errorf("RPC connection failed: %w", err)
    }
    defer client.Close()
    
    // Check KV connection
    kv, err := store.NewMultiNamespaceKVStore()
    if err != nil {
        return fmt.Errorf("KV connection failed: %w", err)
    }
    
    // Check contract balance
    // ... check logic ...
    
    return nil
}
```

---

### 4. Contract Balance Monitoring

**File:** `pkg/monitoring/balance.go` (create new)

```go
package monitoring

import (
    "context"
    "fmt"
    "math/big"
    "time"
    
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/kawai-network/x/constant"
    "github.com/kawai-network/contracts/kawaitoken"
    "github.com/kawai-network/x/alert"
)

// MonitorContractBalances checks if distributors have enough tokens
func MonitorContractBalances() {
    alerter := alert.NewTelegramAlert()
    ticker := time.NewTicker(6 * time.Hour) // Check every 6 hours
    
    go func() {
        for range ticker.C {
            ctx := context.Background()
            
            // Check mining distributor
            if err := checkDistributorBalance(ctx, 
                constant.MiningRewardDistributorAddr, 
                "Mining", 
                big.NewInt(10000), // Alert if <10k KAWAI
                alerter); err != nil {
                alerter.SendAlert("ERROR", "Balance Monitor", 
                    fmt.Sprintf("❌ Failed to check mining balance: %v", err))
            }
            
            // Check cashback distributor
            if err := checkDistributorBalance(ctx, 
                constant.CashbackDistributorAddress, 
                "Cashback", 
                big.NewInt(5000), // Alert if <5k KAWAI
                alerter); err != nil {
                alerter.SendAlert("ERROR", "Balance Monitor", 
                    fmt.Sprintf("❌ Failed to check cashback balance: %v", err))
            }
        }
    }()
}

func checkDistributorBalance(ctx context.Context, distributorAddr, name string, threshold *big.Int, alerter *alert.TelegramAlert) error {
    client, err := ethclient.Dial(constant.MonadRpcUrl)
    if err != nil {
        return err
    }
    defer client.Close()
    
    tokenAddr := common.HexToAddress(constant.KawaiTokenAddress)
    token, err := kawaitoken.NewKawaiToken(tokenAddr, client)
    if err != nil {
        return err
    }
    
    balance, err := token.BalanceOf(nil, common.HexToAddress(distributorAddr))
    if err != nil {
        return err
    }
    
    // Convert to KAWAI (divide by 1e18)
    balanceKAWAI := new(big.Int).Div(balance, big.NewInt(1e18))
    
    if balanceKAWAI.Cmp(threshold) < 0 {
        alerter.SendAlert("WARNING", "Balance Monitor", 
            fmt.Sprintf("⚠️ %s distributor balance LOW!\nCurrent: %s KAWAI\nThreshold: %s KAWAI\nContract: %s", 
                name, balanceKAWAI.String(), threshold.String(), distributorAddr))
    }
    
    return nil
}
```

---

### 5. Gas Price Monitoring

**File:** `pkg/monitoring/gas.go` (create new)

```go
package monitoring

import (
    "context"
    "fmt"
    "math/big"
    "time"
    
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/kawai-network/x/constant"
    "github.com/kawai-network/x/alert"
)

// MonitorGasPrice alerts when gas price is unusually high
func MonitorGasPrice() {
    alerter := alert.NewTelegramAlert()
    ticker := time.NewTicker(30 * time.Minute)
    
    go func() {
        for range ticker.C {
            ctx := context.Background()
            client, err := ethclient.Dial(constant.MonadRpcUrl)
            if err != nil {
                continue
            }
            
            gasPrice, err := client.SuggestGasPrice(ctx)
            client.Close()
            if err != nil {
                continue
            }
            
            // Alert if gas price > 500 gwei
            threshold := big.NewInt(500_000_000_000) // 500 gwei
            if gasPrice.Cmp(threshold) > 0 {
                gasPriceGwei := new(big.Int).Div(gasPrice, big.NewInt(1_000_000_000))
                alerter.SendAlert("WARNING", "Gas Monitor", 
                    fmt.Sprintf("⚠️ Gas price HIGH!\nCurrent: %s gwei\nThreshold: 500 gwei\n\nConsider delaying settlement.", 
                        gasPriceGwei.String()))
            }
        }
    }()
}
```

---

## 🚀 Quick Setup

### Step 1: Configure Telegram Bot (5 minutes)

Already configured in `.env`:
```bash
TELEGRAM_BOT_TOKEN=1593936639:AAE4_UTKemoc7Ib3T2hT9l02uKX2Slk1oIw
TELEGRAM_CHAT_ID=1360992240
```

✅ **Already done!** Your bot is ready to use.

---

### Step 2: Test Alert System (2 minutes)

Create test file: `cmd/dev/test-telegram-alert/main.go`

```go
package main

import (
    "fmt"
    "github.com/kawai-network/x/alert"
)

func main() {
    alerter := alert.NewTelegramAlert()
    
    // Test different alert levels
    fmt.Println("Sending test alerts...")
    
    alerter.SendAlert("INFO", "Test", "ℹ️ This is an info message")
    alerter.SendAlert("SUCCESS", "Test", "✅ This is a success message")
    alerter.SendAlert("WARNING", "Test", "⚠️ This is a warning message")
    alerter.SendAlert("ERROR", "Test", "🚨 This is an error message")
    
    fmt.Println("✅ Test alerts sent! Check your Telegram.")
}
```

Run:
```bash
go run cmd/dev/test-telegram-alert/main.go
```

---

### Step 3: Add Monitoring to Main (10 minutes)

**File:** `main.go`

```go
package main

import (
    "github.com/kawai-network/x/alert"
    "github.com/kawai-network/veridium/pkg/monitoring"
)

func main() {
    // ... existing setup ...
    
    // Send startup notification
    alerter := alert.NewTelegramAlert()
    alerter.SendAlert("SUCCESS", "Backend", "🚀 Kawai backend started successfully!")
    
    // Start monitoring (if in production)
    if os.Getenv("ENVIRONMENT") == "production" {
        monitoring.MonitorContractBalances()
        monitoring.MonitorGasPrice()
        // monitoring.MonitorHealth() // Optional
    }
    
    // ... rest of main ...
}
```

---

## 📊 Recommended Alert Configuration

### Critical Alerts (Immediate) 🚨
- Settlement failed
- Merkle root upload failed
- Contract out of tokens (<1000 KAWAI)
- Backend crashed
- RPC connection failed

### Warning Alerts (1 hour) ⚠️
- High claim failure rate (>5%)
- Low contract balance (<10k KAWAI)
- High gas prices (>500 gwei)
- Slow API response (>2s)

### Info Alerts (Daily) ℹ️
- Settlement completed
- Daily claim summary
- Gas cost report
- User activity stats

---

## 🎯 Production Checklist

### Telegram Setup ✅
- [x] Bot token configured
- [x] Chat ID configured
- [x] Alert system implemented
- [ ] Test alerts sent successfully
- [ ] Monitoring integrated

### Integration Points
- [ ] Settlement monitoring (mining & cashback)
- [ ] Claim failure alerts
- [ ] Contract balance monitoring
- [ ] Gas price monitoring
- [ ] Backend health checks
- [ ] Startup/shutdown notifications

### Testing
- [ ] Test all alert levels (INFO, SUCCESS, WARNING, ERROR)
- [ ] Test alert formatting (Markdown)
- [ ] Test concurrent alerts (no blocking)
- [ ] Verify alerts received in Telegram
- [ ] Test alert rate limiting (avoid spam)

---

## 💡 Best Practices

### 1. Don't Spam
```go
// BAD: Alert on every claim
alerter.SendAlert("INFO", "Claim", "User claimed 10 KAWAI")

// GOOD: Alert on failures only
if err != nil {
    alerter.SendAlert("ERROR", "Claim", fmt.Sprintf("Claim failed: %v", err))
}
```

### 2. Include Context
```go
// BAD: Vague message
alerter.SendAlert("ERROR", "Settlement", "Failed")

// GOOD: Detailed message
alerter.SendAlert("ERROR", "Settlement", 
    fmt.Sprintf("Mining settlement failed!\nPeriod: %d\nError: %v\nTime: %s", 
        period, err, time.Now().Format(time.RFC3339)))
```

### 3. Use Appropriate Levels
```go
// ERROR: Requires immediate action
alerter.SendAlert("ERROR", "Contract", "Out of tokens!")

// WARNING: Requires attention soon
alerter.SendAlert("WARNING", "Balance", "Low balance: 5000 KAWAI")

// INFO: Just for information
alerter.SendAlert("INFO", "Settlement", "Weekly settlement started")

// SUCCESS: Confirmation of success
alerter.SendAlert("SUCCESS", "Settlement", "Settlement completed!")
```

### 4. Rate Limiting
```go
// Avoid sending same alert repeatedly
var lastAlertTime time.Time
var lastAlertMessage string

func sendAlertOnce(alerter *alert.TelegramAlert, level, source, message string) {
    now := time.Now()
    if message == lastAlertMessage && now.Sub(lastAlertTime) < 1*time.Hour {
        return // Skip duplicate alert within 1 hour
    }
    
    alerter.SendAlert(level, source, message)
    lastAlertTime = now
    lastAlertMessage = message
}
```

---

## 🆘 Troubleshooting

### Alerts Not Received
1. Check bot token: `echo $TELEGRAM_BOT_TOKEN`
2. Check chat ID: `echo $TELEGRAM_CHAT_ID`
3. Test manually: `go run cmd/dev/test-telegram-alert/main.go`
4. Check bot is admin in group (if using group chat)
5. Check firewall allows HTTPS to api.telegram.org

### Too Many Alerts
1. Increase alert thresholds
2. Add rate limiting
3. Group similar alerts (daily summary)
4. Use WARNING instead of ERROR for non-critical

### Missing Alerts
1. Check goroutine not blocking
2. Verify error handling (don't fail silently)
3. Add logging: `slog.Info("Sending alert", "message", text)`
4. Test in development first

---

## ✅ Quick Start Checklist

- [x] Telegram bot configured (already done!)
- [ ] Test alert system: `go run cmd/dev/test-telegram-alert/main.go`
- [ ] Add settlement monitoring
- [ ] Add claim failure alerts
- [ ] Add balance monitoring
- [ ] Add gas price monitoring
- [ ] Test in production
- [ ] Verify alerts received
- [ ] Adjust thresholds as needed

---

**Your monitoring system is ready!** 🎉

Just integrate the alerts at critical points and you're good to go.

**Next:** Test the alert system, then integrate into settlement & claim flows.
