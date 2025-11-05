# 🎤 Microphone Integration Guide

## ✅ Quick Start: Record & Transcribe

### Option 1: Using Provided Script (Easiest)

```bash
# Install sox (if not installed)
brew install sox

# Run the script
./record_and_transcribe.sh
```

**How it works:**
1. Press Enter to start recording
2. Speak into your microphone
3. Press Ctrl+C to stop
4. Automatic transcription with Whisper!

### Option 2: Manual Record + Transcribe

```bash
# Record 10 seconds
rec -r 16000 -c 1 -b 16 my_recording.wav trim 0 10

# Then transcribe (edit /tmp/transcribe_recording.go to use my_recording.wav)
export PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig
export CGO_ENABLED=1
go run /tmp/transcribe_recording.go
```

### Option 3: QuickTime Recording

```bash
# 1. Open QuickTime Player
# 2. File > New Audio Recording
# 3. Record your voice
# 4. Export as audio.aiff
# 5. Convert & transcribe:

afconvert audio.aiff -d LEI16 -f WAVE audio.wav

# Edit script to use audio.wav, then:
go run /tmp/transcribe_recording.go
```

---

## 🚀 Real-time Streaming (Advanced)

For real-time transcription, we need to integrate AVAudioEngine or PortAudio.

### Architecture

```
Microphone → Audio Capture → Buffer → Whisper → Text Output
              (AVAudioEngine)   (Chunks)  (Transcribe)
```

### Implementation Steps

#### 1. Add Audio Capture Service

```go
// services/audio_capture_service.go
package services

import (
    "context"
    "time"
)

type AudioCaptureService struct {
    sampleRate int
    channels   int
    bufferSize int
}

func NewAudioCaptureService() *AudioCaptureService {
    return &AudioCaptureService{
        sampleRate: 16000,
        channels:   1,
        bufferSize: 16000, // 1 second buffer
    }
}

// StartCapture starts capturing audio from default microphone
func (s *AudioCaptureService) StartCapture(ctx context.Context) (<-chan []float32, error) {
    // Implementation using AVAudioEngine (macOS) or PortAudio (cross-platform)
    // Returns channel of audio chunks
}
```

#### 2. Real-time Transcription Service

```go
// services/realtime_stt_service.go
package services

type RealtimeSTTService struct {
    whisper *WhisperService
    audio   *AudioCaptureService
}

func (s *RealtimeSTTService) StartTranscription(ctx context.Context) (<-chan string, error) {
    audioChunks, err := s.audio.StartCapture(ctx)
    if err != nil {
        return nil, err
    }
    
    textChan := make(chan string)
    
    go func() {
        for chunk := range audioChunks {
            // Save chunk to temp file
            tempFile := saveTempWav(chunk)
            
            // Transcribe
            text, _ := s.whisper.Transcribe(ctx, "ggml-tiny", tempFile)
            
            if text != "" {
                textChan <- text
            }
        }
        close(textChan)
    }()
    
    return textChan, nil
}
```

### Libraries for Audio Capture

**Option A: AVAudioEngine (macOS only, via CGO)**
```go
/*
#cgo LDFLAGS: -framework AVFoundation
#import <AVFoundation/AVFoundation.h>
*/
import "C"
```

**Option B: PortAudio (cross-platform)**
```go
import "github.com/gordonklaus/portaudio"
```

**Option C: Go-native (using exec)**
```go
cmd := exec.Command("rec", "-r", "16000", "-c", "1", "-")
stdout, _ := cmd.StdoutPipe()
cmd.Start()
// Read from stdout
```

---

## 📊 Comparison

| Method | Pros | Cons | Use Case |
|--------|------|------|----------|
| **Record → Transcribe** | ✅ Simple<br>✅ Reliable<br>✅ Works now | ❌ Not real-time | Best for testing |
| **Streaming (exec)** | ✅ Simple<br>✅ Cross-platform | ⚠️ Medium latency | Good enough |
| **Streaming (PortAudio)** | ✅ Low latency<br>✅ Cross-platform | ❌ Requires C library | Production |
| **Streaming (AVAudioEngine)** | ✅ Native<br>✅ Low latency | ❌ macOS only | macOS production |

---

## 💡 Recommended Approach

### For Now: Use Record → Transcribe

```bash
./record_and_transcribe.sh
```

**Why?**
- ✅ Works immediately
- ✅ No additional dependencies
- ✅ Good for testing
- ✅ Can transcribe longer recordings

### For Production: Implement Streaming

Later, add real-time streaming with one of:
1. **PortAudio** (cross-platform)
2. **AVAudioEngine** (macOS only, best quality)
3. **Exec-based streaming** (quick but hacky)

---

## 🎯 Usage Example

### Current (Record → Transcribe)

```bash
# Terminal 1: Start recording
./record_and_transcribe.sh

# Speak: "Hello, this is a test of microphone input"
# Press Ctrl+C

# Output:
# 🎉 TRANSCRIPTION:
# ============================================================
# Hello, this is a test of microphone input
# ============================================================
```

### Future (Real-time Streaming)

```go
// In your frontend
import { services } from '@wailsio/runtime';

// Start real-time transcription
const stream = await services.RealtimeSTTService.StartTranscription();

stream.on('text', (text) => {
    console.log('Transcribed:', text);
    // Update UI with transcribed text
});

// User speaks → Text appears in real-time
```

---

## 🔧 Quick Setup

1. **Install sox** (for recording):
   ```bash
   brew install sox
   ```

2. **Test microphone**:
   ```bash
   rec -r 16000 -c 1 test.wav trim 0 5
   ```

3. **Run script**:
   ```bash
   ./record_and_transcribe.sh
   ```

4. **Speak and stop** (Ctrl+C)

5. **See transcription!** 🎉

---

## ✅ Next Steps

1. ✅ **Current**: Use record_and_transcribe.sh
2. ⏳ **Phase 2**: Add streaming with exec (simple)
3. ⏳ **Phase 3**: Add PortAudio for production
4. ⏳ **Phase 4**: Add VAD (Voice Activity Detection) for auto-stop

**Mau saya implement Phase 2 (streaming) sekarang?** 
Atau test dulu dengan record_and_transcribe.sh?

