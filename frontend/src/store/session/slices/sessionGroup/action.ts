import { t } from 'i18next';
import { StateCreator } from 'zustand/vanilla';

import { message } from '@/components/AntdStaticMethods';
import { SessionGroupItem } from '@/types/session';
import * as DB from '@/bindings/github.com/kawai-network/veridium/internal/database/generated/queries';
import { getUserId } from '../../helpers';
import { toNullString, toNullInt } from '@/types/database';

import type { SessionStore } from '../../store';
import { SessionGroupsDispatch, sessionGroupsReducer } from './reducer';

/* eslint-disable typescript-sort-keys/interface */
export interface SessionGroupAction {
  addSessionGroup: (name: string) => Promise<string>;
  clearSessionGroups: () => Promise<void>;
  removeSessionGroup: (id: string) => Promise<void>;
  updateSessionGroupName: (id: string, name: string) => Promise<void>;
  updateSessionGroupSort: (items: SessionGroupItem[]) => Promise<void>;
  internal_dispatchSessionGroups: (payload: SessionGroupsDispatch) => void;
}
/* eslint-enable */

export const createSessionGroupSlice: StateCreator<
  SessionStore,
  [['zustand/devtools', never]],
  [],
  SessionGroupAction
> = (set, get) => ({
  addSessionGroup: async (name) => {
    // 🔄 MIGRATED: Direct DB call instead of sessionService.createSessionGroup()
    const userId = getUserId();
    const id = crypto.randomUUID();
    const now = Date.now();
    
    await DB.CreateSessionGroup({
      id,
      name: toNullString(name),
      sort: toNullInt(0),
      userId,
      createdAt: now,
      updatedAt: now,
    });
    
    console.log('[SessionGroup] Created session group via direct DB', { id, name });

    await get().refreshSessions();

    return id;
  },

  clearSessionGroups: async () => {
    // 🔄 MIGRATED: Direct DB call instead of sessionService.removeSessionGroups()
    const userId = getUserId();
    await DB.DeleteAllSessionGroups(userId);
    
    console.log('[SessionGroup] Cleared all session groups via direct DB');
    
    await get().refreshSessions();
  },

  removeSessionGroup: async (id) => {
    // 🔄 MIGRATED: Direct DB call instead of sessionService.removeSessionGroup()
    const userId = getUserId();
    await DB.DeleteSessionGroup({ id, userId });
    
    console.log('[SessionGroup] Deleted session group via direct DB', { id });
    
    await get().refreshSessions();
  },

  updateSessionGroupName: async (id, name) => {
    // 🔄 MIGRATED: Direct DB call instead of sessionService.updateSessionGroup()
    const userId = getUserId();
    const now = Date.now();
    
    await DB.UpdateSessionGroup({
      id,
      userId,
      name: toNullString(name),
      updatedAt: now,
    } as any);
    
    console.log('[SessionGroup] Updated session group name via direct DB', { id, name });
    
    await get().refreshSessions();
  },
  updateSessionGroupSort: async (items) => {
    const sortMap = items.map((item, index) => ({ id: item.id, sort: index }));

    get().internal_dispatchSessionGroups({ sortMap, type: 'updateSessionGroupOrder' });

    message.loading({
      content: t('sessionGroup.sorting', { ns: 'chat' }),
      duration: 0,
      key: 'updateSessionGroupSort',
    });

    // 🔄 MIGRATED: Direct DB call instead of sessionService.updateSessionGroupOrder()
    const userId = getUserId();
    await DB.UpdateSessionGroupOrder({
      userId,
      sortMap: JSON.stringify(sortMap),
    });
    
    console.log('[SessionGroup] Updated session group sort via direct DB', { count: sortMap.length });
    
    message.destroy('updateSessionGroupSort');
    message.success(t('sessionGroup.sortSuccess', { ns: 'chat' }));

    await get().refreshSessions();
  },

  /* eslint-disable sort-keys-fix/sort-keys-fix */
  internal_dispatchSessionGroups: (payload) => {
    const nextSessionGroups = sessionGroupsReducer(get().sessionGroups, payload);
    get().internal_processSessions(get().sessions, nextSessionGroups, 'updateSessionGroups');
  },
});
