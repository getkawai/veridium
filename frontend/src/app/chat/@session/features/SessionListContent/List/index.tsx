import { Button, Empty } from 'antd';
import { createStyles } from 'antd-style';
import { memo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Center } from 'react-layout-kit';
import { SESSION_CHAT_URL } from '@/const/url';
import { useSwitchSession } from '@/hooks/useSwitchSession';
import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
import { useSessionStore } from '@/store/session';
import { sessionSelectors } from '@/store/session/selectors';
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
  hasMore?: boolean;
  onLoadMore?: () => Promise<void>;
}
const SessionList = memo<SessionListProps>(({ dataSource, groupId, showAddButton = true, hasMore, onLoadMore }) => {
  const { t } = useTranslation('chat');
  const { styles } = useStyles();
  const [loadingMore, setLoadingMore] = useState(false);

  const isInit = useSessionStore(sessionSelectors.isSessionListInit);
  const { showCreateSession } = useServerConfigStore(featureFlagsSelectors);

  const switchSession = useSwitchSession();

  // Filter out inbox session (which is always present) before checking if empty
  const nonInboxSessions = dataSource?.filter((session) => session.id !== 'inbox') || [];
  const isEmpty = nonInboxSessions.length === 0;

  const handleLoadMore = async () => {
    if (!onLoadMore) return;
    setLoadingMore(true);
    await onLoadMore();
    setLoadingMore(false);
  };

  return !isInit ? (
    <SkeletonList />
  ) : !isEmpty ? (
    <>
      {nonInboxSessions.map(({ id }) => (
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
      ))}
      {hasMore && onLoadMore && (
        <Center style={{ marginTop: 12, marginBottom: 12 }}>
          <Button loading={loadingMore} onClick={handleLoadMore}>
            {t('loadMore', { defaultValue: 'Load More' })}
          </Button>
        </Center>
      )}
    </>
  ) : showCreateSession ? (
    showAddButton && <AddButton groupId={groupId} />
  ) : (
    <Center>
      <Empty description={t('emptyAgent')} image={Empty.PRESENTED_IMAGE_SIMPLE} />
    </Center>
  );
});

export default SessionList;
