'use client';

import { memo } from 'react';

// import { useGlobalStore } from '@/store/global';
// import { systemStatusSelectors } from '@/store/global/selectors';

// Dummy implementations for development - memoized
const mockGlobalStore = {
  inZenMode: true
};

const useGlobalStore = (selector?: any) => {
  if (selector) {
    return selector(mockGlobalStore);
  }
  return mockGlobalStore;
};

const systemStatusSelectors = {
  inZenMode: (state: any) => state.inZenMode
};

import Toast from './Toast';

const ZenModeToast = memo(() => {
  const inZenMode = useGlobalStore(systemStatusSelectors.inZenMode);

  return inZenMode && <Toast />;
});

export default ZenModeToast;
