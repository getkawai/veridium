import { DispatchInvoke, type ProxyTRPCRequestParams } from './types';

// Import Wails runtime APIs for fallback when Electron is not available
// These will be undefined if Wails runtime is not available
let Window: any;
let Browser: any;

// Try to import Wails runtime - will be undefined if not available
try {
  // Dynamic import to avoid build errors when Wails is not available
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const wailsRuntime = require('@wailsio/runtime');
  if (wailsRuntime) {
    Window = wailsRuntime.Window;
    Browser = wailsRuntime.Browser;
  }
} catch (e) {
  // Wails runtime not available, will use Electron fallback or show error
  // This is expected in non-Wails environments
}

interface StreamerCallbacks {
  onData: (chunk: Uint8Array) => void;
  onEnd: () => void;
  onError: (error: Error) => void;
  onResponse: (response: {
    headers: Record<string, string>;
    status: number;
    statusText: string;
  }) => void;
}

interface IElectronAPI {
  invoke: DispatchInvoke;
  onStreamInvoke: (params: ProxyTRPCRequestParams, callbacks: StreamerCallbacks) => () => void;
}

declare global {
  interface Window {
    electronAPI: IElectronAPI;
  }
}

/**
 * Map of events that can be handled by Wails API when Electron is not available
 */
const WAILS_EVENT_HANDLERS: Record<string, (...args: any[]) => Promise<any>> = {
  closeWindow: async () => {
    if (Window) {
      await Window.Close();
      return;
    }
    throw new Error('Window API not available');
  },
  maximizeWindow: async () => {
    if (Window) {
      await Window.Maximise();
      return;
    }
    throw new Error('Window API not available');
  },
  minimizeWindow: async () => {
    if (Window) {
      await Window.Minimise();
      return;
    }
    throw new Error('Window API not available');
  },
  openExternalLink: async (url: string) => {
    if (Browser) {
      await Browser.OpenURL(url);
      return;
    }
    throw new Error('Browser API not available');
  },
  // Locale and theme updates are handled in frontend, no backend action needed in Wails
  updateLocale: async (locale: string) => {
    // In Wails, locale is managed in frontend only
    // Return success to avoid errors, but no backend action is needed
    console.debug('[electron-client-ipc] updateLocale called (Wails - no-op):', locale);
    return { success: true };
  },
  updateThemeMode: async (themeMode: string) => {
    // In Wails, theme is managed in frontend only
    // Return success to avoid errors, but no backend action is needed
    console.debug('[electron-client-ipc] updateThemeMode called (Wails - no-op):', themeMode);
    return;
  },
};

/**
 * client 端请求 main 端 event 数据的方法
 * 
 * Note: This function supports both Electron IPC and Wails API fallback.
 * For supported events, it will use Wails API when Electron is not available.
 */
export const dispatch: DispatchInvoke = async (event, ...data) => {
  // Check if we're in an Electron environment
  if (window.electronAPI && window.electronAPI.invoke) {
    return window.electronAPI.invoke(event, ...data);
  }

  // Check if this event can be handled by Wails API
  const wailsHandler = WAILS_EVENT_HANDLERS[event as string];
  if (wailsHandler) {
    try {
      return await wailsHandler(...data);
    } catch (error) {
      console.warn(`[electron-client-ipc] Wails handler failed for event "${event}":`, error);
      // Fall through to error message below
    }
  }

  // Event not supported by Wails or Wails not available
  const error = new Error(
    `Electron IPC is not available. This feature requires Electron environment.\n` +
    `Event: ${event}\n` +
    `If you're using Wails, this Electron-specific feature is not supported.`
  );
  
  // Log a warning instead of throwing (to avoid unhandled promise rejections)
  console.warn('[electron-client-ipc]', error.message);
  
  // Cast to satisfy TypeScript - the promise will be rejected anyway
  return Promise.reject(error) as any;
};
