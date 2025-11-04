import { generateApiKey, isApiKeyExpired, validateApiKeyFormat } from '@/utils/apiKey';

import { ApiKeyItem, NewApiKeyItem } from '../schemas';
import {
  DB,
  toNullInt,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

type EncryptAPIKeyVaults = (keyVaults: string) => Promise<string>;
type DecryptAPIKeyVaults = (keyVaults: string) => Promise<{ plaintext: string }>;

const defaultSerialize = (s: string) => s;

export class ApiKeyModel {
  private userId: string;
  private logger = createModelLogger('ApiKey', 'ApiKeyModel', 'database/models/apiKey');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (
    params: Omit<NewApiKeyItem, 'userId' | 'id' | 'key'>,
    encryptor?: EncryptAPIKeyVaults,
  ) => {
    const key = generateApiKey();

    const encrypt = encryptor || defaultSerialize;

    const encryptedKey = await encrypt(key);

    const now = currentTimestampMs();

    const result = await DB.CreateAPIKey({
      name: params.name,
      key: encryptedKey,
      enabled: boolToInt(params.enabled ?? true),
      expiresAt: toNullInt(params.expiresAt ? new Date(params.expiresAt).getTime() : null),
      lastUsedAt: toNullInt(null),
      userId: this.userId,
      createdAt: now,
      updatedAt: now,
    });

    return this.mapApiKey(result);
  };

  delete = async (id: number) => {
    await DB.DeleteAPIKey({
      id: Number(toNullInt(id as any)) as any,
      userId: this.userId,
    });
  };

  deleteAll = async () => {
    await DB.DeleteAllAPIKeys(this.userId);
  };

  query = async (decryptor?: DecryptAPIKeyVaults) => {
    const results = await DB.ListAPIKeys(this.userId);

    const mapped = results.map((r) => this.mapApiKey(r));

    // If no decryptor provided, return raw results
    if (!decryptor) {
      return mapped;
    }

    // Decrypt each API Key's key field
    const decryptedResults = await Promise.all(
      mapped.map(async (apiKey) => {
        const decryptedKey = await decryptor(apiKey.key);
        return {
          ...apiKey,
          key: decryptedKey.plaintext,
        };
      }),
    );

    return decryptedResults;
  };

  findByKey = async (key: string, encryptor?: EncryptAPIKeyVaults) => {
    if (!validateApiKeyFormat(key)) {
      return null;
    }

    const encrypt = encryptor || defaultSerialize;

    const encryptedKey = await encrypt(key);

    try {
      const result = await DB.GetAPIKeyByKey(encryptedKey);
      return this.mapApiKey(result);
    } catch {
      return null;
    }
  };

  validateKey = async (key: string) => {
    const apiKey = await this.findByKey(key);

    if (!apiKey) return false;
    if (!apiKey.enabled) return false;
    if (apiKey.expiresAt && isApiKeyExpired(apiKey.expiresAt)) return false;

    return true;
  };

  update = async (id: number, value: Partial<ApiKeyItem>) => {
    const now = currentTimestampMs();

    await DB.UpdateAPIKey({
      id: Number(toNullInt(id as any)) as any,
      userId: this.userId,
      name: value.name || '',
      enabled: boolToInt(value.enabled ?? true),
      expiresAt: toNullInt(value.expiresAt ? new Date(value.expiresAt).getTime() : null),
      updatedAt: now,
    });
  };

  findById = async (id: number) => {
    try {
      const result = await DB.GetAPIKey({
        id: Number(toNullInt(id as any)) as any,
        userId: this.userId,
      });
      return this.mapApiKey(result);
    } catch {
      return undefined;
    }
  };

  updateLastUsed = async (id: number) => {
    const now = currentTimestampMs();

    await DB.UpdateAPIKeyLastUsed({
      id: Number(toNullInt(id as any)) as any,
      lastUsedAt: toNullInt(now as any),
    });
  };

  // **************** Helper *************** //

  private mapApiKey = (key: any): ApiKeyItem => {
    return {
      id: key.id,
      name: key.name,
      key: key.key,
      enabled: intToBool(key.enabled),
      expiresAt: key.expiresAt ? new Date(key.expiresAt) : null,
      lastUsedAt: key.lastUsedAt ? new Date(key.lastUsedAt) : null,
      userId: key.userId,
      createdAt: new Date(key.createdAt),
      updatedAt: new Date(key.updatedAt),
    } as ApiKeyItem;
  };
}

