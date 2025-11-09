import type { PartialDeep } from 'type-fest';
import type { StateCreator } from 'zustand/vanilla';

import { DEFAULT_PREFERENCE } from '@/const/user';
import { userService } from '@/services/user';
import type { UserStore } from '@/store/user';
import type { GlobalServerConfig } from '@/types/serverConfig';
import { LobeUser, UserInitializationState } from '@/types/user';
import type { UserSettings } from '@/types/user/settings';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';

import { preferenceSelectors } from '../preference/selectors';

const n = setNamespace('common');

const GET_USER_STATE_KEY = 'initUserState';
/**
 * 设置操作
 */
export interface CommonAction {
  refreshUserState: () => Promise<void>;
  updateAvatar: (avatar: string) => Promise<void>;
  internal_checkTrace: (...) => Promise<void>;
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
    await mutate(GET_USER_STATE_KEY);
  },
  updateAvatar: async (avatar) => {
    // 1. 更新服务端/数据库中的头像
    await userService.updateAvatar(avatar);

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
          const data = await userService.getUserState();
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
