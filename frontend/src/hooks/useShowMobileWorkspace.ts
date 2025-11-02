import { useSearchParams } from './useNavigation';

export const useShowMobileWorkspace = () => {
  const searchParams = useSearchParams();
  return searchParams.get('showMobileWorkspace') === 'true';
};
