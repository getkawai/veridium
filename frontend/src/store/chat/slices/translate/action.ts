import { chainTranslate } from '@/prompts';
import { ChatTranslate, TraceNameMap, TracePayload } from '@/types';
import { produce } from 'immer';
import { StateCreator } from 'zustand/vanilla';

import { supportLocales } from '@/locales/resources';
import { chatService } from '@/services/chat';
import { messageService } from '@/services/message';
import { chatSelectors } from '../message/selectors';
import { ChatStore } from '@/store/chat/store';
import { useUserStore } from '@/store/user';
import { systemAgentSelectors } from '@/store/user/selectors';
import { merge } from '@/utils/merge';
import { setNamespace } from '@/utils/storeDebug';
import { TTSService } from '@@/github.com/kawai-network/veridium/internal/tts';

const n = setNamespace('enhance');

/**
 * chat translate
 */
export interface ChatTranslateAction {
  clearTranslate: (id: string) => Promise<void>;
  getCurrentTracePayload: (data: Partial<TracePayload>) => TracePayload;
  translateMessage: (id: string, targetLang: string) => Promise<void>;
  updateMessageTranslate: (id: string, data: Partial<ChatTranslate> | false) => Promise<void>;
}

export const chatTranslate: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatTranslateAction
> = (set, get) => ({
  clearTranslate: async (id) => {
    await get().updateMessageTranslate(id, false);
  },
  getCurrentTracePayload: (data) => ({
    sessionId: get().activeId,
    topicId: get().activeTopicId,
    ...data,
  }),

  translateMessage: async (id, targetLang) => {
    const { internal_toggleChatLoading, updateMessageTranslate, internal_dispatchMessage } = get();

    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    // Get current agent for translation
    const translationSetting = systemAgentSelectors.translation(useUserStore.getState());

    // Detect source language using native language detector (lingua-go)
    // This is much faster and more accurate than using LLM
    let from = '';
    try {
      const detectedLang = await TTSService.DetectLanguageCode(message.content);
      // lingua-go returns locale codes like "en-US", "zh-CN"
      // Check if it's in our supported locales
      if (detectedLang && supportLocales.includes(detectedLang)) {
        from = detectedLang;
      } else {
        // Try to match base language (e.g., "en-US" -> "en")
        const baseLang = detectedLang.split('-')[0];
        const matchedLocale = supportLocales.find(locale => locale.startsWith(baseLang));
        if (matchedLocale) {
          from = matchedLocale;
        }
      }
      console.log('[Translate] Detected source language:', from, 'from:', detectedLang);
    } catch (error) {
      console.error('[Translate] Language detection failed:', error);
      // Fallback: will be empty and translation will proceed without source language
    }

    // create translate extra
    await updateMessageTranslate(id, { content: '', from, to: targetLang });

    internal_toggleChatLoading(true, id, n('translateMessage(start)', { id }));

    let content = '';

    // translate to target language
    await chatService.fetchPresetTaskResult({
      onFinish: async (content) => {
        await updateMessageTranslate(id, { content, from, to: targetLang });
        internal_toggleChatLoading(false, id);
      },
      onMessageHandle: (chunk) => {
        switch (chunk.type) {
          case 'text': {
            internal_dispatchMessage({
              id,
              key: 'translate',
              type: 'updateMessageExtra',
              value: produce({ content: '', from, to: targetLang }, (draft) => {
                content += chunk.text;
                draft.content += content;
              }),
            });
            break;
          }
        }
      },
      params: merge(translationSetting, chainTranslate(message.content, targetLang)),
      trace: get().getCurrentTracePayload({ traceName: TraceNameMap.Translator }),
    });
  },

  updateMessageTranslate: async (id, data) => {
    await messageService.updateMessageTranslate(id, data);

    await get().refreshMessages();
  },
});
