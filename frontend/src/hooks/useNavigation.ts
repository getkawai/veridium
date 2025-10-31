import { useMemo } from 'react';

import { routerSelectors, useRouterStore } from '@/store/router';

/**
 * Custom hook that returns the current pathname from Zustand router store
 * Compatible with next/navigation's usePathname
 */
export const usePathname = () => {
  return useRouterStore(routerSelectors.pathname);
};

/**
 * Custom hook that returns search params with URLSearchParams-like API
 * Compatible with next/navigation's useSearchParams
 */
export const useSearchParams = () => {
  const searchParams = useRouterStore(routerSelectors.searchParams);

  return useMemo(() => {
    // Create a URLSearchParams-like object
    const params = {
      get: (key: string) => searchParams[key] || null,
      has: (key: string) => key in searchParams,
      getAll: (key: string) => {
        const value = searchParams[key];
        return value ? [value] : [];
      },
      toString: () => {
        const entries = Object.entries(searchParams);
        if (entries.length === 0) return '';
        return entries.map(([key, value]) => `${key}=${encodeURIComponent(value)}`).join('&');
      },
      entries: () => Object.entries(searchParams),
      keys: () => Object.keys(searchParams),
      values: () => Object.values(searchParams),
      forEach: (callback: (value: string, key: string) => void) => {
        Object.entries(searchParams).forEach(([key, value]) => callback(value, key));
      },
    };

    return params;
  }, [searchParams]);
};

