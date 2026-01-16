# Deposit Cashback System

Earn 1-2% KAWAI tokens back on every USDT deposit with our tiered cashback system!

## 🎯 Overview

Every time you deposit USDT into your Kawai account, you automatically earn KAWAI tokens as cashback. The more you deposit, the higher your cashback rate!

**Key benefits:**
- ✅ Instant cashback calculation
- ✅ Tiered rates (1-2%)
- ✅ First deposit bonus (2.5%)
- ✅ Weekly claims via Merkle proof
- ✅ No action required to accumulate

## 💎 Cashback Tiers

| Tier | Total Deposits | Rate | Max per Deposit | Visual |
|------|----------------|------|-----------------|--------|
| **Bronze** | < 100 USDT | 1% | 5,000 KAWAI | 🥉 |
| **Silver** | 100-500 USDT | 1.25% | 10,000 KAWAI | 🥈 |
| **Gold** | 500-1,000 USDT | 1.5% | 15,000 KAWAI | 🥇 |
| **Platinum** | 1,000-5,000 USDT | 1.75% | 20,000 KAWAI | 💠 |
| **Diamond** | ≥ 5,000 USDT | 2% | 20,000 KAWAI | 💎 |

### How Tiers Work

**Tier is based on cumulative deposits:**
```
Deposit 1: 50 USDT → Bronze tier
Deposit 2: 60 USDT → Still Bronze (total: 110 USDT)
Deposit 3: 40 USDT → Now Silver! (total: 150 USDT)
All future deposits: Silver rate (1.25%)
```

**Tier applies forward, not backward:**
- Your tier is calculated BEFORE each deposit
- New deposit counts toward next tier
- Once unlocked, tier applies to all future deposits

## 🎁 Special Bonuses

### First Deposit Bonus (2.5%)

**Your first deposit ALWAYS gets 2.5% cashback!**

**Examples:**
```
First deposit: 20 USDT
Normal rate: 1% (Bronze) = 200 KAWAI
First-time bonus: 2.5%
You receive: 500 KAWAI! 🎉

First deposit: 75 USDT
Normal rate: 1% (Bronze) = 750 KAWAI
First-time bonus: 2.5%
You receive: 1,875 KAWAI! 🎉

First deposit: 200 USDT
Normal rate: 1.25% (Silver) = 2,500 KAWAI
First-time bonus: 2.5%
You receive: 5,000 KAWAI! 🎉
```

**This is a one-time bonus** - subsequent deposits use tier rates.

### Tier Progression Strategy

**Optimal deposit sequence:**

**Strategy 1: Maximize first deposit**
```
First: 100 USDT → 2,500 KAWAI (2.5%)
Future: 1.25% rate (Silver tier)
```

**Strategy 2: Tier jumping**
```
First: 50 USDT → 1,250 KAWAI (2.5%)
Second: 500 USDT → 7,500 KAWAI (1.5%, reach Gold)
Future: 1.5% rate (Gold tier)
```

**Strategy 3: Big first deposit**
```
First: 5,000 USDT → 20,000 KAWAI (2.5% = 125,000, capped at Diamond 20k)
Future: 2% rate (Diamond tier) - maximum cashback!
```

## 📊 Cashback Calculation

### Formula

```
Raw cashback = depositAmount × rate
KAWAI amount = (raw cashback × 1e18) / 1e6

Example:
Deposit: 100 USDT (6 decimals)
Rate: 1.5% (Gold tier)
Raw: 100 × 0.015 = 1.5 USDT
KAWAI: (1.5 × 1e18) / 1e6 = 1,500 KAWAI (18 decimals)
```

### Tier Caps

Each tier has a maximum cashback per deposit:

```
Bronze: 5,000 KAWAI max
  → Even if calculation is 10,000, you get 5,000

Diamond: 20,000 KAWAI max
  → Deposit 5,000 USDT × 2% = 100,000 raw
  → But capped at 20,000 KAWAI
```

**Why caps exist:**
- Prevent abuse (very large deposits)
- Sustainable tokenomics (200M allocation)
- Fair distribution (~3 year runway)

## 💰 Example Scenarios

### Scenario 1: Small Depositor

**Timeline:**
```
Week 1: First deposit 30 USDT
  → Tier: Bronze
  → Rate: 2.5% (first-time)
  → Cashback: 750 KAWAI

Week 3: Deposit 25 USDT
  → Tier: Bronze (55 USDT total)
  → Rate: 1%
  → Cashback: 250 KAWAI

Week 6: Deposit 50 USDT
  → Tier: Silver! (105 USDT total)
  → Rate: 1.25%
  → Cashback: 625 KAWAI

Total earned: 1,625 KAWAI from 105 USDT deposits
```

### Scenario 2: Strategic Depositor

**Timeline:**
```
Day 1: First deposit 100 USDT
  → Jump straight to Silver
  → Rate: 2.5% (first-time)
  → Cashback: 2,500 KAWAI
  → Future rate: 1.25% (Silver)

Week 2: Deposit 450 USDT
  → Reach Gold (550 USDT total)
  → Rate: 1.5%
  → Cashback: 6,750 KAWAI
  → Future rate: 1.5% (Gold)

Month 2: Deposit 500 USDT
  → Reach Platinum (1,050 USDT total)
  → Rate: 1.75%
  → Cashback: 8,750 KAWAI
  → Future rate: 1.75% (Platinum)

Total earned: 18,000 KAWAI from 1,050 USDT
```

### Scenario 3: Large Depositor

**Strategy:**
```
First deposit: 5,000 USDT
  → Diamond tier immediately
  → Rate: 2.5% (first-time)
  → Raw: 125,000 KAWAI
  → Capped: 20,000 KAWAI
  → Future rate: 2% (Diamond - maximum!)

Monthly deposit: 1,000 USDT
  → Rate: 2% (Diamond)
  → Cashback: 20,000 KAWAI (capped)
  → Repeatable every month
```
Year 1 earnings:
  First: 20,000 KAWAI
  11 months: 11 × 20,000 = 220,000 KAWAI
  Total: 240,000 KAWAI from 16,000 USDT deposits
```

## 🔄 Cashback Lifecycle

### 1. Deposit (You)

```
1. Go to Wallet → Deposit
2. Enter amount
3. Approve USDT spending (one-time)
4. Confirm deposit transaction
5. Wait for confirmation (~10 seconds)
```

### 2. Sync Deposit (You)

```
6. Copy transaction hash from MetaMask
7. Click "Sync Deposit" in app
8. Paste transaction hash
9. Backend verifies on-chain
10. Your balance updates!
```

### 3. Cashback Calculation (Automatic)

```
Backend automatically:
- Checks your tier (total deposits)
- Applies first-time bonus if applicable
- Calculates cashback amount
- Applies tier cap if needed
- Records in database (off-chain)
- Tracks for weekly settlement
```

### 4. Weekly Settlement (Admin)

```
Every Monday 00:00 UTC:
- Collect all pending cashback
- Generate Merkle tree
- Upload Merkle root to blockchain
- Store individual proofs
- Mark period as settled
```

### 5. Claiming (You)

```
After settlement:
- View claimable periods in Rewards → Cashback
- Click "Claim" on any period
- Sign transaction (pays gas)
- Smart contract verifies proof
- KAWAI tokens minted to your address
```

[Detailed claiming guide →](claiming.md)

## 📈 Tracking Your Cashback

### In the App

**Wallet → Rewards → Cashback tab**

**Stats display:**
- **Total Cashback**: All-time earnings
- **Pending**: Waiting for settlement
- **Claimable**: Ready to claim now
- **Claimed**: Already received

**Deposit History table:**
- Date and transaction hash
- Deposit amount (USDT)
- Cashback earned (KAWAI)
- Tier and rate applied
- Status (Pending/Claimable/Claimed)
- Claim button (when available)

### Tier Progress Indicator

Visual representation of your tier:

```
Bronze [====------] Silver [----------] Gold [----------] Platinum [----------] Diamond
   ^
 You are here (45 USDT)
 
Deposit 55 more USDT to reach Silver!
```

## 💡 Maximizing Cashback

### Tip 1: Optimize First Deposit

**Your first deposit is SPECIAL:**
```
Option A: 20 USDT first
  → Get: 1,000 KAWAI
  → Miss: 4,000 KAWAI potential

Option B: 100 USDT first  
  → Get: 5,000 KAWAI (capped at Silver max)
  → Much better!
```

**Recommendation:** Make first deposit at least 100 USDT to maximize 2.5% bonus.

### Tip 2: Plan Tier Jumps

**Batch deposits to reach tiers:**
```
Instead of:
  50 + 30 + 25 = 105 USDT (spread across Bronze)
  
Do:
  100 + 5 = 105 USDT (immediate Silver)
  
Same total, higher rate sooner!
```

### Tip 3: Mind the Caps

**Don't waste on caps:**
```
Diamond (5%, 20K cap):
  
Bad: Deposit 10,000 USDT
  → Cashback: 500,000 raw → 20,000 capped
  → 480,000 KAWAI wasted!
  
Good: Two deposits of 5,000 USDT
  → Cashback: 20,000 each → 40,000 total
  → No waste!
```

**Cap thresholds:**
- Bronze: Cap at 500 USDT deposit (1% of 500 = 5K)
- Silver: Cap at 800 USDT (1.25% of 800 = 10K)
- Gold: Cap at 1,000 USDT (1.5% of 1,000 = 15K)
- Platinum: Cap at 1,142 USDT (1.75% of 1,142 = ~20K)
- Diamond: Cap at 1,000 USDT (2% of 1,000 = 20K)

**Smart strategy:** Deposit in chunks under cap threshold.

### Tip 4: Compound Your Earnings

**Reinvest cashback strategy:**
```
Month 1: Deposit 1,000 USDT → Earn 20,000 KAWAI
Month 2: Sell 10,000 KAWAI → Get ~50 USDT
Month 2: Deposit 1,000 + 50 = 1,050 USDT → Earn more!

Compound effect over year: 10-15% bonus
```

## 🔐 Security & Safety

### Smart Contract Safety

**DepositCashbackDistributor.sol:**
- ✅ Audited and tested (13/13 tests passing)
- ✅ OpenZeppelin standards
- ✅ ReentrancyGuard protection
- ✅ Role-based access control

**What it does:**
- Verify Merkle proofs
- Mint KAWAI on valid claims
- Track claimed periods
- Enforce allocation cap (200M)

**What it CANNOT do:**
- Change your cashback amounts
- Steal your tokens
- Modify claimed status maliciously

### Verification

**You can verify everything on-chain:**

1. **Check deposit on blockchain:**
   - Go to [Monad Explorer](https://testnet.monad.xyz)
   - Paste your deposit transaction
   - Verify amount and recipient

2. **Check cashback calculation:**
   - Formula is public: `(deposit × rate) / 10000`
   - Tier caps are hardcoded
   - First-time bonus logic is transparent

3. **Check Merkle root:**
   - View current period root on contract
   - Matches your off-chain proof
   - Anyone can verify proofs

## 🆘 Troubleshooting

### "Wrong cashback amount calculated"

**Check these factors:**
```
✓ What's your tier? (Bronze/Silver/Gold/Platinum/Diamond)
✓ Is this your first deposit? (5% override)
✓ Did you hit the tier cap? (5K-20K max)
✓ Account for decimals (18 decimals = big numbers)
```

**Example confusion:**
```
Deposit: 5,000 USDT
Tier: Diamond (5%)
Expected: 250,000 KAWAI
Actually got: 20,000 KAWAI

Explanation: Tier cap applied (20K max for Diamond)
This is correct! ✅
```

### "Cashback not showing in pending"

**Possible causes:**
1. Didn't sync deposit (required step!)
2. Transaction still confirming
3. Wrong network (check Monad)
4. Cache delay (wait 1 minute)

**Solutions:**
- Click "Sync Deposit" and paste transaction hash
- Wait for blockchain confirmation (10-30 seconds)
- Refresh the page
- Check deposit history table

### "Can't claim my cashback"

**Check:**
- Is it after Monday settlement?
- Is status "Claimable" (not "Pending")?
- Do you have MON for gas?
- Is wallet still connected?

**Wait times:**
```
Deposit: Tuesday
  → Status: Pending
  → Next Monday: Settlement runs
  → Monday onwards: Status → Claimable
  → You can claim anytime after Monday
```

## ✅ Quick Reference

### Tier Thresholds
- Bronze: 0-99 USDT
- Silver: 100-499 USDT
- Gold: 500-999 USDT
- Platinum: 1,000-4,999 USDT
- Diamond: 5,000+ USDT

### Cashback Rates
- Bronze: 1%
- Silver: 1.25%
- Gold: 1.5%
- Platinum: 1.75%
- Diamond: 2%
- First deposit: **2.5%** (override)

### Tier Caps
- Bronze: 5,000 KAWAI
- Silver: 10,000 KAWAI
- Gold: 15,000 KAWAI
- Platinum: 20,000 KAWAI
- Diamond: 20,000 KAWAI

## 🚀 Next Steps

1. **[Make Your First Deposit](../user-guide/deposit.md)** - Get that 5% bonus!
2. **[Track Your Rewards](overview.md)** - Monitor all earnings
3. **[Learn to Claim](claiming.md)** - Get your KAWAI tokens
4. **[Understand Tokenomics](../tokenomics/kawai-token.md)** - KAWAI value

---

**Questions?** Check [FAQ](../faq/rewards.md) or [contact support](../support/contact.md).

