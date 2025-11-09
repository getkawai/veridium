import { useEffect } from 'react';

import { useChatStore } from '@/store/chat';
import { useGlobalStore } from '@/store/global';
import { systemStatusSelectors } from '@/store/global/selectors';

export const useFetchThreads = (activeTopicId?: string) => {
  const isDBInited = useGlobalStore(systemStatusSelectors.isDBInited);
  const internal_fetchThreads = useChatStore((s) => s.internal_fetchThreads);

  useEffect(() => {
    if (!isDBInited || !activeTopicId) return;

    internal_fetchThreads(activeTopicId);
  }, [isDBInited, activeTopicId, internal_fetchThreads]);
};
