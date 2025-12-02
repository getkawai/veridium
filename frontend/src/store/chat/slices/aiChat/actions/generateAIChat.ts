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

import { backendAgentChat } from '@/services/backendAgentChat';
import { ChatStore } from '@/store/chat/store';
import { messageMapKey } from '@/store/chat/utils/messageMapKey';
import { setNamespace } from '@/utils/storeDebug';
import { chatSelectors } from '../../../selectors';
import { idGenerator } from '@/database/utils/idGenerator';

// User ID constant for backend calls
const FALLBACK_CLIENT_DB_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

// ================================================================
// API MODE CONFIGURATION
// ================================================================
/**
 * API_MODE controls which backend to use:
 * 
 * - 'REAL': Production mode - calls real backend API with LLM
 * - 'BACKEND_MOCK': Development mode - calls backend mock API (saves to DB, no LLM)
 * - 'FRONTEND_MOCK': UI testing mode - frontend-only mock (no backend, no DB)
 */
type ApiMode = 'REAL' | 'BACKEND_MOCK' | 'FRONTEND_MOCK';

const API_MODE: ApiMode = 'BACKEND_MOCK' as ApiMode;

// ================================================================
// MODE DESCRIPTIONS
// ================================================================
// REAL: 
//   - Uses real LLM (OpenAI, Claude, etc)
//   - Saves to database
//   - Full tool execution
//   - Production-ready
//
// BACKEND_MOCK:
//   - No LLM calls (saves costs)
//   - Saves to database (realistic data flow)
//   - Mock tool results
//   - Good for backend integration testing
//
// FRONTEND_MOCK:
//   - No backend calls
//   - No database
//   - Instant responses
//   - Good for UI/UX development
// ================================================================

const n = setNamespace('ai');

// ================================================================
// HELPER FUNCTIONS FOR DIFFERENT API MODES
// ================================================================

/**
 * Handle REAL API mode - Production with real LLM
 */
async function handleRealAPI(
  get: any,
  set: any,
  context: {
    activeId: string;
    activeTopicId: string | undefined;
    threadId: string | undefined;
    message: string;
    messageUserId: string;
    messageAssistantId: string;
  }
) {
  const { activeId, activeTopicId, threadId, message, messageUserId, messageAssistantId } = context;

  console.log('[Real API] Calling real backend with LLM...');

  // TODO: Implement real API call
  // This will be similar to backend mock but calls the real endpoint
  const response = await backendAgentChat.sendMessage({
    session_id: activeId,
    user_id: FALLBACK_CLIENT_DB_USER_ID,
    message: message,
    topic_id: activeTopicId || undefined,
    thread_id: threadId || undefined,
    message_user_id: messageUserId,
    message_assistant_id: messageAssistantId,
  });

  console.log('[Real API] Response received:', response);

  // Refresh messages from DB only if topic changed
  const finalTopicId = response.topic_id || activeTopicId;

  // If a new topic was created, switch to it
  if (finalTopicId && finalTopicId !== activeTopicId) {
    console.log('[Real API] New topic created, switching to:', finalTopicId);
    await get().refreshTopic();
    await get().switchTopic(finalTopicId);
  } else {
    // For same topic, we rely on stream events to update the UI state
    // No need to fetch messages again
    console.log('[Real API] Staying on same topic, relying on stream updates');
  }

  console.log('[Real API] Complete');
}

/**
 * Handle BACKEND_MOCK mode - Backend mock with DB persistence
 */
async function handleBackendMock(
  get: any,
  set: any,
  context: {
    activeId: string;
    activeTopicId: string | undefined;
    threadId: string | undefined;
    message: string;
    messageUserId: string;
    messageAssistantId: string;
  }
) {
  const { activeId, activeTopicId, threadId, message, messageUserId, messageAssistantId } = context;

  console.log('[Backend Mock] Calling backend mock (saves to DB)...');
  console.log('[Backend Mock] Params:', {
    session_id: activeId,
    user_id: FALLBACK_CLIENT_DB_USER_ID,
    message: message,
    topic_id: activeTopicId,
    thread_id: threadId,
    message_user_id: messageUserId,
    message_assistant_id: messageAssistantId,
  });

  const response = await backendAgentChat.sendMessageMock({
    session_id: activeId,
    user_id: FALLBACK_CLIENT_DB_USER_ID,
    message: message,
    topic_id: activeTopicId || undefined,
    thread_id: threadId || undefined,
    message_user_id: messageUserId,
    message_assistant_id: messageAssistantId,
  });

  console.log('[Backend Mock] Response received:', response);

  // Refresh messages from DB only if topic changed
  const finalTopicId = response.topic_id || activeTopicId;

  // If a new topic was created, switch to it
  if (finalTopicId && finalTopicId !== activeTopicId) {
    console.log('[Backend Mock] New topic created, switching to:', finalTopicId);
    await get().refreshTopic();
    await get().switchTopic(finalTopicId);
  } else {
    // For same topic, we rely on stream events to update the UI state
    // No need to fetch messages again
    console.log('[Backend Mock] Staying on same topic, relying on stream updates');
  }

  console.log('[Backend Mock] Complete');
}

/**
 * Handle FRONTEND_MOCK mode - Frontend-only mock for UI testing
 */
async function handleFrontendMock(
  set: any,
  context: {
    activeId: string;
    activeTopicId: string | undefined;
    threadId: string | undefined;
    message: string;
    mapKey: string;
    messageAssistantId: string;
  }
) {
  const { activeId, activeTopicId, threadId, message, mapKey, messageAssistantId } = context;

  console.log('[Frontend Mock] Simulating AI response (no backend)...');

  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 500));

  const mockResponse = `This is a mock response to: "${message}"\n\nI'm simulating the AI response to test the UI flow without calling the backend.`;

  set(produce((state: ChatStore) => {
    const messages = state.messagesMap[mapKey];
    if (!messages) {
      console.error('[Frontend Mock] Messages not found for key:', mapKey);
      return;
    }

    const assistantMsgIndex = messages.findIndex(m => m.id === messageAssistantId);
    if (assistantMsgIndex === -1) {
      console.error('[Frontend Mock] Assistant message not found');
      return;
    }

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
        fileId: 'file_1',
        filename: 'document.pdf',
        fileType: 'application/pdf',
        fileUrl: '/files/document.pdf',
        text: 'This is a sample chunk from the knowledge base. It contains relevant information about the topic.',
        similarity: 0.95,
      },
      {
        id: 'chunk_2',
        fileId: 'file_2',
        filename: 'guide.md',
        fileType: 'text/markdown',
        fileUrl: '/files/guide.md',
        text: 'Another chunk with more detailed information that was retrieved from the RAG system.',
        similarity: 0.87,
      },
    ];

    // Mock tool calls
    (msg as any).tools = [
      {
        id: 'tool_1',
        identifier: 'lobe-web-browsing',
        apiName: 'search',
        arguments: JSON.stringify({
          query: 'What is the weather today?',
          searchEngines: ['google']
        }),
        type: 'builtin',
        result: {
          id: 'tool_result_1',
          content: JSON.stringify({
            results: [
              {
                title: 'Mock Search Result 1',
                url: 'https://example.com/result1',
                description: 'This is a mock search result for testing purposes.',
              },
              {
                title: 'Mock Search Result 2',
                url: 'https://example.com/result2',
                description: 'Another mock search result with relevant information.',
              },
            ],
          }),
          state: null,
        },
      },
      {
        id: 'tool_2',
        identifier: 'lobe-local-system',
        apiName: 'listLocalFiles',
        arguments: JSON.stringify({
          path: '/home/user/documents'
        }),
        type: 'builtin',
        result: {
          id: 'tool_result_2',
          content: JSON.stringify({
            files: [
              { name: 'document.pdf', size: 1024000, type: 'file' },
              { name: 'images', size: 0, type: 'directory' },
              { name: 'notes.txt', size: 2048, type: 'file' },
            ],
          }),
          state: null,
        },
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

    // Add tool messages (role='tool') for each tool call
    state.messagesMap[mapKey].push({
      id: `tool-msg-${Date.now()}-1`,
      role: 'tool',
      content: JSON.stringify({
        results: [
          {
            title: 'Mock Search Result 1',
            url: 'https://example.com/result1',
            description: 'This is a mock search result for testing purposes.',
          },
          {
            title: 'Mock Search Result 2',
            url: 'https://example.com/result2',
            description: 'Another mock search result with relevant information.',
          },
        ],
      }),
      tool_call_id: 'tool_1',
      sessionId: activeId,
      topicId: activeTopicId,
      threadId,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      meta: {},
    } as UIChatMessage);

    state.messagesMap[mapKey].push({
      id: `tool-msg-${Date.now()}-2`,
      role: 'tool',
      content: JSON.stringify({
        files: [
          { name: 'document.pdf', size: 1024000, type: 'file' },
          { name: 'images', size: 0, type: 'directory' },
          { name: 'notes.txt', size: 2048, type: 'file' },
        ],
      }),
      tool_call_id: 'tool_2',
      sessionId: activeId,
      topicId: activeTopicId,
      threadId,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      meta: {},
    } as UIChatMessage);
  }), false, n('frontendMock/response'));

  console.log('[Frontend Mock] Response complete with full mock data');
}

// ================================================================
// MAIN STORE ACTIONS
// ================================================================

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

    console.log(`[SendMessage] Mode: ${API_MODE}`, {
      activeId,
      activeTopicId,
      threadId,
    });

    // ================================================================
    // STEP 2: CREATE OPTIMISTIC UI
    // ================================================================
    // set({ isCreatingMessage: true }, false, n('creatingMessage/start'));

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
      switch (API_MODE) {
        case 'REAL':
          await handleRealAPI(get, set, {
            activeId,
            activeTopicId,
            threadId,
            message,
            messageUserId,
            messageAssistantId,
          });
          break;

        case 'BACKEND_MOCK':
          await handleBackendMock(get, set, {
            activeId,
            activeTopicId,
            threadId,
            message,
            messageUserId,
            messageAssistantId,
          });
          break;

        case 'FRONTEND_MOCK':
          await handleFrontendMock(set, {
            activeId,
            activeTopicId,
            threadId,
            message,
            mapKey,
            messageAssistantId,
          });
          break;

        default:
          throw new Error(`Unknown API_MODE: ${API_MODE}`);
      }

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
        // Find the loading assistant message by ID (we already have the real ID)
        const msgIndex = messages.findIndex(m => m.id === data.message_id);

        if (msgIndex !== -1) {
          console.log('[Stream] Found assistant message with ID:', data.message_id);

          // Add to loading IDs for animation
          if (!state.chatLoadingIds.includes(data.message_id)) {
            state.chatLoadingIds.push(data.message_id);
          }
        } else {
          console.warn('[Stream] Could not find assistant message for start event:', data.message_id);
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

