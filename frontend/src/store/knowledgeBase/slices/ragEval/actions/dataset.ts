import {
  CreateNewEvalDatasets,
  EvalDatasetRecord,
  RAGEvalDataSetItem,
  insertEvalDatasetRecordSchema,
} from '@/types';
import { StateCreator } from 'zustand/vanilla';

import { notification } from '@/components/AntdStaticMethods';
import { ragEvalService } from '@/services/ragEval';
import { KnowledgeBaseStore } from '@/store/knowledgeBase/store';

export interface RAGEvalDatasetAction {
  createNewDataset: (params: CreateNewEvalDatasets) => Promise<void>;

  importDataset: (file: File, datasetId: number) => Promise<void>;
  refreshDatasetList: () => Promise<void>;
  removeDataset: (id: number) => Promise<void>;
  internal_fetchDatasetRecords: (datasetId: number) => Promise<void>;
  internal_fetchDatasets: (knowledgeBaseId: string) => Promise<void>;
}

export const createRagEvalDatasetSlice: StateCreator<
  KnowledgeBaseStore,
  [['zustand/devtools', never]],
  [],
  RAGEvalDatasetAction
> = (set, get) => ({
  createNewDataset: async (params) => {
    await ragEvalService.createDataset(params);
    await get().refreshDatasetList();
  },

  importDataset: async (file, datasetId) => {
    if (!datasetId) return;
    const fileType = file.name.split('.').pop();

    if (fileType === 'jsonl') {
      // jsonl 文件 需要拆分成单个条，然后逐一校验格式
      const jsonl = await file.text();
      const { default: JSONL } = await import('jsonl-parse-stringify');

      try {
        const items = JSONL.parse(jsonl);

        // check if the items are valid
        insertEvalDatasetRecordSchema.array().parse(items);

        // if valid, send to backend
        await ragEvalService.importDatasetRecords(datasetId, file);
      } catch (e) {
        notification.error({ description: (e as Error).message, message: '文件格式错误' });
      }
    }

    await get().refreshDatasetList();
  },
  refreshDatasetList: async () => {
    const knowledgeBaseId = get().activeId;
    if (knowledgeBaseId) {
      await get().internal_fetchDatasets(knowledgeBaseId);
    }
  },

  removeDataset: async (id) => {
    await ragEvalService.removeDataset(id);
    await get().refreshDatasetList();
  },
  
  internal_fetchDatasetRecords: async (datasetId) => {
    if (!datasetId) return;
    
    try {
      const records = await ragEvalService.getDatasetRecords(datasetId);
      set({ datasetRecords: records }, false, 'internal_fetchDatasetRecords');
    } catch (error) {
      console.error('[internal_fetchDatasetRecords] Error:', error);
    }
  },
  
  internal_fetchDatasets: async (knowledgeBaseId) => {
    try {
      const datasets = await ragEvalService.getDatasets(knowledgeBaseId);
      
      if (!get().initDatasetList) {
        set({ initDatasetList: true, datasets }, false, 'internal_fetchDatasets/init');
      } else {
        set({ datasets }, false, 'internal_fetchDatasets');
      }
    } catch (error) {
      console.error('[internal_fetchDatasets] Error:', error);
      set({ datasets: [] }, false, 'internal_fetchDatasets/error');
    }
  },
});
