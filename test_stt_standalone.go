// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fmodules -fblocks
#cgo LDFLAGS: -framework Speech -framework AVFoundation -framework Foundation
#import <Speech/Speech.h>
#import <AVFoundation/AVFoundation.h>
#import <Foundation/Foundation.h>

typedef void (*TranscriptionCallback)(const char* text, int isFinal, const char* error);

@interface SpeechRecognizer : NSObject <SFSpeechRecognizerDelegate>
@property (nonatomic, strong) SFSpeechRecognizer *recognizer;
@property (nonatomic, strong) SFSpeechRecognitionTask *recognitionTask;
@property (nonatomic, assign) TranscriptionCallback callback;
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

- (BOOL)requestAuthorization:(char **)errorOut {
    __block BOOL authorized = NO;
    __block BOOL completed = NO;
    
    [SFSpeechRecognizer requestAuthorization:^(SFSpeechRecognizerAuthorizationStatus status) {
        authorized = (status == SFSpeechRecognizerAuthorizationStatusAuthorized);
        completed = YES;
    }];
    
    NSDate *timeout = [NSDate dateWithTimeIntervalSinceNow:5.0];
    while (!completed && [timeout timeIntervalSinceNow] > 0) {
        [[NSRunLoop currentRunLoop] runMode:NSDefaultRunLoopMode beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
    }
    
    if (!authorized && errorOut) {
        *errorOut = strdup("Speech recognition not authorized");
    }
    
    return authorized;
}

- (BOOL)transcribeFile:(NSString *)filePath callback:(TranscriptionCallback)callback error:(char **)errorOut {
    self.callback = callback;
    
    NSURL *url = [NSURL fileURLWithPath:filePath];
    if (![[NSFileManager defaultManager] fileExistsAtPath:filePath]) {
        if (errorOut) {
            *errorOut = strdup([[NSString stringWithFormat:@"File not found: %@", filePath] UTF8String]);
        }
        return NO;
    }
    
    SFSpeechURLRecognitionRequest *request = [[SFSpeechURLRecognitionRequest alloc] initWithURL:url];
    request.shouldReportPartialResults = YES;
    
    __block BOOL completed = NO;
    __block BOOL success = NO;
    
    self.recognitionTask = [self.recognizer recognitionTaskWithRequest:request resultHandler:^(SFSpeechRecognitionResult *result, NSError *error) {
        if (error) {
            if (callback) {
                callback(NULL, 0, [[error localizedDescription] UTF8String]);
            }
            completed = YES;
            return;
        }
        
        if (result) {
            NSString *text = result.bestTranscription.formattedString;
            if (callback) {
                callback([text UTF8String], result.isFinal ? 1 : 0, NULL);
            }
            
            if (result.isFinal) {
                completed = YES;
                success = YES;
            }
        }
    }];
    
    NSDate *timeout = [NSDate dateWithTimeIntervalSinceNow:60.0];
    while (!completed && [timeout timeIntervalSinceNow] > 0) {
        [[NSRunLoop currentRunLoop] runMode:NSDefaultRunLoopMode beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
    }
    
    if (!completed && errorOut) {
        *errorOut = strdup("Transcription timeout");
    }
    
    return success;
}

- (BOOL)isAvailable {
    return self.recognizer.isAvailable;
}

- (void)dealloc {
    [self.recognitionTask cancel];
}

@end

void* native_stt_new(const char* locale) {
    @autoreleasepool {
        NSString *localeStr = locale ? [NSString stringWithUTF8String:locale] : @"en-US";
        SpeechRecognizer *recognizer = [[SpeechRecognizer alloc] initWithLocale:localeStr];
        return (__bridge_retained void*)recognizer;
    }
}

int native_stt_request_authorization(void* recognizer, char** error) {
    @autoreleasepool {
        SpeechRecognizer *sr = (__bridge SpeechRecognizer*)recognizer;
        return [sr requestAuthorization:error] ? 1 : 0;
    }
}

int native_stt_transcribe_file(void* recognizer, const char* filePath, TranscriptionCallback callback, char** error) {
    @autoreleasepool {
        SpeechRecognizer *sr = (__bridge SpeechRecognizer*)recognizer;
        NSString *path = [NSString stringWithUTF8String:filePath];
        return [sr transcribeFile:path callback:callback error:error] ? 1 : 0;
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
import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unsafe"
)

func main() {
	fmt.Println("🎤 Native STT Standalone Test")
	fmt.Println("==============================")

	// Check for audio file
	audioFile := "test_voice.aiff"
	if len(os.Args) > 1 {
		audioFile = os.Args[1]
	}

	fmt.Printf("\n1. Checking for audio file: %s\n", audioFile)
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		fmt.Println("\n❌ Audio file not found!")
		fmt.Println("\n📝 Create test audio with:")
		fmt.Println("   say -o test_voice.aiff \"Hello, this is a test\"")
		fmt.Println("\nThen run: go run test_stt_standalone.go test_voice.aiff")
		os.Exit(1)
	}
	fmt.Println("✅ Audio file found")

	// Create recognizer
	fmt.Println("\n2. Creating speech recognizer...")
	cLocale := C.CString("en-US")
	defer C.free(unsafe.Pointer(cLocale))

	recognizer := C.native_stt_new(cLocale)
	if recognizer == nil {
		log.Fatal("❌ Failed to create recognizer")
	}
	defer C.native_stt_free(recognizer)
	fmt.Println("✅ Recognizer created")

	// Request authorization
	fmt.Println("\n3. Requesting authorization...")
	var cError *C.char
	result := C.native_stt_request_authorization(recognizer, &cError)
	if result == 0 {
		if cError != nil {
			errMsg := C.GoString(cError)
			C.native_stt_free_string(cError)
			log.Fatalf("❌ Authorization failed: %s", errMsg)
		}
		log.Fatal("❌ Authorization denied")
	}
	fmt.Println("✅ Authorization granted")

	// Check availability
	fmt.Println("\n4. Checking availability...")
	if C.native_stt_is_available(recognizer) == 0 {
		log.Fatal("❌ STT not available")
	}
	fmt.Println("✅ STT is available")

	// Transcribe
	fmt.Println("\n5. Transcribing audio...")
	fmt.Println("   (This may take a few seconds...)")

	var finalText string
	var transcriptionError error
	var mu sync.Mutex
	done := make(chan struct{})

	// Callback function
	//export goCallback
	goCallback := func(text *C.char, isFinal C.int, err *C.char) {
		mu.Lock()
		defer mu.Unlock()

		if err != nil {
			transcriptionError = fmt.Errorf("%s", C.GoString(err))
			close(done)
			return
		}

		if text != nil {
			finalText = C.GoString(text)
			fmt.Printf("   📝 Partial: %s\n", finalText)
		}

		if isFinal != 0 {
			close(done)
		}
	}

	cPath := C.CString(audioFile)
	defer C.free(unsafe.Pointer(cPath))

	var cTranscribeError *C.char
	C.native_stt_transcribe_file(recognizer, cPath, C.TranscriptionCallback(goCallback), &cTranscribeError)

	// Wait for completion
	<-done

	if transcriptionError != nil {
		log.Fatalf("❌ Transcription failed: %v", transcriptionError)
	}

	if cTranscribeError != nil {
		errMsg := C.GoString(cTranscribeError)
		C.native_stt_free_string(cTranscribeError)
		log.Fatalf("❌ Transcription error: %s", errMsg)
	}

	// Show result
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🎉 TRANSCRIPTION RESULT:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("\n   \"%s\"\n\n", finalText)
	fmt.Println(strings.Repeat("=", 50))
}

