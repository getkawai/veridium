/* eslint-disable sort-keys-fix/sort-keys-fix  */
import {
  blob,
  index,
  integer,
  sqliteTable,
  text,
  uniqueIndex,
} from 'drizzle-orm/sqlite-core';
import { randomUUID } from 'crypto';

import { timestamps } from './_helpers';
import { files } from './file';
import { users } from './user';

export const chunks = sqliteTable(
  'chunks',
  {
    id: text('id')
      .primaryKey()
      .$defaultFn(() => randomUUID()),
    text: text('text'),
    abstract: text('abstract'),
    metadata: text('metadata', { mode: 'json' }),
    index: integer('index'),
    type: text('type'),

    clientId: text('client_id'),
    userId: text('user_id').references(() => users.id, { onDelete: 'cascade' }),

    ...timestamps,
  },
  (t) => [
    uniqueIndex('chunks_client_id_user_id_unique').on(t.clientId, t.userId),
    index('chunks_user_id_idx').on(t.userId),
  ],
);

export type NewChunkItem = typeof chunks.$inferInsert & { fileId?: string };

export const unstructuredChunks = sqliteTable(
  'unstructured_chunks',
  {
    id: text('id')
      .primaryKey()
      .$defaultFn(() => randomUUID()),
    text: text('text'),
    metadata: text('metadata', { mode: 'json' }),
    index: integer('index'),
    type: text('type'),

    ...timestamps,

    parentId: text('parent_id'),
    compositeId: text('composite_id').references(() => chunks.id, { onDelete: 'cascade' }),
    clientId: text('client_id'),
    userId: text('user_id').references(() => users.id, { onDelete: 'cascade' }),
    fileId: text('file_id').references(() => files.id, { onDelete: 'cascade' }),
  },
  (t) => ({
    clientIdUnique: uniqueIndex('unstructured_chunks_client_id_user_id_unique').on(
      t.clientId,
      t.userId,
    ),
  }),
);

export type NewUnstructuredChunkItem = typeof unstructuredChunks.$inferInsert;

// Store embeddings as blob (binary) or text (JSON array)
// We'll use blob for more efficient storage
export const embeddings = sqliteTable(
  'embeddings',
  {
    id: text('id')
      .primaryKey()
      .$defaultFn(() => randomUUID()),
    chunkId: text('chunk_id')
      .references(() => chunks.id, { onDelete: 'cascade' })
      .unique(),
    // Store as blob (Float32Array) or as JSON text
    embeddings: blob('embeddings', { mode: 'buffer' }),
    model: text('model'),
    clientId: text('client_id'),
    userId: text('user_id').references(() => users.id, { onDelete: 'cascade' }),
  },
  (t) => [
    uniqueIndex('embeddings_client_id_user_id_unique').on(t.clientId, t.userId),
    // improve delete embeddings query
    index('embeddings_chunk_id_idx').on(t.chunkId),
  ],
);

export type NewEmbeddingsItem = typeof embeddings.$inferInsert;
export type EmbeddingsSelectItem = typeof embeddings.$inferSelect;
