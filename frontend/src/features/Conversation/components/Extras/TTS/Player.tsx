import { ChatMessageError } from '@/types';
import { ActionIcon, Alert, Button, Highlighter } from '@lobehub/ui';
import { Volume2, TrashIcon } from 'lucide-react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

interface PlayerProps {
  error?: ChatMessageError;
  isLoading?: boolean;
  onDelete: () => void;
  onInitPlay?: () => void;
  onRetry?: () => void;
}

/**
 * Simplified TTS Player for backend audio playback
 * Audio is played directly on system speaker via backend, not in browser
 */
const Player = memo<PlayerProps>(({ onRetry, error, onDelete, isLoading, onInitPlay }) => {
  const { t } = useTranslation('chat');

  return (
    <Flexbox align={'center'} gap={8} horizontal style={{ minWidth: 100 }}>
      {error ? (
        <Alert
          action={
            <Button onClick={onRetry} size={'small'} type={'primary'}>
              {t('retry', { ns: 'common' })}
            </Button>
          }
          closable
          extra={
            error.body && (
              <Highlighter actionIconSize={'small'} language={'json'} variant={'borderless'}>
                {JSON.stringify(error.body, null, 2)}
              </Highlighter>
            )
          }
          message={error.message}
          onClose={onDelete}
          style={{ alignItems: 'center', width: '100%' }}
          type="error"
        />
      ) : (
        <>
          <ActionIcon
            icon={Volume2}
            loading={isLoading}
            onClick={() => {
              console.log('[Player] Volume button clicked, calling onInitPlay');
              onInitPlay?.();
            }}
            size={'small'}
            title={t('tts.play')}
          />
          <ActionIcon icon={TrashIcon} onClick={onDelete} size={'small'} title={t('tts.clear')} />
        </>
      )}
    </Flexbox>
  );
});

export default Player;
