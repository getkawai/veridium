import { useCallback, useEffect, useRef, useState } from 'react';
import * as TTSService from '@@/github.com/kawai-network/veridium/services/ttsservice';

interface UseNativeTTSConfig {
  onError?: (error: any) => void;
  onSuccess?: () => void;
  onUpload?: (voice: string, audioFile: string) => Promise<void>;
}

interface UseNativeTTSReturn {
  audio?: string;
  start: () => void;
  stop: () => void;
  isGlobalLoading: boolean;
  response?: Response;
}

/**
 * Hook for native Text-to-Speech using macOS `say` command via Wails
 */
export const useNativeTTS = (text: string, config: UseNativeTTSConfig): UseNativeTTSReturn => {
  const [isLoading, setIsLoading] = useState(false);
  const [audio, setAudio] = useState<string>();
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const abortControllerRef = useRef<AbortController | null>(null);

  const start = useCallback(async () => {
    if (!text || isLoading) return;

    setIsLoading(true);
    abortControllerRef.current = new AbortController();

    try {
      // Generate speech and get audio data directly from backend
      const audioData = await TTSService.SpeakToAudio(text);

      // Convert to blob URL
      const audioBlob = new Blob([audioData], { type: 'audio/aiff' });
      const audioUrl = URL.createObjectURL(audioBlob);

      setAudio(audioUrl);

      // If onUpload is provided, save the audio data
      if (config.onUpload) {
        // Create a temporary file for upload if needed
        // For now, just pass the audio data as base64 or skip this
        await config.onUpload('default', audioUrl);
      }

      config.onSuccess?.();
    } catch (error) {
      console.error('TTS failed:', error);
      config.onError?.(error);
    } finally {
      setIsLoading(false);
    }
  }, [text, isLoading, config]);

  const stop = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.pause();
      audioRef.current.currentTime = 0;
    }
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    setIsLoading(false);
  }, []);

  // Play audio when available
  useEffect(() => {
    if (audio) {
      audioRef.current = new Audio(audio);
      audioRef.current.play().catch((error) => {
        console.error('Failed to play audio:', error);
        config.onError?.(error);
      });
    }
  }, [audio, config]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      stop();
      if (audio) {
        URL.revokeObjectURL(audio);
      }
    };
  }, [audio, stop]);

  return {
    audio,
    start,
    stop,
    isGlobalLoading: isLoading,
  };
};

