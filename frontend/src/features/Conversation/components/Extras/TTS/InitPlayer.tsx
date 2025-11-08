import { ChatMessageError, ChatTTS } from '@/types';
import { memo, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useNativeTTS } from '@/hooks/useNativeTTS';
import { useChatStore } from '@/store/chat';

import Player from './Player';

export interface TTSProps extends ChatTTS {
  content: string;
  id: string;
  loading?: boolean;
}

const InitPlayer = memo<TTSProps>(({ id, content, contentMd5, file }) => {
  const [isStart, setIsStart] = useState(false);
  const [error, setError] = useState<ChatMessageError>();
  const { t } = useTranslation('chat');

  const [ttsMessage, clearTTS] = useChatStore((s) => [s.ttsMessage, s.clearTTS]);

  const setDefaultError = useCallback(
    (err?: any) => {
      setError({ body: err, message: t('tts.responseError', { ns: 'error' }), type: 500 });
    },
    [t],
  );

  const { isGlobalLoading, start, stop } = useNativeTTS(content, {
    onError: (err) => {
      stop();
      setDefaultError(err);
      setIsStart(false); // Reset state on error
    },
    onSuccess: async () => {
      // Playback completed successfully
      // Mark as played in the store
      ttsMessage(id, { contentMd5, file: 'played', voice: 'default' });
      setIsStart(false); // Reset state after successful playback
    },
  });

  const handleInitStart = useCallback(() => {
    console.log('[InitPlayer] handleInitStart called', { isStart, content });
    if (isStart) {
      console.log('[InitPlayer] Already playing, skipping');
      return;
    }
    console.log('[InitPlayer] Calling start()');
    setIsStart(true);
    start();
  }, [isStart, start, content]);

  const handleDelete = useCallback(() => {
    stop();
    clearTTS(id);
  }, [stop, id]);

  const handleRetry = useCallback(() => {
    setError(undefined);
    start();
  }, [start]);

  // Disabled auto-start TTS - only start when user clicks the button
  // useEffect(() => {
  //   if (file) return;
  //   setTimeout(() => {
  //     handleInitStart();
  //   }, 100);
  // }, [file, handleInitStart]);

  return (
    <Player
      error={error}
      isLoading={isGlobalLoading}
      onDelete={handleDelete}
      onInitPlay={handleInitStart}
      onRetry={handleRetry}
    />
  );
});

export default InitPlayer;

