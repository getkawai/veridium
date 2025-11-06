#!/bin/bash

set -e

echo "🎤 Whisper STT Test Helper"
echo "=========================="
echo ""

# No special environment needed for whisper-cpp CLI

# Step 1: Create test audio with TTS
echo "1️⃣  Creating test audio with macOS TTS..."
say -o test_voice.aiff "Hello, this is a test of the Whisper speech recognition system. The weather is nice today and I am testing speech to text."
echo "✅ Test audio created"
echo ""

# Step 2: Check Whisper service
echo "2️⃣  Checking Whisper service..."
cat > /tmp/test_whisper_check.go << 'GOEOF'
package main

import (
    "fmt"
    "log"
    "github.com/kawai-network/veridium/services"
)

func main() {
    svc, err := services.NewWhisperService()
    if err != nil {
        log.Fatalf("❌ Failed to create service: %v", err)
    }
    defer svc.Close()
    
    fmt.Printf("✅ Whisper service OK\n")
    fmt.Printf("   Models dir: %s\n", svc.GetModelsDirectory())
    
    models, err := svc.ListModels()
    if err != nil {
        log.Fatalf("❌ Failed to list models: %v", err)
    }
    fmt.Printf("   Installed models: %d\n", len(models))

    if len(models) == 0 {
        fmt.Println("\n⚠️  No models installed yet")
        fmt.Println("   Available models:")
        for _, m := range svc.GetAvailableModels() {
            fmt.Printf("   - %s (%s): %s\n", m["name"], m["size"], m["description"])
        }
    } else {
        fmt.Println("   Models:")
        for _, m := range models {
            fmt.Printf("   - %s\n", m)
        }
    }
}
GOEOF

go run /tmp/test_whisper_check.go
echo ""

# Step 3: Download model if needed
echo "3️⃣  Checking for Whisper model..."
cat > /tmp/test_whisper_download.go << 'GOEOF'
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/kawai-network/veridium/services"
)

func main() {
    svc, err := services.NewWhisperService()
    if err != nil {
        log.Fatal(err)
    }
    defer svc.Close()
    
    models, err := svc.ListModels()
    if err != nil {
        log.Fatalf("❌ Failed to list models: %v", err)
    }
    if len(models) > 0 {
        fmt.Printf("✅ Model already installed: %s\n", models[0])
        return
    }
    
    fmt.Println("📥 Downloading ggml-tiny.bin (75 MB)...")
    fmt.Println("   This will take 1-2 minutes...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    err = svc.DownloadModel(ctx, "ggml-tiny.bin")
    if err != nil {
        log.Fatalf("❌ Download failed: %v", err)
    }
    
    fmt.Println("✅ Model downloaded successfully!")
}
GOEOF

go run /tmp/test_whisper_download.go
echo ""

# Step 4: Transcribe!
echo "4️⃣  Transcribing audio..."
cat > /tmp/test_whisper_transcribe.go << 'GOEOF'
package main

import (
    "context"
    "fmt"
    "log"
    "strings"
    "time"
    "github.com/kawai-network/veridium/services"
)

func main() {
    svc, err := services.NewWhisperService()
    if err != nil {
        log.Fatal(err)
    }
    defer svc.Close()
    
    models, err := svc.ListModels()
    if err != nil {
        log.Fatalf("❌ Failed to list models: %v", err)
    }
    if len(models) == 0 {
        log.Fatal("❌ No models found")
    }

    modelId := models[0]
    fmt.Printf("Using model: %s\n", modelId)
    fmt.Println("Transcribing... (this may take 10-30 seconds)")
    
    start := time.Now()
    
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    
    text, err := svc.Transcribe(ctx, modelId, "test_voice.aiff")
    if err != nil {
        log.Fatalf("❌ Transcription failed: %v", err)
    }
    
    duration := time.Since(start)
    
    // Show result
    fmt.Println("\n" + strings.Repeat("=", 60))
    fmt.Println("🎉 TRANSCRIPTION RESULT:")
    fmt.Println(strings.Repeat("=", 60))
    fmt.Printf("\n%s\n\n", strings.TrimSpace(text))
    fmt.Println(strings.Repeat("=", 60))
    fmt.Printf("\n⏱️  Time taken: %.2f seconds\n", duration.Seconds())
}
GOEOF

go run /tmp/test_whisper_transcribe.go
echo ""

echo "✅ Test complete!"
echo ""
echo "💡 To test with your own voice:"
echo "   1. Record with QuickTime: File > New Audio Recording"
echo "   2. Save as my_voice.aiff"
echo "   3. Run: go run /tmp/test_whisper_transcribe.go"
echo "      (Edit the audio file path to 'my_voice.aiff')"
