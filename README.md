# Kawai DeAI Network (Monad)

**Decentralized AI Compute Network on Monad.**
A "Lean Startup" approach to DePIN, leveraging consumer-grade GPUs for `llama.cpp` inference with a sustainable tokenomic model.

---

## 🚀 Core Concept

-   **Service:** Low-cost LLM Inference API (compatible with OpenAI format).
-   **Contributors:** Gamers & Devs running `llama.cpp` nodes.
-   **Rewards:** 
    -   **Contributor:** Earn KAWAI tokens via "Mining" (70% split). Rewards follow a **Halving Schedule** based on total supply (100 -> 50 -> 25 -> 12 per 1M tokens).
    -   **Holder:** Earn 100% of Platform Revenue (USDT) proportional to KAWAI holdings.
-   **Liquidity:** No Initial LP. Value follows the *Hold-to-Earn* utility.
-   **Details:** See [Concept Document](current_concept.md) for full analysis.

---

## ⚙️ How It Works (Summary)

We use a **Hybrid Model** (Off-Chain Accumulation + On-Chain Settlement) to minimize gas fees.
*   **Real-Time:** Rewards are calculated instantly (Usage-based) and split 70/30.
*   **Weekly:** A compressed **Merkle Tree** is uploaded to **Monad** for cheap claiming.

### Bootstrap Phase (Current)
*   **User Client:** Runs local LLM inference (no external contributors yet).
*   **Rewards:** Distributed to treasury pool addresses as "virtual contributors".
*   **Purpose:** Start token economy and accumulate treasury for future liquidity.
*   **Transition:** Ready to switch to distributed network when contributors join.

👉 **[See Full Technical Details in Concept Document](current_concept.md#d-mekanisme-teknis-hybrid-how-it-works)**

---

## 📂 Project Structure & Usage

Veridium consists of two primary components designed for different user roles:

### 1. User Client (Desktop App)
- **Location:** `main.go` (Root)
- **Tech:** [Wails v3](https://wails.io/) + React.
- **Description:** The main interface for AI consumers. It provides a premium desktop experience for chat, knowledge base management, and a Web3 dashboard for USDT deposits and token management.
- **Current Mode:** Local LLM inference with automatic reward recording to treasury pool.
- **How to Run:**
  Refer to the `Makefile` for various development and build commands:
  ```bash
  make dev      # Start fresh (reset DB + full build)
  make build    # Build production binary
  ```

### 2. Contributor Client (CLI)
- **Location:** `cmd/contributor/main.go`
- **Tech:** Go CLI.
- **Description:** The "Miner" application. It wraps `llama.cpp` to provide compute power to the network and earn KAWAI tokens.
- **Status:** Ready for distributed network phase (future).
- **How to Run:**
  ```bash
  go run cmd/contributor/main.go --password YOUR_PASSWORD
  ```

> **Note:** `cmd/server` is currently **deprecated**/not used in the production flow.

---

## 🗺️ Implementation Plan (Roadmap)

This roadmap outlines the path from "Zero" to a fully functional decentralized network.

### Phase 1: Foundation ✅
**Focus:** Infrastructure Setup & Smart Contract Development.
- [x] **Project Initialization:**
    - [x] Go Backend Structure (Middleware).
    - [x] Foundry Environment (Contracts).
- [x] **Auto-Binding Setup:**
    - [x] Makefile for `abigen` (Solidity -> Go).
- [x] **Smart Contracts (MVP):**
    - [x] `KawaiToken.sol`: Standard ERC20 with AccessControl (Mint/Burn).
    - [x] `Escrow.sol`: Simple P2P OTC Market (Orders, Buy, Cancel).
    - [x] `PaymentVault.sol`: Prepaid USDT Deposit system for Consumers.
    - [x] `MerkleDistributor.sol`: Gas-efficient reward claiming system.
- [x] **Middleware (Go):**
    - [x] `pkg/blockchain`: Service to interact with **Monad** (Listen events, Send TXs).
    - [x] `pkg/store`: Persistent storage using Cloudflare Workers KV (Multi-namespace).
    - [x] `pkg/merkle`: Merkle Tree generation logic.

### Phase 2: The "Lean" Launch (Internal Market) ✅
**Focus:** Economic bootstrapping without initial LP.
- [x] **Contributor Client (Go):** 
    - [x] **Wallet Setup:** Generate/Load Contributor Identity (Private Key).
    - [x] Wrapper for `llama.cpp` server (Process Management).
    - [x] Heartbeat system (Proof of Availability).
- [x] **API Gateway & Billing:**
    - [x] OpenAI-compatible `/v1/chat/completions` proxy.
    - [x] Real-time credit deduction system.
        - [x] API key validation & user identification per request.
        - [x] Track token usage per user per request.
        - [x] Deduct USDT credits from user's balance in real-time.
        - [x] Reject requests if insufficient balance.
        - [x] **JSON Storage**: Atomic balance updates using JSON in dedicated User Namespace.
    - [x] **Free Trial System:**
        - [x] **5 USDT Bonus**: Automatic virtual credits for new users.
        - [x] **Anti-Abuse**: Dual-layer protection (Wallet Address + Machine ID).
        - [x] **Atomic Claims**: Prevents race conditions during trial claiming.
    - [x] **Job Reward Recording (Bootstrap Phase):**
        - [x] Record rewards for local LLM inference to treasury pool.
        - [x] Random treasury address selection for fair distribution.
        - [x] Automatic 70/30 split (Treasury/Admin).
        - [x] Phase detection (Mining vs USDT mode).
        - [x] Halving logic (100→50→25→12 KAWAI per 1M tokens).
        - [x] Async execution (no user latency impact).
    - [x] **Deposit Sync Service:**
        - [x] User client syncs deposits via `SyncDeposit(txHash)` after wallet confirmation.
        - [x] Backend verifies transaction onchain before updating balance.
        - [x] Duplicate prevention (transaction can only be synced once).
        - [x] Security: validates transaction status, event data, and user address.
        - [x] Automatic sync after deposit (triggered by user client).
- [x] **Internal P2P Marketplace (Web):**
    - [x] Complete marketplace UI with order book, trading interface, and real-time updates.
    - [x] Order creation, cancellation, and execution functionality.
    - [x] Market statistics dashboard with price tracking and volume data.
    - [x] User order history and trade history with detailed tracking.
    - [x] Real-time event integration using Wails Events API.
    - [x] Granular loading states with progressive data loading.
- [x] **Web3 Dashboard:**
    - [x] Login via Wallet to manage API Keys.
    - [x] USDT Deposit UI (PaymentVault integration).
- [x] **Administration Scripts (Go):**
    - [x] Weekly Snapshot & Dividend Calculator.
    - [x] Contributor Audit & Ban system.
- [x] **Dividend & Reward System (Two-Phase Model):**
    - [x] Weekly Snapshot script.
    - [x] Merkle Airdrop implementation (Pull model) for KAWAI.
    - [x] **Phase 1 -> Phase 2 Transition Logic:**
        - [x] Detect when `totalSupply() == MAX_SUPPLY` in Smart Contract / Middleware.
        - [x] Switch payment mode from KAWAI mining to USDT cost-based.
    - [x] **Phase 2 USDT Payouts:**
        - [x] Implement `COST_RATE_PER_MILLION` as configurable Environment Variable.
        - [x] Calculate Contributor/Admin USDT cost per job in `RecordJobReward`.
        - [x] Distribute remaining USDT Profit to KAWAI Holders weekly.
- [x] **Settlement System (pkg/store):**
    - [x] Multi-namespace KV architecture (Contributors, Proofs, Settlements).
    - [x] Period-specific Merkle proofs with claim status tracking.
    - [x] Atomic settlement with rollback support.
    - [x] Soft delete for contributor lifecycle management.

### Phase 3: Community Liquidity (Growth)
**Focus:** Transition to decentralized liquidity.
- [ ] **Liquidity Incentives:**
    - [ ] Launch "LP Mining" program.
    - [ ] Encourage stakeholders to provide LP on PancakeSwap using their Dividends.
- [x] **Frontend (Web Dashboard):**
    - [x] Initialize Vite + React project.
    - [x] Integrate @lobehub/ui for UI components.
    - [x] Complete P2P marketplace interface with real-time trading capabilities.
    - [x] Wallet integration with multi-wallet support and secure authentication.
    - [x] Market data visualization with statistics and price tracking.
    - [x] Order management system with creation, cancellation, and history tracking.

---

## 🌐 Monad Testnet Deployment

> **Status:** ✅ **UPGRADED v2.0** (2025-12-31) - Partial Fill Support

### **Latest Contracts (v2.0 - Partial Fill)**

| Contract | Address | Description |
|---|---|---|
| **KawaiToken** | `0xF27c5c43a746B329B1c767CE1b319c9EBfE8012E` | Native token (1B Max Supply, Fair Launch) |
| **OTCMarket v2** ⭐ | `0x5b1235038B2F05aC88b791A23814130710eFaaEa` | **NEW:** P2P trading with Partial Fill Support |
| **MockUSDT** | `0xb8cD3f468E9299Fa58B2f4210Fe06fe678d1A1B7` | Test USDT for simulating revenue |
| **PaymentVault** | `0x714238F32A7aE70C0D208D58Cc041D8Dda28e813` | User USDT deposit vault |
| **KAWAI_Distributor** | `0xf4CCb09208cA77153e1681d256247dae0ff119ba` | Merkle distributor for mining rewards |
| **USDT_Distributor** | `0xE964B52D496F37749bd0caF287A356afdC10836C` | Merkle distributor for profit sharing |

### **What's New in v2.0:**
- ✅ **Partial Fill Support** - Buy/sell any amount from orders
- ✅ **Remaining Amount Tracking** - Real-time order status
- ✅ **Efficient View Functions** - Batch queries, pagination
- ✅ **Enhanced Events** - Detailed partial fill events
- ✅ **19 Comprehensive Tests** - Full test coverage
- ✅ **Reconciliation Service** - Auto-sync blockchain ↔ KV store
- ✅ **Rate Limiting** - RPC (10/sec) + KV (100/sec)

### **Old Contracts (v1.0 - Deprecated)**
<details>
<summary>Click to view legacy contracts</summary>

| Contract | Address | Status |
|---|---|---|
| KawaiToken (old) | `0x3EC7A3b85f9658120490d5a76705d4d304f4068D` | ⚠️ Deprecated |
| OTCMarket (old) | `0x134244eDd4349b0B408c5293Ffb4263984F2808C` | ⚠️ Deprecated |

</details>

**Network Info:**
- **Chain ID:** 10143
- **RPC:** `https://testnet-rpc.monad.xyz`
- **Explorer:** `https://testnet.monad.xyz`
- **Deployer:** `0x94D5C06229811c4816107005ff05259f229Eb07b`

---

## 🔧 Environment Variables

Create a `.env` file in the project root with the following:

```bash
# Blockchain Configuration
MONAD_RPC_URL=https://testnet-rpc.monad.xyz
TOKEN_ADDRESS=0x...  # KawaiToken address after deployment
ESCROW_ADDRESS=0x... # OTCMarket address after deployment

# Economic Configuration (Phase 1 & 2)
KAWAI_RATE_PER_MILLION=100      # KAWAI minted per 1M tokens processed
COST_RATE_PER_MILLION=1.0       # USDT cost per 1M tokens (Phase 2)
FREE_TRIAL_AMOUNT_USDT=5.0      # USDT bonus for new users (Default: 5)

# Admin Configuration
ADMIN_ADDRESS=0x...             # Admin wallet for fee collection
ADMIN_PRIVATE_KEY=0x...         # For signing Merkle Root updates

# Cloudflare KV (Multi-Namespace Architecture)
CF_ACCOUNT_ID=...
CF_API_TOKEN=...
CF_KV_CONTRIBUTORS_NAMESPACE_ID=...  # Contributor data
CF_KV_PROOFS_NAMESPACE_ID=...        # Merkle proofs
CF_KV_SETTLEMENTS_NAMESPACE_ID=...   # Settlement metadata
CF_KV_AUTHZ_NAMESPACE_ID=...         # API Keys
CF_KV_USERS_NAMESPACE_ID=...         # User Profiles & Balance (JSON)
```

---

## 🛠️ Tech Stack

-   **Smart Contracts:** Solidity, Foundry, OpenZeppelin.
-   **Backend (Middleware):** Go (Golang), Gin, `go-ethereum`, `cloudflare-go`.
-   **Database:** Cloudflare Workers KV (Multi-namespace).
-   **Frontend:** React 19, Vite, @lobehub/ui, Zustand.
-   **Contributor Node:** Go (Golang), `llama.cpp` (via llamalib).
-   **Blockchain:** Monad (EVM-compatible).
-   **Network Toolkit:** `pkg/jarvis` (Multi-chain support incl. Monad).

---

## 🔒 Quality & Reliability

### Recent Improvements (January 2026)

**P1 Bug Fixes - Critical Issues** ✅
- **Silent Failure Handling**: All errors now properly logged with context
  - Model initialization errors tracked
  - Vector operations failures logged
  - Database operations monitored
- **Input Validation**: Comprehensive validation for user inputs
  - Memory service: validate memory_id and query parameters
  - Vector search: validate query parameters
  - Clear error messages for invalid inputs

**P2/P3 Bug Fixes - Resilience & Configuration** ✅
- **Blockchain Resilience**: Automatic retry with exponential backoff
  - Retry up to 3 times for transient network errors
  - Exponential backoff: 500ms → 1s → 2s → 5s (max)
  - 30-second timeout for blockchain receipt fetching
  - Better error logging for debugging
- **Centralized Configuration**: Single source of truth for timeouts
  - Organized by operation type (blockchain, LLM, database, file, network, cache)
  - Easy to tune performance globally
  - Prepared for environment variable overrides
  - See `internal/constant/timeouts.go` for all timeout values

**Race Condition Fixes** ✅
- **Deposit Double-Spend Prevention**: Idempotency pattern implementation
  - Transaction verification before balance update
  - Duplicate prevention with processed transaction tracking
  - Safe for concurrent deposit sync requests
- **Reward Distribution**: Per-address mutex for concurrent job rewards
  - Prevents lost updates in concurrent scenarios
  - Serializes updates to same contributor address
  - Maintains data consistency
- **Free Trial Atomic Claims**: Implementation of the *Read-Modify-Write* pattern in the KV store to prevent concurrent bypasses of the trial claim logic.

### Code Quality Features
- **Error Handling**: Comprehensive error types and messages
- **Observability**: Structured logging with context
- **Resource Management**: Proper cleanup with defer patterns
- **Rate Limiting**: Built-in rate limiters for RPC and KV operations
- **Idempotency**: Safe retry patterns for critical operations
- **Timeout Management**: Consistent timeout behavior across all operations

---

## 📦 Package Documentation

| Package | Description |
|---------|-------------|
| [`pkg/store`](pkg/store/README.md) | Off-chain KV storage (Contributors, Proofs, Settlements) |
| `pkg/merkle` | Merkle tree generation |
| `pkg/blockchain` | Monad blockchain interaction |
| `pkg/config` | Configuration management |
| [`internal/services`](internal/services/) | Marketplace, Reconciliation, Event Listener services |

## 📚 Additional Documentation

| Document | Description |
|----------|-------------|
| [`DEPLOYMENT_SUMMARY.md`](docs/DEPLOYMENT_SUMMARY.md) | Full deployment details & contract addresses |
| [`MARKETPLACE_UPGRADE_SUMMARY.md`](docs/MARKETPLACE_UPGRADE_SUMMARY.md) | Architecture & partial fill implementation |
| [`CONTRACTS_WORKFLOW.md`](docs/CONTRACTS_WORKFLOW.md) | Smart contract development workflow |
| [`DEPOSIT_SYNC.md`](docs/DEPOSIT_SYNC.md) | Deposit synchronization system & implementation guide |
| [`current_concept.md`](current_concept.md) | Complete project concept & tokenomics |

---

sentryDSN_golang = "https://709dabacc882a777ef059392d056e3da@o4510568649654272.ingest.us.sentry.io/4510568655290368"
sentryDSN_golang_contributor = "https://6d138acbdde2516e32e24f016b472031@o4510620614983680.ingest.us.sentry.io/4510620618850304" // dielzzz89
sentryDSN_react = "https://b66f862d7567c075a44c697757bb8130@o4510618985758720.ingest.us.sentry.io/4510618990804992" // yudapramad