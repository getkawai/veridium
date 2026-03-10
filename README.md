# Kawai DeAI Network (Monad)

**Decentralized AI Compute Network on Monad.**
A "Lean Startup" approach to DePIN, leveraging consumer-grade GPUs for `llama.cpp` inference with a sustainable tokenomic model.

---

## 🚀 Core Concept

-   **Service:** Low-cost AI Generation and API (compatible with OpenAI format).
-   **Contributors:** Gamers & Devs running `llama.cpp` nodes.
-   **Rewards:** 
    -   **Contributor:** Install kawai-contributor, earn KAWAI tokens via "Mining" (85-90% split). Rewards follow a **Halving Schedule** based on total supply (100 -> 50 -> 25 -> 12 per 1M tokens).
    -   **User:** Install kawai-desktop, earn 5% cashback on every AI request (use-to-earn) + 1-5% deposit cashback (tiered).
    -   **Affiliator:** Earn 5% commission from referrals' mining rewards (lifetime passive income).
    -   **Developer:** Earn 5% from mining rewards (distributed to treasury pool).
    -   **Holder:** ✅ Earn 100% of Platform Revenue (USDT) proportional to KAWAI holdings via **Hybrid Holder Registry** (Production Ready).
-   **Liquidity:** No Initial LP. Value follows the *Hold-to-Earn* utility.
-   **Details:** See tokenomics section below for full analysis.

---

## ⚙️ How It Works (Summary)

We use a **Hybrid Model** (Off-Chain Accumulation + On-Chain Settlement) to minimize gas fees.
*   **Real-Time:** Rewards are calculated instantly (Usage-based) and split 70/30.
*   **Weekly:** A compressed **Merkle Tree** is uploaded to **Monad** for cheap claiming.

### Current Phase (Phase 1 - Mining Era)
*   **User Client:** Connects to contributor nodes for AI inference.
*   **Payment:** Users deposit stablecoin (USDC on mainnet, MockUSDT on testnet) to PaymentVault and pay per token usage.
*   **Contributors:** Earn KAWAI tokens (85-90% mining rewards) for providing compute.
*   **Platform Revenue:** 100% of stablecoin payments (contributors paid in KAWAI, not stablecoin).
*   **Revenue Sharing:** Platform profit distributed to KAWAI holders as dividends.

👉 **See Full Technical Details in the sections below**

---

## 📂 Project Structure & Usage

Veridium consists of two primary components designed for different user roles:

### 1. User Client (Desktop App)
- **Location:** `main.go` (Root)
- **Tech:** [Wails v3](https://wails.io/) + React.
- **Description:** The main interface for AI consumers. It provides a premium desktop experience for chat, knowledge base management, and a Web3 dashboard for stablecoin deposits and token management.
- **How It Works:** 
  - Users deposit stablecoin (USDC/MockUSDT) to PaymentVault contract
  - Desktop app connects to contributor nodes for AI inference
  - Stablecoin balance deducted per token usage
  - Contributors earn KAWAI tokens, platform earns USDT revenue
- **How to Run:**
  Refer to the `Makefile` for various development and build commands.

### 2. Contributor Client (CLI)
- **Location:** `cmd/server/` (migrated from `cmd/contributor`)
- **Tech:** Go CLI using Kronk framework.
- **Description:** The "Miner" application. It wraps `llama.cpp` to provide compute power to the network and earn KAWAI tokens.
- **Status:** Migrated to `cmd/server`. See `cmd/server/README.md` for details.

---

## 🌐 Monad Testnet Deployment

All smart contracts are deployed on **Monad Blockchain** (Testnet).

**See:** [`docs/CONTRACTS_OVERVIEW.md`](docs/CONTRACTS_OVERVIEW.md) for complete contract details, addresses, and deployment information.

### **Key Contracts:**
- **KawaiToken** - Native utility token (1B max supply, fair launch)
- **Escrow (OTC Market)** - P2P KAWAI ↔ USDT trading with partial fill support
- **PaymentVault** - User USDT deposit management
- **Reward Distributors** - Mining (85/5/5/5), Cashback (1-5% tiered), Referral rewards
- **MerkleDistributor** - Generic Merkle-based distribution (USDT profit sharing)

**Network Info:**
- **Chain ID:** 10143
- **RPC:** `https://testnet-rpc.monad.xyz`
- **Explorer:** `https://testnet.monad.xyz`
- **Deployer:** `0x94D5C06229811c4816107005ff05259f229Eb07b`

---

## 🛠️ Tech Stack

-   **Smart Contracts:** Solidity, Foundry, OpenZeppelin.
-   **Backend (Middleware):** Go (Golang), Gin, `go-ethereum`, `cloudflare-go`.
-   **Database:** Cloudflare Workers KV (Multi-namespace).
-   **Frontend:** React 19, Vite, @lobehub/ui, Zustand.
-   **Contributor Node:** Go (Golang), `llama.cpp` (via llamalib).
-   **Blockchain:** Monad (EVM-compatible).
-   **Network Toolkit:** `github.com/kawai-network/x/jarvis` (Multi-chain support incl. Monad).
-   **Dependency Injection:** Uber FX (`go.uber.org/fx`) for service lifecycle management.

---

## 💉 Dependency Injection with Uber FX

**New in v2.0** - The application now uses [Uber FX](https://uber-go.github.io/fx/) for dependency injection and lifecycle management.

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      main.go                                 │
│  fx.New(                                                     │
│    app.Module,           // All service providers            │
│    fx.Provide(...),      // App-specific providers           │
│    fx.Invoke(...),       // Service registration             │
│  ).Run()                                                     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  internal/app/module.go                      │
│  fx.Module("app",                                            │
│    fx.Provide(                                               │
│      ProvideDatabase,      // *database.Service              │
│      ProvideKVStore,       // *store.KVStore                 │
│      ProvideWalletService, // *services.WalletService        │
│      ProvideChatModel,     // unillm.LanguageModel           │
│      ...                                                   │
│    ),                                                        │
│  )                                                           │
└─────────────────────────────────────────────────────────────┘
```

### Key Concepts

#### 1. **Providers** (`fx.Provide`)
Functions that create and return dependencies. FX calls them lazily when needed.

```go
func ProvideDatabase(lc fx.Lifecycle) (*database.Service, error) {
    dbService, err := database.NewService()
    if err != nil {
        return nil, err
    }
    
    // Register cleanup hook
    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            return dbService.Close()
        },
    })
    
    return dbService, nil
}
```

#### 2. **Consumers** (`fx.Invoke`)
Functions that receive dependencies. Called after all providers are ready.

```go
func RegisterWailsServices(
    wailsApp *application.App,
    ctx *app.Context,
    fileProcessor *FileProcessorService,
    // ... more dependencies auto-injected
) {
    wailsApp.RegisterService(application.NewService(ctx.Queries))
    // ...
}
```

#### 3. **Named Dependencies**
For multiple instances of the same type (e.g., different language models):

```go
// In module.go
fx.Annotate(ProvideChatModel, fx.ResultTags(`name:"chatModel"`)),
fx.Annotate(ProvideTitleModel, fx.ResultTags(`name:"titleModel"`)),

// In consumer
type Params struct {
    fx.In
    ChatModel  unillm.LanguageModel `name:"chatModel"`
    TitleModel unillm.LanguageModel `name:"titleModel"`
}
```

#### 4. **Lifecycle Hooks**
Automatic cleanup on application shutdown:

```go
lc.Append(fx.Hook{
    OnStart: func(ctx context.Context) error {
        log.Println("Service starting...")
        return nil
    },
    OnStop: func(ctx context.Context) error {
        log.Println("Service stopping...")
        return cleanup()
    },
})
```

### Service Dependency Graph

```
database.Service
    ├── Queries (*db.Queries)
    ├── WalletService
    │   ├── BlockchainClient
    │   └── DeAIService
    ├── VectorSearchService
    │   └── KnowledgeBaseService
    │       └── AgentChatService
    └── ThreadManagementService
        └── TopicService

KVStore (*store.KVStore)
    ├── WalletService
    ├── MarketplaceService
    ├── ReferralService
    └── CashbackService

Language Models (named)
    ├── chatModel → AgentChatService
    ├── titleModel → TopicService
    ├── summaryModel → (future use)
    └── cleanupModel → FileProcessorService
```

### Testing with DI

Mocking services is now straightforward:

```go
func TestMyService(t *testing.T) {
    // Create mock dependencies
    mockDB := &MockDatabase{}
    mockKV := &mocks.KVStore{}
    
    // Inject mocks into service
    svc := services.NewMyService(mockDB, mockKV)
    
    // Setup expectations
    mockDB.On("Query", mock.Anything).Return(expectedData, nil)
    
    // Run test
    result := svc.DoSomething()
    
    // Verify
    mockDB.AssertExpectations(t)
}
```

### Migration Guide (Old → New)

| Old Pattern | New Pattern |
|-------------|-------------|
| `ctx.InitDatabase()` | `ProvideDatabase(lc fx.Lifecycle)` |
| `ctx.InitWalletService()` | `ProvideWalletService(kv *store.KVStore)` |
| Manual `defer cleanup()` | `lc.Append(fx.Hook{OnStop: ...})` |
| `ctx.DB` direct access | Injected via constructor params |

### File Locations

| File | Purpose |
|------|---------|
| `internal/app/module.go` | All FX providers and module definition |
| `internal/app/context.go` | Context struct (now a simple container) |
| `internal/llm/chain.go` | LLM chain building logic |
| `main.go` | FX app initialization and wiring |

### Debugging Tips

1. **Enable FX logging:**
   ```bash
   export FX_VERBOSE=1
   ```

2. **View dependency graph:**
   ```bash
   go run main.go 2>&1 | grep "provided\|invoking"
   ```

3. **Common errors:**
   - `missing dependencies`: Add provider to `fx.Provide()`
   - `duplicate provides`: Check for duplicate type annotations
   - `circular dependency`: Refactor to break the cycle

### Resources

- [Uber FX Documentation](https://uber-go.github.io/fx/)
- [FX GitHub Repository](https://github.com/uber-go/fx)
- [FX Introductory Blog Post](https://eng.uber.com/fx/)

---

## 🧠 Memory & Retrieval Architecture

- **Conversation Memory:** Uses **MuninnDB** only (`internal/services/muninn_memory.go` + `memory_integration.go`).
- **RAG / Vector Retrieval:** Uses **DuckDB** only (`internal/services/duckdb_store.go`, `vector_search.go`, `rag_processor.go`).
- **Important:** DuckDB is **not** used as the conversational memory backend.

For implementation details, see [`docs/memory-architecture.md`](docs/memory-architecture.md).

---

## 🎨 AI Image Generation (NEW)

**Gemini API Integration** - Production-ready image generation with Google's latest models.

### Features
- **Two Models Available:**
  - `gemini-2.5-flash-image` (Nano Banana) - Fast generation, 1024px resolution
  - `gemini-3-pro-image-preview` (Nano Banana Pro) - High quality, up to 4K resolution
  
- **Multi-Provider Support:**
  - **Gemini API** - High-quality image generation with Google's latest models.
  - **Cloudflare Workers AI** - Fast, production-ready Flux models.
  - **Pollinations AI** - Reliable fallback service for diverse model support.
  
- **Priority-Based Selection:**
  - Seamlessly switches between providers (Gemini -> Cloudflare -> Pollinations).
  - 10+ aspect ratios supported (1:1, 16:9, 9:16, 4:3, 3:4, etc.).
  - Quality tiers: Standard (1K), HD (2K), Ultra (4K).
  
- **Production Features:**
  - Thread-safe concurrent generation.
  - API key pool with automatic rotation and obfuscation.
  - Context timeout protection and comprehensive error handling.
  - Enhanced logging with provider-specific prefixes (`[Gemini]`, `[CloudflareGen]`).

### Documentation
- 📖 [Full Implementation Guide](GEMINI_IMAGE_GENERATION.md)
- 📖 [Quick Start Guide](docs/IMAGE_GENERATION_QUICK_START.md)

---

## 📦 Package Documentation

**For Go developers working on backend/blockchain code:**

### Core Veridium Packages
| Package | Description | Key Features |
|---------|-------------|--------------|
| [`pkg/store`](pkg/store/README.md) | Off-chain KV storage (Cloudflare Workers) | Contributors, Merkle Proofs, Settlement automation |
| `pkg/merkle` | Merkle tree generation | Gas-efficient reward distribution |
| `pkg/blockchain` | Monad blockchain interaction | Contract calls, event listening, transaction management |
| `pkg/config` | Configuration management | Environment variables, network settings |
| `pkg/jarvis/contracts` | Smart contract Go bindings | Auto-generated from Solidity contracts |

### Services
| Package | Description |
|---------|-------------|
| [`internal/services`](internal/services/) | Core services: Marketplace, Reconciliation, Event Listener, Rewards |

### Utility Packages
See [`pkg/README.md`](pkg/README.md) for complete list including: `localfs`, `obfuscator`, `nodefs`, `nodepath`, etc.

## 📚 Documentation Guide

**Start here:** This README contains the complete overview of the project.

### System Documentation (Root)
| Document | Purpose | When to Read |
|----------|---------|--------------|
| [`REWARD_SYSTEMS.md`](REWARD_SYSTEMS.md) | Overview & comparison of all reward systems | Understanding reward architecture & current status |
| [`MINING_SYSTEM.md`](MINING_SYSTEM.md) | Complete mining rewards implementation (85/5/5/5 split) | Working on mining/contributor rewards |
| [`CASHBACK_SYSTEM.md`](CASHBACK_SYSTEM.md) | Complete cashback implementation status & guide | Working on deposit cashback features |
| [`REFERRAL_SYSTEM.md`](REFERRAL_SYSTEM.md) | Referral system implementation & status | Working on referral features |
| [`REVENUE_SHARING.md`](REVENUE_SHARING.md) | Complete revenue sharing guide (USDT dividends) | Working on revenue settlement |
| [`MINTER_ROLE_REQUIREMENTS.md`](MINTER_ROLE_REQUIREMENTS.md) | Why MINTER_ROLE is needed for reward distributors | Deploying or debugging reward contracts |

### Technical Deep Dive (docs/)
| Document | Purpose | When to Read |
|----------|---------|--------------|
| [`docs/CONTRACTS_OVERVIEW.md`](docs/CONTRACTS_OVERVIEW.md) | Complete smart contracts overview & architecture | Understanding contract system |
| [`docs/CONTRACTS_WORKFLOW.md`](docs/CONTRACTS_WORKFLOW.md) | Smart contract development & deployment workflow | Developing or deploying contracts |
| [`docs/REFERRAL_CONTRACT_GUIDE.md`](docs/REFERRAL_CONTRACT_GUIDE.md) | Detailed referral contract implementation | Working on referral features |
| [`docs/DEPOSIT_CASHBACK_TOKENOMICS.md`](docs/DEPOSIT_CASHBACK_TOKENOMICS.md) | Economic analysis of cashback tiers | Adjusting cashback parameters |
| [`docs/IMAGE_GENERATION_QUICK_START.md`](docs/IMAGE_GENERATION_QUICK_START.md) | Quick start guide for Gemini image generation | Using image generation API |

### AI Features Documentation
| Document | Purpose | When to Read |
|----------|---------|--------------|
| [`GEMINI_IMAGE_GENERATION.md`](GEMINI_IMAGE_GENERATION.md) | Complete Gemini API implementation guide | Understanding image generation |

---

## 🎯 Quick Commands

### Development
```bash
make dev              # Start fresh (reset DB + full build)
make dev-hot          # Hot reload (keep DB, skip build) - fastest
```

### Reward Settlement (Weekly)
```bash
# All reward types at once (mining, cashback, referral, revenue)
make settle-all

# Or per-type
make settle-mining    # Mining rewards (KAWAI)
make settle-cashback  # Cashback rewards (KAWAI)
make settle-referral  # Referral rewards (KAWAI)
make settle-revenue   # Revenue sharing (USDT dividends)
```

### Smart Contracts
```bash
make contracts-compile    # Compile contracts
make contracts-bindings   # Generate Go bindings
make contracts-update     # Compile + bindings
```

### Documentation
```bash
make docs-serve       # Start docs server (http://localhost:8000)
make docs-build       # Build static docs
```

**See `make help` for complete list of commands.**

---# Trigger workflow - Wed Feb  4 00:02:25 WIB 2026
🐶
