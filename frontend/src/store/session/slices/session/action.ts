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
import { DB, toNullString, toNullInt, boolToInt } from '@/types/database';
import { getUserId, mapSessionFromDB, mapAgentConfigFromDB } from '../../helpers';

import type { SessionStore } from '../../store';
import { settingsSelectors } from '@/store/user/selectors';
import { MetaData } from '@/types/meta';
import {
  LobeAgentSession,
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
    session?: PartialDeep<LobeAgentSession>,
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
    const userId = getUserId();
    const sessions = await DB.ListAllSessions(userId);

    // Delete all sessions (except inbox)
    const sessionsToDelete = sessions.filter(s => s.slug !== 'inbox');
    await Promise.all(
      sessionsToDelete.map((session) => DB.DeleteSession({ id: session.id, userId }))
    );

    console.log('[Session] Cleared all sessions via direct DB', { count: sessionsToDelete.length });

    await get().refreshSessions();
  },

  createSession: async (agent, isSwitchSession = true) => {
    const { switchSession, refreshSessions } = get();

    // merge the defaultAgent in settings
    const defaultAgent = merge(
      DEFAULT_AGENT_LOBE_SESSION,
      settingsSelectors.defaultAgent(useUserStore.getState()),
    );

    const newSession: LobeAgentSession = merge(defaultAgent, agent);

    const userId = getUserId();
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
      userId,
      groupId: toNullString(newSession.group === 'default' ? '' : newSession.group),
      clientId: toNullString(''),
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
      clientId: toNullString(''),
      userId,
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
      userId,
    });

    console.log('[Session] Created session via direct DB', { sessionId });

    // Immediately load agent config for new session
    const dbAgent = await DB.GetAgentBySessionId({ sessionId, userId });
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
    const title = sessionMetaSelectors.getTitle(session.meta);

    const newTitle = t('duplicateSession.title', { ns: 'chat', title: title });

    const messageLoadingKey = 'duplicateSession.loading';

    message.loading({
      content: t('duplicateSession.loading', { ns: 'chat' }),
      duration: 0,
      key: messageLoadingKey,
    });

    try {
      const userId = getUserId();
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
        userId,
      });

      // 2. Duplicate agent
      await DB.DuplicateAgentForSession({
        id: newAgentId,
        createdAt: now,
        updatedAt: now,
        sessionId: id,
        userId,
      });

      // 3. Link agent to session
      await DB.LinkDuplicatedAgentToSession({
        agentId: newAgentId,
        sessionId: newSessionId,
        userId,
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
    const userId = getUserId();
    await DB.DeleteSession({ id: sessionId, userId });

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

    if (session?.type === 'group') {
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
      const userId = getUserId();

      // Fetch sessions and session groups in parallel
      const [dbSessions, dbSessionGroups] = await Promise.all([
        DB.ListAllSessions(userId),
        DB.ListSessionGroups(userId),
      ]);

      // Map database results to frontend types
      const sessions = dbSessions.map(mapSessionFromDB);
      const sessionGroups = dbSessionGroups.map((g: any) => ({
        id: g.id,
        name: g.name || '',
        sort: Number(g.sort) || 0,
        createdAt: new Date(g.createdAt),
        updatedAt: new Date(g.updatedAt),
      }));

      console.log('[Session] Fetched sessions via direct DB', { count: sessions.length });

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
      const groupSessions = sessions.filter((session) => session.type === 'group');
      if (groupSessions.length > 0) {
        const chatGroupStore = getChatGroupStoreState();
        const chatGroups = groupSessions.map((session) => ({
          accessedAt: session.updatedAt.getTime(),
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
          createdAt: session.createdAt.getTime(),
          description: session.meta?.description || '',
          groupId: session.group || null,
          id: session.id,
          pinned: session.pinned || false,
          slug: null,
          title: session.meta?.title || 'Untitled Group',
          updatedAt: session.updatedAt.getTime(),
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
      const userId = getUserId();
      const dbResults = await DB.SearchSessionsByKeyword({
        userId,
        column2: toNullString(keyword),
        column3: toNullString(keyword),
      });
      const results = dbResults.map(mapSessionFromDB);

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
    get().internal_dispatchSessions({ type: 'updateSession', id, value: data });

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
      children: sessions.filter((i) => i.group === item.id && !i.pinned),
    }));

    const defaultGroup = sessions.filter(
      (item) => (!item.group || item.group === 'default') && !item.pinned,
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
      const userId = getUserId();

      // Fetch sessions and session groups in parallel
      const [dbSessions, dbSessionGroups] = await Promise.all([
        DB.ListAllSessions(userId),
        DB.ListSessionGroups(userId),
      ]);

      // Map database results to frontend types
      const sessions = dbSessions.map(mapSessionFromDB);
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
      const groupSessions = sessions.filter((session) => session.type === 'group');
      if (groupSessions.length > 0) {
        const chatGroupStore = getChatGroupStoreState();
        const chatGroups = groupSessions.map((session) => ({
          accessedAt: session.updatedAt.getTime(),
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
          createdAt: session.createdAt.getTime(),
          description: session.meta?.description || '',
          groupId: session.group || null,
          id: session.id,
          pinned: session.pinned || false,
          slug: null,
          title: session.meta?.title || 'Untitled Group',
          updatedAt: session.updatedAt.getTime(),
          userId: '',
        }));

        chatGroupStore.internal_updateGroupMaps(chatGroups);
      }
    } catch (error) {
      console.error('[refreshSessions] Error refreshing sessions:', error);
    }
  },
});
