'use client';

import {
  ConfigProvider,
  FontLoader,
  NeutralColors,
  PrimaryColors,
  ThemeProvider,
} from '@lobehub/ui';
import { ThemeAppearance, createStyles } from 'antd-style';
import 'antd/dist/reset.css';
import {Image} from '@lobehub/ui';
import { ReactNode, memo } from 'react';

import AntdStaticMethods from '@/components/AntdStaticMethods';
import { GlobalStyle } from '@/styles';
import { setCookie } from '@/utils/client/cookie';

const useStyles = createStyles(({ css, token }) => ({
  app: css`
    position: relative;

    overscroll-behavior: none;
    display: flex;
    flex-direction: column;
    align-items: center;

    height: 100%;
    min-height: 100dvh;
    max-height: 100dvh;

    @media (min-device-width: 576px) {
      overflow: hidden;
    }
  `,
  // scrollbar-width and scrollbar-color are supported from Chrome 121
  // https://developer.mozilla.org/en-US/docs/Web/CSS/scrollbar-color
  scrollbar: css`
    scrollbar-color: ${token.colorFill} transparent;
    scrollbar-width: thin;

    #lobe-mobile-scroll-container {
      scrollbar-width: none;

      ::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }
  `,

  // so this is a polyfill for older browsers
  scrollbarPolyfill: css`
    ::-webkit-scrollbar {
      width: 0.75em;
      height: 0.75em;
    }

    ::-webkit-scrollbar-thumb {
      border-radius: 10px;
    }

    :hover::-webkit-scrollbar-thumb {
      border: 3px solid transparent;
      background-color: ${token.colorText};
      background-clip: content-box;
    }

    ::-webkit-scrollbar-track {
      background-color: transparent;
    }
  `,
}));

export interface AppThemeProps {
  children?: ReactNode;
  customFontFamily?: string;
  customFontURL?: string;
  defaultAppearance?: ThemeAppearance;
  defaultNeutralColor?: NeutralColors;
  defaultPrimaryColor?: PrimaryColors;
  globalCDN?: boolean;
}

const AppTheme = memo<AppThemeProps>(
  ({
    children,
    defaultAppearance,
    defaultPrimaryColor,
    defaultNeutralColor,
    globalCDN,
    customFontURL,
    customFontFamily,
  }) => {
    const { styles, cx, theme } = useStyles();

    return (
      <ThemeProvider
        appearance={'auto'}
        className={cx(styles.app, styles.scrollbar, styles.scrollbarPolyfill)}
        customTheme={{
          neutralColor: defaultNeutralColor,
          primaryColor: defaultPrimaryColor,
        }}
        defaultAppearance={defaultAppearance}
        onAppearanceChange={(appearance) => {
          setCookie('LOBE_THEME_APPEARANCE', appearance);
        }}
        theme={{
          cssVar: true,
          token: {
            fontFamily: customFontFamily ? `${customFontFamily},${theme.fontFamily}` : undefined,
            motion: true,
            motionUnit: 0.05,
          },
        }}
        themeMode={'auto'}
      >
        {!!customFontURL && <FontLoader url={customFontURL} />}
        <GlobalStyle />
        <AntdStaticMethods />
        <ConfigProvider
          config={{
            imgAs: Image,
            proxy: globalCDN ? 'unpkg' : undefined,
          }}
        >
          {children}
        </ConfigProvider>
      </ThemeProvider>
    );
  },
);

export default AppTheme;
