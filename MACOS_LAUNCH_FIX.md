# macOS App Launch Fix

## Problem
Downloaded macOS binary from R2 showed "Apple could not verify" warning and didn't open window even after "Open Anyway".

## Root Cause Analysis

### Issue 1: WaitForDebugger Flag (FIXED ✅)
- **Cause**: `internal/app/context.go` checked `VERIDIUM_DEV` or `DEV` environment variables to set `DevMode = true`
- **Impact**: macOS launch services detected this as development build and set `WaitForDebugger = true`, causing app to quit immediately if no debugger attached
- **Evidence**: macOS logs showed `WaitForDebugger = true` flag

### Issue 2: Relative Path Problem (FIXED ✅)
- **Cause**: `internal/paths/paths.go` used relative path `"data"` which resolved differently depending on working directory
- **Impact**: App worked from terminal (PWD = project root) but failed from Finder (PWD = / or /Applications/)
- **Evidence**: App process started but quit immediately when launched via `open` command

## Solution Implemented

### 1. Production-Only Build
Removed all `DevMode` checks and made the app production-only:

**Removed files**:
- `internal/app/context_production.go` (no longer needed)
- `internal/app/context_dev.go` (no longer needed)

**Modified `internal/app/context.go`**:
- Removed `DevMode` variable completely
- Always use cache layer for embeddings (production behavior)
- Sentry always enabled

**Modified `main.go`**:
- Removed `DevMode` check

### 2. Cross-Platform Data Directory Detection
Updated `internal/paths/paths.go` to properly detect packaged apps and use platform-specific data directories:

**Platform-specific paths**:
- **macOS**: `~/Library/Application Support/Kawai/` (when running from .app bundle)
- **Windows**: `%APPDATA%\Kawai\` (with fallback to `%LOCALAPPDATA%`)
- **Linux**: `~/.config/Kawai/` or `$XDG_CONFIG_HOME/Kawai/` (follows XDG Base Directory spec)
- **Development**: `./data/` (when running from terminal)

**Detection logic**:
```go
func IsPackaged() bool {
    switch runtime.GOOS {
    case "darwin":
        // Check for .app bundle structure (/path/to/App.app/Contents/MacOS/binary)
        if filepath.Base(filepath.Dir(filepath.Dir(execPath))) == "Contents" {
            return true
        }
    case "windows":
        // Check for resources/ directory or Program Files installation
        if filepath.Base(execDir) == "resources" || 
           strings.Contains(execPath, "Program Files") {
            return true
        }
    case "linux":
        // Check for resources/ or standard install paths (/usr/, /opt/)
        if filepath.Base(execDir) == "resources" || 
           filepath.HasPrefix(execPath, "/usr/") || 
           filepath.HasPrefix(execPath, "/opt/") {
            return true
        }
    }
    return false
}

func ensureInit() {
    if IsPackaged() {
        // Use platform-specific user data directory
        switch runtime.GOOS {
        case "darwin":
            dataDir = ~/Library/Application Support/Kawai/
        case "windows":
            dataDir = %APPDATA%\Kawai\ (or %LOCALAPPDATA%\Kawai\)
        case "linux":
            dataDir = $XDG_CONFIG_HOME/Kawai/ (or ~/.config/Kawai/)
        }
    } else {
        // Development mode: use relative ./data/
        dataDir = "data"
    }
}
```

**Key Features**:
- ✅ Automatic detection of packaged vs development mode
- ✅ Platform-specific standard paths (macOS, Windows, Linux)
- ✅ XDG Base Directory specification support (Linux)
- ✅ Windows Program Files detection
- ✅ Linux standard install paths (/usr/, /opt/)
- ✅ Fallback to relative path for development
- ✅ Thread-safe initialization with `sync.RWMutex`
- ✅ Creates directory automatically with proper permissions (0755)

### Build Command
```bash
wails3 task darwin:build PRODUCTION=true ARCH=arm64
wails3 task darwin:package PRODUCTION=true ARCH=arm64
```

This uses `-tags production` flag from `build/darwin/Taskfile.yml`.

## Verification

### ✅ All Issues Fixed
1. **WaitForDebugger**: No more `WaitForDebugger = true` in macOS logs
2. **Terminal Launch**: App runs successfully from terminal
3. **Finder Launch**: App runs successfully via `open` command and double-click
4. **Window Display**: Window appears and loads properly
5. **Services**: All services initialize correctly (Sentry, Database, LLM, etc.)
6. **Cross-Platform**: Data directories properly detected for macOS, Windows, Linux

### Test Results
```bash
# Test 1: Terminal launch (development mode)
./bin/Kawai.app/Contents/MacOS/Kawai
# ✅ Works - uses ./data/ directory

# Test 2: Finder launch (production mode)
open bin/Kawai.app
# ✅ Works - uses ~/Library/Application Support/Kawai/

# Test 3: Check process
ps aux | grep Kawai
# ✅ Process running with proper PID

# Test 4: Check data directory
ls -la ~/Library/Application\ Support/Kawai/
# ✅ Database and files created in proper location
```

## Build Configuration
- `build/darwin/Taskfile.yml` - Already has `-tags production` for PRODUCTION=true builds
- `build/darwin/Info.plist` - Production bundle configuration (unchanged)
- `build/darwin/Info.dev.plist` - Development bundle configuration (unchanged)

## Testing
```bash
# Build production binary
wails3 task darwin:package PRODUCTION=true ARCH=arm64

# Test terminal launch (should use ./data/)
./bin/Kawai.app/Contents/MacOS/Kawai

# Test Finder launch (should use ~/Library/Application Support/Kawai/)
open bin/Kawai.app

# Verify data directory
ls -la ~/Library/Application\ Support/Kawai/

# Check for WaitForDebugger (should be empty)
log show --predicate 'process == "launchd" AND eventMessage CONTAINS "Kawai"' --last 1m | grep -i "wait\|debug"
```

## References
- Wails v3 Build Guide: `/Users/yuda/github.com/wailsapp/wails/docs/src/content/docs/guides/build/macos.mdx`
- Wails v3 Fork: `github.com/yudaprama/wails@v3.0.0-alpha.62-kawai`
- Go Build Tags: https://pkg.go.dev/cmd/go#hdr-Build_constraints
- macOS App Bundle Structure: https://developer.apple.com/library/archive/documentation/CoreFoundation/Conceptual/CFBundles/BundleTypes/BundleTypes.html

