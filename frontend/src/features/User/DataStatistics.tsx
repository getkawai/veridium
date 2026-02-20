'use client';

import { Icon, Tooltip } from '@lobehub/ui';
import { Badge } from 'antd';
import { createStyles } from 'antd-style';
import { LoaderCircle } from 'lucide-react';
import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox, FlexboxProps } from 'react-layout-kit';

// import { useClientDataSWR } from '@/libs/swr';
// import { messageService } from '@/services/message';
// import { sessionService } from '@/services/session';
// import { topicService } from '@/services/topic';
// import { useServerConfigStore } from '@/store/serverConfig';
import { formatShortenNumber } from '@/utils/format';
// import { today } from '@/utils/time';

const useStyles = createStyles(({ css, token }) => ({
  card: css`
    padding-block: 6px;
    padding-inline: 8px;
    border-radius: ${token.borderRadius}px;
    background: ${token.colorFillTertiary};

    &:hover {
      background: ${token.colorFillSecondary};
    }
  `,
  count: css`
    font-size: 16px;
    font-weight: bold;
    line-height: 1.2;
  `,
  title: css`
    font-size: 12px;
    line-height: 1.2;
    color: ${token.colorTextDescription};
  `,
  today: css`
    font-size: 12px;
  `,
}));

const DataStatistics = memo<Omit<FlexboxProps, 'children'>>(({ style, ...rest }) => {
  const mobile = false;
  // sessions
  const { data: sessions, isLoading: sessionsLoading } = { data: 0, isLoading: false };
  // topics
  const { data: topics, isLoading: topicsLoading } = { data: 0, isLoading: false };
  // messages
  const { data: { messages, messagesToday } = {}, isLoading: messagesLoading } = { data: { messages: 0, messagesToday: 0 }, isLoading: false };

  const { styles, theme } = useStyles();
  const { t } = useTranslation('common');

  const loading = useMemo(() => <Icon icon={LoaderCircle} spin />, []);

  const items = [
    {
      count: sessionsLoading || sessions === undefined ? loading : sessions,
      key: 'sessions',
      title: t('dataStatistics.sessions'),
    },
    {
      count: topicsLoading || topics === undefined ? loading : topics,
      key: 'topics',
      title: t('dataStatistics.topics'),
    },
    {
      count: messagesLoading || messages === undefined ? loading : messages,
      countToady: messagesToday,
      key: 'messages',
      title: t('dataStatistics.messages'),
    },
  ];

  return (
    <Flexbox
      align={'center'}
      gap={4}
      horizontal
      paddingInline={8}
      style={{ marginBottom: 8, ...style }}
      width={'100%'}
      {...rest}
    >
      {items.map((item) => {
        if (item.key === 'messages') {
          const showBadge = Boolean(item.countToady && item.countToady > 0);
          return (
            <Flexbox
              align={'center'}
              className={styles.card}
              flex={showBadge && !mobile ? 2 : 1}
              gap={4}
              horizontal
              justify={'space-between'}
              key={item.key}
            >
              <Flexbox gap={2}>
                <div className={styles.count}>{formatShortenNumber(item.count)}</div>
                <div className={styles.title}>{item.title}</div>
              </Flexbox>
              {showBadge && (
                <Tooltip title={t('dataStatistics.today')}>
                  <Badge
                    count={`+${item.countToady}`}
                    style={{
                      background: theme.colorSuccess,
                      color: theme.colorSuccessBg,
                      cursor: 'pointer',
                    }}
                  />
                </Tooltip>
              )}
            </Flexbox>
          );
        }

        return (
          <Flexbox className={styles.card} flex={1} gap={2} key={item.key}>
            <Flexbox horizontal>
              <div className={styles.count}>{formatShortenNumber(item.count)}</div>
            </Flexbox>
            <div className={styles.title}>{item.title}</div>
          </Flexbox>
        );
      })}
    </Flexbox>
  );
});

export default DataStatistics;
