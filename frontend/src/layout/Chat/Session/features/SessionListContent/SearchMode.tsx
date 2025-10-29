import { memo, useMemo } from 'react';
import { LobeAgentSession, LobeSessionType, LobeSessions } from '@/types/session';

import SkeletonList from '../SkeletonList';
import SessionList from './List';

// Dummy implementations for development
const useServerConfigStore = (selector: any) => {
  if (selector) {
    return selector({ isMobile: false });
  }
  return { isMobile: false };
};

const serverConfigSelectors = {
  isMobile: (state: any) => state.isMobile,
};

const useSessionStore = (selector?: any, comparator?: any) => {
  if (selector) {
    return selector({
      defaultSessions: [],
      customSessionGroups: [],
      pinnedSessions: [],
      isSearching: false,
      sessionSearchKeywords: '',
      useSearchSessions: (keywords: string) => ({
        data: [],
        isLoading: false,
      }),
    });
  }
  return {
    defaultSessions: [],
    customSessionGroups: [],
    pinnedSessions: [],
    isSearching: false,
    sessionSearchKeywords: '',
    useSearchSessions: (keywords: string) => ({
      data: [],
      isLoading: false,
    }),
  };
};

const SearchMode = memo(() => {
  const [sessionSearchKeywords, useSearchSessions] = useSessionStore((s) => [
    s.sessionSearchKeywords,
    s.useSearchSessions,
  ]);

  const isMobile = useServerConfigStore(serverConfigSelectors.isMobile);

  const { data, isLoading } = useSearchSessions(sessionSearchKeywords);

  const filteredData = useMemo(() => {
    if (!data) return data;

    if (isMobile) {
      return data.filter((session: LobeSessions[0]) => session.type !== LobeSessionType.Group);
    }

    return data.filter(
      (session: LobeSessions[0]) =>
        session.type !== LobeSessionType.Agent || !(session as LobeAgentSession).config?.virtual,
    );
  }, [data, isMobile]);

  return isLoading ? (
    <SkeletonList />
  ) : (
    <SessionList dataSource={filteredData} showAddButton={false} />
  );
});

SearchMode.displayName = 'SessionSearchMode';

export default SearchMode;
