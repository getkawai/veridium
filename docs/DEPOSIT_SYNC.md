# 💰 Deposit Synchronization System

## Overview

Veridium uses **user-triggered sync** for syncing user deposits from the blockchain to the off-chain KV store.

After depositing USDT via wallet, the user client receives a transaction hash and automatically calls `SyncDeposit(txHash)` to update the off-chain balance.

---

## 🎯 Why User-Triggered Sync?

### **Benefits:**
- ✅ **Simple** - No background processes, no WebSocket connections
- ✅ **Reliable** - User has the transaction hash from wallet response
- ✅ **Secure** - Full verification before updating balance
- ✅ **Automatic** - Frontend auto-triggers sync after deposit
- ✅ **Duplicate-safe** - Transaction can only be synced once

### **Why Not Event Listener?**
- ❌ Requires WebSocket connection (Monad RPC is HTTP-only)
- ❌ Redundant (user already has tx hash from wallet)
- ❌ More complex (background process, reconnection logic)
- ❌ May miss events (network issues, connection drops)

---

## 🔄 How It Works

### **Sync Flow**

```
1. User deposits USDT → PaymentVault (on-chain)
   ↓
2. User gets transaction hash from wallet
   ↓
3. User clicks "Sync Deposit" in UI
   ↓
4. Frontend calls: SyncDeposit(txHash, userAddress)
   ↓
5. Backend (Go) verifies transaction:
   - ✅ Transaction exists and succeeded
   - ✅ Transaction contains "Deposited" event
   - ✅ Event user address matches request
   - ✅ Transaction not already processed
   ↓
6. Backend calls KVStore.AddBalance(user, amount)
   ↓
7. Backend marks transaction as processed (prevent duplicates)
   ↓
8. ✅ User balance updated
```

**Implementation:**
- **File:** `internal/services/deposit_sync_service.go`
- **Exposed to frontend:** Via Wails bindings
- **Method:** `SyncDeposit(ctx, SyncDepositRequest) -> SyncDepositResponse`
- **Triggered by:** User client (automatically after wallet deposit confirmation)

---

## 🛡️ Security Features

### **1. Transaction Verification**
```go
// Verify transaction exists and succeeded
receipt, err := client.TransactionReceipt(ctx, txHash)
if receipt.Status != 1 {
    return error("Transaction failed")
}
```

### **2. Event Validation**
```go
// Parse Deposited event from logs
event, err := vault.ParseDeposited(log)
if err != nil {
    return error("No deposit event found")
}
```

### **3. User Address Matching**
```go
// Verify user address matches event
if event.User != requestedUser {
    return error("User address mismatch")
}
```

### **4. Duplicate Prevention**
```go
// Check if transaction already processed
processedKey := fmt.Sprintf("processed_tx:%s", txHash)
if exists(processedKey) {
    return error("Already synced")
}
```

---

## 📝 API Reference

### **DepositSyncService.SyncDeposit**

**Request:**
```typescript
interface SyncDepositRequest {
  txHash: string;      // Transaction hash from wallet
  userAddress: string; // User's wallet address
}
```

**Response:**
```typescript
interface SyncDepositResponse {
  success: boolean;
  message: string;
  amount?: string;       // Deposited amount (USDT)
  newBalance?: string;   // New balance after deposit
  blockNumber?: number;  // Block number of transaction
  alreadySync?: boolean; // True if already synced
}
```

**Example Usage (Frontend):**
```typescript
import { SyncDeposit } from '@/bindings/github.com/kawai-network/veridium/internal/services';

// User deposits via wallet
const tx = await wallet.deposit(amount);

// ✅ Wallet returns tx hash (confirmation that deposit succeeded)
if (tx.hash) {
  // Automatically sync to off-chain balance
  const result = await SyncDeposit({
    txHash: tx.hash,
    userAddress: wallet.address
  });

  if (result.success) {
    console.log(`✅ Deposit synced: ${result.amount} USDT`);
    console.log(`New balance: ${result.newBalance} USDT`);
    // User can now use LLM services immediately
  } else {
    console.error(`❌ Sync failed: ${result.message}`);
    // Show retry button to user
  }
}
```

---

## 🔧 Configuration

### **Contract Addresses**
Defined in `internal/constant/blockchain.go`:

```go
const (
    MonadRpcUrl         = "https://testnet-rpc.monad.xyz"
    PaymentVaultAddress = "0x714238F32A7aE70C0D208D58Cc041D8Dda28e813"
)
```

### **Initialization**
In `internal/app/context.go`:

```go
func (ctx *Context) InitBlockchainClient() {
    // ... blockchain client init ...
    
    // Initialize deposit sync service
    syncService, _ := services.NewDepositSyncService(ctx.KVStore)
    ctx.DepositSyncService = syncService
    ctx.AddCleanup(func() { syncService.Close() })
}
```

---

## 🧪 Testing

### **Test Deposit Sync**
1. Start desktop app
2. Deposit USDT via wallet
3. Wallet returns transaction hash
4. Frontend auto-calls `SyncDeposit(txHash)`
5. Check logs for: `✅ [DepositSync] Deposit synced`
6. Verify balance updated in UI

### **Test Duplicate Prevention**
1. Deposit USDT and sync successfully
2. Try to sync the same transaction again
3. Should receive: `"Deposit already synced"`
4. Balance should not change

### **Test Error Handling**
1. Try to sync with invalid tx hash → `"Transaction not found"`
2. Try to sync failed transaction → `"Transaction failed on blockchain"`
3. Try to sync wrong user address → `"User address mismatch"`

---

## 🐛 Troubleshooting

### **Sync Fails**
- **"Transaction not found"** → Wait for confirmation (1-2 blocks)
- **"Transaction failed"** → Check blockchain explorer
- **"No deposit event found"** → Verify transaction is to PaymentVault
- **"User address mismatch"** → Check wallet address
- **"Already synced"** → Deposit already processed (check balance)

---

## 📊 Monitoring

### **Sync Logs**
```
💰 [DepositSync] Sync request: txHash=0x123..., user=0xABC...
✅ [DepositSync] Deposit synced: user=0xABC..., amount=100 USDT, block=456
```

### **Duplicate Prevention Logs**
```
⚠️  [DepositSync] Transaction already processed: 0x123...
```

---

## 🚀 Future Improvements

- [ ] **Batch sync** - Sync multiple transactions at once
- [ ] **Auto-detect pending deposits** - Scan blockchain for user's deposits
- [ ] **Reconciliation service** - Periodic sync of all deposits
- [ ] **Webhook notifications** - Notify user when deposit is synced
- [ ] **Multi-chain support** - Support deposits from other chains

---

## 📚 Related Documentation

- [Economic Model](README.md#-economic-model)
- [PaymentVault Contract](contracts/PaymentVault.sol)
- [KV Store Package](pkg/store/README.md)
- [Blockchain Client](pkg/blockchain/)

