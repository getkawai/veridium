/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
// Note: To make the code more logic and readable, we just disable the auto sort key eslint rule
// DON'T REMOVE THE FIRST LINE
import { chainSummaryTitle } from '@/prompts';
import { TraceNameMap, UIChatMessage } from '@/types';
import isEqual from 'fast-deep-equal';
import { t } from 'i18next';
import { produce } from 'immer';
import { StateCreator } from 'zustand/vanilla';

import { message } from '@/components/AntdStaticMethods';
import { LOADING_FLAT } from '@/const/message';
import { chatService } from '@/services/chat';
import { messageService } from '@/services/message';
import { topicService } from '@/services/topic';
import { CreateTopicParams } from '@/services/topic/type';
import type { ChatStore } from '@/store/chat';
import type { ChatStoreState } from '@/store/chat/initialState';
import { messageMapKey } from '@/store/chat/utils/messageMapKey';
import { globalHelpers } from '@/store/global/helpers';
import { useSessionStore } from '@/store/session';
import { sessionSelectors } from '@/store/session/selectors';
import { useUserStore } from '@/store/user';
import { systemAgentSelectors } from '@/store/user/selectors';
import { ChatTopic } from '@/types/topic';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';

import { chatSelectors } from '../message/selectors';
import { ChatTopicDispatch, topicReducer } from './reducer';
import { topicSelectors } from './selectors';

const n = setNamespace('t');

export interface ChatTopicAction {
  favoriteTopic: (id: string, favState: boolean) => Promise<void>;
  openNewTopicOrSaveTopic: () => Promise<void>;
  refreshTopic: () => Promise<void>;
  removeAllTopics: () => Promise<void>;
  removeSessionTopics: () => Promise<void>;
  removeGroupTopics: (groupId: string) => Promise<void>;
  removeTopic: (id: string) => Promise<void>;
  removeUnstarredTopic: () => Promise<void>;
  saveToTopic: (sessionId?: string, groupId?: string) => Promise<string | undefined>;
  createTopic: (sessionId?: string, groupId?: string) => Promise<string | undefined>;

  autoRenameTopicTitle: (id: string) => Promise<void>;
  duplicateTopic: (id: string) => Promise<void>;
  summaryTopicTitle: (topicId: string, messages: UIChatMessage[]) => Promise<void>;
  switchTopic: (id?: string, skipRefreshMessage?: boolean) => Promise<void>;
  updateTopicTitle: (id: string, title: string) => Promise<void>;
  internal_fetchTopics: (
    containerId: string,
  ) => Promise<void>;
  internal_searchTopics: (
    keywords?: string,
    sessionId?: string,
    groupId?: string,
  ) => Promise<void>;

  internal_updateTopicTitleInSummary: (id: string, title: string) => void;
  internal_updateTopicLoading: (id: string, loading: boolean) => void;
  internal_createTopic: (params: CreateTopicParams) => Promise<string>;
  internal_updateTopic: (id: string, data: Partial<ChatTopic>) => Promise<void>;
  internal_dispatchTopic: (payload: ChatTopicDispatch, action?: any) => void;
}

export const chatTopic: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatTopicAction
> = (set, get) => ({
  // create
  openNewTopicOrSaveTopic: async () => {
    const { switchTopic, saveToTopic, refreshMessages, activeTopicId } = get();
    const hasTopic = !!activeTopicId;

    if (hasTopic) switchTopic();
    else {
      await saveToTopic();
      refreshMessages();
    }
  },

  createTopic: async (sessionId, groupId) => {
    const { activeId, activeSessionType, internal_createTopic } = get();

    const messages = chatSelectors.activeBaseChats(get());

    set({ creatingTopic: true }, false, n('creatingTopic/start'));
    const topicId = await internal_createTopic({
      title: t('defaultTitle', { ns: 'topic' }),
      messages: messages.map((m) => m.id),
      ...(activeSessionType === 'group'
        ? { groupId: groupId || activeId }
        : { sessionId: sessionId || activeId }),
    });
    set({ creatingTopic: false }, false, n('creatingTopic/end'));

    return topicId;
  },

  saveToTopic: async (sessionId, groupId) => {
    // if there is no message, stop
    const messages = chatSelectors.activeBaseChats(get());
    if (messages.length === 0) return;

    const { activeId, activeSessionType, summaryTopicTitle, internal_createTopic } = get();

    // 1. create topic and bind these messages
    const topicId = await internal_createTopic({
      title: t('defaultTitle', { ns: 'topic' }),
      messages: messages.map((m) => m.id),
      ...(activeSessionType === 'group'
        ? { groupId: groupId || activeId }
        : { sessionId: sessionId || activeId }),
    });

    get().internal_updateTopicLoading(topicId, true);
    // 2. auto summary topic Title
    // we don't need to wait for summary, just let it run async
    summaryTopicTitle(topicId, messages);

    // Clear supervisor todos for temporary topic in current container after saving
    try {
      const { activeId, activeSessionType } = get();
      let isGroupSession = activeSessionType === 'group';
      if (activeSessionType === undefined) {
        const sessionStore = useSessionStore.getState();
        isGroupSession = sessionSelectors.isCurrentSessionGroupSession(sessionStore);
      }

      if (isGroupSession) {
        set(
          produce((state: ChatStoreState) => {
            state.supervisorTodos[messageMapKey(groupId || activeId, null)] = [];
          }),
          false,
          n('resetSupervisorTodosOnSaveToTopic', { groupId: groupId || activeId }),
        );
      }
    } catch (error) {
      if (true) {
        console.error('Failed to reset supervisor todos on save to topic:', error);
      }
    }

    return topicId;
  },

  duplicateTopic: async (id) => {
    const { refreshTopic, switchTopic } = get();

    const topic = topicSelectors.getTopicById(id)(get());
    if (!topic) return;

    const newTitle = t('duplicateTitle', { ns: 'chat', title: topic?.title });

    message.loading({
      content: t('duplicateLoading', { ns: 'topic' }),
      key: 'duplicateTopic',
      duration: 0,
    });

    const newTopicId = await topicService.cloneTopic(id, newTitle);
    await refreshTopic();
    message.destroy('duplicateTopic');
    message.success(t('duplicateSuccess', { ns: 'topic' }));

    await switchTopic(newTopicId);
  },
  // update
  summaryTopicTitle: async (topicId, messages) => {
    const { internal_updateTopicTitleInSummary, internal_updateTopicLoading } = get();
    const topic = topicSelectors.getTopicById(topicId)(get());
    if (!topic) return;

    internal_updateTopicTitleInSummary(topicId, LOADING_FLAT);

    let output = '';

    // Get current agent for topic
    const topicConfig = systemAgentSelectors.topic(useUserStore.getState());

    // Automatically summarize the topic title
    await chatService.fetchPresetTaskResult({
      onError: () => {
        internal_updateTopicTitleInSummary(topicId, topic.title);
      },
      onFinish: async (text) => {
        await get().internal_updateTopic(topicId, { title: text });
      },
      onLoadingChange: (loading) => {
        internal_updateTopicLoading(topicId, loading);
      },
      onMessageHandle: (chunk) => {
        switch (chunk.type) {
          case 'text': {
            output += chunk.text;
          }
        }

        internal_updateTopicTitleInSummary(topicId, output);
      },
      params: merge(topicConfig, chainSummaryTitle(messages, globalHelpers.getCurrentLanguage()), {
        stream: false, // Topic generation doesn't need streaming
      }),
      trace: get().getCurrentTracePayload({ traceName: TraceNameMap.SummaryTopicTitle, topicId }),
    });
  },
  favoriteTopic: async (id, favorite) => {
    await get().internal_updateTopic(id, { favorite });
  },

  updateTopicTitle: async (id, title) => {
    await get().internal_updateTopic(id, { title });
  },

  autoRenameTopicTitle: async (id) => {
    const { activeId: sessionId, summaryTopicTitle, internal_updateTopicLoading } = get();

    internal_updateTopicLoading(id, true);
    const messages = await messageService.getMessages(sessionId, id);

    await summaryTopicTitle(id, messages);
    internal_updateTopicLoading(id, false);
  },

  // query
  /**
   * Fetch topics for a specific container (session or group)
   * Direct Zustand implementation (no SWR) for better performance
   */
  internal_fetchTopics: async (containerId) => {
    if (!containerId) return;

    try {
      console.debug('[internal_fetchTopics] Fetching topics for containerId:', containerId);
      const topics = await topicService.getTopics({ containerId });
      console.debug('[internal_fetchTopics] Fetched topics:', topics.length, 'topics');

      const nextMap = { ...get().topicMaps, [containerId]: topics };

      // no need to update map if the topics have been init and the map is the same
      if (get().topicsInit && isEqual(nextMap, get().topicMaps)) {
        console.debug('[internal_fetchTopics] Skipping update - maps are equal');
        return;
      }

      console.debug('[internal_fetchTopics] Updating topicMaps');
      set(
        { topicMaps: nextMap, topicsInit: true },
        false,
        n('internal_fetchTopics', { containerId }),
      );
    } catch (error) {
      console.error('[internal_fetchTopics] Error fetching topics:', error);
    }
  },
  
  /**
   * Search topics by keywords
   * Direct implementation (no SWR)
   */
  internal_searchTopics: async (keywords, sessionId, groupId) => {
    if (!keywords) {
      set({ searchTopics: [], isSearchingTopic: false }, false, n('internal_searchTopics/clear'));
      return;
    }

    try {
      set({ isSearchingTopic: true }, false, n('internal_searchTopics/start'));
      const data = await topicService.searchTopics(keywords, sessionId, groupId);
      set(
        { searchTopics: data, isSearchingTopic: false },
        false,
        n('internal_searchTopics', { keywords }),
      );
    } catch (error) {
      console.error('[internal_searchTopics] Error searching topics:', error);
      set({ isSearchingTopic: false }, false, n('internal_searchTopics/error'));
    }
  },

  switchTopic: async (id, skipRefreshMessage) => {
    const previousActiveThreadId = get().activeThreadId;
    console.debug('[chatTopic.switchTopic] Switching topic:', {
      previousTopicId: get().activeTopicId,
      newTopicId: id,
      previousActiveThreadId,
      clearingActiveThreadId: true,
    });
    set(
      { activeTopicId: !id ? (null as any) : id, activeThreadId: undefined },
      false,
      n('toggleTopic'),
    );
    console.debug('[chatTopic.switchTopic] After switch:', {
      newTopicId: id,
      activeThreadId: get().activeThreadId,
      wasCleared: previousActiveThreadId !== undefined,
    });

    // Reset supervisor todos when switching topics in group chats
    try {
      const { activeId, activeSessionType, internal_cancelSupervisorDecision } = get();
      // Determine group session robustly (cached flag or from session store)
      let isGroupSession = activeSessionType === 'group';
      if (activeSessionType === undefined) {
        const sessionStore = useSessionStore.getState();
        isGroupSession = sessionSelectors.isCurrentSessionGroupSession(sessionStore);
      }

      if (isGroupSession) {
        const newKey = messageMapKey(activeId, id ?? null);
        set(
          produce((state: ChatStoreState) => {
            state.supervisorTodos[newKey] = [];
          }),
          false,
          n('resetSupervisorTodosOnTopicSwitch', { groupId: activeId, topicId: id ?? null }),
        );

        // Also cancel any pending supervisor decisions tied to this group
        internal_cancelSupervisorDecision?.(activeId);
      }
    } catch {
      // no-op: resetting todos should not block topic switching
    }

    if (skipRefreshMessage) return;
    await get().refreshMessages();
  },
  // delete
  removeSessionTopics: async () => {
    const { switchTopic, activeId, refreshTopic } = get();

    await topicService.removeTopics(activeId);
    await refreshTopic();

    // switch to default topic
    switchTopic();
  },

  removeGroupTopics: async (groupId: string) => {
    const { switchTopic, refreshTopic } = get();

    // Get topics for this specific group from the topic map
    const groupTopics = get().topicMaps[groupId] || [];
    const topicIds = groupTopics.map((t) => t.id);

    if (topicIds.length > 0) {
      await topicService.batchRemoveTopics(topicIds);
    }

    await refreshTopic();

    // switch to default topic
    switchTopic();
  },
  removeAllTopics: async () => {
    const { refreshTopic } = get();

    await topicService.removeAllTopic();
    await refreshTopic();
  },
  removeTopic: async (id) => {
    const { activeId, activeTopicId, switchTopic, refreshTopic } = get();

    // remove messages in the topic
    // TODO: Need to remove because server service don't need to call it
    await messageService.removeMessagesByAssistant(activeId, id);

    // remove topic
    await topicService.removeTopic(id);
    await refreshTopic();

    // switch bach to default topic
    if (activeTopicId === id) switchTopic();
  },
  removeUnstarredTopic: async () => {
    const { refreshTopic, switchTopic } = get();
    const topics = topicSelectors.currentUnFavTopics(get());

    await topicService.batchRemoveTopics(topics.map((t) => t.id));
    await refreshTopic();

    // 切换到默认 topic
    switchTopic();
  },

  // Internal process method of the topics
  internal_updateTopicTitleInSummary: (id, title) => {
    get().internal_dispatchTopic(
      { type: 'updateTopic', id, value: { title } },
      'updateTopicTitleInSummary',
    );
  },
  
  /**
   * Refresh topics from database - direct fetch without SWR cache invalidation
   */
  refreshTopic: async () => {
    const { activeId } = get();
    if (!activeId) return;

    try {
      console.debug('[refreshTopic] Fetching topics for activeId:', activeId);
      const topics = await topicService.getTopics({ containerId: activeId });
      console.debug('[refreshTopic] Fetched topics:', topics.length, 'topics');

      const nextMap = { ...get().topicMaps, [activeId]: topics };
      set(
        { topicMaps: nextMap, topicsInit: true },
        false,
        n('refreshTopic', { activeId }),
      );
    } catch (error) {
      console.error('[refreshTopic] Error refreshing topics:', error);
    }
  },

  internal_updateTopicLoading: (id, loading) => {
    set(
      (state) => {
        if (loading) return { topicLoadingIds: [...state.topicLoadingIds, id] };

        return { topicLoadingIds: state.topicLoadingIds.filter((i) => i !== id) };
      },
      false,
      n('updateTopicLoading'),
    );
  },

  internal_updateTopic: async (id, data) => {
    get().internal_dispatchTopic({ type: 'updateTopic', id, value: data });

    get().internal_updateTopicLoading(id, true);
    await topicService.updateTopic(id, data);
    await get().refreshTopic();
    get().internal_updateTopicLoading(id, false);
  },
  internal_createTopic: async (params) => {
    const tmpId = Date.now().toString();
    console.debug('[internal_createTopic] Creating topic with tmpId:', tmpId, 'params:', params);
    
    get().internal_dispatchTopic(
      { type: 'addTopic', value: { ...params, id: tmpId } },
      'internal_createTopic',
    );

    get().internal_updateTopicLoading(tmpId, true);
    console.debug('[internal_createTopic] Calling topicService.createTopic...');
    const topicId = await topicService.createTopic(params);
    console.debug('[internal_createTopic] Topic created with ID:', topicId);
    get().internal_updateTopicLoading(tmpId, false);

    get().internal_updateTopicLoading(topicId, true);
    console.debug('[internal_createTopic] Calling refreshTopic...');
    await get().refreshTopic();
    console.debug('[internal_createTopic] refreshTopic completed');
    get().internal_updateTopicLoading(topicId, false);

    console.debug('[internal_createTopic] Final topicMaps:', get().topicMaps);
    console.debug('[internal_createTopic] Current topics for activeId:', get().topicMaps[get().activeId]);

    return topicId;
  },

  internal_dispatchTopic: (payload, action) => {
    const nextTopics = topicReducer(topicSelectors.currentTopics(get()), payload);
    const nextMap = { ...get().topicMaps, [get().activeId]: nextTopics };

    // no need to update map if is the same
    if (isEqual(nextMap, get().topicMaps)) return;

    set({ topicMaps: nextMap }, false, action ?? n(`dispatchTopic/${payload.type}`));
  },
});
