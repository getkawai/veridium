import { createOpenAICompatibleRuntime } from '../../core/openaiCompatibleFactory';

// Import Wails bindings
import { StreamFetch } from '@@/github.com/kawai-network/veridium/internal/llama/proxyservice';
import { ProxyRequest } from '@@/github.com/kawai-network/veridium/internal/llama/models';
import { Events } from '@wailsio/runtime';
import { nanoid } from 'nanoid';
import { WailsEvent } from 'node_modules/@wailsio/runtime/types/events';

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

    // Setup event listeners BEFORE creating stream to ensure they're ready
    let isStreamClosed = false;
    let streamEnded = false;
    const encoder = new TextEncoder();
    const chunkBuffer: string[] = [];

    // Listen for response metadata
    const unsubMeta = Events.On(`stream:${requestID}:meta`, (ev: WailsEvent) => {
      // Wails wraps event data in an array, extract first element
      const metadata = Array.isArray(ev.data) ? ev.data[0] : ev.data;
      console.debug('[Kawai] Stream meta received:', metadata);
    });

    // Buffer to accumulate SSE messages
    let sseMessageBuffer = '';
    
    // Listen for data chunks - accumulate into complete SSE messages
    const unsubData = Events.On(`stream:${requestID}:data`, (ev: WailsEvent) => {
      // Wails wraps event data in an array, extract first element
      const chunk = Array.isArray(ev.data) ? ev.data[0] : ev.data;
      
      if (chunk) {
        const chunkStr = typeof chunk === 'string' ? chunk : String(chunk);
        sseMessageBuffer += chunkStr;
        
        // Check if we have a complete SSE message (ends with \n\n)
        const messages = sseMessageBuffer.split('\n\n');
        
        // Last element might be incomplete, keep it in buffer
        sseMessageBuffer = messages.pop() || '';
        
        // Add complete messages to chunk buffer
        for (const msg of messages) {
          if (msg.trim()) {
            // Add back the \n\n separator
            chunkBuffer.push(msg + '\n\n');
          }
        }
      }
    });

    // Listen for stream end - mark as complete
    const unsubEnd = Events.On(`stream:${requestID}:end`, (ev: WailsEvent) => {
      console.debug('[Kawai] Stream end event received, buffered chunks:', chunkBuffer.length);
      streamEnded = true;
    });

    // Create ReadableStream with pull-based approach
    let currentIndex = 0;
    
    const stream = new ReadableStream({
      async start(controller) {
        console.debug('[Kawai] Stream started, waiting for backend...');
        
        // Start the stream request via Wails binding
        const request: ProxyRequest = {
          method: init?.method || 'GET',
          path,
          headers,
          body,
        };

        // Start fetching
        StreamFetch(requestID, request).catch((error) => {
          console.error('[Kawai] Stream error:', error);
          if (!isStreamClosed) {
            controller.error(error);
            isStreamClosed = true;
          }
        });

        // Wait for stream to end (all chunks buffered)
        while (!streamEnded) {
          await new Promise(resolve => setTimeout(resolve, 50));
        }

        console.debug('[Kawai] Backend stream ended, buffered', chunkBuffer.length, 'chunks');
      },
      
      async pull(controller) {
        // Pull-based: only enqueue when SDK requests data
        if (currentIndex < chunkBuffer.length) {
          const chunkStr = chunkBuffer[currentIndex];
          currentIndex++;
          
          try {
            controller.enqueue(encoder.encode(chunkStr));
            console.debug(`[Kawai] Pulled chunk ${currentIndex}/${chunkBuffer.length}`);
          } catch (error) {
            console.error('[Kawai] Failed to enqueue chunk:', error);
            isStreamClosed = true;
            controller.error(error);
          }
        } else if (streamEnded && currentIndex >= chunkBuffer.length) {
          // All chunks sent, close stream
          console.debug('[Kawai] All chunks sent, closing stream');
          if (!isStreamClosed) {
            try {
              controller.close();
              isStreamClosed = true;
            } catch (error) {
              console.error('[Kawai] Failed to close stream:', error);
            }
          }
          // Cleanup listeners
          unsubMeta();
          unsubData();
          unsubEnd();
        }
      },
      
      cancel(reason) {
        console.debug('[Kawai] Stream cancelled by consumer:', reason);
        isStreamClosed = true;
        unsubMeta();
        unsubData();
        unsubEnd();
      },
    });

    // Return Response with streaming body in SSE format
    // OpenAI SDK will parse the SSE format internally
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
