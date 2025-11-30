/**
 * Backend Agent Chat Service
 * 
 * This service wraps the Wails-exposed AgentChatService from the Go backend.
 * It provides a clean TypeScript interface for chat operations.
 * 
 * ALL business logic is handled by the backend:
 * - Session management
 * - Topic auto-generation
 * - Thread branching
 * - Message persistence (SQLite)
 * - Tool calling
 * - RAG with knowledge bases
 * - Context processing
 * 
 * The frontend ONLY handles UI updates.
 */

import { useState } from 'react';

// Import Wails-generated bindings (correct way)
import * as AgentChatService from '@@/github.com/kawai-network/veridium/internal/services/agentchatservice';
import { ChatRequest, type ChatResponse } from '@@/github.com/kawai-network/veridium/internal/services/models';

class BackendAgentChatService {
  /**
   * Send a message to the agent
   * 
   * The backend handles:
   * - Session creation/loading
   * - Topic auto-creation (if first message)
   * - Thread context loading (if thread_id provided)
   * - Message persistence to database
   * - Agent execution with tools & RAG
   * - Returning all relevant IDs for frontend state management
   * 
   * @param params Chat request parameters
   * @returns Chat response with message, IDs, sources, tool calls
   * @throws Error if backend returns an error
   */
  async sendMessage(params: Partial<ChatRequest>): Promise<ChatResponse> {
    try {
      // Create ChatRequest instance from params
      const request = new ChatRequest(params);
      
      const response = await AgentChatService.Chat(request);

      // Check if response is null or has error
      if (!response) {
        throw new Error('Backend returned null response');
      }
      
      if (response.error) {
        throw new Error(response.error);
      }

      return response;
    } catch (error) {
      console.error('[BackendAgentChat] Send message failed:', error);
      throw error;
    }
  }

  /**
   * Send a mock message for testing UI flow
   * 
   * This calls the backend ChatMock method which returns a simple mock response
   * without calling the actual LLM or saving to database.
   * 
   * Note: For full mock with all UI components (reasoning, tools, chunks, etc.),
   * use the frontend mock in generateAIChat.ts instead.
   * 
   * @param params Chat request parameters
   * @returns Mock chat response
   */
  async sendMessageMock(params: Partial<ChatRequest>): Promise<ChatResponse> {
    try {
      const request = new ChatRequest(params);
      
      const response = await AgentChatService.ChatMock(request);

      if (!response) {
        throw new Error('Backend returned null response');
      }
      
      if (response.error) {
        throw new Error(response.error);
      }

      return response;
    } catch (error) {
      console.error('[BackendAgentChat] Send mock message failed:', error);
      throw error;
    }
  }

  /**
   * Stream a message to the agent (not yet implemented in backend)
   * 
   * For now, this delegates to the synchronous sendMessage method.
   * Once backend streaming is implemented via Wails events, this will be updated.
   * 
   * @param params Chat request parameters
   * @returns Chat response
   */
  async sendMessageStream(params: Partial<ChatRequest>): Promise<ChatResponse> {
    console.warn(
      '[BackendAgentChat] Streaming not yet implemented, using synchronous chat',
    );
    return this.sendMessage(params);
  }

  /**
   * Clear a session from backend cache
   * 
   * This removes the session from the in-memory cache but does NOT delete
   * from the database. The session can be reloaded from DB on next access.
   * 
   * @param sessionID Session ID to clear
   */
  async clearSession(sessionID: string): Promise<void> {
    try {
      await AgentChatService.ClearSession(sessionID);
      console.log(`[BackendAgentChat] Cleared session: ${sessionID}`);
    } catch (error) {
      console.error('[BackendAgentChat] Clear session failed:', error);
      throw error;
    }
  }

  /**
   * Helper: Check if Wails bindings are available
   * 
   * Useful for development mode detection
   */
  isAvailable(): boolean {
    return typeof AgentChatService.Chat === 'function';
  }
}

// Export singleton instance
export const backendAgentChat = new BackendAgentChatService();

/**
 * React Hook for using backend agent chat
 * 
 * Example usage:
 * ```tsx
 * const { sendMessage, isLoading, error } = useBackendAgentChat();
 * 
 * const handleSend = async () => {
 *   const response = await sendMessage({
 *     session_id: currentSessionId,
 *     user_id: currentUserId,
 *     message: userMessage,
 *     topic_id: currentTopicId,
 *     thread_id: currentThreadId,
 *   });
 *   
 *   // Update UI with response.message_id, response.topic_id, etc.
 * };
 * ```
 */
export function useBackendAgentChat() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const sendMessage = async (params: Partial<ChatRequest>): Promise<ChatResponse | null> => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await backendAgentChat.sendMessage(params);
      return response;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error';
      setError(errorMessage);
      return null;
    } finally {
      setIsLoading(false);
    }
  };

  const clearSession = async (sessionID: string) => {
    try {
      await backendAgentChat.clearSession(sessionID);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error';
      setError(errorMessage);
    }
  };

  return {
    sendMessage,
    clearSession,
    isLoading,
    error,
    isAvailable: backendAgentChat.isAvailable(),
  };
}

