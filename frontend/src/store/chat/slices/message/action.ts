/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
// Disable the auto sort key eslint rule to make the code more logic and readable
import {
  ChatErrorType,
  ChatImageItem,
  ChatMessageError,
  ChatMessagePluginError,
  CreateMessageParams,
  GroundingSearch,
  MessageMetadata,
  MessageToolCall,
  ModelReasoning,
  TraceEventPayloads,
  TraceEventType,
  UIChatMessage,
  UpdateMessageRAGParams,
} from '@/types';
import { nanoid } from '@/utils';
import { copyToClipboard } from '@lobehub/ui';
import isEqual from 'fast-deep-equal';
import { StateCreator } from 'zustand/vanilla';

import { messageService } from '@/services/message';
import { ChatStore } from '@/store/chat/store';
import { messageMapKey } from '@/store/chat/utils/messageMapKey';
import { useSessionStore } from '@/store/session';
import { sessionSelectors } from '@/store/session/selectors';
import { Action, setNamespace } from '@/utils/storeDebug';

import { DB, toNullString, toNullJSON, currentTimestampMs } from '@/types/database';
import { getUserId } from '@/store/session/helpers';
import { MessageModel } from '@/database/models/message';

import type { ChatStoreState } from '../../initialState';
import { chatSelectors } from '../../selectors';
import { preventLeavingFn, toggleBooleanList } from '../../utils';
import { MessageDispatch, messagesReducer } from './reducer';

const n = setNamespace('m');

// Helper function to parse JSON string or return the value as-is
const parseJSONField = (value: any): any => {
  if (!value) return undefined;
  if (typeof value === 'string') {
    try {
      return JSON.parse(value);
    } catch {
      return value;
    }
  }
  return value;
};

// Helper function to ensure tool arguments is always a string
const normalizeTools = (tools: any): any[] | undefined => {
  if (!tools || !Array.isArray(tools)) return undefined;
  return tools.map((tool: any) => ({
    ...tool,
    // Ensure arguments is always a JSON string (UI expects string, will parse with safeParseJSON)
    arguments: typeof tool.arguments === 'object' 
      ? JSON.stringify(tool.arguments) 
      : tool.arguments,
  }));
};

// Helper function to map messages from DB to UI format
// Helper to get value - handles both NullString format and direct values
const getValue = (val: any): any => {
  if (val === null || val === undefined) return undefined;
  // If it's a NullString format { String: "...", Valid: bool }
  if (typeof val === 'object' && 'String' in val && 'Valid' in val) {
    return val.Valid ? val.String : undefined;
  }
  return val;
};

// Helper to parse value - tries to parse JSON string, otherwise returns as-is
const parseValue = (val: any): any => {
  const v = getValue(val);
  if (v === undefined || v === null) return undefined;
  if (typeof v === 'string') {
    try {
      return JSON.parse(v);
    } catch {
      return v;
    }
  }
  return v;
};

const mapMessagesFromDB = (dbMessages: any[]): UIChatMessage[] => {
  return dbMessages.map((msg: any) => {
    // Parse reasoning - can be string JSON or object
    const reasoning = parseValue(msg.reasoning);
    
    // Parse tools - can be string JSON or object array, then normalize arguments
    const rawTools = parseValue(msg.tools);
    const tools = normalizeTools(rawTools);
    
    // Parse search - can be string JSON or object
    const search = parseValue(msg.search);
    
    // Parse metadata - contains chunksList, imageList, usage, performance
    const metadata = parseValue(msg.metadata);
    
    // Parse error - can be string JSON or object
    const error = parseValue(msg.error);
    
    // Helper to get string value
    const getStr = (val: any): string | undefined => {
      const v = getValue(val);
      return typeof v === 'string' ? v : undefined;
    };

    return {
      id: msg.id,
      content: getStr(msg.content) || '',
      role: msg.role,
      createdAt: typeof msg.createdAt === 'number' ? msg.createdAt : new Date(msg.createdAt).getTime(),
      updatedAt: typeof msg.updatedAt === 'number' ? msg.updatedAt : new Date(msg.updatedAt).getTime(),
      meta: {},
      sessionId: getStr(msg.sessionId),
      topicId: getStr(msg.topicId),
      threadId: getStr(msg.threadId),
      parentId: getStr(msg.parentId),
      groupId: getStr(msg.groupId),
      agentId: getStr(msg.agentId),
      targetId: getStr(msg.targetId),
      // Reasoning data for Reasoning component
      reasoning,
      // Tool calls for Tool component
      tools,
      // Search grounding for SearchGrounding component
      search,
      // Error state
      error,
      // chunksList: prefer from msg.chunksList (from DB query), fallback to metadata.chunksList (from mock)
      chunksList: msg.chunksList?.length > 0 ? msg.chunksList : metadata?.chunksList,
      // imageList: prefer from msg.imageList (from DB query), fallback to metadata.imageList (from mock)
      imageList: msg.imageList?.length > 0 ? msg.imageList : metadata?.imageList,
      // Extra info (model, provider, translate, tts)
      extra: {
        fromModel: getStr(msg.model),
        fromProvider: getStr(msg.provider),
      },
      // Metadata for Usage component
      metadata: {
        usage: metadata?.usage,
        performance: metadata?.performance,
        model: metadata?.model,
        temperature: metadata?.temperature,
      },
    } as UIChatMessage;
  });
};

export interface ChatMessageAction {
  // create
  addAIMessage: () => Promise<void>;
  addUserMessage: (params: { message: string; fileList?: string[] }) => Promise<void>;
  // delete
  /**
   * clear message on the active session
   */
  clearMessage: () => Promise<void>;
  deleteMessage: (id: string) => Promise<void>;
  deleteToolMessage: (id: string) => Promise<void>;
  clearAllMessages: () => Promise<void>;
  // update
  updateInputMessage: (message: string) => void;
  modifyMessageContent: (id: string, content: string) => Promise<void>;
  toggleMessageEditing: (id: string, editing: boolean) => void;
  // query
  internal_fetchMessages: (
    messageContextId: string,
    activeTopicId?: string,
    type?: 'session' | 'group',
  ) => Promise<void>;
  copyMessage: (id: string, content: string) => Promise<void>;
  refreshMessages: () => Promise<void>;
  replaceMessages: (messages: UIChatMessage[]) => void;
  // =========  ↓ Internal Method ↓  ========== //
  // ========================================== //
  // ========================================== //
  internal_updateMessageRAG: (id: string, input: UpdateMessageRAGParams) => Promise<void>;

  /**
   * update message at the frontend
   * this method will not update messages to database
   */
  internal_dispatchMessage: (
    payload: MessageDispatch,
    context?: { topicId?: string | null; sessionId: string },
  ) => void;

  /**
   * update the message content with optimistic update
   * a method used by other action
   */
  internal_updateMessageContent: (
    id: string,
    content: string,
    extra?: {
      toolCalls?: MessageToolCall[];
      reasoning?: ModelReasoning;
      search?: GroundingSearch;
      metadata?: MessageMetadata;
      imageList?: ChatImageItem[];
      model?: string;
      provider?: string;
    },
  ) => Promise<void>;
  /**
   * update the message error with optimistic update
   */
  internal_updateMessageError: (id: string, error: ChatMessageError | null) => Promise<void>;
  internal_updateMessagePluginError: (
    id: string,
    error: ChatMessagePluginError | null,
  ) => Promise<void>;
  /**
   * create a message with optimistic update
   */
  internal_createMessage: (
    params: CreateMessageParams,
    context?: { tempMessageId?: string; skipRefresh?: boolean },
  ) => Promise<string | undefined>;
  /**
   * create a temp message for optimistic update
   * otherwise the message will be too slow to show
   */
  internal_createTmpMessage: (params: CreateMessageParams) => string;
  /**
   * delete the message content with optimistic update
   */
  internal_deleteMessage: (id: string) => Promise<void>;
  internal_traceMessage: (id: string, payload: TraceEventPayloads) => Promise<void>;

  /**
   * method to toggle message create loading state
   * the AI message status is creating -> generating
   * other message role like user and tool , only this method need to be called
   */
  internal_toggleMessageLoading: (loading: boolean, id: string) => void;
  internal_toggleMessageInToolsCalling: (loading: boolean, id: string) => void;
  internal_toggleChatLoading: (loading: boolean, id?: string, action?: Action) => void;

  /**
   * helper to toggle the loading state of the array,used by these three toggleXXXLoading
   */
  internal_toggleLoadingArrays: (
    key: keyof ChatStoreState,
    loading: boolean,
    id?: string,
    action?: Action,
  ) => AbortController | undefined;

  /**
   * Update active session type
   */
  internal_updateActiveSessionType: (sessionType?: 'agent' | 'group') => void;
  /**
   * Update active session ID with cleanup of pending operations
   */
  internal_updateActiveId: (activeId: string) => void;
}

export const chatMessage: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatMessageAction
> = (set, get) => ({
  deleteMessage: async (id) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    let ids = [message.id];

    // if the message is a tool calls, then delete all the related messages
    if (message.tools) {
      const toolMessageIds = message.tools.flatMap((tool) => {
        const messages = chatSelectors
          .activeBaseChats(get())
          .filter((m) => m.tool_call_id === tool.id);

        return messages.map((m) => m.id);
      });
      ids = ids.concat(toolMessageIds);
    }

    get().internal_dispatchMessage({ type: 'deleteMessages', ids });

    await DB.BatchDeleteMessages(ids);

    console.log('[Message] Deleted messages via direct DB', { ids, count: ids.length });

    await get().refreshMessages();
  },

  deleteToolMessage: async (id) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message || message.role !== 'tool') return;

    const removeToolInAssistantMessage = async () => {
      if (!message.parentId) return;
      await get().internal_removeToolToAssistantMessage(message.parentId, message.tool_call_id);
    };

    await Promise.all([
      // 1. remove tool message
      get().internal_deleteMessage(id),
      // 2. remove the tool item in the assistant tools
      removeToolInAssistantMessage(),
    ]);
  },

  clearMessage: async () => {
    const {
      activeId,
      activeTopicId,
      refreshMessages,
      switchTopic,
      activeSessionType,
    } = get();

    // Check if this is a group session - use activeSessionType if available, otherwise check session store
    let isGroupSession = activeSessionType === 'group';
    if (activeSessionType === undefined) {
      // Fallback: check session store directly
      const sessionStore = useSessionStore.getState();
      isGroupSession = sessionSelectors.isCurrentSessionGroupSession(sessionStore);
    }

    if (activeTopicId) {
      // If topic exists, delete messages by topic
      await DB.DeleteMessagesByTopic(toNullString(activeTopicId));
      console.log('[Message] Cleared topic messages via direct DB', { topicId: activeTopicId });
    } else if (isGroupSession) {
      // For group chat without topic, delete by group
      await DB.DeleteMessagesByGroup(toNullString(activeId));
      console.log('[Message] Cleared group messages via direct DB', { groupId: activeId });
    } else {
      // For regular session without topic, delete by session
      await DB.DeleteMessagesBySession(toNullString(activeId));
      console.log('[Message] Cleared session messages via direct DB', { sessionId: activeId });
    }

    if (activeTopicId) {
      await get().removeTopic(activeTopicId);
    }
    // refreshTopic is already called inside removeTopic
    await refreshMessages();

    // after remove topic , go back to default topic
    switchTopic();
  },
  clearAllMessages: async () => {
    const { refreshMessages } = get();

    const userId = getUserId();
    await DB.DeleteAllMessages(userId);

    console.log('[Message] Cleared all messages via direct DB', { userId });

    await refreshMessages();
  },
  addAIMessage: async () => {
    const { internal_createMessage, updateInputMessage, activeTopicId, activeId, inputMessage } =
      get();
    if (!activeId) return;

    await internal_createMessage({
      content: inputMessage,
      role: 'assistant',
      sessionId: activeId,
      // if there is activeTopicId，then add topicId to message
      topicId: activeTopicId,
    });

    updateInputMessage('');
  },
  addUserMessage: async ({ message, fileList }) => {
    const { internal_createMessage, updateInputMessage, activeTopicId, activeId, activeThreadId } =
      get();
    if (!activeId) return;

    await internal_createMessage({
      content: message,
      files: fileList,
      role: 'user',
      sessionId: activeId,
      // if there is activeTopicId，then add topicId to message
      topicId: activeTopicId,
      threadId: activeThreadId,
    });

    updateInputMessage('');
  },
  copyMessage: async (id, content) => {
    await copyToClipboard(content);

    get().internal_traceMessage(id, { eventType: TraceEventType.CopyMessage });
  },
  toggleMessageEditing: (id, editing) => {
    set(
      { messageEditingIds: toggleBooleanList(get().messageEditingIds, id, editing) },
      false,
      'toggleMessageEditing',
    );
  },

  updateInputMessage: (message) => {
    if (isEqual(message, get().inputMessage)) return;

    set({ inputMessage: message }, false, n('updateInputMessage', message));
  },
  modifyMessageContent: async (id, content) => {
    // tracing the diff of update
    // due to message content will change, so we need send trace before update,or will get wrong data
    get().internal_traceMessage(id, {
      eventType: TraceEventType.ModifyMessage,
      nextContent: content,
    });

    await get().internal_updateMessageContent(id, content);
  },

  /**
   * @param enable - whether to enable the fetch
   * @param messageContextId - Can be sessionId or groupId
   */
  internal_fetchMessages: async (messageContextId, activeTopicId, type = 'session') => {
    if (!messageContextId) return;

    const mapKey = messageMapKey(messageContextId, activeTopicId);
    console.log('[internal_fetchMessages] Called with:', {
      messageContextId,
      activeTopicId,
      type,
      mapKey,
    });

    try {
      const userId = getUserId();
      let messages: UIChatMessage[] = [];

      // CRITICAL FIX: Check if there are existing messages in this key
      // Don't overwrite with empty array from DB if messages already exist
      const existingMessages = get().messagesMap[mapKey];
      
      if (existingMessages && existingMessages.length > 0) {
        console.log('[internal_fetchMessages] Skipping - messages already exist:', {
          mapKey,
          count: existingMessages.length,
        });
        return;
      }

      // Default topic (activeTopicId is null/undefined) should always show empty state (Welcome screen)
      // Only fetch messages when a specific topic is selected
      if (activeTopicId) {
        // Use MessageModel.query() to get full message data including plugins, files, chunks
        const messageModel = new MessageModel(DB, userId);
        messages = await messageModel.query({
          topicId: activeTopicId,
          pageSize: 1000,
        });
      }
      // For default topic (no activeTopicId), keep messages as empty array []
      // This ensures Welcome screen is shown in Content.tsx when data.length === 0

      const nextMap = {
        ...get().messagesMap,
        [mapKey]: messages,
      };

      // no need to update map if the messages have been init and the map is the same
      if (get().messagesInit && isEqual(nextMap, get().messagesMap)) {
        console.log('[internal_fetchMessages] Skipping - map unchanged');
        return;
      }

      console.log('[Message] Fetched messages via MessageModel', {
        type,
        messageContextId,
        activeTopicId,
        mapKey,
        count: messages.length,
        isDefaultTopic: !activeTopicId,
        existingKeys: Object.keys(get().messagesMap),
      });

      set(
        { messagesInit: true, messagesMap: nextMap },
        false,
        n('internal_fetchMessages', { messages, messageContextId, activeTopicId, type }),
      );
    } catch (error) {
      console.error('[internal_fetchMessages] Error fetching messages:', error);
    }
  },

  /**
   * Refresh messages from database - direct fetch without SWR cache invalidation
   */
  refreshMessages: async () => {
    const { activeId, activeTopicId } = get();
    if (!activeId) return;

    try {
      const userId = getUserId();
      let messages: UIChatMessage[] = [];

      // Default topic (activeTopicId is null/undefined) should always show empty state (Welcome screen)
      // Only fetch messages when a specific topic is selected
      if (activeTopicId) {
        // Use MessageModel.query() to get full message data including plugins, files, chunks
        const messageModel = new MessageModel(DB, userId);
        messages = await messageModel.query({
          topicId: activeTopicId,
          pageSize: 1000,
        });
      }
      // For default topic (no activeTopicId), keep messages as empty array []
      // This ensures Welcome screen is shown in Content.tsx when data.length === 0

      const nextMap = {
        ...get().messagesMap,
        [messageMapKey(activeId, activeTopicId)]: messages,
      };

      console.log('[Message] Refreshed messages via MessageModel', {
        activeId,
        activeTopicId,
        count: messages.length,
        isDefaultTopic: !activeTopicId
      });

      set(
        { messagesInit: true, messagesMap: nextMap },
        false,
        n('refreshMessages', { activeId, activeTopicId }),
      );
    } catch (error) {
      console.error('[refreshMessages] Error refreshing messages:', error);
    }
  },
  replaceMessages: (messages) => {
    set(
      {
        messagesMap: {
          ...get().messagesMap,
          [messageMapKey(get().activeId, get().activeTopicId)]: messages,
        },
      },
      false,
      'replaceMessages',
    );
  },

  internal_updateMessageRAG: async (id, data) => {
    const { refreshMessages } = get();

    await messageService.updateMessageRAG(id, data);
    await refreshMessages();
  },

  // the internal process method of the AI message
  internal_dispatchMessage: (payload, context) => {
    const activeId = typeof context !== 'undefined' ? context.sessionId : get().activeId;
    const topicId = typeof context !== 'undefined' ? context.topicId : get().activeTopicId;

    const messagesKey = messageMapKey(activeId, topicId);

    const messages = messagesReducer(chatSelectors.getBaseChatsByKey(messagesKey)(get()), payload);

    const nextMap = { ...get().messagesMap, [messagesKey]: messages };

    if (isEqual(nextMap, get().messagesMap)) return;

    set({ messagesMap: nextMap }, false, { type: `dispatchMessage/${payload.type}`, payload });
  },

  internal_updateMessageError: async (id, error) => {
    get().internal_dispatchMessage({ id, type: 'updateMessage', value: { error } });

    // Note: UpdateMessage doesn't support error field directly
    // Error is stored in the message metadata or handled separately
    await DB.UpdateMessage({
      id,
      metadata: toNullString(JSON.stringify({ error: error ? JSON.stringify(error) : undefined })),
      updatedAt: currentTimestampMs(),
    } as any);

    console.log('[Message] Updated message error via direct DB', { id });

    await get().refreshMessages();
  },

  internal_updateMessagePluginError: async (id, error) => {
    // Get current plugin item first
    const item = await DB.GetMessagePlugin(id);

    if (!item) {
      console.error('[Message] Plugin not found for error update', { id });
      return;
    }

    // Update plugin error
    await DB.UpdateMessagePlugin({
      id,
      state: item.state, // Keep existing state
      error: error !== undefined ? toNullJSON(error) : item.error,
    });

    console.log('[Message] Updated message plugin error via direct DB', { id, hasError: !!error });

    await get().refreshMessages();
  },

  internal_updateMessageContent: async (id, content, extra) => {
    const { internal_dispatchMessage, refreshMessages, internal_transformToolCalls } = get();

    // Due to the async update method and refresh need about 100ms
    // we need to update the message content at the frontend to avoid the update flick
    // refs: https://medium.com/@kyledeguzmanx/what-are-optimistic-updates-483662c3e171
    if (extra?.toolCalls) {
      internal_dispatchMessage({
        id,
        type: 'updateMessage',
        value: { tools: internal_transformToolCalls(extra?.toolCalls) },
      });
    } else {
      internal_dispatchMessage({
        id,
        type: 'updateMessage',
        value: { content },
      });
    }

    const updateData: any = {
      id,
      updatedAt: currentTimestampMs(),
    };

    // Add fields that are provided
    if (content !== undefined) updateData.content = toNullString(content);
    if (extra?.toolCalls) updateData.tools = toNullString(JSON.stringify(internal_transformToolCalls(extra.toolCalls)));
    if (extra?.reasoning) updateData.reasoning = toNullString(JSON.stringify(extra.reasoning));
    if (extra?.search) updateData.search = toNullString(JSON.stringify(extra.search));
    if (extra?.metadata) updateData.metadata = toNullString(JSON.stringify(extra.metadata));
    if (extra?.model) updateData.model = toNullString(extra.model);
    if (extra?.provider) updateData.provider = toNullString(extra.provider);
    if (extra?.imageList) updateData.imageList = toNullString(JSON.stringify(extra.imageList));

    try {
      await DB.UpdateMessage(updateData);
      console.log('[Message] Updated message content via direct DB', { id, hasContent: !!content, hasExtra: !!extra });
    } catch (error: any) {
      // If the message doesn't exist in DB (e.g., it's a temp message), skip the update
      // The message will be saved when refreshMessages is called
      if (error?.message?.includes('no rows in result set')) {
        console.log('[Message] Skipping update for temp message (not in DB yet)', { id });
      } else {
        console.error('[Message] Failed to update message:', error);
        throw error;
      }
    }

    await refreshMessages();
  },

  internal_createMessage: async (message, context) => {
    const {
      internal_createTmpMessage,
      refreshMessages,
      internal_toggleMessageLoading,
      internal_dispatchMessage,
    } = get();
    let tempId = context?.tempMessageId;
    if (!tempId) {
      // use optimistic update to avoid the slow waiting
      tempId = internal_createTmpMessage(message);

      internal_toggleMessageLoading(true, tempId);
    }

    try {
      const userId = getUserId();
      const messageModel = new MessageModel(DB, userId);
      const dbMessage = await messageModel.create(message);
      const id = dbMessage.id;

      console.log('[Message] Created message via MessageModel', { id, role: message.role });

      if (!context?.skipRefresh) {
        internal_toggleMessageLoading(true, tempId);
        await refreshMessages();
      }

      internal_toggleMessageLoading(false, tempId);
      return id;
    } catch (e) {
      internal_toggleMessageLoading(false, tempId);
      internal_dispatchMessage({
        id: tempId,
        type: 'updateMessage',
        value: {
          error: { type: ChatErrorType.CreateMessageError, message: (e as Error).message, body: e },
        },
      });
    }
  },

  internal_createTmpMessage: (message) => {
    const { internal_dispatchMessage } = get();

    // use optimistic update to avoid the slow waiting
    const tempId = 'tmp_' + nanoid();
    internal_dispatchMessage({ type: 'createMessage', id: tempId, value: message });

    return tempId;
  },
  internal_deleteMessage: async (id: string) => {
    get().internal_dispatchMessage({ type: 'deleteMessage', id });

    await DB.DeleteMessage(id);

    console.log('[Message] Deleted message via direct DB', { id });

    await get().refreshMessages();
  },
  internal_traceMessage: async (id, payload) => {
    // tracing the diff of update
  },

  // ----- Loading ------- //
  internal_toggleMessageLoading: (loading, id) => {
    set(
      {
        messageLoadingIds: toggleBooleanList(get().messageLoadingIds, id, loading),
      },
      false,
      `internal_toggleMessageLoading/${loading ? 'start' : 'end'}`,
    );
  },
  internal_toggleChatLoading: (loading, id, action) => {
    get().internal_toggleLoadingArrays('messageLoadingIds', loading, id, action);
  },
  internal_toggleMessageInToolsCalling: (loading, id) => {
    get().internal_toggleLoadingArrays('toolsCallingMessageIds', loading, id);
  },
  internal_toggleLoadingArrays: (key, loading, id, action) => {
    const abortControllerKey = `${key}AbortController`;
    if (loading) {
      window.addEventListener('beforeunload', preventLeavingFn);

      const abortController = new AbortController();
      set(
        {
          [abortControllerKey]: abortController,
          [key]: toggleBooleanList(get()[key] as string[], id!, loading),
        },
        false,
        action,
      );

      return abortController;
    } else {
      if (!id) {
        set({ [abortControllerKey]: undefined, [key]: [] }, false, action);
      } else
        set(
          {
            [abortControllerKey]: undefined,
            [key]: toggleBooleanList(get()[key] as string[], id, loading),
          },
          false,
          action,
        );

      window.removeEventListener('beforeunload', preventLeavingFn);
    }
  },
  internal_updateActiveSessionType: (sessionType?: 'agent' | 'group') => {
    if (get().activeSessionType === sessionType) return;

    set({ activeSessionType: sessionType }, false, n('updateActiveSessionType'));
  },

  internal_updateActiveId: (activeId: string) => {
    const currentActiveId = get().activeId;
    if (currentActiveId === activeId) return;

    // Before switching sessions, cancel all pending supervisor decisions
    get().internal_cancelAllSupervisorDecisions();

    set({ activeId }, false, n(`updateActiveId/${activeId}`));
  },
});
