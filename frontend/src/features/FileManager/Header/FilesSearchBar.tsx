'use client';

import { SearchBar } from '@lobehub/ui';
import { memo, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { fileManagerSelectors, useFileStore } from '@/store/file';

import { useUserStore } from '@/store/user';
import { settingsSelectors } from '@/store/user/selectors';
import { HotkeyEnum } from '@/types/hotkey';

const FilesSearchBar = memo<{ mobile?: boolean }>(({ mobile }) => {
  const { t } = useTranslation('file');
  const hotkey = useUserStore(settingsSelectors.getHotkeyById(HotkeyEnum.Search));
  const [query, setQuery] = useFileStore((s) => [
    fileManagerSelectors.searchKeywords(s),
    s.setSearchKeywords,
  ]);
  const [keywords, setKeywords] = useState<string>(query);

  // Sync local state with store state
  useEffect(() => {
    setKeywords(query || '');
  }, [query]);

  return (
    <SearchBar
      allowClear
      enableShortKey={!mobile}
      onChange={(e) => {
        setKeywords(e.target.value);
        if (!e.target.value) setQuery('');
      }}
      onPressEnter={() => setQuery(keywords)}
      placeholder={t('searchFilePlaceholder')}
      shortKey={hotkey}
      spotlight={!mobile}
      style={{ width: 320 }}
      value={keywords}
      variant={'filled'}
    />
  );
});

export default FilesSearchBar;
