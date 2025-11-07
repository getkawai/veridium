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
  AiProvider,
} from '@/types/database';
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';
import { createModelLogger } from '@/utils/logger';

type DecryptUserKeyVaults = (encryptKeyVaultsStr: string | null) => Promise<any>;
type EncryptUserKeyVaults = (keyVaults: string) => Promise<string>;

/**
 * Map keyVaults object to convert any NullString properties to plain strings
 * This ensures that downstream code doesn't need to handle NullString types
 */
const mapKeyVaults = (keyVaults: any): any => {
  if (!keyVaults || typeof keyVaults !== 'object') {
    return keyVaults;
  }

  const mapped: any = {};
  for (const [key, value] of Object.entries(keyVaults)) {
    // Check if value is NullString
    if (value && typeof value === 'object' && 'String' in value && 'Valid' in value) {
      mapped[key] = getNullableString(value as any);
    } else if (value && typeof value === 'object' && !Array.isArray(value)) {
      // Recursively map nested objects
      mapped[key] = mapKeyVaults(value);
    } else {
      mapped[key] = value;
    }
  }
  return mapped;
};

export class AiProviderModel {
  private userId: string;
  private logger = createModelLogger('AiProvider', 'AiProviderModel', 'database/models/aiProvider');

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

    // Map result to standard TypeScript object
    return {
      id: result.id,
      name: getNullableString(result.name as any),
      userId: result.userId,
      sort: result.sort || 0,
      enabled: intToBool(Number(result.enabled) || 0),
      fetchOnClient: intToBool(Number(result.fetchOnClient) || 0),
      checkModel: getNullableString(result.checkModel as any),
      logo: getNullableString(result.logo as any),
      description: getNullableString(result.description as any),
      keyVaults: getNullableString(result.keyVaults as any),
      source: getNullableString(result.source as any),
      settings: parseNullableJSON(result.settings as any),
      config: parseNullableJSON(result.config as any),
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  };

  /**
   * Delete AI provider with atomic transaction
   * ✅ OPTIMIZED: Uses backend transaction for atomicity
   * Deletes provider and all its models atomically
   * All operations succeed or all rollback
   */
  delete = async (id: string) => {
    await DBService.DeleteAIProviderWithModels(id, this.userId);
  };

  deleteAll = async () => {
    const result = await DB.DeleteAllAIProviders(this.userId);
    return result;
  };

  query = async () => {
    const result = await DB.ListAIProviders(this.userId);
    
    return result.map((r) => ({
      id: r.id,
      name: getNullableString(r.name as any),
      userId: r.userId,
      sort: r.sort || 0,
      enabled: intToBool(Number(r.enabled) || 0),
      fetchOnClient: intToBool(Number(r.fetchOnClient) || 0),
      checkModel: getNullableString(r.checkModel as any),
      logo: getNullableString(r.logo as any),
      description: getNullableString(r.description as any),
      keyVaults: getNullableString(r.keyVaults as any),
      source: getNullableString(r.source as any),
      settings: parseNullableJSON(r.settings as any),
      config: parseNullableJSON(r.config as any),
      createdAt: r.createdAt,
      updatedAt: r.updatedAt,
    }));
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
    const result = await DB.GetAIProvider({
      id,
      userId: this.userId,
    });

    if (!result) return undefined;

    return {
      id: result.id,
      name: getNullableString(result.name as any),
      userId: result.userId,
      sort: result.sort || 0,
      enabled: intToBool(Number(result.enabled) || 0),
      fetchOnClient: intToBool(Number(result.fetchOnClient) || 0),
      checkModel: getNullableString(result.checkModel as any),
      logo: getNullableString(result.logo as any),
      description: getNullableString(result.description as any),
      keyVaults: getNullableString(result.keyVaults as any),
      source: getNullableString(result.source as any),
      settings: parseNullableJSON(result.settings as any),
      config: parseNullableJSON(result.config as any),
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  };

  update = async (id: string, value: any) => {
    const result = await DB.UpdateAIProvider({
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

    return {
      id: result.id,
      name: getNullableString(result.name as any),
      userId: result.userId,
      sort: result.sort || 0,
      enabled: intToBool(Number(result.enabled) || 0),
      fetchOnClient: intToBool(Number(result.fetchOnClient) || 0),
      checkModel: getNullableString(result.checkModel as any),
      logo: getNullableString(result.logo as any),
      description: getNullableString(result.description as any),
      keyVaults: getNullableString(result.keyVaults as any),
      source: getNullableString(result.source as any),
      settings: parseNullableJSON(result.settings as any),
      config: parseNullableJSON(result.config as any),
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  };

  updateConfig = async (
    id: string,
    value: any,
    encryptor?: EncryptUserKeyVaults,
  ) => {
    const defaultSerialize = (s: string) => s;
    const encrypt = encryptor ?? defaultSerialize;
    const keyVaults = await encrypt(JSON.stringify(value.keyVaults));

    const result = await DB.UpsertAIProviderConfig({
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

    return {
      id: result.id,
      name: getNullableString(result.name as any),
      userId: result.userId,
      sort: result.sort || 0,
      enabled: intToBool(Number(result.enabled) || 0),
      fetchOnClient: intToBool(Number(result.fetchOnClient) || 0),
      checkModel: getNullableString(result.checkModel as any),
      logo: getNullableString(result.logo as any),
      description: getNullableString(result.description as any),
      keyVaults: getNullableString(result.keyVaults as any),
      source: getNullableString(result.source as any),
      settings: parseNullableJSON(result.settings as any),
      config: parseNullableJSON(result.config as any),
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  };

  toggleProviderEnabled = async (id: string, enabled: boolean) => {
    const result = await DB.ToggleAIProviderEnabled({
      id,
      userId: this.userId,
      enabled: boolToInt(enabled) as any,
      source: toNullString(this.getProviderSource(id)) as any,
      createdAt: currentTimestampMs(),
      updatedAt: currentTimestampMs(),
    });

    return {
      id: result.id,
      name: getNullableString(result.name as any),
      userId: result.userId,
      sort: result.sort || 0,
      enabled: intToBool(Number(result.enabled) || 0),
      fetchOnClient: intToBool(Number(result.fetchOnClient) || 0),
      checkModel: getNullableString(result.checkModel as any),
      logo: getNullableString(result.logo as any),
      description: getNullableString(result.description as any),
      keyVaults: getNullableString(result.keyVaults as any),
      source: getNullableString(result.source as any),
      settings: parseNullableJSON(result.settings as any),
      config: parseNullableJSON(result.config as any),
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
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
        // Map keyVaults to convert any NullString properties to plain strings
        keyVaults = mapKeyVaults(keyVaults);
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
          // Map keyVaults to convert any NullString properties to plain strings
          keyVaults = mapKeyVaults(keyVaults);
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

