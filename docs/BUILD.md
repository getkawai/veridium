# Build & Development Guide

## Prerequisites

- **Go** 1.21+ with CGO enabled
- **Bun** for frontend dependencies
- **Wails v3** CLI (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- **Task** CLI (optional, but recommended: `brew install go-task/tap/go-task`)
- **sox** for audio recording on macOS (`brew install sox`)

## go-whisper Setup

Before building or running in dev mode, you need to compile `go-whisper` (required for STT):

```bash
cd go-whisper
make
cd ..
```

This will:
1. Clone `whisper.cpp` submodule
2. Build `whisper.cpp` with Metal (GPU) support on macOS
3. Install libraries to `go-whisper/build/install/`

## Development

### Using Task (Recommended)

The Taskfile automatically sets `PKG_CONFIG_PATH` for go-whisper:

```bash
# Run in development mode
task dev

# Build for current platform
task build

# Package for distribution
task package
```

### Using Wails CLI Directly

If you prefer to use `wails3` directly, you need to set `PKG_CONFIG_PATH`:

```bash
# Development mode
PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig wails3 dev

# Build
PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig wails3 build
```

## Production Build

### macOS

```bash
# Build for current architecture
task darwin:build

# Build universal binary (arm64 + amd64)
task darwin:build:universal

# Package as .app bundle
task darwin:package

# Package universal .app bundle
task darwin:package:universal
```

### Linux

```bash
task linux:build
task linux:package
```

### Windows

Cross-compilation from macOS to Windows requires:
- MinGW-w64 toolchain
- Wine (for testing)

See `build/windows/Taskfile.yml` for details.

## Environment Variables

The Taskfile uses these variables:

- `PKG_CONFIG_PATH` - Set automatically to `go-whisper/build/install/lib/pkgconfig`
- `VITE_PORT` - Frontend dev server port (default: 9245)
- `PRODUCTION` - Build mode: `true` for production, `false` for dev (default)
- `WAILS_VITE_PORT` - Override vite port from environment

## Troubleshooting

### "Package 'libwhisper' not found"

This means `PKG_CONFIG_PATH` is not set correctly. Solutions:

1. Use `task dev` instead of `wails3 dev` directly
2. Rebuild go-whisper: `cd go-whisper && make && cd ..`
3. Manually set: `export PKG_CONFIG_PATH=$(pwd)/go-whisper/build/install/lib/pkgconfig`

### "sox: command not found"

Install sox for audio recording:

```bash
# macOS
brew install sox

# Linux
sudo apt-get install sox
```

### Build fails with CGO errors

Ensure CGO is enabled:

```bash
export CGO_ENABLED=1
```

For macOS builds, verify Xcode Command Line Tools are installed:

```bash
xcode-select --install
```

## Clean Build

To clean all build artifacts:

```bash
rm -rf bin/
rm -rf frontend/dist/
rm -rf frontend/bindings/
cd go-whisper && make clean && cd ..
```

## Generate Bindings

To regenerate Wails TypeScript bindings:

```bash
wails3 generate bindings -clean=true -ts
```

This is automatically done when running `task dev` or `task build`.

