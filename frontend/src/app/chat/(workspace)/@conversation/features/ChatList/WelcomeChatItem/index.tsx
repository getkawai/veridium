import { memo } from 'react';

// import { sessionSelectors } from '@/store/session/selectors';
// import { useSessionStore } from '@/store/session/store';

// Dummy implementations for UI development
const useSessionStore = (selector?: any) => {
  if (selector) {
    return selector({
      isCurrentSessionGroupSession: false,
    });
  }

  return {
    isCurrentSessionGroupSession: false,
  };
};

const sessionSelectors = {
  isCurrentSessionGroupSession: (state: any) => state.isCurrentSessionGroupSession,
};

import AgentWelcome from './AgentWelcome';
import GroupWelcome from './GroupWelcome';

const WelcomeChatItem = memo(() => {
  const isGroupSession = useSessionStore(sessionSelectors.isCurrentSessionGroupSession);

  if (isGroupSession) return <GroupWelcome />;

  return <AgentWelcome />;
});

export default WelcomeChatItem;
