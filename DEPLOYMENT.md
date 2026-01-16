# Deployment Guide

**Current Environment:** Production on Monad Testnet

For testnet vs mainnet differences, see `.env.README`

---

## CURRENT DEPLOYMENT ADDRESSES
```
MockUSDT:                    0x3AE05118C5B75b1B0b860ec4b7Ec5095188D1CCc
KawaiToken:                  0xf68910e8d19047A309f989FFB515E44FBca5D31A
KAWAI_Distributor:           0x2B11e8385A859Ea75C77E05Bc0D9756A87017E92
USDT_Distributor:            0x896fB97f81ECBEfDBe29DCc3550aC984704932bF
PaymentVault:                0xDA94C8ac2a61eafBd47853EE22702BDCd45B6d93
OTCMarket:                   0x3E597D76B40004c3fC517C404037fD6F16C8fc34
MiningRewardDistributor:     0x1f78c7c472205F1720aAb66a565981561b5EBac0 ⭐ WITH PAUSE
DepositCashbackDistributor:  0x3d5Bfe788782A90ac124096296B45eaFFc43C79B ⭐ WITH PAUSE
ReferralRewardDistributor:   0x1c218602218745B20CE201948CaE836f8E94E111 ⭐ WITH PAUSE
```

**Pause Mechanism Test Results:**
- ✅ Pause status check: Working (all 3 distributors)
- ✅ Pause all distributors: Working (Mining + Cashback + Referral)
- ✅ Unpause all distributors: Working
- ✅ Individual pause/unpause: Working
- ✅ Emergency pause ready for production

**Previous E2E Test Results (Round 4):**
- ✅ Mining rewards: 2 users claimed 450 KAWAI each via UI
- ✅ Cashback rewards: 2 users claimed ~30 KAWAI each via UI
- ✅ Auto-sync: KV automatically syncs with on-chain claimed status
- ✅ Multi-user: Both users can claim simultaneously without conflicts
- ✅ Total test time: ~30 minutes (full fresh deployment)

## GOAL
Deploy semua contracts fresh, inject test data, upload Merkle roots, dan TEST CLAIMING dari UI sampai berhasil.

---

## STEP 1: CLEANUP KV - DELETE ALL DATA (1 menit)
```bash
go run cmd/dev/cleanup-kv-all/main.go --confirm-delete-all
```

**⚠️ DANGER: Ini akan DELETE SEMUA DATA dari SEMUA KV namespaces!**

**Data yang akan di-DELETE:**
- Contributors (job rewards, balances, heartbeat)
- Proofs (Merkle proofs untuk semua periods)
- Settlements (settlement period metadata)
- Cashback (semua cashback records)
- Holders (holder data)
- P2P Marketplace (marketplace data)

**Output yang diharapkan:**
```
🗑️  Cleaning Contributors namespace...
   ✅ Deleted X keys from Contributors
🗑️  Cleaning Proofs namespace...
   ✅ Deleted X keys from Proofs
🗑️  Cleaning Settlements namespace...
   ✅ Deleted X keys from Settlements
🗑️  Cleaning Cashback namespace...
   ✅ Deleted X keys from Cashback
🗑️  Cleaning Holders namespace...
   ✅ Deleted X keys from Holders
🗑️  Cleaning P2P Marketplace namespace...
   ✅ Deleted X keys from P2P Marketplace

✅ COMPLETE CLEANUP FINISHED!
```

**Confirm:** SEMUA data ter-DELETE (bukan mark as settled, tapi REAL deletion via Cloudflare API)

---

## STEP 2: DEPLOY ALL CONTRACTS (5 menit)

### 2.1 Deploy main contracts
```bash
make contracts-deploy-testnet
```

**Output yang diharapkan:**
- 6 contracts deployed (MockUSDT, KawaiToken, KAWAI_Distributor, USDT_Distributor, PaymentVault, OTCMarket)
- Verification otomatis di explorer

**⚠️ IMPORTANT: Catat KawaiToken address dari output!**

### 2.2 Update contracts/.env dengan KawaiToken address
```bash
# Edit contracts/.env, set:
KAWAI_TOKEN_ADDRESS=<KAWAI_TOKEN_ADDRESS_FROM_STEP_2.1>
```

### 2.3 Deploy MiningRewardDistributor (setelah update .env)
```bash
make contracts-deploy-mining-testnet
```

### 2.4 Deploy DepositCashbackDistributor
```bash
make contracts-deploy-cashback-testnet
```

**✅ ACTUAL DEPLOYMENT (2026-01-13):**
```
MockUSDT: 0xa70e7C98331c90C15d3bd6974816BEDb8Da3388a
KawaiToken: 0xD27123B9e723372Fd1bC481bB2Ab9c72C0B160E3
KAWAI_Distributor: 0x2Acc6b2D50979A3891A84887890a7B7Cda3a98c4
USDT_Distributor: 0x44e31f2C0155955E0a5213E2924105FcAC620c2a
PaymentVault: 0x8440e43473F579068bd53768B361327715b94f84
OTCMarket: 0xc9A30Ae0226684B3257FCBbCDc47b073fbacc18B
MiningRewardDistributor: 0xcEDb9c0e7648623e5D8992402970b54Da0DF52ce
DepositCashbackDistributor: 0xEc24361BE5B20d1fe0B5000Af53631C995F7Ac43
```

---

## STEP 3: UPDATE CODE DENGAN ADDRESSES BARU (30 detik)

**Auto-generate dari .env:**
```bash
go run cmd/obfuscator-gen/main.go
```

**Files yang akan di-generate:**
1. `internal/constant/blockchain.go` - backend constants
2. `pkg/jarvis/db/project_tokens.go` - Jarvis token mapping

**Note:** Pastikan `.env` dan `contracts/.env` sudah diupdate dengan addresses dari STEP 2

---

## STEP 4: GRANT MINTER_ROLE (2 menit)

**Grant roles:**
```bash
make contracts-grant-minter-mining
make contracts-grant-minter-cashback
```

**✅ SUDAH GRANTED:**
- Mining: 0xcEDb9c0e7648623e5D8992402970b54Da0DF52ce
- Cashback: 0xEc24361BE5B20d1fe0B5000Af53631C995F7Ac43

**Verify (optional):**
```bash
go run cmd/dev/check-minter-role/main.go
```

---

## STEP 5-7: SETUP TEST USER (1 menit)

**Single command untuk generate wallet + send MON + inject mining data:**
```bash
go run cmd/dev/setup-test-user/main.go
```

**Output:**
```
✅ TEST USER SETUP COMPLETE!

📝 Test User Details:
   Address:     0xA051aB7126E518cD519F7838d1F8D18c2f65886a
   Private Key: 0xd3693013800e0b73360fddeed04aea14590d35e7e367ec9b48b884fbd4f4c0e0
```

**⚠️ PENTING: Simpan private key ini untuk import ke UI nanti!**

---

## STEP 5B: INJECT CASHBACK DATA (Optional - 1 menit)

**Inject test cashback data untuk test user:**
```bash
go run cmd/dev/inject-cashback-data/main.go <TEST_ADDRESS>
```

**Output:**
```
💰 Injecting cashback data for: 0xA051aB7126E518cD519F7838d1F8D18c2f65886a
📅 Current period: 54
📅 Injecting into period: 53 (for settlement)

✅ Deposit 1: 100 USDT
✅ Deposit 2: 500 USDT
✅ Deposit 3: 1000 USDT

✅ Cashback data injected!

📊 Summary:
   Total Deposits: 3
   Total Cashback: ~26 KAWAI
   Pending: 26250000000000000000 wei

Next steps:
  1. Run settlement: make settle-cashback
  2. Upload Merkle root: go run cmd/reward-settlement/main.go upload --type cashback
  3. Test claiming via UI
```

**Note:** Ini akan inject 3 deposits dengan tier berbeda (Tier 2, 3, 4) untuk testing cashback system.

---

## STEP 8: GENERATE SETTLEMENT (2 menit)

### 8.1 Mining Settlement
```bash
go run cmd/reward-settlement/main.go generate --type mining
```

**✅ ACTUAL SETTLEMENT:**
```
Period ID:     1768239231
Merkle Root:   0xabd4226e7f3d22f27d36ce6a3c21d6ee86c5683d0215064c7b25423a5b3f6f06
Contributors:  1
Proofs Saved:  1
```

### 8.2 Cashback Settlement (Optional)
```bash
make settle-cashback
```

**Output:**
```
🌳 Generating cashback settlement...
📊 Cashback Rewards Settlement
─────────────────────────────
Current Period:    54
Settling Period:   53

🔄 [CashbackSettlement] Starting settlement for period 53
📊 [CashbackSettlement] Collected 1 cashback records
🌳 [CashbackSettlement] Merkle root: db8ac17f84d626a24689fec81b27b9f70fc01e71e8a9be7fba1477e905d6b98b
✅ [CashbackSettlement] Stored 1 proofs
📝 [CashbackSettlement] Advance period tx: 0xe2c9be9ba691cf5bdf4bcf948f3cea614a1fcc5d765dbf45c1e26f5587f8e948
✅ [CashbackSettlement] Merkle root set on-chain (block 5776741)
✅ [CashbackSettlement] Settlement complete for period 53
```

**Note:** Cashback settlement sudah otomatis upload Merkle root on-chain (tidak perlu STEP 9 untuk cashback).

---

## STEP 9: UPLOAD MERKLE ROOT (1 menit)
```bash
echo "y" | go run cmd/reward-settlement/main.go upload --type mining
```

**✅ SUDAH UPLOAD:** 
- Transaction: 0x8b5d343444bc5f49a1aa5cd3e5bddfd278bd6dc9fc3b22d48eb17135e9f9f476
- Block: 5746170

**Verify on-chain (optional):**
```bash
cast call 0xcEDb9c0e7648623e5D8992402970b54Da0DF52ce "periodMerkleRoots(uint256)(bytes32)" 1 --rpc-url https://testnet-rpc.monad.xyz
```

Should return: `0xabd4226e7f3d22f27d36ce6a3c21d6ee86c5683d0215064c7b25423a5b3f6f06`

---

## STEP 10: VERIFY CLAIMABLE DATA (1 menit)
```bash
go run cmd/dev/test-claiming-data/main.go 0xbFD83eB024C067889f2FF60Bc2181F9aEc6eAB92
```

**✅ VERIFIED:**
```
✅ Found 1 Merkle proofs:
   1. Type: kawai, Period: 1768239231, Amount: 450000000000000000000, Index: 0

✅ Address has 1 claimable proofs ready!
✅ Ready for on-chain claiming (requires MON for gas)
```

---

## STEP 11: TEST CLAIMING VIA UI (5 menit)

### 11.1 Import Wallet
1. Buka UI Wallet
2. Import private key dari STEP 5-7
3. Verify address matches

### 11.2 Test Mining Rewards

#### Check Claimable Rewards
- Navigate ke "Rewards" → "Mining Rewards"
- Harus show: **450 KAWAI claimable**
- Harus show: **1 unclaimed reward** (Period 1)

#### Claim Rewards
1. Click **"Claim"** button
2. Confirm transaction (gas fee ~0.003 MON)
3. Transaction submitted

#### Verify Success
**Expected behavior:**
- Status: **Unclaimed** → **Confirmed** (langsung, ~1-2 detik)
- Reward muncul di **"Recent Mining Activity"**
- KAWAI balance naik jadi **450 KAWAI**

**Note:** Monad testnet sangat cepat, jadi tidak ada "Pending" state. Tx langsung confirmed.

#### Check Explorer
- Click tx hash link di "Recent Mining Activity"
- Verify status: **Success** ✅
- Verify events: `RewardClaimed` event dengan correct amounts

### 11.3 Test Cashback Rewards (Optional)

#### Check Claimable Cashback
- Navigate ke "Rewards" → "Deposit Cashback"
- Harus show:
  - **Total Cashback Earned:** ~26 KAWAI
  - **Claimable Now:** ~26 KAWAI
  - **Current Tier:** Tier 2-4 (based on deposits)
  - **1 claimable record** in table (Period 53)

#### Claim Cashback
1. Click **"Claim"** button di table
2. Confirm transaction (gas fee ~0.003 MON)
3. Transaction submitted

#### Verify Success
**Expected behavior:**
- Status: **Ready to Claim** → **Claimed** (langsung, ~1-2 detik)
- KAWAI balance naik jadi **476 KAWAI** (450 mining + 26 cashback)
- Claimed cashback muncul di stats

#### Check Explorer
- Click tx hash link
- Verify status: **Success** ✅
- Verify events: `RewardClaimed` event dari DepositCashbackDistributor

---

## STEP 12: VERIFY ON-CHAIN (Optional)

```bash
# Check if claimed
cast call 0x86b11B1A7e4e40D181ac06070a0e98648dBc7859 \
  "hasClaimed(uint256,address)(bool)" \
  1 \
  0x755FE4b121B1945f812513FfBDf54f68fcd54b72 \
  --rpc-url https://testnet-rpc.monad.xyz
```

Expected: `true`

---

```bash
# Check KAWAI balance
export KAWAI_TOKEN=<ADDRESS_FROM_STEP_2>
cast call $KAWAI_TOKEN "balanceOf(address)(uint256)" $TEST_ADDRESS --rpc-url https://testnet-rpc.monad.xyz
```

Harus return: ~450000000000000000000 (450 KAWAI dalam wei)

---

## SUCCESS CRITERIA

### Core Setup
✅ All contracts deployed
✅ MINTER_ROLE granted
✅ Test wallet created dengan MON
✅ Backend running

### Mining Rewards
✅ Mining data injected
✅ Mining settlement generated
✅ Mining Merkle root uploaded
✅ Mining claimable data verified
✅ UI shows mining claimable rewards
✅ Mining claim transaction successful
✅ KAWAI balance increased (mining)

### Cashback Rewards (Optional)
✅ Cashback data injected
✅ Cashback settlement generated (auto-upload Merkle root)
✅ UI shows cashback claimable rewards
✅ Cashback claim transaction successful
✅ KAWAI balance increased (cashback)

### Verification
✅ On-chain verification passed
✅ Explorer shows successful transactions

---

## JIKA GAGAL

### Mining Rewards

#### Claiming gagal dengan "Invalid period"
- Check: Apakah Merkle root ter-upload untuk period tersebut?
- Fix: Run step 9 lagi

#### Claiming gagal dengan "Invalid proof"
- Check: Apakah Merkle proof di KV match dengan on-chain root?
- Fix: Re-generate settlement (step 8-9)

#### UI tidak show claimable rewards
- Check: Apakah backend running?
- Check: Apakah wallet address benar?
- Fix: Restart backend, verify address

### Cashback Rewards

#### Cashback settlement gagal dengan "key not found"
- Check: Apakah cashback data ter-inject ke period yang benar?
- Fix: Re-run inject tool, pastikan inject ke period - 1

#### Cashback settlement gagal dengan "invalid hex character 'x' in private key"
- Check: Private key format di .env
- Fix: Sudah di-fix di `pkg/blockchain/cashback_settlement.go` (auto-strip 0x prefix)

#### UI tidak show cashback rewards
- Check: Apakah backend running?
- Check: Apakah CashbackService API endpoints working?
- Fix: Restart backend, check logs

### General

#### Transaction reverted
- Check: Apakah ada MON untuk gas?
- Check: Apakah MINTER_ROLE granted?
- Fix: Send more MON, re-grant MINTER_ROLE

#### Backend error "failed to get marketplace data from KV"
- Check: Apakah menggunakan `GetCashbackData()` bukan `GetMarketplaceData()`?
- Fix: Sudah di-fix di `pkg/blockchain/cashback_settlement.go`

---

## TOTAL TIME: ~30 menit

**Breakdown:**
- STEP 1: Cleanup KV (1 min)
- STEP 2: Deploy contracts (5 min)
- STEP 3: Update code (30 sec)
- STEP 4: Grant MINTER_ROLE (2 min)
- STEP 5-7: Setup test user (1 min)
- STEP 5B: Inject cashback (1 min) - Optional
- STEP 8: Generate settlement (2 min mining + 2 min cashback)
- STEP 9: Upload Merkle root (1 min)
- STEP 10: Verify claimable (1 min)
- STEP 11: Test claiming UI (5 min)
- STEP 12: Verify on-chain (1 min)

**Total:** ~22 min (without cashback) or ~30 min (with cashback)

## FILES MODIFIED IN THIS DEPLOYMENT

### Auto-Generated (via cmd/obfuscator-gen/main.go)
1. `internal/constant/blockchain.go` - contract addresses
2. `pkg/jarvis/db/project_tokens.go` - token addresses

### Manual Updates Required
3. `.env` - backend config (contract addresses)
4. `contracts/.env` - deployment config (KAWAI_TOKEN_ADDRESS after step 2.1)

### New Tools Created
5. `cmd/dev/setup-test-user/main.go` - All-in-one test user setup
6. `cmd/dev/inject-cashback-data/main.go` - Inject cashback test data
7. `cmd/obfuscator-gen/main.go` - Auto-generate blockchain constants

### Bug Fixes
8. `pkg/blockchain/cashback_settlement.go` - Fixed GetMarketplaceData → GetCashbackData
9. `pkg/blockchain/cashback_settlement.go` - Fixed private key parsing (strip 0x prefix)
10. `cmd/dev/inject-cashback-data/main.go` - Inject into period - 1 for settlement
11. `pkg/store/settlement.go` - Auto-check tx confirmation for mining claims
12. `frontend/src/app/wallet/components/rewards/MiningRewardsSection.tsx` - Auto-refresh when pending claims exist

## COMMANDS SUMMARY
```bash
# 1. Cleanup KV - DELETE ALL DATA
go run cmd/dev/cleanup-kv-all/main.go --confirm-delete-all

# 2. Deploy contracts
make contracts-deploy-testnet
make contracts-deploy-mining-testnet
make contracts-deploy-cashback-testnet

# 3. Update addresses (AUTO-GENERATE)
go run cmd/obfuscator-gen/main.go

# 4. Grant MINTER_ROLE
make contracts-grant-minter-mining
make contracts-grant-minter-cashback

# 5-7. Setup test user (single command)
go run cmd/dev/setup-test-user/main.go

# 5B. Inject cashback data (optional)
go run cmd/dev/inject-cashback-data/main.go <TEST_ADDRESS>

# 8. Generate settlement
go run cmd/reward-settlement/main.go generate --type mining
make settle-cashback  # Optional: cashback settlement (auto-upload)

# 9. Upload Merkle root (mining only)
echo "y" | go run cmd/reward-settlement/main.go upload --type mining

# 10. Verify claimable data
go run cmd/dev/test-claiming-data/main.go <TEST_ADDRESS>

# 11. Start backend
make dev-hot

# 12. Test claiming via UI (import private key)
# - Mining Rewards tab: Claim 450 KAWAI
# - Deposit Cashback tab: Claim ~26 KAWAI (if injected)

# 13. Verify on-chain
make check-balance ADDR=<TEST_ADDRESS>
make check-claim-status TYPE=mining PERIOD=<PERIOD_ID> ADDR=<TEST_ADDRESS>
```

---

## NEXT: EXECUTE STEP BY STEP
Mulai dari step 1, jangan skip, jangan bikin dokumen lagi.
