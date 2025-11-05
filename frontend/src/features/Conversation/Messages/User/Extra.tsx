import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

import ExtraContainer from '@/features/Conversation/components/Extras/ExtraContainer';
import Translate from '@/features/Conversation/components/Extras/Translate';
import { useChatStore } from '@/store/chat';
import { chatSelectors } from '@/store/chat/selectors';

interface UserMessageExtraProps {
  content: string;
  extra: any;
  id: string;
}
export const UserMessageExtra = memo<UserMessageExtraProps>(({ extra, id, content }) => {
  const loading = useChatStore(chatSelectors.isMessageGenerating(id));

  const showTranslate = !!extra?.translate;

  if (!showTranslate) return;

  return (
    <Flexbox gap={8} style={{ marginTop: 8 }}>
      {extra?.translate && (
        <ExtraContainer>
          <Translate id={id} {...extra?.translate} loading={loading} />
        </ExtraContainer>
      )}
    </Flexbox>
  );
});
