# State of Idea: Lean DeAI Network on Monad
**Status:** Active | **Base:** `analisa.md` | **Date:** 2025-12-20

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
*   **Total Mining Reward:** Dihitung berdasarkan **Total Token Processed** (Volume-based: Input + Output Tokens).
    *   **Rate:** 100 KAWAI per 1 Juta Token (1M Tokens).
    *   *Rumus:* `Reward = (Total_Token / 1,000,000) * 100`
*   **Split Ratio (Pembagian):**
    *   **70%** -> Masuk ke Wallet **Contributor** (Pemilik Hardware).
    *   **30%** -> Masuk ke Wallet **Admin/Dev** (Biaya Pengembangan Platform).
*   **Mekanisme Klaim:** Contributor mengklaim porsi 70% mereka menggunakan sistem **Merkle Airdrop** mingguan.

### B. Profit Sharing & Two-Phase Economic Model

User membayar layanan menggunakan **USDT**. Pendapatan ini diklasifikasikan menjadi 2 fase:

#### Fase 1: Mining Era (Supply < 1 Miliar KAWAI)
| Pihak | Pendapatan |
|---|---|
| Contributor | **70% KAWAI** (Mining) |
| Admin/Dev | **30% KAWAI** (Mining) |
| KAWAI Holder | **100% Revenue USDT** |

*Contributor dibayar dengan Token baru (Inflasi). Holder mendapatkan seluruh Revenue USDT.*

#### Fase 2: Post-Mining Era (Supply = 1 Miliar / Max Cap)

Mining berhenti. Contributor dibayar **USDT** berdasarkan volume pekerjaan, bukan persentase Revenue.

**Rumus Biaya (Cost):**
*   **Cost Rate:** `COST_RATE_PER_MILLION` (Default: $1 USDT per 1 Juta Token). *Dapat disesuaikan via Environment Variable.*
*   **Contributor Cost:** `(Total_Token / 1,000,000) * $1 * 70%`
*   **Admin Cost:** `(Total_Token / 1,000,000) * $1 * 30%`
*   **Profit:** `Revenue - Contributor Cost - Admin Cost`

**Contoh Perhitungan:**
*   User membayar: **$100 USDT** untuk job yang memproses 500.000 Token.
*   Cost Rate: 500.000 / 1.000.000 * $1 = **$0.50**
*   Contributor: $0.50 * 70% = **$0.35**
*   Admin: $0.50 * 30% = **$0.15**
*   **Holder Profit:** $100 - $0.35 - $0.15 = **$99.50**

| Pihak | Pendapatan |
|---|---|
| Contributor | Cost-based (USDT) |
| Admin/Dev | Cost-based (USDT) |
| KAWAI Holder | **Profit USDT** (Revenue - Total Cost) |

**Dividen Mingguan (Kedua Fase):**
*   Total USDT dalam seminggu dikumpulkan di `PaymentVault`.
*   Sistem melakukan Snapshot kepemilikan token KAWAI.
*   Sisa USDT (Revenue di F1, Profit di F2) dibagikan proporsional ke Holder.

#### Phase Transition Detection (Deteksi Transisi Fase)

Sistem secara otomatis mendeteksi kapan harus beralih dari Fase 1 ke Fase 2:

*Lokasi Code:* `pkg/blockchain/client.go` -> `GetRewardMode()`

```go
// Pseudo-code:
totalSupply := token.TotalSupply()
maxSupply   := 1_000_000_000 * 1e18  // 1 Billion dengan 18 decimals

if totalSupply >= maxSupply {
    return ModeUSDT  // Fase 2: Bayar USDT
} else {
    return ModeMining  // Fase 1: Mining KAWAI
}
```

**Contoh Skenario Transisi:**
1.  Total Supply saat ini: **999,999,500 KAWAI**.
2.  Contributor menyelesaikan job yang menghasilkan **600 KAWAI** reward.
3.  Middleware memanggil `GetRewardMode()` -> `ModeMining`.
4.  Sistem mint **500 KAWAI** (sampai Max Supply), sisanya **TIDAK** di-mint.
5.  Setelah ini, semua job berikutnya akan menggunakan `ModeUSDT`.


### C. Liquidity Strategy (Strategi Likuiditas)
Karena tidak ada Modal Tim untuk membuat LP di PancakeSwap pada Hari-1:
*   **Fase 1 (Bootstrap): Internal P2P Market (OTC).**
    *   Membangun platform "Bulletin Board" sederhana (mirip Tokopedia/eBay untuk token).
    *   Kontributor yang butuh uang tunai menjual token mereka ke Investor Baru.
    *   **Teknologi:** Smart Contract Escrow Aman (Atomic Swap).
    *   **Harga:** Terbentuk alami oleh Supply & Demand (Market Forces), bukan kurva AMM.
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
3.  **OTC/Escrow Contract:**
    *   Kontrak sederhana untuk memfasilitasi jual-beli P2P (Token <-> USDT) tanpa Slippage AMM.
    *   Fitur: Create Order, Cancel Order, Buy Order.

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

1.  **Pemicu:** Contributor menyelesaikan request LLM -> Server memanggil fungsi `RecordJobReward`.
2.  **Cek Pemilik (Admin Check):**
    *   `IF Contributor_Address == Admin_Address`: 
        *   Contributor (Admin) mendapatkan **100%** Reward langsung ke saldo database-nya.
    *   `ELSE` (Public Contributor):
        *   Contributor mendapatkan **70%** Reward (masuk saldo contributor).
        *   Admin mendapatkan **30%** Fee (masuk saldo admin).
3.  **Akumulasi:** Saldo diupdate seketika di Database (KV Store).
4.  **Mingguan (Weekly):**
    *   Admin menjalankan script `snapshot`.
    *   Script hanya membaca total saldo akhir (tanpa rumus lagi) -> Generate Merkle Root -> Upload ke Blockchain.

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

1.  **Development (MVP):**
    *   Buat Smart Contract Token & Escrow.
    *   Buat Client Script `llama.cpp` sederhana.
    *   Buat Server Middleware (Golang) untuk manajemen job.
2.  **Deployment:**
    *   Deploy kontrak ke Monad Testnet/Mainnet.
    *   Rilis website sederhana untuk Dashboard Worker & P2P Market.
3.  **Launch:**
    *   Undang kontributor awal (Alpha Testers).
    *   Mulai siklus: Kerja -> Poin -> Mingguan Token Dist + USDT Dividen.

---
*Dokumen ini adalah titik acuan untuk pengembangan selanjutnya. Ide yang tidak tercantum di sini dianggap diarsipkan/tidak prioritas.*
