import { ChatMessageError } from '@/types';
import { SpeechRecognitionOptions, useSpeechRecognition } from '@lobehub/tts/react';
import isEqual from 'fast-deep-equal';
import { memo, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { SWRConfiguration } from 'swr';

import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/slices/chat';
import { useChatStore } from '@/store/chat';
import { chatSelectors } from '@/store/chat/selectors';
import { useGlobalStore } from '@/store/global';
import { globalGeneralSelectors } from '@/store/global/selectors';
import { useUserStore } from '@/store/user';
import { settingsSelectors } from '@/store/user/selectors';
import { getMessageError } from '@/utils/fetch';

import CommonSTT from './common';

interface STTConfig extends SWRConfiguration {
  onTextChange: (value: string) => void;
}

const useBrowserSTT = (config: STTConfig) => {
  const ttsSettings = useUserStore(settingsSelectors.currentTTS, isEqual);
  const ttsAgentSettings = useAgentStore(agentSelectors.currentAgentTTS, isEqual);
  const locale = useGlobalStore(globalGeneralSelectors.currentLanguage);

  const autoStop = ttsSettings.sttAutoStop;
  const sttLocale =
    ttsAgentSettings?.sttLocale && ttsAgentSettings.sttLocale !== 'auto'
      ? ttsAgentSettings.sttLocale
      : locale;

  // Check if media devices are available (not available in Wails desktop environment)
  const mediaDevicesAvailable = typeof navigator !== 'undefined' && 
    navigator.mediaDevices && 
    typeof navigator.mediaDevices.getUserMedia === 'function';

  // If media devices are not available, return a mock implementation
  if (!mediaDevicesAvailable) {
    return {
      start: () => {
        // Call onError with SWR-compatible signature (error, key, config)
        config.onError?.(
          new Error('Media devices not available in desktop environment'),
          '',
          { retryCount: 0, dedupe: false } as any
        );
      },
      stop: () => {},
      isLoading: false,
      isRecording: false,
      formattedTime: '00:00',
      time: 0,
      response: undefined,
    };
  }

  return useSpeechRecognition(sttLocale, {
    ...config,
    autoStop,
  } as SpeechRecognitionOptions);
};

const BrowserSTT = memo<{ mobile?: boolean }>(({ mobile }) => {
  const [error, setError] = useState<ChatMessageError>();
  const { t } = useTranslation('chat');

  const [loading, updateInputMessage] = useChatStore((s) => [
    chatSelectors.isAIGenerating(s),
    s.updateInputMessage,
  ]);

  const setDefaultError = useCallback(
    (err?: any) => {
      // Check if this is a media devices error
      const isMediaDeviceError = err?.message?.includes('Media devices not available');
      const errorMessage = isMediaDeviceError 
        ? 'Speech-to-text is not available in the desktop version. Please use the web version or text input.'
        : t('stt.responseError', { ns: 'error' });
      
      setError({ body: err, message: errorMessage, type: 500 });
    },
    [t],
  );

  const { start, isLoading, stop, formattedTime, time, response, isRecording } = useBrowserSTT({
    onError: (err) => {
      stop();
      setDefaultError(err);
    },
    onErrorRetry: (err) => {
      stop();
      setDefaultError(err);
    },
    onSuccess: async () => {
      if (!response) return;
      if (response.status === 200) return;
      const message = await getMessageError(response);
      if (message) {
        setError(message);
      } else {
        setDefaultError();
      }
      stop();
    },
    onTextChange: (text) => {
      if (loading) stop();
      if (text) updateInputMessage(text);
    },
  });

  const desc = t('stt.action');

  const handleTriggerStartStop = useCallback(() => {
    if (loading) return;
    if (!isLoading) {
      start();
    } else {
      stop();
    }
  }, [loading, isLoading, start, stop]);

  const handleCloseError = useCallback(() => {
    setError(undefined);
    stop();
  }, [stop]);

  const handleRetry = useCallback(() => {
    setError(undefined);
    start();
  }, [start]);

  return (
    <CommonSTT
      desc={desc}
      error={error}
      formattedTime={formattedTime}
      handleCloseError={handleCloseError}
      handleRetry={handleRetry}
      handleTriggerStartStop={handleTriggerStartStop}
      isLoading={isLoading}
      isRecording={isRecording}
      mobile={mobile}
      time={time}
    />
  );
});

export default BrowserSTT;
