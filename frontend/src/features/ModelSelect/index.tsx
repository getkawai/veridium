import { Select, type SelectProps } from '@lobehub/ui';
import { createStyles } from 'antd-style';
import { memo, useMemo } from 'react';

import { ModelItemRender, ProviderItemRender, TAG_CLASSNAME } from '@/components/ModelSelect';
// import { useEnabledChatModels } from '@/hooks/useEnabledChatModels';
import { EnabledProviderWithModels } from '@/types/aiProvider';

const useStyles = createStyles(({ css, prefixCls }) => ({
  popup: css`
    &.${prefixCls}-select-dropdown .${prefixCls}-select-item-option-grouped {
      padding-inline-start: 12px;
    }
  `,
  select: css`
    .${prefixCls}-select-selection-item {
      .${TAG_CLASSNAME} {
        display: none;
      }
    }
  `,
}));

interface ModelOption {
  label: any;
  provider: string;
  value: string;
}

interface ModelSelectProps {
  defaultValue?: { model: string; provider?: string };
  onChange?: (props: { model: string; provider: string }) => void;
  requiredAbilities?: (keyof EnabledProviderWithModels['children'][number]['abilities'])[];
  showAbility?: boolean;
  value?: { model: string; provider?: string };
}

const ModelSelect = memo<ModelSelectProps>(
  ({ value, onChange, showAbility = true, requiredAbilities }) => {
    // const enabledList = useEnabledChatModels();
    const enabledList = [
      {
        id: 'kawai',
        name: 'Kawai',
        logo: 'https://example.com/kawai-logo.png',
        source: 'builtin' as const,
        children: [
          {
            id: 'kawai-auto',
            displayName: 'Kawai Auto',
            abilities: { functionCall: true, vision: false }
          }
        ]
      },
      {
        id: 'openai',
        name: 'OpenAI',
        logo: 'https://example.com/openai-logo.png',
        source: 'builtin' as const,
        children: [
          {
            id: 'gpt-4',
            displayName: 'GPT-4',
            abilities: { functionCall: true, vision: true }
          },
          {
            id: 'gpt-3.5-turbo',
            displayName: 'GPT-3.5 Turbo',
            abilities: { functionCall: true, vision: false }
          }
        ]
      }
    ];

    const { styles } = useStyles();

    const options = useMemo<SelectProps['options']>(() => {
      const getChatModels = (provider: EnabledProviderWithModels) => {
        // Safety check: ensure provider.children exists and is an array
        if (!provider.children || !Array.isArray(provider.children)) {
          return [];
        }
        
        const models =
          requiredAbilities && requiredAbilities.length > 0
            ? provider.children.filter((model) =>
                requiredAbilities.every((ability) => Boolean(model.abilities?.[ability])),
              )
            : provider.children;

        return models.map((model) => ({
          label: <ModelItemRender {...model} {...(model.abilities || {})} showInfoTag={showAbility} />,
          provider: provider.id,
          value: `${provider.id}/${model.id}`,
        }));
      };

      if (enabledList.length === 1) {
        const provider = enabledList[0];

        return getChatModels(provider);
      }

      return enabledList
        .map((provider) => {
          const opts = getChatModels(provider);
          if (opts.length === 0) return undefined;

          return {
            label: (
              <ProviderItemRender
                logo={''}
                name={''}
                provider={''}
                source={undefined}
              />
            ),
            options: opts,
          };
        })
        .filter(Boolean) as SelectProps['options'];
    }, [enabledList, requiredAbilities, showAbility]);

    console.log('options', options);
    console.log('enabledList', enabledList);

    return (
      <Select
        className={styles.select}
        classNames={{
          popup: { root: styles.popup },
        }}
        defaultValue={`${value?.provider}/${value?.model}`}
        onChange={(value, option) => {
          const model = value.split('/').slice(1).join('/');
          onChange?.({ model, provider: (option as unknown as ModelOption).provider });
        }}
        options={options}
        popupMatchSelectWidth={false}
        value={`${value?.provider}/${value?.model}`}
      />
    );
  },
);

export default ModelSelect;
