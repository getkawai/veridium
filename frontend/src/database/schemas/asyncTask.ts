/* eslint-disable sort-keys-fix/sort-keys-fix  */
import { integer, sqliteTable, text } from 'drizzle-orm/sqlite-core';
// Use browser-compatible crypto.randomUUID

import { timestamps } from './_helpers';
import { users } from './user';

export const asyncTasks = sqliteTable('async_tasks', {
  id: text('id').$defaultFn(() => crypto.randomUUID()).primaryKey(),
  type: text('type', { mode: 'json' }),

  status: text('status', { mode: 'json' }),
  error: text('error'),

  userId: text('user_id')
    .references(() => users.id, { onDelete: 'cascade' })
    .notNull(),
  duration: integer('duration'),

  ...timestamps,
});

export type NewAsyncTaskItem = typeof asyncTasks.$inferInsert;
export type AsyncTaskSelectItem = typeof asyncTasks.$inferSelect;
