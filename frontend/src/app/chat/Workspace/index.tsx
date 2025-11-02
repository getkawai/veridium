import { Suspense } from 'react';
import { Flexbox } from 'react-layout-kit';
import BrandTextLoading from '@/components/Loading/BrandTextLoading';
import ChatHeader from './ChatHeader';
import PortalPanel from './PortalPanel';
import PortalLayout from './PortalLayout';
import TopicPanel from './TopicPanel';
import TopicLayout from './TopicLayout';
import ChatConversation from './ChatConversation';

const Workspace = () => {
  return (
    <>
      <ChatHeader />
      <Flexbox
        height={'100%'}
        horizontal
        style={{ overflow: 'hidden', position: 'relative' }}
        width={'100%'}
      >
        <Flexbox
          height={'100%'}
          style={{ overflow: 'hidden', position: 'relative' }}
          width={'100%'}
        >
          <ChatConversation />
        </Flexbox>
        <PortalPanel>
          <Suspense fallback={<BrandTextLoading />}>
            <PortalLayout />
          </Suspense>
        </PortalPanel>
        <TopicPanel>
          <TopicLayout />
        </TopicPanel>
      </Flexbox>
    </>
  );
};

export default Workspace;
