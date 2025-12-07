'use client';

import { ChatHeader, ChatHeaderTitle } from '@lobehub/ui/chat';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import HeaderContent from '../../features/HeaderContent';

const Header = memo(() => {
  const { t } = useTranslation('setting');

  return (
    <ChatHeader
      left={<ChatHeaderTitle title={t('header.session')} />}
      right={<HeaderContent />}
      showBackButton
    />
  );
});

export default Header;
