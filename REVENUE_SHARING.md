# Revenue Sharing System

**Complete Guide to USDT Dividend Distribution for KAWAI Holders**

---

## 📖 Overview

The Revenue Sharing system distributes **100% of platform revenue (USDT)** to KAWAI token holders proportionally based on their holdings. This creates a sustainable "Hold-to-Earn" model where token holders benefit directly from platform growth.

**Key Features:**
- ✅ Proportional distribution based on KAWAI holdings
- ✅ Weekly settlement via Merkle tree (gas-efficient)
- ✅ On-chain claiming with cryptographic proofs
- ✅ Automated holder scanning from blockchain
- ✅ Synchronized with other reward systems

---

## 🎯 Quick Start

### For Admins (Weekly Settlement)

```bash
# Revenue settlement only
make settle-revenue

# Or settle all reward types at once (recommended)
make settle-all
```

### For Users (Claiming Dividends)

1. Open Veridium app
2. Go to **Wallet → Rewards → Revenue Share**
3. View your KAWAI balance and share percentage
4. Click "Claim Dividends" when available
5. Sign transaction and receive USDT

---

## 🏗️ Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Revenue Sharing Flow                      │
└─────────────────────────────────────────────────────────────┘

1. Users deposit USDT → PaymentVault
2. Users pay for AI services (KV balance decreases)
3. Weekly: Admin runs settlement
   ├─ Scan all KAWAI holders from blockchain
   ├─ Calculate proportional dividends
   ├─ Generate Merkle tree with proofs
   ├─ Withdraw USDT from PaymentVault
   └─ Upload Merkle root to USDT_Distributor
4. Users claim USDT dividends via frontend
```

### Smart Contracts

| Contract | Address | Purpose |
|----------|---------|---------|
| **PaymentVault** | `0x...` | Holds user USDT deposits |
| **USDT_Distributor** | `0xE964B52D496F37749bd0caF287A356afdC10836C` | Distributes USDT dividends via Merkle proofs |
| **KawaiToken** | `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238` | ERC20 token for holdings calculation |

### Period System

All reward types (mining, cashback, referral, revenue) share the **same weekly period system**:

- **Period 1 starts**: January 1, 2025 (configurable via `CASHBACK_PERIOD_START`)
- **Period increments**: Every Monday 00:00 UTC
- **Settlement**: Always settles `currentPeriod - 1` (previous week)
- **Shared counter**: All types use `kvStore.GetCurrentPeriod()`

This ensures all rewards are synchronized and settled together.

---

## 💰 Economic Model

### Phase 1: Mining Era (Current)

**Revenue Sharing: ACTIVE** ✅

```
User Payment: 1 USDT per 1M tokens
├─ Contributor: 0 USDT (paid in KAWAI tokens via mining)
│  └─ 85-90% KAWAI mining rewards
└─ Platform: 1 USDT (100% revenue)
   └─ Distributed to KAWAI holders as dividend
```

**Key Points:**
- Users pay USDT from Day 1
- Contributors earn KAWAI tokens (mining rewards)
- Platform earns 100% USDT revenue
- All USDT distributed to KAWAI holders
- **Revenue sharing is ACTIVE in Phase 1**

### Phase 2: Post-Mining Era (Future)

**Revenue Sharing: ACTIVE** ✅

```
User Payment: 1 USDT per 1M tokens
├─ Contributor: 0.70 USDT (70% payment)
└─ Platform: 0.30 USDT (30% profit)
   └─ Distributed to KAWAI holders as dividend
```

**Key Points:**
- Mining stops (max supply reached)
- Contributors paid in USDT (70%)
- Platform earns 30% USDT profit
- 30% distributed to KAWAI holders
- **Revenue sharing continues in Phase 2**

**Phase Comparison:**

| Aspect | Phase 1 (Mining Era) | Phase 2 (Post-Mining) |
|--------|---------------------|----------------------|
| User Payment | 1 USDT per 1M tokens | 1 USDT per 1M tokens |
| Contributor Payment | KAWAI tokens (mining) | 0.70 USDT (70%) |
| Platform Revenue | 1 USDT (100%) | 0.30 USDT (30%) |
| Revenue Sharing | ✅ ACTIVE | ✅ ACTIVE |
| Dividend Amount | 100% of payments | 30% of payments |

### Revenue Calculation & Non-Refundable Deposits

**IMPORTANT: User Deposits Are Non-Refundable**

This is a critical aspect of the economic model that is often misunderstood:

**Smart Contract Verification:**
- `PaymentVault.sol` has NO user withdraw function
- Only `onlyOwner` can withdraw (for dividend distribution)
- Users cannot get refunds once deposited
- Deposits are permanent contributions to the platform

**How It Works:**

1. **User Deposits USDT:**
   - User deposits 1000 USDT to PaymentVault
   - Vault balance: 1000 USDT
   - User KV balance: 1000 USDT (off-chain credit tracking)

2. **User Spends Credits:**
   - User uses 100 USDT worth of AI services
   - Vault balance: 1000 USDT (unchanged)
   - User KV balance: 900 USDT (decreased)
   - Platform revenue: 100 USDT (spent amount)

3. **Settlement Time:**
   - Vault balance: 1000 USDT (all deposits)
   - Distributable: 1000 USDT ✅ (correct)
   - User keeps: 900 USDT KV credits for future AI usage
   - Token holders receive: 1000 USDT dividends

**Why This Is Correct:**

- User deposits are **non-refundable contributions** to the platform
- Users receive **AI service credits** in return, not refundable deposits
- Users can spend their KV credits on AI services indefinitely
- All USDT in vault legitimately belongs to the platform
- This model is intentional and verified in smart contract

**Example Scenario:**

```
Week 1:
- Alice deposits 1000 USDT → Vault: 1000, Alice KV: 1000
- Bob deposits 500 USDT → Vault: 1500, Bob KV: 500
- Alice spends 200 USDT on AI → Vault: 1500, Alice KV: 800
- Bob spends 100 USDT on AI → Vault: 1500, Bob KV: 400

Settlement:
- Vault balance: 1500 USDT
- Distributed to KAWAI holders: 1500 USDT ✅
- Alice still has: 800 USDT credits for future AI usage
- Bob still has: 400 USDT credits for future AI usage
- Platform revenue: 1500 USDT (all deposits are non-refundable)
```

**Current Approach:**
- All USDT in PaymentVault = platform revenue
- Query `USDT.balanceOf(PaymentVault)` at settlement time
- No per-request tracking needed
- Simple, efficient, and economically sound

**Rationale:**
- User deposits are **non-refundable** (verified in `PaymentVault.sol`)
- Only `onlyOwner` can withdraw (no user withdraw function)
- Users spend credits via off-chain KV balance tracking
- All USDT in vault = legitimate platform revenue

---

## 🔧 Technical Implementation

### 1. Holder Scanner

**File:** `pkg/blockchain/holder_scanner.go`

**HYBRID APPROACH (Registry + Blockchain Scan):**

To work around Monad testnet's strict 100-block RPC limit for `eth_getLogs`, we use a hybrid approach:

1. **Holder Registry (Primary Source):**
   - Desktop app auto-registers holders on wallet connect
   - CLI contributor auto-registers on wallet unlock
   - Stored in Cloudflare KV (dedicated `holderNamespaceID`)
   - Key format: `holder:{address}`
   - Tracks: address, lastSeen, source (desktop/cli), registered timestamp

2. **Recent Blockchain Scan (Safety Net):**
   - Scans last 90 blocks for Transfer events (under 100-block limit)
   - Catches new holders not yet in registry
   - Ensures no holder is missed

3. **Merge & Deduplicate:**
   - Combines registry + recent scan addresses
   - Queries current balance for each unique address
   - Filters out zero-balance holders

**Code Example:**

```go
// Get holders from registry
holderRegistry := NewHolderRegistry(kvStore)
registryAddresses, _ := holderRegistry.GetAllHolders(ctx)

// Scan recent blockchain (last 90 blocks)
scanner, _ := NewHolderScanner()
currentBlock, _ := client.BlockNumber(ctx)
startBlock := currentBlock - 90
recentHolders, _ := scanner.ScanHoldersFromBlock(ctx, startBlock)

// Merge and deduplicate
holderMap := make(map[common.Address]bool)
for _, addr := range registryAddresses {
    holderMap[addr] = true
}
for _, addr := range recentHolders {
    holderMap[addr] = true
}

// Query current balances
var holders []*KawaiHolder
for addr := range holderMap {
    balance, _ := scanner.GetBalance(ctx, addr)
    if balance.Cmp(big.NewInt(0)) > 0 {
        holders = append(holders, &KawaiHolder{
            Address: addr,
            Balance: balance,
        })
    }
}
```

**Benefits:**
- ✅ Works around RPC 100-block limit
- ✅ Scalable (registry grows with user base)
- ✅ No data loss (all active holders included)
- ✅ Automatic registration (no manual intervention)
- ✅ Mainnet-ready architecture

**Legacy Approach (Deprecated):**

The old approach scanned entire blockchain history from genesis:

```go
func (hs *HolderScanner) ScanHoldersLatest(ctx context.Context) ([]*KawaiHolder, error) {
    // Get all Transfer events with block range filter
    transferIterator, err := hs.kawaiToken.FilterTransfer(&bind.FilterOpts{
        Start: fromBlock,
        End:   toBlock,
    }, nil, nil)
    
    // Collect unique addresses
    addressSet := make(map[common.Address]bool)
    for transferIterator.Next() {
        event := transferIterator.Event
        addressSet[event.From] = true
        addressSet[event.To] = true
    }
    
    // Get current balance for each address
    holders := make([]*KawaiHolder, 0)
    for addr := range addressSet {
        balance, _ := hs.kawaiToken.BalanceOf(nil, addr)
        if balance.Cmp(big.NewInt(0)) > 0 {
            holders = append(holders, &KawaiHolder{
                Address: addr,
                Balance: balance,
            })
        }
    }
    
    return holders, nil
}
```

**Why Deprecated:**
- ❌ Blocked by Monad testnet RPC 100-block limit
- ❌ Not scalable for mainnet (millions of blocks)
- ❌ Slow and resource-intensive

**Features:**
- Block range filtering for efficiency
- Validates holder balances against total supply
- Tracks failed queries for visibility

### 2. Dividend Calculation

**Formula:**
```
holderDividend = (holderBalance / totalSupply) × totalRevenue
```

**Implementation:**
```go
func CalculateHolderShare(holderBalance, totalSupply, totalProfit *big.Int) *big.Int {
    // share = (holderBalance * totalProfit) / totalSupply
    share := new(big.Int).Mul(holderBalance, totalProfit)
    share.Div(share, totalSupply)
    return share
}
```

**Example:**
- Total KAWAI Supply: 1,000,000 tokens
- Holder Balance: 10,000 tokens (1%)
- Total Revenue: 100 USDT
- Holder Dividend: 1 USDT

### 3. Merkle Tree Generation

**File:** `pkg/blockchain/revenue_settlement.go`

Generates Merkle tree with 3-field leaves:

```go
// Merkle leaf structure (same as referral/cashback)
bytes32 leaf = keccak256(
    abi.encodePacked(
        index,       // uint256
        account,     // address
        amount       // uint256 (USDT)
    )
);
```

**Storage:**
- Proofs saved to KV store with `usdt:` prefix
- Distinguishes from KAWAI reward proofs
- Enables frontend to query correct proofs

### 4. Contract Funding

**File:** `pkg/blockchain/revenue_settlement.go`

Withdraws USDT from PaymentVault to USDT_Distributor:

```go
func (rs *RevenueSettlement) WithdrawToDistributor(ctx context.Context, amount *big.Int) error {
    // Get admin private key
    privateKey, err := crypto.HexToECDSA(constant.GetObfuscatedTemp())
    
    // Create transaction options
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
    
    // Call PaymentVault.withdraw(distributorAddr, amount)
    tx, err := rs.paymentVault.Withdraw(auth, rs.distributorAddress, amount)
    
    // Wait for confirmation
    receipt, err := bind.WaitMined(ctx, rs.client, tx)
    
    return nil
}
```

**Features:**
- Requires user confirmation before execution
- Waits for transaction confirmation
- Returns transaction hash for verification

### 5. Merkle Root Upload

**File:** `pkg/blockchain/revenue_settlement.go`

Uploads Merkle root to USDT_Distributor contract:

```go
func (rs *RevenueSettlement) UploadMerkleRoot(ctx context.Context, merkleRoot [32]byte) error {
    // Get admin private key
    privateKey, err := crypto.HexToECDSA(constant.GetObfuscatedTemp())
    
    // Create transaction options
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
    
    // Call USDT_Distributor.setMerkleRoot(merkleRoot)
    tx, err := rs.distributor.SetMerkleRoot(auth, merkleRoot)
    
    // Wait for confirmation
    receipt, err := bind.WaitMined(ctx, rs.client, tx)
    
    return nil
}
```

**Features:**
- Requires user confirmation before execution
- Waits for transaction confirmation
- Enables users to claim dividends after upload

### 6. Phase Detection

**File:** `pkg/config/phase.go`

Detects current phase based on token supply:

```go
func GetPhaseInfo(ctx context.Context) (*PhaseInfo, error) {
    // Get current total supply
    totalSupply, err := kawaiToken.TotalSupply(nil)
    
    // Max supply: 1 billion tokens
    maxSupply := new(big.Int).Mul(big.NewInt(1_000_000_000), big.NewInt(1e18))
    
    // Determine phase
    isPhase2 := totalSupply.Cmp(maxSupply) >= 0
    
    return &PhaseInfo{
        CurrentPhase: phase,
        TotalSupply:  totalSupply,
        MaxSupply:    maxSupply,
        IsPhase2:     isPhase2,
    }, nil
}
```

**Features:**
- Cached phase status (5-minute TTL)
- Force refresh capability
- Thread-safe with mutex

---

## 🚀 Settlement Workflow

### Weekly Settlement Process

```bash
make settle-revenue
```

**Step-by-Step:**

1. **Generate Settlement**
   - Query current period from KV store
   - Calculate settlement period (`currentPeriod - 1`)
   - Scan all KAWAI holders from blockchain
   - Calculate proportional dividends
   - Generate Merkle tree with proofs
   - Save proofs to KV store

2. **Get Vault Balance**
   - Query `USDT.balanceOf(PaymentVault)`
   - Display total revenue amount

3. **Withdraw USDT** (Interactive)
   - Display withdrawal amount
   - Request user confirmation
   - Execute `PaymentVault.withdraw()`
   - Wait for transaction confirmation

4. **Upload Merkle Root** (Interactive)
   - Display Merkle root
   - Request user confirmation
   - Execute `USDT_Distributor.setMerkleRoot()`
   - Wait for transaction confirmation

**Output Example:**
```
📊 Revenue Sharing Settlement (USDT Dividends)
─────────────────────────────

Step 1: Generating revenue settlement...

Current Period:    2
Settling Period:   1

✅ Settlement generated successfully
Merkle Root: 0xabcd...

Step 2: Getting vault balance...

Total Revenue: 1000000000 USDT

Step 3: Withdrawing USDT to distributor...

⚠️  About to withdraw 1000000000 USDT to USDT_Distributor
Continue with withdrawal? (y/n): y

✅ USDT withdrawn successfully

Step 4: Uploading merkle root...

⚠️  About to upload merkle root: 0xabcd...
Continue with upload? (y/n): y

✅ Merkle root uploaded successfully

✅ Revenue settlement completed!
```

---

## 👥 User Claiming Process

### Frontend Integration

**File:** `frontend/src/app/wallet/components/rewards/RevenueShareSection.tsx`

**Features:**
- Real-time KAWAI balance display
- Share percentage calculator
- Claimable USDT table with Merkle proofs
- Phase indicators (Phase 1 vs Phase 2)
- Gas estimation pre-claim
- Transaction confirmation

**User Flow:**

1. User opens Wallet → Rewards → Revenue Share
2. System displays:
   - KAWAI balance (from blockchain)
   - Total supply (from blockchain)
   - Share percentage: `(balance / supply) × 100`
   - Claimable USDT (from KV proofs)
3. User clicks "Claim Dividends"
4. System:
   - Fetches Merkle proof from KV
   - Estimates gas cost
   - Prompts wallet signature
5. User signs transaction
6. System:
   - Calls `USDT_Distributor.claim()`
   - Waits for confirmation
   - Updates UI with success message

---

## 🔍 Verification & Monitoring

### Check Settlement Status

```bash
# View all reward settlement status
make reward-settlement-status
```

### Verify Contract Balance

```bash
# Check USDT balance in distributor
cast call $USDT_DISTRIBUTOR "token()" --rpc-url $RPC_URL
cast call $USDT_TOKEN "balanceOf(address)" $USDT_DISTRIBUTOR --rpc-url $RPC_URL
```

### Check Claim Status

```bash
# Check if user has claimed for a period
make check-claim-status TYPE=revenue PERIOD=<period_id> ADDR=<user_address>
```

### Monitor Transactions

- **Explorer**: https://testnet.monad.xyz
- **Withdrawal TX**: Check PaymentVault events
- **Upload TX**: Check USDT_Distributor events
- **Claim TX**: Check user claim transactions

---

## 📊 API Reference

### Backend APIs

**File:** `internal/services/deai_service.go`

#### Get Revenue Share Stats

```go
GET /api/revenue-share/stats

Response:
{
  "kawai_balance": "10000000000000000000000",  // 10,000 KAWAI
  "total_supply": "1000000000000000000000000", // 1M KAWAI
  "share_percentage": "1.0",                   // 1%
  "phase": "Phase 1",
  "is_phase_2": false
}
```

#### Get Claimable Dividends

```go
GET /api/revenue-share/claimable?address=0x...

Response:
{
  "claimable_records": [
    {
      "period": 1,
      "amount": "1000000",  // 1 USDT (6 decimals)
      "proof": ["0x...", "0x..."],
      "merkle_root": "0x...",
      "claimed": false
    }
  ],
  "total_claimable": "1000000"
}
```

---

## 🛠️ Troubleshooting

### Configuration

**Holder Scan Start Block**

Located in `internal/constant/blockchain.go`:

```go
// HolderScanStartBlock: Starting block for holder scanning
// - Testnet: 0 (scan from genesis)
// - Mainnet: Set to token deployment block
const HolderScanStartBlock = 0
```

**For Mainnet:**
1. Find token deployment block number
2. Update `HolderScanStartBlock` in `internal/constant/blockchain.go`
3. Rebuild application
4. Reduces RPC calls and scan time significantly

### Common Issues

**1. "No period to settle yet"**
- **Cause**: Current period is 1, no previous period exists
- **Solution**: Wait until Period 2 (next Monday)

**2. "Failed to scan holders"**
- **Cause**: RPC connection issue or block range too large
- **Solution**: Check RPC_URL, reduce block range in scanner

**3. "Withdrawal failed"**
- **Cause**: Insufficient permissions or vault balance
- **Solution**: Verify admin key has owner role, check vault balance

**4. "Merkle root upload failed"**
- **Cause**: Insufficient gas or wrong contract address
- **Solution**: Check gas price, verify USDT_Distributor address

**5. "Claim failed - Proof invalid"**
- **Cause**: Merkle root not uploaded or proof mismatch
- **Solution**: Verify settlement completed, check proof in KV

---

## 📚 Related Documentation

- [`REWARD_SYSTEMS.md`](REWARD_SYSTEMS.md) - Overview of all reward systems
- [`MINING_SYSTEM.md`](MINING_SYSTEM.md) - Mining rewards implementation
- [`CASHBACK_SYSTEM.md`](CASHBACK_SYSTEM.md) - Cashback rewards implementation
- [`REFERRAL_SYSTEM.md`](REFERRAL_SYSTEM.md) - Referral rewards implementation
- [`docs/CONTRACTS_GUIDE.md`](docs/CONTRACTS_GUIDE.md) - Smart contract details
- [`TESTING_GUIDE.md`](TESTING_GUIDE.md) - Testing procedures

---

## 🎯 Best Practices

### For Admins

1. **Weekly Settlement**: Run `make settle-all` every Monday
2. **Verify Balance**: Check vault balance before withdrawal
3. **Monitor Claims**: Track claim transactions on explorer
4. **Backup Proofs**: KV store contains all Merkle proofs
5. **Gas Management**: Ensure admin wallet has sufficient MON

### For Developers

1. **Period Sync**: All reward types share same period system
2. **Proof Storage**: Use `usdt:` prefix for revenue proofs
3. **Error Handling**: Track failed queries and proof saves
4. **Block Range**: Use filtered scanning for efficiency
5. **Testing**: Test on testnet before production

### For Users

1. **Hold KAWAI**: More tokens = higher dividend share
2. **Check Weekly**: New dividends available every Monday
3. **Claim Regularly**: Unclaimed dividends accumulate
4. **Gas Costs**: Consider gas when claiming small amounts
5. **Verify TX**: Check explorer for claim confirmation

---

**Last Updated:** January 11, 2026  
**Version:** 1.0  
**Status:** Production Ready (95% - Testnet testing pending)
