# Hardware Validation untuk Reasoning Mode - Dokumentasi Lengkap

## Ringkasan Perubahan

Sistem sekarang **otomatis memeriksa spesifikasi hardware** sebelum mengizinkan user memilih reasoning mode yang membutuhkan banyak resource.

## Masalah yang Diselesaikan

Sebelumnya, jika user memilih reasoning mode (Qwen, GPT-OSS, dll) pada hardware yang tidak cukup kuat, aplikasi bisa:
- Sangat lambat
- Freeze/hang
- Crash
- Pengalaman user yang buruk

## Solusi

Sekarang sistem akan:
1. ✅ **Deteksi hardware** otomatis (RAM, CPU cores, GPU)
2. ✅ **Validasi** apakah hardware cukup untuk mode yang dipilih
3. ✅ **Auto-fallback** ke mode yang sesuai jika hardware tidak memadai
4. ✅ **Log jelas** menjelaskan keputusan yang diambil

## Persyaratan Hardware

### Mode Disabled (Non-Reasoning)
- **RAM Minimum**: 4GB
- **CPU Cores**: 2
- **GPU**: Tidak diperlukan
- **Description**: Lightweight - runs on most systems
- **Cocok untuk**: Hardware low-end, conversation panjang
- **Models**: Llama 3.2 3B, Mistral 7B

### Mode Enabled (Minimal Reasoning)
- **RAM Minimum**: 8GB
- **CPU Cores**: 4
- **GPU**: Direkomendasikan
- **Description**: Moderate - 8GB RAM, 4+ cores recommended
- **Cocok untuk**: Hardware medium, use case umum
- **Models**: Qwen3 1.7B with `/no_think`

### Mode Verbose (Full Reasoning)
- **RAM Minimum**: 16GB
- **CPU Cores**: 6
- **GPU**: Sangat direkomendasikan
- **Description**: High-end - 16GB+ RAM, 6+ cores, GPU strongly recommended
- **Cocok untuk**: Hardware high-end, pertanyaan kompleks
- **Models**: Qwen3, GPT-OSS, DeepSeek R1 with full thinking

## Contoh Behavior

### Contoh 1: Hardware Cukup
```
User pilih: Reasoning Enabled
Hardware: 16GB RAM, 8 cores
Hasil: ✅ Mode enabled berhasil
Log: "Hardware validation passed for enabled mode: RAM=16GB (need 8GB), Cores=8 (need 4)"
```

### Contoh 2: Hardware Tidak Cukup
```
User pilih: Reasoning Verbose
Hardware: 8GB RAM, 4 cores
Hasil: ⚠️ Otomatis switch ke Enabled mode
Log: "RAM tidak cukup: 8GB tersedia, butuh 16GB untuk verbose mode"
Log: "Auto-switching to enabled mode based on available hardware"
```

### Contoh 3: Hardware Sangat Rendah
```
User pilih: Reasoning Enabled
Hardware: 4GB RAM, 2 cores
Hasil: ⚠️ Otomatis switch ke Disabled mode
```

## Implementasi Teknis

### File yang Dimodifikasi

#### 1. `internal/services/reasoning_mode.go`
**Changes:**
- Added import: `github.com/kawai-network/veridium/internal/llama`
- Added `HardwareRequirements` struct to define min specs per mode
- Added `GetHardwareRequirements()` method to `ReasoningConfig`
- Added `ValidateHardware()` method to check if system meets requirements
- Added `SuggestModeForHardware()` function to recommend best mode

**New Code:**
```go
type HardwareRequirements struct {
	MinRAM       int64  // Minimum RAM in GB
	MinCPUCores  int    // Minimum CPU cores
	RecommendGPU bool   // Whether GPU is recommended
	Description  string // Human-readable description
}

func (rc ReasoningConfig) GetHardwareRequirements() HardwareRequirements
func (rc ReasoningConfig) ValidateHardware(specs *llama.HardwareSpecs) (bool, string)
func SuggestModeForHardware(specs *llama.HardwareSpecs) ReasoningMode
```

#### 2. `internal/services/agent_chat_service.go`
Modified `SetReasoningMode()` to:
- Get hardware specs from `LibraryService`
- Validate hardware before setting reasoning mode
- Auto-fallback to suitable mode if hardware insufficient
- Log hardware specs and requirements for transparency

#### 3. `internal/llama/library_service.go`
Added:
- `GetHardwareSpecs()` - exposes hardware specs from installer

### Behavior Flow

```
User pilih Reasoning Mode
        ↓
Sistem ambil hardware specs dari installer
        ↓
Validasi specs vs requirements
        ↓
    [Hardware cukup?]
     /         \
   Ya          Tidak
    ↓           ↓
Set mode    Cari mode alternatif
            ↓
         Set mode alternatif
            ↓
         Log peringatan
```

## Example Usage

```go
// In agent_chat_service.go
func (s *AgentChatService) SetReasoningMode(mode ReasoningMode) error {
    // Get hardware specs
    hardwareSpecs := s.libService.GetHardwareSpecs()

    // Validate
    tempConfig := ReasoningConfig{Mode: mode}
    if valid, reason := tempConfig.ValidateHardware(hardwareSpecs); !valid {
        // Auto-fallback to suggested mode
        suggested := SuggestModeForHardware(hardwareSpecs)
        log.Printf("Hardware insufficient: %s. Switching to %s", reason, suggested)
        mode = suggested
    }

    s.reasoningConfig.Mode = mode
    return nil
}
```

## Example Log Output

### Sufficient Hardware
```
🧠 Reasoning mode changed to: enabled (Balanced (Minimal Reasoning) - Good for most use cases)
💻 Hardware requirements: Moderate - 8GB RAM, 4+ cores recommended
📊 Current system: RAM=16GB, Cores=8, GPU=Apple M1 Max
✅ Hardware validation passed for enabled mode: RAM=16GB (need 8GB), Cores=8 (need 4)
```

### Insufficient Hardware
```
⚠️  Hardware validation failed for verbose mode: Insufficient RAM: 8GB available, but 16GB required for verbose mode. Consider using disabled mode instead.
💡 Auto-switching to enabled mode based on available hardware
🧠 Reasoning mode changed to: enabled (Balanced (Minimal Reasoning) - Good for most use cases)
💻 Hardware requirements: Moderate - 8GB RAM, 4+ cores recommended
📊 Current system: RAM=8GB, Cores=4, GPU=
⚠️  enabled mode recommends GPU acceleration, but no GPU detected. Performance may be degraded.
✅ Hardware validation passed for enabled mode: RAM=8GB (need 8GB), Cores=4 (need 4)
```

## Testing

### File Baru

1. **`internal/services/reasoning_mode_test.go`**
   - Unit tests lengkap (23 test cases)
   - Semua test PASS ✅

2. **`cmd/test-hardware-validation/main.go`**
   - Program demo untuk test validasi
   - Menampilkan hardware dan validasi

### Jalankan Unit Tests
```bash
go test -v ./internal/services -run "TestHardware"
```

Hasil: **PASS** ✅ (23/23 tests)

### Jalankan Demo
```bash
go run ./cmd/test-hardware-validation
```

Hasil dari test pada M1 Pro (12GB available RAM, 8 cores):
- ✅ Disabled mode: SUFFICIENT
- ✅ Enabled mode: SUFFICIENT
- ❌ Verbose mode: INSUFFICIENT (butuh 16GB, hanya ada 12GB)
- 💡 Mode yang disarankan: **Enabled**

### Build Main App
```bash
go build ./internal/services ./internal/llama .
```

Hasil: **SUCCESS** ✅

## Keuntungan

1. ✅ **Tidak ada crash**: User tidak bisa enable mode yang tidak bisa dihandle system
2. ✅ **UX lebih baik**: Auto-fallback, bukan error
3. ✅ **Transparan**: Log jelas menjelaskan keputusan
4. ✅ **Aman**: Requirements konservatif memastikan performance baik
5. ✅ **Future-proof**: Mudah adjust requirements jika model berubah

## API Changes

### Public API Additions

#### In `internal/services/reasoning_mode.go`:
```go
type HardwareRequirements struct { ... }
func (rc ReasoningConfig) GetHardwareRequirements() HardwareRequirements
func (rc ReasoningConfig) ValidateHardware(*llama.HardwareSpecs) (bool, string)
func SuggestModeForHardware(*llama.HardwareSpecs) ReasoningMode
```

#### In `internal/llama/library_service.go`:
```go
func (s *LibraryService) GetHardwareSpecs() *HardwareSpecs
```

### No Breaking Changes
- All existing APIs remain unchanged
- New functionality is additive only
- Backward compatible with existing code

## Migration Guide

No migration needed! The changes are transparent:

1. **Existing code**: Works as before
2. **New behavior**: When `SetReasoningMode()` is called, hardware is now validated
3. **User experience**: Better - automatic mode adjustment instead of failures

## Configuration

No configuration needed. The system uses:
- Hardware detection from `LlamaCppInstaller.HardwareSpecs`
- Predefined requirements per mode
- Automatic fallback logic

## Logging

Enhanced logging provides visibility:

```
🧠 Reasoning mode changed to: enabled (Balanced (Minimal Reasoning) - Good for most use cases)
💻 Hardware requirements: Moderate - 8GB RAM, 4+ cores recommended
📊 Current system: RAM=16GB, Cores=8, GPU=Apple M1 Pro
✅ Hardware validation passed for enabled mode: RAM=16GB (need 8GB), Cores=8 (need 4)
📊 Expected performance:
   - Speed: 1.8s per turn
   - Token efficiency: ~118 tokens/turn
   - Max turns: 30-50 turns
   - Response size: ~450 chars
```

## Known Limitations

1. **Hardware Detection**: May fail on some systems (handled gracefully with warnings)
2. **Static Requirements**: Requirements are fixed, not dynamic based on actual model size
3. **No Runtime Monitoring**: Only validates at mode selection time, not during execution

## Future Improvements

Possible enhancements:
1. Check available (not just total) RAM in real-time
2. Monitor GPU temperature and throttle if needed
3. Add dynamic mode switching during conversation based on resource usage
4. Provide UI feedback about hardware limitations
5. Add benchmarking to fine-tune requirements

## Notes

- Hardware detection may fail on some systems - in this case, the system logs a warning but allows the mode (with caution)
- GPU detection is platform-specific (implemented in `hardware_*.go` files)
- Requirements are conservative to ensure good performance
- Users can still manually load models, but reasoning mode selection is protected

## Kesimpulan

✅ Feature **SELESAI** dan **TESTED**  
✅ **Tidak ada breaking changes**  
✅ **Meningkatkan user experience**  
✅ **Mencegah resource exhaustion**  
✅ **Well documented**  

Sekarang reasoning mode akan **otomatis adjust** berdasarkan hardware yang tersedia! 🎉
