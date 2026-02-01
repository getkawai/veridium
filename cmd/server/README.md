# Kronk Model Server - Migration Complete

## Overview

Kronk server telah berhasil dimigrate dengan **contributor features** dari `cmd/contributor`. Server sekarang adalah **production-grade OpenAI-compatible model server** dengan kemampuan contributor network. 

---

## Architecture Evolution

### Before: cmd/contributor (Deprecated)
- Standalone CLI tool 
- Custom gateway HTTP server
- Direct llamalib integration
- Simple architecture

### After: cmd/server (Current)
- Production-grade Kronk framework
- OpenAI-compatible API
- Model caching + concurrent inference
- **Contributor features integrated**

---

## Features Migrated

✅ **Wallet Management** - Create/unlock/switch wallets
✅ **KV Store Integration** - Cloudflare KV for persistence
✅ **Blockchain Client** - Halving logic & token interactions
✅ **Holder Registry** - Track KAWAI token holders
✅ **Hardware Detection** - Auto-detect CPU/GPU/RAM specs
✅ **Contributor Registry** - Register to network
✅ **Heartbeat System** - 30s keep-alive mechanism
✅ **Tunnel Integration** - Public endpoint via tunnelkit
✅ **Whisper Service** - Speech-to-text capability (go-whisper)
✅ **Stable Diffusion** - Image generation capability

---

## Configuration

### Environment Variables

```bash
# Contributor Features (enabled by default)
KRONK_CONTRIBUTOR_ENABLED=true
KRONK_CONTRIBUTOR_WALLETPASSWORD=your_password
KRONK_CONTRIBUTOR_WALLETADDRESS=0x...  # Optional: switch wallet
KRONK_CONTRIBUTOR_IMPORTMNEMONIC="word1 word2..."  # Optional: import
KRONK_CONTRIBUTOR_HEARTBEATINTERVAL=30s

# Blockchain Configuration (defaults from constants)
KRONK_BLOCKCHAIN_RPCURL=https://rpc.monad.xyz
KRONK_BLOCKCHAIN_TOKENADDRESS=0xBd95bDB3a6FE48CbC2dE3890B8e67Ef96Af65322
KRONK_BLOCKCHAIN_OTCMARKETADDRESS=0x75d1A6CC51035D7E5Cbe88aEc6DCfd6ABEB22bfE
KRONK_BLOCKCHAIN_USDTADDRESS=0x754704bc059f8c67012fed69bc8a327a5aafb603

# Tunnel Configuration
KRONK_TUNNEL_ENABLED=true

# Existing Kronk Config
KRONK_WEB_APIHOST=localhost:8080
KRONK_CACHE_MODELSINCACHE=3
KRONK_CACHE_TTL=20m
```

---

## Usage

### Mode 1: Pure Model Server (Contributor Disabled)
```bash
KRONK_CONTRIBUTOR_ENABLED=false ./server
```

### Mode 2: Contributor Mode (Default)
```bash
# First time - creates wallet
KRONK_CONTRIBUTOR_WALLETPASSWORD=mypass ./server

# Subsequent runs - unlocks wallet
KRONK_CONTRIBUTOR_WALLETPASSWORD=mypass ./server

# Import existing wallet
KRONK_CONTRIBUTOR_WALLETPASSWORD=mypass \
KRONK_CONTRIBUTOR_IMPORTMNEMONIC="word1 word2..." \
./server
```

### Mode 3: No Tunnel (Local Only)
```bash
KRONK_TUNNEL_ENABLED=false \
KRONK_CONTRIBUTOR_WALLETPASSWORD=mypass \
./server
```

---

## Startup Flow

```
main() → kronk.Run() →
  ├─ Init Sentry Logging
  ├─ Download Libraries (llama.cpp)
  ├─ Build Model Index
  ├─ Download Catalog & Templates
  ├─ Init Kronk Cache
  │
  ├─ [CONTRIBUTOR FEATURES]
  │   ├─ Init KV Store
  │   ├─ Init Blockchain Client
  │   ├─ Setup Wallet (create/unlock)
  │   ├─ Register Holder
  │   ├─ Detect Hardware
  │   ├─ Start Tunnel (optional)
  │   ├─ Register Contributor
  │   ├─ Start Heartbeat (30s)
  │   ├─ Init Whisper Service
  │   └─ Init Stable Diffusion
  │
  └─ Start HTTP Server (OpenAI-compatible)
```

---

## API Endpoints

### Kronk (Existing)
- `POST /v1/chat/completions` - Chat completions
- `POST /v1/embeddings` - Text embeddings
- `POST /v1/rerank` - Document reranking
- `GET /v1/models` - List models
- `POST /v1/images/generations` - Image generation (Stable Diffusion)
- `GET /files/*` - Serve generated images

### Contributor (Implemented)
- `POST /v1/audio/transcriptions` - Whisper STT (go-whisper)
- `GET /v1/contributor/status` - Contributor info (TODO)
- `GET /v1/contributor/hardware` - Hardware specs (TODO)

---

## Next Steps

### Phase 1: ✅ Core Migration (DONE)
- [x] Add contributor config
- [x] Integrate KV store
- [x] Integrate wallet service
- [x] Integrate blockchain client
- [x] Add holder registry
- [x] Add hardware detection
- [x] Add tunnel integration
- [x] Add contributor registry
- [x] Add heartbeat system
- [x] Integrate Whisper service
- [x] Integrate Stable Diffusion

### Phase 2: ✅ API Routes (DONE)
- [x] Add image generation endpoint to mux
- [x] Add file server for generated images
- [x] Add Whisper endpoint to mux (go-whisper)
- [ ] Add contributor status endpoint
- [ ] Add hardware info endpoint

### Phase 3: Testing (TODO)
- [ ] Test wallet creation/unlock
- [ ] Test contributor registration
- [ ] Test heartbeat system
- [ ] Test tunnel connectivity
- [ ] Test Whisper transcription
- [ ] Test SD image generation

### Phase 4: Cleanup (TODO)
- [ ] Deprecate `cmd/contributor`
- [ ] Update main documentation
- [ ] Add migration guide
- [ ] Update deployment scripts

---

## Deprecation Notice

**cmd/contributor** is now deprecated. All functionality has been migrated to **cmd/server**.

To migrate:
1. Use `cmd/server` instead of `cmd/contributor`
2. Set environment variables (see Configuration above)
3. Run with `KRONK_CONTRIBUTOR_ENABLED=true`

---

## Development

### Build
```bash
go build -o server ./cmd/server
```

### Run
```bash
KRONK_CONTRIBUTOR_WALLETPASSWORD=test ./server
```

### Test
```bash
go test ./cmd/server/...
```

---

## Release

### Automated Release (GitHub Actions)

Trigger release dengan tag:
```bash
git tag node-v1.0.0
git push origin node-v1.0.0
```

Workflow akan otomatis:
1. Build untuk Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64)
2. Upload ke GitHub Release
3. Upload ke Cloudflare R2 (`/node/v{version}/`)
4. Update symlink `latest`

### Ephemeral Release (Private → Public → Delete)

Untuk release dari private repo tanpa expose source code:

```bash
# Using ephemeral script
./scripts/ephemeral-release-node.sh 1.0.0
```

Script ini akan:
1. Create clean copy dengan hanya `cmd/server`
2. Create temporary public repo
3. Trigger GitHub Actions build
4. Download artifacts
5. Upload ke R2
6. Delete public repo

### Install Script

Users dapat install dengan:
```bash
curl -fsSL https://getkawai.com/node | sh
```

Atau dengan version specific:
```bash
curl -fsSL https://getkawai.com/node | sh -s -- --version 1.0.0
```

### Manual Build & Release

```bash
# Build locally
go build -ldflags "-s -w -X main.version=1.0.0" -o kawai-node ./cmd/server

# Create archive
tar -czf kawai-node-1.0.0-linux-amd64.tar.gz kawai-node

# Upload to R2
aws s3 cp kawai-node-1.0.0-linux-amd64.tar.gz s3://kawai/node/v1.0.0/ \
  --endpoint-url https://your-account.r2.cloudflarestorage.com
```

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Kronk Model Server                       │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   LLM Cache  │  │   Whisper    │  │ Stable Diff  │      │
│  │  (Kronk SDK) │  │   Service    │  │    Engine    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           OpenAI-Compatible HTTP API                  │   │
│  │  /v1/chat/completions  /v1/embeddings  /v1/rerank   │   │
│  │  /v1/audio/transcriptions  /v1/images/generations   │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Contributor Features                     │   │
│  │  • Wallet Management    • Holder Registry            │   │
│  │  • KV Store             • Hardware Detection         │   │
│  │  • Blockchain Client    • Tunnel Integration         │   │
│  │  • Heartbeat System     • Contributor Registry       │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Cloudflare KV   │
                    │  Monad Blockchain│
                    │  Tunnel Network  │
                    └──────────────────┘
```

---

## Notes

- Contributor features are **enabled by default** (`KRONK_CONTRIBUTOR_ENABLED=true`)
- Blockchain config uses constants from `internal/constant/blockchain.go`
- Wallet password is **required** when contributor mode is enabled
- Tunnel is **optional** - can run local-only mode
- All contributor features gracefully handle errors (non-fatal)
- Heartbeat runs in background goroutine
- Cleanup on shutdown (mark contributor offline)

---

## Support

For issues or questions:
- Check logs for detailed error messages
- Verify environment variables are set correctly
- Ensure wallet password is correct
- Check KV store connectivity
- Verify blockchain RPC is accessible