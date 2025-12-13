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

// Import Wails-generated bindings (correct way)
import * as AgentChatService from '@@/github.com/kawai-network/veridium/internal/services/agentchatservice';
import { ChatRequest, type UIChatMessage } from '@@/github.com/kawai-network/veridium/internal/services/models';

class BackendAgentChatService {
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
}

// Export singleton instance
export const backendAgentChat = new BackendAgentChatService();

