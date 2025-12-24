# State of Idea: Lean DeAI Network on BSC
**Status:** Active | **Base:** `analisa.md` | **Date:** 2025-12-20

Dokumen ini merangkum status terakhir dari brainstorming project DePIN AI (Decentralized AI) yang berjalan di jaringan BSC dengan pendekatan *Lean Startup* (Minim Modal).

## 1. Core Concept (Inti Konsep)
Membangun jaringan komputasi AI terdesentralisasi (**DePIN**) di mana kontributor menyewakan GPU mereka untuk menjalankan model LLM (`llama.cpp`), dan dibayar menggunakan kombinasi Token Native & Subsidi Listrik.

*   **Network:** BNB Smart Chain (BSC) — dipilih karena biaya gas murah & ekosistem matang.
*   **Target User:** Pengguna yang membutuhkan API LLM murah/gratis.
*   **Target Worker:** Gamer/Developer dengan GPU menganggur (Consumer Grade).

## 2. Business & Economic Model (Tokenomics)

Strategi utama adalah **"No Initial Liquidity Pool"** untuk menghemat modal awal (Seed Capital $0 untuk LP), digantikan dengan ekonomi berbasis *Real Yield* dan Pasar Internal.

### A. Worker Rewards (Insentif Kontributor)
Kontributor (Node Runners) mendapatkan dua jenis "gaji":
1.  **Native Token (Principal Reward):**
    *   Diberikan sebagai bukti kerja (*Proof of Work/Availability*).
    *   Berfungsi sebagai "Saham" untuk klaim dividen masa depan.
    *   Dihitung berdasarkan Uptime & Tier Hardware.
2.  **USDT (Subsidi Operasional):**
    *   Diberikan **hanya sebagian kecil** (misal 10-20%) untuk menutup biaya listrik harian.
    *   Sumber dana: Fee dari Liquidity Pool (masa depan) atau porsi kecil Revenue.

### B. Holder Rewards (Insentif Investor)
Mengapa orang mau memegang (Hold) token ini jika tidak bisa dijual di PancakeSwap (awal)?
*   **Weekly USDT Dividend:** Revenue nyata dari penggunaan jasa AI dibagikan setiap minggu kepada pemegang token.
*   **Mekanisme:** *Snapshot* mingguan -> *Batch Transfer* (Disperse) USDT ke wallet holders.

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

### A. Blockchain Layer (BSC)
1.  **Token Contract (ERC20/BEP20):**
    *   Standar OpenZeppelin (Aman, Audit-free).
    *   Fitur: Mintable (untuk reward pool), Burnable (untuk deflasi).
2.  **OTC/Escrow Contract:**
    *   Kontrak sederhana untuk memfasilitasi jual-beli P2P (Token <-> USDT) tanpa Slippage AMM.
    *   Fitur: Create Order, Cancel Order, Buy Order.

### B. Off-Chain Layer (Middleware & Nodes)
1.  **AI Nodes (Workers):**
    *   Menjalankan skrip Python yang membungkus `llama.cpp`.
    *   Fungsi: `Pull Job` -> `Inference` -> `Push Result`.
2.  **Central Authority (Middleware Server):**
    *   **Job Dispatcher:** Menerima request user (API) -> Kirim ke Node.
    *   **Proof of Availability:** Melakukan "Ping" berkala ke Node untuk memastikan uptime.
    *   **Verifikasi (Anti-Cheat):** Menggunakan metode "Gold Standard" (menyisipkan pertanyaan jebakan yang jawabannya sudah diketahui) untuk memvalidasi kejujuran Node.
    *   **Accounting:** Mencatat poin kontributor di Database (SQL) sebelum dicarikan menjadi Token (Minting) setiap minggu.

## 4. Roadmap Tahap Awal (Immediate Action Plan)

1.  **Development (MVP):**
    *   Buat Smart Contract Token & Escrow.
    *   Buat Client Script `llama.cpp` sederhana.
    *   Buat Server Middleware (Golang) untuk manajemen job.
2.  **Deployment:**
    *   Deploy kontrak ke BSC Mainnet.
    *   Rilis website sederhana untuk Dashboard Worker & P2P Market.
3.  **Launch:**
    *   Undang kontributor awal (Alpha Testers).
    *   Mulai siklus: Kerja -> Poin -> Mingguan Token Dist + USDT Dividen.

---
*Dokumen ini adalah titik acuan untuk pengembangan selanjutnya. Ide yang tidak tercantum di sini dianggap diarsipkan/tidak prioritas.*
