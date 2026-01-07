# Claiming Rewards

Learn how to claim your earned KAWAI tokens from all reward systems.

## 🎯 Overview

All rewards (mining, cashback, referral) use the same claiming process:

1. **Accumulation** - Rewards earned and tracked off-chain (instant)
2. **Settlement** - Weekly Merkle tree generation (Monday 00:00 UTC)
3. **Claiming** - On-chain claim with proof (anytime after settlement)

**Benefits of this system:**
- ✅ Gas-free accumulation
- ✅ Claim on your schedule
- ✅ Batch multiple periods
- ✅ Cryptographically secure

## 📅 Weekly Settlement Schedule

### Timeline

```
Sunday 23:59 → Period closes
Monday 00:00 → Settlement begins
Monday 00:30 → Merkle tree generated
Monday 01:00 → Proofs available
Monday onwards → You can claim!
```

### What Happens During Settlement

**Admin backend process:**

1. **Collect pending rewards**
   - Mining rewards from all contributors
   - Cashback from all deposits
   - Referral commissions from all affiliators

2. **Generate Merkle tree**
   - Create cryptographic proofs for each user
   - Calculate Merkle root (tree summary)

3. **Upload to blockchain**
   - Submit Merkle root to smart contracts
   - Transaction confirmed (~10 seconds)

4. **Store proofs**
   - Individual proofs saved off-chain
   - Accessible via API/dashboard

5. **Update status**
   - Period marked as "settled"
   - Rewards become "claimable"

## 🎁 Rewards Dashboard

### Accessing Dashboard

1. Open Kawai Desktop App
2. Click **"Wallet"** in sidebar
3. Click **"Rewards"** tab
4. See three sub-tabs:
   - **Mining Rewards**
   - **Cashback Rewards**
   - **Referral Rewards**

### Understanding the Display

#### Stats Section

Each reward type shows:

**Total Earned:**
- All-time cumulative rewards
- Across all periods
- Includes claimed + claimable + pending

**Pending:**
- Current week (not yet settled)
- Status: "Pending settlement"
- Will become claimable Monday

**Claimable:**
- Settled but not claimed
- Status: "Ready to claim"
- Click "Claim" button available

**Claimed:**
- Already withdrawn
- Status: "Claimed"
- Tokens in your wallet

#### History Table

Shows detailed breakdown:

| Column | Description |
|--------|-------------|
| **Period** | Week number (e.g., Period 5) |
| **Amount** | KAWAI tokens earned |
| **Date** | When reward was earned |
| **Status** | Pending/Claimable/Claimed |
| **Action** | Claim button (if claimable) |

## 💎 Claiming Process

### Single Period Claim

**Step-by-step:**

1. **Navigate to rewards**
   ```
   Wallet → Rewards → [Mining/Cashback/Referral]
   ```

2. **Find claimable period**
   - Look for green "Claim" button
   - Check amount and period number

3. **Click "Claim" button**
   - Review transaction details
   - Amount, period, contract

4. **Confirm in wallet**
   - MetaMask (or your wallet) opens
   - Shows transaction:
     - **To:** Reward distributor contract
     - **Function:** claimReward (or similar)
     - **Gas fee:** ~0.001-0.002 MON
   
5. **Wait for confirmation**
   - Usually 10-30 seconds
   - You'll see success notification

6. **Tokens received!**
   - Check your KAWAI balance
   - Should increase by claimed amount
   - Status updates to "Claimed"

### Batch Claiming (Multiple Periods)

**Save gas by claiming multiple periods at once:**

1. **Select periods**
   - Check boxes for multiple periods
   - Or click "Claim All" button

2. **Review batch claim**
   - See total amount
   - See all periods included
   - See estimated gas

3. **Confirm transaction**
   - Single transaction claims all
   - Gas fee: ~0.002-0.003 MON
   - More efficient than separate claims

4. **All tokens received**
   - Total amount added to balance
   - All periods marked "Claimed"

**Example savings:**
```
3 periods, separate claims:
  Gas: 3 × 0.001 = 0.003 MON

3 periods, batch claim:
  Gas: 0.002 MON

Savings: 33%!
```

## 🔍 Verifying Claims

### Check KAWAI Balance

**In the app:**
```
Wallet → Overview → KAWAI Balance
```

**In MetaMask:**
```
1. Add KAWAI token to MetaMask
2. Contract: [KAWAI address]
3. Symbol: KAWAI
4. Decimals: 18
4. Check balance
```

### Check Transaction on Explorer

1. **Copy transaction hash** from wallet
2. **Visit explorer:** [testnet.monad.xyz](https://testnet.monad.xyz)
3. **Paste transaction hash**
4. **Verify:**
   - Status: Success ✅
   - From: Your address
   - To: Distributor contract
   - Logs: TokensClaimed event

### Check Contract State

**Verify claim on-chain:**
```
1. Go to distributor contract on explorer
2. Read contract
3. Call isClaimed(period, yourAddress)
4. Should return: true
```

## 💰 Reward-Specific Claiming

### Mining Rewards

**Contract:** `MiningRewardDistributor`

**Claim includes:**
- Your share (85-90%)
- User share (5%)
- Developer share (5%)
- Affiliator share (0-5%)

**All recipients receive simultaneously.**

**Example:**
```
Job generated: 100 KAWAI total
Click "Claim" once
→ 85 KAWAI to you
→ 5 KAWAI to user
→ 5 KAWAI to developer  
→ 5 KAWAI to referrer (if applicable)
```

### Cashback Rewards

**Contract:** `DepositCashbackDistributor`

**Claim includes:**
- 100% to you (depositor)

**Simple single-recipient claim.**

**Example:**
```
Deposit: 1,000 USDT
Cashback: 20,000 KAWAI
Click "Claim"
→ 20,000 KAWAI to you
```

### Referral Rewards

**Two types:**

**1. One-time bonus:**
- Contract: `MerkleDistributor`
- 5 USDT + 100 KAWAI per referral
- Claimed immediately

**2. Mining commission:**
- Contract: `MiningRewardDistributor`
- 5% of referrals' mining
- Claimed weekly like mining rewards

## 🆘 Troubleshooting

### "No Claimable Rewards"

**Check:**
- Have you earned any rewards?
- Is it after Monday settlement?
- Refresh the page
- Check "Pending" section (wait for settlement)

**If you earned but not showing:**
- Wait until Monday for settlement
- Check specific reward tab (Mining/Cashback/Referral)
- Verify wallet address is correct

### "Transaction Failed"

**Common causes:**

**1. Insufficient gas (MON)**
```
Solution: Get MON from faucet
  → Visit testnet.monad.xyz/faucet
  → Request MON
  → Try claiming again
```

**2. Already claimed**
```
Solution: Check if status is "Claimed"
  → Each period can only be claimed once
  → Check transaction history
```

**3. Wrong network**
```
Solution: Switch to Monad Testnet
  → Open MetaMask
  → Select Monad Testnet
  → Refresh app
```

**4. Contract error**
```
Solution: Check explorer for details
  → Copy transaction hash
  → View on explorer
  → Read error message
```

### "Claim Button Disabled"

**Reasons:**

**Status: "Pending"**
- Wait for Monday settlement
- Cannot claim until settled

**Status: "Claimed"**
- Already withdrawn
- Check your balance/history

**Contract paused**
- Admin maintenance (rare)
- Check Discord for announcements

### "Wrong Amount Received"

**Check these:**

**1. Decimals confusion**
```
KAWAI uses 18 decimals
1 KAWAI = 1000000000000000000 wei
Dashboard shows human-readable
Blockchain shows raw wei
```

**2. Gas fees deducted**
```
Gas fee (MON) is separate
NOT deducted from KAWAI claim
Check MON balance for gas fee
```

**3. Split distribution (mining)**
```
For mining claims:
- You receive your 85-90%
- Others receive their shares
- Check your specific share
```

**4. Tier caps (cashback)**
```
Cashback has tier caps:
- Bronze: 5K max
- Diamond: 20K max
- Check tier at time of deposit
```

## 💡 Best Practices

### Optimal Claiming Strategy

**Option 1: Weekly claiming**
```
Pros:
✅ Regular cash flow
✅ Compounding opportunity
✅ Stay engaged

Cons:
❌ Pay gas every week
❌ Small amounts
```

**Option 2: Monthly claiming**
```
Pros:
✅ Lower gas fees (batch 4 weeks)
✅ Larger amounts
✅ More efficient

Cons:
❌ Delayed liquidity
❌ Need to remember
```

**Option 3: Threshold claiming**
```
Claim when reaching target (e.g., 1,000 KAWAI)

Pros:
✅ Meaningful amounts
✅ Flexible timing
✅ Minimized gas

Cons:
❌ Irregular schedule
❌ May wait long time
```

### Gas Optimization

**1. Batch claims**
- Claim multiple periods together
- Save 30-50% on gas

**2. Time your claims**
- Claim during low network usage
- Usually late night UTC
- Lower gas prices

**3. Group claim types**
- If you have multiple reward types
- Claim all at once
- Fewer wallet interactions

### Record Keeping

**Track your claims:**
```
Date | Period | Type | Amount | TX Hash | Gas Fee
-----|--------|------|---------|---------|--------
Jan 8| 5 | Mining | 250 | 0x123... | 0.001
Jan 8| 5 | Cashback | 5K | 0x456... | 0.001
Jan 15| 6 | All | 300 | 0x789... | 0.002
```

**Benefits:**
- Tax records
- Performance tracking
- Dispute resolution
- Historical analysis

## 📊 Claiming Statistics

### Average Claim Sizes

**By reward type:**

| Type | Average Claim | Frequency |
|------|---------------|-----------|
| Mining | 100-500 KAWAI | Weekly |
| Cashback | 2K-20K KAWAI | Per deposit |
| Referral (bonus) | 100 KAWAI | Per referral |
| Referral (commission) | 50-200 KAWAI | Weekly |

### Gas Costs

**Typical gas fees:**

| Action | Gas (MON) | USD Equivalent* |
|--------|-----------|-----------------|
| Single claim | ~0.001 | ~$0.01 |
| Batch 4 periods | ~0.002 | ~$0.02 |
| Batch 10 periods | ~0.003 | ~$0.03 |

*Assuming MON = $10 (testnet values)

### Optimal Frequency

**ROI calculation:**
```
Claim amount: 100 KAWAI
KAWAI price: $1
Value: $100

Gas cost: 0.001 MON = $0.01
ROI: $100 / $0.01 = 10,000x

Conclusion: Gas is negligible!
```

**Recommendation:** Claim whenever convenient. Gas is very cheap.

## ✅ Claiming Checklist

Before claiming:

- [ ] Rewards are in "Claimable" status
- [ ] Wallet connected and on Monad network
- [ ] Have MON for gas (~0.001-0.003 MON)
- [ ] Verified amount is correct
- [ ] Decided: single or batch claim?

After claiming:

- [ ] Transaction confirmed (green checkmark)
- [ ] KAWAI balance increased
- [ ] Status changed to "Claimed"
- [ ] Transaction saved/recorded
- [ ] Proof checked on explorer

## 🚀 Next Steps

After claiming your rewards:

1. **[Hold for Dividends](../tokenomics/hold-to-earn.md)** - Earn USDT in Phase 2
2. **[Trade on Marketplace](../trading/marketplace.md)** - Convert to USDT
3. **[Reinvest](../user-guide/deposit.md)** - Deposit more, earn more
4. **[Refer More](referral.md)** - Build passive income

---

**Questions?** Check [FAQ](../faq/rewards.md) or [contact support](../support/contact.md).

