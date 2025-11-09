# GGUF Manager vs Model Downloader Comparison

## TL;DR

**Both are needed!** They serve different purposes:

- **`gguf_manager.go`** → **`embedding_manager.go`**: Manages **embedding models** (384-768 dims, 68MB-550MB)
- **`model_downloader.go`**: Manages **LLM models** (0.5B-7B params, GB range)

## Side-by-Side Comparison

| Feature | `embedding_manager.go` (from gguf_manager) | `model_downloader.go` |
|---------|-------------------------------------------|----------------------|
| **Purpose** | Semantic search, RAG embeddings | Text generation, chat |
| **Model Type** | Embedding models | LLM models |
| **Model Size** | 68MB - 550MB | 500MB - 5GB+ |
| **Dimensions** | 384-768 | N/A (generates text) |
| **Download Method** | Direct HTTP from HuggingFace | llama-cli integration |
| **Validation** | ✅ GGUF header validation | ⚠️ File size check only |
| **SHA256 Verification** | ✅ Yes | ❌ No |
| **Progress Tracking** | ✅ Byte-level | ⚠️ Limited |
| **Auto-Download** | ✅ Quality-based | ✅ RAM-based |
| **Model Selection** | Manual + recommended | Automatic based on RAM |
| **Storage Path** | `~/.veridium/models/embeddings/` | `~/.veridium/models/` |
| **Catalog Size** | 5 models | 4 models |
| **Languages** | 12-18+ languages | Multilingual |
| **Use Case** | Vector search, similarity | Chat, generation |

## Models Catalog

### Embedding Models (embedding_manager.go)

| Model | Size | Dims | Quant | Languages | Use Case |
|-------|------|------|-------|-----------|----------|
| **granite-embedding-107m-q6_k_l** | 120MB | 384 | Q6_K_L | 18+ | ✅ **Recommended** |
| granite-embedding-107m-q4_k_m | 123MB | 384 | Q4_K_M | 18+ | Good balance |
| granite-embedding-107m-f16 | 236MB | 384 | F16 | 18+ | Highest quality |
| paraphrase-multilingual-minilm-f16 | 242MB | 384 | F16 | 12+ | Alternative |
| nomic-embed-text-v1.5-q4_k_m | 550MB | 768 | Q4_K_M | 18+ | High quality |

### LLM Models (model_downloader.go)

| Model | Size | Params | Quant | Min RAM | Use Case |
|-------|------|--------|-------|---------|----------|
| Qwen2.5-0.5B-Instruct | ~500MB | 0.5B | Q4_K_M | 2GB | Testing, low-end |
| Qwen2.5-1.5B-Instruct | ~1.2GB | 1.5B | Q4_K_M | 4GB | Lightweight |
| Qwen2.5-3B-Instruct | ~2.5GB | 3B | Q4_K_M | 6GB | ✅ **Recommended** |
| Qwen2.5-7B-Instruct | ~5GB | 7B | Q4_K_M | 10GB | High quality |

## Integration Benefits

### Before Integration
```
❌ Separate, unrelated files
❌ No unified API
❌ Duplicate code for downloads
❌ Inconsistent error handling
```

### After Integration
```
✅ Unified llama.Service
✅ Consistent API for both model types
✅ Shared infrastructure
✅ Better organization
```

## Usage Examples

### Embedding Models

```go
// Initialize service
llamaService, _ := llama.NewService()

// Get embedding manager
embMgr := llamaService.GetEmbeddingManager()

// List available models
for name, model := range embMgr.GetAvailableModels() {
    fmt.Printf("%s: %.1f MB, %d dims\n", 
        name, 
        float64(model.Size)/1024/1024, 
        model.Dimensions)
}

// Download recommended model
modelName := llamaService.GetRecommendedEmbeddingModel()
err := llamaService.DownloadEmbeddingModel(modelName, func(dl, total int64) {
    fmt.Printf("Progress: %.1f%%\n", float64(dl)/float64(total)*100)
})

// Get model path
modelPath, _ := embMgr.GetModelPath(modelName)

// Start llama-server with embedding model
cmd := exec.Command("llama-server",
    "-m", modelPath,
    "--port", "8080",
    "--embedding",
)
cmd.Start()
```

### LLM Models

```go
// Initialize service
llamaService, _ := llama.NewService()

// Auto-download based on RAM
err := llamaService.AutoDownloadRecommendedModel()

// Get available models
models, _ := llamaService.GetAvailableModels()
for _, modelPath := range models {
    fmt.Println("Model:", modelPath)
}

// Start llama-server with LLM model
err = llamaService.StartServer(modelPath, 8081)
```

## When to Use Which?

### Use `embedding_manager.go` when:
- ✅ Building semantic search
- ✅ Implementing RAG (Retrieval Augmented Generation)
- ✅ Creating vector databases
- ✅ Doing similarity matching
- ✅ Need multilingual embeddings

### Use `model_downloader.go` when:
- ✅ Building chat applications
- ✅ Generating text responses
- ✅ Need conversational AI
- ✅ Doing text completion
- ✅ Need reasoning capabilities

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                      llama.Service                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │         LlamaCppReleaseManager                        │ │
│  │  - Manages llama.cpp binaries                         │ │
│  │  - llama-server, llama-cli, llama-embedding           │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │         ModelDownloader (LLM Models)                  │ │
│  │  - Qwen 0.5B-7B                                       │ │
│  │  - Auto-download via llama-cli                        │ │
│  │  - RAM-based selection                                │ │
│  │  - Storage: ~/.veridium/models/                       │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │         EmbeddingManager (Embedding Models)           │ │
│  │  - Granite, Paraphrase, Nomic                         │ │
│  │  - Direct HTTP download                               │ │
│  │  - GGUF validation + SHA256                           │ │
│  │  - Storage: ~/.veridium/models/embeddings/            │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                           │
                           │ provides models to
                           ▼
        ┌──────────────────────────────────────┐
        │     chromem Vector Database          │
        │  - Uses embedding models             │
        │  - Semantic search                   │
        │  - RAG operations                    │
        └──────────────────────────────────────┘
```

## File Structure

```
internal/llama/
├── service.go                    (Main service, integrates both)
├── manager.go                    (LlamaCppReleaseManager)
├── model_downloader.go           (LLM models)
├── embedding_manager.go          (Embedding models) ← NEW
├── proxy.go                      (Proxy for streaming)
└── hardware.go                   (Hardware detection)

~/.veridium/
├── models/
│   ├── qwen2.5-0.5b-instruct-q4_k_m.gguf    (LLM)
│   ├── qwen2.5-3b-instruct-q4_k_m.gguf      (LLM)
│   └── embeddings/
│       ├── granite-embedding-107m-multilingual-Q6_K_L.gguf  ← NEW
│       └── nomic-embed-text-v1.5.Q4_K_M.gguf               ← NEW
└── downloads/
    └── *.tmp (temporary files)
```

## Performance Comparison

### Download Speed

| Method | Embedding Manager | Model Downloader |
|--------|------------------|------------------|
| **Protocol** | Direct HTTP | llama-cli (HTTP) |
| **Speed** | ~10-50 MB/s | ~5-30 MB/s |
| **Resume** | ❌ No | ❌ No |
| **Progress** | ✅ Byte-level | ⚠️ Limited |

### Validation

| Check | Embedding Manager | Model Downloader |
|-------|------------------|------------------|
| **File Size** | ✅ Yes | ✅ Yes |
| **GGUF Header** | ✅ Yes | ❌ No |
| **SHA256** | ✅ Optional | ❌ No |
| **Corruption Detection** | ✅ Yes | ⚠️ Limited |

### Storage Efficiency

| Aspect | Embedding Manager | Model Downloader |
|--------|------------------|------------------|
| **Model Size** | 68MB-550MB | 500MB-5GB |
| **Quantization** | Multiple options | Single (Q4_K_M) |
| **Cleanup** | ✅ Automatic | ⚠️ Manual |
| **Deduplication** | ❌ No | ❌ No |

## Migration Guide

### From Standalone gguf_manager.go

If you were using `gguf_manager.go` standalone:

**Before:**
```go
gm := main.NewGGUFManager()
models := gm.GetAvailableModels()
gm.DownloadModel("granite-embedding-107m-multilingual-q6_k_l", nil)
```

**After:**
```go
llamaService, _ := llama.NewService()
embMgr := llamaService.GetEmbeddingManager()
models := embMgr.GetAvailableModels()
llamaService.DownloadEmbeddingModel("granite-embedding-107m-multilingual-q6_k_l", nil)
```

### Integration Checklist

- [x] Create `embedding_manager.go` from `gguf_manager.go`
- [x] Add `EmbeddingManager` to `llama.Service`
- [x] Add accessor methods to `llama.Service`
- [x] Update package name from `main` to `llama`
- [x] Change storage path to `.veridium`
- [x] Test compilation
- [x] Document integration

## Conclusion

**Both managers are complementary, not redundant:**

1. **`embedding_manager.go`** (from `gguf_manager.go`):
   - ✅ Manages **embedding models** for semantic search
   - ✅ GGUF validation and SHA256 verification
   - ✅ Direct HTTP downloads with progress tracking
   - ✅ 5 high-quality multilingual models

2. **`model_downloader.go`**:
   - ✅ Manages **LLM models** for text generation
   - ✅ llama-cli integration for easy downloads
   - ✅ RAM-based auto-selection
   - ✅ 4 Qwen models (0.5B-7B)

Together, they provide **complete model management** for the Veridium application! 🎉

---

**Date**: 2025-11-09  
**Status**: ✅ Integrated  
**Recommendation**: Keep both, they serve different purposes

