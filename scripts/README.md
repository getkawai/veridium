# Veridium Scripts

## 🎤 Speech Scripts

### record_and_transcribe.sh

Record audio from microphone and transcribe with Whisper.

**Usage:**
```bash
./scripts/record_and_transcribe.sh
```

**Requirements:**
- sox: `brew install sox`
- Whisper model downloaded (script will prompt if needed)

**How it works:**
1. Press Enter to start recording
2. Speak into your microphone
3. Press Ctrl+C to stop
4. Automatic transcription with Whisper

---

### test_whisper_now.sh

Automated test script for Whisper STT.

**Usage:**
```bash
./scripts/test_whisper_now.sh
```

**What it does:**
1. Creates test audio with TTS
2. Checks Whisper service
3. Downloads model if needed
4. Transcribes test audio
5. Shows results

**Output:**
- Model download progress
- GPU information
- Transcription result
- Performance metrics

---

## 📖 Documentation

See [docs/speech/](../docs/speech/) for detailed documentation:
- [MICROPHONE_GUIDE.md](../docs/speech/MICROPHONE_GUIDE.md) - Microphone usage
- [SPEECH_SERVICES_FINAL.md](../docs/speech/SPEECH_SERVICES_FINAL.md) - Complete guide

