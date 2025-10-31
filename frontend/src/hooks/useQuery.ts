import qs from 'query-string';
import { useMemo } from 'react';

import { useSearchParams } from './useNavigation';

export const useQuery = () => {
  const rawQuery = useSearchParams();
  return useMemo(() => qs.parse(rawQuery.toString()), [rawQuery]);
};
