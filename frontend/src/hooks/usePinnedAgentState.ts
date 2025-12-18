import { useMemo } from 'react';
import { useGlobalStore } from '@/store/global';

export const usePinnedAgentState = () => {
  const isPinnedAgent = useGlobalStore((s) => !!s.status.isPinnedAgent);
  const updateSystemStatus = useGlobalStore((s) => s.updateSystemStatus);

  const actions = useMemo(
    () => ({
      pinAgent: () => updateSystemStatus({ isPinnedAgent: true }),
      setIsPinned: (value: boolean | ((prev: boolean) => boolean)) => {
        const newValue = typeof value === 'function' ? value(isPinnedAgent) : value;
        updateSystemStatus({ isPinnedAgent: newValue });
      },
      togglePinAgent: () => {
        updateSystemStatus({ isPinnedAgent: !isPinnedAgent });
      },
      unpinAgent: () => updateSystemStatus({ isPinnedAgent: false }),
    }),
    [isPinnedAgent, updateSystemStatus],
  );

  return [isPinnedAgent, actions] as const;
};
