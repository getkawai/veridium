export enum AsyncTaskType {
  Chunking = 'chunk',
  Embedding = 'embedding',
  ImageGeneration = 'image_generation',
}

export enum AsyncTaskStatus {
  Error = 'error',
  Pending = 'pending',
  Processing = 'processing',
  Success = 'success',
}

export enum AsyncTaskErrorType {
  EmbeddingError = 'EmbeddingError',
  InvalidProviderAPIKey = 'InvalidProviderAPIKey',
  /**
   * Model not found on server
   */
  ModelNotFound = 'ModelNotFound',
  /**
   * the chunk parse result it empty
   */
  NoChunkError = 'NoChunkError',
  ServerError = 'ServerError',
  /**
   * this happens when the task is not trigger successfully
   */
  TaskTriggerError = 'TaskTriggerError',
  Timeout = 'TaskTimeout',
}

export interface IAsyncTaskError {
  body: string | { detail: string };
  name: string;
}

export class AsyncTaskError implements IAsyncTaskError {
  constructor(name: string, message: string) {
    this.name = name;
    this.body = { detail: message };
  }

  name: string;

  body: { detail: string };
}

export interface FileParsingTask {
  chunkCount?: number | null;
  chunkingError?: IAsyncTaskError | null;
  chunkingStatus?: AsyncTaskStatus | null;
  embeddingError?: IAsyncTaskError | null;
  embeddingStatus?: AsyncTaskStatus | null;
  finishEmbedding?: boolean;
}

/**
 * Async task item from database
 * Equivalent to: typeof asyncTasks.$inferSelect
 */
export interface AsyncTaskSelectItem {
  id: string;
  type: string | null;
  status: string | null;
  error: any;
  userId: string;
  duration: number | null;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * New async task (for insert operations)
 * Equivalent to: typeof asyncTasks.$inferInsert
 */
export interface NewAsyncTaskItem {
  id?: string;
  type?: string | null;
  status?: string | null;
  error?: any;
  userId: string;
  duration?: number | null;
  createdAt?: Date;
  updatedAt?: Date;
}
