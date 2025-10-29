import { ActionIcon, Icon } from '@lobehub/ui';
import { Popover, type PopoverProps } from 'antd';
import { Monitor, Moon, Sun } from 'lucide-react';
import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import Menu, { type MenuProps } from '@/components/Menu';
// import { useGlobalStore } from '@/store/global';
// import { systemStatusSelectors } from '@/store/global/selectors';

const themeIcons = {
  auto: Monitor,
  dark: Moon,
  light: Sun,
};

const ThemeButton = memo<{ placement?: PopoverProps['placement'] }>(({ placement = 'right' }) => {
  // const [themeMode, switchThemeMode] = useGlobalStore((s) => [
  //   systemStatusSelectors.themeMode(s),
  //   s.switchThemeMode,
  // ]);

  const { t } = useTranslation('setting');

  const items: MenuProps['items'] = useMemo(
    () => [
      {
        icon: <Icon icon={themeIcons.auto} />,
        key: 'auto',
        label: t('settingCommon.themeMode.auto'),
        onClick: () => {},
      },
      {
        icon: <Icon icon={themeIcons.light} />,
        key: 'light',
        label: t('settingCommon.themeMode.light'),
        onClick: () => {},
      },
      {
        icon: <Icon icon={themeIcons.dark} />,
        key: 'dark',
        label: t('settingCommon.themeMode.dark'),
        onClick: () => {},
      },
    ],
    [t],
  );

  return (
    <Popover
      arrow={false}
      content={<Menu items={items} selectable selectedKeys={['auto']} />}
      placement={placement}
      styles={{
        body: {
          padding: 0,
        },
      }}
      trigger={['click', 'hover']}
    >
      <ActionIcon icon={themeIcons.auto} size={{ blockSize: 32, size: 16 }} />
    </Popover>
  );
});

export default ThemeButton;
