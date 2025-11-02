import { Suspense, lazy } from 'react';

import Loading from '@/components/Loading/BrandTextLoading';

import { PortalHeader } from '@/features/Portal/router';

import Body from './features/Body';

const PortalBody = lazy(() => import('@/features/Portal/router'));

const PortalLayout = () => {
  return (
    <Suspense fallback={<Loading />}>
      <PortalHeader />
      <Body>
        <PortalBody />
      </Body>
    </Suspense>
  );
};

export default PortalLayout;
