import { createOpenAICompatibleRuntime } from '../../core/openaiCompatibleFactory';

// Import Wails bindings
import { StreamFetch } from '@@/github.com/kawai-network/veridium/internal/llama/proxyservice';
import { ProxyRequest } from '@@/github.com/kawai-network/veridium/internal/llama/models';
import { Events } from '@wailsio/runtime';
import { nanoid } from 'nanoid';

/**
 * Custom fetch that uses Wails events for real-time streaming
 * This allows WebView to receive SSE chunks as they arrive from llama-server
 */
const createWailsStreamingFetch = () => {
  return async (url: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const urlString = typeof url === 'string' ? url : url instanceof URL ? url.toString() : String(url);
    
    // Parse URL to get path
    const urlObj = new URL(urlString);
    const path = urlObj.pathname + urlObj.search;

    // Generate unique request ID for this stream
    const requestID = nanoid();

    // Convert headers to plain object
    const headers: Record<string, string> = {};
    if (init?.headers) {
      if (init.headers instanceof Headers) {
        init.headers.forEach((value, key) => {
          headers[key] = value;
        });
      } else if (Array.isArray(init.headers)) {
        init.headers.forEach(([key, value]) => {
          headers[key] = value;
        });
      } else {
        Object.entries(init.headers).forEach(([key, value]) => {
          headers[key] = value;
        });
      }
    }

    // Get body as string
    let body = '';
    if (init?.body) {
      if (typeof init.body === 'string') {
        body = init.body;
      } else if (init.body instanceof ArrayBuffer) {
        body = new TextDecoder().decode(init.body);
      } else if (init.body instanceof Blob) {
        body = await init.body.text();
      } else {
        body = String(init.body);
      }
    }

    console.debug('[Kawai] Starting streaming request:', {
      requestID,
      method: init?.method || 'GET',
      path,
      hasBody: !!body,
    });

    // Create ReadableStream that listens to Wails events
    const stream = new ReadableStream({
      start(controller) {
        const encoder = new TextEncoder();
        let responseMeta: any = null;

        // Listen for response metadata
        const unsubMeta = Events.On(`stream:${requestID}:meta`, (ev: any) => {
          console.debug('[Kawai] Stream meta received:', ev.data);
          responseMeta = ev.data;
        });

        // Listen for data chunks
        const unsubData = Events.On(`stream:${requestID}:data`, (ev: any) => {
          // Enqueue each SSE line as it arrives
          const chunk = ev.data as string;
          controller.enqueue(encoder.encode(chunk));
        });

        // Listen for stream end
        const unsubEnd = Events.On(`stream:${requestID}:end`, (ev: any) => {
          console.debug('[Kawai] Stream ended');
          controller.close();
          
          // Cleanup listeners
          unsubMeta();
          unsubData();
          unsubEnd();
        });

        // Start the stream request via Wails binding
        const request: ProxyRequest = {
          method: init?.method || 'GET',
          path,
          headers,
          body,
        };

        StreamFetch(requestID, request).catch((error) => {
          console.error('[Kawai] Stream error:', error);
          controller.error(error);
          unsubMeta();
          unsubData();
          unsubEnd();
        });
      },
    });

    // Return Response with streaming body
    return new Response(stream, {
      status: 200,
      statusText: 'OK',
      headers: new Headers({ 
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive',
      }),
    });
  };
};

/**
 * Kawai AI Provider - Local LLM via llama.cpp
 * Uses Wails events for real-time streaming from llama-server
 */
export const LobeKawaiAI = createOpenAICompatibleRuntime({
  provider: 'kawai',
  baseURL: 'http://127.0.0.1:8080/v1', // This will be intercepted by custom fetch
  apiKey: 'sk-local', // Placeholder, not used
  
  // Inject custom fetch that uses Wails event streaming
  constructorOptions: {
    fetch: createWailsStreamingFetch(),
  },
  
  debug: {
    chatCompletion: () => false,
  },
});

export default LobeKawaiAI;
