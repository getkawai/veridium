# Audio Recorder Auto-Setup

## Overview
The Audio Recorder service automatically installs platform-specific recording tools in the background.

## Supported Platforms

### macOS
- **Tool**: `sox` (Sound eXchange)
- **Installation**: `brew install sox`
- **Auto-installer**: Uses Homebrew
- **Requirements**: Homebrew must be installed

### Linux
- **Tools** (tries in order):
  1. `arecord` (ALSA utils) - Recommended
  2. `sox` (Sound eXchange) - Alternative
  3. `ffmpeg` (system PATH) - Alternative
  4. `ffmpeg` (~/.local/bin) - Auto-downloaded if none found
- **Installation**: 
  - **Automatic**: Downloads static ffmpeg binary to `~/.local/bin/ffmpeg` (~50 MB)
  - **Manual** (if auto-download fails):
    - Debian/Ubuntu: `sudo apt-get install -y alsa-utils` (or `sox` or `ffmpeg`)
    - RHEL/CentOS: `sudo yum install -y alsa-utils` (or `sox` or `ffmpeg`)
    - Fedora: `sudo dnf install -y alsa-utils` (or `sox` or `ffmpeg`)
    - Arch Linux: `sudo pacman -S alsa-utils` (or `sox` or `ffmpeg`)
    - openSUSE: `sudo zypper install alsa-utils` (or `sox` or `ffmpeg`)
- **Auto-installer**: Downloads static ffmpeg binary from johnvansickle.com (trusted source)
- **Requirements**: Internet connection for auto-download, or manual installation

### Windows
- **Tool**: `ffmpeg`
- **Installation**: Varies by package manager
  - winget: `winget install ffmpeg`
  - chocolatey: `choco install ffmpeg -y`
- **Auto-installer**: Tries winget first, then chocolatey
- **Requirements**: winget (Windows 10+) or chocolatey

## How It Works

### Background Initialization

```go
// Service creation is instant
service := audio_recorder.NewAudioRecorderService(nil)

// Background goroutine starts automatically
go service.initializeInBackground()

// Steps executed in background:
// 1. Check if recording tool is installed
// 2. If not found, attempt auto-installation
// 3. Verify installation
// 4. Log status to console
```

### Console Output

**When tool is already installed:**
```
✅ Audio recording ready with sox
```

**When tool needs installation (macOS):**
```
🔧 sox not found, attempting auto-installation...
   Installing sox via Homebrew...
✅ sox installed successfully
✅ Audio recording ready with sox
```

**When tool needs installation (Linux) - Auto-download:**
```
🔧 No recording tool found
   Attempting to download ffmpeg static binary...
   Downloading ffmpeg for amd64...
   This may take a few minutes (~50 MB)...
   Extracting ffmpeg...
   ffmpeg installed to: /home/user/.local/bin/ffmpeg
   Note: You may need to add ~/.local/bin to your PATH
✅ ffmpeg static binary downloaded successfully
✅ Audio recording ready with ffmpeg (local)
```

**When auto-download fails:**
```
🔧 No recording tool found
   Attempting to download ffmpeg static binary...
⚠️  Auto-download failed: network error

   Please install one of these tools manually:
   $ sudo apt-get install -y alsa-utils    # Recommended (arecord)
   $ sudo apt-get install -y sox           # Alternative
   $ sudo apt-get install -y ffmpeg        # Alternative
```

**When installation fails:**
```
⚠️  Failed to auto-install recording tool: homebrew is not installed
   Audio recording will not be available until tool is installed
```

## Platform-Specific Files

All platform-specific code uses Go build tags:

```
internal/audio_recorder/
├── audio_recorder.go              # Common logic
├── audio_recorder_darwin.go       # macOS: sox + Homebrew
├── audio_recorder_linux.go        # Linux: arecord + apt/yum/pacman
└── audio_recorder_windows.go      # Windows: ffmpeg + winget/choco
```

### Build Tags
- `//go:build darwin` - macOS only
- `//go:build linux` - Linux only
- `//go:build windows` - Windows only

## Recording Format

All platforms record in the same format (optimized for Whisper):
- **Format**: WAV
- **Sample Rate**: 16kHz
- **Channels**: Mono (1)
- **Bit Depth**: 16-bit signed integer
- **Encoding**: PCM

## API Usage

### Check Recording Capabilities
```go
service := audio_recorder.NewAudioRecorderService(nil)

// Check if recording is supported
capabilities := service.CheckRecordingCapabilities()
// Returns:
// {
//   "supported": true,
//   "tool": "sox",
//   "error": ""
// }
```

### Start Recording
```go
ctx := context.Background()
outputPath, err := service.StartRecording(ctx)
// Returns: "/tmp/recording_12345.wav"
```

### Stop Recording
```go
outputPath, err := service.StopRecording()
// Returns path to recorded WAV file
```

## Error Handling

### Installation Failures

**macOS - No Homebrew:**
```
⚠️  Failed to auto-install recording tool: homebrew is not installed
   Please install Homebrew first: https://brew.sh
```

**Linux - Manual installation required:**
```
🔧 No recording tool found
   Audio recording requires one of: arecord, sox, or ffmpeg

   Please run ONE of these commands in your terminal:
   $ sudo apt-get install -y alsa-utils    # Recommended (arecord)
   $ sudo apt-get install -y sox           # Alternative
   $ sudo apt-get install -y ffmpeg        # Alternative
```

**Note**: 
- Linux auto-downloads static ffmpeg binary to `~/.local/bin/` (no sudo required!)
- The service will automatically use any of the available tools
- Priority: arecord > sox > ffmpeg (PATH) > ffmpeg (local)
- Static binary source: https://johnvansickle.com/ffmpeg/ (trusted, widely used)
- Supports: amd64, arm64, armhf, i686

**Windows - No package manager:**
```
⚠️  Failed to auto-install recording tool: no supported package manager found
   Please install ffmpeg manually from: https://ffmpeg.org/download.html
```

## Benefits

1. **✅ Zero Configuration** - Works out of the box
2. **✅ Non-Blocking** - Service starts immediately
3. **✅ Smart Detection** - Auto-detects package managers
4. **✅ Platform-Agnostic** - Same API across all platforms
5. **✅ Transparent** - Clear progress logs
6. **✅ Graceful Fallback** - App continues if installation fails

## Manual Installation

If auto-installation fails, users can install manually:

### macOS
```bash
brew install sox
```

### Linux (Debian/Ubuntu)
```bash
sudo apt-get install alsa-utils
```

### Linux (RHEL/CentOS)
```bash
sudo yum install -y alsa-utils
```

### Linux (Fedora)
```bash
sudo dnf install -y alsa-utils
```

### Linux (Arch)
```bash
sudo pacman -S alsa-utils
```

### Linux (openSUSE)
```bash
sudo zypper install alsa-utils
```

### Windows (winget)
```powershell
winget install ffmpeg
```

### Windows (chocolatey)
```powershell
choco install ffmpeg
```

## Integration with Whisper STT

The audio recorder is designed to work seamlessly with Whisper:

1. **Recording**: Creates WAV file in correct format
2. **Transcription**: Whisper reads the WAV file
3. **Cleanup**: Temp files are managed automatically

```go
// Record audio
outputPath, _ := audioRecorder.StartRecording(ctx)
// ... user speaks ...
outputPath, _ = audioRecorder.StopRecording()

// Transcribe with Whisper
text, _ := whisperService.Transcribe(ctx, "base", outputPath)
```

## Notes

- Recording tools are installed system-wide
- Installation requires internet connection
- First recording may be slower (tool installation)
- Subsequent recordings are instant
- All installations respect system package managers

