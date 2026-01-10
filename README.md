# Kawai DeAI Network (Monad)

**Decentralized AI Compute Network on Monad.**
A "Lean Startup" approach to DePIN, leveraging consumer-grade GPUs for `llama.cpp` inference with a sustainable tokenomic model.

---

## 🚀 Core Concept

-   **Service:** Low-cost LLM Inference API (compatible with OpenAI format) + **AI Image Generation** (Gemini API).
-   **Contributors:** Gamers & Devs running `llama.cpp` nodes.
-   **Rewards:** 
    -   **Contributor:** Earn KAWAI tokens via "Mining" (85-90% split). Rewards follow a **Halving Schedule** based on total supply (100 -> 50 -> 25 -> 12 per 1M tokens).
    -   **User:** Earn 5% cashback on every AI request (use-to-earn) + 1-5% deposit cashback (tiered).
    -   **Affiliator:** Earn 5% commission from referrals' mining rewards (lifetime passive income).
    -   **Developer:** Earn 5% from mining rewards (distributed to treasury pool).
    -   **Holder:** Earn 100% of Platform Revenue (USDT) proportional to KAWAI holdings.
-   **Liquidity:** No Initial LP. Value follows the *Hold-to-Earn* utility.
-   **Details:** See tokenomics section below for full analysis.

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

👉 **See Full Technical Details in the sections below**

---

## 📂 Project Structure & Usage

Veridium consists of two primary components designed for different user roles:

### 1. User Client (Desktop App)
- **Location:** `main.go` (Root)
- **Tech:** [Wails v3](https://wails.io/) + React.
- **Description:** The main interface for AI consumers. It provides a premium desktop experience for chat, knowledge base management, and a Web3 dashboard for USDT deposits and token management.
- **Current Mode:** Local LLM inference with automatic reward recording to treasury pool.
- **How to Run:**
  Refer to the `Makefile` for various development and build commands.

### 2. Contributor Client (CLI)
- **Location:** `cmd/contributor/main.go`
- **Tech:** Go CLI.
- **Description:** The "Miner" application. It wraps `llama.cpp` to provide compute power to the network and earn KAWAI tokens.
- **Status:** Ready for distributed network phase (future).

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
        - [x] **10 USDT with Referral**: Enhanced bonus for referred users (+100%).
        - [x] **Anti-Abuse**: Dual-layer protection (Wallet Address + Machine ID).
        - [x] **Atomic Claims**: Prevents race conditions during trial claiming.
    - [x] **Referral System:**
        - [x] **Viral Growth**: Users earn 5 USDT per successful referral.
        - [x] **Unique Codes**: 6-character alphanumeric referral codes.
        - [x] **URL Detection**: Auto-apply referral from `?ref=CODE` parameter.
        - [x] **Dashboard**: Track referrals, earnings, and share links.
        - [x] **Self-Referral Prevention**: Anti-abuse mechanism.
    - [x] **Deposit Cashback System:**
        - [x] **Tiered Rewards**: 1%-5% KAWAI cashback based on deposit amount.
        - [x] **First-Time Bonus**: 5% cashback for first deposit (overrides base rate).
        - [x] **Tier Caps**: 5K-20K KAWAI per deposit to prevent abuse.
        - [x] **Weekly Claims**: Off-chain accumulation, weekly Merkle settlement.
        - [x] **Batch Claims**: Claim multiple periods in one transaction.
        - [x] **200M Allocation**: 20% of max supply (~3 year runway).
        - [x] **Unlimited Referrals**: No cap on earnings.
    - [x] **Mining Reward Distribution System:**
        - [x] **Referral-Based Splits:** 85/5/5/5 (Contributor/Developer/User/Affiliator) for referral users, 90/5/5 for non-referral.
        - [x] **Job Tracking:** Per-job reward recording with full split details.
        - [x] **9-Field Merkle Trees:** Gas-efficient weekly settlement with flexible developer addresses.
        - [x] **Smart Contract:** `MiningRewardDistributor.sol` deployed to testnet.
        - [x] **Settlement Command:** `mining-settlement` CLI for weekly Merkle generation and upload.
        - [x] **Claim Integration:** Frontend bindings for `ClaimMiningReward()`.
        - [x] Random treasury address selection for fair distribution.
        - [x] Phase detection (Mining vs USDT mode).
        - [x] Halving logic (100→50→25→12 KAWAI per 1M tokens).
        - [x] Async execution (no user latency impact).
    - [x] **Revenue Share UI (Hold-to-Earn):** ✨ **NEW**
        - [x] **Real-Time Blockchain Data:** Fetches KAWAI balance and total supply from chain.
        - [x] **Share Percentage Calculator:** Displays user's proportional ownership (balance/totalSupply × 100).
        - [x] **Claimable USDT Table:** Shows weekly revenue distributions with Merkle proofs.
        - [x] **Phase Indicators:** Clear warnings for Phase 1 (Mining) vs Phase 2 (Revenue Sharing).
        - [x] **Error Handling:** Blockchain failure warnings, robust SetString validation.
        - [x] **Gas Estimation:** Pre-claim gas cost calculation for user transparency.
        - [x] **Comprehensive Docs:** 487-line guide with formulas, examples, and FAQs.
        - [x] **Type-Safe Implementation:** Full TypeScript types, no 'any' types.
        - [x] **Backend Support:** `GetRevenueShareStats()` service for balance/supply queries.
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

### Phase 3: Growth & Expansion ✅
**Focus:** Enhanced features and user experience.
- [x] **Frontend (Web Dashboard):**
    - [x] Initialize Vite + React project.
    - [x] Integrate @lobehub/ui for UI components.
    - [x] Complete P2P marketplace interface with real-time trading capabilities.
    - [x] Wallet integration with multi-wallet support and secure authentication.
    - [x] Market data visualization with statistics and price tracking.
    - [x] Order management system with creation, cancellation, and history tracking.
- [x] **AI Image Generation (Gemini API):** ✨ **NEW**
    - [x] Migration from Pollinations AI to Google Gemini API.
    - [x] Support for 2 models: `gemini-2.5-flash-image` (fast, 1K) & `gemini-3-pro-image-preview` (HD/4K).
    - [x] 10 aspect ratios: 1:1, 2:3, 3:2, 3:4, 4:3, 4:5, 5:4, 9:16, 16:9, 21:9.
    - [x] Quality tiers: Standard (1K), HD (2K), Ultra (4K).
    - [x] API key pool management with automatic rotation (5 keys).
    - [x] Thread-safe concurrent image generation.
    - [x] Context timeout protection (120s).
    - [x] Priority-based model selection (respects user intent).
    - [x] Comprehensive documentation and quick start guide.

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

## 🎨 AI Image Generation (NEW)

**Gemini API Integration** - Production-ready image generation with Google's latest models.

### Features
- **Two Models Available:**
  - `gemini-2.5-flash-image` (Nano Banana) - Fast generation, 1024px resolution
  - `gemini-3-pro-image-preview` (Nano Banana Pro) - High quality, up to 4K resolution
  
- **Flexible Configuration:**
  - 10 aspect ratios supported (1:1, 16:9, 9:16, 4:3, 3:4, etc.)
  - Quality tiers: Standard (1K), HD (2K), Ultra (4K)
  - Priority-based model selection (explicit model > quality > default)
  
- **Production Features:**
  - Thread-safe concurrent generation
  - API key pool with automatic rotation (5 keys)
  - Context timeout protection (120 seconds)
  - Comprehensive error handling
  - Enhanced logging with `[Gemini]` prefix

### Documentation
- 📖 [Full Implementation Guide](GEMINI_IMAGE_GENERATION.md)
- 📖 [Migration Documentation](MIGRATION_POLLINATIONS_TO_GEMINI.md)
- 📖 [Quick Start Guide](docs/IMAGE_GENERATION_QUICK_START.md)
- 📖 [Code Review & Fixes](CODE_REVIEW_FIXES.md)

### Usage Example
```go
opts := image.GenerationOptions{
    Prompt:      "A futuristic city at sunset",
    AspectRatio: "16:9",
    Quality:     "hd",
    Model:       "gemini-3-pro",
}
```

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
| [`MIGRATION_POLLINATIONS_TO_GEMINI.md`](MIGRATION_POLLINATIONS_TO_GEMINI.md) | Migration from Pollinations to Gemini | Understanding the migration |
| [`CODE_REVIEW_FIXES.md`](CODE_REVIEW_FIXES.md) | Code review feedback & fixes | Understanding implementation decisions |