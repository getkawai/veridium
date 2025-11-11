import { ElectronAppState, dispatch } from '@/electron-client-ipc';

/**
 * Service class for interacting with system-level information and actions.
 * Supports both Electron IPC and Wails API.
 * 
 * Window operations (close, maximize, minimize) will automatically use Wails API
 * when Electron is not available.
 */
class ElectronSystemService {
  /**
   * Fetches the application state from the main process.
   * This includes system information (platform, arch) and user-specific paths.
   * 
   * Note: In Wails environments, this requires backend implementation.
   * Falls back to basic platform detection if backend is not available.
   * 
   * @returns {Promise<ElectronAppState>} A promise that resolves with the app state.
   */
  async getAppState(): Promise<ElectronAppState> {
    try {
      // Try to get state from Electron/Wails backend
      return await dispatch('getDesktopAppState');
    } catch (error) {
      // Fallback to basic platform detection if backend is not available
      console.warn('[ElectronSystemService] getAppState failed, using fallback:', error);
      
      const platform = navigator.platform.toLowerCase();
      const isMac = platform.includes('mac');
      const isWindows = platform.includes('win');
      const isLinux = !isMac && !isWindows;
      
      return {
        platform: isMac ? 'darwin' : isWindows ? 'win32' : 'linux',
        isMac,
        isWindows,
        isLinux,
        arch: navigator.userAgent.includes('x64') || navigator.userAgent.includes('x86_64') ? 'x64' : 'arm64',
        systemAppearance: window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light',
        // User paths are not available in browser fallback
        userPath: undefined,
      };
    }
  }

  /**
   * Closes the current window.
   * Uses Wails API when Electron is not available.
   */
  async closeWindow(): Promise<void> {
    return dispatch('closeWindow');
  }

  /**
   * Maximizes the current window.
   * Uses Wails API when Electron is not available.
   */
  async maximizeWindow(): Promise<void> {
    return dispatch('maximizeWindow');
  }

  /**
   * Minimizes the current window.
   * Uses Wails API when Electron is not available.
   */
  async minimizeWindow(): Promise<void> {
    return dispatch('minimizeWindow');
  }

  /**
   * Shows a context menu.
   * Note: Context menu support in Wails requires backend implementation.
   */
  showContextMenu = async (type: string, data?: any) => {
    return dispatch('showContextMenu', { data, type });
  };
}

// Export a singleton instance of the service
export const electronSystemService = new ElectronSystemService();
