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
| **MiningRewardDistributor** | `0xa0dDC59DAcBA9201CC9Ef613707d287b77b2723F` | Mining rewards with referral splits | 1-4 (contributor, developer, user, affiliator) |
| **CashbackDistributor** | `0xcc992d001Bc1963A44212D62F711E502DE162B8E` | Deposit cashback rewards | 1 |
| **KAWAI_Distributor** | `0x988Cbef1F6b9057Cfa7325a7E364543E615f9191` | Legacy referral rewards (MerkleDistributor in mint mode) | 1 |

### Contracts NOT Requiring MINTER_ROLE

| Contract | Address | Purpose | Distribution Method |
|----------|---------|---------|-------------------|
| **USDT_Distributor** | `0xE964B52D496F37749bd0caF287A356afdC10836C` | USDT dividend distribution | Transfer from pre-funded balance |

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
export KAWAI_TOKEN="0x3EC7A3b85f9658120490d5a76705d4d304f4068D"
export MINTER_ROLE="0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6"

# Grant to MiningRewardDistributor
cast send $KAWAI_TOKEN \
  "grantRole(bytes32,address)" \
  $MINTER_ROLE \
  0xa0dDC59DAcBA9201CC9Ef613707d287b77b2723F \
  --private-key $PRIVATE_KEY \
  --rpc-url $RPC_URL

# Grant to CashbackDistributor
cast send $KAWAI_TOKEN \
  "grantRole(bytes32,address)" \
  $MINTER_ROLE \
  0xcc992d001Bc1963A44212D62F711E502DE162B8E \
  --private-key $PRIVATE_KEY \
  --rpc-url $RPC_URL

# Grant to KAWAI_Distributor
cast send $KAWAI_TOKEN \
  "grantRole(bytes32,address)" \
  $MINTER_ROLE \
  0x988Cbef1F6b9057Cfa7325a7E364543E615f9191 \
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

✅ MiningRewardDistributor (0xa0dDC59DAcBA9201CC9Ef613707d287b77b2723F): HAS MINTER_ROLE
✅ CashbackDistributor (0xcc992d001Bc1963A44212D62F711E502DE162B8E): HAS MINTER_ROLE
✅ ReferralRewardDistributor (0x988Cbef1F6b9057Cfa7325a7E364543E615f9191): HAS MINTER_ROLE

All distributors have MINTER_ROLE! ✅
```

**Alternative (Foundry cast):**
```bash
# hasRole(bytes32 role, address account) returns (bool)
cast call 0x3EC7A3b85f9658120490d5a76705d4d304f4068D \
  "hasRole(bytes32,address)" \
  0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6 \
  0xa0dDC59DAcBA9201CC9Ef613707d287b77b2723F \
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

**Allocation:**
- Mining Rewards: 700M KAWAI (70%)
- Cashback Rewards: 200M KAWAI (20%)
- Referral Rewards: 100M KAWAI (10%)

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

