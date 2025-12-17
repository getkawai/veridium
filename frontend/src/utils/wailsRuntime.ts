import { Events } from '@wailsio/runtime';

let isRuntimeReady = false;
let readyPromise: Promise<void> | null = null;
const readyCallbacks: (() => void)[] = [];

/**
 * Check if Wails runtime is ready.
 * The runtime sends 'common:WindowRuntimeReady' event when fully initialized.
 */
export function isWailsRuntimeReady(): boolean {
  return isRuntimeReady;
}

/**
 * Returns a promise that resolves when the Wails runtime is ready.
 * Safe to call multiple times - will return the same promise.
 */
export function waitForWailsRuntime(): Promise<void> {
  if (isRuntimeReady) {
    return Promise.resolve();
  }

  if (readyPromise) {
    return readyPromise;
  }

  readyPromise = new Promise<void>((resolve) => {
    // Check if we're running in a browser without Wails
    if (typeof window === 'undefined' || !window._wails) {
      console.warn('[WailsRuntime] Not running in Wails environment');
      isRuntimeReady = true;
      resolve();
      return;
    }

    // The runtime might already be ready if we loaded late
    // Check by trying to get the environment
    try {
      if (window._wails?.environment?.OS) {
        console.log('[WailsRuntime] Runtime already ready (environment present)');
        isRuntimeReady = true;
        resolve();
        return;
      }
    } catch (e) {
      // Ignore
    }

    // Listen for the runtime ready event
    const handleReady = () => {
      console.log('[WailsRuntime] Runtime ready event received');
      isRuntimeReady = true;
      resolve();
      readyCallbacks.forEach((cb) => cb());
      readyCallbacks.length = 0;
    };

    // Subscribe to the runtime ready event
    Events.On('common:WindowRuntimeReady', handleReady);

    // Also set a timeout in case the event was already fired before we subscribed
    setTimeout(() => {
      if (!isRuntimeReady) {
        try {
          if (window._wails?.environment?.OS) {
            console.log('[WailsRuntime] Runtime ready (detected via timeout check)');
            isRuntimeReady = true;
            resolve();
            return;
          }
        } catch (e) {
          // Ignore
        }

        // Last resort: assume ready after a short delay
        // This handles the case where we're in Wails but the event was already fired
        console.log('[WailsRuntime] Runtime ready (assumed via timeout)');
        isRuntimeReady = true;
        resolve();
      }
    }, 100);
  });

  return readyPromise;
}

/**
 * Register a callback to be called when runtime is ready.
 * If already ready, callback is invoked immediately.
 */
export function onWailsRuntimeReady(callback: () => void): void {
  if (isRuntimeReady) {
    callback();
    return;
  }
  readyCallbacks.push(callback);
  waitForWailsRuntime();
}

// Extend Window type for TypeScript
declare global {
  interface Window {
    _wails?: {
      environment?: {
        OS?: string;
        Arch?: string;
        Debug?: boolean;
      };
      flags?: Record<string, any>;
      invoke?: (message: string) => void;
    };
  }
}
