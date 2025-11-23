/**
 * Backend Library Chat Service
 * 
 * This service wraps the Wails-exposed LibraryChatService from the Go backend.
 * It provides a stateless chat completion interface (OpenAI-compatible).
 * 
 * Unlike AgentChatService, this service:
 * - Does NOT maintain session state in the backend
 * - Does NOT persist messages to the database
 * - Is suitable for one-off tasks like summarization, translation, etc.
 */

import { useState } from 'react';

// Import Wails-generated bindings
import * as LibraryChatService from '@@/github.com/kawai-network/veridium/internal/llama/librarychatservice';
import { ChatCompletionRequest, type ChatCompletionResponse } from '@@/github.com/kawai-network/veridium/internal/llama/models';

class BackendLibraryChatService {
  /**
   * Send a chat completion request
   * 
   * @param params Chat completion request parameters
   * @returns Chat completion response
   */
  async chatCompletion(params: Partial<ChatCompletionRequest>): Promise<ChatCompletionResponse> {
    try {
      // Create ChatCompletionRequest instance from params
      // Ensure required fields are present or have defaults
      const request = new ChatCompletionRequest({
        messages: params.messages || [],
        model: params.model || '',
        max_tokens: params.max_tokens || 2000,
        temperature: params.temperature || 0.7,
        top_p: params.top_p || 0.95,
        stream: false, // We only support non-streaming for now via this method
        ...params
      });

      const response = await LibraryChatService.ChatCompletion(request);

      if (!response) {
        throw new Error('Backend returned null response');
      }

      return response;
    } catch (error) {
      console.error('[BackendLibraryChat] Chat completion failed:', error);
      throw error;
    }
  }

  /**
   * Helper: Check if Wails bindings are available
   */
  isAvailable(): boolean {
    return typeof LibraryChatService.ChatCompletion === 'function';
  }
}

// Export singleton instance
export const backendLibraryChat = new BackendLibraryChatService();

/**
 * React Hook for using backend library chat
 */
export function useBackendLibraryChat() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const chatCompletion = async (params: Partial<ChatCompletionRequest>): Promise<ChatCompletionResponse | null> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await backendLibraryChat.chatCompletion(params);
      return response;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    chatCompletion,
    isLoading,
    error,
    isAvailable: backendLibraryChat.isAvailable(),
  };
}
