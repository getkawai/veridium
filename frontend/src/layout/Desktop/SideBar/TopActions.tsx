import { ActionIcon, ActionIconProps, Hotkey } from '@lobehub/ui';
import { Compass, FolderClosed, MessageSquare, Palette } from 'lucide-react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';
// import { useGlobalStore } from '@/store/global';
// import { SidebarTabKey } from '@/store/global/initialState';
// import { useSessionStore } from '@/store/session';
// import { useUserStore } from '@/store/user';
// import { settingsSelectors } from '@/store/user/selectors';
// import { HotkeyEnum } from '@/types/hotkey';

const ICON_SIZE: ActionIconProps['size'] = {
  blockSize: 40,
  size: 24,
  strokeWidth: 2,
};

export enum SidebarTabKey {
  Chat = 'chat',
  Discover = 'discover',
  Files = 'files',
  Image = 'image',
  Me = 'me',
  Setting = 'settings',
}

export interface TopActionProps {
  isPinned?: boolean | null;
  tab?: SidebarTabKey;
}

//  TODO Change icons
const TopActions = memo<TopActionProps>(({ tab, isPinned }) => {
  const { t } = useTranslation('common');
  // const switchBackToChat = useGlobalStore((s) => s.switchBackToChat);
  // const hotkey = useUserStore(settingsSelectors.getHotkeyById(HotkeyEnum.NavigateToChat));
  // const isChatActive = tab === SidebarTabKey.Chat && !isPinned;
  // const isFilesActive = tab === SidebarTabKey.Files;
  // const isDiscoverActive = tab === SidebarTabKey.Discover;
  // const isImageActive = tab === SidebarTabKey.Image;

  return (
    <Flexbox gap={8}>
      <a
        aria-label={t('tab.chat')}
        href={'/chat'}
        onClick={(e) => {
          // If Cmd key is pressed, let the default link behavior happen (open in new tab)
          if (e.metaKey || e.ctrlKey) {
            return;
          }

          // Otherwise, prevent default and switch session within the current tab
          e.preventDefault();
          // switchBackToChat(useSessionStore.getState().activeId);
        }}
      >
        <ActionIcon
          active={false}
          icon={MessageSquare}
          size={ICON_SIZE}
          title={
            <Flexbox align={'center'} gap={8} horizontal justify={'space-between'}>
              <span>{t('tab.chat')}</span>
              <Hotkey inverseTheme keys={''} />
            </Flexbox>
          }
          tooltipProps={{ placement: 'right' }}
        />
      </a>
      <a aria-label={t('tab.knowledgeBase')} href={'/knowledge'}>
        <ActionIcon
          active={false}
          icon={FolderClosed}
          size={ICON_SIZE}
          title={t('tab.knowledgeBase')}
          tooltipProps={{ placement: 'right' }}
        />
      </a>
      <a aria-label={t('tab.aiImage')} href={'/image'}>
        <ActionIcon
          active={false}
          icon={Palette}
          size={ICON_SIZE}
          title={t('tab.aiImage')}
          tooltipProps={{ placement: 'right' }}
        />
      </a>
      <a aria-label={t('tab.discover')} href={'/discover'}>
        <ActionIcon
          active={false}
          icon={Compass}
          size={ICON_SIZE}
          title={t('tab.discover')}
          tooltipProps={{ placement: 'right' }}
        />
      </a>
    </Flexbox>
  );
});

export default TopActions;
