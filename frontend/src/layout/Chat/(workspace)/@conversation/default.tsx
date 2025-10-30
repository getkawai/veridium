import ChatHydration from './features/ChatHydration';
import ChatInput from './features/ChatInput';
import ChatList from './features/ChatList';
import ChatMinimap from './features/ChatMinimap';
import ThreadHydration from './features/ThreadHydration';
import ZenModeToast from './features/ZenModeToast';

const ChatConversation = async () => {
  const isMobile = false;

  return (
    <>
      <ZenModeToast />
      <ChatList mobile={isMobile} />
      <ChatInput mobile={isMobile} />
      <ChatHydration />
      <ThreadHydration />
      {!isMobile && <ChatMinimap />}
    </>
  );
};

ChatConversation.displayName = 'ChatConversation';

export default ChatConversation;
