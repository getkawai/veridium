import { isEqual } from 'lodash-es';
import { StateCreator } from 'zustand';

import { ListGenerationBatchesWithGenerations } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';
import { Generation, GenerationBatch } from '@/types/generation';
import { setNamespace } from '@/utils/storeDebug';


import { ImageStore } from '../../store';
import { GenerationBatchDispatch, generationBatchReducer } from './reducer';

const n = setNamespace('generationBatch');


export interface GenerationBatchAction {
  setTopicBatchLoaded: (topicId: string) => void;
  internal_dispatchGenerationBatch: (
    topicId: string,
    payload: GenerationBatchDispatch,
    action?: string,
  ) => void;
  removeGeneration: (generationId: string) => Promise<void>;
  internal_deleteGeneration: (generationId: string) => Promise<void>;
  removeGenerationBatch: (batchId: string, topicId: string) => Promise<void>;
  internal_deleteGenerationBatch: (batchId: string, topicId: string) => Promise<void>;
  refreshGenerationBatches: () => Promise<void>;
  internal_fetchGenerationBatches: (topicId: string) => Promise<void>;
  fetchGenerationBatches: (topicId?: string | null) => Promise<GenerationBatch[]>;
}

// ====== action implementation ====== //

export const createGenerationBatchSlice: StateCreator<
  ImageStore,
  [['zustand/devtools', never]],
  [],
  GenerationBatchAction
> = (set, get) => ({
  setTopicBatchLoaded: (topicId: string) => {
    const nextMap = {
      ...get().generationBatchesMap,
      [topicId]: [],
    };

    // no need to update map if the map is the same
    if (isEqual(nextMap, get().generationBatchesMap)) return;

    set(
      {
        generationBatchesMap: nextMap,
      },
      false,
      n('setTopicBatchLoaded'),
    );
  },

  removeGeneration: async (generationId: string) => {
    const { internal_deleteGeneration, activeGenerationTopicId, refreshGenerationBatches } = get();

    await internal_deleteGeneration(generationId);

    // 检查删除后是否有batch变成空的，如果有则删除空batch
    if (activeGenerationTopicId) {
      const updatedBatches = get().generationBatchesMap[activeGenerationTopicId] || [];
      const emptyBatches = updatedBatches.filter((batch) => batch.generations.length === 0);

      // 删除所有空的batch
      for (const emptyBatch of emptyBatches) {
        await get().internal_deleteGenerationBatch(emptyBatch.id, activeGenerationTopicId);
      }

      // 如果删除了空batch，再次刷新数据确保一致性
      if (emptyBatches.length > 0) {
        await refreshGenerationBatches();
      }
    }
  },

  internal_deleteGeneration: async (generationId: string) => {
    const { activeGenerationTopicId, refreshGenerationBatches, internal_dispatchGenerationBatch } =
      get();

    if (!activeGenerationTopicId) return;

    // 找到包含该 generation 的 batch
    const currentBatches = get().generationBatchesMap[activeGenerationTopicId] || [];
    const targetBatch = currentBatches.find((batch) =>
      batch.generations.some((gen) => gen.id === generationId),
    );

    if (!targetBatch) return;

    // 1. 立即更新前端状态（乐观更新）
    internal_dispatchGenerationBatch(
      activeGenerationTopicId,
      { type: 'deleteGenerationInBatch', batchId: targetBatch.id, generationId },
      'internal_deleteGeneration',
    );

    // 2. 调用后端服务删除generation
    // await generationService.deleteGeneration(generationId);

    // 3. 刷新数据确保一致性
    await refreshGenerationBatches();
  },

  removeGenerationBatch: async (batchId: string, topicId: string) => {
    const { internal_deleteGenerationBatch } = get();
    await internal_deleteGenerationBatch(batchId, topicId);
  },

  internal_deleteGenerationBatch: async (batchId: string, topicId: string) => {
    const { internal_dispatchGenerationBatch, refreshGenerationBatches } = get();

    // 1. 立即更新前端状态（乐观更新）
    internal_dispatchGenerationBatch(
      topicId,
      { type: 'deleteBatch', id: batchId },
      'internal_deleteGenerationBatch',
    );

    // 2. 调用后端服务
    // await generationBatchService.deleteGenerationBatch(batchId);

    // 3. 刷新数据确保一致性
    await refreshGenerationBatches();
  },

  internal_dispatchGenerationBatch: (topicId, payload, action) => {
    const currentBatches = get().generationBatchesMap[topicId] || [];
    const nextBatches = generationBatchReducer(currentBatches, payload);

    const nextMap = {
      ...get().generationBatchesMap,
      [topicId]: nextBatches,
    };

    // no need to update map if the map is the same
    if (isEqual(nextMap, get().generationBatchesMap)) return;

    set(
      {
        generationBatchesMap: nextMap,
      },
      false,
      action ?? n(`dispatchGenerationBatch/${payload.type}`),
    );
  },

  refreshGenerationBatches: async () => {
    const { activeGenerationTopicId } = get();
    if (activeGenerationTopicId) {
      await get().internal_fetchGenerationBatches(activeGenerationTopicId);
    }
  },

  internal_fetchGenerationBatches: async (topicId) => {
    if (!topicId) return;

    try {
      const rows = await ListGenerationBatchesWithGenerations(topicId);

      const batchesMap = new Map<string, GenerationBatch>();

      for (const row of rows) {
        if (!batchesMap.has(row.batchId)) {
          batchesMap.set(row.batchId, {
            id: row.batchId,
            generations: [],
            createdAt: new Date(row.batchCreatedAt),
            model: row.model,
            provider: row.provider,
            prompt: row.prompt,
            config: row.config && row.config.Valid ? JSON.parse(row.config.String) : {},
            width: row.width && row.width.Valid ? row.width.Int64 : undefined,
            height: row.height && row.height.Valid ? row.height.Int64 : undefined,
          });
        }

        if (row.genId && row.genId.Valid) {
          const gen: Generation = {
            id: row.genId.String,
            asyncTaskId: row.asyncTaskId && row.asyncTaskId.Valid ? row.asyncTaskId.String : null,
            createdAt: new Date((row.genCreatedAt && row.genCreatedAt.Valid ? row.genCreatedAt.Int64 : row.batchCreatedAt) || Date.now()),
            seed: row.seed && row.seed.Valid ? row.seed.Int64 : null,
            task: {
              id: row.taskId && row.taskId.Valid ? row.taskId.String : '',
              status: (row.taskState && row.taskState.Valid ? row.taskState.String : 'success') as any, // Default to success if no task
              error: row.taskError && row.taskError.Valid ? JSON.parse(row.taskError.String) : undefined,
            },
            asset: row.asset && row.asset.Valid ? JSON.parse(row.asset.String) : null,
          };
          batchesMap.get(row.batchId)!.generations.push(gen);
        }
      }

      const batches = Array.from(batchesMap.values());

      const nextMap = {
        ...get().generationBatchesMap,
        [topicId]: batches,
      };

      // no need to update map if the map is the same
      if (isEqual(nextMap, get().generationBatchesMap)) return;

      set(
        {
          generationBatchesMap: nextMap,
        },
        false,
        n('internal_fetchGenerationBatches(success)', { topicId }),
      );
    } catch (error) {
      console.error('[internal_fetchGenerationBatches] Error:', error);
    }
  },

  fetchGenerationBatches: async (topicId) => {
    if (!topicId) return [];
    await get().internal_fetchGenerationBatches(topicId);
    return get().generationBatchesMap[topicId] || [];
  },
});
