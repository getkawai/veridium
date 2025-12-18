import { PropsWithChildren, Suspense, memo } from 'react';
import { HotkeysProvider } from 'react-hotkeys-hook';
import { Flexbox } from 'react-layout-kit';

import CloudBanner from '@/features/AlertBanner/CloudBanner';
import TitleBar, { TITLE_BAR_HEIGHT } from '@/features/ElectronTitlebar';
import HotkeyHelperPanel from '@/features/HotkeyHelperPanel';
import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
import { HotkeyScopeEnum } from '@/types/hotkey';

import DesktopLayoutContainer from './DesktopLayoutContainer';
import RegisterHotkeys from './RegisterHotkeys';

const DesktopMainLayout = memo<PropsWithChildren>(({ children }) => {
  const { showCloudPromotion } = useServerConfigStore(featureFlagsSelectors);
  return (
    <HotkeysProvider initiallyActiveScopes={[HotkeyScopeEnum.Global]}>
      <TitleBar />
      {showCloudPromotion && <CloudBanner />}
      <Flexbox
        height={`calc(100% - ${TITLE_BAR_HEIGHT}px)`}
        horizontal
        style={{
          position: 'relative',
        }}
        width={'100%'}
      >
        <DesktopLayoutContainer>{children}</DesktopLayoutContainer>
      </Flexbox>
      <HotkeyHelperPanel />
      <Suspense>
        <RegisterHotkeys />
      </Suspense>
    </HotkeysProvider>
  );
});

export default DesktopMainLayout;
