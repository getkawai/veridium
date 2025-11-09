import { useEffect } from 'react';

import { useGlobalStore } from '@/store/global';
import { systemStatusSelectors } from '@/store/global/selectors';
import { useToolStore } from '@/store/tool';

export const useCheckPluginsIsInstalled = (plugins: string[]) => {
  const isDBInited = useGlobalStore(systemStatusSelectors.isDBInited);
  const internal_checkPluginsIsInstalled = useToolStore((s) => s.internal_checkPluginsIsInstalled);

  useEffect(() => {
    internal_checkPluginsIsInstalled(isDBInited, plugins);
  }, [isDBInited, plugins, internal_checkPluginsIsInstalled]);
};
