## Analisis Kelayakan Model Kimi K2.5 untuk cmd/server

Berdasarkan analisis terhadap dokumentasi Unsloth Kimi K2.5 dan codebase `cmd/server`, berikut evaluasi kelayakan model ini:

### Spesifikasi Kimi K2.5

| Parameter | Nilai |
|-----------|-------|
| **Parameter** | 1 Trillion (MoE architecture) |
| **Context Length** | Up to 256K tokens (recommended: 98,304) |
| **Disk Space (1-bit quant)** | 240GB |
| **Disk Space (2-bit quant)** | 375GB (UD-Q2_K_XL) |
| **Disk Space (4-bit quant)** | ~630GB |
| **RAM+VRAM Requirement** | ≥ 240GB (unified memory) |
| **Architecture** | Modified DeepSeek V3 MoE |

### Kompatibilitas dengan cmd/server

#### ✅ **DIDUKUNG** - Fitur yang Kompatibel

1. **OpenAI-Compatible API** - Server sudah mendukung `/v1/chat/completions` endpoint yang kompatibel dengan format OpenAI API yang digunakan oleh llama-server.

2. **llama.cpp Backend** - Server menggunakan llama.cpp library (via `pkg/kronk`) yang mendukung GGUF format. Kimi K2.5 tersedia dalam format GGUF dari Unsloth.

3. **Context Window Configuration** - Server mendukung konfigurasi context window via [`ContextWindow`](pkg/kronk/model/config.go:160) di [`model.Config`](pkg/kronk/model/config.go:154). Default adalah 8K (8192) tapi bisa dikonfigurasi hingga 256K untuk Kimi K2.5.

4. **MoE Support** - Server mendukung MoE models via [`SplitModeRow`](pkg/kronk/model/config.go:131) yang direkomendasikan untuk MoE models.

5. **Sampling Parameters** - Server mendukung parameter sampling yang direkomendasikan Kimi K2.5:
   - `temperature` (0.6 untuk Instant Mode, 1.0 untuk Thinking Mode)
   - `top_p` (0.95)
   - `min_p` (0.01)

#### ⚠️ **PERTIMBANGAN** - Yang Perlu Diperhatikan

1. **Resource Requirements** - Kimi K2.5 membutuhkan resources yang sangat besar:
   - Minimum 240GB disk space (1-bit quant)
   - Minimum 240GB combined RAM+VRAM
   - Untuk performa yang baik (>10 tokens/s): ~256GB RAM

2. **Chat Template** - Kimi K2.5 menggunakan chat template khusus:
   ```
   <|im_system|>system<|im_middle|>...<|im_end|><|im_user|>user<|im_middle|>...<|im_end|><|im_assistant|>assistant<|im_middle|><think>
   ```
   Server menggunakan Jinja template system yang bisa menangani ini jika template tersedia.

3. **Vision Support** - Saat ini belum ada vision support di llama.cpp untuk Kimi K2.5 (walaupun modelnya support vision via MoonViT).

4. **Special Parameters** - Kimi K2.5 memerlukan `rope_scaling.beta_fast = 32.0` (berbeda dari K2 Thinking yang menggunakan 1.0).

#### ❌ **KETERBATASAN**

1. **No Vision in llama.cpp** - Vision encoder (MoonViT) belum didukung di llama.cpp.

2. **Model Size** - Model sangat besar sehingga tidak praktis untuk deployment di hardware consumer.

3. **Quantization Quality** - 1-bit quant (UD-TQ1_0) mungkin mengurangi kualitas signifikan untuk beberapa use cases.

### Rekomendasi Konfigurasi

Jika ingin menggunakan Kimi K2.5 di `cmd/server`, konfigurasi yang disarankan:

```yaml
# Model config untuk Kimi K2.5
kimi-k2.5:
  context-window: 98304      # Recommended context length
  nbatch: 2048               # Default batch size
  nubatch: 512               # Physical batch size
  device: "cuda"             # atau "cpu" jika tidak ada GPU
  flash-attention: 1         # Enable flash attention
  split-mode: 2              # SplitModeRow untuk MoE
  ngpu-layers: -1            # Offload semua layer ke GPU jika memungkinkan
  cache-type-k: q8_0         # KV cache quantization
  cache-type-v: q8_0
```

### Kesimpulan

**Kimi K2.5 TEKNIS KOMPATIBEL** dengan `cmd/server` karena:
- Server menggunakan llama.cpp yang support GGUF
- API endpoint sudah OpenAI-compatible
- Context window configurable
- MoE architecture didukung

Namun, **TIDAK PRAKTIS** untuk deployment umum karena:
- Resource requirements sangat tinggi (240GB+ disk, 240GB+ RAM/VRAM)
- Vision belum support di llama.cpp
- Model size terlalu besar untuk consumer hardware

**Alternatif yang lebih feasible:**
- Gunakan model yang lebih kecil seperti Llama 3.x, Qwen3, atau DeepSeek V3 yang lebih kecil
- Jika memerlukan model besar, pertimbangkan Qwen3-235B-A22B yang lebih efisien

## Analisis Kelayakan Model GLM-4.7-Flash untuk cmd/server

Berdasarkan analisis terhadap dokumentasi Unsloth GLM-4.7-Flash dan codebase `cmd/server`, berikut evaluasi kelayakan model ini:

### Spesifikasi GLM-4.7-Flash

| Parameter | Nilai |
|-----------|-------|
| **Parameter** | 30B MoE (aktif ~3.6B parameters) |
| **Context Length** | Up to 200K tokens (max: 202,752) |
| **RAM/VRAM Requirement** | 18-24GB (4-bit quant), 32GB (full precision) |
| **Architecture** | MoE (Mixture of Experts) |
| **Quantization** | UD-Q4_K_XL (recommended), UD-Q2_K_XL |
| **Disk Size** | ~18GB (4-bit) |

### Kompatibilitas dengan cmd/server

#### ✅ **SANGAT KOMPATIBEL** - Fitur yang Didukung

1. **llama.cpp Native Support** - GLM-4.7-Flash berjalan langsung di llama.cpp (tidak seperti beberapa model lain yang memerlukan workaround). Server menggunakan `pkg/kronk` yang berbasis llama.cpp.

2. **Resource Requirements Realistis** - Hanya membutuhkan 18-24GB RAM/VRAM untuk 4-bit quantization, jauh lebih feasible dibanding Kimi K2.5 yang butuh 240GB+.

3. **Context Window Configurable** - Server mendukung konfigurasi context window via [`ContextWindow`](pkg/kronk/model/config.go:160). GLM-4.7-Flash support hingga 200K tokens.

4. **Jinja Template Support** - Server menggunakan Jinja template system (via `pkg/tools/templates`) yang kompatibel dengan `--jinja` flag yang direkomendasikan untuk GLM-4.7-Flash.

5. **MoE Architecture Support** - Server mendukung MoE models via [`SplitModeRow`](pkg/kronk/model/config.go:131).

6. **Sampling Parameters** - Server mendukung semua parameter sampling yang direkomendasikan:
   - `temperature` (1.0 untuk general, 0.7 untuk tool-calling)
   - `top_p` (0.95 atau 1.0)
   - `min_p` (0.01)
   - `repeat_penalty` (disable atau 1.0)

#### ⚠️ **PERTIMBANGAN** - Yang Perlu Diperhatikan

1. **Chat Template** - GLM-4.7-Flash menggunakan chat template khusus. Perlu memastikan template tersedia di catalog/template system server.

2. **Scoring Function** - llama.cpp perlu versi terbaru (post Jan 21) yang fix bug `scoring_func` dari "softmax" ke "sigmoid". Server menggunakan library system yang bisa diupdate.

3. **Ollama Compatibility Warning** - Unsloth secara eksplisit tidak merekomendasikan penggunaan dengan Ollama karena chat template issues, tapi ini tidak berlaku untuk llama.cpp yang digunakan server.

4. **KV Cache Unified** - Dokumentasi menyebutkan `--kv-unified` bisa meningkatkan performa. Server perlu verifikasi apakah ini didukung di llama.cpp version yang digunakan.

### Perbandingan dengan Kimi K2.5

| Aspek | GLM-4.7-Flash | Kimi K2.5 |
|-------|---------------|-----------|
| **Parameter** | 30B MoE (~3.6B aktif) | 1T MoE |
| **Context** | 200K | 256K |
| **RAM/VRAM** | 18-24GB | 240GB+ |
| **Disk** | ~18GB | 240GB+ |
| **llama.cpp Support** | ✅ Native | ✅ Native |
| **Feasibility** | ✅ Consumer Hardware | ❌ Enterprise Only |

### Rekomendasi Konfigurasi

```yaml
# Model config untuk GLM-4.7-Flash
glm-4.7-flash:
  context-window: 16384        # Bisa dinaikkan hingga 202752
  nbatch: 2048                 # Default batch size
  nubatch: 512                 # Physical batch size
  device: "cuda"               # atau "cpu" jika tidak ada GPU
  flash-attention: 1           # Enable flash attention
  split-mode: 2                # SplitModeRow untuk MoE
  ngpu-layers: -1              # Offload semua layer ke GPU
  cache-type-k: q8_0           # KV cache quantization
  cache-type-v: q8_0
  # Sampling defaults
  temperature: 1.0             # 0.7 untuk tool-calling
  top-p: 0.95
  min-p: 0.01
  repeat-penalty: 1.0          # Disable repeat penalty
```

### Cara Menambahkan ke Server

1. **Download Model**:
   ```bash
   # Model tersedia di HuggingFace
   unsloth/GLM-4.7-Flash-GGUF
   ```

2. **Tambahkan ke Catalog** - Perlu menambahkan entry ke catalog system server dengan:
   - Model ID: `glm-4.7-flash`
   - GGUF URL: `unsloth/GLM-4.7-Flash-GGUF`
   - Chat template: GLM-4.7 specific

3. **Update llama.cpp** - Pastikan menggunakan versi terbaru yang sudah fix `scoring_func` bug.

### Kesimpulan

**GLM-4.7-Flash SANGAT LAYAK** untuk `cmd/server` karena:

1. ✅ **Resource Requirements Realistis** - 18-24GB RAM/VRAM feasible untuk consumer hardware
2. ✅ **Native llama.cpp Support** - Berjalan langsung tanpa workaround
3. ✅ **Context Window Besar** - 200K tokens untuk long-context tasks
4. ✅ **MoE Architecture** - Efisien dengan ~3.6B active parameters
5. ✅ **Performance** - Leading benchmarks di SWE-Bench, GPQA, reasoning/chat
6. ✅ **Server Compatibility** - Semua fitur server (Jinja template, sampling params, MoE support) kompatibel

**Rekomendasi**: GLM-4.7-Flash adalah pilihan yang **jauh lebih baik** dibanding Kimi K2.5 untuk deployment di `cmd/server` karena keseimbangan yang baik antara capability dan resource requirements.

## Analisis Kelayakan Model NVIDIA Nemotron 3 Nano untuk cmd/server

Berdasarkan analisis terhadap dokumentasi Unsloth Nemotron 3 Nano dan codebase `cmd/server`, berikut evaluasi kelayakan model ini:

### Spesifikasi Nemotron 3 Nano

| Parameter | Nilai |
|-----------|-------|
| **Parameter** | 30B hybrid reasoning MoE (~3.6B active parameters) |
| **Context Length** | Up to 1M tokens (default: 262,144 / 256K) |
| **RAM/VRAM Requirement** | 24GB (4-bit quant) |
| **Architecture** | MoE (Mixture of Experts) dengan NoPE (No Positional Embeddings) |
| **Quantization** | UD-Q4_K_XL (recommended), BF16, FP8 |
| **Disk Size** | ~18GB (4-bit) |
| **Reasoning** | Native `<think>` / `</think>` tokens (ID 12/13) |

### Kompatibilitas dengan cmd/server

#### ✅ **SANGAT KOMPATIBEL** - Fitur yang Didukung

1. **llama.cpp Native Support** - Nemotron 3 Nano berjalan langsung di llama.cpp dengan GGUF format. Server menggunakan `pkg/kronk` yang berbasis llama.cpp.

2. **Resource Requirements Realistis** - Hanya membutuhkan 24GB RAM/VRAM untuk 4-bit quantization, sama seperti GLM-4.7-Flash dan feasible untuk consumer hardware high-end.

3. **Context Window Configurable** - Server mendukung konfigurasi context window via [`ContextWindow`](pkg/kronk/model/config.go:160). Nemotron 3 Nano support hingga **1M tokens** - ini adalah context window terbesar di antara ketiga model yang dianalisis!

4. **Jinja Template Support** - Server menggunakan Jinja template system yang kompatibel dengan `--jinja` flag yang direkomendasikan untuk Nemotron 3.

5. **MoE Architecture Support** - Server mendukung MoE models via [`SplitModeRow`](pkg/kronk/model/config.go:131).

6. **Sampling Parameters** - Server mendukung semua parameter sampling yang direkomendasikan:
   - `temperature` (1.0 untuk general, 0.6 untuk tool-calling)
   - `top_p` (1.0 atau 0.95)
   - `--special` flag untuk reasoning tokens

#### ⚠️ **PERTIMBANGAN KHUSUS** - Yang Perlu Diperhatikan

1. **Reasoning Tokens** - Nemotron 3 menggunakan `<think>` (token ID 12) dan `</think>` (token ID 13) untuk reasoning. Ini adalah fitur native yang berbeda dari model lain. Server perlu konfigurasi khusus untuk menangani ini:
   - Gunakan `--special` flag untuk melihat reasoning tokens
   - Perlu `--verbose-prompt` untuk debugging

2. **NoPE (No Positional Embeddings)** - Model trained tanpa positional embeddings eksplisit:
   - **Keuntungan**: Tidak perlu YaRN untuk context window besar
   - Hanya perlu ubah `max_position_embeddings`
   - Lebih efisien untuk long context

3. **Context Window Warning** - Setting ke 1M tokens bisa trigger CUDA OOM:
   - Default direkomendasikan: 262,144 (256K)
   - Perlu hardware yang sangat powerful untuk 1M context

4. **Chat Template** - Nemotron 3 menggunakan format khusus:
   ```
   <|im_start|>system\n<|im_end|>\n<|im_start|>user\n...\n<|im_end|>\n<|im_start|>assistant\n<think></think>...
   ```
   Perlu memastikan template tersedia di catalog system server.

### Perbandingan Ketiga Model

| Aspek | Kimi K2.5 | GLM-4.7-Flash | Nemotron 3 Nano |
|-------|-----------|---------------|-----------------|
| **Parameter** | 1T MoE | 30B MoE (~3.6B) | 30B MoE (~3.6B) |
| **Context** | 256K | 200K | **1M** (default 256K) |
| **RAM/VRAM** | 240GB+ | 18-24GB | 24GB |
| **Disk** | 240GB+ | ~18GB | ~18GB |
| **Architecture** | MoE | MoE | MoE + NoPE |
| **Reasoning** | Native | Native | `<think>` tokens |
| **llama.cpp** | ✅ Native | ✅ Native | ✅ Native |
| **Feasibility** | ❌ Enterprise | ✅ Consumer | ✅ Consumer |

### Rekomendasi Konfigurasi

```yaml
# Model config untuk Nemotron 3 Nano
nemotron-3-nano:
  context-window: 262144      # Default 256K, bisa naik ke 1M jika hardware kuat
  nbatch: 2048                # Default batch size
  nubatch: 512                # Physical batch size
  device: "cuda"              # atau "cpu" jika tidak ada GPU
  flash-attention: 1          # Enable flash attention
  split-mode: 2               # SplitModeRow untuk MoE
  ngpu-layers: -1             # Offload semua layer ke GPU
  cache-type-k: q8_0          # KV cache quantization
  cache-type-v: q8_0
  # Sampling defaults
  temperature: 1.0            # 0.6 untuk tool-calling
  top-p: 1.0                  # 0.95 untuk tool-calling
  # Special flags untuk reasoning
  special: true               # Untuk <think> tokens
```

### Cara Menambahkan ke Server

1. **Download Model**:
   ```bash
   # Model tersedia di HuggingFace
   unsloth/Nemotron-3-Nano-30B-A3B-GGUF
   ```

2. **Tambahkan ke Catalog** - Perlu menambahkan entry ke catalog system server dengan:
   - Model ID: `nemotron-3-nano`
   - GGUF URL: `unsloth/Nemotron-3-Nano-30B-A3B-GGUF`
   - Chat template: Nemotron 3 specific dengan `<think>` support

3. **Update llama.cpp** - Pastikan menggunakan versi terbaru yang support Nemotron 3.

### Keunggulan Nemotron 3 Nano

1. **Context Window Terbesar** - 1M tokens (vs 256K Kimi, 200K GLM)
2. **NoPE Architecture** - Lebih efisien untuk long context, tidak perlu YaRN
3. **Native Reasoning** - `<think>` tokens untuk explicit reasoning
4. **NVIDIA Optimized** - Dari NVIDIA, optimized untuk GPU NVIDIA
5. **Day-Zero Support** - Unsloth memberikan support langsung dari rilis

### Kesimpulan

**Nemotron 3 Nano SANGAT LAYAK** untuk `cmd/server` karena:

1. ✅ **Context Window Terbesar** - 1M tokens untuk ultra-long context tasks
2. ✅ **Resource Requirements Realistis** - 24GB RAM/VRAM feasible untuk consumer hardware high-end
3. ✅ **NoPE Architecture** - Lebih efisien, tidak perlu YaRN
4. ✅ **Native Reasoning** - `<think>` tokens untuk better reasoning visibility
5. ✅ **Native llama.cpp Support** - Berjalan langsung tanpa workaround
6. ✅ **MoE Architecture** - Efisien dengan ~3.6B active parameters
7. ✅ **Server Compatibility** - Semua fitur server kompatibel

**Rekomendasi**: Nemotron 3 Nano adalah pilihan **terbaik** untuk use cases yang memerlukan:
- Ultra-long context (hingga 1M tokens)
- Explicit reasoning visibility
- NVIDIA GPU optimization
- Balance antara capability dan resource requirements

**Peringatan**: Jika hardware terbatas (kurang dari 24GB VRAM), pertimbangkan GLM-4.7-Flash yang bisa berjalan di 18GB.