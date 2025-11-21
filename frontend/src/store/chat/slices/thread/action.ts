/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
// Disable the auto sort key eslint rule to make the code more logic and readable
import { LOADING_FLAT, isDeprecatedEdition } from '@/const';
import { chainSummaryTitle } from '@/prompts';
import {
  CreateMessageParams,
  SendThreadMessageParams,
  ThreadItem,
  ThreadType,
  UIChatMessage,
} from '@/types';
import isEqual from 'fast-deep-equal';
import { StateCreator } from 'zustand/vanilla';

import { chatService } from '@/services/chat';
import { threadService } from '@/services/thread';
import { threadSelectors } from './selectors';
import { ChatStore } from '@/store/chat/store';
import { globalHelpers } from '@/store/global/helpers';
import { useUserStore } from '@/store/user';
import { systemAgentSelectors } from '@/store/user/selectors';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';

import { ThreadDispatch, threadReducer } from './reducer';

const n = setNamespace('thd');

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

    const data = await threadService.createThreadWithMessage({
      topicId,
      sourceMessageId,
      type,
      message,
    });
    set({ isCreatingThread: false }, false, n('creatingThread/end'));

    return data;
  },

  /**
   * Fetch threads for a specific topic
   * Direct Zustand implementation (no SWR) for better performance
   */
  internal_fetchThreads: async (topicId) => {
    if (!topicId || isDeprecatedEdition) return;

    try {
      const threads = await threadService.getThreads(topicId);
      const nextMap = { ...get().threadMaps, [topicId]: threads };

      // no need to update map if the threads have been init and the map is the same
      if (get().threadsInit && isEqual(nextMap, get().threadMaps)) return;

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
   */
  refreshThreads: async () => {
    const topicId = get().activeTopicId;
    if (!topicId) return;

    try {
      const threads = await threadService.getThreads(topicId);
      const nextMap = { ...get().threadMaps, [topicId]: threads };

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
    await threadService.removeThread(id);
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
    await threadService.updateThread(id, data);
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
