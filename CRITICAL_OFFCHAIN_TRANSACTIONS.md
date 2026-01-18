# 🚨 Critical Offchain Transactions Analysis

## Overview

Dokumen ini mengidentifikasi semua operasi offchain yang **CRITICAL** karena jika data tidak konsisten dengan onchain, user tidak bisa claim rewards.

---

## 🔴 TIER 1: CRITICAL - Settlement Operations

### 1. Generate Mining Settlement
**File**: `pkg/store/mining_settlement.go`
**Function**: `GenerateMiningSettlement(ctx, rewardType)`
**Line**: 18

**What it does:**
- Mengambil semua unsettled job rewards
- Generate 9-field Merkle tree
- Save Merkle root ke KV store
- Save proofs untuk setiap contributor

**Critical Data:**
```go
// Merkle Root
key: "settlement:{periodID}:merkle_root"
value: "0xABCD..."

// Proofs
key: "proof:{periodID}:{address}"
value: {
    period, contributor, contributorAmount,
    developerAmount, userAmount, affiliatorAmount,
    developer, user, affiliator, proof
}
```

**Risk if Inconsistent:**
- ❌ User tidak bisa claim (proof invalid)
- ❌ Reward terkunci permanent
- 🔴 **SEVERITY: CRITICAL**

**Location**: Lines 18-200+

---

### 2. Save Settlement Period
**File**: `pkg/store/merkle.go`
**Function**: `SaveSettlementPeriod(ctx, period)`
**Line**: 267

**What it does:**
- Save metadata settlement period
- Track status (pending, completed, failed)
- Store total amounts dan contributor count

**Critical Data:**
```go
key: "settlement:{periodID}"
value: {
    PeriodID, MerkleRoot, RewardType,
    TotalAmount, ContributorCount,
    ProofsSaved, Status, CreatedAt
}
```

**Risk if Inconsistent:**
- ❌ Frontend tidak tahu period sudah settled
- ❌ User tidak bisa query claimable rewards
- 🟡 **SEVERITY: HIGH**

**Location**: Lines 267-287

---

### 3. Deduct Settled Rewards
**File**: `pkg/store/contributor.go`
**Function**: `DeductSettledRewards(ctx, address, rewardType, amount)`
**Line**: 170

**What it does:**
- Deduct pending rewards dari contributor balance
- Prevent double-counting rewards
- Atomic operation dengan lock

**Critical Data:**
```go
key: "contributor:{address}"
value: {
    PendingKawai: "1000" → "0"  // After settlement
}
```

**Risk if Inconsistent:**
- ❌ Contributor bisa claim rewards 2x (offchain + onchain)
- ❌ Double-spending rewards
- 🔴 **SEVERITY: CRITICAL**

**Location**: Lines 170-242

---

### 4. Mark Job Rewards as Settled
**File**: `pkg/store/job_rewards.go`
**Function**: `MarkJobRewardsAsSettled(ctx, address, periodID)`
**Line**: 102

**What it does:**
- Mark job rewards sebagai settled
- Prevent re-inclusion dalam settlement berikutnya

**Critical Data:**
```go
key: "job_rewards:{address}:{jobID}"
value: {
    Settled: false → true,
    SettledPeriod: 0 → 5
}
```

**Risk if Inconsistent:**
- ❌ Job rewards di-settle 2x
- ❌ Contributor dapat reward double
- 🔴 **SEVERITY: CRITICAL**

**Location**: Lines 102-151

---

## 🟡 TIER 2: HIGH - Balance Operations

### 5. Record Job Reward
**File**: `pkg/store/contributor.go`
**Function**: `RecordJobReward(ctx, contributor, user, tokenUsage, referrer)`
**Line**: 357

**What it does:**
- Calculate reward splits (85/5/5/5 atau 90/5/5)
- Update contributor pending balance
- Record job reward untuk settlement

**Critical Data:**
```go
// Contributor balance
key: "contributor:{address}"
value: { PendingKawai: "1000" → "1100" }

// Job reward record
key: "job_rewards:{address}:{timestamp}"
value: {
    ContributorAmount, DeveloperAmount,
    UserAmount, AffiliatorAmount,
    Developer, User, Affiliator
}
```

**Risk if Inconsistent:**
- ❌ Reward calculation salah
- ❌ Settlement amount tidak match dengan job records
- 🟡 **SEVERITY: HIGH**

**Location**: Lines 357-500+

---

### 6. Deduct Balance Atomic
**File**: `pkg/store/balance.go`
**Function**: `DeductBalanceAtomic(ctx, address, amount)`
**Line**: 77

**What it does:**
- Deduct USDT dari user balance
- Atomic operation dengan retry
- Check sufficient balance

**Critical Data:**
```go
key: "balance:{address}"
value: {
    USDTBalance: "10000000" → "9000000",  // -1 USDT
    TrialClaimed: true
}
```

**Risk if Inconsistent:**
- ❌ User bisa spend lebih dari balance
- ❌ Negative balance
- 🟡 **SEVERITY: HIGH**

**Location**: Lines 77-133

---

### 7. Add Balance Atomic
**File**: `pkg/store/balance.go`
**Function**: `AddBalanceAtomic(ctx, address, amount)`
**Line**: 136

**What it does:**
- Add USDT ke user balance (deposit, cashback)
- Atomic operation dengan retry

**Critical Data:**
```go
key: "balance:{address}"
value: {
    USDTBalance: "5000000" → "10000000",  // +5 USDT
    KawaiBalance: "0" → "20000000000000000000"  // +20 KAWAI
}
```

**Risk if Inconsistent:**
- ❌ User kehilangan deposit
- ❌ Cashback tidak tercatat
- 🟡 **SEVERITY: HIGH**

**Location**: Lines 136-187

---

### 8. Claim Free Trial with Referral
**File**: `pkg/store/referral.go`
**Function**: `ClaimFreeTrialWithReferral(ctx, address, machineID, referralCode)`
**Line**: 177

**What it does:**
- Claim trial bonus (5 atau 10 USDT)
- Validate referral code
- Reward referrer
- Mark trial as claimed

**Critical Data:**
```go
// User balance
key: "balance:{address}"
value: {
    USDTBalance: "0" → "10000000",  // 10 USDT
    KawaiBalance: "0" → "200000000000000000000",  // 200 KAWAI
    TrialClaimed: false → true,
    ReferrerAddress: "0x..."
}

// Machine ID
key: "trial_machine:{machineID}"
value: "true"

// Referrer stats
key: "referral:code:{CODE}"
value: {
    TotalReferrals: 5 → 6,
    TotalEarnings: "25000000" → "30000000"
}
```

**Risk if Inconsistent:**
- ❌ User bisa claim trial multiple times
- ❌ Referrer tidak dapat reward
- 🟡 **SEVERITY: HIGH**

**Location**: Lines 177-220

---

## 🟢 TIER 3: MEDIUM - Cashback Operations

### 9. Store Cashback Data
**File**: `pkg/store/cashback_kv.go`
**Function**: `StoreCashbackData(ctx, key, data)`
**Line**: 15

**What it does:**
- Store cashback records untuk deposits
- Track cashback amounts per user

**Critical Data:**
```go
key: "cashback:{address}:{depositID}"
value: {
    DepositAmount: "1000000000",  // 1000 USDT
    CashbackAmount: "20000000000000000000000",  // 20K KAWAI
    Tier: "gold",
    Timestamp: 1234567890
}
```

**Risk if Inconsistent:**
- ❌ User tidak bisa claim cashback
- ❌ Cashback amount salah
- 🟢 **SEVERITY: MEDIUM**

**Location**: Lines 15-25

---

### 10. Add Settled Cashback Period
**File**: `pkg/store/cashback_kv.go`
**Function**: `AddSettledCashbackPeriod(ctx, period)`
**Line**: 92

**What it does:**
- Track which periods sudah settled
- Optimize claim queries

**Critical Data:**
```go
key: "cashback_settled_periods"
value: [1, 2, 3, 4, 5]  // Array of settled periods
```

**Risk if Inconsistent:**
- ❌ Frontend tidak tahu period mana yang claimable
- ❌ User tidak bisa query cashback
- 🟢 **SEVERITY: MEDIUM**

**Location**: Lines 92-115

---

## 🔵 TIER 4: LOW - Marketplace & Metadata

### 11. Store Marketplace Data
**File**: `pkg/store/marketplace.go`
**Function**: `StoreMarketplaceData(ctx, key, data)`
**Line**: 13

**What it does:**
- Store P2P trading orders
- Track order status

**Critical Data:**
```go
key: "order:{orderID}"
value: {
    Seller, Buyer, Amount, Price, Status
}
```

**Risk if Inconsistent:**
- ❌ Order status tidak sync
- ❌ Trading disputes
- 🔵 **SEVERITY: LOW** (tidak affect rewards)

**Location**: Lines 13-24

---

## 📊 Summary Table

| # | Operation | File | Line | Severity | Impact if Inconsistent |
|---|-----------|------|------|----------|------------------------|
| 1 | Generate Mining Settlement | `mining_settlement.go` | 18 | 🔴 CRITICAL | User tidak bisa claim, reward terkunci |
| 2 | Save Settlement Period | `merkle.go` | 267 | 🟡 HIGH | Frontend tidak tahu status settlement |
| 3 | Deduct Settled Rewards | `contributor.go` | 170 | 🔴 CRITICAL | Double-spending rewards |
| 4 | Mark Job Rewards Settled | `job_rewards.go` | 102 | 🔴 CRITICAL | Job di-settle 2x |
| 5 | Record Job Reward | `contributor.go` | 357 | 🟡 HIGH | Reward calculation salah |
| 6 | Deduct Balance Atomic | `balance.go` | 77 | 🟡 HIGH | Negative balance, overspending |
| 7 | Add Balance Atomic | `balance.go` | 136 | 🟡 HIGH | User kehilangan deposit |
| 8 | Claim Free Trial | `referral.go` | 177 | 🟡 HIGH | Multiple trial claims |
| 9 | Store Cashback Data | `cashback_kv.go` | 15 | 🟢 MEDIUM | Cashback tidak bisa claim |
| 10 | Add Settled Period | `cashback_kv.go` | 92 | 🟢 MEDIUM | Query optimization issue |
| 11 | Store Marketplace Data | `marketplace.go` | 13 | 🔵 LOW | Trading disputes |

---

## 🛡️ Protection Mechanisms

### 1. Atomic Operations with Retry
```go
// pkg/store/balance.go:77
maxRetries := 5
backoff := 50 * time.Millisecond

for attempt := 0; attempt < maxRetries; attempt++ {
    // Read-Modify-Write with version check
}
```

### 2. Mutex Locks for Critical Sections
```go
// pkg/store/contributor.go:173
lockInterface, _ := contributorLocks.LoadOrStore(address, &sync.Mutex{})
lock := lockInterface.(*sync.Mutex{})
lock.Lock()
defer lock.Unlock()
```

### 3. Settlement Rollback Support
```go
// pkg/store/settlement.go:100
func PerformSettlementWithConfig(...) {
    // Rollback on failure
    if err != nil {
        s.rollbackSettlement(ctx, periodID)
    }
}
```

### 4. Idempotency Checks
```go
// pkg/store/referral.go:243
if currentData.TrialClaimed {
    return fmt.Errorf("free trial already claimed")
}
```

### 5. 🆕 Double-Verification: Telegram + KV Store
**File**: `pkg/store/contributor.go:357` + `pkg/alert/telegram.go:102`

**What it does:**
- Every job reward is recorded in **TWO independent systems**:
  1. **Cloudflare KV** (primary storage)
  2. **Telegram Channel** (immutable audit trail)
- During settlement, both sources can be cross-checked for consistency
- Telegram messages cannot be edited (if pinned), providing tamper-proof backup

**Implementation:**
```go
// pkg/store/contributor.go:540-565
// 1. Save to KV Store (primary)
if err := s.SaveJobReward(ctx, jobRecord); err != nil {
    slog.Warn("Failed to save job reward record", "error", err)
}

// 2. Send to Telegram (backup + audit)
if s.telegramAlerter != nil {
    s.telegramAlerter.SendJobRewardLog(jobRecord)
}
```

**Benefits:**
- ✅ **Immutable Audit Trail**: Telegram messages provide tamper-proof record
- ✅ **Independent Verification**: Two separate data sources reduce single-point-of-failure
- ✅ **Disaster Recovery**: Can restore from Telegram if KV corrupts
- ✅ **Real-time Monitoring**: Admins can see rewards in real-time
- ✅ **Fraud Detection**: Discrepancies between KV and Telegram indicate issues
- ✅ **Async Operation**: Telegram send is non-blocking (goroutine)

**Rate Limiting:**
- Telegram API: 30 messages/second per bot
- Implementation: Async send (goroutine) prevents blocking
- If rate limit hit: Message queued, doesn't fail job reward recording

**Privacy:**
- Use private Telegram channel for job reward logs
- Only admins have access to channel
- Addresses are shortened in display (0x1234...5678)

### 6. 🆕 Double-Verification: Cashback Tracking
**File**: `pkg/store/cashback.go:143` + `pkg/alert/telegram.go`

**What it does:**
- Every cashback record is logged in **TWO independent systems**:
  1. **Cloudflare KV** (primary storage)
  2. **Telegram Channel** (immutable audit trail)
- During cashback settlement, both sources can be cross-checked for consistency
- Telegram messages provide tamper-proof backup for verification

**Implementation:**
```go
// pkg/store/cashback.go:189-203
// 1. Save to KV Store (primary)
if err := s.StoreCashbackData(ctx, key, data); err != nil {
    return fmt.Errorf("failed to store cashback record: %w", err)
}

// 2. Send to Telegram (backup + audit)
if s.telegramAlerter != nil {
    s.telegramAlerter.SendCashbackLog(&types.CashbackRecord{
        UserAddress:    userAddress,
        TxHash:         txHash,
        DepositAmount:  depositAmount.String(),
        CashbackAmount: cashbackAmount,
        RateBPS:        rate,
        Tier:           tier,
        IsFirstTime:    isFirstTime,
        Period:         period,
        Timestamp:      time.Now(),
    })
}
```

**Message Format:**
```text
💎 Cashback Tracked
2026-01-18 15:30:45 | Period 5 | Tier 2

```json
{"user_address":"0x1234...","tx_hash":"0xabcd...","deposit_amount":"10000000","cashback_amount":"200000000000000000000","rate_bps":1000,"tier":2,"is_first_time":true,"period":5,"timestamp":"2026-01-18T15:30:45Z"}
```

📊 Rate: 10.00% | User: 0x1234...5678 🎁 First-time
```

**Benefits:**
- ✅ **Verify Cashback Calculation**: Cross-check deposit amount → cashback amount
- ✅ **Detect Tier Manipulation**: Ensure tier assignment is correct
- ✅ **Prevent Double Claims**: Verify each deposit only tracked once
- ✅ **Settlement Reconciliation**: Sum Telegram records == KV settlement total
- ✅ **Fraud Detection**: Detect suspicious patterns (multiple first-time bonuses)

### 7. 🆕 Double-Verification: Referral Trial Claims
**File**: `pkg/store/referral.go:177` + `pkg/alert/telegram.go`

**What it does:**
- Every trial claim (with or without referral) is logged in **TWO independent systems**:
  1. **Cloudflare KV** (primary storage)
  2. **Telegram Channel** (immutable audit trail)
- Enables fraud detection (multiple claims, machine ID spoofing)
- Verifies referrer bonus calculations
- Tracks machine IDs for anti-abuse

**Implementation:**
```go
// pkg/store/referral.go:217-233
// After successful trial claim
if s.telegramAlerter != nil {
    s.telegramAlerter.SendReferralTrialLog(&types.ReferralTrialRecord{
        UserAddress:     address,
        ReferrerAddress: referrerAddress,
        ReferralCode:    referralCode,
        TrialUSDT:       fmt.Sprintf("%d", usdtBonus),
        TrialKAWAI:      kawaiBonus,
        ReferrerBonus:   referrerBonus,
        MachineID:       machineID,
        IsReferral:      hasReferral,
        Timestamp:       time.Now(),
    })
}
```

**Message Format:**
```text
🎁 Referral Trial Claimed
2026-01-18 15:30:45

```json
{"user_address":"0x1234...","referrer_address":"0x5678...","referral_code":"ABC123","trial_usdt":"10000000","trial_kawai":"200000000000000000000","referrer_bonus":"5000000","machine_id":"hash...","is_referral":true,"timestamp":"2026-01-18T15:30:45Z"}
```

👤 User: 0x1234...5678
🤝 Referrer: 0x5678...abcd (Code: ABC123)
🔒 Machine: hash1234567890ab...
```

**Benefits:**
- ✅ **Fraud Detection**: Detect multiple claims from same machine ID
- ✅ **Referrer Verification**: Verify referrer bonus calculations
- ✅ **Anti-Abuse**: Track patterns (same user, different machine IDs)
- ✅ **Settlement Reconciliation**: Sum referrer bonuses == referral settlement total
- ✅ **Audit Trail**: Immutable record of all trial claims

### 8. 🆕 Discord Fallback for Telegram Failures
**File**: `pkg/alert/telegram.go` + `pkg/alert/discord.go`

**What it does:**
- If Telegram fails to send audit logs, automatically fallback to Discord webhook
- Ensures audit trail is NEVER lost, even if Telegram is down
- Discord webhook is more reliable (no rate limits, no bot token issues)
- All critical logs (job rewards, cashback, referral trials) have Discord fallback

**Implementation:**
```go
// pkg/alert/telegram.go:27-32
func NewTelegramAlert() *TelegramAlert {
    return &TelegramAlert{
        BotToken:        constant.GetTelegramBotToken(),
        ChatID:          constant.GetTelegramChatId(),
        Client:          &http.Client{Timeout: 10 * time.Second},
        DiscordFallback: NewDiscordAlert(), // Initialize Discord fallback
    }
}

// pkg/alert/telegram.go:145-158
go func() {
    if err := t.SendMessage(text); err != nil {
        slog.Error("Failed to send job reward log to Telegram", "error", err)
        
        // Fallback to Discord if Telegram fails
        if t.DiscordFallback != nil {
            discordMsg := fmt.Sprintf("💰 **Job Reward** (Telegram Fallback)\n```json\n%s\n```", string(jsonData))
            if discordErr := t.DiscordFallback.SendMessage(discordMsg); discordErr != nil {
                slog.Error("Failed to send job reward log to Discord fallback", "error", discordErr)
            } else {
                slog.Info("Job reward log sent to Discord fallback successfully")
            }
        }
    }
}()
```

**Configuration:**
```bash
# .env
DISCORD_WEBHOOK=https://discord.com/api/webhooks/1462463310224953526/...
```

**Benefits:**
- ✅ **High Availability**: Two independent messaging systems (Telegram + Discord)
- ✅ **No Single Point of Failure**: If Telegram down, Discord still works
- ✅ **Automatic Failover**: No manual intervention needed
- ✅ **Audit Trail Guaranteed**: Critical logs NEVER lost
- ✅ **Rate Limit Bypass**: Discord webhook has no rate limits
- ✅ **Async Operation**: Fallback is non-blocking (goroutine)

**Failure Scenarios:**
1. **Telegram Success**: Log sent to Telegram only (normal operation)
2. **Telegram Fails**: Log automatically sent to Discord (fallback)
3. **Both Fail**: Error logged to slog, but KV store still has data (primary storage unaffected)

**Message Format (Discord):**
```text
💰 **Job Reward** (Telegram Fallback)
```json
{"timestamp":"2026-01-18T15:30:45Z","contributor_address":"0x1234...","user_address":"0x5678...","token_usage":1000,"reward_type":"kawai"}
```
```

---

## 🚨 Critical Failure Scenarios

### Scenario 1: Settlement Generated but Upload Failed
```
1. GenerateMiningSettlement() ✅ Success
   - Merkle root saved offchain
   - Proofs saved offchain
   - Contributor balances deducted

2. uploadMiningRoot() ❌ Failed
   - Merkle root NOT uploaded onchain
   - Contract still has old/empty root

Result:
- User has proofs offchain
- Contract has no root to verify against
- User CANNOT claim
- Rewards LOCKED until re-upload
```

**Recovery**: Re-run `reward-settlement upload --type mining`

---

### Scenario 2: Balance Deducted but Service Failed
```
1. DeductBalanceAtomic() ✅ Success
   - User balance: 10 USDT → 9 USDT

2. AI Service Call ❌ Failed
   - Service timeout/error
   - User tidak dapat response

Result:
- User charged 1 USDT
- User tidak dapat service
- Balance inconsistent
```

**Recovery**: Manual refund atau retry mechanism

---

### Scenario 3: Job Reward Recorded but Settlement Failed
```
1. RecordJobReward() ✅ Success
   - Job reward saved
   - Contributor balance updated

2. GenerateMiningSettlement() ❌ Failed
   - Settlement generation error
   - Merkle tree not created

Result:
- Job rewards accumulate
- No settlement for this period
- User must wait for next period
```

**Recovery**: Fix error and re-run settlement

---

### Scenario 4: Trial Claimed but Referrer Not Rewarded
```
1. ClaimFreeTrialWithReferral() ✅ Partial Success
   - User balance updated (10 USDT)
   - Trial marked as claimed
   
2. rewardReferrer() ❌ Failed
   - Referrer reward not saved
   - Referrer stats not updated

Result:
- User got bonus
- Referrer did NOT get reward
- Referrer stats incorrect
```

**Recovery**: Manual compensation for referrer

---

## 💡 Recommendations

### 1. Implement Transaction Log
```go
type TransactionLog struct {
    ID          string
    Type        string  // "settlement", "balance", "reward"
    Status      string  // "pending", "completed", "failed"
    Data        json.RawMessage
    CreatedAt   time.Time
    CompletedAt *time.Time
}

// Log every critical operation
func (s *KVStore) LogTransaction(ctx, txType, data) error {
    // Save to audit log
}
```

### 2. Add Reconciliation Job
```go
// Daily reconciliation
func ReconcileData(ctx) error {
    // 1. Check onchain vs offchain merkle roots
    // 2. Verify claim status sync
    // 3. Validate balance totals
    // 4. Alert on discrepancies
}
```

### 3. Implement Idempotency Keys
```go
func (s *KVStore) RecordJobRewardIdempotent(ctx, idempotencyKey, ...) error {
    // Check if already processed
    if s.IsProcessed(ctx, idempotencyKey) {
        return nil  // Already done
    }
    
    // Process
    err := s.RecordJobReward(ctx, ...)
    
    // Mark as processed
    s.MarkProcessed(ctx, idempotencyKey)
    return err
}
```

### 4. Add Health Checks
```go
func (s *KVStore) HealthCheck(ctx) error {
    // 1. Check KV store connectivity
    // 2. Verify data integrity
    // 3. Check for stuck settlements
    // 4. Validate balance consistency
}
```

---

## 🔍 Monitoring Checklist

- [ ] Monitor settlement success rate
- [ ] Alert on failed settlements
- [ ] Track balance operation failures
- [ ] Monitor claim success rate
- [ ] Check for stuck transactions
- [ ] Verify merkle root uploads
- [ ] Reconcile onchain vs offchain daily
- [ ] Audit trail for all critical ops
- [ ] Backup KV store regularly
- [ ] Test recovery procedures

---

**Last Updated**: 2026-01-18
**Severity Levels**: 🔴 Critical | 🟡 High | 🟢 Medium | 🔵 Low
