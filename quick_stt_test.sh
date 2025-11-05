#!/bin/bash

echo "🎤 Quick Native STT Test"
echo "========================"
echo ""

# Check if audio file exists
if [ ! -f "test_voice.aiff" ]; then
    echo "📝 Creating test audio with TTS..."
    say -o test_voice.aiff "Hello, this is a test of native speech recognition on macOS. The weather is nice today and I am testing the speech to text feature."
    echo "✅ Test audio created: test_voice.aiff"
    echo ""
fi

# Play the audio
echo "🔊 Playing test audio..."
afplay test_voice.aiff &
sleep 1
echo ""

# Create simple Go test
cat > /tmp/quick_stt.go << 'GOEOF'
package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"github.com/kawai-network/veridium/services"
)

func main() {
	fmt.Println("🎤 Transcribing...")
	
	// Get absolute path
	audioPath, _ := filepath.Abs("test_voice.aiff")
	
	// Create STT service
	stt, err := services.NewNativeSTTService("en-US")
	if err != nil {
		log.Fatalf("❌ Failed: %v\n\nℹ️  Note: First time will ask for permission in System Settings", err)
	}
	defer stt.Close()
	
	// Transcribe
	text, err := stt.TranscribeFile(audioPath)
	if err != nil {
		log.Fatalf("❌ Transcription failed: %v", err)
	}
	
	// Show result
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 TRANSCRIPTION RESULT:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n%s\n\n", text)
	fmt.Println(strings.Repeat("=", 60))
}
GOEOF

# Run test
echo "🚀 Running transcription..."
echo ""
cd "$(dirname "$0")"
CGO_ENABLED=1 go run /tmp/quick_stt.go

# Cleanup
rm -f /tmp/quick_stt.go

echo ""
echo "✅ Test complete!"
echo ""
echo "💡 To test with your own voice:"
echo "   1. Record: say -o my_voice.aiff \"your text here\""
echo "   2. Or use QuickTime: File > New Audio Recording"
echo "   3. Run: CGO_ENABLED=1 go run /tmp/quick_stt.go"

