# MINTER_ROLE Requirements for Kawai Network

## 🎯 Overview

**ALL reward distribution contracts require `MINTER_ROLE`** on the `KawaiToken` contract to mint KAWAI rewards on-demand when users claim.

This is a **critical deployment step** - without it, all reward claims will fail with:
```
Error: AccessControl: account 0x... is missing role 0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6
```

---

## 📊 Contracts Requiring MINTER_ROLE

| Contract | Address (Monad Testnet) | Purpose | Mint Calls per Claim |
|----------|------------------------|---------|---------------------|
| **MiningRewardDistributor** | `0x8117D77A219EeF5F7869897C3F0973Afb87d8427` | Mining rewards with referral splits | 1-4 (contributor, developer, user, affiliator) |
| **CashbackDistributor** | `0xdE64f6F5bEe28762c91C76ff762365D553204e35` | Deposit cashback rewards | 1 |
| **KAWAI_Distributor** | `0xaB0DdFbb4bD94d23a32d0C40f9F96d9A61b45463` | Legacy referral rewards (MerkleDistributor in mint mode) | 1 |

### Contracts NOT Requiring MINTER_ROLE

| Contract | Address | Purpose | Distribution Method |
|----------|---------|---------|-------------------|
| **USDT_Distributor** | `0x98a7590406a08Cc64dc074D8698B71e4D997a268` | USDT dividend distribution | Transfer from pre-funded balance |

---

## 🔍 Why MINTER_ROLE?

### Architecture Decision: Mint-on-Demand vs Pre-Funding

**Mint-on-Demand (Current Approach)** ✅
- **Pros:**
  - Gas-efficient deployment (no need to pre-fund contracts)
  - No risk of contract balance running out
  - Simpler accounting (no need to track contract balances)
  - Aligned with KAWAI tokenomics (gradual emission)
- **Cons:**
  - Requires `MINTER_ROLE` grant (one-time setup)

**Pre-Funding (Alternative)**
- **Pros:**
  - No role management needed
- **Cons:**
  - High upfront gas cost (transfer 1B KAWAI to each contract)
  - Risk of contract balance depletion
  - Complex accounting (need to track remaining balance)
  - Not aligned with gradual emission model

**Decision:** Mint-on-demand is superior for our use case.

---

## 📋 Evidence from Smart Contracts

### 1. MiningRewardDistributor.sol

```solidity
// Line 136-153
function claimReward(...) external nonReentrant {
    // ...
    if (contributorAmount > 0) {
        IMintableToken(address(kawaiToken)).mint(msg.sender, contributorAmount);
        totalContributorRewards += contributorAmount;
    }
    
    if (developerAmount > 0 && developer != address(0)) {
        IMintableToken(address(kawaiToken)).mint(developer, developerAmount);
        totalDeveloperRewards += developerAmount;
    }
    
    if (userAmount > 0) {
        IMintableToken(address(kawaiToken)).mint(user, userAmount);
        totalUserRewards += userAmount;
    }
    
    if (affiliatorAmount > 0 && affiliator != address(0)) {
        IMintableToken(address(kawaiToken)).mint(affiliator, affiliatorAmount);
        totalAffiliatorRewards += affiliatorAmount;
    }
    // ...
}
```

**4 mint calls** per claim (contributor, developer, user, affiliator)!

### 2. DepositCashbackDistributor.sol

```solidity
// Line 119-121
function claimCashback(...) external nonReentrant {
    // ...
    // Mint KAWAI tokens
    IMintableToken(address(kawaiToken)).mint(msg.sender, kawaiAmount);
    totalKawaiDistributed += kawaiAmount;
    // ...
}
```

### 3. MerkleDistributor.sol (KAWAI Mode)

```solidity
// Line 85-91
function claim(...) external {
    // ...
    if (mintOnClaim) {
        // Mint new tokens directly to claimant (gas paid by claimant)
        IMintableToken(address(token)).mint(account, amount);
    } else {
        // Transfer from contract's pre-funded balance
        token.safeTransfer(account, amount);
    }
    // ...
}
```

**KAWAI_Distributor** is deployed with `mintOnClaim = true`, so it needs `MINTER_ROLE`.

---

## 🚀 How to Grant MINTER_ROLE

### Option 1: Automated Script (Recommended)

```bash
# Grant to all contracts at once
export PRIVATE_KEY=0x...
./GRANT_ALL_MINTER_ROLES.sh
```

**Features:**
- ✅ Grants to all 3 contracts automatically
- ✅ Checks existing roles (skips if already granted)
- ✅ Verifies each grant
- ✅ Detailed logging with colors
- ✅ Summary report

### Option 2: Individual Scripts

```bash
# Mining rewards only
./GRANT_MINING_MINTER_ROLE.sh

# Cashback rewards only
./GRANT_CASHBACK_MINTER_ROLE.sh

# Legacy referral rewards only
./GRANT_KAWAI_DISTRIBUTOR_MINTER_ROLE.sh
```

### Option 3: Manual (Foundry)

```bash
# Set environment
export PRIVATE_KEY=0x...
export RPC_URL="https://testnet.monad.xyz/"
export KAWAI_TOKEN="0xE32660b39D99988Df4bFdc7e4b68A4DC9D654722"
export MINTER_ROLE="0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6"

# Grant to MiningRewardDistributor
cast send $KAWAI_TOKEN \
  "grantRole(bytes32,address)" \
  $MINTER_ROLE \
  0x8117D77A219EeF5F7869897C3F0973Afb87d8427 \
  --private-key $PRIVATE_KEY \
  --rpc-url $RPC_URL

# Grant to CashbackDistributor
cast send $KAWAI_TOKEN \
  "grantRole(bytes32,address)" \
  $MINTER_ROLE \
  0xdE64f6F5bEe28762c91C76ff762365D553204e35 \
  --private-key $PRIVATE_KEY \
  --rpc-url $RPC_URL

# Grant to KAWAI_Distributor
cast send $KAWAI_TOKEN \
  "grantRole(bytes32,address)" \
  $MINTER_ROLE \
  0xaB0DdFbb4bD94d23a32d0C40f9F96d9A61b45463 \
  --private-key $PRIVATE_KEY \
  --rpc-url $RPC_URL
```

---

## 🔐 Security Considerations

### Access Control

**KawaiToken** uses OpenZeppelin's `AccessControl`:
- **DEFAULT_ADMIN_ROLE**: Can grant/revoke all roles (deployer)
- **MINTER_ROLE**: Can mint new tokens (distributors only)

### Role Verification

Check if all distributor contracts have `MINTER_ROLE`:

```bash
# Check all distributors at once (recommended)
make check-minter-role
```

**Expected Output:**
```
Checking MINTER_ROLE status for all reward distributors...

✅ MiningRewardDistributor (0x8117D77A219EeF5F7869897C3F0973Afb87d8427): HAS MINTER_ROLE
✅ CashbackDistributor (0xdE64f6F5bEe28762c91C76ff762365D553204e35): HAS MINTER_ROLE
✅ ReferralRewardDistributor (0xaB0DdFbb4bD94d23a32d0C40f9F96d9A61b45463): HAS MINTER_ROLE

All distributors have MINTER_ROLE! ✅
```

**Alternative (Foundry cast):**
```bash
# hasRole(bytes32 role, address account) returns (bool)
cast call 0xE32660b39D99988Df4bFdc7e4b68A4DC9D654722 \
  "hasRole(bytes32,address)" \
  0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6 \
  0x8117D77A219EeF5F7869897C3F0973Afb87d8427 \
  --rpc-url https://testnet.monad.xyz/
```

**Result:**
- `0x0000000000000000000000000000000000000000000000000000000000000001` = **Has role** ✅
- `0x0000000000000000000000000000000000000000000000000000000000000000` = **No role** ❌

### Revoking MINTER_ROLE

If needed (e.g., contract upgrade):

```bash
cast send $KAWAI_TOKEN \
  "revokeRole(bytes32,address)" \
  $MINTER_ROLE \
  <CONTRACT_ADDRESS> \
  --private-key $PRIVATE_KEY \
  --rpc-url $RPC_URL
```

---

## 📊 Tokenomics Alignment

From `README.md`:

```markdown
- **KawaiToken.sol**: ERC20 token with AccessControl (Mint/Burn)
```

**Max Supply:** 1,000,000,000 KAWAI (1 billion)

**Distribution:**
- Mining Rewards: Phase-based diminishing rewards (100→50→25→12 KAWAI per 1M tokens)
- Cashback Rewards: Tier-based rates with per-deposit caps
- Referral Rewards: 10% of referee's mining rewards

Note: No fixed allocation caps per distributor. Only total supply cap of 1B KAWAI (ERC20Capped).

**Emission:**
- Gradual emission via claims (mint-on-demand)
- Halving schedule based on total supply
- No pre-minting (except initial liquidity)

**Why Mint-on-Demand Aligns:**
- ✅ Tokens only minted when users claim (real demand)
- ✅ No inflation from unused allocations
- ✅ Transparent on-chain emission
- ✅ Supports halving schedule dynamically

---

## 🧪 Testing After Grant

### 1. Test Mining Reward Claim

```bash
# In frontend (after settlement)
1. Navigate to Wallet → Rewards → Mining
2. Click "Claim" on any pending reward
3. Confirm transaction
4. Verify KAWAI balance increases
```

### 2. Test Cashback Claim

```bash
# In frontend (after settlement)
1. Navigate to Wallet → Rewards → Cashback
2. Click "Claim" on any pending cashback
3. Confirm transaction
4. Verify KAWAI balance increases
```

### 3. Test Referral Claim

```bash
# In frontend (after settlement)
1. Navigate to Wallet → Rewards → Referral
2. Click "Claim" on any pending commission
3. Confirm transaction
4. Verify KAWAI balance increases
```

### Expected Results

**Before MINTER_ROLE Grant:** ❌
```
Error: execution reverted: AccessControl: account 0x... is missing role 0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6
```

**After MINTER_ROLE Grant:** ✅
```
✅ Transaction confirmed
✅ KAWAI minted to your wallet
✅ Balance updated
```

---

## 📝 Deployment Checklist

- [ ] Deploy `KawaiToken`
- [ ] Deploy `MiningRewardDistributor`
- [ ] Deploy `CashbackDistributor`
- [ ] Deploy `KAWAI_Distributor` (MerkleDistributor with `mintOnClaim=true`)
- [ ] **Grant MINTER_ROLE to MiningRewardDistributor** ← **CRITICAL**
- [ ] **Grant MINTER_ROLE to CashbackDistributor** ← **CRITICAL**
- [ ] **Grant MINTER_ROLE to KAWAI_Distributor** ← **CRITICAL**
- [ ] Test mining claim
- [ ] Test cashback claim
- [ ] Test referral claim
- [ ] Setup weekly settlement automation
- [ ] Launch! 🚀

---

## 🎯 Summary

| Question | Answer |
|----------|--------|
| **Do Mining Rewards need MINTER_ROLE?** | ✅ **YES** |
| **Do Cashback Rewards need MINTER_ROLE?** | ✅ **YES** |
| **Do Referral Rewards need MINTER_ROLE?** | ✅ **YES** (KAWAI mode) |
| **Do USDT Dividends need MINTER_ROLE?** | ❌ NO (pre-funded) |
| **Is this aligned with README.md?** | ✅ **YES** (AccessControl design) |
| **Is this aligned with tokenomics?** | ✅ **YES** (gradual emission) |
| **Is this a security risk?** | ❌ NO (standard OpenZeppelin pattern) |

**Conclusion:** `MINTER_ROLE` is **required, secure, and aligned** with the project design. ✅

---

## 📚 References

- **KawaiToken.sol**: Lines 1-50 (AccessControl implementation)
- **MiningRewardDistributor.sol**: Lines 136-153 (mint calls)
- **DepositCashbackDistributor.sol**: Lines 119-121 (mint call)
- **MerkleDistributor.sol**: Lines 85-91 (mint vs transfer)
- **README.md**: Line 82 (AccessControl mention)
- **OpenZeppelin AccessControl**: https://docs.openzeppelin.com/contracts/4.x/access-control

---

**Status:** Ready for deployment ✅  
**Next Step:** Run `./GRANT_ALL_MINTER_ROLES.sh` 🚀

