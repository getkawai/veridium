import { devtools } from 'zustand/middleware';
import { shallow } from 'zustand/shallow';
import { createWithEqualityFn } from 'zustand/traditional';

export interface RouterState {
  hash: string;
  pathname: string;
  searchParams: Record<string, string>;
}

export interface RouterActions {
  back: () => void;
  forward: () => void;
  push: (pathname: string, options?: NavigationOptions) => void;
  removeSearchParam: (key: string) => void;
  replace: (pathname: string, options?: NavigationOptions) => void;
  setSearchParam: (key: string, value: string) => void;
  setSearchParams: (params: Record<string, string>) => void;
}

export interface NavigationOptions {
  hash?: string;
  query?: Record<string, string>;
  replace?: boolean;
}

export type RouterStore = RouterState & RouterActions;

const initialState: RouterState = {
  pathname: '/chat',
  searchParams: {},
  hash: '',
};

export const useRouterStore = createWithEqualityFn<RouterStore>()(
  devtools(
    (set, get) => ({
      ...initialState,

      push: (pathname: string, options?: NavigationOptions) => {
        set(
          {
            pathname,
            searchParams: options?.query || {},
            hash: options?.hash || '',
          },
          false,
          'router/push',
        );
      },

      replace: (pathname: string, options?: NavigationOptions) => {
        set(
          {
            pathname,
            searchParams: options?.query || {},
            hash: options?.hash || '',
          },
          false,
          'router/replace',
        );
      },

      back: () => {
        // In a pure SPA without URL history, we could implement a history stack
        // For now, this is a no-op
        console.warn('Router.back() is not implemented in SPA mode');
      },

      forward: () => {
        // In a pure SPA without URL history, we could implement a history stack
        // For now, this is a no-op
        console.warn('Router.forward() is not implemented in SPA mode');
      },

      setSearchParam: (key: string, value: string) => {
        set(
          (state) => ({
            searchParams: {
              ...state.searchParams,
              [key]: value,
            },
          }),
          false,
          'router/setSearchParam',
        );
      },

      removeSearchParam: (key: string) => {
        set(
          (state) => {
            const newParams = { ...state.searchParams };
            delete newParams[key];
            return { searchParams: newParams };
          },
          false,
          'router/removeSearchParam',
        );
      },

      setSearchParams: (params: Record<string, string>) => {
        set(
          {
            searchParams: params,
          },
          false,
          'router/setSearchParams',
        );
      },
    }),
    { name: 'RouterStore' },
  ),
  shallow,
);

