import { Button, InputPassword } from '@lobehub/ui';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

// import { useChatStore } from '@/store/chat';
// import { useUserStore } from '@/store/user';
// import { keyVaultsConfigSelectors } from '@/store/user/selectors';

// Dummy implementations for UI development
const useUserStore = (selector?: any) => {
  if (selector) {
    return selector({
      password: '',
      updateKeyVaults: (config: any) => console.log('Mock updateKeyVaults called with:', config),
    });
  }

  return {
    password: '',
    updateKeyVaults: (config: any) => console.log('Mock updateKeyVaults called with:', config),
  };
};

const keyVaultsConfigSelectors = {
  password: (state: any) => state.password,
};

const useChatStore = (selector?: any) => {
  if (selector) {
    return selector({
      regenerateMessage: (id: string) => console.log('Mock regenerateMessage called with:', id),
      deleteMessage: (id: string) => console.log('Mock deleteMessage called with:', id),
    });
  }

  return {
    regenerateMessage: (id: string) => console.log('Mock regenerateMessage called with:', id),
    deleteMessage: (id: string) => console.log('Mock deleteMessage called with:', id),
  };
};

import { FormAction } from './style';

interface AccessCodeFormProps {
  id: string;
}

const AccessCodeForm = memo<AccessCodeFormProps>(({ id }) => {
  const { t } = useTranslation('error');
  const [password, updateKeyVaults] = useUserStore((s) => [
    keyVaultsConfigSelectors.password(s),
    s.updateKeyVaults,
  ]);
  const [resend, deleteMessage] = useChatStore((s) => [s.regenerateMessage, s.deleteMessage]);

  return (
    <>
      <FormAction
        avatar={'🗳'}
        description={t('unlock.password.description')}
        title={t('unlock.password.title')}
      >
        <InputPassword
          autoComplete={'new-password'}
          onChange={(e) => {
            updateKeyVaults({ password: e.target.value });
          }}
          placeholder={t('unlock.password.placeholder')}
          value={password}
          variant={'filled'}
        />
      </FormAction>
      <Flexbox gap={12}>
        <Button
          onClick={() => {
            resend(id);
            deleteMessage(id);
          }}
          type={'primary'}
        >
          {t('unlock.confirm')}
        </Button>
        <Button
          onClick={() => {
            deleteMessage(id);
          }}
        >
          {t('unlock.closeMessage')}
        </Button>
      </Flexbox>
    </>
  );
});

export default AccessCodeForm;
