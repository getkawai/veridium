import { Input, Modal, type ModalProps } from '@lobehub/ui';
import { App } from 'antd';
import { MouseEvent, memo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

// Dummy implementations for UI focus
const toggleExpandSessionGroupDummy = (groupId: string, expanded: boolean) => {
  // Dummy implementation
  console.log('Toggle expand group:', groupId, expanded);
};

const updateSessionGroupDummy = (id: string, groupId: string) => {
  console.log('Update session group:', id, 'to', groupId);
  return Promise.resolve();
};

const addSessionGroupDummy = (name: string) => {
  const groupId = `dummy-group-${Date.now()}`;
  console.log('Add group:', name, 'id:', groupId);
  return Promise.resolve(groupId);
};

const dummyGlobalState = {
  toggleExpandSessionGroup: toggleExpandSessionGroupDummy,
};

const dummySessionState = {
  updateSessionGroupId: updateSessionGroupDummy,
  addSessionGroup: addSessionGroupDummy,
};

function useGlobalStore(selector: (s: typeof dummyGlobalState) => any): any {
  return selector(dummyGlobalState);
}

function useSessionStore(selector: (s: typeof dummySessionState) => any): any {
  return selector(dummySessionState);
}

interface CreateGroupModalProps extends ModalProps {
  id: string;
}

const CreateGroupModal = memo<CreateGroupModalProps>(
  ({ id, open, onCancel }: CreateGroupModalProps) => {
    const { t } = useTranslation('chat');

    const toggleExpandSessionGroup = useGlobalStore((s) => s.toggleExpandSessionGroup);
    const { message } = App.useApp();
    const [updateSessionGroup, addCustomGroup] = useSessionStore((s) => [
      s.updateSessionGroupId,
      s.addSessionGroup,
    ]);
    const [input, setInput] = useState('');
    const [loading, setLoading] = useState(false);

    return (
      <div onClick={(e) => e.stopPropagation()}>
        <Modal
          allowFullscreen
          destroyOnHidden
          okButtonProps={{ loading }}
          onCancel={(e) => {
            setInput('');
            onCancel?.(e);
          }}
          onOk={async (e: MouseEvent<HTMLButtonElement>) => {
            if (input.length === 0 || input.length > 20 || input.trim() === '')
              return message.warning(t('sessionGroup.tooLong'));

            setLoading(true);
            const groupId = await addCustomGroup(input);
            await updateSessionGroup(id, groupId);
            toggleExpandSessionGroup(groupId, true);
            setLoading(false);

            message.success(t('sessionGroup.createSuccess'));
            onCancel?.(e);
          }}
          open={open}
          title={t('sessionGroup.createGroup')}
          width={400}
        >
          <Flexbox paddingBlock={16}>
            <Input
              autoFocus
              onChange={(e) => setInput(e.target.value)}
              placeholder={t('sessionGroup.inputPlaceholder')}
              value={input}
            />
          </Flexbox>
        </Modal>
      </div>
    );
  },
);

export default CreateGroupModal;
