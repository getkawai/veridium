'use client';

import { SearchBar } from '@lobehub/ui';
import { useUnmount } from 'ahooks';
import { memo, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useChatStore } from '@/store/chat';
import { useServerConfigStore } from '@/store/serverConfig';

const TopicSearchBar = memo<{ onClear?: () => void }>(({ onClear }) => {
  const { t } = useTranslation('topic');

  const [tempValue, setTempValue] = useState('');
  const [searchKeyword, setSearchKeywords] = useState('');
  const mobile = useServerConfigStore((s) => s.isMobile);
  const [activeSessionId, internal_searchTopics] = useChatStore((s) => [s.activeId, s.internal_searchTopics]);

  useEffect(() => {
    internal_searchTopics(searchKeyword, activeSessionId);
  }, [searchKeyword, activeSessionId, internal_searchTopics]);

  useUnmount(() => {
    useChatStore.setState({ inSearchingMode: false, isSearchingTopic: false });
  });

  const startSearchTopic = () => {
    if (tempValue === searchKeyword) return;

    setSearchKeywords(tempValue);
    useChatStore.setState({ inSearchingMode: !!tempValue, isSearchingTopic: !!tempValue });
  };

  return (
    <SearchBar
      autoFocus
      onBlur={() => {
        if (tempValue === '') {
          onClear?.();

          return;
        }

        startSearchTopic();
      }}
      onChange={(e) => {
        setTempValue(e.target.value);
      }}
      onPressEnter={startSearchTopic}
      placeholder={t('searchPlaceholder')}
      spotlight={!mobile}
      value={tempValue}
      variant={'filled'}
    />
  );
});

export default TopicSearchBar;
