import isEqual from 'fast-deep-equal';
import type { PartialDeep } from 'type-fest';
import type { StateCreator } from 'zustand/vanilla';

import { MESSAGE_CANCEL_FLAT } from '@/const/message';
import { shareService } from '@/services/share';
import type { UserStore } from '@/store/user';
import { DB, toNullString } from '@/types/database';
import { getUserId } from '../../helpers';
import { LobeAgentSettings } from '@/types/session';
import {
  SystemAgentItem,
  UserGeneralConfig,
  UserKeyVaults,
  UserSettings,
  UserSystemAgentConfigKey,
} from '@/types/user/settings';
import { difference } from '@/utils/difference';
import { merge } from '@/utils/merge';

export interface UserSettingsAction {
  importAppSettings: (settings: UserSettings) => Promise<void>;
  importUrlShareSettings: (settingsParams: string | null) => Promise<void>;
  internal_createSignal: () => AbortController;
  resetSettings: () => Promise<void>;
  setSettings: (settings: PartialDeep<UserSettings>) => Promise<void>;
  updateDefaultAgent: (agent: PartialDeep<LobeAgentSettings>) => Promise<void>;
  updateGeneralConfig: (settings: Partial<UserGeneralConfig>) => Promise<void>;
  updateKeyVaults: (settings: Partial<UserKeyVaults>) => Promise<void>;

  updateSystemAgent: (
    key: UserSystemAgentConfigKey,
    value: Partial<SystemAgentItem>,
  ) => Promise<void>;
}

export const createSettingsSlice: StateCreator<
  UserStore,
  [['zustand/devtools', never]],
  [],
  UserSettingsAction
> = (set, get) => ({
  importAppSettings: async (importAppSettings) => {
    const { setSettings } = get();

    await setSettings(importAppSettings);
  },

  /**
   * Import settings from a string in json format
   */
  importUrlShareSettings: async (settingsParams: string | null) => {
    if (settingsParams) {
      const importSettings = shareService.decodeShareSettings(settingsParams);
      if (importSettings?.message || !importSettings?.data) {
        // handle some error
        return;
      }

      await get().setSettings(importSettings.data);
    }
  },

  internal_createSignal: () => {
    const abortController = get().updateSettingsSignal;
    if (abortController && !abortController.signal.aborted)
      abortController.abort(MESSAGE_CANCEL_FLAT);

    const newSignal = new AbortController();

    set({ updateSettingsSignal: newSignal }, false, 'signalForUpdateSettings');

    return newSignal;
  },

  resetSettings: async () => {
    // 🔄 MIGRATED: Direct DB call instead of userService.resetUserSettings()
    const userId = getUserId();
    await DB.DeleteUserSettings(userId);
    
    console.log('[User] Reset user settings via direct DB');
    
    await get().refreshUserState();
  },
  setSettings: async (settings) => {
    const { settings: prevSetting, defaultSettings } = get();

    const nextSettings = merge(prevSetting, settings);

    if (isEqual(prevSetting, nextSettings)) return;

    const diffs = difference(nextSettings, defaultSettings);
    set({ settings: diffs }, false, 'optimistic_updateSettings');

    const abortController = get().internal_createSignal();
    
    // 🔄 MIGRATED: Direct DB call instead of userService.updateUserSettings()
    const userId = getUserId();
    const { keyVaults, ...res } = diffs;
    
    await DB.UpsertUserSettings({
      id: userId,
      tts: toNullString(JSON.stringify(res.tts || {})),
      hotkey: toNullString(JSON.stringify(res.hotkey || {})),
      keyVaults: toNullString(JSON.stringify(keyVaults || {})),
      general: toNullString(JSON.stringify(res.general || {})),
      languageModel: toNullString(JSON.stringify(res.languageModel || {})),
      systemAgent: toNullString(JSON.stringify(res.systemAgent || {})),
      defaultAgent: toNullString(JSON.stringify(res.defaultAgent || {})),
      tool: toNullString(JSON.stringify(res.tool || {})),
      image: toNullString(JSON.stringify(res.image || {})),
    });
    
    console.log('[User] Updated user settings via direct DB');
    
    await get().refreshUserState();
  },
  updateDefaultAgent: async (defaultAgent) => {
    await get().setSettings({ defaultAgent });
  },
  updateGeneralConfig: async (general) => {
    await get().setSettings({ general });
  },
  updateKeyVaults: async (keyVaults) => {
    await get().setSettings({ keyVaults });
  },
  updateSystemAgent: async (key, value) => {
    await get().setSettings({
      systemAgent: { [key]: { ...value } },
    });
  },
});
