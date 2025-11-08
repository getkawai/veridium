import { StateCreator } from 'zustand';

import { ChatStore } from '../../store';
import { ChatTTS } from '@/types';

// TTS functionality moved to backend (native OS TTS)
// These actions kept for backward compatibility but no longer used
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
  // No-op: TTS now handled by backend via useNativeTTS hook
  clearTTS: (id) => {
    // Deprecated: TTS state no longer stored in frontend
  },

  // No-op: TTS now handled by backend via useNativeTTS hook
  ttsMessage: (id, tts) => {
    // Deprecated: TTS state no longer stored in frontend
  },
});

