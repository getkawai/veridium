import { CreateThreadParams, ThreadStatus, ThreadType } from '@/types';
import { nanoid } from 'nanoid';
import { createModelLogger } from '@/utils/logger';

import { ThreadItem } from '@/types/topic/thread';
import {
  DB,
  toNullString,
  getNullableString,
  currentTimestampMs,
  Thread,
} from '@/types/database';

export class ThreadModel {
  private userId: string;
  private logger = createModelLogger('Thread', 'ThreadModel', 'database/models/thread');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: CreateThreadParams) => {
    await this.logger.methodEntry('create', { type: params.type, userId: this.userId });
    
    const now = currentTimestampMs();

    try {
      const result: Thread = await DB.CreateThread({
        id: nanoid(),
        title: params.title ?? 'Untitled',
        type: params.type,
        status: toNullString(ThreadStatus.Active as string), // New threads should always start as 'active'
        topicId: params.topicId,
        sourceMessageId: params.sourceMessageId,
        parentThreadId: toNullString(params.parentThreadId as any),
        clientId: toNullString(''),
        userId: this.userId,
        lastActiveAt: now,
        createdAt: now,
        updatedAt: now,
      });
      
      // With ON CONFLICT DO NOTHING, if conflict occurs, result might be empty/null
      // Check if result is valid
      if (!result || !result.id) {
        await this.logger.warn('Thread create conflict (ON CONFLICT DO NOTHING), returning undefined', { 
          params,
        });
        return undefined;
      }
      
      await this.logger.methodExit('create', { threadId: result.id });
      return this.mapThread(result);
    } catch (error: any) {
      // Check if this is a "no rows in result set" error (ON CONFLICT DO NOTHING with empty RETURNING)
      const isNoRowsError = 
        error?.message?.includes('no rows in result set') ||
        error?.message?.includes('sql: no rows in result set') ||
        error?.code === 'PGRST116'; // PostgREST no rows error code
      
      if (isNoRowsError) {
        // ON CONFLICT DO NOTHING behavior - conflict occurred, no row returned
        await this.logger.warn('Thread create conflict (ON CONFLICT DO NOTHING - no rows returned), returning undefined', { 
          params,
          errorMessage: error?.message,
        });
        return undefined;
      }
      
      // Check if this is a UNIQUE constraint conflict (fallback for databases that throw error)
      const isConflictError = 
        error?.message?.includes('UNIQUE constraint') || 
        error?.message?.includes('2067') ||
        error?.code === 'SQLITE_CONSTRAINT_UNIQUE' ||
        error?.code === 2067;
      
      if (isConflictError) {
        // ON CONFLICT DO NOTHING behavior - return undefined if conflict
        await this.logger.warn('Thread create conflict (UNIQUE constraint), returning undefined', { 
          params,
          errorMessage: error?.message,
        });
        return undefined;
      }
      
      // For other errors, log the full error and re-throw
      await this.logger.error('Thread create failed with non-conflict error', { 
        params,
        error: error?.message || String(error),
        errorCode: error?.code,
        stack: error?.stack,
      });
      throw error;
    }
  };

  delete = async (id: string) => {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    await DB.DeleteThread({
      id,
      userId: this.userId,
    });
    
    await this.logger.methodExit('delete', { id });
  };

  deleteAll = async () => {
    await DB.DeleteAllThreads(this.userId);
  };

  query = async () => {
    const data = await DB.ListAllThreads(this.userId);
    return data.map((t) => this.mapThread(t));
  };

  queryByTopicId = async (topicId: string) => {
    const data = await DB.ListThreadsByTopic({
      topicId: topicId,  // Plain string, not NullString
      userId: this.userId,
    });
    return data.map((t) => this.mapThread(t));
  };

  findById = async (id: string) => {
    try {
      const thread = await DB.GetThread({
        id,
        userId: this.userId,
      });
      return this.mapThread(thread);
    } catch {
      return undefined;
    }
  };

  update = async (id: string, value: Partial<ThreadItem>) => {
    await this.logger.methodEntry('update', { id, value, userId: this.userId });
    
    const now = currentTimestampMs();

    await DB.UpdateThread({
      id,
      userId: this.userId,
      title: value.title ?? '',
      status: toNullString(value.status || ThreadStatus.Active) as any,
      lastActiveAt: value.lastActiveAt ? new Date(value.lastActiveAt).getTime() : now,
      updatedAt: now,
    });
    
    await this.logger.methodExit('update', { id });
  };

  // **************** Helper *************** //

  private mapThread = (thread: Thread): ThreadItem => {
    const statusStr = getNullableString(thread.status as any);
    
    return {
      id: thread.id,
      title: thread.title,
      type: thread.type as ThreadType,
      status: (statusStr as ThreadStatus) || ThreadStatus.Active,
      topicId: thread.topicId,
      sourceMessageId: thread.sourceMessageId,
      parentThreadId: getNullableString(thread.parentThreadId as any),
      userId: thread.userId,
      lastActiveAt: new Date(thread.lastActiveAt),
      createdAt: new Date(thread.createdAt),
      updatedAt: new Date(thread.updatedAt),
    };
  };
}

