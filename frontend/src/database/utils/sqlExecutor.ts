/**
 * SQL Executor - Raw SQL execution via Wails SQLite bindings
 * Replaces Drizzle's db.execute() for repositories
 */

import { Query as WailsQuery, Execute as WailsExecute } from '@@/github.com/wailsapp/wails/v3/pkg/services/sqlite/sqliteservice';

/**
 * Execute a SELECT query and return all rows
 */
export async function executeQuery<T = any>(sql: string, params: any[] = []): Promise<T[]> {
  try {
    const result = await WailsQuery(sql, ...params);
    return (result || []) as T[];
  } catch (error) {
    console.error('SQL Query Error:', sql, params, error);
    throw error;
  }
}

/**
 * Execute a query and return the first row
 */
export async function executeQueryOne<T = any>(sql: string, params: any[] = []): Promise<T | undefined> {
  const rows = await executeQuery<T>(sql, params);
  return rows.length > 0 ? rows[0] : undefined;
}

/**
 * Execute a non-SELECT query (INSERT, UPDATE, DELETE)
 */
export async function executeCommand(sql: string, params: any[] = []): Promise<void> {
  try {
    await WailsExecute(sql, ...params);
  } catch (error) {
    console.error('SQL Command Error:', sql, params, error);
    throw error;
  }
}

/**
 * Execute a batch of commands in sequence
 */
export async function executeBatch(commands: Array<{ sql: string; params?: any[] }>): Promise<void> {
  for (const cmd of commands) {
    await executeCommand(cmd.sql, cmd.params || []);
  }
}

/**
 * Get count from a table
 */
export async function getCount(tableName: string, where?: string, params: any[] = []): Promise<number> {
  const sql = `SELECT COUNT(*) as total FROM ${tableName}${where ? ` WHERE ${where}` : ''}`;
  const result = await executeQueryOne<{ total: number }>(sql, params);
  return result?.total || 0;
}




