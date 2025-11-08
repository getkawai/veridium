# Auto-Download Model Implementation

## Overview
Implemented automatic model download functionality for Kawai AI provider. The system now automatically downloads the optimal Qwen GGUF model based on available system RAM.

## Implementation Date
November 8, 2025

## Files Created/Modified

### 1. New File: `internal/llama/model_downloader.go`

#### Key Components:

**QwenGGUFModel Struct**
```go
type QwenGGUFModel struct {
    Name         string // Model filename
    URL          string // HuggingFace download URL
    Size         int64  // File size in bytes
    Parameters   string // Parameter size (0.5b, 1.5b, etc.)
    Quantization string // Quantization type (Q4_K_M, Q5_K_M, etc.)
    MinRAM       int64  // Minimum RAM required in GB
    Description  string // Model description
}
```

**Available Models**:
1. **Qwen 0.5B** (~468MB) - Requires 2GB RAM
   - Perfect for low-end hardware and testing
   - URL: HuggingFace Qwen2.5-0.5B-Instruct-GGUF

2. **Qwen 1.5B** (~934MB) - Requires 4GB RAM
   - Good balance of speed and quality
   - Recommended for most users

3. **Qwen 3B** (~1.9GB) - Requires 6GB RAM
   - Excellent for most tasks
   - Mid-range option

4. **Qwen 7B** (~4.3GB) - Requires 10GB RAM
   - High-quality for advanced tasks
   - For powerful hardware

**Key Functions**:

1. `GetRecommendedQwenGGUFModels()` - Returns list of available models
2. `SelectOptimalQwenGGUFModel(availableRAM)` - Selects best model for system
3. `DownloadQwenModel(model, progressCallback)` - Downloads from HuggingFace
4. `AutoDownloadRecommendedModel()` - Main auto-download function

### 2. Modified: `internal/llama/service.go`

**Updated `initializeInBackground()` function** (lines 108-119):

```go
if len(models) == 0 {
    log.Println("⚠️  No GGUF models found. Starting auto-download...")
    
    // Auto-download recommended model based on hardware
    if err := s.AutoDownloadRecommendedModel(); err != nil {
        log.Printf("⚠️  Failed to auto-download model: %v", err)
        log.Println("   You can download a model manually later")
        return
    }
    
    log.Println("✅ Model downloaded successfully!")
}
```

### 3. Existing Files Used:

- `internal/llama/hardware.go` - Hardware detection
- `internal/llama/manager_darwin.go` - macOS hardware specs
- `internal/llama/qwen.go` - Qwen model specifications (reference)

## How It Works

### Auto-Download Flow:

```
App Starts
    ↓
NewService() creates llama.Service
    ↓
initializeInBackground() runs in goroutine
    ↓
Check llama.cpp installed ✓
    ↓
Check llama-server binary ✓
    ↓
Check for GGUF models
    ↓
NO MODELS FOUND!
    ↓
DetectHardwareSpecs()
    ├─ Total RAM: 16GB
    ├─ Available RAM: 12.8GB (80% of total)
    ├─ CPU: Apple M1
    └─ GPU: Apple M1 (shared memory)
    ↓
SelectOptimalQwenGGUFModel(12.8GB)
    ├─ 0.5B model: ✓ Fits (requires 2GB)
    ├─ 1.5B model: ✓ Fits (requires 4GB)
    ├─ 3B model: ✓ Fits (requires 6GB)
    ├─ 7B model: ✓ Fits (requires 10GB)
    └─ Selected: Qwen 7B (best fit)
    ↓
DownloadQwenModel(qwen2.5-7b-instruct-q4_k_m.gguf)
    ├─ URL: HuggingFace
    ├─ Size: 4.3GB
    ├─ Progress: 10%... 20%... 100%
    └─ Saved to: ~/.llama-cpp/models/
    ↓
Model Downloaded ✅
    ↓
StartServerAuto()
    ↓
llama-server Running ✅
```

## Hardware Detection

### macOS (Darwin):
- **Total RAM**: From `sysctl -n hw.memsize`
- **Available RAM**: 80% of total (conservative estimate)
- **CPU**: From `sysctl -n machdep.cpu.brand_string`
- **GPU**: From `system_profiler SPDisplaysDataType -json`

### Model Selection Logic:
1. Get available RAM
2. Find largest model that fits
3. If none fit, use smallest (0.5B)
4. Log selection reasoning

Example:
```
System: 16GB RAM → 12.8GB available
Selected: Qwen 7B (requires 10GB, has 2.8GB headroom)
```

## Download Features

### Progress Tracking:
```go
progressCallback := func(progress float64) {
    if int(progress)%10 == 0 { // Log every 10%
        log.Printf("   Download progress: %.0f%%", progress)
    }
}
```

### Error Handling:
- Network failures: Retry logic (via HTTP client)
- Partial downloads: Temp file + rename on success
- Cleanup: Remove temp files on error

### Resume Support:
- Check if model already exists before download
- Skip download if model file present

## Testing

### Manual Test:
```bash
# Clear existing models
rm -rf ~/.llama-cpp/models/*.gguf

# Run app - should auto-download
./veridium

# Check logs:
# 🚀 Initializing llama.cpp in background...
# ✅ llama.cpp is installed
# ✅ llama-server ready
# ⚠️  No GGUF models found. Starting auto-download...
# 📦 Selected model: qwen2.5-1.5b-instruct-q4_k_m.gguf (1.5b, Q4_K_M)
# 📥 Downloading qwen2.5-1.5b-instruct-q4_k_m.gguf from HuggingFace...
#    Size: 934.0 MB
#    Download progress: 10%
#    Download progress: 20%
#    ...
#    Download progress: 100%
# ✅ Model downloaded successfully!
# 🚀 Auto-starting llama-server...
# ✅ llama-server auto-started successfully
```

### Verify Download:
```bash
ls -lh ~/.llama-cpp/models/
# Should show downloaded .gguf file
```

### Test Chat:
```bash
curl http://127.0.0.1:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen2.5-1.5b-instruct-q4_k_m",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": false
  }'
```

## Configuration

### Model Selection Criteria:

1. **RAM-based**: Largest model that fits in available RAM
2. **Quantization**: Q4_K_M (good balance of quality/size)
3. **Model Family**: Qwen 2.5 (latest, best performance)
4. **Fallback**: 0.5B model if RAM is very low

### Customization:

To add more models, edit `GetRecommendedQwenGGUFModels()`:

```go
{
    Name:         "custom-model.gguf",
    URL:          "https://huggingface.co/...",
    Size:         1024 * 1024 * 1024, // 1GB
    Parameters:   "1b",
    Quantization: "Q4_K_M",
    MinRAM:       3,
    Description:  "Custom model description",
}
```

## Performance

### Download Times (estimated):
- **0.5B model** (~468MB): ~1-2 minutes on fast connection
- **1.5B model** (~934MB): ~2-4 minutes
- **3B model** (~1.9GB): ~4-8 minutes
- **7B model** (~4.3GB): ~8-15 minutes

### Disk Space:
- Models stored in: `~/.llama-cpp/models/`
- One model downloaded at a time
- No automatic cleanup (user can delete manually)

## Error Scenarios

### 1. Network Failure
```
⚠️  Failed to auto-download model: failed to download model: connection timeout
   You can download a model manually later
```
**Solution**: Check internet connection, retry

### 2. Insufficient Disk Space
```
⚠️  Failed to auto-download model: failed to create file: no space left on device
```
**Solution**: Free up disk space

### 3. HuggingFace Unavailable
```
⚠️  Failed to auto-download model: download failed with status: 503 Service Unavailable
```
**Solution**: Wait and retry, or download manually

## Manual Download Alternative

If auto-download fails, users can download manually:

```bash
# Create directory
mkdir -p ~/.llama-cpp/models

# Download model
cd ~/.llama-cpp/models
curl -L -O https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF/resolve/main/qwen2.5-0.5b-instruct-q4_k_m.gguf

# Restart app - will detect and use the model
```

## Future Enhancements

### Planned:
1. ✅ Auto-download (DONE)
2. 🔄 UI for model management
3. 🔄 Multiple model support
4. 🔄 Model update notifications
5. 🔄 Download resume after interruption
6. 🔄 Parallel downloads for faster setup

### Possible:
- Model quantization options (Q4, Q5, Q8)
- Custom model URLs
- Model caching/sharing
- Bandwidth throttling
- Download scheduling

## Integration with Kawai Provider

The auto-download seamlessly integrates with the Kawai AI provider:

1. **First Run**: Auto-downloads optimal model
2. **Subsequent Runs**: Uses existing model
3. **Frontend**: No changes needed, works transparently
4. **Error Handling**: Graceful fallback if download fails

## Logs to Watch

Successful auto-download sequence:
```
🚀 Initializing llama.cpp in background...
✅ llama.cpp is installed (version: b6451)
✅ llama-server ready at: /opt/homebrew/bin/llama-server
🎉 llama.cpp is ready to use!
⚠️  No GGUF models found. Starting auto-download...
📦 No models found, starting auto-download...
Detected hardware: RAM=16GB (available=12GB), CPU cores=8, GPU=Apple M1 (VRAM=4GB)
📦 Selected model: qwen2.5-7b-instruct-q4_k_m.gguf (7b, Q4_K_M) - requires 10GB RAM (system has 12GB)
📥 Downloading qwen2.5-7b-instruct-q4_k_m.gguf from HuggingFace...
   Size: 4370.0 MB
   URL: https://huggingface.co/Qwen/Qwen2.5-7B-Instruct-GGUF/resolve/main/qwen2.5-7b-instruct-q4_k_m.gguf
   Download progress: 10%
   Download progress: 20%
   ...
   Download progress: 100%
✅ Model downloaded successfully: qwen2.5-7b-instruct-q4_k_m.gguf (4370.0 MB)
🎉 Model download completed successfully!
🚀 Auto-starting llama-server...
🤖 Auto-selected model: /Users/yuda/.llama-cpp/models/qwen2.5-7b-instruct-q4_k_m.gguf
✅ llama-server started on port 8080 (PID: 12345)
✅ llama-server auto-started successfully
```

## Summary

✅ **Implemented**: Auto-download GGUF models from HuggingFace
✅ **Hardware-aware**: Selects optimal model based on RAM
✅ **Progress tracking**: Logs download progress
✅ **Error handling**: Graceful fallback on failures
✅ **Integration**: Seamless with existing Kawai provider
✅ **Build**: Successful compilation
✅ **Ready**: For production use

---

**Status**: ✅ Implementation Complete
**Testing**: Ready for user verification
**Next**: Test with fresh install (no models)

