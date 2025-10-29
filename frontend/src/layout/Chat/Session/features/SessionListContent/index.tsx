'use client';

import { memo } from 'react';

import DefaultMode from './DefaultMode';
import SearchMode from './SearchMode';

// Dummy implementation for development
const useSessionStore = (selector?: any) => {
  if (selector) {
    return selector({
      isSearching: false,
    });
  }
  return {
    isSearching: false,
  };
};

const SessionListContent = memo(() => {
  const isSearching = useSessionStore((s) => s.isSearching);

  return isSearching ? <SearchMode /> : <DefaultMode />;
});

SessionListContent.displayName = 'SessionListContent';

export default SessionListContent;
