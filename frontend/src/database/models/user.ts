import { UserGuide, UserKeyVaults, UserPreference, UserSettings, TRPCError } from '@/types';
import dayjs from 'dayjs';
import type { JsonValue, PartialDeep } from 'type-fest';

import { merge } from '@/utils/merge';
import { today } from '@/utils/time';
import { createModelLogger } from '@/utils/logger';

import {
  DB,
  type User,
  type UserSetting as UserSettingsDB,
  type CreateUserParams,
  toNullString,
  toNullInt,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';

export interface AdapterAccount {
  type: 'oauth' | 'email' | 'oidc';
  [key: string]: JsonValue | undefined;
}

type DecryptUserKeyVaults = (
  encryptKeyVaultsStr: string | null,
  userId?: string,
) => Promise<UserKeyVaults>;

export class UserNotFoundError extends TRPCError {
  constructor() {
    super({ code: 'UNAUTHORIZED', message: 'user not found' });
  }
}

export class UserModel {
  private userId: string;
  private logger = createModelLogger('User', 'UserModel', 'database/models/user');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  getUserRegistrationDuration = async (): Promise<{
    createdAt: string;
    duration: number;
    updatedAt: string;
  }> => {
    try {
      const user = await DB.GetUser(this.userId);
      
      return {
        createdAt: dayjs(user.createdAt).format('YYYY-MM-DD'),
        duration: dayjs().diff(dayjs(user.createdAt), 'day') + 1,
        updatedAt: today().format('YYYY-MM-DD'),
      };
    } catch {
      return {
        createdAt: today().format('YYYY-MM-DD'),
        duration: 1,
        updatedAt: today().format('YYYY-MM-DD'),
      };
    }
  };

  getUserState = async (decryptor: DecryptUserKeyVaults) => {
    await this.logger.methodEntry('getUserState', { userId: this.userId });
    
    let result;
    
    try {
      result = await DB.GetUserWithSettings(this.userId);
    } catch (error) {
      await this.logger.error('getUserState', 'Failed to get user with settings', { userId: this.userId, error });
      throw new UserNotFoundError();
    }

    if (!result) {
      await this.logger.warn('getUserState', 'User not found', { userId: this.userId });
      throw new UserNotFoundError();
    }

    // Decrypt keyVaults
    let decryptKeyVaults = {};

    try {
      decryptKeyVaults = await decryptor(
        getNullableString(result.settingsKeyVaults as any) || null,
        this.userId,
      );
    } catch {
      /* empty */
    }

    const settings: PartialDeep<UserSettings> = {
      defaultAgent: parseNullableJSON(result.settingsDefaultAgent as any) || {},
      general: parseNullableJSON(result.settingsGeneral as any) || {},
      hotkey: parseNullableJSON(result.settingsHotkey as any) || {},
      image: parseNullableJSON(result.settingsImage as any) || {},
      keyVaults: decryptKeyVaults,
      languageModel: parseNullableJSON(result.settingsLanguageModel as any) || {},
      systemAgent: parseNullableJSON(result.settingsSystemAgent as any) || {},
      tool: parseNullableJSON(result.settingsTool as any) || {},
      tts: parseNullableJSON(result.settingsTts as any) || {},
    };

    const userState = {
      avatar: getNullableString(result.avatar as any) || undefined,
      email: getNullableString(result.email as any) || undefined,
      firstName: getNullableString(result.firstName as any) || undefined,
      fullName: undefined, // Not in schema
      isOnboarded: intToBool(result.isOnboarded),
      lastName: getNullableString(result.lastName as any) || undefined,
      preference: parseNullableJSON(result.preference as any) as UserPreference,
      settings,
      userId: this.userId,
      username: getNullableString(result.username as any) || undefined,
    };

    await this.logger.methodExit('getUserState', userState);
    return userState;
  };

  getUserSSOProviders = async () => {
    const result = await DB.ListNextAuthAccountsByUser(this.userId);
    
    return result.map((account) => ({
      expiresAt: account.expiresAt,
      provider: account.provider,
      providerAccountId: account.providerAccountId,
      scope: getNullableString(account.scope as any),
      type: account.type,
      userId: account.userId,
    })) as unknown as AdapterAccount[];
  };

  getUserSettings = async () => {
    try {
      return await DB.GetUserSettings(this.userId);
    } catch {
      return undefined;
    }
  };

  updateUser = async (value: Partial<User>) => {
    await this.logger.methodEntry('updateUser', { userId: this.userId, value });
    const now = currentTimestampMs();
    
    const result = await DB.UpdateUser({
      id: this.userId,
      username: toNullString(value.username as any),
      email: toNullString(value.email as any),
      avatar: toNullString(value.avatar as any),
      phone: toNullString(value.phone as any),
      firstName: toNullString(value.firstName as any),
      lastName: toNullString(value.lastName as any),
      preference: toNullJSON(value.preference),
      updatedAt: now,
    });

    await this.logger.methodExit('updateUser', result);
    return result;
  };

  deleteSetting = async () => {
    return await DB.DeleteUserSettings(this.userId);
  };

  updateSetting = async (value: Partial<UserSettingsDB>) => {
    return await DB.UpsertUserSettings({
      id: this.userId,
      tts: toNullString(value.tts as any),
      hotkey: toNullString(value.hotkey as any),
      keyVaults: toNullString(value.keyVaults as any),
      general: toNullString(value.general as any),
      languageModel: toNullString(value.languageModel as any),
      systemAgent: toNullString(value.systemAgent as any),
      defaultAgent: toNullString(value.defaultAgent as any),
      tool: toNullString(value.tool as any),
      image: toNullString(value.image as any),
    });
  };

  updatePreference = async (value: Partial<UserPreference>) => {
    let user;
    
    try {
      user = await DB.GetUser(this.userId);
    } catch {
      return;
    }

    if (!user) return;

    const currentPreference = parseNullableJSON(user.preference as any) || {};
    const mergedPreference = merge(currentPreference, value);

    return await DB.UpdateUserPreference({
      id: this.userId,
      preference: toNullJSON(mergedPreference),
      updatedAt: currentTimestampMs(),
    });
  };

  updateGuide = async (value: Partial<UserGuide>) => {
    let user;
    
    try {
      user = await DB.GetUser(this.userId);
    } catch {
      return;
    }

    if (!user) return;

    const prevPreference = (parseNullableJSON(user.preference as any) || {}) as UserPreference;
    const mergedGuide = merge(prevPreference.guide || {}, value);
    
    return await DB.UpdateUserPreference({
      id: this.userId,
      preference: toNullJSON({ 
        ...prevPreference, 
        guide: mergedGuide 
      }),
      updatedAt: currentTimestampMs(),
    });
  };

  // Static methods
  static makeSureUserExist = async (_db: any, userId: string) => {
    const now = currentTimestampMs();
    
    await DB.EnsureUserExists({
      id: userId,
      createdAt: now,
      updatedAt: now,
    });
  };

  static createUser = async (_db: any, params: Partial<CreateUserParams>) => {
    // Check if user already exists
    if (params.id) {
      try {
        const user = await DB.GetUser(params.id);
        if (user) return { duplicate: true };
      } catch {
        // User doesn't exist, continue
      }
    }

    const now = currentTimestampMs();
    
    const user = await DB.CreateUser({
      id: params.id || '',
      username: toNullString(params.username as any),
      email: toNullString(params.email as any),
      avatar: toNullString(params.avatar as any),
      phone: toNullString(params.phone as any),
      firstName: toNullString(params.firstName as any),
      lastName: toNullString(params.lastName as any),
      isOnboarded: boolToInt(false),
      clerkCreatedAt: toNullInt(params.clerkCreatedAt as any),
      emailVerifiedAt: toNullInt(params.emailVerifiedAt as any),
      preference: toNullJSON(params.preference || {}),
      createdAt: now,
      updatedAt: now,
    });

    return { duplicate: false, user };
  };

  static deleteUser = async (_db: any, id: string) => {
    return await DB.DeleteUser(id);
  };

  static findById = async (_db: any, id: string) => {
    try {
      return await DB.GetUser(id);
    } catch {
      return undefined;
    }
  };

  static findByEmail = async (_db: any, email: string) => {
    try {
      return await DB.GetUserByEmail(toNullString(email));
    } catch {
      return undefined;
    }
  };

  static getUserApiKeys = async (_db: any, id: string, decryptor: DecryptUserKeyVaults) => {
    let settings;
    
    try {
      settings = await DB.GetUserSettings(id);
    } catch {
      throw new UserNotFoundError();
    }

    if (!settings) {
      throw new UserNotFoundError();
    }

    // Decrypt keyVaults
    return await decryptor(getNullableString(settings.keyVaults as any) || null, id);
  };
}
