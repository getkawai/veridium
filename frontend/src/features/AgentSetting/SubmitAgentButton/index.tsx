import { ActionIcon, Button } from '@lobehub/ui';
import { Share2 } from 'lucide-react';
import { memo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import SubmitAgentModal from './SubmitAgentModal';

const SubmitAgentButton = memo<{ modal?: boolean }>(({ modal }) => {
  const { t } = useTranslation('setting');
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <>
      {modal ? (
        <Button block icon={Share2} onClick={() => setIsModalOpen(true)} variant={'filled'}>
          {t('submitAgentModal.tooltips')}
        </Button>
      ) : (
        <ActionIcon
          icon={Share2}
          onClick={() => setIsModalOpen(true)}
          size={16}
          title={t('submitAgentModal.tooltips')}
        />
      )}
      <SubmitAgentModal onCancel={() => setIsModalOpen(false)} open={isModalOpen} />
    </>
  );
});

export default SubmitAgentButton;
