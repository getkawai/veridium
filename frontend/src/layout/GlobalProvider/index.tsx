import { ReactNode } from "react";
import { App } from "antd";
import AntdV5MonkeyPatch from "./AntdV5MonkeyPatch";
import AppTheme from "./AppTheme";
import StyleRegistry from "./StyleRegistry";
import StoreInitialization from "./StoreInitialization";
import { ServerConfigStoreProvider } from "@/store/serverConfig";
import Locale from "./Locale";

interface GlobalLayoutProps {
  appearance: string;
  children: ReactNode;
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
  variants,
}: GlobalLayoutProps) => {
  return (
    <StyleRegistry>
      <Locale antdLocale={undefined} defaultLang={userLocale}>
        <AppTheme
          customFontFamily={undefined}
          customFontURL={undefined}
          defaultAppearance={appearance}
          defaultNeutralColor={neutralColor as any}
          defaultPrimaryColor={primaryColor as any}
          globalCDN={false}
        >
          <App>
            <ServerConfigStoreProvider segmentVariants={variants}>
              <StoreInitialization />
              {children}
            </ServerConfigStoreProvider>
          </App>
        </AppTheme>
      </Locale>
      <AntdV5MonkeyPatch />
    </StyleRegistry>
  );
};

export default GlobalLayout;
