import { useMemo } from 'react';

import { useRouterStore } from '@/store/router';

/**
 * Custom hook that returns router actions from Zustand store
 * Compatible with next/navigation's useRouter
 */
export const useRouter = () => {
  const push = useRouterStore((s) => s.push);
  const replace = useRouterStore((s) => s.replace);
  const back = useRouterStore((s) => s.back);
  const forward = useRouterStore((s) => s.forward);

  return useMemo(
    () => ({
      push,
      replace,
      back,
      forward,
      refresh: () => {
        // In SPA mode, refresh is a no-op
        console.warn('Router.refresh() is not needed in SPA mode');
      },
      prefetch: () => {
        // In SPA mode, prefetch is a no-op
      },
    }),
    [push, replace, back, forward],
  );
};

