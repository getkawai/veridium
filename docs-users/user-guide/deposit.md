# Deposit USDT Guide

Learn how to deposit USDT into your Kawai account and start using AI services.

## 🎯 Why Deposit USDT?

USDT is the payment currency for AI services on Kawai Network. With USDT, you can:

- 💬 **Chat with AI**: Pay per request
- 🎨 **Generate Images**: Create AI art
- 🎁 **Earn Cashback**: Get 1-5% KAWAI tokens back
- 💰 **Use-to-Earn**: 5% cashback on every AI request

## 💎 Cashback Tiers

The more you deposit, the higher your cashback rate!

| Tier | Deposit Amount | Base Cashback | Max KAWAI per Deposit | First Deposit Bonus |
|------|----------------|---------------|-----------------------|--------------------|
| 🥉 **Bronze** | < 100 USDT | **1%** | 5,000 KAWAI | **5%** |
| 🥈 **Silver** | 100-500 USDT | **2%** | 10,000 KAWAI | **5%** |
| 🥇 **Gold** | 500-1,000 USDT | **3%** | 15,000 KAWAI | **5%** |
| 💠 **Platinum** | 1,000-5,000 USDT | **4%** | 20,000 KAWAI | **5%** |
| 💎 **Diamond** | ≥ 5,000 USDT | **5%** | 20,000 KAWAI | **5%** |

!!! tip "First Deposit Bonus"
    Your **first deposit always gets 5% cashback**, regardless of amount! Example: Deposit 50 USDT, get 2,500 KAWAI bonus!

## 📋 Prerequisites

Before depositing, make sure you have:

- ✅ **Wallet connected** - See [Wallet Setup](wallet-setup.md)
- ✅ **USDT tokens** on Monad network
- ✅ **MON tokens** for gas fees (~0.001 MON per transaction)

## 💰 Getting USDT on Monad

### Option 1: Bridge from Another Chain

1. Use a bridge service (e.g., LayerZero, Wormhole)
2. Bridge USDT from Ethereum, Polygon, or BSC
3. Select Monad as destination
4. Wait for confirmation (5-10 minutes)

### Option 2: Buy from Testnet Faucet (Testnet Only)

1. Visit [Monad Testnet Faucet](https://testnet.monad.xyz/faucet)
2. Request testnet USDT
3. Use for testing

### Option 3: P2P Marketplace (After Launch)

1. Buy KAWAI from marketplace
2. Trade KAWAI ↔ USDT with other users
3. Use USDT for AI services

## 📝 Step-by-Step Deposit Guide

### Step 1: Open Deposit Page

1. Click **"Wallet"** in sidebar
2. Click **"Deposit"** tab
3. You'll see the deposit interface

### Step 2: Enter Deposit Amount

1. **Enter amount** in USDT
   - Minimum: 1 USDT
   - Maximum: No limit

2. **Check your cashback preview**
   - Shows tier and cashback rate
   - Shows KAWAI tokens you'll earn
   - Shows your new tier (if upgrading)

**Example:**
```
Deposit: 100 USDT
Tier: Silver (2% cashback)
Cashback: 2,000 KAWAI tokens
First deposit: Yes (5% bonus applies)
Actual cashback: 5,000 KAWAI tokens! 🎉
```

### Step 3: Approve USDT Spending

!!! info "One-Time Setup"
    You only need to approve once. After that, you can deposit anytime without approving again.

1. Click **"Approve USDT"**
2. MetaMask opens → Review approval
3. Click **"Approve"**
4. Wait for confirmation (~10 seconds)

**What this does:**
- Allows the PaymentVault contract to move USDT on your behalf
- Secure and standard procedure
- You can revoke approval anytime

### Step 4: Confirm Deposit

1. Click **"Deposit"** button
2. MetaMask opens → Review transaction
   - **From:** Your address
   - **To:** PaymentVault contract
   - **Amount:** Your deposit amount
   - **Gas fee:** ~0.001 MON

3. Click **"Confirm"**
4. Wait for confirmation (~10 seconds)

### Step 5: Sync Deposit

!!! warning "Important"
    After depositing, you must sync the deposit in the app to update your balance.

1. **Wait for transaction confirmation** (green checkmark in MetaMask)
2. **Copy transaction hash** from MetaMask
3. **In Kawai app:** Click "Sync Deposit"
4. **Paste transaction hash**
5. Click **"Sync"**
6. Wait a few seconds
7. **Your balance updates!**

**Why sync?**
- Backend verifies your deposit on-chain
- Prevents fake/invalid deposits
- Updates your internal balance
- Calculates and records cashback

## 💰 Understanding Cashback

### How Cashback Works

```
1. You deposit USDT
   ↓
2. Backend calculates cashback (off-chain)
   - Checks your tier
   - Applies first-time bonus if applicable
   - Applies tier caps (5K-20K KAWAI)
   ↓
3. Cashback recorded in database
   - Status: "Pending"
   - Period: Current week
   ↓
4. Weekly settlement (every Monday)
   - Admin generates Merkle tree
   - Uploads proof to blockchain
   ↓
5. You claim your KAWAI tokens
   - Click "Claim" in Rewards dashboard
   - Receive tokens instantly
```

### When Can I Claim?

- **Accumulation**: Immediate (recorded right after deposit)
- **Settlement**: Every Monday at 00:00 UTC
- **Claiming**: Anytime after settlement

**Example Timeline:**
```
Wednesday 2PM: Deposit 500 USDT → Get 10,000 KAWAI pending
Monday 12AM: Settlement runs → Proof generated
Monday 9AM: You claim → Receive 10,000 KAWAI! 🎉
```

### Checking Cashback Status

1. Go to **Wallet → Rewards → Cashback**
2. You'll see:
   - **Total Cashback**: All-time earnings
   - **Pending**: Waiting for settlement
   - **Claimable**: Ready to claim now
   - **Claimed**: Already received

3. **Deposit History Table** shows:
   - Date and transaction
   - Deposit amount
   - Cashback earned
   - Tier and rate
   - Claim status

## 📊 Deposit Examples

### Example 1: First-Time Small Deposit

**Scenario:**
- First deposit ever
- Amount: 50 USDT

**Result:**
- Normal rate: 1% (Bronze) = 500 KAWAI
- First-time bonus: 5% override
- **You get: 2,500 KAWAI!** 🎉
- Total: 50 USDT balance + 2,500 KAWAI pending

### Example 2: Upgrading Tiers

**Scenario:**
- Already deposited 80 USDT (Bronze tier)
- New deposit: 150 USDT

**Result:**
- New tier: Silver (100+ USDT total)
- This deposit: 2% rate = 3,000 KAWAI
- **You get: 3,000 KAWAI**
- Balance: 230 USDT total
- Future deposits: 2% rate (Silver)

### Example 3: Large Deposit

**Scenario:**
- Deposit: 5,000 USDT

**Result:**
- Tier: Diamond (5%)
- Raw calculation: 5,000 × 5% = 250,000 KAWAI
- **Tier cap applied: 20,000 KAWAI** (Diamond max)
- Still excellent value!

### Example 4: Multiple Small Deposits

**Scenario:**
- Deposit 1: 30 USDT (first-time) → 1,500 KAWAI (5%)
- Deposit 2: 40 USDT → 400 KAWAI (1%, Bronze)
- Deposit 3: 50 USDT → 500 KAWAI (1%, Bronze)

**Result:**
- Total deposited: 120 USDT
- Total cashback: 2,400 KAWAI
- **Next deposit: Silver tier (2%)!**

## 🔐 Security & Safety

### Smart Contract Security

- ✅ **Audited code**: Contracts tested thoroughly
- ✅ **Open source**: View on [GitHub](https://github.com/kawai-network/veridium)
- ✅ **Verified on explorer**: Check contract code
- ✅ **No admin withdrawal**: Only you can withdraw your funds

### Your Funds Are Safe

**What we CAN'T do:**
- ❌ Withdraw your USDT
- ❌ Transfer your funds
- ❌ Lock your account
- ❌ Change your balance

**What you CAN do:**
- ✅ Withdraw anytime (coming soon)
- ✅ Check balance on-chain
- ✅ Verify transactions on explorer

### Transaction Verification

Always verify your deposits:

1. **Check transaction on explorer**:
   - Go to [Monad Explorer](https://testnet.monad.xyz)
   - Paste your transaction hash
   - Verify: status, amount, contract

2. **Check balance on-chain**:
   - Open PaymentVault contract
   - Call `getBalance(yourAddress)`
   - Should match your app balance

## 💡 Tips & Best Practices

### Maximize Your Cashback

1. **First deposit matters**: Make it count (5% bonus!)
2. **Reach tier thresholds**: 100, 500, 1K, 5K USDT
3. **Batch your deposits**: One large deposit vs. many small ones
4. **Don't exceed caps**: Diamond tier cap is 20K KAWAI per deposit

### Optimal Deposit Strategy

**For casual users:**
- Start with 50-100 USDT
- Gets you 2,500-5,000 KAWAI (first-time)
- Enough for hundreds of AI requests

**For active users:**
- Deposit 500-1,000 USDT
- Unlock Gold/Platinum tier
- 3-4% cashback on future deposits

**For power users:**
- Deposit 5,000+ USDT
- Diamond tier (5% cashback)
- Max rewards on every deposit

### Cost Estimates

**AI Chat:**
- Average: 0.01-0.05 USDT per request
- 100 USDT = 2,000-10,000 chats

**Image Generation:**
- Standard: 0.05-0.10 USDT
- HD: 0.10-0.20 USDT
- 100 USDT = 500-2,000 images

## 🆘 Troubleshooting

### "Insufficient MON for gas"

**Solution:**
1. Get MON from [faucet](https://testnet.monad.xyz/faucet)
2. Need ~0.001 MON per transaction
3. 0.1 MON = hundreds of transactions

### "USDT approval failed"

**Solutions:**
- Check you have enough USDT
- Make sure you're on Monad network
- Try increasing gas limit
- Wait and try again (network congestion)

### "Deposit not showing in balance"

**Solutions:**
1. Did you **sync the deposit**?
   - Click "Sync Deposit"
   - Paste transaction hash
   
2. Check transaction status:
   - Open [explorer](https://testnet.monad.xyz)
   - Paste transaction hash
   - Should show "Success"

3. Still not showing?
   - Wait 1-2 minutes
   - Refresh the app
   - Check you're on correct network

### "Wrong cashback amount"

**Check these:**
- Is it your first deposit? (5% applies)
- What's your tier? (1%-5%)
- Did you hit the tier cap? (5K-20K max)
- Account for 18 decimals (1 KAWAI = 1e18 wei)

## ✅ Checklist

Before depositing, verify:

- [ ] Wallet connected and on Monad network
- [ ] Have USDT on Monad
- [ ] Have MON for gas (~0.001 MON)
- [ ] Understand cashback tiers
- [ ] Know your current tier
- [ ] Ready to sync deposit after transaction

## 🚀 Next Steps

After depositing:

1. **[Claim Free Trial](free-trial.md)** - If you haven't yet
2. **[Start Using AI](ai-chat.md)** - Use your USDT balance
3. **[Check Rewards](../rewards/cashback.md)** - Track your cashback
4. **[Refer Friends](../rewards/referral.md)** - Earn more rewards

---

**Questions?** Check [FAQ](../faq/wallet.md) or [contact support](../support/contact.md).

