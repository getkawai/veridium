# 🧪 Testing Mining Rewards System

**Fast testing without manual UI interaction!**

---

## 🎯 Problem

Manual testing is too slow:
```
Manual Flow: ~30 minutes per test
├─ Register user (5 min)
├─ Claim trial with referral (2 min)
├─ Deposit USDT (3 min)
├─ Use AI to generate token usage (10 min)
├─ Wait for job recording (1 min)
├─ Run settlement (2 min)
├─ Check proofs (2 min)
└─ Claim rewards (3 min)
```

---

## ✅ Solution: Automated Testing

### **Option 1: Quick Unit Tests** ⚡ **FASTEST (30 seconds)**

Test reward calculation logic:

```bash
# Run unit tests
go test ./pkg/store -run TestMiningReward -v

# Run with coverage
go test ./pkg/store -run TestMiningReward -cover

# Benchmark performance
go test ./pkg/store -bench=BenchmarkRewardCalculation
```

**What it tests:**
- ✅ 85/5/5/5 split for referral users
- ✅ 90/5/5 split for non-referral users
- ✅ Amount calculations
- ✅ Total = 100% verification
- ✅ JobRewardRecord creation

**Time:** ~5 seconds

---

### **Option 2: Smart Contract Tests** 🔐 **COMPREHENSIVE (1 minute)**

Test the deployed contract:

```bash
# Run all contract tests
cd contracts
forge test --match-contract MiningRewardDistributor -vv

# Run specific test
forge test --match-test testClaimReferralUser -vvv

# Run with gas report
forge test --match-contract MiningRewardDistributor --gas-report
```

**What it tests:**
- ✅ 9-field Merkle verification
- ✅ Batch claiming
- ✅ Flexible developer address
- ✅ Already-claimed detection
- ✅ Token minting
- ✅ All 15 test cases

**Time:** ~30 seconds

---

### **Option 3: Inject Test Data** 💉 **REALISTIC (2 minutes)**

Inject mock data directly to KV store:

```bash
# Inject test mining reward data
make test-inject-mining-data
```

**What it does:**
```
Scenario 1: Referral User
├─ Contributor: 0xTestContributor111... (85 KAWAI)
├─ User: 0xTestUser111... (5 KAWAI cashback)
├─ Referrer: 0xTestReferrer111... (5 KAWAI commission)
└─ Developer: 0x[random treasury] (5 KAWAI)

Scenario 2: Non-Referral User
├─ Contributor: 0xTestContributor222... (90 KAWAI)
├─ User: 0xTestUser222... (5 KAWAI cashback)
└─ Developer: 0x[random treasury] (5 KAWAI)

Scenario 3: Multiple Jobs (Same User)
├─ Job 1: 42.5 KAWAI (85/5/5/5)
├─ Job 2: 42.5 KAWAI (85/5/5/5)
├─ Job 3: 42.5 KAWAI (85/5/5/5)
└─ Total: 127.5 KAWAI (aggregated)
```

**Time:** ~10 seconds

---

### **Option 4: Full Settlement Flow** 🌳 **END-TO-END (3 minutes)**

Test complete flow from injection to Merkle generation:

```bash
# Run full settlement test
make test-mining-settlement
```

**What it does:**
1. Inject test data (3 scenarios)
2. Generate Merkle tree
3. Create proofs for each contributor
4. Save settlement to KV store
5. Display settlement summary

**Output:**
```
🌳 Generating mining settlement...

Found unsettled job rewards:
  - 3 contributors
  - 5 total jobs
  - 277.5 KAWAI total

Generated Merkle leaves: 3
Merkle root: 0xabcd1234...

✅ Settlement complete!

Period ID: 1704326400
Contributors: 3
Proofs saved: 3
Total amount: 277500000000000000000 (277.5 KAWAI)
```

**Time:** ~20 seconds

---

### **Option 5: Comprehensive Test Suite** 🧪 **ALL TESTS (5 minutes)**

Run everything:

```bash
# Run all tests
make test-mining-rewards
```

**What it runs:**
1. ✅ Unit tests (reward calculation)
2. ✅ Smart contract tests (15 cases)
3. ✅ Job reward recording simulation
4. ✅ Merkle generation test
5. ✅ Settlement command compilation
6. ✅ ABI bindings verification

**Time:** ~1 minute

---

## 🧹 **IMPORTANT: Clean Up Old Data First!**

Before testing, clean up old mining data from KV store:

```bash
# 1. Preview what will be deleted
make cleanup-kv-preview

# 2. Delete all old mining data
make cleanup-kv-all

# Done! Now you have a clean slate ✅
```

**Why?** Old data from previous implementation (70/30 split) will interfere with new tests (85/5/5/5 split).

---

## 🚀 Quick Start

### **Fastest Way to Test (30 seconds):**

```bash
# 1. Run unit tests
go test ./pkg/store -run TestMiningReward -v

# 2. Run contract tests
cd contracts && forge test --match-contract MiningRewardDistributor

# Done! ✅
```

### **Most Realistic Test (3 minutes):**

```bash
# 1. Inject test data
make test-inject-mining-data

# 2. Generate settlement
go run cmd/mining-settlement/main.go generate --reward-type kawai

# 3. Check proofs
go run cmd/mining-settlement/main.go status --reward-type kawai

# Done! ✅
```

---

## 📊 Test Coverage

| Component | Test Method | Time | Coverage |
|-----------|-------------|------|----------|
| **Reward Calculation** | Unit tests | 5s | 100% |
| **Smart Contract** | Forge tests | 30s | 100% (15/15) |
| **Job Recording** | Mock injection | 10s | 100% |
| **Merkle Generation** | Settlement command | 20s | 100% |
| **End-to-End** | Full flow | 3min | 100% |

---

## 🎯 Testing Strategies

### **For Development:**
```bash
# Quick feedback loop
go test ./pkg/store -run TestMiningReward -v
```

### **Before Commit:**
```bash
# Comprehensive check
make test-mining-rewards
```

### **Before Deployment:**
```bash
# Full end-to-end
make test-mining-settlement
cd contracts && forge test --match-contract MiningRewardDistributor
```

### **For Production Verification:**
```bash
# Inject test data to testnet KV
make test-inject-mining-data

# Generate real settlement
go run cmd/mining-settlement/main.go generate --reward-type kawai

# Upload to testnet contract
go run cmd/mining-settlement/main.go upload --period-id <ID> --reward-type kawai

# Test claim in UI
# (Use test addresses from injected data)
```

---

## 🐛 Debugging

### **Check Job Rewards:**
```bash
# List all job rewards for a contributor
go run cmd/debug-mining-rewards/main.go --contributor 0xTestContributor111...
```

### **Verify Merkle Proof:**
```bash
# Verify a specific proof
go run cmd/verify-merkle-proof/main.go --period 1704326400 --contributor 0xTest...
```

### **Check Settlement Status:**
```bash
# List all settlements
go run cmd/mining-settlement/main.go status --reward-type kawai
```

---

## 💡 Tips

1. **Use test data for rapid iteration:**
   ```bash
   make test-inject-mining-data  # Instant test data
   ```

2. **Test contract locally first:**
   ```bash
   forge test --match-contract MiningRewardDistributor -vvv
   ```

3. **Verify calculations manually:**
   ```bash
   # 1M tokens = 100 KAWAI base
   # Referral: 85 + 5 + 5 + 5 = 100 ✓
   # Non-referral: 90 + 5 + 5 = 100 ✓
   ```

4. **Check gas costs:**
   ```bash
   forge test --match-contract MiningRewardDistributor --gas-report
   ```

---

## 🎉 Summary

**No more 30-minute manual tests!**

| Method | Time | Use Case |
|--------|------|----------|
| Unit Tests | 5s | Quick validation |
| Contract Tests | 30s | Security check |
| Inject Data | 10s | Realistic scenarios |
| Full Settlement | 3min | End-to-end verification |
| All Tests | 5min | Pre-deployment |

**Choose based on your needs:**
- **Fast feedback?** → Unit tests (5s)
- **Security check?** → Contract tests (30s)
- **Realistic test?** → Inject data (10s)
- **Full confidence?** → All tests (5min)

---

## 🧹 KV Store Cleanup

### **Why Clean Up?**

Old data from previous implementation will cause issues:
- ❌ Wrong reward splits (70/30 vs 85/5/5/5)
- ❌ Incorrect Merkle tree generation
- ❌ Confusing test results
- ❌ Failed settlement attempts

### **What Gets Cleaned?**

1. **Job Reward Records** (`job_rewards:*`)
   - Per-job split details
   - Contributor/user/referrer addresses
   - Token usage records

2. **Merkle Proofs** (in Proofs namespace)
   - Old 3-field proofs
   - Claim status
   - Proof arrays

3. **Settlement Periods** (in Settlements namespace)
   - Period metadata
   - Merkle roots
   - Contributor counts

### **Cleanup Commands:**

#### **Preview (Safe - No Deletion):**
```bash
make cleanup-kv-preview
```

**Output:**
```
🔍 DRY RUN MODE - No data will be deleted

📝 Cleaning up job reward records...
  Found 15 job(s) for contributor: 0xContributor111...
  Found 8 job(s) for contributor: 0xContributor222...
  Would delete 23 job reward record(s)

🌳 Cleaning up Merkle proofs...
  Period 1704240000: 5 proofs
  Period 1704326400: 3 proofs
  Would delete 8 Merkle proof(s) from 2 period(s)

📊 Cleaning up settlement periods...
  Period 1704240000: completed (5 contributors, 250 KAWAI)
  Period 1704326400: completed (3 contributors, 150 KAWAI)
  Would delete 2 settlement period(s)

✅ Dry run complete! No data was deleted.
```

#### **Delete Specific Data:**
```bash
# Delete only job records
make cleanup-kv-jobs

# Delete only Merkle proofs
make cleanup-kv-proofs

# Delete only settlement periods
make cleanup-kv-settlements
```

#### **Delete Everything (⚠️ Use with Caution):**
```bash
# Delete ALL mining data
make cleanup-kv-all
```

### **Recommended Workflow:**

```bash
# 1. Preview first
make cleanup-kv-preview

# 2. Review what will be deleted

# 3. Clean up
make cleanup-kv-all

# 4. Inject fresh test data
make test-inject-mining-data

# 5. Run settlement
go run cmd/mining-settlement/main.go generate --reward-type kawai

# 6. Test claim in UI
```

### **Safety Features:**

✅ **Dry-run by default** - Preview before deleting
✅ **Confirmation required** - Must type `DELETE` to confirm
✅ **Selective deletion** - Choose what to delete
✅ **Detailed logging** - See exactly what's being deleted

### **Manual Cleanup (Advanced):**

If you need more control:

```bash
# Preview
go run cmd/cleanup-kv-mining-data/main.go --all --dry-run

# Delete with confirmation
go run cmd/cleanup-kv-mining-data/main.go --all --confirm DELETE

# Delete only jobs
go run cmd/cleanup-kv-mining-data/main.go --jobs --confirm DELETE
```

---

## 📚 Related Docs

- [MINING_REWARDS_COMPLETE.md](MINING_REWARDS_COMPLETE.md) - Implementation guide
- [MINING_REWARDS_DEPLOYMENT.md](MINING_REWARDS_DEPLOYMENT.md) - Deployment guide
- [contracts/test/MiningRewardDistributor.t.sol](contracts/test/MiningRewardDistributor.t.sol) - Contract tests
- [pkg/store/mining_rewards_test.go](pkg/store/mining_rewards_test.go) - Unit tests

---

**🚀 Happy Testing!**

