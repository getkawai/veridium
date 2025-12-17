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

import { merge } from '@/utils/merge';
import { DB } from '@/types/database';
import { getUserId } from '@/store/session/helpers';
import { mapAgentConfigFromDB } from '@/store/session/helpers';
import { toNullString, toNullJSON, toNullInt, intToBool, getNullableString } from '@/types/database';
import { useState } from 'react';

import type { AgentStore } from '../../store';
import { agentSelectors } from './selectors';
import { KnowledgeItem, KnowledgeType } from '@/types';
import { createServiceLogger } from '@/utils/logger';

const logger = createServiceLogger('AgentChat', 'AgentChatAction', 'store/agent/slices/chat/action.ts');

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
  useFetchFilesAndKnowledgeBases: (agentId?: string, version?: number) => { data: KnowledgeItem[]; error: Error | null; isLoading: boolean };
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
    const agentConfig = agentSelectors.currentAgentConfig(get());
    const agentId = agentConfig?.id;
    if (!agentId || fileIds.length === 0) return;

    const now = Date.now();
    for (const fileId of fileIds) {
      try {
        await DB.LinkAgentToFile({
          agentId,
          fileId,
          enabled: enabled ? 1 : 0,
          createdAt: now,
          updatedAt: now,
        });
      } catch {
        await DB.ToggleAgentFile({
          agentId,
          fileId,
          enabled: enabled ? 1 : 0,
        });
      }
    }

    await get().internal_refreshAgentKnowledge();
  },
  addKnowledgeBaseToAgent: async (knowledgeBaseId) => {
    const agentConfig = agentSelectors.currentAgentConfig(get());
    const agentId = agentConfig?.id;
    if (!agentId) return;

    const now = Date.now();
    try {
      await DB.LinkAgentToKnowledgeBase({
        agentId,
        knowledgeBaseId,
        enabled: 1,
        createdAt: now,
        updatedAt: now,
      });
    } catch {
      await DB.ToggleAgentKnowledgeBase({
        agentId,
        knowledgeBaseId,
        enabled: 1,
      });
    }

    await get().internal_refreshAgentKnowledge();
  },
  removeFileFromAgent: async (fileId) => {
    const agentConfig = agentSelectors.currentAgentConfig(get());
    const agentId = agentConfig?.id;
    if (!agentId) return;

    await DB.UnlinkAgentFromFile({ agentId, fileId });
    await get().internal_refreshAgentKnowledge();
  },
  removeKnowledgeBaseFromAgent: async (knowledgeBaseId) => {
    const agentConfig = agentSelectors.currentAgentConfig(get());
    const agentId = agentConfig?.id;
    if (!agentId) return;

    await DB.UnlinkAgentFromKnowledgeBase({ agentId, knowledgeBaseId });
    await get().internal_refreshAgentKnowledge();
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
      const dbAgent = await DB.GetAgentBySessionId(sessionId);
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

  useFetchFilesAndKnowledgeBases: (agentId?: string, version?: number) => {
    const [data, setData] = useState<KnowledgeItem[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    const activeAgentId = get().activeAgentId;
    const targetId = agentId || activeAgentId;
    // We need to listen to agentConfigInitMap to trigger re-fetch when config changes or reloadsy on mount and targetId change

    // Also we might need to listen to global refresh events if we want this to update after assignment
    // The store has `internal_refreshAgentKnowledge` but that's an action, not a state.
    // However, `internal_refreshAgentKnowledge` just logs for now. 
    // If we want reactive updates, we might need a version signal in the store. 
    // For now, let's implement the fetching logic first.

    useEffect(() => {
      logger.debug('[useFetchFilesAndKnowledgeBases] Hook triggered', { agentId, activeAgentId, targetId });

      if (!targetId) {
        logger.debug('[useFetchFilesAndKnowledgeBases] No targetId, clearing data');
        setData([]);
        return;
      }

      let isMounted = true;
      setIsLoading(true);
      setError(null);

      const fetchData = async () => {
        try {
          logger.info('[useFetchFilesAndKnowledgeBases] Fetching data for', { targetId });
          const [files, knowledgeBases] = await Promise.all([
            DB.GetAgentFilesWithEnabled(targetId),
            DB.GetAgentKnowledgeBases(targetId)
          ]);

          logger.info('[useFetchFilesAndKnowledgeBases] Raw response', { filesLength: files.length, kbLength: knowledgeBases.length, files, knowledgeBases });

          if (!isMounted) return;

          const items: KnowledgeItem[] = [
            ...files.map(f => ({
              id: f.id,
              name: f.name,
              type: KnowledgeType.File,
              enabled: intToBool(f.enabled),
              fileType: f.fileType,
              // Add other fields if needed
            })),
            ...knowledgeBases.map(kb => ({
              id: kb.id,
              name: kb.name,
              type: KnowledgeType.KnowledgeBase,
              enabled: intToBool(kb.enabled),
              description: getNullableString(kb.description),
              avatar: getNullableString(kb.avatar),
            }))
          ];

          setData(items);
        } catch (err) {
          if (isMounted) {
            logger.error('[useFetchFilesAndKnowledgeBases] Error:', err);
            setError(err as Error);
          }
        } finally {
          if (isMounted) {
            setIsLoading(false);
          }
        }
      };

      fetchData();

      return () => {
        isMounted = false;
      };
      return () => {
        isMounted = false;
      };
    }, [targetId, version]);

    // Expose a refetch function effectively by just returning the data which changes
    // But if we want to manually trigger, we'd need to trust the store actions 
    // `internal_refreshAgentKnowledge` logic, which currently does nothing.
    // We should probably subscribe to a refresh signal if we want perfect sync, 
    // but the provided "SWR" style interface implies just data returning.

    return { data, error, isLoading };
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
            .filter((s) => getNullableString(s.type) === 'agent')
            .filter((s) => s.id !== INBOX_SESSION_ID)
            .map((session) =>
              DB.GetAgentBySessionId(session.id)
                .then((dbAgent) => mapAgentConfigFromDB(dbAgent))
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
      const dbAgent = await DB.GetAgentBySessionId(id);
      const data = mapAgentConfigFromDB(dbAgent);

      console.log('[Agent] Refreshed agent config via direct DB', { sessionId: id });

      get().internal_dispatchAgentMap(id, data, 'refresh');
    } catch (error) {
      console.error('[internal_refreshAgentConfig] Error:', error);
    }
  },

  internal_refreshAgentKnowledge: async () => {
    try {
      set({ knowledgeRefreshVersion: get().knowledgeRefreshVersion + 1 }, false, 'internal_refreshAgentKnowledge');
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
