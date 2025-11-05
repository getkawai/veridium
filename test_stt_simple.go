package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/services"
)

func main() {
	fmt.Println("🎤 Native STT Test")
	fmt.Println("==================")

	// Create STT service
	fmt.Println("\n1. Initializing Native STT service...")
	stt, err := services.NewNativeSTTService("en-US")
	if err != nil {
		log.Fatalf("❌ Failed to create STT: %v", err)
	}
	defer stt.Close()

	// Check if available
	fmt.Println("2. Checking availability...")
	if !stt.IsAvailable() {
		log.Fatal("❌ STT not available on this system")
	}
	fmt.Println("✅ STT is available")

	// Get supported locales
	fmt.Println("\n3. Getting supported locales...")
	locales, err := stt.GetSupportedLocales()
	if err != nil {
		log.Printf("⚠️  Warning: Could not get locales: %v", err)
	} else {
		fmt.Printf("✅ Supported locales: %d languages\n", len(locales))
		fmt.Printf("   First 5: %v\n", locales[:min(5, len(locales))])
	}

	// Check for audio file
	audioFile := "test_voice.wav"
	if len(os.Args) > 1 {
		audioFile = os.Args[1]
	}

	fmt.Printf("\n4. Checking for audio file: %s\n", audioFile)
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		fmt.Println("\n❌ Audio file not found!")
		fmt.Println("\n📝 How to create test audio:")
		fmt.Println("   Option 1 - Use TTS to create test audio:")
		fmt.Println("      say -o test_voice.wav \"Hello, this is a test\"")
		fmt.Println("")
		fmt.Println("   Option 2 - Record from mic (5 seconds):")
		fmt.Println("      brew install sox")
		fmt.Println("      sox -d -r 16000 -c 1 test_voice.wav trim 0 5")
		fmt.Println("")
		fmt.Println("   Option 3 - Use QuickTime Player:")
		fmt.Println("      File > New Audio Recording > Record > Save as test_voice.wav")
		fmt.Println("")
		fmt.Println("Then run: go run test_stt_simple.go test_voice.wav")
		os.Exit(1)
	}
	fmt.Println("✅ Audio file found")

	// Transcribe
	fmt.Println("\n5. Transcribing audio...")
	fmt.Println("   (This may take a few seconds...)")
	text, err := stt.TranscribeFile(audioFile)
	if err != nil {
		log.Fatalf("❌ Failed to transcribe: %v", err)
	}

	// Show result
	fmt.Println("\n" + "="*50)
	fmt.Println("🎉 TRANSCRIPTION RESULT:")
	fmt.Println("=" * 50)
	fmt.Printf("\n   \"%s\"\n\n", text)
	fmt.Println("=" * 50)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
