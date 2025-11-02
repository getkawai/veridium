import { MigrationTableItem } from  '@/types';
import { sql } from 'drizzle-orm';

import { LobeChatDatabase } from '../type';

export class DrizzleMigrationModel {
  private db: LobeChatDatabase;

  constructor(db: LobeChatDatabase) {
    this.db = db;
  }

  getTableCounts = async () => {
    // 使用 SQLite 兼容的方式查询用户表数量
    const result = await this.db.execute(
      sql`
        SELECT COUNT(*) as table_count
        FROM sqlite_master
        WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
      `,
    );

    return parseInt((result.rows[0] as any).table_count || '0');
  };

  getMigrationList = async () => {
    const res = await this.db.execute(
      'SELECT * FROM __drizzle_migrations ORDER BY created_at DESC;',
    );

    return res.rows as unknown as MigrationTableItem[];
  };
  getLatestMigrationHash = async () => {
    const res = await this.getMigrationList();

    return res[0].hash;
  };
}
