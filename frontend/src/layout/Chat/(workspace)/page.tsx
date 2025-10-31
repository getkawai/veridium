import { Suspense } from 'react';

import { serverFeatureFlags } from '@/config/featureFlags';
import { isDesktop } from '@/const/version';
import PageTitle from '../features/PageTitle';
import Changelog from './features/ChangelogModal';

const Page = async () => {
  const { hideDocs, showChangelog } = serverFeatureFlags();

  return (
    <>
      <PageTitle />
      {/* <TelemetryNotification mobile={false} /> */}
      {!isDesktop && showChangelog && !hideDocs && (
        <Suspense>
          <Changelog />
        </Suspense>
      )}
    </>
  );
};

Page.displayName = 'Chat';

export default Page;
