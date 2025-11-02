import { ActionIcon } from '@lobehub/ui';
import { useTheme } from 'antd-style';
import { Wand2 } from 'lucide-react';
import { lazy, memo, Suspense } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import { useGlobalStore } from '@/store/global';
import { globalGeneralSelectors } from '@/store/global/selectors';

const EmojiPicker = lazy(() => import('@lobehub/ui/es/EmojiPicker'));

export interface AutoGenerateAvatarProps {
  background?: string;
  canAutoGenerate?: boolean;
  loading?: boolean;
  onChange?: (value: string) => void;
  onGenerate?: () => void;
  value?: string;
}

const AutoGenerateAvatar = memo<AutoGenerateAvatarProps>(
  ({ loading, background, value, onChange, onGenerate, canAutoGenerate }) => {
    const { t } = useTranslation('common');
    const theme = useTheme();
    const locale = useGlobalStore(globalGeneralSelectors.currentLanguage);

    return (
      <Flexbox
        align={'center'}
        flex={'none'}
        gap={2}
        horizontal
        padding={2}
        style={{
          background: theme.colorBgContainer,
          border: `1px solid ${theme.colorBorderSecondary}`,
          borderRadius: 32,
          paddingRight: 8,
          width: 'fit-content',
        }}
      >
        <Suspense fallback={null}>
          <EmojiPicker
            background={background || theme.colorFillTertiary}
            loading={loading}
            locale={locale}
            onChange={onChange}
            size={48}
            style={{
              background: theme.colorFillTertiary,
            }}
            value={value}
          />
        </Suspense>
        <ActionIcon
          disabled={!canAutoGenerate}
          icon={Wand2}
          loading={loading}
          onClick={onGenerate}
          size="small"
          title={!canAutoGenerate ? t('autoGenerateTooltipDisabled') : t('autoGenerate')}
        />
      </Flexbox>
    );
  },
);

export default AutoGenerateAvatar;
