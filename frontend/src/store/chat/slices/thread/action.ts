/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
// Disable the auto sort key eslint rule to make the code more logic and readable
import { isDeprecatedEdition } from '@/const';
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

import { threadSelectors } from './selectors';
import { ChatStore } from '@/store/chat/store';
import { setNamespace } from '@/utils/storeDebug';
import { DB, toNullString, getNullableString, currentTimestampMs, Thread } from '@/types/database';

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
    const {
      activeTopicId,
      activeId,
      threadStartMessageId,
      newThreadMode,
      portalThreadId,
    } = get();

    if (!activeId || !activeTopicId) {
      console.error('[Thread] sendThreadMessage: Missing activeId or activeTopicId');
      return;
    }

    // if message is empty, then stop
    if (!message) return;

    set({ isCreatingThreadMessage: true }, false, n('creatingThreadMessage/start'));

    const newMessage: CreateMessageParams = {
      content: message,
      role: 'user',
      sessionId: activeId,
      topicId: activeTopicId,
      threadId: portalThreadId,
    };

    // if there is no portalThreadId, then create a thread and then append message
    let currentThreadId: string | undefined = portalThreadId;
    if (!portalThreadId) {
      if (!threadStartMessageId) {
        console.error('[Thread] sendThreadMessage: Missing threadStartMessageId');
        set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
        return;
      }

      let threadResult;
      try {
        threadResult = await get().createThread({
          message: newMessage,
          sourceMessageId: threadStartMessageId,
          topicId: activeTopicId,
          type: newThreadMode,
        });
      } catch (error) {
        console.error('[sendThreadMessage] Thread creation threw error:', error);
        set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
        return;
      }

      // Check if thread creation failed
      if (!threadResult || !threadResult.threadId) {
        console.error('[sendThreadMessage] Failed to create thread:', threadResult);
        set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
        return;
      }

      const { threadId } = threadResult;
      currentThreadId = threadId;

      // mark the portal in thread mode
      await get().refreshThreads();
      await get().refreshMessages();

      get().openThreadInPortal(threadId, threadStartMessageId);
    }

    // Determine the thread ID to use
    const finalThreadId = currentThreadId || get().portalThreadId;
    if (!finalThreadId) {
      console.error('[sendThreadMessage] No threadId available');
      set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
      return;
    }

    // 1. Create optimistic user message
    get().internal_createTmpMessage({
      ...newMessage,
      threadId: finalThreadId,
    });

    // 2. Create optimistic assistant message (loading state)
    const tempAssistantMessageId = get().internal_createTmpMessage({
      role: 'assistant',
      content: '...',
      threadId: finalThreadId,
      sessionId: activeId,
      topicId: activeTopicId,
    });
    get().internal_toggleMessageLoading(true, tempAssistantMessageId);

    set({ isCreatingThreadMessage: true }, false, n('sendingMessage/start'));

    // Use backend agent chat to generate AI response
    try {
      // TODO: implement this
      // const backendAgentChat = await import('@/services/backendAgentChat').then(m => m.backendAgentChat);

      // await backendAgentChat.sendMessage({
      //   session_id: activeId,
      //   user_id: userId,
      //   message: message,
      //   topic_id: activeTopicId,
      //   thread_id: finalThreadId,
      //   stream: true, // Enable streaming
      // });
      // Just refresh to get them from DB (this will replace temp messages)
      await get().refreshMessages();

      // Auto-generate thread title if this is a new thread
      if (!portalThreadId) {
        const portalThread = threadSelectors.currentPortalThread(get());

        if (portalThread) {
          const chats = threadSelectors.portalAIChats(get());
          await get().summaryThreadTitle(portalThread.id, chats);
        }
      }

    } catch (error) {
      console.error('[sendThreadMessage] Failed to get AI response:', error);
      // On error, we should probably keep the user message but mark it as failed?
      // For now, let's just remove the loading assistant message
      get().internal_dispatchMessage({
        type: 'deleteMessages',
        ids: [tempAssistantMessageId],
      });
    } finally {
      // Always clear loading state
      set({ isCreatingThreadMessage: false }, false, n('creatingThreadMessage/stop'));
    }
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

    // Don't create message here - backend will save it when sendMessage is called
    // Just return the thread ID
    const data = {
      messageId: undefined as any, // Backend will create the message
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
      const dbThreads = await DB.ListThreadsByTopic(topicId);
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
      const dbThreads = await DB.ListThreadsByTopic(topicId);
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

    await DB.DeleteThread(id);

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
    // const { internal_updateThreadTitleInSummary, internal_updateThreadLoading } = get();
    // const portalThread = threadSelectors.currentPortalThread(get());
    // if (!portalThread) return;

    // internal_updateThreadTitleInSummary(threadId, LOADING_FLAT);

    // let output = '';
    // const threadConfig = systemAgentSelectors.thread(useUserStore.getState());

    // // Limit input messages to prevent AI confusion with long conversations
    // // For title generation, we only need recent context, not entire conversation
    // const limitedMessages = messages.slice(-1); // Last 6 messages max

    // await chatService.fetchPresetTaskResult({
    //   onError: () => {
    //     internal_updateThreadTitleInSummary(threadId, portalThread.title);
    //   },
    //   onFinish: async (text) => {
    //     await get().internal_updateThread(threadId, { title: text });
    //   },
    //   onLoadingChange: (loading) => {
    //     internal_updateThreadLoading(threadId, loading);
    //   },
    //   onMessageHandle: (chunk) => {
    //     switch (chunk.type) {
    //       case 'text': {
    //         output += chunk.text;
    //       }
    //     }

    //     internal_updateThreadTitleInSummary(threadId, output);
    //   },
    //   params: merge(threadConfig, chainSummaryTitle(limitedMessages, globalHelpers.getCurrentLanguage()), {
    //     stream: false, // Thread title generation doesn't need streaming
    //   }),
    // });
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

    const now = currentTimestampMs();

    // Extract string values from potential NullString objects
    const titleValue = data.title
      ? (typeof data.title === 'object' && data.title !== null && 'String' in data.title ? (data.title as any).String : data.title)
      : undefined;

    await DB.UpdateThread({
      id,
      title: titleValue,
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
