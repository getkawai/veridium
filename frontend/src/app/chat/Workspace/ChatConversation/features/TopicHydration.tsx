'use client';

import { useEffect } from 'react';

import { useSearchParams } from '@/hooks/useNavigation';
import { useChatStore } from '@/store/chat';
import { chatSelectors } from '@/store/chat/selectors';

/**
 * TopicHydration component
 * Hydrates activeTopicId from:
 * 1. URL parameter 'topic' on page load/reload
 * 2. Auto-detect from messages if user is already in a topic
 * This ensures that the topic state is correctly set when user reloads the page
 * or when messages are loaded that belong to a topic
 */
const TopicHydration = () => {
  const searchParams = useSearchParams();
  const [activeTopicId, switchTopic, messages] = useChatStore((s) => [
    s.activeTopicId,
    s.switchTopic,
    chatSelectors.activeBaseChats(s)
  ]);

  useEffect(() => {
    console.log('[TopicHydration] Running hydration check...');

    const urlTopicId = searchParams.get('topic');

    console.log('[TopicHydration] Current state:', {
      urlTopicId,
      activeTopicId,
      hasUrl: !!urlTopicId,
      hasActive: !!activeTopicId,
      messageCount: messages.length
    });

    // Priority 1: Hydrate from URL parameter
    if (urlTopicId && urlTopicId !== activeTopicId) {
      console.log('[TopicHydration] Hydrating topic from URL:', urlTopicId);
      switchTopic(urlTopicId, true); // skipRefreshMessage = true to avoid double refresh
      return;
    }

    // Priority 2: Auto-detect topic from messages if not already set
    // Wait for messages to be loaded before checking
    if (!activeTopicId && !urlTopicId && messages.length > 0) {
      console.log('[TopicHydration] Checking messages for topic:', {
        messageCount: messages.length,
        firstMessage: messages[0] ? {
          id: messages[0].id,
          topicId: messages[0].topicId,
          role: messages[0].role
        } : null
      });

      // Find the first message that has a topicId
      const messageWithTopic = messages.find(m => m.topicId);

      if (messageWithTopic?.topicId) {
        console.log('[TopicHydration] Auto-detected topic from messages:', messageWithTopic.topicId);
        switchTopic(messageWithTopic.topicId, true); // skipRefreshMessage = true since messages are already loaded
      } else {
        console.log('[TopicHydration] No topic found in messages');
      }
    }
  }, [searchParams, activeTopicId, switchTopic, messages]);

  return null;
};

export default TopicHydration;
