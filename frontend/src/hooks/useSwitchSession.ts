import { useCallback } from 'react';

import { useChatStore } from '@/store/chat';
import { useGlobalStore } from '@/store/global';
import { useServerConfigStore } from '@/store/serverConfig';
import { useSessionStore } from '@/store/session';

export const useSwitchSession = () => {
  const switchSession = useSessionStore((s) => s.switchSession);
  const togglePortal = useChatStore((s) => s.togglePortal);
  const mobile = useServerConfigStore((s) => s.isMobile);
  const updateSystemStatus = useGlobalStore((s) => s.updateSystemStatus);

  return useCallback(
    (id: string) => {
      switchSession(id);
      togglePortal(false);

      if (mobile) {
        updateSystemStatus({ mobileShowTopic: true });
      }
    },
    [mobile, switchSession, togglePortal, updateSystemStatus],
  );
};
