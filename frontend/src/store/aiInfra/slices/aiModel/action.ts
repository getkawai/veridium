import isEqual from 'fast-deep-equal';
import {
  AiModelSortMap,
  AiProviderModelListItem,
  CreateAiModelParams,
  ToggleAiModelEnableParams,
} from '@/model-bank';
import { StateCreator } from 'zustand/vanilla';

import { aiModelService } from '@/services/aiModel';
import { AIProviderStoreState } from '../../initialState';
import type { AiProviderAction } from '../aiProvider/action';

export interface AiModelAction {
  batchToggleAiModels: (ids: string[], enabled: boolean) => Promise<void>;
  batchUpdateAiModels: (models: AiProviderModelListItem[]) => Promise<void>;
  clearModelsByProvider: (provider: string) => Promise<void>;
  clearRemoteModels: (provider: string) => Promise<void>;
  createNewAiModel: (params: CreateAiModelParams) => Promise<void>;
  fetchRemoteModelList: (providerId: string) => Promise<void>;
  internal_toggleAiModelLoading: (id: string, loading: boolean) => void;

  refreshAiModelList: () => Promise<void>;
  removeAiModel: (id: string, providerId: string) => Promise<void>;
  toggleModelEnabled: (params: Omit<ToggleAiModelEnableParams, 'providerId'>) => Promise<void>;
  updateAiModelsConfig: (
    id: string,
    providerId: string,
    data: Partial<AiProviderModelListItem>,
  ) => Promise<void>;
  updateAiModelsSort: (providerId: string, items: AiModelSortMap[]) => Promise<void>;

  internal_fetchAiProviderModels: (id: string) => Promise<void>;
}

export const createAiModelSlice: StateCreator<
  AIProviderStoreState & AiProviderAction & AiModelAction,
  [['zustand/devtools', never]],
  [],
  AiModelAction
> = (set, get) => ({
  batchToggleAiModels: async (ids, enabled) => {
    const { activeAiProvider } = get();
    if (!activeAiProvider) return;

    await aiModelService.batchToggleAiModels(activeAiProvider, ids, enabled);
    await get().refreshAiModelList();
  },
  batchUpdateAiModels: async (models) => {
    const { activeAiProvider: id } = get();
    if (!id) return;

    await aiModelService.batchUpdateAiModels(id, models);
    await get().refreshAiModelList();
  },
  clearModelsByProvider: async (provider) => {
    await aiModelService.clearModelsByProvider(provider);
    await get().refreshAiModelList();
  },
  clearRemoteModels: async (provider) => {
    await aiModelService.clearRemoteModels(provider);
    await get().refreshAiModelList();
  },
  createNewAiModel: async (data) => {
    await aiModelService.createAiModel(data);
    await get().refreshAiModelList();
  },
  fetchRemoteModelList: async (providerId) => {
    const { modelsService } = await import('@/services/models');

    const data = await modelsService.getModels(providerId);
    if (data) {
      await get().batchUpdateAiModels(
        data.map((model) => ({
          ...model,
          abilities: {
            files: model.files,
            functionCall: model.functionCall,
            imageOutput: model.imageOutput,
            reasoning: model.reasoning,
            search: model.search,
            video: model.video,
            vision: model.vision,
          },
          enabled: model.enabled || false,
          source: 'remote',
          type: model.type || 'chat',
        })),
      );

      await get().refreshAiModelList();
    }
  },
  internal_toggleAiModelLoading: (id, loading) => {
    set(
      (state) => {
        if (loading) return { aiModelLoadingIds: [...state.aiModelLoadingIds, id] };

        return { aiModelLoadingIds: state.aiModelLoadingIds.filter((i) => i !== id) };
      },
      false,
      'toggleAiModelLoading',
    );
  },
  refreshAiModelList: async () => {
    try {
      const activeProvider = get().activeAiProvider;
      if (!activeProvider) return;

      const data = await aiModelService.getAiProviderModelList(activeProvider);
      
      if (!isEqual(data, get().aiProviderModelList)) {
        set({ aiProviderModelList: data, isAiModelListInit: true }, false, 'refreshAiModelList');
      }

      // make refresh provide runtime state async, not block
      get().refreshAiProviderRuntimeState();
    } catch (error) {
      console.error('[refreshAiModelList] Error:', error);
    }
  },
  removeAiModel: async (id, providerId) => {
    await aiModelService.deleteAiModel({ id, providerId });
    await get().refreshAiModelList();
  },
  toggleModelEnabled: async (params) => {
    const { activeAiProvider } = get();
    if (!activeAiProvider) return;

    get().internal_toggleAiModelLoading(params.id, true);

    await aiModelService.toggleModelEnabled({ ...params, providerId: activeAiProvider });
    await get().refreshAiModelList();

    get().internal_toggleAiModelLoading(params.id, false);
  },

  updateAiModelsConfig: async (id, providerId, data) => {
    await aiModelService.updateAiModel(id, providerId, data);
    await get().refreshAiModelList();
  },
  updateAiModelsSort: async (id, items) => {
    await aiModelService.updateAiModelOrder(id, items);
    await get().refreshAiModelList();
  },

  internal_fetchAiProviderModels: async (id) => {
    if (!id) return;

    try {
      const data = await aiModelService.getAiProviderModelList(id);

      // no need to update list if the list have been init and data is the same
      if (get().isAiModelListInit && isEqual(data, get().aiProviderModelList)) return;

      set(
        { aiProviderModelList: data, isAiModelListInit: true },
        false,
        `internal_fetchAiProviderModels/${id}`,
      );
    } catch (error) {
      console.error('[internal_fetchAiProviderModels] Error:', error);
    }
  },
});
