import { useEffect } from 'react';

import { useChatStore } from '@/store/chat';
import { useGlobalStore } from '@/store/global';
import { systemStatusSelectors } from '@/store/global/selectors';
import { useSessionStore } from '@/store/session';

/**
 * Fetch topics for the current session
 */
export const useFetchTopics = () => {
  const [sessionId] = useSessionStore((s) => [s.activeId]);
  const internal_fetchTopics = useChatStore((s) => s.internal_fetchTopics);
  const isDBInited = useGlobalStore(systemStatusSelectors.isDBInited);

  useEffect(() => {
    if (!isDBInited || !sessionId) return;

    internal_fetchTopics(sessionId);
  }, [isDBInited, sessionId, internal_fetchTopics]);
};
