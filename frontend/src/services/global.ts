import type { PartialDeep } from 'type-fest';

import { LobeAgentConfig } from '@/types/agent';
import { GlobalRuntimeConfig } from '@/types/serverConfig';

const VERSION_URL = 'https://registry.npmmirror.com/@lobehub/chat/latest';

class GlobalService {
  /**
   * get latest version from npm
   */
  getLatestVersion = async (): Promise<string> => {
    const res = await fetch(VERSION_URL);
    const data = await res.json();

    return data['version'];
  };

  getGlobalConfig = async (): Promise<GlobalRuntimeConfig> => {
    // Mock global config data
    return {
      serverConfig: {
        aiProvider: {
          openai: {
            enabled: true,
            enabledModels: ['gpt-4o', 'gpt-4o-mini', 'gpt-3.5-turbo'],
            fetchOnClient: false,
          },
          anthropic: {
            enabled: true,
            enabledModels: ['claude-3-opus', 'claude-3-sonnet', 'claude-3-haiku'],
            fetchOnClient: false,
          },
        },
        telemetry: {
          langfuse: false,
        },
        enableUploadFileToServer: true,
        enabledAccessCode: false,
      },
      serverFeatureFlags: {
        showLLM: true,
        showProvider: true,
        showPinList: true,
        showOpenAIApiKey: true,
        showOpenAIProxyUrl: true,
        showApiKeyManage: true,
        enablePlugins: true,
        showDalle: true,
        showAiImage: true,
        showChangelog: true,
        enableCheckUpdates: true,
        showWelcomeSuggest: true,
        enableClerkSignUp: true,
        enableKnowledgeBase: true,
        enableRAGEval: true,
        showCloudPromotion: false,
        showMarket: true,
        enableSTT: true,
        hideGitHub: false,
        hideDocs: false,
        enableGroupChat: true,
        showCreateSession: true,
        isAgentEditable: true,
      },
    };
  };

  getDefaultAgentConfig = async (): Promise<PartialDeep<LobeAgentConfig>> => {
    // Mock default agent config data
    return {
      model: 'gpt-4o-mini',
      provider: 'openai',
      systemRole: 'You are a helpful AI assistant.',
      chatConfig: {
        autoCreateTopicThreshold: 2,
        enableAutoCreateTopic: true,
        enableCompressHistory: true,
        enableHistoryCount: true,
        historyCount: 20,
        inputTemplate: '',
        searchMode: 'auto',
      },
      params: {
        frequency_penalty: 0,
        max_tokens: 4000,
        presence_penalty: 0,
        temperature: 0.7,
        top_p: 1,
      },
      tts: {
        showAllLocaleVoice: false,
        sttLocale: 'auto',
        ttsService: 'openai',
        voice: {
          openai: 'alloy',
        },
      },
    };
  };
}

export const globalService = new GlobalService();
