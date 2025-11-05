// +build darwin

package services

/*
#cgo CFLAGS: -x objective-c -fmodules -fblocks
#cgo LDFLAGS: -framework Speech -framework AVFoundation -framework Foundation
#import <Speech/Speech.h>
#import <AVFoundation/AVFoundation.h>
#import <Foundation/Foundation.h>

// Callback type for transcription results
typedef void (*TranscriptionCallback)(const char* text, int isFinal, const char* error);

@interface SpeechRecognizer : NSObject <SFSpeechRecognizerDelegate>
@property (nonatomic, strong) SFSpeechRecognizer *recognizer;
@property (nonatomic, strong) SFSpeechAudioBufferRecognitionRequest *request;
@property (nonatomic, strong) SFSpeechRecognitionTask *recognitionTask;
@property (nonatomic, strong) AVAudioEngine *audioEngine;
@property (nonatomic, assign) TranscriptionCallback callback;
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
    __block BOOL completed = NO;
    
    [SFSpeechRecognizer requestAuthorization:^(SFSpeechRecognizerAuthorizationStatus status) {
        authorized = (status == SFSpeechRecognizerAuthorizationStatusAuthorized);
        completed = YES;
    }];
    
    // Wait for authorization (max 5 seconds)
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
    
    // Wait for completion (max 60 seconds)
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
	
	// Request authorization
	if err := service.RequestAuthorization(); err != nil {
		C.native_stt_free(recognizer)
		return nil, fmt.Errorf("speech recognition authorization failed: %w", err)
	}
	
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
	
	var finalText string
	var transcriptionError error
	done := make(chan struct{})
	
	// Callback function
	callback := func(text *C.char, isFinal C.int, err *C.char) {
		if err != nil {
			transcriptionError = fmt.Errorf("%s", C.GoString(err))
			close(done)
			return
		}
		
		if text != nil {
			finalText = C.GoString(text)
		}
		
		if isFinal != 0 {
			close(done)
		}
	}
	
	cPath := C.CString(audioPath)
	defer C.free(unsafe.Pointer(cPath))
	
	var cError *C.char
	
	// Start transcription in goroutine
	go func() {
		C.native_stt_transcribe_file(s.recognizer, cPath, C.TranscriptionCallback(callback), &cError)
	}()
	
	// Wait for completion
	<-done
	
	if transcriptionError != nil {
		return "", transcriptionError
	}
	
	if cError != nil {
		errMsg := C.GoString(cError)
		C.native_stt_free_string(cError)
		return "", fmt.Errorf("%s", errMsg)
	}
	
	return finalText, nil
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

