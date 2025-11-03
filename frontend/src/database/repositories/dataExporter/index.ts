import pMap from 'p-map';

import { executeQuery } from '../../utils/sqlExecutor';

interface BaseTableConfig {
  table: string;
  type: 'base';
  userField?: string;
}

export interface RelationTableConfig {
  relations: {
    field: string;
    sourceField?: string;
    sourceTable: string;
  }[];
  table: string;
  type: 'relation';
}

export const DATA_EXPORT_CONFIG = {
  baseTables: [
    // { table: 'users', userField: 'id' },
    { table: 'user_settings', userField: 'id' },
    { table: 'user_installed_plugins' },
    { table: 'agents' },
    // { table: 'agents_files' },
    // { table: 'agents_knowledge_bases' },
    // { table: 'agents_to_sessions' },
    { table: 'ai_models' },
    { table: 'ai_providers' },
    // async tasks should not be included
    // { table: 'async_tasks' },
    // { table: 'chunks' },
    // { table: 'unstructured_chunks' },
    // { table: 'embeddings' },
    // { table: 'files' },
    // { table: 'file_chunks' },
    // { table: 'files_to_sessions' },
    // { table: 'knowledge_bases' },
    // { table: 'knowledge_base_files' },
    { table: 'message_chunks' },
    { table: 'message_plugins' },
    // { table: 'message_query_chunks' },
    // { table: 'message_queries' },
    { table: 'message_translates' },
    // { table: 'message_tts' },
    { table: 'messages' },
    // { table: 'messages_files' },

    // next auth tables won't be included
    // { table: 'nextauth_accounts' },
    // { table: 'nextauth_sessions' },
    // { table: 'nextauth_authenticators' },
    // { table: 'nextauth_verification_tokens' },
    { table: 'session_groups' },
    { table: 'sessions' },
    { table: 'threads' },
    { table: 'topics' },
  ] as BaseTableConfig[],
  relationTables: [
    // {
    //   relations: [{ field: 'hash_id', sourceField: 'file_hash', sourceTable: 'files' }],
    //   table: 'global_files',
    // },
    {
      relations: [
        { field: 'agent_id', sourceField: 'id', sourceTable: 'agents' },
        { field: 'session_id', sourceField: 'id', sourceTable: 'sessions' },
      ],
      table: 'agents_to_sessions',
    },

    // {
    //   relations: [{ field: 'id', sourceField: 'id', sourceTable: 'messages' }],
    //   table: 'message_plugins',
    // },
  ] as RelationTableConfig[],
};

export class DataExporterRepos {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  private removeUserId(data: any[]) {
    return data.map((item) => {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { user_id: _, userId: __, ...rest } = item;
      return rest;
    });
  }

  private async queryTable(config: RelationTableConfig, existingData: Record<string, any[]>) {
    const { table } = config;

    try {
      const whereClauses: string[] = [];
      const params: any[] = [];

      // Build WHERE clauses for each relation
      for (const relation of config.relations) {
        const sourceData = existingData[relation.sourceTable] || [];

        // If source data is empty, this table won't have any data
        if (sourceData.length === 0) {
          console.log(
            `Source table ${relation.sourceTable} has no data, skipping query for ${table}`,
          );
          return [];
        }

        const sourceIds = sourceData.map((item) => item[relation.sourceField || 'id']);
        
        // Build IN clause
        const placeholders = sourceIds.map(() => '?').join(', ');
        whereClauses.push(`${relation.field} IN (${placeholders})`);
        params.push(...sourceIds);
      }

      // Build final query
      const whereClause = whereClauses.length > 0 ? `WHERE ${whereClauses.join(' AND ')}` : '';
      const sql = `SELECT * FROM ${table} ${whereClause}`;

      const result = await executeQuery(sql, params);

      console.log(`Successfully exported table: ${table}, count: ${result.length}`);
      return config.relations ? result : this.removeUserId(result);
    } catch (error) {
      console.error(`Error querying table ${table}:`, error);
      return [];
    }
  }

  private async queryBaseTables(config: BaseTableConfig) {
    const { table } = config;

    try {
      // Use userId or custom userField
      const userField = config.userField || 'user_id';
      const sql = `SELECT * FROM ${table} WHERE ${userField} = ?`;

      const result = await executeQuery(sql, [this.userId]);

      console.log(`Successfully exported table: ${table}, count: ${result.length}`);
      return this.removeUserId(result);
    } catch (error) {
      console.error(`Error querying table ${table}:`, error);
      return [];
    }
  }

  async export(concurrency = 10) {
    const result: Record<string, any[]> = {};

    // 1. Query all base tables in parallel
    console.log('Querying base tables...');
    const baseResults = await pMap(
      DATA_EXPORT_CONFIG.baseTables,
      async (config) => ({ data: await this.queryBaseTables(config), table: config.table }),
      { concurrency },
    );

    // Update result set
    baseResults.forEach(({ table, data }) => {
      result[table] = data;
    });

    // 2. Query all relation tables in parallel
    const relationResults = await pMap(
      DATA_EXPORT_CONFIG.relationTables,
      async (config) => {
        // Check if all source tables have data
        const allSourcesHaveData = config.relations.every(
          (relation) => (result[relation.sourceTable] || []).length > 0,
        );

        if (!allSourcesHaveData) {
          console.log(`Skipping table ${config.table} as some source tables have no data`);
          return { data: [], table: config.table };
        }

        return {
          data: await this.queryTable(config, result),
          table: config.table,
        };
      },
      { concurrency },
    );

    // Update result set
    relationResults.forEach(({ table, data }) => {
      result[table] = data;
    });

    console.log('finalResults:', result);

    return result;
  }
}
