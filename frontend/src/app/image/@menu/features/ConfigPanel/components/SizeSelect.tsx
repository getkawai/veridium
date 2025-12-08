import { memo } from 'react';

import { useGenerationConfigParam } from '@/store/image/slices/generationConfig/hooks';

import Select from '../../../components/SizeSelect';

/**
 * Dropdown selector untuk ukuran gambar preset
 * Contoh: "1024x1024", "512x768", dll
 * Options dinamis berdasarkan model yang dipilih
 */
const SizeSelect = memo(() => {
  const { value, setValue, enumValues } = useGenerationConfigParam('size');
  const options = enumValues!.map((size) => ({
    label: size,
    value: size,
  }));

  return <Select onChange={setValue} options={options} value={value} />;
});

export default SizeSelect;
