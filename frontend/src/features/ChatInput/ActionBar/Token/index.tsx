import { lazy, PropsWithChildren, memo, Suspense } from 'react';

import { useModelHasContextWindowToken } from '@/hooks/useModelHasContextWindowToken';
import { useChatStore } from '@/store/chat';
import { chatSelectors, threadSelectors } from '@/store/chat/selectors';

const LargeTokenContent = lazy(() => import('./TokenTag'));
const LargeTokenContentForGroupChat = lazy(() => import('./TokenTagForGroupChat'));

const Token = memo<PropsWithChildren>(({ children }) => {
  const showTag = useModelHasContextWindowToken();

  return showTag && children;
});

export const MainToken = memo(() => {
  const total = useChatStore(chatSelectors.mainAIChatsMessageString);

  return (
    <Token>
      <Suspense fallback={null}>
        <LargeTokenContent total={total} />
      </Suspense>
    </Token>
  );
});

export const PortalToken = memo(() => {
  const total = useChatStore(threadSelectors.portalDisplayChatsString);

  return (
    <Token>
      <Suspense fallback={null}>
        <LargeTokenContent total={total} />
      </Suspense>
    </Token>
  );
});

export const GroupChatToken = memo(() => {
  const total = useChatStore(chatSelectors.mainAIChatsMessageString);

  return (
    <Token>
      <Suspense fallback={null}>
        <LargeTokenContentForGroupChat total={total} />
      </Suspense>
    </Token>
  );
});
