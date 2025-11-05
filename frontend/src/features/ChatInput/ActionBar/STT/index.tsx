import { memo } from 'react';

import NativeSTT from './native';

const STT = memo<{ mobile?: boolean }>(({ mobile }) => {
  return <NativeSTT mobile={mobile} />;
});

export default STT;

