# Testing Guide

**Last Updated:** March 1, 2026  
**Status:** Active testing procedures

---

## 📋 Overview

This guide provides comprehensive testing procedures for the Kawai Network ecosystem, covering smart contracts, backend services, and frontend integration.

---

## 🧪 Automated Tests

### MINTER_ROLE Status Check

**Purpose:** Verify all reward distributors have MINTER_ROLE granted

**Command:**
```bash
make check-minter-role
```

**Expected Output:**
```
Checking MINTER_ROLE status for all reward distributors...

✅ MiningRewardDistributor (0x...): HAS MINTER_ROLE
✅ CashbackDistributor (0x...): HAS MINTER_ROLE
✅ ReferralRewardDistributor (0x...): HAS MINTER_ROLE

All distributors have MINTER_ROLE! ✅
```

**What It Checks:**
- Connects to Monad RPC
- Queries KawaiToken contract for MINTER_ROLE status
- Verifies all 3 reward distributors have the role
- Exit code 0 if all pass, 1 if any fail

**When to Run:**
- After contract deployment
- Before reward settlement
- When debugging claim failures

---

### Balance Checker

**Purpose:** Check KAWAI token balance for any address

**Command:**
```bash
make check-balance ADDR=0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E
```

**Expected Output:**
```
Address: 0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E
KAWAI Balance: 1000.50 KAWAI
Wei Balance: 1000500000000000000000 wei
```

**What It Checks:**
- Connects to Monad RPC
- Queries KawaiToken.balanceOf() for specified address
- Formats output in KAWAI (18 decimals) and wei

**When to Run:**
- After reward claims
- Before/after transactions
- Debugging balance issues

---

### Contract Tests (Foundry)

**Purpose:** Run smart contract unit tests

**Command:**
```bash
cd contracts
forge test -vvv
```

**Run Specific Test:**
```bash
# Test specific contract
forge test --match-contract MiningRewardDistributor -vvv

# Test specific function
forge test --match-test testClaimReward -vvv

# Test with gas report
forge test --gas-report
```

**Expected Output:**
```
[PASS] testClaimReward() (gas: 123456)
[PASS] testMerkleProofVerification() (gas: 78901)
[PASS] testMINTER_ROLERequired() (gas: 45678)

Test result: OK. 13 passed; 0 failed; 0 skipped
```

**When to Run:**
- Before contract deployment
- After contract changes
- CI/CD pipeline

---

### Backend Tests (Go)

**Purpose:** Run Go unit and integration tests

**Command:**
```bash
# All tests
make test

# Specific package
go test ./internal/services/... -v

# Specific test
go test ./internal/services -run TestClaimMiningReward -v

# With coverage
go test ./... -cover
```

**Expected Output:**
```
=== RUN   TestClaimMiningReward
--- PASS: TestClaimMiningReward (0.5s)
=== RUN   TestGetCashbackStats
--- PASS: TestGetCashbackStats (0.3s)
PASS
coverage: 78.5% of statements
```

**When to Run:**
- Before committing code
- After code changes
- CI/CD pipeline

---

## 🔧 Manual Testing Procedures

### Mining Rewards End-to-End

**Purpose:** Test complete mining reward flow from job completion to claim

**Prerequisites:**
- Backend running
- Frontend accessible
- Test wallet with contributor address

**Steps:**

#### 1. Inject Test Mining Data
```bash
make test-inject-mining-data
```

**Expected Output:**
```
📊 Injected 3 test scenarios:
• Referral user: 85 KAWAI (contributor)
• Non-referral user: 90 KAWAI (contributor)
• Multiple jobs: 127.5 KAWAI (3 jobs aggregated)
```

#### 2. Generate Settlement
```bash
make settle-mining
```

**Expected Output:**
```
📊 Mining Rewards Settlement
Period ID:     1
Merkle Root:   0x6f1fd1fc...
Contributors:  3
Total Amount:  302.5 KAWAI
Proofs Saved:  3
Status:        completed
```

#### 3. Upload Merkle Root
```bash
make upload-merkle-root TYPE=mining
```

**Expected Output:**
```
📤 Uploading mining merkle root...
✅ Merkle root uploaded!
Tx: 0xabc123...
Block: 12345
Gas: 45678
```

#### 4. Test Claim in UI
1. Open frontend: `make dev-hot`
2. Navigate to Wallet → Rewards → Mining
3. Click "Claim" button
4. Confirm transaction in wallet
5. Wait for confirmation

#### 5. Verify Balance
```bash
make check-balance ADDR=<CONTRIBUTOR_ADDRESS>
```

**Expected:** KAWAI balance increased by claimed amount

#### 6. Check Claim Status
```bash
make check-claim-status TYPE=mining PERIOD=1 ADDR=<CONTRIBUTOR_ADDRESS>
```

**Expected Output:**
```
✅ Claimed in period 1
Tx: 0xabc123...
Amount: 126 KAWAI
```

---

### Cashback Rewards End-to-End

**Purpose:** Test complete cashback reward flow from deposit to claim

**Prerequisites:**
- Backend running
- Frontend accessible
- Test wallet with user address

**Steps:**

#### 1. Make Test Deposit
```bash
# Using cast (testnet)
cast send $USDC_ADDRESS "transfer(address,uint256)" \
  $PAYMENT_VAULT \
  1000000000 \
  --private-key $PRIVATE_KEY \
  --rpc-url $TESTNET_RPC
```

**Expected:** Deposit transaction confirmed

#### 2. Check Cashback Tracked
```bash
# Check KV store
go run cmd/dev/check-cashback/main.go <USER_ADDRESS>
```

**Expected Output:**
```
💰 Cashback Records:
• Period 1: 30 KAWAI (Tier 3, 1.5%)
Total: 30 KAWAI
```

#### 3. Generate Settlement
```bash
make settle-cashback
```

**Expected Output:**
```
📊 Cashback Settlement
Period ID:     1
Merkle Root:   0xdef456...
Users:         5
Total Amount:  150 KAWAI
Proofs Saved:  5
```

#### 4. Upload Merkle Root
```bash
make upload-merkle-root TYPE=cashback
```

#### 5. Test Claim in UI
1. Navigate to Wallet → Rewards → Cashback
2. Click "Claim" button
3. Confirm transaction
4. Verify KAWAI balance increased

---

### Referral Rewards End-to-End

**Purpose:** Test complete referral reward flow

**Steps:**

#### 1. Create Referral Code
```typescript
// Frontend console or UI
const code = await CreateReferralCode(userAddress);
console.log("Referral code:", code); // e.g., "ABC123"
```

#### 2. Apply Referral Code
```typescript
// In referral banner
await ApplyReferralCode("ABC123");
```

#### 3. Claim Trial Bonus
```typescript
// Auto-claimed on wallet unlock
const [claimed, usdt, kawai] = await AutoClaimTrialIfNeeded("ABC123");
console.log(`Claimed: ${usdt} USDT + ${kawai} KAWAI`);
```

#### 4. Check Referral Stats
```typescript
const stats = await GetReferralStats(userAddress);
console.log(stats);
// {
//   code: "ABC123",
//   total_referrals: 5,
//   total_earnings_usdt: "25.00",
//   total_earnings_kawai: "5000000000000000000000"
// }
```

#### 5. Generate Settlement
```bash
make settle-referral
```

#### 6. Upload Merkle Root
```bash
make upload-merkle-root TYPE=referral
```

#### 7. Claim Commission
1. Navigate to Wallet → Rewards → Referral
2. Click "Claim Commission"
3. Confirm transaction

---

## 🎯 Integration Testing

### Settlement Automation Test

**Purpose:** Test complete settlement flow for all reward types

**Command:**
```bash
# Settle all reward types at once
make settle-all
```

**Expected Output:**
```
📊 Settlement Summary
─────────────────────
Mining:    Period 1, 302.5 KAWAI, 3 contributors
Cashback:  Period 1, 150 KAWAI, 5 users
Referral:  Period 1, 50 KAWAI, 2 referrers
Revenue:   Period 1, 100 USDT, 10 holders
─────────────────────
Total:     502.5 KAWAI + 100 USDT
```

**What It Does:**
1. Generate mining settlement
2. Generate cashback settlement
3. Generate referral settlement
4. Generate revenue settlement
5. Upload all Merkle roots
6. Display summary

**When to Run:**
- Weekly (production)
- After injecting test data
- Before claim testing

---

### Reward Claim Status Check

**Purpose:** Check claim status for any user/reward type

**Command:**
```bash
make check-claim-status TYPE=mining PERIOD=1 ADDR=<ADDRESS>
```

**Parameters:**
- `TYPE`: `mining`, `cashback`, `referral`, or `revenue`
- `PERIOD`: Settlement period number
- `ADDR`: User/contributor address

**Expected Output:**
```
✅ Claimed in period 1
Tx: 0xabc123...
Amount: 126 KAWAI
Timestamp: 2026-01-15 10:30:45
```

**Or (if not claimed):**
```
❌ Not claimed yet
Claimable: 126 KAWAI
Period: 1
```

---

## 📊 Performance Testing

### Load Testing Settlement

**Purpose:** Test settlement performance with large datasets

**Setup:**
```bash
# Inject 1000 test mining records
go run cmd/dev/test-inject-mining-data/main.go --count 1000
```

**Test:**
```bash
# Time settlement generation
time make settle-mining
```

**Expected Performance:**
- Generation: < 10 seconds for 1000 records
- Merkle root upload: < 30 seconds
- Total: < 1 minute

---

### KV Store Performance

**Purpose:** Test KV store query performance

**Test:**
```bash
# Test cashback stats query
go test ./pkg/store -run TestGetCashbackStats_Performance -bench=.
```

**Expected Performance:**
- First query (no cache): < 1 second
- Cached query: < 100ms
- After claim (cache invalidation): < 1 second

---

## 🔍 Debugging Tools

### KV Store Inspector

**Purpose:** Inspect KV store contents for debugging

**Command:**
```bash
# Inspect mining proofs
go run cmd/dev/inspect-proof/main.go mining 1 <ADDRESS>

# List all KV keys
go run cmd/dev/list-kv-keys/main.go mining

# List job rewards
go run cmd/dev/list-job-rewards/main.go <CONTRIBUTOR>
```

**Expected Output:**
```
Merkle Proof for period 1:
Root: 0x6f1fd1fc...
Proof: [0xabc..., 0xdef..., ...]
Leaf: 0x123...
```

---

### Transaction Debugger

**Purpose:** Debug failed transactions

**Command:**
```bash
# Get transaction details
cast tx <TX_HASH> --rpc-url $RPC_URL

# Get transaction receipt
cast receipt <TX_HASH> --rpc-url $RPC_URL

# Decode input data
cast 4byte-decode <INPUT_DATA>
```

**Expected Output:**
```
Transaction 0xabc123...
Status: Success
Block: 12345
Gas Used: 45678
From: 0x123...
To: 0x456...
Input: claimReward(...)
```

---

## ✅ Testing Checklist

### Pre-Deployment
- [ ] Contract tests pass (`forge test`)
- [ ] Backend tests pass (`make test`)
- [ ] Linting passes (`golangci-lint run`)
- [ ] Frontend linter passes (`npm run lint`)

### Post-Deployment
- [ ] MINTER_ROLE granted (`make check-minter-role`)
- [ ] Contract verified on explorer
- [ ] Backend connects successfully
- [ ] Frontend displays correct addresses

### Weekly Settlement
- [ ] Settlement generates successfully
- [ ] Merkle root uploads succeed
- [ ] Proofs saved in KV store
- [ ] Users can claim rewards
- [ ] Balances update correctly

### Monthly Review
- [ ] Review test coverage
- [ ] Update test procedures
- [ ] Archive old test results
- [ ] Add tests for new features

---

## 📚 Related Documentation

- **Contract Guide:** `docs/CONTRACTS_GUIDE.md`
- **Deployment Guide:** `DEPLOYMENT.md`
- **Reward Systems:** `REWARD_SYSTEMS.md`
- **Mining System:** `MINING_SYSTEM.md`
- **Cashback System:** `CASHBACK_SYSTEM.md`
- **Referral System:** `REFERRAL_SYSTEM.md`

---

**Testing Status:** ✅ Active procedures  
**Last Updated:** March 1, 2026
