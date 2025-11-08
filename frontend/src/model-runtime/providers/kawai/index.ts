import { createOpenAICompatibleRuntime } from '../../core/openaiCompatibleFactory';

// Import Wails proxy binding
import { Fetch } from '@@/github.com/kawai-network/veridium/internal/llama/proxyservice';
import { ProxyResponse } from '@@/github.com/kawai-network/veridium/internal/llama/models';

/**
 * Custom fetch that routes through Wails proxy
 * This allows WebView to call localhost llama-server
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
      // Call Go proxy via Wails binding
      const proxyResponse: ProxyResponse | null = await Fetch({
        method: init?.method || 'GET',
        path,
        headers,
        body,
      });

      // Check if response is null
      if (!proxyResponse) {
        throw new Error('Proxy response is null - llama-server may not be running');
      }

      // Convert Go response to Web Response
      const responseHeaders = new Headers(proxyResponse.headers);
      
      return new Response(proxyResponse.body, {
        status: proxyResponse.status,
        statusText: proxyResponse.statusText,
        headers: responseHeaders,
      });
    } catch (error) {
      console.error('[Kawai] Proxy request failed:', error);
      throw error;
    }
  };
};

/**
 * Kawai AI Provider - Local LLM via llama.cpp
 * Uses Wails proxy to communicate with llama-server
 */
export const LobeKawaiAI = createOpenAICompatibleRuntime({
  provider: 'kawai',
  baseURL: 'http://127.0.0.1:8080/v1', // This will be intercepted by custom fetch
  apiKey: 'sk-local', // Placeholder, not used
  
  // Inject custom fetch that routes through Wails
  constructorOptions: {
    fetch: createWailsProxyFetch(),
  },
  
  debug: {
    chatCompletion: () => false,
  },
});

export default LobeKawaiAI;

