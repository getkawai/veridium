'use client';

import { useEffect } from 'react';

import { useSearchParams } from '@/hooks/useNavigation';
import { useChatStore } from '@/store/chat';
import { useRouterStore } from '@/store/router';

import { useSend } from '../useSend';

const MessageFromUrl = () => {
  const updateInputMessage = useChatStore((s) => s.updateInputMessage);
  const { send: sendMessage } = useSend();
  const searchParams = useSearchParams();
  const removeSearchParam = useRouterStore((s) => s.removeSearchParam);

  useEffect(() => {
    const message = searchParams.get('message');
    if (message) {
      // Remove message from router state
      removeSearchParam('message');

      updateInputMessage(message);
      sendMessage();
    }
  }, [searchParams, updateInputMessage, sendMessage, removeSearchParam]);

  return null;
};

export default MessageFromUrl;
