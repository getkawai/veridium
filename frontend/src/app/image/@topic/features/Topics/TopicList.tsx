'use client';

import { useAutoAnimate } from '@formkit/auto-animate/react';
import { useSize } from 'ahooks';
import { memo, useEffect, useRef } from 'react';
import { Flexbox } from 'react-layout-kit';

import { useImageStore } from '@/store/image';
import { generationTopicSelectors } from '@/store/image/selectors';
import { useUserStore } from '@/store/user';
import { authSelectors } from '@/store/user/slices/auth/selectors';

import NewTopicButton from './NewTopicButton';
import TopicItem from './TopicItem';

// Lists existing generation topics and the “new topic” entry.
const TopicsList = memo(() => {
  const isLogin = useUserStore(authSelectors.isLogin);
  const fetchGenerationTopics = useImageStore((s) => s.fetchGenerationTopics);
  // Only fetch topics after the user logs in.
  useEffect(() => {
    fetchGenerationTopics(!!isLogin);
  }, [isLogin, fetchGenerationTopics]);
  const ref = useRef(null);
  const { width = 80 } = useSize(ref) || {};
  const [parent] = useAutoAnimate();
  const generationTopics = useImageStore(generationTopicSelectors.generationTopics);
  const openNewGenerationTopic = useImageStore((s) => s.openNewGenerationTopic);

  const showMoreInfo = Boolean(width > 120);

  const switchGenerationTopic = useImageStore((s) => s.switchGenerationTopic);
  const activeTopicId = useImageStore(generationTopicSelectors.activeGenerationTopicId);

  // Ref to track if we have performed the initial auto-selection/validation
  const hasInitialized = useRef(false);

  // Auto-select first topic logic:
  // 1. If no topic is active (empty persistence or first load) -> select first.
  // 2. If active topic exists in store but NOT in the fetched list (stale persistence) -> select first (or null if list empty).
  // 3. IMPORTANT: Only do this ONCE on mount/load. If user manually clicks "New Topic" (setting ID to null), do NOT auto-select again.
  useEffect(() => {
    if (!generationTopics) return;

    // 1. Handle Stale/Invalid ID (Always check this to avoid ghost states)
    if (activeTopicId && !generationTopics.some(t => t.id === activeTopicId)) {
      if (generationTopics.length > 0) {
        switchGenerationTopic(generationTopics[0].id);
      } else {
        openNewGenerationTopic();
      }
      return;
    }

    // 2. Initial Auto-Select (Only run if we haven't initialized yet)
    if (!hasInitialized.current) {
      if (!activeTopicId && generationTopics.length > 0) {
        switchGenerationTopic(generationTopics[0].id);
      }
      hasInitialized.current = true;
    }
  }, [generationTopics, activeTopicId, switchGenerationTopic, openNewGenerationTopic]);

  const isEmpty = !generationTopics || generationTopics.length === 0;
  if (isEmpty) {
    // Hide the list until topics are available.
    return null;
  }

  return (
    <Flexbox
      align="center"
      gap={12}
      ref={ref}
      style={{
        maxHeight: '100%',
        overflowY: 'auto',
      }}
      width={'100%'}
    >
      <NewTopicButton
        count={generationTopics?.length}
        onClick={openNewGenerationTopic}
        showMoreInfo={showMoreInfo}
      />
      <Flexbox align="center" gap={12} ref={parent} width={'100%'}>
        {generationTopics.map((topic, index) => (
          <TopicItem
            key={topic.id}
            showMoreInfo={showMoreInfo}
            style={{
              padding:
                // fix the avatar border is clipped by overflow hidden
                generationTopics.length === 1
                  ? '4px 0'
                  : index === generationTopics.length - 1
                    ? '0 0 4px'
                    : '0',
            }}
            topic={topic}
          />
        ))}
      </Flexbox>
    </Flexbox>
  );
});

TopicsList.displayName = 'TopicsList';

export default TopicsList;
