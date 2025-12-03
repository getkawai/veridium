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
import { ChatRequest, type UIChatMessage } from '@@/github.com/kawai-network/veridium/internal/services/models';

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
   * @returns UIChatMessage with all message data (tools, reasoning, usage, etc.)
   * @throws Error if backend returns an error
   */
  async sendMessage(params: Partial<ChatRequest>): Promise<UIChatMessage> {
    try {
      // Create ChatRequest instance from params
      const request = new ChatRequest(params);
      
      const response = await AgentChatService.Chat(request);

      // Check if response is null or has error
      if (!response) {
        throw new Error('Backend returned null response');
      }
      
      if (response.error) {
        throw new Error(response.error.message || 'Unknown error');
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
   * This calls the backend ChatMock method which returns all messages created:
   * [userMessage, assistantMessage, toolMessage1, toolMessage2, ...]
   * 
   * All messages are saved to database for realistic data flow testing.
   * 
   * @param params Chat request parameters
   * @returns Array of UIChatMessage (user + assistant + tool messages)
   */
  async sendMessageMock(params: Partial<ChatRequest>): Promise<UIChatMessage[]> {
    try {
      const request = new ChatRequest(params);
      
      const messages = await AgentChatService.ChatMock(request);

      if (!messages || messages.length === 0) {
        throw new Error('Backend returned empty response');
      }

      return messages as unknown as UIChatMessage[];
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
   * @returns UIChatMessage response
   */
  async sendMessageStream(params: Partial<ChatRequest>): Promise<UIChatMessage> {
    console.warn(
      '[BackendAgentChat] Streaming not yet implemented, using synchronous chat',
    );
    return this.sendMessage(params);
  }

  /**
   * Send a mock message with streaming events
   * 
   * This calls the backend ChatMockStream method which emits events via Wails:
   * - 'start': Generation started
   * - 'reasoning': Thinking content (streamed)
   * - 'chunk': Content chunks (streamed word by word)  
   * - 'tool_call': Tool call initiated
   * - 'tool_result': Tool execution result with pluginState
   * - 'complete': Generation finished
   * 
   * Frontend listens via Events.On('chat:stream', handler) in App.tsx
   * Data is saved to database at the end.
   * 
   * @param params Chat request parameters
   * @returns void - all data comes via events
   */
  async sendMessageMockStream(params: Partial<ChatRequest>): Promise<void> {
    try {
      const request = new ChatRequest(params);
      
      await AgentChatService.ChatMockStream(request);
      
      console.log('[BackendAgentChat] Mock stream completed');
    } catch (error) {
      console.error('[BackendAgentChat] Mock stream failed:', error);
      throw error;
    }
  }

  /**
   * Send a real message with streaming events (REAL LLM)
   * 
   * This calls the backend ChatRealStream method which:
   * - Calls real LLM (e.g., Llama, Qwen, etc.)
   * - Executes real tools
   * - Emits events via Wails for streaming UI:
   *   - 'start': Generation started
   *   - 'reasoning': Thinking content (streamed for reasoning models)
   *   - 'chunk': Content chunks (streamed token by token)  
   *   - 'tool_call': Tool call initiated by LLM
   *   - 'tool_result': Tool execution result with pluginState
   *   - 'complete': Generation finished with full metadata
   * 
   * Frontend listens via Events.On('chat:stream', handler) in App.tsx
   * Data is saved to database at the end.
   * 
   * @param params Chat request parameters
   * @returns void - all data comes via events
   */
  async sendMessageRealStream(params: Partial<ChatRequest>): Promise<void> {
    try {
      const request = new ChatRequest(params);
      
      await AgentChatService.ChatRealStream(request);
      
      console.log('[BackendAgentChat] Real stream completed');
    } catch (error) {
      console.error('[BackendAgentChat] Real stream failed:', error);
      throw error;
    }
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

  const sendMessage = async (params: Partial<ChatRequest>): Promise<UIChatMessage | null> => {
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

