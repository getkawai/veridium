import { PropsWithChildren, Suspense } from 'react';
import { Flexbox } from 'react-layout-kit';

// import { isDesktop } from '@/const/version';
// import ProtocolUrlHandler from '@/features/ProtocolUrlHandler';

import RegisterHotkeys from './RegisterHotkeys';
import SessionPanel from './SessionPanel';
import Session from './@session/default';
import { useTheme } from 'antd-style';
import Workspace from './Workspace';

const Layout = ({ children }: PropsWithChildren) => {
  const theme = useTheme();
  return (
    <>
      <Flexbox
        height={'100%'}
        horizontal
        style={{ maxWidth: '100%', overflow: 'hidden', position: 'relative' }}
        width={'100%'}
      >
        <SessionPanel>
          <Session />
        </SessionPanel>
        <Flexbox
          flex={1}
          style={{
            // background: theme.colorBgContainerSecondary,
            background: theme.colorBgContainer,
            overflow: 'hidden',
            position: 'relative',
          }}
        >
          <Workspace>{children}</Workspace>
        </Flexbox>
      </Flexbox>
      <Suspense>
        <RegisterHotkeys />
      </Suspense>
      {/* {isDesktop && <ProtocolUrlHandler />} */}
    </>
  );
};

Layout.displayName = 'DesktopChatLayout';

export default Layout;
