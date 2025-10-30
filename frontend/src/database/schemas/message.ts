/* eslint-disable sort-keys-fix/sort-keys-fix  */
import { GroundingSearch, ModelReasoning } from  '@/types';
import {
  boolean,
  index,
  jsonb,
  numeric,
  pgTable,
  primaryKey,
  text,
  uniqueIndex,
  uuid,
} from 'drizzle-orm/sqlite-core';
import { createInsertSchema, createSelectSchema } from 'drizzle-zod';

import { idGenerator } from '../utils/idGenerator';
import { timestamps, varchar255 } from './_helpers';
import { agents } from './agent';
import { chatGroups } from './chatGroup';
import { files } from './file';
import { chunks, embeddings } from './rag';
import { sessions } from './session';
import { threads, topics } from './topic';
import { users } from './user';

/**
 * Message groups table for multi-models parallel conversations
 * Allows multiple AI models to respond to the same user message in parallel
 */
// @ts-ignore
export const messageGroups = sqliteTable(
  'message_groups',
  {
    id: varchar255('id')
      .primaryKey()
      .$defaultFn(() => idGenerator('messageGroups'))
      .notNull(),

    // 关联关系 - 只需要 topic 层级
    topicId: text('topic_id').references(() => topics.id, { onDelete: 'cascade' }),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),

    // 支持嵌套结构
    // @ts-ignore
    parentGroupId: varchar255('parent_group_id').references(() => messageGroups.id, {
      onDelete: 'cascade',
    }),

    // 关联的用户消息
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    parentMessageId: text('parent_message_id').references(() => messages.id, {
      onDelete: 'cascade',
    }),

    // 元数据
    title: varchar255('title'),
    description: text('description'),

    clientId: varchar255('client_id'),

    ...timestamps,
  },
  (t) => [uniqueIndex('message_groups_client_id_user_id_unique').on(t.clientId, t.userId)],
);

export const insertMessageGroupSchema = createInsertSchema(messageGroups);

export type NewMessageGroup = typeof messageGroups.$inferInsert;
export type MessageGroupItem = typeof messageGroups.$inferSelect;

// @ts-ignore
export const messages = sqliteTable(
  'messages',
  {
    id: text('id')
      .$defaultFn(() => idGenerator('messages'))
      .primaryKey(),

    role: varchar255('role').notNull(),
    content: text('content', { mode: 'json' }),
    reasoning: text('reasoning').$type<ModelReasoning>(),
    search: text('search').$type<GroundingSearch>(),
    metadata: text('metadata'),

    model: text('model', { mode: 'json' }),
    provider: text('provider', { mode: 'json' }),

    favorite: integer('favorite').default(false),
    error: text('error'),

    tools: text('tools'),

    traceId: text('trace_id', { mode: 'json' }),
    observationId: text('observation_id', { mode: 'json' }),

    clientId: text('client_id'),

    // foreign keys
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
    sessionId: text('session_id').references(() => sessions.id, { onDelete: 'cascade' }),
    topicId: text('topic_id').references(() => topics.id, { onDelete: 'cascade' }),
    threadId: text('thread_id').references(() => threads.id, { onDelete: 'cascade' }),
    // @ts-ignore
    parentId: text('parent_id').references(() => messages.id, { onDelete: 'set null' }),
    quotaId: text('quota_id').references(() => messages.id, { onDelete: 'set null' }),

    // used for group chat
    agentId: text('agent_id').references(() => agents.id, { onDelete: 'set null' }),
    groupId: text('group_id').references(() => chatGroups.id, { onDelete: 'set null' }),
    // targetId can be an agent ID, "user", or null - no FK constraint
    targetId: text('target_id'),

    // used for multi-models parallel
    messageGroupId: varchar255('message_group_id').references(() => messageGroups.id, {
      onDelete: 'cascade',
    }),
    ...timestamps,
  },
  (table) => [
    index('messages_created_at_idx').on(table.createdAt),
    uniqueIndex('message_client_id_user_unique').on(table.clientId, table.userId),
    index('messages_topic_id_idx').on(table.topicId),
    index('messages_parent_id_idx').on(table.parentId),
    index('messages_quota_id_idx').on(table.quotaId),

    index('messages_user_id_idx').on(table.userId),
    index('messages_session_id_idx').on(table.sessionId),
    index('messages_thread_id_idx').on(table.threadId),
  ],
);

// if the message container a plugin
export const messagePlugins = sqliteTable(
  'message_plugins',
  {
    id: text('id')
      .references(() => messages.id, { onDelete: 'cascade' })
      .primaryKey(),

    toolCallId: text('tool_call_id'),
    type: text('type', {
      enum: ['default', 'markdown', 'standalone', 'builtin'],
    }).default('default'),

    apiName: text('api_name'),
    arguments: text('arguments', { mode: 'json' }),
    identifier: text('identifier', { mode: 'json' }),
    state: text('state'),
    error: text('error'),
    clientId: text('client_id'),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
  },
  (t) => ({
    clientIdUnique: uniqueIndex('message_plugins_client_id_user_id_unique').on(
      t.clientId,
      t.userId,
    ),
  }),
);

export type MessagePluginItem = typeof messagePlugins.$inferSelect;
export const updateMessagePluginSchema = createSelectSchema(messagePlugins);

export const messageTTS = sqliteTable(
  'message_tts',
  {
    id: text('id')
      .references(() => messages.id, { onDelete: 'cascade' })
      .primaryKey(),
    contentMd5: text('content_md5'),
    fileId: text('file_id').references(() => files.id, { onDelete: 'cascade' }),
    voice: text('voice'),
    clientId: text('client_id'),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
  },
  (t) => ({
    clientIdUnique: uniqueIndex('message_tts_client_id_user_id_unique').on(t.clientId, t.userId),
  }),
);

export const messageTranslates = sqliteTable(
  'message_translates',
  {
    id: text('id')
      .references(() => messages.id, { onDelete: 'cascade' })
      .primaryKey(),
    content: text('content', { mode: 'json' }),
    from: text('from'),
    to: text('to'),
    clientId: text('client_id'),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
  },
  (t) => ({
    clientIdUnique: uniqueIndex('message_translates_client_id_user_id_unique').on(
      t.clientId,
      t.userId,
    ),
  }),
);

// if the message contains a file
// save the file id and message id
export const messagesFiles = sqliteTable(
  'messages_files',
  {
    fileId: text('file_id')
      .notNull()
      .references(() => files.id, { onDelete: 'cascade' }),
    messageId: text('message_id')
      .notNull()
      .references(() => messages.id, { onDelete: 'cascade' }),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
  },
  (t) => ({
    pk: primaryKey({ columns: [t.fileId, t.messageId] }),
  }),
);

export const messageQueries = sqliteTable(
  'message_queries',
  {
    id: text('id').$defaultFn(() => randomUUID()).primaryKey(),
    messageId: text('message_id')
      .references(() => messages.id, { onDelete: 'cascade' })
      .notNull(),
    rewriteQuery: text('rewrite_query'),
    userQuery: text('user_query'),
    clientId: text('client_id'),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
    embeddingsId: text('embeddings_id').references(() => embeddings.id, { onDelete: 'set null' }),
  },
  (t) => ({
    clientIdUnique: uniqueIndex('message_queries_client_id_user_id_unique').on(
      t.clientId,
      t.userId,
    ),
  }),
);

export type NewMessageQuery = typeof messageQueries.$inferInsert;

export const messageQueryChunks = sqliteTable(
  'message_query_chunks',
  {
    messageId: text('id').references(() => messages.id, { onDelete: 'cascade' }),
    queryId: text('query_id').references(() => messageQueries.id, { onDelete: 'cascade' }),
    chunkId: text('chunk_id').references(() => chunks.id, { onDelete: 'cascade' }),
    similarity: real('similarity', { precision: 6, scale: 5 }),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
  },
  (t) => ({
    pk: primaryKey({ columns: [t.chunkId, t.messageId, t.queryId] }),
  }),
);
export type NewMessageFileChunk = typeof messageQueryChunks.$inferInsert;

// convert message content to the chunks
// then we can use message as the RAG source
export const messageChunks = sqliteTable(
  'message_chunks',
  {
    messageId: text('message_id').references(() => messages.id, { onDelete: 'cascade' }),
    chunkId: text('chunk_id').references(() => chunks.id, { onDelete: 'cascade' }),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
  },
  (t) => ({
    pk: primaryKey({ columns: [t.chunkId, t.messageId] }),
  }),
);
