import {
  AsyncTaskError,
  AsyncTaskErrorType,
  AsyncTaskStatus,
  AsyncTaskType,
  AsyncTaskSelectItem,
  NewAsyncTaskItem,
} from '@/types';
import { nanoid } from 'nanoid';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  toNullInt,
  getNullableString,
  getNullableInt,
  currentTimestampMs,
  AsyncTask,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';
import { NotificationService, NotificationOptions } from '@@/github.com/wailsapp/wails/v3/pkg/services/notifications';

// set timeout to about 5 minutes, and give 2s padding time
export const ASYNC_TASK_TIMEOUT = 298 * 1000;

export class AsyncTaskModel {
  private userId: string;
  private logger = createModelLogger('AsyncTask', 'AsyncTaskModel', 'database/models/asyncTask');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  /**
   * Show error notification to user
   */
  private async showErrorNotification(title: string, message: string) {
    try {
      await NotificationService.SendNotification(
        new NotificationOptions({
          id: `async-task-error-${Date.now()}`,
          title: `Async Task Error: ${title}`,
          body: message,
        })
      );
    } catch (notifError) {
      // Silently fail if notification fails - don't want notification errors to break the app
      console.error('Failed to show notification:', notifError);
    }
  }

  create = async (params: Pick<NewAsyncTaskItem, 'type' | 'status'>): Promise<string> => {
    await this.logger.methodEntry('create', { type: params.type, status: params.status, userId: this.userId });
    
    try {
      const now = currentTimestampMs();
      const id = nanoid();

      await DB.CreateAsyncTask({
        id,
        type: toNullString(params.type as any),
        status: toNullString(params.status as any),
        error: toNullJSON(null),
        userId: this.userId,
        duration: toNullInt(null),
        createdAt: now,
        updatedAt: now,
      });

      await this.logger.methodExit('create', { id });
      return id;
    } catch (error) {
      await this.logger.error('Failed to create async task', { error, params });
      await this.showErrorNotification(
        'Task Creation Failed',
        `Failed to create ${params.type} task. Please try again.`
      );
      throw error;
    }
  };

  delete = async (id: string) => {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    try {
      await DB.DeleteAsyncTask({
        id,
        userId: this.userId,
      });
      
      await this.logger.methodExit('delete', { id });
    } catch (error) {
      await this.logger.error('Failed to delete async task', { error, id });
      // Don't show notification for delete errors - these are usually cleanup operations
      throw error;
    }
  };

  findById = async (id: string) => {
    try {
      const result: AsyncTask = await DB.GetAsyncTask({
        id,
        userId: this.userId,
      });
      return this.mapAsyncTask(result);
    } catch (error) {
      await this.logger.warn('Async task not found', { id, error });
      return undefined;
    }
  };

  update = async (taskId: string, value: Partial<AsyncTaskSelectItem>) => {
    await this.logger.methodEntry('update', { taskId, value, userId: this.userId });
    
    try {
      const now = currentTimestampMs();

      const result = await DB.UpdateAsyncTask({
        id: taskId,
        userId: this.userId,
        status: toNullString(value.status || AsyncTaskStatus.Processing) as any,
        error: toNullJSON(value.error),
        duration: toNullInt(value.duration as any),
        updatedAt: now,
      });
      
      await this.logger.methodExit('update', { taskId });
      return result;
    } catch (error) {
      await this.logger.error('Failed to update async task', { error, taskId, value });
      // Don't show notification for update errors - these are background operations
      throw error;
    }
  };

  findByIds = async (taskIds: string[], type: AsyncTaskType): Promise<AsyncTaskSelectItem[]> => {
    try {
      let chunkTasks: AsyncTaskSelectItem[] = [];

      if (taskIds.length > 0) {
        await this.checkTimeoutTasks(taskIds);
        const results = await DB.GetAsyncTasksByIds({
          ids: taskIds,
          type: toNullString(type) as any,
        });
        chunkTasks = results.map((r) => this.mapAsyncTask(r));
      }

      return chunkTasks;
    } catch (error) {
      await this.logger.error('Failed to find async tasks by IDs', { error, taskIds, type });
      throw error;
    }
  };

  /**
   * make the task status to be `error` if the task is not finished in 5 minutes
   */
  checkTimeoutTasks = async (ids: string[]) => {
    try {
      const timeoutThreshold = Date.now() - ASYNC_TASK_TIMEOUT;

      const taskIds = await DB.GetTimeoutTasks({
        ids,
        status: toNullString(AsyncTaskStatus.Processing) as any,
        createdAt: timeoutThreshold,
      });

      if (taskIds.length > 0) {
        const now = currentTimestampMs();

        await this.logger.warn('Marking tasks as timeout', { taskIds, count: taskIds.length });

        await DB.UpdateTimeoutTasks({
          ids: taskIds,
          status: toNullString(AsyncTaskStatus.Error) as any,
          error: toNullJSON(
            new AsyncTaskError(AsyncTaskErrorType.Timeout, 'task is timeout, please try again'),
          ),
          updatedAt: now,
        });
      }
    } catch (error) {
      await this.logger.error('Failed to check timeout tasks', { error, ids });
      // Don't throw - this is a background check operation
    }
  };

  // **************** Helper *************** //

  private mapAsyncTask = (task: AsyncTask): AsyncTaskSelectItem => {
    return {
      id: task.id,
      type: getNullableString(task.type as any) ?? null,
      status: getNullableString(task.status as any) ?? null,
      error: parseNullableJSON(task.error as any),
      userId: task.userId,
      duration: getNullableInt(task.duration as any) ?? null,
      createdAt: new Date(task.createdAt),
      updatedAt: new Date(task.updatedAt),
    };
  };
}

