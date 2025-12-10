import { useEffect, useRef } from 'react';

import { useChatStore } from '@/store/chat';
import { useSessionStore } from '@/store/session';
import { sessionSelectors } from '@/store/session/selectors';
import { topicSelectors } from '@/store/chat/slices/topic/selectors';

export const useFetchMessages = () => {
  const sessionId = useSessionStore((s) => s.activeId);
  const [activeTopicId, internal_fetchMessages, internal_updateActiveSessionType, topicsInit] = useChatStore((s) => [
    s.activeTopicId,
    s.internal_fetchMessages,
    s.internal_updateActiveSessionType,
    s.topicsInit,
  ]);

  // Get topics for current session to validate activeTopicId belongs to this session
  const currentTopics = useChatStore((s) => topicSelectors.currentTopics(s));

  const [currentSession, isGroupSession] = useSessionStore((s) => [
    sessionSelectors.currentSession(s),
    sessionSelectors.isCurrentSessionGroupSession(s),
  ]);

  // Track previous session to detect session changes
  const prevSessionIdRef = useRef<string | undefined>(undefined);

  // Update active session type when session changes
  useEffect(() => {
    if (currentSession?.type) {
      internal_updateActiveSessionType(currentSession.type as 'agent' | 'group');
    } else {
      internal_updateActiveSessionType(undefined);
    }
  }, [currentSession?.id, currentSession?.type, internal_updateActiveSessionType]);

  // Fetch messages when dependencies change
  useEffect(() => {
    if (!sessionId) return;

    // Detect session change
    const sessionChanged = prevSessionIdRef.current !== undefined && prevSessionIdRef.current !== sessionId;
    prevSessionIdRef.current = sessionId;

    // If session changed, only fetch with null topicId (default topic)
    // This prevents fetching messages with topicId from previous session
    if (sessionChanged) {
      console.log('[useFetchMessages] Session changed, fetching with null topicId', {
        from: prevSessionIdRef.current,
        to: sessionId,
        ignoringTopicId: activeTopicId,
      });
      internal_fetchMessages(
        sessionId,
        undefined, // Force null topicId on session change
        isGroupSession ? 'group' : 'session'
      );
      return;
    }

    // Validate that activeTopicId belongs to current session
    // If topics are initialized and activeTopicId is not in current topics, use null
    let validatedTopicId = activeTopicId;
    if (topicsInit && activeTopicId) {
      const topicExists = currentTopics.some(t => t.id === activeTopicId);
      if (!topicExists) {
        console.log('[useFetchMessages] activeTopicId not in current session topics, using null', {
          activeTopicId,
          sessionId,
          currentTopicsCount: currentTopics.length,
        });
        validatedTopicId = undefined;
      }
    }

    internal_fetchMessages(
      sessionId,
      validatedTopicId,
      isGroupSession ? 'group' : 'session'
    );
  }, [sessionId, activeTopicId, isGroupSession, internal_fetchMessages, topicsInit, currentTopics]);
};
