'use client';

import { SearchBar } from '@lobehub/ui';
import { type ChangeEvent, memo, useCallback, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useUserStore } from '@/store/user';
import { useSessionStore } from '@/store/session';
import { settingsSelectors } from '@/store/user/selectors';
const HotkeyEnum = {
  Search: 'search',
} as const;

import { useDebounce } from '@/hooks/useDebounce';

const SessionSearchBar = memo<{ mobile?: boolean }>(({ mobile }) => {
  const { t } = useTranslation('chat');
  const isUserStateInit = useUserStore((s) => s.isUserStateInit);
  const hotkey = useUserStore(settingsSelectors.getHotkeyById(HotkeyEnum.Search));

  const [sessionSearchKeywords, updateSearchKeywords] = useSessionStore((s) => [
    s.sessionSearchKeywords,
    s.updateSearchKeywords,
  ]);

  const [value, setValue] = useState(sessionSearchKeywords || '');
  const debouncedValue = useDebounce(value, 500);
  const prevKeywords = useRef(sessionSearchKeywords);
  if (prevKeywords.current !== sessionSearchKeywords) {
    prevKeywords.current = sessionSearchKeywords;
    setValue(sessionSearchKeywords || '');
  }

  useEffect(() => {
    if (debouncedValue !== sessionSearchKeywords) {
      updateSearchKeywords(debouncedValue);
    }
  }, [debouncedValue, sessionSearchKeywords, updateSearchKeywords]);

  const handleChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      setValue(e.target.value);
    },
    [],
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
      value={value}
      variant={'filled'}
    />
  );
});

export default SessionSearchBar;
