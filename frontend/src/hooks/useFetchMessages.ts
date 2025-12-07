import { useEffect } from 'react';

import { useChatStore } from '@/store/chat';
import { useSessionStore } from '@/store/session';
import { sessionSelectors } from '@/store/session/selectors';

export const useFetchMessages = () => {
  const sessionId = useSessionStore((s) => s.activeId);
  const [activeTopicId, internal_fetchMessages, internal_updateActiveSessionType] = useChatStore((s) => [
    s.activeTopicId,
    s.internal_fetchMessages,
    s.internal_updateActiveSessionType,
  ]);

  const [currentSession, isGroupSession] = useSessionStore((s) => [
    sessionSelectors.currentSession(s),
    sessionSelectors.isCurrentSessionGroupSession(s),
  ]);

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

    internal_fetchMessages(
      sessionId,
      activeTopicId,
      isGroupSession ? 'group' : 'session'
    );
  }, [sessionId, activeTopicId, isGroupSession, internal_fetchMessages]);
};
