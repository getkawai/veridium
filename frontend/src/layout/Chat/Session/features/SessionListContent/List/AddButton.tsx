import { Button } from '@lobehub/ui';
import { Plus } from 'lucide-react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

// import { useActionSWR } from '@/libs/swr';
// import { useServerConfigStore } from '@/store/serverConfig';
// import { useSessionStore } from '@/store/session';

// Dummy implementations for development
const useSessionStore = (selector?: any) => {
  if (selector) {
    return selector({
      createSession: async (config?: any) => {
        console.log('Mock createSession called with:', config);
        return `session-${Date.now()}`;
      },
    });
  }
  return {
    createSession: async (config?: any) => {
      console.log('Mock createSession called with:', config);
      return `session-${Date.now()}`;
    },
  };
};

const useServerConfigStore = (selector?: any) => {
  if (selector) {
    return selector({ isMobile: false });
  }
  return { isMobile: false };
};

const useActionSWR = (key: any, action: () => Promise<any>) => {
  return {
    mutate: async () => {
      console.log(`Mock mutate called for ${key}`);
      try {
        await action();
      } catch (error) {
        console.error(`Mock action failed for ${key}:`, error);
      }
    },
    isValidating: false,
  };
};

const AddButton = memo<{ groupId?: string }>(({ groupId }) => {
  const { t } = useTranslation('chat');
  const createSession = useSessionStore((s) => s.createSession);
  const mobile = useServerConfigStore((s) => s.isMobile);
  const { mutate, isValidating } = useActionSWR(['session.createSession', groupId], () => {
    return createSession({ group: groupId });
  });

  return (
    <Flexbox flex={1} padding={mobile ? 16 : 0}>
      <Button
        block
        icon={Plus}
        loading={isValidating}
        onClick={() => mutate()}
        style={{
          marginTop: 8,
        }}
        variant={'filled'}
      >
        {t('newAgent')}
      </Button>
    </Flexbox>
  );
});

export default AddButton;
