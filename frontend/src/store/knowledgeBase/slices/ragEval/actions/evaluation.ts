import { CreateNewEvalEvaluation, RAGEvalDataSetItem } from '@/types';
import { StateCreator } from 'zustand/vanilla';

import { ragEvalService } from '@/services/ragEval';
import { KnowledgeBaseStore } from '@/store/knowledgeBase/store';

export interface RAGEvalEvaluationAction {
  checkEvaluationStatus: (id: number) => Promise<void>;
  createNewEvaluation: (params: CreateNewEvalEvaluation) => Promise<void>;
  refreshEvaluationList: () => Promise<void>;
  removeEvaluation: (id: number) => Promise<void>;
  runEvaluation: (id: number) => Promise<void>;
  internal_fetchEvaluationList: (knowledgeBaseId: string) => Promise<void>;
}

export const createRagEvalEvaluationSlice: StateCreator<
  KnowledgeBaseStore,
  [['zustand/devtools', never]],
  [],
  RAGEvalEvaluationAction
> = (set, get) => ({
  checkEvaluationStatus: async (id) => {
    await ragEvalService.checkEvaluationStatus(id);
  },

  createNewEvaluation: async (params) => {
    await ragEvalService.createEvaluation(params);
    await get().refreshEvaluationList();
  },
  refreshEvaluationList: async () => {
    // No-op: handled by useFetchEvaluationList
    console.debug('[refreshEvaluationList] Skipped (handled by useEffect)');
  },

  removeEvaluation: async (id) => {
    await ragEvalService.removeEvaluation(id);
  },

  runEvaluation: async (id) => {
    await ragEvalService.startEvaluationTask(id);
  },

  internal_fetchEvaluationList: async (knowledgeBaseId) => {
    if (!knowledgeBaseId) return;

    try {
      const data = await ragEvalService.getEvaluationList(knowledgeBaseId);

      if (!get().initDatasetList) {
        set({ initDatasetList: true }, false, 'internal_fetchEvaluationList/init');
      }
    } catch (error) {
      console.error('[internal_fetchEvaluationList] Error:', error);
    }
  },
});
