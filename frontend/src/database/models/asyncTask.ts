import {
  AsyncTaskError,
  AsyncTaskErrorType,
  AsyncTaskStatus,
  AsyncTaskType,
} from '@/types';
import { nanoid } from 'nanoid';

import { AsyncTaskSelectItem, NewAsyncTaskItem } from '../schemas';
import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  toNullInt,
  currentTimestampMs,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

// set timeout to about 5 minutes, and give 2s padding time
export const ASYNC_TASK_TIMEOUT = 298 * 1000;

export class AsyncTaskModel {
  private userId: string;
  private logger = createModelLogger('AsyncTask', 'AsyncTaskModel', 'database/models/asyncTask');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: Pick<NewAsyncTaskItem, 'type' | 'status'>): Promise<string> => {
    const now = currentTimestampMs();
    const id = nanoid();

    await DB.CreateAsyncTask({
      id,
      type: params.type,
      status: params.status,
      error: toNullJSON(null),
      userId: this.userId,
      duration: toNullInt(null),
      createdAt: now,
      updatedAt: now,
    });

    return id;
  };

  delete = async (id: string) => {
    await DB.DeleteAsyncTask({
      id,
      userId: this.userId,
    });
  };

  findById = async (id: string) => {
    try {
      const result = await DB.GetAsyncTask({
        id,
        userId: this.userId,
      });
      return this.mapAsyncTask(result);
    } catch {
      return undefined;
    }
  };

  update(taskId: string, value: Partial<AsyncTaskSelectItem>) {
    const now = currentTimestampMs();

    return DB.UpdateAsyncTask({
      id: taskId,
      userId: this.userId,
      status: value.status || AsyncTaskStatus.Processing,
      error: toNullJSON(value.error),
      duration: toNullInt(value.duration as any),
      updatedAt: now,
    });
  }

  findByIds = async (taskIds: string[], type: AsyncTaskType): Promise<AsyncTaskSelectItem[]> => {
    let chunkTasks: AsyncTaskSelectItem[] = [];

    if (taskIds.length > 0) {
      await this.checkTimeoutTasks(taskIds);
      const results = await DB.GetAsyncTasksByIds({
        ids: taskIds,
        type,
      });
      chunkTasks = results.map((r) => this.mapAsyncTask(r));
    }

    return chunkTasks;
  };

  /**
   * make the task status to be `error` if the task is not finished in 5 minutes
   */
  checkTimeoutTasks = async (ids: string[]) => {
    const timeoutThreshold = Date.now() - ASYNC_TASK_TIMEOUT;

    const tasks = await DB.GetTimeoutTasks({
      ids,
      status: AsyncTaskStatus.Processing,
      createdAt: timeoutThreshold,
    });

    if (tasks.length > 0) {
      const now = currentTimestampMs();
      const taskIds = tasks.map((t) => t.id);

      await DB.UpdateTimeoutTasks({
        ids: taskIds,
        status: AsyncTaskStatus.Error,
        error: toNullJSON(
          new AsyncTaskError(AsyncTaskErrorType.Timeout, 'task is timeout, please try again'),
        ),
        updatedAt: now,
      });
    }
  };

  // **************** Helper *************** //

  private mapAsyncTask = (task: any): AsyncTaskSelectItem => {
    return {
      id: task.id,
      type: task.type,
      status: task.status,
      error: parseNullableJSON(task.error as any),
      userId: task.userId,
      duration: task.duration,
      createdAt: new Date(task.createdAt),
      updatedAt: new Date(task.updatedAt),
    } as AsyncTaskSelectItem;
  };
}

