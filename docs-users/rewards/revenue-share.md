# Revenue Share (Hold-to-Earn)

**Earn USDT just by holding KAWAI tokens!**

## 🎯 Overview

Revenue Share, also known as "Hold-to-Earn", is a unique reward system where KAWAI token holders receive 100% of the platform's net profit in USDT, distributed proportionally based on their token holdings.

!!! success "Key Benefits"
    - 💰 **Passive Income** - Earn USDT weekly without any action
    - 📊 **Proportional Distribution** - Fair share based on holdings
    - 🔓 **No Lock Required** - Just hold tokens in your wallet
    - 💵 **Stable Rewards** - Earn in USDT, not volatile tokens
    - 🌱 **Sustainable** - Tied to real platform revenue

## 💡 How It Works

### The Formula

Your weekly USDT reward is calculated using this simple formula:

```
Your USDT = (Your KAWAI Balance / Total KAWAI Supply) × Weekly Net Profit
```

### Example Calculation

**Scenario:**
- Total KAWAI Supply: 100,000,000 KAWAI
- Your KAWAI Balance: 1,000,000 KAWAI (1% of supply)
- Weekly Net Profit: 10,000 USDT

**Your Share:**
```
Your USDT = (1,000,000 / 100,000,000) × 10,000 USDT
          = 0.01 × 10,000 USDT
          = 100 USDT per week
```

**Annual Projection:**
```
Weekly: 100 USDT
Monthly: ~400 USDT
Yearly: ~5,200 USDT

Annual Yield: 5,200 USDT / Investment = ROI%
```

## 📈 Revenue Sources

Platform profit comes from real business activities:

### Phase 1: Mining Era (Current)
- Users pay USDT for AI services
- Platform accumulates treasury
- KAWAI tokens minted as rewards
- Minimal USDT distribution (accumulation phase)

### Phase 2: USDT Era (After 1B Supply)
- No more token minting
- All revenue distributed to holders
- Contributors paid in USDT (not KAWAI)
- Full revenue sharing activated

## 💰 Profit Calculation

### Revenue Breakdown

**Total Revenue (100%):**
```
User Payments (USDT) = Total Revenue
```

**Operating Costs (~90% of gross revenue):**
```
- Contributor Costs: 85-90% of gross revenue (GPU providers)
- Developer Costs: 5% of gross revenue (platform maintenance)
- User Cashback: 5% of gross revenue (use-to-earn)
- Affiliator Commission: 0-5% of gross revenue (referral rewards)

Note: Percentages are of gross revenue. Total operating costs target ~90%.
```

**Net Profit (~10%):**
```
Net Profit = Total Revenue - Operating Costs
           = 100% - 90%
           = 10% of revenue
```

**Distribution:**
```
100% of Net Profit → KAWAI Holders (proportionally)
```

### Example Revenue Flow

**Weekly Platform Activity:**
```
Total AI Requests: 111,000,000 tokens processed
Cost per 1M tokens: 1 USDT
Total Revenue: 111 USDT

Operating Costs (90% of revenue):
- Contributors: 85 USDT (85% of 100)
- Developers: 5 USDT (5% of 100)
- Users: 5 USDT (5% of 100)
- Affiliators: 5 USDT (5% of 100)
Total Costs: 100 USDT

Net Profit: 11 USDT (~10% of revenue → distributed to KAWAI holders)
```

!!! note "Profit Margins"
    Actual profit margins depend on:
    - Platform efficiency
    - Contributor costs
    - Service pricing
    - Network utilization
    
    Target margin: 10-20% net profit

## 📊 Estimated Yields

### Conservative Estimates

Based on network growth projections:

| Holding | Share % | Weekly USDT | Monthly USDT | Annual USDT | Est. APY* |
|---------|---------|-------------|--------------|-------------|-----------|
| 10,000 KAWAI | 0.001% | $0.10 | $0.40 | $5.20 | 5% |
| 100,000 KAWAI | 0.01% | $1.00 | $4.00 | $52.00 | 5% |
| 1,000,000 KAWAI | 0.1% | $10.00 | $40.00 | $520.00 | 5% |
| 10,000,000 KAWAI | 1% | $100.00 | $400.00 | $5,200.00 | 5% |

*Assuming 10,000 USDT weekly profit and stable token price

### Optimistic Estimates

With high network adoption:

| Holding | Share % | Weekly USDT | Monthly USDT | Annual USDT | Est. APY* |
|---------|---------|-------------|--------------|-------------|-----------|
| 10,000 KAWAI | 0.001% | $0.50 | $2.00 | $26.00 | 25% |
| 100,000 KAWAI | 0.01% | $5.00 | $20.00 | $260.00 | 25% |
| 1,000,000 KAWAI | 0.1% | $50.00 | $200.00 | $2,600.00 | 25% |
| 10,000,000 KAWAI | 1% | $500.00 | $2,000.00 | $26,000.00 | 25% |

*Assuming 50,000 USDT weekly profit and stable token price

## 🔄 Distribution Process

### 1. Revenue Collection (Continuous)

```
Users pay USDT → PaymentVault Contract
                → Platform Treasury
```

### 2. Profit Calculation (Weekly)

Every Monday at 00:00 UTC:
```
1. Calculate total revenue
2. Deduct operating costs
3. Determine net profit
4. Snapshot KAWAI holder balances
```

### 3. Merkle Tree Generation

```
For each holder:
  - Calculate share percentage
  - Calculate USDT amount
  - Generate Merkle proof
  - Store proof off-chain
```

### 4. On-Chain Settlement

```
Admin uploads Merkle root → MerkleDistributor Contract
                          → Enables claiming
```

### 5. Claiming (User Action)

```
User clicks "Claim" → Provides Merkle proof
                    → Contract verifies
                    → Transfers USDT to user
```

## 🎁 How to Claim

### Step 1: Navigate to Revenue Share

1. Open Kawai DeAI app
2. Go to **Wallet** tab
3. Click **Rewards**
4. Select **Revenue Share** tab

### Step 2: Check Claimable Amount

You'll see:
- **Total Earned**: Lifetime USDT earned
- **Claimable Now**: Ready to claim
- **Est. Weekly**: Projected weekly earnings

### Step 3: Review Your Share

Check your position:
- Your KAWAI Balance
- Your Share Percentage
- Estimated earnings

### Step 4: Claim Rewards

1. Click **Claim** button on any claimable period
2. Review transaction details:
   - Period number
   - USDT amount
   - Your share percentage
   - Gas fee estimate
3. Click **Confirm & Claim**
4. Wait for transaction confirmation
5. USDT transferred to your wallet!

### Step 5: Track Transaction

- View transaction on Monad Explorer
- Check wallet balance
- See updated stats in dashboard

## 💡 Maximizing Your Revenue Share

### Strategy 1: Accumulate Early

**Why it matters:**
- Lower token price in Phase 1
- Higher percentage of supply
- Compounding effect over time

**Example:**
```
Buy 1M KAWAI at $0.01 = $10,000 investment
If supply is 100M, you own 1%
At 1B supply, you still own 1% (if you hold)
```

### Strategy 2: Hold Long-Term

**Benefits:**
- Weekly USDT income
- No need to sell tokens
- Benefit from network growth
- Potential token appreciation

**Comparison:**
```
Trader: Buy low, sell high (risky, timing-dependent)
Holder: Earn dividends forever (passive, sustainable)
```

### Strategy 3: Reinvest Dividends

**Compound your earnings:**
```
Week 1: Earn 100 USDT → Buy 10,000 KAWAI
Week 2: Earn 101 USDT (slightly more)
Week 3: Earn 102 USDT (compounding effect)
...
Year 1: Significantly higher holdings
```

### Strategy 4: Diversify Holdings

**Balance your portfolio:**
- 50% KAWAI (for dividends)
- 30% USDT (for stability)
- 20% Other assets (for diversification)

## 🔍 Understanding Dilution

### What is Dilution?

As more KAWAI tokens are minted, your percentage of total supply decreases:

```
Initial: 1M KAWAI / 100M supply = 1%
Later: 1M KAWAI / 1B supply = 0.1%

Your share percentage decreased by 10x!
```

### Does Dilution Matter?

**Yes and No:**

**The Bad News:**
- Your % share decreases
- Each token represents smaller ownership

**The Good News:**
- Total profit pool grows with network
- More users = more revenue
- Absolute USDT earnings can stay same or increase

**Example:**
```
Phase 1 (100M supply):
- Your holdings: 1M KAWAI (1%)
- Weekly profit: 1,000 USDT
- Your share: 10 USDT

Phase 2 (1B supply):
- Your holdings: 1M KAWAI (0.1%)
- Weekly profit: 10,000 USDT (10x growth)
- Your share: 10 USDT (same!)
```

### Mitigating Dilution

**Strategies:**
1. **Accumulate Early** - Buy when supply is low
2. **Earn Rewards** - Get free tokens from other systems
3. **Compound** - Reinvest dividends to maintain %
4. **Long-term Hold** - Benefit from network growth

## 📋 Phase Comparison

### Phase 1: Mining Era (Current)

**Token Economics:**
- KAWAI tokens minted as rewards
- Gradual emission (halving schedule)
- Max supply: 1 Billion KAWAI

**Revenue Model:**
- Users pay USDT for services
- Platform accumulates treasury
- Minimal USDT distribution

**Holder Benefits:**
- Token appreciation potential
- Future dividend rights
- Governance rights (planned)

### Phase 2: USDT Era (Future)

**Token Economics:**
- No more minting (1B reached)
- Fixed supply
- Deflationary (if burning implemented)

**Revenue Model:**
- Contributors paid in USDT
- 100% profit to holders
- Weekly USDT distributions

**Holder Benefits:**
- Passive USDT income
- Stable, predictable yields
- Real yield from real revenue

## ❓ Common Questions

### Q: When does revenue sharing start?

**A:** Full revenue sharing begins in Phase 2, after the 1 Billion KAWAI supply is reached. Currently in Phase 1 (Mining Era), minimal distributions may occur as platform accumulates treasury.

### Q: Do I need to lock my tokens?

**A:** No! Just hold KAWAI in your wallet. No staking or locking required. You can trade anytime.

### Q: How often are distributions?

**A:** Weekly settlements every Monday at 00:00 UTC. You can claim anytime after settlement.

### Q: What if I don't claim?

**A:** Your rewards never expire! Claim whenever convenient. However, consider gas fees when claiming small amounts.

### Q: Can I sell tokens after claiming?

**A:** Yes! After claiming, USDT is yours to use however you want. Your KAWAI tokens remain in your wallet.

### Q: Is this sustainable?

**A:** Yes! Unlike ponzi schemes, revenue comes from real AI service payments. As long as people use the platform, there's revenue to distribute.

### Q: What's the minimum holding to earn?

**A:** No minimum! Even 1 KAWAI earns proportional rewards. However, very small holdings may not justify gas fees for claiming.

### Q: How is this different from staking?

**A:** 
- **Staking**: Lock tokens, earn new tokens (inflationary)
- **Revenue Share**: Hold tokens, earn USDT from profit (sustainable)

### Q: What affects my earnings?

**A:**
1. Your KAWAI holdings (more = higher %)
2. Total supply (lower = higher %)
3. Platform revenue (higher = more USDT)
4. Network growth (more users = more profit)

## 🎯 Getting Started

### For New Users

1. **[Get KAWAI Tokens](../trading/buying-kawai.md)**
   - Buy on P2P marketplace
   - Earn through rewards
   - Receive from referrals

2. **[Set Up Wallet](../user-guide/wallet-setup.md)**
   - Create secure wallet
   - Store KAWAI safely
   - Enable notifications

3. **[Monitor Dashboard](../user-guide/dashboard.md)**
   - Check your share %
   - Track earnings
   - Claim rewards

4. **[Claim Rewards](claiming.md)**
   - Weekly after settlement
   - Batch claim multiple periods
   - Minimize gas fees

### For Existing Holders

1. **Check Your Share**
   - Go to Revenue Share tab
   - View current holdings
   - See share percentage

2. **Estimate Earnings**
   - Use calculator (coming soon)
   - Project monthly/yearly
   - Plan accordingly

3. **Optimize Strategy**
   - Accumulate more tokens
   - Compound dividends
   - Hold long-term

## 🚀 Future Enhancements

### Planned Features

- [ ] **Revenue Calculator** - Estimate earnings based on holdings
- [ ] **Auto-Claim** - Automatic claiming (optional)
- [ ] **Compounding** - Auto-reinvest dividends
- [ ] **Analytics Dashboard** - Detailed revenue metrics
- [ ] **Mobile Notifications** - Alert when claimable
- [ ] **Historical Charts** - Track earnings over time

### Governance (Future)

KAWAI holders may vote on:
- Revenue distribution ratios
- Platform fee structures
- Feature priorities
- Treasury management

## 📚 Related Documentation

- **[Rewards Overview](overview.md)** - All reward systems explained
- **[Claiming Guide](claiming.md)** - How to claim any reward
- **[KAWAI Token](../tokenomics/kawai-token.md)** - Token economics
- **[Trading Guide](../trading/p2p-marketplace.md)** - Buy/sell KAWAI
- **[Wallet Setup](../user-guide/wallet-setup.md)** - Secure your tokens

## 🆘 Need Help?

- **[FAQ](../faq/rewards.md)** - Common questions
- **[Support](../support/contact.md)** - Get assistance
- **[Discord](https://discord.gg/SNf3ZEa8Eq)** - Community help

---

**Start earning passive USDT income today!** [Get KAWAI tokens →](../trading/buying-kawai.md)
