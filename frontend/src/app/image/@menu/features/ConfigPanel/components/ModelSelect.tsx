import { Select, type SelectProps } from '@lobehub/ui';
import { createStyles, useTheme } from 'antd-style';
import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import { ModelItemRender, ProviderItemRender } from '@/components/ModelSelect';
import { useAiInfraStore } from '@/store/aiInfra';
import { aiProviderSelectors } from '@/store/aiInfra/slices/aiProvider/selectors';
import { useImageStore } from '@/store/image';
import { imageGenerationConfigSelectors } from '@/store/image/slices/generationConfig/selectors';
import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
import { EnabledProviderWithModels } from '@/types/aiProvider';

const useStyles = createStyles(({ css, prefixCls }) => ({
  popup: css`
    &.${prefixCls}-select-dropdown .${prefixCls}-select-item-option-grouped {
      padding-inline-start: 12px;
    }
  `,
}));

interface ModelOption {
  label: any;
  provider: string;
  value: string;
}

/**
 * Dropdown selector untuk memilih model AI dan provider
 * Fitur:
 * - Grouped by provider
 * - Empty state handling dengan link ke settings
 * - Dynamic options berdasarkan enabled providers
 * - Quick access ke provider settings
 */
const ModelSelect = memo(() => {
  const { styles } = useStyles();
  const { t } = useTranslation('components');
  const theme = useTheme();
  const { showLLM } = useServerConfigStore(featureFlagsSelectors);

  const [currentModel, currentProvider] = useImageStore((s) => [
    imageGenerationConfigSelectors.model(s),
    imageGenerationConfigSelectors.provider(s),
  ]);
  const setModelAndProviderOnSelect = useImageStore((s) => s.setModelAndProviderOnSelect);

  const enabledImageModelList = useAiInfraStore(aiProviderSelectors.enabledImageModelList);

  const options = useMemo<SelectProps['options']>(() => {
    const getImageModels = (provider: EnabledProviderWithModels) => {
      const modelOptions = provider.children.map((model) => ({
        label: <ModelItemRender {...model} {...model.abilities} showInfoTag={false} />,
        provider: provider.id,
        value: `${provider.id}/${model.id}`,
      }));

      return modelOptions;
    };

    return enabledImageModelList.map((provider) => ({
      label: (
        <Flexbox horizontal justify="space-between">
          <ProviderItemRender
            logo={provider.logo}
            name={provider.name}
            provider={provider.id}
            source={provider.source}
          />
        </Flexbox>
      ),
      options: getImageModels(provider),
    }));
  }, [enabledImageModelList, showLLM, t, theme.colorTextTertiary]);

  return (
    <Select
      classNames={{
        root: styles.popup,
      }}
      onChange={(value, option) => {
        // Skip onChange for disabled options (empty states)
        if (value === 'no-provider' || value.includes('/empty')) return;
        const model = value.split('/').slice(1).join('/');
        const provider = (option as unknown as ModelOption).provider;
        if (model !== currentModel || provider !== currentProvider) {
          setModelAndProviderOnSelect(model, provider);
        }
      }}
      options={options}
      shadow
      size={'large'}
      style={{
        width: '100%',
      }}
      value={currentProvider && currentModel ? `${currentProvider}/${currentModel}` : undefined}
    />
  );
});

export default ModelSelect;
