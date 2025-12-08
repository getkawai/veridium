import { memo } from 'react';
import { useTranslation } from 'react-i18next';

import { useGenerationConfigParam } from '@/store/image/slices/generationConfig/hooks';

import InputNumber from '../../../components/SeedNumberInput';

/**
 * Input untuk seed number (random seed)
 * Seed yang sama akan menghasilkan gambar yang sama dengan parameter lain tetap
 * Berguna untuk reproducibility
 */
const SeedNumberInput = memo(() => {
  const { t } = useTranslation('image');
  const { value, setValue } = useGenerationConfigParam('seed');

  return <InputNumber onChange={setValue} placeholder={t('config.seed.random')} value={value} />;
});

export default SeedNumberInput;
