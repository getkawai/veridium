import { useEffect } from 'react'
import { Events, WML } from "@wailsio/runtime";
import GlobalLayout from './layout/GlobalProvider';
import DesktopMainLayout from './layout/Desktop';
import DesktopChatLayout from './app/chat';
import { useChatStore } from './store/chat';

function App() {
  useEffect(() => {
    Events.On('time', (timeValue: any) => {
    });

    Events.On('chat:topic:updated', (data: any) => {
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
    <GlobalLayout appearance={'auto'} isMobile={false} locale={''} neutralColor={undefined} primaryColor={undefined} variants={undefined}>
      <DesktopMainLayout>
        <DesktopChatLayout />
      </DesktopMainLayout>
    </GlobalLayout>
  )
}

export default App
