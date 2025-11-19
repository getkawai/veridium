import { LobeRuntimeAI } from '../BaseAI';
import type { ChatStreamPayload, ChatMethodOptions } from '../../types';
import { ChatCompletion, ChatCompletionStream } from '@@/github.com/kawai-network/veridium/internal/llama/librarychatservice';
import type { ChatCompletionRequest } from '@@/github.com/kawai-network/veridium/internal/llama/models';
import { OpenAIStream } from '../streams';
import { StreamingResponse } from '../../utils/response';
import { Events } from '@wailsio/runtime';
import { OpenAIChatMessage } from '@/types';
import { WailsEvent } from 'node_modules/@wailsio/runtime/types/events';
import { transformResponseToStream } from '../openaiCompatibleFactory/nonStreamToStream';

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
        // Validate and filter messages before sending
        const validMessages = payload.messages.filter((msg: OpenAIChatMessage) => {
          if (!msg.role || !['system', 'user', 'assistant', 'tool'].includes(msg.role)) {
            console.warn(`[${this.provider}] Skipping message with invalid role:`, {
              role: msg.role,
              hasContent: !!msg.content,
              contentPreview: typeof msg.content === 'string' ? msg.content.substring(0, 50) : 'non-string',
            });
            return false;
          }
          return true;
        });

        if (validMessages.length === 0) {
          throw new Error('No valid messages to send after filtering');
        }

        // Build request in OpenAI-compatible format for LibraryChatService
        const request: ChatCompletionRequest = {
          model: payload.model,
          messages: validMessages.map((msg: OpenAIChatMessage) => ({
            role: msg.role,
            content: typeof msg.content === 'string' ? msg.content : JSON.stringify(msg.content),
          })),
          stream: payload.stream ?? true,
          ...(payload.temperature !== undefined && { temperature: payload.temperature }),
          ...(payload.top_p !== undefined && { top_p: payload.top_p }),
          ...(payload.max_tokens !== undefined && { max_tokens: payload.max_tokens }),
        };

        if (this.debug) {
          console.debug(`[${this.provider}] Request:`, request);
        }

        // Check if streaming or non-streaming
        if (request.stream) {
          // Streaming mode using ChatCompletionStream
          const requestID = `chat-${Date.now()}-${Math.random().toString(36).substring(7)}`;
          
          if (this.debug) {
            console.debug(`[${this.provider}] Starting stream with ID:`, requestID);
          }

          // Create a ReadableStream to handle SSE events from Wails
          const stream = new ReadableStream({
            start: async (controller) => {
              try {
                // Listen for stream events
            let streamClosed = false;
            
            const handleData = (data: any) => {
              if (streamClosed) return;
              
              if (this.debug) {
                console.debug(`[${this.provider}] Stream data:`, typeof data, data);
              }
              
              // Data comes as array of SSE strings from Wails events
              const dataArray = Array.isArray(data) ? data : [data];
              
              for (const item of dataArray) {
                const sseString = typeof item === 'string' ? item : String(item);
                
                // Parse SSE data and enqueue
                if (sseString.startsWith('data: ')) {
                  const jsonStr = sseString.substring(6).trim();
                  if (jsonStr === '[DONE]') {
                    if (!streamClosed) {
                      streamClosed = true;
                      controller.close();
                    }
                    return;
                  }
                  try {
                    const chunk = JSON.parse(jsonStr);
                    controller.enqueue(chunk);
                  } catch (e) {
                    console.error('Failed to parse chunk:', jsonStr, e);
                  }
                }
              }
            };

            const handleEnd = () => {
              if (streamClosed) return;
              
              if (this.debug) {
                console.debug(`[${this.provider}] Stream ended`);
              }
              
              if (!streamClosed) {
                streamClosed = true;
                try {
                  controller.close();
                } catch (e) {
                  // Stream already closed, ignore
                  if (this.debug) {
                    console.debug(`[${this.provider}] Stream already closed`);
                  }
                }
              }
            };

                // Register event listeners
                Events.On(`stream:${requestID}:data`, (ev: WailsEvent) => handleData(ev.data));
                Events.On(`stream:${requestID}:end`, () => handleEnd());

                // Start the stream
                await ChatCompletionStream(requestID, request);
                
              } catch (error) {
                console.error(`[${this.provider}] Stream error:`, error);
                controller.error(error);
              }
            },
          });

          // Process through OpenAIStream
          const processedStream = OpenAIStream(stream, {
            callbacks: options?.callback,
            inputStartAt: Date.now(),
          });

          return StreamingResponse(processedStream, {
            headers: options?.headers,
          });

        } else {
          // Non-streaming mode using ChatCompletion
          const response = await ChatCompletion(request);
          
          if (!response) {
            throw new Error('No response from LibraryChatService');
          }

          if (this.debug) {
            console.debug(`[${this.provider}] Response:`, response);
          }

          // Transform non-streaming response to stream format
          // This ensures callbacks (onFinish, onMessageHandle) are properly triggered
          const stream = transformResponseToStream(response as any);

          return StreamingResponse(
            OpenAIStream(stream, {
              callbacks: options?.callback,
              enableStreaming: false,
              inputStartAt: Date.now(),
            }),
            {
              headers: options?.headers,
            },
          );
        }
        
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
