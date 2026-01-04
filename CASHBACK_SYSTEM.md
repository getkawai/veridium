# Deposit Cashback System

## Overview

The Deposit Cashback System rewards users with KAWAI tokens for depositing USDT into the platform. It follows the same architecture as the contributor and referral reward systems, using off-chain accumulation and weekly Merkle-based claims for gas efficiency.

## Architecture

### Flow

```
1. User deposits USDT
   ↓
2. Backend calculates cashback (off-chain)
   - Tiered rates: 1%-5%
   - First-time bonus: 5%
   - Tier caps: 5K-20K KAWAI
   ↓
3. Store in KV: cashback:user:txhash
   ↓
4. Track user for period: cashback_period:N:users
   ↓
5. Weekly: Settlement job runs
   - Collect pending cashback
   - Generate Merkle tree
   - Store proofs in KV
   - Set Merkle root on-chain
   ↓
6. User claims via Merkle proof
   - Single period claim
   - Batch multi-period claim
```

## Components

### 1. Smart Contract (`DepositCashbackDistributor.sol`)

**Features:**
- Period-based Merkle claims (weekly)
- Batch claim support (multiple periods)
- 200M KAWAI allocation cap
- Gas-efficient (users only pay gas when claiming)
- Per-period Merkle roots

**Key Functions:**
- `claimCashback(period, amount, proof)` - Claim for single period
- `claimMultiplePeriods(periods[], amounts[], proofs[][])` - Batch claim
- `advancePeriod(merkleRoot)` - Move to next period (owner only)
- `setMerkleRoot(merkleRoot)` - Update current period root (owner only)

**Test Coverage:**
- ✅ 13/13 tests passing
- Single claim, batch claim, multi-period
- Access control, allocation cap
- Invalid proof rejection

### 2. Backend Tracking (`pkg/store/cashback.go`)

**Cashback Calculation:**

| Tier | Deposit Range | Base Rate | Max KAWAI | First-Time Rate |
|------|--------------|-----------|-----------|-----------------|
| 1 | < 100 USDT | 1% | 5K | 5% |
| 2 | 100-500 USDT | 2% | 10K | 5% |
| 3 | 500-1K USDT | 3% | 15K | 5% |
| 4 | 1K-5K USDT | 4% | 20K | 5% |
| 5 | ≥ 5K USDT | 5% | 20K | 5% |

**Formula:**
```go
cashback = (depositAmount * rate * 1e18) / (10000 * 1e6)
// Converts USDT (6 decimals) to KAWAI (18 decimals)
```

**Key Functions:**
- `CalculateCashback()` - Calculate KAWAI based on USDT deposit
- `TrackCashback()` - Store cashback record in KV
- `GetCashbackStats()` - Retrieve user statistics
- `GetCurrentPeriod()` - Calculate current weekly period

**KV Store Structure:**
```
cashback:user:txhash          → CashbackRecord (individual deposit)
cashback_stats:user           → CashbackStats (aggregated stats)
cashback_period:N:users       → []string (users with pending cashback)
cashback_proof:period:user    → [][]byte (Merkle proof for claim)
```

### 3. Deposit Integration (`internal/services/deposit_sync_service.go`)

**Integration Point:**
- Automatically tracks cashback on every USDT deposit
- Integrated into `SyncDeposit()` flow
- Non-blocking (doesn't fail deposit if cashback tracking fails)

**Flow:**
```go
func (s *DepositSyncService) SyncDeposit(req *SyncDepositRequest) {
    // 1. Verify deposit on-chain
    // 2. Update USDT balance (KV)
    // 3. Track cashback ← NEW
    // 4. Mark transaction as processed
    // 5. Return success
}
```

### 4. Weekly Settlement (`pkg/blockchain/cashback_settlement.go`)

**Settlement Process:**

```go
func SettleCashback(period uint64) error {
    // 1. Collect pending cashback for period
    leaves := collectPendingCashback(period)
    
    // 2. Generate Merkle tree (sorted leaves)
    merkleRoot, proofs := generateMerkleTree(leaves)
    
    // 3. Store proofs in KV for user claims
    storeProofs(period, proofs)
    
    // 4. Set Merkle root on-chain
    if period > currentPeriod {
        distributor.AdvancePeriod(merkleRoot)
    } else {
        distributor.SetMerkleRoot(merkleRoot)
    }
    
    return nil
}
```

**Merkle Tree Generation:**
- Sorted leaves (required by OpenZeppelin)
- Sorted sibling pairs before hashing
- Deterministic tree structure

**Key Functions:**
- `SettleCashback()` - Main settlement function
- `collectPendingCashback()` - Collect all pending records
- `generateMerkleTree()` - Build Merkle tree
- `storeProofs()` - Store proofs in KV
- `setMerkleRoot()` - Set root on-chain

### 5. Frontend Service (`internal/services/cashbackservice.go`)

**Wails-Exposed Functions:**
- `GetCashbackStats(userAddress)` - Get user stats
- `GetCurrentPeriod()` - Get current period

**Response Structure:**
```go
type CashbackStatsResponse struct {
    Total_Cashback   string // Total KAWAI earned (wei)
    Pending_Cashback string // Pending (unclaimed) (wei)
    Claimed_Cashback string // Already claimed (wei)
    Total_Deposits   uint64 // Number of deposits
    First_Deposit_At string // ISO timestamp
    Last_Deposit_At  string // ISO timestamp
}
```

## Deployment

### 1. Deploy Contract

```bash
make contracts-deploy-cashback-testnet
```

This will deploy `DepositCashbackDistributor` to Monad Testnet.

### 2. Grant Minter Role

```bash
# Set CASHBACK_DISTRIBUTOR_ADDRESS in contracts/.env
export CASHBACK_DISTRIBUTOR_ADDRESS=0x...

make contracts-grant-minter-cashback
```

### 3. Update Backend

Update `internal/constant/blockchain.go`:
```go
CashbackDistributorAddress = "0x..." // Deployed address
```

### 4. Run Settlement Job

```go
// Weekly cron job
settlement := blockchain.NewCashbackSettlement(kvStore, privateKey)
period := kvStore.GetCurrentPeriod()
if err := settlement.SettleCashback(ctx, period); err != nil {
    log.Printf("Settlement failed: %v", err)
}
```

## Period Calculation

**Start Date:** 2025-01-01 00:00:00 UTC

**Period Formula:**
```go
weeks := time.Since(startDate).Hours() / (24 * 7)
period := uint64(weeks) + 1 // Period 1, 2, 3, ...
```

**Example:**
- 2025-01-01 to 2025-01-07: Period 1
- 2025-01-08 to 2025-01-14: Period 2
- 2025-01-15 to 2025-01-21: Period 3

## Frontend Integration (Pending)

### Display Pending Cashback

```typescript
import { GetCashbackStats } from '@/bindings/...';

const stats = await GetCashbackStats(userAddress);
const pendingKawai = BigInt(stats.pending_cashback) / BigInt(1e18);
```

### Claim Cashback

```typescript
// 1. Get proof from backend
const proof = await fetch(`/api/cashback/proof/${period}/${userAddress}`);

// 2. Call contract
const distributor = new ethers.Contract(address, abi, signer);
const tx = await distributor.claimCashback(period, amount, proof);
await tx.wait();
```

### Batch Claim

```typescript
// Claim multiple periods at once
const periods = [1, 2, 3];
const amounts = [...]; // Amounts for each period
const proofs = [...];  // Proofs for each period

const tx = await distributor.claimMultiplePeriods(periods, amounts, proofs);
await tx.wait();
```

## Tokenomics

**Allocation:** 200M KAWAI (20% of max supply)

**Depletion Timeline:** ~3 years (based on tier caps and deposit projections)

**See:** `DEPOSIT_CASHBACK_TOKENOMICS.md` for detailed analysis.

## Consistency with Other Reward Systems

| Feature | Contributor | Referral | Cashback |
|---------|------------|----------|----------|
| Accumulation | Off-chain (KV) | Off-chain (KV) | Off-chain (KV) |
| Settlement | Weekly Merkle | Weekly Merkle | Weekly Merkle |
| Claim | Merkle proof | Merkle proof | Merkle proof |
| Batch Claim | ✅ | ✅ | ✅ |
| Gas Cost | User pays on claim | User pays on claim | User pays on claim |
| Allocation | 700M KAWAI | 50M KAWAI | 200M KAWAI |

## Security Considerations

1. **Idempotency:** Deposits tracked only once (tx hash as key)
2. **Atomic Operations:** Balance updates use retry logic
3. **Double Claim Prevention:** On-chain tracking per period
4. **Merkle Proof Verification:** OpenZeppelin's secure implementation
5. **Access Control:** Only owner can set Merkle roots
6. **Allocation Cap:** Hard cap at 200M KAWAI

## Monitoring

**Key Metrics:**
- Total cashback distributed
- Pending cashback (unclaimed)
- Number of unique users
- Average cashback per deposit
- Settlement success rate

**Logs:**
```
✅ [Cashback] Tracked: user=0x..., deposit=1000 USDT, cashback=30K KAWAI, rate=300 bps, tier=3
🔄 [CashbackSettlement] Starting settlement for period 5
📊 [CashbackSettlement] Collected 150 cashback records
🌳 [CashbackSettlement] Merkle root: 0x...
✅ [CashbackSettlement] Settlement complete for period 5
```

## Testing

**Contract Tests:**
```bash
cd contracts
forge test --match-contract DepositCashbackDistributor -vv
```

**Backend Build:**
```bash
go build -o /dev/null .
```

## Status

- ✅ Smart contract (13/13 tests passing)
- ✅ Backend tracking (deposit integration)
- ✅ Settlement job (Merkle generation)
- ✅ On-chain settlement (advance period)
- ✅ Build successful
- ⏳ Frontend (pending)
- ⏳ Deployment (pending)
- ⏳ Cron job setup (pending)

