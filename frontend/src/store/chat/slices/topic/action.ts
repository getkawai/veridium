/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
// Note: To make the code more logic and readable, we just disable the auto sort key eslint rule
// DON'T REMOVE THE FIRST LINE
import { chainSummaryTitle } from '@/prompts';
import { ChatFileItem, TraceNameMap, UIChatMessage } from '@/types';
import isEqual from 'fast-deep-equal';
import { t } from 'i18next';
import { produce } from 'immer';
import { StateCreator } from 'zustand/vanilla';

import { message } from '@/components/AntdStaticMethods';
import { LOADING_FLAT } from '@/const/message';
import { chatService } from '@/services/chat';
import { messageService } from '@/services/message';
import type { ChatStore } from '@/store/chat';

// 🔄 MIGRATED: Direct imports for message operations
import { MessageModel } from '@/database/models/message';

// Local type definition (migrated from @/services/topic/type)
interface CreateTopicParams {
  favorite?: boolean;
  groupId?: string | null;
  messages?: string[];
  sessionId?: string | null;
  title: string;
}
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
import { getUserId, mapTopicFromDB, toDbSessionId } from '@/store/topic/helpers';
import { toNullString } from '@/types/database';
import { DB } from '@/types/database';
import { chatSelectors } from '../message/selectors';
import { ChatTopicDispatch, topicReducer } from './reducer';
import { topicSelectors } from './selectors';
import { fileService } from '@/services/file';

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

  // MASIH DIGUNAKAN: Membuat topic baru dengan title default "默认话题" (defaultTitle)
  // Fungsi ini digunakan untuk:
  // - Membuat topic baru secara manual (belum ada pemanggilan langsung dari UI saat ini)
  // - Menggunakan title default dari translation 'defaultTitle'
  // - Tidak melakukan auto-summary title (berbeda dengan saveToTopic)
  // - Mengikat semua messages dari temporary chat ke topic yang baru dibuat
  // Note: Fungsi ini berbeda dengan saveToTopic yang juga melakukan auto-summary title
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

  // MASIH DIGUNAKAN: Menyimpan temporary chat ke topic permanen
  // Dipanggil dari:
  // 1. SaveTopic button (features/ChatInput/ActionBar/SaveTopic/index.tsx) melalui openNewTopicOrSaveTopic
  // 2. Hotkey untuk save topic (HotkeyEnum.SaveTopic)
  // 
  // Alur kerja:
  // 1. Membuat topic baru dengan title default "默认话题" (defaultTitle)
  // 2. Mengikat semua messages dari temporary chat ke topic yang baru dibuat
  // 3. Melakukan auto-summary title secara async (menggunakan LLM untuk generate title yang lebih deskriptif)
  // 4. Membersihkan supervisor todos untuk temporary topic di group session
  // 
  // Perbedaan dengan createTopic:
  // - saveToTopic melakukan auto-summary title (async)
  // - saveToTopic membersihkan supervisor todos
  // - saveToTopic lebih sering digunakan dari UI (via SaveTopic button)
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

    const userId = getUserId();
    const newTopicId = crypto.randomUUID();
    const now = Date.now();
    
    await DB.DuplicateTopic({
      id: newTopicId,
      title: toNullString(newTitle),
      createdAt: now,
      updatedAt: now,
      id2: id,
      userId,
    });
    
    console.log('[Topic] Duplicated topic via direct DB', { sourceId: id, newId: newTopicId });
    
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
    const messageModel = new MessageModel(DB, getUserId());
    const messages = await messageModel.query({ sessionId, topicId: id });

    const fileList = (await Promise.all(
      messages
        .flatMap((item) => item.files)
        .filter(Boolean)
        .map(async (id) => fileService.getFile(id!)),
    )) as ChatFileItem[];

    const result = messages.map((item) => ({
      ...item,
      imageList: fileList
        .filter((file) => item.files?.includes(file.id) && file.fileType.startsWith('image'))
        .map((file) => ({
          alt: file.name,
          id: file.id,
          url: file.url,
        })),
    }));

    await summaryTopicTitle(id, result);
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
      const userId = getUserId();
      const dbSessionId = toDbSessionId(containerId);
      
      const dbTopics = await DB.ListTopics({
        userId,
        sessionId: toNullString(dbSessionId),
        limit: 1000,
        offset: 0,
      });
      
      const topics = dbTopics.map(mapTopicFromDB);
      
      console.log('[Topic] Fetched topics via direct DB', { containerId, count: topics.length });

      const nextMap = { ...get().topicMaps, [containerId]: topics };

      // no need to update map if the topics have been init and the map is the same
      if (get().topicsInit && isEqual(nextMap, get().topicMaps)) {
        console.debug('[internal_fetchTopics] Skipping update - maps are equal');
        return;
      }

      console.debug('[internal_fetchTopics] Updating topicMaps');
      set(
        { topicMaps: nextMap as Record<string, ChatTopic[]>, topicsInit: true },
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
      
      const userId = getUserId();
      const searchPattern = `%${keywords}%`;
      const containerId = sessionId || groupId || '';
      
      const dbTopics = await DB.SearchTopicsByTitle({
        userId,
        title: toNullString(searchPattern),
        column3: containerId,
        sessionId: toNullString(toDbSessionId(sessionId)),
        groupId: toNullString(groupId),
      });
      
      const data = dbTopics.map(mapTopicFromDB);
      
      console.log('[Topic] Searched topics via direct DB', { keywords, count: data.length });
      
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

    const userId = getUserId();
    const dbSessionId = toDbSessionId(activeId);
    
    await DB.DeleteTopicsBySession({
      userId,
      sessionId: toNullString(dbSessionId),
    });
    
    console.log('[Topic] Deleted session topics via direct DB', { sessionId: activeId });
    
    await refreshTopic();

    // switch to default topic
    switchTopic();
  },

  removeGroupTopics: async (groupId: string) => {
    const { switchTopic, refreshTopic } = get();

    const userId = getUserId();
    
    await DB.DeleteTopicsByGroup({
      userId,
      groupId: toNullString(groupId),
    });
    
    console.log('[Topic] Deleted group topics via direct DB', { groupId });

    await refreshTopic();

    // switch to default topic
    switchTopic();
  },
  removeAllTopics: async () => {
    const { refreshTopic } = get();

    const userId = getUserId();
    
    await DB.DeleteAllTopics(userId);
    
    console.log('[Topic] Deleted all topics via direct DB');
    
    await refreshTopic();
  },
  removeTopic: async (id) => {
    const { activeId, activeTopicId, switchTopic, refreshTopic } = get();

    // remove messages in the topic
    // TODO: Need to remove because server service don't need to call it
    await messageService.removeMessagesByAssistant(activeId, id);

    const userId = getUserId();
    
    await DB.DeleteTopic({ id, userId });
    
    console.log('[Topic] Deleted topic via direct DB', { id });
    await refreshTopic();

    // switch bach to default topic
    if (activeTopicId === id) switchTopic();
  },
  removeUnstarredTopic: async () => {
    const { refreshTopic, switchTopic } = get();
    const topics = topicSelectors.currentUnFavTopics(get());

    const userId = getUserId();
    
    // Delete all unstarred topics
    await Promise.all(
      topics.map((topic) => DB.DeleteTopic({ id: topic.id, userId }))
    );
    
    console.log('[Topic] Deleted unstarred topics via direct DB', { count: topics.length });
    
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
      
      const userId = getUserId();
      const dbTopics = await DB.ListTopics({
        userId,
        sessionId: toNullString(activeId),
        limit: 1000,
        offset: 0,
      });
      const topics = dbTopics.map(mapTopicFromDB);
      
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
    
    const userId = getUserId();
    const now = Date.now();
    
    await DB.UpdateTopic({
      id,
      userId,
      title: data.title ? toNullString(data.title) : undefined,
      historySummary: data.historySummary !== undefined ? toNullString(data.historySummary) : undefined,
      metadata: data.metadata ? toNullString(JSON.stringify(data.metadata)) : undefined,
      updatedAt: now,
    } as any);
    
    console.log('[Topic] Updated topic via direct DB', { id });
    
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
    
    const userId = getUserId();
    const topicId = crypto.randomUUID();
    const now = Date.now();
    
    await DB.CreateTopic({
      id: topicId,
      title: toNullString(params.title || 'Untitled'),
      favorite: 0,
      sessionId: toNullString(toDbSessionId(params.sessionId)),
      groupId: toNullString(params.groupId),
      userId,
      clientId: toNullString(''),
      historySummary: toNullString(''),
      metadata: toNullString(JSON.stringify({})),
      createdAt: now,
      updatedAt: now,
    });
    
    // Update messages with topic ID
    if (params.messages && params.messages.length > 0) {
      await DB.UpdateMessagesTopicId({
        topicId: toNullString(topicId),
        userId,
        ids: params.messages,
      });
    }
    
    console.log('[Topic] Created topic via direct DB', { topicId });
    
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
