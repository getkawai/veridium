import { integer, text } from 'drizzle-orm/sqlite-core';

// SQLite stores timestamps as integers (Unix timestamp in milliseconds)
export const timestamptz = (name: string) =>
  integer(name, { mode: 'timestamp_ms' }).notNull().default(new Date());

export const varchar255 = (name: string) => text(name);

export const createdAt = () =>
  integer('created_at', { mode: 'timestamp_ms' })
    .notNull()
    .$defaultFn(() => new Date());

export const updatedAt = () =>
  integer('updated_at', { mode: 'timestamp_ms' })
    .notNull()
    .$defaultFn(() => new Date())
    .$onUpdateFn(() => new Date());

export const accessedAt = () =>
  integer('accessed_at', { mode: 'timestamp_ms' })
    .notNull()
    .$defaultFn(() => new Date());

// columns.helpers.ts
export const timestamps = {
  accessedAt: accessedAt(),
  createdAt: createdAt(),
  updatedAt: updatedAt(),
};
