# Package Store

Off-chain storage layer untuk Kawai Network menggunakan Cloudflare Workers KV.

## Overview

Package ini menyediakan:
- **Contributor Management**: Data user, balances, heartbeat
- **Merkle Proofs**: Period-specific proofs untuk reward claims
- **Settlement Automation**: Weekly settlement dengan rollback support
- **Claim Tracking**: Status tracking untuk prevent token loss

---

## Architecture

### Multi-Namespace Design

```
┌─────────────────────────────────────────────────────────────────┐
│                     Cloudflare Workers KV                        │
├──────────────────────┬──────────────────────┬───────────────────┤
│  kawai-contributors  │    kawai-proofs      │ kawai-settlements │
│                      │                      │                   │
│  Key: {address}      │  Key: {addr}:{pid}   │  Key: {periodID}  │
├──────────────────────┼──────────────────────┼───────────────────┤
│  • User profiles     │  • Weekly proofs     │  • Period metadata│
│  • Balances          │  • Claim status      │  • Settlement log │
│  • Heartbeat         │  • Proof data        │  • Total amounts  │
└──────────────────────┴──────────────────────┴───────────────────┘
```

### Why Multiple Namespaces?

| Benefit | Description |
|---------|-------------|
| **Isolation** | Each data type completely isolated |
| **No Collision** | Simple keys without prefixes |
| **Independent Rate Limits** | 3x throughput capacity |
| **Easier Maintenance** | Can clear proofs without affecting users |
| **Faster Scans** | List operations only scan relevant data |

---

## Quick Start

### Initialize Store

```go
import (
    "github.com/kawai-network/veridium/internal/constant"
    "github.com/kawai-network/veridium/pkg/store"
)

// Create multi-namespace store
kv, err := store.NewMultiNamespaceKVStore(store.MultiNamespaceConfig{
    APIToken:                constant.GetCfApiToken(),
    AccountID:               constant.GetCfAccountId(),
    ContributorsNamespaceID: constant.GetCfKvContributorsNamespaceId(),
    ProofsNamespaceID:       constant.GetCfKvProofsNamespaceId(),
    SettlementsNamespaceID:  constant.GetCfKvSettlementsNamespaceId(),
})
```

### Weekly Settlement Flow

```go
// 1. Generate unique period ID
periodID := store.GenerateUniquePeriodID()

// 2. Get snapshots (sorted for consistent Merkle tree)
snapshots, _ := kv.GetSettlementSnapshots(ctx, "kawai")

// 3. Generate Merkle tree (external library)
merkleRoot, proofs := GenerateMerkleTree(snapshots)

// 4. Perform settlement (atomic with rollback)
period, err := kv.PerformSettlement(ctx, periodID, merkleRoot, "kawai", proofs)
// Automatically: saves proofs → resets balances → saves metadata
```

### Claim Flow (Safe)

```go
// 1. Get claimable rewards
claimable, _ := kv.GetClaimableRewards(ctx, address)

// 2. Mark pending BEFORE sending TX
kv.MarkClaimPending(ctx, address, periodID, txHash)

// 3. Submit to smart contract
contract.Claim(proof.Proof, proof.Amount, proof.MerkleRoot)

// 4. Confirm or fail based on TX result
if txSuccess {
    kv.ConfirmClaim(ctx, address, periodID)
} else {
    kv.MarkClaimFailed(ctx, address, periodID, "reason")
}
```

---

## Data Structures

### ContributorData

```go
type ContributorData struct {
    WalletAddress      string    `json:"wallet_address"`
    EndpointURL        string    `json:"endpoint_url"`
    HardwareSpecs      string    `json:"hardware_specs"`
    RegisteredAt       time.Time `json:"registered_at"`
    LastSeen           time.Time `json:"last_seen"`
    Status             string    `json:"status"`              // online/offline/deleted
    AccumulatedRewards string    `json:"accumulated_rewards"` // KAWAI (wei)
    AccumulatedUSDT    string    `json:"accumulated_usdt"`    // USDT (micro)
    IsActive           bool      `json:"is_active"`           // Soft delete flag
    DeletedAt          time.Time `json:"deleted_at,omitempty"`
    IsAdmin            bool      `json:"is_admin,omitempty"`
}
```

**Namespace:** `kawai-contributors`  
**Key:** `{address}` (lowercase)  
**Example:** `0x742d35cc6634c0532925a3b844bc454e4438f44e`

### MerkleProofData

```go
type MerkleProofData struct {
    Index         uint64      `json:"index"`          // Leaf index
    Amount        string      `json:"amount"`         // Claimable amount (wei)
    Proof         []string    `json:"proof"`          // Merkle proof hashes
    MerkleRoot    string      `json:"merkle_root"`    // Root for this period
    PeriodID      int64       `json:"period_id"`      // Settlement period
    CreatedAt     time.Time   `json:"created_at"`
    RewardType    string      `json:"reward_type"`    // "kawai" or "usdt"
    ClaimStatus   ClaimStatus `json:"claim_status"`   // unclaimed/pending/confirmed/failed
    ClaimTxHash   string      `json:"claim_tx_hash"`
    ClaimAttempts int         `json:"claim_attempts"`
    ClaimedAt     time.Time   `json:"claimed_at,omitempty"`
    Address       string      `json:"address"`        // Contributor address
}

type ClaimStatus string
const (
    ClaimStatusUnclaimed ClaimStatus = "unclaimed"
    ClaimStatusPending   ClaimStatus = "pending"
    ClaimStatusConfirmed ClaimStatus = "confirmed"
    ClaimStatusFailed    ClaimStatus = "failed"
)
```

**Namespace:** `kawai-proofs`  
**Key:** `{address}:{periodID}`  
**Example:** `0x742d35cc...:1704067200000000000`

### SettlementPeriod

```go
type SettlementPeriod struct {
    PeriodID         int64            `json:"period_id"`
    MerkleRoot       string           `json:"merkle_root"`
    StartDate        time.Time        `json:"start_date"`
    EndDate          time.Time        `json:"end_date"`
    TotalAmount      string           `json:"total_amount"`
    RewardType       string           `json:"reward_type"`
    Status           SettlementStatus `json:"status"`
    ContributorCount int              `json:"contributor_count"`
    ProofsSaved      int              `json:"proofs_saved"`
    BalancesReset    int              `json:"balances_reset"`
    StartedAt        time.Time        `json:"started_at"`
    CompletedAt      time.Time        `json:"completed_at,omitempty"`
    Error            string           `json:"error,omitempty"`
}

type SettlementStatus string
const (
    SettlementStatusPending       SettlementStatus = "pending"
    SettlementStatusProofsSaved   SettlementStatus = "proofs_saved"
    SettlementStatusBalancesReset SettlementStatus = "balances_reset"
    SettlementStatusCompleted     SettlementStatus = "completed"
    SettlementStatusFailed        SettlementStatus = "failed"
)
```

**Namespace:** `kawai-settlements`  
**Key:** `{periodID}`  
**Example:** `1704067200000000000`

---

## API Reference

### Contributor Operations

| Method | Description |
|--------|-------------|
| `RegisterContributor(ctx, address, endpoint, specs)` | Register new contributor |
| `GetContributor(ctx, address)` | Get contributor data |
| `ListContributors(ctx)` | List all contributors |
| `ListActiveContributors(ctx)` | List active only |
| `ListContributorsWithBalance(ctx, rewardType)` | List with balance > 0 |
| `UpdateHeartbeat(ctx, address)` | Update LastSeen |
| `SoftDeleteContributor(ctx, address)` | Soft delete (preserve data) |
| `RestoreContributor(ctx, address)` | Restore deleted |
| `ResetAccumulatedRewards(ctx, address, type)` | Reset balance to 0 |

### Settlement Operations

| Method | Description |
|--------|-------------|
| `GenerateUniquePeriodID()` | Generate unique period ID |
| `GetSettlementSnapshots(ctx, rewardType)` | Get sorted balances |
| `PerformSettlement(ctx, periodID, root, type, proofs)` | Full settlement |
| `PerformSettlementWithConfig(ctx, ..., config)` | With custom config |
| `PerformSettlementParallel(ctx, ..., workers)` | Parallel processing |
| `ResumeSettlement(ctx, periodID, proofs, config)` | Resume failed |
| `SaveSettlementPeriod(ctx, period)` | Save metadata |
| `GetSettlementPeriod(ctx, periodID)` | Get metadata |
| `ListSettlementPeriods(ctx)` | List all periods |

### Proof Operations

| Method | Description |
|--------|-------------|
| `SaveMerkleProofForPeriod(ctx, addr, periodID, proof)` | Save proof |
| `GetMerkleProofForPeriod(ctx, addr, periodID)` | Get proof |
| `ListMerkleProofs(ctx, addr)` | List all proofs for address |
| `DeleteMerkleProof(ctx, addr, periodID)` | Delete proof |
| `CleanupOldProofs(ctx, duration)` | Cleanup old proofs |

### Claim Operations

| Method | Description |
|--------|-------------|
| `GetClaimableRewards(ctx, addr)` | Get all claimable |
| `MarkClaimPending(ctx, addr, periodID, txHash)` | Mark as pending |
| `ConfirmClaim(ctx, addr, periodID)` | Confirm success |
| `MarkClaimFailed(ctx, addr, periodID, reason)` | Mark as failed |
| `RetryFailedClaim(ctx, addr, periodID)` | Reset for retry |
| `GetPendingClaims(ctx)` | List all pending |

### Admin Operations

| Method | Description |
|--------|-------------|
| `EnsureAdminExists(ctx, adminAddress)` | Auto-register admin |
| `RecordJobReward(ctx, contributor, tokens, admin, mode)` | Record reward (70/30 split) |

---

## Configuration

### SettlementConfig

```go
type SettlementConfig struct {
    BatchSize      int           // Contributors per batch (default: 50)
    BatchDelay     time.Duration // Delay between batches (default: 100ms)
    MaxRetries     int           // Retries per operation (default: 3)
    EnableRollback bool          // Rollback on failure (default: true)
}

// Default
config := store.DefaultSettlementConfig()

// Custom for large scale
config := &store.SettlementConfig{
    BatchSize:      100,
    BatchDelay:     200 * time.Millisecond,
    MaxRetries:     5,
    EnableRollback: true,
}
```

---

## Safety Features

### 1. Claim Status Tracking
Prevents token loss from TX reverts.

```go
// BEFORE TX: Mark pending (don't delete proof)
kv.MarkClaimPending(ctx, addr, periodID, txHash)

// AFTER TX success: Confirm
kv.ConfirmClaim(ctx, addr, periodID)

// AFTER TX fail: Mark failed (can retry)
kv.MarkClaimFailed(ctx, addr, periodID, reason)
```

### 2. Settlement Rollback
Atomic settlement - all or nothing.

```go
config := &SettlementConfig{EnableRollback: true}
// If any proof fails to save → all saved proofs deleted
```

### 3. Consistent Merkle Ordering
Snapshots sorted by lowercase address.

```go
snapshots := kv.GetSettlementSnapshots(ctx, "kawai")
// Returns sorted by strings.ToLower(Address)
```

### 4. Unique Period ID
Prevents collision from multiple servers.

```go
periodID := store.GenerateUniquePeriodID()
// Uses UnixNano + random suffix
```

### 5. Soft Delete
Preserves data for unregistered contributors.

```go
kv.SoftDeleteContributor(ctx, addr)
// IsActive=false, but data preserved for settlement
```

### 6. Resumable Settlement
Continue from last successful step.

```go
kv.ResumeSettlement(ctx, periodID, proofs, config)
// Checks SettlementStatus and continues appropriately
```

---

## Edge Cases Handled

| Issue | Solution | Status |
|-------|----------|--------|
| Double claim | On-chain mapping + ClaimStatus | ✅ |
| TX revert after off-chain update | ClaimStatus tracking | ✅ |
| Settlement partial failure | Rollback mechanism | ✅ |
| Merkle tree ordering mismatch | Sorted snapshots | ✅ |
| Large scale (1000+ contributors) | Batch + parallel processing | ✅ |
| Contributor unregister | Soft delete | ✅ |
| Period ID collision | Unique ID generation | ✅ |
| Settlement interruption | Resumable settlement | ✅ |
| Admin not registered | Auto-register | ✅ |
| Zero balance contributors | Skip in snapshots | ✅ |

---

## Performance

### Cloudflare KV Characteristics

| Metric | Value |
|--------|-------|
| Read Latency | ~50ms (global edge) |
| Write Latency | ~200ms |
| Global Consistency | ~60 seconds |
| Max Key Size | 512 bytes |
| Max Value Size | 25 MB |

### With 3 Namespaces

| Metric | Per Namespace | Total |
|--------|---------------|-------|
| Write/sec | 1,000 | 3,000 |
| Read/sec | 100,000 | 300,000 |

### Optimization Tips

1. Use `PerformSettlementParallel()` for 1000+ contributors
2. Set appropriate `BatchSize` (50-100)
3. Run `CleanupOldProofs()` monthly
4. Cache `GetOnlineContributors()` for 30s

---

## Files

```
pkg/store/
├── contributor.go      # ContributorData, KVStore, multi-namespace
├── merkle.go           # MerkleProofData, ClaimStatus, SettlementPeriod
├── settlement.go       # Settlement automation, claim tracking
├── keys.go             # Key generation functions
├── example_settlement.go # Usage examples
└── README.md           # This file
```

---

## Examples

See `example_settlement.go` for complete examples:
- `ExampleWeeklySettlement()` - Basic settlement
- `ExampleSettlementWithConfig()` - Custom config
- `ExampleParallelSettlement()` - Large scale
- `ExampleClaimFlow()` - Safe claim process
- `ExampleSoftDeleteContributor()` - Soft delete
- `ExampleResumeSettlement()` - Resume failed

---

## License

Copyright © 2025 Kawai Network
