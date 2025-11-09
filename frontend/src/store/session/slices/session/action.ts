import isEqual from 'fast-deep-equal';
import { t } from 'i18next';
import type { PartialDeep } from 'type-fest';
import { StateCreator } from 'zustand/vanilla';

import { message } from '@/components/AntdStaticMethods';
import { MESSAGE_CANCEL_FLAT } from '@/const/message';
import { DEFAULT_AGENT_LOBE_SESSION, INBOX_SESSION_ID } from '@/const/session';
import { DEFAULT_CHAT_GROUP_CHAT_CONFIG } from '@/const/settings';
import { chatGroupService } from '@/services/chatGroup';
import { sessionService } from '@/services/session';
import { useAgentStore } from '@/store/agent';
import { getChatGroupStoreState } from '@/store/chatGroup';
import { useUserStore } from '@/store/user';

import type { SessionStore } from '../../store';
import { settingsSelectors } from '@/store/user/selectors';
import { MetaData } from '@/types/meta';
import {
  ChatSessionList,
  LobeAgentSession,
  LobeSessionGroups,
  LobeSessionType,
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

  internal_fetchSessions: (...) => Promise<void>;
  internal_searchSessions: (...) => Promise<void>;

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
    await sessionService.removeAllSessions();
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

    const id = await sessionService.createSession(LobeSessionType.Agent, newSession);

    // Immediately load agent config for new session
    const config = await sessionService.getSessionConfig(id);
    const agentStore = useAgentStore.getState();
    agentStore.internal_dispatchAgentMap(id, config, 'createSession');
    
    // Mark config as loaded in agentConfigInitMap
    agentStore.internal_updateAgentConfigInitMap(id, true);

    await refreshSessions();

    // Whether to goto  to the new session after creation, the default is to switch to
    if (isSwitchSession) switchSession(id);

    return id;
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

    const newId = await sessionService.cloneSession(id, newTitle);

    // duplicate Session Error
    if (!newId) {
      message.destroy(messageLoadingKey);
      message.error(t('copyFail', { ns: 'common' }));
      return;
    }

    await refreshSessions();
    message.destroy(messageLoadingKey);
    message.success(t('duplicateSession.success', { ns: 'chat' }));

    switchSession(newId);
  },
  pinSession: async (id, pinned) => {
    await get().internal_updateSession(id, { pinned });
  },
  removeSession: async (sessionId) => {
    await sessionService.removeSession(sessionId);
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
  updateSessionGroupId: async (sessionId, group) => {
    const session = sessionSelectors.getSessionById(sessionId)(get());

    if (session?.type === 'group') {
      // For group sessions (chat groups), use the chat group service
      await chatGroupService.updateGroup(sessionId, {
        groupId: group === 'default' ? null : group,
      });
      await get().refreshSessions();
    } else {
      // For regular agent sessions, use the existing session service
      await get().internal_updateSession(sessionId, { group });
    }
  },

  updateSessionMeta: async (meta) => {
    const session = sessionSelectors.currentSession(get());
    if (!session) return;

    const { activeId, refreshSessions } = get();

    const abortController = get().signalSessionMeta as AbortController;
    if (abortController) abortController.abort(MESSAGE_CANCEL_FLAT);
    const controller = new AbortController();
    set({ signalSessionMeta: controller }, false, 'updateSessionMetaSignal');

    await sessionService.updateSessionMeta(activeId, meta, controller.signal);
    await refreshSessions();
  },

  internal_fetchSessions: (enabled, isLogin) => {
    useEffect(() => {
      if (!enabled) return;

      const fetchSessions = async () => {
        try {
          const data = await sessionService.getGroupedSessions();

          // Skip update if data hasn't changed
          if (
            get().isSessionsFirstFetchFinished &&
            isEqual(get().sessions, data.sessions) &&
            isEqual(get().sessionGroups, data.sessionGroups)
          ) {
            return;
          }

          get().internal_processSessions(
            data.sessions,
            data.sessionGroups,
            n('useFetchSessions/updateData') as any,
          );

          // Sync chat groups from group sessions to chat store
          const groupSessions = data.sessions.filter((session) => session.type === 'group');
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
              },
              createdAt: session.createdAt,
              description: session.meta?.description || '',
              groupId: session.group || null,
              id: session.id,
              pinned: session.pinned || false,
              slug: null,
              title: session.meta?.title || 'Untitled Group',
              updatedAt: session.updatedAt,
              userId: '',
            }));

            chatGroupStore.internal_updateGroupMaps(chatGroups);
          }

          set({ isSessionsFirstFetchFinished: true }, false, n('useFetchSessions/onSuccess'));
        } catch (error) {
          console.error('[useFetchSessions] Error fetching sessions:', error);
        }
      };

      fetchSessions();
    }, [enabled, isLogin]);
  },

  internal_searchSessions: (keyword) => {
    useEffect(() => {
      const searchSessions = async () => {
        if (!keyword) return;

        try {
          const results = await sessionService.searchSessions(keyword);
          console.debug('[useSearchSessions] Search results:', results.length);
        } catch (error) {
          console.error('[useSearchSessions] Error searching sessions:', error);
        }
      };

      searchSessions();
    }, [keyword]);
  },

  /* eslint-disable sort-keys-fix/sort-keys-fix */
  internal_dispatchSessions: (payload) => {
    const nextSessions = sessionsReducer(get().sessions, payload);
    get().internal_processSessions(nextSessions, get().sessionGroups);
  },
  internal_updateSession: async (id, data) => {
    get().internal_dispatchSessions({ type: 'updateSession', id, value: data });

    await sessionService.updateSession(id, data);
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
      const data = await sessionService.getGroupedSessions();
      get().internal_processSessions(data.sessions, data.sessionGroups);

      // Sync chat groups
      const groupSessions = data.sessions.filter((session) => session.type === 'group');
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
          },
          createdAt: session.createdAt,
          description: session.meta?.description || '',
          groupId: session.group || null,
          id: session.id,
          pinned: session.pinned || false,
          slug: null,
          title: session.meta?.title || 'Untitled Group',
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
