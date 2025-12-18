import { useMemo } from 'react';

import { useAgentStore } from '@/store/agent';

export const useOpenChatSettings = () => {
  return useMemo(() => {
    return () => {
      useAgentStore.setState({ showAgentSetting: true });
    };
  }, []);
};
