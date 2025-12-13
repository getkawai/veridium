import isEqual from 'fast-deep-equal';
import { t } from 'i18next';
import type { PartialDeep } from 'type-fest';
import { StateCreator } from 'zustand/vanilla';

import { message } from '@/components/AntdStaticMethods';
import { MESSAGE_CANCEL_FLAT } from '@/const/message';
import { DEFAULT_AGENT_LOBE_SESSION, INBOX_SESSION_ID } from '@/const/session';
import { DEFAULT_CHAT_GROUP_CHAT_CONFIG } from '@/const/settings';
import { useAgentStore } from '@/store/agent';
import { getChatGroupStoreState } from '@/store/chatGroup';
import { useUserStore } from '@/store/user';
import { DB, toNullString, toNullInt, boolToInt, Session, getNullableString } from '@/types/database';
import { getUserId, mapAgentConfigFromDB } from '../../helpers';

import type { SessionStore } from '../../store';
import { settingsSelectors } from '@/store/user/selectors';
import { MetaData } from '@/types/meta';
import {
  LobeSession,
  LobeSessionGroups,
  LobeSessions,
  UpdateSessionParams,
} from '@/types/session';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';

import { SessionDispatch, sessionsReducer } from './reducers';
import { sessionSelectors } from './selectors';
import { sessionMetaSelectors } from './selectors/meta';

const n = setNamespace('session');

/* eslint-disable typescript-sort-keys/interface */
export interface SessionAction {
  /**
   * switch the session
   */
  switchSession: (sessionId: string) => void;
  /**
   * reset sessions to default
   */
  clearSessions: () => Promise<void>;
  /**
   * create a new session
   * @param agent
   * @returns sessionId
   */
  createSession: (
    session?: any,
    isSwitchSession?: boolean,
  ) => Promise<string>;

  duplicateSession: (id: string) => Promise<void>;
  triggerSessionUpdate: (id: string) => Promise<void>;
  updateSessionGroupId: (sessionId: string, groupId: string) => Promise<void>;
  updateSessionMeta: (meta: Partial<MetaData>) => void;

  /**
   * Pins or unpins a session.
   */
  pinSession: (id: string, pinned: boolean) => Promise<void>;
  /**
   * re-fetch the data
   */
  refreshSessions: () => Promise<void>;

  /**
   * load more sessions
   */
  loadMoreSessions: () => Promise<void>;

  /**
   * remove session
   * @param id - sessionId
   */
  removeSession: (id: string) => Promise<void>;

  updateSearchKeywords: (keywords: string) => void;

  internal_fetchSessions: (enabled: boolean, isLogin: boolean | undefined) => Promise<void>;
  internal_searchSessions: (keyword?: string) => Promise<void>;

  internal_dispatchSessions: (payload: SessionDispatch) => void;
  internal_updateSession: (id: string, data: Partial<UpdateSessionParams>) => Promise<void>;
  internal_processSessions: (
    sessions: LobeSessions,
    customGroups: LobeSessionGroups,
    actions?: string,
  ) => void;
  /* eslint-enable */
}

export const createSessionSlice: StateCreator<
  SessionStore,
  [['zustand/devtools', never]],
  [],
  SessionAction
> = (set, get) => ({
  clearSessions: async () => {
    const sessions = await DB.ListAllSessions();

    // Delete all sessions (except inbox)
    const sessionsToDelete = sessions.filter(s => s.slug !== 'inbox');
    await Promise.all(
      sessionsToDelete.map((session) => DB.DeleteSession(session.id))
    );

    console.log('[Session] Cleared all sessions via direct DB', { count: sessionsToDelete.length });

    await get().refreshSessions();
  },

  loadMoreSessions: async () => {
    const { sessionsPage, sessionsHasMore, internal_processSessions, sessions } = get();
    if (!sessionsHasMore) return;

    const nextPage = sessionsPage + 1;
    const limit = 10;
    const offset = nextPage * limit;

    try {
      // Use DB.ListSessions with limit and offset
      const dbSessions: Session[] = await DB.ListSessions({ limit, offset });
      const newSessions: Session[] = dbSessions;

      // Append new sessions
      const nextSessions = [...sessions, ...newSessions];

      // Update state
      set({
        sessionsPage: nextPage,
        sessionsHasMore: newSessions.length >= limit
      }, false, n('loadMoreSessions'));

      // Process for groups with all sessions (existing + new)
      get().internal_processSessions(nextSessions as any, get().sessionGroups);

    } catch (error) {
      console.error('[loadMoreSessions] Error:', error);
    }
  },

  createSession: async (agent, isSwitchSession = true) => {
    const { switchSession, refreshSessions } = get();

    // merge the defaultAgent in settings
    const defaultAgent = merge(
      DEFAULT_AGENT_LOBE_SESSION,
      settingsSelectors.defaultAgent(useUserStore.getState()),
    );

    const newSession: any = merge(defaultAgent, agent);

    const sessionId = crypto.randomUUID();
    const agentId = crypto.randomUUID();
    const now = Date.now();

    // Create session with agent
    await DB.CreateSession({
      id: sessionId,
      slug: '',
      title: toNullString(newSession.meta?.title),
      description: toNullString(newSession.meta?.description),
      avatar: toNullString(newSession.meta?.avatar),
      backgroundColor: toNullString(newSession.meta?.backgroundColor),
      type: toNullString('agent'),
      groupId: toNullString(newSession.group === 'default' ? '' : newSession.group),
      pinned: newSession.pinned ? 1 : 0,
      createdAt: now,
      updatedAt: now,
    });

    // Create agent config
    await DB.CreateAgent({
      id: agentId,
      slug: toNullString(''),
      title: toNullString(''),
      description: toNullString(''),
      tags: toNullString(''),
      avatar: toNullString(''),
      backgroundColor: toNullString(''),
      plugins: toNullString(JSON.stringify(newSession.config?.plugins || [])),
      chatConfig: toNullString(JSON.stringify(newSession.config?.chatConfig || {})),
      fewShots: toNullString(JSON.stringify(newSession.config?.fewShots || [])),
      model: toNullString(newSession.config?.model),
      params: toNullString(JSON.stringify(newSession.config?.params || {})),
      provider: toNullString(newSession.config?.provider),
      systemRole: toNullString(newSession.config?.systemRole),
      tts: toNullString(''),
      virtual: newSession.config?.virtual ? 1 : 0,
      openingMessage: toNullString(newSession.config?.openingMessage),
      openingQuestions: toNullString(JSON.stringify(newSession.config?.openingQuestions || [])),
      createdAt: now,
      updatedAt: now,
    });

    // Link agent to session
    await DB.LinkAgentToSession({
      agentId,
      sessionId,
    });

    console.log('[Session] Created session via direct DB', { sessionId });

    // Immediately load agent config for new session
    const dbAgent = await DB.GetAgentBySessionId(sessionId);
    const config = mapAgentConfigFromDB(dbAgent);
    const agentStore = useAgentStore.getState();
    agentStore.internal_dispatchAgentMap(sessionId, config, 'createSession');

    // Mark config as loaded in agentConfigInitMap
    agentStore.internal_updateAgentConfigInitMap(sessionId, true);

    await refreshSessions();

    // Whether to goto  to the new session after creation, the default is to switch to
    if (isSwitchSession) switchSession(sessionId);

    return sessionId;
  },

  duplicateSession: async (id) => {
    const { switchSession, refreshSessions } = get();
    const session = sessionSelectors.getSessionById(id)(get());

    if (!session) return;
    const title = getNullableString(session.title);

    const newTitle = t('duplicateSession.title', { ns: 'chat', title: title });

    const messageLoadingKey = 'duplicateSession.loading';

    message.loading({
      content: t('duplicateSession.loading', { ns: 'chat' }),
      duration: 0,
      key: messageLoadingKey,
    });

    try {
      const newSessionId = crypto.randomUUID();
      const newAgentId = crypto.randomUUID();
      const now = Date.now();

      // 1. Duplicate session
      await DB.DuplicateSession({
        id: newSessionId,
        title: toNullString(newTitle),
        createdAt: now,
        updatedAt: now,
        id2: id,
      });

      // 2. Duplicate agent
      await DB.DuplicateAgentForSession({
        id: newAgentId,
        createdAt: now,
        updatedAt: now,
        sessionId: id,
      });

      // 3. Link agent to session
      await DB.LinkDuplicatedAgentToSession({
        agentId: newAgentId,
        sessionId: newSessionId,
      });

      console.log('[Session] Duplicated session via direct DB', {
        sourceId: id,
        newId: newSessionId
      });

      await refreshSessions();
      message.destroy(messageLoadingKey);
      message.success(t('duplicateSession.success', { ns: 'chat' }));

      switchSession(newSessionId);
    } catch (error) {
      console.error('[duplicateSession] Error:', error);
      message.destroy(messageLoadingKey);
      message.error(t('copyFail', { ns: 'common' }));
    }
  },
  pinSession: async (id, pinned) => {
    await get().internal_updateSession(id, { pinned });
  },
  removeSession: async (sessionId) => {
    await DB.DeleteSession(sessionId);

    console.log('[Session] Deleted session via direct DB', { sessionId });

    await get().refreshSessions();

    // If the active session deleted, switch to the inbox session
    if (sessionId === get().activeId) {
      get().switchSession(INBOX_SESSION_ID);
    }
  },

  switchSession: (sessionId) => {
    if (get().activeId === sessionId) return;

    set({ activeId: sessionId }, false, n(`activeSession/${sessionId}`));
  },

  triggerSessionUpdate: async (id) => {
    await get().internal_updateSession(id, { updatedAt: new Date() });
  },

  updateSearchKeywords: (keywords) => {
    set(
      { isSearching: !!keywords, sessionSearchKeywords: keywords },
      false,
      n('updateSearchKeywords'),
    );
  },
  updateSessionGroupId: async (sessionId, groupId) => {
    const session = sessionSelectors.getSessionById(sessionId)(get());

    if (getNullableString(session?.type) === 'group') {
      // For group sessions (chat groups), use the chat group service
      // await chatGroupService.updateGroup(sessionId, {
      //   groupId: groupId === 'default' ? null : groupId,
      // });
      // await get().refreshSessions();
    } else {
      // For regular agent sessions, use the existing session service
      await get().internal_updateSession(sessionId, { group: groupId });
    }
  },

  updateSessionMeta: async (meta) => {
    const session = sessionSelectors.currentSession(get());
    if (!session) return;

    const { activeId, refreshSessions } = get();

    // Skip inbox session (cannot modify meta)
    if (activeId === INBOX_SESSION_ID) return;

    const abortController = get().signalSessionMeta as AbortController;
    if (abortController) abortController.abort(MESSAGE_CANCEL_FLAT);
    const controller = new AbortController();
    set({ signalSessionMeta: controller }, false, 'updateSessionMetaSignal');

    const userId = getUserId();
    const now = Date.now();

    await DB.UpdateSession({
      id: activeId,
      userId,
      title: meta.title ? toNullString(meta.title) : undefined,
      description: meta.description ? toNullString(meta.description) : undefined,
      avatar: meta.avatar ? toNullString(meta.avatar) : undefined,
      backgroundColor: meta.backgroundColor ? toNullString(meta.backgroundColor) : undefined,
      updatedAt: now,
    } as any);

    console.log('[Session] Updated session meta via direct DB', { id: activeId });

    await refreshSessions();
  },

  internal_fetchSessions: async (enabled, isLogin) => {
    if (!enabled) return;

    try {
      // Fetch sessions and session groups in parallel
      // Initial fetch: Limit 20, Offset 0
      const limit = 10;
      const offset = 0;

      const [dbSessions, dbSessionGroups] = await Promise.all([
        DB.ListSessions({ limit, offset }),
        DB.ListSessionGroups(),
      ]);

      // Reset pagination state
      set({ sessionsPage: 0, sessionsHasMore: dbSessions.length >= limit }, false, n('internal_fetchSessions/resetPage'));

      // Map database results to frontend types
      const sessions = dbSessions;
      const sessionGroups = dbSessionGroups.map((g: any) => ({
        id: g.id,
        name: g.name || '',
        sort: Number(g.sort) || 0,
        createdAt: new Date(g.createdAt),
        updatedAt: new Date(g.updatedAt),
      }));

      // Skip update if data hasn't changed
      if (
        get().isSessionsFirstFetchFinished &&
        isEqual(get().sessions, sessions) &&
        isEqual(get().sessionGroups, sessionGroups)
      ) {
        return;
      }

      get().internal_processSessions(
        sessions as any,
        sessionGroups,
        n('internal_fetchSessions/updateData') as any,
      );

      // Sync chat groups from group sessions to chat store
      const groupSessions = sessions.filter((session) => getNullableString(session.type) === 'group');
      if (groupSessions.length > 0) {
        const chatGroupStore = getChatGroupStoreState();
        const chatGroups = groupSessions.map((session) => ({
          accessedAt: session.updatedAt,
          clientId: null,
          config: {
            maxResponseInRow: 3,
            orchestratorModel: 'gpt-4',
            orchestratorProvider: 'openai',
            responseOrder: 'sequential' as const,
            responseSpeed: 'medium' as const,
            scene: DEFAULT_CHAT_GROUP_CHAT_CONFIG.scene,
            allowDM: DEFAULT_CHAT_GROUP_CHAT_CONFIG.allowDM,
            enableSupervisor: DEFAULT_CHAT_GROUP_CHAT_CONFIG.enableSupervisor,
            revealDM: DEFAULT_CHAT_GROUP_CHAT_CONFIG.revealDM,
          },
          createdAt: session.createdAt,
          description: getNullableString(session.description) || '',
          groupId: getNullableString(session.groupId) || null,
          id: session.id,
          pinned: Boolean(session.pinned),
          slug: null,
          title: getNullableString(session.title) || 'Untitled Group',
          updatedAt: session.updatedAt,
          userId: '',
        }));

        chatGroupStore.internal_updateGroupMaps(chatGroups);
      }

      set({ isSessionsFirstFetchFinished: true }, false, n('internal_fetchSessions/onSuccess'));
    } catch (error) {
      console.error('[internal_fetchSessions] Error fetching sessions:', error);
    }
  },

  internal_searchSessions: async (keyword) => {
    if (!keyword) return;

    try {
      const dbResults = await DB.SearchSessionsByKeyword({
        column1: toNullString(keyword),
        column2: toNullString(keyword),
      });
      const results = dbResults;

      console.log('[Session] Searched sessions via direct DB', { keyword, count: results.length });
    } catch (error) {
      console.error('[internal_searchSessions] Error searching sessions:', error);
    }
  },

  /* eslint-disable sort-keys-fix/sort-keys-fix */
  internal_dispatchSessions: (payload) => {
    const nextSessions = sessionsReducer(get().sessions, payload);
    get().internal_processSessions(nextSessions, get().sessionGroups);
  },
  internal_updateSession: async (id, data) => {
    // Convert boolean pinned to number for internal store update (Session expects number)
    const value: any = { ...data };
    if (data.pinned !== undefined) {
      value.pinned = data.pinned ? 1 : 0;
    }
    get().internal_dispatchSessions({ type: 'updateSession', id, value });

    const userId = getUserId();
    const now = Date.now();

    // Map UpdateSessionParams to DB params
    const meta = data as any; // Cast to access meta properties

    await DB.UpdateSession({
      id,
      userId,
      title: meta.title ? toNullString(meta.title) : undefined,
      description: meta.description ? toNullString(meta.description) : undefined,
      avatar: meta.avatar ? toNullString(meta.avatar) : undefined,
      backgroundColor: meta.backgroundColor ? toNullString(meta.backgroundColor) : undefined,
      groupId: data.group !== undefined ? toNullString(data.group === 'default' ? '' : data.group) : undefined,
      pinned: data.pinned !== undefined ? toNullInt(boolToInt(data.pinned)) : undefined,
      updatedAt: data.updatedAt ? data.updatedAt.getTime() : now,
    } as any);

    console.log('[Session] Updated session via direct DB', { id });

    await get().refreshSessions();
  },
  internal_processSessions: (sessions, sessionGroups) => {
    const customGroups = sessionGroups.map((item) => ({
      ...item,
      children: sessions.filter((i) => getNullableString(i.groupId) === item.id && !i.pinned),
    }));

    const defaultGroup = sessions.filter(
      (item) => (!getNullableString(item.groupId) || getNullableString(item.groupId) === 'default') && !item.pinned,
    );
    const pinnedGroup = sessions.filter((item) => item.pinned);

    set(
      {
        customSessionGroups: customGroups,
        defaultSessions: defaultGroup,
        pinnedSessions: pinnedGroup,
        sessionGroups,
        sessions,
      },
      false,
      n('processSessions'),
    );
  },
  refreshSessions: async () => {
    try {
      // When refreshing, we want to keep the current amount of data loaded
      // So limit should be (page + 1) * 20
      const { sessionsPage } = get();
      const currentLimit = (sessionsPage + 1) * 20;

      // Fetch sessions and session groups in parallel
      // We use ListSessions with larger limit to cover everything currently loaded
      const [dbSessions, dbSessionGroups] = await Promise.all([
        DB.ListSessions({ limit: currentLimit, offset: 0 }),
        DB.ListSessionGroups(),
      ]);

      // Map database results to frontend types
      const sessions = dbSessions;
      const sessionGroups = dbSessionGroups.map((g: any) => ({
        id: g.id,
        name: g.name || '',
        sort: Number(g.sort) || 0,
        createdAt: new Date(g.createdAt),
        updatedAt: new Date(g.updatedAt),
      }));

      console.log('[Session] Refreshed sessions via direct DB', { count: sessions.length });

      get().internal_processSessions(sessions as any, sessionGroups);

      // Sync chat groups
      const groupSessions = sessions.filter((session) => getNullableString(session.type) === 'group');
      if (groupSessions.length > 0) {
        const chatGroupStore = getChatGroupStoreState();
        const chatGroups = groupSessions.map((session) => ({
          accessedAt: session.updatedAt,
          clientId: null,
          config: {
            maxResponseInRow: 3,
            orchestratorModel: 'gpt-4',
            orchestratorProvider: 'openai',
            responseOrder: 'sequential' as const,
            responseSpeed: 'medium' as const,
            scene: DEFAULT_CHAT_GROUP_CHAT_CONFIG.scene,
            allowDM: DEFAULT_CHAT_GROUP_CHAT_CONFIG.allowDM,
            enableSupervisor: DEFAULT_CHAT_GROUP_CHAT_CONFIG.enableSupervisor,
            revealDM: DEFAULT_CHAT_GROUP_CHAT_CONFIG.revealDM,
          },
          createdAt: session.createdAt,
          description: getNullableString(session.description) || '',
          groupId: getNullableString(session.groupId) || null,
          id: session.id,
          pinned: Boolean(session.pinned),
          slug: null,
          title: getNullableString(session.title) || 'Untitled Group',
          updatedAt: session.updatedAt,
          userId: '',
        }));

        chatGroupStore.internal_updateGroupMaps(chatGroups);
      }
    } catch (error) {
      console.error('[refreshSessions] Error refreshing sessions:', error);
    }
  },
});
