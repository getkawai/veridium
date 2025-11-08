import { useCallback, useState } from 'react';
import { TTSService } from '@@/github.com/kawai-network/veridium/internal/tts';

interface UseNativeTTSConfig {
  onError?: (error: any) => void;
  onSuccess?: () => void;
  voice?: string;
}

interface UseNativeTTSReturn {
  start: () => void;
  stop: () => void;
  isGlobalLoading: boolean;
}

/**
 * Hook for native Text-to-Speech using system speaker via Wails
 * Audio is played directly on the backend using native OS TTS (e.g., macOS `say` command)
 */
export const useNativeTTS = (text: string, config: UseNativeTTSConfig): UseNativeTTSReturn => {
  const [isLoading, setIsLoading] = useState(false);

  const start = useCallback(async () => {
    console.log('[useNativeTTS] start() called', { text, isLoading });
    
    if (!text || isLoading) {
      console.log('[useNativeTTS] Skipping - no text or already loading');
      return;
    }

    setIsLoading(true);
    console.log('[useNativeTTS] Calling TTSService.Speak()...');

    try {
      // Play speech directly on backend (blocks until playback completes)
      if (config.voice) {
        console.log('[useNativeTTS] Using voice:', config.voice);
        await TTSService.SpeakWithVoice(text, config.voice);
      } else {
        console.log('[useNativeTTS] Using default voice');
        await TTSService.Speak(text);
      }

      console.log('[useNativeTTS] Playback completed successfully');
      // Playback completed successfully
      config.onSuccess?.();
    } catch (error) {
      console.error('[useNativeTTS] TTS failed:', error);
      config.onError?.(error);
    } finally {
      setIsLoading(false);
      console.log('[useNativeTTS] Loading state cleared');
    }
  }, [text, isLoading, config]);

  const stop = useCallback(async () => {
    try {
      // Stop any ongoing speech on backend
      await TTSService.Stop();
      setIsLoading(false);
    } catch (error) {
      console.error('Failed to stop TTS:', error);
    }
  }, []);

  return {
    start,
    stop,
    isGlobalLoading: isLoading,
  };
};

