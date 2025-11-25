# Veridium Documentation

## 📚 Documentation Index

### Core Architecture

- **[EINO_ARCHITECTURE.md](EINO_ARCHITECTURE.md)** - Complete Eino integration architecture (98% complete)
- **[PROPER_MIGRATION_STRATEGY.md](PROPER_MIGRATION_STRATEGY.md)** - Frontend to backend migration strategy (5 phases)
- **[HISTORY_SUMMARY.md](HISTORY_SUMMARY.md)** - Auto-summarization for long conversations (production ready)

### Build & Development

- **[BUILD.md](BUILD.md)** - Build instructions, development guide, and troubleshooting
- **[AUDIO_RECORDER_AUTO_SETUP.md](AUDIO_RECORDER_AUTO_SETUP.md)** - Audio recorder auto-installation guide

### Performance & Optimization

- **[LLM_OPTIMIZATION_GUIDE.md](LLM_OPTIMIZATION_GUIDE.md)** - LLM performance optimization guide
- **[HARDWARE_VALIDATION_COMPLETE.md](HARDWARE_VALIDATION_COMPLETE.md)** - Hardware validation for reasoning modes

### Speech Services

Comprehensive documentation for Text-to-Speech and Speech-to-Text features:

- **[STT_SUCCESS.md](speech/STT_SUCCESS.md)** - ✅ **STT Implementation Success Report** (READ THIS FIRST!)
- **[SPEECH_SERVICES_FINAL.md](speech/SPEECH_SERVICES_FINAL.md)** - Complete guide to TTS and STT services
- **[WHISPER_SETUP.md](speech/WHISPER_SETUP.md)** - Whisper STT setup instructions
- **[WHISPER_TEST_RESULTS.md](speech/WHISPER_TEST_RESULTS.md)** - Performance benchmarks and test results
- **[MICROPHONE_GUIDE.md](speech/MICROPHONE_GUIDE.md)** - How to use microphone for voice input
- **[NATIVE_SPEECH_SETUP.md](speech/NATIVE_SPEECH_SETUP.md)** - Native speech setup (archived)
- **[QUICK_STT_TEST.md](speech/QUICK_STT_TEST.md)** - Quick testing guide

### Quick Links

- [Main README](../README.md)
- [Scripts](../scripts/)

## 🚀 Getting Started

### For Development
1. Read [BUILD.md](BUILD.md) for build setup
2. Read [EINO_ARCHITECTURE.md](EINO_ARCHITECTURE.md) for architecture overview
3. Check [PROPER_MIGRATION_STRATEGY.md](PROPER_MIGRATION_STRATEGY.md) for current migration status

### For Speech Features
1. Read [SPEECH_SERVICES_FINAL.md](speech/SPEECH_SERVICES_FINAL.md) for overview
2. Follow [WHISPER_SETUP.md](speech/WHISPER_SETUP.md) for setup
3. Try [microphone integration](speech/MICROPHONE_GUIDE.md)

## 📁 Project Structure

```
docs/
├── EINO_ARCHITECTURE.md           # Core architecture
├── PROPER_MIGRATION_STRATEGY.md   # Migration plan
├── HISTORY_SUMMARY.md             # Auto-summarization
├── BUILD.md                       # Build guide
├── LLM_OPTIMIZATION_GUIDE.md      # Performance guide
├── HARDWARE_VALIDATION_COMPLETE.md # Hardware validation
├── AUDIO_RECORDER_AUTO_SETUP.md   # Audio setup
└── speech/                        # Speech services docs
    ├── SPEECH_SERVICES_FINAL.md
    ├── WHISPER_SETUP.md
    ├── WHISPER_TEST_RESULTS.md
    ├── MICROPHONE_GUIDE.md
    ├── NATIVE_SPEECH_SETUP.md
    ├── STT_SUCCESS.md
    └── QUICK_STT_TEST.md

scripts/                           # Utility scripts
├── record_and_transcribe.sh
└── test_whisper_now.sh
```

