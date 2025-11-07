import { UserGuide, UserKeyVaults, UserPreference, UserSettings, TRPCError } from '@/types';
import dayjs from 'dayjs';
import type { JsonValue, PartialDeep } from 'type-fest';

import { merge } from '@/utils/merge';
import { today } from '@/utils/time';

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

export class UserNotFoundError extends TRPCError {
  constructor() {
    super({ code: 'UNAUTHORIZED', message: 'user not found' });
  }
}

export class UserModel {
  private userId: string;

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
    let result;
    
    try {
      result = await DB.GetUserWithSettings(this.userId);
    } catch {
      throw new UserNotFoundError();
    }

    if (!result) {
      throw new UserNotFoundError();
    }

    // Decrypt keyVaults
    let decryptKeyVaults = {};

    try {
      decryptKeyVaults = await decryptor(
        getNullableString(result.settingsKeyVaults as any) || null,
        this.userId,
      );
      // Map keyVaults to convert any NullString properties to plain strings
      decryptKeyVaults = mapKeyVaults(decryptKeyVaults);
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

    return {
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
      const result = await DB.GetUserSettings(this.userId);
      if (!result) return undefined;

      return {
        id: result.id,
        tts: getNullableString(result.tts as any),
        hotkey: getNullableString(result.hotkey as any),
        keyVaults: getNullableString(result.keyVaults as any),
        general: getNullableString(result.general as any),
        languageModel: getNullableString(result.languageModel as any),
        systemAgent: getNullableString(result.systemAgent as any),
        defaultAgent: getNullableString(result.defaultAgent as any),
        tool: getNullableString(result.tool as any),
        image: getNullableString(result.image as any),
      };
    } catch {
      return undefined;
    }
  };

  updateUser = async (value: Partial<User>) => {
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

    return {
      id: result.id,
      username: getNullableString(result.username as any),
      email: getNullableString(result.email as any),
      avatar: getNullableString(result.avatar as any),
      phone: getNullableString(result.phone as any),
      firstName: getNullableString(result.firstName as any),
      lastName: getNullableString(result.lastName as any),
      preference: parseNullableJSON(result.preference as any),
      isOnboarded: intToBool(result.isOnboarded),
      clerkCreatedAt: result.clerkCreatedAt,
      emailVerifiedAt: result.emailVerifiedAt,
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  };

  deleteSetting = async () => {
    return await DB.DeleteUserSettings(this.userId);
  };

  updateSetting = async (value: Partial<UserSettingsDB>) => {
    const result = await DB.UpsertUserSettings({
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

    return {
      id: result.id,
      tts: getNullableString(result.tts as any),
      hotkey: getNullableString(result.hotkey as any),
      keyVaults: getNullableString(result.keyVaults as any),
      general: getNullableString(result.general as any),
      languageModel: getNullableString(result.languageModel as any),
      systemAgent: getNullableString(result.systemAgent as any),
      defaultAgent: getNullableString(result.defaultAgent as any),
      tool: getNullableString(result.tool as any),
      image: getNullableString(result.image as any),
    };
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

    const result = await DB.UpdateUserPreference({
      id: this.userId,
      preference: toNullJSON(mergedPreference),
      updatedAt: currentTimestampMs(),
    });

    return {
      id: result.id,
      username: getNullableString(result.username as any),
      email: getNullableString(result.email as any),
      avatar: getNullableString(result.avatar as any),
      phone: getNullableString(result.phone as any),
      firstName: getNullableString(result.firstName as any),
      lastName: getNullableString(result.lastName as any),
      preference: parseNullableJSON(result.preference as any),
      isOnboarded: intToBool(result.isOnboarded),
      clerkCreatedAt: result.clerkCreatedAt,
      emailVerifiedAt: result.emailVerifiedAt,
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
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
    
    const result = await DB.UpdateUserPreference({
      id: this.userId,
      preference: toNullJSON({ 
        ...prevPreference, 
        guide: mergedGuide 
      }),
      updatedAt: currentTimestampMs(),
    });

    return {
      id: result.id,
      username: getNullableString(result.username as any),
      email: getNullableString(result.email as any),
      avatar: getNullableString(result.avatar as any),
      phone: getNullableString(result.phone as any),
      firstName: getNullableString(result.firstName as any),
      lastName: getNullableString(result.lastName as any),
      preference: parseNullableJSON(result.preference as any),
      isOnboarded: intToBool(result.isOnboarded),
      clerkCreatedAt: result.clerkCreatedAt,
      emailVerifiedAt: result.emailVerifiedAt,
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
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
    
    const result = await DB.CreateUser({
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

    const user = {
      id: result.id,
      username: getNullableString(result.username as any),
      email: getNullableString(result.email as any),
      avatar: getNullableString(result.avatar as any),
      phone: getNullableString(result.phone as any),
      firstName: getNullableString(result.firstName as any),
      lastName: getNullableString(result.lastName as any),
      preference: parseNullableJSON(result.preference as any),
      isOnboarded: intToBool(result.isOnboarded),
      clerkCreatedAt: result.clerkCreatedAt,
      emailVerifiedAt: result.emailVerifiedAt,
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };

    return { duplicate: false, user };
  };

  static deleteUser = async (_db: any, id: string) => {
    return await DB.DeleteUser(id);
  };

  static findById = async (_db: any, id: string) => {
    try {
      const result = await DB.GetUser(id);
      if (!result) return undefined;

      return {
        id: result.id,
        username: getNullableString(result.username as any),
        email: getNullableString(result.email as any),
        avatar: getNullableString(result.avatar as any),
        phone: getNullableString(result.phone as any),
        firstName: getNullableString(result.firstName as any),
        lastName: getNullableString(result.lastName as any),
        preference: parseNullableJSON(result.preference as any),
        isOnboarded: intToBool(result.isOnboarded),
        clerkCreatedAt: result.clerkCreatedAt,
        emailVerifiedAt: result.emailVerifiedAt,
        createdAt: result.createdAt,
        updatedAt: result.updatedAt,
      };
    } catch {
      return undefined;
    }
  };

  static findByEmail = async (_db: any, email: string) => {
    try {
      const result = await DB.GetUserByEmail(toNullString(email));
      if (!result) return undefined;

      return {
        id: result.id,
        username: getNullableString(result.username as any),
        email: getNullableString(result.email as any),
        avatar: getNullableString(result.avatar as any),
        phone: getNullableString(result.phone as any),
        firstName: getNullableString(result.firstName as any),
        lastName: getNullableString(result.lastName as any),
        preference: parseNullableJSON(result.preference as any),
        isOnboarded: intToBool(result.isOnboarded),
        clerkCreatedAt: result.clerkCreatedAt,
        emailVerifiedAt: result.emailVerifiedAt,
        createdAt: result.createdAt,
        updatedAt: result.updatedAt,
      };
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
    let keyVaults = await decryptor(getNullableString(settings.keyVaults as any) || null, id);
    // Map keyVaults to convert any NullString properties to plain strings
    return mapKeyVaults(keyVaults);
  };
}
