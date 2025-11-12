import { UIChatMessage } from '@/types';
import { ActionIconGroup, type ActionIconGroupEvent, ActionIconGroupItemType } from '@lobehub/ui';
import { App } from 'antd';
import { memo, use, useCallback, useContext, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useSearchParams } from '@/hooks/useNavigation';
import { useNativeTTS } from '@/hooks/useNativeTTS';

import ShareMessageModal from '@/features/Conversation/components/ShareMessageModal';
import { VirtuosoContext } from '@/features/Conversation/components/VirtualizedList/VirtuosoContext';
import { useChatStore } from '@/store/chat';
import { threadSelectors } from '@/store/chat/selectors';
import { useSessionStore } from '@/store/session';
import { sessionSelectors } from '@/store/session/selectors';

import { InPortalThreadContext } from '../../../context/InPortalThreadContext';
import { useChatListActionsBar } from '../../../hooks/useChatListActionsBar';
import { ErrorActionsBar } from './Error';

interface AssistantActionsProps {
  data: UIChatMessage;
  id: string;
  index: number;
}
export const AssistantActionsBar = memo<AssistantActionsProps>(({ id, data, index }) => {
  const { error, tools } = data;
  
  const hasThreadSelector = useMemo(
    () => threadSelectors.hasThreadBySourceMsgId(id),
    [id],
  );
  
  const [isThreadMode, hasThread] = useChatStore((s) => [
    !!s.activeThreadId,
    hasThreadSelector(s),
  ]);
  const isGroupSession = useSessionStore(sessionSelectors.isCurrentSessionGroupSession);
  const [showShareModal, setShareModal] = useState(false);

  const {
    regenerate,
    edit,
    delAndRegenerate,
    copy,
    divider,
    del,
    branching,
    // export: exportPDF,
    share,
    translate,
    tts,
  } = useChatListActionsBar({ hasThread });

  // TTS hook
  const { isGlobalLoading: isTTSLoading, start: startTTS } = useNativeTTS(data.content, {
    onError: (err) => {
      console.error('[AssistantActions] TTS failed:', err);
    },
    onSuccess: () => {
      console.log('[AssistantActions] TTS completed');
    },
  });

  const hasTools = !!tools;

  const inPortalThread = useContext(InPortalThreadContext);
  const inThread = isThreadMode || inPortalThread;

  const items = useMemo(() => {
    // Add TTS button with loading state
    const ttsWithLoading = {
      ...tts,
      loading: isTTSLoading,
    };

    if (hasTools) return [delAndRegenerate, copy, ttsWithLoading];

    return [edit, copy, ttsWithLoading, inThread || isGroupSession ? null : branching].filter(
      Boolean,
    ) as ActionIconGroupItemType[];
  }, [inThread, hasTools, isGroupSession, isTTSLoading]);

  const { t } = useTranslation('common');
  const searchParams = useSearchParams();
  const urlTopic = searchParams.get('topic');
  const [
    activeTopicId,
    deleteMessage,
    regenerateMessage,
    translateMessage,
    delAndRegenerateMessage,
    copyMessage,
    openThreadCreator,
    resendThreadMessage,
    delAndResendThreadMessage,
    toggleMessageEditing,
  ] = useChatStore((s) => [
    s.activeTopicId,
    s.deleteMessage,
    s.regenerateMessage,
    s.translateMessage,
    s.delAndRegenerateMessage,
    s.copyMessage,
    s.openThreadCreator,
    s.resendThreadMessage,
    s.delAndResendThreadMessage,
    s.toggleMessageEditing,
  ]);

  const topic = urlTopic || activeTopicId;
  const { message } = App.useApp();
  const virtuosoRef = use(VirtuosoContext);

  const onActionClick = useCallback(
    async (action: ActionIconGroupEvent) => {
      switch (action.key) {
        case 'edit': {
          toggleMessageEditing(id, true);

          virtuosoRef?.current?.scrollIntoView({ align: 'start', behavior: 'auto', index });
        }
      }
      if (!data) return;

      switch (action.key) {
        case 'copy': {
          await copyMessage(id, data.content);
          message.success(t('copySuccess', { defaultValue: 'Copy Success' }));
          break;
        }
        case 'branching': {
          if (!topic) {
            message.warning(t('branchingRequiresSavedTopic'));
            break;
          }
          openThreadCreator(id);
          break;
        }

        case 'del': {
          deleteMessage(id);
          break;
        }

        case 'regenerate': {
          if (inPortalThread) {
            resendThreadMessage(id);
          } else regenerateMessage(id);

          // if this message is an error message, we need to delete it
          if (data.error) deleteMessage(id);
          break;
        }

        case 'delAndRegenerate': {
          if (inPortalThread) {
            delAndResendThreadMessage(id);
          } else {
            delAndRegenerateMessage(id);
          }
          break;
        }

        // case 'export': {
        //   setModal(true);
        //   break;
        // }

        case 'share': {
          setShareModal(true);
          break;
        }

        case 'tts': {
          console.log('[AssistantActions] TTS button clicked');
          startTTS();
          break;
        }
      }

      if (action.keyPath.at(-1) === 'translate') {
        // click the menu data with translate data, the result is:
        // key: 'en-US'
        // keyPath: ['en-US','translate']
        const lang = action.keyPath[0];
        translateMessage(id, lang);
      }
    },
    [data.content, topic, startTTS],
  );

  if (error) return <ErrorActionsBar onActionClick={onActionClick} />;

  return (
    <>
      <ActionIconGroup
        items={items}
        menu={{
          items: [
            edit,
            copy,
            divider,
            translate,
            divider,
            share,
            // exportPDF,
            divider,
            regenerate,
            delAndRegenerate,
            del,
          ],
        }}
        onActionClick={onActionClick}
      />
      {/*{showModal && (*/}
      {/*  <ExportPreview content={data.content} onClose={() => setModal(false)} open={showModal} />*/}
      {/*)}*/}
      <ShareMessageModal
        message={data!}
        onCancel={() => {
          setShareModal(false);
        }}
        open={showShareModal}
      />
    </>
  );
});
