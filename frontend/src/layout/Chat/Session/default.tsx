import { Suspense } from 'react';

import CircleLoading from '@/components/Loading/CircleLoading';
import SessionListContent from './features/SessionListContent';
import SkeletonList from './features/SkeletonList';
import DesktopLayout from './layout';

const Session = (props: any) => {
  return (
    <Suspense fallback={<CircleLoading />}>
      <DesktopLayout {...props}>
        <Suspense fallback={<SkeletonList />}>
          <SessionListContent />
        </Suspense>
      </DesktopLayout>
    </Suspense>
  );
};

Session.displayName = 'Session';

export default Session;
