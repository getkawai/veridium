import { Flexbox } from 'react-layout-kit';

import ImagePanel from '@/features/ImageSidePanel';
import ImageTopicPanel from '@/features/ImageTopicPanel';

import Container from './_layout/Desktop/Container';
import RegisterHotkeys from './_layout/Desktop/RegisterHotkeys';
import ImageWorkspace from './features/ImageWorkspace';
import Topics from './@topic/features/Topics';
import ConfigPanel from './@menu/features/ConfigPanel';

const DesktopImageLayout = () => {
  return (
    <>
      <Flexbox
        height={'100%'}
        horizontal
        style={{ maxWidth: '100%', overflow: 'hidden', position: 'relative' }}
        width={'100%'}
      >
        <ImagePanel>
          <ConfigPanel />
        </ImagePanel>
        <Container>
          <ImageWorkspace />
        </Container>
        <ImageTopicPanel>
          <Topics />
        </ImageTopicPanel>
      </Flexbox>
      <RegisterHotkeys />
    </>
  );
};

export default DesktopImageLayout;