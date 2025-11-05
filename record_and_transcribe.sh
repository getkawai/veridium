#!/bin/bash

echo "🎤 Record & Transcribe"
echo "====================="
echo ""
echo "Siap record? Press Ctrl+C untuk stop recording"
echo ""
read -p "Press Enter to start recording..."
echo ""
echo "🔴 RECORDING... (Press Ctrl+C when done)"

# Record dari mic ke WAV
# -t wav: output format
# -: output to stdout, then pipe to file
rec -r 16000 -c 1 -b 16 recording.wav

echo ""
echo "✅ Recording saved to recording.wav"
echo ""
echo "🎯 Transcribing..."

# Set environment
export PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig
export CGO_ENABLED=1

# Transcribe
cat > /tmp/transcribe_recording.go << 'GOEOF'
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
    
    models := svc.ListModels()
    if len(models) == 0 {
        log.Fatal("❌ No models found")
    }
    
    modelId := models[0].Id
    fmt.Printf("Using model: %s\n\n", modelId)
    
    start := time.Now()
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    
    text, err := svc.Transcribe(ctx, modelId, "recording.wav")
    if err != nil {
        log.Fatalf("❌ Failed: %v", err)
    }
    
    duration := time.Since(start)
    
    fmt.Println(strings.Repeat("=", 60))
    fmt.Println("🎉 TRANSCRIPTION:")
    fmt.Println(strings.Repeat("=", 60))
    fmt.Printf("\n%s\n\n", strings.TrimSpace(text))
    fmt.Println(strings.Repeat("=", 60))
    fmt.Printf("⏱️  Time: %.2f seconds\n", duration.Seconds())
}
GOEOF

go run /tmp/transcribe_recording.go

echo ""
echo "✅ Done!"

