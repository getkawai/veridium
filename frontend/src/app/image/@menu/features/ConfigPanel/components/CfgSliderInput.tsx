import { SliderWithInput } from '@lobehub/ui';
import { memo } from 'react';

import { useGenerationConfigParam } from '@/store/image/slices/generationConfig/hooks';

/**
 * Slider input untuk CFG (Classifier-Free Guidance) scale
 * Mengontrol seberapa kuat model mengikuti prompt
 * Nilai lebih tinggi = lebih strict mengikuti prompt
 */
const CfgSliderInput = memo(() => {
  const { value, setValue, min, max } = useGenerationConfigParam('cfg');
  return <SliderWithInput max={max} min={min} onChange={setValue} value={value} />;
});

export default CfgSliderInput;
