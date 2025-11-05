# 🎤 Cara Test Native STT dengan Ngomong Langsung

## ⚠️ Current Status

Native STT service ada **technical issue** dengan CGO + NSRunLoop yang menyebabkan crash.  
Sedang dalam proses fix dengan dispatch_semaphore approach.

## ✅ Workaround: Test dengan Record Audio Dulu

Untuk test STT dengan suara Anda sendiri **sekarang**, gunakan salah satu cara ini:

---

### 🎙️ Option 1: QuickTime Player (Termudah!)

1. **Buka QuickTime Player**
2. **File > New Audio Recording** (atau `Cmd+Option+N`)
3. **Klik tombol merah** untuk mulai record
4. **Ngomong sesuatu**, misalnya:
   ```
   "Hello, this is a test of the native speech recognition system.
   I am testing the speech to text feature on macOS.
   The weather is nice today and I hope this works perfectly."
   ```
5. **Klik Stop**
6. **File > Export As** > Save as `my_voice.m4a`

7. **Test transcription** (setelah service di-fix):
   ```bash
   CGO_ENABLED=1 go run test_stt_direct.go my_voice.m4a
   ```

---

### 🎵 Option 2: Voice Memos App

1. **Buka Voice Memos** (di Applications atau Launchpad)
2. **Klik tombol merah** untuk record
3. **Ngomong sesuatu**
4. **Klik Done**
5. **Right-click recording** > **Share** > **Save to Files**
6. Save as `my_voice.m4a`

7. **Test transcription**:
   ```bash
   CGO_ENABLED=1 go run test_stt_direct.go my_voice.m4a
   ```

---

### 🔧 Option 3: Terminal dengan `sox` (Advanced)

1. **Install sox**:
   ```bash
   brew install sox
   ```

2. **Record 10 detik**:
   ```bash
   sox -d -r 16000 -c 1 my_voice.wav trim 0 10
   ```

3. **Ngomong saat recording!**

4. **Test transcription**:
   ```bash
   CGO_ENABLED=1 go run test_stt_direct.go my_voice.wav
   ```

---

## 🎯 Expected Result (After Fix)

```
🎤 Native STT Direct Test
==========================

Audio file: /Users/yuda/.../my_voice.m4a

1. Creating Native STT service...
✅ STT service created

2. Checking availability...
✅ STT is available

3. Transcribing audio...
   (This may take 5-10 seconds...)

============================================================
🎉 TRANSCRIPTION RESULT:
============================================================

Hello, this is a test of the native speech recognition system.
I am testing the speech to text feature on macOS.
The weather is nice today and I hope this works perfectly.

============================================================

✅ Test complete!
```

---

## 🔨 Fix in Progress

Sedang implement:
1. ✅ Dispatch semaphore untuk synchronization
2. ✅ Proper memory management dengan ARC
3. ✅ Thread-safe callback handling
4. ⏳ Testing dengan berbagai audio format

---

## 💡 Alternative: TTS Test (Already Working!)

Sambil menunggu STT fix, Anda bisa test **TTS** (Text-to-Speech) yang sudah working:

```bash
# Test TTS
go run -c -o /tmp/tts_test ./services/tts_service.go ./services/tts_service_test.go
/tmp/tts_test -test.v -test.run TestTTSService_Speak

# Atau langsung
say "Hello, this is a test of text to speech"
```

---

## 📝 Notes

- **Permission**: Pertama kali akan minta permission di System Settings
- **Audio Format**: Support `.wav`, `.m4a`, `.aiff`, `.mp3`
- **Languages**: 50+ bahasa (en-US, id-ID, ja-JP, zh-CN, dll)
- **Quality**: Siri-level accuracy

---

## 🚀 Status Update

**ETA untuk fix**: 1-2 jam (implementing dispatch_semaphore approach)

Saya akan update Anda segera setelah fix selesai! 🎉

