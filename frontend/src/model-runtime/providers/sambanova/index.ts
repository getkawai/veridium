import { ModelProvider } from '@/model-bank';

import { createOpenAICompatibleRuntime } from '../../core/openaiCompatibleFactory';

export const LobeSambaNovaAI = createOpenAICompatibleRuntime({
  baseURL: 'https://api.sambanova.ai/v1',
  debug: {
    chatCompletion: () => false,
  },
  provider: ModelProvider.SambaNova,
});
