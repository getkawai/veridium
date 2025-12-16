import { useGlobalStore } from '@/store/global';

/**
 * Returns the active tab key (chat/market/settings/...)
 */
export const useActiveTabKey = () => {
  return useGlobalStore((s) => s.sidebarKey);
};
