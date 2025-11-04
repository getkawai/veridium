/* eslint-disable sort-keys-fix/sort-keys-fix  */

import { LobeChatPluginManifest } from '@/chat-plugin-sdk';
import { CustomPluginParams } from '@/types/tool/plugin';

/**
 * Plugin type enum
 */
export type PluginType = 'plugin' | 'customPlugin';

/**
 * Installed plugin item (database record)
 * Equivalent to: typeof userInstalledPlugins.$inferSelect
 */
export interface InstalledPluginItem {
  userId: string;
  identifier: string;
  type: PluginType;
  manifest?: LobeChatPluginManifest;
  settings?: any;
  customParams?: CustomPluginParams;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * New installed plugin (for insert operations)
 * Equivalent to: typeof userInstalledPlugins.$inferInsert
 */
export interface NewInstalledPlugin {
  userId: string;
  identifier: string;
  type: PluginType;
  manifest?: LobeChatPluginManifest;
  settings?: any;
  customParams?: CustomPluginParams;
  createdAt?: Date;
  updatedAt?: Date;
}

