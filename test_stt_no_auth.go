package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c -fmodules -fblocks
#cgo LDFLAGS: -framework Speech -framework AVFoundation -framework Foundation
#import <Speech/Speech.h>
#import <AVFoundation/AVFoundation.h>
#import <Foundation/Foundation.h>

@interface SpeechRecognizer : NSObject <SFSpeechRecognizerDelegate>
@property (nonatomic, strong) SFSpeechRecognizer *recognizer;
@end

@implementation SpeechRecognizer

- (instancetype)initWithLocale:(NSString *)localeIdentifier {
    self = [super init];
    if (self) {
        NSLocale *locale = [NSLocale localeWithLocaleIdentifier:localeIdentifier];
        _recognizer = [[SFSpeechRecognizer alloc] initWithLocale:locale];
        _recognizer.delegate = self;
    }
    return self;
}

- (char *)transcribeFileSync:(NSString *)filePath error:(char **)errorOut {
    NSURL *url = [NSURL fileURLWithPath:filePath];
    if (![[NSFileManager defaultManager] fileExistsAtPath:filePath]) {
        if (errorOut) {
            *errorOut = strdup([[NSString stringWithFormat:@"File not found: %@", filePath] UTF8String]);
        }
        return NULL;
    }
    
    SFSpeechURLRecognitionRequest *request = [[SFSpeechURLRecognitionRequest alloc] initWithURL:url];
    request.shouldReportPartialResults = NO;
    
    __block BOOL completed = NO;
    __block NSString *finalText = nil;
    __block NSError *finalError = nil;
    
    [self.recognizer recognitionTaskWithRequest:request resultHandler:^(SFSpeechRecognitionResult *result, NSError *error) {
        if (error) {
            finalError = error;
            completed = YES;
            return;
        }
        
        if (result && result.isFinal) {
            finalText = result.bestTranscription.formattedString;
            completed = YES;
        }
    }];
    
    NSDate *timeout = [NSDate dateWithTimeIntervalSinceNow:60.0];
    while (!completed && [timeout timeIntervalSinceNow] > 0) {
        [[NSRunLoop currentRunLoop] runMode:NSDefaultRunLoopMode beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
    }
    
    if (!completed) {
        if (errorOut) {
            *errorOut = strdup("Transcription timeout");
        }
        return NULL;
    }
    
    if (finalError) {
        if (errorOut) {
            *errorOut = strdup([[finalError localizedDescription] UTF8String]);
        }
        return NULL;
    }
    
    if (finalText) {
        return strdup([finalText UTF8String]);
    }
    
    return NULL;
}

- (BOOL)isAvailable {
    return self.recognizer.isAvailable;
}

@end

void* native_stt_new(const char* locale) {
    @autoreleasepool {
        NSString *localeStr = locale ? [NSString stringWithUTF8String:locale] : @"en-US";
        SpeechRecognizer *recognizer = [[SpeechRecognizer alloc] initWithLocale:localeStr];
        return (__bridge_retained void*)recognizer;
    }
}

char* native_stt_transcribe_file_sync(void* recognizer, const char* filePath, char** error) {
    @autoreleasepool {
        SpeechRecognizer *sr = (__bridge SpeechRecognizer*)recognizer;
        NSString *path = [NSString stringWithUTF8String:filePath];
        return [sr transcribeFileSync:path error:error];
    }
}

int native_stt_is_available(void* recognizer) {
    @autoreleasepool {
        SpeechRecognizer *sr = (__bridge SpeechRecognizer*)recognizer;
        return [sr isAvailable] ? 1 : 0;
    }
}

void native_stt_free(void* recognizer) {
    if (recognizer) {
        CFRelease(recognizer);
    }
}

void native_stt_free_string(char* str) {
    if (str) {
        free(str);
    }
}
*/
import "C"

func main() {
	if runtime.GOOS != "darwin" {
		log.Fatal("This test only works on macOS")
	}

	fmt.Println("🎤 Native STT Test (No Auth)")
	fmt.Println("==============================\n")

	audioFile := "test_voice.aiff"
	if len(os.Args) > 1 {
		audioFile = os.Args[1]
	}

	absPath, _ := filepath.Abs(audioFile)
	fmt.Printf("Audio file: %s\n", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Println("\n❌ Audio file not found!")
		os.Exit(1)
	}

	// Create recognizer
	fmt.Println("\n1. Creating speech recognizer...")
	cLocale := C.CString("en-US")
	defer C.free(unsafe.Pointer(cLocale))

	recognizer := C.native_stt_new(cLocale)
	if recognizer == nil {
		log.Fatal("❌ Failed to create recognizer")
	}
	defer C.native_stt_free(recognizer)
	fmt.Println("✅ Recognizer created")

	// Check availability
	fmt.Println("\n2. Checking availability...")
	if C.native_stt_is_available(recognizer) == 0 {
		log.Fatal("❌ STT not available\n\n" +
			"ℹ️  Grant permission in:\n" +
			"   System Settings > Privacy & Security > Speech Recognition\n" +
			"   Enable for Terminal")
	}
	fmt.Println("✅ STT is available")

	// Transcribe
	fmt.Println("\n3. Transcribing audio...")
	fmt.Println("   (This may take 5-10 seconds...)")

	cPath := C.CString(absPath)
	defer C.free(unsafe.Pointer(cPath))

	var cError *C.char
	cText := C.native_stt_transcribe_file_sync(recognizer, cPath, &cError)

	if cError != nil {
		errMsg := C.GoString(cError)
		C.native_stt_free_string(cError)
		log.Fatalf("❌ Transcription failed: %s", errMsg)
	}

	if cText == nil {
		log.Fatal("❌ Transcription returned no text")
	}

	text := C.GoString(cText)
	C.native_stt_free_string(cText)

	// Show result
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 TRANSCRIPTION RESULT:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n%s\n\n", text)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n✅ Test complete!")
}

