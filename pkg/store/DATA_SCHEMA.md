# Cloudflare KV Data Schema

Dokumentasi ini menjelaskan schema data untuk setiap namespace yang digunakan dalam `pkg/store`.

---

## 1. Namespace: `CONTRIBUTORS`
**Deskripsi**: Menyimpan profil contributor (miner) dan log reward pekerjaan (mining/inference).
**Implementasi**: `pkg/store/contributor.go`, `pkg/store/job_rewards.go`

### A. Contributor Profile
*   **Key Format**: `{wallet_address}` (lowercase, e.g., `0x123...abc`)
*   **Value Format**: JSON (`ContributorData`)

```json
{
  "wallet_address": "0x123...abc",
  "endpoint_url": "https://...",
  "hardware_specs": "RTX 3090...",
  "registered_at": "2024-01-01T00:00:00Z",
  "last_seen": "2024-01-02T12:00:00Z",
  "status": "online",             // "online", "offline", "deleted", "admin"
  "accumulated_rewards": "1000",  // KAWAI (wei) - Phase 1
  "accumulated_usdt": "500000",   // USDT (micro) - Phase 2
  "is_active": true,
  "deleted_at": "0001-01-01T00:00:00Z",
  "is_admin": false,
  "version": 1,

  // Discovery & Load Balancing Metadata
  "region": "us-west",
  "available_models": ["llama-2-70b", "mistral-7b"],
  "active_requests": 5,
  "total_requests": 1500,
  "avg_response_time": 1.5,
  "success_rate": 0.98,
  "last_health_check": "2024-01-02T12:00:00Z",

  // Parsed Hardware Fields
  "cpu_cores": 16,
  "total_ram": 64,
  "available_ram": 48,
  "gpu_model": "RTX 3090",
  "gpu_memory": 24
}
```

### B. Job Rewards (Log Mining)
*   **Key Format**: `job_rewards:{contributor_address}:{timestamp_unix}`
*   **Value Format**: JSON (`JobRewardRecord`)

```json
{
  "timestamp": "2024-01-02T12:00:00Z",
  "contributor_address": "0x123...",
  "user_address": "0x456...",
  "referrer_address": "0x789...",
  "developer_address": "0xdev...",
  "contributor_amount": "850",
  "developer_amount": "50",
  "user_amount": "50",
  "affiliator_amount": "50",
  "token_usage": 1500,
  "reward_type": "mining",       // "mining", "cashback", "referral", "revenue"
  "has_referrer": true,
  "is_settled": false,           // true jika sudah masuk settlement period
  "settled_period_id": 0
}
```

---

## 2. Namespace: `PROOFS`
**Deskripsi**: Menyimpan Merkle Proofs untuk reward yang sudah disettle (mingguan).
**Implementasi**: `pkg/store/merkle.go`

*   **Key Format**: `{address}:{period_id}` (e.g., `0x123...:1704067200000000000` - period_id dalam nanoseconds)
*   **Value Format**: JSON (`MerkleProofData`)

```json
{
  "index": 1,
  "amount": "1000000000000000000", // Total claimable amount (wei)
  "proof": ["0xabc...", "0xdef..."],
  "merkle_root": "0xroot...",
  "period_id": 1704067200000000000, // Unix timestamp nanoseconds
  "created_at": "2024-01-08T00:00:00Z",
  "reward_type": "mining",         // "mining", "cashback", "referral", "revenue"
  "claim_status": "unclaimed",     // "unclaimed", "pending", "confirmed", "failed"
  "claim_tx_hash": "0xtx...",
  "claim_attempts": 0,
  "claimed_at": "2024-01-09T00:00:00Z",
  "address": "0x123...",
  
  // Breakdown fields (for verification)
  "contributor_amount": "850...",
  "developer_amount": "50...",
  "user_amount": "50...",
  "affiliator_amount": "50...",
  "developer_address": "0x...",
  "user_address": "0x...",
  "affiliator_address": "0x..."
}
```

---

## 3. Namespace: `SETTLEMENTS`
**Deskripsi**: Metadata periode settlement mingguan.
**Implementasi**: `pkg/store/settlement.go`

*   **Key Format**: `{period_id}` (Unix timestamp nanoseconds)
*   **Value Format**: JSON (`SettlementPeriod`)

```json
{
  "period_id": 1704067200000000000,
  "merkle_root": "0xroot...",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-08T00:00:00Z",
  "total_amount": "5000000...",
  "reward_type": "mining",
  "status": "completed", // "pending", "proofs_saved", "balances_reset", "completed", "failed"
  "contributor_count": 150,
  "proofs_saved": 150,
  "balances_reset": 150,
  "started_at": "2024-01-08T00:00:00Z",
  "completed_at": "2024-01-08T00:05:00Z",
  "error": ""
}
```

---

## 4. Namespace: `USERS`
**Deskripsi**: Menyimpan balance user, status trial, dan referral.
**Implementasi**: `pkg/store/balance.go`, `pkg/store/balance_trial.go`

*   **Key Format**: `balance:{address}`
*   **Value Format**: JSON (`UserBalance`)

```json
{
  "address": "0xuser...",
  "usdt_balance": "1000000",     // Micro USDT (6 decimals)
  "kawai_balance": "0",          // Wei (18 decimals)
  "trial_claimed": true,
  "referrer_address": "0xref..." // Address referrer (optional)
}
```

---

## 5. Namespace: `APIKEY`
**Deskripsi**: Mapping API Key ke Wallet Address (Forward Index).
**Implementasi**: `pkg/store/apikey.go`

*   **Key Format**: `{api_key}` (e.g., `vk-a1b2c3...`)
*   **Value Format**: String (Raw Wallet Address)

```text
0x1234567890abcdef1234567890abcdef12345678
```

---

## 6. Namespace: `AUTHZ`
**Deskripsi**: Mapping Wallet Address ke API Key (Reverse Index).
**Implementasi**: `pkg/store/apikey.go`

*   **Key Format**: `{wallet_address}`
*   **Value Format**: String (Raw API Key)

```text
vk-a1b2c3d4e5f6...
```

---

## 7. Namespace: `CASHBACK`
**Deskripsi**: Menyimpan data cashback deposit user.
**Implementasi**: `pkg/store/cashback.go`

### A. Cashback Record
*   **Key Format**: `cashback:{user_address}:{tx_hash}`
*   **Value Format**: JSON (`CashbackRecord`)

```json
{
  "user_address": "0xuser...",
  "deposit_tx_hash": "0xtx...",
  "deposit_amount": "1000000",   // USDT (6 decimals)
  "cashback_amount": "500...",   // KAWAI (18 decimals)
  "rate": 200,                   // Basis points (2%)
  "tier": 5,                     // Tier 1-5
  "is_first_time": true,
  "created_at": "2024-01-01T...",
  "period": 1,                   // Settlement period sequence
  "claimed": false,
  "proof": ["..."],
  "merkle_root": "..."
}
```

### B. Cashback Stats
*   **Key Format**: `cashback_stats:{user_address}`
*   **Value Format**: JSON (`CashbackStats`)

```json
{
  "total_cashback": "5000...",
  "pending_cashback": "1000...",
  "claimed_cashback": "4000...",
  "total_deposits": 5,
  "total_deposit_amount": "5000000",
  "first_deposit_at": "...",
  "last_deposit_at": "..."
}
```

### C. Settlement Metadata (Period Tracking)
*   **Key**: `cashback_period:{id}:users` -> JSON List `["0x...", "0x..."]`
*   **Key**: `cashback_period:{id}:merkle_root` -> JSON String `"0xroot..."`
*   **Key**: `cashback_proof:{id}:{address}` -> JSON Struct (Proof details)

---

## 8. Namespace: `HOLDER`
**Deskripsi**: Registry pemegang token/point.
**Implementasi**: `pkg/store/holder.go`

*   **Key Format**: `holder:{address}`
*   **Value Format**: JSON (`HolderInfo`)

```json
{
  "address": "0x...",
  "lastSeen": 1704153600,
  "source": "cli",  // "desktop", "cli", "transfer"
  "registered": 1704000000
}
```

---

## 9. Namespace: `P2PMARKETPLACE`
**Deskripsi**: Penyimpanan data marketplace P2P.
**Implementasi**: `pkg/store/marketplace.go`

*   **Key Format**: Arbitrary Strings (defined by application logic)
*   **Value Format**: Raw Bytes / JSON String
*   **Feature**: Mendukung TTL (Time-To-Live) secara native via Cloudflare KV expiration.
