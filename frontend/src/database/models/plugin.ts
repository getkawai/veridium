import { LobeTool } from '@/types';

import { InstalledPluginItem, NewInstalledPlugin } from '../schemas';
import {
  DB,
  toNullJSON,
  parseNullableJSON,
  currentTimestampMs,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

export class PluginModel {
  private userId: string;
  private logger = createModelLogger('Plugin', 'PluginModel', 'database/models/plugin');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (
    params: Pick<
      NewInstalledPlugin,
      'type' | 'identifier' | 'manifest' | 'customParams' | 'settings'
    >,
  ) => {
    const now = currentTimestampMs();

    const result = await DB.UpsertPlugin({
      identifier: params.identifier,
      type: params.type,
      manifest: toNullJSON(params.manifest),
      customParams: toNullJSON(params.customParams),
      settings: toNullJSON(params.settings),
      userId: this.userId,
      createdAt: now,
      updatedAt: now,
    });

    return this.mapPlugin(result);
  };

  delete = async (id: string) => {
    await DB.DeletePlugin({
      identifier: id,
      userId: this.userId,
    });
  };

  deleteAll = async () => {
    await DB.DeleteAllPlugins(this.userId);
  };

  query = async () => {
    const data = await DB.ListPlugins(this.userId);

    return data.map<LobeTool>((item) => {
      const manifest = parseNullableJSON(item.manifest as any);
      return {
        customParams: parseNullableJSON(item.customParams as any),
        identifier: item.identifier,
        manifest: manifest,
        settings: parseNullableJSON(item.settings as any),
        source: item.type as any,
        type: item.type as any,
        runtimeType: manifest?.type || 'default',
      };
    });
  };

  findById = async (id: string) => {
    try {
      const plugin = await DB.GetPlugin({
        identifier: id,
        userId: this.userId,
      });
      return this.mapPlugin(plugin);
    } catch {
      return undefined;
    }
  };

  update = async (id: string, value: Partial<InstalledPluginItem>) => {
    const now = currentTimestampMs();

    await DB.UpdatePlugin({
      identifier: id,
      userId: this.userId,
      type: value.type || '',
      manifest: toNullJSON(value.manifest),
      customParams: toNullJSON(value.customParams),
      settings: toNullJSON(value.settings),
      updatedAt: now,
    });
  };

  // **************** Helper *************** //

  private mapPlugin = (plugin: any): InstalledPluginItem => {
    return {
      identifier: plugin.identifier,
      type: plugin.type,
      manifest: parseNullableJSON(plugin.manifest as any),
      customParams: parseNullableJSON(plugin.customParams as any),
      settings: parseNullableJSON(plugin.settings as any),
      userId: plugin.userId,
      createdAt: new Date(plugin.createdAt),
      updatedAt: new Date(plugin.updatedAt),
    } as InstalledPluginItem;
  };
}

