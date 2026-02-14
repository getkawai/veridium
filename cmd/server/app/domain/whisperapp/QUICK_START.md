# Quick Start: Migrasi ke kawai-network/whisper

## Ringkasan Cepat

Migrasi dari `veridium/pkg/whisper` ke `kawai-network/whisper` hanya membutuhkan 3 langkah utama:

1. ✅ Update dependencies
2. ✅ Setup directory dan download model
3. ✅ Update konfigurasi

## Langkah 1: Update Dependencies

```bash
cd /path/to/veridium
go get github.com/kawai-network/whisper
go mod tidy
```

## Langkah 2: Quick Setup

### Opsi A: Setup Interaktif (Disarankan untuk Pertama Kali)

```bash
# Jalankan setup interaktif
go run cmd/server/main.go setup whisper
```

Setup ini akan:
- Membuat directory yang diperlukan
- Mengonversi model lama ke format baru
- Mendownload model yang diperlukan
- Memverifikasi setup

### Opsi B: Setup Cepat (Non-Interaktif)

```bash
# Setup dengan model base saja
go run cmd/server/main.go setup whisper --model base
```

### Opsi C: Manual Setup

```bash
# 1. Buat directories (path otomatis dari internal/paths)
mkdir -p ./data/models/whisper
mkdir -p ./data/libraries/whisper

# 2. Download model (pilih salah satu)
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin \
  -o ./data/models/whisper/base.bin

# Note: Path dikelola otomatis oleh internal/paths, tidak perlu export env vars
```

## Langkah 3: Konfigurasi (Otomatis)

Path untuk Whisper models dan library sudah dikonfigurasi secara otomatis menggunakan `internal/paths`:

- **Development**: `./data/models/whisper` dan `./data/libraries/whisper`
- **Production (macOS)**: `~/Library/Application Support/Kawai/models/whisper`
- **Production (Windows)**: `%APPDATA%\Kawai\models\whisper`
- **Production (Linux)**: `~/.config/Kawai/models/whisper`

Tidak perlu set environment variables manual lagi!

## Memulai Server

```bash
# Build server
go build -o bin/server cmd/server/main.go

# Jalankan server
./bin/server start
```

## Testing

### Test Endpoint

```bash
# Test transcription endpoint
curl -X POST http://localhost:8080/v1/audio/transcriptions \
  -F "file=@test_audio.wav" \
  -F "model=base" \
  -F "language=id"
```

### Test dengan Go Code

```go
import "github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"

// Setup whisper
ctx := context.Background()
err := whisperapp.QuickSetup(ctx, "base")
if err != nil {
    log.Fatal(err)
}

// Model sudah siap untuk digunakan!
```

## Konversi Model dari Format Lama

Jika Anda sudah punya model di format lama (`ggml-{name}.bin`):

```go
import "github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"

// Konversi satu model
err := whisperapp.ConvertModel("./data/models/whisper", "base")

// Konversi semua model
converted, errors := whisperapp.ConvertAllModels("./data/models/whisper", false)
fmt.Printf("Converted %d models\n", len(converted))
```

## Perintah CLI yang Tersedia

```bash
# Setup lengkap (interaktif)
./bin/server setup whisper

# Setup cepat dengan model tertentu
./bin/server setup whisper --model base

# Setup dengan multiple models
./bin/server setup whisper --models tiny,base,small

# Setup untuk production
./bin/server setup whisper --production --model base

# Diagnose setup
./bin/server setup whisper --diagnose

# Convert models
./bin/server setup whisper --convert

# List downloaded models
./bin/server setup whisper --list-models
```

## Perbedaan Utama dalam Code

### Import

**Sebelum:**
```go
import (
    "github.com/kawai-network/veridium/pkg/whisper"
    "github.com/kawai-network/veridium/pkg/whisper/model"
)
```

**Sesudah:**
```go
import (
    whisper "github.com/kawai-network/whisper"
    "github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"
)
```

### Inisialisasi

**Sebelum:**
```go
cfg := whisper.Config{
    ModelPath: modelPath,
    UseGPU:    true,
}
w, err := whisper.New(cfg)
defer w.Free()
```

**Sesudah:**
```go
w, err := whisper.New(libDir)
defer w.Close()
err := w.Load(modelPath)
```

### Transkripsi

**Sebelum:**
```go
opts := []whisper.TranscribeOption{
    whisper.WithThreads(4),
    whisper.WithLanguage("id"),
}
result, err := w.Transcribe(audioPath, opts...)
```

**Sesudah:**
```go
opts := whisper.TranscriptionOptions{
    Threads:   4,
    Language:  "id",
    Translate: false,
    Diarize:   false,
}
result, err := w.Transcribe(audioFile, opts)
```

## Fitur Baru

### 1. Voice Activity Detection (VAD)

```go
// Load VAD model
err := w.LoadVAD(modelPath)

// Perform VAD
vadSegments, err := w.VAD(audioSamples)
```

### 2. Speaker Diarization

```go
opts := whisper.TranscriptionOptions{
    Threads:   4,
    Language:  "",
    Diarize:   true,  // Enable speaker detection
}
```

### 3. Detailed Segments

```go
result, err := w.Transcribe(audioFile, opts)
for _, segment := range result.Segments {
    fmt.Printf("[%d] %s (%d - %d ns)\n", 
        segment.Id, segment.Text, segment.Start, segment.End)
}
```

## Troubleshooting Cepat

### Masalah: Library tidak ditemukan

**Error:**
```
Library not found in ./data/libraries/whisper, attempting to download...
```

**Solusi:**
Pastikan internet connection aktif. Library akan didownload otomatis ke path yang dikelola oleh `internal/paths`.

### Masalah: Model tidak ditemukan

**Error:**
```
model base not found at ./data/models/whisper/base.bin
```

**Solusi:**
```bash
# Download model
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin \
  -o ./data/models/whisper/base.bin
```

### Masalah: FFmpeg tidak ditemukan

**Solusi:**
```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg
```

### Masalah: Out of Memory

**Solusi:**
- Gunakan model yang lebih kecil (tiny/base)
- Kurangi jumlah threads
- Tambah RAM atau swap

```bash
# Gunakan model tiny (path otomatis dari internal/paths)
# Download model tiny
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin \
  -o ./data/models/whisper/tiny.bin
```

## Model Recommendation

| Use Case | Model yang Disarankan |
|----------|---------------------|
| Testing/Development | `tiny` atau `base` |
| Production (Standar) | `base` atau `small` |
| High Quality | `medium` |
| Best Quality | `large-v3` |
| Balance Quality/Speed | `large-v3-turbo` |

### Download Model Specific

```bash
# Path otomatis: ./data/models/whisper (development)
# Tiny (39MB, fastest)
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin \
  -o ./data/models/whisper/tiny.bin

# Base (141MB, recommended)
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin \
  -o ./data/models/whisper/base.bin

# Small (465MB, better quality)
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin \
  -o ./data/models/whisper/small.bin
```

## Checklist Sebelum Production

- [ ] Setup sudah dijalankan dan berhasil
- [ ] Model yang diperlukan sudah didownload
- [ ] FFmpeg sudah terinstall
- [ ] Server bisa dimulai tanpa error
- [ ] Endpoint API sudah diuji
- [ ] Performance sudah diukur
- [ ] Logging sudah dikonfigurasi
- [ ] Error handling sudah diuji
- [ ] Backup sudah dibuat
- [ ] Path management verified (using internal/paths)

## Next Steps

1. **Testing**: Jalankan comprehensive tests
   ```bash
   go test ./app/domain/whisperapp/...
   ```

2. **Performance Testing**: Ukur performance transkripsi
   ```bash
   ab -n 100 -c 10 -p test_audio.json -T multipart/form-data \
     http://localhost:8080/v1/audio/transcriptions
   ```

3. **Monitoring**: Setup monitoring untuk production
   - Memory usage
   - Transcription latency
   - Error rates
   - GPU usage (jika applicable)

4. **Documentation**: Update API documentation
5. **Deployment**: Deploy ke staging environment

## Bantuan Tambahan

- Lihat `MIGRATION_GUIDE.md` untuk panduan lengkap
- Lihat `model_helper.go` untuk API helper functions
- Lihat `setup.go` untuk setup utilities
- Lihat `convert_models.go` untuk konversi model

## Contoh Penggunaan Lengkap

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    whisper "github.com/kawai-network/whisper"
    "github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"
)

func main() {
    ctx := context.Background()
    
    // Setup
    err := whisperapp.QuickSetup(ctx, "base")
    if err != nil {
        fmt.Printf("Setup failed: %v\n", err)
        os.Exit(1)
    }
    
    // Create whisper instance
    libDir := "./data/lib/whisper"
    w, err := whisper.New(libDir)
    if err != nil {
        fmt.Printf("Failed to create whisper: %v\n", err)
        os.Exit(1)
    }
    defer w.Close()
    
    // Load model
    modelPath := "./data/models/whisper/base.bin"
    if err := w.Load(modelPath); err != nil {
        fmt.Printf("Failed to load model: %v\n", err)
        os.Exit(1)
    }
    
    // Transcribe
    opts := whisper.TranscriptionOptions{
        Threads:   4,
        Language:  "id",
        Translate: false,
        Diarize:   false,
    }
    
    result, err := w.Transcribe("audio.wav", opts)
    if err != nil {
        fmt.Printf("Transcription failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Transcription: %s\n", result.Text)
}
```

---

**Selamat migrasi! 🚀**

Jika Anda mengalami masalah, cek `MIGRATION_GUIDE.md` untuk troubleshooting lebih detail.