# Test Resume Download

Tool untuk membuktikan bahwa grab library support resume download.

## Build

```bash
go build -o bin/test-resume-download ./cmd/test-resume-download
```

## Test Scenario

### Test 1: Resume Berfungsi (Tanpa Cleanup)

```bash
# Download file besar (akan di-interrupt)
./bin/test-resume-download \
  "https://huggingface.co/bartowski/Qwen_Qwen3-VL-8B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-8B-Instruct-Q4_K_M.gguf" \
  "/tmp/test-qwen.gguf"

# Tunggu beberapa detik, lalu tekan Ctrl+C
# File partial akan tersimpan di /tmp/test-qwen.gguf

# Jalankan lagi command yang sama
./bin/test-resume-download \
  "https://huggingface.co/bartowski/Qwen_Qwen3-VL-8B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-8B-Instruct-Q4_K_M.gguf" \
  "/tmp/test-qwen.gguf"

# ✅ Seharusnya muncul "RESUMED" dan melanjutkan dari byte terakhir
```

### Test 2: Resume Tidak Berfungsi (Dengan Cleanup)

```bash
# Download dengan cleanup (simulasi behavior saat ini)
./bin/test-resume-download \
  "https://huggingface.co/bartowski/Qwen_Qwen3-VL-8B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-8B-Instruct-Q4_K_M.gguf" \
  "/tmp/test-qwen2.gguf"

# Tekan Ctrl+C setelah beberapa detik

# Hapus file (simulasi cleanup)
rm /tmp/test-qwen2.gguf

# Jalankan lagi
./bin/test-resume-download \
  "https://huggingface.co/bartowski/Qwen_Qwen3-VL-8B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-8B-Instruct-Q4_K_M.gguf" \
  "/tmp/test-qwen2.gguf"

# ❌ Download mulai dari 0 lagi (tidak ada "RESUMED")
```

## Expected Output

### Saat Resume Berfungsi:
```
📦 Found existing file: /tmp/test-qwen.gguf (512.50 MB)
🔄 Will attempt to resume download...

🚀 Downloading: https://huggingface.co/...
📁 Destination: /tmp/test-qwen.gguf

🔄 RESUMED | 15.2% | 800.5/5245.2 MB | 12.34 MB/s | ETA: 6m 12s
```

### Saat Resume Tidak Berfungsi:
```
📥 Starting fresh download...

🚀 Downloading: https://huggingface.co/...
📁 Destination: /tmp/test-qwen2.gguf

📥 DOWNLOADING | 5.0% | 262.3/5245.2 MB | 10.50 MB/s | ETA: 8m 30s
```

## Indicators

- `🔄 RESUMED` = Resume berfungsi, melanjutkan dari byte terakhir
- `📥 DOWNLOADING` = Download fresh dari awal
- `resp.DidResume` = Flag internal grab yang menunjukkan resume terjadi
