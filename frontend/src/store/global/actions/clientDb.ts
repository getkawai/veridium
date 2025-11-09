import { useEffect } from 'react';
import type { StateCreator } from 'zustand/vanilla';

import { initializeDB } from '@/database/client/db';
import type { GlobalStore } from '@/store/global';
import { DatabaseLoadingState, OnStageChange } from '@/types/clientDB';

type InitClientDBParams = { onStateChange: OnStageChange };
/**
 * 设置操作
 */
export interface GlobalClientDBAction {
  initializeClientDB: (params?: InitClientDBParams) => Promise<void>;
  markPgliteEnabled: () => void;
  useInitClientDB: (params?: InitClientDBParams) => void;
}

export const clientDBSlice: StateCreator<
  GlobalStore,
  [['zustand/devtools', never]],
  [],
  GlobalClientDBAction
> = (set, get) => ({
  initializeClientDB: async (params) => {
    // if the db has started initialized or not error, just skip.
    if (
      get().initClientDBStage !== DatabaseLoadingState.Idle &&
      get().initClientDBStage !== DatabaseLoadingState.Error
    )
      return;

    await initializeDB({
      onError: ({ error }) => {
        set({
          initClientDBError: error,
        });
      },
      onProgress: (data) => {
        set({ initClientDBProcess: data });
      },
      onStateChange: (state) => {
        set({ initClientDBStage: state });
        params?.onStateChange?.(state);
      },
    });
  },
  markPgliteEnabled: async () => {
    get().updateSystemStatus({ isEnablePglite: true });

    if (navigator.storage && !!navigator.storage.persist) {
      // 1. Check if persistent permission has been obtained
      const isPersisted = await navigator.storage.persisted();

      // 2. If the persistent permission has not been obtained, request permission
      if (!isPersisted) {
        await navigator.storage.persist();
      }
    }
  },
  useInitClientDB: (params) => {
    useEffect(() => {
      get().initializeClientDB(params);
    }, []); // Run once on mount
  },
});
