'use client';

import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { createStoreUpdater } from 'zustand-utils';

import { enableNextAuth } from '@/const/auth';
import { useFetchAiProviderRuntimeState } from '@/hooks/useFetchAiProviderRuntimeState';
import { useIsMobile } from '@/hooks/useIsMobile';
import { useAgentStore } from '@/store/agent';
import { useAiInfraStore } from '@/store/aiInfra';
import { useGlobalStore } from '@/store/global';
import { systemStatusSelectors } from '@/store/global/selectors';
import { useRouterStore } from '@/store/router';
import { useServerConfigStore } from '@/store/serverConfig';
import { serverConfigSelectors } from '@/store/serverConfig/selectors';
import { useUserStore } from '@/store/user';
import { authSelectors } from '@/store/user/selectors';

const StoreInitialization = memo(() => {
  // prefetch error ns to avoid don't show error content correctly
  useTranslation('error');

  const routerPush = useRouterStore((s) => s.push);
  const [isLogin, isSignedIn, useInitUserState] = useUserStore((s) => [
    authSelectors.isLogin(s),
    s.isSignedIn,
    s.useInitUserState,
  ]);

  const { serverConfig } = useServerConfigStore();

  const useInitSystemStatus = useGlobalStore((s) => s.useInitSystemStatus);
  const useInitClientDB = useGlobalStore((s) => s.useInitClientDB);

  const useInitInboxAgentStore = useAgentStore((s) => s.useInitInboxAgentStore);
  const useLoadAllAgentConfigs = useAgentStore((s) => s.useLoadAllAgentConfigs);

  // init the system preference
  useInitSystemStatus();

  // init the client database (connects to backend-initialized database)
  useInitClientDB();

  // fetch server config
  const useFetchServerConfig = useServerConfigStore((s) => s.useInitServerConfig);
  useFetchServerConfig();

  // Update NextAuth status
  const useUserStoreUpdater = createStoreUpdater(useUserStore);
  const oAuthSSOProviders = useServerConfigStore(serverConfigSelectors.oAuthSSOProviders);
  useUserStoreUpdater('oAuthSSOProviders', oAuthSSOProviders);

  /**
   * The store function of `isLogin` will both consider the values of `enableAuth` and `isSignedIn`.
   * But during initialization, the value of `enableAuth` might be incorrect cause of the async fetch.
   * So we need to use `isSignedIn` only to determine whether request for the default agent config and user state.
   *
   * IMPORTANT: Explicitly convert to boolean to avoid passing null/undefined downstream,
   * which would cause unnecessary API requests with invalid login state.
   */
  const isDBInited = useGlobalStore(systemStatusSelectors.isDBInited);
  const isLoginOnInit = isDBInited ? Boolean(enableNextAuth ? isSignedIn : isLogin) : false;

  // init inbox agent and default agent config
  useInitInboxAgentStore(isLoginOnInit, serverConfig.defaultAgent?.config);

  // batch load all agent configs after sessions are loaded
  useLoadAllAgentConfigs(isDBInited && isLoginOnInit, isLoginOnInit);

  // init user provider key vaults
  useFetchAiProviderRuntimeState(isLoginOnInit);

  // init user state
  useInitUserState(isLoginOnInit, serverConfig, {
    onSuccess: (state) => {
      if (state.isOnboard === false) {
        routerPush('/onboard');
      }
    },
  });

  const useStoreUpdater = createStoreUpdater(useGlobalStore);

  const mobile = useIsMobile();

  useStoreUpdater('isMobile', mobile);

  return null;
});

export default StoreInitialization;
