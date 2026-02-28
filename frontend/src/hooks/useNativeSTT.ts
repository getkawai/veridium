import { useCallback, useEffect, useRef, useState } from 'react';
import { Events } from '@wailsio/runtime';
import * as AudioRecorderService from '@@/github.com/kawai-network/veridium/internal/audio_recorder/audiorecorderservice';

interface UseNativeSTTConfig {
  onTextChange: (text: string) => void;
  onError?: (error: any) => void;
  onSuccess?: () => void;
}

interface UseNativeSTTReturn {
  start: () => void;
  stop: () => void;
  isLoading: boolean;
  isRecording: boolean;
  formattedTime: string;
  time: number;
  response?: Response;
}

/**
 * Hook for native Speech-to-Text using Whisper via Wails
 * Uses native OS audio recording (sox/arecord/ffmpeg)
 */
export const useNativeSTT = (config: UseNativeSTTConfig): UseNativeSTTReturn => {
  const [isRecording, setIsRecording] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [time, setTime] = useState(0);
  const timerRef = useRef<number | null>(null);
  const isStoppingRef = useRef<boolean>(false);

  const formattedTime = `${Math.floor(time / 60)
    .toString()
    .padStart(2, '0')}:${(time % 60).toString().padStart(2, '0')}`;

  const start = useCallback(async () => {
    try {
      // Start native audio recording first
      await AudioRecorderService.StartRecording();
      
      // Update state after recording started successfully
      setIsRecording(true);
      setTime(0);

      // Start timer
      timerRef.current = window.setInterval(() => {
        setTime((prev) => prev + 1);
      }, 1000);
    } catch (error) {
      console.error('Failed to start recording:', error);
      config.onError?.(error);
      setIsRecording(false);
    }
  }, [config]);

  const stop = useCallback(async () => {
    // Prevent multiple simultaneous stop calls
    if (isStoppingRef.current) {
      console.warn('Stop already in progress, ignoring duplicate call');
      return;
    }

    if (!isRecording) {
      console.warn('Stop called but not recording');
      return;
    }

    // Set flag to prevent re-entry
    isStoppingRef.current = true;

    // Stop timer first
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }

    // Update UI state
    setIsRecording(false);
    setIsLoading(true);

    try {
      // Stop native audio recording and get the file path
      const audioPath: string = await AudioRecorderService.StopRecording();

      console.log('Audio file saved to:', audioPath);
      throw new Error('Native STT backend is not available in this build.');
    } catch (error) {
      console.error('Transcription failed:', error);
      config.onError?.(error);
    } finally {
      setIsLoading(false);
      // Reset flag after stop completes
      isStoppingRef.current = false;
    }
  }, [isRecording, config]);

  // Listen to recording events
  useEffect(() => {
    const unsubscribeStarted = Events.On('audio:recording:started', (ev: any) => {
      const path = ev.data as string;
      console.log('Recording started:', path);
    });

    const unsubscribeStopped = Events.On('audio:recording:stopped', (ev: any) => {
      const path = ev.data as string;
      console.log('Recording stopped:', path);
    });

    const unsubscribeCancelled = Events.On('audio:recording:cancelled', () => {
      console.log('Recording cancelled');
      setIsRecording(false);
      setIsLoading(false);
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    });

    return () => {
      unsubscribeStarted();
      unsubscribeStopped();
      unsubscribeCancelled();
    };
  }, []);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
      // Only cancel if actually recording
      // Don't cancel if already stopped (prevents "not recording" error)
      if (isRecording && !isStoppingRef.current) {
        AudioRecorderService.CancelRecording().catch((err) => {
          // Ignore "not recording" errors
          if (!err.message?.includes('not recording')) {
            console.error('Failed to cancel recording:', err);
          }
        });
      }
    };
  }, [isRecording]);

  return {
    start,
    stop,
    isLoading,
    isRecording,
    formattedTime,
    time,
  };
};
