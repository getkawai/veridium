import { useEffect } from 'react';

import { useAiInfraStore } from '@/store/aiInfra';

/**
 * Custom hook to fetch AI provider runtime state
 * Replaces the old SWR-based useFetchAiProviderRuntimeState
 */
export function useFetchAiProviderRuntimeState(isLogin: boolean | null | undefined) {
  const internal_fetchAiProviderRuntimeState = useAiInfraStore(
    (s) => s.internal_fetchAiProviderRuntimeState,
  );

  useEffect(() => {
    internal_fetchAiProviderRuntimeState(isLogin);
  }, [isLogin, internal_fetchAiProviderRuntimeState]);
}

