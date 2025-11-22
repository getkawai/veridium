import isEqual from 'fast-deep-equal';
import { produce } from 'immer';
import { StateCreator } from 'zustand/vanilla';

import { INBOX_SESSION_ID } from '@/const/session';
import { DEFAULT_CHAT_GROUP_CHAT_CONFIG } from '@/const/settings';
import type { ChatGroupItem, NewChatGroup } from '@/types/chatGroup';
import type { ChatStoreState } from '@/store/chat/initialState';
import { useChatStore } from '@/store/chat/store';
import { getSessionStoreState } from '@/store/session';
import { setNamespace } from '@/utils/storeDebug';
import { DB, toNullString, toNullJSON, toNullInt, boolToInt } from '@/types/database';
import { getUserId, mapChatGroupFromDB } from './helpers';

import {
  ChatGroupAction,
  ChatGroupState,
  ChatGroupStore,
  initialChatGroupState,
} from './initialState';
import { ChatGroupReducer, chatGroupReducers } from './reducers';
import { chatGroupSelectors } from './selectors';

const n = setNamespace('chatGroup');

const syncChatStoreGroupMap = (groupMap: Record<string, ChatGroupItem>) => {
  useChatStore.setState(
    produce((state: ChatStoreState) => {
      state.groupMaps = groupMap;
      state.groupsInit = true;
    }),
    false,
    n('syncGroupMap/chat'),
  );
};

export const chatGroupAction: StateCreator<
  ChatGroupStore,
  [['zustand/devtools', never]],
  [],
  ChatGroupAction
> = (set, get) => {
  const dispatch: ChatGroupAction['internal_dispatchChatGroup'] = (payload) => {
    set(
      produce((draft: ChatGroupState) => {
        const reducer = chatGroupReducers[
          payload.type as keyof typeof chatGroupReducers
        ] as ChatGroupReducer;
        if (reducer) {
          // Apply the reducer and return the new state
          return reducer(draft, payload);
        }
      }),
      false,
      payload,
    );
  };

  return {
    ...initialChatGroupState,

    addAgentsToGroup: async (groupId, agentIds) => {
      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.addAgentsToGroup()
      const userId = getUserId();
      const now = Date.now();
      
      await Promise.all(
        agentIds.map((agentId, index) =>
          DB.LinkChatGroupToAgent({
            chatGroupId: groupId,
            agentId,
            userId,
            enabled: 1,
            sortOrder: toNullInt(index),
            role: toNullString('member'),
            createdAt: now,
            updatedAt: now,
          })
        )
      );
      
      console.log('[ChatGroup] Added agents to group via direct DB', { groupId, count: agentIds.length });
      
      await get().internal_refreshGroups();
    },

    /**
     * @param silent - if true, do not switch to the new group session
     */
    createGroup: async (newGroup: Omit<NewChatGroup, 'userId'>, agentIds?: string[], silent = false) => {
      const { switchSession } = getSessionStoreState();

      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.createGroup()
      const userId = getUserId();
      const groupId = crypto.randomUUID();
      const now = Date.now();
      
      const dbGroup = await DB.CreateChatGroup({
        id: groupId,
        title: toNullString(newGroup.title || 'Untitled Group'),
        description: toNullString(newGroup.description || ''),
        config: toNullJSON(newGroup.config || {}),
        clientId: toNullString(newGroup.clientId || ''),
        userId,
        groupId: toNullString(newGroup.groupId || ''),
        pinned: Number(boolToInt(newGroup.pinned || false)),
        createdAt: now,
        updatedAt: now,
      });
      
      const group = mapChatGroupFromDB(dbGroup);
      
      console.log('[ChatGroup] Created group via direct DB', { groupId });


      if (agentIds && agentIds.length > 0) {
        // Already migrated above, use the same logic
        await Promise.all(
          agentIds.map((agentId, index) =>
            DB.LinkChatGroupToAgent({
              chatGroupId: group.id,
              agentId,
              userId,
              enabled: 1,
              sortOrder: toNullInt(index),
              role: toNullString('member'),
              createdAt: now,
              updatedAt: now,
            })
          )
        );

        // Wait a brief moment to ensure database transactions are committed
        // This prevents race condition where loadGroups() executes before member addition is fully persisted
        await new Promise<void>((resolve) => {
          setTimeout(resolve, 100);
        });
      }

      dispatch({ payload: group, type: 'addGroup' });

      await get().loadGroups();
      await getSessionStoreState().refreshSessions();

      if (!silent) {
        switchSession(group.id);
      }

      return group.id;
    },
    deleteGroup: async (id) => {
      // First, get all group members to identify virtual members
      // Note: ChatGroupAgentItem type is incorrectly defined in schema as agents table type
      // but getGroupAgents actually returns chatGroupsAgents junction table entries
      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.getGroupAgents()
      const userId = getUserId();
      const groupAgents = (await DB.GetChatGroupAgentLinks({ chatGroupId: id, userId })) as unknown as Array<{
        agentId: string;
        chatGroupId: string;
      }>;

      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.deleteGroup()
      await DB.DeleteChatGroup({ id, userId });
      
      console.log('[ChatGroup] Deleted group via direct DB', { id });
      dispatch({ payload: id, type: 'deleteGroup' });

      // Now delete virtual members (agents with virtual: true)
      const sessionStore = getSessionStoreState();
      const sessions = sessionStore.sessions || [];

      // Find and delete all virtual sessions that were members of this group
      const virtualMemberDeletions = groupAgents
        .map((groupAgent) => {
          // groupAgent has agentId property from the junction table
          const session = sessions.find((s) => {
            // Type guard: check if it's an agent session
            if (s.type === 'agent') {
              return s.config?.id === groupAgent.agentId;
            }
            return false;
          });

          // Only delete if the session exists and has virtual flag set to true
          if (session && session.type === 'agent' && session.config?.virtual) {
            return sessionStore.removeSession(session.id);
          }
          return null;
        })
        .filter(Boolean);

      // Wait for all virtual member deletions to complete
      await Promise.all(virtualMemberDeletions);

      await get().loadGroups();
      await getSessionStoreState().refreshSessions();

      // If the active session is the deleted group, switch to the inbox session
      if (sessionStore.activeId === id) {
        sessionStore.switchSession(INBOX_SESSION_ID);
      }
    },

    internal_dispatchChatGroup: dispatch,

    internal_refreshGroups: async () => {
      await get().loadGroups();

      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.getGroups()
      const userId = getUserId();
      const dbGroups = await DB.ListChatGroups(userId);
      const groups = dbGroups.map(mapChatGroupFromDB);
      
      console.log('[ChatGroup] Fetched groups via direct DB', { count: groups.length });
      const nextGroupMap = groups.reduce(
        (map, group) => {
          map[group.id] = group;
          return map;
        },
        {} as Record<string, ChatGroupItem>,
      );

      if (!isEqual(get().groupMap, nextGroupMap)) {
        set(
          {
            groupMap: nextGroupMap,
            groupsInit: true,
            isGroupsLoading: false,
          },
          false,
          n('internal_refreshGroups/updateGroupMap'),
        );

        syncChatStoreGroupMap(nextGroupMap);
      }

      // Refresh sessions so session-related group info stays up to date
      await getSessionStoreState().refreshSessions();
    },

    internal_updateGroupMaps: (groups) => {
      // Build a candidate map from incoming groups
      const incomingMap = groups.reduce(
        (map, group) => {
          map[group.id] = group;
          return map;
        },
        {} as Record<string, ChatGroupItem>,
      );

      // Merge with existing map, preserving existing config if present
      const mergedMap = produce(get().groupMap, (draft) => {
        for (const id of Object.keys(incomingMap)) {
          const incoming = incomingMap[id];
          const existing = draft[id];
          if (existing) {
            draft[id] = {
              ...existing,
              ...incoming,
              // Keep existing config (authoritative) if present; do not overwrite
              config: existing.config || incoming.config,
            } as ChatGroupItem;
          } else {
            draft[id] = incoming;
          }
        }
      });

      set(
        {
          groupMap: mergedMap,
          groupsInit: true,
          isGroupsLoading: false,
        },
        false,
        n('internal_updateGroupMaps/chatGroup'),
      );

      syncChatStoreGroupMap(mergedMap);
    },

    loadGroups: async () => {
      dispatch({ payload: true, type: 'setGroupsLoading' });
      
      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.getGroups()
      const userId = getUserId();
      const dbGroups = await DB.ListChatGroups(userId);
      const groups = dbGroups.map(mapChatGroupFromDB);
      
      console.log('[ChatGroup] Loaded groups via direct DB', { count: groups.length });
      
      dispatch({ payload: groups, type: 'loadGroups' });
    },

    pinGroup: async (id, pinned) => {
      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.updateGroup()
      const userId = getUserId();
      const now = Date.now();
      
      await DB.UpdateChatGroup({
        id,
        userId,
        title: toNullString(''), // Will be ignored by update
        description: toNullString(''),
        config: toNullString(''),
        pinned: Number(boolToInt(pinned)),
        updatedAt: now,
      });
      
      console.log('[ChatGroup] Pinned group via direct DB', { id, pinned });
      
      dispatch({ payload: { id, pinned }, type: 'updateGroup' });
      await get().internal_refreshGroups();
    },

    refreshGroupDetail: async (groupId: string) => {
      try {
        // 🔄 MIGRATED: Direct DB call instead of chatGroupService.getGroup()
        const userId = getUserId();
        const dbGroup = await DB.GetChatGroup({ id: groupId, userId });
        if (!dbGroup) throw new Error(`Group ${groupId} not found`);
        
        const group = mapChatGroupFromDB(dbGroup);
        
        console.log('[ChatGroup] Refreshed group detail via direct DB', { groupId });
        
        const currentGroup = get().groupMap[group.id];
        if (isEqual(currentGroup, group)) return;

        set(
          {
            groupMap: { ...get().groupMap, [group.id]: group },
          },
          false,
          n('refreshGroupDetail'),
        );
      } catch (error) {
        console.error('[refreshGroupDetail] Error:', error);
      }
    },

    refreshGroups: async () => {
      try {
        // 🔄 MIGRATED: Direct DB call instead of chatGroupService.getGroups()
        const userId = getUserId();
        const dbGroups = await DB.ListChatGroups(userId);
        const groups = dbGroups.map(mapChatGroupFromDB);
        
        console.log('[ChatGroup] Refreshed groups via direct DB', { count: groups.length });
        const incomingMap = groups.reduce(
          (map, group) => {
            map[group.id] = group;
            return map;
          },
          {} as Record<string, ChatGroupItem>,
        );

        if (!isEqual(get().groupMap, incomingMap)) {
          set({ groupMap: incomingMap, groups }, false, n('refreshGroups'));
          syncChatStoreGroupMap(incomingMap);
        }
      } catch (error) {
        console.error('[refreshGroups] Error:', error);
      }
    },

    removeAgentFromGroup: async (groupId, agentId) => {
      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.removeAgentsFromGroup()
      const userId = getUserId();
      
      await DB.UnlinkChatGroupFromAgent({
        chatGroupId: groupId,
        agentId,
        userId,
      });
      
      console.log('[ChatGroup] Removed agent from group via direct DB', { groupId, agentId });
      
      await get().internal_refreshGroups();
    },

    reorderGroupMembers: async (groupId, orderedAgentIds) => {
      console.log('REORDER GROUP MEMBERS', groupId, orderedAgentIds);

      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.updateAgentInGroup()
      const userId = getUserId();
      const now = Date.now();
      
      await Promise.all(
        orderedAgentIds.map((agentId, index) =>
          DB.UpdateChatGroupAgentOrder({
            chatGroupId: groupId,
            agentId,
            userId,
            sortOrder: toNullInt(index),
            updatedAt: now,
          })
        ),
      );
      
      console.log('[ChatGroup] Reordered group members via direct DB', { groupId, count: orderedAgentIds.length });

      await get().internal_refreshGroups();
    },

    toggleGroupSetting: (open) => {
      set({ showGroupSetting: open }, false, 'toggleGroupSetting');
    },

    toggleThread: (agentId) => {
      set({ activeThreadAgentId: agentId }, false, 'toggleThread');
    },

    updateGroup: async (id: string, value: Partial<ChatGroupItem>) => {
      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.updateGroup()
      const userId = getUserId();
      const now = Date.now();
      
      await DB.UpdateChatGroup({
        id,
        userId,
        title: value.title ? toNullString(value.title) : toNullString(''),
        description: value.description ? toNullString(value.description) : toNullString(''),
        config: value.config ? toNullJSON(value.config) : toNullString(''),
        pinned: value.pinned !== undefined ? Number(boolToInt(value.pinned)) : 0,
        updatedAt: now,
      });
      
      console.log('[ChatGroup] Updated group via direct DB', { id });
      
      dispatch({ payload: { id, value }, type: 'updateGroup' });
      await get().internal_refreshGroups();
    },

    updateGroupConfig: async (config) => {
      const group = chatGroupSelectors.currentGroup(get());
      if (!group) return;

      const mergedConfig = {
        ...DEFAULT_CHAT_GROUP_CHAT_CONFIG,
        ...group.config,
        ...config,
      };

      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.updateGroup()
      const userId = getUserId();
      const now = Date.now();
      
      await DB.UpdateChatGroup({
        id: group.id,
        userId,
        title: toNullString(group.title || ''),
        description: toNullString(group.description || ''),
        config: toNullJSON(mergedConfig),
        pinned: Number(boolToInt(group.pinned || false)),
        updatedAt: now,
      });
      
      console.log('[ChatGroup] Updated group config via direct DB', { id: group.id });

      // Immediately update the local store to ensure configuration is available
      // Note: reducer expects payload: { id, value }
      dispatch({
        payload: { id: group.id, value: { config: mergedConfig } },
        type: 'updateGroup',
      });

      // Also update the chat store's groupMaps to keep it in sync
      useChatStore.setState(
        produce((draft: ChatStoreState) => {
          const existing = draft.groupMaps[group.id];
          if (existing) {
            draft.groupMaps[group.id] = {
              ...existing,
              config: mergedConfig,
            };
          }
        }),
        false,
        n('updateGroupConfig/syncChatStore'),
      );

      // Refresh groups to ensure consistency
      await get().internal_refreshGroups();
    },

    updateGroupMeta: async (meta: Partial<ChatGroupItem>) => {
      const group = chatGroupSelectors.currentGroup(get());
      if (!group) return;

      const id = group.id;

      // 🔄 MIGRATED: Direct DB call instead of chatGroupService.updateGroup()
      const userId = getUserId();
      const now = Date.now();
      
      await DB.UpdateChatGroup({
        id,
        userId,
        title: meta.title ? toNullString(meta.title) : toNullString(group.title || ''),
        description: meta.description ? toNullString(meta.description) : toNullString(group.description || ''),
        config: meta.config ? toNullJSON(meta.config) : toNullJSON(group.config || {}),
        pinned: meta.pinned !== undefined ? Number(boolToInt(meta.pinned)) : Number(boolToInt(group.pinned || false)),
        updatedAt: now,
      });
      
      console.log('[ChatGroup] Updated group meta via direct DB', { id });
      // Keep local store in sync immediately
      dispatch({ payload: { id, value: meta }, type: 'updateGroup' });
      await get().internal_refreshGroups();
    },

    internal_fetchGroupDetail: async (enabled: boolean, groupId: string) => {
      if (!enabled || !groupId) return;

      try {
        // 🔄 MIGRATED: Direct DB call instead of chatGroupService.getGroup()
        const userId = getUserId();
        const dbGroup = await DB.GetChatGroup({ id: groupId, userId });
        if (!dbGroup) throw new Error(`Group ${groupId} not found`);
        
        const group = mapChatGroupFromDB(dbGroup);

        const currentGroup = get().groupMap[group.id];
        if (isEqual(currentGroup, group)) return;

        const nextGroupMap = {
          ...get().groupMap,
          [group.id]: group,
        };

        set({ groupMap: nextGroupMap }, false, n('internal_fetchGroupDetail/onSuccess'));
        syncChatStoreGroupMap(nextGroupMap);
      } catch (error) {
        console.error('[internal_fetchGroupDetail] Error:', error);
      }
    },

    internal_fetchGroups: async (enabled: boolean, isLogin: boolean) => {
      if (!enabled) return;

      try {
        // 🔄 MIGRATED: Direct DB call instead of chatGroupService.getGroups()
        const userId = getUserId();
        const dbGroups = await DB.ListChatGroups(userId);
        const groups = dbGroups.map(mapChatGroupFromDB);
        const incomingMap = groups.reduce(
          (map, group) => {
            map[group.id] = group;
            return map;
          },
          {} as Record<string, ChatGroupItem>,
        );

        const currentMap = get().groupMap;
        const nextGroupMap = { ...currentMap, ...incomingMap };

        if (get().groupsInit && isEqual(currentMap, nextGroupMap)) {
          return;
        }

        set(
          {
            groupMap: nextGroupMap,
            groups,
            groupsInit: true,
            isGroupsLoading: false,
          },
          false,
          n('internal_fetchGroups/onSuccess'),
        );

        syncChatStoreGroupMap(nextGroupMap);
      } catch (error) {
        console.error('[internal_fetchGroups] Error:', error);
      }
    },

    useFetchGroupDetail: async (enabled: boolean, groupId: string) => {
      const store = get();
      await (store as any).internal_fetchGroupDetail(enabled, groupId);
    },

    useFetchGroups: async (enabled: boolean, isLogin: boolean) => {
      const store = get();
      await (store as any).internal_fetchGroups(enabled, isLogin);
    },
  };
};
