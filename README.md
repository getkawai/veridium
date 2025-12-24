# Kawai DeAI Network (Monad)

**Decentralized AI Compute Network on Monad.**
A "Lean Startup" approach to DePIN, leveraging consumer-grade GPUs for `llama.cpp` inference with a sustainable tokenomic model.

---

## 🚀 Core Concept

-   **Service:** Low-cost LLM Inference API (compatible with OpenAI format).
-   **Workers:** Gamers & Devs running `llama.cpp` nodes.
-   **Rewards:** 
    -   **Worker:** Earn KAWAI tokens via "Mining" (70% to Worker, 30% to Dev).
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
- [ ] **Worker Client (Go):** 
    - [x] **Wallet Setup:** Generate/Load Worker Identity (Private Key).
    - [x] Wrapper for `llama.cpp` server (Process Management).
    - [x] Heartbeat system (Proof of Availability).
- [ ] **Verification System (Middleware):**
    - [ ] "Gold Standard" verification (random trap questions) to prevent cheat nodes.
- [ ] **Consumer API (User Client):**
    - [ ] Web3 Dashboard: Login via Wallet to manage API Keys.
    - [x] OpenAI-compatible `/v1/chat/completions` proxy (Base foundation).
    - [ ] Real-time credit deduction system.
- [ ] **Internal P2P Marketplace (Web):**
    - [ ] UI for Workers to list their Token rewards.
    - [ ] UI for Investors to buy Tokens with USDT.
- [x] **Administration Scripts (Go):**
    - [x] Weekly Snapshot & Dividend Calculator.
    - [x] Worker Audit & Ban system.
- [ ] **Dividend & Reward System (Two-Phase Model):**
    - [x] Weekly Snapshot script.
    - [x] Merkle Airdrop implementation (Pull model) for KAWAI.
    - [ ] **Phase 1 -> Phase 2 Transition Logic:**
        - [ ] Detect when `totalSupply() == MAX_SUPPLY` in Smart Contract / Middleware.
        - [ ] Switch payment mode from KAWAI mining to USDT cost-based.
    - [ ] **Phase 2 USDT Payouts:**
        - [ ] Implement `COST_RATE_PER_MILLION` as configurable Environment Variable.
        - [ ] Calculate Worker/Admin USDT cost per job in `RecordJobReward`.
        - [ ] Distribute remaining USDT Profit to KAWAI Holders weekly.

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

## 🛠️ Tech Stack

-   **Smart Contracts:** Solidity, Hardhat, OpenZeppelin.
-   **Backend (Middleware):** Go (Golang), `net/http`, `go-ethereum`, `cloudflare-go`.
-   **Database:** Cloudflare Workers KV.
-   **Frontend:** React, Vite, Ant Design Web3.
-   **Worker Node:** Go (Golang), `llama.cpp`.
-   **Blockchain:** Monad (EVM-compatible).