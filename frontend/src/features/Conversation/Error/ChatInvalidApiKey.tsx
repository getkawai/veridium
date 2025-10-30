import { memo } from 'react';
import { useTranslation } from 'react-i18next';

import InvalidAPIKey from '@/components/InvalidAPIKey';
// import { useProviderName } from '@/hooks/useProviderName';
// import { useChatStore } from '@/store/chat';
import { GlobalLLMProviderKey } from '@/types/user/settings/modelProvider';

// Dummy implementations for UI development
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

const useProviderName = (provider: string) => {
  const providerNames: Record<string, string> = {
    openai: 'OpenAI',
    anthropic: 'Anthropic',
    google: 'Google',
    azure: 'Azure OpenAI',
    bedrock: 'Amazon Bedrock',
    ollama: 'Ollama',
    default: provider.charAt(0).toUpperCase() + provider.slice(1),
  };

  return providerNames[provider] || providerNames.default;
};

interface ChatInvalidAPIKeyProps {
  id: string;
  provider?: string;
}
const ChatInvalidAPIKey = memo<ChatInvalidAPIKeyProps>(({ id, provider }) => {
  const { t } = useTranslation('modelProvider');
  const { t: modelProviderErrorT } = useTranslation(['modelProvider', 'error']);
  const [resend, deleteMessage] = useChatStore((s) => [s.regenerateMessage, s.deleteMessage]);
  const providerName = useProviderName(provider as GlobalLLMProviderKey);

  return (
    <InvalidAPIKey
      bedrockDescription={t('bedrock.unlock.description')}
      description={modelProviderErrorT(`unlock.apiKey.description`, {
        name: providerName,
        ns: 'error',
      })}
      id={id}
      onClose={() => {
        deleteMessage(id);
      }}
      onRecreate={() => {
        resend(id);
        deleteMessage(id);
      }}
      provider={provider}
    />
  );
});

export default ChatInvalidAPIKey;
