import { Modal } from '@lobehub/ui';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';

import { useGlobalStore } from '@/store/global';

const ChangelogModal = memo(() => {
  const { t } = useTranslation('common');
  const [open, toggleChangelog] = useGlobalStore((s) => [
    s.status.isShowChangelog,
    s.toggleChangelog,
  ]);

  return (
    <Modal
      footer={null}
      onCancel={() => toggleChangelog(false)}
      open={open}
      title={t('changelog')}
    >
      <div style={{ padding: 24 }}>Changelog Content Placeholder</div>
    </Modal>
  );
});

export default ChangelogModal;
