import { useGlobalStore } from '@/store/global';

export const useShowMobileWorkspace = () => {
  return useGlobalStore((s) => !!s.status.mobileShowTopic);
};
