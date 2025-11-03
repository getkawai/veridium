import pMap from 'p-map';

import {
  FilterCondition,
  PaginationParams,
  TableBasicInfo,
  TableColumnInfo,
} from '@/types/tableViewer';

import { executeQuery, executeQueryOne, executeCommand } from '../../utils/sqlExecutor';

export class TableViewerRepo {
  constructor(_db: any, _userId: string) {
    // No longer needed, kept for API compatibility
  }

  /**
   * Get all tables in the database
   */
  async getAllTables(): Promise<TableBasicInfo[]> {
    const sqlString = `
      SELECT
        name,
        type
      FROM sqlite_master
      WHERE type IN ('table', 'view')
        AND name NOT LIKE 'sqlite_%'
      ORDER BY name;
    `;

    const tables = await executeQuery<{ name: string; type: string }>(sqlString);
    const tableNames = tables.map((row) => row.name);

    const counts = await pMap(tableNames, async (name) => this.getTableCount(name), {
      concurrency: 10,
    });

    return tables.map((row, index) => ({
      count: counts[index],
      name: row.name,
      type: row.type as 'BASE TABLE' | 'VIEW',
    }));
  }

  /**
   * Get detailed structure info for a table
   */
  async getTableDetails(tableName: string): Promise<TableColumnInfo[]> {
    const sqlString = `PRAGMA table_info(${tableName})`;
    const columns = await executeQuery<any>(sqlString);

    return columns.map((col) => ({
      defaultValue: col.dflt_value,
      foreignKey: undefined, // SQLite doesn't provide foreign key info in PRAGMA table_info
      isPrimaryKey: !!col.pk,
      name: col.name,
      nullable: !col.notnull,
      type: col.type,
    }));
  }

  /**
   * Get table data with pagination, sorting, and filtering
   */
  async getTableData(tableName: string, pagination: PaginationParams, filters?: FilterCondition[]) {
    const offset = (pagination.page - 1) * pagination.pageSize;

    // Build base query
    let selectClause = `SELECT * FROM ${tableName}`;
    let whereClause = '';
    let orderClause = '';
    const params: any[] = [];

    // Add filters
    if (filters && filters.length > 0) {
      const whereConditions: string[] = [];
      filters.forEach((filter) => {
        switch (filter.operator) {
          case 'equals': {
            whereConditions.push(`${filter.column} = ?`);
            params.push(filter.value);
            break;
          }
          case 'contains': {
            whereConditions.push(`UPPER(${filter.column}) LIKE UPPER(?)`);
            params.push(`%${filter.value}%`);
            break;
          }
          case 'startsWith': {
            whereConditions.push(`UPPER(${filter.column}) LIKE UPPER(?)`);
            params.push(`${filter.value}%`);
            break;
          }
          case 'endsWith': {
            whereConditions.push(`UPPER(${filter.column}) LIKE UPPER(?)`);
            params.push(`%${filter.value}`);
            break;
          }
        }
      });

      if (whereConditions.length > 0) {
        whereClause = ` WHERE ${whereConditions.join(' AND ')}`;
      }
    }

    // Add sorting
    if (pagination.sortBy) {
      const direction = pagination.sortOrder === 'desc' ? 'DESC' : 'ASC';
      orderClause = ` ORDER BY ${pagination.sortBy} ${direction}`;
    }

    // Add pagination
    const limitClause = ` LIMIT ? OFFSET ?`;
    params.push(pagination.pageSize, offset);

    const query = selectClause + whereClause + orderClause + limitClause;

    // Get total count
    let countQuery = `SELECT COUNT(*) as total FROM ${tableName}`;
    let countParams: any[] = [];
    if (whereClause) {
      countQuery += whereClause;
      countParams = params.slice(0, -2); // Remove LIMIT and OFFSET params
    }

    // Execute queries
    const [data, countResult] = await Promise.all([
      executeQuery(query, params),
      executeQueryOne<{ total: number }>(countQuery, countParams),
    ]);

    return {
      data,
      pagination: {
        page: pagination.page,
        pageSize: pagination.pageSize,
        total: Number(countResult?.total || 0),
      },
    };
  }

  /**
   * Update a row in the table
   */
  async updateRow(
    tableName: string,
    id: string,
    primaryKeyColumn: string,
    data: Record<string, any>,
  ) {
    const setParts = Object.keys(data).map((key) => `${key} = ?`);
    const values = Object.values(data);

    const sqlString = `UPDATE ${tableName} SET ${setParts.join(', ')} WHERE ${primaryKeyColumn} = ?`;
    values.push(id);

    await executeCommand(sqlString, values);

    // Get the updated row
    const getRow = `SELECT * FROM ${tableName} WHERE ${primaryKeyColumn} = ?`;
    return executeQueryOne(getRow, [id]);
  }

  /**
   * Delete a row from the table
   */
  async deleteRow(tableName: string, id: string, primaryKeyColumn: string) {
    const sqlString = `DELETE FROM ${tableName} WHERE ${primaryKeyColumn} = ?`;
    await executeCommand(sqlString, [id]);
  }

  /**
   * Insert a new row
   */
  async insertRow(tableName: string, data: Record<string, any>) {
    const columns = Object.keys(data);
    const placeholders = columns.map(() => '?').join(', ');
    const values = Object.values(data);

    const sqlString = `INSERT INTO ${tableName} (${columns.join(', ')}) VALUES (${placeholders})`;

    await executeCommand(sqlString, values);

    // For SQLite, get the last inserted row
    const getLastRow = `SELECT * FROM ${tableName} WHERE rowid = last_insert_rowid()`;
    return executeQueryOne(getLastRow);
  }

  /**
   * Get total count for a table
   */
  async getTableCount(tableName: string): Promise<number> {
    const sqlString = `SELECT COUNT(*) as total FROM ${tableName}`;
    const result = await executeQueryOne<{ total: number }>(sqlString);
    return Number(result?.total || 0);
  }

  /**
   * Batch delete rows
   */
  async batchDelete(tableName: string, ids: string[], primaryKeyColumn: string) {
    const inClause = ids.map(() => '?').join(', ');
    const sqlString = `DELETE FROM ${tableName} WHERE ${primaryKeyColumn} IN (${inClause})`;

    await executeCommand(sqlString, ids);
  }

  /**
   * Export table data (supports paginated export)
   */
  async exportTableData(
    tableName: string,
    pagination?: PaginationParams,
    filters?: FilterCondition[],
  ) {
    return this.getTableData(tableName, pagination || { page: 1, pageSize: 1000 }, filters);
  }
}
