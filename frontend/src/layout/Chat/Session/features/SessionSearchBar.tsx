'use client';

import { SearchBar } from '@lobehub/ui';
import { type ChangeEvent, memo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';

// import { useSessionStore } from '@/store/session';
// import { useUserStore } from '@/store/user';
// import { settingsSelectors } from '@/store/user/selectors';
import { HotkeyEnum } from '@/types/hotkey';

const dummyUserState = {
  isLoaded: true,
};

const settingsSelectors = {
  getHotkeyById: (id: any) => (s: any) => '',
};

const useUserStore = (selector?: (s: typeof dummyUserState) => any) => {
  if (selector) return selector(dummyUserState);
  return dummyUserState;
};

const dummySessionState = {
  sessionSearchKeywords: '',
  useSearchSessions: (keywords: string) => ({ isValidating: false }),
  updateSearchKeywords: (value: string) => {},
};

const useSessionStore = (selector?: (s: typeof dummySessionState) => any): any => {
  if (selector) return selector(dummySessionState);
  return dummySessionState;
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
