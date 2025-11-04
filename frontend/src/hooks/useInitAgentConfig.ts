import { useAgentStore } from '@/store/agent';
import { useSessionStore } from '@/store/session';

/**
 * @deprecated This hook is deprecated in favor of centralized data loading.
 * Agent configs are now loaded at app startup via StoreInitialization.
 * Components should read directly from agentStore instead of triggering fetches.
 *
 * This hook remains for backward compatibility but will be removed in future versions.
 *
 * Migration Guide:
 * Instead of:
 *   const { isLoading } = useInitAgentConfig();
 *
 * Use:
 *   const activeId = useSessionStore((s) => s.activeId);
 *   const isConfigLoaded = useAgentStore((s) => !!s.agentConfigInitMap[activeId]);
 *   const isLoading = !isConfigLoaded;
 *
 * If a targetAgentId is provided, use it to fetch the agent config directly.
 * Otherwise, use the active session id to fetch the config.
 */
export const useInitAgentConfig = (targetAgentId?: string) => {
  // Return loading state based on agentConfigInitMap instead of fetching
  const activeId = useSessionStore((s) => s.activeId);
  const sessionId = targetAgentId || activeId;

  const isConfigLoaded = useAgentStore((s) => !!s.agentConfigInitMap[sessionId]);

  return {
    isLoading: !isConfigLoaded,
    data: undefined, // Components should read from store directly
  };
};
