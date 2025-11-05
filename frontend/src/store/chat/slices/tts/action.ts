import { StateCreator } from 'zustand';

import { ChatStore } from '../../store';
import { ChatTTS } from '@/types';

export interface TTSAction {
  clearTTS: (id: string) => void;
  ttsMessage: (id: string, tts: ChatTTS) => void;
}

export const createTTSSlice: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  TTSAction
> = (set, get) => ({
  clearTTS: (id) => {
    const { messagesMap } = get();
    const message = messagesMap[id];
    if (!message) return;

    set(
      {
        messagesMap: {
          ...messagesMap,
          [id]: {
            ...message,
            extra: {
              ...message.extra,
              tts: undefined,
            },
          },
        },
      },
      false,
      'clearTTS',
    );
  },

  ttsMessage: (id, tts) => {
    const { messagesMap } = get();
    const message = messagesMap[id];
    if (!message) return;

    set(
      {
        messagesMap: {
          ...messagesMap,
          [id]: {
            ...message,
            extra: {
              ...message.extra,
              tts,
            },
          },
        },
      },
      false,
      'ttsMessage',
    );
  },
});

