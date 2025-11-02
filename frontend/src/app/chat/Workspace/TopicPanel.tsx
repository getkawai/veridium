'use client';

import { DraggablePanel, DraggablePanelContainer } from '@lobehub/ui';
import { createStyles, useResponsive } from 'antd-style';
import isEqual from 'fast-deep-equal';
import { memo } from 'react';

import { CHAT_SIDEBAR_WIDTH } from '@/const/layoutTokens';
import TopicLayout from './TopicLayout';
// import { useChatStore } from '@/store/chat';
// import { chatPortalSelectors } from '@/store/chat/slices/portal/selectors';
// import { useGlobalStore } from '@/store/global';
// import { systemStatusSelectors } from '@/store/global/selectors';

// Dummy implementations for development - memoized
const mockChatStore = {
  showPortal: false,
};

const useChatStore = (selector?: any) => {
  if (selector) {
    return selector(mockChatStore);
  }
  return mockChatStore;
};

const chatPortalSelectors = {
  showPortal: (state: any) => state.showPortal,
};

const mockGlobalStore = {
  showChatSideBar: true,
  toggleChatSideBar: (value: boolean) => {
    console.log('Mock toggleChatSideBar called with:', value);
  },
};

const useGlobalStore = (selector: any) => {
  if (selector) {
    return selector(mockGlobalStore);
  }
  return mockGlobalStore;
};

const systemStatusSelectors = {
  showChatSideBar: (state: any) => state.showChatSideBar,
};

const useStyles = createStyles(({ css, token }) => ({
  content: css`
    display: flex;
    flex-direction: column;
    height: 100% !important;
  `,
  drawer: css`
    z-index: 20;
    background: ${token.colorBgContainer};
  `,
  header: css`
    border-block-end: 1px solid ${token.colorBorderSecondary};
  `,
}));

const TopicPanel = memo(() => {
  const { styles } = useStyles();
  const { md = true } = useResponsive();
  const [showTopic, toggleConfig] = useGlobalStore((s) => [
    systemStatusSelectors.showChatSideBar(s),
    s.toggleChatSideBar,
  ]);
  const showPortal = useChatStore(chatPortalSelectors.showPortal);

  const handleExpand = (expand: boolean) => {
    if (isEqual(expand, Boolean(showTopic))) return;
    toggleConfig(expand);
  };

  // Temporarily commented out to prevent potential infinite loop during UI development
  // useEffect(() => {
  //   if (lg && cacheExpand) toggleConfig(true);
  //   if (!lg) toggleConfig(false);
  // }, [lg, cacheExpand]);

  return (
    <DraggablePanel
      className={styles.drawer}
      classNames={{
        content: styles.content,
      }}
      expand={showTopic && !showPortal}
      minWidth={CHAT_SIDEBAR_WIDTH}
      mode={md ? 'fixed' : 'float'}
      onExpandChange={handleExpand}
      placement={'right'}
      showHandleWhenCollapsed={false}
      showHandleWideArea={false}
      styles={{
        handle: { display: 'none' },
      }}
    >
      <DraggablePanelContainer
        style={{
          flex: 'none',
          height: '100%',
          maxHeight: '100vh',
          minWidth: CHAT_SIDEBAR_WIDTH,
        }}
      >
        <TopicLayout />
      </DraggablePanelContainer>
    </DraggablePanel>
  );
});

export default TopicPanel;
