import { createWithEqualityFn } from 'zustand/traditional';
import { StateCreator } from 'zustand/vanilla';

import { type GlobalGeneralAction, generalActionSlice } from './actions/general';
import { type GlobalWorkspacePaneAction, globalWorkspaceSlice } from './actions/workspacePane';
import { type GlobalState, initialState } from './initialState';

//  ===============  聚合 createStoreFn ============ //

export interface GlobalStore
  extends GlobalState,
    GlobalWorkspacePaneAction,
    GlobalGeneralAction {
  /* empty */
}

const createStore: StateCreator<GlobalStore> = (...parameters) => ({
  ...initialState,
  ...globalWorkspaceSlice(...parameters),
  ...generalActionSlice(...parameters),
});

//  ===============  实装 useStore ============ //

export const useGlobalStore = createWithEqualityFn<GlobalStore>()(createStore);
