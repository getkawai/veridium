import { useEffect } from 'react';
import { StateCreator } from 'zustand/vanilla';

import { globalService } from '@/services/global';
import { GlobalRuntimeConfig } from '@/types/serverConfig';

import type { ServerConfigStore } from './store';

export interface ServerConfigAction {
  useInitServerConfig: () => void;
}

export const createServerConfigSlice: StateCreator<
  ServerConfigStore,
  [['zustand/devtools', never]],
  [],
  ServerConfigAction
> = (set) => ({
  useInitServerConfig: () => {
    useEffect(() => {
      const initConfig = async () => {
        try {
          const data = await globalService.getGlobalConfig();
          set(
            { featureFlags: data.serverFeatureFlags, serverConfig: data.serverConfig },
            false,
            'initServerConfig',
          );
        } catch (error) {
          console.error('[useInitServerConfig] Error:', error);
        }
      };

      initConfig();
    }, []);
  },
});
