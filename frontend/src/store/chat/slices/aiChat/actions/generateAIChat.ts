/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
/**
 * BIG BANG MIGRATION: Backend-First Chat Flow
 * 
 * This is a RADICAL simplification:
 * - Frontend handles ONLY UI state
 * - Backend handles ALL business logic (LLM, tools, RAG, persistence)
 * - ~150 lines vs original 1144 lines (87% reduction)
 * 
 * What's REMOVED:
 * - internal_coreProcessMessage (460 lines)
 * - internal_fetchAIChatMessage (200 lines)  
 * - Tool orchestration (100 lines)
 * - RAG workflow (80 lines)
 * - Context engineering (50 lines)
 * - Topic auto-creation logic (70 lines)
 * 
 * What's KEPT:
 * - State management (messagesMap, activeIds)
 * - UI updates (refreshMessages, refreshTopic)
 * - Optimistic UI (temp messages)
 */

import { LOADING_FLAT } from '@/const';
import {
  SendMessageParams,
  UIChatMessage,
  ChatErrorType,
} from '@/types';
import { produce } from 'immer';
import { StateCreator } from 'zustand/vanilla';

import { ChatStore } from '@/store/chat/store';
import { useSessionStore } from '@/store/session/store';
import { messageMapKey } from '@/store/chat/utils/messageMapKey';
import { setNamespace } from '@/utils/storeDebug';
import { getResolvedUserId } from '@/utils/userId';
import { chatSelectors } from '../../../selectors';
import { idGenerator } from '@/database/utils/idGenerator';
import {
  StreamEventPayload,
  ChatRequest,
} from '@@/github.com/kawai-network/veridium/internal/services';
import { ChatRealStream } from '@@/github.com/kawai-network/veridium/internal/services/agentchatservice';

import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/slices/chat';

// Re-export StreamEventPayload for consumers
export type { StreamEventPayload };

// ================================================================
// API MODE CONFIGURATION
// ================================================================

const n = setNamespace('ai');

export interface AIGenerateAction {
  sendMessage: (params: SendMessageParams) => Promise<void>;
  regenerateMessage: (id: string) => Promise<void>;
  delAndRegenerateMessage: (id: string) => Promise<void>;
  stopGenerateMessage: () => void;
  clearSendMessageError: () => void;

  // Internal action to handle stream events from App.tsx
  internal_handleStreamEvent: (data: any) => void;
}

export const generateAIChat: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  AIGenerateAction
> = (set, get) => ({

  /**
   * MAIN FUNCTION: Send Message
   * 
   * Simplified flow:
   * 1. Validate input
   * 2. Create optimistic UI (temp messages)
   * 3. Call backend (handles everything)
   * 4. Update UI with real data
   * 5. Refresh from DB
   */
  sendMessage: async (params) => {
    const { message, files } = params;
    const {
      activeId,
      activeTopicId,
      activeThreadId,
    } = get();

    // ================================================================
    // STEP 1: VALIDATION
    // ================================================================
    if (!activeId) {
      console.error('[SendMessage] No active session');
      return;
    }

    if (!message?.trim() && (!files || files.length === 0)) {
      console.error('[SendMessage] No message content');
      return;
    }

    const threadId = activeThreadId;
    const mapKey = messageMapKey(activeId, activeTopicId);
    const userID = getResolvedUserId();

    if (!userID) {
      set(produce((state: ChatStore) => {
        if (!state.messagesMap[mapKey]) {
          state.messagesMap[mapKey] = [];
        }
        state.messagesMap[mapKey].push({
          id: `error-${Date.now()}`,
          role: 'assistant',
          content: '❌ Error: Wallet is not connected. Please connect wallet first.',
          sessionId: activeId,
          topicId: activeTopicId,
          error: {
            type: ChatErrorType.CreateMessageError,
            message: 'Wallet is not connected',
          },
          createdAt: Date.now(),
          updatedAt: Date.now(),
          meta: {},
        } as UIChatMessage);
      }), false, n('error/noWallet'));
      return;
    }

    // Extract file IDs from uploaded files (already processed by FileProcessorService)
    const fileIds = files?.map(f => f.id).filter(Boolean) || [];

    // ================================================================
    // STEP 2: CREATE OPTIMISTIC UI
    // ================================================================
    // Flag ini digunakan oleh selector isSendButtonDisabledByMessage 
    // untuk menonaktifkan tombol kirim selama proses pembuatan pesan, 
    // mencegah user mengirim pesan berulang saat yang pertama belum 
    // selesai diproses.
    set({ isCreatingMessage: true }, false, n('creatingMessage/start'));

    // Generate actual message IDs that will be registered in backend and database
    const messageUserId = idGenerator('messages');
    const messageAssistantId = idGenerator('messages');

    // Create optimistic user message
    set(produce((state: ChatStore) => {
      if (!state.messagesMap[mapKey]) {
        state.messagesMap[mapKey] = [];
      }
      state.messagesMap[mapKey].push({
        id: messageUserId,
        role: 'user',
        content: message,
        sessionId: activeId,
        topicId: activeTopicId,
        threadId,
        files: files?.map(f => f.id),
        createdAt: Date.now(),
        updatedAt: Date.now(),
        meta: {},
      } as UIChatMessage);
    }), false, n('optimistic/userMessage'));

    // Create optimistic assistant message (loading state)
    set(produce((state: ChatStore) => {
      state.messagesMap[mapKey].push({
        id: messageAssistantId,
        role: 'assistant',
        content: LOADING_FLAT,
        sessionId: activeId,
        topicId: activeTopicId,
        threadId,
        createdAt: Date.now(),
        updatedAt: Date.now(),
        loading: true,
        meta: {},
      } as UIChatMessage);
    }), false, n('optimistic/assistantMessage'));

    // ================================================================
    // STEP 3: CALL API BASED ON MODE
    // ================================================================
    try {
      // console.log('[Backend Real Stream] Starting real LLM streaming...');

      // Get enabled tools from agent store
      // Reasoning mode is now auto-detected by backend based on loaded model
      const agentState = useAgentStore.getState();
      const enabledPlugins = agentSelectors.currentAgentPlugins(agentState);

      // console.log('[Backend Real Stream] Tools:', enabledPlugins);

      // Call real streaming - uses real LLM with streaming events
      // Events are handled by internal_handleStreamEvent (called from App.tsx)
      const request = new ChatRequest({
        session_id: activeId,
        user_id: userID,
        message: message,
        topic_id: activeTopicId || undefined,
        thread_id: threadId || undefined,
        message_user_id: messageUserId,
        message_assistant_id: messageAssistantId,
        file_ids: fileIds.length > 0 ? fileIds : undefined,
        // User-controlled tools (reasoning auto-detected by backend)
        tools: enabledPlugins.length > 0 ? enabledPlugins : undefined,
      });

      await ChatRealStream(request);

      // Note: User message is also created via streaming events
      // Real LLM response comes token by token via events

      // console.log('[Backend Real Stream] Streaming complete, data came via events');

    } catch (error) {
      console.error('[SendMessage] Failed:', error);
      console.error('[SendMessage] Error details:', {
        message: error instanceof Error ? error.message : String(error),
        stack: error instanceof Error ? error.stack : undefined,
      });

      // Remove optimistic messages on error
      set(produce((state: ChatStore) => {
        state.messagesMap[mapKey] = state.messagesMap[mapKey].filter(
          (msg) => msg.id !== messageUserId && msg.id !== messageAssistantId
        );
      }), false, n('optimistic/error'));

      // Show error message
      set(produce((state: ChatStore) => {
        state.messagesMap[mapKey].push({
          id: `error-${Date.now()}`,
          role: 'assistant',
          content: `❌ Error: ${error instanceof Error ? error.message : 'Unknown error'}`,
          sessionId: activeId,
          topicId: activeTopicId,
          error: {
            type: ChatErrorType.CreateMessageError,
            message: error instanceof Error ? error.message : 'Unknown error',
          },
          createdAt: Date.now(),
          updatedAt: Date.now(),
          meta: {},
        } as UIChatMessage);
      }), false, n('error/show'));

    } finally {
      set({ isCreatingMessage: false }, false, n('creatingMessage/stop'));
    }
  },

  /**
   * Regenerate a message
   * 
   * Find parent message and resend
   */
  regenerateMessage: async (id) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    // Find parent user message  
    const messages = chatSelectors.activeBaseChats(get());
    const messageIndex = messages.findIndex((m) => m.id === id);

    if (messageIndex > 0) {
      const userMessage = messages[messageIndex - 1];
      if (userMessage.role === 'user') {
        await get().sendMessage({
          message: userMessage.content,
        });
      }
    }
  },

  /**
   * Delete and regenerate a message
   */
  delAndRegenerateMessage: async (id) => {
    await get().internal_deleteMessage(id);
    await get().regenerateMessage(id);
  },

  /**
   * Stop message generation
   * 
   * Stop ongoing chat generation
   */
  stopGenerateMessage: () => {
    set({ isCreatingMessage: false }, false, n('generating/stop'));
  },

  /**
   * Clear send message error
   * 
   * Clear the error message for the current session/topic
   */
  clearSendMessageError: () => {
    const { activeId, activeTopicId } = get();
    const mapKey = messageMapKey(activeId, activeTopicId);

    set(produce((state: ChatStore) => {
      if (state.mainSendMessageOperations[mapKey]) {
        delete state.mainSendMessageOperations[mapKey].inputSendErrorMsg;
      }
    }), false, n('error/clear'));
  },

  /**
   * Handle stream events with smart batching
   * - Batches rapid events using RAF for smooth streaming
   * - Processes complete events immediately to prevent infinite loading
   */
  internal_handleStreamEvent: (() => {
    let rafId: number | null = null;
    const pendingEvents: StreamEventPayload[] = [];
    
    const processStreamEvent = (data: StreamEventPayload) => {
      const { activeId, activeTopicId } = get();
      const currentMapKey = messageMapKey(activeId, activeTopicId);

      // Check if backend created a new topic (first message scenario)
      const newTopicId = data.topic_id;
      const isNewTopic = newTopicId && newTopicId !== activeTopicId;

      // Track if we need to update activeTopicId after set() completes
      let shouldUpdateActiveTopicId = false;
      let shouldRefreshSessions = false;
      let shouldScheduleTopicRefresh = false;

      set(produce((state: ChatStore) => {
        // Determine which mapKey to use for finding messages
        // Messages might still be in the old mapKey if topic was just created
        let mapKey = currentMapKey;
        let messages = state.messagesMap[mapKey];

        // If messages not found in current mapKey and we have a new topic,
        // try the old mapKey (without topic)
        if (!messages && isNewTopic) {
          const oldMapKey = messageMapKey(activeId, undefined);
          messages = state.messagesMap[oldMapKey];
          if (messages) {
            mapKey = oldMapKey;
          }
        }

      if (!messages) {
        console.warn('[Stream] No messages found for mapKey:', mapKey);
        return;
      }

      if (data.type === 'start') {
        // Find the loading assistant message by ID
        const msgIndex = messages.findIndex(m => m.id === data.message_id);

        if (msgIndex !== -1) {
          // console.log('[Stream] Start - found assistant message:', data.message_id);

          // Add to loading IDs for animation
          if (!state.chatLoadingIds.includes(data.message_id)) {
            state.chatLoadingIds.push(data.message_id);
          }

          // Handle new topic creation - move messages to new mapKey
          if (isNewTopic) {
            // console.log('[Stream] New topic created, moving messages:', { oldTopicId: activeTopicId, newTopicId });

            // Update topicId on all messages in current conversation
            messages.forEach(m => {
              m.topicId = newTopicId;
            });

            // Move messages to new mapKey
            const newMapKey = messageMapKey(activeId, newTopicId);
            state.messagesMap[newMapKey] = messages;

            // Clear old mapKey
            delete state.messagesMap[mapKey];

            // Mark that we need to update activeTopicId AFTER set() completes
            shouldUpdateActiveTopicId = true;

            // Add new topic to topicMaps (optimistic update)
            if (!state.topicMaps[activeId]) {
              state.topicMaps[activeId] = [];
            }
            const topicExists = state.topicMaps[activeId].some(t => t.id === newTopicId);
            if (!topicExists) {
              state.topicMaps[activeId].unshift({
                id: newTopicId,
                title: 'New Conversation', // Will be updated via chat:topic:updated event
                sessionId: activeId,
                favorite: false,
                createdAt: Date.now(),
                updatedAt: Date.now(),
              });
            }
          } else if (data.topic_id) {
            // Just update topic_id on the message
            messages[msgIndex].topicId = data.topic_id;
          }
        } else {
          console.warn('[Stream] Start - could not find assistant message:', data.message_id);
        }
      } else if (data.type === 'reasoning') {
        // Reasoning/thinking content streamed
        const msg = messages.find(m => m.id === data.message_id);
        if (msg) {
          // Support both flat full_content (legacy) and nested reasoning.content (standard UIChatMessage)
          const content = data.reasoning?.content || data.full_content || data.content;

          (msg as any).reasoning = {
            content: content,
            duration: 0,
          };
          msg.updatedAt = Date.now();
        }
      } else if (data.type === 'chunk') {
        // Content chunk streamed
        const msg = messages.find(m => m.id === data.message_id);
        if (msg) {
          // Support both full_content (legacy) and content (standard UIChatMessage)
          msg.content = data.content || data.full_content || msg.content || '';
          msg.updatedAt = Date.now();
        }
      } else if (data.type === 'tool_call') {
        // Tool call initiated - update tools array on assistant message
        // console.log('[Stream] Tool call:', data.tool?.apiName || 'unknown');

        const msg = messages.find(m => m.id === data.message_id);
        if (msg) {
          // Set/update tools array so Tool component renders
          (msg as any).tools = data.tools?.map((t: any) => ({
            id: t.id || t.ID,
            identifier: t.identifier || t.Identifier,
            apiName: t.apiName || t.APIName,
            arguments: t.arguments || t.Arguments || '{}',
            type: t.type || t.Type || 'builtin',
          })) || [];
          msg.updatedAt = Date.now();
        }
      } else if (data.type === 'tool_result') {
        // Tool execution result - add tool message to messagesMap
        // console.log('[Stream] Tool result:', data.tool_call_id);

        // Use the current activeTopicId from state (may have been updated by start event)
        const currentTopicId = state.activeTopicId;

        // Create tool message
        const toolMessage: UIChatMessage = {
          id: data.tool_msg_id,
          role: 'tool',
          content: typeof data.content === 'string' ? data.content : JSON.stringify(data.content || ''),
          tool_call_id: data.tool_call_id,
          sessionId: activeId,
          topicId: currentTopicId || newTopicId || activeTopicId,
          pluginState: data.pluginState,
          plugin: data.plugin ? {
            apiName: data.plugin.apiName,
            arguments: data.plugin.arguments,
            identifier: data.plugin.identifier,
            type: data.plugin.type,
          } : undefined,
          meta: {},
          createdAt: Date.now(),
          updatedAt: Date.now(),
        } as UIChatMessage;

        // Add to messages if not already exists
        if (!messages.find(m => m.id === toolMessage.id)) {
          messages.push(toolMessage);
        }
        } else if (data.type === 'complete') {
          // Generation finished
          // console.log('[Stream] Complete event - before processing, chatLoadingIds:', state.chatLoadingIds);
          const msg = messages.find(m => m.id === data.message_id);
          if (msg) {
            msg.content = data.content || msg.content;
            msg.updatedAt = Date.now();
            (msg as any).loading = false;

            // Update additional fields if provided
            if (data.reasoning) {
              (msg as any).reasoning = data.reasoning;
            }
            if (data.search) {
              (msg as any).search = data.search;
            }
            if (data.chunksList) {
              (msg as any).chunksList = data.chunksList;
            }
            if (data.imageList) {
              (msg as any).imageList = data.imageList;
            }
            if (data.usage) {
              (msg as any).usage = data.usage;
            }
            if (data.performance) {
              (msg as any).performance = data.performance;
            }
          }

          // Remove from loading IDs
          // const beforeLength = state.chatLoadingIds.length;
          state.chatLoadingIds = state.chatLoadingIds.filter(id => id !== data.message_id);
          // const afterLength = state.chatLoadingIds.length;
          // console.log('[Stream] Complete - removed from loading:', data.message_id, 'before:', beforeLength, 'after:', afterLength);
          // console.log('[Stream] Complete - chatLoadingIds now:', state.chatLoadingIds);

          // Mark that we need to refresh sessions AFTER set() completes
          shouldRefreshSessions = true;
          
          // Mark if we need to schedule topic refresh
          if (isNewTopic) {
            shouldScheduleTopicRefresh = true;
          }
        }
    }), false, n('streamEvent'));

    // Execute side effects AFTER state update completes to avoid infinite loops
    if (shouldUpdateActiveTopicId && newTopicId) {
      // Update activeTopicId in a separate set() call
      set({ activeTopicId: newTopicId }, false, n('updateActiveTopicId'));
    }

    if (shouldRefreshSessions) {
      // Refresh session list to update sort order (Last Active)
      // Use setTimeout to ensure this happens after current event loop
      setTimeout(() => {
        useSessionStore.getState().refreshSessions();
      }, 0);
    }

    if (shouldScheduleTopicRefresh) {
      // After complete event, schedule a topic refresh to get LLM-generated title
      // This is a fallback in case chat:topic:updated event is not received
      // console.log('[Stream] Scheduling topic refresh for new topic title...');
      setTimeout(() => {
        get().refreshTopic();
      }, 5000); // Wait 5 seconds for backend to generate title
    }
      
      // console.log('[Stream] Event processing completed for:', data.type, data.message_id);
    };
    
    const flushPendingEvents = () => {
      if (rafId !== null) {
        cancelAnimationFrame(rafId);
        rafId = null;
      }
      
      if (pendingEvents.length === 0) return;
      
      const eventsToProcess = [...pendingEvents];
      pendingEvents.length = 0;
      
      // console.log('[Stream] Processing batch of', eventsToProcess.length, 'events');
      for (const event of eventsToProcess) {
        processStreamEvent(event);
      }
    };
    
    return (data: StreamEventPayload) => {
      // console.log('[Stream] Event received:', data.type, data.message_id);
      
      // Add to pending queue
      pendingEvents.push(data);
      
      // If it's a complete event, process immediately to prevent infinite loading
      if (data.type === 'complete') {
        // console.log('[Stream] Complete event detected, flushing immediately');
        flushPendingEvents();
        return;
      }
      
      // For other events, batch with RAF
      if (rafId !== null) {
        return; // Already scheduled
      }
      
      rafId = requestAnimationFrame(() => {
        rafId = null;
        flushPendingEvents();
      });
    };
  })(),
});
