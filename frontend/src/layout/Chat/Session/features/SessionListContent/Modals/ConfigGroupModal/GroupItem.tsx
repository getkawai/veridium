import { ActionIcon, EditableText, SortableList } from '@lobehub/ui';
import { App } from 'antd';
import { createStyles } from 'antd-style';
import { PencilLine, Trash } from 'lucide-react';
import { memo, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface DummyGroupItem {
  id: string;
  name: string;
}

// Dummy implementations for development - memoized
const mockSessionStore = {
  updateSessionGroupName: async (id: string, name: string) => {
    console.log('Mock updateSessionGroupName called with:', id, name);
  },
  removeSessionGroup: async (id: string) => {
    console.log('Mock removeSessionGroup called with:', id);
  },
};

const useSessionStore = (selector?: any) => {
  if (selector) {
    return selector(mockSessionStore);
  }
  return mockSessionStore;
};

const useStyles = createStyles(({ css }) => ({
  content: css`
    position: relative;
    overflow: hidden;
    flex: 1;
  `,
  title: css`
    flex: 1;
    height: 28px;
    line-height: 28px;
    text-align: start;
  `,
}));

const GroupItem = memo<DummyGroupItem>(({ id, name }) => {
  const { t } = useTranslation('chat');
  const { styles } = useStyles();
  const { message, modal } = App.useApp();

  const [editing, setEditing] = useState(false);
  const [updateSessionGroupName, removeSessionGroup] = useSessionStore((s) => [
    s.updateSessionGroupName,
    s.removeSessionGroup,
  ]);

  return (
    <>
      <SortableList.DragHandle />
      {!editing ? (
        <>
          <span className={styles.title}>{name}</span>
          <ActionIcon icon={PencilLine} onClick={() => setEditing(true)} size={'small'} />
          <ActionIcon
            icon={Trash}
            onClick={() => {
              modal.confirm({
                centered: true,
                okButtonProps: {
                  danger: true,
                  type: 'primary',
                },
                onOk: async () => {
                  await removeSessionGroup(id);
                },
                title: t('sessionGroup.confirmRemoveGroupAlert'),
              });
            }}
            size={'small'}
          />
        </>
      ) : (
        <EditableText
          editing={editing}
          onChangeEnd={async (input) => {
            if (name !== input) {
              if (!input) return;
              if (input.length === 0 || input.length > 20 || input.trim() === '')
                return message.warning(t('sessionGroup.tooLong'));

              await updateSessionGroupName(id, input);
              message.success(t('sessionGroup.renameSuccess'));
            }
            setEditing(false);
          }}
          onEditingChange={(e) => setEditing(e)}
          showEditIcon={false}
          style={{ height: 28 }}
          value={name}
        />
      )}
    </>
  );
});

export default GroupItem;
