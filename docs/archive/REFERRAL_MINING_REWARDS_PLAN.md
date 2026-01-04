# Referral-Based Mining Reward Distribution

**Status:** 🚧 Implementation Plan  
**Created:** 2026-01-04  
**Branch:** `feature/referral-mining-rewards`

---

## 📋 Overview

This document outlines the implementation plan for a new mining reward distribution model that incentivizes the referral program by giving affiliators (referrers) a share of their referees' mining rewards.

### Current Model (Existing)

All mining rewards are split:
- **90% Contributor** (miner)
- **10% Developer** (protocol treasury)

### New Model (Proposed)

Mining rewards are split based on whether the user was referred.

#### **For Referral Users:**
- **85% Contributor** (miner) - reduced to fund affiliator
- **5% Developer** (protocol treasury)
- **5% User** (requester/consumer) ⭐ NEW
- **5% Affiliator** (referrer) ⭐ NEW

**Total: 100%**

#### **For Non-Referral Users:**
- **90% Contributor** (miner)
- **5% Developer** (protocol treasury)
- **5% User** (requester/consumer) ⭐ NEW

**Total: 100%**

---

## 🎯 Goals & Benefits

### 1. **Stronger Referral Incentive**
- Affiliators earn **lifetime passive income** from their referees' mining activity
- More attractive than one-time referral bonuses
- Creates long-term engagement for affiliators

### 2. **User Incentive (NEW)**
- Users get **5% cashback** on every request they make
- Encourages platform usage and loyalty
- Creates a "use-to-earn" model alongside mining
- Users benefit from both sides: cheap AI + token rewards

### 3. **Sustainable Growth Model**
- Referred users tend to be more engaged (personal invitation)
- Affiliators motivated to actively promote the platform
- Users motivated to use the platform more (5% cashback)
- Network effect: more users = more liquidity = more value

### 4. **Fair Distribution**
- Contributors get 85-90% (still majority of rewards)
- Developers get 5% (consistent across all scenarios)
- Users get 5% cashback (incentivizes usage)
- Affiliators get 5% commission from referrals (contributor sacrifices 5% for growth)
- Trade-off: Contributors earn slightly less (85% vs 90%) but benefit from larger network

### 5. **Simple & Transparent**
- Only 2 scenarios, easy for users to understand
- No complex tier/level systems
- Clear reward expectations

---

## 🏗️ Architecture Changes

### 1. Smart Contract Updates

#### **A. New Contract: `MiningRewardDistributor.sol`** ⭐

**Decision:** Create a new dedicated contract instead of reusing the existing `MerkleDistributor.sol`.

**Why?**
- ✅ Full on-chain transparency of reward breakdown
- ✅ Better analytics and audit trail via events
- ✅ Single claim transaction distributes to all parties (gas efficient)
- ✅ Easier to maintain and extend in the future
- ✅ Clear separation of concerns
- ✅ Clear naming: focuses on mining rewards (referral mechanism is implementation detail)

**Contract Design:**

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title MiningRewardDistributor
 * @notice Distributes mining rewards with referral-based splits
 * @dev Uses Merkle proofs for gas-efficient batch distributions
 */
contract MiningRewardDistributor is AccessControl, ReentrancyGuard {
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    
    IERC20 public immutable kawaiToken;
    
    // Merkle root for each period
    mapping(uint256 => bytes32) public merkleRoots;
    
    // Track claimed rewards: period => user => claimed
    mapping(uint256 => mapping(address => bool)) public hasClaimed;
    
    // Reward breakdown for transparency
    struct RewardBreakdown {
        uint256 contributorAmount;  // 85% (referral) or 90% (non-referral)
        uint256 developerAmount;    // 5%
        uint256 userAmount;         // 5%
        uint256 affiliatorAmount;   // 5% (referral only) or 0%
        address affiliator;         // referrer address (zero if non-referral)
        address user;               // requester address
    }
    
    event MerkleRootSet(uint256 indexed period, bytes32 merkleRoot);
    event RewardClaimed(
        uint256 indexed period,
        address indexed contributor,
        address indexed user,
        uint256 contributorAmount,
        uint256 developerAmount,
        uint256 userAmount,
        uint256 affiliatorAmount,
        address affiliator
    );
    
    constructor(address _kawaiToken) Ownable(msg.sender) {
        require(_kawaiToken != address(0), "Invalid KAWAI address");
        kawaiToken = IERC20(_kawaiToken);
        currentPeriod = 1;
    }
    
    // Note: Developer addresses are NOT fixed in constructor
    // Each Merkle leaf includes the developer address (from GetRandomTreasuryAddress)
    // This allows flexible distribution to various treasury addresses per job
    
    /**
     * @notice Set Merkle root for a period
     * @param period The settlement period ID
     * @param merkleRoot The Merkle root hash
     */
    function setMerkleRoot(uint256 period, bytes32 merkleRoot) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        require(merkleRoots[period] == bytes32(0), "Root already set");
        merkleRoots[period] = merkleRoot;
        emit MerkleRootSet(period, merkleRoot);
    }
    
    /**
     * @notice Claim mining rewards with referral split
     * @param period The settlement period
     * @param contributorAmount Amount for contributor (85% referral, 90% non-referral)
     * @param developerAmount Amount for developer (5%)
     * @param userAmount Amount for user/requester (5%)
     * @param affiliatorAmount Amount for affiliator (5% referral only, 0% non-referral)
     * @param user User/requester address
     * @param affiliator Referrer address (zero address if non-referral)
     * @param merkleProof Merkle proof for verification
     */
    function claimReward(
        uint256 period,
        uint256 contributorAmount,
        uint256 developerAmount,
        uint256 userAmount,
        uint256 affiliatorAmount,
        address user,
        address affiliator,
        bytes32[] calldata merkleProof
    ) external nonReentrant {
        require(!hasClaimed[period][msg.sender], "Already claimed");
        require(merkleRoots[period] != bytes32(0), "Invalid period");
        
        // Verify Merkle proof
        bytes32 leaf = keccak256(abi.encodePacked(
            msg.sender,
            contributorAmount,
            developerAmount,
            userAmount,
            affiliatorAmount,
            user,
            affiliator
        ));
        require(
            MerkleProof.verify(merkleProof, merkleRoots[period], leaf),
            "Invalid proof"
        );
        
        // Mark as claimed
        hasClaimed[period][msg.sender] = true;
        
        // Transfer rewards
        require(
            kawaiToken.transfer(msg.sender, contributorAmount),
            "Contributor transfer failed"
        );
        
        if (developerAmount > 0) {
            address developer = getRoleAdmin(ADMIN_ROLE); // or treasury address
            require(
                kawaiToken.transfer(developer, developerAmount),
                "Developer transfer failed"
            );
        }
        
        if (userAmount > 0 && user != address(0)) {
            require(
                kawaiToken.transfer(user, userAmount),
                "User transfer failed"
            );
        }
        
        if (affiliatorAmount > 0 && affiliator != address(0)) {
            require(
                kawaiToken.transfer(affiliator, affiliatorAmount),
                "Affiliator transfer failed"
            );
        }
        
        emit RewardClaimed(
            period,
            msg.sender,
            user,
            contributorAmount,
            developerAmount,
            userAmount,
            affiliatorAmount,
            affiliator
        );
    }
    
    /**
     * @notice Batch claim for multiple periods
     */
    function claimMultiplePeriods(
        uint256[] calldata periods,
        uint256[] calldata contributorAmounts,
        uint256[] calldata developerAmounts,
        uint256[] calldata userAmounts,
        uint256[] calldata affiliatorAmounts,
        address[] calldata users,
        address[] calldata affiliators,
        bytes32[][] calldata merkleProofs
    ) external {
        require(
            periods.length == contributorAmounts.length &&
            periods.length == developerAmounts.length &&
            periods.length == userAmounts.length &&
            periods.length == affiliatorAmounts.length &&
            periods.length == users.length &&
            periods.length == affiliators.length &&
            periods.length == merkleProofs.length,
            "Array length mismatch"
        );
        
        for (uint256 i = 0; i < periods.length; i++) {
            claimReward(
                periods[i],
                contributorAmounts[i],
                developerAmounts[i],
                userAmounts[i],
                affiliatorAmounts[i],
                users[i],
                affiliators[i],
                merkleProofs[i]
            );
        }
    }
}
```

**Key Features:**
- ✅ Merkle-based distribution (gas efficient)
- ✅ Supports both referral and non-referral users
- ✅ 5% user cashback (use-to-earn incentive)
- ✅ Transparent reward breakdown in events
- ✅ Single claim distributes to all parties (contributor, developer, user, affiliator)
- ✅ Batch claiming for multiple periods
- ✅ Reentrancy protection
- ✅ Full on-chain audit trail

**Claim Flow:**
1. **Contributor initiates claim** (they have the most incentive - 85-90% reward)
2. Contract verifies Merkle proof with all parameters
3. Contract distributes rewards to all parties in one transaction:
   - Contributor gets 85-90%
   - Developer gets 5%
   - User gets 5% (cashback)
   - Affiliator gets 5% (if referral user)
4. Event emitted with full breakdown for analytics

---

### 2. Backend Service Updates

#### **A. Update `pkg/store/contributor.go`**

Modify `RecordJobReward()` to track referral information:

```go
// RewardRecord now includes referral and user info
type RewardRecord struct {
    ContributorAddress string  `json:"contributor_address"`
    UserAddress       string  `json:"user_address"`        // NEW: requester address
    Amount            float64  `json:"amount"`
    TokenUsage        int64    `json:"token_usage"`
    Timestamp         int64    `json:"timestamp"`
    Mode              string   `json:"mode"` // "mining" or "usdt"
    
    // NEW: Referral tracking
    HasReferrer       bool     `json:"has_referrer"`
    ReferrerAddress   string   `json:"referrer_address,omitempty"`
    
    // NEW: Reward breakdown
    ContributorReward float64  `json:"contributor_reward"` // 85% (referral) or 90% (non-referral)
    DeveloperReward   float64  `json:"developer_reward"`   // 5%
    UserReward        float64  `json:"user_reward"`        // 5%
    AffiliatorReward  float64  `json:"affiliator_reward"`  // 5% (referral) or 0% (non-referral)
}

// Update RecordJobReward to accept user and referrer info
func (s *ContributorStore) RecordJobReward(
    ctx context.Context,
    contributorAddress string,
    userAddress string,        // NEW parameter
    tokenUsage int64,
    referrerAddress string,    // NEW parameter
) error {
    // ... existing supply check logic ...
    
    // Determine split based on referrer
    hasReferrer := referrerAddress != "" && referrerAddress != "0x0000000000000000000000000000000000000000"
    
    var contributorReward, developerReward, userReward, affiliatorReward float64
    
    if mode == config.ModeMining {
        baseReward := calculateMiningReward(tokenUsage, currentSupply)
        
        developerReward = baseReward * 0.05    // Always 5%
        userReward = baseReward * 0.05         // Always 5%
        
        if hasReferrer {
            // Referral user: 85/5/5/5 split
            contributorReward = baseReward * 0.85
            affiliatorReward = baseReward * 0.05
        } else {
            // Non-referral user: 90/5/5 split
            contributorReward = baseReward * 0.90
            affiliatorReward = 0
        }
        
        // Total: Always 100%
    } else {
        // Phase 2 (USDT mode) - similar logic
        // ...
    }
    
    // Save reward record with breakdown
    record := RewardRecord{
        ContributorAddress: contributorAddress,
        UserAddress:       userAddress,
        TokenUsage:        tokenUsage,
        Timestamp:         time.Now().Unix(),
        Mode:              string(mode),
        HasReferrer:       hasReferrer,
        ReferrerAddress:   referrerAddress,
        ContributorReward: contributorReward,
        DeveloperReward:   developerReward,
        UserReward:        userReward,
        AffiliatorReward:  affiliatorReward,
    }
    
    // Update contributor balance
    if err := s.updateBalance(ctx, contributorAddress, contributorReward); err != nil {
        return err
    }
    
    // Update developer balance
    if err := s.updateBalance(ctx, developerAddress, developerReward); err != nil {
        return err
    }
    
    // Update user balance (cashback)
    if err := s.updateBalance(ctx, userAddress, userReward); err != nil {
        return err
    }
    
    // Update affiliator balance if applicable
    if hasReferrer && affiliatorReward > 0 {
        if err := s.updateBalance(ctx, referrerAddress, affiliatorReward); err != nil {
            return err
        }
    }
    
    return nil
}
```

#### **B. Update API Gateway**

Modify the job completion handler to pass referrer information:

```go
// In internal/api/gateway.go or similar

func (g *Gateway) handleJobCompletion(ctx context.Context, job *Job) error {
    // Get user info
    user, err := g.userStore.GetUser(ctx, job.UserAddress)
    if err != nil {
        return err
    }
    
    // Get referrer address (if any)
    referrerAddress := ""
    if user.ReferredBy != "" {
        referrer, err := g.userStore.GetUserByReferralCode(ctx, user.ReferredBy)
        if err == nil && referrer != nil {
            referrerAddress = referrer.WalletAddress
        }
    }
    
    // Record reward with user and referrer info
    return g.contributorStore.RecordJobReward(
        ctx,
        job.ContributorAddress,
        job.UserAddress,      // Pass user address
        job.TokenUsage,
        referrerAddress,      // Pass referrer address
    )
}
```

---

### 3. Settlement & Merkle Generation Updates

#### **A. Update Weekly Settlement Script**

Modify `cmd/admin/snapshot.go` to generate Merkle tree with referral splits:

```go
type MerkleLeaf struct {
    ContributorAddress string
    ContributorAmount  *big.Int
    DeveloperAmount    *big.Int
    UserAmount         *big.Int
    AffiliatorAmount   *big.Int
    UserAddress        string
    AffiliatorAddress  string
}

func generateMerkleTree(rewards []RewardRecord) (*merkle.Tree, error) {
    leaves := make([][]byte, 0)
    
    // Group rewards by contributor
    contributorRewards := groupRewardsByContributor(rewards)
    
    for contributorAddress, rewardList := range contributorRewards {
        var totalContributor, totalDeveloper, totalUser, totalAffiliator float64
        var userAddress, affiliatorAddress string
        
        for _, reward := range rewardList {
            totalContributor += reward.ContributorReward
            totalDeveloper += reward.DeveloperReward
            totalUser += reward.UserReward
            totalAffiliator += reward.AffiliatorReward
            
            userAddress = reward.UserAddress
            if reward.HasReferrer {
                affiliatorAddress = reward.ReferrerAddress
            }
        }
        
        // Create leaf
        leaf := crypto.Keccak256(
            common.HexToAddress(contributorAddress).Bytes(),
            toWei(totalContributor).Bytes(),
            toWei(totalDeveloper).Bytes(),
            toWei(totalUser).Bytes(),
            toWei(totalAffiliator).Bytes(),
            common.HexToAddress(userAddress).Bytes(),
            common.HexToAddress(affiliatorAddress).Bytes(),
        )
        
        leaves = append(leaves, leaf)
    }
    
    return merkle.NewTree(leaves)
}
```

---

### 4. Frontend UI Updates

#### **A. Affiliator Commission Dashboard**

Create a new section in the Rewards Dashboard to show affiliator earnings:

**Location:** `frontend/src/app/wallet/components/rewards/AffiliatorRewardsSection.tsx`

```typescript
import { Card, Statistic, Table, Button, Tooltip } from 'antd';
import { Users, TrendingUp, Award } from 'lucide-react';
import { useState, useEffect } from 'react';
import { ReferralService } from '@@/github.com/kawai-network/veridium/internal/services';

export const AffiliatorRewardsSection = ({ currentNetwork, theme, styles, onRefresh }) => {
  const [stats, setStats] = useState(null);
  const [miningCommissions, setMiningCommissions] = useState([]);
  const [loading, setLoading] = useState(false);

  const loadData = async (force = false) => {
    setLoading(true);
    try {
      // Get referral stats
      const statsData = await ReferralService.GetReferralStats();
      setStats(statsData);
      
      // Get mining commissions (NEW API)
      const commissions = await ReferralService.GetMiningCommissions();
      setMiningCommissions(commissions);
    } catch (error) {
      console.error('Failed to load affiliator data:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
    onRefresh?.(() => loadData(true));
  }, []);

  return (
    <div>
      {/* Stats Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16, marginBottom: 24 }}>
        <Card>
          <Statistic
            title="Total Referrals"
            value={stats?.totalReferrals || 0}
            prefix={<Users size={20} />}
          />
        </Card>
        
        <Card>
          <Statistic
            title="Mining Commission (Lifetime)"
            value={(stats?.totalMiningCommission || 0).toFixed(2)}
            suffix="KAWAI"
            prefix={<Award size={20} />}
            valueStyle={{ color: '#3f8600' }}
          />
        </Card>
        
        <Card>
          <Statistic
            title="Active Miners"
            value={stats?.activeMiners || 0}
            prefix={<TrendingUp size={20} />}
          />
        </Card>
      </div>

      {/* Commission History */}
      <Card title="Mining Commission History">
        <Table
          dataSource={miningCommissions}
          loading={loading}
          columns={[
            {
              title: 'Period',
              dataIndex: 'period',
              key: 'period',
            },
            {
              title: 'Referee',
              dataIndex: 'refereeAddress',
              key: 'refereeAddress',
              render: (addr) => `${addr.slice(0, 6)}...${addr.slice(-4)}`,
            },
            {
              title: 'Their Mining',
              dataIndex: 'refereeMiningAmount',
              key: 'refereeMiningAmount',
              render: (amount) => `${amount.toFixed(2)} KAWAI`,
            },
            {
              title: 'Your Commission (5%)',
              dataIndex: 'commissionAmount',
              key: 'commissionAmount',
              render: (amount) => (
                <span style={{ color: '#3f8600', fontWeight: 600 }}>
                  +{amount.toFixed(2)} KAWAI
                </span>
              ),
            },
            {
              title: 'Status',
              dataIndex: 'status',
              key: 'status',
              render: (status) => (
                <span style={{ 
                  color: status === 'claimed' ? '#3f8600' : '#faad14',
                  fontWeight: 500 
                }}>
                  {status === 'claimed' ? '✓ Claimed' : '⏳ Pending'}
                </span>
              ),
            },
          ]}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {/* Info Banner */}
      <Card style={{ marginTop: 16, background: theme.colorInfoBg }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Award size={24} style={{ color: theme.colorInfo }} />
          <div>
            <div style={{ fontWeight: 600, marginBottom: 4 }}>
              Lifetime Passive Income
            </div>
            <div style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              You earn 5% commission from all mining rewards your referrals generate, forever.
              The more they mine, the more you earn!
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
};
```

#### **B. Update Rewards Dashboard Tabs**

Add "Affiliator Commission" tab to `RewardsContent.tsx`:

```typescript
// In frontend/src/app/wallet/RewardsContent.tsx

const tabs = [
  {
    key: 'mining',
    label: <span><Coins size={16} /> Mining Rewards</span>,
    children: <MiningRewardsSection {...props} />,
  },
  {
    key: 'cashback',
    label: <span><Award size={16} /> Deposit Cashback</span>,
    children: <CashbackRewardsSection {...props} />,
  },
  {
    key: 'referral',
    label: <span><Users size={16} /> Referral Rewards</span>,
    children: <ReferralRewardsSection {...props} />,
  },
  // NEW TAB
  {
    key: 'affiliator',
    label: <span><TrendingUp size={16} /> Mining Commission</span>,
    children: <AffiliatorRewardsSection {...props} />,
  },
];
```

---

## ⏰ Weekly Settlement & Claim Flow

### **Off-Chain Accumulation (Monday - Saturday)**

Every time a job completes:
```go
// Backend records reward breakdown
RecordJobReward(
    contributorAddress,
    userAddress,
    tokenUsage,
    referrerAddress,
)

// Updates balances in KV store:
// - Contributor: +850 or +900 KAWAI
// - Developer: +50 KAWAI
// - User: +50 KAWAI
// - Affiliator: +50 KAWAI (if referral)
```

**Benefits:**
- ⚡ Instant recording (no blockchain latency)
- 💰 Zero gas fees
- 📈 Scales to millions of jobs

### **On-Chain Settlement (Sunday)**

Admin runs weekly settlement script:

```bash
# 1. Generate Merkle Tree
go run cmd/admin/snapshot.go

# Merkle leaves structure:
# Each leaf = (contributor, contributorAmount, developer, developerAmount, user, userAmount, affiliator, affiliatorAmount)
```

**What happens:**
1. **Snapshot:** Read all accumulated rewards from KV store
2. **Group:** Group by contributor (one leaf per contributor per week)
3. **Generate:** Create Merkle tree with all reward data
4. **Upload:** Send Merkle root to `MiningRewardDistributor.sol`
5. **Store Proofs:** Save individual proofs to KV store for claiming

**Cost:**
- 1 transaction per week (~$0.01 on Monad)
- Handles unlimited number of jobs/users

### **Claim Process (Anytime)**

#### **Option A: Contributor-Initiated Claim** (Recommended)

Contributor claims for everyone:

```typescript
// Frontend: Rewards Dashboard
const claimReward = async (period: number) => {
  const proof = await getProof(period, contributorAddress);
  
  await MiningRewardDistributor.claimReward(
    period,
    proof.contributorAmount,  // 850 KAWAI
    proof.developerAmount,    // 50 KAWAI
    proof.userAmount,         // 50 KAWAI
    proof.affiliatorAmount,   // 50 KAWAI
    proof.userAddress,
    proof.affiliatorAddress,
    proof.merkleProof
  );
  
  // ✅ All parties receive their rewards in ONE transaction
};
```

**Benefits:**
- ✅ One transaction distributes to all parties
- ✅ Contributor pays gas (they get 85-90%, worth it)
- ✅ User and Affiliator get rewards automatically
- ✅ Simpler UX

#### **Option B: Individual Claims** (Alternative)

Each party claims their own:

```typescript
// User claims their cashback
await claimUserReward(period, amount, proof);

// Affiliator claims their commission
await claimAffiliatorReward(period, amount, proof);

// Contributor claims their mining reward
await claimContributorReward(period, amount, proof);
```

**Trade-offs:**
- ❌ 3-4 separate transactions = more gas
- ✅ Each party controls their own claim timing
- ❌ More complex UX

**Decision: Use Option A** (Contributor-initiated claim)

### **Batch Claiming**

Claim multiple weeks at once:

```typescript
const claimMultipleWeeks = async (periods: number[]) => {
  const proofs = await Promise.all(
    periods.map(p => getProof(p, contributorAddress))
  );
  
  await MiningRewardDistributor.claimMultiplePeriods(
    periods,
    proofs.map(p => p.contributorAmount),
    proofs.map(p => p.developerAmount),
    proofs.map(p => p.userAmount),
    proofs.map(p => p.affiliatorAmount),
    proofs.map(p => p.userAddress),
    proofs.map(p => p.affiliatorAddress),
    proofs.map(p => p.merkleProof)
  );
};
```

**Example:**
```
Claim Week 1-4 together:
├─ Total Contributor: 3,400 KAWAI
├─ Total Developer: 200 KAWAI
├─ Total User: 200 KAWAI
└─ Total Affiliator: 200 KAWAI

Gas saved: ~75% vs claiming individually
```

---

## 🧪 Testing Plan

### 1. Smart Contract Tests

**File:** `contracts/test/MiningRewardDistributor.t.sol`

Test cases:
- ✅ Set Merkle root (admin only)
- ✅ Claim referral user rewards (85/5/5/5 split)
- ✅ Claim non-referral user rewards (90/5/5 split)
- ✅ Verify Merkle proof validation with all parameters
- ✅ Verify all parties receive correct amounts
- ✅ Prevent double claiming
- ✅ Batch claim multiple periods
- ✅ Event emission with correct data
- ✅ Edge cases (zero amounts, invalid proofs, zero address affiliator)
- ✅ Gas optimization tests

### 2. Backend Tests

**File:** `pkg/store/contributor_test.go`

Test cases:
- ✅ Record reward with referrer (85/5/5/5 split)
- ✅ Record reward without referrer (90/5/5 split)
- ✅ Verify balance updates for all parties (contributor, developer, user, affiliator)
- ✅ Test halving logic with referral splits
- ✅ Test Phase 1 → Phase 2 transition
- ✅ Verify user cashback always 5%
- ✅ Verify developer always 5%
- ✅ Verify contributor gets 85% (referral) or 90% (non-referral)

### 3. Integration Tests

**File:** `internal/api/gateway_test.go`

Test cases:
- ✅ End-to-end job completion with referral
- ✅ Verify reward distribution to contributor, developer, affiliator
- ✅ Test with multiple referral users
- ✅ Test mixed scenario (referral + non-referral users)

---

## 📊 Monitoring & Analytics

### Key Metrics to Track

1. **Referral Conversion Rate**
   - % of referred users who actually mine
   - Average mining volume per referred user

2. **Affiliator Retention**
   - % of affiliators still active after 30/60/90 days
   - Average commission per affiliator per month

3. **Network Growth**
   - Total mining volume: referral vs non-referral users
   - New user acquisition rate
   - Referral program ROI (commission cost vs user LTV)

4. **Economic Impact**
   - Total KAWAI distributed to affiliators
   - Developer revenue impact (10% → 5% reduction)
   - Overall network value growth

### Dashboard Queries

```sql
-- Total mining commission paid to affiliators
SELECT 
  SUM(affiliator_reward) as total_commission,
  COUNT(DISTINCT referrer_address) as unique_affiliators
FROM reward_records
WHERE has_referrer = true;

-- Top affiliators by commission
SELECT 
  referrer_address,
  COUNT(DISTINCT contributor_address) as total_referrals,
  SUM(affiliator_reward) as total_commission
FROM reward_records
WHERE has_referrer = true
GROUP BY referrer_address
ORDER BY total_commission DESC
LIMIT 10;

-- Referral vs Non-Referral mining volume
SELECT 
  has_referrer,
  COUNT(*) as job_count,
  SUM(token_usage) as total_tokens,
  SUM(contributor_reward + developer_reward + affiliator_reward) as total_rewards
FROM reward_records
GROUP BY has_referrer;
```

---

## 🚀 Deployment Plan

### Phase 1: Smart Contract Deployment

1. **Deploy `MiningRewardDistributor.sol`**
   ```bash
   cd contracts
   forge script script/DeployReferralMiningDistributor.s.sol --rpc-url $MONAD_RPC_URL --broadcast
   ```

2. **Verify contract on explorer**
   ```bash
   forge verify-contract <CONTRACT_ADDRESS> MiningRewardDistributor --chain monad-testnet
   ```

3. **Update `.env` with new contract address**
   ```bash
   REFERRAL_MINING_DISTRIBUTOR_ADDRESS=0x...
   ```

### Phase 2: Backend Updates

1. **Update Go bindings**
   ```bash
   make bindings
   ```

2. **Deploy backend changes**
   - Update `pkg/store/contributor.go`
   - Update API gateway handlers
   - Update settlement scripts

3. **Database migration** (if needed)
   - Add referral tracking fields to reward records

### Phase 3: Frontend Updates

1. **Generate TypeScript bindings**
   ```bash
   cd frontend
   npm run bindings
   ```

2. **Implement UI components**
   - `AffiliatorRewardsSection.tsx`
   - Update `RewardsContent.tsx`

3. **Add new API endpoints**
   - `ReferralService.GetMiningCommissions()`

### Phase 4: Testing & Rollout

1. **Internal testing** (1 week)
   - Test with treasury accounts
   - Verify all splits are correct
   - Test settlement & claiming

2. **Beta testing** (1 week)
   - Invite 10-20 early affiliators
   - Monitor metrics closely
   - Gather feedback

3. **Full launch**
   - Announce new referral mining rewards
   - Update documentation
   - Marketing campaign

---

## 📝 Documentation Updates

### Files to Update

1. **`README.md`**
   - Update mining reward split section
   - Add referral mining commission info

2. **`current_concept.md`**
   - Update "Split Ratio" section (lines 27-32)
   - Add referral-based distribution explanation

3. **`REFERRAL_SYSTEM.md`**
   - Add "Mining Commission" section
   - Explain lifetime passive income model

4. **Create `REFERRAL_MINING_COMMISSION.md`**
   - Detailed guide for affiliators
   - How to maximize commission earnings
   - FAQ section

---

## ⚠️ Risks & Mitigation

### 1. **Sybil Attack (Self-Referral)**

**Risk:** Users create multiple accounts to refer themselves and earn 5% commission.

**Mitigation:**
- ✅ Already implemented: Machine ID + Wallet Address tracking
- ✅ Minimum mining threshold before commission kicks in (e.g., 1000 KAWAI mined)
- ✅ KYC requirement for large affiliators (optional)

### 2. **Developer Revenue Impact**

**Risk:** Developer revenue drops from 10% → 5% for referral users.

**Mitigation:**
- ✅ Trade-off is acceptable: 5% cost to acquire users with higher LTV
- ✅ Monitor referral user retention and mining volume
- ✅ If referral users mine 2x more than non-referral, net revenue is positive

### 3. **Smart Contract Bugs**

**Risk:** Bugs in new distributor contract could lock funds or allow exploits.

**Mitigation:**
- ✅ Comprehensive test coverage (>95%)
- ✅ Use OpenZeppelin battle-tested contracts
- ✅ Internal audit before mainnet deployment
- ✅ Start with small amounts, gradually increase

### 4. **Gas Cost Increase**

**Risk:** Additional affiliator transfer increases gas costs.

**Mitigation:**
- ✅ Still using Merkle tree (gas efficient)
- ✅ Batch claiming available
- ✅ Monad has very low gas fees (~$0.01 per tx)

---

## 🎯 Success Criteria

### Short-term (1 month)

- ✅ 50+ active affiliators
- ✅ 30% of new users come from referrals
- ✅ Referral users mine 1.5x more than non-referral users
- ✅ Zero critical bugs in production

### Medium-term (3 months)

- ✅ 200+ active affiliators
- ✅ 50% of new users come from referrals
- ✅ Top 10 affiliators earn >1000 KAWAI/month in commission
- ✅ Positive ROI on referral program

### Long-term (6 months)

- ✅ 500+ active affiliators
- ✅ 70% of new users come from referrals
- ✅ Referral program is primary growth driver
- ✅ Network effect: affiliators recruit sub-affiliators

---

## 🔄 Future Enhancements

### 1. **Multi-Level Referral (MLM-lite)**

Allow affiliators to earn from their referrals' referrals:
- Level 1 (direct): 5% commission
- Level 2 (indirect): 1% commission
- Cap at 2 levels to prevent pyramid scheme concerns

### 2. **Affiliator Tiers**

Reward high-performing affiliators:
- **Bronze** (5+ referrals): 5% commission
- **Silver** (20+ referrals): 6% commission
- **Gold** (50+ referrals): 7% commission
- **Platinum** (100+ referrals): 8% commission

### 3. **Bonus Campaigns**

Temporary promotions:
- "Double Commission Weekend" (10% instead of 5%)
- "Top Affiliator Contest" (extra KAWAI prizes)
- "Referral Sprint" (bonus for most referrals in 1 week)

### 4. **Affiliator NFTs**

Gamification:
- NFT badges for achievement milestones
- Exclusive perks for NFT holders
- Tradeable on marketplace

---

## 📚 References

- [Current Concept Document](current_concept.md)
- [Referral System Documentation](REFERRAL_SYSTEM.md)
- [Cashback System Documentation](CASHBACK_SYSTEM.md)
- [Merkle Distributor Implementation](contracts/contracts/MerkleDistributor.sol)
- [Referral Reward Distributor](contracts/contracts/ReferralRewardDistributor.sol)

---

## ✅ Implementation Checklist

### Smart Contracts
- [ ] Create `MiningRewardDistributor.sol`
- [ ] Write comprehensive tests (>95% coverage)
- [ ] Deploy to Monad testnet
- [ ] Verify on block explorer
- [ ] Generate Go bindings

### Backend
- [ ] Update `RewardRecord` struct with referral fields
- [ ] Modify `RecordJobReward()` to accept referrer address
- [ ] Update API gateway to pass referrer info
- [ ] Modify settlement script for new Merkle tree format
- [ ] Add `GetMiningCommissions()` API endpoint
- [ ] Update reward calculation logic

### Frontend
- [ ] Generate TypeScript bindings
- [ ] Create `AffiliatorRewardsSection.tsx`
- [ ] Add "Mining Commission" tab to Rewards Dashboard
- [ ] Implement commission history table
- [ ] Add stats cards for affiliator metrics
- [ ] Update claim flow for new contract

### Documentation
- [ ] Update `README.md` mining reward section
- [ ] Update `current_concept.md` split ratio section
- [ ] Update `REFERRAL_SYSTEM.md` with mining commission
- [ ] Create `REFERRAL_MINING_COMMISSION.md` guide
- [ ] Update API documentation

### Testing
- [ ] Smart contract unit tests
- [ ] Backend integration tests
- [ ] Frontend UI tests
- [ ] End-to-end tests
- [ ] Load testing (simulate 1000+ users)

### Deployment
- [ ] Deploy contracts to testnet
- [ ] Deploy backend updates
- [ ] Deploy frontend updates
- [ ] Internal testing (1 week)
- [ ] Beta testing (1 week)
- [ ] Full production launch

### Monitoring
- [ ] Set up analytics dashboard
- [ ] Configure alerts for anomalies
- [ ] Track key metrics (conversion, retention, ROI)
- [ ] Weekly performance reports

---

**Last Updated:** 2026-01-04  
**Status:** 🚧 Ready for Implementation  
**Estimated Timeline:** 3-4 weeks

