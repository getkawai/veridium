import { Empty } from 'antd';
import { createStyles } from 'antd-style';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Center } from 'react-layout-kit';
import { SESSION_CHAT_URL } from '@/const/url';
// import { useSwitchSession } from '@/hooks/useSwitchSession';
// import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
// import { useSessionStore } from '@/store/session';
// import { sessionSelectors } from '@/store/session/selectors';

// Dummy implementations for development
const useSwitchSession = () => {
  return (sessionId: string) => {
    console.log('Mock switchSession called with:', sessionId);
  };
};

const featureFlagsSelectors = {
  showCreateSession: true,
};

const useServerConfigStore = (selector: any) => {
  // Handle the case where featureFlagsSelectors object is passed directly
  if (selector && typeof selector === 'object' && selector.showCreateSession !== undefined) {
    return selector;
  }
  // Handle selector function
  if (typeof selector === 'function') {
    return selector(featureFlagsSelectors);
  }
  return featureFlagsSelectors;
};

const mockSessionStore = {
  isSessionListInit: true,
};

const useSessionStore = (selector?: any) => {
  if (selector) {
    return selector(mockSessionStore);
  }
  return mockSessionStore;
};

const sessionSelectors = {
  isSessionListInit: (state: any) => state.isSessionListInit,
};

import { LobeSessions } from '@/types/session';

import SkeletonList from '../../SkeletonList';
import AddButton from './AddButton';
import SessionItem from './Item';

const useStyles = createStyles(
  ({ css }) => css`
    min-height: 70px;
  `,
);
interface SessionListProps {
  dataSource?: LobeSessions;
  groupId?: string;
  showAddButton?: boolean;
}
const SessionList = memo<SessionListProps>(({ dataSource, groupId, showAddButton = true }) => {
  const { t } = useTranslation('chat');
  const { styles } = useStyles();

  const isInit = useSessionStore(sessionSelectors.isSessionListInit);
  const { showCreateSession } = useServerConfigStore(featureFlagsSelectors);

  const switchSession = useSwitchSession();

  const isEmpty = !dataSource || dataSource.length === 0;
  return !isInit ? (
    <SkeletonList />
  ) : !isEmpty ? (
    dataSource.map(({ id }) => (
      <div className={styles} key={id}>
        <a
          aria-label={id}
          href={SESSION_CHAT_URL(id, false)}
          onClick={(e) => {
            e.preventDefault();
            switchSession(id);
          }}
        >
          <SessionItem id={id} />
        </a>
      </div>
    ))
  ) : showCreateSession ? (
    showAddButton && <AddButton groupId={groupId} />
  ) : (
    <Center>
      <Empty description={t('emptyAgent')} image={Empty.PRESENTED_IMAGE_SIMPLE} />
    </Center>
  );
});

export default SessionList;
