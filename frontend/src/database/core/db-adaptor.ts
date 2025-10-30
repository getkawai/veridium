import { clientDB, initializeDB } from '../client/db';
import { LobeChatDatabase } from '../type';

/**
 * 懒加载数据库实例
 * 避免每次模块导入时都初始化数据库
 * 
 * Note: Now using Wails SQLite for all environments
 */
let cachedDB: LobeChatDatabase | null = null;

export const getServerDB = async (): Promise<LobeChatDatabase> => {
  // 如果已经有缓存的实例，直接返回
  if (cachedDB) return cachedDB;

  try {
    // Initialize Wails SQLite database
    await initializeDB();
    cachedDB = clientDB as LobeChatDatabase;
    return cachedDB;
  } catch (error) {
    console.error('❌ Failed to initialize database:', error);
    throw error;
  }
};

// Export the clientDB as serverDB for consistency
export const serverDB = clientDB as LobeChatDatabase;
