# KawaiBSC DeAI Network

**Decentralized AI Compute Network on BNB Smart Chain.**
A "Lean Startup" approach to DePIN, leveraging consumer-grade GPUs for `llama.cpp` inference with a sustainable tokenomic model.

---

## 🚀 Core Concept

-   **Service:** Low-cost LLM Inference API (compatible with OpenAI format).
-   **Workers:** Gamers & Devs running `llama.cpp` nodes.
-   **Rewards:** Hybrid model (Native Token for "Proof of Availability" + USDT for "Operational Subsidy").
-   **Liquidity:** No Initial LP. Uses **Weekly USDT Dividends** and an **Internal P2P Market** to bootstrap value.
-   **Details:** See [Concept Document](current_concept.md) for full analysis.

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
- [x] **Middleware (Go):**
    - [x] `pkg/blockchain`: Service to interact with BSC (Listen events, Send TXs).
    - [x] `pkg/store`: Persistent storage using Cloudflare Workers KV.

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
- [ ] **Dividend System:**
    - [x] Weekly Snapshot script.
    - [ ] Batch Transfer (Disperse) implementation.

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
-   **Blockchain:** BNB Smart Chain (BSC).