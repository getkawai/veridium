import { StateCreator } from 'zustand/vanilla';

import { knowledgeBaseService } from '@/services/knowledgeBase';
import { KnowledgeBaseStore } from '@/store/knowledgeBase/store';
import { CreateKnowledgeBaseParams, KnowledgeBaseItem } from '@/types/knowledgeBase';

export interface KnowledgeBaseCrudAction {
  createNewKnowledgeBase: (params: CreateKnowledgeBaseParams) => Promise<string>;
  internal_toggleKnowledgeBaseLoading: (id: string, loading: boolean) => void;
  refreshKnowledgeBaseList: () => Promise<void>;

  removeKnowledgeBase: (id: string) => Promise<void>;
  updateKnowledgeBase: (id: string, value: CreateKnowledgeBaseParams) => Promise<void>;

  internal_fetchKnowledgeBaseItem: (id: string) => Promise<void>;
  fetchKnowledgeBaseList: () => Promise<void>;
}

export const createCrudSlice: StateCreator<
  KnowledgeBaseStore,
  [['zustand/devtools', never]],
  [],
  KnowledgeBaseCrudAction
> = (set, get) => ({
  createNewKnowledgeBase: async (params) => {
    const id = await knowledgeBaseService.createKnowledgeBase(params);

    await get().refreshKnowledgeBaseList();

    return id;
  },
  internal_toggleKnowledgeBaseLoading: (id, loading) => {
    set(
      (state) => {
        if (loading) return { knowledgeBaseLoadingIds: [...state.knowledgeBaseLoadingIds, id] };

        return { knowledgeBaseLoadingIds: state.knowledgeBaseLoadingIds.filter((i) => i !== id) };
      },
      false,
      'toggleKnowledgeBaseLoading',
    );
  },
  refreshKnowledgeBaseList: async () => {
    await get().fetchKnowledgeBaseList();
  },
  removeKnowledgeBase: async (id) => {
    await knowledgeBaseService.deleteKnowledgeBase(id);
    await get().refreshKnowledgeBaseList();
  },
  updateKnowledgeBase: async (id, value) => {
    get().internal_toggleKnowledgeBaseLoading(id, true);
    await knowledgeBaseService.updateKnowledgeBaseList(id, value);
    await get().refreshKnowledgeBaseList();

    get().internal_toggleKnowledgeBaseLoading(id, false);
  },

  internal_fetchKnowledgeBaseItem: async (id) => {
    try {
      const item = await knowledgeBaseService.getKnowledgeBaseById(id);

      if (item) {
        set({
          activeKnowledgeBaseId: id,
          activeKnowledgeBaseItems: {
            ...get().activeKnowledgeBaseItems,
            [id]: item,
          },
        }, false, 'internal_fetchKnowledgeBaseItem');
      }
    } catch (error) {
      console.error('[internal_fetchKnowledgeBaseItem] Error:', error);
    }
  },

  fetchKnowledgeBaseList: async () => {
    try {
      set({ isFetchingList: true });
      const list = await knowledgeBaseService.getKnowledgeBaseList();

      if (!get().initKnowledgeBaseList) {
        set({ initKnowledgeBaseList: true, isFetchingList: false, knowledgeBaseList: list }, false, 'fetchKnowledgeBaseList/init');
      } else {
        set({ isFetchingList: false, knowledgeBaseList: list }, false, 'fetchKnowledgeBaseList');
      }
    } catch (error) {
      console.error('[fetchKnowledgeBaseList] Error:', error);
      set({ isFetchingList: false, knowledgeBaseList: [] }, false, 'fetchKnowledgeBaseList/error');
    }
  },
});
