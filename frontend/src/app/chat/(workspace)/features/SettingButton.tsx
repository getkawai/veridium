'use client';

import { ActionIcon } from '@lobehub/ui';
import { AlignJustify } from 'lucide-react';
// import dynamic from 'next/dynamic';
import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { DESKTOP_HEADER_ICON_SIZE, MOBILE_HEADER_ICON_SIZE } from '@/const/layoutTokens';
// import { useOpenChatSettings } from '@/hooks/useInterceptingRoutes';
// import { useChatGroupStore } from '@/store/chatGroup';
// import { useSessionStore } from '@/store/session';
// import { sessionSelectors } from '@/store/session/selectors';
// import { useUserStore } from '@/store/user';
// import { settingsSelectors } from '@/store/user/selectors';
// import { HotkeyEnum } from '@/types/hotkey';

// const AgentSettings = dynamic(() => import('./AgentSettings'), {
//   ssr: false,
// });

// const AgentTeamSettings = dynamic(() => import('./AgentTeamSettings'), {
//   ssr: false,
// });

// Dummy components to replace dynamic imports
// const AgentSettings = memo(() => null);
// const AgentTeamSettings = memo(() => null);
import AgentSettings from './AgentSettings';
import AgentTeamSettings from './AgentTeamSettings';

const SettingButton = memo<{ mobile?: boolean }>(({ mobile }) => {
  // Dummy data for hotkey using useMemo to prevent infinite loops
  const hotkey = useMemo(() => 'Ctrl+Shift+S', []);

  const { t } = useTranslation('common');

  // Dummy data for session using useMemo to prevent infinite loops
  const id = useMemo(() => 'dummy-session-id', []);
  const isGroupSession = useMemo(() => false, []);

  // Dummy functions using useMemo to prevent infinite loops
  const openChatSettings = useMemo(() => () => console.log('Open chat settings'), []);
  const toggleGroupSetting = useMemo(() => () => console.log('Toggle group setting'), []);

  return (
    <>
      <ActionIcon
        icon={AlignJustify}
        onClick={() => (isGroupSession ? toggleGroupSetting() : openChatSettings())}
        size={mobile ? MOBILE_HEADER_ICON_SIZE : DESKTOP_HEADER_ICON_SIZE}
        title={t('openChatSettings.title', { ns: 'hotkey' })}
        tooltipProps={{
          hotkey,
          placement: 'bottom',
        }}
      />

      {isGroupSession ? <AgentTeamSettings key={id} /> : <AgentSettings key={id} />}
    </>
  );
});

export default SettingButton;
