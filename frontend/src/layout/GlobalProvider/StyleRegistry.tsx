'use client';

import { StyleProvider, extractStaticStyle } from 'antd-style';
import { PropsWithChildren } from 'react';

/**
 * StyleRegistry adapted for Vite/CSR
 * In Vite, we don't need SSR hydration, so this is simplified
 */
const StyleRegistry = ({ children }: PropsWithChildren) => {
  return <StyleProvider cache={extractStaticStyle.cache}>{children}</StyleProvider>;
};

export default StyleRegistry;