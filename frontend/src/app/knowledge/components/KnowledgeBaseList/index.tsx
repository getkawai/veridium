import React from 'react';
import { Flexbox } from 'react-layout-kit';
import { Virtuoso } from 'react-virtuoso';

import { knowledgeBaseSelectors, useKnowledgeBaseStore } from '@/store/knowledgeBase';

import KnowledgeBaseItem from '../KnowledgeBaseItem';
import EmptyStatus from './EmptyStatus';
import { SkeletonList } from './SkeletonList';

const KnowledgeBaseList = () => {
  const [data, isFetching] = useKnowledgeBaseStore((s) => [
    knowledgeBaseSelectors.knowledgeBaseList(s),
    knowledgeBaseSelectors.isFetchingList(s),
  ]);
  const fetchKnowledgeBaseList = useKnowledgeBaseStore((s) => s.fetchKnowledgeBaseList);

  React.useEffect(() => {
    fetchKnowledgeBaseList();
  }, []);

  if (isFetching) return <SkeletonList />;

  if (data?.length === 0) return <EmptyStatus />;

  return (
    <Flexbox height={'100%'}>
      <Virtuoso
        data={data}
        fixedItemHeight={36}
        itemContent={(index, item) => (
          <KnowledgeBaseItem id={item.id} key={item.id} name={item.name} />
        )}
      />
    </Flexbox>
  );
};

export default KnowledgeBaseList;
