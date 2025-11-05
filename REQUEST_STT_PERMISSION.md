# 🔐 Request Speech Recognition Permission

## Problem

Native STT needs permission, but requesting permission programmatically crashes in CGO context.

**Error**: `Speech recognition not authorized (status: 0)`  
**Status 0** = "Not Determined" (permission never requested)

## ✅ Solution: Manual Permission Grant

### Option 1: Via System Settings (Easiest)

1. **Open System Settings**
2. Go to **Privacy & Security**
3. Click **Speech Recognition**
4. **Enable** for **Terminal** (or your app)
5. Run test again

### Option 2: Via Terminal Command

```bash
# This will trigger permission dialog
say "Testing speech recognition"

# Then try STT test
CGO_ENABLED=1 go run test_stt_direct.go test_voice.aiff
```

### Option 3: Build as App Bundle

Create proper macOS app with Info.plist:

```xml
<key>NSSpeechRecognitionUsageDescription</key>
<string>This app needs speech recognition for transcription</string>
```

## 🎯 After Permission Granted

Once permission is granted (status will be 3), Native STT should work perfectly!

```bash
CGO_ENABLED=1 go run test_stt_direct.go test_voice.aiff
```

Expected output:
```
============================================================
🎉 TRANSCRIPTION RESULT:
============================================================

Hello, this is a test of native speech recognition...

============================================================
```

## 🔄 Alternative: Use Whisper

If permission is too complicated, we can restore Whisper which doesn't need permission:

```bash
git merge feature/whisper-hybrid-stt
```

Whisper works immediately without any permission dialogs!

