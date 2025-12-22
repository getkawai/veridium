'use client';

import { SearchBar, SearchBarProps } from '@lobehub/ui';
import { createStyles } from 'antd-style';
import { memo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { withSuspense } from '@/components/withSuspense';

export const useStyles = createStyles(({ css, prefixCls, token }) => ({
  active: css`
    box-shadow: ${token.boxShadow};
  `,
  bar: css`
    .${prefixCls}-input-group-wrapper {
      padding: 0;
    }
  `,
}));

interface StoreSearchBarProps extends SearchBarProps {
  mobile?: boolean;
}

const Search = memo<StoreSearchBarProps>(() => {
  const { t } = useTranslation('discover');
  const [word, setWord] = useState<string>('');
  const handleSearch = (value: string) => {
    // TODO: implement search
  };

  return (
    <SearchBar
      data-testid="search-bar"
      defaultValue={''}
      enableShortKey
      onInputChange={(v) => {
        setWord(v);
        if (!v) handleSearch('');
      }}
      onSearch={handleSearch}
      placeholder={t('search.placeholder')}
      style={{
        width: 'min(720px,100%)',
      }}
      value={word}
      variant={'outlined'}
    />
  );
});

export default withSuspense(Search);
