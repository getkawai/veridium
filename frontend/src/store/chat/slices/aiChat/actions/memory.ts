import { StateCreator } from 'zustand/vanilla';

import { ChatStore } from '@/store/chat';

export interface ChatMemoryAction {
  // No actions currently
}

export const chatMemory: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatMemoryAction
> = (set, get) => ({
  // No memory actions currently implemented
});
