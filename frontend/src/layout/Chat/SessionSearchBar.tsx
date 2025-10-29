'use client';

import { SearchBar } from '@lobehub/ui';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';

// import { useSessionStore } from '@/store/session';
// import { useUserStore } from '@/store/user';
// import { settingsSelectors } from '@/store/user/selectors';
// import { HotkeyEnum } from '@/types/hotkey';

const SessionSearchBar = memo<{ mobile?: boolean }>(({ mobile }) => {
  const { t } = useTranslation('chat');
  // const isLoaded = useUserStore((s) => s.isLoaded);
  const isLoaded = true;
  // const hotkey = useUserStore(settingsSelectors.getHotkeyById(HotkeyEnum.Search));
  const hotkey = '';

  // const [keywords, useSearchSessions, updateSearchKeywords] = useSessionStore((s) => [
  //   s.sessionSearchKeywords,
  //   s.useSearchSessions,
  //   s.updateSearchKeywords,
  // ]);

  // const { isValidating } = useSearchSessions(keywords);
  const isValidating = false;

  // const handleChange = useCallback(
  //   (e: ChangeEvent<HTMLInputElement>) => {
  //     updateSearchKeywords(e.target.value);
  //   },
  //   [updateSearchKeywords],
  // );

  return (
    <SearchBar
      allowClear
      enableShortKey={!mobile}
      loading={!isLoaded || isValidating}
      onChange={(e) => {}}
      placeholder={t('searchAgentPlaceholder')}
      shortKey={hotkey}
      spotlight={!mobile}
      value={''}
      variant={'filled'}
    />
  );
});

export default SessionSearchBar;
