import isEqual from 'fast-deep-equal';
import { useContext } from 'react';

import { isDeprecatedEdition } from '@/const/version';
// import { aiProviderSelectors, useAiInfraStore } from '@/store/aiInfra';
// import { useUserStore } from '@/store/user';
// import { keyVaultsConfigSelectors } from '@/store/user/selectors';

// Dummy implementations for UI development
const useUserStore = (selector?: any) => {
  if (selector) {
    return selector({
      updateKeyVaultConfig: (provider: string, config: any) => console.log('Mock updateKeyVaultConfig called with:', provider, config),
    });
  }

  return {
    updateKeyVaultConfig: (provider: string, config: any) => console.log('Mock updateKeyVaultConfig called with:', provider, config),
  };
};

const keyVaultsConfigSelectors = {
  getVaultByProvider: (provider: string) => (state: any) => ({
    apiKey: 'mock-api-key',
    baseURL: 'https://api.example.com',
  }),
};

const aiProviderSelectors = {
  providerConfigById: (provider: string) => (state: any) => ({
    keyVaults: {
      apiKey: 'mock-api-key',
      baseURL: 'https://api.example.com',
    },
  }),
};

const useAiInfraStore = (selector?: any, comparator?: any) => {
  const mockProviderConfig = {
    keyVaults: {
      apiKey: 'mock-api-key',
      baseURL: 'https://api.example.com',
    },
  };

  if (selector) {
    const result = selector({
      updateAiProviderConfig: (id: string, config: any) => console.log('Mock updateAiProviderConfig called with:', id, config),
      aiProviderRuntimeConfig: {
        openai: mockProviderConfig,
        anthropic: mockProviderConfig,
        google: mockProviderConfig,
      },
    });

    return result;
  }

  return {
    updateAiProviderConfig: (id: string, config: any) => console.log('Mock updateAiProviderConfig called with:', id, config),
  };
};

import { LoadingContext } from './LoadingContext';

export const useApiKey = (provider: string) => {
  const [apiKey, baseURL, setConfig] = useUserStore((s) => [
    keyVaultsConfigSelectors.getVaultByProvider(provider as any)(s)?.apiKey,
    keyVaultsConfigSelectors.getVaultByProvider(provider as any)(s)?.baseURL,
    s.updateKeyVaultConfig,
  ]);
  const { setLoading } = useContext(LoadingContext);
  const updateAiProviderConfig = useAiInfraStore((s) => s.updateAiProviderConfig);
  const data = useAiInfraStore(aiProviderSelectors.providerConfigById(provider), isEqual);

  // TODO: remove this in V2
  if (isDeprecatedEdition) return { apiKey, baseURL, setConfig };
  //

  return {
    apiKey: data?.keyVaults.apiKey,
    baseURL: data?.keyVaults?.baseURL,
    setConfig: async (id: string, params: Record<string, string>) => {
      const next = { ...data?.keyVaults, ...params };
      if (isEqual(data?.keyVaults, next)) return;

      setLoading(true);
      await updateAiProviderConfig(id, {
        keyVaults: { ...data?.keyVaults, ...params },
      });
      setLoading(false);
    },
  };
};
