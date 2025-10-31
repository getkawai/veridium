import { useCallback } from 'react';

import { useQueryRoute } from '@/hooks/useQueryRoute';
import { useChatStore } from '@/store/chat';
import { routerSelectors, useRouterStore } from '@/store/router';
import { useServerConfigStore } from '@/store/serverConfig';
import { useSessionStore } from '@/store/session';

export const useSwitchSession = () => {
  const switchSession = useSessionStore((s) => s.switchSession);
  const togglePortal = useChatStore((s) => s.togglePortal);
  const mobile = useServerConfigStore((s) => s.isMobile);
  const router = useQueryRoute();
  const pathname = useRouterStore(routerSelectors.pathname);

  return useCallback(
    (id: string) => {
      switchSession(id);
      togglePortal(false);

      const chatPath = '/chat';
      if (mobile || pathname !== chatPath) {
        setTimeout(() => {
          router.push(chatPath, {
            query: { session: id, showMobileWorkspace: 'true' },
          });
        }, 50);
      }
    },
    [mobile, pathname, router, switchSession, togglePortal],
  );
};
