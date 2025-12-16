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
        enableGroupChat: false,
        showCreateSession: true,
        isAgentEditable: true,
      },
    };
  };
}

export const globalService = new GlobalService();
