import { Suspense } from 'react';
import SkeletonList from './features/SkeletonList';
import Topic from './features/Topic';
import ConfigSwitcher from './features/ConfigSwitcher';
import { Flexbox } from 'react-layout-kit';

const TopicLayout = () => {
  return (
    <Flexbox height={'100%'} style={{ overflow: 'hidden', position: 'relative' }} width={'100%'}>
      <Suspense fallback={<SkeletonList />}>
        <ConfigSwitcher />
      </Suspense>
      <Topic />
    </Flexbox>
  );
};

export default TopicLayout;
