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
    await aiProviderService.toggleProviderEnabled(id, enabled);
    await get().refreshAiProviderList();

    get().internal_toggleAiProviderLoading(id, false);
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
    await aiProviderService.updateAiProviderOrder(items);
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

    try {
      const data = await aiProviderService.getAiProviderList();

      if (!get().initAiProviderList) {
        set(
          { aiProviderList: data, initAiProviderList: true },
          false,
          'internal_fetchAiProviderList/init',
        );
        return;
      }

      set({ aiProviderList: data }, false, 'internal_fetchAiProviderList/refresh');
    } catch (error) {
      console.error('[internal_fetchAiProviderList] Error:', error);
    }
  },

  internal_fetchAiProviderRuntimeState: async (isLogin) => {
    const isAuthLoaded = authSelectors.isLoaded(useUserStore.getState());
    const shouldFetch =
      isAuthLoaded && !isDeprecatedEdition && isLogin !== null && isLogin !== undefined;

    if (!shouldFetch) return;

    try {
      const [{ LOBE_DEFAULT_MODEL_LIST: builtinAiModelList }] = await Promise.all([
        import('@/model-bank'),
      ]);

      if (isLogin) {
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
          'internal_fetchAiProviderRuntimeState/login',
        );
      } else {
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
