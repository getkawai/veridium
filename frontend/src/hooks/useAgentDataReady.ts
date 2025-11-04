import { useAgentStore } from '@/store/agent';
import { useGlobalStore } from '@/store/global';
import { systemStatusSelectors } from '@/store/global/selectors';
import { useSessionStore } from '@/store/session';

/**
 * Hook to check if all required data is loaded and ready
 * Use this in components that need agent data
 */
export const useAgentDataReady = () => {
  const isDBInited = useGlobalStore(systemStatusSelectors.isDBInited);
  const isSessionsLoaded = useSessionStore((s) => s.isSessionsFirstFetchFinished);
  const isAllAgentConfigsLoaded = useAgentStore((s) => s.isAllAgentConfigsLoaded);

  return isDBInited && isSessionsLoaded && isAllAgentConfigsLoaded;
};

