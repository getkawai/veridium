'use client';

import { memo } from 'react';

import DefaultMode from './DefaultMode';
import SearchMode from './SearchMode';

import { useSessionStore } from '@/store/session';

const SessionListContent = memo(() => {
  const isSearching = useSessionStore((s) => s.isSearching);

  return isSearching ? <SearchMode /> : <DefaultMode />;
});

SessionListContent.displayName = 'SessionListContent';

export default SessionListContent;
