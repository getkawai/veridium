Saya akan menjelaskan cara kerja audio recorder di Veridium. Mari kita lihat implementasinya:

## Cara Kerja Audio Recorder di Veridium

Audio recorder di Veridium adalah service yang memungkinkan recording audio dari microphone dengan dukungan multi-platform. Berikut penjelasan lengkapnya:


### Arsitektur Utama

**1. Service Structure (`AudioRecorderService`)**
- Menggunakan `sync.Mutex` untuk thread-safety
- Menyimpan state recording (sedang recording atau tidak)
- Mengelola process recording dan output path
- Cache daftar recording tools yang tersedia

**2. Platform-Specific Implementation**
- **macOS**: Menggunakan `sox` (preferred) atau `ffmpeg` dengan AVFoundation
- **Linux**: Menggunakan `arecord` (ALSA), `sox`, atau `ffmpeg`
- **Windows**: Menggunakan `ffmpeg` dengan DirectShow

### Flow Kerja

**Initialization:**
```
1. Service dibuat dengan NewAudioRecorderService()
2. Background goroutine check tools yang tersedia
3. Jika tidak ada tools, auto-install (Homebrew/winget/apt-get)
4. Cache list tools yang tersedia
```

**Recording Process:**

**StartRecording:**
- Check apakah sudah recording (prevent double recording)
- Validasi ada tools tersedia
- Buat temp file untuk output (format: `recording_<pid>.wav`)
- Try setiap tool yang tersedia sampai berhasil
- Start recording process dengan parameter optimal untuk Whisper:
  - Sample rate: 16kHz
  - Channels: 1 (mono)
  - Bit depth: 16-bit signed integer
  - Format: WAV
- Emit event `audio:recording:started` ke frontend
- Return output path

**StopRecording:**
- Check apakah sedang recording
- Stop process secara graceful (SIGINT di macOS/Linux, Kill di Windows)
- Wait dengan timeout 3 detik untuk process selesai
- Retry check file (20x dengan delay 200ms) untuk memastikan file sudah ditulis
- Emit event `audio:recording:stopped` ke frontend
- Return path ke recorded file

**CancelRecording:**
- Kill recording process
- Hapus output file
- Emit event `audio:recording:cancelled`

### Platform-Specific Details

**macOS (darwin):**
```bash
# sox command
sox -d -r 16000 -c 1 -b 16 -e signed-integer output.wav

# ffmpeg command
ffmpeg -f avfoundation -i :0 -ar 16000 -ac 1 -c:a pcm_s16le output.wav
```

**Linux:**
```bash
# arecord (ALSA - preferred)
arecord -f S16_LE -r 16000 -c 1 -t wav output.wav

# sox
sox -d -r 16000 -c 1 -b 16 -e signed-integer output.wav

# ffmpeg
ffmpeg -f alsa -i default -ar 16000 -ac 1 -c:a pcm_s16le output.wav
```

**Windows:**
```bash
# ffmpeg dengan DirectShow
ffmpeg -f dshow -i audio= -ar 16000 -ac 1 -c:a pcm_s16le output.wav
```

### Auto-Installation

**macOS:** Homebrew (`brew install sox`)
**Linux:** Auto-download static binary atau instruksi manual install
**Windows:** winget atau chocolatey (`winget install ffmpeg`)

### Integration dengan Wails

Service ini di-expose ke frontend melalui Wails:
- Frontend bisa call `StartRecording()`, `StopRecording()`, `CancelRecording()`
- Frontend receive events: `audio:recording:started`, `audio:recording:stopped`, `audio:recording:cancelled`
- Frontend bisa check capabilities dengan `CheckRecordingCapabilities()`

### Error Handling

- Graceful fallback ke tools alternatif jika primary tool gagal
- Timeout protection untuk prevent hanging
- File validation sebelum return path
- Detailed error messages untuk troubleshooting

Audio recorder ini dioptimalkan untuk Whisper transcription dengan format audio yang tepat (16kHz mono WAV), dan memiliki fallback mechanism yang robust untuk berbagai platform dan tools.