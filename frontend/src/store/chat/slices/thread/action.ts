/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
// Disable the auto sort key eslint rule to make the code more logic and readable
import { LOADING_FLAT, isDeprecatedEdition } from '@/const';
import { chainSummaryTitle } from '@/prompts';
import {
  CreateMessageParams,
  SendThreadMessageParams,
  ThreadItem,
  ThreadStatus,
  ThreadType,
  UIChatMessage,
} from '@/types';
import isEqual from 'fast-deep-equal';
import { StateCreator } from 'zustand/vanilla';
import { nanoid } from 'nanoid';

import { clientDB } from '@/database/client/db';
import { chatService } from '@/services/chat';
import { threadSelectors } from './selectors';
import { ChatStore } from '@/store/chat/store';
import { globalHelpers } from '@/store/global/helpers';
import { useUserStore } from '@/store/user';
import { systemAgentSelectors } from '@/store/user/selectors';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';
import { DB, toNullString, getNullableString, currentTimestampMs, Thread } from '@/types/database';
import { getUserId } from '@/store/session/helpers';
import { MessageModel } from '@/database/models/message';

import { ThreadDispatch, threadReducer } from './reducer';

const n = setNamespace('thd');

const mapThread = (thread: Thread): ThreadItem => {
  const statusStr = getNullableString(thread.status as any);

  return {
    id: thread.id,
    title: thread.title,
    type: thread.type as ThreadType,
    status: (statusStr as ThreadStatus) || ThreadStatus.Active,
    topicId: thread.topicId,
    sourceMessageId: thread.sourceMessageId,
    parentThreadId: getNullableString(thread.parentThreadId as any),
    userId: thread.userId,
    lastActiveAt: new Date(thread.lastActiveAt),
    createdAt: new Date(thread.createdAt),
    updatedAt: new Date(thread.updatedAt),
  };
};

export interface ChatThreadAction {
  // update
  updateThreadInputMessage: (message: string) => void;
  refreshThreads: () => Promise<void>;
  /**
   * Sends a new thread message to the AI chat system
   */
  sendThreadMessage: (params: SendThreadMessageParams) => Promise<void>;
  resendThreadMessage: (messageId: string) => Promise<void>;
  delAndResendThreadMessage: (messageId: string) => Promise<void>;
  createThread: (params: {
    message: CreateMessageParams;
    sourceMessageId: string;
    topicId: string;
    type: ThreadType;
  }) => Promise<{ threadId: string; messageId: string }>;
  openThreadCreator: (messageId: string) => void;
  openThreadInPortal: (threadId: string, sourceMessageId: string) => void;
  closeThreadPortal: () => void;
  internal_fetchThreads: (topicId: string) => Promise<void>;
  summaryThreadTitle: (threadId: string, messages: UIChatMessage[]) => Promise<void>;
  updateThreadTitle: (id: string, title: string) => Promise<void>;
  removeThread: (id: string) => Promise<void>;
  switchThread: (id: string) => void;

  internal_updateThreadTitleInSummary: (id: string, title: string) => void;
  internal_updateThreadLoading: (id: string, loading: boolean) => void;
  internal_updateThread: (id: string, data: Partial<ThreadItem>) => Promise<void>;
  internal_dispatchThread: (payload: ThreadDispatch, action?: any) => void;
}

export const chatThreadMessage: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatThreadAction
> = (set, get) => ({
  updateThreadInputMessage: (message) => {
    if (isEqual(message, get().threadInputMessage)) return;

    set({ threadInputMessage: message }, false, n(`updateThreadInputMessage`, message));
  },

  openThreadCreator: (messageId) => {
    set(
      { threadStartMessageId: messageId, portalThreadId: undefined, startToForkThread: true },
      false,
      'openThreadCreator',
    );
    get().togglePortal(true);
  },
  openThreadInPortal: (threadId, sourceMessageId) => {
    set(
      { portalThreadId: threadId, threadStartMessageId: sourceMessageId, startToForkThread: false },
      false,
      'openThreadInPortal',
    );
    get().togglePortal(true);
  },

  closeThreadPortal: () => {
    set(
      { threadStartMessageId: undefined, portalThreadId: undefined, startToForkThread: undefined },
      false,
      'closeThreadPortal',
    );
    get().togglePortal(false);
  },
  sendThreadMessage: async ({ message }) => {
    // TODO: MIGRATE TO BACKEND
    // internal_coreProcessMessage was removed in big bang migration
    // Thread feature needs to be migrated to use backend agent chat
    console.error('[Thread] sendThreadMessage not yet migrated to backend');
    return;
    
    /* DISABLED - NEEDS MIGRATION
    const {
      internal_coreProcessMessage,
      activeTopicId,
      activeId,
      threadStartMessageId,
      newThreadMode,
      portalThreadId,
    } = get();
    if (!activeId || !activeTopicId) return;

    // if message is empty or no files, then stop
    if (!message) return;

    set({ isCreatingThreadMessage: true }, false, n('creatingThreadMessage/start'));

    const newMessage: CreateMessageParams = {
      content: message,
      // if message has attached with files, then add files to message and the agent
      // files: fileIdList,
      role: 'user',
      sessionId: activeId,
      // if there is activeTopicId，then add topicId to message
      topicId: activeTopicId,
      threadId: portalThreadId,
    };

    let parentMessageId: string | undefined = undefined;
    let tempMessageId: string | undefined = undefined;

    // if there is no portalThreadId, then create a thread and then append message
    let currentThreadId: string | undefined = portalThreadId;
    if (!portalThreadId) {
      if (!threadStartMessageId) return;
      // we need to create a temp message for optimistic update
      tempMessageId = get().internal_createTmpMessage({
        ...newMessage,
        threadId: THREAD_DRAFT_ID,
      });
      get().internal_toggleMessageLoading(true, tempMessageId);

      let threadResult;
      try {
        threadResult = await get().createThread({
          message: newMessage,
          sourceMessageId: threadStartMessageId,
          topicId: activeTopicId,
          type: newThreadMode,
        });
      } catch (error) {
        // Thread creation threw an error (non-conflict error, e.g., database error)
        console.error('[sendThreadMessage] Thread creation threw error:', error);
        
        // Clean up temp message
        get().internal_toggleMessageLoading(false, tempMessageId);
        if (tempMessageId) {
          get().internal_dispatchMessage({ type: 'deleteMessage', id: tempMessageId });
        }
        
        set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
        return;
      }

      // Check if thread creation failed (conflict - returns undefined)
      if (!threadResult || !threadResult.threadId) {
        console.error('[sendThreadMessage] Failed to create thread (conflict or undefined):', threadResult);
        
        // Clean up orphaned message if it was created without a thread
        if (threadResult?.messageId) {
          console.warn('[sendThreadMessage] Cleaning up orphaned message:', threadResult.messageId);
          try {
            await get().internal_deleteMessage(threadResult.messageId);
            // Refresh messages to remove the orphaned message from UI
            await get().refreshMessages();
          } catch (error) {
            console.error('[sendThreadMessage] Failed to delete orphaned message:', error);
          }
        }
        
        // Clean up temp message
        get().internal_toggleMessageLoading(false, tempMessageId);
        if (tempMessageId) {
          get().internal_dispatchMessage({ type: 'deleteMessage', id: tempMessageId });
        }
        
        set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
        return;
      }

      // Ensure messageId exists
      if (!threadResult.messageId) {
        console.error('[sendThreadMessage] Thread created but messageId is missing:', threadResult);
        get().internal_toggleMessageLoading(false, tempMessageId);
        set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
        return;
      }

      const { threadId, messageId } = threadResult;
      parentMessageId = messageId;
      currentThreadId = threadId;

      // mark the portal in thread mode
      await get().refreshThreads();
      await get().refreshMessages();

      get().openThreadInPortal(threadId, threadStartMessageId);
    } else {
      // if there is a thread, just append message
      // we need to create a temp message for optimistic update
      tempMessageId = get().internal_createTmpMessage(newMessage);
      get().internal_toggleMessageLoading(true, tempMessageId);

      parentMessageId = await get().internal_createMessage(newMessage, { tempMessageId });
      
      // CRITICAL: Refresh messages after creating the message so portalAIChats can find it
      await get().refreshMessages();
    }

    get().internal_toggleMessageLoading(false, tempMessageId);

    if (!parentMessageId) return;
    //  update assistant update to make it rerank
    useSessionStore.getState().triggerSessionUpdate(get().activeId);

    // Get the current messages to generate AI response
    // Use currentThreadId to ensure we have the correct threadId even if store hasn't updated yet
    let messages = threadSelectors.portalAIChats(get());
    
    // Double-check: if messages are empty but we have a threadId, try refreshing again
    if (messages.length === 0 && currentThreadId) {
      console.warn('[sendThreadMessage] Messages array is empty, refreshing again...', {
        currentThreadId,
        portalThreadId: get().portalThreadId,
      });
      await get().refreshMessages();
      messages = threadSelectors.portalAIChats(get());
      if (messages.length > 0) {
        console.debug('[sendThreadMessage] Messages found after refresh:', messages.length);
      } else {
        console.error('[sendThreadMessage] Messages still empty after refresh. This may cause the AI request to fail.');
      }
    }

    // Ensure we have a valid threadId
    const finalThreadId = currentThreadId || get().portalThreadId;
    if (!finalThreadId) {
      console.error('[sendThreadMessage] No threadId available:', {
        currentThreadId,
        portalThreadId: get().portalThreadId,
      });
      set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
      return;
    }

    await internal_coreProcessMessage(messages, parentMessageId, {
      ragQuery: get().internal_shouldUseRAG() ? message : undefined,
      threadId: finalThreadId,
      inPortalThread: true,
    });

    set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));

    // 说明是在新建 thread，需要自动总结标题
    if (!portalThreadId) {
      const portalThread = threadSelectors.currentPortalThread(get());

      if (!portalThread) return;

      const chats = threadSelectors.portalAIChats(get());
      await get().summaryThreadTitle(portalThread.id, chats);
    }
    */ // END DISABLED
  },
  resendThreadMessage: async (messageId) => {
    // TODO: MIGRATE TO BACKEND
    // internal_resendMessage was removed in big bang migration
    console.error('[Thread] resendThreadMessage not yet migrated to backend');
    return;
  },
  delAndResendThreadMessage: async (id) => {
    get().resendThreadMessage(id);
    get().deleteMessage(id);
  },
  createThread: async ({ message, sourceMessageId, topicId, type }) => {
    set({ isCreatingThread: true }, false, n('creatingThread/start'));

    const userId = getUserId();
    const now = currentTimestampMs();
    let thread;

    try {
      // Create thread directly with DB call
      const result: Thread = await DB.CreateThread({
        id: nanoid(),
        title: message.content.slice(0, 20),
        type: type,
        status: toNullString(ThreadStatus.Active),
        topicId: topicId,
        sourceMessageId: sourceMessageId,
        parentThreadId: toNullString(''),
        clientId: toNullString(''),
        userId: userId,
        lastActiveAt: now,
        createdAt: now,
        updatedAt: now,
      });

      // With ON CONFLICT DO NOTHING, if conflict occurs, result might be empty/null
      if (!result || !result.id) {
        console.error('[createThread] Thread creation failed (conflict), aborting message creation');
        set({ isCreatingThread: false }, false, n('creatingThread/end'));
        return { messageId: undefined as any, threadId: undefined as any };
      }

      thread = mapThread(result);
      console.log('[Thread] Created thread via direct DB', { threadId: thread.id });
    } catch (error: any) {
      // Check if this is a "no rows in result set" error (ON CONFLICT DO NOTHING with empty RETURNING)
      const isNoRowsError =
        error?.message?.includes('no rows in result set') ||
        error?.message?.includes('sql: no rows in result set') ||
        error?.code === 'PGRST116';

      if (isNoRowsError) {
        console.error('[createThread] Thread creation conflict (ON CONFLICT DO NOTHING - no rows returned), aborting message creation');
        set({ isCreatingThread: false }, false, n('creatingThread/end'));
        return { messageId: undefined as any, threadId: undefined as any };
      }

      // Check if this is a UNIQUE constraint conflict
      const isConflictError =
        error?.message?.includes('UNIQUE constraint') ||
        error?.message?.includes('2067') ||
        error?.code === 'SQLITE_CONSTRAINT_UNIQUE' ||
        error?.code === 2067;

      if (isConflictError) {
        console.error('[createThread] Thread creation conflict (UNIQUE constraint), aborting message creation');
        set({ isCreatingThread: false }, false, n('creatingThread/end'));
        return { messageId: undefined as any, threadId: undefined as any };
      }

      // For other errors, log and re-throw
      console.error('[createThread] Thread creation failed with non-conflict error:', error);
      set({ isCreatingThread: false }, false, n('creatingThread/end'));
      throw error;
    }

    // Create message using MessageModel (not yet migrated to direct DB)
    const messageModel = new MessageModel(clientDB as any, userId);
    const dbMessage = await messageModel.create({
      ...message,
      sessionId: message.sessionId || '',
      threadId: thread.id,
    });

    // If message creation failed, we still return the threadId but no messageId
    if (!dbMessage?.id) {
      console.error('[createThread] Message creation failed after thread creation');
    }

    const data = {
      messageId: dbMessage?.id || (undefined as any),
      threadId: thread.id
    };

    set({ isCreatingThread: false }, false, n('creatingThread/end'));
    return data;
  },

  /**
   * Fetch threads for a specific topic
   * 🔄 MIGRATED: Direct DB call instead of threadService.getThreads()
   */
  internal_fetchThreads: async (topicId) => {
    if (!topicId || isDeprecatedEdition) return;

    try {
      const userId = getUserId();
      const dbThreads = await DB.ListThreadsByTopic({
        topicId: topicId,  // Plain string, not NullString
        userId,
      });
      const threads = dbThreads.map(mapThread);
      const nextMap = { ...get().threadMaps, [topicId]: threads };

      // no need to update map if the threads have been init and the map is the same
      if (get().threadsInit && isEqual(nextMap, get().threadMaps)) return;

      console.log('[Thread] Fetched threads via direct DB', { topicId, count: threads.length });

      set(
        { threadMaps: nextMap, threadsInit: true },
        false,
        n('internal_fetchThreads', { topicId }),
      );
    } catch (error) {
      console.error('[internal_fetchThreads] Error fetching threads:', error);
    }
  },

  /**
   * Refresh threads from database - direct fetch without SWR cache invalidation
   * 🔄 MIGRATED: Direct DB call instead of threadService.getThreads()
   */
  refreshThreads: async () => {
    const topicId = get().activeTopicId;
    if (!topicId) return;

    try {
      const userId = getUserId();
      const dbThreads = await DB.ListThreadsByTopic({
        topicId: topicId,
        userId,
      });
      const threads = dbThreads.map(mapThread);
      const nextMap = { ...get().threadMaps, [topicId]: threads };

      console.log('[Thread] Refreshed threads via direct DB', { topicId, count: threads.length });

      set(
        { threadMaps: nextMap, threadsInit: true },
        false,
        n('refreshThreads', { topicId }),
      );
    } catch (error) {
      console.error('[refreshThreads] Error refreshing threads:', error);
    }
  },
  removeThread: async (id) => {
    const currentActiveThreadId = get().activeThreadId;
    console.debug('[chatThread.removeThread] Removing thread:', {
      threadId: id,
      currentActiveThreadId,
      willClearActiveThreadId: currentActiveThreadId === id,
    });

    const userId = getUserId();
    await DB.DeleteThread({
      id,
      userId,
    });

    console.log('[Thread] Deleted thread via direct DB', { threadId: id });

    await get().refreshThreads();

    if (get().activeThreadId === id) {
      console.debug('[chatThread.removeThread] Clearing activeThreadId because removed thread was active');
      set({ activeThreadId: undefined });
    }
  },
  switchThread: async (id) => {
    const previousActiveThreadId = get().activeThreadId;
    console.debug('[chatThread.switchThread] Switching thread:', {
      previousActiveThreadId,
      newThreadId: id,
      activeTopicId: get().activeTopicId,
    });
    set({ activeThreadId: id }, false, n('toggleTopic'));
    console.debug('[chatThread.switchThread] After switch:', {
      activeThreadId: get().activeThreadId,
    });
  },
  updateThreadTitle: async (id, title) => {
    await get().internal_updateThread(id, { title });
  },

  summaryThreadTitle: async (threadId, messages) => {
    const { internal_updateThreadTitleInSummary, internal_updateThreadLoading } = get();
    const portalThread = threadSelectors.currentPortalThread(get());
    if (!portalThread) return;

    internal_updateThreadTitleInSummary(threadId, LOADING_FLAT);

    let output = '';
    const threadConfig = systemAgentSelectors.thread(useUserStore.getState());

    // Limit input messages to prevent AI confusion with long conversations
    // For title generation, we only need recent context, not entire conversation
    const limitedMessages = messages.slice(-1); // Last 6 messages max

    await chatService.fetchPresetTaskResult({
      onError: () => {
        internal_updateThreadTitleInSummary(threadId, portalThread.title);
      },
      onFinish: async (text) => {
        await get().internal_updateThread(threadId, { title: text });
      },
      onLoadingChange: (loading) => {
        internal_updateThreadLoading(threadId, loading);
      },
      onMessageHandle: (chunk) => {
        switch (chunk.type) {
          case 'text': {
            output += chunk.text;
          }
        }

        internal_updateThreadTitleInSummary(threadId, output);
      },
      params: merge(threadConfig, chainSummaryTitle(limitedMessages, globalHelpers.getCurrentLanguage()), {
        stream: false, // Thread title generation doesn't need streaming
      }),
    });
  },

  // Internal process method of the topics
  internal_updateThreadTitleInSummary: (id, title) => {
    get().internal_dispatchThread(
      { type: 'updateThread', id, value: { title } },
      'updateThreadTitleInSummary',
    );
  },

  internal_updateThreadLoading: (id, loading) => {
    set(
      (state) => {
        if (loading) return { threadLoadingIds: [...state.threadLoadingIds, id] };

        return { threadLoadingIds: state.threadLoadingIds.filter((i) => i !== id) };
      },
      false,
      n('updateThreadLoading'),
    );
  },

  internal_updateThread: async (id, data) => {
    get().internal_dispatchThread({ type: 'updateThread', id, value: data });

    get().internal_updateThreadLoading(id, true);

    // 🔄 MIGRATED: Direct DB call instead of threadService.updateThread()
    const userId = getUserId();
    const now = currentTimestampMs();

    await DB.UpdateThread({
      id,
      userId,
      title: data.title ? toNullString(data.title) : undefined,
      status: data.status ? toNullString(data.status) : undefined,
      lastActiveAt: data.lastActiveAt ? new Date(data.lastActiveAt).getTime() : now,
      updatedAt: now,
    } as any);

    console.log('[Thread] Updated thread via direct DB', { id });

    await get().refreshThreads();
    get().internal_updateThreadLoading(id, false);
  },

  internal_dispatchThread: (payload, action) => {
    const nextThreads = threadReducer(threadSelectors.currentTopicThreads(get()), payload);
    const nextMap = { ...get().threadMaps, [get().activeTopicId!]: nextThreads };

    // no need to update map if is the same
    if (isEqual(nextMap, get().threadMaps)) return;

    set({ threadMaps: nextMap }, false, action ?? n(`dispatchThread/${payload.type}`));
  },
});
