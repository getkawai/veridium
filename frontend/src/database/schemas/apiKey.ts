/* eslint-disable sort-keys-fix/sort-keys-fix  */
import { boolean, integer, sqliteTable, text, varchar } from 'drizzle-orm/sqlite-core';
import { createInsertSchema } from 'drizzle-zod';

import { timestamps, timestamptz } from './_helpers';
import { users } from './user';

export const apiKeys = sqliteTable('api_keys', {
  id: integer('id').primaryKey().generatedByDefaultAsIdentity(), // auto-increment primary key
  name: text('name').notNull(), // name of the API key
  key: text('key').notNull().unique(), // API key
  enabled: integer('enabled').default(true), // whether the API key is enabled
  expiresAt: timestamptz('expires_at'), // expires time
  lastUsedAt: timestamptz('last_used_at'), // last used time
  userId: text('user_id')
    .references(() => users.id, { onDelete: 'cascade' })
    .notNull(), // belongs to user, when user is deleted, the API key will be deleted

  ...timestamps,
});

export const insertApiKeySchema = createInsertSchema(apiKeys);

export type ApiKeyItem = typeof apiKeys.$inferSelect;
export type NewApiKeyItem = typeof apiKeys.$inferInsert;
