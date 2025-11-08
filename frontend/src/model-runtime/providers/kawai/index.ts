import { createWailsCompatibleRuntime } from '../../core/wailsCompatibleFactory';

/**
 * Kawai AI Provider - Native Wails Implementation
 * 
 * Uses direct Wails bindings instead of OpenAI SDK for better control
 * over SSE parsing and streaming behavior. This approach:
 * 
 * 1. Calls Wails Fetch binding directly (no HTTP from browser)
 * 2. Receives complete SSE response from llama-server
 * 3. Parses SSE format ourselves
 * 4. Returns stream of parsed JSON objects
 * 5. Processes through OpenAIStream for protocol conversion
 * 
 * This is simpler and more reliable than fighting with OpenAI SDK's
 * expectations about HTTP streaming and SSE format.
 */
export const LobeKawaiAI = createWailsCompatibleRuntime({
  provider: 'kawai',
  baseURL: 'http://127.0.0.1:8080/v1', // Not actually used for HTTP, just for reference
  debug: {
    chatCompletion: () => true, // Enable debug logging
  },
});

export default LobeKawaiAI;
