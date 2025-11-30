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

// import { backendAgentChat } from '@/services/backendAgentChat';
import { ChatStore } from '@/store/chat/store';
import { messageMapKey } from '@/store/chat/utils/messageMapKey';
import { setNamespace } from '@/utils/storeDebug';
import { chatSelectors } from '../../../selectors';

// User ID constant for backend calls
// const FALLBACK_CLIENT_DB_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

const n = setNamespace('ai');

export interface AIGenerateAction {
  sendMessage: (params: SendMessageParams) => Promise<void>;
  regenerateMessage: (id: string) => Promise<void>;
  delAndRegenerateMessage: (id: string) => Promise<void>;
  stopGenerateMessage: () => void;

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
      // refreshTopic,
    } = get();

    // Validation
    if (!activeId) {
      console.error('[BigBang] No active session');
      return;
    }

    if (!message?.trim() && (!files || files.length === 0)) {
      console.error('[BigBang] No message content');
      return;
    }

    const threadId = activeThreadId;

    // ================================================================
    // MAIN FLOW: User Message + AI Response
    // ================================================================

    set({ isCreatingMessage: true }, false, n('creatingMessage/start'));

    // Step 1: Create optimistic user message
    const mapKey = messageMapKey(activeId, activeTopicId);
    const tempUserId = `temp-user-${Date.now()}`;

    console.log('[MOCK] Creating messages with key:', mapKey, {
      activeId,
      activeTopicId,
      threadId,
    });

    set(produce((state: ChatStore) => {
      if (!state.messagesMap[mapKey]) {
        state.messagesMap[mapKey] = [];
      }
      console.log('[MOCK] messagesMap before push:', Object.keys(state.messagesMap));
      state.messagesMap[mapKey].push({
        id: tempUserId,
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

    // Step 2: Create optimistic assistant message (loading state)
    const tempAssistantId = `temp-assistant-${Date.now()}`;

    set(produce((state: ChatStore) => {
      state.messagesMap[mapKey].push({
        id: tempAssistantId,
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

    try {
      // ================================================================
      // MOCK RESPONSE - Simulate streaming with full data
      // ================================================================
      console.log('[MOCK] Simulating AI response...');
      
      // Simulate delay
      await new Promise(resolve => setTimeout(resolve, 500));

      // Step 3: Update assistant message with FULL mock data
      const mockResponse = `This is a mock response to: "${message}"\n\nI'm simulating the AI response to test the UI flow without calling the backend.`;
      
      set(produce((state: ChatStore) => {
        console.log('[MOCK] Updating message with key:', mapKey);
        console.log('[MOCK] messagesMap keys:', Object.keys(state.messagesMap));
        
        const messages = state.messagesMap[mapKey];
        if (!messages) {
          console.error('[MOCK] Messages not found for key:', mapKey);
          return;
        }

        console.log('[MOCK] Found messages:', messages.length);
        const assistantMsgIndex = messages.findIndex(m => m.id === tempAssistantId);
        if (assistantMsgIndex !== -1) {
          const msg = messages[assistantMsgIndex];
          
          // Update content
          msg.content = mockResponse;
          msg.updatedAt = Date.now();
          (msg as any).loading = false;

          // Mock reasoning data
          (msg as any).reasoning = {
            content: 'Let me think about this step by step:\n1. First, I need to understand the question\n2. Then, I will formulate a response\n3. Finally, I will provide a clear answer',
            status: 'complete',
          };

          // Mock RAG chunks data
          (msg as any).chunksList = [
            {
              id: 'chunk_1',
              content: 'This is a sample chunk from the knowledge base. It contains relevant information about the topic.',
              metadata: {
                source: 'document.pdf',
                page: 1,
              },
            },
            {
              id: 'chunk_2',
              content: 'Another chunk with more detailed information that was retrieved from the RAG system.',
              metadata: {
                source: 'guide.md',
                section: 'Introduction',
              },
            },
          ];

          // Mock tool calls
          (msg as any).tools = [
            {
              id: 'tool_1',
              identifier: 'search',
              apiName: 'search',
              arguments: JSON.stringify({ query: 'test query' }),
              type: 'builtin',
            },
            {
              id: 'tool_2',
              identifier: 'calculator',
              apiName: 'calculator',
              arguments: JSON.stringify({ expression: '2+2' }),
              type: 'builtin',
            },
          ];

          // Mock search grounding
          (msg as any).search = {
            citations: [
              {
                id: 'citation_1',
                title: 'Wikipedia - Example Article',
                url: 'https://en.wikipedia.org/wiki/Example',
              },
              {
                id: 'citation_2',
                title: 'GitHub Documentation',
                url: 'https://docs.github.com/en',
              },
            ],
            searchQueries: ['test query', 'related query'],
          };

          // Mock image list
          (msg as any).imageList = [
            {
              id: 'img_1',
              url: 'https://via.placeholder.com/300x200',
              alt: 'Sample image 1',
            },
          ];

          // Mock usage
          (msg as any).usage = {
            prompt_tokens: 150,
            completion_tokens: 80,
            total_tokens: 230,
          };

          // Mock performance
          (msg as any).performance = {
            total_tokens: 230,
            duration: 1500,
          };

          // Mock metadata
          (msg as any).metadata = {
            model: 'mock-model',
            temperature: 0.7,
          };
        }
      }), false, n('mock/response'));

      console.log('[MOCK] Response complete with full data');

      // COMMENTED OUT: Real backend call
      /*
      const response = await backendAgentChat.sendMessage({
        session_id: activeId,
        user_id: FALLBACK_CLIENT_DB_USER_ID,
        message: message,
        topic_id: activeTopicId || undefined,
        thread_id: threadId || undefined,
        tools: [], // TODO: Get enabled tools from agent store
        knowledge_base_id: undefined, // TODO: Get from KB state
        temperature: 0.7,
        max_tokens: 2000,
        stream: true, // Enable streaming
      });
      */

      // MOCK: No backend response to process
      console.log('[MOCK] Skipping backend response processing');

    } catch (error) {
      console.error('[BigBang] Failed:', error);

      // Remove temp messages on error
      set(produce((state: ChatStore) => {
        state.messagesMap[mapKey] = state.messagesMap[mapKey].filter(
          (msg) => !msg.id.startsWith('temp-')
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
   * Handle stream events globally
   * Called from App.tsx when a stream event is received
   */
  internal_handleStreamEvent: (data: any) => {
    const { activeId, activeTopicId } = get();
    const mapKey = messageMapKey(activeId, activeTopicId);

    set(produce((state: ChatStore) => {
      const messages = state.messagesMap[mapKey];
      if (!messages) return;

      if (data.type === 'start') {
        // Find the latest temporary assistant message
        // Support both prefixes: 'temp-assistant-' (main chat) and 'tmp_' (thread)
        let tempMsgIndex = -1;
        for (let i = messages.length - 1; i >= 0; i--) {
          const msg = messages[i];
          // Check loading state from either message property (main chat) or store state (thread)
          const isLoading = (msg as any).loading || state.messageLoadingIds.includes(msg.id);

          const isTemp = (msg.id.startsWith('temp-assistant-') || msg.id.startsWith('tmp_')) &&
            isLoading &&
            msg.role === 'assistant';

          if (isTemp) {
            tempMsgIndex = i;
            break;
          }
        }

        if (tempMsgIndex !== -1) {
          // Update ID to real ID
          messages[tempMsgIndex].id = data.message_id;
          console.log('[Stream] Linked temp message to real ID:', data.message_id);

          // Add to loading IDs for animation
          if (!state.chatLoadingIds.includes(data.message_id)) {
            state.chatLoadingIds.push(data.message_id);
          }
        } else {
          console.warn('[Stream] Could not find temp assistant message for start event');
        }
      } else if (data.type === 'chunk') {
        // Find message by REAL ID
        const msg = messages.find(m => m.id === data.message_id);
        if (msg) {
          msg.content = data.full_content;
          msg.updatedAt = Date.now();
        }
      } else if (data.type === 'tool_calling') {
        // Tool is being called - update tools array for Tool component rendering
        console.log('[Stream] Tool calling event received:', data);
        
        const msg = messages.find(m => m.id === data.message_id);
        if (msg) {
          // Set tools array so Tool component renders
          (msg as any).tools = data.tools?.map((t: any) => ({
            id: t.id,
            identifier: t.identifier,
            apiName: t.apiName,
            arguments: t.arguments || '{}',
            type: t.type || 'builtin',
          })) || [];
          msg.updatedAt = Date.now();
          
          console.log('[Stream] Updated message with tools:', msg);
        } else {
          console.warn('[Stream] Could not find message for tool_calling:', data.message_id);
        }
      } else if (data.type === 'complete') {
        const msg = messages.find(m => m.id === data.message_id);
        if (msg) {
          msg.content = data.content;
          msg.updatedAt = Date.now();
          (msg as any).loading = false;
        }

        // Remove from loading IDs
        state.chatLoadingIds = state.chatLoadingIds.filter(id => id !== data.message_id);
      }
    }), false, n('streamEvent'));
  },
});

