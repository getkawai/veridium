import { Suspense } from 'react';
import Desktop from './_layout/Desktop';
import SkeletonList from './features/SkeletonList';
import Topic from './features/Topic';
import ConfigSwitcher from './features/ConfigSwitcher';

const TopicLayout = () => {
  return (
    <Desktop>
      <Suspense fallback={<SkeletonList />}>
        <ConfigSwitcher />
      </Suspense>
      <Topic />
    </Desktop>
  );
};

export default TopicLayout;
