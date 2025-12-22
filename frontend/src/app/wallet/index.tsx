import NProgress from '@/components/NProgress';
import Container from './Container';
import Header from './Header';
import DesktopWalletLayout from './wallet';

const WalletLayout = () => {
  return (
    <>
      <NProgress />
      <Container>
        <Header />
        <DesktopWalletLayout />
      </Container>
    </>
  );
};

export default WalletLayout;
