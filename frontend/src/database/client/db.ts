import { sql } from 'drizzle-orm';
import { BaseSQLiteDatabase } from 'drizzle-orm/sqlite-core';
import { Md5 } from 'ts-md5';

import {
  ClientDBLoadingProgress,
  DatabaseLoadingState,
  MigrationSQL,
  MigrationTableItem,
} from '@/types/clientDB';
import { sleep } from '@/utils/sleep';

import migrations from '../core/migrations.json';
import { DrizzleMigrationModel } from '../models/drizzleMigration';
import * as schema from '../schemas';
import { WailsSQLiteDriver, createDrizzleWailsSQLite } from './wails-sqlite-driver';
import { initWailsSQLite } from './wails-sqlite';

const sqliteSchemaHashCache = 'VERIDIUM_SQLITE_SCHEMA_HASH';

const DB_NAME = 'veridium';
type DrizzleInstance = BaseSQLiteDatabase<'sync', any, typeof schema>;

interface onErrorState {
  error: Error;
  migrationTableItems: MigrationTableItem[];
  migrationsSQL: MigrationSQL[];
}

export interface DatabaseLoadingCallbacks {
  onError?: (error: onErrorState) => void;
  onProgress?: (progress: ClientDBLoadingProgress) => void;
  onStateChange?: (state: DatabaseLoadingState) => void;
}

export class DatabaseManager {
  private static instance: DatabaseManager;
  private dbInstance: DrizzleInstance | null = null;
  private driver: WailsSQLiteDriver | null = null;
  private initPromise: Promise<DrizzleInstance> | null = null;
  private callbacks?: DatabaseLoadingCallbacks;
  private isLocalDBSchemaSynced = false;

  private constructor() {}

  static getInstance() {
    if (!DatabaseManager.instance) {
      DatabaseManager.instance = new DatabaseManager();
    }
    return DatabaseManager.instance;
  }

  // 数据库迁移方法
  private async migrate(skipMultiRun = false): Promise<DrizzleInstance> {
    if (this.isLocalDBSchemaSynced && skipMultiRun) return this.db;

    let hash: string | undefined;
    if (typeof localStorage !== 'undefined') {
      const cacheHash = localStorage.getItem(sqliteSchemaHashCache);
      hash = Md5.hashStr(JSON.stringify(migrations));
      // if hash is the same, no need to migrate
      if (hash === cacheHash) {
        try {
          const drizzleMigration = new DrizzleMigrationModel(this.db as any);

          // 检查数据库中是否存在表
          const tableCount = await drizzleMigration.getTableCounts();

          // 如果表数量大于0，则认为数据库已正确初始化
          if (tableCount > 0) {
            this.isLocalDBSchemaSynced = true;
            return this.db;
          }
        } catch (error) {
          console.warn('Error checking table existence, proceeding with migration', error);
          // 如果查询失败，继续执行迁移以确保安全
        }
      }
    }

    const start = Date.now();
    try {
      this.callbacks?.onStateChange?.(DatabaseLoadingState.Migrating);

      // Apply migrations using the Wails SQLite driver
      if (this.driver && migrations) {
        for (const migration of migrations as any[]) {
          if (migration.sql) {
            // Execute each SQL statement in the migration
            const statements = migration.sql.split('--> statement-breakpoint');
            for (const statement of statements) {
              const trimmed = statement.trim();
              if (trimmed) {
                await this.driver.execute(trimmed);
              }
            }
          }
        }
      }

      if (typeof localStorage !== 'undefined' && hash) {
        localStorage.setItem(sqliteSchemaHashCache, hash);
      }

      this.isLocalDBSchemaSynced = true;

      console.info(`🗂 Migration success, take ${Date.now() - start}ms`);
    } catch (cause) {
      console.error('❌ Local database schema migration failed', cause);
      throw cause;
    }

    return this.db;
  }

  // 初始化数据库
  async initialize(callbacks?: DatabaseLoadingCallbacks): Promise<DrizzleInstance> {
    if (this.initPromise) return this.initPromise;

    this.callbacks = callbacks;

    this.initPromise = (async () => {
      try {
        if (this.dbInstance) return this.dbInstance;

        const time = Date.now();
        // 初始化数据库
        this.callbacks?.onStateChange?.(DatabaseLoadingState.Initializing);

        // Initialize Wails SQLite connection
        this.callbacks?.onProgress?.({
          phase: 'dependencies',
          progress: 50,
        });

        this.driver = await initWailsSQLite();

        this.callbacks?.onProgress?.({
          costTime: Date.now() - time,
          phase: 'dependencies',
          progress: 100,
        });

        // Create Drizzle instance with Wails SQLite driver
        const db = createDrizzleWailsSQLite(this.driver, schema);

        // Create BaseSQLiteDatabase instance
        this.dbInstance = new BaseSQLiteDatabase('sync', db.dialect as any, db.session as any, undefined);

        await this.migrate(true);

        this.callbacks?.onStateChange?.(DatabaseLoadingState.Finished);
        console.log(`✅ Database initialized in ${Date.now() - time}ms`);

        await sleep(50);

        this.callbacks?.onStateChange?.(DatabaseLoadingState.Ready);

        return this.dbInstance as DrizzleInstance;
      } catch (e) {
        this.initPromise = null;
        this.callbacks?.onStateChange?.(DatabaseLoadingState.Error);
        const error = e as Error;

        // 查询迁移表数据
        let migrationsTableData: MigrationTableItem[] = [];
        try {
          // 尝试查询迁移表
          const drizzleMigration = new DrizzleMigrationModel(this.db as any);
          migrationsTableData = await drizzleMigration.getMigrationList();
        } catch (queryError) {
          console.error('Failed to query migrations table:', queryError);
        }

        this.callbacks?.onError?.({
          error: {
            message: error.message,
            name: error.name,
            stack: error.stack,
          },
          migrationTableItems: migrationsTableData,
          migrationsSQL: migrations,
        });

        console.error(error);
        throw error;
      }
    })();

    return this.initPromise;
  }

  // 获取数据库实例
  get db(): DrizzleInstance {
    if (!this.dbInstance) {
      throw new Error('Database not initialized. Please call initialize() first.');
    }
    return this.dbInstance;
  }

  // 创建代理对象
  createProxy(): DrizzleInstance {
    return new Proxy({} as DrizzleInstance, {
      get: (target, prop) => {
        return this.db[prop as keyof DrizzleInstance];
      },
    });
  }

  async resetDatabase(): Promise<void> {
    // 1. Close the Wails SQLite connection
    if (this.driver) {
      try {
        const { closeWailsSQLite } = await import('./wails-sqlite');
        await closeWailsSQLite();
        console.log('Wails SQLite connection closed successfully.');
      } catch (e) {
        console.error('Error closing Wails SQLite connection:', e);
      }
    }

    // 2. Reset database instance and initialization state
    this.dbInstance = null;
    this.driver = null;
    this.initPromise = null;
    this.isLocalDBSchemaSynced = false;

    // 3. Clear schema hash cache
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem(sqliteSchemaHashCache);
    }

    console.log(`✅ Database '${DB_NAME}' reset successfully`);
  }
}

// 导出单例
const dbManager = DatabaseManager.getInstance();

// 保持原有的 clientDB 导出不变
export const clientDB = dbManager.createProxy();

// 导出初始化方法，供应用启动时使用
export const initializeDB = (callbacks?: DatabaseLoadingCallbacks) =>
  dbManager.initialize(callbacks);

export const resetClientDatabase = async () => {
  await dbManager.resetDatabase();
};

export const updateMigrationRecord = async (migrationHash: string) => {
  await clientDB.run(
    sql`INSERT INTO "drizzle"."__drizzle_migrations" ("hash", "created_at") VALUES (${migrationHash}, ${Date.now()});`,
  );

  await initializeDB();
};
