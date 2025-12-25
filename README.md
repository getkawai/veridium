# Kawai DeAI Network (Monad)

**Decentralized AI Compute Network on Monad.**
A "Lean Startup" approach to DePIN, leveraging consumer-grade GPUs for `llama.cpp` inference with a sustainable tokenomic model.

---

## 🚀 Core Concept

-   **Service:** Low-cost LLM Inference API (compatible with OpenAI format).
-   **Contributors:** Gamers & Devs running `llama.cpp` nodes.
-   **Rewards:** 
    -   **Contributor:** Earn KAWAI tokens via "Mining" (70% to Contributor, 30% to Dev).
    -   **Holder:** Earn 100% of Platform Revenue (USDT) proportional to KAWAI holdings.
-   **Liquidity:** No Initial LP. Value follows the *Hold-to-Earn* utility.
-   **Details:** See [Concept Document](current_concept.md) for full analysis.

---

## ⚙️ How It Works (Summary)

We use a **Hybrid Model** (Off-Chain Accumulation + On-Chain Settlement) to minimize gas fees.
*   **Real-Time:** Rewards are calculated instantly (Usage-based) and split 70/30.
*   **Weekly:** A compressed **Merkle Tree** is uploaded to **Monad** for cheap claiming.

👉 **[See Full Technical Details in Concept Document](current_concept.md#d-mekanisme-teknis-hybrid-how-it-works)**

---

## 🗺️ Implementation Plan (Roadmap)

This roadmap outlines the path from "Zero" to a fully functional decentralized network.

### Phase 1: Foundation (Current Status)
**Focus:** Infrastructure Setup & Smart Contract Development.
- [x] **Project Initialization:**
    - [x] Go Backend Structure (Middleware).
    - [x] Hardhat Environment (Contracts).
- [x] **Auto-Binding Setup:**
    - [x] Makefile for `abigen` (Solidity -> Go).
- [x] **Smart Contracts (MVP):**
    - [x] `KawaiToken.sol`: Standard ERC20 with AccessControl (Mint/Burn).
    - [x] `Escrow.sol`: Simple P2P OTC Market (Orders, Buy, Cancel).
    - [x] `PaymentVault.sol`: Prepaid USDT Deposit system for Consumers.
    - [x] `MerkleDistributor.sol`: Gas-efficient reward claiming system.
- [x] **Middleware (Go):**
    - [x] `pkg/blockchain`: Service to interact with **Monad** (Listen events, Send TXs).
    - [x] `pkg/store`: Persistent storage using Cloudflare Workers KV.
    - [x] `pkg/merkle`: Merkle Tree generation logic.

### Phase 2: The "Lean" Launch (Internal Market)
**Focus:** Economic bootstrapping without initial LP.
- [ ] **Contributor Client (Go):** 
    - [x] **Wallet Setup:** Generate/Load Contributor Identity (Private Key).
    - [x] Wrapper for `llama.cpp` server (Process Management).
    - [x] Heartbeat system (Proof of Availability).
- [ ] **Verification System (Middleware):**
    - [ ] "Gold Standard" verification (random trap questions) to prevent cheat nodes.
- [ ] **Consumer API (User Client):**
    - [ ] Web3 Dashboard: Login via Wallet to manage API Keys.
    - [x] OpenAI-compatible `/v1/chat/completions` proxy (Base foundation).
    - [ ] Real-time credit deduction system.
- [ ] **Internal P2P Marketplace (Web):**
    - [ ] UI for Contributors to list their Token rewards.
    - [ ] UI for Investors to buy Tokens with USDT.
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

### Phase 3: Community Liquidity (Growth)
**Focus:** Transition to decentralized liquidity.
- [ ] **Liquidity Incentives:**
    - [ ] Launch "LP Mining" program.
    - [ ] Encourage stakeholders to provide LP on PancakeSwap using their Dividends.
- [ ] **Consumer API (User Client):**
    - [ ] Web3 Dashboard: Login via Wallet to manage API Keys.
    - [x] OpenAI-compatible `/v1/chat/completions` proxy (Base).
    - [ ] Real-time credit deduction system.
- [ ] **Frontend (Web Dashboard):**
    - [x] Initialize Vite + React project.
    - [x] Integrate Ant Design Web3 for Wallet connection.
    - [ ] Implement USDT Deposit UI (PaymentVault integration).

---

## 🌐 Monad Testnet Deployment

> **Status:** Pending Gas (Faucet Cooldown)

After deployment, contract addresses will be listed here:

| Contract | Address | Description |
|---|---|---|
| MockUSDT | `TBD` | Test USDT for simulating revenue |
| KawaiToken | `TBD` | Native token (1B Max Supply, Fair Launch) |
| KAWAI_Distributor | `TBD` | Merkle distributor for mining rewards |
| USDT_Distributor | `TBD` | Merkle distributor for profit sharing |
| PaymentVault | `TBD` | User USDT deposit vault |
| OTCMarket | `TBD` | P2P trading escrow |

**Network Info:**
- **Chain ID:** 10143
- **RPC:** `https://testnet-rpc.monad.xyz`
- **Explorer:** `https://testnet.monad.xyz` (if available)

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

# Admin Configuration
ADMIN_ADDRESS=0x...             # Admin wallet for fee collection
ADMIN_PRIVATE_KEY=0x...         # For signing Merkle Root updates

# Cloudflare KV (Worker Data Storage)
CF_ACCOUNT_ID=...
CF_API_TOKEN=...
CF_KV_NAMESPACE_ID=...
```

---

## 🛠️ Tech Stack

-   **Smart Contracts:** Solidity, Foundry, OpenZeppelin.
-   **Backend (Middleware):** Go (Golang), `net/http`, `go-ethereum`, `cloudflare-go`.
-   **Database:** Cloudflare Workers KV.
-   **Frontend:** React, Vite, Ant Design Web3.
-   **Contributor Node:** Go (Golang), `llama.cpp`.
-   **Blockchain:** Monad (EVM-compatible).
-   **Network Toolkit:** `pkg/jarvis` (Multi-chain support incl. Monad).