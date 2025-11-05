//go:build darwin
// +build darwin

package services

/*
#cgo CFLAGS: -x objective-c -fmodules -fblocks
#cgo LDFLAGS: -framework Speech -framework AVFoundation -framework Foundation
#import <Speech/Speech.h>
#import <AVFoundation/AVFoundation.h>
#import <Foundation/Foundation.h>

@interface SpeechRecognizer : NSObject <SFSpeechRecognizerDelegate>
@property (nonatomic, strong) SFSpeechRecognizer *recognizer;
@property (nonatomic, strong) SFSpeechAudioBufferRecognitionRequest *request;
@property (nonatomic, strong) SFSpeechRecognitionTask *recognitionTask;
@property (nonatomic, strong) AVAudioEngine *audioEngine;
@end

@implementation SpeechRecognizer

- (instancetype)initWithLocale:(NSString *)localeIdentifier {
    self = [super init];
    if (self) {
        NSLocale *locale = [NSLocale localeWithLocaleIdentifier:localeIdentifier];
        _recognizer = [[SFSpeechRecognizer alloc] initWithLocale:locale];
        _recognizer.delegate = self;
        _audioEngine = [[AVAudioEngine alloc] init];
    }
    return self;
}

- (BOOL)requestAuthorization:(char **)errorOut {
    __block BOOL authorized = NO;

    // Use dispatch_semaphore for reliable synchronization
    dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);

    [SFSpeechRecognizer requestAuthorization:^(SFSpeechRecognizerAuthorizationStatus status) {
        authorized = (status == SFSpeechRecognizerAuthorizationStatusAuthorized);
        dispatch_semaphore_signal(semaphore);
    }];

    // Wait for authorization (max 10 seconds)
    dispatch_time_t timeout = dispatch_time(DISPATCH_TIME_NOW, 10 * NSEC_PER_SEC);
    long result = dispatch_semaphore_wait(semaphore, timeout);

    if (result != 0) {
        // Timeout
        if (errorOut) {
            *errorOut = strdup("Authorization request timeout");
        }
        return NO;
    }

    if (!authorized && errorOut) {
        *errorOut = strdup("Speech recognition not authorized");
    }

    return authorized;
}

- (char *)transcribeFileSync:(NSString *)filePath error:(char **)errorOut {
    NSURL *url = [NSURL fileURLWithPath:filePath];
    if (![[NSFileManager defaultManager] fileExistsAtPath:filePath]) {
        if (errorOut) {
            *errorOut = strdup([[NSString stringWithFormat:@"File not found: %@", filePath] UTF8String]);
        }
        return NULL;
    }

    // Check authorization
    SFSpeechRecognizerAuthorizationStatus authStatus = [SFSpeechRecognizer authorizationStatus];
    if (authStatus != SFSpeechRecognizerAuthorizationStatusAuthorized) {
        if (errorOut) {
            NSString *msg = [NSString stringWithFormat:@"Speech recognition not authorized (status: %ld). Enable in System Settings > Privacy & Security > Speech Recognition", (long)authStatus];
            *errorOut = strdup([msg UTF8String]);
        }
        return NULL;
    }

    SFSpeechURLRecognitionRequest *request = [[SFSpeechURLRecognitionRequest alloc] initWithURL:url];
    request.shouldReportPartialResults = NO; // Only get final result

    __block NSString *finalText = nil;
    __block NSError *finalError = nil;
    __block BOOL callbackCalled = NO;

    // Use dispatch_semaphore for reliable synchronization in CGO context
    dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);

    self.recognitionTask = [self.recognizer recognitionTaskWithRequest:request resultHandler:^(SFSpeechRecognitionResult *result, NSError *error) {
        callbackCalled = YES;

        if (error) {
            finalError = error;
            dispatch_semaphore_signal(semaphore);
            return;
        }

        if (result && result.isFinal) {
            finalText = result.bestTranscription.formattedString;
            dispatch_semaphore_signal(semaphore);
        }
    }];

    // Wait for completion with 60 second timeout
    dispatch_time_t timeout = dispatch_time(DISPATCH_TIME_NOW, 60 * NSEC_PER_SEC);
    long waitResult = dispatch_semaphore_wait(semaphore, timeout);

    if (waitResult != 0) {
        // Timeout
        [self.recognitionTask cancel];
        if (errorOut) {
            if (callbackCalled) {
                *errorOut = strdup("Transcription incomplete (callback called but not final)");
            } else {
                *errorOut = strdup("Transcription timeout (callback never called - check permission)");
            }
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

- (NSArray<NSString *> *)supportedLocales {
    NSMutableArray *locales = [NSMutableArray array];

    // Common supported locales
    NSArray *commonLocales = @[
        @"en-US", @"en-GB", @"en-AU", @"en-IN",
        @"zh-CN", @"zh-TW", @"zh-HK",
        @"ja-JP", @"ko-KR",
        @"es-ES", @"es-MX",
        @"fr-FR", @"fr-CA",
        @"de-DE", @"it-IT", @"pt-BR", @"pt-PT",
        @"ru-RU", @"ar-SA", @"id-ID",
        @"th-TH", @"vi-VN", @"tr-TR",
        @"nl-NL", @"pl-PL", @"sv-SE",
        @"da-DK", @"fi-FI", @"no-NO",
        @"he-IL", @"ro-RO", @"uk-UA"
    ];

    for (NSString *localeId in commonLocales) {
        NSLocale *locale = [NSLocale localeWithLocaleIdentifier:localeId];
        SFSpeechRecognizer *testRecognizer = [[SFSpeechRecognizer alloc] initWithLocale:locale];
        if (testRecognizer.isAvailable) {
            [locales addObject:localeId];
        }
    }

    return locales;
}

- (void)dealloc {
    [self.recognitionTask cancel];
    [self.audioEngine stop];
}

@end

// C wrapper functions

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

char** native_stt_supported_locales(void* recognizer, int* count) {
    @autoreleasepool {
        SpeechRecognizer *sr = (__bridge SpeechRecognizer*)recognizer;
        NSArray<NSString *> *locales = [sr supportedLocales];

        *count = (int)[locales count];
        char** result = (char**)malloc(sizeof(char*) * (*count));

        for (int i = 0; i < *count; i++) {
            result[i] = strdup([locales[i] UTF8String]);
        }

        return result;
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

void native_stt_free_string_array(char** arr, int count) {
    if (arr) {
        for (int i = 0; i < count; i++) {
            if (arr[i]) {
                free(arr[i]);
            }
        }
        free(arr);
    }
}
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

// NativeSTTService provides speech-to-text using macOS Speech Framework
type NativeSTTService struct {
	recognizer unsafe.Pointer
	locale     string
	mu         sync.Mutex
}

// TranscriptionResult contains the transcription result
type TranscriptionResult struct {
	Text    string
	IsFinal bool
	Error   error
}

// NewNativeSTTService creates a new native STT service
// locale: language code (e.g., "en-US", "id-ID", "ja-JP")
func NewNativeSTTService(locale string) (*NativeSTTService, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("native STT only supported on macOS")
	}

	if locale == "" {
		locale = "en-US"
	}

	cLocale := C.CString(locale)
	defer C.free(unsafe.Pointer(cLocale))

	recognizer := C.native_stt_new(cLocale)
	if recognizer == nil {
		return nil, fmt.Errorf("failed to create speech recognizer")
	}

	service := &NativeSTTService{
		recognizer: recognizer,
		locale:     locale,
	}

	// Note: Authorization will be requested automatically on first use
	// Skipping explicit authorization request to avoid CGO/main thread issues

	return service, nil
}

// RequestAuthorization requests user permission for speech recognition
func (s *NativeSTTService) RequestAuthorization() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var cError *C.char
	result := C.native_stt_request_authorization(s.recognizer, &cError)

	if result == 0 {
		if cError != nil {
			errMsg := C.GoString(cError)
			C.native_stt_free_string(cError)
			return fmt.Errorf("%s", errMsg)
		}
		return fmt.Errorf("authorization denied")
	}

	return nil
}

// TranscribeFile transcribes an audio file to text
func (s *NativeSTTService) TranscribeFile(audioPath string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cPath := C.CString(audioPath)
	defer C.free(unsafe.Pointer(cPath))

	var cError *C.char
	cText := C.native_stt_transcribe_file_sync(s.recognizer, cPath, &cError)

	if cError != nil {
		errMsg := C.GoString(cError)
		C.native_stt_free_string(cError)
		return "", fmt.Errorf("%s", errMsg)
	}

	if cText == nil {
		return "", fmt.Errorf("transcription returned no text")
	}

	text := C.GoString(cText)
	C.native_stt_free_string(cText)

	return text, nil
}

// IsAvailable checks if speech recognition is available
func (s *NativeSTTService) IsAvailable() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := C.native_stt_is_available(s.recognizer)
	return result != 0
}

// GetSupportedLocales returns list of supported language codes
func (s *NativeSTTService) GetSupportedLocales() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count C.int
	cLocales := C.native_stt_supported_locales(s.recognizer, &count)

	if cLocales == nil {
		return nil, fmt.Errorf("failed to get supported locales")
	}

	// Convert C array to Go slice
	locales := make([]string, int(count))
	cArray := (*[1 << 30]*C.char)(unsafe.Pointer(cLocales))[:count:count]

	for i := 0; i < int(count); i++ {
		locales[i] = C.GoString(cArray[i])
	}

	C.native_stt_free_string_array(cLocales, count)

	return locales, nil
}

// GetLocale returns the current locale
func (s *NativeSTTService) GetLocale() string {
	return s.locale
}

// Close releases resources
func (s *NativeSTTService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.recognizer != nil {
		C.native_stt_free(s.recognizer)
		s.recognizer = nil
	}

	return nil
}

// GetRecommendedLocales returns recommended language codes with descriptions
func GetRecommendedLocales() map[string]string {
	return map[string]string{
		"en-US": "English (United States)",
		"en-GB": "English (United Kingdom)",
		"en-AU": "English (Australia)",
		"en-IN": "English (India)",
		"zh-CN": "Chinese (Simplified)",
		"zh-TW": "Chinese (Traditional)",
		"ja-JP": "Japanese",
		"ko-KR": "Korean",
		"id-ID": "Indonesian",
		"es-ES": "Spanish (Spain)",
		"es-MX": "Spanish (Mexico)",
		"fr-FR": "French (France)",
		"de-DE": "German",
		"it-IT": "Italian",
		"pt-BR": "Portuguese (Brazil)",
		"ru-RU": "Russian",
		"ar-SA": "Arabic",
		"th-TH": "Thai",
		"vi-VN": "Vietnamese",
		"tr-TR": "Turkish",
		"nl-NL": "Dutch",
		"pl-PL": "Polish",
		"sv-SE": "Swedish",
	}
}
