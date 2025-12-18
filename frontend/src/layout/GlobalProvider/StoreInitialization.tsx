'use client';

import { memo, useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { createStoreUpdater } from 'zustand-utils';

import { enableNextAuth } from '@/const/auth';
import { useFetchAiProviderRuntimeState } from '@/hooks/useFetchAiProviderRuntimeState';
import { useAgentStore } from '@/store/agent';
import { useGlobalStore } from '@/store/global';
import { useServerConfigStore } from '@/store/serverConfig';
import { serverConfigSelectors } from '@/store/serverConfig/selectors';
import { useUserStore } from '@/store/user';
import { authSelectors } from '@/store/user/selectors';
import { waitForWailsRuntime } from '@/utils/wailsRuntime';

const StoreInitialization = memo(() => {
  const [isRuntimeReady, setIsRuntimeReady] = useState(false);

  // Wait for Wails runtime to be ready before initializing stores
  useEffect(() => {
    waitForWailsRuntime().then(() => {
      setIsRuntimeReady(true);
    });
  }, []);

  // prefetch error ns to avoid don't show error content correctly
  useTranslation('error');

  const [isLogin, isSignedIn, useInitUserState] = useUserStore((s) => [
    authSelectors.isLogin(s),
    s.isSignedIn,
    s.useInitUserState,
  ]);

  const { serverConfig } = useServerConfigStore();

  const useInitSystemStatus = useGlobalStore((s) => s.useInitSystemStatus);

  const useInitInboxAgentStore = useAgentStore((s) => s.useInitInboxAgentStore);
  const useLoadAllAgentConfigs = useAgentStore((s) => s.useLoadAllAgentConfigs);

  // init the system preference
  useInitSystemStatus();

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
   *
   * Also ensure runtime is ready before making any DB calls to avoid
   * "Unable to parse request body as JSON" errors during initialization.
   */
  const isLoginOnInit = Boolean(enableNextAuth ? isSignedIn : isLogin) && isRuntimeReady;

  // init inbox agent and default agent config
  useInitInboxAgentStore(isLoginOnInit, serverConfig.defaultAgent?.config);

  // batch load all agent configs after sessions are loaded
  useLoadAllAgentConfigs(isLoginOnInit);

  // init user provider key vaults
  useFetchAiProviderRuntimeState(isLoginOnInit);

  // init user state
  useInitUserState(isLoginOnInit, serverConfig, {
    onSuccess: (state) => {
      if (state.isOnboard === false) {
        // routerPush('/onboard');
      }
    },
  });

  const useStoreUpdater = createStoreUpdater(useGlobalStore);

  useStoreUpdater('isMobile', false);

  return null;
});

export default StoreInitialization;
