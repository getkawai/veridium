import { fileManagerSelectors, useFileStore } from '@/store/file';

/**
 * Hook to manage file category filter in Zustand store
 * Replaces previous URL search params implementation
 */
export const useFileCategory = (): [string, (value: string) => void] => {
  const category = useFileStore(fileManagerSelectors.currentCategory);
  const setCategory = useFileStore((s) => s.setCategory);

  return [category, setCategory];
};
