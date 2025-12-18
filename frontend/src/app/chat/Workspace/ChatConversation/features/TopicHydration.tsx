'use client';

/**
 * @deprecated This component is not used in desktop app (Veridium)
 * Kept for reference only - was part of LobeChat web app architecture
 * Desktop apps initialize state from database/messages, not URL
 * See: useInitTopicState hook in hooks/useInitTopicState.ts
 * Integrated in: layout/GlobalProvider/StoreInitialization.tsx
 */

import { useEffect } from 'react';

import { useChatStore } from '@/store/chat';
import { chatSelectors } from '@/store/chat/selectors';

/**
 * TopicHydration component
 * Hydrates activeTopicId from:
 * 1. Auto-detect from messages if user is already in a topic
 * This ensures that the topic state is correctly set when messages are loaded that belong to a topic
 */
const TopicHydration = () => {
  const [activeTopicId, switchTopic, messages] = useChatStore((s) => [
    s.activeTopicId,
    s.switchTopic,
    chatSelectors.activeBaseChats(s)
  ]);

  useEffect(() => {
    // Priority: Auto-detect topic from messages if not already set
    // Wait for messages to be loaded before checking
    if (!activeTopicId && messages.length > 0) {
      // Find the first message that has a topicId
      const messageWithTopic = messages.find(m => m.topicId);

      if (messageWithTopic?.topicId) {
        switchTopic(messageWithTopic.topicId, true); // skipRefreshMessage = true since messages are already loaded
      }
    }
  }, [activeTopicId, switchTopic, messages]);

  return null;
};

export default TopicHydration;
