import { PropsWithChildren } from 'react';

import NProgress from '@/components/NProgress';

import Container from './Container';
import Header from './Header';

const WalletLayout = ({ children }: PropsWithChildren) => {
  return (
    <>
      <NProgress />
      <Container>
        <Header />
        {children}
      </Container>
    </>
  );
};

export default WalletLayout;
