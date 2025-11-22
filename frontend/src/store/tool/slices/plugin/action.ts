import { Schema, ValidationResult } from '@cfworker/json-schema';
import { StateCreator } from 'zustand/vanilla';

import { MESSAGE_CANCEL_FLAT } from '@/const/message';
import { pluginService } from '@/services/plugin';
import { merge } from '@/utils/merge';
import { DB, toNullString, toNullJSON, parseNullableJSON, currentTimestampMs } from '@/types/database';
import { getUserId } from '@/store/session/helpers';

import { ToolStore } from '../../store';
import { pluginStoreSelectors } from '../oldStore/selectors';
import { pluginSelectors } from './selectors';

/**
 * 插件接口
 */
export interface PluginAction {
  checkPluginsIsInstalled: (plugins: string[]) => Promise<void>;
  removeAllPlugins: () => Promise<void>;
  updateInstallMcpPlugin: (id: string, value: any) => Promise<void>;
  updatePluginSettings: <T>(
    id: string,
    settings: Partial<T>,
    options?: { override?: boolean },
  ) => Promise<void>;
  internal_checkPluginsIsInstalled: (enable: boolean, plugins: string[]) => Promise<void>;
  validatePluginSettings: (identifier: string) => Promise<ValidationResult | undefined>;
}

export const createPluginSlice: StateCreator<
  ToolStore,
  [['zustand/devtools', never]],
  [],
  PluginAction
> = (set, get) => ({
  checkPluginsIsInstalled: async (plugins) => {
    // if there is no plugins, just skip.
    if (plugins.length === 0) return;

    const { loadPluginStore, installPlugins } = get();

    // check if the store is empty
    // if it is, we need to load the plugin store
    if (pluginStoreSelectors.onlinePluginStore(get()).length === 0) {
      await loadPluginStore();
    }

    await installPlugins(plugins);
  },
  removeAllPlugins: async () => {
    const userId = getUserId();
    await DB.DeleteAllPlugins(userId);

    console.log('[Plugin] Removed all plugins via direct DB', { userId });

    await get().refreshPlugins();
  },

  updateInstallMcpPlugin: async (id, value) => {
    const installedPlugin = pluginSelectors.getInstalledPluginById(id)(get());

    if (!installedPlugin) return;

    const userId = getUserId();
    const now = currentTimestampMs();

    await DB.UpdatePlugin({
      identifier: id,
      userId,
      customParams: toNullJSON({ mcp: merge(installedPlugin.customParams?.mcp, value) }),
      updatedAt: now,
    });

    console.log('[Plugin] Updated MCP plugin via direct DB', { id, hasMcpValue: !!value });

    await get().refreshPlugins();
  },

  updatePluginSettings: async (id, settings, { override } = {}) => {
    const signal = get().updatePluginSettingsSignal;
    if (signal) signal.abort(MESSAGE_CANCEL_FLAT);

    const newSignal = new AbortController();

    const previousSettings = pluginSelectors.getPluginSettingsById(id)(get());
    const nextSettings = override ? settings : merge(previousSettings, settings);

    set({ updatePluginSettingsSignal: newSignal }, false, 'create new Signal');

    const userId = getUserId();
    const now = currentTimestampMs();

    await DB.UpdatePlugin({
      identifier: id,
      userId,
      settings: toNullJSON(nextSettings),
      updatedAt: now,
    });

    console.log('[Plugin] Updated plugin settings via direct DB', { id, hasSettings: !!nextSettings });

    await get().refreshPlugins();
  },
  internal_checkPluginsIsInstalled: async (enable, plugins) => {
    if (!enable || plugins.length === 0) return;

    await get().checkPluginsIsInstalled(plugins);
  },
  validatePluginSettings: async (identifier) => {
    const manifest = pluginSelectors.getToolManifestById(identifier)(get());
    if (!manifest || !manifest.settings) return;
    const settings = pluginSelectors.getPluginSettingsById(identifier)(get());

    // validate the settings
    const { Validator } = await import('@cfworker/json-schema');
    const validator = new Validator(manifest.settings as Schema);
    const result = validator.validate(settings);

    if (!result.valid) return { errors: result.errors, valid: false };

    return { errors: [], valid: true };
  },
});
