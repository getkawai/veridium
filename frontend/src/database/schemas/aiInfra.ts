/* eslint-disable sort-keys-fix/sort-keys-fix  */
import { boolean, integer, jsonb, pgTable, primaryKey, text, varchar } from 'drizzle-orm/sqlite-core';

import { AiProviderConfig, AiProviderSettings } from '@/types/aiProvider';

import { timestamps } from './_helpers';
import { users } from './user';

export const aiProviders = sqliteTable(
  'ai_providers',
  {
    id: text('id').notNull(),
    name: text('name'),

    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),

    sort: integer('sort', { mode: 'boolean' }),
    enabled: integer('enabled'),
    fetchOnClient: integer('fetch_on_client'),
    checkModel: text('check_model'),
    logo: text('logo'),
    description: text('description'),

    // need to be encrypted
    keyVaults: text('key_vaults'),
    source: text('source', { enum: ['builtin', 'custom'], length: 20 }),
    settings: text('settings')
      .$defaultFn(() => ({}))
      .$type<AiProviderSettings>(),

    config: text('config')
      .$defaultFn(() => ({}))
      .$type<AiProviderConfig>(),

    ...timestamps,
  },
  (table) => [primaryKey({ columns: [table.id, table.userId] })],
);

export type NewAiProviderItem = Omit<typeof aiProviders.$inferInsert, 'userId'>;
export type AiProviderSelectItem = typeof aiProviders.$inferSelect;

export const aiModels = sqliteTable(
  'ai_models',
  {
    id: text('id').notNull(),
    displayName: text('display_name'),
    description: text('description'),
    organization: text('organization'),
    enabled: integer('enabled'),
    providerId: text('provider_id').notNull(),
    type: text('type').default('chat').notNull(),
    sort: integer('sort', { mode: 'boolean' }),

    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
    pricing: text('pricing'),
    parameters: text('parameters').$defaultFn(() => ({})),
    config: text('config'),
    abilities: text('abilities').$defaultFn(() => ({})),
    contextWindowTokens: integer('context_window_tokens'),
    source: text('source', { enum: ['remote', 'custom', 'builtin'], length: 20 }),
    releasedAt: text('released_at'),

    ...timestamps,
  },
  (table) => [primaryKey({ columns: [table.id, table.providerId, table.userId] })],
);

export type NewAiModelItem = Omit<typeof aiModels.$inferInsert, 'userId'>;
export type AiModelSelectItem = typeof aiModels.$inferSelect;
