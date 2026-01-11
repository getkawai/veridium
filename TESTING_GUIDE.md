# 🧪 Comprehensive Testing Guide

Complete end-to-end testing guide for all 3 reward systems (Mining, Cashback, Referral).

---

## 📋 Pre-Requisites

### 1. Environment Setup

- ✅ Monad Testnet RPC configured
- ✅ Admin wallet with testnet MON for gas
- ✅ MINTER_ROLE granted to all distributors
- ✅ Cloudflare KV namespaces configured (credentials in env)
- ✅ Wails desktop app built and running (no separate server needed)

### 1.5. Optional: Cleanup Old Test Data

**⚠️ Important:** If you've run tests before, cleanup old data for fresh testing:

```bash
# This will guide you through cleaning:
# • Mining job rewards (unsettled)
# • Cashback records (test data)
# • Merkle proofs (all periods)
# • Settlement periods (all)
#
# Preserves:
# ✅ User profiles
# ✅ API keys
# ✅ Authentication data

make cleanup-test-data
```

**When to cleanup:**
- ✅ Before running full E2E tests
- ✅ After failed test runs
- ✅ When switching between test scenarios
- ❌ Not needed for first-time testing

### 2. Contract Addresses (Monad Testnet)

- **MiningRewardDistributor:** `0x8117D77A219EeF5F7869897C3F0973Afb87d8427`
- **DepositCashbackDistributor:** `0xdE64f6F5bEe28762c91C76ff762365D553204e35`
- **ReferralRewardDistributor:** `0xaB0DdFbb4bD94d23a32d0C40f9F96d9A61b45463`
- **KAWAI Token:** `0xE32660b39D99988Df4bFdc7e4b68A4DC9D654722`

### 3. Check MINTER_ROLE Status

```bash
# Check if all distributors have MINTER_ROLE
make check-minter-role

# Expected output:
# ═══════════════════════════════════════════════════════════
# 🔐 MINTER_ROLE Status Check
# ═══════════════════════════════════════════════════════════
# 
# MiningRewardDistributor:       ✅ GRANTED
#    Address: 0x8117D77A219EeF5F7869897C3F0973Afb87d8427
# 
# DepositCashbackDistributor:    ✅ GRANTED
#    Address: 0xdE64f6F5bEe28762c91C76ff762365D553204e35
# 
# ReferralRewardDistributor:     ✅ GRANTED
#    Address: 0xaB0DdFbb4bD94d23a32d0C40f9F96d9A61b45463
# 
# ═══════════════════════════════════════════════════════════
# ✅ All distributors have MINTER_ROLE!
#    Ready for reward claims.
# ═══════════════════════════════════════════════════════════
```

---

## 🧪 Test 1: Mining Rewards (End-to-End)

**Duration:** ~15 minutes  
**Status:** ✅ 100% Ready

### Step 1: Inject Test Mining Data

```bash
# Generate fake mining rewards for testing
make test-inject-mining-data

# Expected output:
# ✅ Injected 10 jobs for 3 contributors
# ✅ Total rewards: ~1000 KAWAI
```

### Step 2: Generate Mining Settlement

```bash
# Generate Merkle tree and proofs
make settle-mining

# Expected output:
# 🌳 Generating Merkle tree...
# ✅ Merkle Root: 0x...
# 📊 Contributors: 3
# 📊 Total Amount: 1000000000000000000000 KAWAI
# 📊 Proofs Saved: 3
```

### Step 3: Upload Merkle Root

```bash
# Get the Merkle root from previous step
MERKLE_ROOT="0x..."  # From settle-mining output

# Upload to contract using Go tool
make upload-merkle-root TYPE=mining ROOT=$MERKLE_ROOT

# Expected output:
# ═══════════════════════════════════════
# Contract Type: mining
# Merkle Root: 0x...
# ═══════════════════════════════════════
# 
# 🚀 Uploading Merkle root to contract...
# 
# ✅ Merkle root uploaded successfully!
# ═══════════════════════════════════════
# Transaction Hash: 0x...
# Explorer: https://explorer.monad.xyz/tx/0x...
# ═══════════════════════════════════════
```

### Step 4: Start Desktop App

```bash
# Start the Wails desktop app in dev mode
make dev-hot

# This will:
# - Build the Go backend (embedded in app)
# - Start the React frontend
# - Launch the desktop application
# - Data stored in Cloudflare KV (no separate server needed)
```

### Step 5: Test Claim in UI

1. Open app: http://localhost:5173
2. Connect wallet (use test contributor address)
3. Go to: **Wallet → Rewards → Mining Rewards**
4. You should see:
   - ✅ Claimable amount: ~333 KAWAI (1000/3)
   - ✅ "Claim" button enabled
5. Click "Claim"
6. Sign transaction
7. Wait for confirmation (now with tx wait!)
8. Expected:
   - ✅ Success message
   - ✅ Transaction hash displayed
   - ✅ Balance updated
   - ✅ "Claimed" status shown

### Step 6: Verify On-Chain

```bash
# Check KAWAI balance increased
make check-balance ADDR=<CONTRIBUTOR_ADDRESS>

# Expected output:
# ═══════════════════════════════════════
# Address: 0x...
# KAWAI Balance: 333 KAWAI
# Wei Balance: 333000000000000000000 wei
# ═══════════════════════════════════════

# Check claim status
make check-claim-status TYPE=mining PERIOD=<PERIOD_ID> ADDR=<CONTRIBUTOR_ADDRESS>

# Expected output:
# ═══════════════════════════════════════
# Contract: mining
# Period: 1767549424
# Address: 0x...
# Has Claimed: true
# Status: ✅ Already Claimed
# ═══════════════════════════════════════
```

### ✅ Success Criteria

- Claim transaction succeeds
- KAWAI balance increases
- hasClaimed returns true
- UI shows "Claimed" status

---

## 🧪 Test 2: Cashback Rewards (End-to-End)

**Duration:** ~20 minutes  
**Status:** ✅ 100% Ready

### Step 1: Make USDT Deposit

1. Open app: http://localhost:5173
2. Connect wallet
3. Go to: **Wallet → Deposit**
4. Deposit amount: 100 USDT (or any amount)
5. Approve USDT
6. Confirm deposit
7. Wait for confirmation

### Step 2: Verify Cashback Tracked

**Note:** Cashback tracking happens automatically in the desktop app (Wails). Data is stored in Cloudflare KV.

1. Go to: **Wallet → Rewards → Cashback Rewards** in the app
2. You should see:
   - ✅ Total Earned: ~2000 KAWAI (2% of 100 USDT for Tier 2)
   - ✅ Pending: 2000 KAWAI
   - ✅ Claimed: 0 KAWAI
   - ✅ Total Deposits: 1

**Alternative:** Check KV directly using Cloudflare dashboard or settlement tool output.

### Step 3: Wait for Settlement Period

```bash
# For testing, manually trigger settlement:
make settle-cashback

# Expected output:
# 📊 Cashback Rewards Settlement
# Current Period: 53
# Settling Period: 52
# 🌳 Generating Merkle tree...
# ✅ Merkle Root: 0x...
# 📊 Users: 1
# 📊 Total Cashback: 2000 KAWAI
```

### Step 4: Upload Merkle Root

```bash
MERKLE_ROOT="0x..."  # From settle-cashback output

# Upload to contract using Go tool
make upload-merkle-root TYPE=cashback ROOT=$MERKLE_ROOT
```

### Step 5: Test Claim in UI

1. Go to: **Wallet → Rewards → Cashback Rewards**
2. You should see:
   - ✅ Current Tier: Tier 2 (100-500 USDT)
   - ✅ Cashback Rate: 1.25%
   - ✅ Total Earned: 2000 KAWAI
   - ✅ Claimable: 2000 KAWAI
   - ✅ Deposit History table with Period 52
3. Click "Claim" on the period
4. Sign transaction
5. Wait for confirmation
6. Expected:
   - ✅ Success message
   - ✅ Transaction hash
   - ✅ Balance updated
   - ✅ Status: "Claimed"

### Step 6: Verify On-Chain

```bash
# Check KAWAI balance increased
make check-balance ADDR=<USER_ADDRESS>

# Check claim status
make check-claim-status TYPE=cashback PERIOD=52 ADDR=<USER_ADDRESS>

# Expected: Has Claimed: true
```

### ✅ Success Criteria

- Deposit tracked correctly
- Cashback calculated correctly (tier-based)
- Settlement generates valid Merkle tree
- Claim transaction succeeds
- KAWAI balance increases
- UI updates correctly

---

## 🧪 Test 3: Referral Rewards (End-to-End)

**Duration:** ~25 minutes  
**Status:** ✅ 100% Ready

### Step 1: Create Referral Code

1. Open app as Referrer
2. Go to: **Wallet → Rewards → Referral Rewards**
3. Click "Generate Referral Code"
4. Copy code (e.g., "ABC123")

### Step 2: Refer New User

1. Open app in incognito/different browser
2. During wallet setup, enter referral code
3. Complete wallet setup
4. Expected:
   - ✅ New user gets 10 USDT + 200 KAWAI (bonus)
   - ✅ Referrer gets 5 USDT + 100 KAWAI (one-time)

### Step 3: New User Mines

1. New user runs contributor client
2. Completes some jobs
3. Mining rewards accumulate
4. 5% commission goes to referrer (AffiliatorAmount)

### Step 4: Generate Referral Settlement

```bash
# After mining period ends
make settle-referral

# Expected output:
# 🤝 Starting referral commission settlement
# Found 1 referrers with commissions
# 1. 0x... (Referrer): 50 KAWAI (10 jobs)
# 🌳 Generating Merkle tree...
# ✅ Merkle Root: 0x...
```

### Step 5: Upload Merkle Root

```bash
MERKLE_ROOT="0x..."  # From settle-referral output

# Upload to contract using Go tool
make upload-merkle-root TYPE=referral ROOT=$MERKLE_ROOT
```

### Step 6: Test Claim in UI

1. Open app as Referrer
2. Go to: **Wallet → Rewards → Referral Rewards**
3. You should see:
   - ✅ Total Referrals: 1
   - ✅ Total Earnings: 150 KAWAI (100 one-time + 50 commission)
   - ✅ Claimable Commission: 50 KAWAI
4. Click "Claim Commission"
5. Sign transaction
6. Wait for confirmation
7. Expected:
   - ✅ Success message
   - ✅ Transaction hash
   - ✅ Balance updated

### ✅ Success Criteria

- Referral code generation works
- One-time bonuses distributed correctly
- Mining commission tracked (5%)
- Settlement aggregates commissions
- Claim transaction succeeds
- KAWAI balance increases

---

## 🧪 Test 4: Revenue Sharing (USDT Dividends)

**Duration:** ~15 minutes  
**Status:** ✅ 95% Ready (Needs testnet testing)

### Step 1: Ensure USDT in PaymentVault

```bash
# Check PaymentVault balance
# Users should have deposited USDT for AI usage
# All USDT in vault = platform revenue
```

### Step 2: Generate Revenue Settlement

```bash
# Settle revenue sharing (USDT dividends)
make settle-revenue

# Expected output:
# 📊 Revenue Sharing Settlement (USDT Dividends)
# 
# Step 1: Generating revenue settlement...
# Current Period:    2
# Settling Period:   1
# ✅ Settlement generated successfully
# Merkle Root: 0xabcd...
# 
# Step 2: Getting vault balance...
# Total Revenue: 1000000000 USDT
# 
# Step 3: Withdrawing USDT to distributor...
# ⚠️  About to withdraw 1000000000 USDT to USDT_Distributor
# Continue with withdrawal? (y/n): y
# ✅ USDT withdrawn successfully
# 
# Step 4: Uploading merkle root...
# ⚠️  About to upload merkle root: 0xabcd...
# Continue with upload? (y/n): y
# ✅ Merkle root uploaded successfully
# ✅ Revenue settlement completed!
```

### Step 3: Test Claim in UI

1. Open app as KAWAI holder
2. Go to: **Wallet → Rewards → Revenue Share**
3. You should see:
   - ✅ KAWAI Balance: X tokens
   - ✅ Share Percentage: Y%
   - ✅ Claimable USDT: Z USDT
4. Click "Claim Dividends"
5. Sign transaction
6. Wait for confirmation
7. Expected:
   - ✅ Success message
   - ✅ USDT balance increased

### ✅ Success Criteria

- Revenue settlement generates valid Merkle tree
- USDT withdrawn from vault successfully
- Merkle root uploaded successfully
- Claim transaction succeeds
- USDT balance increases proportionally

---

## 🧪 Test 5: Unified Settlement (All 4 Types)

**Duration:** ~15 minutes  
**Status:** ✅ 100% Ready

### Step 1: Settle All Types at Once

```bash
# One command to settle all 4 reward types
make settle-all

# Expected output:
# 🚀 Settling All Reward Types
# 1️⃣  Mining Rewards
# ✅ Mining settlement completed!
# 2️⃣  Cashback Rewards
# ✅ Cashback settlement completed!
# 3️⃣  Referral Rewards
# ✅ Referral settlement completed!
# 4️⃣  Revenue Sharing (USDT Dividends)
# [Interactive confirmations for withdraw + upload]
# ✅ Revenue settlement completed!
# 🎉 All settlements completed successfully!
```

### Step 2: Check Status

```bash
make reward-settlement-status

# Expected output:
# ⛏️  MINING REWARDS
# Period 1767549424 | proofs_saved | 1000 KAWAI | 3
# 💰 CASHBACK REWARDS
# Current Period: 53
# 🤝 REFERRAL REWARDS
# Status: Implemented
# 💵 REVENUE SHARING
# Status: Implemented
```

### ✅ Success Criteria

- All 4 settlements run without errors
- Merkle roots generated for each type
- Revenue settlement requires 2 confirmations
- Status command shows all settlements

---

## 🐛 Debugging Tips

### If Claim Fails

1. Check MINTER_ROLE is granted
2. Check Merkle root is uploaded
3. Check proof is valid
4. Check period ID is correct
5. Check user hasn't claimed already
6. Check gas price/limit
7. Check wallet has MON for gas

### If Settlement Fails

1. Check Cloudflare KV credentials
2. Check data exists in KV
3. Check period calculation
4. Check no data for period (empty is OK)

### If UI Doesn't Show Rewards

1. Check desktop app is running (Wails)
2. Check Cloudflare KV credentials are configured
3. Check wallet is connected
4. Check correct network (Monad Testnet)
5. Check browser console for errors in app DevTools

### Common Issues

- **"Proof invalid"** → Merkle root not uploaded yet
- **"Already claimed"** → User claimed in previous test
- **"Insufficient gas"** → Add MON to wallet
- **"No rewards"** → Settlement not run yet
- **"Transaction reverted"** → Check MINTER_ROLE

---

## 📊 Testing Checklist

### Mining Rewards

- [ ] Test data injection works
- [ ] Settlement generates valid Merkle tree
- [ ] Merkle root upload succeeds
- [ ] UI shows claimable rewards
- [ ] Claim transaction succeeds
- [ ] Balance increases correctly
- [ ] hasClaimed returns true
- [ ] UI updates to "Claimed"

### Cashback Rewards

- [ ] Deposit tracking works
- [ ] Cashback calculation correct (tier-based)
- [ ] Settlement generates valid Merkle tree
- [ ] Merkle root upload succeeds
- [ ] UI shows claimable cashback
- [ ] Claim transaction succeeds
- [ ] Balance increases correctly
- [ ] hasClaimed returns true
- [ ] UI updates correctly

### Referral Rewards

- [ ] Referral code generation works
- [ ] Code validation works
- [ ] One-time bonuses distributed
- [ ] Mining commission tracked (5%)
- [ ] Settlement aggregates commissions
- [ ] Merkle root upload succeeds
- [ ] UI shows claimable commission
- [ ] Claim transaction succeeds
- [ ] Balance increases correctly

### Revenue Sharing (USDT Dividends)

- [ ] USDT in PaymentVault tracked
- [ ] Holder scanning works
- [ ] Dividend calculation correct (proportional)
- [ ] Settlement generates valid Merkle tree
- [ ] USDT withdrawal succeeds (with confirmation)
- [ ] Merkle root upload succeeds (with confirmation)
- [ ] UI shows claimable USDT
- [ ] Claim transaction succeeds
- [ ] USDT balance increases correctly

### Unified Tool

- [ ] settle-all runs all 4 types
- [ ] Revenue settlement requires 2 confirmations
- [ ] Status command shows all settlements
- [ ] No errors or warnings
- [ ] All Merkle roots valid

### Code Quality

- [ ] No linter errors
- [ ] Transaction confirmation works
- [ ] Input validation works
- [ ] Negative balance check works
- [ ] Error messages clear
- [ ] Logs informative

---

## 🎯 Next Steps After Testing

1. Fix any bugs found during testing
2. Implement automatic Merkle root upload
3. Setup monitoring for settlements
4. Setup automated weekly settlement cron job
5. Load testing with multiple users
6. Security audit
7. Mainnet deployment

---

**Ready to start testing?** 🚀

Run: `make test-inject-mining-data`

