# Performance Analysis: Rewards Tab Loading

**Date:** Jan 6, 2026  
**Issue:** Slow initial load when Rewards tab is clicked (~20+ seconds)  
**Status:** 🔴 Performance Bottleneck Identified

---

## 🐌 Problem Summary

When the Rewards tab is clicked for the first time, users experience a **20+ second loading delay**. This is caused by:

1. **Sequential KV API calls** in cashback loading (up to 104 calls)
2. **Multiple parallel service calls** (Mining, Cashback, Referral)
3. **High latency** per Cloudflare KV API call (~100-300ms each)

---

## 📊 Detailed Analysis

### 1. Mining Rewards Loading

**File:** `internal/services/deai_service.go:539`  
**Function:** `GetClaimableRewards()`

**Process Flow:**
```
DeAIService.GetClaimableRewards()
  → s.kv.GetClaimableRewards(ctx, userAddr)
    → ListMerkleProofs(ctx, address)
      → Cloudflare KV: ListWorkersKVKeys()        // API call #1
      → For each proof key:
          → Cloudflare KV: GetWorkersKV()         // API calls #2, #3, #4...
    → GetContributor(ctx, address)
      → Cloudflare KV: GetWorkersKV()             // API call #N
```

**Performance:**
- **API Calls:** 1 LIST + N GET (where N = number of periods with proofs)
- **Example:** User with 10 periods = **11 API calls**
- **Latency:** 11 × 200ms = **~2.2 seconds**

**Impact:** ⚠️ Moderate (acceptable for now)

---

### 2. Cashback Rewards Loading ⚠️ **CRITICAL BOTTLENECK**

**File:** `pkg/store/cashback.go:327`  
**Function:** `GetClaimableCashbackRecords()`

**Process Flow:**
```go
GetClaimableCashbackRecords(ctx, userAddress)
  → currentPeriod := s.GetCurrentPeriod()         // e.g., 52
  
  // ⚠️ SEQUENTIAL LOOP - BOTTLENECK HERE
  for period := uint64(1); period < currentPeriod; period++ {
    // Check merkle root
    rootKey := fmt.Sprintf("cashback_period:%d:merkle_root", period)
    rootData, err := s.GetCashbackData(ctx, rootKey)  // API call #1, #2, #3...
    
    // Check user proof
    proofKey := fmt.Sprintf("cashback_proof:%d:%s", period, userAddress)
    proofData, err := s.GetCashbackData(ctx, proofKey) // API call #53, #54, #55...
  }
```

**Performance:**
- **API Calls:** 2 × currentPeriod (worst case)
- **Example:** Period 52 = **2 × 52 = 104 API calls** 😱
- **Latency:** 104 × 200ms = **~20.8 seconds**

**Impact:** 🔴 **CRITICAL** - This is the main bottleneck!

**Why Sequential?**
- Current implementation loops from period 1 to current period
- Each iteration waits for 2 API calls to complete before moving to next
- No parallelization or caching

---

### 3. Referral Rewards Loading

**File:** `internal/services/referralservice.go:53`  
**Function:** `GetReferralStats()`

**Process Flow:**
```
ReferralService.GetReferralStats(userAddress)
  → s.kvStore.GetReferralCodeByAddress(ctx, userAddress)
    → Cloudflare KV: GetWorkersKV()               // API call #1
    → s.GetReferralData(ctx, code)
      → Cloudflare KV: GetWorkersKV()             // API call #2
```

**Performance:**
- **API Calls:** 2 GET calls
- **Latency:** 2 × 200ms = **~400ms**

**Impact:** ✅ Fast (acceptable)

**Error Case:** New users without referral codes
- Error: `"no referral code for this address: get: 'key not found' (10009)"`
- This is **expected behavior** for new users
- UX could be improved (see UX Issues section)

---

## 🔥 Root Cause: Sequential Period Scanning

**Location:** `pkg/store/cashback.go:327-358`

```go
// ⚠️ BOTTLENECK: Sequential loop through all periods
for period := uint64(1); period < currentPeriod; period++ {
    // 2 API calls per iteration
    // No parallelization
    // No early exit
    // No caching
}
```

**Why This is Bad:**
1. **Linear growth:** As weeks pass, loading time increases linearly
2. **Wasted calls:** Checks periods that have no data
3. **No optimization:** Every page load repeats the same scan
4. **Poor UX:** Users wait 20+ seconds for data that could load in 2-3 seconds

---

## 💡 Solutions (Prioritized)

### **Solution A: Parallel Loading** ⭐ QUICK WIN

**Complexity:** Low  
**Time to Implement:** ~2 hours  
**Performance Gain:** 20s → 2-3s (10x faster)

**Implementation:**
```go
func (s *KVStore) GetClaimableCashbackRecords(ctx context.Context, userAddress string) ([]CashbackRecord, error) {
    currentPeriod := s.GetCurrentPeriod()
    
    // Use goroutines for parallel fetching
    var wg sync.WaitGroup
    resultsChan := make(chan CashbackRecord, currentPeriod)
    errorsChan := make(chan error, currentPeriod)
    
    // Parallel fetch with worker pool (limit concurrency)
    semaphore := make(chan struct{}, 10) // Max 10 concurrent requests
    
    for period := uint64(1); period < currentPeriod; period++ {
        wg.Add(1)
        go func(p uint64) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release
            
            // Check root + proof in parallel
            if record, err := s.fetchCashbackForPeriod(ctx, p, userAddress); err == nil {
                resultsChan <- record
            }
        }(period)
    }
    
    // Wait and collect
    go func() {
        wg.Wait()
        close(resultsChan)
        close(errorsChan)
    }()
    
    // Collect results
    var records []CashbackRecord
    for record := range resultsChan {
        records = append(records, record)
    }
    
    return records, nil
}
```

**Pros:**
- ✅ Easy to implement
- ✅ Immediate performance improvement
- ✅ No schema changes required
- ✅ Backward compatible

**Cons:**
- ⚠️ Still makes many API calls (but in parallel)
- ⚠️ Cloudflare rate limits may apply

---

### **Solution B: Settled Periods Index** ⭐⭐ BEST LONG-TERM

**Complexity:** Medium  
**Time to Implement:** ~4 hours  
**Performance Gain:** 20s → <1s (20x+ faster)

**Implementation:**
```go
// Store settled periods explicitly
// Key: "cashback_settled_periods"
// Value: [1, 2, 5, 10, 15, ...] (only periods with data)

func (s *KVStore) GetClaimableCashbackRecords(ctx context.Context, userAddress string) ([]CashbackRecord, error) {
    // Get list of settled periods (1 API call)
    settledPeriods, err := s.GetSettledPeriods(ctx)
    if err != nil {
        return nil, err
    }
    
    // Only check periods that have been settled (5-10 periods typically)
    var records []CashbackRecord
    for _, period := range settledPeriods {
        // 2 API calls per settled period only
        if record, err := s.fetchCashbackForPeriod(ctx, period, userAddress); err == nil {
            records = append(records, record)
        }
    }
    
    return records, nil
}

// Update during settlement
func (s *KVStore) SettleCashbackPeriod(ctx context.Context, period uint64, ...) error {
    // ... settlement logic ...
    
    // Add to settled periods list
    s.AddSettledPeriod(ctx, period)
}
```

**Pros:**
- ✅ Minimal API calls (1 + 2N where N = settled periods)
- ✅ Scales well (only checks relevant periods)
- ✅ Fast even with 100+ total periods
- ✅ Easy to maintain

**Cons:**
- ⚠️ Requires schema change
- ⚠️ Need migration for existing data
- ⚠️ Settlement code needs update

---

### **Solution C: In-Memory Cache** ⭐ NICE TO HAVE

**Complexity:** Low  
**Time to Implement:** ~1 hour  
**Performance Gain:** Subsequent loads instant

**Implementation:**
```go
type CashbackCache struct {
    records   []CashbackRecord
    expiresAt time.Time
    mu        sync.RWMutex
}

var cache = make(map[string]*CashbackCache)

func (s *KVStore) GetClaimableCashbackRecords(ctx context.Context, userAddress string) ([]CashbackRecord, error) {
    // Check cache
    if cached := getCached(userAddress); cached != nil && !cached.IsExpired() {
        return cached.records, nil
    }
    
    // Fetch from KV
    records, err := s.fetchCashbackRecordsFromKV(ctx, userAddress)
    if err != nil {
        return nil, err
    }
    
    // Cache for 5 minutes
    setCache(userAddress, records, 5*time.Minute)
    
    return records, nil
}
```

**Pros:**
- ✅ Instant subsequent loads
- ✅ Reduces KV API calls
- ✅ Easy to implement
- ✅ No schema changes

**Cons:**
- ⚠️ First load still slow (unless combined with A or B)
- ⚠️ Memory usage (minimal for typical use)
- ⚠️ Cache invalidation complexity

---

### **Solution D: Lazy Loading** ⚠️ NOT RECOMMENDED

**Implementation:**
```typescript
// Only load when tab is active
useEffect(() => {
  if (activeTab === 'cashback') {
    loadCashbackStats();
  }
}, [activeTab]);
```

**Pros:**
- ✅ Initial page load faster

**Cons:**
- ❌ First cashback tab click still slow
- ❌ Poor UX (user waits when switching tabs)
- ❌ Doesn't solve the root problem

---

## 🎯 Recommended Implementation Plan

### Phase 1: Quick Win (Week 1)
1. **Implement Solution A (Parallel Loading)**
   - File: `pkg/store/cashback.go`
   - Function: `GetClaimableCashbackRecords()`
   - Time: 2 hours
   - Gain: 10x faster (20s → 2-3s)

### Phase 2: Long-Term Fix (Week 2)
2. **Implement Solution B (Settled Periods Index)**
   - Files: `pkg/store/cashback.go`, `pkg/blockchain/cashback_settlement.go`
   - Time: 4 hours
   - Gain: 20x+ faster (20s → <1s)

### Phase 3: Optimization (Week 3)
3. **Add Solution C (Cache Layer)**
   - File: `pkg/store/cashback.go`
   - Time: 1 hour
   - Gain: Instant subsequent loads

---

## 🐛 UX Issues Identified

### Issue 1: Referral Error for New Users

**Location:** `frontend/src/app/wallet/components/rewards/ReferralRewardsSection.tsx:137-150`

**Error:**
```
RuntimeError: no referral code for this address: get: 'key not found' (10009)
```

**Root Cause:**
- `GetReferralStats()` calls `GetReferralCodeByAddress()`
- Returns error if user hasn't created a referral code yet
- This is **expected behavior** but poor UX

**Current UX:**
```typescript
if (error) {
  return (
    <div className={styles.placeholderCard}>
      <h3 style={{ color: theme.colorError }}>Error Loading Referral Data</h3>
      <p>{error}</p> // Shows technical error message
      <Button onClick={retry}>Retry</Button>
    </div>
  );
}
```

**Improved UX:**
```typescript
// Detect "no referral code" error specifically
const isNewUser = error?.includes('no referral code');

if (isNewUser) {
  return (
    <div className={styles.placeholderCard}>
      <Users size={48} />
      <h3>Create Your Referral Code</h3>
      <p>Start earning rewards by referring friends!</p>
      <Button onClick={createReferralCode} type="primary">
        Generate Referral Code
      </Button>
    </div>
  );
}

if (error) {
  return (
    <div className={styles.placeholderCard}>
      <h3 style={{ color: theme.colorError }}>Error Loading Referral Data</h3>
      <p>{error}</p>
      <Button onClick={retry}>Retry</Button>
    </div>
  );
}
```

**Benefits:**
- ✅ Clear call-to-action for new users
- ✅ No confusing error messages
- ✅ Better onboarding experience

---

## 📈 Expected Performance Improvements

| Solution | API Calls | Latency | Improvement |
|----------|-----------|---------|-------------|
| **Current** | 104 (sequential) | ~20s | Baseline |
| **A: Parallel** | 104 (parallel) | ~2-3s | **10x faster** |
| **B: Index** | 10-20 | <1s | **20x+ faster** |
| **C: Cache** | 0 (cached) | <50ms | **400x faster** |
| **A+B+C** | 10-20 (first) / 0 (cached) | <1s / <50ms | **Best** |

---

## 🚀 Next Steps

1. **Create feature branch:** `feature/performance-cashback-loading`
2. **Implement Solution A** (Parallel Loading)
3. **Test with real data** (multiple periods)
4. **Measure performance improvement**
5. **Deploy to testnet**
6. **Plan Solution B** (Index) for next iteration

---

## 📝 Related Files

### Backend
- `pkg/store/cashback.go:315-380` - Cashback loading logic
- `internal/services/cashbackservice.go:62-99` - Service layer
- `internal/services/deai_service.go:539-604` - Mining rewards
- `internal/services/referralservice.go:51-70` - Referral stats

### Frontend
- `frontend/src/app/wallet/RewardsContent.tsx:58-120` - Tab management
- `frontend/src/app/wallet/components/rewards/CashbackRewardsSection.tsx:95-123` - Cashback loading
- `frontend/src/app/wallet/components/rewards/MiningRewardsSection.tsx:64-91` - Mining loading
- `frontend/src/app/wallet/components/rewards/ReferralRewardsSection.tsx:30-63` - Referral loading

---

**Status:** 📋 Documented, ready for implementation  
**Priority:** 🔴 High (affects user experience)  
**Assigned:** TBD  
**Estimated Time:** 7 hours total (2h + 4h + 1h)

