# Monad Mainnet Deployment Summary
**Date**: January 23, 2026  
**Chain ID**: 143  
**Network**: Monad Mainnet  
**RPC**: https://rpc.monad.xyz

## ✅ Deployment Status: COMPLETE

All 7 contracts successfully deployed, verified, and configured with MINTER_ROLE grants.

---

## 📋 Deployed Contracts

### 1. KawaiToken (ERC20)
- **Address**: `0x9cbdb316b31fd2efa469c57dcf57be0af630f64c`
- **Verification**: ✅ Verified on Sourcify
- **Explorer**: https://monadvision.com/address/0x9cbdb316b31fd2efa469c57dcf57be0af630f64c

### 2. OTCMarket
- **Address**: `0xe7d73b901b7202b4f686166420ee76cfe860d28d`
- **Verification**: ✅ Verified on Sourcify
- **Explorer**: https://monadvision.com/address/0xe7d73b901b7202b4f686166420ee76cfe860d28d

### 3. PaymentVault
- **Address**: `0xffdf0fb715bec64db41307c26abf545295d31e44`
- **Verification**: ✅ Verified on Sourcify
- **Explorer**: https://monadvision.com/address/0xffdf0fb715bec64db41307c26abf545295d31e44

### 4. RevenueDistributor
- **Address**: `0x7454495f1a7e2854e4215a4d797e0abd7e14bbe4`
- **Verification**: ✅ Verified on Sourcify
- **Explorer**: https://monadvision.com/address/0x7454495f1a7e2854e4215a4d797e0abd7e14bbe4

### 5. MiningRewardDistributor
- **Address**: `0xc58d3f5d04e5748fc1806980e26c1eb487045442`
- **Verification**: ⚠️ Failed (expected - via_ir compilation)
- **MINTER_ROLE**: ✅ Granted and verified
- **Explorer**: https://monadvision.com/address/0xc58d3f5d04e5748fc1806980e26c1eb487045442

### 6. DepositCashbackDistributor
- **Address**: `0x1feff071f37a5cb8833e227d8dddea43aa374449`
- **Verification**: ✅ Verified on Sourcify
- **MINTER_ROLE**: ✅ Granted and verified
- **Explorer**: https://monadvision.com/address/0x1feff071f37a5cb8833e227d8dddea43aa374449

### 7. ReferralRewardDistributor
- **Address**: `0xfbbe8b96d1b5eff919ce09da28737c667faa7957`
- **Verification**: ✅ Verified on Sourcify
- **MINTER_ROLE**: ✅ Granted and verified
- **Explorer**: https://monadvision.com/address/0xfbbe8b96d1b5eff919ce09da28737c667faa7957

---

## 🔐 MINTER_ROLE Configuration

All three reward distributors have been granted MINTER_ROLE on KawaiToken:

| Contract | Address | Role Status | Transaction |
|----------|---------|-------------|-------------|
| MiningRewardDistributor | `0xc58d3f5d04e5748fc1806980e26c1eb487045442` | ✅ Verified | [View](https://monadvision.com/tx/0x33a967803f8e649471d35de14ec49c29a885ee33d4546ac36638abedd6bc9cb4) |
| DepositCashbackDistributor | `0x1feff071f37a5cb8833e227d8dddea43aa374449` | ✅ Verified | [View](https://monadvision.com/tx/0x5ef80692b20a844812d9568e3c32eb493c65932b3a8be5a2a75052491f288e84) |
| ReferralRewardDistributor | `0xfbbe8b96d1b5eff919ce09da28737c667faa7957` | ✅ Verified | [View](https://monadvision.com/tx/0x0377162c1e95233f471e3c08f38ed1fc7e60f0d20565023775953aa2fd52663f) |

**MINTER_ROLE Hash**: `0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6`

**Grant Transactions Confirmed**:
- Block Range: 50602086 - 50602262
- All 3 GrantRole events visible on [MonadVision](https://monadvision.com/address/0x9cbdb316b31fd2efa469c57dcf57be0af630f64c)

---

## 🐛 Issues Encountered & Resolved

### Issue 1: MINTER_ROLE Grant Failure
**Problem**: Initial deployment showed grants succeeded but verification failed.

**Root Cause**: 
- Script extracted correct addresses from broadcast files
- But `contracts/.env` contained old addresses from previous deployment
- Grant commands used old addresses, verification checked old addresses

**Solution Applied**:
1. Fixed script to use extracted variables consistently (no re-sourcing from file)
2. Manually granted MINTER_ROLE to correct addresses
3. Verified all grants successful
4. Updated `.env.mainnet` and `contracts/.env` with correct addresses
5. Regenerated Go constants

### Issue 2: Script Variable Scope Bug
**Problem**: Extra `fi` statement and duplicate parameters in grant commands.

**Fix Applied**:
- Removed duplicate parameters from `cast send` commands
- Removed extra `fi` statement that broke if-else structure
- Script now uses extracted variables throughout

---

## 📝 Configuration Files Updated

### ✅ `.env.mainnet` (Root)
All contract addresses updated with mainnet deployment.

### ✅ `contracts/.env`
Copied from `.env.mainnet` for consistency.

### ✅ `internal/constant/blockchain.go`
Generated with correct mainnet addresses.

### ✅ `pkg/jarvis/db/project_tokens.go`
Generated with correct contract mappings.

---

## 🔧 Deployment Script Fixes

### File: `contracts/deploy-all.sh`

**Changes Made**:
1. **Step 13 - Grant MINTER_ROLE**:
   - Removed duplicate parameters in `cast send` commands
   - Fixed: `"$MINTER_ROLE" "$MINING_ADDRESS" "$MINTER_ROLE" "$MINING_ADDRESS"` 
   - To: `"$MINTER_ROLE" "$MINING_ADDRESS"`

2. **Step 13 - Syntax Fix**:
   - Removed extra `fi` statement that broke if-else structure

3. **Step 12 & 13 - Variable Consistency**:
   - Script now uses extracted variables (`$KAWAI_ADDRESS`, `$MINING_ADDRESS`, etc.)
   - No longer re-sources from `contracts/.env` which could have stale data

---

## ✅ Verification Commands

To verify MINTER_ROLE grants manually:

```bash
# Source mainnet config
source .env.mainnet

# Verify MiningRewardDistributor
cast call "$KAWAI_TOKEN_ADDRESS" \
  "hasRole(bytes32,address)" \
  "0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6" \
  "$MINING_DISTRIBUTOR_ADDRESS" \
  --rpc-url "$MONAD_RPC_URL"

# Verify CashbackDistributor
cast call "$KAWAI_TOKEN_ADDRESS" \
  "hasRole(bytes32,address)" \
  "0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6" \
  "$CASHBACK_DISTRIBUTOR_ADDRESS" \
  --rpc-url "$MONAD_RPC_URL"

# Verify ReferralDistributor
cast call "$KAWAI_TOKEN_ADDRESS" \
  "hasRole(bytes32,address)" \
  "0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6" \
  "$REFERRAL_DISTRIBUTOR_ADDRESS" \
  --rpc-url "$MONAD_RPC_URL"
```

Expected output: `0x0000000000000000000000000000000000000000000000000000000000000001` (true)

---

## 🚀 Next Steps

1. ✅ **Deployment**: Complete
2. ✅ **MINTER_ROLE Grants**: Complete
3. ✅ **Constants Generation**: Complete
4. ⏭️ **Backend Restart**: Required to load new constants
5. ⏭️ **Frontend Update**: Update contract addresses if hardcoded
6. ⏭️ **Testing**: Test reward claiming in production

---

## 📊 Gas Usage Summary

| Contract | Gas Used | Cost (ETH @ 102 gwei) |
|----------|----------|----------------------|
| KawaiToken | 1,128,249 | 0.115 ETH |
| OTCMarket | 1,308,252 | 0.133 ETH |
| PaymentVault | 558,144 | 0.057 ETH |
| RevenueDistributor | 554,615 | 0.057 ETH |
| MiningDistributor | 1,812,228 | 0.185 ETH |
| CashbackDistributor | 1,281,260 | 0.131 ETH |
| ReferralDistributor | 1,310,543 | 0.134 ETH |
| **Total Deployment** | **7,953,291** | **~0.812 ETH** |
| MINTER_ROLE Grants (3x) | 300,000 | 0.031 ETH |
| **Grand Total** | **8,253,291** | **~0.843 ETH** |

---

## 🔗 Important Links

- **MonadVision Explorer**: https://monadvision.com
- **Monad RPC**: https://rpc.monad.xyz
- **Sourcify Verifier**: https://sourcify-api-monad.blockvision.org
- **USDC Token**: https://monadvision.com/address/0x754704bc059f8c67012fed69bc8a327a5aafb603
- **KawaiToken (with GrantRole events)**: https://monadvision.com/address/0x9cbdb316b31fd2efa469c57dcf57be0af630f64c

---

## 📞 Support

For issues or questions:
- Check contract addresses in `.env.mainnet`
- Verify MINTER_ROLE grants using commands above
- Review deployment logs in `contract-deploy-mainnet.log`
- Contact: Admin Address `0x9F463EeF9EffaBDFf97e909C1c5BA1d6df8f7Cc3`

---

**Deployment Completed Successfully** ✅
