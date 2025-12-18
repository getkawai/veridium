// import { useRouter } from 'next/navigation';
import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

import BrandWatermark from '@/components/BrandWatermark';
import Menu from '@/components/Menu';
import { enableAuth, enableNextAuth } from '@/const/auth';
import { isDeprecatedEdition } from '@/const/version';
// import { useUserStore } from '@/store/user';
// import { authSelectors } from '@/store/user/selectors';

import DataStatistics from '../DataStatistics';
import UserInfo from '../UserInfo';
import UserLoginOrSignup from '../UserLoginOrSignup';
import LangButton from './LangButton';
import ThemeButton from './ThemeButton';
import { ProfileTabs } from '@/store/global/initialState';
import { useGlobalStore } from '@/store/global';

import { useMenu } from './useMenu';

const PanelContent = memo<{ closePopover: () => void }>(({ closePopover }) => {
  // const router = useRouter();
  // const isLoginWithAuth = useUserStore(authSelectors.isLoginWithAuth);
  const isLoginWithAuth = false;
  // const [openSignIn, signOut] = useUserStore((s) => [s.openLogin, s.logout]);
  const openSignIn = () => { };
  const signOut = () => { };
  const { mainItems, logoutItems } = useMenu();

  const handleSignIn = () => {
    openSignIn();
    closePopover();
  };

  const handleSignOut = () => {
    signOut();
    closePopover();
    // NextAuth doesn't need to redirect to login page
    if (enableNextAuth) return;
    // router.push('/login');
  };

  return (
    <Flexbox gap={2} style={{ minWidth: 300 }}>
      {!enableAuth || (enableAuth && isLoginWithAuth) ? (
        <>
          <UserInfo avatarProps={{ clickable: false }} />
          {!isDeprecatedEdition && (
            <div
              onClick={() => useGlobalStore.getState().toggleUserProfile(true, ProfileTabs.Stats)}
              style={{ color: 'inherit', cursor: 'pointer' }}
            >
              <DataStatistics />
            </div>
          )}
        </>
      ) : (
        <UserLoginOrSignup onClick={handleSignIn} />
      )}

      <Menu items={mainItems} onClick={closePopover} />
      <Flexbox
        align={'center'}
        horizontal
        justify={'space-between'}
        style={isLoginWithAuth ? { paddingRight: 6 } : { padding: '6px 6px 6px 16px' }}
      >
        {isLoginWithAuth ? (
          <Menu items={logoutItems} onClick={handleSignOut} />
        ) : (
          <BrandWatermark />
        )}
        <Flexbox align={'center'} flex={'none'} gap={2} horizontal>
          <LangButton />
          <ThemeButton />
        </Flexbox>
      </Flexbox>
    </Flexbox>
  );
});

export default PanelContent;
