# Cara Kerja Teknis: User Deposit USDC di Kawai Desktop

## ⚠️ PENTING: Network Requirement

**USDC deposit HARUS menggunakan Monad Network!**

- **Mainnet:** USDC contract `0x754704bc059f8c67012fed69bc8a327a5aafb603` (Monad Mainnet)
- **Testnet:** MockUSDT contract `0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc` (Monad Testnet)

User **TIDAK BISA** langsung deposit USDC dari exchange (Binance, Coinbase, dll) karena:
1. Exchange belum support Monad Network withdrawal
2. USDC di network lain (Ethereum, BSC, dll) adalah contract yang berbeda
3. Kirim ke network yang salah = **dana hilang**

**Solusi:** User harus bridge USDC dari network lain ke Monad dulu. Lihat [Deposit from Exchange Guide](../docs-users/user-guide/deposit-from-exchange.md) untuk panduan lengkap.

---

## 📋 Overview

User deposit USDC (atau MockUSDT di testnet) ke PaymentVault smart contract untuk mendapatkan service credits yang bisa digunakan untuk AI services.

---

## 🔄 Flow Diagram

```
User Wallet (USDC)
    ↓
[1] Check Allowance
    ↓
[2] Approve (if needed) → USDC Contract
    ↓
[3] Deposit → PaymentVault Contract
    ↓
[4] Off-chain Balance Update → Cloudflare KV
    ↓
User can use AI services
```

---

## 🔧 Technical Flow (Step by Step)

### Frontend (`frontend/src/app/wallet/wallet.tsx`)

```typescript
// User clicks "Deposit" button
const handleDeposit = async () => {
  // 1. Convert amount to raw format (6 decimals for USDC)
  const rawAmount = Math.floor(amount * 1_000_000).toString();
  
  // 2. Call backend deposit function
  const txHash = await DeAIService.DepositToVault(rawAmount);
  
  // 3. Show success message
  message.success(`Deposit Successful! TX: ${txHash}`);
}
```

**Input:**
- User enters amount: `100` USDC
- Converted to raw: `100000000` (100 * 10^6)

---

### Backend (`internal/services/deai_service.go`)

#### Step 1: Validate & Prepare

```go
func (s *DeAIService) DepositToVault(amountStr string) (string, error) {
    // Check wallet is unlocked
    if s.wallet.currentAccount == nil {
        return "", fmt.Errorf("wallet is locked")
    }
    
    // Convert amount string to big.Int
    amount := new(big.Int)
    amount.SetString(amountStr, 10) // "100000000" → big.Int
    
    // Resolve addresses (auto-switches testnet/mainnet)
    stablecoinAddr := common.HexToAddress(constant.UsdtTokenAddress)
    // Testnet: 0x3AE...1CCc (MockUSDT)
    // Mainnet: 0x754...603 (USDC)
    
    vaultAddr, _ := contracts.ResolveAddress("PaymentVault")
}
```

#### Step 2: Check Allowance

```go
// Check if PaymentVault can spend user's USDC
allowance, err := s.GetUSDTAllowance(
    s.wallet.currentAccount.AddressHex(), 
    "PaymentVault"
)

allowanceBig := new(big.Int)
allowanceBig.SetString(allowance, 10)
```

**Why?** ERC-20 tokens require approval before contract can transfer them.

#### Step 3: Approve (if needed)

```go
if allowanceBig.Cmp(amount) < 0 {
    // Allowance insufficient, need approval
    
    // Load USDC contract
    stablecoin, _ := binding.KawaiToken(stablecoinAddr.Hex(), s.reader)
    
    // Approve PaymentVault to spend amount
    tx, _ := stablecoin.Approve(opts, vaultAddr, amount)
    
    // Wait for approval transaction to be mined
    receipt, _ := bind.WaitMined(ctx, s.reader.Client(), tx)
    
    if receipt.Status == 0 {
        return "", fmt.Errorf("approval transaction failed")
    }
}
```

**Blockchain Transaction 1 (if needed):**
```
From: User Address
To: USDC Contract
Function: approve(spender, amount)
  - spender: PaymentVault address
  - amount: 100000000 (100 USDC)
Gas: ~50,000
```

#### Step 4: Deposit to Vault

```go
// Load PaymentVault contract
vault, _ := binding.Vault("PaymentVault", s.reader)

// Call deposit function
tx, _ := vault.Deposit(opts, amount)

return tx.Hash().Hex(), nil
```

**Blockchain Transaction 2:**
```
From: User Address
To: PaymentVault Contract
Function: deposit(amount)
  - amount: 100000000 (100 USDC)
Gas: ~150,000
```

---

### Smart Contract (`contracts/contracts/PaymentVault.sol`)

```solidity
contract PaymentVault {
    IERC20 public immutable stablecoin; // USDC or MockUSDT
    
    function deposit(uint256 _amount) external nonReentrant {
        require(_amount > 0, "Amount must be > 0");
        
        // Transfer USDC from user to vault
        stablecoin.safeTransferFrom(msg.sender, address(this), _amount);
        
        // Emit event for off-chain tracking
        emit Deposited(msg.sender, _amount);
    }
}
```

**What Happens:**
1. ✅ Validates amount > 0
2. ✅ Transfers USDC from user wallet to PaymentVault
3. ✅ Emits `Deposited` event with user address and amount
4. ✅ User's USDC now in vault, ready for off-chain tracking

---

### Off-Chain Balance Tracking

**After deposit transaction is confirmed:**

1. **Backend monitors blockchain events** (or user triggers balance update)
2. **Cloudflare KV Store** updated:
   ```json
   {
     "wallet_address": "0x123...",
     "usdt_balance": "100000000",  // Raw amount
     "formatted_balance": "100.00", // Display amount
     "last_updated": "2026-01-21T..."
   }
   ```

3. **User can now use AI services:**
   - Each API call deducts from KV balance
   - No blockchain transaction needed for usage
   - Fast and gas-free

---

## 💰 Example: Deposit 100 USDC

### User Action
```
User has: 500 USDC in wallet
User wants: 100 USDC service credits
```

### Transaction 1: Approve (if first time)
```
Function: USDC.approve(PaymentVault, 100 USDC)
Gas Cost: ~0.0005 MON (~$0.01)
Time: ~2 seconds
Result: PaymentVault can now spend 100 USDC
```

### Transaction 2: Deposit
```
Function: PaymentVault.deposit(100 USDC)
Gas Cost: ~0.0015 MON (~$0.03)
Time: ~2 seconds
Result: 
  - User wallet: 400 USDC (500 - 100)
  - PaymentVault: +100 USDC
  - User KV balance: 100 USDC credits
```

### After Deposit
```
User can use AI services:
- Chat: -0.01 USDC per message
- Image: -0.05 USDC per image
- No gas fees for usage
- Instant deduction from KV balance
```

---

## 🔐 Security Features

### Smart Contract Level
1. ✅ **ReentrancyGuard** - Prevents reentrancy attacks
2. ✅ **SafeERC20** - Safe token transfers
3. ✅ **Ownable** - Only owner can withdraw
4. ✅ **Immutable stablecoin** - Cannot change token address

### Backend Level
1. ✅ **Wallet unlock check** - Ensures user authorized
2. ✅ **Amount validation** - Prevents invalid amounts
3. ✅ **Transaction confirmation** - Waits for mining
4. ✅ **Error handling** - Graceful failure handling

### Off-Chain Level
1. ✅ **Event monitoring** - Tracks all deposits
2. ✅ **Balance reconciliation** - Matches on-chain state
3. ✅ **Audit trail** - All transactions logged

---

## 🌐 Network Differences

### Testnet (Monad Testnet)
```
Stablecoin: MockUSDT
Address: 0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc
Features: Has mint() function for testing
Gas: Very cheap (testnet MON)
```

### Mainnet (Monad Mainnet)
```
Stablecoin: USDC (Circle)
Address: 0x754704bc059f8c67012fed69bc8a327a5aafb603
Features: Real USDC, no mint()
Gas: Real cost (mainnet MON)
```

**Code automatically switches based on environment!**

---

## 📊 Data Flow

### On-Chain (Blockchain)
```
User Wallet → USDC Contract → PaymentVault Contract
```
- **Permanent** - Cannot be reversed
- **Transparent** - Anyone can verify
- **Costs gas** - User pays transaction fees

### Off-Chain (Cloudflare KV)
```
PaymentVault Event → Backend Monitor → KV Store → User Balance
```
- **Fast** - No blockchain delay
- **Flexible** - Can be updated instantly
- **Gas-free** - No transaction costs for usage

---

## 🔄 Balance Lifecycle

### 1. Deposit (On-Chain)
```
User deposits 100 USDC → PaymentVault
```

### 2. Credit (Off-Chain)
```
KV Store: user_balance = 100 USDC
```

### 3. Usage (Off-Chain)
```
User uses AI service:
- Chat message: -0.01 USDC
- KV Store: user_balance = 99.99 USDC
```

### 4. Revenue (On-Chain)
```
Weekly settlement:
- Owner withdraws from PaymentVault
- Distributes to KAWAI holders
```

---

## 🚨 Error Handling

### Common Errors

1. **"Wallet is locked"**
   - User needs to unlock wallet first
   - Solution: Enter password

2. **"Insufficient balance"**
   - User doesn't have enough USDC
   - Solution: Get more USDC first

3. **"Approval failed"**
   - Approval transaction reverted
   - Solution: Check gas, try again

4. **"Deposit failed"**
   - Deposit transaction reverted
   - Solution: Check allowance, balance, gas

### Transaction Failures

If approval succeeds but deposit fails:
- ✅ Approval is still valid
- ✅ Can retry deposit without re-approving
- ✅ No USDC lost

---

## 💡 Key Points

1. **Two Transactions (First Time)**
   - Approve: Allow vault to spend USDC
   - Deposit: Transfer USDC to vault

2. **One Transaction (Subsequent)**
   - Only deposit (if allowance sufficient)

3. **Automatic Network Detection**
   - Code switches between MockUSDT/USDC
   - No manual configuration needed

4. **Off-Chain Balance**
   - Fast usage without gas fees
   - Tracked in Cloudflare KV
   - Reconciled with on-chain state

5. **Non-Refundable**
   - Deposits are service credits
   - Cannot withdraw back to wallet
   - Used for AI services only

---

## 📝 Code References

### Backend
- `internal/services/deai_service.go` - DepositToVault()
- `pkg/jarvis/contracts/wrapper.go` - Contract wrappers
- `internal/constant/blockchain.go` - Addresses

### Frontend
- `frontend/src/app/wallet/wallet.tsx` - Deposit UI
- `frontend/src/config/network.ts` - Network config

### Smart Contracts
- `contracts/contracts/PaymentVault.sol` - Vault contract
- `contracts/contracts/MockUSDT.sol` - Test token

---

## 🎯 Summary

**User deposits USDC → PaymentVault smart contract → Off-chain balance tracking → User can use AI services**

**Benefits:**
- ✅ Secure (smart contract + audited)
- ✅ Fast (off-chain usage)
- ✅ Gas-efficient (only 1-2 transactions)
- ✅ Transparent (on-chain verification)
- ✅ Automatic (network detection)

**User Experience:**
1. Click "Deposit"
2. Enter amount
3. Approve (if first time)
4. Confirm deposit
5. Wait ~4 seconds
6. Start using AI services! 🎉
