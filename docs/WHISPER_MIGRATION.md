# Whisper Service Migration

## Overview
Migrated from `go-whisper` (CGO bindings) to `whisper-cpp` CLI-based implementation.

## Why Migrate?

### Previous Implementation (go-whisper)
- ❌ Complex CGO dependencies
- ❌ Requires building `whisper.cpp` from source
- ❌ Platform-specific build issues
- ❌ Large binary size due to embedded libraries
- ❌ Difficult to cross-compile

### New Implementation (whisper-cpp CLI)
- ✅ Simple CLI wrapper
- ✅ No CGO dependencies
- ✅ Easy installation via package managers
- ✅ Smaller binary size
- ✅ Easy to cross-compile
- ✅ Better separation of concerns

## Architecture

### File Structure
```
internal/whisper/
├── manager.go              # Core manager logic (platform-agnostic)
├── manager_darwin.go       # macOS installation (Homebrew)
├── manager_linux.go        # Linux installation instructions
├── manager_windows.go      # Windows installation instructions
├── manager_unsupported.go  # Fallback for other platforms
└── service.go              # Wails-compatible service wrapper
```

### Build Tags
Each platform-specific file uses Go build tags:
- `//go:build darwin` - macOS only
- `//go:build linux` - Linux only
- `//go:build windows` - Windows only
- `//go:build !darwin && !linux && !windows` - Unsupported platforms

Go compiler automatically selects the correct file based on `GOOS`.

## Installation

### macOS
```bash
brew install whisper-cpp
```

### Linux
```bash
# Option 1: Homebrew for Linux
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
brew install whisper-cpp

# Option 2: Build from source
git clone https://github.com/ggml-org/whisper.cpp
cd whisper.cpp
make
```

### Windows
Download pre-built binaries from: https://github.com/ggml-org/whisper.cpp/releases

## Usage

### Automatic Setup (Recommended)
The service automatically handles installation and model download in the background:

```go
// Just create the service - it handles everything automatically!
service, _ := whisper.NewService()
// ✅ Auto-installs whisper-cpp if not found
// ✅ Auto-downloads recommended model (base) if no models exist
// ✅ All happens in background, non-blocking
```

### Manual Model Management (Optional)
```go
service, _ := whisper.NewService()

// Download a specific model
ctx := context.Background()
err := service.DownloadModel(ctx, "small")

// List available models
models := service.GetAvailableModels()
// Returns: tiny, base, small, medium, large-v3, etc.

// Check if model is downloaded
isDownloaded := service.IsModelDownloaded("base")
```

### Transcribe Audio
```go
// Simple transcription
text, err := service.Transcribe(ctx, "base", "/path/to/audio.wav")

// With custom options
options := map[string]interface{}{
    "language": "en",    // Language code
    "threads":  4,       // Number of threads
    "translate": false,  // Translate to English
}
text, err := service.TranscribeWithOptions(ctx, "base", "/path/to/audio.wav", options)
```

## API Comparison

### Old API (go-whisper)
```go
// Complex initialization
w, err := whisper.New(modelsDir, whisper.OptMaxConcurrent(2))

// Model management
models := w.ListModels()  // Returns []*schema.Model
model := w.GetModelById("ggml-base")

// Transcription requires loading audio samples
samples, _ := loadAudioSamples(audioPath)
w.WithModel(model, func(t *task.Context) error {
    return t.Transcribe(ctx, 0, samples, callback)
})
```

### New API (whisper-cpp CLI)
```go
// Simple initialization
service, err := whisper.NewService()

// Model management
models, _ := service.ListModels()  // Returns []string
isInstalled := service.IsModelDownloaded("base")

// Direct transcription
text, err := service.Transcribe(ctx, "base", audioPath)
```

## Benefits

1. **Simpler Codebase**: No CGO, no complex audio loading
2. **Easier Deployment**: Just install `whisper-cpp` CLI
3. **Better Portability**: Works on any platform with `whisper-cpp`
4. **Smaller Binary**: No embedded libraries
5. **Easier Testing**: Can test with mock CLI
6. **Platform-Specific Code**: Clean separation using build tags
7. **Auto-Setup**: Automatically installs whisper-cpp and downloads models in background
8. **Zero Configuration**: Works out of the box, no manual setup required

## Migration Checklist

- [x] Create new `Manager` with platform-specific files
- [x] Create `Service` wrapper for Wails compatibility
- [x] Update `main.go` to use new service
- [x] Remove `go-whisper` dependencies from `go.mod`
- [x] Remove old test files
- [x] Add unsupported platform stub
- [x] Test compilation on macOS
- [x] Add automatic whisper-cpp installation
- [x] Add automatic model download in background
- [x] Update documentation
- [ ] Test on Linux
- [ ] Test on Windows
- [ ] Update frontend bindings
- [ ] Test STT feature end-to-end

## Notes

- Audio files must be in WAV format (whisper-cpp requirement)
- Models are stored in `~/.kawai-agent/whisper-models/`
- First transcription may be slower (model loading)
- GPU acceleration is automatic if available

