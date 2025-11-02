import { useEffect } from 'react';

import { lobeChat } from '../lobeChat';
import { PluginPayload } from '../lobeChat';

export const useOnStandalonePluginInit = <T = any>(
  callback: (payload: PluginPayload<T>) => void,
) => {
  useEffect(() => {
    lobeChat.getPluginPayload().then((e) => {
      if (!e) return;

      callback(e);
    });
  }, []);
};
