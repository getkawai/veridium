import { ElectronAppState } from '@/electron-client-ipc';
import { useEffect } from 'react';
import { StateCreator } from 'zustand/vanilla';

// Import for type usage
import { electronSystemService } from '@/services/electron/system';
import { globalAgentContextManager } from '@/utils/client/GlobalAgentContextManager';
import { merge } from '@/utils/merge';

import type { ElectronStore } from '../store';

// ======== Action Interface ======== //

export interface ElectronAppAction {
  updateElectronAppState: (state: ElectronAppState) => void;

  /**
   * Initializes the basic Electron application state, including system info and special paths.
   * Should be called once when the application starts.
   */
  useInitElectronAppState: () => void;
}

// ======== Action Implementation ======== //

export const createElectronAppSlice: StateCreator<
  ElectronStore,
  [['zustand/devtools', never]],
  [],
  ElectronAppAction
> = (set, get) => ({
  updateElectronAppState: (state: ElectronAppState) => {
    const prevState = get().appState;
    set({ appState: merge(prevState, state) });
  },

  useInitElectronAppState: () => {
    useEffect(() => {
      const initAppState = async () => {
        try {
          const result = await electronSystemService.getAppState();
          set({ appState: result, isAppStateInit: true }, false, 'initElectronAppState');

          // Update the global agent context manager with relevant paths
          globalAgentContextManager.updateContext({
            desktopPath: result.userPath!.desktop,
            documentsPath: result.userPath!.documents,
            downloadsPath: result.userPath!.downloads,
            homePath: result.userPath!.home,
            musicPath: result.userPath!.music,
            picturesPath: result.userPath!.pictures,
            userDataPath: result.userPath!.userData,
            videosPath: result.userPath!.videos,
          });
        } catch (error) {
          console.error('[useInitElectronAppState] Error:', error);
        }
      };

      initAppState();
    }, []);
  },
});
