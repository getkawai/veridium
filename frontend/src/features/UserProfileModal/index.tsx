'use client';

import { Modal } from '@lobehub/ui';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';

import { useGlobalStore } from '@/store/global';

const UserProfileModal = memo(() => {
  const { t } = useTranslation('common');
  const [open, toggleUserProfile] = useGlobalStore((s) => [
    s.status.isShowUserProfile,
    s.toggleUserProfile,
  ]);

  return (
    <Modal
      footer={null}
      onCancel={() => toggleUserProfile(false)}
      open={open}
      title={t('userPanel.profile')}
    >
      <div style={{ padding: 24 }}>User Profile Content Placeholder</div>
    </Modal>
  );
});

export default UserProfileModal;
