'use client';

import { ChatHeader } from '@lobehub/ui/chat';
import { memo } from 'react';

import { ProductLogo } from '@/components/Branding';

import StoreSearchBar from './Search';

const Header = memo(() => {
  return (
    <ChatHeader
      left={<ProductLogo extra={'Wallet'} size={36} type={'text'} />}
      right={<ProductLogo extra={'Veridium'} size={36} type={'text'} />}
      style={{
        position: 'relative',
        zIndex: 10,
      }}
      styles={{
        center: { flex: 1, maxWidth: 1440 },
        left: { flex: 1, maxWidth: 240 },
        right: { flex: 1, maxWidth: 240 },
      }}
    >
      <StoreSearchBar />
    </ChatHeader>
  );
});

export default Header;
