import { LobeRuntimeAI } from '../BaseAI';
import type { ChatStreamPayload, ChatMethodOptions } from '../../types';
import { Fetch } from '@@/github.com/kawai-network/veridium/internal/llama/proxyservice';
import { handleChatStream } from './chat';
import { OpenAIStream } from '../streams';
import { StreamingResponse } from '../../utils/response';

export interface WailsCompatibleFactoryOptions {
  provider: string;
  baseURL?: string;
  debug?: {
    chatCompletion?: () => boolean;
  };
}

/**
 * Create a Wails-native runtime that uses direct Wails bindings
 * instead of OpenAI SDK. This provides better control over SSE parsing
 * and streaming behavior.
 */
export const createWailsCompatibleRuntime = ({
  provider,
  baseURL = 'http://127.0.0.1:8080/v1',
  debug: debugParams,
}: WailsCompatibleFactoryOptions) => {
  
  return class LobeWailsCompatibleAI implements LobeRuntimeAI {
    baseURL: string;
    private provider: string;
    private debug: boolean;

    constructor(options: any = {}) {
      this.baseURL = options.baseURL || baseURL;
      this.provider = provider;
      this.debug = debugParams?.chatCompletion?.() || false;
      
      if (this.debug) {
        console.debug(`[${this.provider}] Wails-native runtime initialized:`, {
          baseURL: this.baseURL,
        });
      }
    }

    async chat(payload: ChatStreamPayload, options?: ChatMethodOptions): Promise<Response> {
      if (this.debug) {
        console.debug(`[${this.provider}] Chat request:`, {
          model: payload.model,
          messagesCount: payload.messages.length,
          stream: payload.stream ?? true,
          temperature: payload.temperature,
        });
      }

      try {
        // Build request body in OpenAI-compatible format
        const requestBody = {
          model: payload.model,
          messages: payload.messages,
          stream: payload.stream ?? true,
          ...(payload.temperature !== undefined && { temperature: payload.temperature }),
          ...(payload.top_p !== undefined && { top_p: payload.top_p }),
          ...(payload.max_tokens !== undefined && { max_tokens: payload.max_tokens }),
          ...(payload.frequency_penalty !== undefined && { frequency_penalty: payload.frequency_penalty }),
          ...(payload.presence_penalty !== undefined && { presence_penalty: payload.presence_penalty }),
          ...(payload.stream && { stream_options: { include_usage: true } }),
        };

        if (this.debug) {
          console.debug(`[${this.provider}] Request body:`, requestBody);
        }

        // Direct Wails binding call - no OpenAI SDK, no custom fetch
        if (this.debug) {
          console.debug(`[${this.provider}] Calling Wails Fetch binding...`);
        }
        
        const proxyResponse = await Fetch({
          method: 'POST',
          path: '/v1/chat/completions',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(requestBody),
        }).catch((error) => {
          console.error(`[${this.provider}] Wails Fetch error:`, error);
          throw error;
        });

        if (!proxyResponse) {
          console.error(`[${this.provider}] No response from Wails proxy`);
          throw new Error('No response from Wails proxy');
        }

        if (this.debug) {
          console.debug(`[${this.provider}] Response received:`, {
            status: proxyResponse.status,
            bodyLength: proxyResponse.body?.length || 0,
            bodyPreview: proxyResponse.body?.substring(0, 100) || '',
          });
        }

        // Convert SSE string to ReadableStream of parsed objects
        const parsedStream = handleChatStream(proxyResponse, this.debug);
        
        // Process through OpenAIStream to convert to protocol format
        // OpenAIStream expects a ReadableStream of ChatCompletionChunk objects
        const processedStream = OpenAIStream(parsedStream, {
          callbacks: options?.callback,
          inputStartAt: Date.now(),
        });
        
        // Return as StreamingResponse
        return StreamingResponse(processedStream, {
          headers: options?.headers,
        });
        
      } catch (error) {
        console.error(`[${this.provider}] Chat error:`, error);
        throw error;
      }
    }

    async models(): Promise<any> {
      // TODO: Implement via Wails binding to GetAvailableModels
      if (this.debug) {
        console.debug(`[${this.provider}] Models list requested (not implemented yet)`);
      }
      return [];
    }

    // Optional methods - can be implemented later as needed
    async embeddings(): Promise<any> {
      throw new Error('Embeddings not implemented for Wails runtime');
    }

    async generateObject(): Promise<any> {
      throw new Error('Generate object not implemented for Wails runtime');
    }
  };
};
