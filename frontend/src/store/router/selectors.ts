import type { RouterStore } from './store';

export const routerSelectors = {
  pathname: (s: RouterStore) => s.pathname,
  searchParams: (s: RouterStore) => s.searchParams,
  hash: (s: RouterStore) => s.hash,

  // Get a specific search param
  getSearchParam: (key: string) => (s: RouterStore) => s.searchParams[key],

  // Get pathname segments
  pathSegments: (s: RouterStore) => s.pathname.split('/').filter(Boolean),

  // Check if current path matches
  isPath: (path: string) => (s: RouterStore) => s.pathname === path,

  // Check if current path starts with
  isPathPrefix: (prefix: string) => (s: RouterStore) => s.pathname.startsWith(prefix),

  // Get the first segment (active tab)
  activeTab: (s: RouterStore) => s.pathname.split('/').filter(Boolean)[0] || 'chat',
};

