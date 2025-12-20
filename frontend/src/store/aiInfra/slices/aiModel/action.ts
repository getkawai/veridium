import isEqual from 'fast-deep-equal';
import {
  AiModelSortMap,
  AiProviderModelListItem,
  CreateAiModelParams,
  ToggleAiModelEnableParams,
} from '@/model-bank';
import { StateCreator } from 'zustand/vanilla';

// ✅ MIGRATION COMPLETE: All operations now use direct DB calls
import { AIProviderStoreState } from '../../initialState';
import type { AiProviderAction } from '../aiProvider/action';
import { DB } from '@/types/database';
import { getUserId, toNullInt64, boolToInt, toNullString, mapModelFromDB } from '../aiProvider/helpers';

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

    const userId = getUserId();
    const now = Date.now();

    // Batch toggle all models in parallel
    await Promise.all(
      ids.map((id) =>
        DB.ToggleAIModelEnabled({
          id,
          providerId: activeAiProvider,
          enabled: toNullInt64(boolToInt(enabled)),
          type: 'chat',
          source: toNullString('custom'),
          createdAt: now,
          updatedAt: now,
        }),
      ),
    );

    console.log(`[AI Model] Batch toggled ${ids.length} models to ${enabled} via direct DB`);

    await get().refreshAiModelList();
  },
  batchUpdateAiModels: async (models) => {
    const { activeAiProvider: id } = get();
    if (!id) return;

    const userId = getUserId();
    const now = Date.now();

    // Batch update all models in parallel
    await Promise.all(
      models.map(async (model) => {
        // Get current model to merge
        const current = await DB.GetAIModel({ id: model.id, providerId: id });
        if (!current) {
          console.warn(`[AI Model] Model ${model.id} not found, skipping update`);
          return;
        }

        // Merge updates (only fields that exist in UpdateAIModelParams)
        await DB.UpdateAIModel({
          id: model.id,
          providerId: id,
          displayName: model.displayName ? toNullString(model.displayName) : current.displayName,
          description: current.description,
          enabled: model.enabled !== undefined ? toNullInt64(boolToInt(model.enabled)) : current.enabled,
          sort: current.sort, // Keep current sort (use updateAiModelsSort for sort changes)
          pricing: current.pricing,
          parameters: current.parameters,
          config: current.config,
          abilities: model.abilities ? toNullString(JSON.stringify(model.abilities)) : current.abilities,
          updatedAt: now,
        });
      }),
    );

    console.log(`[AI Model] Batch updated ${models.length} models via direct DB`);

    await get().refreshAiModelList();
  },
  clearModelsByProvider: async (provider) => {
    // Delete all models for this provider
    await DB.DeleteAIModelsByProvider(provider);

    console.log(`[AI Model] Cleared all models for provider ${provider} via direct DB`);

    await get().refreshAiModelList();
  },
  clearRemoteModels: async (provider) => {
    // Delete all remote models for this provider
    // Note: We use DeleteAIModelsByProvider and filter by source in the future
    // For now, get all models and delete remote ones
    const models = await DB.ListAIModelsByProvider(provider);
    const remoteModels = models.filter((m: any) => m.source?.String === 'remote');
    
    await Promise.all(
      remoteModels.map((m: any) =>
        DB.DeleteAIModel({
          id: m.id,
          providerId: provider,
        }),
      ),
    );

    console.log(`[AI Model] Cleared remote models for provider ${provider} via direct DB`);

    await get().refreshAiModelList();
  },
  createNewAiModel: async (data) => {
    const userId = getUserId();
    const now = Date.now();

    // Validation
    if (!data.id || !data.displayName || !data.providerId) {
      throw new Error('Model ID, display name, and provider ID are required');
    }

    // Check if already exists
    try {
      const existing = await DB.GetAIModel({ id: data.id, providerId: data.providerId });
      if (existing) {
        throw new Error(`Model ${data.id} already exists`);
      }
    } catch (e: any) {
      // Not found error is OK
      if (!e.message?.includes('not found')) {
        throw e;
      }
    }

    // Create model
    await DB.CreateAIModel({
      id: data.id,
      displayName: toNullString(data.displayName || data.id),
      description: toNullString(''),
      organization: toNullString(''),
      enabled: toNullInt64(1), // New models enabled by default
      providerId: data.providerId,
      type: data.type || 'chat',
      sort: toNullInt64(0),
      pricing: toNullString('{}'),
      parameters: toNullString('{}'),
      config: toNullString('{}'),
      abilities: toNullString(JSON.stringify(data.abilities || {})),
      contextWindowTokens: toNullInt64(data.contextWindowTokens || 0),
      source: toNullString('custom'),
      releasedAt: toNullString(data.releasedAt || ''),
      createdAt: now,
      updatedAt: now,
    });

    console.log(`[AI Model] Created model ${data.id} via direct DB`);

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

      const dbModels = await DB.ListAIModelsByProvider(activeProvider);

      const data = dbModels.map(mapModelFromDB) as any;
      
      if (!isEqual(data, get().aiProviderModelList)) {
        set({ aiProviderModelList: data, isAiModelListInit: true }, false, 'refreshAiModelList/directDB');
      }

      // make refresh provide runtime state async, not block
      get().refreshAiProviderRuntimeState();
    } catch (error) {
      console.error('[refreshAiModelList] Error:', error);
    }
  },
  removeAiModel: async (id, providerId) => {
    await DB.DeleteAIModel({
      id,
      providerId,
    });

    console.log(`[AI Model] Deleted model ${id} via direct DB`);

    await get().refreshAiModelList();
  },
  toggleModelEnabled: async (params) => {
    const { activeAiProvider } = get();
    if (!activeAiProvider) return;

    get().internal_toggleAiModelLoading(params.id, true);

    try {
      const now = Date.now();

      await DB.ToggleAIModelEnabled({
        id: params.id,
        providerId: activeAiProvider,
        enabled: toNullInt64(boolToInt(params.enabled)),
        type: 'chat', // Required field
        source: toNullString('custom'), // Required field
        createdAt: now,
        updatedAt: now,
      });

      console.log(`[AI Model] Toggled ${params.id} to ${params.enabled} via direct DB`);

      await get().refreshAiModelList();
    } finally {
      get().internal_toggleAiModelLoading(params.id, false);
    }
  },

  updateAiModelsConfig: async (id, providerId, data) => {
    const now = Date.now();

    // Get current model to merge
    const current = await DB.GetAIModel({ id, providerId });
    if (!current) {
      throw new Error(`Model ${id} not found`);
    }

    // Update model with merged data
    await DB.UpdateAIModel({
      id,
      providerId,
      displayName: data.displayName ? toNullString(data.displayName) : current.displayName,
      description: current.description,
      enabled: data.enabled !== undefined ? toNullInt64(boolToInt(data.enabled)) : current.enabled,
      sort: current.sort,
      pricing: current.pricing,
      parameters: data.parameters ? toNullString(JSON.stringify(data.parameters)) : current.parameters,
      config: data.config ? toNullString(JSON.stringify(data.config)) : current.config,
      abilities: data.abilities ? toNullString(JSON.stringify(data.abilities)) : current.abilities,
      updatedAt: now,
    });

    console.log(`[AI Model] Updated config for ${id} via direct DB`);

    await get().refreshAiModelList();
  },
  updateAiModelsSort: async (id, items) => {
    const now = Date.now();

    // Batch update all model sorts in parallel
    await Promise.all(
      items.map(({ id: modelId, sort }) =>
        DB.UpdateAIModelSort({
          id: modelId,
          providerId: id,
          sort: toNullInt64(sort),
          type: 'chat', // Required field
          enabled: toNullInt64(1), // Required field
          source: toNullString('custom'), // Required field
          updatedAt: now,
          createdAt: now, // Required field
        }),
      ),
    );

    console.log(`[AI Model] Updated sort order for ${items.length} models via direct DB`);

    await get().refreshAiModelList();
  },

  internal_fetchAiProviderModels: async (id) => {
    if (!id) return;

    try {
      const dbModels = await DB.ListAIModelsByProvider(id);

      const data = dbModels.map(mapModelFromDB) as any;

      // no need to update list if the list have been init and data is the same
      if (get().isAiModelListInit && isEqual(data, get().aiProviderModelList)) return;

      set(
        { aiProviderModelList: data, isAiModelListInit: true },
        false,
        `internal_fetchAiProviderModels/${id}/directDB`,
      );
    } catch (error) {
      console.error('[internal_fetchAiProviderModels] Error:', error);
    }
  },
});
