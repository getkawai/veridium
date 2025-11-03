import { nanoid } from 'nanoid';
import { ModelProvider } from '@/model-bank';

import { DEFAULT_MODEL_PROVIDER_LIST } from '@/config/modelProviders';
import { merge } from '@/utils/merge';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';

type DecryptUserKeyVaults = (encryptKeyVaultsStr: string | null) => Promise<any>;
type EncryptUserKeyVaults = (keyVaults: string) => Promise<string>;

export class AiProviderModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (
    { keyVaults: userKey, ...params }: any,
    encryptor?: EncryptUserKeyVaults,
  ) => {
    const defaultSerialize = (s: string) => s;
    const encrypt = encryptor ?? defaultSerialize;
    const keyVaults = await encrypt(JSON.stringify(userKey));

    const now = currentTimestampMs();
    const id = params.id || nanoid();

    const result = await DB.CreateAIProvider({
      id,
      name: toNullString(params.name) as any,
      userId: this.userId,
      sort: params.sort ?? 0,
      enabled: boolToInt(true) as any,
      fetchOnClient: boolToInt(params.fetchOnClient ?? false) as any,
      checkModel: toNullString(params.checkModel) as any,
      logo: toNullString(params.logo) as any,
      description: toNullString(params.description) as any,
      keyVaults: toNullString(keyVaults) as any,
      source: toNullString(params.source || 'custom') as any,
      settings: toNullJSON(params.settings) as any,
      config: toNullJSON(params.config) as any,
      createdAt: now,
      updatedAt: now,
    });

    return result;
  };

  /**
   * NOTE: Also deletes all models of the provider
   * Uses separate queries (no transaction)
   */
  delete = async (id: string) => {
    // 1. Delete all models of the provider
    await DB.DeleteModelsByProvider({
      providerId: toNullString(id) as any,
      userId: this.userId,
    });

    // 2. Delete the provider
    await DB.DeleteAIProvider({
      id,
      userId: this.userId,
    });
  };

  deleteAll = async () => {
    return await DB.DeleteAllAIProviders(this.userId);
  };

  query = async () => {
    return await DB.ListAIProviders(this.userId);
  };

  getAiProviderList = async (): Promise<any[]> => {
    const result = await DB.GetAIProviderListSimple(this.userId);

    return result.map((r) => ({
      description: getNullableString(r.description as any),
      enabled: intToBool(Number(r.enabled) || 0),
      id: r.id,
      logo: getNullableString(r.logo as any),
      name: getNullableString(r.name as any),
      sort: r.sort || 0,
      source: getNullableString(r.source as any),
    }));
  };

  findById = async (id: string) => {
    return await DB.GetAIProvider({
      id,
      userId: this.userId,
    });
  };

  update = async (id: string, value: any) => {
    return await DB.UpdateAIProvider({
      id,
      userId: this.userId,
      name: toNullString(value.name) as any,
      sort: value.sort ?? 0,
      enabled: boolToInt(value.enabled ?? true) as any,
      fetchOnClient: boolToInt(value.fetchOnClient ?? false) as any,
      checkModel: toNullString(value.checkModel) as any,
      logo: toNullString(value.logo) as any,
      description: toNullString(value.description) as any,
      keyVaults: toNullString(value.keyVaults) as any,
      settings: toNullJSON(value.settings) as any,
      config: toNullJSON(value.config) as any,
      updatedAt: currentTimestampMs(),
    });
  };

  updateConfig = async (
    id: string,
    value: any,
    encryptor?: EncryptUserKeyVaults,
  ) => {
    const defaultSerialize = (s: string) => s;
    const encrypt = encryptor ?? defaultSerialize;
    const keyVaults = await encrypt(JSON.stringify(value.keyVaults));

    return await DB.UpsertAIProviderConfig({
      id,
      userId: this.userId,
      keyVaults: toNullString(keyVaults) as any,
      config: toNullJSON(value.config) as any,
      fetchOnClient: boolToInt(value.fetchOnClient ?? false) as any,
      checkModel: toNullString(value.checkModel) as any,
      source: toNullString(this.getProviderSource(id)) as any,
      createdAt: currentTimestampMs(),
      updatedAt: currentTimestampMs(),
    });
  };

  toggleProviderEnabled = async (id: string, enabled: boolean) => {
    return await DB.ToggleAIProviderEnabled({
      id,
      userId: this.userId,
      enabled: boolToInt(enabled) as any,
      source: toNullString(this.getProviderSource(id)) as any,
      createdAt: currentTimestampMs(),
      updatedAt: currentTimestampMs(),
    });
  };

  /**
   * NOTE: No transaction support - updates sequentially
   */
  updateOrder = async (sortMap: { id: string; sort: number }[]) => {
    await Promise.all(
      sortMap.map(({ id, sort }) =>
        DB.UpsertAIProvider({
          id,
          name: toNullString('') as any,
          userId: this.userId,
          sort: sort as any,
          enabled: boolToInt(true) as any,
          fetchOnClient: boolToInt(false) as any,
          checkModel: toNullString('') as any,
          logo: toNullString('') as any,
          description: toNullString('') as any,
          keyVaults: toNullString('') as any,
          source: toNullString(this.getProviderSource(id)) as any,
          settings: toNullJSON(null) as any,
          config: toNullJSON(null) as any,
          createdAt: currentTimestampMs(),
          updatedAt: currentTimestampMs(),
        }),
      ),
    );
  };

  getAiProviderById = async (
    id: string,
    decryptor?: DecryptUserKeyVaults,
  ): Promise<any | undefined> => {
    let result = await DB.GetAIProviderDetail({
      id,
      userId: this.userId,
    });

    if (!result) {
      // If the provider is builtin but not init, insert it
      if (this.isBuiltInProvider(id)) {
        await DB.CreateAIProvider({
          id,
          name: toNullString('') as any,
          userId: this.userId,
          sort: 0 as any,
          enabled: boolToInt(true) as any,
          fetchOnClient: boolToInt(false) as any,
          checkModel: toNullString('') as any,
          logo: toNullString('') as any,
          description: toNullString('') as any,
          keyVaults: toNullString('') as any,
          source: toNullString('builtin') as any,
          settings: toNullJSON(null) as any,
          config: toNullJSON(null) as any,
          createdAt: currentTimestampMs(),
          updatedAt: currentTimestampMs(),
        });

        result = await DB.GetAIProviderDetail({
          id,
          userId: this.userId,
        });
      }

      if (!result) return undefined;
    }

    const decrypt = decryptor ?? JSON.parse;

    let keyVaults = {};

    const keyVaultsStr = getNullableString(result.keyVaults as any);
    if (keyVaultsStr) {
      try {
        keyVaults = await decrypt(keyVaultsStr);
      } catch {
        /* empty */
      }
    }

    const fetchOnClientVal = result.fetchOnClient;
    const fetchOnClient =
      typeof fetchOnClientVal === 'number'
        ? intToBool(fetchOnClientVal)
        : undefined;

    return {
      id: result.id,
      name: getNullableString(result.name as any),
      logo: getNullableString(result.logo as any),
      description: getNullableString(result.description as any),
      enabled: intToBool(Number(result.enabled) || 0),
      source: getNullableString(result.source as any),
      keyVaults,
      settings: parseNullableJSON(result.settings as any) || undefined,
      config: parseNullableJSON(result.config as any),
      fetchOnClient,
      checkModel: getNullableString(result.checkModel as any),
    };
  };

  getAiProviderRuntimeConfig = async (decryptor?: DecryptUserKeyVaults) => {
    const result = await DB.GetAIProviderRuntimeConfigs(this.userId);

    const decrypt = decryptor ?? JSON.parse;
    let runtimeConfig: Record<string, any> = {};

    for (const item of result) {
      const builtin = DEFAULT_MODEL_PROVIDER_LIST.find((provider) => provider.id === item.id);

      const userSettings = parseNullableJSON(item.settings as any) || {};

      let keyVaults = {};
      const keyVaultsStr = getNullableString(item.keyVaults as any);
      if (keyVaultsStr) {
        try {
          keyVaults = await decrypt(keyVaultsStr);
        } catch {
          /* empty */
        }
      }

      const fetchOnClientVal = item.fetchOnClient;
      const fetchOnClient =
        typeof fetchOnClientVal === 'number'
          ? intToBool(fetchOnClientVal)
          : undefined;

      runtimeConfig[item.id] = {
        config: parseNullableJSON(item.config as any) || {},
        fetchOnClient,
        keyVaults,
        settings: builtin ? merge(builtin.settings, userSettings) : userSettings,
      };
    }

    return runtimeConfig;
  };

  private isBuiltInProvider = (id: string) => Object.values(ModelProvider).includes(id as any);

  private getProviderSource = (id: string) => (this.isBuiltInProvider(id) ? 'builtin' : 'custom');
}

