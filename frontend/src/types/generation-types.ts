// Types moved from lambda routers (no longer needed for Wails)

import { AsyncTaskError, AsyncTaskStatus, Generation } from '@/types';

export type GetGenerationStatusResult = {
  error: AsyncTaskError | null;
  generation: Generation | null;
  status: AsyncTaskStatus;
};

export type UpdateTopicValue = {
  title?: string;
  coverUrl?: string;
};
