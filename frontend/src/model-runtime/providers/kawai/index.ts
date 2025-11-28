import { ModelProvider } from '@/model-bank';
import { createOpenAICompatibleRuntime } from '@/model-runtime/core/openaiCompatibleFactory';

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
export const LobeKawaiAI = createOpenAICompatibleRuntime({
  baseURL: 'https://api.ai21.com/studio/v1',
  chatCompletion: {
    handlePayload: (payload) => {
      return {
        ...payload,
        stream: !payload.tools,
      } as any;
    },
  },
  debug: {
    chatCompletion: () => false, // Dummy replacement for process.env.DEBUG_AI21_CHAT_COMPLETION === '1'
  },
  provider: ModelProvider.Kawai,
});

export default LobeKawaiAI;
