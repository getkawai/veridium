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
