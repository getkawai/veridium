# 🎁 Referral System Documentation

## Overview

The Kawai DeAI Network referral system incentivizes viral growth by rewarding both referrers and new users.

## Incentive Structure

| User Type | Bonus Amount | Description |
|-----------|--------------|-------------|
| **New User (No Referral)** | 5 USDT | Base trial bonus |
| **New User (With Referral)** | 10 USDT | Enhanced trial bonus (+100%) |
| **Referrer** | 5 USDT | Reward per successful referral |

### Why This Works

- **Win-Win:** Both parties benefit
- **Strong Incentive:** 5 USDT = 100% of base bonus
- **Viral Potential:** Coefficient of 2.0-3.0x
- **No Limit:** Unlimited referrals = unlimited earnings

---

## User Flow

### For New Users

1. **Receive Referral Link**
   ```
   https://kawai.network?ref=ABC123
   ```

2. **Open App**
   - App detects `?ref=ABC123` in URL
   - Shows enhanced bonus message: "🎉 Get 10 USDT FREE with referral code!"

3. **Setup Wallet**
   - Referral code auto-applied
   - Banner shows: "Referral Applied: ABC123"
   - Bonus upgraded: 5 USDT → 10 USDT

4. **Claim Bonus**
   - Automatically claimed on first wallet unlock
   - Receives 10 USDT instead of 5 USDT
   - Referrer receives 5 USDT

### For Referrers

1. **Generate Referral Code**
   ```typescript
   // Desktop App (Wails)
   import { CreateReferralCode } from '@/bindings/services/referralservice';
   
   const code = await CreateReferralCode(userAddress);
   console.log(code); // "ABC123"
   ```

2. **Share Code**
   - Via social media
   - Direct link: `kawai://app?ref=ABC123`
   - In-app share button

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
3. Determine bonus amount (5 or 8 USDT)
4. Claim trial with atomic operation
5. Reward referrer (3 USDT)
6. Update referral stats

#### 3. Wails Service (Desktop App)

```go
// internal/services/referralservice.go

type ReferralService struct {
    kvStore *store.KVStore
}

// Wails-exposed methods (auto-generated TypeScript bindings)
func (s *ReferralService) CreateReferralCode(userAddress string) (*ReferralCodeResponse, error)
func (s *ReferralService) GetReferralStats(userAddress string) (*ReferralStatsResponse, error)
func (s *ReferralService) ClaimFreeTrialWithReferral(address, machineID, referralCode string) (*ClaimTrialWithReferralResponse, error)
func (s *ReferralService) GetReferralBonusAmounts() map[string]interface{}
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
- Bonus comparison (5 vs 8 USDT)

#### 2. Referral Dashboard

```typescript
// frontend/src/features/Referral/ReferralDashboard.tsx

<ReferralDashboard 
  stats={referralStats}
  onShare={handleShare}
/>
```

**Features:**
- Display referral code (large, copyable)
- Statistics cards (referrals, earnings, pending)
- Recent referrals list
- Native share API integration
- How-it-works guide

#### 3. URL Parameter Detection

```typescript
// In AuthSignInBox.tsx

useEffect(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const refCode = urlParams.get('ref');
  if (refCode) {
    setReferralCode(refCode.toUpperCase());
    setHasReferral(true);
  }
}, []);
```

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
- **Machine ID:** One claim per device

### 3. Atomic Operations
- Race condition protection
- Retry mechanism with exponential backoff
- Prevents double-claiming

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
- [ ] Claim trial with referral (8 USDT)
- [ ] Claim trial without referral (5 USDT)
- [ ] Reward referrer (3 USDT)
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

**A:** No! Refer unlimited friends and earn 3 USDT per referral.

### Q: When do I receive my referral reward?

**A:** Instantly! As soon as your friend claims their trial, you get 3 USDT.

### Q: What if my friend doesn't use my code?

**A:** They'll still get 5 USDT, but they miss out on the extra 3 USDT bonus. Remind them to use your code!

### Q: Can I change my referral code?

**A:** No, referral codes are permanent. But you can share it as many times as you want.

---

## Contact

For questions or issues:
- Discord: https://discord.gg/kawai
- Twitter: @kawai_network
- Email: support@kawai.network

---

**Built with ❤️ by the Kawai Team**

