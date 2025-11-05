import { ChatMessageError } from '@/types';
import { memo, useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useNativeSTT } from '@/hooks/useNativeSTT';
import { useChatStore } from '@/store/chat';
import { chatSelectors } from '@/store/chat/selectors';
import { useChatInputStore } from '../../store';

import CommonSTT from './common';

const NativeSTT = memo<{ mobile?: boolean }>(({ mobile }) => {
  const [error, setError] = useState<ChatMessageError>();
  const { t } = useTranslation('chat');

  const loading = useChatStore(chatSelectors.isAIGenerating);
  const [editor, getMarkdownContent] = useChatInputStore((s) => [
    s.editor,
    s.getMarkdownContent,
  ]);

  const setDefaultError = useCallback(
    (err?: any) => {
      const errorMessage = err?.message || t('stt.responseError', { ns: 'error' });
      setError({ body: err, message: errorMessage, type: 500 });
    },
    [t],
  );

  const { start, isLoading, stop, formattedTime, time, isRecording } = useNativeSTT({
    onError: (err) => {
      // Don't call stop() here - it's already stopped in the hook
      setDefaultError(err);
    },
    onSuccess: async () => {
      // Success handled in onTextChange
    },
    onTextChange: (text) => {
      // Don't call stop() here - transcription is already complete
      console.log('STT onTextChange called with:', text, 'length:', text?.length);
      
      if (!editor) {
        console.warn('Editor not initialized');
        return;
      }
      
      if (text) {
        // Get current content from editor
        const currentContent = getMarkdownContent();
        console.log('Current editor content:', currentContent);
        
        // Append to existing content if there's any, otherwise replace
        const newContent = currentContent ? `${currentContent} ${text}` : text;
        console.log('Setting editor content to:', newContent);
        
        // Update editor content
        editor.setDocument('markdown', newContent);
      } else {
        console.warn('Text is empty or falsy, not updating input');
      }
    },
  });

  const desc = t('stt.action');

  const handleTriggerStartStop = useCallback(() => {
    if (loading) return;
    if (!isLoading && !isRecording) {
      start();
    } else {
      stop();
    }
  }, [loading, isLoading, isRecording, start, stop]);

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

export default NativeSTT;

