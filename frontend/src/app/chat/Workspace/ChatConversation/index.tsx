// import ChatHydration from './features/ChatHydration';
import ChatInput from './features/ChatInput';
import ChatList from './features/ChatList';
import ChatMinimap from './features/ChatMinimap';
// import ThreadHydration from './features/ThreadHydration';
// TopicHydration removed - replaced by useInitTopicState in StoreInitialization
// import TopicHydration from './features/TopicHydration';
import ZenModeToast from './features/ZenModeToast';

const ChatConversation = () => {
  const isMobile = false;

  return (
    <>
      {/* TopicHydration removed - state initialized in StoreInitialization.tsx */}
      <ZenModeToast />
      <ChatList mobile={isMobile} />
      <ChatInput mobile={isMobile} />
      {/* <ChatHydration /> */}
      {/* <ThreadHydration /> */}
      {!isMobile && <ChatMinimap />}
    </>
  );
};

export default ChatConversation;
