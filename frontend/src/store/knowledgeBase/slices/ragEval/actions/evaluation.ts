import { CreateNewEvalEvaluation, RAGEvalDataSetItem } from '@/types';
import { useEffect } from 'react';
import { StateCreator } from 'zustand/vanilla';

import { ragEvalService } from '@/services/ragEval';
import { KnowledgeBaseStore } from '@/store/knowledgeBase/store';

export interface RAGEvalEvaluationAction {
  checkEvaluationStatus: (id: number) => Promise<void>;
  createNewEvaluation: (params: CreateNewEvalEvaluation) => Promise<void>;
  refreshEvaluationList: () => Promise<void>;
  removeEvaluation: (id: number) => Promise<void>;
  runEvaluation: (id: number) => Promise<void>;
  useFetchEvaluationList: (knowledgeBaseId: string) => void;
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

  useFetchEvaluationList: (knowledgeBaseId) => {
    useEffect(() => {
      if (!knowledgeBaseId) return;

      const fetchEvaluationList = async () => {
        try {
          const data = await ragEvalService.getEvaluationList(knowledgeBaseId);

          if (!get().initDatasetList) {
            set({ initDatasetList: true }, false, 'useFetchDatasets/init');
          }
        } catch (error) {
          console.error('[useFetchEvaluationList] Error:', error);
        }
      };

      fetchEvaluationList();
    }, [knowledgeBaseId]);
  },
});
