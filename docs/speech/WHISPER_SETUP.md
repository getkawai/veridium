# Whisper Integration Setup

## ✅ Setup Complete

go-whisper telah berhasil diintegrasikan ke Veridium sebagai submodule.

## 📁 Struktur

```
veridium/
├── go-whisper/                    # Git submodule
│   ├── build/install/lib/         # Compiled libraries
│   ├── third_party/whisper.cpp/   # whisper.cpp source
│   └── ...
├── services/
│   └── whisper_service.go         # Whisper service untuk Wails
└── main.go                        # WhisperService registered
```

## 🔧 Build Instructions

### 1. Clone dengan Submodules

```bash
git clone --recursive https://github.com/kawai-network/veridium.git
cd veridium
```

Atau jika sudah clone:

```bash
git submodule update --init --recursive
```

### 2. Build whisper.cpp Libraries

```bash
cd go-whisper
make libwhisper libffmpeg
cd ..
```

### 3. Build Veridium

```bash
# Set PKG_CONFIG_PATH untuk whisper libraries
export PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig

# Build dengan CGO enabled
CGO_ENABLED=1 go build -o bin/veridium
```

### 4. Atau gunakan Task

```bash
# Build untuk macOS
task darwin:build

# Build production
task darwin:package
```

## 📦 WhisperService API

### Available Methods

```go
// List all downloaded models
ListModels() []*schema.Model

// Get specific model
GetModel(id string) *schema.Model

// Download model from HuggingFace
DownloadModel(ctx context.Context, modelName string) error

// Delete model
DeleteModel(id string) error

// Transcribe audio file to text
Transcribe(ctx context.Context, modelId, audioPath string) (string, error)

// Transcribe with timestamps
TranscribeWithSegments(ctx context.Context, modelId, audioPath string) ([]*schema.Segment, error)

// Get list of available models for download
GetAvailableModels() []map[string]interface{}

// Get models directory path
GetModelsDirectory() string
```

### Frontend Usage (TypeScript)

```typescript
import { services } from '@@/github.com/kawai-network/veridium/services';

// List models
const models = await services.WhisperService.ListModels();

// Download a model
await services.WhisperService.DownloadModel('ggml-base.bin');

// Transcribe audio
const text = await services.WhisperService.Transcribe(
  'ggml-base',
  '/path/to/audio.wav'
);

// Get available models
const availableModels = await services.WhisperService.GetAvailableModels();
```

## 📝 Recommended Models

| Model | Size | Speed | Accuracy | Use Case |
|-------|------|-------|----------|----------|
| `ggml-tiny.bin` | 75 MB | Fastest | Low | Quick testing |
| `ggml-base.bin` | 142 MB | Fast | Good | **Recommended** |
| `ggml-small.bin` | 466 MB | Medium | Better | High quality |
| `ggml-medium.bin` | 1.5 GB | Slow | High | Professional |
| `ggml-large-v3.bin` | 3.1 GB | Very slow | Best | Maximum accuracy |

## 🎯 Audio Requirements

- **Format**: WAV (16-bit PCM)
- **Sample Rate**: 16kHz (auto-resampled)
- **Channels**: Mono (auto-converted)

For other formats (MP3, M4A, etc.), convert to WAV first using ffmpeg:

```bash
ffmpeg -i input.mp3 -ar 16000 -ac 1 output.wav
```

## 🚀 GPU Acceleration

### macOS (Metal)
GPU acceleration is **enabled by default** on Apple Silicon (M1/M2/M3).

### Linux (CUDA)
```bash
cd go-whisper
GGML_CUDA=1 make libwhisper libffmpeg
cd ..
```

### Linux (Vulkan)
```bash
cd go-whisper
GGML_VULKAN=1 make libwhisper libffmpeg
cd ..
```

## 📍 Models Storage Location

Models are stored in:
- **macOS**: `~/Library/Application Support/veridium/whisper-models/`
- **Linux**: `~/.config/veridium/whisper-models/`
- **Windows**: `%APPDATA%\veridium\whisper-models\`

## 🔍 Troubleshooting

### Error: "whisper.cpp not found"
```bash
cd go-whisper
make libwhisper libffmpeg
```

### Error: "pkg-config not found"
```bash
# macOS
brew install pkg-config

# Linux
sudo apt install pkg-config
```

### Error: "CMake not found"
```bash
# macOS
brew install cmake

# Linux
sudo apt install cmake
```

### Build fails with CGO errors
Make sure PKG_CONFIG_PATH is set:
```bash
export PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig
```

## 📚 Resources

- [go-whisper GitHub](https://github.com/mutablelogic/go-whisper)
- [whisper.cpp GitHub](https://github.com/ggerganov/whisper.cpp)
- [Whisper Models (HuggingFace)](https://huggingface.co/ggerganov/whisper.cpp)

## ✅ Integration Status

- [x] go-whisper setup as submodule
- [x] whisper.cpp libraries built
- [x] WhisperService created
- [x] Service registered in main.go
- [x] TypeScript bindings generated
- [ ] Frontend UI for transcription (TODO)
- [ ] Model management UI (TODO)

## 🎉 Next Steps

1. Generate TypeScript bindings:
   ```bash
   make bindings-generate
   ```

2. Create frontend UI for:
   - Model download/management
   - Audio file upload
   - Transcription display
   - Real-time transcription progress

3. Add features:
   - Batch transcription
   - Language detection
   - Speaker diarization
   - Export to SRT/VTT subtitles

