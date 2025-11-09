import { useEffect } from 'react';

import { useGlobalStore } from '@/store/global';
import { systemStatusSelectors } from '@/store/global/selectors';
import { useSessionStore } from '@/store/session';
import { useUserStore } from '@/store/user';
import { authSelectors } from '@/store/user/slices/auth/selectors';

export const useFetchSessions = () => {
  const isDBInited = useGlobalStore(systemStatusSelectors.isDBInited);
  const isLogin = useUserStore(authSelectors.isLogin);
  const internal_fetchSessions = useSessionStore((s) => s.internal_fetchSessions);

  useEffect(() => {
    if (isDBInited) {
      internal_fetchSessions(true, isLogin);
    }
  }, [isDBInited, isLogin, internal_fetchSessions]);
};
