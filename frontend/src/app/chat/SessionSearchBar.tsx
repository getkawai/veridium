'use client';

import { SearchBar } from '@lobehub/ui';
import { memo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { type ChangeEvent } from 'react';

import { useUserStore } from '@/store/user';
import { useSessionStore } from '@/store/session';
import { settingsSelectors } from '@/store/user/selectors';
const HotkeyEnum = {
  Search: 'search',
} as const;

const SessionSearchBar = memo<{ mobile?: boolean }>(({ mobile }) => {
  const { t } = useTranslation('chat');
  const isUserStateInit = useUserStore((s) => s.isUserStateInit);
  const hotkey = useUserStore(settingsSelectors.getHotkeyById(HotkeyEnum.Search));

  const [keywords, updateSearchKeywords] = useSessionStore((s) => [
    s.sessionSearchKeywords,
    s.updateSearchKeywords,
  ]);

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
      loading={!isUserStateInit}
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
