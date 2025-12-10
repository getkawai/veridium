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
import {
  StreamEventPayload,
  UIChatMessage as BackendUIChatMessage,
} from '@@/github.com/kawai-network/veridium/internal/services/models';

// Re-export StreamEventPayload for consumers
export type { StreamEventPayload };

// User ID constant for backend calls
const FALLBACK_CLIENT_DB_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

// ================================================================
// API MODE CONFIGURATION
// ================================================================
/**
 * API_MODE controls which backend to use:
 * 
 * - 'REAL': Production mode - calls real backend API with LLM (non-streaming)
 * - 'BACKEND_MOCK': Development mode - calls backend mock API (saves to DB, no LLM)
 * - 'BACKEND_MOCK_STREAM': Streaming mock - emits events with delays for realistic UI testing
 * - 'BACKEND_REAL_STREAM': Production streaming - real LLM with event streaming
 * - 'FRONTEND_MOCK': UI testing mode - frontend-only mock (no backend, no DB)
 */
type ApiMode = 'REAL' | 'BACKEND_MOCK' | 'BACKEND_MOCK_STREAM' | 'BACKEND_REAL_STREAM' | 'FRONTEND_MOCK';

const API_MODE: ApiMode = 'BACKEND_REAL_STREAM' as ApiMode;

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
//   - Returns all messages at once
//   - Good for backend integration testing
//
// BACKEND_MOCK_STREAM:
//   - No LLM calls (saves costs)
//   - Saves to database (realistic data flow)
//   - Mock tool results
//   - Emits events progressively with delays (simulates real streaming)
//   - Events: start, reasoning, chunk, tool_call, tool_result, complete
//   - Good for testing streaming UI
//
// BACKEND_REAL_STREAM:
//   - Uses REAL local LLM (Llama, Qwen, etc.)
//   - Real tool execution
//   - Saves to database
//   - Emits events progressively (real streaming from LLM)
//   - Events: start, reasoning, chunk, tool_call, tool_result, complete
//   - Production-ready with streaming UI
//
// FRONTEND_MOCK:
//   - No backend calls
//   - No database
//   - Instant responses
//   - Good for UI/UX development
// ================================================================

const n = setNamespace('ai');

// ================================================================
// HELPER FUNCTIONS
// ================================================================

/**
 * Map backend UIChatMessage to frontend UIChatMessage
 * Backend uses camelCase consistently now after the refactor
 */
function mapBackendToFrontendMessage(
  backendMsg: BackendUIChatMessage,
  existingMsg: UIChatMessage,
  defaults: { activeId: string; activeTopicId?: string; threadId?: string }
): void {
  const { activeId, activeTopicId, threadId } = defaults;

  // Map core fields
  existingMsg.id = backendMsg.id || existingMsg.id;
  existingMsg.content = backendMsg.content || existingMsg.content;
  existingMsg.sessionId = backendMsg.sessionId || activeId;
  existingMsg.topicId = backendMsg.topicId || activeTopicId;
  existingMsg.threadId = backendMsg.threadId || threadId;
  existingMsg.createdAt = backendMsg.createdAt || Date.now();
  existingMsg.updatedAt = backendMsg.updatedAt || Date.now();
  (existingMsg as any).loading = false;

  // Map optional fields directly (backend now uses same structure as frontend)
  if (backendMsg.tools?.length) (existingMsg as any).tools = backendMsg.tools;
  if (backendMsg.children?.length) (existingMsg as any).children = backendMsg.children;
  if (backendMsg.chunksList?.length) (existingMsg as any).chunksList = backendMsg.chunksList;
  if (backendMsg.imageList?.length) (existingMsg as any).imageList = backendMsg.imageList;
  if (backendMsg.fileList?.length) (existingMsg as any).fileList = backendMsg.fileList;
  if (backendMsg.videoList?.length) (existingMsg as any).videoList = backendMsg.videoList;
  if (backendMsg.reasoning) (existingMsg as any).reasoning = backendMsg.reasoning;
  if (backendMsg.search) (existingMsg as any).search = backendMsg.search;
  if (backendMsg.usage) (existingMsg as any).usage = backendMsg.usage;
  if (backendMsg.performance) (existingMsg as any).performance = backendMsg.performance;
  if (backendMsg.metadata) (existingMsg as any).metadata = backendMsg.metadata;
  if (backendMsg.extra) (existingMsg as any).extra = backendMsg.extra;
  if (backendMsg.plugin) (existingMsg as any).plugin = backendMsg.plugin;
  if (backendMsg.meta) existingMsg.meta = backendMsg.meta as any;
}

// ================================================================
// API MODE HANDLERS
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
    mapKey: string;
    fileIds?: string[];
  }
) {
  const { activeId, activeTopicId, threadId, message, messageUserId, messageAssistantId, mapKey, fileIds } = context;

  console.log('[Real API] Calling real backend with LLM...');

  const response = await backendAgentChat.sendMessage({
    session_id: activeId,
    user_id: FALLBACK_CLIENT_DB_USER_ID,
    message: message,
    topic_id: activeTopicId || undefined,
    thread_id: threadId || undefined,
    message_user_id: messageUserId,
    message_assistant_id: messageAssistantId,
    file_ids: fileIds && fileIds.length > 0 ? fileIds : undefined,
  });

  console.log('[Real API] Response received:', response);

  const finalTopicId = response.topicId || activeTopicId;

  // Update the assistant message with backend response
  set(produce((state: ChatStore) => {
    const messages = state.messagesMap[mapKey];
    if (!messages) return;

    const assistantMsgIndex = messages.findIndex((m: UIChatMessage) => m.id === messageAssistantId);
    if (assistantMsgIndex === -1) {
      console.error('[Real API] Assistant message not found');
      return;
    }

    mapBackendToFrontendMessage(response, messages[assistantMsgIndex], { activeId, activeTopicId, threadId });
    console.log('[Real API] Updated assistant message with backend data');
  }), false, n('realAPI/updateAssistantMessage'));

  // If a new topic was created, switch to it
  if (finalTopicId && finalTopicId !== activeTopicId) {
    console.log('[Real API] New topic created, switching to:', finalTopicId);
    await get().refreshTopic();
    await get().switchTopic(finalTopicId);
  } else {
    console.log('[Real API] Staying on same topic');
  }

  console.log('[Real API] Complete');
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

    // Extract file IDs from uploaded files (already processed by FileProcessorService)
    const fileIds = files?.map(f => f.id).filter(Boolean) || [];

    console.log(`[SendMessage] Mode: ${API_MODE}`, {
      activeId,
      activeTopicId,
      threadId,
      fileIds,
    });

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
      switch (API_MODE) {
        case 'REAL':
          await handleRealAPI(get, set, {
            activeId,
            activeTopicId,
            threadId,
            message,
            messageUserId,
            messageAssistantId,
            mapKey,
            fileIds,
          });
          break;

        case 'BACKEND_MOCK': {
          console.log('[Backend Mock] Calling backend mock (saves to DB)...');

          // Backend returns array: [userMsg, assistantMsg, toolMsg1, toolMsg2, ...]
          const mockMessages = await backendAgentChat.sendMessageMock({
            session_id: activeId,
            user_id: FALLBACK_CLIENT_DB_USER_ID,
            message: message,
            topic_id: activeTopicId || undefined,
            thread_id: threadId || undefined,
            message_user_id: messageUserId,
            message_assistant_id: messageAssistantId,
            file_ids: fileIds.length > 0 ? fileIds : undefined,
          });

          console.log('[Backend Mock] Received', mockMessages.length, 'messages');

          // Find assistant message to get topicId
          const assistantMsg = mockMessages.find(m => m.role === 'assistant');
          const mockFinalTopicId = assistantMsg?.topicId || activeTopicId;

          // Replace optimistic messages with backend messages
          set(produce((state: ChatStore) => {
            const messages = state.messagesMap[mapKey];
            if (!messages) return;

            // Remove optimistic user and assistant messages
            const filteredMessages = messages.filter(
              m => m.id !== messageUserId && m.id !== messageAssistantId
            );

            // Add all messages from backend (user, assistant, tools)
            for (const msg of mockMessages) {
              filteredMessages.push(msg as unknown as UIChatMessage);
            }

            state.messagesMap[mapKey] = filteredMessages;
            console.log('[Backend Mock] Replaced optimistic messages with', mockMessages.length, 'backend messages');
          }), false, n('backendMock/replaceMessages'));

          // If a new topic was created, switch to it
          if (mockFinalTopicId && mockFinalTopicId !== activeTopicId) {
            console.log('[Backend Mock] New topic created, switching to:', mockFinalTopicId);
            await get().refreshTopic();
            await get().switchTopic(mockFinalTopicId);
          }

          console.log('[Backend Mock] Complete');
          break;
        }

        case 'BACKEND_MOCK_STREAM': {
          console.log('[Backend Mock Stream] Starting streaming mock...');

          // Call streaming mock - all updates come via events
          // Events are handled by internal_handleStreamEvent (called from App.tsx)
          await backendAgentChat.sendMessageMockStream({
            session_id: activeId,
            user_id: FALLBACK_CLIENT_DB_USER_ID,
            message: message,
            topic_id: activeTopicId || undefined,
            thread_id: threadId || undefined,
            message_user_id: messageUserId,
            message_assistant_id: messageAssistantId,
            file_ids: fileIds.length > 0 ? fileIds : undefined,
          });

          // Note: User message is also created via streaming events
          // The optimistic user message will be updated with real data

          console.log('[Backend Mock Stream] Streaming complete, data came via events');
          break;
        }

        case 'BACKEND_REAL_STREAM': {
          console.log('[Backend Real Stream] Starting real LLM streaming...');

          // Call real streaming - uses real LLM with streaming events
          // Events are handled by internal_handleStreamEvent (called from App.tsx)
          await backendAgentChat.sendMessageRealStream({
            session_id: activeId,
            user_id: FALLBACK_CLIENT_DB_USER_ID,
            message: message,
            topic_id: activeTopicId || undefined,
            thread_id: threadId || undefined,
            message_user_id: messageUserId,
            message_assistant_id: messageAssistantId,
            file_ids: fileIds.length > 0 ? fileIds : undefined,
          });

          // Note: User message is also created via streaming events
          // Real LLM response comes token by token via events

          console.log('[Backend Real Stream] Streaming complete, data came via events');
          break;
        }

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
   * 
   * Event types:
   * - start: Generation started, add to loading IDs, handle new topic creation
   * - reasoning: Thinking content (update reasoning field)
   * - chunk: Content chunks (update content field)
   * - tool_call: Tool call initiated (update tools array)
   * - tool_result: Tool execution result (add tool message to messagesMap)
   * - complete: Generation finished (finalize message, remove from loading)
   */
  internal_handleStreamEvent: (data: StreamEventPayload) => {
    const { activeId, activeTopicId } = get();
    const currentMapKey = messageMapKey(activeId, activeTopicId);

    // Check if backend created a new topic (first message scenario)
    const newTopicId = data.topic_id;
    const isNewTopic = newTopicId && newTopicId !== activeTopicId;

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

      if (!messages) return;

      if (data.type === 'start') {
        // Find the loading assistant message by ID
        const msgIndex = messages.findIndex(m => m.id === data.message_id);

        if (msgIndex !== -1) {
          console.log('[Stream] Start - found assistant message:', data.message_id);

          // Add to loading IDs for animation
          if (!state.chatLoadingIds.includes(data.message_id)) {
            state.chatLoadingIds.push(data.message_id);
          }

          // Handle new topic creation - move messages to new mapKey
          if (isNewTopic) {
            console.log('[Stream] New topic created, moving messages:', { oldTopicId: activeTopicId, newTopicId });

            // Update topicId on all messages in current conversation
            messages.forEach(m => {
              m.topicId = newTopicId;
            });

            // Move messages to new mapKey
            const newMapKey = messageMapKey(activeId, newTopicId);
            state.messagesMap[newMapKey] = messages;

            // Clear old mapKey
            delete state.messagesMap[mapKey];

            // Update activeTopicId
            state.activeTopicId = newTopicId;

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
        console.log('[Stream] Tool call:', data.tool?.apiName || 'unknown');

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
        console.log('[Stream] Tool result:', data.tool_call_id);

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
        state.chatLoadingIds = state.chatLoadingIds.filter(id => id !== data.message_id);
        console.log('[Stream] Complete - message finalized:', data.message_id);
      }
    }), false, n('streamEvent'));

    // After complete event, schedule a topic refresh to get LLM-generated title
    // This is a fallback in case chat:topic:updated event is not received
    if (data.type === 'complete' && isNewTopic) {
      console.log('[Stream] Scheduling topic refresh for new topic title...');
      setTimeout(() => {
        get().refreshTopic();
      }, 5000); // Wait 5 seconds for backend to generate title
    }
  },
});

