import { useEffect, useRef } from 'react';

import { useChatStore } from '@/store/chat';
import { useSessionStore } from '@/store/session';

/**
 * Fetch topics for the current session
 * Also syncs chatStore.activeId with sessionStore.activeId
 * and resets activeTopicId when session changes to prevent showing
 * topics from the previous session
 */
export const useFetchTopics = () => {
  const [sessionId] = useSessionStore((s) => [s.activeId]);
  const [chatActiveId, internal_fetchTopics, switchTopic, internal_updateActiveId] = useChatStore((s) => [
    s.activeId,
    s.internal_fetchTopics,
    s.switchTopic,
    s.internal_updateActiveId,
  ]);

  // Track previous session to detect session changes
  const prevSessionIdRef = useRef<string | undefined>(undefined);

  useEffect(() => {
    if (!sessionId) return;

    // CRITICAL: Sync chatStore.activeId with sessionStore.activeId
    // This ensures topicMaps and messagesMap use the correct session key
    if (chatActiveId !== sessionId) {
      console.log('[useFetchTopics] Syncing chatStore.activeId with session', {
        chatActiveId,
        sessionId,
      });
      internal_updateActiveId(sessionId);
    }

    // Check if session changed (not first mount)
    const sessionChanged = prevSessionIdRef.current !== undefined && prevSessionIdRef.current !== sessionId;

    // Update ref BEFORE any state changes to avoid issues
    const previousSession = prevSessionIdRef.current;
    prevSessionIdRef.current = sessionId;

    if (sessionChanged) {
      // Reset activeTopicId to null when session changes
      // This ensures we start fresh on the "Default Topic" for the new session
      console.log('[useFetchTopics] Session changed, resetting activeTopicId', {
        from: previousSession,
        to: sessionId,
      });
      switchTopic(undefined, true); // Reset to default topic, skip message refresh (will be handled separately)
    }

    // Fetch topics for the new session
    internal_fetchTopics(sessionId);
  }, [sessionId, chatActiveId, internal_fetchTopics, switchTopic, internal_updateActiveId]);
};
