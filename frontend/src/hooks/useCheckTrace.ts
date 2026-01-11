import { useCallback } from 'react';

import { useUserStore } from '@/store/user';
import { preferenceSelectors } from '@/store/user/selectors';

interface CheckTraceResult {
  data: {
    canEnableTrace: boolean;
    showModal: boolean;
  };
  mutate: (telemetry: boolean) => void;
}

export const useCheckTrace = (isPreferenceInit: boolean): CheckTraceResult => {
  const [updatePreference] = useUserStore((s) => s.updatePreference);
  const userAllowTrace = useUserStore(preferenceSelectors.userAllowTrace);

  const canEnableTrace = typeof userAllowTrace === 'boolean' ? userAllowTrace : false;

  const mutate = useCallback(
    (telemetry: boolean) => {
      updatePreference({ telemetry: { enabled: telemetry } });
    },
    [updatePreference],
  );

  return {
    data: {
      canEnableTrace,
      showModal: isPreferenceInit && typeof userAllowTrace !== 'boolean',
    },
    mutate,
  };
};
