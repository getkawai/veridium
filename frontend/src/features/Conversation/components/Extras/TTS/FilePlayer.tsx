import { memo, useCallback, useState } from 'react';

import { useNativeTTS } from '@/hooks/useNativeTTS';
import { useChatStore } from '@/store/chat';

import { TTSProps } from './InitPlayer';
import Player from './Player';

/**
 * FilePlayer for replaying TTS
 * Since we use backend playback, we regenerate TTS on each play
 * (TTS generation is fast, no need for file caching)
 */
const FilePlayer = memo<TTSProps>(({ content, id }) => {
  const [clearTTS] = useChatStore((s) => [s.clearTTS]);
  const [error, setError] = useState<any>();

  const { isGlobalLoading, start, stop } = useNativeTTS(content, {
    onError: (err) => {
      stop();
      setError(err);
    },
    onSuccess: async () => {
      // Playback completed successfully
    },
  });

  const handleDelete = useCallback(() => {
    clearTTS(id);
  }, [id]);

  const handlePlay = useCallback(() => {
    console.log('[FilePlayer] handlePlay called', { content });
    setError(undefined);
    start();
  }, [start, content]);

  const handleRetry = useCallback(() => {
    setError(undefined);
    start();
  }, [start]);

  return (
    <Player
      error={error}
      isLoading={isGlobalLoading}
      onDelete={handleDelete}
      onInitPlay={handlePlay}
      onRetry={handleRetry}
    />
  );
});

export default FilePlayer;
