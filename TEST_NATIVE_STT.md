# Test Native STT - Quick Guide

## 🎤 Option 1: Manual Recording (Paling Mudah)

### Step 1: Record Audio dengan macOS

```bash
# Record 5 detik audio dari microphone
# Tekan Ctrl+C untuk stop sebelum 5 detik
sox -d -r 16000 -c 1 test_voice.wav trim 0 5

# Atau gunakan QuickTime Player:
# 1. Buka QuickTime Player
# 2. File > New Audio Recording
# 3. Klik record, ngomong sesuatu
# 4. Stop, save as test_voice.wav
```

### Step 2: Test dengan Go

```bash
# Buat test file
cat > test_stt.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "github.com/kawai-network/veridium/services"
)

func main() {
    // Create STT service
    stt, err := services.NewNativeSTTService("en-US")
    if err != nil {
        log.Fatalf("Failed to create STT: %v", err)
    }
    defer stt.Close()
    
    // Check if available
    if !stt.IsAvailable() {
        log.Fatal("STT not available")
    }
    
    fmt.Println("🎤 Transcribing audio...")
    
    // Transcribe
    text, err := stt.TranscribeFile("test_voice.wav")
    if err != nil {
        log.Fatalf("Failed to transcribe: %v", err)
    }
    
    fmt.Printf("✅ Result: %s\n", text)
}
EOF

# Run
CGO_ENABLED=1 go run test_stt.go
```

## 🎙️ Option 2: Real-time Mic Input (Advanced)

Untuk real-time mic input, perlu extend `native_stt_service.go`:

### Add Real-time Method

```go
// Add to native_stt_service.go

// StartRealTimeTranscription starts transcribing from microphone
func (s *NativeSTTService) StartRealTimeTranscription(callback func(text string, isFinal bool)) error {
    // Implementation needed - requires AVAudioEngine setup
    // This is more complex, needs audio buffer handling
}
```

## 🚀 Quick Test Commands

### Test dengan TTS (Text-to-Speech) dulu

```bash
# Buat audio file dengan TTS
say -o test_voice.wav "Hello, this is a test of speech recognition"

# Test transcribe
go run test_stt.go
```

### Test dengan Voice Anda

```bash
# Install sox (jika belum)
brew install sox

# Record 5 detik
sox -d -r 16000 -c 1 my_voice.wav trim 0 5

# Ngomong sesuatu saat recording!
# Lalu transcribe
CGO_ENABLED=1 go run test_stt.go
```

## ⚠️ Requirements

1. **Permission**: Pertama kali akan minta permission di System Settings
2. **Audio Format**: WAV, 16kHz, Mono
3. **macOS**: 10.15+ required

## 📝 Expected Output

```
🎤 Transcribing audio...
✅ Result: Hello, this is a test of speech recognition
```

## 🔧 Troubleshooting

### Permission Denied
```bash
# Grant permission:
# System Settings > Privacy & Security > Speech Recognition
# Enable for Terminal/Your App
```

### Audio Format Error
```bash
# Convert audio to correct format
ffmpeg -i input.mp3 -ar 16000 -ac 1 output.wav
```

