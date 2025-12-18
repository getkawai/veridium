import { produce } from 'immer';
import type { StateCreator } from 'zustand/vanilla';

import type { GlobalStore } from '@/store/global';
import { ProfileTabs, SettingsTabs, SidebarTabKey } from '@/store/global/initialState';
import { setNamespace } from '@/utils/storeDebug';

const n = setNamespace('w');

export interface GlobalWorkspacePaneAction {
  switchBackToChat: (sessionId?: string) => void;
  switchSideBar: (key: SidebarTabKey) => void;
  switchToImage: () => void;
  toggleAgentSystemRoleExpand: (agentId: string, expanded?: boolean) => void;
  toggleChatSideBar: (visible?: boolean) => void;
  toggleExpandInputActionbar: (expand?: boolean) => void;
  toggleExpandSessionGroup: (id: string, expand: boolean) => void;
  toggleMobilePortal: (visible?: boolean) => void;
  toggleChangelog: (visible?: boolean) => void;
  toggleSettings: (visible?: boolean, tab?: SettingsTabs) => void;
  toggleSystemRole: (visible?: boolean) => void;
  toggleUserProfile: (visible?: boolean, tab?: ProfileTabs) => void;
  toggleWideScreen: (enable?: boolean) => void;
  toggleZenMode: () => void;
}

export const globalWorkspaceSlice: StateCreator<
  GlobalStore,
  [['zustand/devtools', never]],
  [],
  GlobalWorkspacePaneAction
> = (set, get) => ({
  switchBackToChat: (sessionId) => {
    const isMobile = get().isMobile;
    set({ sidebarKey: SidebarTabKey.Chat }, false, n('switchBackToChat'));

    if (isMobile) {
      get().updateSystemStatus({ mobileShowTopic: true });
    }
  },

  switchSideBar: (key) => {
    set({ sidebarKey: key }, false, n('switchSideBar', key));
  },

  switchToImage: () => {
    set({ sidebarKey: SidebarTabKey.Image }, false, n('switchToImage'));
  },

  toggleAgentSystemRoleExpand: (agentId, expanded) => {
    const { status } = get();
    const systemRoleExpandedMap = status.systemRoleExpandedMap || {};
    const nextExpanded = typeof expanded === 'boolean' ? expanded : !systemRoleExpandedMap[agentId];

    get().updateSystemStatus(
      {
        systemRoleExpandedMap: {
          ...systemRoleExpandedMap,
          [agentId]: nextExpanded,
        },
      },
      n('toggleAgentSystemRoleExpand', { agentId, expanded: nextExpanded }),
    );
  },
  toggleChatSideBar: (newValue) => {
    const showChatSideBar =
      typeof newValue === 'boolean' ? newValue : !get().status.showChatSideBar;

    get().updateSystemStatus({ showChatSideBar }, n('toggleAgentPanel', newValue));
  },
  toggleExpandInputActionbar: (newValue) => {
    const expandInputActionbar =
      typeof newValue === 'boolean' ? newValue : !get().status.expandInputActionbar;

    get().updateSystemStatus({ expandInputActionbar }, n('toggleExpandInputActionbar', newValue));
  },
  toggleExpandSessionGroup: (id, expand) => {
    const { status } = get();
    const nextExpandSessionGroup = produce(status.expandSessionGroupKeys, (draft: string[]) => {
      if (expand) {
        if (draft.includes(id)) return;
        draft.push(id);
      } else {
        const index = draft.indexOf(id);
        if (index !== -1) draft.splice(index, 1);
      }
    });
    get().updateSystemStatus({ expandSessionGroupKeys: nextExpandSessionGroup });
  },
  toggleMobilePortal: (newValue) => {
    const mobileShowPortal =
      typeof newValue === 'boolean' ? newValue : !get().status.mobileShowPortal;

    get().updateSystemStatus({ mobileShowPortal }, n('toggleMobilePortal', newValue));
  },
  toggleMobileTopic: (newValue) => {
    const mobileShowTopic =
      typeof newValue === 'boolean' ? newValue : !get().status.mobileShowTopic;

    get().updateSystemStatus({ mobileShowTopic }, n('toggleMobileTopic', newValue));
  },
  toggleChangelog: (newValue) => {
    const isShowChangelog = typeof newValue === 'boolean' ? newValue : !get().status.isShowChangelog;

    get().updateSystemStatus({ isShowChangelog }, n('toggleChangelog', newValue));
  },
  toggleSettings: (newValue, tab) => {
    const isShowSettings = typeof newValue === 'boolean' ? newValue : !get().status.isShowSettings;

    get().updateSystemStatus(
      { isShowSettings, settingsTab: tab || get().status.settingsTab },
      n('toggleSettings', newValue),
    );
  },
  toggleSystemRole: (newValue) => {
    const showSystemRole = typeof newValue === 'boolean' ? newValue : !get().status.mobileShowTopic;

    get().updateSystemStatus({ showSystemRole }, n('toggleMobileTopic', newValue));
  },
  toggleUserProfile: (newValue, tab) => {
    const isShowUserProfile =
      typeof newValue === 'boolean' ? newValue : !get().status.isShowUserProfile;

    get().updateSystemStatus(
      { isShowUserProfile, profileTab: tab || get().status.profileTab },
      n('toggleUserProfile', newValue),
    );
  },
  toggleWideScreen: (newValue) => {
    const wideScreen = typeof newValue === 'boolean' ? newValue : !get().status.noWideScreen;

    get().updateSystemStatus({ noWideScreen: wideScreen }, n('toggleWideScreen', newValue));
  },
  toggleZenMode: () => {
    const { status } = get();
    const nextZenMode = !status.zenMode;

    get().updateSystemStatus({ zenMode: nextZenMode }, n('toggleZenMode'));
  },
});
