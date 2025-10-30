/* eslint-disable sort-keys-fix/sort-keys-fix  */
import { LobeChatPluginManifest } from '@lobehub/chat-plugin-sdk';
import { integer, primaryKey, sqliteTable, text } from 'drizzle-orm/sqlite-core';

import { DEFAULT_PREFERENCE } from '@/const/user';
import { CustomPluginParams } from '@/types/tool/plugin';

import { timestamps, timestamptz } from './_helpers';

export const users = sqliteTable('users', {
  id: text('id').primaryKey().notNull(),
  username: text('username').unique(),
  email: text('email'),

  avatar: text('avatar'),
  phone: text('phone'),
  firstName: text('first_name'),
  lastName: text('last_name'),
  fullName: text('full_name'),

  isOnboarded: integer('is_onboarded', { mode: 'boolean' }).default(false),
  // Time user was created in Clerk
  clerkCreatedAt: timestamptz('clerk_created_at'),

  // Required by nextauth, all null allowed
  emailVerifiedAt: timestamptz('email_verified_at'),

  preference: text('preference', { mode: 'json' }).$type<typeof DEFAULT_PREFERENCE>().$defaultFn(() => DEFAULT_PREFERENCE),

  ...timestamps,
});

export type NewUser = typeof users.$inferInsert;
export type UserItem = typeof users.$inferSelect;

export const userSettings = sqliteTable('user_settings', {
  id: text('id')
    .references(() => users.id, { onDelete: 'cascade' })
    .primaryKey(),

  tts: text('tts', { mode: 'json' }),
  hotkey: text('hotkey', { mode: 'json' }),
  keyVaults: text('key_vaults'),
  general: text('general', { mode: 'json' }),
  languageModel: text('language_model', { mode: 'json' }),
  systemAgent: text('system_agent', { mode: 'json' }),
  defaultAgent: text('default_agent', { mode: 'json' }),
  tool: text('tool', { mode: 'json' }),
  image: text('image', { mode: 'json' }),
});
export type UserSettingsItem = typeof userSettings.$inferSelect;

export const userInstalledPlugins = sqliteTable(
  'user_installed_plugins',
  {
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),

    identifier: text('identifier').notNull(),
    type: text('type', { enum: ['plugin', 'customPlugin'] }).notNull(),
    manifest: text('manifest', { mode: 'json' }).$type<LobeChatPluginManifest>(),
    settings: text('settings', { mode: 'json' }),
    customParams: text('custom_params', { mode: 'json' }).$type<CustomPluginParams>(),

    ...timestamps,
  },
  (self) => ({
    id: primaryKey({ columns: [self.userId, self.identifier] }),
  }),
);

export type NewInstalledPlugin = typeof userInstalledPlugins.$inferInsert;
export type InstalledPluginItem = typeof userInstalledPlugins.$inferSelect;
