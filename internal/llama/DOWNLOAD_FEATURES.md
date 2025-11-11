# Download Features - Grab Integration

## 🎉 Successfully Implemented!

Veridium now uses the **grab** library for model downloads (GGUF files), providing production-ready download capabilities with advanced features.

## 📦 Download Strategy

Veridium uses **two different download approaches** for different file types:

### 1. **Model Downloads (GGUF files)** → Uses `pkg/yzma/download.GetWithProgress()`
- Large files (100MB - 5GB+)
- Direct downloads (no extraction needed)
- Benefits from resume, progress tracking, rate limiting
- Examples: Qwen models, embedding models
- Unified API with binary downloads

### 2. **Binary Downloads (llama.cpp)** → Uses `pkg/yzma/download`
- Pre-built binaries in ZIP format
- Requires extraction after download
- Platform-specific URLs (darwin/metal, linux/cuda, etc.)
- Smaller files (~50-200MB)
- **Now uses grab internally** with custom ZIP extraction
- Automatic retry logic
- Resume support for binary downloads too!

This hybrid approach gives us the best of both worlds - grab for all downloads, with custom ZIP handling for binaries!

## ✅ Implemented Features

### 1. **Resume Download** 📦
- Automatically resumes interrupted downloads
- Uses HTTP Range requests
- No data loss if app closes during download
- Works seamlessly across restarts

```go
// Automatic resume is enabled by default
opts := DefaultDownloadOptions()
opts.ResumeIfPossible = true // Default
```

### 2. **Progress Tracking** 📊
- Real-time progress percentage
- Download speed in MB/s
- ETA (Estimated Time of Arrival)
- Stuck download detection

```go
// Progress updates every 2 seconds by default
📥 Progress: 45.2% (3.45 MB/s, ETA: 2m15s)
📥 Progress: 67.8% (3.52 MB/s, ETA: 1m23s)
✅ Download complete: model.gguf (468.2 MB)
```

### 3. **Automatic Retry** 🔄
- Retries up to 3 times on failure
- Exponential backoff (2s, 4s, 6s)
- Cleans up partial downloads between retries
- Clear error messages

```go
opts := DefaultDownloadOptions()
opts.MaxRetries = 3 // Default

// Example output:
⚠️  Attempt 1 failed: context deadline exceeded
🔄 Retry attempt 2/3 (waiting 4s)...
✅ Download complete: model.gguf (468.2 MB)
```

### 4. **Bandwidth Throttling** 🔧
- Optional rate limiting in MB/s
- Prevents network saturation
- Easy to configure

```go
// Limit to 5 MB/s
opts := WithRateLimit(5)

// Or manually:
opts := DefaultDownloadOptions()
opts.RateLimitMBps = 5

// Unlimited (default):
opts.RateLimitMBps = 0
```

### 5. **Batch Downloads** ⚡
- Download multiple files concurrently
- Configurable worker count
- Progress tracking per file

```go
urls := []string{
    "https://example.com/model1.gguf",
    "https://example.com/model2.gguf",
    "https://example.com/model3.gguf",
}

respch, _ := DownloadBatch(urls, "/download/dir", 3)
for resp := range respch {
    if err := resp.Err(); err != nil {
        log.Printf("Error: %v", err)
    } else {
        log.Printf("Downloaded: %s", resp.Filename)
    }
}
```

## 📝 API Usage

### Basic Download

```go
// Use default options (retry, resume, progress)
opts := DefaultDownloadOptions()
err := DownloadWithGrab(url, destPath, opts)
```

### Custom Options

```go
opts := DownloadOptions{
    MaxRetries:       5,              // More retries
    ShowProgress:     true,           // Enable progress
    ResumeIfPossible: true,           // Enable resume
    ProgressInterval: 1 * time.Second, // Update every second
    RateLimitMBps:    10,             // Limit to 10 MB/s
}
err := DownloadWithGrab(url, destPath, opts)
```

### Quick Rate Limiting

```go
// Limit to 3 MB/s
opts := WithRateLimit(3)
err := DownloadWithGrab(url, destPath, opts)
```

## 🔧 Integration Points

### In `installer.go`

Both `DownloadChatModel` and `DownloadEmbeddingModel` now use grab:

```go
func (lcm *LlamaCppInstaller) DownloadChatModel(modelSpec QwenModelSpec) error {
    // ... setup ...
    
    // Download using grab with automatic retry, resume, and progress tracking
    opts := DefaultDownloadOptions()
    if err := DownloadWithGrab(modelSpec.URL, tempModelPath, opts); err != nil {
        lcm.cleanupTempFile(tempModelPath)
        return fmt.Errorf("failed to download model: %w", err)
    }
    
    // ... validation ...
}
```

## 🎯 User Experience Improvements

### Before (with go-getter):
```
📥 Downloading model...
❌ Error: context deadline exceeded
```

### After (with grab):
```
📥 Downloading chat model: qwen2.5-0.5b-instruct-q4_k_m
   URL: https://huggingface.co/...
   Expected size: 468.0 MB
   This may take several minutes depending on network speed...
📥 Progress: 12.3% (2.45 MB/s, ETA: 2m45s)
📥 Progress: 34.7% (3.12 MB/s, ETA: 1m52s)
📥 Progress: 58.9% (3.45 MB/s, ETA: 1m05s)
📥 Progress: 89.2% (3.52 MB/s, ETA: 0m18s)
✅ Download complete: qwen2.5-0.5b-instruct-q4_k_m.gguf (468.2 MB)
```

### If Download Interrupted:
```
📥 Progress: 45.2% (3.45 MB/s, ETA: 2m15s)
^C (app closed)

// Next run:
📥 Downloading chat model: qwen2.5-0.5b-instruct-q4_k_m
   URL: https://huggingface.co/...
   Expected size: 468.0 MB
   This may take several minutes depending on network speed...
📥 Progress: 45.2% (3.45 MB/s, ETA: 2m15s)  ← Resumed from where it stopped!
   📦 Download resumed successfully
✅ Download complete: qwen2.5-0.5b-instruct-q4_k_m.gguf (468.2 MB)
```

### With Rate Limiting:
```
🔧 Rate limit: 5 MB/s
📥 Progress: 23.4% (5.00 MB/s, ETA: 3m12s)  ← Capped at 5 MB/s
📥 Progress: 56.7% (5.00 MB/s, ETA: 1m45s)
```

## 🧪 Testing

All tests passed successfully:
```bash
✅ PASSED: 20 tests
⏭️  SKIPPED: 33 tests (integration tests)
❌ FAILED: 0 tests

🎉 ALL TESTS PASSED!
```

## 📦 Dependencies

- `github.com/kawai-network/veridium/pkg/grab` - Main download library
- `golang.org/x/time/rate` - Rate limiting support

## 🚀 Performance

- **Throughput**: Full network speed (unless rate limited)
- **Memory**: Efficient streaming (32KB buffer by default)
- **CPU**: Minimal overhead
- **Disk**: Writes directly to file (no memory buffering)

## 🔒 Safety Features

1. **Atomic Operations**: Download to `.tmp` file, rename on success
2. **Cleanup**: Automatic cleanup of partial downloads on failure
3. **Validation**: File size, checksum, and GGUF format validation
4. **Stuck Detection**: Detects and fails stuck downloads
5. **Context-Aware**: Respects context cancellation

## 📚 Additional Resources

- [grab GitHub](https://github.com/cavaliergopher/grab)
- [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)

## 🎉 Summary

Veridium now has **production-ready download capabilities** with:
- ✅ Resume support for interrupted downloads
- ✅ Real-time progress tracking with speed and ETA
- ✅ Automatic retry with exponential backoff
- ✅ Optional bandwidth throttling
- ✅ Batch download support
- ✅ Comprehensive error handling
- ✅ All tests passing

No more lost downloads! No more manual retries! 🚀

