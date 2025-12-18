'use client';

import { Modal } from '@lobehub/ui';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import { useGlobalStore } from '@/store/global';

const SettingsModal = memo(() => {
  const { t } = useTranslation('setting');
  const [open, tab, toggleSettings] = useGlobalStore((s) => [
    s.status.isShowSettings,
    s.status.settingsTab,
    s.toggleSettings,
  ]);

  return (
    <Modal
      allowFullscreen
      footer={null}
      onCancel={() => toggleSettings(false)}
      open={open}
      title={t('header.title')}
      width={'min(90%, 1280px)'}
    >
      <Flexbox height={'70vh'} horizontal width={'100%'}>
        <Flexbox style={{ borderRight: '1px solid var(--color-border-secondary)' }} width={200}>
          {/* Category Placeholder */}
          <div style={{ padding: 16 }}>{tab}</div>
        </Flexbox>
        <Flexbox flex={1} padding={24}>
          {/* Content Placeholder */}
          <div>Settings Content for {tab}</div>
        </Flexbox>
      </Flexbox>
    </Modal>
  );
});

export default SettingsModal;
