import { Open as OpenSQLite, Close as CloseSQLite } from '@@/github.com/wailsapp/wails/v3/pkg/services/sqlite/sqliteservice';

import { WailsSQLiteDriver } from './wails-sqlite-driver';

let isInitialized = false;

/**
 * Initialize the Wails SQLite connection
 */
export async function initWailsSQLite(): Promise<WailsSQLiteDriver> {
  if (!isInitialized) {
    try {
      await OpenSQLite();
      isInitialized = true;
      console.log('✅ Wails SQLite connection opened');
    } catch (error) {
      console.error('❌ Failed to open Wails SQLite connection:', error);
      throw error;
    }
  }

  const driver = new WailsSQLiteDriver({
    onQuery: (query, params) => {
      // Optional: Log queries in development
      if (process.env.NODE_ENV === 'development') {
        console.log('[SQLite Query]', query, params);
      }
    },
  });

  return driver;
}

/**
 * Close the Wails SQLite connection
 */
export async function closeWailsSQLite(): Promise<void> {
  if (isInitialized) {
    try {
      await CloseSQLite();
      isInitialized = false;
      console.log('✅ Wails SQLite connection closed');
    } catch (error) {
      console.error('❌ Failed to close Wails SQLite connection:', error);
      throw error;
    }
  }
}

