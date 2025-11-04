import isEqual from 'fast-deep-equal';
import { memo } from 'react';

import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
import { useUserStore } from '@/store/user';
import { settingsSelectors } from '@/store/user/selectors';

import BrowserSTT from './browser';
import OpenaiSTT from './openai';

// Check if media devices are available (not in Wails desktop environment)
const isBrowserSTTAvailable = () => {
  return typeof navigator !== 'undefined' && 
    navigator.mediaDevices && 
    typeof navigator.mediaDevices.getUserMedia === 'function';
};

const STT = memo<{ mobile?: boolean }>(({ mobile }) => {
  const { sttServer } = useUserStore(settingsSelectors.currentTTS, isEqual);

  const { enableSTT } = useServerConfigStore(featureFlagsSelectors);
  if (!enableSTT) return;

  switch (sttServer) {
    case 'openai': {
      return <OpenaiSTT mobile={mobile} />;
    }
  }
  
  // Only show browser STT if media devices are available
  if (!isBrowserSTTAvailable()) return null;
  
  return <BrowserSTT mobile={mobile} />;
});

export default STT;
