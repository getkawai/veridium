import { ActionIcon, ActionIconProps } from '@lobehub/ui';
import { Browser } from '@wailsio/runtime';
import { FlaskConical, Github } from 'lucide-react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';
import { GITHUB } from '@/const/url';

const ICON_SIZE: ActionIconProps['size'] = {
  blockSize: 36,
  size: 20,
  strokeWidth: 1.5,
};

const BottomActions = memo(() => {
  const { t } = useTranslation('common');

  // const { hideGitHub } = useServerConfigStore(featureFlagsSelectors);
  const hideGitHub = false;

  return (
    <Flexbox gap={8}>
      {!hideGitHub && (
        <a
          aria-label={'GitHub'}
          href={GITHUB}
          onClick={(e) => {
            e.preventDefault();
            Browser.OpenURL(GITHUB);
          }}
          target={'_blank'}
        >
          <ActionIcon
            icon={Github}
            size={ICON_SIZE}
            title={'GitHub'}
            tooltipProps={{ placement: 'right' }}
          />
        </a>
      )}
      <a aria-label={t('labs')} href={'/labs'}>
        <ActionIcon
          icon={FlaskConical}
          size={ICON_SIZE}
          title={t('labs')}
          tooltipProps={{ placement: 'right' }}
        />
      </a>
    </Flexbox>
  );
});

export default BottomActions;
