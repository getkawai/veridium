import { PropsWithChildren, Suspense } from 'react';
import { Flexbox } from 'react-layout-kit';

// import { isDesktop } from '@/const/version';
// import ProtocolUrlHandler from '@/features/ProtocolUrlHandler';

import RegisterHotkeys from './RegisterHotkeys';
import SessionPanel from './SessionPanel';
import Workspace from './Workspace';
import Session from './@session/default';

const Layout = ({ children }: PropsWithChildren) => {
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
        <Workspace>{children}</Workspace>
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
