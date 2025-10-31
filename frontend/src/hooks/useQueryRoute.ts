import { type ParsedQuery } from 'query-string';
import { useMemo } from 'react';

import { routerSelectors, useRouterStore } from '@/store/router';

interface QueryRouteOptions {
  hash?: string;
  query?: ParsedQuery;
  replace?: boolean;
  replaceHash?: boolean;
  withHash?: boolean;
}

export const useQueryRoute = () => {
  const push = useRouterStore((s) => s.push);
  const replace = useRouterStore((s) => s.replace);
  const currentSearchParams = useRouterStore(routerSelectors.searchParams);

  return useMemo(
    () => ({
      push: (url: string, options: QueryRouteOptions = {}) => {
        const mergedQuery = options.replace
          ? options.query || {}
          : { ...currentSearchParams, ...(options.query || {}) };

        push(url, {
          query: mergedQuery as Record<string, string>,
          hash: options.hash,
        });
      },
      replace: (url: string, options: QueryRouteOptions = {}) => {
        const mergedQuery = options.replace
          ? options.query || {}
          : { ...currentSearchParams, ...(options.query || {}) };

        replace(url, {
          query: mergedQuery as Record<string, string>,
          hash: options.hash,
        });
      },
    }),
    [push, replace, currentSearchParams],
  );
};
