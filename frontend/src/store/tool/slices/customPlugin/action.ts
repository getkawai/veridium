import { LobeChatPluginManifest } from '@/chat-plugin-sdk';
import { t } from 'i18next';
import { merge } from 'lodash-es';
import { StateCreator } from 'zustand/vanilla';

import { notification } from '@/components/AntdStaticMethods';
import { mcpService } from '@/services/mcp';
import { pluginService } from '@/services/plugin';
import { toolService } from '@/services/tool';
import { pluginHelpers } from '@/store/tool/helpers';
import { DB, toNullString, toNullJSON, currentTimestampMs } from '@/types/database';
import { getUserId } from '@/store/session/helpers';
import { LobeToolCustomPlugin, PluginInstallError } from '@/types/tool/plugin';
import { setNamespace } from '@/utils/storeDebug';

import { ToolStore } from '../../store';
import { pluginSelectors } from '../plugin/selectors';
import { defaultCustomPlugin } from './initialState';

const n = setNamespace('customPlugin');

export interface CustomPluginAction {
  installCustomPlugin: (value: LobeToolCustomPlugin) => Promise<void>;
  reinstallCustomPlugin: (id: string) => Promise<void>;
  uninstallCustomPlugin: (id: string) => Promise<void>;
  updateCustomPlugin: (id: string, value: LobeToolCustomPlugin) => Promise<void>;
  updateNewCustomPlugin: (value: Partial<LobeToolCustomPlugin>) => void;
}

export const createCustomPluginSlice: StateCreator<
  ToolStore,
  [['zustand/devtools', never]],
  [],
  CustomPluginAction
> = (set, get) => ({
  installCustomPlugin: async (value) => {
    const userId = getUserId();
    const now = currentTimestampMs();

    await DB.UpsertPlugin({
      identifier: value.identifier,
      userId,
      type: 'customPlugin',
      manifest: toNullJSON(value.manifest),
      customParams: toNullJSON(value.customParams),
      settings: toNullJSON(value.settings),
      createdAt: now,
      updatedAt: now,
    });

    console.log('[Plugin] Created custom plugin via direct DB', { identifier: value.identifier });

    await get().refreshPlugins();
    set({ newCustomPlugin: defaultCustomPlugin }, false, n('saveToCustomPluginList'));
  },
  reinstallCustomPlugin: async (id) => {
    const plugin = pluginSelectors.getCustomPluginById(id)(get());
    if (!plugin) return;

    const { refreshPlugins, updateInstallLoadingState } = get();

    try {
      updateInstallLoadingState(id, true);
      let manifest: LobeChatPluginManifest;
      // mean this is a mcp plugin
      if (!!plugin.customParams?.mcp) {
        const url = plugin.customParams?.mcp?.url;
        if (!url) return;

        manifest = await mcpService.getStreamableMcpServerManifest({
          auth: plugin.customParams.mcp.auth,
          headers: plugin.customParams.mcp.headers,
          identifier: plugin.identifier,
          metadata: {
            avatar: plugin.customParams.avatar,
            description: plugin.customParams.description,
          },
          url,
        });
      } else {
        manifest = await toolService.getToolManifest(
          plugin.customParams?.manifestUrl,
          plugin.customParams?.useProxy,
        );
      }
      updateInstallLoadingState(id, false);

      const userId = getUserId();
      const now = currentTimestampMs();

      await DB.UpdatePlugin({
        identifier: id,
        userId,
        manifest: toNullJSON(manifest),
        updatedAt: now,
      });

      console.log('[Plugin] Updated plugin manifest via direct DB', { id });

      await refreshPlugins();
    } catch (error) {
      updateInstallLoadingState(id, false);

      console.error(error);
      const err = error as PluginInstallError;

      const meta = pluginSelectors.getPluginMetaById(id)(get());
      const name = pluginHelpers.getPluginTitle(meta);

      notification.error({
        description: t(`error.${err.message}`, { error: err.cause, ns: 'plugin' }),
        message: t('error.reinstallError', { name, ns: 'plugin' }),
      });
    }
  },
  uninstallCustomPlugin: async (id) => {
    const userId = getUserId();
    await DB.DeletePlugin({
      identifier: id,
      userId,
    });

    console.log('[Plugin] Uninstalled custom plugin via direct DB', { id });

    await get().refreshPlugins();
  },

  updateCustomPlugin: async (id, value) => {
    const { reinstallCustomPlugin } = get();

    // 1. Update plugin info
    const userId = getUserId();
    const now = currentTimestampMs();

    await DB.UpdatePlugin({
      identifier: id,
      userId,
      type: value.type || '',
      manifest: toNullJSON(value.manifest),
      customParams: toNullJSON(value.customParams),
      settings: toNullJSON(value.settings),
      updatedAt: now,
    });

    console.log('[Plugin] Updated custom plugin via direct DB', { id });

    // 2. 重新安装插件
    await reinstallCustomPlugin(id);
  },
  updateNewCustomPlugin: (newCustomPlugin) => {
    set(
      { newCustomPlugin: merge({}, get().newCustomPlugin, newCustomPlugin) },
      false,
      n('updateNewDevPlugin'),
    );
  },
});
