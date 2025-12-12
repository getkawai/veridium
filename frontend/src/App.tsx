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

function App() {
  const sidebarKey = useGlobalStore((s) => s.sidebarKey);

  useEffect(() => {
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

  return (
    <GlobalLayout appearance={'dark'} isMobile={false} locale={''} neutralColor={undefined} primaryColor={undefined} variants={undefined}>
      <DesktopMainLayout>
        {sidebarKey === SidebarTabKey.Chat && <DesktopChatLayout />}
        {sidebarKey === SidebarTabKey.Image && <DesktopImageLayout />}
        {sidebarKey === SidebarTabKey.Files && <KnowledgeHomePage />}
      </DesktopMainLayout>
    </GlobalLayout>
  )
}

export default App
