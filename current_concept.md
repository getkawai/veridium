# State of Idea: Lean DeAI Network on Monad
**Status:** Active | **Base:** `analisa.md` | **Date:** 2026-01-01

Dokumen ini merangkum status terakhir dari brainstorming project DePIN AI (Decentralized AI) yang berjalan di jaringan **Monad** dengan pendekatan *Lean Startup* (Minim Modal).

## 1. Core Concept (Inti Konsep)
Membangun jaringan komputasi AI terdesentralisasi (**DePIN**) di mana kontributor menyewakan GPU mereka untuk menjalankan model LLM (`llama.cpp`), dan dibayar menggunakan **Token Native (KAWAI)** serta mendapatkan **Bagi Hasil USDT (Real Yield)**.

*   **Network:** **Monad** — dipilih karena throughput tinggi, biaya gas sangat murah & EVM-compatible.
*   **Target User:** Pengguna yang membutuhkan API LLM murah/gratis.
*   **Target Contributor:** Gamer/Developer dengan GPU menganggur (Consumer Grade).

## 2. Business & Economic Model (Tokenomics)

Strategi utama adalah **"No Initial Liquidity Pool"** untuk menghemat modal awal (Seed Capital $0 untuk LP), digantikan dengan ekonomi berbasis *Real Yield* dan Pasar Internal.

### A. Mining & Token Generation (Penciptaan Token)
*   **Fair Launch Policy:** Tidak ada Pre-Mint, tidak ada alokasi Investor/VC. Supply awal = **0**.
*   **Max Supply (Hard Cap):** **1 Miliar (1,000,000,000) KAWAI**. Setelah angka ini tercapai, mining berhenti.
*   **Mining Only:** Satu-satunya cara token baru tercipta (Minting) adalah melalui **Proof of Computation** (Menyelesaikan unit pekerjaan LLM).
*   **Dynamic Mining Difficulty (Halving):** Reward KAWAI berkurang secara berkala berdasarkan total supply yang sudah di-mint untuk menjaga kelangkaan:
    *   **Phase 1A (< 500M KAWAI):** 100 KAWAI per 1M Token.
    *   **Phase 1B (500M - 750M KAWAI):** 50 KAWAI per 1M Token (Halving 1).
    *   **Phase 1C (750M - 875M KAWAI):** 25 KAWAI per 1M Token (Halving 2).
    *   **Phase 1D (> 875M KAWAI):** 12 KAWAI per 1M Token (Halving 3).
*   **Implementation Note:** Logika ini dieksekusi secara otomatis di fungsi `RecordJobReward` (`pkg/store/contributor.go`) yang mengambil `totalSupply` terkini dari blockchain melalui `SupplyQuerier`.
*   **Split Ratio (Pembagian):**
    
    **Untuk User dengan Referral (Total: 100%):**
    *   **85%** -> Masuk ke Wallet **Contributor** (Pemilik Hardware/Miner).
    *   **5%** -> Masuk ke Wallet **Developer** (Biaya Pengembangan Platform).
    *   **5%** -> Masuk ke Wallet **User** (Cashback untuk requester).
    *   **5%** -> Masuk ke Wallet **Affiliator** (Referrer yang mengajak user tersebut).
    
    **Untuk User tanpa Referral (Total: 100%):**
    *   **90%** -> Masuk ke Wallet **Contributor** (Pemilik Hardware/Miner).
    *   **5%** -> Masuk ke Wallet **Developer** (Biaya Pengembangan Platform).
    *   **5%** -> Masuk ke Wallet **User** (Cashback untuk requester).
    
    *   **Admin Selection:** Admin address dipilih secara random dari treasury pool (`internal/constant/treasury.go`) setiap job.
    *   **Auto-Registration:** Admin accounts otomatis dibuat/diupdate jika belum terdaftar sebagai contributor.
    *   **Use-to-Earn:** User mendapat 5% cashback dari setiap request mereka, mendorong usage.
    *   **Lifetime Commission:** Affiliator mendapat 5% dari mining reward referral mereka (contributor mengorbankan 5% untuk growth).
*   **Mekanisme Klaim:** Contributor mengklaim porsi 85-90% mereka, User mengklaim 5% cashback, dan Affiliator mengklaim 5% komisi menggunakan sistem **Merkle Airdrop** mingguan.

### B. Profit Sharing & Two-Phase Economic Model

User membayar layanan menggunakan **USDT**. Pendapatan ini diklasifikasikan menjadi 2 fase:

#### Fase 1: Mining Era (Supply < 1 Miliar KAWAI)

**User dengan Referral (Total: 100%):**
| Pihak | Pendapatan |
|---|---|
| Contributor | **85% KAWAI** (Mining) |
| Developer | **5% KAWAI** (Mining) |
| User | **5% KAWAI** (Cashback) |
| Affiliator | **5% KAWAI** (Commission) |
| KAWAI Holder | **100% Revenue USDT** |

**User tanpa Referral (Total: 100%):**
| Pihak | Pendapatan |
|---|---|
| Contributor | **90% KAWAI** (Mining) |
| Developer | **5% KAWAI** (Mining) |
| User | **5% KAWAI** (Cashback) |
| KAWAI Holder | **100% Revenue USDT** |

*Contributor dibayar dengan Token baru (Inflasi). User mendapat 5% cashback untuk mendorong usage. Affiliator mendapat komisi lifetime dari mining referral mereka (contributor mengorbankan 5%). Holder mendapatkan seluruh Revenue USDT.*

#### Fase 2: Post-Mining Era (Supply = 1 Miliar / Max Cap)

Mining berhenti. Contributor dibayar **USDT** berdasarkan volume pekerjaan, bukan persentase Revenue.

**Rumus Biaya (Cost):**
*   **Cost Rate:** `COST_RATE_PER_MILLION` (Default: $1 USDT per 1 Juta Token). *Dapat disesuaikan via Environment Variable.*

**User dengan Referral (Total Cost: 100%):**
*   **Contributor Cost:** `(Total_Token / 1,000,000) * $1 * 85%`
*   **Developer Cost:** `(Total_Token / 1,000,000) * $1 * 5%`
*   **User Cashback:** `(Total_Token / 1,000,000) * $1 * 5%`
*   **Affiliator Cost:** `(Total_Token / 1,000,000) * $1 * 5%`
*   **Profit:** `Revenue - Contributor Cost - Developer Cost - User Cashback - Affiliator Cost`

**User tanpa Referral (Total Cost: 100%):**
*   **Contributor Cost:** `(Total_Token / 1,000,000) * $1 * 90%`
*   **Developer Cost:** `(Total_Token / 1,000,000) * $1 * 5%`
*   **User Cashback:** `(Total_Token / 1,000,000) * $1 * 5%`
*   **Profit:** `Revenue - Contributor Cost - Developer Cost - User Cashback`

**Contoh Perhitungan (User dengan Referral):**
*   User membayar: **$100 USDT** untuk job yang memproses 500.000 Token.
*   Cost Rate: 500.000 / 1.000.000 * $1 = **$0.50**
*   Contributor: $0.50 * 85% = **$0.425**
*   Developer: $0.50 * 5% = **$0.025**
*   User Cashback: $0.50 * 5% = **$0.025**
*   Affiliator: $0.50 * 5% = **$0.025**
*   **Holder Profit:** $100 - $0.425 - $0.025 - $0.025 - $0.025 = **$99.50**

| Pihak | Pendapatan |
|---|---|
| Contributor | Cost-based (USDT) |
| Developer | Cost-based (USDT) |
| KAWAI Holder | **Profit USDT** (Revenue - Total Cost) |

**Dividen Mingguan (Kedua Fase):**
*   Total USDT dalam seminggu dikumpulkan di `PaymentVault`.
*   Sistem melakukan Snapshot kepemilikan token KAWAI.
*   Sisa USDT (Revenue di F1, Profit di F2) dibagikan proporsional ke Holder.

#### Phase Transition Detection (Deteksi Transisi Fase)

Sistem secara otomatis mendeteksi kapan harus beralih dari Fase 1 ke Fase 2:

*Lokasi Code:* `pkg/store/contributor.go` -> `RecordJobReward()` (inline detection)

```go
// Simplified inline check:
mode := config.ModeMining
if s.supplyQuerier != nil {
    currentSupply, _ := s.supplyQuerier.GetTotalSupply(ctx)
    maxSupply, _ := s.supplyQuerier.GetMaxSupply(ctx)
    if currentSupply != nil && maxSupply != nil && currentSupply.Cmp(maxSupply) >= 0 {
        mode = config.ModeUSDT  // Max supply reached, switch to USDT
    }
}
```

**Key Features:**
*   **Simple & Direct:** Langsung cek `totalSupply >= maxSupply` tanpa wrapper functions.
*   **Blockchain Source:** `maxSupply` diambil langsung dari smart contract ABI (`KawaiToken.MAX_SUPPLY()`).
*   **Defensive:** Safe nil checks dengan fallback ke `ModeMining` jika blockchain unavailable.

**Contoh Skenario Transisi:**
1.  Total Supply saat ini: **999,999,500 KAWAI**.
2.  Contributor menyelesaikan job yang menghasilkan **600 KAWAI** reward.
3.  `RecordJobReward()` cek supply -> `ModeMining` (masih < 1B).
4.  Sistem mint **600 KAWAI** -> Total supply sekarang **1,000,000,100 KAWAI**.
5.  Job berikutnya akan otomatis detect `ModeUSDT` karena supply >= maxSupply.


### C. Liquidity Strategy (Strategi Likuiditas)
Karena tidak ada Modal Tim untuk membuat LP di PancakeSwap pada Hari-1:
*   **Fase 1 (Bootstrap): Internal P2P Market (OTC) ✅ UPGRADED v2.0**
    *   ✅ **Partial Fill Support** - Users dapat buy/sell sebagian order (flexibility maksimal).
    *   ✅ **Remaining Amount Tracking** - Real-time tracking sisa order yang belum terisi.
    *   ✅ Platform "Bulletin Board" lengkap dengan order book real-time dan trading interface.
    *   ✅ Kontributor dapat membuat sell order untuk menjual token KAWAI mereka.
    *   ✅ Investor dapat membeli token langsung dari order book dengan eksekusi instan atau partial.
    *   ✅ **Teknologi:** Smart Contract Escrow v2 dengan Atomic Swap + Partial Fill.
    *   ✅ **Harga:** Terbentuk alami oleh Supply & Demand (Market Forces) dengan order book.
    *   ✅ **UI Features:** Market statistics, price tracking, order history, trade history.
    *   ✅ **Real-time Updates:** Event-driven updates menggunakan Wails Events API.
    *   ✅ **Reconciliation Service:** Auto-sync blockchain ↔ KV store setiap 5 menit.
    *   ✅ **Offline Resilience:** Chunked event replay untuk catch-up setelah offline.
*   **Fase 2 (Growth): Community-Owned Liquidity.**
    *   Setelah profit stabil, stakeholder didorong untuk membuat LP sendiri di PancakeSwap.
    *   Insentif: LP Provider mendapatkan bonus Token tambahan + Trading Fees.

## 3. Technical Architecture (Arsitektur Teknis)

Sistem menggunakan pendekatan **Hybrid (On-Chain Settlement + Off-Chain Verification)** untuk menghindari gas fee mahal.

### A. Blockchain Layer (Monad)
1.  **Token Contract (ERC20/BEP20):**
    *   Standar OpenZeppelin (Aman, Audit-free).
    *   Fitur: Mintable (untuk reward pool), Burnable (untuk deflasi).
2.  **Merkle Distributor Contract:**
    *   Kontrak untuk distribusi reward massal dengan biaya gas murah.
    *   Admin hanya meng-upload "Root" (Hash) mingguan. Worker melakukan klaim dengan bukti kriptografi.
3.  **OTC/Escrow Contract ✅ UPGRADED v2.0 (Deployed 2025-12-31):**
    *   Kontrak lengkap untuk memfasilitasi jual-beli P2P (Token <-> USDT) tanpa Slippage AMM.
    *   ✅ **NEW:** Partial Fill Support - Buy/sell sebagian order dengan `buyOrderPartial()`.
    *   ✅ **NEW:** Remaining Amount Tracking - Field `remainingAmount` untuk tracking real-time.
    *   ✅ **NEW:** Efficient View Functions - Batch queries, pagination, filter by seller.
    *   ✅ **NEW:** Enhanced Events - `OrderPartiallyFilled` event untuk partial fills.
    *   ✅ Fitur: Create Order, Cancel Order, Buy Order (Full/Partial).
    *   ✅ Real-time market data dan order book management.
    *   ✅ Event emission untuk real-time UI updates.
    *   ✅ Comprehensive order and trade history tracking.
    *   ✅ 19 comprehensive tests passing.
    *   📍 **Address:** `0x5b1235038B2F05aC88b791A23814130710eFaaEa`

### B. Off-Chain Layer (Middleware & Nodes)
1.  **AI Nodes (Contributors):**
    *   Menjalankan skrip Python yang membungkus `llama.cpp`.
    *   Fungsi: `Pull Job` -> `Inference` -> `Push Result`.
2.  **Central Authority (Middleware Server):**
    *   **Job Dispatcher:** Menerima request user (API) -> Kirim ke Node.
    *   **Proof of Availability:** Melakukan "Ping" berkala ke Node untuk memastikan uptime.
    *   **Verifikasi (Anti-Cheat):** Menggunakan metode "Gold Standard" (menyisipkan pertanyaan jebakan yang jawabannya sudah diketahui) untuk memvalidasi kejujuran Node.
    *   **Accounting & Merkle Generator:** Mencatat poin -> Generate Merkle Tree -> Upload Root ke Blockchain -> Simpan Proof di KV Store untuk diklaim Contributor.

### C. Off-Chain Storage Architecture (Cloudflare KV)

Data off-chain disimpan di **Cloudflare Workers KV** dengan arsitektur **Multi-Namespace**:

| Namespace | Deskripsi | Key Format |
|-----------|-----------|------------|
| `contributors` | Data profil & saldo contributor | `{wallet_address}` |
| `proofs` | Bukti Merkle per periode | `{address}:{period_id}` |
| `settlements` | Metadata settlement | `{period_id}` |

👉 **Detail lengkap:** [`pkg/store/README.md`](pkg/store/README.md)

### D. Logic Implementasi Pembagian (Reward Algorithm)
*Lokasi Code:* `pkg/store/contributor.go` -> `RecordJobReward()`

Logika pembagian 70/30 dieksekusi secara **Real-Time (Per Job)** oleh Middleware saat job selesai:

1.  **Pemicu:** Contributor menyelesaikan request LLM -> Server memanggil fungsi `RecordJobReward(contributorAddress, tokenUsage)`.
2.  **Auto-Detection:**
    *   **Admin Selection:** Random admin address dipilih dari treasury pool (`constant.GetRandomTreasuryAddress()`).
    *   **Reward Mode:** Otomatis detect `ModeMining` vs `ModeUSDT` berdasarkan `totalSupply >= maxSupply`.
    *   **Halving Rate:** Otomatis adjust rate (100/50/25/12) berdasarkan supply thresholds.
3.  **Pembagian Reward:**
    *   Contributor mendapatkan **85% (referral) atau 90% (non-referral)** Reward (KAWAI atau USDT tergantung mode).
    *   Developer mendapatkan **5%** Fee.
    *   User mendapatkan **5%** Cashback (use-to-earn incentive).
    *   Affiliator mendapatkan **5%** Commission (jika user punya referrer, diambil dari porsi contributor).
    *   **Total:** Selalu 100%.
    *   **Auto-Registration:** Jika admin belum terdaftar, otomatis dibuat via `EnsureAdminExists()`.
4.  **Akumulasi:** Saldo diupdate seketika di Database (KV Store).
5.  **Mingguan (Weekly):**
    *   Admin menjalankan script `snapshot`.
    *   Script hanya membaca total saldo akhir (tanpa rumus lagi) -> Generate Merkle Root -> Upload ke Blockchain.

**CLI Tools:**
*   `make admin-register` - Bulk register semua treasury addresses sebagai admin contributors.
*   `make admin-register-dry` - Dry-run untuk preview tanpa eksekusi.

### E. Mekanisme Teknis Hybrid (How It Works)

Agar jaringan tetap "Lean" (hemat biaya), kami menggunakan model **Off-Chain Accumulation + On-Chain Settlement**.

#### 1. Senin - Sabtu: Akumulasi (Off-Chain)
*   **Aksi:** Contributor memproses job AI (LLM Inference).
*   **Pencatatan:** Poin kinerja dicatat di **Database Terpusat** (KV Store).
    *   *Code Ref:* `pkg/store/contributor.go` -> `SaveContributor()` & `UpdateHeartbeat()`
*   **Biaya:** $0 Gas Fees. Kecepatan instan.

#### 2. Minggu: Settlement (Weekly Batch)
*   **Kalkulasi:** Aturan **70/30 Split** diterapkan secara **Real-Time** oleh middleware setiap kali job selesai.
    *   *Code Ref:* `pkg/store/contributor.go` -> `RecordJobReward()`
*   **Kompresi:** Ribuan transaksi dikompres menjadi satu **Merkle Tree**.
*   **Blockchain:** Admin hanya mengirim **Merkle Root** (hash kecil) ke Smart Contract.
*   **Biaya:** 1 Transaksi murah per minggu (~$0.01 di Monad).

#### 3. Settlement Process (Atomic)
*   **Flow:** Snapshot → Generate Merkle Tree → Save Proofs → Reset Balances
*   **Safety:** Rollback otomatis jika ada kegagalan di tengah proses
*   **Resumable:** Settlement yang terganggu bisa dilanjutkan
    *   *Code Ref:* `pkg/store/settlement.go` -> `PerformSettlement()`

#### 4. Klaim (User Action)
*   **Interface:** Contributor menghubungkan Wallet ke Dashboard Web.
*   **Bukti:** Website mengambil "Bukti Kriptografis" (Proof) dari database.
*   **Withdraw:** Smart Contract memverifikasi Proof terhadap Root dan merilis token.
*   **Tracking:** Status klaim dilacak: `unclaimed` → `pending` → `confirmed`/`failed`
    *   *Code Ref:* `pkg/store/merkle.go` -> `MarkClaimPending()`, `ConfirmClaim()`

## 4. Roadmap Tahap Awal (Immediate Action Plan)

1.  **Development (MVP) ✅ COMPLETED:**
    *   ✅ Buat Smart Contract Token & Escrow.
    *   ✅ Buat Client Script `llama.cpp` sederhana.
    *   ✅ Buat Server Middleware (Golang) untuk manajemen job.
    *   ✅ **P2P Marketplace:** Complete trading interface dengan real-time updates.
2.  **Deployment ✅ UPGRADED (2026-01-01):**
    *   ✅ Deploy kontrak ke Monad Testnet/Mainnet.
    *   ✅ **NEW Contracts v2.0:** OTCMarket dengan Partial Fill Support deployed & verified.
    *   ✅ Rilis website lengkap untuk Dashboard Worker & P2P Market.
    *   ✅ **Live Contracts:** All contracts deployed dan terintegrasi dengan UI.
    *   ✅ **Backend Services:** Reconciliation, Event Replay, Rate Limiting implemented.
    *   ✅ **Reward System:** Auto-detection untuk admin selection, reward mode, dan halving rates.
    *   ✅ **Admin Management:** CLI tools untuk bulk admin registration dari treasury pool.
    *   📚 **Documentation:** DEPLOYMENT_SUMMARY.md, MARKETPLACE_UPGRADE_SUMMARY.md, CONTRACTS_WORKFLOW.md
3.  **Launch (READY):**
    *   Undang kontributor awal (Alpha Testers).
    *   Mulai siklus: Kerja -> Poin -> Mingguan Token Dist + USDT Dividen.
    *   **P2P Trading:** Users dapat langsung mulai trading KAWAI tokens.

### Next Phase: Community Growth
*   **Marketing & Onboarding:** Attract contributors dan investors ke platform.
*   **Liquidity Incentives:** Encourage trading activity di P2P marketplace.
*   **Performance Optimization:** Monitor dan optimize real-time trading experience.

---
*Dokumen ini adalah titik acuan untuk pengembangan selanjutnya. Ide yang tidak tercantum di sini dianggap diarsipkan/tidak prioritas.*
