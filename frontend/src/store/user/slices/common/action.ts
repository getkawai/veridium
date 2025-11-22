import { useEffect } from 'react';
import type { PartialDeep } from 'type-fest';
import type { StateCreator } from 'zustand/vanilla';

import { DEFAULT_PREFERENCE } from '@/const/user';
import type { UserStore } from '@/store/user';
import { DB, toNullString, parseNullableJSON, getNullableString } from '@/types/database';
import { getUserId } from '../../helpers';
import { AsyncLocalStorage } from '@/utils/localStorage';

const preferenceStorage = new AsyncLocalStorage('LOBE_PREFERENCE');
import type { GlobalServerConfig } from '@/types/serverConfig';
import { LobeUser, UserInitializationState } from '@/types/user';
import type { UserSettings } from '@/types/user/settings';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';

import { preferenceSelectors } from '../preference/selectors';

const n = setNamespace('common');

/**
 * 设置操作
 */
export interface CommonAction {
  refreshUserState: () => Promise<void>;
  updateAvatar: (avatar: string) => Promise<void>;
  internal_checkTrace: (shouldFetch: boolean) => Promise<boolean | undefined>;
  useInitUserState: (
    isLogin: boolean | undefined,
    serverConfig: GlobalServerConfig,
    options?: {
      onSuccess: (data: UserInitializationState) => void;
    },
  ) => void;
}

export const createCommonSlice: StateCreator<
  UserStore,
  [['zustand/devtools', never]],
  [],
  CommonAction
> = (set, get) => ({
  refreshUserState: async () => {
    // No-op: handled by useInitUserState
    console.debug('[refreshUserState] Skipped (handled by useEffect)');
  },
  updateAvatar: async (avatar) => {
    // 🔄 MIGRATED: Direct DB call instead of userService.updateAvatar()
    const userId = getUserId();
    const now = Date.now();
    
    await DB.UpdateUser({
      id: userId,
      username: toNullString(''),
      email: toNullString(''),
      avatar: toNullString(avatar),
      phone: toNullString(''),
      firstName: toNullString(''),
      lastName: toNullString(''),
      preference: toNullString(''),
      updatedAt: now,
    });
    
    console.log('[User] Updated avatar via direct DB');

    await get().refreshUserState();
  },

  internal_checkTrace: async (shouldFetch) => {
    if (!shouldFetch) return;

    const userAllowTrace = preferenceSelectors.userAllowTrace(get());

    // if user have set the trace, return false
    if (typeof userAllowTrace === 'boolean') return false;

    return get().isUserCanEnableTrace;
  },

  useInitUserState: (isLogin, serverConfig, options) => {
    useEffect(() => {
      if (!isLogin) return;

      const initUserState = async () => {
        try {
          // 🔄 MIGRATED: Direct DB call instead of userService.getUserState()
          const userId = getUserId();
          
          // Ensure user exists
          await DB.EnsureUserExists({
            id: userId,
            createdAt: Date.now(),
            updatedAt: Date.now(),
          });
          
          // Get user with settings
          const dbUser = await DB.GetUserWithSettings(userId);
          
          // Count messages and sessions
          const [messageCount, sessionCount] = await Promise.all([
            DB.CountMessages(userId),
            DB.CountSessions(userId),
          ]);
          
          // Get preference from LocalStorage
          const preference = await preferenceStorage.getFromLocalStorage();
          
          // Map to UserInitializationState
          const data: UserInitializationState = {
            userId: dbUser.id,
            username: getNullableString(dbUser.username) || undefined,
            email: getNullableString(dbUser.email) || undefined,
            avatar: getNullableString(dbUser.avatar) || '',
            firstName: getNullableString(dbUser.firstName) || undefined,
            lastName: getNullableString(dbUser.lastName) || undefined,
            fullName: [getNullableString(dbUser.firstName), getNullableString(dbUser.lastName)]
              .filter(Boolean)
              .join(' ') || undefined,
            isOnboard: Boolean(dbUser.isOnboarded),
            canEnablePWAGuide: messageCount >= 4,
            canEnableTrace: messageCount >= 4,
            hasConversation: messageCount > 0 || sessionCount > 0,
            preference: (preference || { telemetry: { enabled: false } }) as any,
            settings: {
              tts: parseNullableJSON(dbUser.settingsTts),
              hotkey: parseNullableJSON(dbUser.settingsHotkey),
              keyVaults: parseNullableJSON(dbUser.settingsKeyVaults),
              general: parseNullableJSON(dbUser.settingsGeneral),
              languageModel: parseNullableJSON(dbUser.settingsLanguageModel),
              systemAgent: parseNullableJSON(dbUser.settingsSystemAgent),
              defaultAgent: parseNullableJSON(dbUser.settingsDefaultAgent),
              tool: parseNullableJSON(dbUser.settingsTool),
              image: parseNullableJSON(dbUser.settingsImage),
            } as any,
          };
          
          console.log('[User] Fetched user state via direct DB');
          
          options?.onSuccess?.(data);

          if (data) {
            // merge settings
            const serverSettings: PartialDeep<UserSettings> = {
              defaultAgent: serverConfig.defaultAgent,
              image: serverConfig.image,
              languageModel: serverConfig.languageModel,
              systemAgent: serverConfig.systemAgent,
            };

            const defaultSettings = merge(get().defaultSettings, serverSettings);

            // merge preference
            const isEmpty = Object.keys(data.preference || {}).length === 0;
            const preference = isEmpty ? DEFAULT_PREFERENCE : data.preference;

            // if there is avatar or userId (from client DB), update it into user
            const user =
              data.avatar || data.userId
                ? merge(get().user, {
                    avatar: data.avatar,
                    email: data.email,
                    firstName: data.firstName,
                    fullName: data.fullName,
                    id: data.userId,
                    latestName: data.lastName,
                    username: data.username,
                  } as LobeUser)
                : get().user;

            set(
              {
                defaultSettings,
                isOnboard: data.isOnboard,
                isShowPWAGuide: data.canEnablePWAGuide,
                isUserCanEnableTrace: data.canEnableTrace,
                isUserHasConversation: data.hasConversation,
                isUserStateInit: true,
                preference,
                serverLanguageModel: serverConfig.languageModel,
                settings: data.settings || {},
                subscriptionPlan: data.subscriptionPlan,
                user,
              },
              false,
              n('initUserState'),
            );
            get().refreshDefaultModelProviderList({ trigger: 'fetchUserState' });
          }
        } catch (error) {
          console.error('[useInitUserState] Error:', error);
        }
      };

      initUserState();
    }, [isLogin]);
  },
});
