import { Query as WailsQuery, Execute as WailsExecute } from '@/bindings/github.com/wailsapp/wails/v3/pkg/services/sqlite/sqliteservice';
import type { Row, Rows } from '@/bindings/github.com/wailsapp/wails/v3/pkg/services/sqlite/models';

/**
 * Custom Drizzle driver for Wails SQLite
 * Wraps the Wails SQLite service bindings to work with Drizzle ORM
 */

export interface WailsSQLiteDriverConfig {
  onQuery?: (query: string, params: any[]) => void;
}

export class WailsSQLiteDriver {
  constructor(private config?: WailsSQLiteDriverConfig) {}

  async query<T = Rows>(query: string, params?: any[]): Promise<T> {
    this.config?.onQuery?.(query, params || []);
    
    try {
      const result = await WailsQuery(query, ...(params || []));
      return result as T;
    } catch (error) {
      console.error('Wails SQLite query error:', error);
      throw error;
    }
  }

  async execute(query: string, params?: any[]): Promise<{ rows: Row[]; rowCount?: number }> {
    this.config?.onQuery?.(query, params || []);

    try {
      // For SELECT queries, use query method
      if (query.trim().toUpperCase().startsWith('SELECT')) {
        const rows = await this.query<Rows>(query, params);
        return { rows: rows || [], rowCount: rows?.length || 0 };
      } else {
        // For non-SELECT queries, use execute method
        await WailsExecute(query, ...(params || []));
        return { rows: [], rowCount: undefined };
      }
    } catch (error) {
      console.error('Wails SQLite execute error:', error);
      throw error;
    }
  }

  async run(query: string, params?: any[]): Promise<void> {
    await this.execute(query, params);
  }

  async all(query: string, params?: any[]): Promise<Row[]> {
    return this.query<Rows>(query, params);
  }

  async get<T = Row>(query: string, params?: any[]): Promise<T | undefined> {
    const rows = await this.query<Rows>(query, params);
    return (rows.length > 0 ? rows[0] : undefined) as T | undefined;
  }

  async values<T = any[]>(query: string, params?: any[]): Promise<T[]> {
    const rows = await this.query<Rows>(query, params);
    return rows.map(row => Object.values(row)) as T[];
  }
}

/**
 * Drizzle session implementation for Wails SQLite
 */
export class WailsSQLiteSession {
  constructor(
    private driver: WailsSQLiteDriver,
    _schema: Record<string, unknown>,
  ) {}

  prepareQuery(query: string, fields?: any[], params?: any[], customResultMapper?: any) {
    return {
      execute: async (params?: any[]) => {
        const rows = await this.driver.query(query, params);
        return { rows };
      },
      all: async (params?: any[]) => {
        return this.driver.all(query, params);
      },
      get: async (params?: any[]) => {
        return this.driver.get(query, params);
      },
      values: async (params?: any[]) => {
        return this.driver.values(query, params);
      },
      run: async (params?: any[]) => {
        return this.driver.run(query, params);
      },
    };
  }

  prepareOneTimeQuery(query: string, fields?: any[], params?: any[], customResultMapper?: any) {
    return this.prepareQuery(query, fields, params, customResultMapper);
  }

  async run(query: string, params?: any[]): Promise<any> {
    return this.driver.run(query, params);
  }

  async all(query: string, params?: any[]): Promise<Row[]> {
    return this.driver.all(query, params);
  }

  async get<T = Row>(query: string, params?: any[]): Promise<T | undefined> {
    return this.driver.get<T>(query, params);
  }

  async values<T = any[]>(query: string, params?: any[]): Promise<T[]> {
    return this.driver.values<T>(query, params);
  }

  async count(query: string, params?: any[]): Promise<number> {
    const result = await this.get(query, params);
    if (result && typeof result === 'object') {
      const values = Object.values(result);
      return typeof values[0] === 'number' ? values[0] : 0;
    }
    return 0;
  }

  async transaction<T>(
    transaction: (tx: WailsSQLiteSession) => Promise<T>,
  ): Promise<T> {
    await this.driver.execute('BEGIN');

    try {
      const result = await transaction(this);
      await this.driver.execute('COMMIT');
      return result;
    } catch (error) {
      await this.driver.execute('ROLLBACK');
      throw error;
    }
  }

  get dialect() {
    return {
      migrate: async (migrations: any, session: any, config: any) => {
        // Migration logic will be handled separately
        console.log('Running migrations through dialect...');
        
        for (const migration of migrations) {
          const queries = migration.sql || [];
          for (const query of queries) {
            await this.driver.execute(query);
          }
        }
      },
    };
  }
}

/**
 * Create a Drizzle-compatible database instance
 */
export function createDrizzleWailsSQLite<TSchema extends Record<string, unknown>>(
  driver: WailsSQLiteDriver,
  schema: TSchema,
) {
  const session = new WailsSQLiteSession(driver, schema);

  return {
    _: {
      schema,
    },
    query: schema,
    session,
    dialect: session.dialect,
    
    // Drizzle ORM methods
    select: (fields?: any) => {
      // This will be implemented by Drizzle's query builder
      throw new Error('Use drizzle() function to create proper instance');
    },
    
    insert: (table: any) => {
      throw new Error('Use drizzle() function to create proper instance');
    },
    
    update: (table: any) => {
      throw new Error('Use drizzle() function to create proper instance');
    },
    
    delete: (table: any) => {
      throw new Error('Use drizzle() function to create proper instance');
    },
    
    execute: async (query: any) => {
      if (typeof query === 'string') {
        return driver.execute(query);
      }
      // Handle SQL query objects
      throw new Error('Complex query execution not yet implemented');
    },
    
    transaction: <T>(
      transaction: (tx: any) => Promise<T>,
    ): Promise<T> => {
      return session.transaction(transaction);
    },
    
    $with: (name: string) => {
      throw new Error('CTE not yet implemented');
    },
  };
}
