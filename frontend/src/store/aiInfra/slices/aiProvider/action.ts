import { isDeprecatedEdition, isDesktop, isUsePgliteDB } from '@/const';
import { getModelPropertyWithFallback } from '@/model-runtime';
import { uniqBy } from 'lodash-es';
import {
  AIImageModelCard,
  EnabledAiModel,
  LobeDefaultAiModelListItem,
  ModelAbilities,
} from '@/model-bank';
import { StateCreator } from 'zustand/vanilla';
import { aiProviderService } from '@/services/aiProvider';
import { DEFAULT_MODEL_PROVIDER_LIST } from '@/config/modelProviders';
import { AIProviderStoreState } from '../../initialState';
import type { AiModelAction } from '../aiModel/action';
import { useUserStore } from '@/store/user';
import { authSelectors } from '@/store/user/selectors';
import {
  AiProviderDetailItem,
  AiProviderListItem,
  AiProviderRuntimeState,
  AiProviderSortMap,
  AiProviderSourceEnum,
  CreateAiProviderParams,
  EnabledProvider,
  EnabledProviderWithModels,
  UpdateAiProviderConfigParams,
  UpdateAiProviderParams,
} from '@/types/aiProvider';
import { DB } from '@/types/database';
import {
  getUserId,
  mapProviderFromDB,
  mapModelFromDB,
  mapRuntimeConfigFromDB,
  toNullString,
  toNullInt64,
  boolToInt,
} from './helpers';

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
        parameters:
          (model as AIImageModelCard).parameters ||
          (await getModelPropertyWithFallback(model.id, 'parameters')),
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

enum AiProviderSwrKey {
  fetchAiProviderItem = 'FETCH_AI_PROVIDER_ITEM',
  fetchAiProviderList = 'FETCH_AI_PROVIDER',
  fetchAiProviderRuntimeState = 'FETCH_AI_PROVIDER_RUNTIME_STATE',
}

type AiProviderRuntimeStateWithBuiltinModels = AiProviderRuntimeState & {
  builtinAiModelList: LobeDefaultAiModelListItem[];
  enabledChatModelList?: EnabledProviderWithModels[];
  enabledImageModelList?: EnabledProviderWithModels[];
};

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
    await aiProviderService.createAiProvider({ ...params, source: AiProviderSourceEnum.Custom });
    await get().refreshAiProviderList();
  },
  deleteAiProvider: async (id: string) => {
    await aiProviderService.deleteAiProvider(id);

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

      const data = await aiProviderService.getAiProviderById(activeProvider);
      if (data) {
        set({ aiProviderDetail: data }, false, 'refreshAiProviderDetail');
      }
      await get().refreshAiProviderRuntimeState();
    } catch (error) {
      console.error('[refreshAiProviderDetail] Error:', error);
    }
  },
  refreshAiProviderList: async () => {
    try {
      const data = await aiProviderService.getAiProviderList();
      set({ aiProviderList: data }, false, 'refreshAiProviderList');
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
    await aiProviderService.deleteAiProvider(id);
    await get().refreshAiProviderList();
  },

  toggleProviderEnabled: async (id: string, enabled: boolean) => {
    get().internal_toggleAiProviderLoading(id, true);

    // 🚀 PHASE 2 MIGRATION: Feature flag for rollback
    const USE_DIRECT_DB_CALLS = true;

    try {
      if (USE_DIRECT_DB_CALLS) {
        // ✅ NEW: Direct DB call (Phase 2 - Simple Write)
        const userId = getUserId();
        const now = Date.now();

        await DB.ToggleAIProviderEnabled({
          id,
          userId,
          enabled: toNullInt64(boolToInt(enabled)),
          source: toNullString('custom'),
          createdAt: now,
          updatedAt: now,
        });

        console.log(`[AI Provider] Toggled ${id} to ${enabled} via direct DB`);
      } else {
        // ⏳ OLD: Service layer (Phase 2 - Fallback)
        await aiProviderService.toggleProviderEnabled(id, enabled);
      }

      await get().refreshAiProviderList();
    } finally {
      get().internal_toggleAiProviderLoading(id, false);
    }
  },

  updateAiProvider: async (id, value) => {
    get().internal_toggleAiProviderLoading(id, true);
    await aiProviderService.updateAiProvider(id, value);
    await get().refreshAiProviderList();
    await get().refreshAiProviderDetail();

    get().internal_toggleAiProviderLoading(id, false);
  },

  updateAiProviderConfig: async (id, value) => {
    get().internal_toggleAiProviderConfigUpdating(id, true);
    await aiProviderService.updateAiProviderConfig(id, value);
    await get().refreshAiProviderDetail();

    get().internal_toggleAiProviderConfigUpdating(id, false);
  },

  updateAiProviderSort: async (items) => {
    // 🚀 PHASE 2 MIGRATION: Feature flag for rollback
    const USE_DIRECT_DB_CALLS = true;

    if (USE_DIRECT_DB_CALLS) {
      // ✅ NEW: Direct DB call (Phase 2 - Batch Write)
      const userId = getUserId();
      const now = Date.now();

      // Batch update all sorts in parallel
      await Promise.all(
        items.map(({ id, sort }) =>
          DB.UpdateAIProvider({
            id,
            userId,
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
    } else {
      // ⏳ OLD: Service layer (Phase 2 - Fallback)
      await aiProviderService.updateAiProviderOrder(items);
    }

    await get().refreshAiProviderList();
  },
  internal_fetchAiProviderItem: async (id) => {
    if (!id) return;

    try {
      const data = await aiProviderService.getAiProviderById(id);
      if (!data) return;

      set({ activeAiProvider: id, aiProviderDetail: data }, false, 'internal_fetchAiProviderItem');
    } catch (error) {
      console.error('[internal_fetchAiProviderItem] Error:', error);
    }
  },

  internal_fetchAiProviderList: async (opts) => {
    if (opts?.enabled === false) return;

    // 🚀 PHASE 1 MIGRATION: Feature flag for rollback
    const USE_DIRECT_DB_CALLS = true;

    try {
      if (USE_DIRECT_DB_CALLS) {
        // ✅ NEW: Direct DB call (Phase 1 - Optimized)
        const userId = getUserId();
        const dbProviders = opts?.enabled
          ? await DB.ListEnabledAIProviders(userId)
          : await DB.ListAIProviders(userId);

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
      } else {
        // ⏳ OLD: Service layer (Phase 1 - Fallback)
        const data = await aiProviderService.getAiProviderList();

        if (!get().initAiProviderList) {
          set(
            { aiProviderList: data, initAiProviderList: true },
            false,
            'internal_fetchAiProviderList/init/service',
          );
          return;
        }

        set({ aiProviderList: data }, false, 'internal_fetchAiProviderList/refresh/service');
      }
    } catch (error) {
      console.error('[internal_fetchAiProviderList] Error:', error);
    }
  },

  internal_fetchAiProviderRuntimeState: async (isLogin) => {
    const isAuthLoaded = authSelectors.isLoaded(useUserStore.getState());
    const shouldFetch =
      isAuthLoaded && !isDeprecatedEdition && isLogin !== null && isLogin !== undefined;

    if (!shouldFetch) return;

    // 🚀 PHASE 1 MIGRATION: Feature flag for rollback
    const USE_DIRECT_DB_CALLS = true;

    try {
      const [{ LOBE_DEFAULT_MODEL_LIST: builtinAiModelList }] = await Promise.all([
        import('@/model-bank'),
      ]);

      if (isLogin) {
        if (USE_DIRECT_DB_CALLS) {
          // ✅ NEW: Direct DB calls (Phase 1 - Optimized)
          const startTime = performance.now();
          const userId = getUserId();

          // Parallel fetch from database
          const [dbProviders, dbModels, dbConfigs] = await Promise.all([
            DB.ListEnabledAIProviders(userId),
            DB.ListEnabledAIModels(userId),
            DB.GetAIProviderRuntimeConfigs(userId),
          ]);

          const loadTime = performance.now() - startTime;
          console.log(`[AI Provider] Direct DB load completed in ${loadTime.toFixed(2)}ms`);

          // Transform DB results
          const enabledAiProviders: EnabledProvider[] = dbProviders.map((p) => ({
            id: p.id,
            name: mapProviderFromDB(p).name,
            source: mapProviderFromDB(p).source as any,
          }));

          const enabledAiModels: EnabledAiModel[] = dbModels.map(mapModelFromDB) as any;

          // Build runtime config
          const aiProviderRuntimeConfig: Record<string, any> = {};
          dbConfigs.forEach((c) => {
            const mapped = mapRuntimeConfigFromDB(c);
            aiProviderRuntimeConfig[mapped.id] = {
              keyVaults: mapped.keyVaults,
              settings: mapped.settings,
              config: mapped.config,
              fetchOnClient: mapped.fetchOnClient,
            };
          });

          // Filter providers by model type
          const enabledChatAiProviders = enabledAiProviders.filter((provider) =>
            enabledAiModels.some((m) => m.providerId === provider.id && m.type === 'chat'),
          );

          const enabledImageAiProviders = enabledAiProviders.filter((provider) =>
            enabledAiModels.some((m) => m.providerId === provider.id && m.type === 'image'),
          );

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
            'internal_fetchAiProviderRuntimeState/login/directDB',
          );
        } else {
          // ⏳ OLD: Service layer (Phase 1 - Fallback)
          const data = await aiProviderService.getAiProviderRuntimeState();

          const [enabledChatModelList, enabledImageModelList] = await Promise.all([
            buildProviderModelLists(data.enabledChatAiProviders, data.enabledAiModels, 'chat'),
            buildProviderModelLists(data.enabledImageAiProviders, data.enabledAiModels, 'image'),
          ]);

          set(
            {
              aiProviderRuntimeConfig: data.runtimeConfig,
              builtinAiModelList,
              enabledAiModels: data.enabledAiModels,
              enabledAiProviders: data.enabledAiProviders,
              enabledChatModelList,
              enabledImageModelList,
              isInitAiProviderRuntimeState: true,
            },
            false,
            'internal_fetchAiProviderRuntimeState/login/service',
          );
        }
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
