import { useMemo } from 'react';

import { routerSelectors, useRouterStore } from '@/store/router';

export const usePinnedAgentState = () => {
  const searchParams = useRouterStore(routerSelectors.searchParams);
  const setSearchParam = useRouterStore((s) => s.setSearchParam);
  const removeSearchParam = useRouterStore((s) => s.removeSearchParam);

  const isPinned = searchParams.pinned === 'true';

  const actions = useMemo(
    () => ({
      pinAgent: () => setSearchParam('pinned', 'true'),
      setIsPinned: (value: boolean | ((prev: boolean) => boolean)) => {
        const newValue = typeof value === 'function' ? value(isPinned) : value;
        if (newValue) {
          setSearchParam('pinned', 'true');
        } else {
          removeSearchParam('pinned');
        }
      },
      togglePinAgent: () => {
        if (isPinned) {
          removeSearchParam('pinned');
        } else {
          setSearchParam('pinned', 'true');
        }
      },
      unpinAgent: () => removeSearchParam('pinned'),
    }),
    [isPinned, setSearchParam, removeSearchParam],
  );

  return [isPinned, actions] as const;
};
