import { ActionIcon, Button, Dropdown, type MenuProps } from '@lobehub/ui';
import { HardDriveDownload } from 'lucide-react';
import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { configService } from '@/services/config';
import { useSessionStore } from '@/store/session';

import SubmitAgentButton from './SubmitAgentButton';

export const HeaderContent = memo<{ modal?: boolean }>(({ modal }) => {
  const { t } = useTranslation('setting');
  const id = useSessionStore((s) => s.activeId);

  const items = useMemo<MenuProps['items']>(
    () => [
      {
        key: 'agent',
        label: <div>{t('exportType.agent', { ns: 'common' })}</div>,
        onClick: () => {
          if (!id) return;

          configService.exportSingleAgent(id);
        },
      },
      {
        key: 'agentWithMessage',
        label: <div>{t('exportType.agentWithMessage', { ns: 'common' })}</div>,
        onClick: () => {
          if (!id) return;

          configService.exportSingleSession(id);
        },
      },
    ],
    [id, t],
  );

  return (
    <>
      <SubmitAgentButton modal={modal} />
      <Dropdown arrow={false} menu={{ items }} trigger={['click']}>
        {modal ? (
          <Button block icon={HardDriveDownload} variant={'filled'}>
            {t('export', { ns: 'common' })}
          </Button>
        ) : (
          <ActionIcon
            icon={HardDriveDownload}
            size={16}
            title={t('export', { ns: 'common' })}
          />
        )}
      </Dropdown>
    </>
  );
});

export default HeaderContent;
