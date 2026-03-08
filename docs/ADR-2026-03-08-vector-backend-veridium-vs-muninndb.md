# ADR-2026-03-08: Keputusan Migrasi Vector Backend (Veridium -> MuninnDB)

- Status: Accepted
- Tanggal: 2026-03-08
- Owner: Veridium Core
- Scope: RAG/vector retrieval untuk file/knowledge base (bukan conversation memory)

## Konteks

Arsitektur runtime Veridium saat ini adalah split:

1. Conversation memory memakai MuninnDB.
2. RAG/vector retrieval file memakai DuckDB + SQLite.

Implementasi sekarang:

1. Ingestion RAG menyimpan chunk ke SQLite dan embedding ke DuckDB.
2. Query RAG mencari top-K vector di DuckDB lalu hydrate konten chunk dari SQLite.
3. Filtering per file (fileIDs) dilakukan di layer Veridium setelah hasil vector didapat.

## Problem Statement

Perlu keputusan apakah vector retrieval Veridium saat ini sebaiknya diganti penuh ke MuninnDB untuk menyatukan backend memory + vector.

## Opsi

1. Opsi A: Pertahankan split sekarang (Muninn untuk memory percakapan, DuckDB untuk RAG).
2. Opsi B: Big-bang migration sekarang (ganti total vector retrieval ke MuninnDB).
3. Opsi C: Migrasi bertahap via `VectorStore` interface + adapter Muninn + A/B test.

## Kriteria Penilaian

Bobot total 100:

1. Retrieval quality (Recall@K, nDCG): 30
2. Latency p95 end-to-end: 20
3. Kompleksitas implementasi/migrasi: 15
4. Risiko regresi produksi: 15
5. Operasional (backup/recovery/observability): 10
6. Risiko legal/lisensi: 10

Skor 1-5 (5 terbaik), nilai akhir = skor * bobot.

## Matriks Skor


| Kriteria             | Bobot   | Opsi A  | Opsi B  | Opsi C  |
| -------------------- | ------- | ------- | ------- | ------- |
| Retrieval quality    | 30      | 4       | 3       | 4       |
| Latency p95          | 20      | 4       | 3       | 4       |
| Kompleksitas migrasi | 15      | 5       | 2       | 4       |
| Risiko regresi       | 15      | 4       | 2       | 4       |
| Operasional          | 10      | 3       | 4       | 4       |
| Legal/lisensi        | 10      | 4       | 2       | 3       |
| **Total tertimbang** | **100** | **405** | **270** | **390** |


Interpretasi:

1. Opsi A tertinggi untuk stabilitas jangka dekat.
2. Opsi C mendekati Opsi A dan paling baik untuk roadmap unifikasi.
3. Opsi B terendah karena risiko teknis + legal saat ini.

## Keputusan

Dipilih kombinasi:

1. Jangka dekat: Opsi A (tetap split).
2. Jangka menengah: Opsi C (migrasi bertahap terukur).
3. Opsi B (big-bang sekarang): ditolak.

## Alasan Teknis Utama

1. Path RAG Veridium saat ini ketat pada pola DuckDB(top-K vector) + SQLite(chunk metadata/content) + file-level filtering.
2. Embedded wrapper Muninn yang dipakai Veridium belum mengekspos kontrak retrieval spesifik yang dibutuhkan untuk parity RAG Veridium (terutama contract adapter yang sekarang dipakai aplikasi).
3. Perlu adapter dan validasi kualitas/latensi dulu sebelum replace jalur produksi.
4. Ada risiko lisensi BSL untuk penggunaan komersial tertentu, sehingga unifikasi penuh perlu legal clearance.

## Risiko yang Harus Ditangani Sekarang

Risiko prioritas tinggi pada DuckDB saat ini:

1. Inisialisasi melakukan `DROP TABLE IF EXISTS vectors`, berpotensi menghapus index/data vector pada startup.

Tindakan wajib:

1. Ganti mekanisme ini dengan migrasi schema yang aman (idempotent, tanpa destructive default).
2. Tambah smoke test startup untuk memastikan data vector persisten lintas restart.

## Rencana Implementasi (Bertahap)

Tahap 1 - Hardening baseline (wajib sebelum eksperimen):

1. Hilangkan destructive init pada DuckDB.
2. Tambah benchmark baseline: Recall@5/10, nDCG@10, p95 retrieval latency, ingestion throughput.
3. Freeze dataset evaluasi agar hasil antar backend bisa dibandingkan apple-to-apple.

Tahap 2 - Abstraction:

1. Introduce `VectorStore` interface di Veridium.
2. Adapter pertama: `DuckDBVectorStore` (existing behavior).
3. Adapter kedua: `MuninnVectorStore` (fitur minimum untuk parity RAG Veridium).

Tahap 3 - A/B Validation:

1. Jalankan dual-run (shadow mode) tanpa mengubah jawaban user.
2. Bandingkan metrik kualitas dan latensi.
3. Dokumentasikan gap retrieval berdasarkan use case (FAQ, code docs, long PDF, multi-file scope).

Tahap 4 - Rollout:

1. Feature flag per workspace/KB.
2. Canary internal.
3. Rollout bertahap dengan rollback one-click.

## Go/No-Go Gates

Muninn boleh jadi default vector backend hanya jika seluruh syarat terpenuhi:

1. Recall@10 minimal setara baseline DuckDB (delta >= -1%).
2. nDCG@10 minimal setara baseline DuckDB (delta >= -2%).
3. p95 latency retrieval tidak lebih lambat dari 15% baseline.
4. Tidak ada P1/P2 regresi fungsional selama 14 hari canary.
5. Legal approval tertulis untuk model distribusi/komersialisasi produk.

Jika salah satu gagal, tetap di split architecture dan lanjutkan iterasi adapter.

## Dampak

Positif:

1. Mengurangi risiko regresi produksi jangka pendek.
2. Tetap membuka jalur unifikasi backend yang terukur.
3. Memaksa keputusan berbasis metrik, bukan preferensi arsitektur.

Negatif:

1. Kompleksitas sementara meningkat karena dua backend berjalan paralel.
2. Membutuhkan investasi tooling benchmark dan observability.

## Catatan Validasi Saat ADR Ditulis

1. Test subset `veridium/internal/services` untuk area Muninn/Vector/RAG lulus.
2. Test `muninndb/pkg/embedded` dan `muninndb/pkg/engine/activation` lulus.
3. Ada warning linker macOS version mismatch; tidak memblokir test result.
