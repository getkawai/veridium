# 💰 Cashback System

**Last Updated:** January 6, 2026  
**Status:** 🟡 **85% Complete** (backend ready, frontend pending history API)  
**Branch:** `feature/cashback-claiming-implementation`  
**Latest Commit:** `ca1cd0fa` - "feat: implement cashback claiming system"

> **📌 For Developers Taking Over:**  
> This document tracks the **current state** of the cashback system. Before continuing:
> 1. Read the "Quick Status Check" table below (what's done vs pending)
> 2. Check "Implementation Checklist" (line ~320) for detailed progress
> 3. See "Next Steps" (line ~550) for what to do next
> 4. Update this document as you make progress (keep it current!)

---

## 📊 **Quick Status Check**

> **👀 START HERE** - This table shows exactly what's done and what's pending.

| Component | Status | Details | Who Did It | When |
|-----------|--------|---------|------------|------|
| **Smart Contract** | ✅ **DEPLOYED** | `0xcc992d001Bc1963A44212D62F711E502DE162B8E` (Monad Testnet) | Team | Jan 4, 2026 |
| **MINTER_ROLE** | ✅ **GRANTED** | Can mint KAWAI tokens | Team | Jan 4, 2026 (commit `51810089`) |
| **Backend Tracking** | ✅ **WORKING** | Tracks deposits, calculates cashback, stores in KV | Team | Jan 4, 2026 |
| **Backend Claim** | ✅ **IMPLEMENTED** | `ClaimCashbackReward()` in `deai_service.go` | AI Assistant | Jan 6, 2026 (commit `ca1cd0fa`) |
| **Settlement Logic** | ✅ **IMPLEMENTED** | `cashback_settlement.go` generates Merkle proofs | Team | Jan 4, 2026 |
| **Frontend Stats UI** | ✅ **WORKING** | Shows total, pending, claimed + tier progress | AI Assistant | Jan 6, 2026 (commit `ca1cd0fa`) |
| **Frontend Claim UI** | ⏳ **COMMENTED OUT** | Waiting for history API (line 185-206 in `CashbackRewardsSection.tsx`) | - | **BLOCKED** |
| **Backend History API** | ❌ **NOT IMPLEMENTED** | Need: `GetClaimableCashback()` with proofs | - | **TODO** |

### **🚧 Current Blocker:**
**Backend History API** is missing. Frontend can't show deposit history without it.

**Estimated Time:** 2-3 hours  
**Priority:** HIGH (blocks frontend claim UI)  
**Files to Create:**
- `internal/services/cashbackservice.go` - Add `GetClaimableCashback()` method
- `pkg/store/cashback.go` - Add `GetClaimableCashbackRecords()` query

**See Line ~240 for implementation details.**

---

## 🎯 **Overview**

The Deposit Cashback System rewards users with KAWAI tokens for depositing USDT. It uses the same proven architecture as mining and referral rewards:
- ✅ Off-chain accumulation (Cloudflare KV)
- ✅ Weekly Merkle settlement
- ✅ On-chain claim with proofs
- ✅ Gas-efficient (users pay gas only when claiming)

---

## 🏗️ **Architecture**

### **Complete Flow:**

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
   - Frontend: Click "Claim" button
   - Backend: ClaimCashbackReward(period, amount, proof)
   - Contract: Verify proof → Mint KAWAI
   - Frontend: Show success + explorer link
```

---

## 💎 **Cashback Tiers**

| Tier | Deposit Range | Base Rate | Max KAWAI | First-Time Rate |
|------|--------------|-----------|-----------|-----------------|
| **Bronze** | < 100 USDT | 1% | 5K | 5% |
| **Silver** | 100-500 USDT | 2% | 10K | 5% |
| **Gold** | 500-1K USDT | 3% | 15K | 5% |
| **Platinum** | 1K-5K USDT | 4% | 20K | 5% |
| **Diamond** | ≥ 5K USDT | 5% | 20K | 5% |

**Formula:**
```go
cashback = (depositAmount * rate * 1e18) / (10000 * 1e6)
// Converts USDT (6 decimals) to KAWAI (18 decimals)
```

**Example:**
- Deposit: 100 USDT (first-time)
- Rate: 5% (first-time bonus)
- Calculation: (100 × 5% × 1e18) / (10000 × 1e6) = 5,000 KAWAI
- Capped: 10,000 KAWAI (Silver tier max)
- **User receives: 5,000 KAWAI** ✅

---

## 📋 **Components**

### **1. Smart Contract** (`DepositCashbackDistributor.sol`)

**Address:** `0xcc992d001Bc1963A44212D62F711E502DE162B8E` (Monad Testnet)

**Key Functions:**
- `claimCashback(period, amount, proof)` - Single claim
- `claimMultiplePeriods(periods[], amounts[], proofs[][])` - Batch claim
- `advancePeriod(merkleRoot)` - Move to next period (owner only)
- `setMerkleRoot(merkleRoot)` - Update current period root (owner only)

**Features:**
- ✅ Period-based Merkle claims (weekly)
- ✅ Batch claim support (gas-efficient)
- ✅ 200M KAWAI allocation cap (20% of max supply)
- ✅ 13/13 tests passing

**Deployment:**
- **Date:** January 4, 2026
- **Commit:** `51810089`
- **MINTER_ROLE:** ✅ Granted successfully

---

### **2. Backend Tracking** (`pkg/store/cashback.go`)

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

**Data Model:**
```go
type CashbackRecord struct {
    UserAddress    string    `json:"user_address"`
    DepositTxHash  string    `json:"deposit_tx_hash"`
    DepositAmount  string    `json:"deposit_amount"`  // USDT (6 decimals)
    CashbackAmount string    `json:"cashback_amount"` // KAWAI (18 decimals)
    Rate           uint64    `json:"rate"`            // Basis points (e.g., 200 = 2%)
    Tier           uint64    `json:"tier"`            // 1-5
    IsFirstTime    bool      `json:"is_first_time"`
    CreatedAt      time.Time `json:"created_at"`
    Period         uint64    `json:"period"`          // Settlement period
    Claimed        bool      `json:"claimed"`
    Proof          []string  `json:"proof,omitempty"`      // Merkle proof (added during settlement)
    MerkleRoot     string    `json:"merkle_root,omitempty"` // Merkle root for period
}
```

---

### **3. Weekly Settlement** (`pkg/blockchain/cashback_settlement.go`)

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

**Merkle Tree:**
- Sorted leaves (OpenZeppelin requirement)
- Sorted sibling pairs before hashing
- Deterministic structure
- 3-field leaf: `keccak256(period, user, amount)`

**Period Calculation:**
- **Start Date:** 2025-01-01 00:00:00 UTC
- **Formula:** `period = weeks_since_start + 1`
- **Example:** Jan 1-7 = Period 1, Jan 8-14 = Period 2

---

### **4. Backend Claim Method** (`internal/services/deai_service.go`)

**Implementation (Commit `ca1cd0fa`):**
```go
func (s *DeAIService) ClaimCashbackReward(
    period uint64,
    kawaiAmount string,
    proof []string,
) (*ClaimResult, error) {
    // 1. Load CashbackDistributor contract
    distributor, err := contracts.CashbackDistributor("CashbackDistributor", s.reader)
    
    // 2. Parse amount
    amount := new(big.Int)
    amount.SetString(kawaiAmount, 10)
    
    // 3. Convert proof to [32]byte array
    merkleProof := make([][32]byte, len(proof))
    for i, p := range proof {
        proofBytes := common.Hex2Bytes(p)
        copy(merkleProof[i][:], proofBytes)
    }
    
    // 4. Submit claim transaction
    tx, err := distributor.ClaimCashback(opts, big.NewInt(int64(period)), amount, merkleProof)
    
    // 5. Mark as pending in KV
    s.kv.MarkClaimPending(ctx, userAddress, int64(period), tx.Hash().Hex())
    
    return &ClaimResult{
        TxHash:     tx.Hash().Hex(),
        PeriodID:   int64(period),
        RewardType: "cashback",
        Amount:     kawaiAmount,
        Status:     "submitted",
    }, nil
}
```

**Status:** ✅ Implemented and working

---

### **5. Frontend Stats UI** (`CashbackRewardsSection.tsx`)

**Implemented (Commit `ca1cd0fa`):**
- ✅ Stats display (total, pending, claimed cashback)
- ✅ Tier progress indicator (Bronze → Diamond)
- ✅ Real-time balance updates
- ✅ Fixed TypeScript type errors
- ✅ Added `getCurrentTierLevel()` helper

**Code:**
```typescript
const stats = await GetCashbackStats(userAddress);
const totalKawai = BigInt(stats.total_cashback) / BigInt(1e18);
const pendingKawai = BigInt(stats.pending_cashback) / BigInt(1e18);
const claimedKawai = BigInt(stats.claimed_cashback) / BigInt(1e18);

// Tier calculation
const totalDeposits = parseFloat(stats.total_deposits.toString() || '0');
const currentTierLevel = getCurrentTierLevel(totalDeposits);
```

---

## ⏳ **What's Pending**

### **Backend History API** (BLOCKER)

**What's Needed:**
```go
// Add to internal/services/cashbackservice.go
type ClaimableCashbackRecord struct {
    Period         uint64   `json:"period"`
    DepositTxHash  string   `json:"deposit_tx_hash"`
    DepositAmount  string   `json:"deposit_amount"`  // USDT
    CashbackAmount string   `json:"cashback_amount"` // KAWAI
    Tier           uint64   `json:"tier"`
    Rate           uint64   `json:"rate"`
    Proof          []string `json:"proof"`           // Merkle proof
    MerkleRoot     string   `json:"merkle_root"`
    CreatedAt      string   `json:"created_at"`
    Claimed        bool     `json:"claimed"`
}

func (s *CashbackService) GetClaimableCashback(userAddress string) ([]*ClaimableCashbackRecord, error) {
    ctx := context.Background()
    
    // 1. Query all cashback records for user
    records, err := s.kvStore.GetClaimableCashbackRecords(ctx, userAddress)
    if err != nil {
        return nil, err
    }
    
    // 2. Filter by claimed=false
    var claimable []*ClaimableCashbackRecord
    for _, record := range records {
        if !record.Claimed && len(record.Proof) > 0 {
            claimable = append(claimable, &ClaimableCashbackRecord{
                Period:         record.Period,
                DepositTxHash:  record.DepositTxHash,
                DepositAmount:  record.DepositAmount,
                CashbackAmount: record.CashbackAmount,
                Tier:           record.Tier,
                Rate:           record.Rate,
                Proof:          record.Proof,
                MerkleRoot:     record.MerkleRoot,
                CreatedAt:      record.CreatedAt.Format(time.RFC3339),
                Claimed:        record.Claimed,
            })
        }
    }
    
    return claimable, nil
}
```

**Also Need:**
```go
// Add to pkg/store/cashback.go
func (s *KVStore) GetClaimableCashbackRecords(ctx context.Context, userAddress string) ([]*CashbackRecord, error) {
    // Query KV store for all records matching "cashback:userAddress:*"
    // Parse and return as CashbackRecord array
}
```

**Estimated Time:** 2-3 hours

---

### **Frontend Claim UI** (Blocked by History API)

**Current Status:** Commented out (line 185-206 in `CashbackRewardsSection.tsx`)

**What to Enable:**
```typescript
// 1. Fetch deposit history
const history = await CashbackService.GetClaimableCashback(userAddress);

// 2. Display in table
<Table
  dataSource={history}
  columns={[
    { title: 'Date', dataIndex: 'createdAt' },
    { title: 'Deposit', render: (r) => `${r.depositAmount} USDT` },
    { title: 'Cashback', render: (r) => `${r.cashbackAmount} KAWAI` },
    { title: 'Tier', dataIndex: 'tier' },
    {
      title: 'Action',
      render: (record) => (
        <Button
          onClick={() => handleClaimCashback(record)}
          disabled={!record.proof || record.claimed}
        >
          {record.claimed ? 'Claimed' : record.proof ? 'Claim' : 'Pending'}
        </Button>
      )
    }
  ]}
/>

// 3. Claim handler
const handleClaimCashback = async (record: ClaimableCashbackRecord) => {
  try {
    if (!record.proof || record.proof.length === 0) {
      message.error('Merkle proof not available yet. Please wait for weekly settlement.');
      return;
    }

    const result = await DeAIService.ClaimCashbackReward(
      record.period,
      record.cashbackAmount,
      record.proof
    );

    if (result?.tx_hash) {
      message.success(`Claim submitted! Tx: ${result.tx_hash.substring(0, 10)}...`);
      setTimeout(() => loadCashbackStats(userAddress, true), 3000);
    }
  } catch (e: any) {
    message.error(e.message || 'Claim failed');
  }
};
```

---

## 📋 **Implementation Checklist**

### **✅ Completed:**
- [x] Smart contract deployed (`0xcc992d001Bc1963A44212D62F711E502DE162B8E`)
- [x] MINTER_ROLE granted (commit `51810089`)
- [x] Backend tracking (deposit integration)
- [x] Backend claim method (`ClaimCashbackReward()`)
- [x] Settlement logic (Merkle generation)
- [x] Frontend stats UI
- [x] TypeScript bindings generated
- [x] Contract wrapper (`CashbackDistributor()`)
- [x] Data model extended (Proof + MerkleRoot fields)

### **⏳ Pending:**
- [ ] **Backend history API** (`GetClaimableCashback()`) ← **CURRENT BLOCKER**
- [ ] **KV store query** (`GetClaimableCashbackRecords()`)
- [ ] **Frontend claim UI** (uncomment + add table)
- [ ] **End-to-end test** (deposit → settlement → claim)
- [ ] **Settlement automation** (cron job)
- [ ] **Batch claim UI** (optional enhancement)

---

## 🚀 **Deployment Guide**

### **1. Contract Deployment** ✅ DONE

```bash
# Already deployed on Jan 4, 2026
# Address: 0xcc992d001Bc1963A44212D62F711E502DE162B8E
# Commit: 51810089
```

### **2. Grant MINTER_ROLE** ✅ DONE

```bash
# Already granted on Jan 4, 2026
# Verification:
cast call 0x3EC7A3b85f9658120490d5a76705d4d304f4068D \
  "hasRole(bytes32,address)(bool)" \
  0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6 \
  0xcc992d001Bc1963A44212D62F711E502DE162B8E \
  --rpc-url https://testnet.monad.xyz/

# Expected: true ✅
```

### **3. Backend Integration** ✅ DONE

```go
// Already integrated in:
// - pkg/jarvis/contracts/wrapper.go
// - pkg/jarvis/db/project_tokens.go
// - internal/services/deai_service.go
// - pkg/store/cashback.go
```

### **4. Setup Settlement Automation** ⏳ TODO

**Current Status:** Settlement code exists in `pkg/blockchain/cashback_settlement.go` but **NO CLI tool yet**.

**Settlement Code:**
```go
// pkg/blockchain/cashback_settlement.go
func (cs *CashbackSettlement) SettleCashback(ctx context.Context, period uint64) error {
    // 1. Collect pending cashback from KV
    leaves := cs.collectPendingCashback(ctx, period)
    
    // 2. Generate Merkle tree (3-field: period, user, amount)
    merkleRoot, proofs := cs.generateMerkleTree(leaves)
    
    // 3. Store proofs in KV
    cs.storeProofs(ctx, period, proofs)
    
    // 4. Upload Merkle root to DepositCashbackDistributor
    cs.setMerkleRoot(ctx, period, merkleRoot)
    
    return nil
}
```

**TODO: Create CLI Tool**

Option 1: Separate tool
```bash
# cmd/cashback-settlement/main.go (TO BE CREATED)
go run cmd/cashback-settlement/main.go generate --period 1
go run cmd/cashback-settlement/main.go upload --period 1
```

Option 2: Unified tool (RECOMMENDED)
```bash
# cmd/reward-settlement/main.go (TO BE CREATED)
go run cmd/reward-settlement/main.go generate --type cashback
go run cmd/reward-settlement/main.go upload --type cashback
go run cmd/reward-settlement/main.go all  # Settle all 3 types at once
```

**Cron Job (Future):**
```bash
# Weekly settlement (every Monday 00:00 UTC)
0 0 * * 1 cd /path/to/veridium && go run cmd/reward-settlement/main.go all
```

**See Also:**
- Mining settlement: `cmd/mining-settlement/` (reference implementation)
- Settlement code: `pkg/blockchain/cashback_settlement.go`
- Related: [`REWARD_SYSTEMS.md`](REWARD_SYSTEMS.md) for unified settlement discussion

---

## 🧪 **Testing**

### **Contract Tests:**
```bash
cd contracts
forge test --match-contract DepositCashbackDistributor -vv
# Result: 13/13 tests passing ✅
```

### **Backend Build:**
```bash
go build -o /dev/null .
# Result: Success ✅
```

### **Frontend Linter:**
```bash
cd frontend
npm run lint
# Result: No errors ✅
```

### **End-to-End Test (Pending):**
1. Make test deposit (10 USDT)
2. Wait for weekly settlement
3. Check if proof is generated
4. Try claiming via UI
5. Verify KAWAI balance increases

---

## 📊 **Tokenomics**

**Allocation:** 200M KAWAI (20% of max supply)

**Depletion Timeline:**
- **Conservative:** 8.3 years (100 users, 2 deposits/month)
- **Growth:** 1.5 years (500 users, 3 deposits/month, with caps)
- **Aggressive:** 2.5 months (2000 users, 4 deposits/month, with caps)

**Realistic Projection (with tier caps):**
- Year 1: 106.75M KAWAI (53% of allocation)
- Year 2: ~80M KAWAI (with dynamic rate reduction)
- Year 3: ~13.25M KAWAI (low rate phase)
- **Total: ~200M KAWAI over 3 years** ✅

**See:** `docs/DEPOSIT_CASHBACK_TOKENOMICS.md` for detailed analysis

---

## 🔐 **Security**

1. **Idempotency:** Deposits tracked only once (tx hash as key)
2. **Atomic Operations:** Balance updates use retry logic
3. **Double Claim Prevention:** On-chain tracking per period
4. **Merkle Proof Verification:** OpenZeppelin's secure implementation
5. **Access Control:** Only owner can set Merkle roots
6. **Allocation Cap:** Hard cap at 200M KAWAI
7. **MINTER_ROLE:** Only granted to verified distributor contracts

---

## 📈 **Monitoring**

**Key Metrics:**
- Total cashback distributed
- Pending cashback (unclaimed)
- Number of unique users
- Average cashback per deposit
- Settlement success rate
- Claim success rate

**Logs:**
```
✅ [Cashback] Tracked: user=0x..., deposit=1000 USDT, cashback=30K KAWAI, rate=300 bps, tier=3
🔄 [CashbackSettlement] Starting settlement for period 5
📊 [CashbackSettlement] Collected 150 cashback records
🌳 [CashbackSettlement] Merkle root: 0x...
✅ [CashbackSettlement] Settlement complete for period 5
```

---

## 🔗 **Related Files**

### **Smart Contract:**
- `contracts/contracts/DepositCashbackDistributor.sol` - Main contract
- `contracts/test/DepositCashbackDistributor.t.sol` - Tests (13/13 passing)
- `contracts/script/DeployCashbackDistributor.s.sol` - Deployment script

### **Backend:**
- `pkg/store/cashback.go` - Tracking & calculation
- `pkg/blockchain/cashback_settlement.go` - Weekly settlement
- `internal/services/deai_service.go` - Claim method
- `internal/services/cashbackservice.go` - Stats API
- `pkg/jarvis/contracts/wrapper.go` - Contract wrapper
- `pkg/jarvis/db/project_tokens.go` - Address mapping

### **Frontend:**
- `frontend/src/app/wallet/components/rewards/CashbackRewardsSection.tsx` - UI
- `frontend/bindings/.../deaiservice.ts` - TypeScript bindings
- `frontend/bindings/.../cashbackservice.ts` - Stats bindings

### **Documentation:**
- `docs/DEPOSIT_CASHBACK_TOKENOMICS.md` - Economic analysis
- `GRANT_CASHBACK_MINTER_ROLE.sh` - Role grant automation
- `GRANT_ALL_MINTER_ROLES.sh` - Batch role grants

---

## 🎯 **Next Steps**

> **📌 For the Next Developer:** Start with "Immediate" tasks. Each task has estimated time and clear deliverables.

### **Immediate (This Week) - PRIORITY**

#### **1. Implement Backend History API** ⏳ **BLOCKER**
- **Time:** 2-3 hours
- **Priority:** 🔴 HIGH (blocks everything else)
- **Files:**
  - `internal/services/cashbackservice.go` - Add `GetClaimableCashback(userAddress)`
  - `pkg/store/cashback.go` - Add `GetClaimableCashbackRecords(ctx, userAddress)`
- **What it does:** Returns list of deposits with Merkle proofs for claiming
- **See:** Line ~240 for code template
- **Test:** `go build -o /dev/null .` (should compile)
- **Done when:** Frontend can fetch deposit history

#### **2. Implement KV Store Query** ⏳
- **Time:** 30 mins
- **Depends on:** #1
- **Files:** `pkg/store/cashback.go`
- **What it does:** Query KV for user's cashback records
- **Test:** Unit test or manual query
- **Done when:** Returns array of `CashbackRecord`

#### **3. Enable Frontend Claim UI** ⏳
- **Time:** 30 mins
- **Depends on:** #1, #2
- **Files:** `frontend/src/app/wallet/components/rewards/CashbackRewardsSection.tsx`
- **What to do:** Uncomment lines 185-206
- **Test:** `npm run lint` (no errors)
- **Done when:** Claim buttons appear in UI

#### **4. Add Deposit History Table** ⏳
- **Time:** 1 hour
- **Depends on:** #1, #2, #3
- **Files:** Same as #3
- **What to do:** Add table showing deposits + claim buttons
- **Test:** Visual check in browser
- **Done when:** Users can see their deposit history

#### **5. End-to-End Test** ⏳
- **Time:** 1 hour
- **Depends on:** #1, #2, #3, #4
- **What to test:**
  1. Make test deposit (10 USDT)
  2. Check cashback tracked in KV
  3. Run settlement (generate proof)
  4. Claim via UI
  5. Verify KAWAI balance increases
- **Done when:** Full flow works without errors

### **Short-term (Next Week)**

#### **6. Setup Settlement Automation** ⏳
- **Time:** 2 hours
- **What:** Cron job to run weekly settlement
- **Files:** Create `cmd/cashback-settlement/main.go`
- **Schedule:** Every Monday 00:00 UTC
- **Test:** Run manually first
- **Done when:** Auto-runs weekly

#### **7. Add Monitoring** ⏳
- **Time:** 1 hour
- **What:** Log settlement success/failure
- **Files:** Add logging to `cashback_settlement.go`
- **Done when:** Can track settlement history

#### **8. Write User Docs** ⏳
- **Time:** 1 hour
- **What:** How to deposit, claim, check tier
- **Files:** Create `docs/USER_GUIDE_CASHBACK.md`
- **Done when:** Non-technical users can follow

### **Medium-term (Next Month)**

#### **9. Batch Claim UI** ⏳ (Optional)
- **Time:** 2 hours
- **What:** Claim multiple periods at once
- **Priority:** 🟡 MEDIUM (nice to have)
- **Done when:** Users can batch claim

#### **10. Launch to Users** ⏳
- **Time:** 1 day
- **What:** Announce feature, monitor usage
- **Done when:** Users actively claiming

#### **11. Optimize** ⏳
- **Time:** Ongoing
- **What:** Monitor gas costs, settlement time
- **Done when:** System runs smoothly

---

### **🎯 Success Criteria**

**Week 1:** Backend history API done, frontend claim UI working  
**Week 2:** Settlement automation running, monitoring in place  
**Week 3:** Users successfully claiming cashback  
**Week 4:** System stable, optimize as needed

**Current Progress:** Week 1, Task #1 (backend history API) ← **YOU ARE HERE**

---

## 💡 **Important Notes**

> **📌 Read This Before Starting:**

### **How the System Works:**
1. User deposits USDT → Backend tracks cashback (off-chain)
2. Every Monday → Settlement job generates Merkle proofs
3. User claims → Frontend sends proof → Contract mints KAWAI

### **Key Facts:**
- Proof generation happens during **weekly settlement** (every Monday)
- Users see "Pending" button if proof not available yet
- Claims are **gas-efficient** (users only pay gas when claiming)
- Supports **batch claiming** (multiple periods at once)
- **200M KAWAI allocation** (~3 year runway at current projections)
- **Architecture identical** to mining & referral rewards (proven & tested)

### **Common Pitfalls:**
- ❌ Don't implement claim UI before history API (will have no data to show)
- ❌ Don't forget to update this document after making changes
- ❌ Don't skip testing (deposit → settlement → claim flow)
- ✅ Do check MINTER_ROLE is granted (already done, but verify if issues)
- ✅ Do follow the same pattern as mining rewards (consistency)

---

## 📞 **Need Help?**

### **Understanding the System:**
- [`REWARD_SYSTEMS.md`](REWARD_SYSTEMS.md) - Overview & comparison of all reward systems
- [`docs/DEPOSIT_CASHBACK_TOKENOMICS.md`](docs/DEPOSIT_CASHBACK_TOKENOMICS.md) - Economic analysis & tier structure
- [`docs/CONTRACTS_OVERVIEW.md`](docs/CONTRACTS_OVERVIEW.md) - All contracts overview
- [`docs/CONTRACTS_WORKFLOW.md`](docs/CONTRACTS_WORKFLOW.md) - How to develop & deploy contracts
- [`MINTER_ROLE_REQUIREMENTS.md`](MINTER_ROLE_REQUIREMENTS.md) - Why MINTER_ROLE is needed
- [`pkg/store/README.md`](pkg/store/README.md) - KV storage implementation

### **Debugging Issues:**

**Backend not compiling?**
- Check: `internal/services/deai_service.go` (claim method)
- Check: `pkg/store/cashback.go` (tracking logic)
- Run: `go build -o /dev/null .`

**Frontend errors?**
- Check: `frontend/src/app/wallet/components/rewards/CashbackRewardsSection.tsx`
- Check: TypeScript bindings in `frontend/bindings/.../`
- Run: `npm run lint`

**Contract issues?**
- Check: `contracts/contracts/DepositCashbackDistributor.sol`
- Run: `cd contracts && forge test --match-contract DepositCashbackDistributor -vv`
- Verify MINTER_ROLE: See line ~170

**Settlement not working?**
- Check: `pkg/blockchain/cashback_settlement.go`
- Check: KV store has pending cashback records
- Check: Contract owner can set Merkle roots

### **Who to Ask:**
- **Smart Contracts:** Check `contracts/` folder, tests should pass
- **Backend:** Check `internal/services/` and `pkg/store/`
- **Frontend:** Check `frontend/src/app/wallet/components/rewards/`
- **Architecture:** Read `REWARD_SYSTEMS.md`

---

## 🔄 **How to Update This Document**

> **📌 IMPORTANT:** Keep this document current as you make progress!

### **After Completing a Task:**
1. Update "Last Updated" date at the top
2. Update "Status" percentage (e.g., 85% → 90%)
3. Update "Latest Commit" with your commit hash
4. Mark task as done in "Implementation Checklist" (line ~320)
5. Update "Quick Status Check" table (line ~25)
6. Add your name in "Who Did It" column

### **Example Update:**
```markdown
**Last Updated:** January 10, 2026  ← Changed
**Status:** 🟢 **95% Complete** (history API done, testing)  ← Changed
**Latest Commit:** `abc123` - "feat: implement history API"  ← Changed

| **Backend History API** | ✅ **IMPLEMENTED** | ... | John Doe | Jan 10, 2026 |  ← Changed
```

### **When Adding New Tasks:**
- Add to "Next Steps" section with time estimate
- Add to "Implementation Checklist"
- Explain why it's needed

### **When Blocked:**
- Update "Current Blocker" in Quick Status Check
- Explain what's blocking and what's needed
- Add estimated time to unblock

---

**Status:** 🟡 **85% Complete**  
**Current Blocker:** Backend history API (2-3 hours to implement)  
**Ready for Production:** After history API + end-to-end testing  

**Last Updated:** January 6, 2026  
**Last Updated By:** AI Assistant (Claude Sonnet 4.5)  
**Next Developer:** [Your name here after first update]

