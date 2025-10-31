import type { BaseSQLiteDatabase } from 'drizzle-orm/sqlite-core';

import * as schema from './schemas';

export type LobeChatDatabaseSchema = typeof schema;

export type LobeChatDatabase = BaseSQLiteDatabase<'sync', any, LobeChatDatabaseSchema> & {
  execute: (query: string, params?: any[]) => Promise<{ rows: any[]; rowCount?: number }>;
};

export type Transaction = Parameters<Parameters<LobeChatDatabase['transaction']>[0]>[0];
