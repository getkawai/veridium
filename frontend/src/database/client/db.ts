/**
 * Database Client - Wails Bindings Only
 * Database initialization is handled by Go backend.
 * This file only provides connection status and callbacks.
 */

import {
  ClientDBLoadingProgress,
  DatabaseLoadingState,
} from '@/types/clientDB';
import { sleep } from '@/utils/sleep';

// Import Wails bindings check
import { DB } from '@/types/database';

const DB_NAME = 'veridium';

interface onErrorState {
  error: Error;
  message: string;
}

export interface DatabaseLoadingCallbacks {
  onError?: (error: onErrorState) => void;
  onProgress?: (progress: ClientDBLoadingProgress) => void;
  onStateChange?: (state: DatabaseLoadingState) => void;
}

export class DatabaseManager {
  private static instance: DatabaseManager;
  private isInitialized = false;
  private initPromise: Promise<void> | null = null;
  private callbacks?: DatabaseLoadingCallbacks;

  private constructor() {}

  static getInstance() {
    if (!DatabaseManager.instance) {
      DatabaseManager.instance = new DatabaseManager();
    }
    return DatabaseManager.instance;
  }

  /**
   * Initialize database connection check.
   * Database is already initialized by Go backend, we just verify it works.
   */
  async initialize(callbacks?: DatabaseLoadingCallbacks): Promise<void> {
    if (this.initPromise) return this.initPromise;

    this.callbacks = callbacks;

    this.initPromise = (async () => {
      try {
        if (this.isInitialized) return;

        const time = Date.now();
        
        this.callbacks?.onStateChange?.(DatabaseLoadingState.Initializing);
        this.callbacks?.onProgress?.({
          phase: 'dependencies',
          progress: 30,
        });

        // Verify Wails bindings are available
        if (typeof DB === 'undefined') {
          throw new Error('Wails DB bindings not available');
        }

        this.callbacks?.onProgress?.({
          phase: 'dependencies',
          progress: 60,
        });

        // Test database connection with a simple query
        try {
          // Try to ping the database with a simple count query
          // This will fail if database is not initialized
          await DB.CountMessages('test-connection-check');
          console.log('✅ Database connection verified');
        } catch (e) {
          console.error('❌ Database connection failed:', e);
          throw new Error('Failed to connect to database. Make sure backend is running.');
        }

        this.callbacks?.onProgress?.({
          costTime: Date.now() - time,
          phase: 'dependencies',
          progress: 100,
        });

        this.isInitialized = true;
        this.callbacks?.onStateChange?.(DatabaseLoadingState.Finished);
        console.log(`✅ Database ready in ${Date.now() - time}ms`);

        await sleep(50);

        this.callbacks?.onStateChange?.(DatabaseLoadingState.Ready);
      } catch (e) {
        this.initPromise = null;
        this.callbacks?.onStateChange?.(DatabaseLoadingState.Error);
        const error = e as Error;

        this.callbacks?.onError?.({
          error: {
            message: error.message,
            name: error.name,
            stack: error.stack,
          },
          message: 'Database connection failed. Please restart the application.',
        });

        console.error('Database initialization error:', error);
        throw error;
      }
    })();

    return this.initPromise;
  }

  /**
   * Check if database is initialized
   */
  get isReady(): boolean {
    return this.isInitialized;
  }

  /**
   * Reset database state (for testing/debugging)
   */
  async resetDatabase(): Promise<void> {
    this.isInitialized = false;
    this.initPromise = null;
    console.log(`✅ Database '${DB_NAME}' state reset`);
  }
}

// Export singleton instance
const dbManager = DatabaseManager.getInstance();

/**
 * Legacy clientDB export - now just a marker object.
 * All database operations should use direct Wails bindings (DB.*).
 * This is kept for backward compatibility with existing code.
 */
export const clientDB = {
  _type: 'wails-binding',
  _note: 'Use DB.* from @/types/database for all database operations',
} as any;

/**
 * Initialize database connection check
 */
export const initializeDB = (callbacks?: DatabaseLoadingCallbacks) =>
  dbManager.initialize(callbacks);

/**
 * Reset database state
 */
export const resetClientDatabase = async () => {
  await dbManager.resetDatabase();
};

/**
 * Check if database is ready
 */
export const isDatabaseReady = () => dbManager.isReady;

/**
 * Get database configuration
 */
export const getClientDBConfig = () => ({
  mode: 'client' as const,
  driver: 'wails',
  initialized: dbManager.isReady,
});

/**
 * Legacy migration method - no longer needed.
 * Database is initialized and migrated by Go backend.
 * Kept for backward compatibility.
 */
export const updateMigrationRecord = async (migrationHash: string) => {
  console.warn('updateMigrationRecord is deprecated - migrations handled by Go backend');
  // Do nothing - migrations are handled by Go backend
};
