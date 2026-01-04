# Cashback System UI Implementation Plan

**Branch:** `feature/cashback-ui-implementation`  
**Created:** 2026-01-04  
**Estimated Time:** 3-4 hours  
**Priority:** P0 - Critical (Missing Major Feature)

---

## 📋 Executive Summary

Implementasi UI untuk Cashback System yang sudah tersedia di backend. Backend bindings sudah lengkap dengan 2 methods dari `CashbackService`, namun belum ada UI component yang menggunakannya.

**Backend Status:** ✅ Complete (2/2 methods)  
**UI Status:** ❌ Missing (0% implemented)  
**Integration Status:** ❌ Not integrated

---

## 🎯 Goals & Objectives

### Primary Goals
1. ✅ Create comprehensive Cashback Dashboard UI
2. ✅ Display tiered cashback rates (1%-5%)
3. ✅ Show user's cashback statistics
4. ✅ Implement cashback claim interface
5. ✅ Integrate with existing Rewards system

### Success Metrics
- [ ] Users can view their cashback stats
- [ ] Users can see tier progression
- [ ] Users can claim accumulated cashback
- [ ] UI matches design system (Ant Design + @lobehub/ui)
- [ ] Real-time updates via Wails Events

---

## 📊 Backend Analysis

### Available Services

#### 1. CashbackService.GetCashbackStats()
```typescript
// Input
userAddress: string

// Output: CashbackStatsResponse
{
  totalDeposits: string;           // Total USDT deposited
  totalCashbackEarned: string;     // Total KAWAI earned
  currentTier: number;             // 0-4 (tier index)
  nextTierAmount: string;          // USDT needed for next tier
  claimableAmount: string;         // KAWAI ready to claim
  pendingAmount: string;           // KAWAI accumulating (not settled)
  lastClaimPeriod: number;         // Last claimed period ID
  depositHistory: Array<{
    amount: string;                // USDT deposited
    cashbackAmount: string;        // KAWAI earned
    cashbackRate: number;          // Rate applied (1-5%)
    timestamp: string;             // ISO date
    txHash: string;                // Transaction hash
    claimed: boolean;              // Claim status
  }>;
}
```

#### 2. CashbackService.GetCurrentPeriod()
```typescript
// Output
periodID: number  // Current settlement period
```

### Cashback Tiers (from README.md)

| Tier | Deposit Range | Cashback Rate | Cap per Deposit |
|------|---------------|---------------|-----------------|
| 0 | $0 - $99 | 1% | 5,000 KAWAI |
| 1 | $100 - $499 | 2% | 10,000 KAWAI |
| 2 | $500 - $999 | 3% | 15,000 KAWAI |
| 3 | $1,000 - $4,999 | 4% | 18,000 KAWAI |
| 4 | $5,000+ | 5% | 20,000 KAWAI |

**Special:** First-time deposit gets 5% regardless of amount (overrides base rate).

---

## 🎨 UI Design Specification

### Component Structure

```
CashbackContent.tsx (Main Component)
├── CashbackSummaryCards (Stats Overview)
│   ├── TotalCashbackCard
│   ├── ClaimableCard
│   └── PendingCard
├── TierProgressSection (Visual Tier Display)
│   ├── CurrentTierBadge
│   ├── TierProgressBar
│   └── NextTierInfo
├── DepositHistoryTable (Transaction List)
│   ├── DepositRow
│   └── ClaimButton
└── HowItWorksCard (Info Section)
```

### Visual Mockup

```
┌─────────────────────────────────────────────────────────────┐
│ Cashback Dashboard                            [Refresh] [?] │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ 💰 Total     │  │ ✅ Claimable │  │ ⏳ Pending   │      │
│  │ 12,500 KAWAI │  │ 8,000 KAWAI  │  │ 4,500 KAWAI  │      │
│  │ from 5 deps  │  │ Ready now    │  │ This period  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│ Your Cashback Tier                                           │
│                                                               │
│  Current: Tier 2 (3% Cashback) 🎯                           │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  $750 / $1,000                                               │
│  Deposit $250 more to unlock Tier 3 (4% Cashback)          │
│                                                               │
│  Tier 0  Tier 1  Tier 2  Tier 3  Tier 4                    │
│    1%     2%     3%✓     4%      5%                         │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│ Deposit History                                              │
│                                                               │
│  Date       Amount    Rate  Cashback    Status    Action    │
│  ────────────────────────────────────────────────────────── │
│  Jan 3     $500      3%    1,500 KAWAI  ✅ Claimed  [View] │
│  Jan 2     $250      5%*   1,250 KAWAI  ⏳ Pending  [Claim]│
│  Dec 30    $100      2%    200 KAWAI    ✅ Claimed  [View] │
│                                                               │
│  * First deposit bonus (5%)                                  │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│ How Cashback Works                                           │
│  • Earn 1-5% KAWAI cashback on every USDT deposit           │
│  • First deposit always gets 5% bonus                        │
│  • Higher deposits unlock better rates                       │
│  • Claim weekly after settlement                             │
│  • 200M KAWAI allocated (~3 year runway)                     │
└─────────────────────────────────────────────────────────────┘
```

---

## 🛠️ Implementation Tasks

### Phase 1: Component Creation (1.5 hours)

#### Task 1.1: Create CashbackContent.tsx
- [ ] Create main component file
- [ ] Import necessary dependencies (Ant Design, icons, services)
- [ ] Set up component state management
- [ ] Implement data fetching from CashbackService

**File:** `frontend/src/app/wallet/CashbackContent.tsx`

```typescript
import { Card, Table, Progress, Tag, Statistic, Row, Col, Button, Empty, Skeleton } from 'antd';
import { Gift, TrendingUp, Clock, Award, Info } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { useState, useEffect } from 'react';
import { CashbackService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { CashbackStatsResponse } from '@@/github.com/kawai-network/veridium/internal/services/models';
import type { CashbackContentProps } from './types';

const CashbackContent = ({ styles, theme, currentNetwork }: CashbackContentProps) => {
  // State management
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<CashbackStatsResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  
  // ... implementation
};
```

#### Task 1.2: Create Summary Cards Component
- [ ] Total Cashback Earned card
- [ ] Claimable Amount card
- [ ] Pending Amount card
- [ ] Loading skeletons for each card

```typescript
const CashbackSummaryCards = ({ stats, loading, theme }) => {
  return (
    <Row gutter={16}>
      <Col span={8}>
        <Card size="small" style={{ background: 'linear-gradient(135deg, #667eea20, #764ba220)' }}>
          {loading ? <Skeleton /> : (
            <Statistic
              title="Total Cashback Earned"
              value={stats?.totalCashbackEarned || '0'}
              suffix="KAWAI"
              prefix={<Gift size={20} />}
            />
          )}
        </Card>
      </Col>
      {/* Claimable & Pending cards */}
    </Row>
  );
};
```

#### Task 1.3: Create Tier Progress Component
- [ ] Current tier badge with visual indicator
- [ ] Progress bar showing deposits toward next tier
- [ ] Next tier unlock information
- [ ] All 5 tiers visualization

```typescript
const TierProgressSection = ({ stats, theme }) => {
  const tiers = [
    { level: 0, min: 0, max: 99, rate: 1, cap: 5000 },
    { level: 1, min: 100, max: 499, rate: 2, cap: 10000 },
    { level: 2, min: 500, max: 999, rate: 3, cap: 15000 },
    { level: 3, min: 1000, max: 4999, rate: 4, cap: 18000 },
    { level: 4, min: 5000, max: Infinity, rate: 5, cap: 20000 },
  ];
  
  const currentTier = tiers[stats?.currentTier || 0];
  const nextTier = tiers[Math.min((stats?.currentTier || 0) + 1, 4)];
  const progress = calculateProgress(stats?.totalDeposits, currentTier, nextTier);
  
  return (
    <Card title="Your Cashback Tier" size="small">
      {/* Tier visualization */}
    </Card>
  );
};
```

#### Task 1.4: Create Deposit History Table
- [ ] Table with columns: Date, Amount, Rate, Cashback, Status, Action
- [ ] Claim button for unclaimed deposits
- [ ] Transaction hash link to explorer
- [ ] First deposit indicator (5% bonus)
- [ ] Pagination support

```typescript
const depositHistoryColumns = [
  {
    title: 'Date',
    dataIndex: 'timestamp',
    render: (date: string) => new Date(date).toLocaleDateString(),
  },
  {
    title: 'Deposit Amount',
    dataIndex: 'amount',
    render: (amount: string) => `$${parseFloat(amount).toFixed(2)} USDT`,
  },
  {
    title: 'Rate',
    dataIndex: 'cashbackRate',
    render: (rate: number, record: any) => (
      <span>
        {rate}%
        {record.isFirstDeposit && <Tag color="gold" size="small">First Deposit Bonus</Tag>}
      </span>
    ),
  },
  // ... more columns
];
```

### Phase 2: Integration (1 hour)

#### Task 2.1: Add to Wallet Layout
- [ ] Add Cashback tab to wallet menu
- [ ] Update wallet.tsx to include CashbackContent
- [ ] Add icon for Cashback menu item
- [ ] Update types.ts with CashbackContentProps

**File:** `frontend/src/app/wallet/wallet.tsx`

```typescript
// Add to MenuContent
const menuItems = [
  { key: 'home', icon: <Icon icon={Home} />, label: 'Home' },
  { key: 'otc', icon: <Icon icon={ShoppingCart} />, label: 'OTC Market' },
  { key: 'rewards', icon: <Icon icon={Gift} />, label: 'Rewards' },
  { key: 'cashback', icon: <Icon icon={Award} />, label: 'Cashback' }, // NEW
  { key: 'settings', icon: <Icon icon={Settings} />, label: 'Settings' },
];

// Add to renderContent()
case 'cashback':
  return <CashbackContent styles={styles} theme={theme} currentNetwork={currentNetwork} />;
```

#### Task 2.2: Update Types
- [ ] Add CashbackContentProps interface
- [ ] Update MenuKey type to include 'cashback'

**File:** `frontend/src/app/wallet/types.ts`

```typescript
export type MenuKey = 'home' | 'otc' | 'rewards' | 'cashback' | 'settings';

export interface CashbackContentProps {
  styles: any;
  theme: any;
  currentNetwork: NetworkInfo | null;
}
```

#### Task 2.3: Connect to Backend Services
- [ ] Implement loadCashbackStats() function
- [ ] Add error handling with user-friendly messages
- [ ] Implement refresh functionality
- [ ] Add loading states

```typescript
const loadCashbackStats = async (userAddress: string) => {
  setLoading(true);
  setError(null);
  try {
    const result = await CashbackService.GetCashbackStats(userAddress);
    if (!result) {
      setError('Failed to load cashback data');
      return;
    }
    setStats(result);
  } catch (e: any) {
    console.error('Failed to load cashback stats:', e);
    setError(e.message || 'Failed to load cashback data');
  } finally {
    setLoading(false);
  }
};
```

### Phase 3: Claim Integration (1 hour)

#### Task 3.1: Integrate with Rewards Claim System
- [ ] Reuse claim logic from RewardsContent
- [ ] Add cashback-specific claim handling
- [ ] Implement claim confirmation modal
- [ ] Add gas estimation for claim transactions

```typescript
const handleClaimCashback = async (depositTxHash: string, amount: string) => {
  // Reuse DeAIService.ClaimKawaiReward() with cashback-specific parameters
  // Show confirmation modal with gas estimate
  // Execute claim transaction
  // Update UI after successful claim
};
```

#### Task 3.2: Add Real-time Updates
- [ ] Subscribe to cashback-related Wails Events
- [ ] Update stats on deposit events
- [ ] Update stats on claim events
- [ ] Auto-refresh after period settlement

```typescript
useEffect(() => {
  // Subscribe to deposit events
  const unsubscribeDeposit = Events.On('wallet:deposit_confirmed', handleDepositConfirmed);
  
  // Subscribe to cashback claim events
  const unsubscribeClaim = Events.On('cashback:claim_completed', handleClaimCompleted);
  
  return () => {
    unsubscribeDeposit();
    unsubscribeClaim();
  };
}, []);
```

### Phase 4: Polish & Testing (0.5 hours)

#### Task 4.1: Add Visual Enhancements
- [ ] Gradient backgrounds for tier cards
- [ ] Animated progress bars
- [ ] Hover effects on interactive elements
- [ ] Responsive design for mobile

#### Task 4.2: Add Empty States
- [ ] No deposits yet state
- [ ] No claimable cashback state
- [ ] Error state with retry button

```typescript
if (!stats?.depositHistory || stats.depositHistory.length === 0) {
  return (
    <Empty
      image={Empty.PRESENTED_IMAGE_SIMPLE}
      description={
        <span>
          No deposits yet. Make your first deposit to start earning cashback!
          <br />
          <strong>First deposit gets 5% bonus!</strong>
        </span>
      }
    >
      <Button type="primary" onClick={() => setModalType('deposit')}>
        Make First Deposit
      </Button>
    </Empty>
  );
}
```

#### Task 4.3: Add Info Tooltips
- [ ] Tier system explanation
- [ ] First deposit bonus tooltip
- [ ] Cap per deposit explanation
- [ ] Settlement period info

```typescript
<Tooltip title="First deposit always gets 5% cashback regardless of amount">
  <Info size={14} style={{ cursor: 'help' }} />
</Tooltip>
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] Component renders without crashing
- [ ] Handles null/undefined stats gracefully
- [ ] Calculates tier progress correctly
- [ ] Formats amounts correctly (USDT/KAWAI)
- [ ] Handles loading states properly

### Integration Tests
- [ ] Backend service calls work correctly
- [ ] Error handling displays user-friendly messages
- [ ] Refresh functionality updates data
- [ ] Navigation to/from cashback tab works

### User Acceptance Tests
- [ ] User can view total cashback earned
- [ ] User can see current tier and progress
- [ ] User can view deposit history
- [ ] User can claim available cashback
- [ ] First deposit bonus is clearly indicated
- [ ] Mobile responsive design works

### Edge Cases
- [ ] No wallet connected
- [ ] No deposits made yet
- [ ] All cashback already claimed
- [ ] Network error during load
- [ ] Claim transaction fails

---

## 📁 File Structure

```
frontend/src/app/wallet/
├── wallet.tsx                    # Modified: Add cashback menu item
├── types.ts                      # Modified: Add CashbackContentProps
├── CashbackContent.tsx           # NEW: Main cashback component
└── components/                   # NEW: Cashback sub-components
    ├── CashbackSummaryCards.tsx
    ├── TierProgressSection.tsx
    ├── DepositHistoryTable.tsx
    └── ClaimCashbackModal.tsx
```

---

## 🎨 Design System Compliance

### Colors (from existing theme)
- **Primary Gradient:** `linear-gradient(135deg, #667eea, #764ba2)` (KAWAI theme)
- **Success:** `#22c55e` (Claimable amounts)
- **Warning:** `#f59e0b` (Pending amounts)
- **Info:** `#3b82f6` (Tier progress)
- **Card Background:** `theme.colorBgContainer`
- **Border:** `theme.colorBorderSecondary`

### Typography
- **Heading:** 20px, font-weight: 600
- **Stat Value:** 24px, font-weight: 700
- **Body:** 13px, font-weight: 400
- **Caption:** 11px, color: colorTextTertiary

### Spacing
- **Card Gap:** 16px
- **Section Gap:** 20px
- **Inner Padding:** 16px (Card), 24px (Main content)

---

## 🔗 Dependencies

### Existing Dependencies (No new installs needed)
- ✅ `antd` - UI components
- ✅ `@lobehub/ui` - Enhanced UI components
- ✅ `lucide-react` - Icons
- ✅ `react-layout-kit` - Flexbox utilities
- ✅ `antd-style` - Styling utilities

### Backend Services
- ✅ `CashbackService` - Already available in bindings
- ✅ `DeAIService` - For claim transactions (reuse existing)
- ✅ `JarvisService` - For gas estimation (reuse existing)

---

## 📊 Success Metrics

### Completion Criteria
- [x] Branch created: `feature/cashback-ui-implementation`
- [ ] All 4 phases completed
- [ ] All testing checklist items passed
- [ ] Code review approved
- [ ] Merged to master

### Performance Targets
- Load time: < 500ms
- Smooth animations: 60fps
- No layout shifts (CLS < 0.1)
- Mobile responsive: All breakpoints

### User Experience Goals
- Intuitive tier progression visualization
- Clear call-to-action for claiming
- Helpful tooltips and info sections
- Error messages are actionable

---

## 🚀 Deployment Plan

### Pre-deployment
1. [ ] Run linter: `npm run lint`
2. [ ] Build frontend: `npm run build`
3. [ ] Test in dev mode: `make dev`
4. [ ] Test production build: `make build`

### Deployment Steps
1. [ ] Commit changes with descriptive message
2. [ ] Push branch to remote
3. [ ] Create Pull Request with this plan attached
4. [ ] Request code review
5. [ ] Address review comments
6. [ ] Merge to master after approval

### Post-deployment
1. [ ] Monitor Sentry for errors
2. [ ] Collect user feedback
3. [ ] Track cashback claim success rate
4. [ ] Iterate based on analytics

---

## 📝 Notes & Considerations

### Design Decisions
1. **Placement:** Added as separate tab in wallet (not sub-section of Rewards)
   - **Rationale:** Cashback is distinct from mining rewards, deserves own space
   
2. **Tier Visualization:** Horizontal progress bar with all 5 tiers shown
   - **Rationale:** Users can see full progression path at a glance
   
3. **Claim Integration:** Reuse existing claim infrastructure from Rewards
   - **Rationale:** Consistent UX, avoid code duplication

### Future Enhancements (Out of Scope)
- [ ] Cashback calculator (preview cashback before deposit)
- [ ] Tier comparison modal (detailed tier benefits)
- [ ] Cashback leaderboard (gamification)
- [ ] Push notifications for claimable cashback
- [ ] Cashback referral bonus (earn extra for referred deposits)

### Known Limitations
- Cashback data depends on backend settlement (weekly)
- Claims require gas fees (MON on Monad Testnet)
- First deposit detection relies on backend logic

---

## 👥 Team & Resources

### Assignee
- **Developer:** TBD
- **Reviewer:** TBD
- **Designer:** N/A (using existing design system)

### Reference Materials
- [Cashback System Documentation](CASHBACK_SYSTEM.md)
- [Deposit Cashback Tokenomics](DEPOSIT_CASHBACK_TOKENOMICS.md)
- [README.md](README.md) - Phase 2 implementation details
- [Existing Rewards UI](frontend/src/app/wallet/RewardsContent.tsx) - Reference for claim flow

### Support Channels
- **Backend Issues:** Check `internal/services/cashback_service.go`
- **UI Issues:** Reference `OTCContent.tsx` and `RewardsContent.tsx`
- **Design Questions:** Follow Ant Design + @lobehub/ui patterns

---

## ✅ Acceptance Criteria

### Must Have (P0)
- [x] Branch created successfully
- [ ] CashbackContent component created and functional
- [ ] Summary cards display correct stats
- [ ] Tier progress visualization works
- [ ] Deposit history table populated
- [ ] Claim functionality integrated
- [ ] Error handling implemented
- [ ] Loading states work correctly
- [ ] Mobile responsive

### Should Have (P1)
- [ ] Real-time updates via events
- [ ] Gas estimation for claims
- [ ] Info tooltips for user guidance
- [ ] Empty states for no data
- [ ] Smooth animations

### Nice to Have (P2)
- [ ] Cashback calculator preview
- [ ] Tier comparison modal
- [ ] Export deposit history (CSV)
- [ ] Share cashback stats (social)

---

## 📅 Timeline

| Phase | Duration | Start | End |
|-------|----------|-------|-----|
| Phase 1: Component Creation | 1.5h | Day 1 | Day 1 |
| Phase 2: Integration | 1h | Day 1 | Day 1 |
| Phase 3: Claim Integration | 1h | Day 1 | Day 1 |
| Phase 4: Polish & Testing | 0.5h | Day 1 | Day 1 |
| **Total** | **4h** | - | - |

**Target Completion:** Same day (single sprint)

---

## 🎯 Next Steps

1. ✅ **Create branch** - `feature/cashback-ui-implementation` (DONE)
2. ✅ **Create implementation plan** - This document (DONE)
3. ⏳ **Start Phase 1** - Create CashbackContent.tsx component
4. ⏳ **Continue with Phase 2-4** - Follow task checklist
5. ⏳ **Testing & Review** - Complete all checklist items
6. ⏳ **Merge to master** - After approval

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-04  
**Status:** 🟢 Ready to Start

