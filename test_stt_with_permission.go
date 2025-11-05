package main

import (
"fmt"
"log"
"os"
"path/filepath"
"strings"
"unsafe"
)

/*
#cgo CFLAGS: -x objective-c -fmodules -fblocks
#cgo LDFLAGS: -framework Speech -framework AVFoundation -framework Foundation
#import <Speech/Speech.h>
#import <Foundation/Foundation.h>

int check_authorization_status() {
    SFSpeechRecognizerAuthorizationStatus status = [SFSpeechRecognizer authorizationStatus];
    return (int)status;
}

const char* get_authorization_status_string() {
    SFSpeechRecognizerAuthorizationStatus status = [SFSpeechRecognizer authorizationStatus];
    switch (status) {
        case SFSpeechRecognizerAuthorizationStatusNotDetermined:
            return "Not Determined";
        case SFSpeechRecognizerAuthorizationStatusDenied:
            return "Denied";
        case SFSpeechRecognizerAuthorizationStatusRestricted:
            return "Restricted";
        case SFSpeechRecognizerAuthorizationStatusAuthorized:
            return "Authorized";
        default:
            return "Unknown";
    }
}
*/
import "C"

func main() {
fmt.Println("🎤 Native STT Permission Check")
fmt.Println("================================\n")

// Check current authorization status
status := C.check_authorization_status()
statusStr := C.GoString(C.get_authorization_status_string())

fmt.Printf("Authorization Status: %s (%d)\n\n", statusStr, status)

if status == 3 { // Authorized
fmt.Println("✅ Speech recognition is authorized!")
fmt.Println("\nYou can now test transcription:")
fmt.Println("  CGO_ENABLED=1 go run test_stt_direct.go test_voice.aiff")
} else if status == 0 { // Not Determined
fmt.Println("⚠️  Permission not yet requested")
fmt.Println("\nℹ️  To grant permission:")
fmt.Println("   1. Run the transcription test once")
fmt.Println("   2. System will show permission dialog")
fmt.Println("   3. Click 'OK' to authorize")
fmt.Println("\n   OR manually enable in:")
fmt.Println("   System Settings > Privacy & Security > Speech Recognition")
} else if status == 2 { // Denied
fmt.Println("❌ Speech recognition is denied")
fmt.Println("\nℹ️  To enable:")
fmt.Println("   System Settings > Privacy & Security > Speech Recognition")
fmt.Println("   Enable for Terminal or your app")
} else if status == 1 { // Restricted
fmt.Println("❌ Speech recognition is restricted (parental controls?)")
}
}
