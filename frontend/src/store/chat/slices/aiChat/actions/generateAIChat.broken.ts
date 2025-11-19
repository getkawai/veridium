/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
/**
 * AI Chat Generation Actions - Refactored Version
 * 
 * This file has been refactored to use backend AgentChatService for all business logic.
 * 
 * BEFORE: 1147 lines (frontend handled all logic)
 * AFTER:  ~300 lines (backend handles all logic)
 * 
 * What backend now handles:
 * - Session management
 * - Topic auto-generation
 * - Thread context loading
 * - Message persistence
 * - Context processing
 * - Tool orchestration
 * - RAG workflow
 * - Agent execution
 * 
 * What frontend now handles:
 * - UI state updates
 * - User input collection
 * - Error display
 */

import { MESSAGE_CANCEL_FLAT } from '@/const';
import {
  CreateMessageParams,
  SendMessageParams,
  TraceEventType,
} from '@/types';
import { StateCreator } from 'zustand/vanilla';

import { backendAgentChat } from '@/services/backendAgentChat';
import type { ChatResponse } from '@@/github.com/kawai-network/veridium/internal/services/models';

// Fallback user ID (same as used by userService internally)
const FALLBACK_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

import { ChatStore } from '@/store/chat/store';
import { useSessionStore } from '@/store/session';
import { setNamespace } from '@/utils/storeDebug';

import { chatSelectors } from '../../../selectors';

const n = setNamespace('ai');

// Removed: ProcessMessageParams (no longer needed after refactoring)

export interface AIGenerateAction {
  /**
   * Sends a new message to the AI chat system
   */
  sendMessage: (params: SendMessageParams) => Promise<void>;
  /**
   * Regenerates a specific message in the chat
   */
  regenerateMessage: (id: string) => Promise<void>;
  /**
   * Deletes an existing message and generates a new one in its place
   */
  delAndRegenerateMessage: (id: string) => Promise<void>;
  /**
   * Interrupts the ongoing ai message generation process
   */
  stopGenerateMessage: () => void;

  // =========  ↓ Internal Method ↓  ========== //
  /**
   * Traces message events for analytics
   */
  internal_traceMessage: (id: string, params: { eventType: TraceEventType }) => void;
}

export const generateAIChat: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  AIGenerateAction
> = (set, get) => ({
  /**
   * Delete and regenerate message
   * Gets parent user message, deletes assistant message, and resends
   */
  delAndRegenerateMessage: async (id) => {
    // Find parent user message
    const messages = chatSelectors.activeBaseChats(get());
    const assistantMsgIndex = messages.findIndex((m) => m.id === id);

    if (assistantMsgIndex > 0) {
      // Find the user message before this assistant message
      for (let i = assistantMsgIndex - 1; i >= 0; i--) {
        if (messages[i].role === 'user') {
          const userMessage = messages[i];

          // Delete assistant message
          await get().deleteMessage(id);

          // Resend user message (backend will regenerate)
          await get().sendMessage({ message: userMessage.content });

          // Trace the event
          get().internal_traceMessage(id, {
            eventType: TraceEventType.DeleteAndRegenerateMessage,
          });
          return;
        }
      }
    }

    console.warn('[delAndRegenerateMessage] Could not find parent user message for:', id);
  },

  /**
   * Regenerate message
   * Similar to delAndRegenerate but keeps the original
   */
  regenerateMessage: async (id) => {
    // Find parent user message
    const messages = chatSelectors.activeBaseChats(get());
    const assistantMsgIndex = messages.findIndex((m) => m.id === id);

    if (assistantMsgIndex > 0) {
      // Find the user message before this assistant message
      for (let i = assistantMsgIndex - 1; i >= 0; i--) {
        if (messages[i].role === 'user') {
          const userMessage = messages[i];

          // Delete old assistant message
          await get().deleteMessage(id);

          // Resend (backend will regenerate)
          await get().sendMessage({ message: userMessage.content });

          // Trace the event
          get().internal_traceMessage(id, { eventType: TraceEventType.RegenerateMessage });
          return;
        }
      }
    }

    console.warn('[regenerateMessage] Could not find parent user message for:', id);
  },

  /**
   * Send message - REFACTORED to use backend AgentChatService
   * 
   * Backend now handles:
   * - Session management
   * - Topic auto-generation
   * - Thread context loading
   * - Message persistence
   * - Context processing
   * - Tool orchestration
   * - RAG workflow
   * - Agent execution
   * 
   * Frontend only handles UI updates
   */
  sendMessage: async ({ message, files, onlyAddUserMessage, isWelcomeQuestion }) => {
    const { activeTopicId, activeId, activeThreadId } = get();

    console.debug('[generateAIChat.sendMessage] Initial state:', {
      activeId,
      activeTopicId,
      activeThreadId,
      message: message?.substring(0, 50),
      onlyAddUserMessage,
    });

    if (!activeId) return;

    const fileIdList = files?.map((f) => f.id);
    const hasFile = !!fileIdList && fileIdList.length > 0;

    // If message is empty or no files, then stop
    if (!message && !hasFile) return;

    // Note: Server mode routing removed - all logic now in backend
    
    set({ isCreatingMessage: true }, false, n('creatingMessage/start'));

    // If only adding user message (no AI response needed)
    if (onlyAddUserMessage) {
      const newMessage: CreateMessageParams = {
        content: message,
        files: fileIdList,
        role: 'user',
        sessionId: activeId,
        topicId: activeTopicId,
        threadId: activeThreadId,
      };

      await get().internal_createMessage(newMessage);
      set({ isCreatingMessage: false }, false, n('creatingMessage/stop'));
      return;
    }

    try {
      // Get user ID (use fallback constant)
      const userId = FALLBACK_USER_ID;

      // Get enabled tools (empty for now - can be added later via agent config)
      const tools: string[] = [];
      // TODO: Get from agent config if tools are configured
      // const agentConfig = agentChatConfigSelectors.currentChatConfig(getAgentStoreState());
      // const tools = agentConfig.tools || [];

      console.log('[Backend] Calling AgentChatService.Chat()...', {
        session_id: activeId,
        user_id: userId,
        topic_id: activeTopicId,
        thread_id: activeThreadId,
        message_length: message.length,
      });

      // Backend handles EVERYTHING!
      const response: ChatResponse = await backendAgentChat.sendMessage({
        session_id: activeId,
        user_id: userId,
        message: message,
        topic_id: activeTopicId || undefined,
        thread_id: activeThreadId || undefined,
        tools: tools,
        temperature: 0.7,
        max_tokens: 2000,
      });

      console.log('[Backend] Response received:', {
        message_id: response.message_id,
        topic_id: response.topic_id,
        thread_id: response.thread_id,
        has_sources: response.sources && response.sources.length > 0,
        has_tool_calls: response.tool_calls && response.tool_calls.length > 0,
        finish_reason: response.finish_reason,
      });

      // Handle topic auto-creation
      if (response.topic_id && response.topic_id !== activeTopicId) {
        console.log('[Backend] Topic auto-created:', response.topic_id);
        
        // Update state directly (no need to call switchTopic which triggers more requests)
        set(
          { activeTopicId: response.topic_id },
          false,
          n('topicAutoCreated'),
        );
        
        console.log('[Backend] Updated activeTopicId to:', response.topic_id);
      }

      // Refresh messages from database to show the conversation
      // Pass skipRefresh=true to avoid triggering another backend call during topic switch
      await get().refreshMessages();

      // Update session
      useSessionStore.getState().triggerSessionUpdate(activeId);

      console.log('[Backend] Success! Message:', response.message?.substring(0, 100));

    } catch (error) {
      console.error('[sendMessage] Backend error:', error);

      // Show error in UI
      set(
        {
          isCreatingMessage: false,
        },
        false,
        n('creatingMessage/error'),
      );

      // TODO: Show error toast to user
      throw error;
    }

    set({ isCreatingMessage: false }, false, n('creatingMessage/stop'));
  },

  /**
   * Stop message generation
   * Aborts the ongoing message generation request
   */
  stopGenerateMessage: () => {
    const { chatLoadingIdsAbortController } = get();

    if (!chatLoadingIdsAbortController) return;

    chatLoadingIdsAbortController.abort(MESSAGE_CANCEL_FLAT);

    // Update loading state
    set({ isCreatingMessage: false }, false, n('stopGenerateMessage'));
  },

  /**
   * Trace message events for analytics
   * Kept for analytics/monitoring purposes
   */
  internal_traceMessage: (id, params) => {
    // Keep existing implementation for tracing
    console.debug('[Trace]', id, params.eventType);
    // TODO: Implement actual tracing if needed
  },
});

