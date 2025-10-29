import { useEffect } from 'react'
import {Events, WML} from "@wailsio/runtime";
import CircleLoader from './components/CircleLoader';
import GlobalLayout from './layout/GlobalProvider';
import DesktopMainLayout from './layout/Desktop';
import ChatLayout from './layout/Chat';

function App() {
  useEffect(() => {
    Events.On('time', (timeValue: any) => {
    });
    // Reload WML so it picks up the wml tags
    WML.Reload();
  }, []);

  return (
    <GlobalLayout appearance={'auto'} isMobile={false} locale={''} neutralColor={undefined} primaryColor={undefined} variants={undefined}>
      <DesktopMainLayout>
        <ChatLayout session={<div>Session</div>}>
          <CircleLoader/>
        </ChatLayout>
      </DesktopMainLayout>
    </GlobalLayout>
  )
}

export default App
