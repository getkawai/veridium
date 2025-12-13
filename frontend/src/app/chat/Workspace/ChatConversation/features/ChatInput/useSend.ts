import { useCallback, useMemo } from 'react';

import { useGeminiChineseWarning } from '@/hooks/useGeminiChineseWarning';
import { getAgentStoreState } from '@/store/agent';
import { agentSelectors } from '@/store/agent/selectors';
import { getChatStoreState, useChatStore } from '@/store/chat';
import { aiChatSelectors, chatSelectors } from '@/store/chat/selectors';
import { fileChatSelectors, useFileStore } from '@/store/file';
import { mentionSelectors, useMentionStore } from '@/store/mention';
import { useSessionStore } from '@/store/session';
import { sessionMetaSelectors } from '@/store/session/selectors';

export interface UseSendMessageParams {
  message?: string;
  isWelcomeQuestion?: boolean;
  onlyAddAIMessage?: boolean;
  onlyAddUserMessage?: boolean;
}

export type UseSendGroupMessageParams = UseSendMessageParams & {
  targetMemberId?: string;
};

export const useSend = () => {
  const [
    isContentEmpty,
    sendMessage,
    addAIMessage,
    stopGenerateMessage,
    generating,
    isSendButtonDisabledByMessage,
  ] = useChatStore((s) => [
    !s.inputMessage,
    s.sendMessage,
    s.addAIMessage,
    s.stopGenerateMessage,
    chatSelectors.isAIGenerating(s),
    chatSelectors.isSendButtonDisabledByMessage(s),
    aiChatSelectors.isCurrentSendMessageLoading(s),
  ]);
  const checkGeminiChineseWarning = useGeminiChineseWarning();

  // 使用订阅以保持最新文件列表
  const reactiveFileList = useFileStore(fileChatSelectors.chatUploadFileList);
  const [isUploadingFiles, clearChatUploadFileList] = useFileStore((s) => [
    fileChatSelectors.isUploadingFiles(s),
    s.clearChatUploadFileList,
  ]);

  const isInputEmpty = isContentEmpty && reactiveFileList.length === 0;

  const canNotSend =
    isInputEmpty || isUploadingFiles || isSendButtonDisabledByMessage;

  const handleSend = async (params: UseSendMessageParams = {}) => {
    const store = useChatStore.getState();

    // Use provided message or fall back to store's inputMessage
    const messageToSend = params.message !== undefined ? params.message : store.inputMessage;
    const fileList = fileChatSelectors.chatUploadFileList(useFileStore.getState());

    // For welcome questions or when message is provided, skip canNotSend check
    const isProvidedMessage = params.message !== undefined;
    if (!isProvidedMessage && canNotSend) return;

    const mainInputEditor = store.mainInputEditor;

    // Only require mainInputEditor if not a provided message (welcome question)
    if (!isProvidedMessage && !mainInputEditor) {
      console.warn('not found mainInputEditor instance');
      return;
    }

    if (chatSelectors.isAIGenerating(store)) return;

    // if there is no message and no image, then we should not send the message
    if (!messageToSend && fileList.length === 0) return;

    // Check for Chinese text warning with Gemini model
    const agentStore = getAgentStoreState();
    const currentModel = agentSelectors.currentAgentModel(agentStore);
    const shouldContinue = await checkGeminiChineseWarning({
      model: currentModel,
      prompt: messageToSend,
      scenario: 'chat',
    });

    if (!shouldContinue) return;

    if (params.onlyAddAIMessage) {
      addAIMessage();
    } else {
      sendMessage({ files: fileList, message: messageToSend, ...params });
    }

    clearChatUploadFileList();
    if (mainInputEditor) {
      mainInputEditor.setExpand(false);
      mainInputEditor.clearContent();
      mainInputEditor.focus();
    }
  };

  const stop = () => {
    const store = getChatStoreState();
    const generating = chatSelectors.isAIGenerating(store);

    if (generating) {
      stopGenerateMessage();
      return;
    }

    const isCreatingMessage = aiChatSelectors.isCurrentSendMessageLoading(store);

    if (isCreatingMessage) {
      // cancelSendMessageInServer();
    }
  };

  return useMemo(
    () => ({
      disabled: canNotSend,
      generating: generating,
      send: handleSend,
      stop,
    }),
    [canNotSend, generating, stop, handleSend],
  );
};

export const useSendGroupMessage = () => {
  const [
    isContentEmpty,
    sendGroupMessage,
    updateInputMessage,
    stopGenerateMessage,
    isSendButtonDisabledByMessage,
    isCreatingMessage,
  ] = useChatStore((s) => [
    !s.inputMessage,
    s.sendGroupMessage,
    s.updateInputMessage,
    s.stopGenerateMessage,
    chatSelectors.isSendButtonDisabledByMessage(s),
    chatSelectors.isCreatingMessage(s),
  ]);

  const isSupervisorThinking = useChatStore((s) =>
    chatSelectors.isSupervisorLoading(s.activeId)(s),
  );
  const checkGeminiChineseWarning = useGeminiChineseWarning();

  const fileList = fileChatSelectors.chatUploadFileList(useFileStore.getState());
  const [isUploadingFiles, clearChatUploadFileList] = useFileStore((s) => [
    fileChatSelectors.isUploadingFiles(s),
    s.clearChatUploadFileList,
  ]);

  const isInputEmpty = isContentEmpty && fileList.length === 0;

  const canNotSend =
    isInputEmpty ||
    isUploadingFiles ||
    isSendButtonDisabledByMessage ||
    isCreatingMessage ||
    isSupervisorThinking;

  const handleSend = useCallback(
    async (params: UseSendGroupMessageParams = {}) => {
      if (canNotSend) return;

      const store = useChatStore.getState();
      if (!store.activeId) return;

      const mainInputEditor = store.mainInputEditor;
      if (!mainInputEditor) {
        console.warn('not found mainInputEditor instance');
        return;
      }

      if (
        chatSelectors.isSupervisorLoading(store.activeId)(store) ||
        chatSelectors.isCreatingMessage(store)
      )
        return;

      const inputMessage = store.inputMessage;

      // if there is no message and no files, then we should not send the message
      if (!inputMessage && fileList.length === 0) return;

      // Check for Chinese text warning with Gemini model
      const agentStore = getAgentStoreState();
      const currentModel = agentSelectors.currentAgentModel(agentStore);
      const shouldContinue = await checkGeminiChineseWarning({
        model: currentModel,
        prompt: inputMessage,
        scenario: 'chat',
      });

      if (!shouldContinue) return;

      // Append mentioned users as plain text like "@userName"
      const mentionState = useMentionStore.getState();
      const mentioned = mentionSelectors.mentionedUsers(mentionState);
      const sessionState = useSessionStore.getState();
      const mentionText =
        mentioned.length > 0
          ? ` ${mentioned
            .map((id) => sessionMetaSelectors.getAgentMetaByAgentId(id)(sessionState).title || id)
            .map((name) => `@${name}`)
            .join(' ')}`
          : '';
      const messageWithMentions = `${inputMessage}${mentionText}`.trim();

      sendGroupMessage({
        files: fileList,
        groupId: store.activeId,
        message: messageWithMentions,
        targetMemberId: params.targetMemberId,
        ...params,
      });

      clearChatUploadFileList();
      mainInputEditor.setExpand(false);
      mainInputEditor.clearContent();
      mainInputEditor.focus();
      updateInputMessage('');
      // clear mentioned users after sending
      mentionState.clearMentionedUsers();
    },
    [
      canNotSend,
      fileList,
      clearChatUploadFileList,
      updateInputMessage,
      checkGeminiChineseWarning,
    ],
  );

  const stop = useCallback(() => {
    const store = getChatStoreState();
    const isAgentGenerating = chatSelectors.isAIGenerating(store);
    const isCreating = chatSelectors.isCreatingMessage(store);

    if (isAgentGenerating) {
      stopGenerateMessage();
      return;
    }

    if (isCreating) {
      // For group messages, we don't have a separate cancel method like in single chat
      // The isCreatingMessage state will be reset when the operation completes
      // We can potentially add a cancel group message functionality in the future
      console.warn('Group message creation in progress, cannot cancel');
    }
  }, [stopGenerateMessage]);

  return useMemo(
    () => ({
      disabled: canNotSend,
      generating: isSupervisorThinking || isCreatingMessage,
      send: handleSend,
      stop,
      updateInputMessage,
    }),
    [canNotSend, isSupervisorThinking, isCreatingMessage, handleSend, stop, updateInputMessage],
  );
};
