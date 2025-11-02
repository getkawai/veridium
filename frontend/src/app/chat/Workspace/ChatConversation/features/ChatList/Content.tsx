'use client';

import { memo, useCallback } from 'react';

import { SkeletonList, VirtualizedList } from '@/features/Conversation';
import WideScreenContainer from '@/features/Conversation/components/WideScreenContainer';
// import { useFetchMessages } from '@/hooks/useFetchMessages';
// import { useChatStore } from '@/store/chat';
// import { chatSelectors } from '@/store/chat/selectors';

// Dummy implementations for UI development
const useFetchMessages = () => {
  // Mock hook - does nothing in UI development
};

const useChatStore = (selector?: any) => {
  if (selector) {
    return selector({
      isCurrentChatLoaded: true,
      mainDisplayChatIDs: [],
    });
  }

  return {
    isCurrentChatLoaded: true,
    mainDisplayChatIDs: [],
  };
};

const chatSelectors = {
  isCurrentChatLoaded: (state: any) => state.isCurrentChatLoaded,
  mainDisplayChatIDs: (state: any) => state.mainDisplayChatIDs,
};

import MainChatItem from './ChatItem';
import Welcome from './WelcomeChatItem';

interface ListProps {
  mobile?: boolean;
}

const Content = memo<ListProps>(({ mobile }) => {
  const [isCurrentChatLoaded] = useChatStore((s) => [chatSelectors.isCurrentChatLoaded(s)]);

  useFetchMessages();
  const data = useChatStore(chatSelectors.mainDisplayChatIDs);

  const itemContent = useCallback(
    (index: number, id: string) => <MainChatItem id={id} index={index} />,
    [mobile],
  );

  if (!isCurrentChatLoaded) return <SkeletonList mobile={mobile} />;

  if (data.length === 0)
    return (
      <WideScreenContainer flex={1} height={'100%'}>
        <Welcome />
      </WideScreenContainer>
    );

  return <VirtualizedList dataSource={data} itemContent={itemContent} mobile={mobile} />;
});

Content.displayName = 'ChatListRender';

export default Content;
