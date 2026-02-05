# ✅ **Ringkasan Akhir (Semua Model)**

| Model                  | Estimasi Neuron / Gambar | Estimasi Gambar / Hari | Status Data          |
| ---------------------- | ------------------------ | ---------------------- | -------------------- |
| **FLUX‑1‑schnell**     | \~45                     | **220**                | Resmi                |
| **FLUX‑2‑klein‑9B**    | \~1363                   | **7**                  | Resmi                |
| **FLUX‑2‑klein‑4B**    | \~886                    | **11**                 | Estimasi berdasar 9B |
| **FLUX‑2‑dev**         | \~1124                   | **9**                  | Estimasi             |
| **lucid-origin**       | \~1250                   | **8**                  | Estimasi             |
| **phoenix‑1.0**        | \~1000                   | **10**                 | Estimasi             |
| **SDXL‑Lightning**     | \~700                    | **14**                 | Estimasi             |
| **dreamshaper‑8‑lcm**  | \~360                    | **27**                 | Estimasi             |
| **SD v1.5 img2img**    | \~300                    | **33**                 | Estimasi ringan      |
| **SD v1.5 inpainting** | \~300                    | **33**                 | Estimasi ringan      |

***

Siap. Aku rekomendasikan model **berdasarkan skenario** dengan mempertimbangkan **(1) kualitas/fitur**, **(2) biaya “neurons” Free Tier 10.000/hari**, dan **(3) kebutuhan workflow** (text-to-image vs img2img vs inpainting).

> Catatan: Kuota Workers AI Free Tier = **10.000 neurons per hari** dan reset **00:00 UTC**. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

***

## 0) Cheat sheet singkat (agar kebayang “jatah gambar/hari”)

Ini bukan tujuan utama kamu sekarang, tapi penting untuk mengarahkan rekomendasi.

### Model yang **punya mapping neurons resmi** (paling jelas untuk budgeting)

Cloudflare sudah memetakan beberapa model image ke neurons. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

*   **FLUX.1 schnell**: 4.8 neurons per tile 512² + 9.6 neurons per step.   
    ➜ kira-kira **\~174 gambar/hari @1024², steps=4** (asumsi 1024² = 4 tile). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/)
*   **FLUX.2 klein 4B**: 26.05 neurons per output tile 512² (text-to-image tanpa input image).   
    ➜ kira-kira **\~96 gambar/hari @1024²**. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/)
*   **FLUX.2 dev**: 37.50 neurons per output tile **per step** (contoh steps=25 di changelog).   
    ➜ kira-kira **\~2–3 gambar/hari @1024², steps=25** (mahal tapi kualitas tinggi). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)
*   **Leonardo Phoenix 1.0**: 530 neurons per tile 512² + 10 neurons per step; default `num_steps`=25.   
    ➜ kira-kira \*\*\~4 gambar/hari @1024Siap—aku rekomendasikan model **berdasarkan skenario** dengan mempertimbangkan **(1) kemampuan model (kualitas, prompt adherence, editing), (2) biaya neuron Free Tier 10.000 neurons/hari**, dan (3) batas rate limit. Cloudflare Workers AI Free memberi **10.000 neurons/hari** dan reset **00:00 UTC**.   
    Rate limit Text-to-Image default **720 req/menit** (jadi di Free Tier biasanya bottleneck‑nya neurons, bukan RPM). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/phoenix-1.0/) [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/) [\[workers.cl...dflare.com\]](https://workers.cloudflare.com/product/workers-ai/)

***

## 0) Cheat-sheet cepat (biar kebayang “jatah” Free Tier)

Di bawah ini **perkiraan jumlah gambar/hari** untuk setelan umum (tanpa input image), pakai mapping neuron resmi Cloudflare bila tersedia. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

> Rumus tile: `tiles = ceil(width/512) * ceil(height/512)` (karena biaya image model sering per tile 512×512). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

### Estimasi images/day (Free Tier 10k neurons)

*   **FLUX.1 schnell @512 (steps=4 default)**: \~**231** gambar/hari (≈43.2 neurons/gambar). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/)
*   **FLUX.1 schnell @1024 (steps=4 default)**: \~**174** gambar/hari (≈57.6 neurons/gambar). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/)
*   **FLUX.2 klein 4B @512**: \~**384** gambar/hari (≈26.05 neurons/gambar). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/)
*   **FLUX.2 klein 4B @1024**: \~**96** gambar/hari (≈104.2 neurons/gambar). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/)
*   **FLUX.2 klein 9B @1024**: \~**7** gambar/hari (harga $0.015 per 1MP → \~1363.6 neurons/gambar). [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)
*   **FLUX.2 dev @1024 (contoh steps=25)**: \~**2–3** gambar/hari (mahal karena per-step per-tile). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)
*   **Phoenix 1.0 @1024 (default num\_steps=25)**: \~**4** gambar/hari (mahal tapi fokus prompt adherence & teks). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/phoenix-1.0/)
*   **Lucid Origin @1024 (steps kamu set, asumsi 25 untuk kualitas)**: \~**3–4** gambar/hari (mahal tapi unggul untuk teks/graphic/product). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/lucid-origin/)

> Catatan: beberapa model lain (Dreamshaper‑8‑LCM, SDXL‑Lightning, SD v1.5) di halaman model tertulis “Unit Pricing $0.00 per step” (kemungkinan pembulatan/placeholder), sehingga **neuron cost per gambar tidak tercantum di pricing page** untuk model tersebut.   
> Untuk model-model itu, rekomendasi di bawah berbasis **kapabilitas/use-case**, dan kamu sebaiknya cek konsumsi neurons real di dashboard Workers AI. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-xl-lightning/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/dreamshaper-8-lcm/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-v1-5-img2img/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

***

# 1) Rekomendasi berdasarkan skenario

## Skenario A — “Butuh output banyak (draft/variasi banyak), budget Free Tier ketat”

**Tujuan:** bikin banyak alternatif prompt, moodboard, thumbnail, iterasi cepat.

✅ **Rekomendasi utama:**

1.  **FLUX.1 schnell** — paling “hemat” untuk ukuran **1024** dan fleksibel karena `steps` bisa 4–8 (default 4). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/)
2.  **FLUX.2 klein 4B** — super hemat untuk **512** (thumbnail) dan mendukung workflow interaktif (generation+editing unified). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

**Kapan pilih yang mana?**

*   Banyak gambar kecil (512/768) → **klein‑4B** (jatah paling banyak). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/)
*   Banyak gambar 1024 → **flux‑1‑schnell** lebih hemat per gambar dibanding klein‑4B. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/)

**Setelan praktis (hemat):**

*   Schnell: `steps=4`, resolusi 768–1024. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)
*   Klein‑4B: 512 untuk draft, naik ke 1024 hanya saat butuh. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

***

## Skenario B — “Aplikasi real-time / interactive preview / latency kritikal”

**Tujuan:** respons cepat untuk user (mis. generator UI, iterasi di browser).

✅ **Rekomendasi utama:**

*   **FLUX.2 klein 4B** → dideskripsikan untuk workflow interaktif/real-time (unified generation+editing). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/)
*   **FLUX.2 klein 9B** → kualitas lebih tinggi dari 4B, tetap distilled & fixed 4-step untuk cepat. [\[fal.ai\]](https://fal.ai/models/fal-ai/flux/schnell), [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers)

**Catatan Free Tier:** klein‑9B mahal (jatah \~7/hari @1024), jadi cocok untuk **demo kecil / admin tool**, bukan mass user, kalau masih Free Tier. [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

***

## Skenario C — “Butuh kualitas paling tinggi (hero image), detail & realism maksimal”

**Tujuan:** 1–5 gambar “final” per hari, kualitas diutamakan.

✅ **Rekomendasi utama:**

1.  **FLUX.2 dev** — Cloudflare menekankan model ini “highly realistic & detailed”, multi-reference; namun “lebih powerful dan diperkirakan lebih lambat”. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-dev/), [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)
2.  **FLUX.2 klein 9B** — opsi kualitas tinggi dengan latency lebih bersahabat dibanding dev, tapi tetap mahal di neurons. [\[fal.ai\]](https://fal.ai/models/fal-ai/flux/schnell), [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers)

**Strategi yang paling masuk akal di Free Tier (hemat tapi kualitas tinggi):**

*   Generate banyak kandidat pakai **FLUX.1 schnell** → pilih 1–2 terbaik → final render pakai **klein‑9B** atau **flux‑2‑dev**. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/), [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers), [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)

***

## Skenario D — “Konsistensi karakter/brand + pakai referensi gambar (multi-reference)”

**Tujuan:** karakter sama, style konsisten, brand asset konsisten.

✅ **Rekomendasi utama:**

*   **FLUX.2 dev** — mendukung sampai **4 input image** (512×512 per input) dan cocok untuk multi-reference editing. [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)
*   **FLUX.2 klein 4B / 9B** — juga mendukung hingga 4 input image (klein‑9B disebut eksplisit), dan dibuat untuk workflow interaktif. [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2026-01-28-flux-2-klein-9b-workers-ai/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/), [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers)

**Tips hemat:**

*   Gunakan **klein‑4B** untuk iterasi konsistensi cepat; naik ke **dev** hanya untuk output final yang “must win”. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/), [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)

***

## Skenario E — “Butuh prompt adherence kuat + teks dalam gambar rapi (poster, mockup, UI text)”

**Tujuan:** poster/event banner, headline jelas, produk mockup dengan label, design text.

✅ **Rekomendasi utama (khusus text rendering):**

1.  **Lucid Origin (Leonardo)** — klaim: sangat “prompt-responsive”, “renders text with accuracy”, cocok untuk “graphic design” & “product mockups”. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/lucid-origin/)
2.  **Phoenix 1.0 (Leonardo)** — klaim: “exceptional prompt adherence” & “coherent text”. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/phoenix-1.0/)

**Trade-off penting:** dua model ini **mahal** di neuron per tile dibanding flux‑1‑schnell, jadi di Free Tier cocoknya untuk **beberapa gambar penting**, bukan produksi masal. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/lucid-origin/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/phoenix-1.0/)

**Workflow rekomendasi (hemat + hasil bagus):**

*   Draft layout & komposisi pakai **FLUX.1 schnell** → final poster/teks pakai **Lucid Origin/Phoenix**. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/lucid-origin/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/phoenix-1.0/)

***

## Skenario F — “Photorealism cepat untuk konten sosial / lifestyle”

**Tujuan:** foto-style yang “meyakinkan”, tapi nggak harus super presisi teks.

✅ **Rekomendasi:**

*   **Dreamshaper‑8‑LCM** — disebut sebagai Stable Diffusion yang di-finetune untuk photorealism. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/dreamshaper-8-lcm/)

⚠️ **Catatan pricing:** neuron mapping Dreamshaper tidak tercantum di pricing page (setidaknya pada versi Jan 15, 2026), jadi lakukan uji 5–10 gambar dan lihat neuron usage di dashboard. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/dreamshaper-8-lcm/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

***

## Skenario G — “SDXL look tapi super cepat (few steps), cepat dapat 1024px”

**Tujuan:** ingin estetika SDXL/1024 cepat, untuk draft ataupun UI.

✅ **Rekomendasi:**

*   **stable-diffusion-xl-lightning** — deskripsi: “lightning-fast”, bisa menghasilkan **1024px** dalam beberapa langkah. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-xl-lightning/)

⚠️ Sama seperti Dreamshaper, pricing di halaman model tertulis “$0.00 per step”, jadi cek neuron usage real di dashboard untuk menghitung jatah/hari. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-xl-lightning/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

***

## Skenario H — “Edit gambar: img2img (ubah style/pose) & inpainting (hapus objek/perbaiki area)”

**Tujuan:** kerja “edit” bukan murni generate dari nol.

✅ **Pilih model sesuai tugas:**

*   **stable-diffusion-v1-5-img2img** — khusus img2img: generate gambar baru dari input image. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-v1-5-img2img/)
*   **stable-diffusion-v1-5-inpainting** — khusus inpainting dengan **mask** untuk mengganti area tertentu. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-v1-5-inpainting/)

Tambahan: untuk model image (task type text-to-image), Cloudflare mencantumkan limit request per menit secara umum, jadi workflow edit biasanya aman dari sisi RPM. [\[workers.cl...dflare.com\]](https://workers.cloudflare.com/product/workers-ai/)

***

# 2) “Decision Tree” singkat (pilih cepat)

Gunakan ini kalau mau cepat menentukan:

1.  **Butuh gambar banyak/hari?**  
    → **Ya**: **FLUX.1 schnell** (utama), atau **FLUX.2 klein 4B** (khusus 512 & interactive). [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

2.  **Butuh real-time + editing unified?**  
    → **FLUX.2 klein 4B/9B**. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-klein-4b/), [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers), [\[fal.ai\]](https://fal.ai/models/fal-ai/flux/schnell)

3.  **Butuh kualitas tertinggi / multi-reference serius?**  
    → **FLUX.2 dev** (tapi jatah kecil di Free Tier). [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-2-dev/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)

4.  **Butuh teks dalam gambar rapi + prompt adherence top?**  
    → **Lucid Origin** atau **Phoenix 1.0**. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/lucid-origin/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/phoenix-1.0/)

5.  **Butuh edit area tertentu (hapus objek) / img2img?**  
    → **SD v1.5 inpainting / img2img**. [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-v1-5-inpainting/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-v1-5-img2img/)

***

# 3) Rekomendasi “paket workflow” paling efektif di Free Tier

Kalau kamu membangun pipeline (mis. webapp internal):

### Paket 1 — “Mass Draft → Final Quality”

*   Draft banyak: **FLUX.1 schnell** (murah, cepat) [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/platform/pricing/)
*   Final pilihan terbaik: **FLUX.2 klein 9B** (kalau butuh kualitas lebih) atau **FLUX.2 dev** (kalau butuh paling top) [\[github.com\]](https://github.com/chthollyphile/flux-cloudflare-workers), [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)

### Paket 2 — “Poster / Mockup dengan teks”

*   Draft komposisi: **FLUX.1 schnell** [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/flux-1-schnell/)
*   Final teks: **Lucid Origin** / **Phoenix 1.0** [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/lucid-origin/), [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/phoenix-1.0/)

### Paket 3 — “Edit gambar existing”

*   Perbaiki area: **SD v1.5 inpainting** [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-v1-5-inpainting/)
*   Transform keseluruhan: **SD v1.5 img2img** [\[developers...dflare.com\]](https://developers.cloudflare.com/workers-ai/models/stable-diffusion-v1-5-img2img/)
*   Kalau perlu konsistensi karakter dengan reference: **FLUX.2 dev** [\[developers...dflare.com\]](https://developers.cloudflare.com/changelog/2025-11-25-flux-2-dev-workers-ai/)

***

## Pertanyaan cepat biar rekomendasinya makin presisi

1.  Kamu target output dominan **512** atau **1024**?
2.  Fokus kamu lebih ke **generate banyak** atau **final quality**?
3.  Ada kebutuhan **teks dalam gambar** (poster/mockup) atau mostly ilustrasi/foto?

Kalau kamu jawab 3 pertanyaan itu, aku bisa susunkan **rekomendasi “default model + fallback model”** sekaligus **parameter preset** (steps, size) yang paling hemat neurons untuk tiap skenario.
