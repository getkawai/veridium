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
  internal_fetchKnowledgeBaseList: (params: { suspense?: boolean }) => Promise<void>;
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
    await get().internal_fetchKnowledgeBaseList();
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

  internal_fetchKnowledgeBaseList: async (params = {}) => {
    try {
      const list = await knowledgeBaseService.getKnowledgeBaseList();
      
      if (!get().initKnowledgeBaseList) {
        set({ initKnowledgeBaseList: true, knowledgeBaseList: list }, false, 'internal_fetchKnowledgeBaseList/init');
      } else {
        set({ knowledgeBaseList: list }, false, 'internal_fetchKnowledgeBaseList');
      }
    } catch (error) {
      console.error('[internal_fetchKnowledgeBaseList] Error:', error);
      set({ knowledgeBaseList: [] }, false, 'internal_fetchKnowledgeBaseList/error');
    }
  },
});
