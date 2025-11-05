# 🎤 Quick STT Test - Ngomong Langsung!

## ✅ Whisper Restored & Working!

Whisper STT sudah di-restore dan siap digunakan!

## 🎙️ Cara Test dengan Ngomong Langsung

### Step 1: Record Audio Anda

**Option A: QuickTime Player (Termudah)**
```bash
# 1. Buka QuickTime Player
# 2. File > New Audio Recording (Cmd+Option+N)
# 3. Klik tombol merah, ngomong sesuatu
# 4. Stop, lalu Export As > my_voice.m4a
```

**Option B: Terminal dengan sox**
```bash
# Install sox
brew install sox

# Record 10 detik (ngomong saat recording!)
sox -d -r 16000 -c 1 my_voice.wav trim 0 10
```

**Option C: Voice Memos App**
```bash
# 1. Buka Voice Memos
# 2. Record
# 3. Share > Save to Files > my_voice.m4a
```

### Step 2: Download Whisper Model

```bash
# Download tiny model (75MB, fastest)
go run test_whisper_download.go
```

Atau buat script:
```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/kawai-network/veridium/services"
)

func main() {
    svc, _ := services.NewWhisperService()
    defer svc.Close()
    
    fmt.Println("Downloading ggml-tiny.bin...")
    err := svc.DownloadModel(context.Background(), "ggml-tiny.bin")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("✅ Model downloaded!")
}
```

### Step 3: Transcribe!

```bash
# Buat test script
cat > test_whisper_quick.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/kawai-network/veridium/services"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run test_whisper_quick.go <audio_file>")
        os.Exit(1)
    }
    
    audioFile := os.Args[1]
    
    fmt.Println("🎤 Whisper STT Test")
    fmt.Println("===================\n")
    
    // Create service
    svc, err := services.NewWhisperService()
    if err != nil {
        log.Fatal(err)
    }
    defer svc.Close()
    
    // Check model
    models := svc.ListModels()
    if len(models) == 0 {
        fmt.Println("❌ No models found!")
        fmt.Println("\nDownload a model first:")
        fmt.Println("  go run test_whisper_download.go")
        os.Exit(1)
    }
    
    modelId := models[0].Id
    fmt.Printf("Using model: %s\n\n", modelId)
    
    // Transcribe
    fmt.Println("Transcribing... (may take 10-30 seconds)")
    text, err := svc.Transcribe(context.Background(), modelId, audioFile)
    if err != nil {
        log.Fatal(err)
    }
    
    // Show result
    fmt.Println("\n" + "="*60)
    fmt.Println("🎉 TRANSCRIPTION RESULT:")
    fmt.Println("="*60)
    fmt.Printf("\n%s\n\n", text)
    fmt.Println("="*60)
}
EOF

# Run it!
PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig \
CGO_ENABLED=1 go run test_whisper_quick.go my_voice.wav
```

## 🚀 All-in-One Script

```bash
#!/bin/bash

echo "🎤 Whisper STT - All-in-One Test"
echo "================================="
echo ""

# Check if audio file exists
if [ ! -f "my_voice.wav" ]; then
    echo "📝 No audio file found. Let's create one!"
    echo ""
    echo "Option 1: Use QuickTime (press Enter when done)"
    echo "Option 2: Record now with sox (10 seconds)"
    echo ""
    read -p "Choose (1/2): " choice
    
    if [ "$choice" = "2" ]; then
        echo "Recording in 3... 2... 1... SPEAK NOW!"
        sox -d -r 16000 -c 1 my_voice.wav trim 0 10
        echo "✅ Recording saved!"
    else
        echo "Open QuickTime > New Audio Recording"
        echo "Record, then export as my_voice.wav"
        read -p "Press Enter when ready..."
    fi
fi

# Download model if needed
if [ ! -d "$(go run -e 'svc,_:=services.NewWhisperService();fmt.Print(svc.GetModelsDirectory())')/ggml-tiny.bin" ]; then
    echo ""
    echo "📥 Downloading Whisper model..."
    go run test_whisper_download.go
fi

# Transcribe
echo ""
echo "🎤 Transcribing..."
PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig \
CGO_ENABLED=1 go run test_whisper_quick.go my_voice.wav
```

## 📊 Expected Performance

- **Tiny model**: ~5-10 seconds for 10 sec audio
- **Base model**: ~10-20 seconds for 10 sec audio
- **GPU acceleration**: Automatic on M1/M2 Macs

## 🎯 What to Say

Try these test phrases:
```
"Hello, this is a test of the Whisper speech recognition system."

"Halo, ini adalah tes sistem pengenalan suara Whisper."

"今日はいい天気ですね。" (Japanese)

"你好，这是语音识别测试。" (Chinese)
```

Whisper supports 99 languages! 🌍

## ✅ Advantages vs Native STT

| Feature | Whisper | Native STT |
|---------|---------|------------|
| **Permission** | ✅ None needed | ❌ Requires System Settings |
| **Languages** | ✅ 99 languages | ⚠️ 50+ languages |
| **Offline** | ✅ 100% offline | ✅ Offline |
| **Quality** | ✅ Excellent | ✅ Siri-level |
| **Speed** | ⚠️ 5-30 sec | ✅ Real-time capable |
| **Platform** | ✅ Cross-platform | ❌ macOS only |

## 🎉 Ready to Test!

Sekarang Anda bisa ngomong langsung dan lihat hasilnya! 🚀

