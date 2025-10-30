'use client';

import { SearchBar } from '@lobehub/ui';
import { memo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { type ChangeEvent } from 'react';

// Dummy implementations for UI focus
const dummyUserState = {
  isLoaded: true,
};

const settingsSelectors = {
  getHotkeyById: (id: string) => () => (id === 'search' ? '⌘K' : ''),
};

const HotkeyEnum = {
  Search: 'search',
} as const;

const useUserStore = (selector: (s: typeof dummyUserState) => any): any => {
  return selector(dummyUserState);
};

const dummySessionState = {
  sessionSearchKeywords: '',
  useSearchSessions: (keywords: string) => ({
    isValidating: keywords.length > 0 && Math.random() > 0.7, // Dummy validation
  }),
  updateSearchKeywords: (value: string) => {
    // Dummy update
  },
};

const useSessionStore = (selector: (s: typeof dummySessionState) => any): any => {
  return selector(dummySessionState);
};

const SessionSearchBar = memo<{ mobile?: boolean }>(({ mobile }) => {
  const { t } = useTranslation('chat');
  const isLoaded = useUserStore((s) => s.isLoaded);
  const hotkey = useUserStore(settingsSelectors.getHotkeyById(HotkeyEnum.Search));

  const [keywords, useSearchSessions, updateSearchKeywords] = useSessionStore((s) => [
    s.sessionSearchKeywords,
    s.useSearchSessions,
    s.updateSearchKeywords,
  ]);

  const { isValidating } = useSearchSessions(keywords);

  const handleChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      updateSearchKeywords(e.target.value);
    },
    [updateSearchKeywords],
  );

  return (
    <SearchBar
      allowClear
      enableShortKey={!mobile}
      loading={!isLoaded || isValidating}
      onChange={handleChange}
      placeholder={t('searchAgentPlaceholder')}
      shortKey={hotkey}
      spotlight={!mobile}
      value={keywords}
      variant={'filled'}
    />
  );
});

export default SessionSearchBar;
