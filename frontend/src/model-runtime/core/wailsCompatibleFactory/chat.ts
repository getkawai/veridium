import type { ProxyResponse } from '@@/github.com/kawai-network/veridium/internal/llama/models';

/**
 * Parse SSE format line and extract JSON data
 */
function parseSSELine(line: string): any | null {
  const trimmed = line.trim();
  
  // Skip empty lines
  if (!trimmed) return null;
  
  // Handle [DONE] marker
  if (trimmed === 'data: [DONE]' || trimmed === '[DONE]') return null;
  
  // Parse data: {...} format
  if (trimmed.startsWith('data: ')) {
    const jsonStr = trimmed.substring(6).trim();
    if (jsonStr === '[DONE]') return null;
    
    try {
      return JSON.parse(jsonStr);
    } catch (error) {
      console.error('[Wails] Failed to parse SSE JSON:', jsonStr, error);
      return null;
    }
  }
  
  return null;
}

/**
 * Convert SSE response body to ReadableStream of parsed chunks
 * This handles the complete SSE response from llama-server and converts it
 * to a stream of parsed JSON objects that the UI can consume
 */
export function handleChatStream(
  proxyResponse: ProxyResponse,
  debug: boolean = false
): ReadableStream {
  const sseText = proxyResponse.body || '';
  
  // Split into lines and parse each SSE message
  const lines = sseText.split('\n');
  const chunks: any[] = [];
  
  for (const line of lines) {
    const parsed = parseSSELine(line);
    if (parsed) {
      chunks.push(parsed);
    }
  }
  
  if (debug) {
    console.debug('[Wails] Parsed', chunks.length, 'chunks from SSE response');
    if (chunks.length > 0) {
      console.debug('[Wails] First chunk:', chunks[0]);
      console.debug('[Wails] Last chunk:', chunks[chunks.length - 1]);
    }
  }
  
  // Create ReadableStream that yields parsed chunks one by one
  let currentIndex = 0;
  const stream = new ReadableStream({
    pull(controller) {
      if (currentIndex < chunks.length) {
        const chunk = chunks[currentIndex];
        currentIndex++;
        
        if (debug && currentIndex <= 3) {
          console.debug('[Wails] Streaming chunk', currentIndex, '/', chunks.length, ':', 
            chunk.choices?.[0]?.delta?.content || '[no content]');
        }
        
        // Enqueue the parsed JSON object directly
        controller.enqueue(chunk);
      } else {
        if (debug) {
          console.debug('[Wails] Stream complete, sent', chunks.length, 'chunks');
        }
        controller.close();
      }
    },
  });
  
  // Return the ReadableStream directly
  // This stream yields parsed JSON objects (ChatCompletionChunk format)
  return stream;
}
