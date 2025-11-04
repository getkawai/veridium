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
      const result = await DB.CreateThread({
        id: nanoid(),
        title: toNullString(params.title as any),
        type: params.type,
        status: params.type as any, // ThreadStatus is string, need to pass as NullString
        topicId: toNullString(params.topicId) as any,
        sourceMessageId: toNullString(params.sourceMessageId) as any,
        parentThreadId: toNullString(params.parentThreadId as any),
        clientId: toNullString(''),
        userId: this.userId,
        lastActiveAt: now,
        createdAt: now,
        updatedAt: now,
      });
      
      await this.logger.methodExit('create', { threadId: result.id });
      return this.mapThread(result);
    } catch (error) {
      // ON CONFLICT DO NOTHING behavior - return undefined if conflict
      await this.logger.warn('Thread create conflict, returning undefined', { params });
      return undefined;
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
      title: toNullString(value.title as any),
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
      title: getNullableString(thread.title as any) || '',
      type: thread.type as ThreadType,
      status: (statusStr as ThreadStatus) || ThreadStatus.Active,
      topicId: getNullableString(thread.topicId as any) || '',
      sourceMessageId: getNullableString(thread.sourceMessageId as any) || '',
      parentThreadId: getNullableString(thread.parentThreadId as any),
      userId: thread.userId,
      lastActiveAt: new Date(thread.lastActiveAt),
      createdAt: new Date(thread.createdAt),
      updatedAt: new Date(thread.updatedAt),
    };
  };
}

