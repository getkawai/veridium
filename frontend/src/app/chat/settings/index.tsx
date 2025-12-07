import SafeSpacing from '@/components/SafeSpacing';
import { HEADER_HEIGHT } from '@/const/layoutTokens';
import Footer from '@/features/Setting/Footer';
import SettingContainer from '@/features/Setting/SettingContainer';
import Header from './_layout/Desktop/Header';
import EditPage from './page';

const Layout = () => (
  <>
    <Header />
    <SettingContainer addonAfter={<Footer />} addonBefore={<SafeSpacing height={HEADER_HEIGHT} />}>
      <EditPage />
    </SettingContainer>
  </>
);

Layout.displayName = 'DesktopSessionSettingsLayout';

export default Layout;
