import { useMemo } from 'react';

import { useAgentStore } from '@/store/agent';

export const useOpenChatSettings = (defaultTab?: any) => {
  return useMemo(() => {
    return (tab?: any) => {
      useAgentStore.setState({ showAgentSetting: true });
    };
  }, []);
};
