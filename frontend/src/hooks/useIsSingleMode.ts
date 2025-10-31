import { useSearchParams } from './useNavigation';

/**
 * Hook to check if the current page is in single mode
 * Single mode is used for standalone windows in desktop app
 * @returns boolean indicating if the current page is in single mode
 */
export const useIsSingleMode = (): boolean => {
  const searchParams = useSearchParams();
  return searchParams.get('mode') === 'single';
};