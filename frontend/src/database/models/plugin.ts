import { LobeTool } from '@/types';

import { InstalledPluginItem, NewInstalledPlugin } from '@/types/plugin/installedPlugin';
import {
  DB,
  toNullJSON,
  parseNullableJSON,
  currentTimestampMs,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';
import { NotificationService, NotificationOptions } from '@@/github.com/wailsapp/wails/v3/pkg/services/notifications';

export class PluginModel {
  private userId: string;
  private logger = createModelLogger('Plugin', 'PluginModel', 'database/models/plugin');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  /**
   * Show error notification to user
   */
  private async showErrorNotification(title: string, message: string) {
    try {
      await NotificationService.SendNotification(
        new NotificationOptions({
          id: `plugin-error-${Date.now()}`,
          title: `Plugin Error: ${title}`,
          body: message,
        })
      );
    } catch (notifError) {
      // Silently fail if notification fails - don't want notification errors to break the app
      console.error('Failed to show notification:', notifError);
    }
  }

  create = async (
    params: Pick<
      NewInstalledPlugin,
      'type' | 'identifier' | 'manifest' | 'customParams' | 'settings'
    >,
  ) => {
    await this.logger.methodEntry('create', { identifier: params.identifier, userId: this.userId });
    
    try {
      const now = currentTimestampMs();

      const result = await DB.UpsertPlugin({
        identifier: params.identifier,
        type: params.type,
        manifest: toNullJSON(params.manifest),
        customParams: toNullJSON(params.customParams),
        settings: toNullJSON(params.settings),
        createdAt: now,
        updatedAt: now,
      });

      await this.logger.methodExit('create', { identifier: result.identifier });
      return this.mapPlugin(result);
    } catch (error) {
      await this.logger.error('Failed to create plugin', { error, params });
      await this.showErrorNotification(
        'Install Failed',
        `Failed to install plugin "${params.identifier}". Please try again.`
      );
      throw error;
    }
  };

  delete = async (id: string) => {
    await this.logger.methodEntry('delete', { identifier: id, userId: this.userId });
    
    try {
      await DB.DeletePlugin(id);
      
      await this.logger.methodExit('delete', { identifier: id });
    } catch (error) {
      await this.logger.error('Failed to delete plugin', { error, id });
      await this.showErrorNotification(
        'Uninstall Failed',
        `Failed to uninstall plugin. Please try again.`
      );
      throw error;
    }
  };

  deleteAll = async () => {
    try {
      await DB.DeleteAllPlugins();
    } catch (error) {
      await this.logger.error('Failed to delete all plugins', { error });
      await this.showErrorNotification(
        'Uninstall All Failed',
        `Failed to uninstall all plugins. Please try again.`
      );
      throw error;
    }
  };

  query = async () => {
    try {
      const data = await DB.ListPlugins();

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
    } catch (error) {
      await this.logger.error('Failed to query plugins', { error });
      throw error;
    }
  };

  findById = async (id: string) => {
    try {
      const plugin = await DB.GetPlugin(id);
      return this.mapPlugin(plugin);
    } catch (error) {
      await this.logger.warn('Plugin not found', { id, error });
      return undefined;
    }
  };

  update = async (id: string, value: Partial<InstalledPluginItem>) => {
    await this.logger.methodEntry('update', { identifier: id, userId: this.userId });
    
    try {
      const now = currentTimestampMs();

      await DB.UpdatePlugin({
        identifier: id,
        type: value.type || '',
        manifest: toNullJSON(value.manifest),
        customParams: toNullJSON(value.customParams),
        settings: toNullJSON(value.settings),
        updatedAt: now,
      });
      
      await this.logger.methodExit('update', { identifier: id });
    } catch (error) {
      await this.logger.error('Failed to update plugin', { error, id, value });
      await this.showErrorNotification(
        'Update Failed',
        `Failed to update plugin "${id}". Please try again.`
      );
      throw error;
    }
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

