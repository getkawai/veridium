import pMap from 'p-map';

import {
  FilterCondition,
  PaginationParams,
  TableBasicInfo,
  TableColumnInfo,
} from '@/types/tableViewer';

import { LobeChatDatabase } from '../../type';

export class TableViewerRepo {
  private db: LobeChatDatabase;

  constructor(db: LobeChatDatabase, _userId: string) {
    this.db = db;
  }

  /**
   * 获取数据库中所有的表
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

    const tables = await this.db.execute(sqlString);

    const tableNames = tables.rows.map((row) => row.name) as string[];

    const counts = await pMap(tableNames, async (name) => this.getTableCount(name), {
      concurrency: 10,
    });

    return tables.rows.map((row, index) => ({
      count: counts[index],
      name: row.name,
      type: row.type,
    })) as TableBasicInfo[];
  }

  /**
   * 获取指定表的详细结构信息
   */
  async getTableDetails(tableName: string): Promise<TableColumnInfo[]> {
    const sqlString = `PRAGMA table_info(${tableName})`;

    const columns = await this.db.execute(sqlString);

    return columns.rows.map((col: any) => ({
      defaultValue: col.dflt_value,
      foreignKey: undefined, // SQLite doesn't provide foreign key info in PRAGMA table_info
      isPrimaryKey: !!col.pk,
      name: col.name,
      nullable: !col.notnull,
      type: col.type,
    }));
  }

  /**
   * 获取表数据，支持分页、排序和筛选
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
    const [data, count] = await Promise.all([
      this.db.execute(query, params),
      this.db.execute(countQuery, countParams)
    ]);

    return {
      data: data.rows,
      pagination: {
        page: pagination.page,
        pageSize: pagination.pageSize,
        total: Number(count.rows[0].total),
      },
    };
  }

  /**
   * 更新表中的一行数据
   */
  async updateRow(
    tableName: string,
    id: string,
    primaryKeyColumn: string,
    data: Record<string, any>,
  ) {
    const setParts = Object.keys(data).map(key => `${key} = ?`);
    const values = Object.values(data);

    const sqlString = `UPDATE ${tableName} SET ${setParts.join(', ')} WHERE ${primaryKeyColumn} = ?`;
    values.push(id);

    await this.db.execute(sqlString, values);

    // Get the updated row
    const getRow = `SELECT * FROM ${tableName} WHERE ${primaryKeyColumn} = ?`;
    const result = await this.db.execute(getRow, [id]);
    return result.rows[0];
  }

  /**
   * 删除表中的一行数据
   */
  async deleteRow(tableName: string, id: string, primaryKeyColumn: string) {
    const sqlString = `DELETE FROM ${tableName} WHERE ${primaryKeyColumn} = ?`;
    await this.db.execute(sqlString, [id]);
  }

  /**
   * 插入新行数据
   */
  async insertRow(tableName: string, data: Record<string, any>) {
    const columns = Object.keys(data);
    const placeholders = columns.map(() => '?').join(', ');
    const values = Object.values(data);

    const sqlString = `INSERT INTO ${tableName} (${columns.join(', ')}) VALUES (${placeholders})`;

    await this.db.execute(sqlString, values);

    // For SQLite, we need to get the last inserted row differently
    const getLastRow = `SELECT * FROM ${tableName} WHERE rowid = last_insert_rowid()`;
    const result = await this.db.execute(getLastRow);
    return result.rows[0];
  }

  /**
   * 获取表的总记录数
   */
  async getTableCount(tableName: string): Promise<number> {
    const sqlString = `SELECT COUNT(*) as total FROM ${tableName}`;
    const result = await this.db.execute(sqlString);
    return Number(result.rows[0].total);
  }

  /**
   * 批量删除数据
   */
  async batchDelete(tableName: string, ids: string[], primaryKeyColumn: string) {
    const inClause = ids.map(() => '?').join(', ');
    const sqlString = `DELETE FROM ${tableName} WHERE ${primaryKeyColumn} IN (${inClause})`;

    await this.db.execute(sqlString, ids);
  }

  /**
   * 导出表数据（支持分页导出）
   */
  async exportTableData(
    tableName: string,
    pagination?: PaginationParams,
    filters?: FilterCondition[],
  ) {
    return this.getTableData(tableName, pagination || { page: 1, pageSize: 1000 }, filters);
  }
}
