# Panduan Migrasi: veridium/pkg/whisper → kawai-network/whisper

## Ringkasan

Dokumen ini membahas cara migrasi dari paket `veridium/pkg/whisper` (berbasis whisper.cpp) ke paket `kawai-network/whisper` (berbasis gowhisper C++ library) dalam server Veridium.

## Perbandingan Singkat

| Fitur | veridium/pkg/whisper | kawai-network/whisper |
|-------|---------------------|---------------------|
| Backend | whisper.cpp | gowhisper (custom C++) |
| Library Download | Manual | Auto-download |
| Model Format | `ggml-{name}.bin` | `{name}.bin` |
| Model Management | Built-in | Custom helper |
| VAD Support | ❌ | ✅ |
| Speaker Diarization | ❌ | ✅ |
| GPU Support | ✅ | ✅ |
| FFmpeg Dependency | ✅ | ✅ |

## Alasan Migrasi

1. **Fitur Tambahan**: Voice Activity Detection (VAD) dan Speaker Diarization
2. **Performance Optimized**: gowhisper library yang lebih teroptimasi
3. **Independensi**: Paket standalone yang tidak bergantung pada struktur veridium
4. **Auto-download**: Library gowhisper didownload otomatis jika tidak ditemukan

## Prasyarat

### Software yang Diperlukan

- FFmpeg (untuk konversi audio)
  - macOS: `brew install ffmpeg`
  - Linux: `sudo apt install ffmpeg`
  - Windows: Download dari https://ffmpeg.org/

### Hardware

- CPU: Multi-core disarankan
- RAM: Minimal 2GB (untuk model base)
- GPU: Opsional (NVIDIA/AMD dengan driver yang sesuai)

## Langkah-Langkah Migrasi

### 1. Update Dependencies

```bash
cd /path/to/veridium
go get github.com/kawai-network/whisper
```

### 2. Update Import di whisperapp.go

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
)
```

### 3. Perubahan API Utama

#### Inisialisasi Whisper

**Sebelum (veridium/pkg/whisper):**
```go
cfg := whisper.Config{
    ModelPath: modelPath,
    UseGPU:    true,
}
w, err := whisper.New(cfg)
defer w.Free()
```

**Sesudah (kawai-network/whisper):**
```go
// Library akan didownload otomatis jika tidak ditemukan
w, err := whisper.New(libDir)
defer w.Close()

// Load model secara terpisah
err := w.Load(modelPath)
```

#### Konfigurasi Transkripsi

**Sebelum:**
```go
opts := []whisper.TranscribeOption{
    whisper.WithThreads(4),
    whisper.WithLanguage("id"),
    whisper.WithTranslate(),
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
    Prompt:    "",
}
result, err := w.Transcribe(audioFile, opts)
```

#### Hasil Transkripsi

**Sebelum:**
```go
type Result struct {
    Text     string
    Segments []Segment
    Language string
}

type Segment struct {
    Text  string
    Start time.Duration
    End   time.Duration
}
```

**Sesudah:**
```go
type TranscriptionResult struct {
    Segments []*Segment
    Text     string
}

type Segment struct {
    Id     int32
    Text   string
    Start  int64  // dalam nanosecond
    End    int64  // dalam nanosecond
    Tokens []int32
}
```

### 4. Format File Model

**Perubahan penting:** Format nama file model berbeda!

| Sebelum | Sesudah | Perintah Konversi |
|---------|---------|-------------------|
| `ggml-base.bin` | `base.bin` | `mv ggml-base.bin base.bin` |
| `ggml-small.bin` | `small.bin` | `mv ggml-small.bin small.bin` |

### 5. Setup Environment (Otomatis)

Path untuk Whisper sudah dikonfigurasi otomatis menggunakan `internal/paths`:

```go
// Di code, gunakan paths.WhisperModels() dan paths.WhisperLib()
modelsDir := paths.WhisperModels()  // e.g., ./data/models/whisper
libDir := paths.WhisperLib()        // e.g., ./data/libraries/whisper
```

**Development vs Production:**
- **Development**: `./data/models/whisper` dan `./data/libraries/whisper`
- **Production (macOS)**: `~/Library/Application Support/Kawai/models/whisper`
- **Production (Windows)**: `%APPDATA%\Kawai\models\whisper`  
- **Production (Linux)**: `~/.config/Kawai/models/whisper`

Tidak perlu set environment variables WHISPER_MODELS_DIR dan WHISPER_LIB_DIR lagi!

### 6. Download Model yang Diperlukan

Gunakan helper yang sudah disediakan:

```go
// Di setup script atau command
import "github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"

// Download model dengan progress
err := whisperapp.DownloadModelWithLogger(context.Background(), "base", "./data/models/whisper")
if err != nil {
    log.Fatal(err)
}
```

Atau gunakan CLI:

```bash
# Download model base
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin -o ./data/models/whisper/base.bin

# Atau gunakan script setup (jika ada)
./scripts/setup-whisper.sh base
```

## Fitur Baru yang Tersedia

### 1. Voice Activity Detection (VAD)

```go
// Load model VAD (gunakan model yang sama)
err := w.LoadVAD(modelPath)
if err != nil {
    return err
}

// Convert audio ke float32 samples
samples, err := convertAudioToFloat32(audioFile)

// Jalankan VAD
vadSegments, err := w.VAD(samples)
for _, seg := range vadSegments {
    fmt.Printf("Voice: %.2fs - %.2fs\n", seg.Start, seg.End)
}
```

### 2. Speaker Diarization

```go
opts := whisper.TranscriptionOptions{
    Threads:   4,
    Language:  "",
    Translate: false,
    Diarize:   true,  // Enable speaker diarization
    Prompt:    "",
}

result, err := w.Transcribe(audioFile, opts)

// Check untuk speaker turn
for _, segment := range result.Segments {
    if strings.Contains(segment.Text, "[SPEAKER_TURN]") {
        fmt.Println("Speaker berubah!")
    }
}
```

### 3. Detailed Token Information

```go
result, err := w.Transcribe(audioFile, opts)

for _, segment := range result.Segments {
    fmt.Printf("Segment %d:\n", segment.Id)
    fmt.Printf("  Text: %s\n", segment.Text)
    fmt.Printf("  Start: %d ns\n", segment.Start)
    fmt.Printf("  End: %d ns\n", segment.End)
    fmt.Printf("  Tokens: %v\n", segment.Tokens)
}
```

## Testing Setelah Migrasi

### 1. Unit Tests

Buat test untuk memastikan fungsi bekerja:

```go
func TestWhisperApp_Transcribe(t *testing.T) {
    // Setup
    cfg := Config{Log: logger.New()}
    app := newApp(cfg)
    
    // Test transcription dengan file audio
    file, err := os.Open("test_audio.wav")
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()
    
    result, err := app.transcribe(context.Background(), "model_path", file, "id")
    if err != nil {
        t.Fatalf("Transcribe failed: %v", err)
    }
    
    if result == "" {
        t.Error("Expected non-empty result")
    }
}
```

### 2. Integration Tests

```bash
# Test endpoint
curl -X POST http://localhost:8080/v1/audio/transcriptions \
  -F "file=@test_audio.wav" \
  -F "model=base" \
  -F "language=id"
```

### 3. Performance Testing

Bandingkan performance sebelum dan sesudah migrasi:

```go
start := time.Now()
result, err := w.Transcribe(audioFile, opts)
duration := time.Since(start)

log.Printf("Transcription took %v", duration)
log.Printf("Result length: %d characters", len(result.Text))
```

## Troubleshooting

### Masalah: Library gowhisper tidak ditemukan

**Error:**
```
Library not found in ./data/libraries/whisper, attempting to download...
failed to open library: ...
```

**Solusi:**
Pastikan internet connection aktif untuk auto-download. Library akan didownload otomatis ke path dari `paths.WhisperLib()`.

```bash
# Library akan didownload otomatis
# Path dikelola oleh internal/paths:
# - Dev: ./data/libraries/whisper
# - Prod: ~/Library/Application Support/Kawai/libraries/whisper (macOS)
```

### Masalah: Model tidak ditemukan

**Error:**
```
model base not found at ./data/models/whisper/base.bin
```

**Solusi:**
Download model ke path dari `paths.WhisperModels()`:

```bash
# Path otomatis: ./data/models/whisper (dev) atau ~/Library/Application Support/Kawai/models/whisper (prod)
curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin \
  -o ./data/models/whisper/base.bin
```

### Masalah: FFmpeg tidak ditemukan

**Error:**
```
ffmpeg failed: exec: "ffmpeg": executable file not found in $PATH
```

**Solusi:**
Install FFmpeg:

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt update
sudo apt install ffmpeg

# Verifikasi
ffmpeg -version
```

### Masalah: Out of Memory

**Error:**
```
failed to create whisper: out of memory
```

**Solusi:**
- Gunakan model yang lebih kecil (tiny/base)
- Kurangi jumlah threads
- Tambah swap space jika perlu

```go
opts := whisper.TranscriptionOptions{
    Threads:   2,  // Kurangi threads
    Language:  "",
    Translate: false,
    Diarize:   false,
}
```

### Masalah: Hasil transkripsi kosong

**Possible Causes:**
1. File audio tidak valid
2. Model tidak di-load dengan benar
3. Format audio tidak didukung

**Solusi:**
```go
// Validasi file audio sebelum transkripsi
info, err := os.Stat(audioFile)
if err != nil || info.Size() == 0 {
    return "", fmt.Errorf("invalid audio file")
}

// Pastikan model berhasil di-load
if err := w.Load(modelPath); err != nil {
    return "", fmt.Errorf("failed to load model: %w", err)
}

// Tambahkan logging untuk debugging
log.Printf("Audio file size: %d bytes", info.Size())
log.Printf("Model path: %s", modelPath)
```

## Rollback Plan

Jika terjadi masalah setelah migrasi, Anda bisa rollback dengan mudah:

### 1. Restore Import

```go
import (
    "github.com/kawai-network/veridium/pkg/whisper"
    "github.com/kawai-network/veridium/pkg/whisper/model"
)
```

### 2. Restore Code

Gunakan git untuk restore versi sebelumnya:

```bash
git checkout HEAD~1 cmd/server/app/domain/whisperapp/whisperapp.go
```

### 3. Restore Dependencies

```bash
go mod tidy
```

## Checklist Migrasi

- [ ] Update go.mod dengan dependency baru
- [ ] Update imports di whisperapp.go
- [ ] Update inisialisasi Whisper
- [ ] Update konfigurasi transkripsi
- [ ] Update path file model (ggml-{name}.bin → {name}.bin)
- [ ] Verify path management (using internal/paths, no manual env vars needed)
- [ ] Download model yang diperlukan
- [ ] Test unit tests
- [ ] Test integration tests
- [ ] Test endpoint API
- [ ] Verify performance
- [ ] Update documentation
- [ ] Deploy ke staging environment
- [ ] Monitor logs untuk error
- [ ] Deploy ke production

## Contact & Support

Jika Anda mengalami masalah selama migrasi:

1. Cek log files untuk error details
2. Verifikasi semua environment variables
3. Pastikan semua dependencies terinstall
4. Jalankan tests untuk memastikan fungsi bekerja

## Appendix

### A. Daftar Model yang Tersedia

- `tiny` - 39M params, 1GB RAM
- `tiny.en` - 39M params, 1GB RAM (English only)
- `base` - 74M params, 2GB RAM
- `base.en` - 74M params, 2GB RAM (English only)
- `small` - 244M params, 4GB RAM
- `small.en` - 244M params, 4GB RAM (English only)
- `medium` - 769M params, 8GB RAM
- `medium.en` - 769M params, 8GB RAM (English only)
- `large-v1` - 1550M params, 16GB RAM
- `large-v2` - 1550M params, 16GB RAM
- `large-v3` - 1550M params, 16GB RAM (recommended)
- `large-v3-turbo` - 809M params, 8GB RAM (faster)

### B. Environment Variables Reference

**Catatan:** Environment variables WHISPER_MODELS_DIR dan WHISPER_LIB_DIR sudah **tidak diperlukan lagi**. Path dikelola otomatis oleh `internal/paths`.

| Variable | Default | Description | Status |
|----------|---------|-------------|--------|
| `WHISPER_LIB_DIR` | `./data/libraries/whisper` | Direktori library gowhisper | **Deprecated** - Gunakan `paths.WhisperLib()` |
| `WHISPER_MODELS_DIR` | `./data/models/whisper` | Direktori model files | **Deprecated** - Gunakan `paths.WhisperModels()` |
| `WHISPER_THREADS` | `4` | Jumlah threads untuk transkripsi | **Opsional** - Bisa set via kode |

### C. Contoh Konfigurasi Production

```yaml
# config.yaml
whisper:
  # Path dikelola otomatis oleh internal/paths
  # Development: ./data/models/whisper, ./data/libraries/whisper
  # Production: ~/Library/Application Support/Kawai/models/whisper (macOS)
  default_model: "base"
  threads: 4
  enable_gpu: true
  
  # Model download settings
  auto_download: false
  verify_checksums: true
  
  # Performance tuning
  cache_enabled: true
  max_concurrent: 3
  timeout: "10m"
```

---

**Dokumen ini diperbarui pada: 2025-01-15**
**Versi: 1.0.0**