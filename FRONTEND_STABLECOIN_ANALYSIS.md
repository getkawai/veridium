# Frontend Stablecoin Migration Analysis

## 🔍 Current State

Setelah memeriksa `frontend/src`, ditemukan **referensi USDT yang perlu dipertimbangkan** untuk update.

---

## 📊 Findings

### ✅ SUDAH BENAR (Tidak Perlu Diubah)

#### 1. **Network Configuration** (`frontend/src/config/network.ts`)
```typescript
{
  address: config.contracts.usdt,  // ✅ Menggunakan backend config
  name: 'Tether USD',
  symbol: 'USDT',
  decimals: 6,
}
```
**Status:** ✅ **CORRECT**
- Menggunakan `config.contracts.usdt` dari backend
- Backend sudah otomatis switch antara MockUSDT (testnet) dan USDC (mainnet)
- Tidak perlu diubah karena address sudah dinamis

#### 2. **Reward Type Filtering** (`RevenueShareSection.tsx`)
```typescript
.filter((p): p is ClaimableReward => p !== null && p.reward_type === 'usdt');
```
**Status:** ✅ **CORRECT** (dengan catatan)
- Backend sudah support both `"usdt"` dan `"stablecoin"` reward types
- Filter ini akan tetap bekerja untuk backward compatibility
- **OPTIONAL:** Bisa ditambahkan `|| p.reward_type === 'stablecoin'` untuk future-proofing

---

## 🟡 PERLU DIPERTIMBANGKAN (User-Facing Text)

### Area yang Menampilkan "USDT" ke User

#### 1. **Token Display Name** (`HomeContent.tsx:827`)
```typescript
<div style={{ fontWeight: 600 }}>USDT</div>
<div style={{ fontSize: 12, color: theme.colorTextSecondary }}>
  Tether USD
</div>
```

**Pertanyaan:**
- Apakah tetap tampilkan "USDT" atau ganti ke "Stablecoin"?
- Apakah perlu dynamic label berdasarkan network?

**Opsi:**
1. **Keep as "USDT"** - Lebih familiar untuk user
2. **Dynamic Label** - "MockUSDT" (testnet) vs "USDC" (mainnet)
3. **Generic "Stablecoin"** - Lebih akurat tapi kurang familiar

#### 2. **Referral Bonus Messages**
```typescript
// ReferralBanner.tsx:62
message.success('Referral code applied! You\'ll get 10 USDT + 200 KAWAI bonus 🎉');

// AuthSignInBox.tsx:335
? '🎉 Get 10 USDT + 200 KAWAI FREE!' 
: 'Get 5 USDT + 100 KAWAI FREE'
```

**Pertanyaan:**
- Apakah bonus masih dalam "USDT" atau berubah ke "USDC" di mainnet?
- Apakah perlu dynamic text?

#### 3. **Revenue Share UI**
```typescript
// RevenueShareSection.tsx
suffix="USDT"
"Earn 100% of platform profit (USDT) proportional to your KAWAI holdings"
"No claimable revenue yet. Keep holding KAWAI to earn weekly USDT!"
```

**Pertanyaan:**
- Apakah tetap "USDT" atau ganti ke "Stablecoin"?
- Apakah perlu clarification untuk mainnet users?

---

## 💡 Rekomendasi

### Opsi A: **MINIMAL CHANGES** (Recommended)
**Filosofi:** "USDT" sebagai generic term untuk stablecoin

**Alasan:**
- User lebih familiar dengan "USDT"
- Secara teknis, backend sudah handle switching
- Tidak perlu confuse user dengan "MockUSDT" vs "USDC"
- Lebih simple dan clean

**Changes Needed:**
1. ✅ **NONE** - Keep semua text as "USDT"
2. ✅ Backend sudah handle address switching
3. ✅ Reward type filtering sudah backward compatible

**Pros:**
- Zero frontend changes needed
- User experience consistent
- Familiar terminology

**Cons:**
- Technically inaccurate on mainnet (it's USDC, not USDT)
- Might confuse advanced users

---

### Opsi B: **DYNAMIC LABELS** (More Accurate)
**Filosofi:** Show actual token name based on network

**Changes Needed:**

#### 1. Update Token Display (`HomeContent.tsx`)
```typescript
// Before
<div style={{ fontWeight: 600 }}>USDT</div>
<div style={{ fontSize: 12, color: theme.colorTextSecondary }}>
  Tether USD
</div>

// After
<div style={{ fontWeight: 600 }}>
  {isTestnet ? 'MockUSDT' : 'USDC'}
</div>
<div style={{ fontSize: 12, color: theme.colorTextSecondary }}>
  {isTestnet ? 'Mock Tether USD (Testnet)' : 'USD Coin'}
</div>
```

#### 2. Update Referral Messages
```typescript
const stablecoinName = isTestnet ? 'USDT' : 'USDC';
message.success(`Referral code applied! You'll get 10 ${stablecoinName} + 200 KAWAI bonus 🎉`);
```

#### 3. Update Revenue Share Text
```typescript
const stablecoinName = isTestnet ? 'USDT' : 'USDC';
suffix={stablecoinName}
```

**Pros:**
- Technically accurate
- Transparent to users
- Educational

**Cons:**
- More code changes
- Might confuse casual users
- Need to pass network context everywhere

---

### Opsi C: **GENERIC "STABLECOIN"** (Most Generic)
**Filosofi:** Use generic term to avoid confusion

**Changes Needed:**
- Replace all "USDT" with "Stablecoin" in UI
- Add tooltip explaining it's MockUSDT (testnet) or USDC (mainnet)

**Pros:**
- Accurate and generic
- Future-proof if we add more stablecoins

**Cons:**
- Less familiar to users
- "Stablecoin" sounds technical
- More confusing than "USDT"

---

## 🎯 Final Recommendation

### **PILIHAN: Opsi A (MINIMAL CHANGES)**

**Reasoning:**
1. ✅ Backend sudah handle semua logic switching
2. ✅ "USDT" adalah term yang paling familiar untuk users
3. ✅ Tidak perlu confuse users dengan technical details
4. ✅ Zero frontend changes = less risk
5. ✅ Reward type filtering sudah backward compatible

**What to Do:**
1. ✅ **NOTHING** - Keep frontend as-is
2. ✅ Backend sudah support both `"usdt"` dan `"stablecoin"` reward types
3. ✅ Address switching sudah otomatis via backend config
4. ✅ User experience tetap consistent

**Optional Enhancement (Future):**
- Add tooltip on token display: "USDT (MockUSDT on testnet, USDC on mainnet)"
- Add info icon with explanation
- Add to FAQ/Help section

---

## 🔧 If You Choose Opsi B (Dynamic Labels)

### Files to Update:

1. **`frontend/src/config/network.ts`**
   - Add helper function `getStablecoinName(config: BackendConfig): string`

2. **`frontend/src/app/wallet/HomeContent.tsx`**
   - Update token display to use dynamic name
   - Update balance calculations (already correct)

3. **`frontend/src/app/wallet/components/rewards/RevenueShareSection.tsx`**
   - Update all "USDT" text to use dynamic name
   - Update filter to include `'stablecoin'` reward type

4. **`frontend/src/features/Referral/ReferralBanner.tsx`**
   - Update bonus messages to use dynamic name

5. **`frontend/src/app/wallet/AuthSignInBox.tsx`**
   - Update free trial messages to use dynamic name

### Implementation Example:

```typescript
// frontend/src/config/network.ts
export function getStablecoinName(config: BackendConfig): string {
  return config.network.isTestnet ? 'MockUSDT' : 'USDC';
}

export function getStablecoinDisplayName(config: BackendConfig): string {
  return config.network.isTestnet 
    ? 'Mock Tether USD (Testnet)' 
    : 'USD Coin';
}

// Usage in components
const config = await getBackendNetworkConfig();
const stablecoinName = getStablecoinName(config);
```

---

## 📝 Summary

### Current Status: ✅ **FRONTEND ALREADY WORKS**

**Why?**
- Backend config (`config.contracts.usdt`) sudah dinamis
- Address switching otomatis (MockUSDT testnet, USDC mainnet)
- Reward type filtering backward compatible
- No breaking changes

### Decision Needed:

**Question:** Apakah user-facing text perlu diubah dari "USDT" ke nama yang lebih akurat?

**Options:**
- **A. Keep "USDT"** - Simple, familiar, no changes needed ✅ **RECOMMENDED**
- **B. Dynamic Labels** - Accurate, transparent, requires changes
- **C. Generic "Stablecoin"** - Generic, future-proof, less familiar

**My Recommendation:** **Opsi A** - Keep as "USDT" for simplicity and familiarity.

---

## 🚀 Action Items

### If Choosing Opsi A (Recommended):
- [ ] ✅ **NO ACTION NEEDED** - Frontend already correct
- [ ] Optional: Add tooltip/FAQ explaining testnet vs mainnet tokens

### If Choosing Opsi B (Dynamic Labels):
- [ ] Update `network.ts` with helper functions
- [ ] Update `HomeContent.tsx` token display
- [ ] Update `RevenueShareSection.tsx` text and filter
- [ ] Update `ReferralBanner.tsx` messages
- [ ] Update `AuthSignInBox.tsx` messages
- [ ] Test on both testnet and mainnet
- [ ] Update user documentation

### If Choosing Opsi C (Generic):
- [ ] Replace all "USDT" with "Stablecoin"
- [ ] Add tooltips explaining actual token
- [ ] Update documentation
- [ ] Test user comprehension

---

**Conclusion:** Frontend sudah bekerja dengan baik karena menggunakan backend config. Perubahan hanya diperlukan jika ingin update user-facing text untuk akurasi atau transparansi.
