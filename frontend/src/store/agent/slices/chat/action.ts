import isEqual from 'fast-deep-equal';
import { produce } from 'immer';
import { useEffect } from 'react';
import type { PartialDeep } from 'type-fest';
import { StateCreator } from 'zustand/vanilla';

import { MESSAGE_CANCEL_FLAT } from '@/const/message';
import { INBOX_SESSION_ID } from '@/const/session';
// import { agentService } from '@/services/agent';
import { AgentState } from '@/store/agent/slices/chat/initialState';
import { getSessionStoreState, useSessionStore } from '@/store/session';
import { LobeAgentChatConfig, LobeAgentConfig } from '@/types/agent';
import { KnowledgeItem } from '@/types/knowledgeBase';
import { merge } from '@/utils/merge';
import { DB } from '@/types/database';
import { getUserId } from '@/store/session/helpers';
import { mapAgentConfigFromDB } from '@/store/session/helpers';
import { toNullString, toNullJSON, toNullInt } from '@/types/database';

import type { AgentStore } from '../../store';
import { agentSelectors } from './selectors';

/**
 * 助手接口
 */
export interface AgentChatAction {
  addFilesToAgent: (fileIds: string[], boolean?: boolean) => Promise<void>;
  addKnowledgeBaseToAgent: (knowledgeBaseId: string) => Promise<void>;
  internal_createAbortController: (key: keyof AgentState) => AbortController;

  internal_dispatchAgentMap: (
    id: string,
    config: PartialDeep<LobeAgentConfig>,
    actions?: string,
  ) => void;
  internal_refreshAgentConfig: (id: string) => Promise<void>;
  internal_refreshAgentKnowledge: () => Promise<void>;
  internal_updateAgentConfig: (
    id: string,
    data: PartialDeep<LobeAgentConfig>,
    signal?: AbortSignal,
  ) => Promise<void>;
  internal_updateAgentConfigInitMap: (id: string, loaded: boolean) => void;
  internal_updateActiveId: (id: string) => void;
  removeFileFromAgent: (fileId: string) => Promise<void>;
  removeKnowledgeBaseFromAgent: (knowledgeBaseId: string) => Promise<void>;

  removePlugin: (id: string) => void;
  toggleFile: (id: string, open?: boolean) => Promise<void>;
  toggleKnowledgeBase: (id: string, open?: boolean) => Promise<void>;

  togglePlugin: (id: string, open?: boolean) => Promise<void>;
  updateAgentChatConfig: (config: Partial<LobeAgentChatConfig>) => Promise<void>;
  updateAgentConfig: (config: PartialDeep<LobeAgentConfig>) => Promise<void>;
  internal_fetchAgentConfig: (isLogin: boolean | undefined, sessionId: string) => Promise<void>;
  internal_fetchFilesAndKnowledgeBases: () => Promise<void>;
  useInitInboxAgentStore: (
    isLogin: boolean | undefined,
    defaultAgentConfig?: PartialDeep<LobeAgentConfig>,
  ) => void;
  useLoadAllAgentConfigs: (isLogin: boolean) => void;
}

export const createChatSlice: StateCreator<
  AgentStore,
  [['zustand/devtools', never]],
  [],
  AgentChatAction
> = (set, get) => ({
  addFilesToAgent: async (fileIds, enabled) => {
    const { activeAgentId, internal_refreshAgentConfig, internal_refreshAgentKnowledge } = get();
    if (!activeAgentId) return;
    if (fileIds.length === 0) return;

    // await agentService.createAgentFiles(activeAgentId, fileIds, enabled);
    await internal_refreshAgentConfig(get().activeId);
    await internal_refreshAgentKnowledge();
  },
  addKnowledgeBaseToAgent: async (knowledgeBaseId) => {
    const { activeAgentId, internal_refreshAgentConfig, internal_refreshAgentKnowledge } = get();
    if (!activeAgentId) return;

    // await agentService.createAgentKnowledgeBase(activeAgentId, knowledgeBaseId, true);
    await internal_refreshAgentConfig(get().activeId);
    await internal_refreshAgentKnowledge();
  },
  removeFileFromAgent: async (fileId) => {
    const { activeAgentId, internal_refreshAgentConfig, internal_refreshAgentKnowledge } = get();
    if (!activeAgentId) return;

    // await agentService.deleteAgentFile(activeAgentId, fileId);
    await internal_refreshAgentConfig(get().activeId);
    await internal_refreshAgentKnowledge();
  },
  removeKnowledgeBaseFromAgent: async (knowledgeBaseId) => {
    const { activeAgentId, internal_refreshAgentConfig, internal_refreshAgentKnowledge } = get();
    if (!activeAgentId) return;

    // await agentService.deleteAgentKnowledgeBase(activeAgentId, knowledgeBaseId);
    await internal_refreshAgentConfig(get().activeId);
    await internal_refreshAgentKnowledge();
  },

  removePlugin: async (id) => {
    await get().togglePlugin(id, false);
  },
  toggleFile: async (id, open) => {
    const { activeAgentId, internal_refreshAgentConfig } = get();
    if (!activeAgentId) return;

    // await agentService.toggleFile(activeAgentId, id, open);

    await internal_refreshAgentConfig(get().activeId);
  },
  toggleKnowledgeBase: async (id, open) => {
    const { activeAgentId, internal_refreshAgentConfig } = get();
    if (!activeAgentId) return;

    // await agentService.toggleKnowledgeBase(activeAgentId, id, open);

    await internal_refreshAgentConfig(get().activeId);
  },
  togglePlugin: async (id, open) => {
    const originConfig = agentSelectors.currentAgentConfig(get());

    const config = produce(originConfig, (draft) => {
      draft.plugins = produce(draft.plugins || [], (plugins) => {
        const index = plugins.indexOf(id);
        const shouldOpen = open !== undefined ? open : index === -1;

        if (shouldOpen) {
          // 如果 open 为 true 或者 id 不存在于 plugins 中，则添加它
          if (index === -1) {
            plugins.push(id);
          }
        } else {
          // 如果 open 为 false 或者 id 存在于 plugins 中，则移除它
          if (index !== -1) {
            plugins.splice(index, 1);
          }
        }
      });
    });

    await get().updateAgentConfig(config);
  },
  updateAgentChatConfig: async (config) => {
    const { activeId } = get();

    if (!activeId) return;

    await get().updateAgentConfig({ chatConfig: config });
  },
  updateAgentConfig: async (config) => {
    const { activeId } = get();

    if (!activeId) return;

    const controller = get().internal_createAbortController('updateAgentConfigSignal');

    await get().internal_updateAgentConfig(activeId, config, controller.signal);
  },
  internal_fetchAgentConfig: async (isLogin, sessionId) => {
    if (isLogin !== true || sessionId.startsWith('cg_')) return;

    try {
      const userId = getUserId();
      const dbAgent = await DB.GetAgentBySessionId({ sessionId, userId });
      const data = mapAgentConfigFromDB(dbAgent);

      console.log('[Agent] Fetched agent config via direct DB', { sessionId });

      get().internal_dispatchAgentMap(sessionId, data, 'fetch');

      set(
        {
          activeAgentId: data?.id || undefined,
          agentConfigInitMap: { ...get().agentConfigInitMap, [sessionId]: true },
        },
        false,
        'fetchAgentConfig',
      );
    } catch (error) {
      console.error('[internal_fetchAgentConfig] Error fetching agent config:', error);
    }
  },

  internal_fetchFilesAndKnowledgeBases: async () => {
    const activeAgentId = get().activeAgentId;
    if (!activeAgentId) return;

    try {
      // TODO: Implement when agentService is available
      // const data = await agentService.getFilesAndKnowledgeBases(activeAgentId);
    } catch (error) {
      console.error('[internal_fetchFilesAndKnowledgeBases] Error:', error);
    }
  },

  useInitInboxAgentStore: (isLogin, defaultAgentConfig) => {
    useEffect(() => {
      if (isLogin !== true) return;
      if (get().isInboxAgentConfigInit) return; // Only fetch once

      const initInboxAgent = async () => {
        try {
          const dbAgent = await DB.GetAgentBySessionId(INBOX_SESSION_ID);
          const data = mapAgentConfigFromDB(dbAgent);

          set(
            {
              defaultAgentConfig: merge(get().defaultAgentConfig, defaultAgentConfig),
              isInboxAgentConfigInit: true,
              agentConfigInitMap: { ...get().agentConfigInitMap, [INBOX_SESSION_ID]: true },
            },
            false,
            'initDefaultAgent',
          );

          if (data) {
            get().internal_dispatchAgentMap(INBOX_SESSION_ID, data, 'initInbox');
          }
        } catch (error) {
          console.error('[useInitInboxAgentStore] Error loading inbox config:', error);
          // Inbox should always exist (created by backend at startup)
          // If this fails, it indicates a serious issue
        }
      };

      initInboxAgent();
    }, [isLogin, defaultAgentConfig]);
  },

  useLoadAllAgentConfigs: (isLogin) => {
    useEffect(() => {
      if (!isLogin) return;
      if (get().isAllAgentConfigsLoaded) return; // Only fetch once

      const loadAllConfigs = async () => {
        try {
          const sessionStore = getSessionStoreState();
          const sessions = sessionStore.sessions;

          // Batch load all agent configs
          // Skip inbox session as it's already handled by useInitInboxAgentStore
          // This prevents race condition and redundant API calls
          const configPromises = sessions
            .filter((s) => s.type === 'agent')
            .filter((s) => s.id !== INBOX_SESSION_ID)
            .map((session) =>
              sessionService
                .getSessionConfig(session.id)
                .then((config) => ({ sessionId: session.id, config }))
                .catch(() => null),
            );

          const results = await Promise.all(configPromises);

          // Populate agentMap and agentConfigInitMap
          const agentConfigInitMap = { ...get().agentConfigInitMap };

          results.forEach((result) => {
            if (result) {
              get().internal_dispatchAgentMap(result.sessionId, result.config, 'batchLoad');
              agentConfigInitMap[result.sessionId] = true;
            }
          });

          set({ agentConfigInitMap }, false, 'batchLoadConfigInitMap');

          const count = results.filter((r) => r !== null).length;
          console.info(`[AgentStore] Loaded ${count} agent configs`);
          set({ isAllAgentConfigsLoaded: true }, false, 'allAgentConfigsLoaded');
        } catch (error) {
          console.error('[AgentStore] Failed to batch load agent configs:', error);
        }
      };

      loadAllConfigs();
    }, [isLogin]);
  },
  /* eslint-disable sort-keys-fix/sort-keys-fix */

  internal_dispatchAgentMap: (id, config, actions) => {
    const agentMap = produce(get().agentMap, (draft) => {
      if (!draft[id]) {
        draft[id] = config;
      } else {
        draft[id] = merge(draft[id], config);
      }
    });

    if (isEqual(get().agentMap, agentMap)) return;

    set({ agentMap }, false, 'dispatchAgent' + (actions ? `/${actions}` : ''));
  },

  internal_updateAgentConfigInitMap: (id, loaded) => {
    set(
      { agentConfigInitMap: { ...get().agentConfigInitMap, [id]: loaded } },
      false,
      'updateAgentConfigInitMap',
    );
  },

  internal_updateActiveId: (id) => {
    set({ activeId: id }, false, 'updateActiveId');
  },

  internal_updateAgentConfig: async (id, data, signal) => {
    const prevModel = agentSelectors.currentAgentModel(get());
    // optimistic update at frontend
    get().internal_dispatchAgentMap(id, data, 'optimistic_updateAgentConfig');

    const userId = getUserId();
    const now = Date.now();

    await DB.UpdateAgent({
      sessionId: id,
      userId,
      model: data.model ? toNullString(data.model) : undefined,
      systemRole: data.systemRole ? toNullString(data.systemRole) : undefined,
      plugins: data.plugins ? toNullJSON(data.plugins) : undefined,
      chatConfig: data.chatConfig ? toNullJSON(data.chatConfig) : undefined,
      params: data.params ? toNullJSON(data.params) : undefined,
      openingMessage: data.openingMessage ? toNullString(data.openingMessage) : undefined,
      openingQuestions: data.openingQuestions ? toNullJSON(data.openingQuestions) : undefined,
      fewShots: data.fewShots ? toNullJSON(data.fewShots) : undefined,
      virtual: data.virtual !== undefined ? toNullInt(data.virtual ? 1 : 0) : undefined,
      provider: data.provider ? toNullString(data.provider) : undefined,
      updatedAt: now,
    } as any);

    console.log('[Agent] Updated agent config via direct DB', { sessionId: id });

    await get().internal_refreshAgentConfig(id);

    // refresh sessions to update the agent config if the model has changed
    if (prevModel !== data.model) await useSessionStore.getState().refreshSessions();
  },

  internal_refreshAgentConfig: async (id) => {
    try {
      const userId = getUserId();
      const dbAgent = await DB.GetAgentBySessionId({ sessionId: id, userId });
      const data = mapAgentConfigFromDB(dbAgent);

      console.log('[Agent] Refreshed agent config via direct DB', { sessionId: id });

      get().internal_dispatchAgentMap(id, data, 'refresh');
    } catch (error) {
      console.error('[internal_refreshAgentConfig] Error:', error);
    }
  },

  internal_refreshAgentKnowledge: async () => {
    try {
      // TODO: Implement when agentService is available
      // const data = await agentService.getFilesAndKnowledgeBases(get().activeAgentId);
      console.debug('[internal_refreshAgentKnowledge] Refreshed');
    } catch (error) {
      console.error('[internal_refreshAgentKnowledge] Error:', error);
    }
  },
  internal_createAbortController: (key) => {
    const abortController = get()[key] as AbortController;
    if (abortController) abortController.abort(MESSAGE_CANCEL_FLAT);
    const controller = new AbortController();
    set({ [key]: controller }, false, 'internal_createAbortController');

    return controller;
  },
});
