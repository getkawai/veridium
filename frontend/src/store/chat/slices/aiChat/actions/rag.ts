import { chainRewriteQuery } from '@/prompts';
import { StateCreator } from 'zustand/vanilla';

import { chatService } from '@/services/chat';

// 🔄 MIGRATED: Direct DB imports for RAG operations
import { DB, toNullString } from '@/types/database';
import { getUserId } from '@/store/session/helpers';
import { nanoid } from '@/utils';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/selectors';
import { ChatStore } from '@/store/chat';
import { chatSelectors } from '../../message/selectors';
import { toggleBooleanList } from '@/store/chat/utils';
import { useUserStore } from '@/store/user';
import { systemAgentSelectors } from '@/store/user/selectors';
import { ChatSemanticSearchChunk } from '@/types/chunk';
import { ragService } from '@/services/rag';

export interface ChatRAGAction {
  deleteUserMessageRagQuery: (id: string) => Promise<void>;
  /**
   * Retrieve chunks from semantic search
   */
  internal_retrieveChunks: (
    id: string,
    userQuery: string,
    messages: string[],
  ) => Promise<{ chunks: ChatSemanticSearchChunk[]; queryId?: string; rewriteQuery?: string }>;
  /**
   * Rewrite user content to better RAG query
   */
  internal_rewriteQuery: (id: string, content: string, messages: string[]) => Promise<string>;

  /**
   * Check if we should use RAG
   */
  internal_shouldUseRAG: () => boolean;
  internal_toggleMessageRAGLoading: (loading: boolean, id: string) => void;
  rewriteQuery: (id: string) => Promise<void>;
}

const knowledgeIds = () => agentSelectors.currentKnowledgeIds(useAgentStore.getState());
const hasEnabledKnowledge = () => agentSelectors.hasEnabledKnowledge(useAgentStore.getState());

export const chatRag: StateCreator<ChatStore, [['zustand/devtools', never]], [], ChatRAGAction> = (
  set,
  get,
) => ({
  deleteUserMessageRagQuery: async (id) => {
    const message = chatSelectors.getMessageById(id)(get());

    if (!message || !message.ragQueryId) return;

    // optimistic update the message's ragQuery
    get().internal_dispatchMessage({
      id,
      type: 'updateMessage',
      value: { ragQuery: null },
    });

    // 🔄 MIGRATED: Direct DB call instead of ragService.deleteMessageRagQuery()
    const userId = getUserId();
    await DB.DeleteMessageQuery({
      id: message.ragQueryId,
      userId,
    });

    console.log('[RAG] Deleted message RAG query via direct DB', { queryId: message.ragQueryId });
    await get().refreshMessages();
  },

  internal_retrieveChunks: async (id, userQuery, messages) => {
    get().internal_toggleMessageRAGLoading(true, id);

    const message = chatSelectors.getMessageById(id)(get());

    // 1. get the rewrite query
    let rewriteQuery = message?.ragQuery as string | undefined;

    // if there is no ragQuery and there is a chat history
    // we need to rewrite the user message to get better results
    if (!message?.ragQuery && messages.length > 0) {
      rewriteQuery = await get().internal_rewriteQuery(id, userQuery, messages);
    }

    // 2. retrieve chunks from semantic search
    const files = chatSelectors.currentUserFiles(get()).map((f) => f.id);
    try {
      const chunks = await ragService.semanticSearch(
        rewriteQuery || userQuery,
        knowledgeIds().fileIds.concat(files),
      );

      // Create message query to get queryId
      const userId = getUserId();
      const messageQuery = await DB.CreateMessageQuery({
        id: nanoid(),
        messageId: id,
        rewriteQuery: toNullString(rewriteQuery || userQuery),
        userQuery: toNullString(userQuery),
        clientId: toNullString(''),
        userId,
        embeddingsId: toNullString(null),
      });

      const queryId = messageQuery.id;

      get().internal_toggleMessageRAGLoading(false, id);

      return { chunks, queryId, rewriteQuery };
    } catch {
      get().internal_toggleMessageRAGLoading(false, id);

      return { chunks: [] };
    }
  },
  internal_rewriteQuery: async (id, content, messages) => {
    let rewriteQuery = content;

    const queryRewriteConfig = systemAgentSelectors.queryRewrite(useUserStore.getState());
    if (!queryRewriteConfig.enabled) return content;

    const rewriteQueryParams = {
      model: queryRewriteConfig.model,
      provider: queryRewriteConfig.provider,
      ...chainRewriteQuery(
        content,
        messages,
        !!queryRewriteConfig.customPrompt ? queryRewriteConfig.customPrompt : undefined,
      ),
    };

    let ragQuery = '';
    await chatService.fetchPresetTaskResult({
      onFinish: async (text) => {
        rewriteQuery = text;
      },

      onMessageHandle: (chunk) => {
        if (chunk.type !== 'text') return;
        ragQuery += chunk.text;

        get().internal_dispatchMessage({
          id,
          type: 'updateMessage',
          value: { ragQuery },
        });
      },
      params: rewriteQueryParams,
    });

    return rewriteQuery;
  },
  internal_shouldUseRAG: () => {
    //  if there is enabled knowledge, try with ragQuery
    return hasEnabledKnowledge();
  },

  internal_toggleMessageRAGLoading: (loading, id) => {
    set(
      {
        messageRAGLoadingIds: toggleBooleanList(get().messageRAGLoadingIds, id, loading),
      },
      false,
      'internal_toggleMessageLoading',
    );
  },

  rewriteQuery: async (id) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    // delete the current ragQuery
    await get().deleteUserMessageRagQuery(id);

    const chats = chatSelectors.mainAIChatsWithHistoryConfig(get());

    await get().internal_rewriteQuery(
      id,
      message.content,
      chats.map((m) => m.content),
    );
  },
});
