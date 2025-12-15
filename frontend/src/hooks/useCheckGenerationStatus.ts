import { useEffect, useRef } from 'react';

import { GetGenerationWithAsyncTask } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';
import { AsyncTaskStatus } from '@/types/asyncTask';
import { useImageStore } from '@/store/image';
import { generationTopicSelectors } from '@/store/image/slices/generationTopic/selectors';
import { GetGenerationWithAsyncTaskRow } from '@/types/database';

export const useCheckGenerationStatus = (
  generationId: string,
  asyncTaskId: string | null | undefined,
  topicId: string,
  enable = true,
) => {
  const timeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);
  const requestCountRef = useRef(0);
  const isErrorRef = useRef(false);

  useEffect(() => {
    // Reset state when inputs change
    requestCountRef.current = 0;
    isErrorRef.current = false;

    if (!enable || !generationId || !asyncTaskId || !topicId) {
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
      return;
    }

    const checkStatus = async () => {
      try {
        requestCountRef.current += 1;

        // Fetch status
        const result: GetGenerationWithAsyncTaskRow = await GetGenerationWithAsyncTask(generationId);

        if (!result) return;

        const status = result.asyncTaskStatus && result.asyncTaskStatus.Valid
          ? result.asyncTaskStatus.String
          : 'pending';

        // Update if finalized
        if (status === AsyncTaskStatus.Success || status === AsyncTaskStatus.Error) {
          const store = useImageStore.getState();

          // Find target batch
          const currentBatches = store.generationBatchesMap[topicId] || [];
          const targetBatch = currentBatches.find((batch) =>
            batch.generations.some((gen) => gen.id === generationId),
          );

          if (targetBatch) {
            // Update store logic adapted from action.ts
            const generationUpdate = {
              id: result.id,
              asyncTaskId: result.asyncTaskId && result.asyncTaskId.Valid ? result.asyncTaskId.String : null,
              createdAt: new Date(result.createdAt || Date.now()), // createdAt is number
              seed: result.seed && result.seed.Valid ? result.seed.Int64 : null,
              task: {
                id: asyncTaskId,
                status: status as any,
                error: result.asyncTaskError && result.asyncTaskError.Valid ? JSON.parse(result.asyncTaskError.String) : undefined,
              },
              asset: result.asset && result.asset.Valid ? JSON.parse(result.asset.String) : null,
            };

            store.internal_dispatchGenerationBatch(
              topicId,
              {
                type: 'updateGenerationInBatch',
                batchId: targetBatch.id,
                generationId,
                value: generationUpdate,
              },
              `useCheckGenerationStatus/${status}`
            );

            // Update topic cover if needed
            if (status === AsyncTaskStatus.Success && generationUpdate.asset?.thumbnailUrl) {
              const currentTopic = generationTopicSelectors.getGenerationTopicById(topicId)(store);

              // If current topic has no coverUrl, update it with this generation's thumbnail
              if (currentTopic && !currentTopic.coverUrl) {
                await store.updateGenerationTopicCover(
                  topicId,
                  generationUpdate.asset.thumbnailUrl,
                );
              }
            }

            // Refresh batches
            await store.refreshGenerationBatches();
          }

          // Stop polling
          return;
        }

        // Continue polling with exponential backoff
        const baseInterval = 1000;
        const maxInterval = 5000; // Cap at 5s for responsiveness
        const backoffMultiplier = Math.floor(requestCountRef.current / 5);
        let dynamicInterval = Math.min(
          baseInterval * Math.pow(1.5, backoffMultiplier),
          maxInterval,
        );

        if (isErrorRef.current) {
          dynamicInterval = Math.min(dynamicInterval * 2, maxInterval);
        }

        timeoutRef.current = setTimeout(checkStatus, dynamicInterval);

      } catch (error) {
        console.error('Failed to check generation status:', error);
        isErrorRef.current = true;
        // Retry with delay
        timeoutRef.current = setTimeout(checkStatus, 5000);
      }
    };

    // Start polling
    checkStatus();

    return () => {
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
    };
  }, [enable, generationId, asyncTaskId, topicId]);
};
