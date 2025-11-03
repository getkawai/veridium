# Browser TTS Implementation

## Overview
Implementasi Browser TTS menggunakan Web Speech API yang tersedia di browser modern. Ini memberikan alternatif gratis dan efisien untuk layanan TTS eksternal seperti OpenAI, Edge Speech, atau Microsoft Speech.

## Fitur
- ✅ Tidak memerlukan API eksternal atau biaya tambahan
- ✅ Bekerja offline (setelah halaman dimuat)
- ✅ Mendukung berbagai bahasa dan voice tergantung browser
- ✅ Integrasi sempurna dengan UI yang ada
- ✅ Preview voice sebelum memilih

## Cara Menggunakan

### 1. Buka Settings TTS
Navigasi ke: **Settings > Agent > TTS**

### 2. Pilih Browser TTS
Pilih "Browser TTS" dari dropdown **TTS Service**

### 3. Pilih Voice
Setelah memilih Browser TTS, akan muncul dropdown untuk memilih voice. Voice yang tersedia tergantung pada:
- Browser yang digunakan (Chrome, Firefox, Safari, dll.)
- Sistem operasi (Windows, macOS, Linux, iOS, Android)
- Bahasa yang terinstall di sistem

### 4. Test Voice
Klik tombol play di samping dropdown voice untuk mendengar preview

## Browser Support

### Fully Supported
- ✅ Chrome/Chromium (Desktop & Mobile)
- ✅ Edge (Desktop & Mobile)
- ✅ Safari (Desktop & Mobile)
- ✅ Firefox (Desktop & Mobile)
- ✅ Opera

### Limitations
- Voice quality dan pilihan voice bervariasi antar browser
- Beberapa browser mungkin memiliki voice yang terbatas
- Voice synthesis mungkin berbeda di setiap platform

## Technical Details

### Files Modified/Created
1. **`src/hooks/useTTS.ts`** - Added `useBrowserTTS` hook
2. **`src/utils/browserTTS.ts`** - Utility functions for browser TTS
3. **`src/hooks/useBrowserVoices.ts`** - Hook to fetch available voices
4. **`src/types/agent/tts.ts`** - Added 'browser' to TTSServer type
5. **`src/features/AgentSetting/AgentTTS/options.tsx`** - Added Browser TTS option
6. **`src/features/AgentSetting/AgentTTS/index.tsx`** - Added browser voice selector
7. **`src/store/agent/slices/chat/selectors/agent.ts`** - Added browser voice selector

### Web Speech API
Browser TTS menggunakan `SpeechSynthesisUtterance` dan `speechSynthesis` API yang merupakan bagian dari [Web Speech API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Speech_API).

## Troubleshooting

### Voice tidak muncul
- Refresh halaman, beberapa browser memerlukan waktu untuk memuat voice list
- Pastikan browser support Web Speech API
- Cek console untuk error messages

### Voice tidak jelas atau robotik
- Coba voice lain dari dropdown
- Beberapa system voices lebih baik daripada yang lain
- Update sistem operasi untuk mendapatkan voice yang lebih baik

### TTS tidak berfungsi sama sekali
- Pastikan browser support Web Speech API
- Cek permission browser untuk audio/speech
- Coba browser lain jika masalah berlanjut

## Advantages vs Other TTS Services

### Browser TTS
- ✅ Gratis
- ✅ Offline capable
- ✅ Low latency
- ✅ No API key required
- ❌ Voice quality bervariasi
- ❌ Limited voice options

### OpenAI TTS
- ✅ High quality voices
- ✅ Consistent across platforms
- ❌ Requires API key & costs money
- ❌ Requires internet
- ❌ Higher latency

### Edge/Microsoft Speech
- ✅ Good quality
- ✅ Many voices
- ❌ Requires internet
- ❌ May require authentication
- ❌ Platform dependent

## Future Improvements
- [ ] Voice caching untuk performa lebih baik
- [ ] Custom voice settings (pitch, rate, volume)
- [ ] Voice download progress indicator
- [ ] Better error handling dan user feedback

