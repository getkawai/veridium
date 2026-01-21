# 🎁 Referral System Documentation

## Overview

The Kawai DeAI Network referral system incentivizes viral growth by rewarding both referrers and new users.

## Incentive Structure

### One-Time Bonuses

| User Type | Bonus Amount | Description |
|-----------|--------------|-------------|
| **New User (No Referral)** | 5 USDT | Base trial bonus |
| **New User (With Referral)** | 10 USDT | Enhanced trial bonus (+100%) |
| **Referrer** | 5 USDT | One-time reward per successful referral |

### Lifetime Mining Commission (NEW) 🌟

| User Type | Commission Rate | Description |
|-----------|----------------|-------------|
| **Affiliator (Referrer)** | **5%** | Lifetime commission from referral's mining rewards |

**Example:**
- Your referral mines 1000 KAWAI
- You automatically earn **50 KAWAI** (5% commission)
- This continues **forever** as long as they keep mining
- The more they mine, the more you earn!

### Why This Works

- **Win-Win:** Both parties benefit
- **Strong Incentive:** 5 USDT one-time + 5% lifetime commission
- **Passive Income:** Earn while your referrals mine
- **Viral Potential:** Coefficient of 2.0-3.0x
- **No Limit:** Unlimited referrals = unlimited earnings
- **Long-term Value:** Affiliators have incentive to support their referrals

---

## User Flow

### For New Users

1. **Receive Referral Code**
   - Friend shares their 6-digit code: `ABC123`
   - Via social media, messaging, or word of mouth

2. **Open App & Enter Referral Code**
   - Click "Have a Referral Code?" banner
   - Manually enter the 6-digit code
   - Code is validated (format: 6 alphanumeric characters)
   - Code is saved locally for later use

3. **Code Applied**
   - Banner shows: "🎉 Bonus Upgraded!"
   - Bonus upgraded: 5 USDT + 100 KAWAI → 10 USDT + 200 KAWAI
   - Success message confirms application
   - Referral code stored in browser localStorage

4. **Create/Setup Wallet**
   - Generate new mnemonic or import existing wallet
   - Create wallet with password
   - Wallet is created successfully

5. **Auto-Claim Bonus on First Unlock**
   - **Automatically triggered** when unlocking wallet for the first time
   - Backend validates referral code and machine ID
   - Trial bonus claimed with referral upgrade
   - Receives 10 USDT + 200 KAWAI (instead of 5 USDT + 100 KAWAI)
   - Referrer receives 5 USDT + 100 KAWAI
   - Referral code cleared from localStorage after successful claim

### For Referrers

1. **Generate Referral Code**
   ```typescript
   // Desktop App (Wails)
   import { CreateReferralCode } from '@/bindings/services/referralservice';
   
   const code = await CreateReferralCode(userAddress);
   console.log(code); // "ABC123"
   ```

2. **Share Code**
   - Copy code button in Rewards > Referral tab
   - Share via social media (native share API)
   - Direct message to friends
   - Users manually enter the code

3. **Track Earnings**
   ```typescript
   // Desktop App (Wails)
   import { GetReferralStats } from '@/bindings/services/referralservice';
   
   const stats = await GetReferralStats(userAddress);
   console.log(stats.code);                // "ABC123"
   console.log(stats.total_referrals);     // 15
   console.log(stats.total_earnings_usdt);  // 45.00
   console.log(stats.total_earnings_kawai); // "1500000000000000000000"
   ```

---

## Technical Implementation

### Backend (Go)

#### 1. Referral Code Generation

```go
// pkg/store/referral.go

// Generates unique 6-character alphanumeric code
func GenerateReferralCode() (string, error)

// Creates referral code for user
func (s *KVStore) CreateReferralCode(ctx context.Context, address string) (*ReferralData, error)
```

**Features:**
- 6 characters (UPPERCASE + NUMBERS)
- Excludes confusing chars (0, O, 1, I)
- Collision detection
- Stored in KV: `referral:code:{CODE}`
- Address mapping: `referral:address:{ADDRESS}`

#### 2. Trial Claim with Referral

```go
// pkg/store/referral.go

func (s *KVStore) ClaimFreeTrialWithReferral(
    ctx context.Context,
    address string,
    machineID string,
    referralCode string,
) (int64, error)
```

**Logic:**
1. Validate referral code (if provided)
2. Check for self-referral
3. Determine bonus amount (5 or 10 USDT)
4. Claim trial with atomic operation
5. Reward referrer (5 USDT)
6. Update referral stats

#### 2. Auto-Claim Trial on Wallet Unlock (NEW ✨)

```go
// internal/services/wallet_service.go

func (s *WalletService) AutoClaimTrialIfNeeded(referralCode string) (bool, float64, string, error)
```

**Features:**
- Called automatically after successful wallet unlock
- Checks if trial already claimed (address + machine ID)
- Validates referral code if provided
- Claims trial with referral bonus
- Returns: `(claimed bool, usdtAmount float64, kawaiAmount string, error)`

**Anti-Abuse Protection:**
- Requires machine ID (fail-closed if unavailable)
- Prevents duplicate claims by address
- Prevents duplicate claims by machine ID
- Prevents self-referral

**Flow:**
1. Get machine ID (required for anti-abuse)
2. Check if trial already claimed
3. Validate referral code (if provided)
4. Claim trial with `ClaimFreeTrialWithReferral`
5. Return claimed amounts

#### 3. Wails Service (Desktop App)

```go
// internal/services/referralservice.go

type ReferralService struct {
    kvStore *store.KVStore
}

// Wails-exposed methods (auto-generated TypeScript bindings)
func (s *ReferralService) CreateReferralCode(userAddress string) (*ReferralCodeResponse, error)
func (s *ReferralService) GetReferralStats(userAddress string) (*ReferralStatsResponse, error)
func (s *ReferralService) GetReferralBonusAmounts() map[string]interface{}

// Auto-claim is handled by WalletService
func (s *WalletService) AutoClaimTrialIfNeeded(referralCode string) (bool, float64, string, error)
```

**Note:** For backend API endpoints (contributor/gateway), see `pkg/gateway/handler_referral.go`

### Frontend (React + TypeScript)

#### 1. Referral Banner Component

```typescript
// frontend/src/features/Referral/ReferralBanner.tsx

<ReferralBanner 
  onReferralApplied={(code) => {
    setReferralCode(code);
    setHasReferral(true);
  }}
/>
```

**Features:**
- Collapsible input field
- Code validation (6 alphanumeric)
- Success animation
- Bonus comparison (5 vs 10 USDT)
- **Saves code to localStorage** for auto-claim

#### 2. Auto-Claim on Unlock (NEW ✨)

```typescript
// frontend/src/store/user/slices/wallet/action.ts

unlockWallet: async (password: string) => {
  await WalletService.UnlockWallet(password);
  await get().refreshWalletStatus();
  
  // Auto-claim trial if needed
  const pendingReferralCode = localStorage.getItem('pendingReferralCode') || '';
  const [claimed, usdtAmount, kawaiAmount] = await WalletService.AutoClaimTrialIfNeeded(pendingReferralCode);
  
  if (claimed) {
    localStorage.removeItem('pendingReferralCode');
    console.log(`🎉 Free trial claimed: ${usdtAmount} USDT + ${kawaiAmount} KAWAI`);
  }
}
```

**Features:**
- Retrieves referral code from localStorage
- Calls `AutoClaimTrialIfNeeded` after unlock
- Clears referral code after successful claim
- Non-blocking: doesn't fail unlock if claim fails
- Logs success message with claimed amounts

#### 3. Referral Rewards Section

```typescript
// frontend/src/app/wallet/components/rewards/ReferralRewardsSection.tsx
// Integrated in unified Rewards Dashboard

<ReferralRewardsSection
  currentNetwork={currentNetwork}
  theme={theme}
  styles={styles}
  onRefresh={(refreshFn) => { /* callback */ }}
/>
```

**Features:**
- Display referral code (large, copyable)
- Statistics cards (referrals, USDT earned, KAWAI earned)
- Copy code button with success feedback
- Native share API integration
- Benefits breakdown (friend gets / you get)
- Step-by-step "How It Works" guide
- High-precision KAWAI formatting (18 decimals)

**Location:**
- Accessible via: Wallet → Rewards → Referral Rewards tab
- Part of unified rewards dashboard (Mining | Cashback | Referral)

#### 4. Storage & Persistence

**localStorage Keys:**
- `pendingReferralCode`: Stores referral code until trial is claimed
- Automatically cleared after successful claim
- Persists across app restarts until used

**Flow:**
```
User applies code → localStorage.setItem('pendingReferralCode', code)
User unlocks wallet → Read from localStorage
Trial claimed → localStorage.removeItem('pendingReferralCode')
```

```typescript
// In AuthSignInBox.tsx
// Users manually enter referral code via ReferralBanner

<ReferralBanner 
  onReferralApplied={(code) => {
    setReferralCode(code);
    setHasReferral(true);
  }}
/>
```

**Note:** URL parameter detection (`?ref=CODE`) has been removed for simplicity. Users now manually enter the 6-digit code during wallet setup.

---

## Data Models

### ReferralData (Go)

```go
type ReferralData struct {
    Code           string    `json:"code"`            // ABC123
    OwnerAddress   string    `json:"owner_address"`   // 0x1234...
    TotalReferrals int       `json:"total_referrals"` // 15
    TotalEarnings  int64     `json:"total_earnings"`  // 45000000 (micro USDT)
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

### ReferralClaim (Go)

```go
type ReferralClaim struct {
    ReferralCode   string     `json:"referral_code"`
    ReferredUser   string     `json:"referred_user"`
    ReferrerReward int64      `json:"referrer_reward"`
    Status         string     `json:"status"` // "pending", "completed"
    CreatedAt      time.Time  `json:"created_at"`
    CompletedAt    *time.Time `json:"completed_at,omitempty"`
}
```

---

## Storage Schema (Cloudflare KV)

```
# Referral Code → Data
referral:code:ABC123 → {ReferralData JSON}

# Address → Referral Code
referral:address:0x1234... → "ABC123"

# Referral Claims
referral:claim:ABC123:0x5678... → {ReferralClaim JSON}

# Machine ID (Anti-Abuse)
trial_machine:MACHINE_ID_123 → "true"
```

---

## Anti-Abuse Mechanisms

### 1. Self-Referral Prevention
```go
if strings.EqualFold(referralData.OwnerAddress, address) {
    return fmt.Errorf("cannot use your own referral code")
}
```

### 2. Dual-Layer Trial Claim Protection
- **Wallet Address:** One claim per address
- **Machine ID:** One claim per device (required, fail-closed)

### 3. Machine ID Requirement (NEW ✨)
```go
// Fail-closed: require machine ID for anti-abuse
machineID, err := s.getMachineID()
if err != nil {
    return false, 0, "0", fmt.Errorf("machine id unavailable: %w", err)
}
```

**Benefits:**
- Prevents bypass of machine-level checks
- More secure than soft-fail approach
- Appropriate for valuable one-time bonuses
- Machine ID should always be available in desktop app

### 4. Atomic Operations
- Race condition protection
- Retry mechanism with exponential backoff
- Prevents double-claiming

### 5. Auto-Claim Validation
- Checks trial status before claiming
- Validates referral code format and existence
- Verifies referrer is not the same as user
- All validations happen server-side

---

## Marketing Copy

### Landing Page

```
🎁 Get 5 USDT + 100 KAWAI FREE to Chat with AI
→ Or get 10 USDT + 200 KAWAI FREE with a referral code!

✅ No credit card
✅ No email required
✅ Instant access

[Claim Your Bonus]
```

### Social Media Share

```
Join Kawai DeAI Network and get 10 USDT + 200 KAWAI FREE! 
Use my code: ABC123

🤖 Decentralized AI
💰 No credit card
⚡ Instant access

kawai://app?ref=ABC123
```

### In-App Notifications

```
🎉 Referral Bonus Applied!
You'll receive 10 USDT + 200 KAWAI instead of 5 USDT + 100 KAWAI

Your friend ABC123 will also earn 5 USDT + 100 KAWAI
```

---

## Analytics & Tracking

### Key Metrics

1. **Referral Conversion Rate**
   ```
   (Users with referral / Total signups) * 100
   Target: >40%
   ```

2. **Viral Coefficient**
   ```
   Average referrals per user
   Target: >1.0 (viral growth)
   ```

3. **Referrer Engagement**
   ```
   % of users who create referral code
   Target: >60%
   ```

4. **Top Referrers**
   ```
   Leaderboard of users by referral count
   Reward top 10 with bonus KAWAI tokens
   ```

### Events to Track

```typescript
analytics.track('referral_code_created', { user_address });
analytics.track('referral_link_shared', { code, method: 'twitter' });
analytics.track('referral_applied', { code, new_user_address });
analytics.track('referral_bonus_claimed', { code, amount: 8 });
analytics.track('referrer_rewarded', { code, amount: 3 });
```

---

## Testing Checklist

### Backend Tests

- [ ] Generate unique referral codes
- [ ] Prevent duplicate codes
- [ ] Validate code format
- [ ] Prevent self-referral
- [ ] Claim trial with referral (10 USDT)
- [ ] Claim trial without referral (5 USDT)
- [ ] Reward referrer (5 USDT)
- [ ] Update referral stats
- [ ] Handle invalid referral code
- [ ] Prevent double-claiming

### Frontend Tests

- [ ] Detect referral code from URL
- [ ] Display referral banner
- [ ] Apply referral code manually
- [ ] Show bonus upgrade (5→8 USDT)
- [ ] Generate referral code
- [ ] Display referral dashboard
- [ ] Copy referral code
- [ ] Share referral link
- [ ] Track referral stats

---

## Deployment Steps

### 1. Backend

```bash
# Deploy referral system
cd /Users/yuda/github.com/kawai-network/veridium-1

# Test locally
go test ./pkg/store/... -v -run TestReferral
go test ./pkg/gateway/... -v -run TestReferral

# Build
make build

# Deploy
./bin/veridium
```

### 2. Frontend

```bash
cd frontend

# Install dependencies (if needed)
bun install

# Test
bun run test

# Build
bun run build

# Deploy
# (Wails will bundle frontend automatically)
```

### 3. Environment Variables

```bash
# No new env vars needed!
# Uses existing:
# - FREE_TRIAL_AMOUNT_USDT (base: 5.0)
# - CF_KV_USERS_NAMESPACE_ID (for storage)
```

---

## Future Enhancements

### Phase 2 (Month 2)

1. **Tiered Referral Rewards**
   ```
   1-10 referrals: 3 USDT each
   11-50 referrals: 4 USDT each
   51+ referrals: 5 USDT each
   ```

2. **Referral Leaderboard**
   - Top 10 referrers displayed publicly
   - Monthly prizes (KAWAI tokens)
   - Badge system (Bronze, Silver, Gold)

3. **Social Verification**
   - Tweet about Kawai → Extra 2 USDT
   - Join Discord → Extra 1 USDT
   - GitHub star → Extra 1 USDT

4. **Referral Contests**
   ```
   "Refer 10 friends in January → Win 100 KAWAI"
   ```

### Phase 3 (Month 3+)

1. **Multi-Level Referrals**
   - Level 1: 3 USDT (direct referral)
   - Level 2: 1 USDT (referral's referral)

2. **Referral NFTs**
   - Mint NFT for top referrers
   - NFT holders get permanent bonus

3. **API for Partners**
   - Allow partners to create custom referral codes
   - Track conversions per partner
   - Revenue sharing model

---

## Support & FAQ

### Q: How do I get a referral code?

**A:** Click "Wallet" → "Referral" → "Generate Code". You'll get a unique 6-character code.

### Q: Can I use my own referral code?

**A:** No, self-referrals are not allowed to prevent abuse.

### Q: Is there a limit to referrals?

**A:** No! Refer unlimited friends and earn 5 USDT per referral.

### Q: When do I receive my referral reward?

**A:** Instantly! As soon as your friend claims their trial, you get 5 USDT.

### Q: What if my friend doesn't use my code?

**A:** They'll still get 5 USDT, but they miss out on the extra 5 USDT bonus (total 10 USDT). Remind them to use your code!

### Q: Can I change my referral code?

**A:** No, referral codes are permanent. But you can share it as many times as you want.

### Q: How does the 5% mining commission work?

**A:** When your referral mines KAWAI tokens, you automatically earn 5% of their mining rewards as commission. This is lifetime passive income - as long as they keep mining, you keep earning!

### Q: Is there a limit to mining commission?

**A:** No limit! You earn 5% from all mining activity of all your referrals, forever.

### Q: How is the mining reward split?

**A:** For users with referrals:
- **90%** goes to the miner (contributor)
- **5%** goes to the developer (protocol)
- **5%** goes to you (affiliator)

For users without referrals:
- **90%** goes to the miner
- **10%** goes to the developer

### Q: When can I claim my mining commission?

**A:** Mining commissions are accumulated off-chain and settled weekly via Merkle tree distribution. You can claim them from the "Mining Commission" tab in the Rewards Dashboard.

### Q: Can I see how much my referrals are mining?

**A:** Yes! The "Mining Commission" tab shows detailed statistics including:
- Total commission earned
- Number of active mining referrals
- Commission history per referral
- Claimable vs claimed amounts

---

## 💰 Mining Commission Deep Dive

### Overview

The mining commission system creates a **win-win-win** scenario:
1. **Miners** still get 90% of rewards (majority)
2. **Affiliators** earn 5% lifetime passive income
3. **Protocol** gets sustainable user acquisition

### How It Works

```
User mines → Generate reward → Check referral status

IF user has referrer:
  ├─ 90% → Miner (contributor)
  ├─ 5%  → Developer (protocol)
  └─ 5%  → Affiliator (referrer) ✨

ELSE (no referrer):
  ├─ 90% → Miner
  └─ 10% → Developer
```

### Example Scenario

**Month 1:**
- You refer Alice
- Alice mines 1,000 KAWAI
- You earn **50 KAWAI** (5% commission)

**Month 2:**
- Alice mines 2,000 KAWAI
- You earn **100 KAWAI** (5% commission)

**Month 3:**
- Alice mines 1,500 KAWAI
- You earn **75 KAWAI** (5% commission)

**Total: 225 KAWAI passive income** from just one referral!

### Maximizing Your Commission

**1. Quality Over Quantity**
- Focus on referring active users who will mine regularly
- Engaged miners = consistent commission

**2. Support Your Referrals**
- Help them set up mining
- Share tips and best practices
- Active referrals mine more = you earn more

**3. Build a Community**
- Create a Discord/Telegram group for your referrals
- Share mining strategies
- Foster engagement and retention

**4. Track Performance**
- Monitor your commission dashboard
- Identify top-performing referrals
- Learn what makes them successful

### Commission Claiming

**Weekly Settlement:**
1. Every Sunday, referral rewards are settled
2. Merkle tree is generated with all commissions
3. Merkle root is uploaded to blockchain

**Current Status:** ✅ Settlement code is **COMPLETE** (implemented in `pkg/blockchain/referral_settlement.go`)

**Settlement Implementation:**
```go
// docs/REFERRAL_CONTRACT_GUIDE.md (lines 98-126)
// Pseudo-code - needs full implementation

func GenerateMerkleTree() ([]byte, error) {
    // 1. Get all referrers with pending rewards
    referrers := GetAllReferrersWithPendingRewards()
    
    // 2. Create leaves (3-field: period, account, amount)
    var leaves [][]byte
    for _, ref := range referrers {
        leaf := keccak256(period, ref.Address, ref.PendingKawai)
        leaves = append(leaves, leaf)
    }
    
    // 3. Build Merkle tree
    tree := merkletree.New(leaves)
    
    // 4. Store proofs in KV
    for i, ref := range referrers {
        proof := tree.GetProof(i)
        StoreProof(ref.Address, period, proof)
    }
    
    return tree.Root(), nil
}
```

**✅ Settlement Tool Available**

The unified settlement tool supports referral commissions:

```bash
# Generate referral settlement
make settle-referral

# Or settle all reward types at once
make settle-all

# Check status
make reward-settlement-status

# Advanced usage
go run cmd/reward-settlement/main.go generate --type referral
go run cmd/reward-settlement/main.go upload --type referral
go run cmd/reward-settlement/main.go all  # Settle all 3 types at once
```

**Implementation:** `cmd/reward-settlement/main.go` (unified tool for all 3 reward types)

**Claiming Process:**
1. Go to Wallet → Rewards → Referral Commission tab
2. View your claimable commission
3. Click "Claim" button
4. Sign transaction with your wallet
5. Receive KAWAI tokens instantly

**Batch Claiming:**
- Claim multiple weeks at once (supported by contract)
- Saves on gas fees
- More efficient for large commissions

**See Also:**
- Mining settlement: `cmd/mining-settlement/` (reference implementation)
- Cashback settlement: `pkg/blockchain/cashback_settlement.go` (reference code)
- Related: [`REWARD_SYSTEMS.md`](REWARD_SYSTEMS.md) for unified settlement discussion

### Economics

**Why 5%?**
- Sustainable for protocol (developer only sacrifices 5%)
- Attractive for affiliators (lifetime passive income)
- Fair for miners (still get 90%)

**Lifetime Value:**
- Average miner: ~500 KAWAI/month
- Your commission: 25 KAWAI/month per referral
- 10 active referrals = 250 KAWAI/month passive income
- 100 active referrals = 2,500 KAWAI/month passive income

**ROI for Protocol:**
- Cost: 5% of mining rewards
- Benefit: Referred users have 2x higher retention
- Net result: Positive ROI on user acquisition

---

## 📚 Related Documentation

- **Overview:** [`REWARD_SYSTEMS.md`](REWARD_SYSTEMS.md) - Overview & comparison of all reward systems
- **Contract Details:** [`docs/CONTRACTS_OVERVIEW.md`](docs/CONTRACTS_OVERVIEW.md) - All contracts overview
- **Contract Guide:** [`docs/REFERRAL_CONTRACT_GUIDE.md`](docs/REFERRAL_CONTRACT_GUIDE.md) - Detailed referral contract implementation
- **Contract Development:** [`docs/CONTRACTS_WORKFLOW.md`](docs/CONTRACTS_WORKFLOW.md) - How to develop & deploy contracts
- **MINTER_ROLE:** [`MINTER_ROLE_REQUIREMENTS.md`](MINTER_ROLE_REQUIREMENTS.md) - Why MINTER_ROLE is needed
- **Backend Store:** [`pkg/store/README.md`](pkg/store/README.md) - KV storage implementation

---

## Contact

For questions or issues:
- Discord: https://discord.gg/SNf3ZEa8Eq
- Twitter: @kawai_network
- Email: support@kawai.network

---

**Built with ❤️ by the Kawai Team**

