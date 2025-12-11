import { ChatTranslate } from '@/types';
import { ActionIcon, Icon, Markdown, Tag, copyToClipboard } from '@lobehub/ui';
import { Browser } from '@wailsio/runtime';
import { App } from 'antd';
import { useTheme } from 'antd-style';
import { ChevronDown, ChevronUp, ChevronsRight, CopyIcon, TrashIcon } from 'lucide-react';
import { memo, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import BubblesLoading from '@/components/BubblesLoading';
import { useChatStore } from '@/store/chat';

interface TranslateProps extends ChatTranslate {
  id: string;
  loading?: boolean;
}

const Translate = memo<TranslateProps>(({ content = '', from, to, id, loading }) => {
  const theme = useTheme();
  const { t } = useTranslation('common');
  const [show, setShow] = useState(true);
  const clearTranslate = useChatStore((s) => s.clearTranslate);

  const { message } = App.useApp();

  // Custom components untuk desktop app link handling
  const markdownComponents = useMemo(
    () => ({
      a: ({ href, children, ...props }: any) => (
        <a
          {...props}
          href={href}
          onClick={(e) => {
            e.preventDefault();
            if (href) Browser.OpenURL(href);
          }}
          rel="noopener noreferrer"
          target="_blank"
        >
          {children}
        </a>
      ),
    }),
    []
  );

  return (
    <Flexbox gap={8}>
      <Flexbox align={'center'} horizontal justify={'space-between'}>
        <div>
          <Flexbox gap={4} horizontal>
            <Tag style={{ margin: 0 }}>
              {from && from !== '' ? t(`lang.${from}` as any) : '...'}
            </Tag>
            <Icon color={theme.colorTextTertiary} icon={ChevronsRight} />
            <Tag>
              {t(`lang.${to}` as any)}
            </Tag>
          </Flexbox>
        </div>
        <Flexbox horizontal>
          <ActionIcon
            icon={CopyIcon}
            onClick={async () => {
              await copyToClipboard(content);
              message.success(t('copySuccess'));
            }}
            size={'small'}
            title={t('copy')}
          />
          <ActionIcon
            icon={TrashIcon}
            onClick={() => {
              clearTranslate(id);
            }}
            size={'small'}
            title={t('translate.clear', { ns: 'chat' })}
          />
          <ActionIcon
            icon={show ? ChevronDown : ChevronUp}
            onClick={() => {
              setShow(!show);
            }}
            size={'small'}
          />
        </Flexbox>
      </Flexbox>
      {!show ? null : loading && !content ? (
        <BubblesLoading />
      ) : (
        <Markdown components={markdownComponents} variant={'chat'}>
          {content}
        </Markdown>
      )}
    </Flexbox>
  );
});

export default Translate;
