import { ModelProvider } from '@/model-bank';

import { createOpenAICompatibleRuntime } from '../../core/openaiCompatibleFactory';

export const LobeUpstageAI = createOpenAICompatibleRuntime({
  baseURL: 'https://api.upstage.ai/v1/solar',
  debug: {
    chatCompletion: () => false,
  },
  provider: ModelProvider.Upstage,
});
