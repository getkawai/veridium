import { CreateThreadParams, ThreadStatus } from '@/types';
import { nanoid } from 'nanoid';

import { ThreadItem } from '../schemas';
import {
  DB,
  toNullString,
  getNullableString,
  currentTimestampMs,
} from '@/types/database';

export class ThreadModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: CreateThreadParams) => {
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

      return this.mapThread(result);
    } catch (error) {
      // ON CONFLICT DO NOTHING behavior - return undefined if conflict
      return undefined;
    }
  };

  delete = async (id: string) => {
    await DB.DeleteThread({
      id,
      userId: this.userId,
    });
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
    const now = currentTimestampMs();

    await DB.UpdateThread({
      id,
      userId: this.userId,
      title: toNullString(value.title as any),
      status: value.status || ThreadStatus.Active,
      lastActiveAt: value.lastActiveAt ? new Date(value.lastActiveAt).getTime() : now,
      updatedAt: now,
    });
  };

  // **************** Helper *************** //

  private mapThread = (thread: any): ThreadItem => {
    return {
      id: thread.id,
      title: getNullableString(thread.title as any),
      type: thread.type,
      status: thread.status,
      topicId: getNullableString(thread.topicId as any),
      sourceMessageId: getNullableString(thread.sourceMessageId as any),
      parentThreadId: getNullableString(thread.parentThreadId as any),
      clientId: getNullableString(thread.clientId as any),
      userId: thread.userId,
      lastActiveAt: new Date(thread.lastActiveAt),
      createdAt: new Date(thread.createdAt),
      updatedAt: new Date(thread.updatedAt),
    } as ThreadItem;
  };
}

