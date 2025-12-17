import { Modal } from '@lobehub/ui';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';
import List from './List';

interface AssignKnowledgeBaseProps {
  open?: boolean;
  setOpen: (open: boolean) => void;
}
export const AssignKnowledgeBaseModal = memo<AssignKnowledgeBaseProps>(({ setOpen, open }) => {
  const { t } = useTranslation('chat');

  return (
    <Modal
      allowFullscreen
      footer={null}
      onCancel={() => {
        setOpen(false);
      }}
      open={open}
      styles={{ body: { overflow: 'hidden' } }}
      title={t('knowledgeBase.library.title')}
      width={600}
    >
      <Flexbox
        gap={16}
        style={{ maxHeight: 'inherit' }}
        width={'100%'}
      >
        <List />
      </Flexbox>
    </Modal>
  );
});
