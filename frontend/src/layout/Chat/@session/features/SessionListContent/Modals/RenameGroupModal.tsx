import { Input, Modal, type ModalProps } from '@lobehub/ui';
import { App } from 'antd';
import { memo, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface RenameGroupModalProps extends ModalProps {
  id: string;
}

const RenameGroupModal = memo<RenameGroupModalProps>(({ id, open, onCancel }) => {
  const { t } = useTranslation('chat');

  // Dummy update function for UI focus
  const updateSessionGroupName = (groupId: string, newName: string) => {
    // Dummy implementation - no actual update
    console.log(`Dummy update: Group ${groupId} renamed to ${newName}`);
  };

  // Dummy group data
  const group = {
    id,
    name: `Group ${id}`,
  };

  const [input, setInput] = useState<string>('');
  const [loading, setLoading] = useState(false);

  const { message } = App.useApp();

  useEffect(() => {
    setInput(group?.name ?? '');
  }, [group]);

  return (
    <Modal
      allowFullscreen
      destroyOnHidden
      okButtonProps={{ loading }}
      onCancel={(e) => {
        setInput(group?.name ?? '');
        onCancel?.(e);
      }}
      onOk={async (e) => {
        if (input.length === 0 || input.length > 20)
          return message.warning(t('sessionGroup.tooLong'));
        setLoading(true);
        await updateSessionGroupName(id, input);
        message.success(t('sessionGroup.renameSuccess'));
        setLoading(false);

        onCancel?.(e);
      }}
      open={open}
      title={t('sessionGroup.rename')}
      width={400}
    >
      <Input
        autoFocus
        defaultValue={group?.name}
        onChange={(e) => setInput(e.target.value)}
        placeholder={t('sessionGroup.inputPlaceholder')}
        value={input}
      />
    </Modal>
  );
});

export default RenameGroupModal;
