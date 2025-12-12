import { fileManagerSelectors, useFileStore } from '@/store/file';

/**
 * Hook to get the active file modal ID from Zustand store
 * Replaces previous URL search params implementation
 */
export const useFileModalId = (): string | undefined => {
  return useFileStore(fileManagerSelectors.activeFileId);
};

/**
 * Hook to set the file modal ID in the Zustand store
 */
export const useSetFileModalId = () => {
  return useFileStore((s) => s.setActiveFileId);
};

/**
 * Standalone function to set file modal ID
 * Now just returns the store action
 */
export const createSetFileModalId = (/* unused arg */) => {
  return (id?: string) => {
    useFileStore.getState().setActiveFileId(id);
  };
};
