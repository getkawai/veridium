/* eslint-disable sort-keys-fix/sort-keys-fix  */
import {
  index,
  integer,
  primaryKey,
  sqliteTable,
  text,
  uniqueIndex,
} from 'drizzle-orm/sqlite-core';
import { createInsertSchema } from 'drizzle-zod';

import { LobeAgentChatConfig, LobeAgentTTSConfig } from '@/types/agent';

import { idGenerator, randomSlug } from '../utils/idGenerator';
import { timestamps } from './_helpers';
import { files, knowledgeBases } from './file';
import { users } from './user';

// Agent table is the main table for storing agents
// agent is a model that represents the assistant that is created by the user
// agent can have its own knowledge base and files

export const agents = sqliteTable(
  'agents',
  {
    id: text('id')
      .primaryKey()
      .$defaultFn(() => idGenerator('agents'))
      .notNull(),
    slug: text('slug')
      .$defaultFn(() => randomSlug(4))
      .unique(),
    title: text('title'),
    description: text('description'),
    tags: text('tags', { mode: 'json' }).$type<string[]>().$defaultFn(() => []),
    avatar: text('avatar'),
    backgroundColor: text('background_color'),

    plugins: text('plugins', { mode: 'json' }).$type<string[]>().$defaultFn(() => []),

    clientId: text('client_id'),

    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),

    chatConfig: text('chat_config', { mode: 'json' }).$type<LobeAgentChatConfig>(),

    fewShots: text('few_shots', { mode: 'json' }),
    model: text('model'),
    params: text('params', { mode: 'json' }).$defaultFn(() => ({})),
    provider: text('provider'),
    systemRole: text('system_role'),
    tts: text('tts', { mode: 'json' }).$type<LobeAgentTTSConfig>(),

    virtual: integer('virtual', { mode: 'boolean' }).default(false),

    openingMessage: text('opening_message'),
    openingQuestions: text('opening_questions', { mode: 'json' }).$type<string[]>().$defaultFn(() => []),

    ...timestamps,
  },
  (t) => ({
    clientIdUnique: uniqueIndex('client_id_user_id_unique').on(t.clientId, t.userId),
    titleIndex: index('agents_title_idx').on(t.title),
    descriptionIndex: index('agents_description_idx').on(t.description),
  }),
);

export const insertAgentSchema = createInsertSchema(agents);

export type NewAgent = typeof agents.$inferInsert;
export type AgentItem = typeof agents.$inferSelect;

export const agentsKnowledgeBases = sqliteTable(
  'agents_knowledge_bases',
  {
    agentId: text('agent_id')
      .references(() => agents.id, { onDelete: 'cascade' })
      .notNull(),
    knowledgeBaseId: text('knowledge_base_id')
      .references(() => knowledgeBases.id, { onDelete: 'cascade' })
      .notNull(),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
    enabled: integer('enabled', { mode: 'boolean' }).default(true),

    ...timestamps,
  },
  (t) => ({
    pk: primaryKey({ columns: [t.agentId, t.knowledgeBaseId] }),
  }),
);

export const agentsFiles = sqliteTable(
  'agents_files',
  {
    fileId: text('file_id')
      .notNull()
      .references(() => files.id, { onDelete: 'cascade' }),
    agentId: text('agent_id')
      .notNull()
      .references(() => agents.id, { onDelete: 'cascade' }),
    enabled: integer('enabled', { mode: 'boolean' }).default(true),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),

    ...timestamps,
  },
  (t) => ({
    pk: primaryKey({ columns: [t.fileId, t.agentId, t.userId] }),
  }),
);
