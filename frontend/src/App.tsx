import { useEffect } from 'react'
import { Events, WML } from "@wailsio/runtime";
import GlobalLayout from './layout/GlobalProvider';
import DesktopMainLayout from './layout/Desktop';
import DesktopChatLayout from './app/chat';
import { useChatStore } from './store/chat';
import DesktopImageLayout from './app/image';
import { useGlobalStore } from './store/global';
import { SidebarTabKey } from './store/global/initialState';
import KnowledgeHomePage from './app/knowledge/routes/KnowledgeHome';
import SettingsModal from './features/SettingsModal';
import UserProfileModal from './features/UserProfileModal';
import ChangelogModal from './features/ChangelogModal';
import { useUserStore } from './store/user';
import AuthSignInBox from './features/User/AuthSignInBox';
import { Flex, Spin } from 'antd';
import DesktopWalletLayout from './app/wallet/wallet';

function App() {
  const sidebarKey = useGlobalStore((s) => s.sidebarKey);
  const { isWalletLocked, isWalletLoaded, refreshWalletStatus } = useUserStore((s) => ({
    isWalletLocked: s.isWalletLocked,
    isWalletLoaded: s.isWalletLoaded,
    refreshWalletStatus: s.refreshWalletStatus,
  }));

  useEffect(() => {
    refreshWalletStatus();

    Events.On('chat:topic:updated', (ev: any) => {
      const data = ev.data;
      if (data && data.topic_id && data.title) {
        useChatStore.getState().internal_dispatchTopic({
          type: 'updateTopic',
          id: data.topic_id,
          value: { title: data.title }
        }, 'updateTopicTitleFromEvent');
      }
    });

    // Global subscription to chat stream events
    Events.On('chat:stream', (ev: any) => {
      const data = ev.data;
      const activeId = useChatStore.getState().activeId;

      // Only handle events for the current active session
      if (data && data.session_id === activeId) {
        useChatStore.getState().internal_handleStreamEvent(data);
      }
    });

    // Reload WML so it picks up the wml tags
    WML.Reload();
  }, []);

  if (!isWalletLoaded) {
    return (
      <Flex align="center" justify="center" style={{ height: '100vh', width: '100vw' }}>
        <Spin size="large" />
      </Flex>
    );
  }

  return (
    <GlobalLayout appearance={'dark'} locale={''} neutralColor={undefined} primaryColor={undefined} variants={undefined}>
      {isWalletLocked ? (
        <Flex align="center" justify="center" style={{ height: '100vh', width: '100vw', background: 'var(--color-bg-layout)' }}>
          <AuthSignInBox />
        </Flex>
      ) : (
        <>
          <DesktopMainLayout>
            {sidebarKey === SidebarTabKey.Chat && <DesktopChatLayout />}
            {sidebarKey === SidebarTabKey.Wallet && <DesktopWalletLayout />}
            {sidebarKey === SidebarTabKey.Image && <DesktopImageLayout />}
            {sidebarKey === SidebarTabKey.Files && <KnowledgeHomePage />}
          </DesktopMainLayout>
          <SettingsModal />
          <UserProfileModal />
          <ChangelogModal />
        </>
      )}
    </GlobalLayout>
  )
}

export default App
