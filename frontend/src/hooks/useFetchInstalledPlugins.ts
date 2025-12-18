import { useEffect } from 'react';

import { useToolStore } from '@/store/tool';

export const useFetchInstalledPlugins = () => {
  const [useFetchInstalledPlugins, isLoading] = useToolStore((s) => [
    s.useFetchInstalledPlugins,
    s.loadingInstallPlugins,
  ]);

  useEffect(() => {
    useFetchInstalledPlugins(true);
  }, []);

  return { isLoading };
};
