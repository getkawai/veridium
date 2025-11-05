import { ChatMessageError, ChatTTS } from '@/types';
import { memo, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useNativeTTS } from '@/hooks/useNativeTTS';
import { useChatStore } from '@/store/chat';
import { useFileStore } from '@/store/file';

import Player from './Player';

export interface TTSProps extends ChatTTS {
  content: string;
  id: string;
  loading?: boolean;
}

const InitPlayer = memo<TTSProps>(({ id, content, contentMd5, file }) => {
  const [isStart, setIsStart] = useState(false);
  const [error, setError] = useState<ChatMessageError>();
  const uploadTTS = useFileStore((s) => s.uploadTTSByArrayBuffers);
  const { t } = useTranslation('chat');

  const [ttsMessage, clearTTS] = useChatStore((s) => [s.ttsMessage, s.clearTTS]);

  const setDefaultError = useCallback(
    (err?: any) => {
      setError({ body: err, message: t('tts.responseError', { ns: 'error' }), type: 500 });
    },
    [t],
  );

  const { isGlobalLoading, audio, start, stop } = useNativeTTS(content, {
    onError: (err) => {
      stop();
      setDefaultError(err);
    },
    onSuccess: async () => {
      // Success
    },
    onUpload: async (currentVoice, audioFilePath) => {
      // Read audio file and upload
      try {
        const audioData = await window.go.main.NodeFsService.ReadFile(audioFilePath);
        const fileID = await uploadTTS(id, [audioData]);
        ttsMessage(id, { contentMd5, file: fileID, voice: currentVoice });
      } catch (err) {
        console.error('Failed to upload TTS audio:', err);
      }
    },
  });

  const handleInitStart = useCallback(() => {
    if (isStart) return;
    start();
    setIsStart(true);
  }, [isStart, start]);

  const handleDelete = useCallback(() => {
    stop();
    clearTTS(id);
  }, [stop, id]);

  const handleRetry = useCallback(() => {
    setError(undefined);
    start();
  }, [start]);

  useEffect(() => {
    if (file) return;
    setTimeout(() => {
      handleInitStart();
    }, 100);
  }, [file, handleInitStart]);

  return (
    <Player
      audio={audio}
      error={error}
      isLoading={isGlobalLoading}
      onDelete={handleDelete}
      onInitPlay={handleInitStart}
      onRetry={handleRetry}
    />
  );
});

export default InitPlayer;

