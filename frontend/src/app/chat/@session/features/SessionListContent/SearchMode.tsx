import { memo, useMemo } from 'react';
import { getNullableString } from '@/types/database';
import { LobeSessionType } from '@/types/session';

import SkeletonList from '../SkeletonList';
import SessionList from './List';
import { useServerConfigStore } from '@/store/serverConfig';
import { serverConfigSelectors } from '@/store/serverConfig/selectors';
import { useSessionStore } from '@/store/session';

const SearchMode = memo(() => {
  const [searchResults] = useSessionStore((s) => [
    s.searchResults,
  ]);

  const isMobile = useServerConfigStore(serverConfigSelectors.isMobile);

  const data = searchResults || [];
  const isLoading = false; // Add real loading state if needed from store

  const filteredData = useMemo(() => {
    if (!data) return data;

    if (isMobile) {
      return data.filter((session) => getNullableString(session.type) !== LobeSessionType.Group);
    }

    return data.filter((session) => {
      // For now, just exclude groups if we are desktop? Or maybe show everything?
      // Original logic: session.type !== 'agent' || !virtual
      // This implies showing groups AND non-virtual agents.
      // Let's just return true for all results from search for now as the backend search should handle relevance.
      return true;
    });
  }, [data, isMobile]);

  return isLoading ? (
    <SkeletonList />
  ) : (
    <SessionList dataSource={filteredData} showAddButton={false} />
  );
});

SearchMode.displayName = 'SessionSearchMode';

export default SearchMode;
