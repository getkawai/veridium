import { ReactNode } from 'react';
import AntdV5MonkeyPatch from './AntdV5MonkeyPatch';
import AppTheme from './AppTheme';
import StyleRegistry from './StyleRegistry';
import StoreInitialization from './StoreInitialization';
import { ServerConfigStoreProvider } from '@/store/serverConfig';

interface GlobalLayoutProps {
  appearance: string;
  children: ReactNode;
  isMobile: boolean;
  locale: string;
  neutralColor?: string;
  primaryColor?: string;
  variants?: string;
}

const GlobalLayout = ({
  children,
  neutralColor,
  primaryColor,
  locale: userLocale,
  appearance,
  isMobile,
  variants,
}: GlobalLayoutProps) => {
  return (
    <StyleRegistry>
        <AppTheme
          customFontFamily={undefined}
          customFontURL={undefined}
          defaultAppearance={appearance}
          defaultNeutralColor={neutralColor as any}
          defaultPrimaryColor={primaryColor as any}
          globalCDN={false}>
          <ServerConfigStoreProvider isMobile={isMobile} segmentVariants={variants}>
            <StoreInitialization />
            {children}
          </ServerConfigStoreProvider>
        </AppTheme>
      <AntdV5MonkeyPatch />
    </StyleRegistry>
  );
};

export default GlobalLayout;
