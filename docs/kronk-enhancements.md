# Kronk Enhancements Implementation

## Overview
This document describes the features implemented from the Kronk repository (commits a92184d9..HEAD) into the Veridium server.

## Branch
`feature/kronk-enhancements`

## Implemented Features

### 1. Model Metadata Reader (`pkg/tools/models/info.go`)

**Purpose**: Read GGUF model metadata directly from model files without loading them into memory.

**Key Functions**:
- `ModelInformation(modelID string) (ModelInfo, error)` - Reads model metadata from GGUF file
- Supports all GGUF metadata value types (uint8, int8, uint16, int16, uint32, int32, float32, float64, bool, string, array)

**Benefits**:
- Fast model information retrieval
- No need to load model into memory
- Extracts model architecture, size, and metadata

**Example Usage**:
```go
info, err := models.ModelInformation("qwen2-audio-7b")
// Returns: ModelInfo with ID, Size, Metadata, IsGPTModel, IsEmbedModel, etc.
```

---

### 2. VRAM Calculator (`pkg/tools/models/vram.go`)

**Purpose**: Calculate VRAM requirements for running models based on configuration parameters.

**Key Functions**:
- `CalculateVRAM(modelID string, cfg VRAMConfig) (VRAM, error)` - Calculate VRAM for local models
- `CalculateVRAMFromHuggingFace(ctx context.Context, modelURL string, cfg VRAMConfig) (VRAM, error)` - Calculate VRAM for HuggingFace models

**Configuration Parameters** (`VRAMConfig`):
- `ContextWindow` - Context window size (e.g., 8192, 131072)
- `BytesPerElement` - Cache type: 1 for q8_0, 2 for f16
- `Slots` - Number of concurrent sequences (n_seq_max)
- `CacheSequences` - Additional cache sequences: 0=none, 1=FMC or SPC, 2=FMC+SPC

**Constants Provided**:
```go
// Context Windows
ContextWindow8K   = 8192
ContextWindow128K = 131072

// Bytes Per Element
BytesPerElementQ8_0 = 1
BytesPerElementF16  = 2

// Slots
Slots2 = 2
Slots4 = 4

// Cache Sequences
CacheSequenceNone   = 0
CacheSequenceSingle = 1  // FMC or SPC
CacheSequenceBoth   = 2  // FMC + SPC
```

**VRAM Calculation Formula**:
```
KV_per_token_per_layer = head_count_kv × (key_length + value_length) × bytes_per_element
KV_per_slot = n_ctx × n_layers × KV_per_token_per_layer
Total_slots = slots + cache_sequences
Slot_memory = total_slots × KV_per_slot
Total_VRAM = model_size + slot_memory
```

**Example Calculation**:
For Qwen3-Coder-30B (36GB, 128k context, q8_0 cache):
- No caching (2 slots): ~48.8GB total VRAM
- FMC (3 slots): ~55.2GB total VRAM
- FMC+SPC (4 slots): ~61.6GB total VRAM

**Benefits**:
- Helps users understand memory requirements before loading models
- Supports different cache configurations
- Works with both local and HuggingFace models

---

### 3. API Endpoint

**Endpoint**: `POST /v1/vram/calculate`

**Authentication**: Wallet-based authentication required

**Request Body**:
```json
{
  "model_id": "qwen3-coder-30b-a3b-instruct-ud-q8_k_xl",
  "context_window": 131072,
  "bytes_per_element": 1,
  "slots": 2,
  "cache_sequences": 2
}
```

**Response**:
```json
{
  "model_id": "qwen3-coder-30b-a3b-instruct-ud-q8_k_xl",
  "model_size_bytes": 38654705664,
  "context_window": 131072,
  "block_count": 48,
  "head_count_kv": 4,
  "key_length": 128,
  "value_length": 128,
  "bytes_per_element": 1,
  "slots": 2,
  "cache_sequences": 2,
  "kv_per_token_per_layer": 1024,
  "kv_per_slot": 6442450944,
  "total_slots": 4,
  "slot_memory": 25769803776,
  "total_vram": 64424509440
}
```

**Example cURL**:
```bash
curl -X POST http://localhost:8080/v1/vram/calculate \
  -H "Content-Type: application/json" \
  -H "X-Wallet-Address: 0x..." \
  -H "X-Signature: ..." \
  -H "X-Message: ..." \
  -d '{
    "model_id": "qwen3-coder-30b",
    "context_window": 131072,
    "bytes_per_element": 1,
    "slots": 2,
    "cache_sequences": 2
  }'
```

---

### 4. Helper Functions

**`NormalizeHuggingFaceDownloadURL(url string) string`**

Converts short HuggingFace URLs to full download URLs.

**Input**: `mradermacher/Qwen2-Audio-7B-GGUF/Qwen2-Audio-7B.Q8_0.gguf`

**Output**: `https://huggingface.co/mradermacher/Qwen2-Audio-7B-GGUF/resolve/main/Qwen2-Audio-7B.Q8_0.gguf`

---

### 5. Additional Enhancements

Features implemented for robustness and developer experience:

- **Defensive Batch Processing (`pkg/kronk/model/batch.go`)**: 
  - Checks for batch overflow before processing
  - Logs detailed per-slot state on overflow
  - Fails slots gracefully instead of crashing

- **Sampling Parameters (`pkg/kronk/model/params.go`)**:
  - Exported default constants (e.g., `DefTemp`, `DefTopK`)
  - Improved documentation for all sampling parameters

- **Enhanced Shutdown Logic**:
  - Improved `stop` method in batch engine to properly wait for goroutines
  - (Note: Main cache shutdown logic was already robust in Veridium)

---

## Files Modified

1. **`pkg/tools/models/info.go`** (new) - Model metadata reader
2. **`pkg/tools/models/vram.go`** (new) - VRAM calculator
3. **`pkg/tools/models/models.go`** - Added `NormalizeHuggingFaceDownloadURL`
4. **`cmd/server/app/domain/toolapp/model.go`** - Added VRAM request/response models
5. **`cmd/server/app/domain/toolapp/toolapp.go`** - Added `calculateVRAM` handler
6. **`cmd/server/app/domain/toolapp/route.go`** - Added VRAM endpoint route
7. **`pkg/kronk/model/batch.go`** - Added defensive batch processing & shutdown fixes
8. **`pkg/kronk/model/params.go`** - Refactored constants and improved documentation

**Total Changes**: 1000+ insertions across 8 files

---

## Testing

To test the implementation:

1. **Build the server**:
   ```bash
   go build ./cmd/server
   ```

2. **Start the server**:
   ```bash
   ./server
   ```

3. **Test VRAM calculation**:
   ```bash
   curl -X POST http://localhost:8080/v1/vram/calculate \
     -H "Content-Type: application/json" \
     -d '{
       "model_id": "your-model-id",
       "context_window": 8192,
       "bytes_per_element": 1,
       "slots": 2,
       "cache_sequences": 0
     }'
   ```

---

## Future Enhancements (Not Yet Implemented)

No significant features remain unimplemented from the Kronk analysis scope.

---

## References

- Kronk repository commits: a92184d9..HEAD (22 commits analyzed)
- Key commits:
  - `4e04550` - Loading model information directly from model file
  - `bc645e7` - VRAM calculator
  - `7db1e80` - Race and shutdown fixes
  - `044c489` - Defensive logging improvements

---

## Commit

```
commit 480ceb01
feat: Add VRAM calculator and model metadata reader

Implemented features from kronk repository (commits a92184d9..HEAD)
```
