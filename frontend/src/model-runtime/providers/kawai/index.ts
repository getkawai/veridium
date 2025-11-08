import { createOpenAICompatibleRuntime } from '../../core/openaiCompatibleFactory';

// Import Wails bindings
import { Fetch } from '@@/github.com/kawai-network/veridium/internal/llama/proxyservice';

/**
 * Custom fetch that proxies requests through Wails binding
 * Uses simple request-response for reliability
 */
const createWailsProxyFetch = () => {
  return async (url: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const urlString = typeof url === 'string' ? url : url instanceof URL ? url.toString() : String(url);
    
    // Parse URL to get path
    const urlObj = new URL(urlString);
    const path = urlObj.pathname + urlObj.search;

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

    console.debug('[Kawai] Proxying request through Wails:', {
      method: init?.method || 'GET',
      path,
      hasBody: !!body,
    });

    try {
      // Make blocking request through Wails binding
      const proxyResponse = await Fetch({
        method: init?.method || 'GET',
        path,
        headers,
        body,
      });

      if (!proxyResponse) {
        throw new Error('No response from proxy');
      }

      console.debug('[Kawai] Proxy response received:', {
        status: proxyResponse.status,
        contentLength: proxyResponse.body?.length || 0,
        contentPreview: proxyResponse.body?.substring(0, 100) || '',
      });

      // Convert SSE string to ReadableStream for OpenAI SDK
      const sseText = proxyResponse.body || '';
      const encoder = new TextEncoder();
      
      const stream = new ReadableStream({
        start(controller) {
          // Split by double newline to get complete SSE messages
          const messages = sseText.split('\n\n').filter(msg => msg.trim());
          
          console.debug('[Kawai] Converting', messages.length, 'SSE messages to stream');
          
          // Enqueue each SSE message
          for (const message of messages) {
            if (message.trim()) {
              // Add back the double newline separator
              controller.enqueue(encoder.encode(message + '\n\n'));
            }
          }
          
          controller.close();
        }
      });

      // Return Response with streaming body
      // OpenAI SDK will parse the SSE format from the stream
      return new Response(stream, {
        status: proxyResponse.status,
        statusText: proxyResponse.statusText,
        headers: new Headers(proxyResponse.headers),
      });
    } catch (error) {
      console.error('[Kawai] Proxy request failed:', error);
      throw error;
    }
  };
};

/**
 * Kawai AI Provider - Local LLM via llama.cpp
 * Uses Wails binding to proxy requests to llama-server
 */
export const LobeKawaiAI = createOpenAICompatibleRuntime({
  provider: 'kawai',
  baseURL: 'http://127.0.0.1:8080/v1', // This will be intercepted by custom fetch
  apiKey: 'sk-local', // Placeholder, not used
  
  // Inject custom fetch that uses Wails proxy
  constructorOptions: {
    fetch: createWailsProxyFetch(),
  },
  
  debug: {
    chatCompletion: () => false,
  },
});

export default LobeKawaiAI;
