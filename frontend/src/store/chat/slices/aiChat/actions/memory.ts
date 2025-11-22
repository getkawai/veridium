import { chainSummaryHistory } from '@/prompts';
import { TraceNameMap, UIChatMessage } from '@/types';
import { StateCreator } from 'zustand/vanilla';

import { chatService } from '@/services/chat';
import { ChatStore } from '@/store/chat';
import { useUserStore } from '@/store/user';
import { systemAgentSelectors } from '@/store/user/selectors';
import { DB, toNullString } from '@/types/database';

const DEFAULT_USER_ID = 'DEFAULT_LOBE_CHAT_USER';
const getUserId = () => DEFAULT_USER_ID;

export interface ChatMemoryAction {
  internal_summaryHistory: (messages: UIChatMessage[]) => Promise<void>;
}

export const chatMemory: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatMemoryAction
> = (set, get) => ({
  internal_summaryHistory: async (messages) => {
    const topicId = get().activeTopicId;
    if (messages.length <= 1 || !topicId) return;

    const { model, provider } = systemAgentSelectors.historyCompress(useUserStore.getState());

    let historySummary = '';
    await chatService.fetchPresetTaskResult({
      onFinish: async (text) => {
        historySummary = text;
      },
      params: { ...chainSummaryHistory(messages), model, provider, stream: false },
      trace: {
        sessionId: get().activeId,
        topicId: get().activeTopicId,
        traceName: TraceNameMap.SummaryHistoryMessages,
      },
    });

    // 🔄 MIGRATED: Direct DB call instead of topicService.updateTopic()
    const userId = getUserId();
    await DB.UpdateTopic({
      id: topicId,
      userId,
      historySummary: toNullString(historySummary),
      metadata: toNullString(JSON.stringify({ model, provider })),
      updatedAt: Date.now(),
    });
    
    console.log('[Memory] Updated topic history summary via direct DB', { topicId });
    
    await get().refreshTopic();
    await get().refreshMessages();
  },
});
