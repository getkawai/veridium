import { uniqBy } from 'lodash-es';
import {
  AIImageModelCard,
  EnabledAiModel,
  ModelAbilities,
} from '@/model-bank';
import { StateCreator } from 'zustand/vanilla';
// ✅ MIGRATION COMPLETE: All operations now use direct DB calls
import { DEFAULT_MODEL_PROVIDER_LIST, LOBE_DEFAULT_MODEL_LIST } from '@/config/modelProviders';
import { AIProviderStoreState } from '../../initialState';
import type { AiModelAction } from '../aiModel/action';
import { useUserStore } from '@/store/user';
import { authSelectors } from '@/store/user/selectors';
import {
  AiProviderListItem,
  AiProviderSortMap,
  AiProviderSourceEnum,
  CreateAiProviderParams,
  EnabledProvider,
  UpdateAiProviderConfigParams,
  UpdateAiProviderParams,
} from '@/types/aiProvider';
import { DB } from '@/types/database';
import {
  mapProviderFromDB,
  toNullString,
  toNullInt64,
  boolToInt,
  parseNullJSON,
} from './helpers';
import type { AiProviderRuntimeConfig, ResponseAnimationStyle } from '@/types/aiProvider';

/**
 * Static Kawai AI Provider & Model Configuration
 * These are hardcoded because:
 * 1. Model is handled internally by backend (llama.cpp)
 * 2. User cannot configure or add models
 * 3. Data never changes - always "kawai" provider with "kawai-auto" model
 */
const STATIC_KAWAI_PROVIDER: EnabledProvider = {
  id: 'kawai',
  name: 'Kawai',
  source: AiProviderSourceEnum.Builtin,
};

const STATIC_KAWAI_MODEL: EnabledAiModel = {
  id: 'kawai-auto',
  displayName: 'Kawai Auto',
  providerId: 'kawai',
  type: 'chat',
  abilities: { functionCall: true, vision: true, files: true },
  contextWindowTokens: 128000,
};

const STATIC_KAWAI_RUNTIME_CONFIG: Record<string, AiProviderRuntimeConfig> = {
  kawai: {
    keyVaults: {},
    settings: {
      defaultShowBrowserRequest: true,
      proxyUrl: { placeholder: 'https://node.getkawai.com/v1' },
      responseAnimation: { speed: 2, text: 'smooth' as ResponseAnimationStyle },
      showApiKey: false,
      showModelFetcher: false,
    },
    config: {},
    fetchOnClient: false,
  },
};

/**
 * Get models by provider ID and type, with proper formatting and deduplication
 */
export const getModelListByType = async (
  enabledAiModels: EnabledAiModel[],
  providerId: string,
  type: string,
) => {
  const filteredModels = enabledAiModels.filter(
    (model) => model.providerId === providerId && model.type === type,
  );

  const models = await Promise.all(
    filteredModels.map(async (model) => ({
      abilities: (model.abilities || {}) as ModelAbilities,
      contextWindowTokens: model.contextWindowTokens,
      displayName: model.displayName ?? '',
      id: model.id,
      ...(model.type === 'image' && {
        parameters: (model as AIImageModelCard).parameters || {},
      }),
    })),
  );

  return uniqBy(models, 'id');
};

/**
 * Build provider model lists with proper async handling
 */
const buildProviderModelLists = async (
  providers: EnabledProvider[],
  enabledAiModels: EnabledAiModel[],
  type: 'chat' | 'image',
) => {
  return Promise.all(
    providers.map(async (provider) => ({
      ...provider,
      children: await getModelListByType(enabledAiModels, provider.id, type),
      name: provider.name || provider.id,
    })),
  );
};

// enum AiProviderSwrKey {
//   fetchAiProviderItem = 'FETCH_AI_PROVIDER_ITEM',
//   fetchAiProviderList = 'FETCH_AI_PROVIDER',
//   fetchAiProviderRuntimeState = 'FETCH_AI_PROVIDER_RUNTIME_STATE',
// }

// type AiProviderRuntimeStateWithBuiltinModels = AiProviderRuntimeState & {
//   builtinAiModelList: LobeDefaultAiModelListItem[];
//   enabledChatModelList?: EnabledProviderWithModels[];
//   enabledImageModelList?: EnabledProviderWithModels[];
// };

export interface AiProviderAction {
  createNewAiProvider: (params: CreateAiProviderParams) => Promise<void>;
  deleteAiProvider: (id: string) => Promise<void>;
  internal_toggleAiProviderConfigUpdating: (id: string, loading: boolean) => void;
  internal_toggleAiProviderLoading: (id: string, loading: boolean) => void;
  refreshAiProviderDetail: () => Promise<void>;
  refreshAiProviderList: () => Promise<void>;
  refreshAiProviderRuntimeState: () => Promise<void>;
  removeAiProvider: (id: string) => Promise<void>;
  toggleProviderEnabled: (id: string, enabled: boolean) => Promise<void>;
  updateAiProvider: (id: string, value: UpdateAiProviderParams) => Promise<void>;
  updateAiProviderConfig: (id: string, value: UpdateAiProviderConfigParams) => Promise<void>;
  updateAiProviderSort: (items: AiProviderSortMap[]) => Promise<void>;

  internal_fetchAiProviderItem: (id: string) => Promise<void>;
  internal_fetchAiProviderList: (opts?: { enabled?: boolean }) => Promise<void>;
  internal_fetchAiProviderRuntimeState: (isLogin: boolean | null | undefined) => Promise<void>;
}

export const createAiProviderSlice: StateCreator<
  AIProviderStoreState & AiProviderAction & AiModelAction,
  [['zustand/devtools', never]],
  [],
  AiProviderAction
> = (set, get) => ({
  createNewAiProvider: async (params) => {
    const now = Date.now();

    // Validation
    if (!params.id || !params.name) {
      throw new Error('Provider ID and name are required');
    }

    // Check if already exists
    try {
      const existing = await DB.GetAIProvider(params.id);
      if (existing) {
        throw new Error(`Provider ${params.id} already exists`);
      }
    } catch (e: any) {
      // Not found error is OK, means we can create
      if (!e.message?.includes('not found')) {
        throw e;
      }
    }

    // Create provider
    await DB.CreateAIProvider({
      id: params.id,
      name: toNullString(params.name),
      sort: toNullInt64(0), // Default sort
      enabled: toNullInt64(1), // New providers enabled by default
      fetchOnClient: toNullInt64(0), // Default false
      checkModel: toNullString(''),
      logo: toNullString(params.logo || ''),
      description: toNullString(params.description || ''),
      keyVaults: toNullString(JSON.stringify(params.keyVaults || {})),
      source: toNullString('custom'),
      settings: toNullString(JSON.stringify(params.settings || {})),
      config: toNullString(JSON.stringify(params.config || {})),
      createdAt: now,
      updatedAt: now,
    });

    console.log(`[AI Provider] Created provider ${params.id} via direct DB`);

    await get().refreshAiProviderList();
  },
  deleteAiProvider: async (id: string) => {
    // Delete provider (backend handles cascade delete of models)
    await DB.DeleteAIProvider(id);

    console.log(`[AI Provider] Deleted provider ${id} via direct DB`);

    await get().refreshAiProviderList();
  },
  internal_toggleAiProviderConfigUpdating: (id, loading) => {
    set(
      (state) => {
        if (loading)
          return { aiProviderConfigUpdatingIds: [...state.aiProviderConfigUpdatingIds, id] };

        return {
          aiProviderConfigUpdatingIds: state.aiProviderConfigUpdatingIds.filter((i) => i !== id),
        };
      },
      false,
      'toggleAiProviderLoading',
    );
  },
  internal_toggleAiProviderLoading: (id, loading) => {
    set(
      (state) => {
        if (loading) return { aiProviderLoadingIds: [...state.aiProviderLoadingIds, id] };

        return { aiProviderLoadingIds: state.aiProviderLoadingIds.filter((i) => i !== id) };
      },
      false,
      'toggleAiProviderLoading',
    );
  },
  refreshAiProviderDetail: async () => {
    try {
      const activeProvider = get().activeAiProvider;
      if (!activeProvider) return;

      const dbProvider = await DB.GetAIProviderDetail(activeProvider);

      if (dbProvider) {
        const data = mapProviderFromDB(dbProvider);
        set({ aiProviderDetail: data as any }, false, 'refreshAiProviderDetail/directDB');
      }
      await get().refreshAiProviderRuntimeState();
    } catch (error) {
      console.error('[refreshAiProviderDetail] Error:', error);
    }
  },
  refreshAiProviderList: async () => {
    try {
      // Just call the already-migrated internal_fetchAiProviderList
      await get().internal_fetchAiProviderList();
      await get().refreshAiProviderRuntimeState();
    } catch (error) {
      console.error('[refreshAiProviderList] Error:', error);
    }
  },
  refreshAiProviderRuntimeState: async () => {
    // Runtime state refresh is handled by useFetchAiProviderRuntimeState
    // This is a no-op now as we don't use SWR cache invalidation
    console.debug('[refreshAiProviderRuntimeState] Skipped (handled by useEffect)');
  },
  removeAiProvider: async (id) => {
    // Just call the already-migrated deleteAiProvider
    await get().deleteAiProvider(id);
  },

  toggleProviderEnabled: async (id: string, enabled: boolean) => {
    get().internal_toggleAiProviderLoading(id, true);

    try {
      const now = Date.now();

      await DB.ToggleAIProviderEnabled({
        id,
        enabled: toNullInt64(boolToInt(enabled)),
        source: toNullString('custom'),
        createdAt: now,
        updatedAt: now,
      });

      console.log(`[AI Provider] Toggled ${id} to ${enabled} via direct DB`);

      await get().refreshAiProviderList();
    } finally {
      get().internal_toggleAiProviderLoading(id, false);
    }
  },

  updateAiProvider: async (id, value) => {
    get().internal_toggleAiProviderLoading(id, true);

    try {
      const now = Date.now();

      // Get current provider to merge with updates
      const current = await DB.GetAIProvider(id);
      if (!current) {
        throw new Error(`Provider ${id} not found`);
      }

      // Merge updates with current values
      await DB.UpdateAIProvider({
        id,
        name: value.name ? toNullString(value.name) : current.name,
        sort: current.sort, // Keep current sort (use updateAiProviderSort for sort changes)
        enabled: current.enabled, // Keep current enabled (use toggleProviderEnabled for enable changes)
        fetchOnClient: current.fetchOnClient,
        checkModel: current.checkModel,
        logo: value.logo ? toNullString(value.logo) : current.logo,
        description: value.description ? toNullString(value.description) : current.description,
        keyVaults: current.keyVaults, // Keep current keyVaults (use updateAiProviderConfig for config changes)
        settings: value.settings ? toNullString(JSON.stringify(value.settings)) : current.settings,
        config: value.config ? toNullString(JSON.stringify(value.config)) : current.config,
        updatedAt: now,
      });

      console.log(`[AI Provider] Updated provider ${id} via direct DB`);

      await get().refreshAiProviderList();
      await get().refreshAiProviderDetail();
    } finally {
      get().internal_toggleAiProviderLoading(id, false);
    }
  },

  updateAiProviderConfig: async (id, value) => {
    get().internal_toggleAiProviderConfigUpdating(id, true);

    try {
      const now = Date.now();

      // Get current provider
      const current = await DB.GetAIProvider(id);
      if (!current) {
        throw new Error(`Provider ${id} not found`);
      }

      // Parse current config and merge with updates
      const currentConfig = parseNullJSON(current.config, {});
      const mergedConfig = { ...currentConfig, ...(value.config || {}) };

      const currentSettings = parseNullJSON(current.settings, {});
      // Note: value.settings is not in UpdateAiProviderConfigParams, so we keep current settings
      const mergedSettings = currentSettings;

      const currentKeyVaults = parseNullJSON(current.keyVaults, {});
      const mergedKeyVaults = { ...currentKeyVaults, ...(value.keyVaults || {}) };

      // Update only config-related fields
      await DB.UpdateAIProvider({
        id,
        config: toNullString(JSON.stringify(mergedConfig)),
        settings: toNullString(JSON.stringify(mergedSettings)),
        keyVaults: toNullString(JSON.stringify(mergedKeyVaults)),
        updatedAt: now,
        // Keep other fields unchanged
        name: current.name,
        sort: current.sort,
        enabled: current.enabled,
        fetchOnClient: current.fetchOnClient,
        checkModel: current.checkModel,
        logo: current.logo,
        description: current.description,
      });

      console.log(`[AI Provider] Updated config for ${id} via direct DB`);

      await get().refreshAiProviderDetail();
    } finally {
      get().internal_toggleAiProviderConfigUpdating(id, false);
    }
  },

  updateAiProviderSort: async (items) => {
    const now = Date.now();

    // Batch update all sorts in parallel
    await Promise.all(
      items.map(({ id, sort }) =>
        DB.UpdateAIProvider({
          id,
          sort: toNullInt64(sort),
          updatedAt: now,
          // Required fields (empty = no change)
          name: toNullString(''),
          enabled: toNullInt64(0),
          fetchOnClient: toNullInt64(0),
          checkModel: toNullString(''),
          logo: toNullString(''),
          description: toNullString(''),
          keyVaults: toNullString(''),
          settings: toNullString(''),
          config: toNullString(''),
        }),
      ),
    );

    console.log(`[AI Provider] Updated sort order for ${items.length} providers via direct DB`);

    await get().refreshAiProviderList();
  },
  internal_fetchAiProviderItem: async (id) => {
    if (!id) return;

    try {
      const dbProvider = await DB.GetAIProviderDetail(id);

      if (!dbProvider) return;

      const data = mapProviderFromDB(dbProvider);
      set({ activeAiProvider: id, aiProviderDetail: data as any }, false, 'internal_fetchAiProviderItem/directDB');
    } catch (error) {
      console.error('[internal_fetchAiProviderItem] Error:', error);
    }
  },

  internal_fetchAiProviderList: async (opts) => {
    if (opts?.enabled === false) return;

    try {
      const dbProviders = opts?.enabled
        ? await DB.ListEnabledAIProviders()
        : await DB.ListAIProviders();

      const data: AiProviderListItem[] = dbProviders.map((p) => {
        const mapped = mapProviderFromDB(p);
        return {
          id: mapped.id,
          name: mapped.name,
          enabled: mapped.enabled,
          sort: mapped.sort,
          source: mapped.source as any,
          logo: mapped.logo,
          description: mapped.description,
        };
      });

      if (!get().initAiProviderList) {
        set(
          { aiProviderList: data, initAiProviderList: true },
          false,
          'internal_fetchAiProviderList/init/directDB',
        );
        return;
      }

      set({ aiProviderList: data }, false, 'internal_fetchAiProviderList/refresh/directDB');
    } catch (error) {
      console.error('[internal_fetchAiProviderList] Error:', error);
    }
  },

  internal_fetchAiProviderRuntimeState: async (isLogin) => {
    const isAuthLoaded = authSelectors.isLoaded(useUserStore.getState());
    // For desktop app, allow fetch even if auth is not fully loaded
    const shouldFetch = isLogin !== null && isLogin !== undefined && (isAuthLoaded !== false);

    if (!shouldFetch) return;

    try {
      const builtinAiModelList = LOBE_DEFAULT_MODEL_LIST;

      if (isLogin) {
        // Use static Kawai provider & model - no DB calls needed
        // Model is handled internally by backend (llama.cpp), user cannot configure
        console.log('[AI Provider] Using static Kawai configuration (no DB call)');

        const enabledAiProviders: EnabledProvider[] = [STATIC_KAWAI_PROVIDER];
        const enabledAiModels: EnabledAiModel[] = [STATIC_KAWAI_MODEL];
        const aiProviderRuntimeConfig = STATIC_KAWAI_RUNTIME_CONFIG;

        // Kawai only has chat models, no image models
        const enabledChatAiProviders = [STATIC_KAWAI_PROVIDER];
        const enabledImageAiProviders: EnabledProvider[] = [];

        // Build model lists
        const [enabledChatModelList, enabledImageModelList] = await Promise.all([
          buildProviderModelLists(enabledChatAiProviders, enabledAiModels, 'chat'),
          buildProviderModelLists(enabledImageAiProviders, enabledAiModels, 'image'),
        ]);

        set(
          {
            aiProviderRuntimeConfig,
            builtinAiModelList,
            enabledAiModels,
            enabledAiProviders,
            enabledChatModelList,
            enabledImageModelList,
            isInitAiProviderRuntimeState: true,
          },
          false,
          'internal_fetchAiProviderRuntimeState/login/static',
        );
      } else {
        // No login: Use builtin models only
        const enabledAiProviders: EnabledProvider[] = DEFAULT_MODEL_PROVIDER_LIST.filter(
          (provider) => provider.enabled,
        ).map((item) => ({
          id: item.id,
          name: item.name,
          source: AiProviderSourceEnum.Builtin,
        }));

        const enabledChatAiProviders = enabledAiProviders.filter((provider) => {
          return builtinAiModelList.some(
            (model) => model.providerId === provider.id && model.type === 'chat',
          );
        });

        const enabledImageAiProviders = enabledAiProviders.filter((provider) => {
          return builtinAiModelList.some(
            (model) => model.providerId === provider.id && model.type === 'image',
          );
        });

        const enabledAiModels = builtinAiModelList.filter((m) => m.enabled);
        const [enabledChatModelList, enabledImageModelList] = await Promise.all([
          buildProviderModelLists(enabledChatAiProviders, enabledAiModels, 'chat'),
          buildProviderModelLists(enabledImageAiProviders, enabledAiModels, 'image'),
        ]);

        set(
          {
            aiProviderRuntimeConfig: {},
            builtinAiModelList,
            enabledAiModels,
            enabledAiProviders,
            enabledChatModelList,
            enabledImageModelList,
            isInitAiProviderRuntimeState: true,
          },
          false,
          'internal_fetchAiProviderRuntimeState/noLogin',
        );
      }
    } catch (error) {
      console.error('[internal_fetchAiProviderRuntimeState] Error:', error);
    }
  },
});
