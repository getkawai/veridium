# Unified Rewards Dashboard - Refactor Plan

**Branch:** `feature/cashback-ui-implementation` (reusing)  
**Date:** 2026-01-04  
**Status:** 🔄 Refactoring in Progress  
**Reason:** Better architecture - integrate Cashback into existing Rewards Dashboard

---

## 🎯 Goal

Refactor from **separate Cashback dashboard** to **unified Rewards dashboard with tabs**:
- Mining Rewards (existing)
- Cashback Rewards (new)
- Referral Rewards (future-ready)

---

## 📊 Architecture Comparison

### Before (Separate Dashboards)
```
Wallet Menu:
├── Home
├── OTC Market
├── Rewards ← Mining only
├── Cashback ← Separate dashboard
└── Settings

Issues:
❌ 2 menu items for rewards
❌ Duplicated claim logic
❌ Inconsistent UX
❌ Hard to maintain
```

### After (Unified Dashboard)
```
Wallet Menu:
├── Home
├── OTC Market
├── Rewards ← All reward types
│   ├── [Tab] Mining Rewards
│   ├── [Tab] Cashback Rewards
│   └── [Tab] Referral Rewards (future)
└── Settings

Benefits:
✅ Single rewards entry point
✅ Shared claim infrastructure
✅ Consistent UX
✅ Easy to maintain
✅ Future-proof for new reward types
```

---

## 🗂️ New File Structure

```
frontend/src/app/wallet/
├── RewardsContent.tsx                    # Modified: Add tabs
├── types.ts                              # Modified: Remove CashbackContentProps
├── wallet.tsx                            # Modified: Remove cashback menu
├── CashbackContent.tsx                   # DELETE: Move to section
└── components/                           # NEW: Modular sections
    └── rewards/
        ├── MiningRewardsSection.tsx      # NEW: Extract from RewardsContent
        ├── CashbackRewardsSection.tsx    # NEW: From CashbackContent
        ├── ReferralRewardsSection.tsx    # NEW: Placeholder for future
        └── shared/
            ├── RewardSummaryCards.tsx    # NEW: Shared component
            ├── ClaimableRewardsList.tsx  # NEW: Shared component
            └── ClaimConfirmModal.tsx     # NEW: Shared component
```

---

## 🔧 Implementation Tasks

### Phase 1: Extract & Modularize (1 hour)

#### Task 1.1: Create Shared Components Directory
```bash
mkdir -p frontend/src/app/wallet/components/rewards/shared
```

#### Task 1.2: Extract MiningRewardsSection
**File:** `components/rewards/MiningRewardsSection.tsx`

Extract current RewardsContent logic into reusable section:
```typescript
interface MiningRewardsSectionProps {
  currentNetwork: NetworkInfo | null;
  transactions: WalletTransaction[];
}

export const MiningRewardsSection = ({ currentNetwork, transactions }: Props) => {
  // All existing RewardsContent logic here
  return (
    <Flexbox gap={20}>
      {/* Summary Cards */}
      {/* Unclaimed Rewards */}
      {/* Pending Claims */}
      {/* Recent Activity */}
    </Flexbox>
  );
};
```

#### Task 1.3: Create CashbackRewardsSection
**File:** `components/rewards/CashbackRewardsSection.tsx`

Move CashbackContent.tsx logic here:
```typescript
interface CashbackRewardsSectionProps {
  currentNetwork: NetworkInfo | null;
}

export const CashbackRewardsSection = ({ currentNetwork }: Props) => {
  // All CashbackContent logic here
  return (
    <Flexbox gap={20}>
      {/* Summary Cards */}
      {/* Tier Progress */}
      {/* Deposit History */}
      {/* Info Section */}
    </Flexbox>
  );
};
```

#### Task 1.4: Create Shared Summary Cards
**File:** `components/rewards/shared/RewardSummaryCards.tsx`

Reusable summary card component:
```typescript
interface RewardSummaryCardsProps {
  cards: Array<{
    title: string;
    value: string;
    suffix: string;
    icon: React.ReactNode;
    gradient: string;
  }>;
  loading?: boolean;
}

export const RewardSummaryCards = ({ cards, loading }: Props) => {
  return (
    <Row gutter={16}>
      {cards.map((card) => (
        <Col key={card.title} xs={24} sm={8}>
          <Card size="small" style={{ background: card.gradient }}>
            {loading ? <Skeleton /> : (
              <Statistic
                title={card.title}
                value={card.value}
                suffix={card.suffix}
                prefix={card.icon}
              />
            )}
          </Card>
        </Col>
      ))}
    </Row>
  );
};
```

### Phase 2: Update RewardsContent (30 minutes)

#### Task 2.1: Add Tabs to RewardsContent
**File:** `RewardsContent.tsx`

```typescript
import { Tabs } from 'antd';
import { Coins, Award, Users } from 'lucide-react';
import { MiningRewardsSection } from './components/rewards/MiningRewardsSection';
import { CashbackRewardsSection } from './components/rewards/CashbackRewardsSection';

const RewardsContent = ({ styles, theme, currentNetwork, transactions }: Props) => {
  const [activeTab, setActiveTab] = useState<'mining' | 'cashback'>('mining');

  return (
    <Flexbox style={{ width: '100%' }} gap={20}>
      {/* Header */}
      <Flexbox horizontal justify="space-between" align="center">
        <div>
          <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Rewards</h2>
          <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>
            Claim your KAWAI rewards from mining, deposits, and referrals
          </span>
        </div>
      </Flexbox>

      {/* Tabs */}
      <Tabs
        activeKey={activeTab}
        onChange={(key) => setActiveTab(key as any)}
        items={[
          {
            key: 'mining',
            label: (
              <span>
                <Coins size={16} style={{ marginRight: 8, verticalAlign: 'middle' }} />
                Mining Rewards
              </span>
            ),
            children: (
              <MiningRewardsSection
                currentNetwork={currentNetwork}
                transactions={transactions}
              />
            ),
          },
          {
            key: 'cashback',
            label: (
              <span>
                <Award size={16} style={{ marginRight: 8, verticalAlign: 'middle' }} />
                Deposit Cashback
              </span>
            ),
            children: (
              <CashbackRewardsSection currentNetwork={currentNetwork} />
            ),
          },
        ]}
      />
    </Flexbox>
  );
};

export default RewardsContent;
```

### Phase 3: Cleanup (30 minutes)

#### Task 3.1: Remove Cashback Menu Item
**File:** `wallet.tsx`

```typescript
// Remove Award import
// Remove 'cashback' from MenuKey type
// Remove cashback menu item from menuItems array
// Remove cashback case from renderContent()

const menuItems = [
  { key: 'home', icon: <Icon icon={Home} />, label: 'Home' },
  { key: 'otc', icon: <Icon icon={ShoppingCart} />, label: 'OTC Market' },
  { key: 'rewards', icon: <Icon icon={Gift} />, label: 'Rewards' }, // ← Single entry
  { key: 'settings', icon: <Icon icon={Settings} />, label: 'Settings' },
];
```

#### Task 3.2: Update Types
**File:** `types.ts`

```typescript
// Remove 'cashback' from MenuKey
export type MenuKey = 'home' | 'otc' | 'rewards' | 'settings';

// Remove CashbackContentProps (no longer needed)
```

#### Task 3.3: Delete CashbackContent.tsx
```bash
rm frontend/src/app/wallet/CashbackContent.tsx
```

---

## 🧪 Testing Checklist

### Functional Tests
- [ ] Mining tab displays correctly
- [ ] Cashback tab displays correctly
- [ ] Tab switching works smoothly
- [ ] All data loads correctly in both tabs
- [ ] Claim functionality works in both tabs
- [ ] Real-time updates work in both tabs

### UI/UX Tests
- [ ] Tabs are visually consistent
- [ ] Icons display correctly
- [ ] Mobile responsive (tabs work on small screens)
- [ ] Loading states work in both tabs
- [ ] Error states work in both tabs
- [ ] Empty states work in both tabs

### Integration Tests
- [ ] Backend services called correctly
- [ ] State management works across tabs
- [ ] Navigation to/from rewards works
- [ ] No console errors

---

## 📊 Code Metrics

### Before Refactor
```
Files: 3 (RewardsContent, CashbackContent, wallet)
Lines: ~1,074 lines
  - RewardsContent.tsx: 612 lines
  - CashbackContent.tsx: 462 lines
Duplication: High (claim logic, summary cards)
Maintainability: Medium
```

### After Refactor
```
Files: 7 (modular structure)
Lines: ~900 lines (15% reduction)
  - RewardsContent.tsx: 150 lines (tabs container)
  - MiningRewardsSection.tsx: 400 lines
  - CashbackRewardsSection.tsx: 350 lines
  - Shared components: ~100 lines
Duplication: Low (shared components)
Maintainability: High
```

**Improvement:**
- ✅ 15% less code
- ✅ Better modularity
- ✅ Easier to test
- ✅ Easier to extend

---

## 🚀 Deployment Plan

### Step 1: Implement Refactor (2 hours)
```bash
# Already on feature/cashback-ui-implementation branch
git status

# Create component directories
mkdir -p frontend/src/app/wallet/components/rewards/shared

# Implement tasks 1.1 - 3.3
# ...

# Commit refactored code
git add -A
git commit -m "refactor: unify rewards dashboard with tabs"
```

### Step 2: Test Locally
```bash
# Build and test
npm run build
make dev

# Test all functionality
# - Tab switching
# - Data loading
# - Claim functionality
```

### Step 3: Create New PR
```bash
# Push changes
git push origin feature/cashback-ui-implementation

# Create PR
gh pr create \
  --base master \
  --head feature/cashback-ui-implementation \
  --title "feat: Add Cashback to Unified Rewards Dashboard" \
  --body "See UNIFIED_REWARDS_REFACTOR_PLAN.md"
```

---

## 📝 Commit Messages

```bash
# Commit 1: Create modular structure
git commit -m "refactor: extract MiningRewardsSection from RewardsContent

- Create components/rewards/ directory structure
- Extract mining rewards logic to separate section
- Prepare for unified rewards dashboard"

# Commit 2: Add cashback section
git commit -m "feat: add CashbackRewardsSection component

- Move CashbackContent logic to section component
- Integrate with unified rewards structure
- Add tier visualization and deposit history"

# Commit 3: Add shared components
git commit -m "refactor: create shared reward components

- Add RewardSummaryCards for reusability
- Add ClaimConfirmModal for consistent UX
- Reduce code duplication across sections"

# Commit 4: Update RewardsContent with tabs
git commit -m "feat: implement unified rewards dashboard with tabs

- Add Tabs component to RewardsContent
- Integrate Mining and Cashback sections
- Update header to reflect all reward types"

# Commit 5: Cleanup
git commit -m "refactor: remove separate cashback menu item

- Remove cashback from wallet menu
- Update types to remove CashbackContentProps
- Delete standalone CashbackContent.tsx
- Simplify navigation structure"
```

---

## 🎯 Success Criteria

### Must Have
- [x] PR #41 closed with explanation
- [ ] Unified dashboard with 2 tabs working
- [ ] All existing functionality preserved
- [ ] No regressions in mining rewards
- [ ] Cashback section fully functional
- [ ] Code is more maintainable
- [ ] Mobile responsive

### Should Have
- [ ] Shared components reduce duplication
- [ ] Consistent UX across tabs
- [ ] Smooth tab transitions
- [ ] Loading states work correctly

### Nice to Have
- [ ] Animated tab transitions
- [ ] Tab state persists in URL
- [ ] Keyboard navigation (arrow keys)
- [ ] Tab badges showing counts

---

## 📚 Documentation Updates

### Files to Update
1. ✅ `UNIFIED_REWARDS_REFACTOR_PLAN.md` (this file)
2. ⏳ `CASHBACK_UI_IMPLEMENTATION_PLAN.md` (mark as superseded)
3. ⏳ `CASHBACK_IMPLEMENTATION_SUMMARY.md` (update with refactor notes)
4. ⏳ `README.md` (update if needed)

---

## 💡 Future Enhancements

### Phase 4: Add Referral Rewards Tab (Future)
```typescript
// Easy to add third tab
{
  key: 'referral',
  label: <span><Users size={16} /> Referral Rewards</span>,
  children: <ReferralRewardsSection />,
}
```

### Phase 5: Add Rewards Overview (Future)
```typescript
// Optional: Add "All Rewards" tab showing combined view
{
  key: 'overview',
  label: <span><BarChart size={16} /> Overview</span>,
  children: <RewardsOverviewSection />,
}
```

---

## ⏱️ Timeline

| Task | Duration | Status |
|------|----------|--------|
| Close PR #41 | 5 min | ✅ Done |
| Create refactor plan | 15 min | ✅ Done |
| Extract MiningRewardsSection | 30 min | ⏳ Next |
| Create CashbackRewardsSection | 30 min | ⏳ Pending |
| Create shared components | 30 min | ⏳ Pending |
| Update RewardsContent | 30 min | ⏳ Pending |
| Cleanup & remove old code | 30 min | ⏳ Pending |
| Testing | 30 min | ⏳ Pending |
| Create new PR | 10 min | ⏳ Pending |
| **Total** | **3 hours** | **5% Done** |

---

## 🎉 Summary

**Status:** 🔄 Refactoring Started  
**Approach:** Unified Rewards Dashboard with Tabs  
**Benefits:** Better UX, less code, easier maintenance  
**Timeline:** ~3 hours total  
**Next:** Extract MiningRewardsSection component

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-04  
**Status:** 🟢 In Progress

