# Contributor Migration: cmd/contributor → cmd/server

## Status: ✅ COMPLETE (Phase 1)

Migration dari `cmd/contributor` ke `cmd/server` telah selesai. Semua fitur contributor telah diintegrasikan ke Kronk model server.

---

## Summary

**Before:**
- `cmd/contributor` - Standalone CLI dengan custom gateway
- `cmd/server` - Pure Kronk model server

**After:**
- `cmd/server` - Kronk model server **+ contributor features**
- `cmd/contributor` - **DEPRECATED**

---

## What Was Migrated

### ✅ Core Infrastructure
- **KV Store** - Cloudflare KV integration
- **Wallet Service** - Create/unlock/switch wallets
- **Blockchain Client** - Halving logic & token interactions
- **Holder Registry** - KAWAI token holder tracking

### ✅ Contributor Features
- **Hardware Detection** - Auto-detect CPU/GPU/RAM specs
- **Contributor Registry** - Register to network via KV
- **Heartbeat System** - 30s keep-alive mechanism
- **Tunnel Integration** - Public endpoint via tunnelkit

### ✅ AI Services
- **Whisper Service** - Speech-to-text capability
- **Stable Diffusion** - Image generation capability
- **LLM Service** - Already in Kronk (improved)

---

## Configuration Changes

### Old (cmd/contributor)
```bash
./contributor \
  --password mypass \
  --wallet 0x... \
  --import-mnemonic "word1 word2..."
```

### New (cmd/server)
```bash
# Environment variables
export KRONK_CONTRIBUTOR_ENABLED=true
export KRONK_CONTRIBUTOR_WALLETPASSWORD=mypass
export KRONK_CONTRIBUTOR_WALLETADDRESS=0x...
export KRONK_CONTRIBUTOR_IMPORTMNEMONIC="word1 word2..."

./server
```

---

## Code Changes

### Files Modified
- `cmd/server/api/services/kronk/kronk.go` - Added contributor initialization

### New Imports
```go
"github.com/ethereum/go-ethereum/common"
"github.com/kawai-network/x/constant"
"github.com/kawai-network/veridium/internal/image"
"github.com/kawai-network/veridium/internal/services"
"github.com/kawai-network/veridium/internal/whisper"
"github.com/kawai-network/veridium/pkg/blockchain"
"github.com/kawai-network/veridium/pkg/hardware"
"github.com/kawai-network/veridium/pkg/store"
"github.com/kawai-network/veridium/pkg/tunnelkit"
```

### New Config Sections
```go
Contributor struct {
    Enabled           bool
    WalletPassword    string
    WalletAddress     string
    ImportMnemonic    string
    HeartbeatInterval time.Duration
}

Blockchain struct {
    RPCUrl           string
    TokenAddress     string
    OTCMarketAddress string
    USDTAddress      string
}

Tunnel struct {
    Enabled bool
}
```

---

## Startup Flow Comparison

### Old (cmd/contributor)
```
main()
├─ Init Sentry
├─ Init KV Store
├─ Setup Wallet
├─ Register Holder
├─ Detect Hardware
├─ Start Tunnel
├─ Register Contributor
├─ Start Heartbeat
├─ Init LLM (llamalib)
├─ Init Whisper
├─ Init Stable Diffusion
└─ Start Gateway Server
```

### New (cmd/server)
```
main() → kronk.Run()
├─ Init Sentry
├─ Download Libraries (llama.cpp)
├─ Build Model Index
├─ Download Catalog & Templates
├─ Init Kronk Cache
│
├─ [IF CONTRIBUTOR_ENABLED]
│   ├─ Init KV Store
│   ├─ Init Blockchain Client
│   ├─ Setup Wallet
│   ├─ Register Holder
│   ├─ Detect Hardware
│   ├─ Start Tunnel
│   ├─ Register Contributor
│   ├─ Start Heartbeat
│   ├─ Init Whisper
│   └─ Init Stable Diffusion
│
└─ Start HTTP Server (OpenAI-compatible)
```

---

## Benefits of Migration

### 1. Better Architecture
- Layered design (foundation/app/api)
- Separation of concerns
- Testable components

### 2. Production Features
- Model caching with TTL
- Concurrent inference
- OpenAI-compatible API
- Structured logging + tracing

### 3. Flexibility
- Can disable contributor features
- Can run local-only (no tunnel)
- Environment-based configuration

### 4. Maintainability
- Single codebase to maintain
- Consistent error handling
- Better observability

---

## Migration Guide for Users

### Step 1: Stop Old Contributor
```bash
# Stop cmd/contributor if running
pkill -f contributor
```

### Step 2: Set Environment Variables
```bash
# Create .env or export
export KRONK_CONTRIBUTOR_ENABLED=true
export KRONK_CONTRIBUTOR_WALLETPASSWORD=your_password

# Optional: blockchain config (uses defaults if not set)
export KRONK_BLOCKCHAIN_RPCURL=https://rpc.monad.xyz
export KRONK_BLOCKCHAIN_TOKENADDRESS=0xBd95bDB3a6FE48CbC2dE3890B8e67Ef96Af65322
```

### Step 3: Run New Server
```bash
./server
```

### Step 4: Verify
```bash
# Check logs for:
# ✓ Connected to Cloudflare KV
# ✓ Wallet unlocked
# ✓ Holder registered
# ✓ Hardware detected
# ✓ Tunnel started
# ✓ Contributor registered
# ✓ Heartbeat started
```

---

## Testing Checklist

### Phase 1: Core Features ✅
- [x] Config parsing
- [x] KV store initialization
- [x] Wallet creation
- [x] Wallet unlock
- [x] Blockchain client init
- [x] Holder registration
- [x] Hardware detection
- [x] Contributor registration
- [x] Heartbeat system
- [x] Tunnel integration
- [x] Whisper service init
- [x] Stable Diffusion init

### Phase 2: API Endpoints (TODO)
- [ ] Whisper transcription endpoint
- [ ] Stable Diffusion generation endpoint
- [ ] Contributor status endpoint
- [ ] Hardware info endpoint

### Phase 3: Integration Tests (TODO)
- [ ] End-to-end contributor flow
- [ ] Wallet operations
- [ ] Heartbeat reliability
- [ ] Tunnel connectivity
- [ ] AI service requests

---

## Known Issues

### 1. go-ethereum Dependency
```
Error: table.SetHeader undefined
```
**Status:** Pre-existing issue, not related to migration
**Workaround:** Will be fixed in dependency update

### 2. API Routes Not Added Yet
**Status:** Phase 2 work
**Impact:** Whisper/SD endpoints not exposed yet
**ETA:** Next sprint

---

## Next Steps

### Immediate (Phase 2)
1. Add Whisper endpoint to mux
2. Add Stable Diffusion endpoint to mux
3. Add contributor info endpoints
4. Add static file serving for SD outputs

### Short-term (Phase 3)
1. Write integration tests
2. Test wallet operations
3. Test heartbeat reliability
4. Test tunnel connectivity

### Long-term (Phase 4)
1. Deprecate `cmd/contributor` officially
2. Remove `cmd/contributor` code
3. Update all documentation
4. Update deployment scripts
5. Announce migration to users

---

## Rollback Plan

If issues arise, rollback is simple:

```bash
# Use old contributor binary
./contributor --password mypass
```

All data (wallets, KV entries) remain compatible.

---

## Support

For migration issues:
1. Check logs for detailed errors
2. Verify environment variables
3. Test wallet password
4. Check KV connectivity
5. Verify blockchain RPC

---

## Credits

- **Original contributor:** `cmd/contributor/main.go`
- **Kronk framework:** Ardan Labs
- **Migration:** Integrated into `cmd/server/api/services/kronk/kronk.go`
- **Date:** 2026-01-31
