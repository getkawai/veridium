# Native STT Issue

## Problem

Crash (SIGSEGV) saat memanggil `native_stt_transcribe_file_sync` dari Go code.

## Root Cause

1. **NSRunLoop + CGO**: NSRunLoop tidak reliable di CGO context
2. **Memory Management**: Objective-C objects tidak ter-retain dengan benar
3. **Thread Safety**: Speech recognition callbacks berjalan di background thread

## Solution Options

### Option 1: Use Dispatch Semaphore (Recommended)

Ganti NSRunLoop dengan dispatch_semaphore untuk synchronization yang lebih robust di CGO.

### Option 2: Use Real-time Recognition

Implement real-time mic input dengan AVAudioEngine (lebih complex).

### Option 3: Use Command Line Tool

Buat wrapper CLI tool yang call Speech Framework, lalu call dari Go via exec.Command.

## Temporary Workaround

Untuk test STT dengan ngomong langsung, gunakan **QuickTime Player** atau **macOS Voice Memos**:

1. Record audio dengan QuickTime atau Voice Memos
2. Export as `.m4a` or `.aiff`
3. Transcribe dengan service (setelah fix)

## Next Steps

1. Implement dispatch_semaphore version
2. Add proper memory management with ARC
3. Test with various audio formats

