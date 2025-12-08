'use client';

import { useQueryState } from 'nuqs';
import { useLayoutEffect } from 'react';
import { createStoreUpdater } from 'zustand-utils';

import { useImageStore } from '@/store/image';

/**
 * Two-way sync between the `topic` query param and `activeGenerationTopicId`.
 */
const TopicUrlSync = () => {
  const useStoreUpdater = createStoreUpdater(useImageStore);

  const [topic, setTopic] = useQueryState('topic', { history: 'replace', throttleMs: 500 });
  useStoreUpdater('activeGenerationTopicId', topic);

  useLayoutEffect(() => {
    const unsubscribeTopic = useImageStore.subscribe(
      (s) => s.activeGenerationTopicId,
      (state) => {
        setTopic(state || null);
      },
    );

    return () => {
      unsubscribeTopic();
    };
  }, [setTopic]);

  // This component renders nothing; it only keeps URL and store in sync.
  return null;
};

TopicUrlSync.displayName = 'TopicUrlSync';

export default TopicUrlSync;
