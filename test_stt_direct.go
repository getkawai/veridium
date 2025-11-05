package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kawai-network/veridium/services"
)

func main() {
	fmt.Println("🎤 Native STT Direct Test")
	fmt.Println("==========================\n")

	// Check for audio file
	audioFile := "test_voice.aiff"
	if len(os.Args) > 1 {
		audioFile = os.Args[1]
	}

	absPath, err := filepath.Abs(audioFile)
	if err != nil {
		log.Fatalf("❌ Failed to get absolute path: %v", err)
	}

	fmt.Printf("Audio file: %s\n", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Println("\n❌ Audio file not found!")
		fmt.Println("\nCreate test audio:")
		fmt.Println("  say -o test_voice.aiff \"Hello, this is a test\"")
		os.Exit(1)
	}

	// Create STT service
	fmt.Println("\n1. Creating Native STT service...")
	stt, err := services.NewNativeSTTService("en-US")
	if err != nil {
		log.Fatalf("❌ Failed to create STT: %v\n\n"+
			"ℹ️  This might be because:\n"+
			"   - First time: System Settings > Privacy & Security > Speech Recognition\n"+
			"   - Enable for Terminal or your app\n", err)
	}
	defer stt.Close()
	fmt.Println("✅ STT service created")

	// Check availability
	fmt.Println("\n2. Checking availability...")
	if !stt.IsAvailable() {
		log.Fatal("❌ STT not available for this locale")
	}
	fmt.Println("✅ STT is available")

	// Transcribe
	fmt.Println("\n3. Transcribing audio...")
	fmt.Println("   (This may take 5-10 seconds...)")
	
	text, err := stt.TranscribeFile(absPath)
	if err != nil {
		log.Fatalf("❌ Transcription failed: %v", err)
	}

	// Show result
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 TRANSCRIPTION RESULT:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n%s\n\n", text)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n✅ Test complete!")
}

