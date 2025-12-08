import { SliderWithInput } from '@lobehub/ui';
import { memo } from 'react';

import { useGenerationConfigParam } from '@/store/image/slices/generationConfig/hooks';

/**
 * Slider input untuk jumlah inference steps
 * Lebih banyak steps = kualitas lebih baik tapi lebih lambat
 */
const StepsSliderInput = memo(() => {
  const { value, setValue, min, max } = useGenerationConfigParam('steps');
  return <SliderWithInput max={max} min={min} onChange={setValue} value={value} />;
});

export default StepsSliderInput;
